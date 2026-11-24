<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poVisibility.ts","sha256":"087b6ae52cfe12116298df28195819dd06d9c69c16754ce1fe9f6edad1a5abf8","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poVisibility.ts

## Purchase order lists come in fixed scopes

Lines 4-84

Asked as: filter purchase orders by status; purchase order list scope; show only rejected POs; stale purchase orders list; expiring POs filter

Turbo lists visible purchase orders by scope: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. The stale scope takes a stale_before date and the expiring scope takes an expiring_before date. Other scopes ignore that date. A list can be limited to a number of results.

## Which purchase orders count as revisions

Lines 14-26

Asked as: PO shows as unapproved; edited purchase order marked revision; what counts as a revision; legacy manual entry PO; unapproved PO with number

A purchase order is a revision (under revision) when its status is Unapproved, it already has a PO number (po_number) and it is not a legacy manual entry (legacy_manual_entry). An Unapproved purchase order without a PO number is a new one, not a revision. A legacy manual entry is never a revision.

## Who a revised purchase order waits on

Lines 28-44

Asked as: who is my PO waiting on; revised PO stuck; PO waiting for approval; rejected PO waiting on creator; second approval pending

For a revision, Turbo shows who must act next. If it was rejected, it waits on the creator (or "The creator" when no name) to resubmit or cancel it. If it was approved once, it waits on the second approver (or "A second approver") for second approval. Otherwise it waits on the approver (or "The approver") for first approval.

## Who may cancel or close a purchase order

Lines 106-111

Asked as: who can cancel a PO; close my purchase order; can't cancel PO; PO owner cancel rights; payables admin close PO

Only an Active purchase order can be cancelled or closed. Its owner or payables staff (the payables_admin claim) may do either. The server makes the same checks.

## Cancel is offered only with no committed expenses

Lines 112-116

Asked as: cancel button missing; PO has expenses cannot cancel; cancel greyed out; pending expense blocks cancel; why won't PO cancel

Turbo offers Cancel only when the committed expense count is 0. A purchase order with expenses can't be cancelled. The lists know only the committed count, so the server also refuses a purchase order that has an expense that is not committed.

## Manual close needs a non-One-Time PO with a committed expense

Lines 117-121

Asked as: close button not showing; manually close one-time PO; close PO with no expenses; recurring PO close option; cumulative PO can't close

Turbo offers Close only for a Recurring or Cumulative purchase order (not One-Time) that has at least 1 committed expense. The server checks pending expenses and says what remains.
