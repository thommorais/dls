package pb

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"

	"dls/dls-core/domain"
	"dls/dls-core/seed"
)

// Counts is what the seed command prints, so a run says what it wrote.
type Counts struct {
	Songs       int
	People      int
	Types       int
	Episodes    int
	Appearances int
	Moments     int
	Openings    int
	Archive     int
}

func (c Counts) String() string {
	return fmt.Sprintf(
		"%d songs, %d people, %d moment types, %d episodes, %d appearances, %d moments, %d openings, %d archive entries",
		c.Songs, c.People, c.Types, c.Episodes, c.Appearances, c.Moments, c.Openings, c.Archive,
	)
}

// Seed writes a development dataset. It refuses to run over existing data
// unless reset is set: the command exists to fill a fresh database, and
// running it twice by accident should not double every counter.
func Seed(app core.App, data seed.Dataset, reset bool) (Counts, error) {
	var counts Counts

	existing, err := app.FindAllRecords(ColEpisodes)
	if err != nil {
		return counts, err
	}
	if len(existing) > 0 {
		if !reset {
			return counts, fmt.Errorf("%d episodes already exist: pass --reset to replace them", len(existing))
		}
		if err := wipe(app); err != nil {
			return counts, err
		}
	}

	songIDs := map[string]string{}
	for _, song := range data.Songs {
		id, err := insert(app, ColSongs, func(rec *core.Record) {
			rec.Set("title", song.Title)
			rec.Set("artist", song.Artist)
			rec.Set("genre", song.Genre)
			rec.Set("year", song.Year)
		})
		if err != nil {
			return counts, fmt.Errorf("song %s: %w", song.Key, err)
		}
		songIDs[song.Key] = id
		counts.Songs++
	}

	personIDs := map[string]string{}
	for _, person := range data.People {
		id, err := insert(app, ColPeople, func(rec *core.Record) {
			rec.Set("slug", person.Slug)
			rec.Set("name", person.Name)
			rec.Set("kind", string(person.Kind))
			rec.Set("gender", string(person.Gender))
			setJSON(rec, "socials", map[string]string{})
		})
		if err != nil {
			return counts, fmt.Errorf("person %s: %w", person.Slug, err)
		}
		personIDs[person.Slug] = id
		counts.People++
	}

	typeIDs := map[string]string{}
	for _, momentType := range data.Types {
		id, err := insert(app, ColMomentTypes, func(rec *core.Record) {
			rec.Set("slug", momentType.Slug)
			rec.Set("label", momentType.Label)
			rec.Set("description", momentType.Description)
			rec.Set("color", momentType.Color)
			rec.Set("position", momentType.Position)
		})
		if err != nil {
			return counts, fmt.Errorf("moment type %s: %w", momentType.Slug, err)
		}
		typeIDs[momentType.Slug] = id
		counts.Types++
	}

	episodeIDs := map[string]string{}
	for _, episode := range data.Episodes {
		episodeID, err := insert(app, ColEpisodes, func(rec *core.Record) {
			rec.Set("youtube_id", episode.YouTubeID)
			rec.Set("number", episode.Number)
			rec.Set("slug", episode.Slug)
			rec.Set("title", episode.Title)
			setDate(rec, "published_at", episode.PublishedAt)
			rec.Set("duration_seconds", episode.DurationSeconds)
			rec.Set("thumbnail_url", episode.ThumbnailURL)
			rec.Set("description", episode.Description)
		})
		if err != nil {
			return counts, fmt.Errorf("episode %s: %w", episode.Slug, err)
		}
		episodeIDs[episode.Slug] = episodeID
		counts.Episodes++

		for _, appearance := range episode.Appearances {
			if _, err := insert(app, ColAppearances, func(rec *core.Record) {
				rec.Set("episode", episodeID)
				rec.Set("person", personIDs[appearance.PersonSlug])
				rec.Set("role", string(appearance.Role))
				rec.Set("is_interview", appearance.IsInterview)
			}); err != nil {
				return counts, fmt.Errorf("appearance %s/%s: %w", episode.Slug, appearance.PersonSlug, err)
			}
			counts.Appearances++
		}

		for _, moment := range episode.Moments {
			if _, err := insert(app, ColMoments, func(rec *core.Record) {
				rec.Set("episode", episodeID)
				rec.Set("type", typeIDs[moment.TypeSlug])
				setDate(rec, "at", episode.PublishedAt)
				rec.Set("video_timestamp", moment.VideoTimestamp)
				rec.Set("actor", personIDs[moment.ActorSlug])
				rec.Set("summary", moment.Summary)
				rec.Set("trigger_word", moment.TriggerWord)
				if id, ok := songIDs[moment.SongKey]; ok {
					rec.Set("song", id)
				}
				// Everything in the seed was typed by a human, which is what
				// the real MVP looks like too.
				rec.Set("source", string(domain.SourceManual))
				setJSON(rec, "payload", map[string]any{})
			}); err != nil {
				return counts, fmt.Errorf("moment %s/%s: %w", episode.Slug, moment.TypeSlug, err)
			}
			counts.Moments++
		}

		for _, opening := range episode.Openings {
			if _, err := insert(app, ColOpenings, func(rec *core.Record) {
				rec.Set("title", opening.Title)
				rec.Set("author_name", opening.AuthorName)
				rec.Set("author_handle", opening.AuthorHandle)
				rec.Set("genre", opening.Genre)
				rec.Set("episode", episodeID)
				rec.Set("at_seconds", opening.AtSeconds)
				setDate(rec, "aired_at", episode.PublishedAt)
			}); err != nil {
				return counts, fmt.Errorf("opening %s: %w", opening.Title, err)
			}
			counts.Openings++
		}
	}

	for _, entry := range data.Archive {
		if _, err := insert(app, ColArchive, func(rec *core.Record) {
			rec.Set("slug", entry.Slug)
			rec.Set("title", entry.Title)
			rec.Set("summary", entry.Summary)
			rec.Set("body", entry.Body)
			rec.Set("kind", string(entry.Kind))
			setJSON(rec, "tags", entry.Tags)
			if id, ok := episodeIDs[entry.EpisodeSlug]; ok {
				rec.Set("episode", id)
			}
			if id, ok := personIDs[entry.PersonSlug]; ok {
				rec.Set("person", id)
			}
			if id, ok := typeIDs[entry.MomentTypeSlug]; ok {
				rec.Set("moment_type", id)
			}
			rec.Set("source_url", entry.SourceURL)
			setDate(rec, "published_at", entry.PublishedAt)
		}); err != nil {
			return counts, fmt.Errorf("archive %s: %w", entry.Slug, err)
		}
		counts.Archive++
	}

	return counts, nil
}

func insert(app core.App, collection string, fill func(*core.Record)) (string, error) {
	c, err := app.FindCollectionByNameOrId(collection)
	if err != nil {
		return "", mapErr(err)
	}
	rec := core.NewRecord(c)
	fill(rec)
	if err := app.Save(rec); err != nil {
		return "", err
	}
	return rec.Id, nil
}

// wipe clears the seeded collections in dependency order, children first, so
// a cascade never fires against a row that is about to be deleted anyway.
func wipe(app core.App) error {
	for _, name := range []string{
		ColArchive, ColOpenings, ColMoments, ColAppearances,
		ColEpisodes, ColMomentTypes, ColPeople, ColSongs,
	} {
		records, err := app.FindAllRecords(name)
		if err != nil {
			return err
		}
		for _, rec := range records {
			if err := app.Delete(rec); err != nil {
				return fmt.Errorf("delete %s/%s: %w", name, rec.Id, err)
			}
		}
	}
	return nil
}
