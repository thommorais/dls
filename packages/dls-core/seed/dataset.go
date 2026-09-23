package seed

import (
	"fmt"
	"math/rand/v2"
	"time"

	"dls/dls-core/domain"
)

// Relations are carried by slug rather than by id: the dataset is written
// before anything exists, and the writer resolves them as it inserts.

type Moment struct {
	TypeSlug       string
	VideoTimestamp int
	ActorSlug      string
	Summary        string
	TriggerWord    string
	SongKey        string
}

type Appearance struct {
	PersonSlug  string
	Role        domain.AppearanceRole
	IsInterview bool
}

type Opening struct {
	Title        string
	AuthorName   string
	AuthorHandle string
	Genre        string
	AtSeconds    int
	SentAt       time.Time
	Status       domain.OpeningStatus
}

type Episode struct {
	Number          int
	Slug            string
	Title           string
	YouTubeID       string
	ThumbnailURL    string
	Description     string
	PublishedAt     time.Time
	DurationSeconds int

	Moments     []Moment
	Appearances []Appearance
	Openings    []Opening
}

type ArchiveEntry struct {
	Slug           string
	Title          string
	Summary        string
	Body           string
	Kind           domain.ArchiveKind
	Tags           []string
	EpisodeSlug    string
	PersonSlug     string
	MomentTypeSlug string
	SourceURL      string
	PublishedAt    time.Time
}

type Dataset struct {
	Types    []domain.MomentType
	People   []domain.Person
	Songs    []Song
	Episodes []Episode
	Archive  []ArchiveEntry
}

// firstEpisode is a Wednesday, and the show goes out weekly.
var firstEpisode = time.Date(2025, 1, 8, 20, 0, 0, 0, time.UTC)

var questions = []string{
	"Dona Neide, tem café",
	"Dona Neide, que horas acaba isso",
	"Dona Neide, o senhor lá fora é convidado",
	"Dona Neide, cadê o cabo",
	"Dona Neide, a gente já falou disso",
	"Dona Neide, isso pode ir pro corte",
	"Dona Neide, quem marcou essa entrevista",
	"Dona Neide, tá gravando",
}

var theories = []string{
	"Elevador só demora quando alguém está atrasado",
	"Segunda-feira começa no domingo às seis da tarde",
	"Toda fila de padaria tem uma pessoa que não sabe o que quer",
	"Bolo de aniversário no trabalho é uma forma de reunião",
	"Ninguém nunca viu dois pombos brigando de verdade",
	"Todo bairro tem exatamente uma loja que nunca tem cliente",
}

var fights = []string{
	"Pastel é salgado ou lanche",
	"Se pode colocar ketchup em churrasco",
	"Qual o lado certo de descer do ônibus",
	"Se bolo de rolo é bolo",
	"Quem paga a conta quando só um pediu entrada",
}

var brands = []string{
	"a padaria da esquina",
	"uma marca de chinelo",
	"o sabonete que a avó usava",
	"uma lanchonete de estrada",
	"a cerveja mais barata do mercado",
}

var cries = []string{
	"Lembrou do cachorro de infância",
	"Ouviu a própria mãe em um áudio",
	"A convidada falou do pai dela",
	"Propaganda de fim de ano no intervalo",
}

// typeFrequency is the per-episode range for each type, which is what makes
// the charts uneven enough to be worth looking at.
var typeFrequency = map[string][2]int{
	"pergunta-pra-neide":  {2, 7},
	"saiu-cantando":       {0, 5},
	"teoria-maluca":       {1, 4},
	"briga-de-bar":        {0, 3},
	"chorou-ao-vivo":      {0, 1},
	"propaganda-nao-paga": {0, 2},
}

// Hand-placed episodes, so the records and streaks have a story behind them
// instead of being whatever the generator happened to roll.
var scripted = map[int]map[string]int{
	// "O recorde: dezoito músicas em duas horas"
	16: {"saiu-cantando": 18},
	// "O dia em que a Belinha chorou com propaganda"
	20: {"chorou-ao-vivo": 2, "propaganda-nao-paga": 4},
	// "O episódio em que ninguém concordou com nada"
	22: {"briga-de-bar": 7, "saiu-cantando": 0},
	// "O dia em que todo mundo cantou junto"
	31: {"saiu-cantando": 9},
	// "O episódio sem nenhum assunto"
	34: {"teoria-maluca": 6},
}

// Quiet weeks, which is what puts a gap in a streak.
var silentSongs = map[int]bool{3: true, 9: true, 22: true}

func guests() []domain.Person {
	out := []domain.Person{}
	for _, person := range people {
		if person.Kind == domain.PersonGuest {
			out = append(out, person)
		}
	}
	return out
}

// Build returns the dataset. It is deterministic: the same seed produces the
// same show every time, so a number on screen can be checked against a row.
func Build() Dataset {
	r := rand.New(rand.NewPCG(42, 7))
	guestPool := guests()

	episodes := make([]Episode, 0, len(episodeTitles))
	for i, title := range episodeTitles {
		published := firstEpisode.AddDate(0, 0, 7*i)
		duration := 3600 + r.IntN(5400)

		episode := Episode{
			Number:          i + 1,
			Slug:            fmt.Sprintf("ep-%03d", i+1),
			Title:           title,
			YouTubeID:       fmt.Sprintf("ccc%08d", 1000+i*37),
			ThumbnailURL:    fmt.Sprintf("https://img.cafecomcaos.test/ep-%03d.jpg", i+1),
			Description:     "Episódio semanal do " + ChannelName + ". " + title + ".",
			PublishedAt:     published,
			DurationSeconds: duration,
		}

		episode.Appearances = buildAppearances(r, guestPool)
		episode.Moments = buildMoments(r, i, duration, episode.Appearances)
		episodes = append(episodes, episode)
	}

	attachOpenings(r, episodes)

	return Dataset{
		Types:    momentTypes,
		People:   people,
		Songs:    songs,
		Episodes: episodes,
		Archive:  archiveEntries(),
	}
}

func buildAppearances(r *rand.Rand, guestPool []domain.Person) []Appearance {
	out := []Appearance{
		{PersonSlug: "tico", Role: domain.RoleHost},
		{PersonSlug: "belinha", Role: domain.RoleHost},
	}

	// One guest most weeks, two now and then, none when the hosts are alone.
	count := 1
	switch n := r.IntN(10); {
	case n < 2:
		count = 0
	case n < 8:
		count = 1
	default:
		count = 2
	}

	taken := map[string]bool{}
	for len(taken) < count {
		guest := guestPool[r.IntN(len(guestPool))]
		if taken[guest.Slug] {
			continue
		}
		taken[guest.Slug] = true
		out = append(out, Appearance{
			PersonSlug: guest.Slug,
			Role:       domain.RoleGuest,
			// Now and then someone shows up and never gets interviewed,
			// which is exactly the case the counter has to exclude.
			IsInterview: r.IntN(10) > 0,
		})
	}

	if r.IntN(6) == 0 {
		out = append(out, Appearance{PersonSlug: "dona-neide", Role: domain.RoleRemote})
	}
	return out
}

func buildMoments(r *rand.Rand, index, duration int, appearances []Appearance) []Moment {
	out := []Moment{}

	for _, momentType := range momentTypes {
		count := rollCount(r, index, momentType.Slug)
		for range count {
			out = append(out, buildMoment(r, momentType.Slug, duration, appearances))
		}
	}
	return out
}

func rollCount(r *rand.Rand, index int, slug string) int {
	if counts, ok := scripted[index]; ok {
		if count, ok := counts[slug]; ok {
			return count
		}
	}
	if slug == "saiu-cantando" && silentSongs[index] {
		return 0
	}

	span := typeFrequency[slug]
	return span[0] + r.IntN(span[1]-span[0]+1)
}

func buildMoment(r *rand.Rand, slug string, duration int, appearances []Appearance) Moment {
	moment := Moment{
		TypeSlug:       slug,
		VideoTimestamp: 60 + r.IntN(max(duration-120, 60)),
		ActorSlug:      pickActor(r, appearances),
	}

	switch slug {
	case "pergunta-pra-neide":
		moment.Summary = questions[r.IntN(len(questions))] + "?"
	case "saiu-cantando":
		word := triggerWords[r.IntN(len(triggerWords))]
		song := songs[r.IntN(len(songs))]
		moment.TriggerWord = word
		moment.SongKey = song.Key
		moment.Summary = fmt.Sprintf("Alguém disse %q e o estúdio emendou %s", word, song.Title)
	case "teoria-maluca":
		moment.Summary = theories[r.IntN(len(theories))]
	case "briga-de-bar":
		moment.Summary = fights[r.IntN(len(fights))]
	case "chorou-ao-vivo":
		moment.Summary = cries[r.IntN(len(cries))]
	case "propaganda-nao-paga":
		moment.Summary = "Elogiou " + brands[r.IntN(len(brands))] + " por sete minutos sem ganhar nada"
	}
	return moment
}

// pickActor keeps the hosts in front: they carry the show, and a guest only
// occasionally sets something off.
func pickActor(r *rand.Rand, appearances []Appearance) string {
	if r.IntN(10) == 0 {
		for _, appearance := range appearances {
			if appearance.Role == domain.RoleGuest {
				return appearance.PersonSlug
			}
		}
	}
	if r.IntN(100) < 55 {
		return "tico"
	}
	return "belinha"
}

// attachOpenings airs most of the library and leaves the tail in the inbox,
// so the openings page has both states to show.
func attachOpenings(r *rand.Rand, episodes []Episode) {
	aired := len(openings) - 6

	for i := 0; i < aired; i++ {
		opening := openings[i]
		episode := &episodes[(i*7)%len(episodes)]
		episode.Openings = append(episode.Openings, Opening{
			Title:        opening.Title,
			AuthorName:   opening.Author,
			AuthorHandle: opening.Handle,
			Genre:        opening.Genre,
			AtSeconds:    5 + r.IntN(90),
			SentAt:       episode.PublishedAt.AddDate(0, 0, -(3 + r.IntN(20))),
			Status:       domain.OpeningAired,
		})
	}

	// The unaired ones belong to no episode, which is why the inbox is read
	// through its own query rather than through an episode.
	last := &episodes[len(episodes)-1]
	for _, opening := range openings[aired:] {
		last.Openings = append(last.Openings, Opening{
			Title:        opening.Title,
			AuthorName:   opening.Author,
			AuthorHandle: opening.Handle,
			Genre:        opening.Genre,
			SentAt:       last.PublishedAt.AddDate(0, 0, -r.IntN(14)),
			Status:       domain.OpeningReceived,
		})
	}
}
