package pb

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// placeholderTitle matches the "Episódio 648" stand-in the static seed uses
// when it has a video id but no real title. The ingest tool overwrites those,
// never a title a human has since typed in.
var placeholderTitle = regexp.MustCompile(`^Epis[oó]dio \d+$`)

// DraftEpisode is what the ingest tool knows about a video before a human has
// looked at it: real metadata pulled from the YouTube Data API, but no show
// number unless the title happened to carry one.
type DraftEpisode struct {
	YouTubeID       string
	Number          int
	Title           string
	PublishedAt     time.Time
	DurationSeconds int
	ThumbnailURL    string
	Description     string
}

// DraftMoment and DraftOpening are what an LLM pass over a transcript
// produces: a guess, not a fact, which is what Confidence is for. Source is
// always "llm" here; a human clears that by editing the record once they've
// checked it against the video.
type DraftMoment struct {
	TypeSlug       string
	VideoTimestamp int
	Summary        string
	Confidence     float64
}

type DraftOpening struct {
	Title        string
	AuthorName   string
	AuthorHandle string
	Genre        string
	AtSeconds    int
	Confidence   float64
}

// UpsertEpisode finds the episode for a YouTube video, or creates it. An
// episode that already exists only has its blank or placeholder fields
// filled in (a real title replaces "Episódio 648", a zero duration gets a
// real one), so a field a human has since corrected is never touched.
func UpsertEpisode(app core.App, ep DraftEpisode) (string, error) {
	existing, err := app.FindRecordsByFilter(ColEpisodes, "youtube_id = {:yid}", "", 1, 0, dbx.Params{"yid": ep.YouTubeID})
	if err != nil {
		return "", mapErr(err)
	}
	if len(existing) > 0 {
		rec := existing[0]
		changed := false

		if title := rec.GetString("title"); (title == "" || placeholderTitle.MatchString(title)) && ep.Title != "" {
			rec.Set("title", ep.Title)
			changed = true
		}
		if rec.GetDateTime("published_at").IsZero() && !ep.PublishedAt.IsZero() {
			setDate(rec, "published_at", ep.PublishedAt)
			changed = true
		}
		if rec.GetInt("duration_seconds") == 0 && ep.DurationSeconds > 0 {
			rec.Set("duration_seconds", ep.DurationSeconds)
			changed = true
		}
		if rec.GetString("thumbnail_url") == "" && ep.ThumbnailURL != "" {
			rec.Set("thumbnail_url", ep.ThumbnailURL)
			changed = true
		}
		if rec.GetString("description") == "" && ep.Description != "" {
			rec.Set("description", ep.Description)
			changed = true
		}
		if changed {
			if err := app.Save(rec); err != nil {
				return "", err
			}
		}
		return rec.Id, nil
	}

	slug := fmt.Sprintf("yt-%s", ep.YouTubeID)
	if ep.Number > 0 {
		slug = fmt.Sprintf("ep-%d", ep.Number)
	}

	return insert(app, ColEpisodes, func(rec *core.Record) {
		rec.Set("youtube_id", ep.YouTubeID)
		rec.Set("number", ep.Number)
		rec.Set("slug", slug)
		rec.Set("title", ep.Title)
		setDate(rec, "published_at", ep.PublishedAt)
		rec.Set("duration_seconds", ep.DurationSeconds)
		rec.Set("thumbnail_url", ep.ThumbnailURL)
		rec.Set("description", ep.Description)
	})
}

// InsertDraftMoment adds one LLM-sourced moment, skipping it if a moment
// already sits at that exact second of that type: re-running ingest on a
// video it has already seen should not double every count.
func InsertDraftMoment(app core.App, episodeID string, m DraftMoment) error {
	typeID, ok, err := findTypeID(app, m.TypeSlug)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("moment type %q does not exist: add it to the vocabulary first", m.TypeSlug)
	}

	dup, err := app.FindRecordsByFilter(
		ColMoments, "episode = {:episode} && type = {:type} && video_timestamp = {:at}", "", 1, 0,
		dbx.Params{"episode": episodeID, "type": typeID, "at": m.VideoTimestamp},
	)
	if err != nil {
		return mapErr(err)
	}
	if len(dup) > 0 {
		return nil
	}

	_, err = insert(app, ColMoments, func(rec *core.Record) {
		rec.Set("episode", episodeID)
		rec.Set("type", typeID)
		rec.Set("video_timestamp", m.VideoTimestamp)
		rec.Set("summary", m.Summary)
		rec.Set("source", "llm")
		rec.Set("confidence", m.Confidence)
		setJSON(rec, "payload", map[string]any{})
	})
	return err
}

// InsertDraftOpening adds one LLM-sourced opening, deduplicated the same way
// as a moment: same episode, same second.
func InsertDraftOpening(app core.App, episodeID string, o DraftOpening) error {
	dup, err := app.FindRecordsByFilter(
		ColOpenings, "episode = {:episode} && at_seconds = {:at}", "", 1, 0,
		dbx.Params{"episode": episodeID, "at": o.AtSeconds},
	)
	if err != nil {
		return mapErr(err)
	}
	if len(dup) > 0 {
		return nil
	}

	_, err = insert(app, ColOpenings, func(rec *core.Record) {
		rec.Set("title", o.Title)
		rec.Set("author_name", o.AuthorName)
		rec.Set("author_handle", o.AuthorHandle)
		rec.Set("genre", o.Genre)
		rec.Set("episode", episodeID)
		rec.Set("at_seconds", o.AtSeconds)
		rec.Set("source", "llm")
		rec.Set("confidence", o.Confidence)
	})
	return err
}

func findTypeID(app core.App, slug string) (string, bool, error) {
	records, err := app.FindRecordsByFilter(ColMomentTypes, "slug = {:slug}", "", 1, 0, dbx.Params{"slug": slug})
	if err != nil {
		return "", false, mapErr(err)
	}
	if len(records) == 0 {
		return "", false, nil
	}
	return records[0].Id, true, nil
}
