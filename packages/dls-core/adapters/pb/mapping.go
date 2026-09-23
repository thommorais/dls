package pb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"

	"dls/dls-core/domain"
)

// domainNotFound is the sentinel every repository returns for an absent row.
var domainNotFound = domain.ErrNotFound

// mapErr translates PocketBase storage errors into the domain sentinels the
// services and the HTTP adapter understand. PocketBase surfaces a missing
// record as sql.ErrNoRows from the underlying query.
func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domainNotFound
	}
	return err
}

// strSlice reads a JSON string-array field. A field that is empty, null or
// unparseable yields nil rather than an error: a malformed tag list should
// not make a record unreadable.
func strSlice(rec *core.Record, field string) []string {
	raw := rec.GetString(field)
	if raw == "" || raw == "null" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func jsonMap(rec *core.Record, field string) map[string]any {
	raw := rec.GetString(field)
	if raw == "" || raw == "null" {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func setJSON(rec *core.Record, field string, v any) {
	encoded, err := json.Marshal(v)
	if err != nil {
		rec.Set(field, "[]")
		return
	}
	rec.Set(field, string(encoded))
}

// dateString renders a time the way PocketBase stores one, for use as a query
// parameter. The zero time yields an empty string, which callers treat as
// "no bound".
func dateString(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	dt, err := types.ParseDateTime(*t)
	if err != nil {
		return ""
	}
	return dt.String()
}

func setDate(rec *core.Record, field string, t time.Time) {
	if t.IsZero() {
		rec.Set(field, "")
		return
	}
	dt, err := types.ParseDateTime(t)
	if err != nil {
		rec.Set(field, "")
		return
	}
	rec.Set(field, dt)
}

func timeOf(rec *core.Record, field string) time.Time {
	return rec.GetDateTime(field).Time()
}

func timePtrOf(rec *core.Record, field string) *time.Time {
	dt := rec.GetDateTime(field)
	if dt.IsZero() {
		return nil
	}
	t := dt.Time()
	return &t
}
