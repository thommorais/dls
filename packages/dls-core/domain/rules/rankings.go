package rules

import (
	"sort"
	"strings"

	"dls/dls-core/domain"
)

// TopSongs ranks the songs the hosts broke into. A moment that names no song
// is not a ranking row, and a song that has since been deleted from the
// library still counts, labelled by its id rather than vanishing.
func TopSongs(moments []domain.Moment, songs []domain.Song, limit int) []domain.Count {
	titles := make(map[domain.SongID]string, len(songs))
	for _, song := range songs {
		titles[song.ID] = song.Title
	}

	tally := map[string]int{}
	labels := map[string]string{}
	for _, moment := range moments {
		if moment.SongID == "" {
			continue
		}
		key := string(moment.SongID)
		tally[key]++
		labels[key] = fallback(titles[moment.SongID], key)
	}

	return rank(tally, labels, limit)
}

// TopTriggerWords ranks the words that set the hosts off. The same word typed
// with different case or stray spaces is one word: these are transcribed by
// hand, and the counter would otherwise split on the typing.
func TopTriggerWords(moments []domain.Moment, limit int) []domain.Count {
	tally := map[string]int{}
	labels := map[string]string{}
	for _, moment := range moments {
		word := strings.ToLower(strings.TrimSpace(moment.TriggerWord))
		if word == "" {
			continue
		}
		tally[word]++
		labels[word] = word
	}

	return rank(tally, labels, limit)
}

// ByActor ranks the people doing the counted thing, which is how Cauê and
// Load get compared. Types are named by slug because they are rows, not
// constants; an empty slice counts every type.
func ByActor(moments []domain.Moment, typeSlugs []string, people []domain.Person, limit int) []domain.Count {
	wanted := make(map[string]bool, len(typeSlugs))
	for _, slug := range typeSlugs {
		wanted[slug] = true
	}

	byID := peopleByID(people)
	tally := map[string]int{}
	labels := map[string]string{}
	for _, moment := range moments {
		if moment.ActorID == "" {
			continue
		}
		if len(wanted) > 0 && !wanted[moment.Type.Slug] {
			continue
		}
		person, ok := byID[moment.ActorID]
		key := fallback(person.Slug, string(moment.ActorID))
		tally[key]++
		if ok {
			labels[key] = fallback(person.Name, key)
		} else {
			labels[key] = key
		}
	}

	return rank(tally, labels, limit)
}

// rank turns the counting maps into the ranked slice. Ties break on label and
// then key so the same rows always come back in the same order: a chart that
// reshuffles between two identical loads reads as a bug.
func rank(tally map[string]int, labels map[string]string, limit int) []domain.Count {
	out := make([]domain.Count, 0, len(tally))
	for key, count := range tally {
		out = append(out, domain.Count{Key: key, Label: labels[key], Count: count})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].Label != out[j].Label {
			return out[i].Label < out[j].Label
		}
		return out[i].Key < out[j].Key
	})

	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

func fallback(value, or string) string {
	if value == "" {
		return or
	}
	return value
}

// TopGuests ranks the people who came on the show. It reads "interviewed" the
// same way the women counter does, so the two stats can never disagree about
// what an interview is.
func TopGuests(appearances []domain.Appearance, people []domain.Person, limit int) []domain.Count {
	byID := peopleByID(people)
	tally := map[string]int{}
	labels := map[string]string{}

	for _, appearance := range appearances {
		if !appearance.IsInterview || appearance.Role != domain.RoleGuest {
			continue
		}
		person, ok := byID[appearance.PersonID]
		key := fallback(person.Slug, string(appearance.PersonID))
		tally[key]++
		if ok {
			labels[key] = fallback(person.Name, key)
		} else {
			labels[key] = key
		}
	}

	return rank(tally, labels, limit)
}

// OrderTypes puts the moment types in the order the site shows them: the
// position given in the admin, then the label. Sorting here rather than in
// SQL means every page and every chart legend agrees without each one
// remembering to ask for the same order.
func OrderTypes(types []domain.MomentType) []domain.MomentType {
	out := make([]domain.MomentType, len(types))
	copy(out, types)

	sort.Slice(out, func(i, j int) bool {
		if out[i].Position != out[j].Position {
			return out[i].Position < out[j].Position
		}
		return out[i].Label < out[j].Label
	})
	return out
}
