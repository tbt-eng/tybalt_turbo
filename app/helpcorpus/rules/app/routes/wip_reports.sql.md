<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/wip_reports.sql","sha256":"13bc02883e310bad0a29acdb33818710cd8d817aff3dca965555805a86302607","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/wip_reports.sql

## WIP report includes only Active jobs that are not proposals

Lines 1-4

Asked as: proposals in WIP report; WIP missing jobs; only active jobs WIP; P jobs excluded

The WIP report selects only jobs with status 'Active'. It leaves out any job whose number starts with P (proposals).

## My WIP report lists only jobs you manage

Lines 5-5

Asked as: my WIP only my jobs; WIP showing other managers jobs; jobs I manage WIP; my WIP report empty

In 'my' mode, the WIP report keeps only jobs whose manager (manager) is the current user (uid). In other modes this filter does not apply.

## Branch WIP report filters by branch unless branch is empty

Lines 6-6

Asked as: branch WIP all branches; WIP blank branch; branch filter not working; WIP every branch

In 'branch' mode, the WIP report keeps only jobs in the chosen branch. If the branch is an empty string, Turbo applies no branch filter and includes jobs from all branches.

## Division WIP report requires a time allocation for the division

Lines 7-8

Asked as: division WIP missing jobs; no time allocation division; WIP division jobs not showing; division WIP filter

In 'division' mode, the WIP report includes a job only if the job has a time allocation (job_time_allocations) for the chosen division.

## Division WIP can require logged time up to the end date

Lines 9-14

Asked as: division WIP needs time logged; WIP require time; job hours before end date; WIP zero hours not counting

In 'division' mode, when require_time is not 0, a job must also have a time entry from the chosen division with hours above 0 and a date on or before the end date (end_date). The end date itself counts. Zero-hour entries do not count. When require_time is 0, Turbo skips this check.

## WIP pricing skips jobs without a project value

Lines 15-17

Asked as: WIP jobs with no value; project value zero WIP; WIP pricing excluded jobs; jobs without project value

Only jobs with a project value (project_value) above 0 go into the selected jobs used for WIP pricing. Jobs with a value of 0 or less are excluded from pricing and counted only in the exclusion count.
