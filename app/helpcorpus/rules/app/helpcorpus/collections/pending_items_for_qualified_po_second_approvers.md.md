<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/pending_items_for_qualified_po_second_approvers.md","sha256":"81d88d34d0a1e920b3f6ac8d135b4b3d13e5d718aaa5c907857e1646c3a2bd00","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/pending_items_for_qualified_po_second_approvers.md

## Nobody can list second-approver pending counts through the API

Lines 7-9

Asked as: list second approver counts; can't list pending counts; api list view records; who can read pending approvers view

The view pending_items_for_qualified_po_second_approvers has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can view one second-approver pending count through the API

Lines 11-13

Asked as: view one pending count; can't open approver count record; api view single record; second approver count not visible

The view pending_items_for_qualified_po_second_approvers has no view rule. Nobody can view one of its records through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Second-approval wait defaults to 24 hours

Lines 18-33

Asked as: how long before second approval; 24 hour wait for PO; second_stage_timeout_hours; when does PO need second approver

The wait before a purchase order (PO) needs a second approver comes from second_stage_timeout_hours in the 'purchase_orders' app_config value. Turbo uses 24 hours when that value is missing, is not valid JSON, is not a number, or is 0 or less.

## Only active PO approvers can be qualified second approvers

Lines 34-49

Asked as: who is a qualified second approver; not showing as PO approver; po_approver claim missing; inactive approver not listed

A person counts only if they have an active admin profile and the po_approver claim with its limits (po_approver_props).

## Which POs count as waiting for a second approval

Lines 50-84

Asked as: why PO not waiting for approval; PO missing from second approval queue; approved PO not needing second approver; which POs need second sign off

A PO counts only if it is approved (approved is set), not rejected, has status 'Unapproved', and has no second approval yet. Its approval must be older than the wait (24 hours by default). Its approval_total must be strictly above the second approval threshold, and that threshold must be above 0.

## Second approval threshold falls back to the kind by job

Lines 56-67

Asked as: no threshold for expenditure kind; PO kind default threshold; project or capital fallback limit; PO never asks for second approval

The threshold comes from the PO's expenditure kind (second_approval_threshold). With no kind, Turbo uses the kind named 'project' if the PO has a job, or 'capital' if it has none. If no threshold is found, it is 0 and the PO never needs a second approval.

## Approver's limit depends on the PO's kind

Lines 91-99

Asked as: approver limit by PO kind; which max amount applies; approver spending limit per kind; project_max or max_amount

Turbo picks the approver's limit by the PO's kind. Capital uses max_amount, project uses project_max, sponsorship uses sponsorship_max, staff_and_social uses staff_and_social_max, media_and_event uses media_and_event_max, and computer uses computer_max. A missing limit counts as 0, and any other kind gives 0.

## Approver's divisions must include the PO's division

Lines 102-108

Asked as: approver division not matching; PO division approver can't see; divisions list approver wrong; approver can approve all divisions

An approver matches a PO when their divisions list is empty (all divisions) or contains the PO's division. If the divisions value is not valid JSON, the approver matches no PO.

## Qualification limit must exceed threshold and cover the total

Lines 110-117

Asked as: approver limit too low for PO; equal limit approve PO; approver can't cover PO total; why approver not qualified amount

An approver qualifies for a PO only when their limit for that kind is strictly above the second approval threshold. The limit must also be at least the PO's approval_total, so an equal limit passes.

## View returns count of qualifying POs per approver

Lines 118-122

Asked as: count of POs per approver; how many POs approver can second; num_pos_qualified meaning; approver with zero qualifying POs

The view returns one row per qualified approver. The row has the user's id and num_pos_qualified, the number of POs that person qualifies to second-approve. Approvers with no qualifying PO get no row.
