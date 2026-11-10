<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/visible_po_by_id.sql","sha256":"5fed736401899bcad5689f35c6199de34d509ffa8fc590f28eb8b2ab0fc7a345","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/visible_po_by_id.sql

## Who can open a purchase order by its id

Lines 1-3

Asked as: can't open purchase order; PO not showing; no access to PO; purchase order not found; which POs can I see

Turbo returns a purchase order by id only when it passes at least one visibility test. The tests are: active (is_active_visible), closed or cancelled (is_closed_cancelled_visible), unapproved direct (is_unapproved_direct_visible), unapproved second stage eligible (is_unapproved_second_stage_eligible), revision (is_revision_visible), legacy (is_legacy_visible) and group (is_group_visible). If none passes, Turbo returns no purchase order. The tests themselves are defined in the visibility base (__PO_VISIBILITY_BASE__), which this file does not show.

## Purchase order lookup returns at most one record

Lines 4-98

Asked as: look up PO by id; PO details; single purchase order record; show PO totals; PO remaining amount

The lookup matches one purchase order id ({:id}) and returns a single row with the PO's details, such as status, total, approvers, vendor, job, branch, category, currency, expenses totals and remaining amount.

## Revisable flag is only a pre-check

Lines 77-85

Asked as: revise button says yes but fails; revisable but rejected; revision still blocked; can revise then error; revise check passed

The revisable flag only says whether the revise route can accept the purchase order. The check that the user is the creator, and the rest of the validation, still run when the revision is submitted. A true flag does not guarantee the revision is accepted.

## A purchase order can be revised only when it is Active

Lines 79-80

Asked as: revise PO not active; edit approved purchase order; change PO after approval; PO status not allowed to revise; can't revise draft PO

The revisable flag is true only when the purchase order's status is 'Active'. For any other status, Turbo marks the PO as not revisable.

## Legacy manual entry POs cannot be revised

Lines 81-81

Asked as: revise manual entry PO; legacy PO can't change; manual PO revision blocked; old manual purchase order edit

A purchase order with legacy_manual_entry other than 0 is not revisable.

## Imported purchase orders cannot be revised

Lines 82-82

Asked as: imported PO can't edit; change imported purchase order; revise PO from import; imported order locked

A purchase order that was imported (_imported = 1) is not revisable.

## A purchase order with expenses cannot be revised

Lines 83-83

Asked as: PO with expenses edit; change PO after expense added; can't revise PO linked expense; expense stops PO change; amend PO with claims

If any expense is linked to the purchase order (expenses.purchase_order), the PO is not revisable. A single expense is enough.

## A purchase order with a non-cancelled child PO cannot be revised

Lines 84-84

Asked as: PO has child order; revise parent purchase order; child PO blocks revision; split PO can't change; parent PO locked by child

If another purchase order has this PO as its parent (parent_po) and that child's status is not 'Cancelled', the PO is not revisable. Cancelled children do not block a revision.
