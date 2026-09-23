package httpapi

import (
	"time"

	"dls/dls-core/domain"
	"dls/dls-core/ports"
)

// The view types are the API contract. They are deliberately separate from
// the domain models, so a field rename inside the hexagon does not silently
// break every client.

type tallyView struct {
	WomenInterviewed int            `json:"women_interviewed"`
	OpeningsAired    int            `json:"openings_aired"`
	Moments          map[string]int `json:"moments"`
	TotalMoments     int            `json:"total_moments"`
}

func toTallyView(tally domain.Tally) tallyView {
	moments := tally.Moments
	if moments == nil {
		moments = map[string]int{}
	}
	return tallyView{
		WomenInterviewed: tally.WomenInterviewed,
		OpeningsAired:    tally.OpeningsAired,
		Moments:          moments,
		TotalMoments:     tally.TotalMoments(),
	}
}

type momentTypeView struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Position    int    `json:"position"`
}

func toMomentTypeView(t domain.MomentType) momentTypeView {
	return momentTypeView{
		ID: string(t.ID), Slug: t.Slug, Label: t.Label,
		Description: t.Description, Color: t.Color, Position: t.Position,
	}
}

func toMomentTypeViews(types []domain.MomentType) []momentTypeView {
	out := make([]momentTypeView, 0, len(types))
	for _, t := range types {
		out = append(out, toMomentTypeView(t))
	}
	return out
}

type episodeView struct {
	ID              string    `json:"id"`
	Number          int       `json:"number"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	YouTubeID       string    `json:"youtube_id"`
	PublishedAt     time.Time `json:"published_at"`
	DurationSeconds int       `json:"duration_seconds"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Description     string    `json:"description,omitempty"`
	Tally           tallyView `json:"tally"`
}

func toEpisodeView(e domain.Episode) episodeView {
	return episodeView{
		ID: string(e.ID), Number: e.Number, Slug: e.Slug, Title: e.Title,
		YouTubeID: e.YouTubeID, PublishedAt: e.PublishedAt,
		DurationSeconds: e.DurationSeconds, ThumbnailURL: e.ThumbnailURL,
		Description: e.Description, Tally: toTallyView(e.Tally),
	}
}

func toEpisodeViews(episodes []domain.Episode) []episodeView {
	out := make([]episodeView, 0, len(episodes))
	for _, e := range episodes {
		out = append(out, toEpisodeView(e))
	}
	return out
}

func toEpisodeViewPtr(e *domain.Episode) *episodeView {
	if e == nil {
		return nil
	}
	view := toEpisodeView(*e)
	return &view
}

type momentView struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	TypeLabel      string `json:"type_label"`
	Color          string `json:"color"`
	VideoTimestamp int    `json:"video_timestamp"`
	Actor          string `json:"actor,omitempty"`
	Summary        string `json:"summary,omitempty"`
	TriggerWord    string `json:"trigger_word,omitempty"`
	Song           string `json:"song,omitempty"`
	ClipURL        string `json:"clip_url,omitempty"`
	Source         string `json:"source,omitempty"`
}

func toMomentViews(moments []domain.Moment) []momentView {
	out := make([]momentView, 0, len(moments))
	for _, m := range moments {
		out = append(out, momentView{
			ID: string(m.ID), Type: m.Type.Slug, TypeLabel: m.Type.Label, Color: m.Type.Color,
			VideoTimestamp: m.VideoTimestamp, Actor: string(m.ActorID), Summary: m.Summary,
			TriggerWord: m.TriggerWord, Song: string(m.SongID), ClipURL: m.ClipURL,
			Source: string(m.Source),
		})
	}
	return out
}

type personView struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Gender string `json:"gender"`
	Photo  string `json:"photo,omitempty"`
}

func toPersonViews(people []domain.Person) []personView {
	out := make([]personView, 0, len(people))
	for _, p := range people {
		out = append(out, personView{
			ID: string(p.ID), Slug: p.Slug, Name: p.Name,
			Kind: string(p.Kind), Gender: string(p.Gender), Photo: p.Photo,
		})
	}
	return out
}

type appearanceView struct {
	Person      string `json:"person"`
	Role        string `json:"role"`
	IsInterview bool   `json:"is_interview"`
}

func toAppearanceViews(appearances []domain.Appearance) []appearanceView {
	out := make([]appearanceView, 0, len(appearances))
	for _, a := range appearances {
		out = append(out, appearanceView{
			Person: string(a.PersonID), Role: string(a.Role), IsInterview: a.IsInterview,
		})
	}
	return out
}

type songView struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist,omitempty"`
	Genre  string `json:"genre,omitempty"`
	Year   int    `json:"year,omitempty"`
	Link   string `json:"link,omitempty"`
}

func toSongViews(songs []domain.Song) []songView {
	out := make([]songView, 0, len(songs))
	for _, s := range songs {
		out = append(out, songView{
			ID: string(s.ID), Title: s.Title, Artist: s.Artist,
			Genre: s.Genre, Year: s.Year, Link: s.Link,
		})
	}
	return out
}

type openingView struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	AuthorName   string     `json:"author_name"`
	AuthorHandle string     `json:"author_handle,omitempty"`
	Genre        string     `json:"genre,omitempty"`
	Episode      string     `json:"episode,omitempty"`
	AtSeconds    int        `json:"at_seconds,omitempty"`
	SentAt       time.Time  `json:"sent_at"`
	AiredAt      *time.Time `json:"aired_at,omitempty"`
	MediaURL     string     `json:"media_url,omitempty"`
	Status       string     `json:"status"`
}

func toOpeningViews(openings []domain.Opening) []openingView {
	out := make([]openingView, 0, len(openings))
	for _, o := range openings {
		out = append(out, openingView{
			ID: string(o.ID), Title: o.Title, AuthorName: o.AuthorName,
			AuthorHandle: o.AuthorHandle, Genre: o.Genre, Episode: string(o.EpisodeID),
			AtSeconds: o.AtSeconds, SentAt: o.SentAt, AiredAt: o.AiredAt,
			MediaURL: o.MediaURL, Status: string(o.Status),
		})
	}
	return out
}

type archiveView struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary,omitempty"`
	Body        string    `json:"body,omitempty"`
	Kind        string    `json:"kind"`
	Tags        []string  `json:"tags"`
	Episode     string    `json:"episode,omitempty"`
	Person      string    `json:"person,omitempty"`
	MomentType  string    `json:"moment_type,omitempty"`
	SourceURL   string    `json:"source_url,omitempty"`
	PublishedAt time.Time `json:"published_at"`
}

func toArchiveView(e domain.ArchiveEntry) archiveView {
	tags := e.Tags
	if tags == nil {
		tags = []string{}
	}
	return archiveView{
		ID: string(e.ID), Slug: e.Slug, Title: e.Title, Summary: e.Summary,
		Body: e.Body, Kind: string(e.Kind), Tags: tags,
		Episode: string(e.EpisodeID), Person: string(e.PersonID),
		MomentType: string(e.MomentTypeID), SourceURL: e.SourceURL,
		PublishedAt: e.PublishedAt,
	}
}

// The list view drops the body: an index of twenty entries should not ship
// twenty articles.
func toArchiveViews(entries []domain.ArchiveEntry) []archiveView {
	out := make([]archiveView, 0, len(entries))
	for _, entry := range entries {
		view := toArchiveView(entry)
		view.Body = ""
		out = append(out, view)
	}
	return out
}

type countView struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

func toCountViews(counts []domain.Count) []countView {
	out := make([]countView, 0, len(counts))
	for _, c := range counts {
		out = append(out, countView{Key: c.Key, Label: c.Label, Count: c.Count})
	}
	return out
}

type recordView struct {
	Type    string `json:"type"`
	Count   int    `json:"count"`
	Episode string `json:"episode"`
	Title   string `json:"title"`
}

type streakView struct {
	Type   string `json:"type"`
	Length int    `json:"length"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type averageView struct {
	Type       string  `json:"type"`
	Total      int     `json:"total"`
	Episodes   int     `json:"episodes"`
	PerEpisode float64 `json:"per_episode"`
}

func toRankingsView(r ports.Rankings) map[string]any {
	actors := make(map[string][]countView, len(r.Actors))
	for slug, board := range r.Actors {
		actors[slug] = toCountViews(board)
	}

	records := make([]recordView, 0, len(r.Records))
	for _, record := range r.Records {
		records = append(records, recordView{
			Type: record.TypeSlug, Count: record.Count,
			Episode: string(record.Episode), Title: record.Title,
		})
	}

	streaks := make([]streakView, 0, len(r.Streaks))
	for _, streak := range r.Streaks {
		streaks = append(streaks, streakView{
			Type: streak.TypeSlug, Length: streak.Length,
			From: string(streak.From), To: string(streak.To),
		})
	}

	averages := make([]averageView, 0, len(r.Averages))
	for _, average := range r.Averages {
		averages = append(averages, averageView{
			Type: average.TypeSlug, Total: average.Total,
			Episodes: average.Episodes, PerEpisode: average.PerEpisode,
		})
	}

	return map[string]any{
		"types":         toMomentTypeViews(r.Types),
		"songs":         toCountViews(r.Songs),
		"trigger_words": toCountViews(r.TriggerWords),
		"guests":        toCountViews(r.Guests),
		"actors":        actors,
		"records":       records,
		"streaks":       streaks,
		"averages":      averages,
	}
}
