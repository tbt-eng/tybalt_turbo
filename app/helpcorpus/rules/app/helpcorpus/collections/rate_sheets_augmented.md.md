<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"72b3aeda3c71e5a6d722c52b229be977bf913875343cd2d0e6f3fdf9bc800d75","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Nobody can create, update or delete augmented rate sheets through the API

Lines 5-5

Asked as: cannot edit rate sheet via API; change augmented rate sheet; delete rate sheet record API; read-only rate sheets view

The rate sheets augmented view (rate_sheets_augmented) is read-only. It has no create, update or delete rules, so nobody can change its records through the API. Turbo changes rate sheets only through the rate_sheets collection.

## Who can list rate sheets in the augmented view

Lines 7-11

Asked as: list rate sheets empty; rate sheets list no results; who can see rate sheets; signed in list rate sheets

Any signed-in person can list records in the rate sheets augmented view (rate_sheets_augmented). A request without a signed-in user gets nothing.

## Who can view one rate sheet in the augmented view

Lines 13-17

Asked as: view single rate sheet refused; open one rate sheet forbidden; rate sheet detail access denied; logged out rate sheet view

Any signed-in person can view one record in the rate sheets augmented view (rate_sheets_augmented). A request without a signed-in user is refused.

## Job count includes jobs on current and historical rate sheets

Lines 31-39

Asked as: job count wrong; how many jobs on rate sheet; rate sheet shows 0 jobs; job count double counted

The job count (job_count) counts each job that uses the rate sheet in one of two ways: the job's current rate sheet (jobs.rate_sheet, when it is not empty) or a rate sheet period row (job_rate_sheets). A job that has both links to the same rate sheet counts once, because the query uses UNION. A rate sheet with no jobs shows 0.

## Job count is for one revision but the lists cover the whole code

Lines 31-65

Asked as: job count differs from job list; jobs count mismatch revisions; count and job numbers disagree; old revision jobs listed

The job count (job_count) counts the jobs of this rate sheet record only. The job numbers and client names lists match on the rate sheet code, so they can include jobs of other revisions with the same code. The count and the lists can therefore disagree.

## Job numbers and client names are lists by rate sheet code

Lines 32-65

Asked as: job numbers list blank; client names not showing; which clients use rate sheet; jobs listed for other revision

The job numbers (job_numbers) and client names (client_names) fields hold one item per line, separated by a line feed. They are built by rate sheet code (rs.code), so they include jobs and clients that use any revision with the same code. A rate sheet with no jobs or clients shows an empty string.

## Rate sheets that belong to a project request are hidden

Lines 66-66

Asked as: missing rate sheet from list; project request rate sheet hidden; rate sheet not showing up; proposed rates not in view

The view lists only rate sheets that have no project request (request is empty or null). A rate sheet that proposed rates in a project request (rate_sheets.request) does not appear in this view.

## Rate sheet code format in the view

Lines 75-77

Asked as: rate sheet code format; invalid rate sheet number; effective date format wrong; rate sheet code pattern

The code of a rate sheet must look like 2026-11-00004: four digits for the year, a month from 01 to 12, then five digits. The effective date must look like YYYY-MM-DD.
