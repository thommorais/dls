package domain

import "time"

type Episode struct {
	ID              EpisodeID
	YouTubeID       string
	Number          int
	Slug            string
	Title           string
	PublishedAt     time.Time
	DurationSeconds int
	ThumbnailURL    string
	Description     string

	// Derived on read, never persisted.
	Tally Tally
}
