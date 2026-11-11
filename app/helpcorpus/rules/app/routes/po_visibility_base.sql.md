<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/po_visibility_base.sql","sha256":"39872a5879772f48b6bed76bbaee082b67b0b6575cc7c3d34ded10581d7935e3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/po_visibility_base.sql

## Pending queue is narrower than purchase order visibility

Lines 31-35

Asked as: PO shows in list but not in my queue; why isn't PO in my pending; purchase order not in my approvals; unapproved PO I can see but can't approve

A person can see more Unapproved purchase orders than they can act on. The pending queue holds only the first-stage approver's items, the priority second approver's items, and the second-stage fallback items after the timeout. Seeing a purchase order does not put it in that person's pending queue.

## Seeing a purchase order does not show all its linked expenses

Lines 41-58

Asked as: can't see all expenses on PO; missing expenses on purchase order; can't open expense after PO closed; expense count doesn't match what I see

Purchase order visibility and expense visibility follow separate rules. A person who can open a purchase order may not see every expense linked to it. The committed expense count on the purchase order (committed_expenses_count) can be larger than the list of expenses that person sees. An expense owner can still view their own expense after the purchase order becomes Closed or Cancelled, even when they cannot view that purchase order.

## Report claim holders see Closed and Cancelled purchase orders

Lines 64-75

Asked as: can see closed POs but shouldn't; why can I view cancelled purchase orders; report access closed POs; see all closed purchase orders

The caller context (caller_context) sets has_report_claim when the user holds the report claim. That claim lets the user see every Closed or Cancelled purchase order (via is_closed_cancelled_visible).

## Approver limits depend on the PO's kind

Lines 98-150

Asked as: approval limit wrong for PO; my approver limit not applied; limit depends on PO type; why am I not eligible to approve

The limit used for an approver comes from the PO's expenditure kind. Capital uses max_amount, project uses project_max, and sponsorship, staff and social, media and event, and computer each use their own max. If the PO has no kind, it counts as project when it has a job and capital when it has none. An unknown kind gives no limit, so the approver is not eligible. Approver limits count only while the approver's admin profile is active.

## PO approver limits count only with an active admin profile

Lines 98-115

Asked as: approver limit missing; can't approve no limit showing; admin profile inactive approver; approver claim but no limit

Turbo reads a person's purchase order limits (caller_approver_props) only when the person holds the PO approver claim and the admin profile is active (ap.active = 1). Without both, no limit exists. The person then cannot qualify as an approver in the second-stage visibility or pending checks.

## Second approval threshold falls back to the job-based kind

Lines 125-150

Asked as: second approval not required; PO skipped second approver; what threshold for second approval; no second approval needed on PO

The PO's second approval threshold comes from its expenditure kind. If the PO has no kind, it uses the threshold of the project kind when the PO has a job, or the capital kind when it has none. If none is found, the threshold is 0, so no second approval is needed.

## Revisions are visible to every signed-in user

Lines 170-176

Asked as: rejected PO still showing; revised PO visible to everyone; rejection reason blank; resubmitted PO shows to others

A revision is an Unapproved PO that has a po_number. Any authenticated user can see it, so it does not look deleted while it waits for approval again. This does not make it actionable for them. The rejection reason is blank for everyone except the creator, the approver and the rejector.

## Print-limit and approval fields show 0 or blank when unset

Lines 194-194

Asked as: budget field shows zero; blank approver name on PO; empty fields on purchase order; covered budget shows 0

Turbo shows a missing 'covered within project budget' value as 0 (covered_within_project_budget). Names such as the user, approver, vendor, job, client, division, category and branch show as blank when the record has none.

## Approval total uses the home-currency amount when set

Lines 197-200

Asked as: approval uses wrong amount; home currency total not used; approval total mismatch; which amount is used for approval

Approval checks compare the PO's approval_total_home to limits and thresholds. If approval_total_home is 0 or empty, they use approval_total.

## Currency and display defaults

Lines 201-206

Asked as: PO shows CAD by default; no currency on purchase order; exchange rate shows 1; PO kind label wrong

A PO with no currency shows CAD as its code and symbol, and a rate of 1. Kind name and label fall back to project when the PO has a job and capital when it has none, if the PO kind is blank.

## Committed expense count counts only committed expenses

Lines 207-207

Asked as: committed expense count too high; why count differs from my expenses; count of committed expenses; expenses count incorrect on PO

A PO's committed expense count includes only expenses that have a commit date. The count can exceed the expenses a given user may see.

## Expenses total counts all expenses, committed or not

Lines 208-208

Asked as: expenses total includes uncommitted; total expenses on PO; sum of expenses wrong; total zero with no expenses

The expenses total (expenses_total) adds up the total of every expense linked to the purchase order, committed or not. It is 0 when there are none.

## Remaining amount calculation

Lines 209-222

Asked as: remaining amount wrong; PO balance calculation; how much left on PO; recurring PO remaining balance

A Recurring PO's remaining amount is its approval_total minus all its expenses. Other PO types use total minus all expenses. Expenses count whether committed or not. The remaining occurrences of a Recurring PO are its occurrences minus its committed expenses, never below 0.

## Recurring purchase orders show remaining occurrences, never below zero

Lines 209-214

Asked as: remaining occurrences negative; recurring PO occurrences left; how many occurrences remain; occurrences showing zero

For a Recurring purchase order, Turbo shows the expected occurrences (occurrences). It also shows the remaining occurrences: occurrences minus committed expenses, with a floor of 0. Other PO types show 0 for both.

## Maximum amount to print for the vendor

Lines 223-232

Asked as: how much can be printed; print amount for vendor; maximum amount to print; cumulative PO print total

The most that can be printed is the PO total, rounded to 2 decimals. For a Cumulative PO it is the total minus committed expenses only. Printing uses the PO currency.

## Names and lookup fields show blank when the record is missing

Lines 233-251

Asked as: PO missing creator name; blank vendor on purchase order; PO shows no job number; names blank on PO list

Turbo lists every purchase order even when a related record is missing. A person's name is their given name and surname joined with a space. The creator, approvers, rejector, canceller and closer names show '' when the profile is missing. The same blank applies to the parent PO number, vendor name and alias, job number and description, division code and name, category name, and branch code and name.

## Display kind falls back to job-based kind

Lines 252-253

Asked as: kind blank on PO; PO kind shows project; why kind says capital; kind name empty

When a purchase order has no kind (kind blank), Turbo shows the kind 'project' if the PO has a job and 'capital' if it has none (kind_name, kind_label). Without a match the name is blank.

## Active purchase orders are visible to every signed-in user

Lines 256-264

Asked as: can open active PO but not expenses; see active purchase orders; anyone can view active PO; active PO visible to all

Any authenticated user can see a purchase order whose status is Active (is_active_visible). Linked expenses follow separate expense visibility rules, so a user may open the PO and still not see all its expenses.

## Who can see Closed or Cancelled purchase orders

Lines 266-288

Asked as: who can view closed PO; cancelled PO visibility; why can't I see closed PO; job manager sees cancelled PO

A Closed or Cancelled purchase order is visible to its creator, its first approver and its second approver. It is also visible to the job's manager or alternate manager when the PO has a job, and to holders of the report claim. Expense owners are not included unless they fall into one of these groups.

## Direct visibility of Unapproved purchase orders

Lines 290-321

Asked as: creator can't see rejected PO; who sees my unapproved PO; rejected PO visibility; approver can't see new PO

The creator always sees their own Unapproved PO, including a rejected one. The rejector sees a rejected PO they rejected. While the PO is not rejected and has no first approval, the assigned approver sees it. After first approval and before second approval, the first approver and the priority second approver see it.

## Second-stage eligible approvers can see Unapproved purchase orders

Lines 323-356

Asked as: second approver can see PO; why can I view PO waiting second approval; approver sees PO but can't approve; PO over threshold visible to me

An approver can see an Unapproved, non-rejected PO that has its first approval and no second approval when all of these hold. Their approver profile is active and has a limit for the PO's kind. Their divisions list is empty or includes the PO's division. The PO's approval total exceeds the second approval threshold, which must be above 0. Their limit exceeds the threshold and covers the full total (limit >= total). The timeout is not applied to visibility.

## Division list limits which purchase orders an approver can act on

Lines 338-346

Asked as: divisions restrict my approvals; can't approve PO from my division; division limit on approver; approver restricted to divisions

An approver's division list (divisions) restricts the purchase orders that approver can vet or approve at the second stage. An empty list means no restriction. A non-empty list must contain the PO's division.

## Second-stage eligibility needs a limit strictly above the threshold

Lines 347-352

Asked as: limit equal to threshold; why am I not second approver; limit too low for second approval; threshold equals my limit

To be an eligible second approver, the person's limit must be strictly above the second-approval threshold. It must also be at least the approval amount, which must exceed a positive threshold. A limit equal to the threshold does not qualify.

## Revision visibility flag

Lines 358-365

Asked as: revision shows for everyone; who can see revised PO; revised PO visibility; revision not actionable

A revision is visible to any signed-in user (is_revision_visible). Visibility does not affect who can act on it.

## Legacy manual-entry purchase orders need the legacy claim

Lines 367-373

Asked as: legacy PO not visible; manual entry PO hidden; legacy purchase order access; can't see old manual PO

A PO flagged legacy_manual_entry is visible to holders of the legacy_po_create_update claim.

## PO group members are visible to super_po and admin holders

Lines 375-382

Asked as: group PO not showing; super PO can't see closed; PO group members hidden; see unapproved group POs

Holders of the super_po or admin claim see every PO that belongs to a PO group, in every status. This lets the people who create groups see the unapproved and closed ones, which belong to the staff member they were created for.

## Pending queue: first-stage approver

Lines 391-440

Asked as: why is PO pending for me; PO waiting for my approval; first approver queue; pending PO not showing

An Unapproved, non-rejected PO with no first approval is pending for its assigned approver only. The approver's divisions list must be empty or include the PO's division, and they must have a limit for the PO's kind. When the PO needs two approvals, they must have a positive limit below the PO total. In the single-stage case, any limit works.

## Pending actionability only applies to unrejected Unapproved purchase orders

Lines 391-396

Asked as: rejected PO still in queue; pending status wrong; PO shows pending incorrectly; who must act on PO now

A purchase order is pending for a person (is_unapproved_actionable_now) only when its status is Unapproved and it has no rejection. Then one of three paths must match: first-stage approver, priority second approver, or fallback after timeout.

## First-stage approver needs a limit and a matching division

Lines 397-413

Asked as: assigned approver can't approve; first stage approval blocked; limit for kind missing; division blocks my approval

For a stage-1 pending PO, the person must be the assigned approver and the PO must not yet be approved. The person's limit for the PO's kind must exist, and the division rule must pass.

## Two-stage first-stage approver needs a positive limit below the amount

Lines 414-422

Asked as: first approver can't act on two-stage PO; limit above amount first approver; two stage PO not pending; limit too high to approve first

When the amount exceeds a positive second-approval threshold, the assigned first approver sees the PO as pending only if their limit is above 0 and below the full approval amount.

## Pending queue: requester exception for two-stage purchase orders

Lines 423-430

Asked as: creator approver can't approve own PO; I created PO can't approve; self approval over threshold; requester approving own purchase order

When the PO exceeds the second approval threshold, the assigned first approver can act on it even when their limit covers the full amount, provided they are also the PO's creator and their limit is above 0.

## Single-stage purchase orders are pending for the assigned approver at any limit

Lines 431-437

Asked as: single stage PO not pending; why does low limit still approve; no second approval PO pending; one approver PO queue

When the second-approval threshold does not apply (threshold is 0, or the amount does not exceed it), the assigned first approver sees the PO as pending. Their limit only needs to exist, and it can be lower than the amount.

## Pending queue: priority second approver

Lines 441-446

Asked as: PO waiting for priority approver; second approver queue; first approved PO next step; who approves after first approval

After first approval and before second approval, an Unapproved non-rejected PO is pending for its priority second approver. No timeout applies to them.

## Pending queue: second-stage fallback after timeout

Lines 447-471

Asked as: other approvers can approve after wait; timeout before other approvers; PO stuck waiting second approval; fallback approver after timeout

If the first approval is older than the timeout window (timeoutHours), other eligible second approvers can act on the PO. They need a division match, a threshold above 0, a PO total above the threshold, and a limit above the threshold that covers the full total. Before the timeout, only the priority second approver can act.

## Fallback second-stage timeout is strictly older than the timeout

Lines 447-451

Asked as: how long until fallback approver; timeout hours second approval; when others can approve PO; first approval too recent

The fallback second-stage path requires the first approval (approved) to be older than the timeout window in hours (timeoutHours). Only then can any eligible approver see the PO as pending.

## Fallback second approver needs a limit above the threshold covering the amount

Lines 452-471

Asked as: fallback approver limit too low; why I can't take over PO; limit covers amount fallback; timeout approver not eligible

For the fallback path, the approver's limit must be greater than the second-approval threshold and at least the approval amount. The threshold must be above 0 and the amount above the threshold. The division rule also applies.

## Client of a purchase order comes from its job

Lines 486-489

Asked as: client missing on PO; PO client blank; which client is on purchase order; client comes from job

Turbo takes a purchase order's client name and client ID (client_name, client_id) from the client of its job. A purchase order with no job, or whose job has no client, shows a blank client. The job's manager and alternate manager also come from this job link.
