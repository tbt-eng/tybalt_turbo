<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_wip_api.go","sha256":"12e71d402f0e854becc67a1ea925b03a32cbcc533ad2cc9278bead0b010fe5ff","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_wip_api.go

## Job WIP report returns its figures and warning counts

Lines 17-34

A successful WIP request returns one row for the job. The row has the project value, time value, hours, estimated hours, expense value and PO value. It also has the as-of date, the missing-role hours, the missing-role rate hours and a no-rate-sheet flag. It has counts or hours of unpriced and estimated items: unpriced hours, unpriced expenses, unpriced POs and estimated POs.

## WIP expense deduction from active POs ignores inclusion controls

Lines 17-18

Inclusion controls never change the expense deduction from active POs. This prevents the same commitment from being counted twice.

## Job WIP report is current, not a historical snapshot

Lines 39-45

The job WIP (work in progress) report is calculated as of today (UTC), not as of an earlier date. Time and expenses count only dates through today. Active POs can cover future work and still count.

## Job WIP report has no start date

Lines 41-45

The job WIP report sets its start date (start_date) to empty. It counts time and expenses from the beginning of the job through today.

## Job WIP report refuses an unknown job

Lines 46-48

When the job id matches no job, Turbo refuses the WIP request with the message "Job not found".

## Job WIP report fails when the figures cannot load

Lines 49-51

When the WIP query fails for any reason other than a missing job, Turbo shows the message "Failed to load WIP".
