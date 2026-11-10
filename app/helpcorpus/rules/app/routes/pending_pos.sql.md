<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/pending_pos.sql","sha256":"29c9a32ff3b37ce691bf4c20b2700445e139ab696897b3c1e4adc040e84d843a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/pending_pos.sql

## Pending POs list starts from the PO visibility rules

Lines 1-3

Asked as: pending purchase orders list; POs waiting on me; list of POs to approve; which POs need my action

The pending purchase orders list starts from the set of POs that the PO visibility base (__PO_VISIBILITY_BASE__) returns. Turbo then keeps only the POs flagged as needing action now. The visibility rules themselves come from another file and are not defined here.

## Pending POs list shows approval, currency and spending details

Lines 4-76

Asked as: PO approver details; PO total in home currency; PO currency rate shown; pending PO amounts and remaining

Each PO in the pending list carries its approver, second approver (second_approver) and priority second approver (priority_second_approver). It also carries the approval total in the PO's currency and in the home currency (approval_total, approval_total_home), the currency rate and its date, the committed expenses count and total, and the remaining amount (remaining_amount). The recurring expected and remaining occurrences and the print maximum amount (print_max_amount) appear too. The list adds display names for people, vendor, job, client, division, category, branch and kind.

## Pending POs list shows only POs that need action now

Lines 77-78

Asked as: why PO not in pending list; PO missing from pending; only POs needing action; pending POs filter

The pending purchase orders list returns only POs whose flag is_unapproved_actionable_now equals 1. The visibility base query (__PO_VISIBILITY_BASE__, defined elsewhere) decides which POs get this flag. This file does not show what the flag means.

## Pending POs list sorts by most recently updated first

Lines 79-79

Asked as: pending POs order; newest PO first; latest updated POs on top; sort pending purchase orders

Turbo sorts the pending purchase orders list by the updated date, newest first.
