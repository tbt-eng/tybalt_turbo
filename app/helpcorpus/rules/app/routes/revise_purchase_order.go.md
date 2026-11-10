<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/revise_purchase_order.go","sha256":"7a78ed620224d25581c6a6f06eb1ed8a04626ed7a9a8f7d5139ae997a0f5c59e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/revise_purchase_order.go

## Revision keeps the purchase order number

Lines 1-5

Asked as: po number changes after edit; same po number after revising; revised order new number; purchase order number lost

A revised purchase order goes back to Unapproved and keeps its po_number. Approving it again re-activates the same number. The change is saved in one transaction, so a failed revision changes nothing.

## Revision request lets the creator submit new approvers

Lines 24-31

Asked as: change approver on existing po; add new approver when editing; update approvers with new total; swap approver on purchase order

The revision request carries the new total, an optional currency, an optional approver and an optional second approver (priority_second_approver). A currency or approver that is left out keeps its current value. The creator can send new approvers with the new total because the new total may need different approvers.

## Revision with an unreadable request body is refused

Lines 38-44

Asked as: invalid request body error; revise po request failing; revision fails bad data

When the request body cannot be read, Turbo refuses the revision with the message 'invalid request body' (code invalid_request_body).

## Revising a purchase order that does not exist is refused

Lines 49-56

Asked as: purchase order not found; can't find po to revise; error fetching purchase order

When the purchase order cannot be found by its id, Turbo refuses the revision with code po_not_found. The message is 'error fetching purchase order:' followed by the lookup error.

## Only active purchase orders can be revised

Lines 58-64

Asked as: cannot edit draft po; revise po not active; edit pending purchase order

Turbo refuses to revise a purchase order whose status (status) is not Active. The message is "only active purchase orders can be revised".

## Revision refusal messages for each blocking condition

Lines 58-121

Asked as: why can't I revise this; revision refused message; revise button error; po edit blocked reason

Turbo refuses a revision with these messages: 'only active purchase orders can be revised', 'only the creator can revise this purchase order', 'legacy purchase orders cannot be revised', 'this purchase order has associated expenses and cannot be revised' and 'this purchase order has child purchase orders and cannot be revised'.

## Only the creator can revise a purchase order

Lines 66-72

Asked as: edit someone else's purchase order; not my po can't change; edit po I didn't create

Turbo refuses a revision from anyone other than the creator of the purchase order (uid). The message is "only the creator can revise this purchase order". The status check runs before this check.

## Legacy and imported purchase orders cannot be revised

Lines 74-83

Asked as: old imported po edit; manual entry purchase order change; legacy po can't update

Turbo refuses to revise a purchase order that is a legacy manual entry (legacy_manual_entry) or imported (_imported). The message is "legacy purchase orders cannot be revised". These orders never went through Turbo approval, so they cannot return to it.

## Purchase order with expenses cannot be revised

Lines 85-103

Asked as: po has expenses can't edit; expense attached to po blocks change; linked expenses stop revision

Turbo refuses to revise a purchase order when any expense uses it. The message is "this purchase order has associated expenses and cannot be revised".

## Purchase order with child orders cannot be revised

Lines 105-121

Asked as: child po blocks edit; sub purchase order stops revision; parent po can't change; cancelled child still blocking

Turbo refuses to revise a purchase order that has a child purchase order (parent_po) whose status is not Cancelled. The message is "this purchase order has child purchase orders and cannot be revised". Cancelled children do not block a revision.

## A revision must change the total or currency

Lines 123-134

Asked as: same total nothing changed; save revision without change; revise with no changes

Turbo refuses a revision when the new total equals the current total and the effective currency is unchanged. The message is "a revision must change the total or the currency". A changed total or a changed currency is enough. The currency is compared by its effective currency code, after trimming spaces from the submitted value.

## Omitted currency and approvers keep their current values

Lines 123-144

Asked as: blank currency on revision; leave approver empty when revising; fields left empty revision

When the request leaves out the currency, the purchase order keeps its current currency. When it leaves out the approver or the priority second approver, each keeps its current value. A submitted value is trimmed of spaces before it is set. The total is always replaced by the submitted total.

## Currency change is compared by effective currency code

Lines 123-128

Asked as: currency change same code; switch currency no change detected; currency spelling different same

Turbo trims the submitted currency and compares its effective currency code (EffectiveCurrencyCode) with the code of the current currency. The revision counts as a currency change only when these codes differ. A different currency value with the same effective code does not count. The total is compared with the current total for exact equality.

## Revision failure with an error message

Lines 128-134

Asked as: revision failed no notification; why did revision fail; revision error total unchanged

The revision message reads 'a revision must change the total or the currency' (code total_unchanged) when the new total equals the current total and the effective currency is the same. Other errors from the approval checks and from saving are also returned to the caller. Turbo does not send a notification when the revision fails.

## Approver is notified only when one is set

Lines 136-204

Asked as: no notification after revise; email approver when po edited; revision email not sent

After a successful revision, Turbo sends an immediate revision approval notification only when the purchase order has an approver. The notification includes the previous approval total and the creator. A failed notification is logged and does not fail the revision. The response returns the revised purchase order.

## Revision clears approvals and sets the second approver fields

Lines 137-148

Asked as: approval reset after edit; second approver removed on revise; approvals cleared when changed

A revision sets the status to Unapproved. It clears the approved and second_approval timestamps and the second_approver. A new approver or second approver, if submitted, is trimmed and saved.

## Revision returns the order to Unapproved and clears approvals

Lines 145-169

Asked as: po back to unapproved; need approval again after edit; revised order needs re-approval

A revision sets the purchase order status to Unapproved. It clears the approval (approved), the second approval (second_approval) and the second approver (second_approver), so the order needs approval again. The purchase order keeps its po_number. The change is saved in one transaction, so a failed check saves nothing.

## Revision re-checks approval totals, approvers and division

Lines 150-161

Asked as: approval totals recalculated after edit; approver policy check on revise; rejection cleared on revision

After a revision Turbo recomputes the approval totals and checks the approvers against the approval policy (CleanPurchaseOrder, ValidatePurchaseOrder). CleanPurchaseOrder also clears any rejection fields. Turbo refuses the revision when the division is not active (EnsureActiveDivision). The exact rules and messages come from those functions in the hooks package.

## Division must be active on revision

Lines 156-161

Asked as: inactive division revision; division disabled cannot edit po; division not active error

Turbo checks that the purchase order's division is active (EnsureActiveDivision) before it saves the revision. The validation runs with the approver policy checks (ValidatePurchaseOrder) and its errors are returned to the caller.

## Approver is notified of a revision

Lines 197-202

Asked as: email approver about revision; approver alert after edit; notification failed but revision saved

After a successful revision, if the order has an approver, Turbo sends that approver an immediate revision approval notification that includes the previous total. When the notification fails, Turbo logs the error and the revision still succeeds. An order with no approver sends no notification.
