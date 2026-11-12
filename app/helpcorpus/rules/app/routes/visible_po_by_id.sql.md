<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/visible_po_by_id.sql","sha256":"aac9074b04cc169ec000e9ee1b13f5612e1c2f1b89208c281c3f4f906874d225","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/visible_po_by_id.sql

## PO detail includes calculated amounts, names and approval fields

Lines 4-76

Asked as: PO detail fields; purchase order totals; remaining amount on PO; approver names on PO; PO calculated values

The PO detail returns the stored PO fields together with calculated values: expenses total, remaining amount, print max amount, approval totals, committed expense count and recurring occurrence counts. It also returns the names of the people, vendor, client, job, division, branch, category and kind. These values come from the visibility base query, which this file does not show.

## Purchase order can be revised only when all conditions hold

Lines 77-86

Asked as: why can't I revise my PO; revise button missing; edit an active PO; PO not revisable; imported PO revision

A purchase order (PO) is marked revisable (revisable) only when its status is 'Active' and it is not a legacy manual entry (legacy_manual_entry = 0). It must also have no PO group (po_group is empty) and must not be an imported record (_imported = 1 excludes it). The revise route still runs its creator check and other validation when the revision is submitted.

## PO with any expense cannot be revised

Lines 84-84

Asked as: PO has expense can't revise; revise blocked by expense; expense attached to PO; why PO locked after expense

A purchase order is not revisable once any expense links to it (an expense whose purchase_order is this PO). One expense is enough to block revision.

## PO with a non-cancelled child PO cannot be revised

Lines 85-85

Asked as: child PO blocks revision; parent PO can't be changed; sub PO stops revise; linked child purchase order

A purchase order is not revisable when another PO has it as parent (parent_po) and that child's status is not 'Cancelled'. A cancelled child does not block revision.

## Single PO lookup returns it only if the person may see it

Lines 88-99

Asked as: can't see a PO; PO not found by id; purchase order missing for me; PO lookup returns nothing; no access to this PO

Looking up one purchase order by id returns it only when the visibility base marks it visible to the person asking. The visibility cases are active, closed or cancelled, unapproved direct, unapproved second-stage eligible, revision, legacy and group. If none applies, nothing is returned. The visibility rules themselves come from the __PO_VISIBILITY_BASE__ query, which is not in this file.
