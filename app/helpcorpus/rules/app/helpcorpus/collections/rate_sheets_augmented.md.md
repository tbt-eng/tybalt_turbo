<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"aff577fb4a1ca985709614e7539ea417c28c12748a41b119b456b4fe9b2d379e","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Nobody can create, update or delete rate sheets in this view

Lines 5-17

Asked as: cannot edit rate sheet; read only rate sheets view; cannot delete rate sheet via API; rate sheets view is locked

The rate sheets view (rate_sheets_augmented) has only a list rule and a view rule. Nobody can create, update or delete its records through the API, because it is a read-only view. Turbo's own server code can still change the underlying rate sheets.

## Who can list rate sheets in the augmented view

Lines 9-11

Asked as: list rate sheets without login; who can see rate sheet list; rate sheets not loading signed out; anyone can browse rate sheets

Any signed-in person can list the records of the rate sheets view (rate_sheets_augmented). A request without a signed-in user cannot list them.

## Who can view one rate sheet in the augmented view

Lines 15-17

Asked as: open one rate sheet; view single rate sheet permission; rate sheet details not showing; signed out rate sheet access

Any signed-in person can view one record of the rate sheets view (rate_sheets_augmented). A request without a signed-in user cannot view it.

## Job count counts each job once, from current and past use

Lines 31-37

Asked as: how many jobs use rate sheet; job count wrong; jobs counted twice rate sheet; rate sheet usage total

The job count (job_count) is the number of distinct jobs that use the rate sheet. It counts a job whose current rate sheet (jobs.rate_sheet) is this one, and a job that has this rate sheet in its history (job_rate_sheets). A job that appears in both is counted once, because the query uses UNION. A job with no rate sheet is not counted.

## Unused rate sheets still appear with a job count of 0

Lines 31-37

Asked as: unused rate sheet showing; rate sheet with zero jobs; rate sheet missing from list no jobs; why is job count zero

The view joins jobs with a LEFT JOIN. A rate sheet that no job uses still appears, and its job count (job_count) is 0.

## Rate sheets that belong to a project request are hidden

Lines 38-38

Asked as: project proposed rates not showing; rate sheet missing from list; project request rates hidden; proposal rate sheet not visible

The view lists only rate sheets whose request field (rs.request) is empty or null. A rate sheet that a project request owns (proposed rates) does not appear in it.

## Rate sheet code format

Lines 47-47

Asked as: rate sheet code invalid; wrong rate sheet number format; rate sheet code error; code not accepted rate sheet

The code of a rate sheet is required and must match year-month-number: four digits, a month from 01 to 12, and five digits, such as 2026-11-00004.

## Effective date format of a rate sheet

Lines 49-49

Asked as: rate sheet date error; effective date rejected; wrong date format rate sheet; start date of rate sheet invalid

The effective date (effective_date) of a rate sheet is required and has the pattern YYYY-MM-DD (four digits, two digits, two digits).
