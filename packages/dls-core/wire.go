// Package dls assembles the hexagon. It is the only place that knows both
// the services and the adapters: everything else depends on ports.
package dls

import (
	pbcore "github.com/pocketbase/pocketbase/core"

	"dls/dls-core/adapters/pb"
	"dls/dls-core/ports"
	"dls/dls-core/services"
)

type App struct {
	Stats ports.StatsUseCase
}

func New(app pbcore.App) *App {
	types := pb.NewMomentTypeRepository(app)

	return &App{
		Stats: services.NewStatsService(
			pb.NewEpisodeRepository(app),
			types,
			pb.NewMomentRepository(app, types),
			pb.NewAppearanceRepository(app),
			pb.NewOpeningRepository(app),
			pb.NewPersonRepository(app),
			pb.NewSongRepository(app),
			pb.NewArchiveRepository(app),
		),
	}
}

// Migrate installs the schema. It is safe on every boot.
func Migrate(app pbcore.App) error {
	return pb.Register(app)
}
