package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	kkyoutube "github.com/kkdai/youtube/v2"
	"github.com/pocketbase/pocketbase"
	"google.golang.org/api/googleapi"
	youtubedata "google.golang.org/api/youtube/v3"

	"dls/dls-core/adapters/pb"
)

type ingestRun struct {
	app          *pocketbase.PocketBase
	ytData       *youtubedata.Service
	ytdl         kkyoutube.Client
	llm          openAIClient
	dryRun       bool
	metadataOnly bool
	lang         string
}

// video ingests one video. Every error it returns already says which stage
// failed and, where the cause is known, what to do about it.
func (r ingestRun) video(ctx context.Context, videoID string) error {
	log.Printf("%s: fetching metadata", videoID)
	episode, err := fetchMetadata(ctx, r.ytData, videoID)
	if err != nil {
		return fmt.Errorf("metadata: %w", explainYouTubeAPIError(err))
	}

	var result extraction
	if !r.metadataOnly {
		log.Printf("%s: fetching transcript", videoID)
		video, err := r.ytdl.GetVideo(videoID)
		if err != nil {
			return fmt.Errorf("video info: %w", err)
		}
		transcript, err := fetchTranscript(ctx, video, r.lang)
		if err != nil {
			return fmt.Errorf("transcript: %w", err)
		}
		if len(transcript) == 0 {
			return fmt.Errorf("transcript is empty; captions may still be generating")
		}

		log.Printf("%s: categorizing %d transcript segments", videoID, len(transcript))
		if result, err = categorize(ctx, r.llm, transcript); err != nil {
			return fmt.Errorf("categorize: %w", err)
		}
		if len(result.Moments) == 0 && len(result.Openings) == 0 {
			log.Printf("%s: WARNING: the model found no moments or openings, check the transcript language and the model", videoID)
		}
	}

	if r.dryRun {
		printDraft(videoID, episode, result)
		return nil
	}

	episodeID, err := pb.UpsertEpisode(r.app, episode)
	if err != nil {
		return fmt.Errorf("saving episode: %w", err)
	}

	// A row that fails to save is skipped and counted rather than aborting
	// the rest: entries are deduplicated by second, so a re-run fills the gaps.
	var saveErrs []error
	saved := 0
	for _, m := range result.Moments {
		if err := pb.InsertDraftMoment(r.app, episodeID, pb.DraftMoment{
			TypeSlug:       m.TypeSlug,
			VideoTimestamp: m.AtSeconds,
			Summary:        m.Summary,
			Confidence:     m.Confidence,
		}); err != nil {
			saveErrs = append(saveErrs, fmt.Errorf("moment at %ds: %w", m.AtSeconds, err))
			continue
		}
		saved++
	}
	for _, o := range result.Openings {
		// title is a required field; genre is the closest thing the
		// transcript gives when the show never names the piece.
		if o.Title == "" {
			o.Title = o.Genre
		}
		if err := pb.InsertDraftOpening(r.app, episodeID, pb.DraftOpening{
			Title:        o.Title,
			AuthorName:   o.AuthorName,
			AuthorHandle: o.AuthorHandle,
			Genre:        o.Genre,
			AtSeconds:    o.AtSeconds,
			Confidence:   o.Confidence,
		}); err != nil {
			saveErrs = append(saveErrs, fmt.Errorf("opening at %ds: %w", o.AtSeconds, err))
			continue
		}
		saved++
	}

	log.Printf("%s: saved %d of %d entries (%d moments, %d openings found)",
		videoID, saved, len(result.Moments)+len(result.Openings), len(result.Moments), len(result.Openings))
	if len(saveErrs) > 0 {
		return fmt.Errorf("%d entries not saved: %w", len(saveErrs), errors.Join(saveErrs...))
	}
	return nil
}

func explainYouTubeAPIError(err error) error {
	var apiErr *googleapi.Error
	if !errors.As(err, &apiErr) {
		return err
	}
	for _, item := range apiErr.Errors {
		switch item.Reason {
		case "quotaExceeded", "dailyLimitExceeded", "rateLimitExceeded":
			return fmt.Errorf("YouTube API quota exhausted, try again after it resets: %w", err)
		case "keyInvalid", "API_KEY_INVALID":
			return fmt.Errorf("YOUTUBE_API_KEY was rejected: %w", err)
		case "accessNotConfigured", "forbidden":
			return fmt.Errorf("YouTube Data API v3 is not enabled for this key's project, or the key is restricted: %w", err)
		}
	}
	return err
}
