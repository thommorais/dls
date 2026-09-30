package services_test

import (
	"errors"
	"testing"
	"time"

	"dls/dls-core/domain"
	"dls/dls-core/services"
)

var (
	jan = time.Date(2026, 1, 5, 20, 0, 0, 0, time.UTC)
	feb = time.Date(2026, 2, 10, 20, 0, 0, 0, time.UTC)
)

// The three types the fixture knows about. Dance exists and has never
// happened, which is the case a hardcoded counter could not express.
var (
	question = domain.MomentType{ID: "t-question", Slug: "ana_question", Label: "Pergunta pra Ana", Color: "#6366F1", Position: 1}
	sing     = domain.MomentType{ID: "t-sing", Slug: "sing", Label: "Saiu cantando", Color: "#10B981", Position: 2}
	dance    = domain.MomentType{ID: "t-dance", Slug: "dance", Label: "Saiu dançando", Color: "#F59E0B", Position: 3}
)

func momentOf(episode domain.EpisodeID, momentType domain.MomentType, actor domain.PersonID) domain.Moment {
	return domain.Moment{EpisodeID: episode, TypeID: momentType.ID, Type: momentType, ActorID: actor}
}

func singing(episode domain.EpisodeID, actor domain.PersonID, word string, song domain.SongID) domain.Moment {
	m := momentOf(episode, sing, actor)
	m.TriggerWord, m.SongID = word, song
	return m
}

type fixture struct {
	svc         *services.StatsService
	archive     *fakeArchive
	types       *fakeMomentTypes
	episodes    *fakeEpisodes
	moments     *fakeMoments
	appearances *fakeAppearances
	openings    *fakeOpenings
	people      *fakePeople
	songs       *fakeSongs
}

// Two episodes of show: three questions to Ana, three songs broken into, two
// women interviewed, two openings aired.
func newFixture() *fixture {
	f := &fixture{
		episodes: &fakeEpisodes{items: []domain.Episode{
			{ID: "e1", Slug: "ep-1", Number: 1, Title: "Episódio 1", PublishedAt: jan},
			{ID: "e2", Slug: "ep-2", Number: 2, Title: "Episódio 2", PublishedAt: feb},
		}},
		types: &fakeMomentTypes{items: []domain.MomentType{question, sing, dance}},
		moments: &fakeMoments{items: []domain.Moment{
			momentOf("e1", question, "p-caue"),
			momentOf("e1", question, "p-load"),
			singing("e1", "p-load", "pizza", "s-evid"),
			momentOf("e2", question, "p-caue"),
			singing("e2", "p-caue", "bolo", "s-evid"),
			singing("e2", "p-load", "café", "s-anun"),
		}},
		appearances: &fakeAppearances{items: []domain.Appearance{
			{ID: "a1", EpisodeID: "e1", PersonID: "p-guest", Role: domain.RoleGuest, IsInterview: true},
			{ID: "a2", EpisodeID: "e1", PersonID: "p-dude", Role: domain.RoleGuest, IsInterview: true},
			{ID: "a3", EpisodeID: "e2", PersonID: "p-guest", Role: domain.RoleGuest, IsInterview: true},
		}},
		openings: &fakeOpenings{items: []domain.Opening{
			{ID: "o1", EpisodeID: "e1", Genre: "forró"},
			{ID: "o2", EpisodeID: "e2", Genre: "rock"},
		}},
		people: &fakePeople{items: []domain.Person{
			{ID: "p-caue", Slug: "caue", Name: "Cauê", Kind: domain.PersonHost, Gender: domain.GenderMan},
			{ID: "p-load", Slug: "load", Name: "Load", Kind: domain.PersonHost, Gender: domain.GenderMan},
			{ID: "p-guest", Slug: "guest", Name: "Guest", Kind: domain.PersonGuest, Gender: domain.GenderWoman},
			{ID: "p-dude", Slug: "dude", Name: "Dude", Kind: domain.PersonGuest, Gender: domain.GenderMan},
		}},
		archive: &fakeArchive{items: []domain.ArchiveEntry{
			{ID: "k1", Slug: "saiu-cantando", Title: "Saiu cantando", Kind: domain.ArchiveGlossary},
			{ID: "k2", Slug: "o-primeiro-episodio", Title: "O primeiro episódio", Kind: domain.ArchiveMilestone},
		}},
		songs: &fakeSongs{items: []domain.Song{
			{ID: "s-evid", Title: "Evidências"},
			{ID: "s-anun", Title: "Anunciação"},
		}},
	}
	f.svc = services.NewStatsService(f.episodes, f.types, f.moments, f.appearances, f.openings, f.people, f.songs, f.archive)
	return f
}

func TestOverviewCountsEveryEpisodeWhenUnfiltered(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Overview(t.Context(), domain.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if got.Tally.Of("ana_question") != 3 || got.Tally.Of("sing") != 3 {
		t.Errorf("moments = %v, want 3 questions and 3 songs", got.Tally.Moments)
	}
	if got.Tally.WomenInterviewed != 2 || got.Tally.OpeningsAired != 2 {
		t.Errorf("women = %d, openings = %d, want 2 and 2", got.Tally.WomenInterviewed, got.Tally.OpeningsAired)
	}
	if got.Episodes != 2 {
		t.Errorf("episodes = %d, want 2", got.Episodes)
	}
	if got.Latest == nil || got.Latest.Slug != "ep-2" {
		t.Errorf("latest = %+v, want ep-2", got.Latest)
	}
	if got.First == nil || got.First.Slug != "ep-1" {
		t.Errorf("first = %+v, want ep-1", got.First)
	}
}

func TestOverviewReadsEachTableWholeWhenUnfiltered(t *testing.T) {
	f := newFixture()

	if _, err := f.svc.Overview(t.Context(), domain.Filter{}); err != nil {
		t.Fatal(err)
	}

	if !f.moments.scoped {
		t.Fatal("moments were never read")
	}
	if f.moments.lastScope != nil {
		t.Errorf("scope = %v, want nil: an unfiltered read should not name every episode", f.moments.lastScope)
	}
}

func TestOverviewScopesToTheFilteredEpisodes(t *testing.T) {
	f := newFixture()
	from := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	got, err := f.svc.Overview(t.Context(), domain.Filter{From: &from})
	if err != nil {
		t.Fatal(err)
	}

	if got.Tally.Of("ana_question") != 1 || got.Tally.Of("sing") != 2 {
		t.Errorf("moments = %v, want 1 question and 2 songs", got.Tally.Moments)
	}
	if got.Tally.WomenInterviewed != 1 || got.Tally.OpeningsAired != 1 {
		t.Errorf("women = %d, openings = %d, want 1 and 1", got.Tally.WomenInterviewed, got.Tally.OpeningsAired)
	}
	if len(f.moments.lastScope) != 1 || f.moments.lastScope[0] != "e2" {
		t.Errorf("scope = %v, want [e2]", f.moments.lastScope)
	}
}

func TestOverviewOfNoEpisodesIsZero(t *testing.T) {
	f := newFixture()
	from := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	got, err := f.svc.Overview(t.Context(), domain.Filter{From: &from})
	if err != nil {
		t.Fatal(err)
	}

	if got.Tally.TotalMoments() != 0 || got.Tally.WomenInterviewed != 0 || got.Episodes != 0 || got.Latest != nil {
		t.Errorf("got %+v, want an empty overview", got)
	}
	if len(got.Types) != 3 {
		t.Errorf("types = %v, want all three even with nothing counted", got.Types)
	}
}

func TestOverviewPropagatesAStorageFailure(t *testing.T) {
	f := newFixture()
	f.moments.failOn = "ListByEpisodes"

	if _, err := f.svc.Overview(t.Context(), domain.Filter{}); err == nil {
		t.Fatal("want an error when the moments table cannot be read")
	}
}

func TestEpisodesCarryTheirOwnTally(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Episodes(t.Context(), domain.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d episodes, want 2", len(got))
	}
	if got[0].Tally.Of("ana_question") != 2 || got[0].Tally.Of("sing") != 1 || got[0].Tally.WomenInterviewed != 1 {
		t.Errorf("ep-1 tally = %+v, want 2 questions, 1 song, 1 woman", got[0].Tally)
	}
	if got[1].Tally.Of("ana_question") != 1 || got[1].Tally.Of("sing") != 2 || got[1].Tally.OpeningsAired != 1 {
		t.Errorf("ep-2 tally = %+v, want 1 question, 2 songs, 1 opening", got[1].Tally)
	}
}

func TestEpisodeGathersOnlyThatEpisode(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Episode(t.Context(), "ep-1")
	if err != nil {
		t.Fatal(err)
	}

	if got.Episode.Slug != "ep-1" {
		t.Fatalf("episode = %s, want ep-1", got.Episode.Slug)
	}
	if len(got.Moments) != 3 {
		t.Errorf("got %d moments, want 3", len(got.Moments))
	}
	if len(got.Appearances) != 2 {
		t.Errorf("got %d appearances, want 2", len(got.Appearances))
	}
	if len(got.Openings) != 1 {
		t.Errorf("got %d openings, want 1", len(got.Openings))
	}
	if got.Episode.Tally.Of("ana_question") != 2 || got.Episode.Tally.Of("sing") != 1 {
		t.Errorf("tally = %+v, want 2 questions and 1 song", got.Episode.Tally)
	}
}

func TestEpisodeReportsAMissingSlugAsNotFound(t *testing.T) {
	f := newFixture()

	_, err := f.svc.Episode(t.Context(), "nope")

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestRankingsOrderEveryBoard(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Rankings(t.Context(), domain.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Songs) != 2 || got.Songs[0].Label != "Evidências" || got.Songs[0].Count != 2 {
		t.Errorf("songs = %+v, want Evidências on top with 2", got.Songs)
	}
	if len(got.TriggerWords) != 3 {
		t.Errorf("trigger words = %+v, want 3", got.TriggerWords)
	}
	if got.Actors["sing"][0].Key != "load" || got.Actors["sing"][0].Count != 2 {
		t.Errorf("singers = %+v, want Load on top with 2", got.Actors["sing"])
	}
	if got.Actors["ana_question"][0].Key != "caue" || got.Actors["ana_question"][0].Count != 2 {
		t.Errorf("ana askers = %+v, want Cauê on top with 2", got.Actors["ana_question"])
	}
	board, ok := got.Actors["dance"]
	if !ok || len(board) != 0 {
		t.Errorf("dance board = %v (present: %v), want present and empty", board, ok)
	}
	if got.Guests[0].Key != "guest" || got.Guests[0].Count != 2 {
		t.Errorf("guests = %+v, want the returning guest on top", got.Guests)
	}
}

func TestRankingsCapEveryBoardAtTheLimit(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Rankings(t.Context(), domain.Filter{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}

	for name, board := range map[string][]domain.Count{
		"songs": got.Songs, "words": got.TriggerWords, "guests": got.Guests,
		"singers": got.Actors["sing"], "ana askers": got.Actors["ana_question"],
	} {
		if len(board) != 1 {
			t.Errorf("%s has %d rows, want 1", name, len(board))
		}
	}
}

func TestOpeningsPassTheGenreToStorage(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Openings(t.Context(), domain.Filter{Genre: "rock"})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Errorf("got %d openings, want the 1 rock one", len(got))
	}
	if f.openings.lastFilter.Genre != "rock" {
		t.Errorf("filter = %+v, want the genre passed through", f.openings.lastFilter)
	}
}

func TestOverviewCarriesEveryTypeInOrder(t *testing.T) {
	f := newFixture()
	// However the storage hands them over.
	f.types.items = []domain.MomentType{dance, sing, question}

	got, err := f.svc.Overview(t.Context(), domain.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"ana_question", "sing", "dance"}
	for i, slug := range want {
		if got.Types[i].Slug != slug {
			t.Fatalf("types = %+v, want %v", got.Types, want)
		}
	}
	if got.Tally.Of("dance") != 0 {
		t.Errorf("dance = %d, want 0: a type nobody has done yet still has a counter", got.Tally.Of("dance"))
	}
}

func TestOverviewPropagatesAFailureReadingTheTypes(t *testing.T) {
	f := newFixture()
	f.types.failOn = "List"

	if _, err := f.svc.Overview(t.Context(), domain.Filter{}); err == nil {
		t.Fatal("want an error when the types table cannot be read")
	}
}

func TestRankingsCarryTheBoxScore(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Rankings(t.Context(), domain.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	records := map[string]domain.Record{}
	for _, record := range got.Records {
		records[record.TypeSlug] = record
	}
	// ep-1 has two questions, ep-2 has two songs.
	if records["ana_question"].Count != 2 || records["ana_question"].Episode != "e1" {
		t.Errorf("question record = %+v, want 2 in e1", records["ana_question"])
	}
	if records["sing"].Count != 2 || records["sing"].Episode != "e2" {
		t.Errorf("sing record = %+v, want 2 in e2", records["sing"])
	}

	streaks := map[string]domain.Streak{}
	for _, streak := range got.Streaks {
		streaks[streak.TypeSlug] = streak
	}
	if streaks["sing"].Length != 2 {
		t.Errorf("sing streak = %+v, want both episodes", streaks["sing"])
	}

	averages := map[string]domain.Average{}
	for _, average := range got.Averages {
		averages[average.TypeSlug] = average
	}
	if averages["ana_question"].PerEpisode != 1.5 {
		t.Errorf("question average = %+v, want 1.5 per episode", averages["ana_question"])
	}
	if averages["dance"].PerEpisode != 0 {
		t.Errorf("dance average = %+v, want 0", averages["dance"])
	}
}

func TestArchivePassesTheKindToStorage(t *testing.T) {
	f := newFixture()

	got, err := f.svc.Archive(t.Context(), domain.Filter{Kind: "glossary"})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].Slug != "saiu-cantando" {
		t.Errorf("got %+v, want the one glossary entry", got)
	}
	if f.archive.lastFilter.Kind != "glossary" {
		t.Errorf("filter = %+v, want the kind passed through", f.archive.lastFilter)
	}
}

func TestArchiveEntryReportsAMissingSlugAsNotFound(t *testing.T) {
	f := newFixture()

	if _, err := f.svc.ArchiveEntry(t.Context(), "nope"); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
