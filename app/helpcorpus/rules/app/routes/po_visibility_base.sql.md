<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/po_visibility_base.sql","sha256":"d07549eec8d9db276bc3c9640d674c4758a0a690ae829d55acd0c0569887d304","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/po_visibility_base.sql

## PO visibility differs from linked expense visibility

Lines 41-58

The ability to open a purchase order does not give access to every expense linked to it. Expense visibility follows separate rules. The count of committed expenses and the expense total on a PO can therefore be larger than the list of expenses a user sees. An expense owner can view their own expense even when they can no longer see the PO after it is Closed or Cancelled.

## Report claim holders are identified for visibility

Lines 64-75

Turbo checks whether the caller holds the report claim (name 'report'). That claim lets the caller see Closed and Cancelled purchase orders they have no other link to.

## Legacy PO claim holders are identified

Lines 76-85

Turbo treats a user as a legacy PO editor when the user holds the claim named 'legacy_po_create_update'. The flag has_legacy_po_create_update_claim is 1 for them and 0 for everyone else. Legacy manual PO visibility uses this flag.

## super_po and admin claims give PO group access

Lines 86-95

A user who holds the 'super_po' claim or the 'admin' claim gets has_po_group_claim = 1. Any other user gets 0. Turbo uses this flag to show purchase order group members.

## Approver limits count only with active admin profile

Lines 98-115

Turbo reads a user's PO approver limits (po_approver_props) only when the user holds the PO approver claim. The user's admin profile must also be active (active = 1). Otherwise the user has no limits and cannot qualify as a second-stage approver or as a limit-based approver.

## Which limit applies to a purchase order

Lines 125-150

Turbo picks the approver's limit from the PO's expenditure kind: capital (max_amount), project, sponsorship, staff and social, media and event, or computer. When the PO has no kind, Turbo treats it as project if it has a job and as capital if it has no job. An unknown kind gives no limit (NULL), so the person cannot approve.

## Second approval threshold falls back to default kind

Lines 125-150

Turbo takes the second approval threshold from the PO's expenditure kind. When the PO has no kind, it takes the threshold of project (job present) or capital (no job). When none is found, the threshold is 0, which means no second approval.

## Revision rejection reason is hidden from others

Lines 170-176

On a revision, the rejection reason is blank unless the viewer is the creator, the approver or the rejector.

## Project budget coverage flag defaults to 0

Lines 194-194

When a purchase order has no value for covered_within_project_budget, Turbo reports 0, which means not covered.

## Approval total uses home currency when set

Lines 197-200

For threshold and limit comparisons Turbo uses the approval total in home currency. When that is 0 or empty, it uses approval_total.

## Currency defaults to CAD

Lines 201-206

When a PO has no currency record, Turbo shows code CAD, symbol CAD, rate 1 and no icon or rate date.

## Committed expenses count and expense total

Lines 207-208

Turbo counts only committed expenses as committed expenses. The expense total on a PO sums all linked expenses, committed or not. These figures cover every expense, not only those the viewer can see.

## Recurring PO occurrences remaining

Lines 209-214

A Recurring PO shows its expected occurrences. Remaining occurrences equal expected minus committed expenses and never go below 0. Other PO types show 0 for both.

## Remaining amount calculation

Lines 215-222

For a Recurring PO the remaining amount is the approval total minus all linked expenses. For other types it is the PO total minus all linked expenses.

## Print maximum amount for Cumulative POs

Lines 223-229

For a Cumulative PO the print maximum is the total minus committed expenses only, rounded to 2 decimals. For other types it is the PO total, rounded to 2 decimals.

## Missing names and display values show blank

Lines 230-250

When a related person, vendor, job, client, division, category, branch or parent PO is missing, Turbo shows an empty text. The kind name and label fall back to project (job present) or capital (no job) when the PO has no kind.

## Active purchase orders are visible to everyone

Lines 253-261

Any signed-in user can see a purchase order whose status is Active (is_active_visible).

## Who sees Closed or Cancelled purchase orders

Lines 263-285

A Closed or Cancelled PO is visible only to its creator, its first approver, its second approver, the job's manager, the job's alternate manager, and holders of the report claim. Expense owners do not see it unless they are in one of these groups.

## Direct visibility of Unapproved purchase orders

Lines 287-318

The creator always sees their own Unapproved PO, including a rejected one. The rejector sees a rejected PO they rejected. For a PO that is not rejected and has no first approval, only the assigned approver sees it. After first approval and before second approval, the first approver and the priority second approver see it.

## Second-stage approvers can see Unapproved purchase orders

Lines 320-353

Once the first approval is done, a PO is visible to anyone who could be its second approver, with no timeout. The PO must not be rejected and must have no second approval. The PO's approval total (home currency when set) must exceed the second approval threshold, which must be above 0. The person's limit must be above the threshold and at least the total. Their division list must be empty or include the PO's division.

## Revisions stay visible to all users

Lines 355-362

A revision is an Active PO sent back to Unapproved with its PO number. Every signed-in user can see it, so it does not look deleted. Revisions do not change who can act on the PO.

## Legacy manual POs visible to legacy claim holders

Lines 364-370

A PO entered manually as legacy (legacy_manual_entry = 1) is visible to holders of the legacy_po_create_update claim.

## PO group members visible to super_po and admin

Lines 372-379

Holders of the super_po or admin claim see every purchase order that belongs to a PO group, in any status.

## Stage 1 pending: assigned approver within limits

Lines 388-437

An Unapproved, not rejected PO with no first approval is pending for its assigned approver only. The approver's division list must be empty or include the PO's division, and the approver must have a limit for the PO's kind. One of three cases must apply. (1) The threshold is exceeded and the approver's limit is above 0 but below the total (first-stage vetting). (2) The threshold is exceeded, the approver is also the requester, and the limit is above 0. (3) The threshold is not exceeded, so the PO needs one approval. In case 3 the check does not compare the limit with the total.

## Stage 2 pending for priority second approver

Lines 438-443

After first approval, and with no second approval, the PO is pending for its priority second approver. This applies even when the person's limit would not cover the PO. The PO must not be rejected.

## Stage 2 fallback after timeout

Lines 444-468

Once the first approval is older than the timeout in hours (timeoutHours), any eligible second approver sees the PO as pending. The PO must not be rejected and must have no second approval. The approver's division must match, the PO's total must exceed the threshold above 0, and the approver's limit must be above the threshold and at least the total.

## Blank PO kind shows default kind name and label

Lines 489-492

When a purchase order's kind is empty, Turbo shows the kind name and label of the 'project' kind if the PO has a job. If the PO has no job, it shows the 'capital' kind. A PO with a kind set shows its own kind.
