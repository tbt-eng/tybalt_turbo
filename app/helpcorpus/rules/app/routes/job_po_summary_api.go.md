<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_po_summary_api.go","sha256":"502c9c46878510ca0778d2446aca7635314cf02c0791ff931d70d2ae2f067d82","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_po_summary_api.go

## Summary covers active and under-revision POs only

Lines 28-53

The job purchase order summary is for the active and under-revision purchase orders of one job. Which purchase orders count is decided by job_po_summary.sql, not by this file.

## Job PO summary requires a job id

Lines 35-38

The job purchase order summary refuses a request that has no job id. Turbo shows the message "id is required".

## Job PO summary can be filtered by four options

Lines 40-53

The summary of a job's purchase orders accepts the optional filters division, branch, type (purchase order type) and uid (user id). Turbo passes each filter to the summary query (job_po_summary.sql). The filters are empty when the request omits them. What an empty filter does depends on that SQL file, which this file does not show.

## Empty result returns an empty summary

Lines 54-56

When the summary query returns no row, Turbo returns an empty result with no fields. It does not return zeros or an error.

## Summary query failure shows an error message

Lines 57-57

When the summary query fails for any reason other than finding no row, Turbo shows the message "failed to execute query: " followed by the database error text.

## Missing text values become empty strings

Lines 60-65

When the query gives no value for earliest PO, latest PO, branches, divisions, types or names, Turbo returns an empty string for that field.

## Missing or unreadable total shows 0

Lines 67-76

Turbo returns a total amount of 0 when the query gives no total. It also returns 0 when the total cannot be read as a number, and it gives no error. A missing revision count shows as 0 too.

## Summary fields returned for a job's POs

Lines 74-85

A found summary returns the total amount (total_amount), the revision count (revision_count), the earliest PO (earliest_po), the latest PO (latest_po), the branches, the divisions, the types and the names. The values come from the summary query.
