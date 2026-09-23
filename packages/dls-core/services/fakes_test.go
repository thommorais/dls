package services_test

import (
	"context"
	"fmt"

	"dls/dls-core/domain"
)

// In-memory doubles for the driven ports. They enforce only what storage
// enforces (the filter, the scope, existence), never a counting rule, so a
// test failing here means the service under test got the reading wrong.

type failer struct {
	// failOn forces the named method to return an error, to test propagation.
	failOn string
}

func (f failer) fail(method string) error {
	if f.failOn == method {
		return fmt.Errorf("storage exploded in %s", method)
	}
	return nil
}

func inScope(episode domain.EpisodeID, scope []domain.EpisodeID) bool {
	if scope == nil {
		return true
	}
	for _, id := range scope {
		if id == episode {
			return true
		}
	}
	return false
}

type fakeEpisodes struct {
	failer
	items      []domain.Episode
	lastFilter domain.Filter
}

func (r *fakeEpisodes) List(_ context.Context, filter domain.Filter) ([]domain.Episode, error) {
	r.lastFilter = filter
	if err := r.fail("List"); err != nil {
		return nil, err
	}
	out := []domain.Episode{}
	for _, episode := range r.items {
		if filter.From != nil && episode.PublishedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && episode.PublishedAt.After(*filter.To) {
			continue
		}
		out = append(out, episode)
	}
	return out, nil
}

func (r *fakeEpisodes) GetBySlug(_ context.Context, slug string) (domain.Episode, error) {
	if err := r.fail("GetBySlug"); err != nil {
		return domain.Episode{}, err
	}
	for _, episode := range r.items {
		if episode.Slug == slug {
			return episode, nil
		}
	}
	return domain.Episode{}, domain.ErrNotFound
}

type fakeMomentTypes struct {
	failer
	items []domain.MomentType
}

func (r *fakeMomentTypes) List(_ context.Context) ([]domain.MomentType, error) {
	if err := r.fail("List"); err != nil {
		return nil, err
	}
	return r.items, nil
}

type fakeMoments struct {
	failer
	items     []domain.Moment
	lastScope []domain.EpisodeID
	scoped    bool
}

func (r *fakeMoments) ListByEpisodes(_ context.Context, episodes []domain.EpisodeID) ([]domain.Moment, error) {
	r.lastScope, r.scoped = episodes, true
	if err := r.fail("ListByEpisodes"); err != nil {
		return nil, err
	}
	out := []domain.Moment{}
	for _, moment := range r.items {
		if inScope(moment.EpisodeID, episodes) {
			out = append(out, moment)
		}
	}
	return out, nil
}

type fakeAppearances struct {
	failer
	items []domain.Appearance
}

func (r *fakeAppearances) ListByEpisodes(_ context.Context, episodes []domain.EpisodeID) ([]domain.Appearance, error) {
	if err := r.fail("ListByEpisodes"); err != nil {
		return nil, err
	}
	out := []domain.Appearance{}
	for _, appearance := range r.items {
		if inScope(appearance.EpisodeID, episodes) {
			out = append(out, appearance)
		}
	}
	return out, nil
}

type fakeOpenings struct {
	failer
	items      []domain.Opening
	lastFilter domain.Filter
}

func (r *fakeOpenings) ListByEpisodes(_ context.Context, episodes []domain.EpisodeID) ([]domain.Opening, error) {
	if err := r.fail("ListByEpisodes"); err != nil {
		return nil, err
	}
	out := []domain.Opening{}
	for _, opening := range r.items {
		if inScope(opening.EpisodeID, episodes) {
			out = append(out, opening)
		}
	}
	return out, nil
}

func (r *fakeOpenings) List(_ context.Context, filter domain.Filter) ([]domain.Opening, error) {
	r.lastFilter = filter
	if err := r.fail("List"); err != nil {
		return nil, err
	}
	out := []domain.Opening{}
	for _, opening := range r.items {
		if filter.Genre != "" && opening.Genre != filter.Genre {
			continue
		}
		out = append(out, opening)
	}
	return out, nil
}

type fakePeople struct {
	failer
	items []domain.Person
}

func (r *fakePeople) List(_ context.Context) ([]domain.Person, error) {
	if err := r.fail("List"); err != nil {
		return nil, err
	}
	return r.items, nil
}

type fakeSongs struct {
	failer
	items []domain.Song
}

func (r *fakeSongs) List(_ context.Context) ([]domain.Song, error) {
	if err := r.fail("List"); err != nil {
		return nil, err
	}
	return r.items, nil
}
