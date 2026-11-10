<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expenses.sql","sha256":"947582cb53efc5801efffbabe5dca3e3d023b1d603f79e8bb945b17d34f2623f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expenses.sql

## Missing currency falls back to CAD

Lines 7-10

When an expense has no matching currency, the list shows the currency code 'CAD' and the symbol 'CAD'. The currency id and the currency icon are blank in that case.

## Missing branch, division or person shows blank

Lines 14-23

When an expense has no branch, division or staff profile, the list shows empty text for the branch code, division code, surname and given name. The expense is still listed.

## Missing category shows as 'No Category'

Lines 19-19

When an expense has no matching category, the job expense list shows the category name 'No Category'.

## Job expense list shows only committed expenses

Lines 26-27

The job expense list returns only expenses of the requested job that have been committed (committed is not empty). Expenses that are not yet committed do not appear.

## Optional filters for branch, division, payment type and staff

Lines 28-31

The job expense list can be filtered by branch, division, payment type and staff member (uid). A filter that is missing or empty is ignored and does not narrow the list.

## Category filter has a 'none' value for uncategorized expenses

Lines 32-38

A category filter of 'none' shows only expenses that have no category. Any other category value shows only expenses in that category. A missing or empty category filter shows all expenses.

## Job expense list order and paging

Lines 39-40

Expenses are listed by date, newest first (e.date DESC). The list returns at most the requested limit of rows, starting at the requested offset.
