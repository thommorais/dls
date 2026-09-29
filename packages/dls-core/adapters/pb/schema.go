package pb

import (
	"github.com/pocketbase/pocketbase/core"
)

func strPtr(s string) *string { return &s }

func autodates() []core.Field {
	return []core.Field{
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	}
}

const slugPattern = `^[a-z0-9]+(-[a-z0-9]+)*$`

func ensureSongs(app core.App) error {
	if _, ok := find(app, ColSongs); ok {
		return nil
	}
	c := core.NewBaseCollection(ColSongs)
	c.Fields.Add(
		&core.TextField{Name: "title", Required: true, Max: 200, Presentable: true},
		&core.TextField{Name: "artist", Max: 200, Presentable: true},
		&core.TextField{Name: "genre", Max: 80},
		&core.NumberField{Name: "year"},
		&core.URLField{Name: "link"},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_songs_title", false, "title", "")

	return app.Save(c)
}

func ensurePeople(app core.App) error {
	if _, ok := find(app, ColPeople); ok {
		return nil
	}
	c := core.NewBaseCollection(ColPeople)
	c.Fields.Add(
		&core.TextField{Name: "slug", Required: true, Max: 80, Pattern: slugPattern},
		&core.TextField{Name: "name", Required: true, Max: 120, Presentable: true},
		&core.SelectField{Name: "kind", Required: true, MaxSelect: 1, Values: []string{"host", "guest", "staff"}},
		// Not required, and "unknown" is the default the counter refuses to
		// read as a woman, so an unfilled row is never counted by accident.
		&core.SelectField{Name: "gender", MaxSelect: 1, Values: []string{"woman", "man", "nonbinary", "unknown"}},
		&core.TextField{Name: "photo", Max: 300},
		&core.JSONField{Name: "socials", MaxSize: 20000},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_people_slug", true, "slug", "")

	return app.Save(c)
}

func ensureMomentTypes(app core.App) error {
	if _, ok := find(app, ColMomentTypes); ok {
		return nil
	}
	c := core.NewBaseCollection(ColMomentTypes)
	c.Fields.Add(
		&core.TextField{Name: "slug", Required: true, Max: 80, Pattern: slugPattern},
		&core.TextField{Name: "label", Required: true, Max: 120, Presentable: true},
		&core.TextField{Name: "description", Max: 500},
		// The colour lives with the type so a stat keeps the same colour in
		// every chart on every page.
		&core.TextField{Name: "color", Max: 20},
		&core.NumberField{Name: "position"},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_moment_types_slug", true, "slug", "")

	return app.Save(c)
}

func ensureEpisodes(app core.App) error {
	if _, ok := find(app, ColEpisodes); ok {
		return nil
	}
	c := core.NewBaseCollection(ColEpisodes)
	c.Fields.Add(
		&core.TextField{Name: "youtube_id", Max: 40},
		&core.NumberField{Name: "number"},
		&core.TextField{Name: "slug", Required: true, Max: 120, Pattern: slugPattern},
		&core.TextField{Name: "title", Required: true, Max: 300, Presentable: true},
		&core.DateField{Name: "published_at"},
		&core.NumberField{Name: "duration_seconds"},
		&core.TextField{Name: "thumbnail_url", Max: 500},
		&core.EditorField{Name: "description"},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_episodes_slug", true, "slug", "")
	c.AddIndex("idx_dls_episodes_published", false, "published_at", "")

	return app.Save(c)
}

func ensureAppearances(app core.App) error {
	if _, ok := find(app, ColAppearances); ok {
		return nil
	}
	episodes, err := app.FindCollectionByNameOrId(ColEpisodes)
	if err != nil {
		return err
	}
	people, err := app.FindCollectionByNameOrId(ColPeople)
	if err != nil {
		return err
	}

	c := core.NewBaseCollection(ColAppearances)
	c.Fields.Add(
		&core.RelationField{Name: "episode", Required: true, CollectionId: episodes.Id, CascadeDelete: true, MaxSelect: 1},
		&core.RelationField{Name: "person", Required: true, CollectionId: people.Id, CascadeDelete: true, MaxSelect: 1},
		&core.SelectField{Name: "role", Required: true, MaxSelect: 1, Values: []string{"host", "guest", "remote"}},
		&core.BoolField{Name: "is_interview"},
	)
	c.Fields.Add(autodates()...)
	// One row per person per episode: two rows would double every count that
	// reads this table.
	c.AddIndex("idx_dls_appearances_unique", true, "episode, person", "")
	c.AddIndex("idx_dls_appearances_episode", false, "episode", "")

	return app.Save(c)
}

func ensureMoments(app core.App) error {
	if _, ok := find(app, ColMoments); ok {
		return nil
	}
	episodes, err := app.FindCollectionByNameOrId(ColEpisodes)
	if err != nil {
		return err
	}
	types, err := app.FindCollectionByNameOrId(ColMomentTypes)
	if err != nil {
		return err
	}
	people, err := app.FindCollectionByNameOrId(ColPeople)
	if err != nil {
		return err
	}
	songs, err := app.FindCollectionByNameOrId(ColSongs)
	if err != nil {
		return err
	}

	c := core.NewBaseCollection(ColMoments)
	c.Fields.Add(
		&core.RelationField{Name: "episode", Required: true, CollectionId: episodes.Id, CascadeDelete: true, MaxSelect: 1},
		// The type is a row, so a new running joke is an insert here rather
		// than a migration.
		&core.RelationField{Name: "type", Required: true, CollectionId: types.Id, MaxSelect: 1},
		&core.DateField{Name: "at"},
		&core.NumberField{Name: "video_timestamp"},
		&core.RelationField{Name: "actor", CollectionId: people.Id, MaxSelect: 1},
		&core.TextField{Name: "summary", Max: 500, Presentable: true},
		&core.TextField{Name: "trigger_word", Max: 120},
		&core.RelationField{Name: "song", CollectionId: songs.Id, MaxSelect: 1},
		&core.URLField{Name: "clip_url"},
		&core.SelectField{Name: "source", MaxSelect: 1, Values: []string{"manual", "llm", "algo"}},
		&core.NumberField{Name: "confidence"},
		// Whatever a future type needs before it earns a column of its own.
		&core.JSONField{Name: "payload", MaxSize: 50000},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_moments_episode", false, "episode", "")
	c.AddIndex("idx_dls_moments_type", false, "type", "")
	c.AddIndex("idx_dls_moments_actor", false, "actor", "")

	return app.Save(c)
}

func ensureOpenings(app core.App) error {
	if _, ok := find(app, ColOpenings); ok {
		return nil
	}
	episodes, err := app.FindCollectionByNameOrId(ColEpisodes)
	if err != nil {
		return err
	}

	c := core.NewBaseCollection(ColOpenings)
	c.Fields.Add(
		&core.TextField{Name: "title", Required: true, Max: 200, Presentable: true},
		&core.TextField{Name: "author_name", Max: 120},
		&core.TextField{Name: "author_handle", Max: 120},
		&core.TextField{Name: "genre", Max: 80},
		// Empty until it airs, which is what separates the inbox from the
		// library.
		&core.RelationField{Name: "episode", CollectionId: episodes.Id, MaxSelect: 1},
		&core.NumberField{Name: "at_seconds"},
		&core.DateField{Name: "sent_at"},
		&core.DateField{Name: "aired_at"},
		&core.TextField{Name: "media_url", Max: 500},
		&core.SelectField{Name: "status", Required: true, MaxSelect: 1, Values: []string{"received", "aired", "archived"}},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_openings_episode", false, "episode", "")
	c.AddIndex("idx_dls_openings_status", false, "status", "")

	return app.Save(c)
}

func ensureArchive(app core.App) error {
	if _, ok := find(app, ColArchive); ok {
		return nil
	}
	episodes, err := app.FindCollectionByNameOrId(ColEpisodes)
	if err != nil {
		return err
	}
	people, err := app.FindCollectionByNameOrId(ColPeople)
	if err != nil {
		return err
	}
	types, err := app.FindCollectionByNameOrId(ColMomentTypes)
	if err != nil {
		return err
	}

	c := core.NewBaseCollection(ColArchive)
	c.Fields.Add(
		&core.TextField{Name: "slug", Required: true, Max: 120, Pattern: slugPattern},
		&core.TextField{Name: "title", Required: true, Max: 200, Presentable: true},
		&core.TextField{Name: "summary", Max: 500},
		&core.EditorField{Name: "body"},
		&core.SelectField{Name: "kind", Required: true, MaxSelect: 1, Values: []string{"about", "glossary", "segment", "trivia", "bio", "milestone"}},
		&core.JSONField{Name: "tags", MaxSize: 20000},
		// The optional relations are what let a stat link to its own
		// explanation.
		&core.RelationField{Name: "episode", CollectionId: episodes.Id, MaxSelect: 1},
		&core.RelationField{Name: "person", CollectionId: people.Id, MaxSelect: 1},
		&core.RelationField{Name: "moment_type", CollectionId: types.Id, MaxSelect: 1},
		&core.TextField{Name: "source_url", Max: 500},
		&core.DateField{Name: "published_at"},
	)
	c.Fields.Add(autodates()...)
	c.AddIndex("idx_dls_archive_slug", true, "slug", "")
	c.AddIndex("idx_dls_archive_kind", false, "kind", "")

	return app.Save(c)
}

// applyRules opens every collection for reading and closes it for writing.
// The site is public and read-only; everything is typed in the admin UI,
// which authenticates as a superuser and ignores these rules.
func applyRules(app core.App) error {
	for _, name := range []string{
		ColSongs, ColPeople, ColMomentTypes, ColEpisodes,
		ColAppearances, ColMoments, ColOpenings, ColArchive,
	} {
		c, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			return err
		}
		c.ListRule = strPtr("")
		c.ViewRule = strPtr("")
		c.CreateRule = nil
		c.UpdateRule = nil
		c.DeleteRule = nil
		if err := app.Save(c); err != nil {
			return err
		}
	}
	return nil
}
