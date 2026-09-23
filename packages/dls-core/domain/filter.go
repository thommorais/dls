package domain

import "time"

// Filter narrows every read the API serves. The zero value means everything,
// and each field is applied only when set, so one type covers the overview,
// the episode list, the openings library and the rankings.
type Filter struct {
	// From and To bound the episode publish date, inclusive. Moments and
	// appearances are scoped through the episode they belong to, because the
	// date a reader means is the date of the show, not the date someone typed
	// the row in.
	From *time.Time
	To   *time.Time

	// Genre narrows the openings library. It is ignored everywhere else.
	Genre string

	// Kind narrows the archive to one shelf. It is ignored everywhere else.
	Kind string

	// Limit caps the rows of each ranking. Zero means no cap.
	Limit int
}

// Scoped reports whether the filter narrows the set of episodes. An unscoped
// read asks each repository for everything in one call instead of naming
// every episode id it just loaded.
func (f Filter) Scoped() bool {
	return f.From != nil || f.To != nil
}
