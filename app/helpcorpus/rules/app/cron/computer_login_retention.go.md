<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/computer_login_retention.go","sha256":"d25b1892f506fb4733a52408c9e3050b26bb87d13158cdc36027741ad56c1623","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/computer_login_retention.go

## Computer logins are deleted after the retention period

Lines 13-42

Asked as: old computer logins deleted; login records expire; how long logins are kept; login retention period; computer login cleanup

A daily job deletes computer login records (computer_logins) created before the cutoff. The cutoff is now minus COMPUTER_LOGIN_RETENTION_DAYS days, in UTC. The number of days is set in the constants file, which this file does not show.

## Raw logins with a computer name never expire

Lines 20-27

Asked as: raw login never deleted; old raw logins kept; radiator keeps newest login; computer name login not expiring; why login still there

The retention job does not delete a raw login that has a computer name, however old it is. The radiator endpoint keeps the newest one per computer name and deletes it when that computer sends a valid report.

## Retention job leaves computers unchanged

Lines 29-31

Asked as: computer last login after cleanup; last user cleared on computer; computer record unchanged by job; login deletion effect on computer; computer keeps last login

The retention job never changes computers. A computer keeps its last login (last_login) and its last user after all its login records are deleted.

## Retention job deletes directly in the database

Lines 32-37

Asked as: retention job no hooks; cleanup runs raw SQL; login delete bypasses hooks; cutoff date format text; database delete login cleanup

The retention job deletes rows with SQL, not record by record. No hooks run on these deletes. The job compares the created time with the cutoff as text, so it formats the cutoff in the PocketBase date layout.

## Retention job stops at the first error

Lines 43-53

Asked as: cleanup job failed halfway; raw logins not deleted after error; partial login deletion; retention job error log; cleanup stopped early

If deleting logins fails, the job does not try to delete the raw logins and logs the error. If deleting raw logins fails, the logins already deleted stay deleted. On success the job logs the cutoff and both counts.

## Raw logins without a computer name expire like logins

Lines 46-49

Asked as: raw reports without computer name; empty computer name raw login; unmatched raw login deleted; old raw reports removed; raw login daily cleanup

The daily job deletes raw logins (computer_raw_logins) that have no computer name (an empty or null computer_name) and were created before the cutoff. These reports cannot be matched to a computer, so only this job removes them.
