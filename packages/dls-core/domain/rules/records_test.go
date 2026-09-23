package rules_test

import (
	"testing"
	"time"

	"dls/dls-core/domain"
	"dls/dls-core/domain/rules"
)

func day(n int) time.Time {
	return time.Date(2026, 1, n, 20, 0, 0, 0, time.UTC)
}

// Episodes carrying tallies, which is what the service hands these folds.
func episode(id domain.EpisodeID, title string, n int, counts map[string]int) domain.Episode {
	return domain.Episode{
		ID:          id,
		Slug:        string(id),
		Title:       title,
		PublishedAt: day(n),
		Tally:       domain.Tally{Moments: counts},
	}
}

func TestRecordsFindTheBestEpisodeForEachType(t *testing.T) {
	episodes := []domain.Episode{
		episode("e1", "Um", 1, map[string]int{"sing": 2, "question": 1}),
		episode("e2", "Dois", 2, map[string]int{"sing": 5}),
		episode("e3", "Três", 3, map[string]int{"question": 4}),
	}
	types := []domain.MomentType{{Slug: "sing"}, {Slug: "question"}}

	got := rules.Records(episodes, types)

	if len(got) != 2 {
		t.Fatalf("got %d records, want 2", len(got))
	}
	byType := map[string]domain.Record{}
	for _, record := range got {
		byType[record.TypeSlug] = record
	}
	if byType["sing"].Count != 5 || byType["sing"].Episode != "e2" {
		t.Errorf("sing record = %+v, want 5 in e2", byType["sing"])
	}
	if byType["question"].Count != 4 || byType["question"].Title != "Três" {
		t.Errorf("question record = %+v, want 4 in Três", byType["question"])
	}
}

func TestRecordsGiveATieToTheEarlierEpisode(t *testing.T) {
	episodes := []domain.Episode{
		// Handed over in the wrong order on purpose.
		episode("e2", "Dois", 2, map[string]int{"sing": 3}),
		episode("e1", "Um", 1, map[string]int{"sing": 3}),
	}

	got := rules.Records(episodes, []domain.MomentType{{Slug: "sing"}})

	if got[0].Episode != "e1" {
		t.Errorf("record = %+v, want the earlier episode to keep it", got[0])
	}
}

func TestRecordsSkipATypeThatNeverHappened(t *testing.T) {
	episodes := []domain.Episode{episode("e1", "Um", 1, map[string]int{"sing": 1})}
	types := []domain.MomentType{{Slug: "sing"}, {Slug: "never"}}

	got := rules.Records(episodes, types)

	if len(got) != 1 || got[0].TypeSlug != "sing" {
		t.Errorf("got %+v, want only the type that happened", got)
	}
}

func TestStreaksCountConsecutiveEpisodes(t *testing.T) {
	episodes := []domain.Episode{
		episode("e1", "Um", 1, map[string]int{"sing": 1}),
		episode("e2", "Dois", 2, map[string]int{"sing": 2}),
		episode("e3", "Três", 3, map[string]int{}),
		episode("e4", "Quatro", 4, map[string]int{"sing": 1}),
		episode("e5", "Cinco", 5, map[string]int{"sing": 1}),
		episode("e6", "Seis", 6, map[string]int{"sing": 9}),
	}

	got := rules.Streaks(episodes, []domain.MomentType{{Slug: "sing"}})

	if len(got) != 1 {
		t.Fatalf("got %d streaks, want 1", len(got))
	}
	if got[0].Length != 3 || got[0].From != "e4" || got[0].To != "e6" {
		t.Errorf("streak = %+v, want 3 episodes from e4 to e6", got[0])
	}
}

func TestStreaksReadInPublishOrderNotInputOrder(t *testing.T) {
	episodes := []domain.Episode{
		episode("e3", "Três", 3, map[string]int{"sing": 1}),
		episode("e1", "Um", 1, map[string]int{"sing": 1}),
		episode("e2", "Dois", 2, map[string]int{"sing": 1}),
	}

	got := rules.Streaks(episodes, []domain.MomentType{{Slug: "sing"}})

	if got[0].Length != 3 || got[0].From != "e1" || got[0].To != "e3" {
		t.Errorf("streak = %+v, want 3 episodes from e1 to e3", got[0])
	}
}

func TestStreaksSkipATypeThatNeverHappened(t *testing.T) {
	episodes := []domain.Episode{episode("e1", "Um", 1, map[string]int{})}

	got := rules.Streaks(episodes, []domain.MomentType{{Slug: "sing"}})

	if len(got) != 0 {
		t.Errorf("got %+v, want nothing", got)
	}
}

func TestAveragesDivideByEveryEpisodeInRange(t *testing.T) {
	episodes := []domain.Episode{
		episode("e1", "Um", 1, map[string]int{"sing": 3}),
		episode("e2", "Dois", 2, map[string]int{}),
		episode("e3", "Três", 3, map[string]int{"sing": 2}),
		episode("e4", "Quatro", 4, map[string]int{"sing": 5}),
	}

	got := rules.Averages(episodes, []domain.MomentType{{Slug: "sing"}})

	if len(got) != 1 {
		t.Fatalf("got %d averages, want 1", len(got))
	}
	// Ten over four episodes: the silent episode counts in the denominator.
	if got[0].Total != 10 || got[0].Episodes != 4 || got[0].PerEpisode != 2.5 {
		t.Errorf("average = %+v, want 10 over 4 at 2.5", got[0])
	}
}

func TestAveragesOfNoEpisodesIsZeroNotADivideByZero(t *testing.T) {
	got := rules.Averages(nil, []domain.MomentType{{Slug: "sing"}})

	if len(got) != 1 || got[0].PerEpisode != 0 {
		t.Errorf("got %+v, want a zero average", got)
	}
}
