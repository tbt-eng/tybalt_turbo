<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/litestream_check.go","sha256":"33b7f56e62802b02c504477e97583c90670912b183b13bb79609802687c82395","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/litestream_check.go

## Litestream check defaults when settings are missing or invalid

Lines 22-97

Turbo reads the litestream settings from app_config (key "litestream"). It checks the S3 prefix "litestream/0009/" and treats a snapshot as stale after 35 minutes. If the setting is missing or cannot be read, Turbo uses these defaults.

## Litestream snapshot prefix setting

Lines 99-101

Turbo uses snapshot_prefix from the litestream setting as the S3 prefix to check. An empty or non-text value is ignored and the default prefix "litestream/0009/" stays.

## Litestream stale threshold must be above zero

Lines 103-105

Turbo uses stale_threshold_minutes as the age in minutes before it alerts. A value that is zero, negative or not a number is ignored and the 35-minute default stays. The value is converted with a whole-number duration, so a fraction of a minute may be truncated.

## Litestream alert email and restart settings

Lines 107-113

Turbo sends alert emails only to the address in alert_email, and an empty value means no address is set. restart_on_stale is a true/false setting that defaults to false. Any value that is not true or false is ignored.

## Healthy check clears alert memory and sends nothing

Lines 126-129

When replication is healthy, Turbo writes no log entry and sends no alert. It also clears the remembered alert, so the next problem triggers a new email.

## Stale replication can restart the app

Lines 131-137

If the status is stale_snapshot and restart_on_stale is true, Turbo logs an error and exits the process with code 1 so Fly restarts it. This happens after the alert step. Other unhealthy statuses never restart the app.

## Litestream S3 region default and environment settings

Lines 140-153

Turbo reads the bucket, access key, secret key and endpoint from LITESTREAM_* environment variables. The region comes from LITESTREAM_REGION and defaults to "ca-central-1" when empty.

## Litestream check needs bucket and credentials

Lines 155-181

The check needs LITESTREAM_BUCKET, LITESTREAM_ACCESS_KEY_ID and LITESTREAM_SECRET_ACCESS_KEY. If any is empty, Turbo does not contact S3 and reports missing_credentials. The missing names are listed in alphabetical order.

## Litestream S3 listing failure

Lines 183-194

If Turbo cannot list the S3 objects, it reports s3_list_failed. The listing has a 30-second time limit.

## Litestream finds no snapshots

Lines 196-202

If no object under the prefix has a modified time, Turbo reports no_snapshots_found.

## Snapshot is stale only when older than the threshold

Lines 204-218

Turbo measures the age of the newest snapshot from its last-modified time. A snapshot whose age equals the threshold is still healthy. Turbo reports stale_snapshot only when the age is greater than the threshold.

## Newest snapshot is chosen by last-modified time

Lines 227-255

Turbo lists every object under the prefix across all pages. It picks the one with the latest last-modified time, and objects without one are skipped. A custom endpoint is used only if it starts with "http", and it then uses path-style addressing.

## Unhealthy check logs the problem

Lines 262-290

For each unhealthy status, Turbo writes a log entry. Missing credentials is logged as a warning and the other statuses as errors.

## Alert email is sent only when an address is set

Lines 292-295

Turbo sends an alert email only when alert_email is set. Without it, Turbo only logs the problem.

## Duplicate alerts are suppressed for the same condition

Lines 297-311

Turbo sends one alert email per condition and remembers its fingerprint. The fingerprint holds the status plus the missing variable names or the prefix. While the same condition continues, Turbo does not send another email. A different condition sends a new email. If sending fails, Turbo does not remember it and tries again at the next check.

## Alert email sender and content by status

Lines 313-361

The alert email goes from the sender name and address in the app's settings. Its subject and body depend on the status: missing credentials, S3 listing error, no snapshots found or stale. Each body says the application will continue running. The stale body says "restart automatically" when restart_on_stale is on.

## Litestream alert memory holds only the last alert sent

Lines 364-383

Turbo remembers one value: the fingerprint (lastSentFingerprint) of the last Litestream alert email it sent. A new successful send replaces it, so a different problem is not suppressed and gets its own email. A condition that returns after another problem was alerted in between also gets a new email. The value is held in memory in this file, and Reset sets it to empty.
