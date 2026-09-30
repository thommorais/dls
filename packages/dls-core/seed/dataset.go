package seed

import (
	"fmt"
	"sort"
	"time"

	"dls/dls-core/domain"
)

// Relations are carried by slug rather than by id: the dataset is written
// before anything exists, and the writer resolves them as it inserts.

type Moment struct {
	TypeSlug       string
	VideoTimestamp int
	ActorSlug      string
	Summary        string
	TriggerWord    string
	SongKey        string
}

type Appearance struct {
	PersonSlug  string
	Role        domain.AppearanceRole
	IsInterview bool
}

type Opening struct {
	Title        string
	AuthorName   string
	AuthorHandle string
	Genre        string
	AtSeconds    int
}

type Episode struct {
	Number          int
	Slug            string
	Title           string
	YouTubeID       string
	ThumbnailURL    string
	Description     string
	PublishedAt     time.Time
	DurationSeconds int

	Moments     []Moment
	Appearances []Appearance
	Openings    []Opening
}

type ArchiveEntry struct {
	Slug           string
	Title          string
	Summary        string
	Body           string
	Kind           domain.ArchiveKind
	Tags           []string
	EpisodeSlug    string
	PersonSlug     string
	MomentTypeSlug string
	SourceURL      string
	PublishedAt    time.Time
}

type Dataset struct {
	Types    []domain.MomentType
	People   []domain.Person
	Songs    []Song
	Episodes []Episode
	Archive  []ArchiveEntry
}

// Build returns the dataset. Episodes, their moments and their openings come
// from data.go, transcribed from the show's own spreadsheet (see
// apps/site/public). Titles, publish dates and durations aren't in that
// transcription yet, so they stay blank rather than invented.
func Build() Dataset {
	return Dataset{
		Types:    MomentTypes,
		People:   people,
		Songs:    songs,
		Episodes: buildEpisodes(),
		Archive:  archiveEntries(),
	}
}

func buildEpisodes() []Episode {
	byNumber := make(map[int]*Episode, len(rawEpisodes))
	episodes := make([]Episode, len(rawEpisodes))
	for i, raw := range rawEpisodes {
		episodes[i] = Episode{
			Number:       raw.Number,
			Slug:         fmt.Sprintf("ep-%d", raw.Number),
			Title:        fmt.Sprintf("Episódio %d", raw.Number),
			YouTubeID:    raw.YouTubeID,
			ThumbnailURL: fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", raw.YouTubeID),
		}
		byNumber[raw.Number] = &episodes[i]
	}

	for _, raw := range rawMoments {
		episode, ok := byNumber[raw.EpisodeNumber]
		if !ok {
			continue
		}
		episode.Moments = append(episode.Moments, Moment{
			TypeSlug:       raw.TypeSlug,
			VideoTimestamp: raw.AtSeconds,
			Summary:        raw.Summary,
		})
	}

	for _, raw := range rawOpenings {
		episode, ok := byNumber[raw.EpisodeNumber]
		if !ok {
			continue
		}
		episode.Openings = append(episode.Openings, Opening{
			Title:        raw.Title,
			AuthorName:   raw.Author,
			AuthorHandle: raw.Handle,
			Genre:        raw.Genre,
			AtSeconds:    raw.AtSeconds,
		})
	}

	sort.Slice(episodes, func(i, j int) bool { return episodes[i].Number < episodes[j].Number })
	return episodes
}
