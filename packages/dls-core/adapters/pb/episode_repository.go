package pb

import (
	"context"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"dls/dls-core/domain"
	"dls/dls-core/ports"
)

type EpisodeRepository struct {
	app core.App
}

func NewEpisodeRepository(app core.App) *EpisodeRepository {
	return &EpisodeRepository{app: app}
}

var _ ports.EpisodeRepository = (*EpisodeRepository)(nil)

func toEpisode(rec *core.Record) domain.Episode {
	return domain.Episode{
		ID:              domain.EpisodeID(rec.Id),
		YouTubeID:       rec.GetString("youtube_id"),
		Number:          rec.GetInt("number"),
		Slug:            rec.GetString("slug"),
		Title:           rec.GetString("title"),
		PublishedAt:     timeOf(rec, "published_at"),
		DurationSeconds: rec.GetInt("duration_seconds"),
		ThumbnailURL:    rec.GetString("thumbnail_url"),
		Description:     rec.GetString("description"),
	}
}

func (r *EpisodeRepository) List(_ context.Context, filter domain.Filter) ([]domain.Episode, error) {
	// An always-true base keeps the expression valid when nothing is bound,
	// which is the common case.
	where := []string{"id != ''"}
	params := dbx.Params{}

	if from := dateString(filter.From); from != "" {
		where = append(where, "published_at >= {:from}")
		params["from"] = from
	}
	if to := dateString(filter.To); to != "" {
		where = append(where, "published_at <= {:to}")
		params["to"] = to
	}

	records, err := r.app.FindRecordsByFilter(ColEpisodes, strings.Join(where, " && "), "published_at", 0, 0, params)
	if err != nil {
		return nil, mapErr(err)
	}

	out := make([]domain.Episode, 0, len(records))
	for _, rec := range records {
		out = append(out, toEpisode(rec))
	}
	return out, nil
}

func (r *EpisodeRepository) GetBySlug(_ context.Context, slug string) (domain.Episode, error) {
	rec, err := r.app.FindFirstRecordByData(ColEpisodes, "slug", slug)
	if err != nil {
		return domain.Episode{}, mapErr(err)
	}
	return toEpisode(rec), nil
}
