<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/close_purchase_order.go","sha256":"fc41753dc89038172e92f52e2d6a2bbfdc750a04a2a9086efd89ed060443e5ac","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/close_purchase_order.go

## Closing a purchase order that does not exist is refused

Lines 22-31

Asked as: close a PO that doesn't exist; purchase order not found; po_not_found error; error fetching purchase order

If the purchase order id is not found, Turbo refuses to close it. The error code is po_not_found and the message starts 'error fetching purchase order:'.

## Who may close a purchase order

Lines 33-51

Asked as: who can close a PO; can't close someone else's purchase order; not authorized to close PO; payables admin close permission

The person who owns the purchase order (its uid) may close it. A holder of the payables admin claim (payables_admin) may also close it. Anyone else gets "you are not authorized to close purchase orders".

## Order of checks when closing a purchase order

Lines 35-104

Asked as: wrong error when closing PO; why does closing show a different message; order of checks closing PO; which error first closing purchase order

Turbo checks in this order: who may close, One-Time type, at least one committed expense, Active status, then pending expenses. It stops at the first check that fails. A non-Active PO with no committed expense gets the 'only cumulative or recurring purchase orders with at least one associated expense may be closed manually' message, not 'only active purchase orders can be closed'.

## One-Time purchase orders cannot be closed manually

Lines 53-61

Asked as: cannot close one-time PO; close a one-time order; one-time purchase order won't close; cancel one-time PO instead

Turbo refuses to close a One-Time purchase order and says "One-Time purchase orders may be cancelled but not manually closed". The same refusal applies to the owner and to payables admins.

## Closing needs at least one committed expense

Lines 63-80

Asked as: no committed expense to close PO; close PO with no expenses; purchase order has no committed expenses; cancel the purchase order instead

A purchase order can be closed only if at least one of its expenses is committed. Otherwise Turbo says "only cumulative or recurring purchase orders with at least one associated expense may be closed manually. Cancel the purchase order instead."

## Only Active purchase orders can be closed

Lines 82-89

Asked as: close PO not active; closing a draft or pending PO; only active purchase orders closable; can't close PO status

Turbo refuses to close a purchase order whose status is not Active and says "only active purchase orders can be closed". This check runs after the type check and the committed-expense check.

## Pending expenses block manual closure

Lines 91-104

Asked as: PO won't close pending expenses; expense stopping me closing PO; waiting for approval blocks closing; close button greyed out expense

Turbo refuses to close a purchase order while any of its expenses is submitted but neither committed nor rejected. Drafts, recalled expenses and rejected expenses do not block closure. Approved expenses that wait for commitment block it too.

## Closing a purchase order records closer, time and Closed status

Lines 106-120

Asked as: what happens when PO is closed; record who closed purchase order; PO closed date and closer saved; status changed to closed

When all checks pass, Turbo sets the purchase order's status to Closed (status). It also stores the current time (closed) and the person who closed it (closer). The change is saved in one transaction.

## Successful closure returns the updated purchase order

Lines 137-142

Asked as: closed PO shows wrong details; view purchase order after closing; confirm PO closed; returned closed purchase order record

When a manual close succeeds, Turbo reads the purchase order again and returns the saved record. The record shows status Closed, the closed time and the closer.

## Pending expense message names approvers and approved expenses awaiting commitment

Lines 150-196

Asked as: why can't I close PO awaiting approval; expenses awaiting approval message; approved expense waiting commitment; which approver is holding PO

Closing counts expenses on the PO that are submitted, not committed and not rejected. Unapproved ones give 'This PO has one or more expenses awaiting approval by <approver names>.' Approved ones give 'This PO has one or more approved expenses awaiting commitment.' When both exist, both sentences appear in that order, joined by a space.

## Pending-expense approver names are sorted and listed once

Lines 156-169

Asked as: list of approvers blocking PO; duplicate approver names closing; who must approve before closing; approver names order message

The message lists each approver once per approval state. Names are sorted by approver name, then by approver id.

## Draft, recalled and rejected expenses do not block closing

Lines 158-169

Asked as: draft expense blocks closing; rejected expense stops PO closing; recalled expense closing PO; approved expense stops close

Only expenses that are submitted, not committed and not rejected block manual closure. Drafts, recalled expenses and rejected expenses do not block it. Approved expenses that are not yet committed do block it.

## Approver with no name still blocks closing

Lines 174-188

Asked as: approver name missing closing PO; blank approver name blocks close; unknown approver blocking PO; approver profile has no name

If an approver's profile name is missing or blank, the expense still blocks closure. The message lists 'an approver whose name is unavailable' in place of the name, after the named approvers.
