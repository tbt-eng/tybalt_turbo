<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/stats_api.go","sha256":"a0c82bb7d0620c51fc799be9e7f7edc292cb2f8e32d35d3f09052da529fd82ff","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/stats_api.go

## Stats require a signed-in user

Lines 22-25

Asked as: not signed in; stats unauthorized error; login required for stats; session expired stats; can't see stats logged out

Turbo refuses the stats request when the user is not signed in. The message is "unauthorized".

## Only admins can view stats

Lines 27-33

Asked as: stats permission denied; not an admin stats; admin claim required; can't open stats as user; error checking admin claim

Turbo checks that the user holds the admin claim (admin) before it returns stats. A user without that claim gets the message "admin claim required". If Turbo cannot check the claim, it returns the message "error checking admin claim".

## Stats show three counts from the stats query

Lines 35-44

Asked as: purchase order count stats; approved expenses count; distinct users stats; stats query failed; only first row of stats

An admin gets one result with three counts: qualifying purchase orders (qualifying_po_count), approved expenses (approved_expense_count) and distinct users (distinct_user_count). The counts come from the stats.sql query, and this file does not show how it defines each one. Turbo returns only the first row of the query result. If the query fails, the message is "failed to execute stats query: " followed by the error.

## Empty stats result returns zero counts

Lines 40-42

Asked as: stats showing zeros; no data in stats; empty stats page; stats all zero; blank stats counts

If the stats query returns no rows, Turbo returns all three counts as 0 and does not show an error.
