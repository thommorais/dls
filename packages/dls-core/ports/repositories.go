// Package ports declares the interfaces the hexagon is wired through: driven
// ports the adapters implement, driving ports the services implement.
package ports

import (
	"context"

	"dls/dls-core/domain"
)

// Every repository returns domain models. Nothing above the adapter layer
// ever sees a PocketBase record.
//
// The ListByEpisodes methods take nil to mean every episode, which is how an
// unfiltered read gets each table in a single query.

type EpisodeRepository interface {
	List(ctx context.Context, filter domain.Filter) ([]domain.Episode, error)
	GetBySlug(ctx context.Context, slug string) (domain.Episode, error)
}

type MomentTypeRepository interface {
	List(ctx context.Context) ([]domain.MomentType, error)
}

type MomentRepository interface {
	ListByEpisodes(ctx context.Context, episodes []domain.EpisodeID) ([]domain.Moment, error)
}

type AppearanceRepository interface {
	ListByEpisodes(ctx context.Context, episodes []domain.EpisodeID) ([]domain.Appearance, error)
}

type OpeningRepository interface {
	ListByEpisodes(ctx context.Context, episodes []domain.EpisodeID) ([]domain.Opening, error)
	List(ctx context.Context, filter domain.Filter) ([]domain.Opening, error)
}

type PersonRepository interface {
	List(ctx context.Context) ([]domain.Person, error)
}

type SongRepository interface {
	List(ctx context.Context) ([]domain.Song, error)
}
