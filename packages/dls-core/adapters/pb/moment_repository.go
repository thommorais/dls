package pb

import (
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"dls/dls-core/domain"
	"dls/dls-core/ports"
)

type MomentTypeRepository struct {
	app core.App
}

func NewMomentTypeRepository(app core.App) *MomentTypeRepository {
	return &MomentTypeRepository{app: app}
}

var _ ports.MomentTypeRepository = (*MomentTypeRepository)(nil)

func toMomentType(rec *core.Record) domain.MomentType {
	return domain.MomentType{
		ID:          domain.MomentTypeID(rec.Id),
		Slug:        rec.GetString("slug"),
		Label:       rec.GetString("label"),
		Description: rec.GetString("description"),
		Color:       rec.GetString("color"),
		Position:    rec.GetInt("position"),
	}
}

func (r *MomentTypeRepository) List(_ context.Context) ([]domain.MomentType, error) {
	records, err := r.app.FindAllRecords(ColMomentTypes)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]domain.MomentType, 0, len(records))
	for _, rec := range records {
		out = append(out, toMomentType(rec))
	}
	return out, nil
}

type MomentRepository struct {
	app   core.App
	types *MomentTypeRepository
}

func NewMomentRepository(app core.App, types *MomentTypeRepository) *MomentRepository {
	return &MomentRepository{app: app, types: types}
}

var _ ports.MomentRepository = (*MomentRepository)(nil)

func (r *MomentRepository) ListByEpisodes(ctx context.Context, episodes []domain.EpisodeID) ([]domain.Moment, error) {
	records, err := r.readScoped(ColMoments, episodes)
	if err != nil {
		return nil, err
	}

	// The type is resolved here so nothing above this package joins. The
	// table is a handful of rows, so one read covers every moment.
	types, err := r.types.List(ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[domain.MomentTypeID]domain.MomentType, len(types))
	for _, momentType := range types {
		byID[momentType.ID] = momentType
	}

	out := make([]domain.Moment, 0, len(records))
	for _, rec := range records {
		typeID := domain.MomentTypeID(rec.GetString("type"))
		out = append(out, domain.Moment{
			ID:             domain.MomentID(rec.Id),
			EpisodeID:      domain.EpisodeID(rec.GetString("episode")),
			TypeID:         typeID,
			Type:           byID[typeID],
			At:             timeOf(rec, "at"),
			VideoTimestamp: rec.GetInt("video_timestamp"),
			ActorID:        domain.PersonID(rec.GetString("actor")),
			Summary:        rec.GetString("summary"),
			TriggerWord:    rec.GetString("trigger_word"),
			SongID:         domain.SongID(rec.GetString("song")),
			ClipURL:        rec.GetString("clip_url"),
			Source:         domain.Source(rec.GetString("source")),
			Confidence:     rec.GetFloat("confidence"),
			Payload:        jsonMap(rec, "payload"),
		})
	}
	return out, nil
}

// readScoped is the shared shape of every by-episode read: nil means the
// whole table, a list means those episodes, and an empty list means nothing.
func (r *MomentRepository) readScoped(collection string, episodes []domain.EpisodeID) ([]*core.Record, error) {
	return readScoped(r.app, collection, episodes)
}

func readScoped(app core.App, collection string, episodes []domain.EpisodeID) ([]*core.Record, error) {
	if episodes == nil {
		records, err := app.FindAllRecords(collection)
		return records, mapErr(err)
	}
	if len(episodes) == 0 {
		return nil, nil
	}

	values := make([]any, 0, len(episodes))
	for _, id := range episodes {
		values = append(values, string(id))
	}
	records, err := app.FindAllRecords(collection, dbx.In("episode", values...))
	return records, mapErr(err)
}

type AppearanceRepository struct {
	app core.App
}

func NewAppearanceRepository(app core.App) *AppearanceRepository {
	return &AppearanceRepository{app: app}
}

var _ ports.AppearanceRepository = (*AppearanceRepository)(nil)

func (r *AppearanceRepository) ListByEpisodes(_ context.Context, episodes []domain.EpisodeID) ([]domain.Appearance, error) {
	records, err := readScoped(r.app, ColAppearances, episodes)
	if err != nil {
		return nil, err
	}

	out := make([]domain.Appearance, 0, len(records))
	for _, rec := range records {
		out = append(out, domain.Appearance{
			ID:          domain.AppearanceID(rec.Id),
			EpisodeID:   domain.EpisodeID(rec.GetString("episode")),
			PersonID:    domain.PersonID(rec.GetString("person")),
			Role:        domain.AppearanceRole(rec.GetString("role")),
			IsInterview: rec.GetBool("is_interview"),
		})
	}
	return out, nil
}
