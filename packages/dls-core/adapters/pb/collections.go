// Package pb is the PocketBase adapter: the schema, and the repositories that
// read it. Records are mapped to domain models at this boundary, so nothing
// above this package ever sees a *core.Record.
package pb

import "github.com/pocketbase/pocketbase/core"

const (
	ColEpisodes    = "dls_episodes"
	ColPeople      = "dls_people"
	ColAppearances = "dls_appearances"
	ColMomentTypes = "dls_moment_types"
	ColMoments     = "dls_moments"
	ColSongs       = "dls_songs"
	ColOpenings    = "dls_openings"
	ColArchive     = "dls_archive"
)

func find(app core.App, name string) (*core.Collection, bool) {
	c, err := app.FindCollectionByNameOrId(name)
	if err != nil || c == nil {
		return nil, false
	}
	return c, true
}
