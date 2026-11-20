<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/computer_login_retention.go","sha256":"493658497b0372af15103107bf1800238c3fdf95927508dcae3d18f3644c37e2","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/computer_login_retention.go

## Old computer logins and raw logins are deleted after the retention period

Lines 13-40

Asked as: old computer logins deleted; login retention period; delete old logins daily; raw logins cleanup; login history purge

A daily job deletes every computer login (computer_logins) and every raw login (computer_raw_logins) that was created before the cutoff. The cutoff is now (UTC) minus COMPUTER_LOGIN_RETENTION_DAYS days. The number of days is set in the constants file, not here.

## Login cleanup never changes computers

Lines 19-40

Asked as: computer lost last login; last user cleared on computer; computer login removed after cleanup; last_login missing

The cleanup deletes only logins and raw logins. It does not touch computers, so a computer keeps its last login (last_login) and its last user even after all its logins are deleted.

## Login cleanup deletes with SQL and skips record hooks

Lines 19-53

Asked as: login deletion no hooks; bulk delete logins sql; deleted logins without callbacks; cleanup skips hooks

The cleanup deletes rows directly with SQL, not record by record. No record hook runs for these deletions.

## Login cleanup stops at the first failure

Lines 28-40

Asked as: login cleanup failed halfway; raw logins not deleted; partial login deletion; cleanup error stops midway

The job deletes logins first, then raw logins. If deleting logins fails, Turbo does not try to delete raw logins. If deleting raw logins fails, the logins already deleted stay deleted.

## A login created exactly at the cutoff is kept

Lines 45-53

Asked as: login exactly at cutoff; login same time as cutoff; why login not deleted at cutoff; boundary login kept

Turbo deletes a login or raw login only when its created time is strictly earlier than the cutoff. A row created exactly at the cutoff stays.
