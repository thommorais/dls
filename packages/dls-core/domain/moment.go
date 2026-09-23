package domain

import "time"

// MomentType is a row, not a constant. Every running joke the show has is one
// record here, so adding "they argue about football" costs an insert rather
// than a migration, and every filter, label and chart colour reads from the
// same place instead of drifting apart.
type MomentType struct {
	ID          MomentTypeID
	Slug        string
	Label       string
	Description string
	// Color is served by the API so a type keeps the same colour in every
	// chart on every page.
	Color string
	// Position orders the type on the site. Ties fall back to the label.
	Position int
}

// Source records who spotted the moment. Manual is a human typing it into the
// admin; the others are the crons that will write into the same table later,
// and Confidence is only meaningful for those.
type Source string

const (
	SourceManual Source = "manual"
	SourceLLM    Source = "llm"
	SourceAlgo   Source = "algo"
)

// Moment is one countable thing that happened at one second of one episode.
// The spine (episode, type, when, who, summary) is the same for every type;
// TriggerWord and SongID are filled only by the types that have a song in
// them, and Payload holds whatever a future type needs before it earns a
// column of its own.
type Moment struct {
	ID        MomentID
	EpisodeID EpisodeID
	TypeID    MomentTypeID

	// Type is resolved by the adapter. Nothing above that layer joins.
	Type MomentType

	// At is when it happened. It is usually the episode's own date, but a
	// show recorded on one day and published on another keeps both honest.
	At time.Time
	// VideoTimestamp is seconds into the video, which is what the site turns
	// into a YouTube deep link.
	VideoTimestamp int

	ActorID PersonID
	Summary string

	TriggerWord string
	SongID      SongID

	ClipURL    string
	Source     Source
	Confidence float64

	Payload map[string]any
}
