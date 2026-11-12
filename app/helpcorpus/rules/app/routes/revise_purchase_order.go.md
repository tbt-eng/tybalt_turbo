<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/revise_purchase_order.go","sha256":"60ce8ede216163206b2ee9c3c584e444a9b7067faad59eed988752221519f3be","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/revise_purchase_order.go

## Revising an active purchase order returns it to Unapproved

Lines 1-5

Asked as: edit approved PO; PO back to unapproved; changed PO lost approval; need to approve again

A revision of an active purchase order (status Active) sets its status to Unapproved and clears approved, second_approval and second_approver. The purchase order keeps its po_number, so a new approval re-activates the same number. The change is saved in one transaction.

## Revision request leaves unset fields unchanged

Lines 24-40

Asked as: fields left blank; missing total becomes zero; partial update purchase order; blank field wipes value

A field that the revision request omits keeps its current value. The request can send a new total, currency, approver, priority second approver, job, description, division, category or budget confirmation. The total is not optional: an omitted total is read as 0.

## Unreadable revision request body is refused

Lines 47-53

Asked as: invalid request body; revision won't parse; bad json revise PO

When the revision request body cannot be read, Turbo refuses the revision with the message "invalid request body" (code invalid_request_body). Turbo checks the body before it looks up the purchase order.

## Revision is saved all together or not at all

Lines 57-57

Asked as: revision half saved; partial changes on error; rollback PO changes

Turbo looks up, checks and saves a purchase order revision in one transaction. When any check fails, the purchase order keeps its previous details and its Active status.

## Revision of a missing purchase order is refused

Lines 58-65

Asked as: purchase order not found; PO does not exist; po_not_found

When the purchase order does not exist, Turbo refuses the revision with the code po_not_found.

## Only active purchase orders can be revised

Lines 67-73

Asked as: can't edit PO status; PO not active; po_not_active

Turbo refuses to revise a purchase order whose status is not Active. The message is "only active purchase orders can be revised" (po_not_active).

## Only the creator can revise a purchase order

Lines 75-81

Asked as: someone else's PO; not my purchase order; not the creator edit

Turbo refuses a revision from anyone other than the person who created the purchase order (uid). The message is "only the creator can revise this purchase order".

## Legacy and imported purchase orders cannot be revised

Lines 83-92

Asked as: imported PO can't edit; legacy manual entry locked; old PO cannot change

Turbo refuses to revise a purchase order that is a legacy manual entry (legacy_manual_entry) or imported (_imported). The message is "legacy purchase orders cannot be revised". These orders never went through Turbo approval, so they cannot return to it.

## Purchase order with expenses cannot be revised

Lines 94-112

Asked as: PO has expenses can't edit; expense attached to purchase order; linked expenses block change

Turbo refuses to revise a purchase order if any expense uses it. The message is "this purchase order has associated expenses and cannot be revised". A purchase order with no expenses can be revised.

## Purchase order with child purchase orders cannot be revised

Lines 114-130

Asked as: child purchase order blocks edit; sub PO revise; parent PO can't change

Turbo refuses to revise a purchase order that has a child purchase order (parent_po) whose status is not Cancelled. The message is "this purchase order has child purchase orders and cannot be revised". Cancelled children do not block a revision.

## Job change checked by RevisionJobChangeRefusal

Lines 132-141

Asked as: new job rejected; can't switch job on PO; job change not allowed

Turbo asks hooks.RevisionJobChangeRefusal whether the new job is allowed and refuses the revision with that code and message when it objects. If the request has no job, the current job is kept and checked.

## Budget confirmation changes only with a new job

Lines 142-149

Asked as: change covered within budget; budget checkbox won't change; within project budget toggle

A revision cannot change the budget confirmation (covered_within_project_budget) unless the job also changes. Turbo refuses with "the budget confirmation can change only with a new job" (revision_requires_job_change). Sending the same value as the current one is accepted.

## New job clears category and budget confirmation

Lines 159-166

Asked as: category reset when job changed; budget confirmation reset to false; switching job wipes category

When a revision changes the job, Turbo clears the category unless the request sets one. It sets the budget confirmation (covered_within_project_budget) to true only when the request sends true, otherwise to false. When the job is unchanged, the category stays unless the request sets one.

## Revision must change something

Lines 167-181

Asked as: no changes saved; nothing changed revision error; same values revise PO

Turbo refuses a revision that leaves the total, currency, job, description, division and category all the same. The message is: "a revision must change the total, currency, job, description, division or category". Turbo compares currencies by their effective currency code (EffectiveCurrencyCode).

## Revision with a new job sets the job and budget confirmation

Lines 188-196

Asked as: moving PO to another job; job not saved on purchase order; budget flag saved with job

Only a job change makes Turbo write the job and the budget confirmation on the purchase order. Division, description, category, total and currency are always written from the request or the current values.

## Revision can replace approvers

Lines 197-206

Asked as: change approver on PO; new second approver; swap approvers purchase order

The revision request can set a new approver and a new priority second approver. Turbo trims the values. Turbo always clears the earlier approval and the second approver.

## Revision recomputes approval totals and validates approvers and division

Lines 208-219

Asked as: approval totals wrong after edit; division not allocated to job; approver not valid for division

After the changes, Turbo recalculates the approval totals and clears any rejection fields (CleanPurchaseOrder). It requires the division to be active (EnsureActiveDivision). It then validates the purchase order (ValidatePurchaseOrder), which checks that the division is allocated to a new job, that the job and category match, and the approvers against the new division. A failure refuses the revision.

## Revision requires an active division

Lines 214-216

Asked as: division inactive PO; closed division purchase order; division not active revise

After it recomputes the approval totals, Turbo checks that the revised purchase order's division is active. When it is not, Turbo refuses the revision with the error that EnsureActiveDivision returns for the division field.

## Failed saving of a revision is refused

Lines 221-227

Asked as: error updating purchase order; PO save failed; could not save revision

When Turbo cannot save the revised purchase order, it refuses the revision with the message "error updating purchase order: ..." (code error_updating_purchase_order). The purchase order stays unchanged.

## Revision errors from validation are returned to the user

Lines 232-247

Asked as: validation error PO; approval policy failed revision; revision fields failing validation

When the approval policy or division checks in a revision fail, Turbo returns that error to the requester. The purchase order does not change. A validation error comes back as "Validation error" with the failing fields.

## Revision refusals show a message and code

Lines 232-241

Asked as: why was revision refused; revision error message; revise PO refused code

Each refusal of a revision, such as not Active, not the creator or no change, returns its own message and code.

## Approver is notified of a revision

Lines 250-262

Asked as: approver not notified; no approval email after edit; notification after revising PO

After a revision is saved, Turbo sends an immediate approval notification to the approver when the purchase order has an approver. A failed notification is logged and does not undo the revision.
