package domain

import "time"

// ArchiveKind separates the shelves of the knowledge base. A glossary entry
// explains a running joke, a segment entry explains a recurring bit, trivia
// is a fact worth knowing, a bio is a person, and a milestone is a date the
// channel remembers.
type ArchiveKind string

const (
	// About is the channel explaining itself: what the show is, and how the
	// numbers on the site are counted.
	ArchiveAbout     ArchiveKind = "about"
	ArchiveGlossary  ArchiveKind = "glossary"
	ArchiveSegment   ArchiveKind = "segment"
	ArchiveTrivia    ArchiveKind = "trivia"
	ArchiveBio       ArchiveKind = "bio"
	ArchiveMilestone ArchiveKind = "milestone"
)

// ArchiveEntry is one page of the channel's knowledge base. The optional
// relations are what let a stat link to its own explanation: a counter on the
// home page points at the glossary entry for its moment type, and an episode
// page shows the milestones that happened in it.
type ArchiveEntry struct {
	ID      ArchiveID
	Slug    string
	Title   string
	Summary string
	Body    string
	Kind    ArchiveKind
	Tags    []string

	EpisodeID    EpisodeID
	PersonID     PersonID
	MomentTypeID MomentTypeID

	SourceURL   string
	PublishedAt time.Time
}
