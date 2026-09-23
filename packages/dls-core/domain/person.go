package domain

// PersonKind separates the people who make the show from the people who visit
// it. A guest who later joins the cast changes kind; their past appearances
// keep whatever role they were recorded with.
type PersonKind string

const (
	PersonHost  PersonKind = "host"
	PersonGuest PersonKind = "guest"
	PersonStaff PersonKind = "staff"
)

// Gender is what the "women interviewed" counter reads. Unknown is the zero
// value on purpose: an unfilled field must never be counted as a woman.
type Gender string

const (
	GenderUnknown   Gender = "unknown"
	GenderWoman     Gender = "woman"
	GenderMan       Gender = "man"
	GenderNonBinary Gender = "nonbinary"
)

type Person struct {
	ID     PersonID
	Slug   string
	Name   string
	Kind   PersonKind
	Gender Gender
	Photo  string
}
