<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poVisibility.ts","sha256":"aaa1bdf4e31e234e540e7b7f93ce4491258722b1e949f0467885b24b94e79c6e","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poVisibility.ts

## What counts as a purchase order revision

Lines 14-26

Turbo treats a purchase order as a revision (under revision) when its status is Unapproved, it already has a PO number (po_number), and it is not a legacy manual entry (legacy_manual_entry). An Unapproved PO with no PO number is a new PO, not a revision. A legacy manual entry is never a revision.

## Who a revised purchase order waits on after rejection

Lines 34-36

When a purchase order revision has been rejected (rejected is set), Turbo says it waits on the creator (uid_name), or "The creator" if no name is known. The reason shown is "To resubmit or cancel it".

## Who a revision waits on after first approval

Lines 37-42

When a revision is not rejected but already has a first approval (approved is set), Turbo says it waits on the priority second approver (priority_second_approver_name). If that name is missing it shows "A second approver". The reason is "Second approval".

## Who a revision waits on before any approval

Lines 43-43

When a revision is neither rejected nor approved, Turbo says it waits on the approver (approver_name), or "The approver" if no name is known. The reason is "First approval". The rejected check runs before the approved check.

## Purchase order list scopes and date filters

Lines 4-84

Turbo can list visible purchase orders by scope: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. The stale cutoff date (stale_before) applies only to the stale scope. The expiring cutoff date (expiring_before) applies only to the expiring scope. A date given with any other scope is ignored. An optional limit caps the number of results.

## Revisable flag is set only for a single purchase order

Lines 61-62

The revisable field on a visible purchase order is set only when one purchase order is fetched by its id. Lists of purchase orders do not carry it.

## Purchase order group lists use visible purchase orders only

Lines 86-92

Fetching a purchase order group returns only the purchase orders in that group that the person is allowed to see.
