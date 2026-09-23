package domain

// Record is the best single episode for one moment type: the box-score line
// that gets quoted. Ties go to the earliest episode, because the record
// belongs to whoever did it first.
type Record struct {
	TypeSlug string
	Count    int
	Episode  EpisodeID
	Title    string
}

// Streak is the longest run of consecutive episodes where a type happened at
// least once. Consecutive means consecutive in publish order, so a gap in
// episode numbering does not break a run that reads unbroken to a viewer.
type Streak struct {
	TypeSlug string
	Length   int
	From     EpisodeID
	To       EpisodeID
}

// Average is the per-episode rate of a type, which is what makes two ranges
// of different length comparable.
type Average struct {
	TypeSlug   string
	Total      int
	Episodes   int
	PerEpisode float64
}
