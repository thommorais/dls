package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	kkyoutube "github.com/kkdai/youtube/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
	googleoption "google.golang.org/api/option"
	youtubedata "google.golang.org/api/youtube/v3"

	"dls/dls-core/adapters/pb"
	"dls/dls-core/seed"
)

// newIngestCommand turns a YouTube video into draft data: real metadata from
// the YouTube Data API, and moments/openings an LLM picked out of the
// transcript. Everything it writes carries source "llm" and a confidence
// score, so it sits in the same tables as the hand-checked data but stays
// visibly a draft until a human confirms it.
func newIngestCommand(app *pocketbase.PocketBase) *cobra.Command {
	var dryRun bool
	var metadataOnly bool
	var lang string

	cmd := &cobra.Command{
		Use:   "ingest [video-id-or-url...]",
		Short: "Fetch a video's metadata and transcript, and draft its moments",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			youtubeKey := os.Getenv("YOUTUBE_API_KEY")
			if youtubeKey == "" {
				return fmt.Errorf("YOUTUBE_API_KEY is not set")
			}

			ctx := context.Background()
			ytData, err := youtubedata.NewService(ctx, googleoption.WithAPIKey(youtubeKey))
			if err != nil {
				return fmt.Errorf("youtube data api client: %w", err)
			}

			var llm openAIClient
			if !metadataOnly {
				var err error
				if llm, err = newOpenAIClientFromEnv(); err != nil {
					return fmt.Errorf("%w (use --metadata-only to skip categorization)", err)
				}
			}

			ytdl := kkyoutube.Client{}

			for _, arg := range args {
				videoID := extractVideoID(arg)
				log.Printf("%s: fetching metadata", videoID)

				episode, err := fetchMetadata(ctx, ytData, videoID)
				if err != nil {
					return fmt.Errorf("%s: metadata: %w", videoID, err)
				}

				var result extraction
				if !metadataOnly {
					log.Printf("%s: fetching transcript", videoID)
					video, err := ytdl.GetVideo(videoID)
					if err != nil {
						return fmt.Errorf("%s: video info: %w", videoID, err)
					}
					transcript, err := ytdl.GetTranscript(video, lang)
					if err != nil {
						return fmt.Errorf("%s: transcript: %w", videoID, err)
					}

					log.Printf("%s: categorizing %d transcript segments", videoID, len(transcript))
					result, err = categorize(ctx, llm, transcript)
					if err != nil {
						return fmt.Errorf("%s: categorize: %w", videoID, err)
					}
				}

				if dryRun {
					printDraft(videoID, episode, result)
					continue
				}

				episodeID, err := pb.UpsertEpisode(app, episode)
				if err != nil {
					return fmt.Errorf("%s: upsert episode: %w", videoID, err)
				}

				for _, m := range result.Moments {
					if err := pb.InsertDraftMoment(app, episodeID, pb.DraftMoment{
						TypeSlug:       m.TypeSlug,
						VideoTimestamp: m.AtSeconds,
						Summary:        m.Summary,
						Confidence:     m.Confidence,
					}); err != nil {
						return fmt.Errorf("%s: moment at %ds: %w", videoID, m.AtSeconds, err)
					}
				}
				for _, o := range result.Openings {
					// title is a required field; genre is the closest thing the
					// transcript gives when the show never names the piece.
					if o.Title == "" {
						o.Title = o.Genre
					}
					if err := pb.InsertDraftOpening(app, episodeID, pb.DraftOpening{
						Title:        o.Title,
						AuthorName:   o.AuthorName,
						AuthorHandle: o.AuthorHandle,
						Genre:        o.Genre,
						AtSeconds:    o.AtSeconds,
						Confidence:   o.Confidence,
					}); err != nil {
						return fmt.Errorf("%s: opening at %ds: %w", videoID, o.AtSeconds, err)
					}
				}
				log.Printf("%s: drafted %d moments, %d openings", videoID, len(result.Moments), len(result.Openings))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "fetch and categorize, but do not write to the database")
	cmd.Flags().BoolVar(&metadataOnly, "metadata-only", false, "only fetch title/duration/publish date, skip the transcript and LLM step")
	cmd.Flags().StringVar(&lang, "lang", "pt", "transcript language code")
	return cmd
}

// extraction is what one categorization pass returns: guesses, each with the
// model's own confidence, not facts.
type extraction struct {
	Moments  []momentGuess  `json:"moments"`
	Openings []openingGuess `json:"openings"`
}

type momentGuess struct {
	TypeSlug   string  `json:"type_slug"`
	AtSeconds  int     `json:"at_seconds"`
	Summary    string  `json:"summary"`
	Confidence float64 `json:"confidence"`
}

type openingGuess struct {
	AtSeconds    int     `json:"at_seconds"`
	Title        string  `json:"title"`
	AuthorName   string  `json:"author_name"`
	AuthorHandle string  `json:"author_handle"`
	Genre        string  `json:"genre"`
	Confidence   float64 `json:"confidence"`
}

var (
	videoIDPattern       = regexp.MustCompile(`(?:v=|youtu\.be/|/embed/|/live/|/shorts/)([A-Za-z0-9_-]{6,})`)
	durationPattern      = regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	episodeNumberPattern = regexp.MustCompile(`(?i)DLSHOW\s*#?\s*(\d+)`)
)

// extractVideoID accepts a bare id, a watch/embed/short url, and returns just
// the id either way.
func extractVideoID(arg string) string {
	if m := videoIDPattern.FindStringSubmatch(arg); len(m) == 2 {
		return m[1]
	}
	return arg
}

func fetchMetadata(ctx context.Context, svc *youtubedata.Service, videoID string) (pb.DraftEpisode, error) {
	resp, err := svc.Videos.List([]string{"snippet", "contentDetails"}).Id(videoID).Context(ctx).Do()
	if err != nil {
		return pb.DraftEpisode{}, err
	}
	if len(resp.Items) == 0 {
		return pb.DraftEpisode{}, fmt.Errorf("no video found for id %q", videoID)
	}

	item := resp.Items[0]
	publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

	thumb := ""
	if t := item.Snippet.Thumbnails; t != nil {
		switch {
		case t.Maxres != nil:
			thumb = t.Maxres.Url
		case t.High != nil:
			thumb = t.High.Url
		case t.Default != nil:
			thumb = t.Default.Url
		}
	}

	return pb.DraftEpisode{
		YouTubeID:       videoID,
		Number:          parseEpisodeNumber(item.Snippet.Title),
		Title:           item.Snippet.Title,
		PublishedAt:     publishedAt,
		DurationSeconds: parseISO8601Duration(item.ContentDetails.Duration),
		ThumbnailURL:    thumb,
		Description:     item.Snippet.Description,
	}, nil
}

func parseISO8601Duration(s string) int {
	m := durationPattern.FindStringSubmatch(s)
	if m == nil {
		return 0
	}
	h, _ := strconv.Atoi(m[1])
	mnt, _ := strconv.Atoi(m[2])
	sec, _ := strconv.Atoi(m[3])
	return h*3600 + mnt*60 + sec
}

func parseEpisodeNumber(title string) int {
	m := episodeNumberPattern.FindStringSubmatch(title)
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

// categorize sends the whole transcript to the model in one pass and forces a
// function call, so the response is always the structured shape below rather
// than prose that has to be parsed back out.
func categorize(ctx context.Context, llm openAIClient, transcript kkyoutube.VideoTranscript) (extraction, error) {
	var out extraction

	typeSlugs := make([]string, 0, len(seed.MomentTypes))
	var typeDocs strings.Builder
	for _, t := range seed.MomentTypes {
		typeSlugs = append(typeSlugs, t.Slug)
		fmt.Fprintf(&typeDocs, "- %s: %s\n", t.Slug, t.Description)
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"moments": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type_slug":  map[string]any{"type": "string", "enum": typeSlugs},
						"at_seconds": map[string]any{"type": "integer"},
						"summary":    map[string]any{"type": "string"},
						"confidence": map[string]any{"type": "number"},
					},
					"required": []string{"type_slug", "at_seconds", "summary", "confidence"},
				},
			},
			"openings": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"at_seconds":    map[string]any{"type": "integer"},
						"title":         map[string]any{"type": "string", "description": "how the show introduces or describes this opening; empty string if it does not"},
						"author_name":   map[string]any{"type": "string"},
						"author_handle": map[string]any{"type": "string"},
						"genre":         map[string]any{"type": "string"},
						"confidence":    map[string]any{"type": "number"},
					},
					"required": []string{"at_seconds", "author_name", "genre", "confidence"},
				},
			},
		},
		"required": []string{"moments", "openings"},
	}

	prompt := fmt.Sprintf(
		"This transcript is from a Brazilian YouTube talk show, in Portuguese. Each line is prefixed with its "+
			"start time in seconds, like [90s].\n\n"+
			"Find every moment that fits one of these recurring categories:\n%s\n"+
			"Also find every listener-submitted opening/vinheta: a piece of music or a jingle sent in by a "+
			"viewer that the show plays, usually announced with the sender's name. It can happen at any point "+
			"in the video, often more than 30 minutes in, so read the whole transcript.\n\n"+
			"Only extract things that actually happen in the transcript, never invent one to fill a category. "+
			"Copy at_seconds from the bracketed time of the line where it happens. confidence is your own 0-1 "+
			"estimate of how sure you are this is a real instance of the category.\n\n%s",
		typeDocs.String(), transcriptText(transcript),
	)

	raw, err := llm.callFunction(ctx, "record_extraction",
		"Record every categorized moment and listener-submitted opening found in this transcript.",
		schema, prompt)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("parsing model output: %w", err)
	}
	return out, nil
}

func transcriptText(t kkyoutube.VideoTranscript) string {
	var b strings.Builder
	for _, seg := range t {
		fmt.Fprintf(&b, "[%ds] %s\n", seg.StartMs/1000, strings.TrimSpace(seg.Text))
	}
	return b.String()
}

func printDraft(videoID string, episode pb.DraftEpisode, result extraction) {
	fmt.Printf("=== %s ===\n", videoID)
	encoded, _ := json.MarshalIndent(struct {
		Episode pb.DraftEpisode `json:"episode"`
		extraction
	}{episode, result}, "", "  ")
	fmt.Println(string(encoded))
}
