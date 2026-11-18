<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_jobs.go","sha256":"2909f1d6aaa8b906daa4625751f0313d43e7c737f5e8acd8cd21b90740fa0048","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_jobs.go

## Each listed job shows whether it uses the revision now

Lines 14-24

Asked as: job current flag; which jobs use this rate sheet; job still on revision; current job marker

Each job in the list shows its ID, number, description, client and a current flag (current). The flag is true when the job uses the revision now. It is false when only a period of the job's rate sheet history used the revision.

## Rate sheet jobs list covers one exact revision

Lines 26-41

Asked as: jobs for one revision; rate sheet revision job list; other revisions jobs; jobs sorted by client

The rate sheet jobs list returns the jobs of the one revision of a rate sheet that is requested by its record ID. Other revisions of the same rate sheet code have their own lists. Turbo sorts the jobs by client name and job number.

## Any signed-in user can read a rate sheet's jobs

Lines 26-30

Asked as: who can view rate sheet jobs; permission to see rate sheet jobs; signed in user access jobs; no claim needed to read

Any signed-in user can read the jobs of a rate sheet revision, because they can already list jobs and clients. No special claim is needed.

## Unknown rate sheet gives 'rate sheet not found'

Lines 33-35

Asked as: rate sheet not found; missing rate sheet error; wrong rate sheet id; rate sheet doesn't exist

When no rate sheet record has the requested ID, Turbo refuses the request with the message "rate sheet not found". It does not return an empty list in this case.

## Revision with no jobs returns an empty list

Lines 36-40

Asked as: revision has no jobs; empty jobs list; failed to load jobs; no jobs showing for rate sheet

When no job uses or used the revision, Turbo returns an empty list and not an error. If the jobs fail to load, Turbo shows "failed to load jobs".
