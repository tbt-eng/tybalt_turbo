<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_api.go","sha256":"0925790ad6a9290ec58d753d68afe2a245a1423c77a7128333308970fbbe47ac","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_api.go

## Job rows show client, alias and invoicing profile flag

Lines 41-58

Asked as: job missing invoicing profile; no invoicing profile flag; job client alias blank; job list outstanding balance

Each job row has the client name and ID, the client alias (blank when none), branch, manager, outstanding balance and its date. It also has a flag (no_invoicing_profile) that is true when the job has no invoicing profile.

## Latest jobs list size per group

Lines 60-123

Asked as: how many proposals listed; latest jobs limit; group name on jobs; latest jobs 50 rows

The latest jobs list returns up to 50 proposals and up to 50 projects (latestJobsLimit). Each row carries a group name (group_name) that labels its group.

## Single job lookup returns not found

Lines 78-100

Asked as: job not found error; job id lookup fails; no job matches id; get job by id

When a job ID is given and no job matches, Turbo answers "job not found". Without an ID, Turbo returns the whole jobs list. This job lookup and list shows no claim check in this file.

## Latest jobs list needs the job claim

Lines 107-114

Asked as: cannot see latest jobs; not authorized latest jobs; latest jobs access denied; who can view latest jobs

Only procurement (the job claim) can view the latest jobs list. Anyone else gets the message "you are not authorized to view latest jobs".

## Unused jobs search defaults to this year's proposals

Lines 130-139

Asked as: unused jobs this year; unused proposals default prefix; search unused jobs P25; job number prefix missing

The unused jobs search takes a job number prefix from the query, or else from the path. When neither is given, it uses P plus the last two digits of the current UTC year, such as P25. Blank or whitespace-only prefixes count as missing.

## Unused jobs search returns 50 rows by default

Lines 141-146

Asked as: unused jobs row limit; how many unused jobs shown; limit unused jobs search; unused jobs result cap

The unused jobs search returns at most 50 jobs unless a limit is given. A limit that is not a whole number above 0 is ignored, and 50 applies.

## Unused jobs are stale jobs with age zero

Lines 148-157

Asked as: jobs with no references; find unused jobs zero age; last reference of unused job; jobs never referenced

The unused jobs search runs the stale jobs query with an age of 0, to find jobs with no references. Each row shows the last reference and its type. This file shows no claim check for this search.

## Stale jobs search defaults to this year's proposals

Lines 164-172

Asked as: stale jobs this year; stale proposals default prefix; search stale jobs P25; stale jobs prefix missing

The stale jobs search takes a job number prefix from the query, or else from the path. When neither is given, it uses P plus the last two digits of the current UTC year.

## Stale jobs default age is 180 days

Lines 174-180

Asked as: stale jobs age default; how many days stale; stale jobs 180 days; change stale age limit

The stale jobs search uses an age of 180 days unless a different age is given. An age that is not a whole number above 0 is ignored, and 180 applies.

## Stale jobs search returns 50 rows by default

Lines 182-198

Asked as: stale jobs row limit; how many stale jobs shown; stale jobs result cap; last reference stale job

The stale jobs search returns at most 50 jobs unless a limit is given. A limit that is not a whole number above 0 is ignored, and 50 applies. Each row shows the last reference and its type. This file shows no claim check for this search.
