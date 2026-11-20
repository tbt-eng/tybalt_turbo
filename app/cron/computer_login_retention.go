package cron

import (
	"time"

	"tybalt/constants"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// computerLoginRetentionCutoff is the oldest created time that a login or a
// raw login can have and stay. A row created before it is deleted.
func computerLoginRetentionCutoff(now time.Time) time.Time {
	return now.UTC().AddDate(0, 0, -constants.COMPUTER_LOGIN_RETENTION_DAYS)
}

// deleteExpiredComputerLogins deletes the computer_logins and
// computer_raw_logins rows created more than COMPUTER_LOGIN_RETENTION_DAYS
// before now. It never touches computers: the last_login and the last user of
// a computer stay when all its logins are gone.
//
// It deletes with SQL, not record by record. No hook watches these
// collections and no row refers to a login or a raw login, so a record delete
// would have nothing more to do. created is stored as text in the PocketBase
// layout, so the cutoff is formatted the same way and compared as text.
func deleteExpiredComputerLogins(app core.App, now time.Time) (logins int64, rawLogins int64, err error) {
	cutoff := computerLoginRetentionCutoff(now).Format(types.DefaultDateLayout)

	logins, err = deleteCreatedBefore(app, "computer_logins", cutoff)
	if err != nil {
		return 0, 0, err
	}
	rawLogins, err = deleteCreatedBefore(app, "computer_raw_logins", cutoff)
	if err != nil {
		return logins, 0, err
	}
	return logins, rawLogins, nil
}

// deleteCreatedBefore deletes the rows of table created before cutoff and
// returns how many it deleted. table is always a constant name from this file,
// never input.
func deleteCreatedBefore(app core.App, table string, cutoff string) (int64, error) {
	result, err := app.DB().NewQuery(
		"DELETE FROM " + table + " WHERE created < {:cutoff}",
	).Bind(dbx.Params{"cutoff": cutoff}).Execute()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// runComputerLoginRetention is the body of the daily cron job.
func runComputerLoginRetention(app core.App) {
	now := time.Now()
	cutoff := computerLoginRetentionCutoff(now).Format(types.DefaultDateLayout)
	logins, rawLogins, err := deleteExpiredComputerLogins(app, now)
	if err != nil {
		app.Logger().Error(
			"computer login retention failed",
			"error", err,
			"cutoff", cutoff,
			"deleted_logins", logins,
			"deleted_raw_logins", rawLogins,
		)
		return
	}
	app.Logger().Info(
		"deleted expired computer logins",
		"cutoff", cutoff,
		"deleted_logins", logins,
		"deleted_raw_logins", rawLogins,
	)
}
