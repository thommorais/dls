package rules_test

import (
	"testing"

	"dls/dls-core/domain"
	"dls/dls-core/domain/rules"
)

var songs = []domain.Song{
	{ID: "s-evidencias", Title: "Evidências", Artist: "Chitãozinho & Xororó"},
	{ID: "s-anunciacao", Title: "Anunciação", Artist: "Alceu Valença"},
	{ID: "s-zeca", Title: "Deixa a Vida Me Levar", Artist: "Zeca Pagodinho"},
}

func trigger(word string, song domain.SongID, actor domain.PersonID) domain.Moment {
	m := moment("e1", sing, actor)
	m.TriggerWord, m.SongID = word, song
	return m
}

func keys(counts []domain.Count) []string {
	out := make([]string, 0, len(counts))
	for _, count := range counts {
		out = append(out, count.Key)
	}
	return out
}

func equal(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestTopSongsRanksByCountAndBreaksTiesByLabel(t *testing.T) {
	moments := []domain.Moment{
		trigger("evidência", "s-evidencias", caue.ID),
		trigger("prova", "s-evidencias", load.ID),
		trigger("vida", "s-zeca", load.ID),
		trigger("anúncio", "s-anunciacao", caue.ID),
	}

	got := rules.TopSongs(moments, songs, 0)

	// Evidências leads on count; the two one-offs tie and sort by title.
	if want := []string{"s-evidencias", "s-anunciacao", "s-zeca"}; !equal(keys(got), want) {
		t.Fatalf("keys = %v, want %v", keys(got), want)
	}
	if got[0].Count != 2 || got[0].Label != "Evidências" {
		t.Errorf("top = %+v, want Evidências counted twice", got[0])
	}
}

func TestTopSongsIgnoresMomentsThatNameNoSong(t *testing.T) {
	moments := []domain.Moment{
		trigger("evidência", "s-evidencias", caue.ID),
		trigger("nada", "", caue.ID),
		moment("e1", question, caue.ID),
	}

	got := rules.TopSongs(moments, songs, 0)

	if len(got) != 1 || got[0].Count != 1 {
		t.Errorf("got %+v, want only Evidências once", got)
	}
}

func TestTopSongsFallsBackToTheIDWhenTheSongIsGone(t *testing.T) {
	moments := []domain.Moment{trigger("some", "s-deleted", caue.ID)}

	got := rules.TopSongs(moments, songs, 0)

	if len(got) != 1 || got[0].Label != "s-deleted" {
		t.Errorf("got %+v, want the id as the label", got)
	}
}

func TestTopTriggerWordsFoldsCaseAndSurroundingSpace(t *testing.T) {
	moments := []domain.Moment{
		trigger("Pizza", "s-zeca", caue.ID),
		trigger(" pizza ", "s-zeca", load.ID),
		trigger("PIZZA", "s-zeca", load.ID),
		trigger("café", "s-zeca", caue.ID),
	}

	got := rules.TopTriggerWords(moments, 0)

	if want := []string{"pizza", "café"}; !equal(keys(got), want) {
		t.Fatalf("keys = %v, want %v", keys(got), want)
	}
	if got[0].Count != 3 {
		t.Errorf("pizza counted %d times, want 3", got[0].Count)
	}
}

func TestTopTriggerWordsIgnoresEmptyWords(t *testing.T) {
	moments := []domain.Moment{
		trigger("", "s-zeca", caue.ID),
		trigger("   ", "s-zeca", caue.ID),
		trigger("pizza", "s-zeca", caue.ID),
	}

	got := rules.TopTriggerWords(moments, 0)

	if len(got) != 1 {
		t.Errorf("got %+v, want only pizza", got)
	}
}

func TestTopTriggerWordsRespectsTheLimit(t *testing.T) {
	moments := []domain.Moment{
		trigger("pizza", "s-zeca", caue.ID),
		trigger("pizza", "s-zeca", caue.ID),
		trigger("café", "s-zeca", caue.ID),
		trigger("bolo", "s-zeca", caue.ID),
	}

	got := rules.TopTriggerWords(moments, 2)

	if want := []string{"pizza", "bolo"}; !equal(keys(got), want) {
		t.Errorf("keys = %v, want %v", keys(got), want)
	}
}

func TestByActorCountsOnlyTheGivenTypes(t *testing.T) {
	moments := []domain.Moment{
		trigger("pizza", "s-zeca", caue.ID),
		trigger("café", "s-zeca", load.ID),
		trigger("bolo", "s-zeca", load.ID),
		moment("e1", question, caue.ID),
		moment("e1", question, caue.ID),
	}

	singers := rules.ByActor(moments, []string{sing.Slug}, cast(), 0)

	if want := []string{"load", "caue"}; !equal(keys(singers), want) {
		t.Fatalf("keys = %v, want %v", keys(singers), want)
	}
	if singers[0].Count != 2 || singers[0].Label != "Load" {
		t.Errorf("top singer = %+v, want Load with 2", singers[0])
	}
}

func TestByActorWithNoTypesCountsEverything(t *testing.T) {
	moments := []domain.Moment{
		trigger("pizza", "s-zeca", caue.ID),
		moment("e1", question, caue.ID),
		moment("e1", sing, load.ID),
	}

	got := rules.ByActor(moments, nil, cast(), 0)

	if want := []string{"caue", "load"}; !equal(keys(got), want) {
		t.Errorf("keys = %v, want %v", keys(got), want)
	}
}

func TestByActorIgnoresMomentsWithNoActor(t *testing.T) {
	moments := []domain.Moment{
		moment("e1", question, ""),
		moment("e1", question, caue.ID),
	}

	got := rules.ByActor(moments, nil, cast(), 0)

	if len(got) != 1 {
		t.Errorf("got %+v, want only Cauê", got)
	}
}

func TestByActorNamesAnUnknownTypeNoOne(t *testing.T) {
	moments := []domain.Moment{
		trigger("pizza", "s-zeca", caue.ID),
		moment("e1", question, load.ID),
	}

	got := rules.ByActor(moments, []string{"not_a_type"}, cast(), 0)

	if len(got) != 0 {
		t.Errorf("got %+v, want nothing", got)
	}
}

func TestTopGuestsRanksInterviewedGuests(t *testing.T) {
	appearances := []domain.Appearance{
		{EpisodeID: "e1", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
		{EpisodeID: "e2", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
		{EpisodeID: "e2", PersonID: dude.ID, Role: domain.RoleGuest, IsInterview: true},
	}

	got := rules.TopGuests(appearances, cast(), 0)

	if want := []string{"guest", "dude"}; !equal(keys(got), want) {
		t.Fatalf("keys = %v, want %v", keys(got), want)
	}
	if got[0].Count != 2 {
		t.Errorf("top guest counted %d times, want 2", got[0].Count)
	}
}

func TestTopGuestsIgnoresHostsAndWalkOns(t *testing.T) {
	appearances := []domain.Appearance{
		{EpisodeID: "e1", PersonID: caue.ID, Role: domain.RoleHost, IsInterview: true},
		// On the episode, never interviewed.
		{EpisodeID: "e1", PersonID: dude.ID, Role: domain.RoleGuest, IsInterview: false},
		{EpisodeID: "e1", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
	}

	got := rules.TopGuests(appearances, cast(), 0)

	if len(got) != 1 || got[0].Key != "guest" {
		t.Errorf("got %+v, want only the interviewed guest", got)
	}
}

func TestOrderTypesSortsByPositionThenLabel(t *testing.T) {
	types := []domain.MomentType{
		{Slug: "c", Label: "Zebra", Position: 1},
		{Slug: "a", Label: "Abacaxi", Position: 2},
		{Slug: "b", Label: "Banana", Position: 1},
		// Never positioned, so it lands first and sorts by label among its kind.
		{Slug: "d", Label: "Dado"},
	}

	got := rules.OrderTypes(types)

	want := []string{"d", "b", "c", "a"}
	slugs := make([]string, 0, len(got))
	for _, momentType := range got {
		slugs = append(slugs, momentType.Slug)
	}
	if !equal(slugs, want) {
		t.Errorf("slugs = %v, want %v", slugs, want)
	}
}
