package pb

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// Register installs the dls schema. It is idempotent: an existing collection
// is left alone, so the app can call it on every boot.
//
// Order matters only because of the relations: a collection has to exist
// before another one can point at it.
func Register(app core.App) error {
	steps := []struct {
		name string
		run  func(core.App) error
	}{
		{"songs", ensureSongs},
		{"people", ensurePeople},
		{"moment types", ensureMomentTypes},
		{"episodes", ensureEpisodes},
		{"episode facts", ensureEpisodeFacts},
		{"appearances", ensureAppearances},
		{"moments", ensureMoments},
		{"openings", ensureOpenings},
		{"archive", ensureArchive},
	}

	for _, step := range steps {
		if err := step.run(app); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	if err := applyRules(app); err != nil {
		return fmt.Errorf("rules: %w", err)
	}
	return nil
}
