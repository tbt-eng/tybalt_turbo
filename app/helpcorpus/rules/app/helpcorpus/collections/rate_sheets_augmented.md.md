<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"5ef7d9ed7b71b485090ddb8a17aad44e7051699862bf01bd651f6e6f63f2f147","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Nobody can create, update or delete through the rate sheets view

Lines 5-5

Asked as: cannot edit rate sheet; rate sheets read only; can't update rate sheet api; delete rate sheet view

The rate sheets view (rate_sheets_augmented) has no create, update or delete rules, so nobody can do these actions on it through the API. It is a saved query that Turbo reads and never writes.

## Who can list rate sheets with job counts

Lines 7-11

Asked as: list rate sheets; who can see rate sheets; rate sheet list access; signed out rate sheets refused

Any signed-in person can list the records of the rate sheets view (rate_sheets_augmented). A request with no signed-in person is refused.

## Who can view one rate sheet with its job count

Lines 13-17

Asked as: view one rate sheet; open a rate sheet; single rate sheet access; rate sheet details permission

Any signed-in person can view one record of the rate sheets view (rate_sheets_augmented). A request with no signed-in person is refused.

## Rate sheet fields shown in the view and their patterns

Lines 22-52

Asked as: rate sheet date format; effective date format; rate sheet id characters; rate sheet fields shown

Each record shows the rate sheet's name, effective date, revision, active flag, created and updated times, and job_count. The effective_date must be text in the form YYYY-MM-DD, and the id must contain only lowercase letters and digits.

## Job count counts each job once, from current and past rate sheet use

Lines 30-36

Asked as: job count wrong; jobs counted twice rate sheet; how job count works; rate sheet job total

job_count is the number of distinct jobs that use the rate sheet now (the job's rate_sheet field) or used it before (the job_rate_sheets rows, the rate sheet periods). The same job counts once, even if it appears in both places. A job with no rate sheet set is not counted from the jobs table.

## Rate sheets used by no job still appear with a count of 0

Lines 30-36

Asked as: unused rate sheet missing; rate sheet zero jobs; rate sheet no jobs showing; empty rate sheet not listed

The view joins jobs to rate sheets with a LEFT JOIN, so a rate sheet that no job uses still appears, and its job_count is 0.

## Rate sheets view leaves out rate sheets that belong to a project request

Lines 37-37

Asked as: project request rate sheet missing; proposed rates not showing; rate sheet hidden from view; rate sheet attached to request

The view lists only rate sheets whose request field (rs.request) is empty or not set. A rate sheet that is attached to a project request (for example proposed rates) does not appear in it.
