package domain

// Tally is the set of running counters the site exists to show. It is derived
// on read, never persisted, so a correction to one row is reflected
// everywhere without a backfill.
//
// Moments are counted by type slug rather than by named fields: the types are
// rows, so the counters cannot be known at compile time.
type Tally struct {
	WomenInterviewed int
	OpeningsAired    int
	Moments          map[string]int
}

// Of reads one moment counter. A type with nothing in range reads as zero
// rather than being absent, so a caller never has to check the map first.
func (t Tally) Of(slug string) int {
	return t.Moments[slug]
}

// TotalMoments is every counted moment regardless of type.
func (t Tally) TotalMoments() int {
	total := 0
	for _, count := range t.Moments {
		total += count
	}
	return total
}

// Count is one row of a ranking: a key that identifies the thing, a label to
// print, and how many times it happened.
type Count struct {
	Key   string
	Label string
	Count int
}
