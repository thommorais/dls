// Package rules holds the folds that turn stored rows into the show's
// counters. Everything here is a pure function over slices: the repositories
// read, these functions count, and nothing in between touches a database.
package rules

import "dls/dls-core/domain"

// Tally counts every running joke across the rows it is given. Scoping (one
// episode, one date range) is the caller's job: pass the rows you want
// counted.
func Tally(
	moments []domain.Moment,
	appearances []domain.Appearance,
	people []domain.Person,
	openings []domain.Opening,
) domain.Tally {
	tally := domain.Tally{Moments: map[string]int{}}
	byID := peopleByID(people)

	for _, moment := range moments {
		// A moment whose type row is gone cannot be counted under any type.
		// Counting it under the empty slug would invent a counter nothing
		// can label.
		if moment.Type.Slug == "" {
			continue
		}
		tally.Moments[moment.Type.Slug]++
	}

	for _, appearance := range appearances {
		if isWomanInterviewed(appearance, byID) {
			tally.WomenInterviewed++
		}
	}

	tally.OpeningsAired = len(openings)

	return tally
}

// TallyByEpisode splits the same counters per episode, for the episode list
// and its charts. An episode with no rows at all is absent from the map
// rather than present and zero.
func TallyByEpisode(
	moments []domain.Moment,
	appearances []domain.Appearance,
	people []domain.Person,
	openings []domain.Opening,
) map[domain.EpisodeID]domain.Tally {
	byID := peopleByID(people)
	out := map[domain.EpisodeID]domain.Tally{}

	at := func(episode domain.EpisodeID) domain.Tally {
		tally, ok := out[episode]
		if !ok {
			tally = domain.Tally{Moments: map[string]int{}}
			out[episode] = tally
		}
		return tally
	}

	for _, moment := range moments {
		if moment.Type.Slug == "" {
			continue
		}
		// The map inside the tally is shared, so writing through it needs no
		// write back.
		at(moment.EpisodeID).Moments[moment.Type.Slug]++
	}

	for _, appearance := range appearances {
		if !isWomanInterviewed(appearance, byID) {
			continue
		}
		tally := at(appearance.EpisodeID)
		tally.WomenInterviewed++
		out[appearance.EpisodeID] = tally
	}

	for _, opening := range openings {
		tally := at(opening.EpisodeID)
		tally.OpeningsAired++
		out[opening.EpisodeID] = tally
	}

	return out
}

// isWomanInterviewed is the whole definition of the headline stat: a guest,
// actually interviewed, whose gender was recorded as woman. An appearance
// whose person is missing counts for nothing, and so does an unfilled gender.
func isWomanInterviewed(appearance domain.Appearance, people map[domain.PersonID]domain.Person) bool {
	if !appearance.IsInterview || appearance.Role != domain.RoleGuest {
		return false
	}
	person, ok := people[appearance.PersonID]
	return ok && person.Gender == domain.GenderWoman
}

func peopleByID(people []domain.Person) map[domain.PersonID]domain.Person {
	out := make(map[domain.PersonID]domain.Person, len(people))
	for _, person := range people {
		out[person.ID] = person
	}
	return out
}
