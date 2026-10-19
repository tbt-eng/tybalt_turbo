package routes

import _ "embed"

// singleJobScopeQuery also keeps the Staff and Divisions summaries on the shared pricing path.
const singleJobScopeQuery = `WITH selected_jobs AS (SELECT * FROM jobs WHERE id = {:job_id})`

//go:embed job_priced_time_entries.sql
var jobPricedTimeEntriesQuery string

// JobTimeValuation identifies employee-rate estimates and unpriced hours in
// the API and CSV exports. Percent is null for incomplete or zero totals.
type JobTimeValuation struct {
	// RateSheets names the sheets that priced the row's entries, such as
	// "2026 Standard Rates (rev. 0)", joined by "; ". Empty when none applied.
	RateSheets          string   `db:"rate_sheets" json:"rate_sheets"`
	Value               float64  `db:"value" json:"value"`
	Total               float64  `db:"total" json:"total"`
	Percent             *float64 `db:"percent" json:"percent"`
	EstimatedHours      float64  `db:"estimated_hours" json:"estimated_hours"`
	TotalEstimatedHours float64  `db:"total_estimated_hours" json:"total_estimated_hours"`
	UnpricedHours       float64  `db:"unpriced_hours" json:"unpriced_hours"`
	TotalUnpricedHours  float64  `db:"total_unpriced_hours" json:"total_unpriced_hours"`
	ValueStatus         string   `db:"value_status" json:"value_status"`
	TotalStatus         string   `db:"total_status" json:"total_status"`
}
