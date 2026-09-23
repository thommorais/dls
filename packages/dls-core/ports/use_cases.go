package ports

import (
	"context"

	"dls/dls-core/domain"
)

// StatsUseCase is the whole driving surface. Every method reads: the MVP is
// written in the PocketBase admin, so nothing here mutates.
type StatsUseCase interface {
	Overview(ctx context.Context, filter domain.Filter) (Overview, error)
	Episodes(ctx context.Context, filter domain.Filter) ([]domain.Episode, error)
	Episode(ctx context.Context, slug string) (EpisodeDetail, error)
	Openings(ctx context.Context, filter domain.Filter) ([]domain.Opening, error)
	Rankings(ctx context.Context, filter domain.Filter) (Rankings, error)
}

// Overview is the home page: the counters, plus enough context to say how
// much show they were counted over.
type Overview struct {
	Tally domain.Tally
	// Types is every moment type there is, ordered, so the page can label and
	// colour a counter that has nothing in it yet instead of hiding it.
	Types    []domain.MomentType
	Episodes int
	First    *domain.Episode
	Latest   *domain.Episode
}

// EpisodeDetail carries everything one episode page renders, so the page is
// one request. People and Songs are the rows the moments and appearances
// point at, not the whole library.
type EpisodeDetail struct {
	Episode     domain.Episode
	Moments     []domain.Moment
	Appearances []domain.Appearance
	People      []domain.Person
	Songs       []domain.Song
	Openings    []domain.Opening
}

// Rankings is the leaderboard page. Each board is already sorted and capped.
type Rankings struct {
	Songs        []domain.Count
	TriggerWords []domain.Count
	Guests       []domain.Count
	// Actors is one leaderboard per moment type, keyed by type slug: who asks
	// the most, who sings the most, and whatever type gets added next. A type
	// nobody has done yet is present with no rows.
	Actors map[string][]domain.Count
	Types  []domain.MomentType
}
