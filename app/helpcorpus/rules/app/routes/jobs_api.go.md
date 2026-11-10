<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_api.go","sha256":"25c427258b8aeef3bb394ecf28b7044d9bda8ac941522d1553d2210a94e4277a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_api.go

## Latest jobs returns 50 proposals and 50 projects

Lines 56-57

Asked as: latest jobs how many; only 50 proposals; jobs list limit; group name on jobs

The latest jobs list returns at most 50 proposals and at most 50 projects (latestJobsLimit = 50). Each row carries a group name (group_name) that says which group it belongs to.

## Single job lookup says job not found

Lines 74-96

Asked as: job not found; can't find job by id; job id lookup empty; search job returns nothing

When a job id is given and the query finds no row, Turbo answers "job not found". Without an id, Turbo returns the whole list of jobs from the jobs query. This file does not check any claim for either case.

## Latest jobs list requires the job claim

Lines 103-110

Asked as: cannot see latest jobs; not authorized latest jobs; job claim missing; jobs page access denied

Turbo shows the latest jobs only to staff who hold the job claim. Anyone else gets the message "you are not authorized to view latest jobs".

## Latest jobs query uses the fixed limit

Lines 112-119

Asked as: change jobs list size; raise latest jobs limit; jobs limit fixed; more than 50 latest jobs

The latest jobs endpoint passes the fixed limit of 50 (latestJobsLimit) to the jobs_latest.sql query. A caller cannot change it.

## Unused jobs default to current year's proposal prefix

Lines 126-135

Asked as: unused jobs prefix; unused jobs wrong year; P25 prefix default; unused jobs filter by year

The unused jobs list filters by job number prefix. It takes the prefix from the query parameter, or else from the path. If both are empty or blank, it uses P plus the two-digit current UTC year, such as P25. Turbo trims spaces from the prefix.

## Unused jobs default limit is 50

Lines 137-142

Asked as: unused jobs only 50; change unused jobs count; unused jobs limit parameter; show more unused jobs

The unused jobs list returns 50 jobs by default. A positive whole-number limit in the query overrides this. A limit that is zero, negative or not a number is ignored and the default of 50 stays.

## Unused jobs are stale jobs with age 0

Lines 144-153

Asked as: jobs with no references; unused jobs meaning; jobs never referenced; which jobs are unused

Turbo finds unused jobs (jobs with no references) by running the stale-jobs query with age 0. The result includes the last reference and its type (last_reference, last_reference_type) for each job. This file checks no claim for this list.

## Stale jobs default prefix is current year's proposal prefix

Lines 160-168

Asked as: stale jobs prefix; stale jobs wrong year; stale jobs current year; stale jobs filter prefix

The stale jobs list filters by job number prefix. It takes the prefix from the query parameter, or else from the path. If both are empty, it uses P plus the two-digit current UTC year, such as P25.

## Stale jobs default age is 180 days

Lines 170-176

Asked as: stale jobs age; days before job stale; change stale days; stale after 180 days

The stale jobs list uses an age of 180 days by default. A positive whole-number age in the query overrides it. An age that is zero, negative or not a number is ignored and 180 stays.

## Stale jobs default limit is 50

Lines 178-183

Asked as: stale jobs only 50; stale jobs count limit; how many stale jobs shown; stale jobs list size

The stale jobs list returns 50 jobs by default. A positive whole-number limit in the query overrides it. An invalid or non-positive limit is ignored.

## Stale jobs list shows last reference

Lines 185-194

Asked as: last reference on job; stale job last used; job last reference type; when job was last referenced

Each stale job row includes its last reference and the type of that reference (last_reference, last_reference_type). The prefix, age and limit go to the jobs_stale.sql query. This file checks no claim for this list.
