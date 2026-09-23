package rules_test

import (
	"testing"

	"dls/dls-core/domain"
	"dls/dls-core/domain/rules"
)

// The cast used across the tally tests: two hosts, a woman guest, a man guest,
// and one person whose gender was never filled in.
var (
	caue    = domain.Person{ID: "p-caue", Slug: "caue", Name: "Cauê Moura", Kind: domain.PersonHost, Gender: domain.GenderMan}
	load    = domain.Person{ID: "p-load", Slug: "load", Name: "Load", Kind: domain.PersonHost, Gender: domain.GenderMan}
	ana     = domain.Person{ID: "p-ana", Slug: "ana-bonasa", Name: "Ana Bonasa", Kind: domain.PersonStaff, Gender: domain.GenderWoman}
	guest   = domain.Person{ID: "p-guest", Slug: "guest", Name: "Guest", Kind: domain.PersonGuest, Gender: domain.GenderWoman}
	dude    = domain.Person{ID: "p-dude", Slug: "dude", Name: "Dude", Kind: domain.PersonGuest, Gender: domain.GenderMan}
	nameles = domain.Person{ID: "p-nameless", Slug: "nameless", Name: "Nameless", Kind: domain.PersonGuest, Gender: domain.GenderUnknown}
)

// Two moment types, exactly as they would be read from the types collection.
var (
	question = domain.MomentType{ID: "t-question", Slug: "ana_question", Label: "Pergunta pra Ana"}
	sing     = domain.MomentType{ID: "t-sing", Slug: "sing", Label: "Saiu cantando"}
)

func cast() []domain.Person {
	return []domain.Person{caue, load, ana, guest, dude, nameles}
}

func moment(episode domain.EpisodeID, momentType domain.MomentType, actor domain.PersonID) domain.Moment {
	return domain.Moment{EpisodeID: episode, TypeID: momentType.ID, Type: momentType, ActorID: actor}
}

func TestTallyCountsEveryTypeOfMoment(t *testing.T) {
	moments := []domain.Moment{
		moment("e1", question, caue.ID),
		moment("e1", question, load.ID),
		moment("e1", sing, load.ID),
	}
	appearances := []domain.Appearance{
		{ID: "a1", EpisodeID: "e1", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
	}
	openings := []domain.Opening{
		{ID: "o1", EpisodeID: "e1", Status: domain.OpeningAired},
		{ID: "o2", Status: domain.OpeningReceived},
	}

	got := rules.Tally(moments, appearances, cast(), openings)

	if got.Of("ana_question") != 2 {
		t.Errorf("ana_question = %d, want 2", got.Of("ana_question"))
	}
	if got.Of("sing") != 1 {
		t.Errorf("sing = %d, want 1", got.Of("sing"))
	}
	if got.WomenInterviewed != 1 || got.OpeningsAired != 1 {
		t.Errorf("women = %d, openings aired = %d, want 1 and 1", got.WomenInterviewed, got.OpeningsAired)
	}
	if got.TotalMoments() != 3 {
		t.Errorf("total moments = %d, want 3", got.TotalMoments())
	}
}

func TestTallyReadsAnUncountedTypeAsZero(t *testing.T) {
	got := rules.Tally(nil, nil, nil, nil)

	if got.Of("never_happened") != 0 {
		t.Errorf("got %d, want 0", got.Of("never_happened"))
	}
	if got.TotalMoments() != 0 {
		t.Errorf("total = %d, want 0", got.TotalMoments())
	}
}

func TestTallyIgnoresAMomentWhoseTypeIsGone(t *testing.T) {
	// The type row was deleted, so the adapter resolved nothing.
	orphan := domain.Moment{EpisodeID: "e1", TypeID: "t-deleted"}

	got := rules.Tally([]domain.Moment{orphan, moment("e1", sing, caue.ID)}, nil, nil, nil)

	if got.TotalMoments() != 1 {
		t.Errorf("total = %d, want 1: a moment with no type cannot be counted under one", got.TotalMoments())
	}
}

func TestWomenInterviewedCountsOnlyInterviewedWomenGuests(t *testing.T) {
	appearances := []domain.Appearance{
		{ID: "a1", EpisodeID: "e1", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
		// Walked through the shot, never interviewed.
		{ID: "a2", EpisodeID: "e1", PersonID: ana.ID, Role: domain.RoleGuest, IsInterview: false},
		// A woman behind the desk is not a guest being interviewed.
		{ID: "a3", EpisodeID: "e2", PersonID: ana.ID, Role: domain.RoleHost, IsInterview: true},
		{ID: "a4", EpisodeID: "e2", PersonID: dude.ID, Role: domain.RoleGuest, IsInterview: true},
		// An unfilled gender must never be counted as a woman.
		{ID: "a5", EpisodeID: "e2", PersonID: nameles.ID, Role: domain.RoleGuest, IsInterview: true},
	}

	got := rules.Tally(nil, appearances, cast(), nil)

	if got.WomenInterviewed != 1 {
		t.Errorf("women interviewed = %d, want 1", got.WomenInterviewed)
	}
}

func TestWomenInterviewedCountsTheSameGuestOncePerEpisode(t *testing.T) {
	appearances := []domain.Appearance{
		{ID: "a1", EpisodeID: "e1", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
		{ID: "a2", EpisodeID: "e2", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
	}

	got := rules.Tally(nil, appearances, cast(), nil)

	if got.WomenInterviewed != 2 {
		t.Errorf("women interviewed = %d, want 2: the stat counts interviews, not people", got.WomenInterviewed)
	}
}

func TestWomenInterviewedIgnoresAnAppearanceWithNoPerson(t *testing.T) {
	appearances := []domain.Appearance{
		{ID: "a1", EpisodeID: "e1", PersonID: "p-deleted", Role: domain.RoleGuest, IsInterview: true},
	}

	got := rules.Tally(nil, appearances, cast(), nil)

	if got.WomenInterviewed != 0 {
		t.Errorf("women interviewed = %d, want 0", got.WomenInterviewed)
	}
}

func TestTallyByEpisodeSplitsTheCountersPerEpisode(t *testing.T) {
	moments := []domain.Moment{
		moment("e1", question, caue.ID),
		moment("e2", question, caue.ID),
		moment("e2", sing, load.ID),
	}
	appearances := []domain.Appearance{
		{ID: "a1", EpisodeID: "e2", PersonID: guest.ID, Role: domain.RoleGuest, IsInterview: true},
	}
	openings := []domain.Opening{
		{ID: "o1", EpisodeID: "e1", Status: domain.OpeningAired},
	}

	got := rules.TallyByEpisode(moments, appearances, cast(), openings)

	if len(got) != 2 {
		t.Fatalf("got %d episodes, want 2", len(got))
	}
	if got["e1"].Of("ana_question") != 1 || got["e1"].OpeningsAired != 1 || got["e1"].Of("sing") != 0 {
		t.Errorf("e1 = %+v, want one question and one opening", got["e1"])
	}
	if got["e2"].Of("ana_question") != 1 || got["e2"].Of("sing") != 1 || got["e2"].WomenInterviewed != 1 {
		t.Errorf("e2 = %+v, want a question, a song and a woman interviewed", got["e2"])
	}
	if got["e1"].OpeningsAired != 1 || got["e2"].OpeningsAired != 0 {
		t.Errorf("openings landed on the wrong episode: %+v", got)
	}
}
