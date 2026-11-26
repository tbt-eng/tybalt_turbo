<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_api.go","sha256":"e227d6a16183acc6b365ba1dfdac4489b4ece3c55fbf8ede866952cfe42d6ce6","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_api.go

## Jobs list flags jobs that have no invoicing profile

Lines 54-56

Asked as: job without invoicing profile; missing invoicing profile flag; job has no invoice profile; jobs list invoicing warning

Each job in the jobs list carries a flag (no_invoicing_profile) that is true when the job has no invoicing profile. Only the jobs list query (jobs.sql) reads this flag. The SQL file decides how Turbo sets it, and this file does not show that.

## Latest jobs list is capped at 50

Lines 59-122

Asked as: latest jobs limit; why only 50 jobs; latest proposals and projects; group name on jobs

The latest jobs list returns up to 50 proposals and up to 50 projects (latestJobsLimit = 50). Each row carries a group name (group_name) that says which group it belongs to.

## Single job lookup returns 'job not found' when missing

Lines 77-99

Asked as: job not found error; cannot find job by id; job id returns nothing; single job lookup failing

When a job id is given, Turbo returns that one job. If no job matches the id, Turbo refuses with the message "job not found". Without an id, Turbo returns the whole list of jobs.

## Latest jobs list needs the job claim

Lines 106-113

Asked as: cannot view latest jobs; not authorized latest jobs; latest jobs permission; job claim missing

Only staff with the job claim can view the latest jobs. Anyone else gets the message "you are not authorized to view latest jobs".

## Unused jobs list defaults to this year's proposal prefix

Lines 129-138

Asked as: default job prefix; unused jobs wrong year prefix; P25 prefix; jobs prefix empty

The unused jobs list takes the job number prefix from the query parameter prefix. If it is empty, Turbo uses the prefix in the path. If that is also empty, Turbo uses P plus the last two digits of the current UTC year, such as P25. Turbo trims spaces from the prefix.

## Unused jobs list limit defaults to 50

Lines 140-145

Asked as: unused jobs only 50 results; change unused jobs limit; limit parameter ignored; unused jobs count cap

The unused jobs list returns 50 jobs unless the limit query parameter is a whole number above 0. A limit that is zero, negative or not a number is ignored and 50 applies.

## Unused jobs are jobs with no references

Lines 147-157

Asked as: jobs with no references; which jobs are unused; last reference of job; jobs never referenced

Turbo finds unused jobs by running the stale jobs query with an age of 0, so it returns jobs that have no references. Each row includes the last reference and its type (last_reference, last_reference_type). Turbo does not check a claim for this list in this file.

## Stale jobs list defaults to this year's proposal prefix

Lines 163-171

Asked as: stale jobs wrong prefix; old jobs list year; stale jobs default prefix; prefix for stale jobs

The stale jobs list takes the job number prefix from the query parameter prefix, then from the path. If both are empty, it uses P plus the last two digits of the current UTC year, such as P25. Turbo trims spaces from the prefix.

## Stale jobs age defaults to 180 days

Lines 173-179

Asked as: stale jobs older than 180 days; change stale age; how old job is stale; stale days setting

The stale jobs list uses an age of 180 days unless the age query parameter is a whole number above 0. A zero, negative or non-numeric age is ignored and 180 applies.

## Stale jobs list limit defaults to 50

Lines 181-197

Asked as: stale jobs limit 50; stale jobs result cap; how many stale jobs shown; stale list size

The stale jobs list returns 50 jobs unless the limit query parameter is a whole number above 0. An invalid limit is ignored and 50 applies. Turbo does not check a claim for this list in this file.
