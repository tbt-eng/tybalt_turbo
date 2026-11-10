<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_expense_summary.sql","sha256":"ce7dda0a4320f0cbca655c13d37df628c1ee8c827cd5ee47fd133505eee588a8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_expense_summary.sql

## Expense total uses the settled total when it is not zero

Lines 4-4

The summary's total amount (total_amount) adds up each expense's settled total (settled_total). When the settled total is 0, Turbo uses the expense's total (total) instead.

## Summary shows earliest and latest expense dates

Lines 5-6

The summary gives the earliest (earliest_expense) and latest (latest_expense) expense date among the matching expenses.

## Summary lists the distinct branches, divisions, payment types and staff

Lines 7-26

The summary lists each distinct branch (id and code), division (id and code), payment type and person (uid and full name, given name then surname) found in the matching expenses. A branch, division or person missing from its table (LEFT JOIN) still produces an entry with empty values.

## Expenses without a category show as No Category

Lines 19-21

In the categories list, an expense with no category appears with the id 'none' and the name 'No Category'.

## Job expense summary counts only committed expenses

Lines 27-28

The job expense summary includes only expenses of the given job that have been committed (committed is not empty). Expenses that are not yet committed are left out of the totals, dates and lists.

## Optional filters by branch, division, payment type and person

Lines 29-32

The summary can be narrowed by branch, division, payment type and person (uid). A filter that is missing or empty is ignored. Otherwise the expense must match the value exactly.

## Category filter supports 'none' for uncategorized expenses

Lines 33-38

A missing or empty category filter is ignored. The filter value 'none' keeps only expenses with no category (empty category). Any other value keeps only expenses in that category.
