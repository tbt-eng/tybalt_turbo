<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/tallies.sql","sha256":"33de4b39d84456ac779aea51696d3cb447a0e9807180bbf9647f1488de1a7b2b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/tallies.sql

## Missing names, submitted and committer show as defaults

Lines 14-25

When a name is missing, the tallies query shows an empty string for the owner's given name and surname and for the approver, committer and rejector names. An unset committer shows as an empty string, and an unset submitted value shows as 0.

## Shared reviewer count per time sheet

Lines 26-26

Each tallies row shows how many reviewers a time sheet is shared with (shared_reviewer_count). A sheet with no reviewers shows 0.

## Work hours split by whether an entry has a job

Lines 27-29

Work hours count time entries with time type R or RT. Entries with no job go in work_hours, and entries with a job go in work_job_hours. work_total_hours adds both.

## Hours totals per time type

Lines 30-38

The tallies query adds up hours separately for time types OP (op_hours), OV (ov_hours), OH (oh_hours), OB (ob_hours), OS (os_hours) and RB (rb_hours). It also sums meals hours (meals_hours) and payout request amounts (payout_request_amount).

## Non-work total excludes R, RT, OR, OW, OTO and RB

Lines 37-37

non_work_total_hours adds the hours of all entries whose time type is not R, RT, OR, OW, OTO or RB. Off rotation, off week, payout request and bank entries therefore do not count as non-work hours, and neither does regular work.

## Distinct job, division and time type lists

Lines 39-43

Each tallies row lists the distinct job numbers, division codes and names, and time type codes and names used in the sheet. Empty values are left out.

## Date lists for off rotation, off week, payout and bank entries

Lines 44-47

Each tallies row lists the dates of entries with time type OR (off rotation), OW (off week), OTO (payout request) and RB (bank entry). Dates are not de-duplicated.

## Sheets with no time entries are not listed

Lines 48-49

The tallies query starts from time entries and joins them to time sheets with an inner join. A time sheet with no time entries does not appear in the tallies. Time entries without a time sheet are also excluded.

## Tallies list time sheets by owner, approver or shared reviewer

Lines 64-73

The tallies query lists time sheets according to the caller's role. Role 'uid' lists the caller's own time entries' sheets (te.uid). Role 'approver' lists sheets where the caller is the approver (ts.approver). Role 'reviewer' lists sheets where the caller is explicitly listed as a reviewer in time_sheet_reviewers.

## Unknown role returns no time sheets

Lines 64-73

The tallies query matches only the roles 'uid', 'approver' and 'reviewer'. Any other role value matches no time sheet, so the list is empty.

## Pending-only filter shows unapproved sheets

Lines 74-74

When pendingOnly is not 0, the tallies query returns only time sheets whose approved value is empty (ts.approved = ''). When pendingOnly is 0, this filter does not apply.

## Approved-only filter shows approved sheets

Lines 75-75

When approvedOnly is not 0, the tallies query returns only time sheets whose approved value is not empty (ts.approved != ''). When approvedOnly is 0, this filter does not apply.

## One row per time sheet, newest week first

Lines 76-77

The tallies query groups time entries by time sheet to give one row per sheet. It sorts the rows by week ending, newest first.
