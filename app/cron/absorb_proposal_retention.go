package cron

import (
	"time"

	"tybalt/constants"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// absorbProposalRetentionCutoff is the oldest decided_at that a decided
// absorb proposal can have and stay. One decided before it is deleted.
func absorbProposalRetentionCutoff(now time.Time) time.Time {
	return now.UTC().AddDate(0, 0, -constants.ABSORB_PROPOSAL_RETENTION_DAYS)
}

// deleteExpiredAbsorbProposals deletes the absorb proposals that were decided
// (accepted, declined, withdrawn or obsolete) more than
// ABSORB_PROPOSAL_RETENTION_DAYS before now, with their
// absorb_proposal_records rows. It returns how many proposals it deleted. A
// pending proposal has no decided_at and always stays, however old it is.
//
// The record of a committed absorb is its absorb_history row, so a closed
// proposal is not kept as history. It deletes with SQL in one transaction, not record
// by record: no hook watches these collections, and plain SQL does not run a
// cascade delete, so the record rows are deleted first. decided_at is stored
// as text in the PocketBase layout, so the cutoff is compared as text.
func deleteExpiredAbsorbProposals(app core.App, now time.Time) (int64, error) {
	params := dbx.Params{"cutoff": absorbProposalRetentionCutoff(now).Format(types.DefaultDateLayout)}
	const expired = "status != 'pending' AND decided_at != '' AND decided_at < {:cutoff}"

	var deleted int64
	err := app.RunInTransaction(func(txApp core.App) error {
		if _, err := txApp.DB().NewQuery(
			"DELETE FROM absorb_proposal_records WHERE proposal IN (SELECT id FROM absorb_proposals WHERE " + expired + ")",
		).Bind(params).Execute(); err != nil {
			return err
		}
		result, err := txApp.DB().NewQuery(
			"DELETE FROM absorb_proposals WHERE " + expired,
		).Bind(params).Execute()
		if err != nil {
			return err
		}
		deleted, err = result.RowsAffected()
		return err
	})
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// runAbsorbProposalRetention is the body of the daily cron job.
func runAbsorbProposalRetention(app core.App) {
	now := time.Now()
	cutoff := absorbProposalRetentionCutoff(now).Format(types.DefaultDateLayout)
	deleted, err := deleteExpiredAbsorbProposals(app, now)
	if err != nil {
		app.Logger().Error("absorb proposal retention failed", "error", err, "cutoff", cutoff)
		return
	}
	app.Logger().Info("deleted old absorb proposals", "cutoff", cutoff, "deleted", deleted)
}
