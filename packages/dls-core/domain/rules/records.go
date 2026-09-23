package rules

import (
	"sort"

	"dls/dls-core/domain"
)

// Records is the box score: the single best episode for each moment type.
// Every fold here reads the tallies already attached to the episodes, so the
// rows are counted once and reused.
func Records(episodes []domain.Episode, types []domain.MomentType) []domain.Record {
	ordered := byPublishDate(episodes)
	out := make([]domain.Record, 0, len(types))

	for _, momentType := range types {
		best := domain.Record{TypeSlug: momentType.Slug}
		for _, episode := range ordered {
			// Strictly greater, walking in publish order, leaves a tie with
			// the episode that did it first.
			if count := episode.Tally.Of(momentType.Slug); count > best.Count {
				best.Count, best.Episode, best.Title = count, episode.ID, episode.Title
			}
		}
		if best.Count > 0 {
			out = append(out, best)
		}
	}
	return out
}

// Streaks is the longest run of consecutive episodes in which a type happened
// at least once. Consecutive means consecutive in publish order: a viewer
// counts shows, not episode numbers.
func Streaks(episodes []domain.Episode, types []domain.MomentType) []domain.Streak {
	ordered := byPublishDate(episodes)
	out := make([]domain.Streak, 0, len(types))

	for _, momentType := range types {
		best, current := domain.Streak{TypeSlug: momentType.Slug}, domain.Streak{TypeSlug: momentType.Slug}
		for _, episode := range ordered {
			if episode.Tally.Of(momentType.Slug) == 0 {
				current = domain.Streak{TypeSlug: momentType.Slug}
				continue
			}
			if current.Length == 0 {
				current.From = episode.ID
			}
			current.Length, current.To = current.Length+1, episode.ID
			if current.Length > best.Length {
				best = current
			}
		}
		if best.Length > 0 {
			out = append(out, best)
		}
	}
	return out
}

// Averages is the per-episode rate of each type. Every episode in range is in
// the denominator, silent ones included: an average that only counted the
// episodes where it happened would never drop below one.
func Averages(episodes []domain.Episode, types []domain.MomentType) []domain.Average {
	out := make([]domain.Average, 0, len(types))

	for _, momentType := range types {
		average := domain.Average{TypeSlug: momentType.Slug, Episodes: len(episodes)}
		for _, episode := range episodes {
			average.Total += episode.Tally.Of(momentType.Slug)
		}
		if average.Episodes > 0 {
			average.PerEpisode = float64(average.Total) / float64(average.Episodes)
		}
		out = append(out, average)
	}
	return out
}

func byPublishDate(episodes []domain.Episode) []domain.Episode {
	out := make([]domain.Episode, len(episodes))
	copy(out, episodes)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].PublishedAt.Before(out[j].PublishedAt)
	})
	return out
}
