package domain

// Named types (not aliases) for every identifier: passing an EpisodeID where a
// PersonID is expected is a compile error, which catches wiring mistakes that
// string-typed IDs would let through silently.
type (
	EpisodeID    string
	PersonID     string
	MomentID     string
	MomentTypeID string
	SongID       string
	OpeningID    string
	AppearanceID string
)
