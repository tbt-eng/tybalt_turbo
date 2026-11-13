<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"94ba6bc4bd190b3a75c7e1aaa324ce8eaaf71018bf026d239f2057db8863b16f","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Augmented rate sheets cannot be created, changed or deleted through the API

Lines 5-17

Asked as: cannot edit rate sheet; can't change augmented rate sheet; read only rate sheets; api create rate sheet fails; delete rate sheet not allowed

The rate_sheets_augmented view has no create, update or delete rules. Nobody can create, update or delete its records through the API, because it is a read-only view. Turbo's own server code is not limited by these rules.

## Any signed-in person can list augmented rate sheets

Lines 7-11

Asked as: list rate sheets not showing; see all rate sheets; rate sheets list signed in; anonymous rate sheet access; who can view rate sheets

Any signed-in person can list records in the rate_sheets_augmented view (list rule: @request.auth.id != ""). Anyone who is not signed in cannot.

## Any signed-in person can view one augmented rate sheet

Lines 13-17

Asked as: open single rate sheet; view one rate sheet; rate sheet details not loading; can't see rate sheet record; rate sheet access denied

Any signed-in person can view one record in the rate_sheets_augmented view (view rule: @request.auth.id != ""). Anyone who is not signed in cannot.

## Job count counts jobs that use a rate sheet now or in the past

Lines 30-36

Asked as: job count wrong; how many jobs use rate sheet; rate sheet usage count; jobs using old rate sheet; past jobs rate sheet count

For each rate sheet, job_count counts the distinct jobs that use it. A job counts if its current rate sheet (jobs.rate_sheet) is this sheet, or if job_rate_sheets has a row for this sheet and that job. The UNION removes duplicate pairs, so a job in both places counts once per rate sheet.

## Rate sheets with no jobs show a job count of zero

Lines 30-37

Asked as: rate sheet zero jobs; unused rate sheet count; job count shows 0; rate sheet missing from list; empty rate sheet job count

The view uses a LEFT JOIN, so a rate sheet that no job uses still appears, with job_count 0. Jobs whose current rate sheet is empty are left out of the count.
