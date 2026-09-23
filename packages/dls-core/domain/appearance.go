package domain

// AppearanceRole is how a person took part in one episode. It is recorded per
// appearance rather than read from Person.Kind, because a host can show up as
// a guest and a guest can end up behind the desk.
type AppearanceRole string

const (
	RoleHost   AppearanceRole = "host"
	RoleGuest  AppearanceRole = "guest"
	RoleRemote AppearanceRole = "remote"
)

// Appearance links a person to an episode. IsInterview marks the ones that
// were actually interviewed, so a guest who only walked through the shot does
// not inflate the count.
type Appearance struct {
	ID          AppearanceID
	EpisodeID   EpisodeID
	PersonID    PersonID
	Role        AppearanceRole
	IsInterview bool
}
