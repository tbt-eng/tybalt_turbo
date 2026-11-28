package cron

import (
	"time"

	"tybalt/constants"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// computerLoginRetentionCutoff is the oldest created time that a login, or a
// raw login without a computer name, can have and stay. A row created before
// it is deleted.
func computerLoginRetentionCutoff(now time.Time) time.Time {
	return now.UTC().AddDate(0, 0, -constants.COMPUTER_LOGIN_RETENTION_DAYS)
}

// deleteExpiredComputerLogins deletes the computer_logins rows, and the
// computer_raw_logins rows without a computer name, created more than
// COMPUTER_LOGIN_RETENTION_DAYS before now.
//
// A raw login with a computer name never expires. The radiator endpoint keeps
// one per name, the newest report of that computer, and deletes it when the
// computer sends a valid report. A raw login without a name can't be matched
// to a computer, so only this job deletes it.
//
// It never touches computers: the last_login and the last user of a computer
// stay when all its logins are gone.
//
// It deletes with SQL, not record by record. No hook watches these
// collections and no row refers to a login or a raw login, so a record delete
// would have nothing more to do. created is stored as text in the PocketBase
// layout, so the cutoff is formatted the same way and compared as text.
func deleteExpiredComputerLogins(app core.App, now time.Time) (logins int64, rawLogins int64, err error) {
	cutoff := computerLoginRetentionCutoff(now).Format(types.DefaultDateLayout)

	logins, err = deleteCreatedBefore(app,
		"DELETE FROM computer_logins WHERE created < {:cutoff}",
		cutoff,
	)
	if err != nil {
		return 0, 0, err
	}
	rawLogins, err = deleteCreatedBefore(app,
		"DELETE FROM computer_raw_logins WHERE created < {:cutoff} AND COALESCE(computer_name, '') = ''",
		cutoff,
	)
	if err != nil {
		return logins, 0, err
	}
	return logins, rawLogins, nil
}

// deleteCreatedBefore runs a DELETE statement with cutoff bound and returns
// how many rows it deleted. statement is always a constant from this file,
// never input.
func deleteCreatedBefore(app core.App, statement string, cutoff string) (int64, error) {
	result, err := app.DB().NewQuery(statement).Bind(dbx.Params{"cutoff": cutoff}).Execute()
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
