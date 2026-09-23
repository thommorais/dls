package domain

// Song is the library of music the hosts break into. It is a fact table: the
// same song is referenced by every moment that triggered it.
type Song struct {
	ID     SongID
	Title  string
	Artist string
	Genre  string
	Year   int
	Link   string
}
