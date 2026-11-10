<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/litestream_check.go","sha256":"33b7f56e62802b02c504477e97583c90670912b183b13bb79609802687c82395","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/litestream_check.go

## Litestream check defaults when settings are missing or invalid

Lines 22-97

Asked as: litestream settings missing; default prefix and threshold; litestream config not read; litestream app_config defaults

Turbo reads the litestream settings from app_config (key "litestream"). It checks the S3 prefix "litestream/0009/" and treats a snapshot as stale after 35 minutes. If the setting is missing or cannot be read, Turbo uses these defaults.

## Litestream snapshot prefix setting

Lines 99-101

Asked as: wrong litestream prefix; change snapshot prefix; litestream folder path setting; empty snapshot prefix

Turbo uses snapshot_prefix from the litestream setting as the S3 prefix to check. An empty or non-text value is ignored and the default prefix "litestream/0009/" stays.

## Litestream stale threshold must be above zero

Lines 103-105

Asked as: stale threshold minutes; how long before litestream alerts; litestream alert after how many minutes; zero threshold litestream

Turbo uses stale_threshold_minutes as the age in minutes before it alerts. A value that is zero, negative or not a number is ignored and the 35-minute default stays. The value is converted with a whole-number duration, so a fraction of a minute may be truncated.

## Litestream alert email and restart settings

Lines 107-113

Asked as: where litestream alert emails go; restart on stale setting; litestream alert address empty; turn on litestream auto restart

Turbo sends alert emails only to the address in alert_email, and an empty value means no address is set. restart_on_stale is a true/false setting that defaults to false. Any value that is not true or false is ignored.

## Healthy check clears alert memory and sends nothing

Lines 126-129

Asked as: replication healthy no alert; litestream alert went away; clear litestream alert memory; backup OK so no email

When replication is healthy, Turbo writes no log entry and sends no alert. It also clears the remembered alert, so the next problem triggers a new email.

## Stale replication can restart the app

Lines 131-137

Asked as: litestream restarts the app; app exits when backup stale; fly restarts turbo litestream; stale snapshot app crash

If the status is stale_snapshot and restart_on_stale is true, Turbo logs an error and exits the process with code 1 so Fly restarts it. This happens after the alert step. Other unhealthy statuses never restart the app.

## Litestream S3 region default and environment settings

Lines 140-153

Asked as: litestream region setting; litestream bucket env vars; LITESTREAM_REGION blank; litestream s3 endpoint config

Turbo reads the bucket, access key, secret key and endpoint from LITESTREAM_* environment variables. The region comes from LITESTREAM_REGION and defaults to "ca-central-1" when empty.

## Litestream check needs bucket and credentials

Lines 155-181

Asked as: missing_credentials litestream; litestream bucket not set; litestream access key missing; litestream check not contacting s3

The check needs LITESTREAM_BUCKET, LITESTREAM_ACCESS_KEY_ID and LITESTREAM_SECRET_ACCESS_KEY. If any is empty, Turbo does not contact S3 and reports missing_credentials. The missing names are listed in alphabetical order.

## Litestream S3 listing failure

Lines 183-194

Asked as: cannot list s3 backups; litestream s3 timeout; s3_list_failed error; backup bucket unreachable

If Turbo cannot list the S3 objects, it reports s3_list_failed. The listing has a 30-second time limit.

## Litestream finds no snapshots

Lines 196-202

Asked as: no litestream backups found; empty backup bucket; no_snapshots_found; litestream nothing in bucket

If no object under the prefix has a modified time, Turbo reports no_snapshots_found.

## Snapshot is stale only when older than the threshold

Lines 204-218

Asked as: when is snapshot stale; backup exactly at threshold; litestream old backup warning; snapshot too old alert

Turbo measures the age of the newest snapshot from its last-modified time. A snapshot whose age equals the threshold is still healthy. Turbo reports stale_snapshot only when the age is greater than the threshold.

## Newest snapshot is chosen by last-modified time

Lines 227-255

Asked as: which litestream backup is newest; litestream pagination objects; custom s3 endpoint ignored; litestream path style url

Turbo lists every object under the prefix across all pages. It picks the one with the latest last-modified time, and objects without one are skipped. A custom endpoint is used only if it starts with "http", and it then uses path-style addressing.

## Unhealthy check logs the problem

Lines 262-290

Asked as: litestream log warning; litestream problem not logged; unhealthy backup log entry; litestream error in logs

For each unhealthy status, Turbo writes a log entry. Missing credentials is logged as a warning and the other statuses as errors.

## Alert email is sent only when an address is set

Lines 292-295

Asked as: no email for litestream problem; litestream alert not sent; litestream only logs problem; no alert address litestream

Turbo sends an alert email only when alert_email is set. Without it, Turbo only logs the problem.

## Duplicate alerts are suppressed for the same condition

Lines 297-311

Asked as: repeated litestream alert emails; stop litestream spam emails; same backup problem emailed again; litestream duplicate email

Turbo sends one alert email per condition and remembers its fingerprint. The fingerprint holds the status plus the missing variable names or the prefix. While the same condition continues, Turbo does not send another email. A different condition sends a new email. If sending fails, Turbo does not remember it and tries again at the next check.

## Alert email sender and content by status

Lines 313-361

Asked as: litestream alert email subject; who sends litestream emails; litestream email body wording; restart automatically email text

The alert email goes from the sender name and address in the app's settings. Its subject and body depend on the status: missing credentials, S3 listing error, no snapshots found or stale. Each body says the application will continue running. The stale body says "restart automatically" when restart_on_stale is on.

## Litestream alert memory holds only the last alert sent

Lines 364-383

Asked as: litestream alert after other problem; alert not resent when problem returns; litestream last alert memory; new email for returning backup problem

Turbo remembers one value: the fingerprint (lastSentFingerprint) of the last Litestream alert email it sent. A new successful send replaces it, so a different problem is not suppressed and gets its own email. A condition that returns after another problem was alerted in between also gets a new email. The value is held in memory in this file, and Reset sets it to empty.
