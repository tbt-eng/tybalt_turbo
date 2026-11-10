<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/stats.sql","sha256":"17ab7ec4c5a9f8f9e323c9bb2395c6331617052a5206d601936653c629114f21","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/stats.sql

## Which purchase orders count in the statistics

Lines 1-4

Asked as: which purchase orders are counted; PO statistics status filter; active or closed POs; qualifying PO count; purchase orders left out of stats

The statistics count only purchase orders with status Active or Closed. Purchase orders with any other status are left out. The count is returned as qualifying_po_count (qualifying_pos).

## Imported purchase orders are excluded from statistics

Lines 5-5

Asked as: imported PO not counted; imported purchase order stats; why is my import missing; _imported flag statistics; imported POs excluded

The statistics skip a purchase order marked as imported (_imported = true). A purchase order with _imported empty or false is counted if its status qualifies.

## Which expenses count in the statistics

Lines 7-13

Asked as: expense not in statistics; approved expense count; rejected expense excluded; expense without purchase order stats; why expense not counted

The statistics count an expense only if it is linked to a qualifying purchase order (Active or Closed, not imported). The expense must have an approved value (approved is not empty) and no rejected value (rejected is empty). Expenses with no purchase order are not counted. The count is returned as approved_expense_count.

## How distinct users are counted in the statistics

Lines 14-29

Asked as: distinct user count; how many users counted; same person counted twice; unique users in stats; approver counted as user

The distinct user count includes each person once, however many roles they have. It counts the creator (uid), approver and second approver of qualifying purchase orders. It also counts the creator (uid), approver and committer of approved expenses. Empty values are ignored. The result is returned as distinct_user_count.

## Statistics query returns three totals

Lines 30-33

Asked as: statistics totals row; stats query output columns; three totals returned; stats summary numbers; what does stats query return

The query returns one row with qualifying_po_count, approved_expense_count and distinct_user_count.
