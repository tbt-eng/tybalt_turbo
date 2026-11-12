<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/revise_purchase_order.go","sha256":"d28fce3ebb50e03181f490d5b3867826242a8abcbb40b893543b81a6a7c2ce77","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/revise_purchase_order.go

## Revised purchase order keeps its PO number

Lines 1-4

Asked as: keep same PO number; revised PO number; po number after revision; reactivate same number

A revised purchase order keeps its po_number. When an approver approves the revision, the purchase order becomes Active again under the same number.

## Attachment cannot change in a revision

Lines 24-29

Asked as: change attachment on revision; replace attachment PO; upload new file revision; attachment cannot be changed

The revise action cannot change the purchase order's attachment, because the request is JSON. The creator replaces the attachment while the revision awaits approval.

## Revision request must always include the total

Lines 30-40

Asked as: total missing in revision; revision total zero; omitted total becomes 0; revise without total

The total is not optional in a revision request. Turbo always takes it from the request and never keeps the current total. A request that omits it is treated as a total of 0. Other fields, such as currency, job, description, division and category, keep their current values when left out.

## Revision refused when the request body is invalid

Lines 47-53

Asked as: invalid request body; revision request broken; bad revise request; malformed revision data

When Turbo cannot read the revision request body, it refuses the revision. The message is "invalid request body" (code invalid_request_body). Turbo makes no change to the purchase order.

## A revision succeeds or fails as a whole

Lines 57-57

Asked as: revision partially saved; revision failed rollback; PO changed after failed revise; all or nothing revision

Turbo does all the checks and the save of a purchase order revision in one transaction. If any check refuses or the save fails, the purchase order stays as it was, Active with its old values.

## Revising a missing purchase order is refused

Lines 58-65

Asked as: revise missing PO; purchase order not found; error fetching purchase order; PO does not exist revise

If the purchase order does not exist, Turbo refuses the revision. The message is "error fetching purchase order: ..." followed by the lookup error (code po_not_found).

## Only Active purchase orders can be revised

Lines 67-73

Asked as: cannot revise closed PO; revise non-active PO; revise draft purchase order; why can't I revise PO status

If the purchase order is not Active, Turbo refuses the revision with the message "only active purchase orders can be revised".

## Only the creator can revise a purchase order

Lines 75-81

Asked as: revise someone else's PO; not my purchase order; revise colleague's PO; not the PO creator

If the person revising is not the creator of the purchase order (uid), Turbo refuses with the message "only the creator can revise this purchase order".

## Legacy and imported purchase orders cannot be revised

Lines 83-92

Asked as: revise imported PO; legacy PO edit; manual legacy purchase order; cannot edit old imported PO

A purchase order with legacy_manual_entry set or imported (_imported) cannot be revised. Turbo refuses with the message "legacy purchase orders cannot be revised".

## Purchase orders from a PO group cannot be revised

Lines 94-104

Asked as: revise group PO; PO group revision blocked; grouped purchase order edit; po group can't change

If the purchase order belongs to a PO group (po_group is not blank), Turbo refuses the revision with the message "a purchase order from a PO group cannot be revised".

## A purchase order with expenses cannot be revised

Lines 106-124

Asked as: PO has expenses revise; expense linked to PO; change PO with expenses; expenses block PO revision

If any expense uses the purchase order, Turbo refuses the revision with the message "this purchase order has associated expenses and cannot be revised". Without expenses, the revision replaces the approval.

## A purchase order with child purchase orders cannot be revised

Lines 126-142

Asked as: PO has child POs; child purchase order blocks; sub PO revision; cancel child PO first

If the purchase order has a child purchase order (parent_po) that is not Cancelled, Turbo refuses with the message "this purchase order has child purchase orders and cannot be revised". Cancelled children do not block a revision.

## New job is checked by a shared job-change rule

Lines 144-153

Asked as: change job on PO; new job revision refused; switch job purchase order; job change not allowed

When a revision names a job, Turbo checks it with hooks.RevisionJobChangeRefusal and refuses with that function's code and message. The result depends on that function, which is in another file. Turbo trims spaces around the job, and a job left out keeps the current job.

## Budget confirmation changes only with a new job

Lines 154-161

Asked as: budget confirmation won't change; covered within budget edit; change budget checkbox PO; budget confirmation refused

A revision that changes covered_within_project_budget without changing the job is refused with the message "the budget confirmation can change only with a new job". The same value as the current one is accepted.

## Fields left out keep their current values

Lines 163-214

Asked as: blank fields revision; fields reset on revise; missing values cleared; keep old currency description

A revision keeps the current currency, description, division and category when the request leaves them out. It keeps the approver and priority second approver the same way. The total is always taken from the request. Text values are trimmed of spaces, except the description.

## A new job resets budget confirmation and category

Lines 171-178

Asked as: category cleared after job change; budget confirmation reset; new job resets category; job switch clears category

When the job changes, Turbo sets the budget confirmation to the request's value, or to false if the request gives none. It also clears the category unless the request sets one.

## A revision must change something

Lines 179-193

Asked as: nothing changed revision; no change PO revise; same values revision refused; currency spelling same revision

If the total, currency, job, description, division and category all stay the same, Turbo refuses with the message "a revision must change the total, currency, job, description, division or category". Turbo compares currencies by their effective currency code (EffectiveCurrencyCode). A different spelling of the same effective currency is not a change.

## A revision is checked again like a new purchase order

Lines 195-231

Asked as: PO checks on revision; division inactive PO; approver policy check revise; revision recalculated approvals

Before saving, Turbo recalculates the approval totals and clears any rejection fields (CleanPurchaseOrder). It checks that the division is active and that the approvers fit the policy (ValidatePurchaseOrder). It also checks the division against the job. The exact rules are in other files.

## Approvers can be changed with a revision

Lines 209-218

Asked as: change approver on PO; new second approver; replace approver revision; approval cleared when revising

A revision can set a new approver and a new priority second approver. A value that is left out keeps the current one. Turbo always clears the previous approval, the second approval and the second approver (second_approver).

## Revising a purchase order returns it to Unapproved

Lines 215-218

Asked as: PO back to unapproved; revision status unapproved; approval removed after revise; PO needs reapproval

A revision of an Active purchase order sets its status to Unapproved in the same transaction. It keeps its po_number, so approving it again re-activates the same number. It clears the approved date, the second approval and the second approver (approved, second_approval, second_approver).

## Revision refused when the purchase order cannot be saved

Lines 233-239

Asked as: cannot save PO revision; error updating purchase order; PO save failed revise; revision database error

If Turbo cannot save the revised purchase order, it refuses with "error updating purchase order: ..." and the cause (code error_updating_purchase_order). The purchase order is not changed.

## Validation failures on a revision are shown as a validation error

Lines 244-260

Asked as: validation error revise PO; field errors PO revision; revision rejected by rule; which fields are wrong

When the purchase order checks fail with field errors, Turbo answers "Validation error" and lists the fields with a problem. When a purchase order rule (hook) refuses the revision, Turbo returns that rule's own error and status. Turbo does not save the revision in either case.

## Approver is notified after a revision

Lines 267-272

Asked as: approver notification revision; email approver after revise; notify approver PO change; no notification sent revision

After a revision, if the purchase order has an approver, Turbo sends the approver a revision notification immediately. If the notification fails, the revision still stands.
