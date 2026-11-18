<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"8594252b1c81cd30a53af15caaf2254d3809ea08279f690aa18c58e993a28a70","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Nobody can create, update or delete augmented rate sheet records

Lines 5-17

Asked as: cannot edit augmented rate sheets; read only rate sheet view; why can't I change rate sheet records; delete augmented rate sheet

The rate_sheets_augmented view has only a list rule and a view rule. Nobody can create, update or delete its records through the API, because it is a read-only view. Turbo's own server code does not write to it either.

## Who can list rate sheets in the augmented view

Lines 7-11

Asked as: who can see rate sheet list; list rate sheets not showing; rate sheets list access; signed in list of rate sheets

Any signed-in person can list the records of the rate_sheets_augmented view. A request without a signed-in user fails the filter (`@request.auth.id != ""`).

## Who can view one rate sheet in the augmented view

Lines 13-17

Asked as: open single rate sheet; view one rate sheet; rate sheet details not loading; can't open a rate sheet

Any signed-in person can view a single record of the rate_sheets_augmented view. A request without a signed-in user fails the filter.

## Job count counts jobs that use the rate sheet directly or by period

Lines 31-39

Asked as: how many jobs use rate sheet; job count wrong; rate sheet jobs total; jobs counted twice on rate sheet

The job count (job_count) counts the distinct jobs that use the rate sheet. A job counts when its current rate sheet is this one (jobs.rate_sheet), or when this rate sheet appears in the job's rate sheet periods (job_rate_sheets). A job that matches both ways counts once, because the query uses UNION. A rate sheet with no jobs has a count of 0.

## Job numbers and client names are shared by all revisions of a code

Lines 32-53

Asked as: same jobs on every revision; job numbers list missing; client names empty on rate sheet; old revision shows different clients

The job numbers (job_numbers) and client names (client_names) are grouped by the rate sheet code, not by the single record. Every revision of a rate sheet with the same code shows the same lists of jobs and clients. The lists contain each value once, and are empty text when no job uses the code. A job without a client adds no client name.

## Rate sheets that belong to a project request are left out

Lines 54-54

Asked as: project request rates missing; proposed rates not showing; rate sheet hidden from view; why is my rate sheet gone

The view includes only rate sheets with no project request (`request` is empty or null). Proposed rates in a project request do not appear in it.

## Rate sheet code format in the view

Lines 63-63

Asked as: rate sheet code wrong format; rate sheet number invalid; code won't save pattern; rate sheet code month

The rate sheet code must match the pattern year-month-five digits, such as 2026-11-00004. The month runs from 01 to 12.

## Effective date format in the view

Lines 65-65

Asked as: effective date rejected; date format error rate sheet; effective date required; wrong date on rate sheet

The effective date (effective_date) is required and must have the form YYYY-MM-DD.
