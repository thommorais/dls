package domain

import "time"

// OpeningStatus tracks an audience-sent opening from inbox to air.
type OpeningStatus string

const (
	OpeningReceived OpeningStatus = "received"
	OpeningAired    OpeningStatus = "aired"
	OpeningArchived OpeningStatus = "archived"
)

// Opening is one intro sent in by a viewer. The ones that aired carry the
// episode and the second they played at, so the site can link into the video.
type Opening struct {
	ID           OpeningID
	Title        string
	AuthorName   string
	AuthorHandle string
	Genre        string
	EpisodeID    EpisodeID
	AtSeconds    int
	SentAt       time.Time
	AiredAt      *time.Time
	MediaURL     string
	Status       OpeningStatus
}
