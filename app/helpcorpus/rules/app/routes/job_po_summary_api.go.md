<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_po_summary_api.go","sha256":"502c9c46878510ca0778d2446aca7635314cf02c0791ff931d70d2ae2f067d82","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_po_summary_api.go

## Summary covers active and under-revision POs only

Lines 28-53

Asked as: active purchase orders only; which POs are counted; under revision POs; job PO totals scope

The job purchase order summary is for the active and under-revision purchase orders of one job. Which purchase orders count is decided by job_po_summary.sql, not by this file.

## Job PO summary requires a job id

Lines 35-38

Asked as: id is required; job PO summary no job; missing job id error; summary without job

The job purchase order summary refuses a request that has no job id. Turbo shows the message "id is required".

## Job PO summary can be filtered by four options

Lines 40-53

Asked as: filter job POs by division; filter PO summary by branch; purchase order type filter; summary filter by user

The summary of a job's purchase orders accepts the optional filters division, branch, type (purchase order type) and uid (user id). Turbo passes each filter to the summary query (job_po_summary.sql). The filters are empty when the request omits them. What an empty filter does depends on that SQL file, which this file does not show.

## Empty result returns an empty summary

Lines 54-56

Asked as: empty PO summary result; job has no purchase orders; summary returns nothing; blank job PO summary

When the summary query returns no row, Turbo returns an empty result with no fields. It does not return zeros or an error.

## Summary query failure shows an error message

Lines 57-57

Asked as: failed to execute query; PO summary database error; job PO summary crashed; summary query error message

When the summary query fails for any reason other than finding no row, Turbo shows the message "failed to execute query: " followed by the database error text.

## Missing text values become empty strings

Lines 60-65

Asked as: blank earliest PO; latest PO missing; no branch names shown; empty text in PO summary

When the query gives no value for earliest PO, latest PO, branches, divisions, types or names, Turbo returns an empty string for that field.

## Missing or unreadable total shows 0

Lines 67-76

Asked as: total amount shows zero; PO total is 0; wrong total on job; revision count zero

Turbo returns a total amount of 0 when the query gives no total. It also returns 0 when the total cannot be read as a number, and it gives no error. A missing revision count shows as 0 too.

## Summary fields returned for a job's POs

Lines 74-85

Asked as: what fields job PO summary returns; job PO total and latest PO; branches and divisions for job; summary of job purchase orders

A found summary returns the total amount (total_amount), the revision count (revision_count), the earliest PO (earliest_po), the latest PO (latest_po), the branches, the divisions, the types and the names. The values come from the summary query.
