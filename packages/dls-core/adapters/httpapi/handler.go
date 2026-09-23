// Package httpapi mounts the read API on PocketBase's router. The browser
// reads these routes rather than the collections, because every counter is
// computed here and the site should not have to know how.
package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"

	"dls/dls-core/domain"
	"dls/dls-core/ports"
)

const BasePath = "/api/dls"

type Handler struct {
	stats ports.StatsUseCase
}

func New(stats ports.StatsUseCase) *Handler {
	return &Handler{stats: stats}
}

// Mount registers the routes. Nothing here binds auth: the site is public
// and read-only, and the collections themselves refuse writes.
func (h *Handler) Mount(e *core.ServeEvent) {
	g := e.Router.Group(BasePath)

	g.GET("/overview", h.overview)
	g.GET("/episodes", h.episodes)
	g.GET("/episodes/{slug}", h.episode)
	g.GET("/openings", h.openings)
	g.GET("/rankings", h.rankings)
	g.GET("/archive", h.archive)
	g.GET("/archive/{slug}", h.archiveEntry)
}

func (h *Handler) overview(e *core.RequestEvent) error {
	overview, err := h.stats.Overview(e.Request.Context(), filterOf(e))
	if err != nil {
		return fail(e, err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"tally":    toTallyView(overview.Tally),
		"types":    toMomentTypeViews(overview.Types),
		"episodes": overview.Episodes,
		"first":    toEpisodeViewPtr(overview.First),
		"latest":   toEpisodeViewPtr(overview.Latest),
	})
}

func (h *Handler) episodes(e *core.RequestEvent) error {
	episodes, err := h.stats.Episodes(e.Request.Context(), filterOf(e))
	if err != nil {
		return fail(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{"episodes": toEpisodeViews(episodes)})
}

func (h *Handler) episode(e *core.RequestEvent) error {
	detail, err := h.stats.Episode(e.Request.Context(), e.Request.PathValue("slug"))
	if err != nil {
		return fail(e, err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"episode":     toEpisodeView(detail.Episode),
		"moments":     toMomentViews(detail.Moments),
		"appearances": toAppearanceViews(detail.Appearances),
		"people":      toPersonViews(detail.People),
		"songs":       toSongViews(detail.Songs),
		"openings":    toOpeningViews(detail.Openings),
	})
}

func (h *Handler) openings(e *core.RequestEvent) error {
	openings, err := h.stats.Openings(e.Request.Context(), filterOf(e))
	if err != nil {
		return fail(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{"openings": toOpeningViews(openings)})
}

func (h *Handler) rankings(e *core.RequestEvent) error {
	rankings, err := h.stats.Rankings(e.Request.Context(), filterOf(e))
	if err != nil {
		return fail(e, err)
	}
	return e.JSON(http.StatusOK, toRankingsView(rankings))
}

func (h *Handler) archive(e *core.RequestEvent) error {
	entries, err := h.stats.Archive(e.Request.Context(), filterOf(e))
	if err != nil {
		return fail(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{"entries": toArchiveViews(entries)})
}

func (h *Handler) archiveEntry(e *core.RequestEvent) error {
	entry, err := h.stats.ArchiveEntry(e.Request.Context(), e.Request.PathValue("slug"))
	if err != nil {
		return fail(e, err)
	}
	return e.JSON(http.StatusOK, toArchiveView(entry))
}

// filterOf reads the query string every list route shares. An unparseable
// value is ignored rather than rejected: a bookmarked URL with a stale
// parameter should still show the page.
func filterOf(e *core.RequestEvent) domain.Filter {
	query := e.Request.URL.Query()

	filter := domain.Filter{
		Genre: query.Get("genre"),
		Kind:  query.Get("kind"),
	}
	if from, ok := parseDate(query.Get("from")); ok {
		filter.From = &from
	}
	if to, ok := parseDate(query.Get("to")); ok {
		filter.To = &to
	}
	if limit, err := strconv.Atoi(query.Get("limit")); err == nil && limit > 0 {
		filter.Limit = limit
	}
	return filter
}

// parseDate takes a plain date or a full timestamp, because a filter in a
// URL is typed by people as often as it is generated.
func parseDate(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func fail(e *core.RequestEvent, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return e.NotFoundError("not found", nil)
	case errors.Is(err, domain.ErrValidation):
		return e.BadRequestError(err.Error(), nil)
	default:
		return e.InternalServerError("request failed", err)
	}
}
