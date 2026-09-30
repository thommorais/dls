package pb

import (
	"context"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"dls/dls-core/domain"
	"dls/dls-core/ports"
)

type PersonRepository struct {
	app core.App
}

func NewPersonRepository(app core.App) *PersonRepository {
	return &PersonRepository{app: app}
}

var _ ports.PersonRepository = (*PersonRepository)(nil)

func (r *PersonRepository) List(_ context.Context) ([]domain.Person, error) {
	records, err := r.app.FindAllRecords(ColPeople)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]domain.Person, 0, len(records))
	for _, rec := range records {
		out = append(out, domain.Person{
			ID:     domain.PersonID(rec.Id),
			Slug:   rec.GetString("slug"),
			Name:   rec.GetString("name"),
			Kind:   domain.PersonKind(rec.GetString("kind")),
			Gender: domain.Gender(rec.GetString("gender")),
			Photo:  rec.GetString("photo"),
		})
	}
	return out, nil
}

type SongRepository struct {
	app core.App
}

func NewSongRepository(app core.App) *SongRepository {
	return &SongRepository{app: app}
}

var _ ports.SongRepository = (*SongRepository)(nil)

func (r *SongRepository) List(_ context.Context) ([]domain.Song, error) {
	records, err := r.app.FindAllRecords(ColSongs)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]domain.Song, 0, len(records))
	for _, rec := range records {
		out = append(out, domain.Song{
			ID:     domain.SongID(rec.Id),
			Title:  rec.GetString("title"),
			Artist: rec.GetString("artist"),
			Genre:  rec.GetString("genre"),
			Year:   rec.GetInt("year"),
			Link:   rec.GetString("link"),
		})
	}
	return out, nil
}

type OpeningRepository struct {
	app core.App
}

func NewOpeningRepository(app core.App) *OpeningRepository {
	return &OpeningRepository{app: app}
}

var _ ports.OpeningRepository = (*OpeningRepository)(nil)

func toOpening(rec *core.Record) domain.Opening {
	return domain.Opening{
		ID:           domain.OpeningID(rec.Id),
		Title:        rec.GetString("title"),
		AuthorName:   rec.GetString("author_name"),
		AuthorHandle: rec.GetString("author_handle"),
		Genre:        rec.GetString("genre"),
		EpisodeID:    domain.EpisodeID(rec.GetString("episode")),
		AtSeconds:    rec.GetInt("at_seconds"),
		AiredAt:      timePtrOf(rec, "aired_at"),
		MediaURL:     rec.GetString("media_url"),
	}
}

func (r *OpeningRepository) ListByEpisodes(_ context.Context, episodes []domain.EpisodeID) ([]domain.Opening, error) {
	records, err := readScoped(r.app, ColOpenings, episodes)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Opening, 0, len(records))
	for _, rec := range records {
		out = append(out, toOpening(rec))
	}
	return out, nil
}

func (r *OpeningRepository) List(_ context.Context, filter domain.Filter) ([]domain.Opening, error) {
	where := []string{"id != ''"}
	params := dbx.Params{}

	if filter.Genre != "" {
		where = append(where, "genre = {:genre}")
		params["genre"] = filter.Genre
	}
	if from := dateString(filter.From); from != "" {
		where = append(where, "aired_at >= {:from}")
		params["from"] = from
	}
	if to := dateString(filter.To); to != "" {
		where = append(where, "aired_at <= {:to}")
		params["to"] = to
	}

	records, err := r.app.FindRecordsByFilter(ColOpenings, strings.Join(where, " && "), "-aired_at", 0, 0, params)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]domain.Opening, 0, len(records))
	for _, rec := range records {
		out = append(out, toOpening(rec))
	}
	return out, nil
}

type ArchiveRepository struct {
	app core.App
}

func NewArchiveRepository(app core.App) *ArchiveRepository {
	return &ArchiveRepository{app: app}
}

var _ ports.ArchiveRepository = (*ArchiveRepository)(nil)

func toArchiveEntry(rec *core.Record) domain.ArchiveEntry {
	return domain.ArchiveEntry{
		ID:           domain.ArchiveID(rec.Id),
		Slug:         rec.GetString("slug"),
		Title:        rec.GetString("title"),
		Summary:      rec.GetString("summary"),
		Body:         rec.GetString("body"),
		Kind:         domain.ArchiveKind(rec.GetString("kind")),
		Tags:         strSlice(rec, "tags"),
		EpisodeID:    domain.EpisodeID(rec.GetString("episode")),
		PersonID:     domain.PersonID(rec.GetString("person")),
		MomentTypeID: domain.MomentTypeID(rec.GetString("moment_type")),
		SourceURL:    rec.GetString("source_url"),
		PublishedAt:  timeOf(rec, "published_at"),
	}
}

func (r *ArchiveRepository) List(_ context.Context, filter domain.Filter) ([]domain.ArchiveEntry, error) {
	where := []string{"id != ''"}
	params := dbx.Params{}

	if filter.Kind != "" {
		where = append(where, "kind = {:kind}")
		params["kind"] = filter.Kind
	}

	records, err := r.app.FindRecordsByFilter(ColArchive, strings.Join(where, " && "), "title", 0, 0, params)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]domain.ArchiveEntry, 0, len(records))
	for _, rec := range records {
		out = append(out, toArchiveEntry(rec))
	}
	return out, nil
}

func (r *ArchiveRepository) GetBySlug(_ context.Context, slug string) (domain.ArchiveEntry, error) {
	rec, err := r.app.FindFirstRecordByData(ColArchive, "slug", slug)
	if err != nil {
		return domain.ArchiveEntry{}, mapErr(err)
	}
	return toArchiveEntry(rec), nil
}
