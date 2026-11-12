package hooks

import (
	"fmt"
	"math"
	"net/http"
	"sort"

	"tybalt/errs"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Near-match limits for callers that show "differs from X in N roles". A sheet
// that differs in more roles than this is not shown as close.
const (
	RateSheetNearMatchMaxRoles = 5
	RateSheetNearMatchLimit    = 3
)

// RateSheetEntryValues is one role's rates on a rate sheet, or on a set of
// entries that is not saved yet (a new sheet, or a draft in a request).
type RateSheetEntryValues struct {
	Role         string  `db:"role" json:"role"`
	Rate         float64 `db:"rate" json:"rate"`
	OvertimeRate float64 `db:"overtime_rate" json:"overtime_rate"`
}

// RateSheetMatch is an active rate sheet compared with a set of entries.
// DifferingRoles counts the roles whose rate or overtime rate differ, plus the
// roles that are on one side only. It is 0 for an exact match.
type RateSheetMatch struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Revision       int    `json:"revision"`
	EffectiveDate  string `json:"effective_date"`
	DifferingRoles int    `json:"differing_roles"`
}

// RateSheetComparison is the result of CompareWithActiveRateSheets. Exact holds
// the active sheets with identical entries. Near holds every other active sheet
// that was compared, closest first; use Closest to trim it.
type RateSheetComparison struct {
	Exact []RateSheetMatch `json:"exact"`
	Near  []RateSheetMatch `json:"near"`
}

// Closest returns at most limit near matches that differ in at most maxRoles
// roles, closest first.
func (c RateSheetComparison) Closest(maxRoles int, limit int) []RateSheetMatch {
	result := []RateSheetMatch{}
	for _, match := range c.Near {
		if len(result) >= limit {
			break
		}
		if match.DifferingRoles > maxRoles {
			break // Near is sorted, so no later match is closer
		}
		result = append(result, match)
	}
	return result
}

// rateSheetEntryRow is one row of the active sheets query. A sheet with no
// entries gives one row with an empty role.
type rateSheetEntryRow struct {
	SheetID       string  `db:"sheet_id"`
	Name          string  `db:"name"`
	Revision      int     `db:"revision"`
	EffectiveDate string  `db:"effective_date"`
	Role          string  `db:"role"`
	Rate          float64 `db:"rate"`
	OvertimeRate  float64 `db:"overtime_rate"`
}

// rateCents is a rate in whole cents, so that the comparison does not depend
// on how a float was rounded on its way in.
type rateCents struct {
	rate     int64
	overtime int64
}

func toRateCents(rate float64, overtime float64) rateCents {
	return rateCents{
		rate:     int64(math.Round(rate * 100)),
		overtime: int64(math.Round(overtime * 100)),
	}
}

// CompareWithActiveRateSheets compares entries with every active rate sheet
// except those named excludeName, so that a revision is never compared with
// its own earlier revisions. Two sheets are identical when they have the same
// set of roles and every role's rate and overtime rate are equal. Inactive
// sheets are never compared. app may be a transaction app.
func CompareWithActiveRateSheets(app core.App, entries []RateSheetEntryValues, excludeName string) (RateSheetComparison, error) {
	result := RateSheetComparison{Exact: []RateSheetMatch{}, Near: []RateSheetMatch{}}

	wanted := make(map[string]rateCents, len(entries))
	for _, entry := range entries {
		wanted[entry.Role] = toRateCents(entry.Rate, entry.OvertimeRate)
	}

	rows := []rateSheetEntryRow{}
	err := app.DB().NewQuery(`
		SELECT
			rs.id AS sheet_id,
			rs.name AS name,
			rs.revision AS revision,
			rs.effective_date AS effective_date,
			COALESCE(rse.role, '') AS role,
			COALESCE(rse.rate, 0) AS rate,
			COALESCE(rse.overtime_rate, 0) AS overtime_rate
		FROM rate_sheets rs
		LEFT JOIN rate_sheet_entries rse ON rse.rate_sheet = rs.id
		WHERE rs.active = true AND rs.name != {:excludeName}
		ORDER BY rs.id
	`).Bind(dbx.Params{
		"excludeName": excludeName,
	}).All(&rows)
	if err != nil {
		return result, err
	}

	type sheet struct {
		match   RateSheetMatch
		entries map[string]rateCents
	}
	sheets := []*sheet{}
	byID := map[string]*sheet{}
	for _, row := range rows {
		s, ok := byID[row.SheetID]
		if !ok {
			s = &sheet{
				match: RateSheetMatch{
					ID:            row.SheetID,
					Name:          row.Name,
					Revision:      row.Revision,
					EffectiveDate: row.EffectiveDate,
				},
				entries: map[string]rateCents{},
			}
			byID[row.SheetID] = s
			sheets = append(sheets, s)
		}
		if row.Role != "" {
			s.entries[row.Role] = toRateCents(row.Rate, row.OvertimeRate)
		}
	}

	for _, s := range sheets {
		differing := 0
		for role, want := range wanted {
			if have, ok := s.entries[role]; !ok || have != want {
				differing++
			}
		}
		for role := range s.entries {
			if _, ok := wanted[role]; !ok {
				differing++
			}
		}
		s.match.DifferingRoles = differing
		if differing == 0 {
			result.Exact = append(result.Exact, s.match)
		} else {
			result.Near = append(result.Near, s.match)
		}
	}

	sortMatches(result.Exact)
	sortMatches(result.Near)
	return result, nil
}

func sortMatches(matches []RateSheetMatch) {
	sort.SliceStable(matches, func(i, j int) bool {
		a, b := matches[i], matches[j]
		if a.DifferingRoles != b.DifferingRoles {
			return a.DifferingRoles < b.DifferingRoles
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.Revision < b.Revision
	})
}

// LoadRateSheetEntryValues returns the stored entries of a rate sheet, in the
// form that CompareWithActiveRateSheets takes.
func LoadRateSheetEntryValues(app core.App, rateSheetId string) ([]RateSheetEntryValues, error) {
	entries := []RateSheetEntryValues{}
	err := app.DB().NewQuery(`
		SELECT role, rate, overtime_rate
		FROM rate_sheet_entries
		WHERE rate_sheet = {:rateSheetId}
	`).Bind(dbx.Params{
		"rateSheetId": rateSheetId,
	}).All(&entries)
	return entries, err
}

// DuplicateRateSheetError is the error for a sheet whose rates are identical
// to an active sheet. field is the key the UI shows the error under. The
// error's data lists every match, so the UI can link to the existing sheet.
func DuplicateRateSheetError(field string, matches []RateSheetMatch) *errs.HookError {
	first := matches[0]
	message := fmt.Sprintf(
		"the rates are identical to active rate sheet %q (rev. %d); use that sheet instead",
		first.Name, first.Revision,
	)
	return &errs.HookError{
		Status:  http.StatusConflict,
		Message: "a rate sheet with identical rates is already active",
		Data: map[string]errs.CodeError{
			field: {
				Code:    "duplicate_rate_sheet",
				Message: message,
				Data:    map[string]any{"matches": matches},
			},
		},
	}
}
