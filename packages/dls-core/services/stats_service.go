// Package services implements the driving ports. Every method here is a read:
// it scopes the episodes, loads the rows that belong to them, and hands them
// to the pure folds in domain/rules.
package services

import (
	"context"

	"dls/dls-core/domain"
	"dls/dls-core/domain/rules"
	"dls/dls-core/ports"
)

// StatsService is the only implementation of the driving port.
var _ ports.StatsUseCase = (*StatsService)(nil)

type StatsService struct {
	episodes    ports.EpisodeRepository
	momentTypes ports.MomentTypeRepository
	moments     ports.MomentRepository
	appearances ports.AppearanceRepository
	openings    ports.OpeningRepository
	people      ports.PersonRepository
	songs       ports.SongRepository
}

func NewStatsService(
	episodes ports.EpisodeRepository,
	momentTypes ports.MomentTypeRepository,
	moments ports.MomentRepository,
	appearances ports.AppearanceRepository,
	openings ports.OpeningRepository,
	people ports.PersonRepository,
	songs ports.SongRepository,
) *StatsService {
	return &StatsService{
		episodes:    episodes,
		momentTypes: momentTypes,
		moments:     moments,
		appearances: appearances,
		openings:    openings,
		people:      people,
		songs:       songs,
	}
}

func (s *StatsService) Overview(ctx context.Context, filter domain.Filter) (ports.Overview, error) {
	episodes, rows, err := s.load(ctx, filter)
	if err != nil {
		return ports.Overview{}, err
	}
	types, err := s.types(ctx)
	if err != nil {
		return ports.Overview{}, err
	}

	return ports.Overview{
		Tally:    rules.Tally(rows.moments, rows.appearances, rows.people, rows.openings),
		Types:    types,
		Episodes: len(episodes),
		First:    edge(episodes, earliest),
		Latest:   edge(episodes, latest),
	}, nil
}

func (s *StatsService) Episodes(ctx context.Context, filter domain.Filter) ([]domain.Episode, error) {
	episodes, rows, err := s.load(ctx, filter)
	if err != nil {
		return nil, err
	}

	tallies := rules.TallyByEpisode(rows.moments, rows.appearances, rows.people, rows.openings)
	out := make([]domain.Episode, 0, len(episodes))
	for _, episode := range episodes {
		episode.Tally = tallies[episode.ID]
		out = append(out, episode)
	}
	return out, nil
}

func (s *StatsService) Episode(ctx context.Context, slug string) (ports.EpisodeDetail, error) {
	episode, err := s.episodes.GetBySlug(ctx, slug)
	if err != nil {
		return ports.EpisodeDetail{}, err
	}

	scope := []domain.EpisodeID{episode.ID}
	rows, err := s.gather(ctx, scope)
	if err != nil {
		return ports.EpisodeDetail{}, err
	}
	songs, err := s.songs.List(ctx)
	if err != nil {
		return ports.EpisodeDetail{}, err
	}

	episode.Tally = rules.Tally(rows.moments, rows.appearances, rows.people, rows.openings)

	return ports.EpisodeDetail{
		Episode:     episode,
		Moments:     rows.moments,
		Appearances: rows.appearances,
		Openings:    rows.openings,
		People:      peopleOf(rows, rows.people),
		Songs:       songsOf(rows.moments, songs),
	}, nil
}

func (s *StatsService) Openings(ctx context.Context, filter domain.Filter) ([]domain.Opening, error) {
	return s.openings.List(ctx, filter)
}

func (s *StatsService) Rankings(ctx context.Context, filter domain.Filter) (ports.Rankings, error) {
	_, rows, err := s.load(ctx, filter)
	if err != nil {
		return ports.Rankings{}, err
	}
	songs, err := s.songs.List(ctx)
	if err != nil {
		return ports.Rankings{}, err
	}
	types, err := s.types(ctx)
	if err != nil {
		return ports.Rankings{}, err
	}

	// One board per type, including the types nobody has done yet: a counter
	// the site advertises should not disappear from the rankings just because
	// it is still at zero.
	actors := make(map[string][]domain.Count, len(types))
	for _, momentType := range types {
		actors[momentType.Slug] = rules.ByActor(rows.moments, []string{momentType.Slug}, rows.people, filter.Limit)
	}

	return ports.Rankings{
		Songs:        rules.TopSongs(rows.moments, songs, filter.Limit),
		TriggerWords: rules.TopTriggerWords(rows.moments, filter.Limit),
		Guests:       rules.TopGuests(rows.appearances, rows.people, filter.Limit),
		Actors:       actors,
		Types:        types,
	}, nil
}

// types reads the moment types in the order the site shows them.
func (s *StatsService) types(ctx context.Context) ([]domain.MomentType, error) {
	types, err := s.momentTypes.List(ctx)
	if err != nil {
		return nil, err
	}
	return rules.OrderTypes(types), nil
}

// rowset is everything the folds need for one scope.
type rowset struct {
	moments     []domain.Moment
	appearances []domain.Appearance
	openings    []domain.Opening
	people      []domain.Person
}

// load resolves the filter to a set of episodes and reads the rows that
// belong to them.
func (s *StatsService) load(ctx context.Context, filter domain.Filter) ([]domain.Episode, rowset, error) {
	episodes, err := s.episodes.List(ctx, filter)
	if err != nil {
		return nil, rowset{}, err
	}
	rows, err := s.gather(ctx, scopeOf(filter, episodes))
	if err != nil {
		return nil, rowset{}, err
	}
	return episodes, rows, nil
}

func (s *StatsService) gather(ctx context.Context, scope []domain.EpisodeID) (rowset, error) {
	moments, err := s.moments.ListByEpisodes(ctx, scope)
	if err != nil {
		return rowset{}, err
	}
	appearances, err := s.appearances.ListByEpisodes(ctx, scope)
	if err != nil {
		return rowset{}, err
	}
	openings, err := s.openings.ListByEpisodes(ctx, scope)
	if err != nil {
		return rowset{}, err
	}
	people, err := s.people.List(ctx)
	if err != nil {
		return rowset{}, err
	}
	return rowset{moments: moments, appearances: appearances, openings: openings, people: people}, nil
}

// scopeOf names the episodes to read rows for. An unfiltered read returns nil,
// which the repositories take as "every episode": naming all of them would
// send the whole table back as a query parameter for no gain.
func scopeOf(filter domain.Filter, episodes []domain.Episode) []domain.EpisodeID {
	if !filter.Scoped() {
		return nil
	}
	ids := make([]domain.EpisodeID, 0, len(episodes))
	for _, episode := range episodes {
		ids = append(ids, episode.ID)
	}
	return ids
}

// peopleOf keeps the people an episode actually involved, so the detail
// payload does not carry the whole cast of the show.
func peopleOf(rows rowset, people []domain.Person) []domain.Person {
	wanted := map[domain.PersonID]bool{}
	for _, moment := range rows.moments {
		wanted[moment.ActorID] = true
	}
	for _, appearance := range rows.appearances {
		wanted[appearance.PersonID] = true
	}

	out := make([]domain.Person, 0, len(wanted))
	for _, person := range people {
		if wanted[person.ID] {
			out = append(out, person)
		}
	}
	return out
}

func songsOf(moments []domain.Moment, songs []domain.Song) []domain.Song {
	wanted := map[domain.SongID]bool{}
	for _, moment := range moments {
		wanted[moment.SongID] = true
	}

	out := make([]domain.Song, 0, len(wanted))
	for _, song := range songs {
		if wanted[song.ID] {
			out = append(out, song)
		}
	}
	return out
}

func earliest(a, b domain.Episode) bool { return a.PublishedAt.Before(b.PublishedAt) }
func latest(a, b domain.Episode) bool   { return a.PublishedAt.After(b.PublishedAt) }

// edge picks the episode at one end of the range. It returns nil rather than
// a zero episode, so an empty range reads as absent instead of as a show
// published at the zero time.
func edge(episodes []domain.Episode, better func(a, b domain.Episode) bool) *domain.Episode {
	if len(episodes) == 0 {
		return nil
	}
	pick := episodes[0]
	for _, episode := range episodes[1:] {
		if better(episode, pick) {
			pick = episode
		}
	}
	return &pick
}
