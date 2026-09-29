package seed

import (
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
	SentAt       time.Time
	Status       domain.OpeningStatus
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

// Build returns the dataset. There is no real episode data transcribed yet,
// so it seeds only the moment type vocabulary and leaves episodes, people,
// songs and archive entries empty until the real show's records are added.
func Build() Dataset {
	return Dataset{
		Types:    momentTypes,
		People:   people,
		Songs:    songs,
		Episodes: []Episode{},
		Archive:  archiveEntries(),
	}
}
