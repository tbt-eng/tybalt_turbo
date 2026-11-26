package cron

import (
	"time"

	"tybalt/absorb"
	"tybalt/constants"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// absorbHistoryRetentionCutoff is the oldest commit time that a row of
// absorb_history can have and stay. A row created before it is deleted.
func absorbHistoryRetentionCutoff(now time.Time) time.Time {
	return now.UTC().AddDate(0, 0, -constants.ABSORB_HISTORY_RETENTION_DAYS)
}

// deleteExpiredAbsorbHistory deletes the absorb_history rows of absorbs that
// were committed more than ABSORB_HISTORY_RETENTION_DAYS before now, and
// returns how many it deleted. created is the commit time, stored as text in
// the PocketBase layout, so the cutoff is compared as text. No hook watches
// the collection and no row depends on it, so it deletes with plain SQL.
func deleteExpiredAbsorbHistory(app core.App, now time.Time) (int64, error) {
	result, err := app.DB().NewQuery(
		"DELETE FROM " + absorb.HistoryCollection + " WHERE created < {:cutoff}",
	).Bind(dbx.Params{
		"cutoff": absorbHistoryRetentionCutoff(now).Format(types.DefaultDateLayout),
	}).Execute()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// runAbsorbHistoryRetention is the body of the daily cron job.
func runAbsorbHistoryRetention(app core.App) {
	now := time.Now()
	cutoff := absorbHistoryRetentionCutoff(now).Format(types.DefaultDateLayout)
	deleted, err := deleteExpiredAbsorbHistory(app, now)
	if err != nil {
		app.Logger().Error("absorb history retention failed", "error", err, "cutoff", cutoff)
		return
	}
	app.Logger().Info("deleted old absorb history", "cutoff", cutoff, "deleted", deleted)
}
