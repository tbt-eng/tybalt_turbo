<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/user_view_grants.md","sha256":"a661355c046fbdb66762e926a75e0237266c20fe4a537cf854bb8594b15ed3f6","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/user_view_grants.md

## Nobody can list user view grants through the API

Lines 7-9

Asked as: list user view grants; see all view grants; API cannot list grants; no access to grants list

The user_view_grants view has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can view one user view grant through the API

Lines 11-13

Asked as: open one view grant; view single grant record; API grant record not visible; cannot read grant by id

The user_view_grants view has no view rule. Nobody can open a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Staff on an active purchase order become visible to everyone

Lines 24-34

Asked as: purchase order makes staff visible; active PO visible to everyone; why can everyone see my profile; PO requester visibility grant

A user who is the requester (uid), the approver or the second approver on a purchase order with status 'Active' gets a grant with reason 'active_po'. The grant has an empty viewer (viewer), so it is not tied to one person. Purchase orders with another status give no grant. A blank uid, approver or second approver gives no grant.

## Expense approver may view the staff member's profile

Lines 35-42

Asked as: expense approver sees my profile; approver can view staff profile; committed expense approver access; who can see my profile after expense

When a person has an expense with a commit (committed is not empty) and an approver, the approver gets a grant (reason 'expense_approver') to view that person. The grant has one row for each pair of person and approver. Expenses that are not committed, or that have no person or no approver, give no grant.

## Grants only exist for users that exist

Lines 43-45

Asked as: grant missing for user; deleted user still has grant; grant not showing up; grant without user record

The view joins each grant to the users table on the uid. A grant for a uid that is not a user record does not appear. The viewer is read from users by a left join, so the 'active_po' grants keep an empty viewer.
