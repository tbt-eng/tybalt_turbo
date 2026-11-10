<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/stats.sql","sha256":"17ab7ec4c5a9f8f9e323c9bb2395c6331617052a5206d601936653c629114f21","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/stats.sql

## Which purchase orders count in the statistics

Lines 1-4

The statistics count only purchase orders with status Active or Closed. Purchase orders with any other status are left out. The count is returned as qualifying_po_count (qualifying_pos).

## Imported purchase orders are excluded from statistics

Lines 5-5

The statistics skip a purchase order marked as imported (_imported = true). A purchase order with _imported empty or false is counted if its status qualifies.

## Which expenses count in the statistics

Lines 7-13

The statistics count an expense only if it is linked to a qualifying purchase order (Active or Closed, not imported). The expense must have an approved value (approved is not empty) and no rejected value (rejected is empty). Expenses with no purchase order are not counted. The count is returned as approved_expense_count.

## How distinct users are counted in the statistics

Lines 14-29

The distinct user count includes each person once, however many roles they have. It counts the creator (uid), approver and second approver of qualifying purchase orders. It also counts the creator (uid), approver and committer of approved expenses. Empty values are ignored. The result is returned as distinct_user_count.

## Statistics query returns three totals

Lines 30-33

The query returns one row with qualifying_po_count, approved_expense_count and distinct_user_count.
