package domain

import "time"

// Opening is one intro sent in by a viewer that aired. It carries the episode
// and the second it played at, so the site can link into the video.
type Opening struct {
	ID           OpeningID
	Title        string
	AuthorName   string
	AuthorHandle string
	Genre        string
	EpisodeID    EpisodeID
	AtSeconds    int
	AiredAt      *time.Time
	MediaURL     string
}
