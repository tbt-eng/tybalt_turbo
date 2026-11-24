<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders.go","sha256":"8d9d949ea89fc1921993e040ce7b6e56b5d1264f4520022f06a9e864e687f746","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders.go

## Revised purchase order is Unapproved with a PO number

Lines 23-29

Asked as: edited PO still waiting; revised purchase order number; changed PO keeps number; revision awaiting approval

Turbo treats a purchase order as under revision when its status is Unapproved, it already has a PO number (po_number) and it is not a legacy manual entry (legacy_manual_entry = 0). A revision keeps its number while it waits for approval.

## Second approver lookup status messages

Lines 179-214

Asked as: why no second approver; second approver message; nobody can final approve; second approval status

The second approver lookup reports one of four states. 'not_required' means second approval is not required. 'requester_qualifies' means the requester can do the second approval. 'candidates_available' means eligible second approvers exist. Otherwise the status is 'required_no_candidates': "Second approval is required, but no second-stage approver can final-approve this amount."

## Second approval timeout comes from a setting

Lines 216-234

Asked as: second stage timeout hours; how long before second approver; when others can approve; approval timeout setting

The visibility queries for pending and visible purchase orders pass the second-stage timeout in hours (timeoutHours). Turbo reads it from the purchase order second-stage timeout setting. The same value appears in the second approver lookup as second_stage_timeout_hours.

## Single visible purchase order hidden if not visible

Lines 236-255

Asked as: PO disappeared; can't find purchase order; PO not showing; missing PO no access

Turbo looks up one purchase order through the same visibility rules as the list. If the purchase order is missing or the caller may not see it, the caller gets no row. Turbo does not say which of the two applies.

## Only unapproved purchase orders can be approved

Lines 281-297

Asked as: can't approve this PO; approve button not working; approve rejected order; approve already rejected

Approving a purchase order fails unless its status is Unapproved. The message is "only unapproved purchase orders can be approved". A rejected PO also cannot be approved: "rejected purchase orders cannot be approved".

## Approval uses the home-currency total and the kind

Lines 299-323

Asked as: kind no longer exists; approval fails invalid kind; PO kind removed; who can approve amount

Turbo works out who may approve from the purchase order's division, its approval total in home currency (EffectiveApprovalTotalHome), its kind and whether it has a job. It normalizes the kind with the job flag. If the kind is invalid or no longer exists, the approval fails with "purchase order kind is invalid or no longer exists".

## Second approval requirement is set by the policy

Lines 325-333

Asked as: why need second approval; when second approval needed; second sign off required; needs two approvers

Turbo decides whether second approval is required from the approval policy, not from a stored flag. It records whether the first approval already exists (approved) and whether the second approval exists (second_approval). Those facts pick the stage of the approval.

## Assigned approver may be valid through a bypass

Lines 333-333

Asked as: assigned approver not valid; approver bypass owner; approver can't approve my PO; assigned approver invalid

The assigned approver is valid for first-stage approval if they are in the first-stage pool, or if they are valid through the bypass for the purchase order owner (IsValidFirstStageViaBypass). If the assigned approver is blank or neither, approval fails with "assigned approver is not valid for first-stage approval".

## PO number assigned on activation

Lines 335-350

Asked as: PO number missing when active; no number after approval; when PO number assigned; active PO without number

When a PO becomes Active without a PO number, Turbo generates one. A PO that already has a number, such as a revision, keeps it.

## Combined approval by a qualified second approver on a fresh PO

Lines 353-371

Asked as: approve as both approvers; one person both stages; second approver approves first; approve both stages

When second approval is required and nothing is approved yet, a person who is a second-stage approver can approve as both first and second approver in one action. Turbo sets both approvers to that person and both approval times to now, and the PO becomes Active. The person's limit must cover the approval total, or Turbo refuses with "you do not have sufficient approval limit to complete final approval".

## Combined approval sets both approvers

Lines 355-371

Asked as: same approver both fields; one approval both stages; limit too low to approve; single approval both approvers

When second approval is required and the caller is a second-stage approver, one approval does both stages. Turbo sets approver and second_approver to the caller, sets approved and second_approval to the current time, and activates the purchase order. A caller whose limit is too low gets "you do not have sufficient approval limit to complete final approval".

## First approval only by the assigned approver

Lines 372-379

Asked as: not authorized to approve; why can't I approve; only assigned approver approves; approve someone else's PO

On a PO with no first approval, only the approver assigned to it can approve. Anyone else gets "you are not authorized to approve this purchase order". This applies unless the caller took the combined-approval path as a second-stage approver.

## First approval checks when second approval is required

Lines 381-413

Asked as: second approval pool empty; priority second approver required; contact administrator approval pool; first approval refused second approver

When a PO needs second approval, first approval is refused in these cases. Turbo shows "second approval pool is empty for this purchase order; contact an administrator" if no second-stage approvers exist. It shows "priority second approver is required when second approval is required" if none is set. It shows "priority second approver is not valid for second-stage approval" if the chosen person is not a second-stage approver. It shows "first approval pool is empty for this purchase order; contact an administrator" if there are no first-stage approvers and the assigned approver is not valid through the bypass.

## Dual-approval first approval needs a priority second approver

Lines 389-405

Asked as: pick second approver; no priority approver set; priority approver not valid; choose second approver first

When second approval is required, the first approver must have set a priority second approver (priority_second_approver) on the purchase order. If none is set, approval fails with "priority second approver is required when second approval is required". If the one set is not a second-stage approver, it fails with "priority second approver is not valid for second-stage approval".

## Assigned approver must be valid for first stage

Lines 415-421

Asked as: approver not valid first stage; no assigned approver; approver can't give first approval; assigned approver missing

First approval is refused with "assigned approver is not valid for first-stage approval" if the PO has no assigned approver. The same happens if that approver is neither a first-stage approver nor valid through the bypass.

## First approval activates single-stage POs

Lines 423-432

Asked as: approved but still unapproved; when does PO become active; after first approval status; single stage activates

After first approval, Turbo records the approver and approval time. If second approval is not required, the PO becomes Active at once. If second approval is required, the PO stays Unapproved and waits for the second approver.

## Final approval by a second-stage approver

Lines 436-464

Asked as: final approval not authorized; second sign off limit; second approval refused; final approver limit

On a PO that is first-approved and needs second approval, the second approval must come from a second-stage approver. Otherwise Turbo refuses with "you are not authorized to perform second approval on this purchase order". The approver's limit must cover the approval total, or Turbo refuses with "you do not have sufficient approval limit to complete final approval". On success the PO becomes Active, with the second approver and time recorded. A PO that already has a second approval is refused with "only unapproved purchase orders can be approved".

## Second approval already given is refused

Lines 436-450

Asked as: already second approved; approve again error; second approval already given; only unapproved approved error

If a purchase order that needs second approval has already been second-approved, a new approval attempt fails with "only unapproved purchase orders can be approved". A caller who is not a second-stage approver gets "you are not authorized to perform second approval on this purchase order".

## Final approval needs a sufficient limit

Lines 451-464

Asked as: final limit not enough; approval limit too low; limit does not cover total; final approval limit

A second-stage approver whose final limit does not cover the approval total gets "you do not have sufficient approval limit to complete final approval". If the limit covers it, Turbo records second_approver and second_approval and activates the purchase order.

## Activating a first-approved single-stage PO

Lines 465-485

Asked as: finish approving PO; activate first approved PO; complete single stage approval; approve first approved order

A first-approved PO that does not need second approval can still be approved. Its assigned approver must be a first-stage approver, or Turbo shows "assigned approver is not valid for first-stage approval". The caller must be the assigned approver or another first-stage approver, or Turbo shows "you are not authorized to approve this purchase order". The PO then becomes Active.

## Approval returns the updated PO with expanded details

Lines 488-522

Asked as: what approve returns; PO details after approval; updated PO after approving; approval response details

After approval Turbo saves the changes in one transaction and returns the updated purchase order. It includes the creator, the approver, the division and the job.

## Approval notifications

Lines 527-580

Asked as: email after approval; priority approver emailed; no approval email; notification after approving

After first approval of a PO that needs second approval, with a priority second approver set and the PO not Active, Turbo emails the priority second approver. For a PO in a group that has no number yet, one email covers the whole group. When a PO becomes Active and the caller is not its creator, Turbo notifies the creator that it is active.

## Priority second approver is told after first approval

Lines 527-556

Asked as: email second approver; told to second approve; priority approver notification; who gets emailed

When a purchase order has just been first-approved, needs second approval, has a priority second approver and is not yet Active, Turbo emails that person (po_priority_second_approval_required). A purchase order in a group with no number sends one email for the whole group. A revised purchase order, which has a number, gets the single email.

## Creator is told when PO becomes Active

Lines 557-580

Asked as: email when PO active; notify creator approved; creator told approved; active notification sent

When an approval activates the purchase order, Turbo sends the creator the po_active message with the approver's name. It does not send it when the creator is the person who approved.

## Pending purchase orders list

Lines 587-640

Asked as: what I need to approve; approvals waiting for me; pending approvals list; PO not in pending

Turbo returns the POs the caller can approve now. These are stage-1 POs assigned to the caller, and second-stage POs where the caller is the priority second approver. After the second-stage timeout, they also include second-stage POs that any eligible second approver can approve. A single pending PO lookup returns "pending purchase order not found" if the PO is not in that list.

## Visible purchase orders scope and parameters

Lines 642-707

Asked as: scope not allowed; list filter scope; stale_before required; expiring_before missing

The visible POs list accepts a scope: all (the default), mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. Another scope is refused with "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". Scope group needs a group ("group is required when scope=group"). Scope stale needs stale_before ("stale_before is required when scope=stale"). Scope expiring needs expiring_before ("expiring_before is required when scope=expiring"). A limit must be a positive integer, or Turbo says "limit must be a positive integer".

## Visible purchase orders scope values and required parameters

Lines 646-681

Asked as: invalid scope list; which scopes available; group scope needs group; scope values list

The scope defaults to all and must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group. Otherwise Turbo says "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". Scope group needs group ("group is required when scope=group"). Scope stale needs stale_before ("stale_before is required when scope=stale"). Scope expiring needs expiring_before ("expiring_before is required when scope=expiring").

## Visible purchase orders limit must be positive

Lines 682-695

Asked as: limit must be positive; limit number error; how many results; limit not accepted

The optional limit must be a whole number above zero, or Turbo says "limit must be a positive integer". With no limit there is none.

## Single visible purchase order lookup

Lines 709-736

Asked as: PO not found or not visible; cannot open purchase order; open single PO error; view one purchase order

Fetching one PO returns it only if it is visible to the caller. Otherwise Turbo says "purchase order not found or not visible", and it does not say whether the PO exists.

## Searchable purchase orders list

Lines 738-754

Asked as: search purchase orders; find PO search; PO search results; search list empty

The search list returns the purchase orders that the caller may see, under the same visibility rules as the other lists, with no scope, filter or limit.

## Rejection reason is required and at least 5 characters

Lines 771-803

Asked as: reject needs reason; rejection reason too short; reason required to reject; reject reason length

Rejecting a PO needs a rejection reason. A missing reason gives "you must provide a rejection reason". A reason shorter than 5 characters after trimming spaces gives "rejection reason must be at least 5 characters".

## Only unapproved, unrejected POs can be rejected

Lines 821-837

Asked as: reject already rejected; can't reject this PO; cannot reject approved; reject status error

A PO that is already rejected cannot be rejected again: "rejected purchase orders cannot be rejected again". A PO whose status is not Unapproved cannot be rejected: "only unapproved purchase orders can be rejected".

## Rejected PO cannot be rejected again; checked before status

Lines 821-828

Asked as: rejected again message; reject twice error; already rejected check; reject order checked first

When someone rejects a purchase order, Turbo first checks whether it is already rejected. If so, it refuses with "rejected purchase orders cannot be rejected again". Only then does it check the status.

## Rejection uses approver policy and bypass

Lines 839-884

Asked as: who can reject; not authorized to reject; reject as approver; reject with bypass

The caller counts as an approver if they are a first-stage approver. They also count if they are the assigned approver and valid through the bypass for the owner. A qualified second approver may also reject. Anyone else gets "you are not authorized to reject this purchase order".

## Who may reject a purchase order

Lines 863-884

Asked as: reject first approved PO; reject before active; reject approved not active; rejection permission

A first-stage approver, a qualified second-stage approver, or the assigned approver who is valid through the bypass may reject. Anyone else gets "you are not authorized to reject this purchase order". A PO that is first-approved but not yet Active can still be rejected by these people.

## Rejection records details and notifies the creator

Lines 886-943

Asked as: rejection notifies creator; reject email failed; rejection saved details; who rejected PO

Rejection saves the rejection time, the reason and the rejector. Turbo notifies the PO creator. If the notification fails, the rejection still stands.

## Cancelling a rejected revision

Lines 968-980

Asked as: cancel rejected revision; delete rejected PO; only creator cancels; rejected revision cancel

The creator of a rejected revision can cancel it even though it is not Active. Anyone else gets "only the creator can cancel a rejected revision".

## Only active purchase orders can be cancelled

Lines 981-1009

Asked as: cancel PO button; cannot cancel PO; payables admin cancel; only active cancelled

Any other PO cannot be cancelled, and Turbo says "only active purchase orders can be cancelled". The PO's creator or a payables admin (the payables_admin claim) may cancel it. Anyone else gets "you are not authorized to cancel this purchase order".

## POs with expenses cannot be cancelled

Lines 1012-1035

Asked as: cancel PO with expenses; expense blocks cancel; cancel revision expenses; expenses stop cancellation

A PO with any expense, committed or not, cannot be cancelled. Turbo says "this purchase order has associated expenses and cannot be cancelled". This applies to revisions too. Cancelling sets the status to Cancelled and records the time and the canceller.

## Cancelling a PO is refused when it has any expenses

Lines 1012-1030

Asked as: linked expenses cancel error; can't cancel attached expenses; expenses attached to PO; cancel refused expenses

Turbo counts the expenses linked to the purchase order. If there is at least one, it refuses with "this purchase order has associated expenses and cannot be cancelled". This applies to the owner and to payables admins.

## Cancelling sets status, time and canceller

Lines 1032-1035

Asked as: cancelled time canceller; who cancelled PO; status after cancel; record of cancellation

A successful cancel sets the PO status to Cancelled. It also records the time (cancelled) and the person who cancelled it (canceller).

## Converting to Cumulative: status and type

Lines 1099-1116

Asked as: make cumulative; convert PO cumulative; one-time to cumulative; change PO type

Only an Active PO, or one closed by the system, can be converted. Imported POs closed by the system are excluded. Otherwise Turbo says "only active or automatically closed purchase orders can be converted to Cumulative". The PO must be One-Time, or Turbo says "only One-Time purchase orders can be converted to Cumulative".

## Auto-closed imported POs cannot convert

Lines 1099-1107

Asked as: closed by system convert; imported PO convert; auto closed cannot convert; convert closed PO

A Closed purchase order qualifies for conversion only if the system closed it (closed_by_system) and it was not imported (_imported). Otherwise the status must be Active. Others get "only active or automatically closed purchase orders can be converted to Cumulative".

## Conversion check order

Lines 1099-1133

Asked as: conversion error order; why convert refused; which conversion check fails; convert check order

Turbo checks the status first, then that the type is One-Time ("only One-Time purchase orders can be converted to Cumulative"), then the payables_admin claim ("you are not authorized to convert this purchase order to Cumulative").

## Converting to Cumulative needs payables admin

Lines 1118-1133

Asked as: payables admin convert; not allowed cumulative; who can make cumulative; convert permission

Only a payables admin (the payables_admin claim) can convert a PO to Cumulative. Anyone else gets "you are not authorized to convert this purchase order to Cumulative".

## Conversion refused if expenses exceed the PO total

Lines 1135-1154

Asked as: expenses more than PO total; convert too many expenses; expenses exceed total; reduce expense before convert

Turbo totals all expenses on the PO, committed or not, and compares whole cents. If that total is above the PO total, conversion is refused with "the expenses of this purchase order total more than the purchase order total, so it can't be Cumulative. Reduce or remove an expense first". A total equal to the PO total passes this check.

## Converting an auto-closed PO reopens it

Lines 1156-1180

Asked as: reopen closed PO; auto closed PO reopen; reactivate closed PO; closed PO back active

For a PO closed by the system, committed expenses must be below the PO total, compared in cents. If they are equal or higher, Turbo says "committed expenses already use the full purchase order total, so there is no room for another invoice". Otherwise the PO becomes Active again, with the closed time and closer cleared and closed_by_system set to false. Its type then becomes Cumulative.

## Conversion needs room for another invoice when reopening

Lines 1156-1177

Asked as: no room for invoice; full PO total used; reopen needs room; committed expenses full

When an auto-closed purchase order is converted, Turbo compares committed expenses to the total in whole cents. If they are equal or greater, it fails with "committed expenses already use the full purchase order total, so there is no room for another invoice". Otherwise it reopens the purchase order: status Active, closed and closer cleared, closed_by_system false.

## Conversion sets the type to Cumulative

Lines 1179-1210

Asked as: PO type changed; cumulative after convert; type set cumulative; converted PO result

After the checks pass, Turbo sets the type to Cumulative and returns the updated purchase order.

## Parent PO number format and sequence

Lines 1214-1359

Asked as: PO number format; how PO numbers generated; unable unique PO number; next PO number

A new parent PO number is YYMM-NNNN, using the current year and month, such as 2401-0001. Turbo takes the highest number of that month below 5000 and adds one, skipping numbers already used. Numbers from 5000 up are reserved for manual or imported POs. After 4999 Turbo cannot generate a number and says "unable to generate a unique PO number".

## Child PO number format and limit

Lines 1229-1296

Asked as: child PO number; 99 children limit; sub PO number; child suffix

A child PO number is its parent's number plus a two-digit suffix, such as 2401-0001-01. The suffix is one more than the highest existing child suffix, starting at 01. Turbo fails if the parent is missing or has no number. It also fails after 99 children, with "maximum number of child POs reached (99)".

## Parent PO numbers skip the 5000 range

Lines 1297-1358

Asked as: PO number 5000; manual number reserved; skip 5000 numbers; numbers above 5000

Turbo numbers a new purchase order as YYMM-NNNN, using the current year and month. It starts after the highest number of the month below 5000 and takes the first free one. Numbers from YYMM-5000 up are reserved for manual or imported numbers. If none below 5000 is free, it fails with "unable to generate a unique PO number".

## Approver lookup input checks

Lines 1361-1493

Asked as: amount must be number; division required approver; has job boolean error; approver lookup inputs

The approver lookup needs a valid amount, or Turbo says "Amount must be a valid number". A has_job value that is not a boolean gives "has_job must be a valid boolean". A missing division gives "division is required". The kind defaults to the capital expenditure kind. An unknown kind gives "kind is invalid or no longer exists".

## Approver lookup for recurring and foreign currency POs

Lines 1432-1473

Asked as: recurring PO approvers; exchange rate missing; foreign currency approver; currency not found

For a Recurring PO, Turbo works out the total value from the amount, dates and frequency, and uses it to find approvers. Bad recurrence input is returned as field errors. An unknown currency gives "currency is invalid or no longer exists". A foreign currency with no exchange rate gives "selected foreign currency is missing an exchange rate". The amount is converted to the home currency before approvers are chosen.

## Recurring approver lookup uses the total value

Lines 1432-1458

Asked as: recurring total calculation; payments count approvers; recurrence field errors; recurring approver total

For a Recurring purchase order, the approver lookup turns the amount, dates and frequency into the number of payments and the total value, and finds approvers for that total. Date or frequency errors come back as field errors (invalid_recurrence). Other errors read "Error calculating recurring PO total".

## Approver lookup converts foreign currency to home currency

Lines 1460-1473

Asked as: exchange rate error approver; currency conversion approvers; foreign currency home total; currency rate zero

The amount is multiplied by the currency rate before the approvers are found. An unknown currency gives "currency is invalid or no longer exists". A foreign currency with no rate above zero gives "selected foreign currency is missing an exchange rate".

## Approver lookup rejects an unknown kind

Lines 1475-1493

Asked as: kind not exist approver; approver lookup kind error; unknown kind lookup; kind removed lookup

If the kind does not exist, the lookup says "kind is invalid or no longer exists".

## Second approver lookup

Lines 1495-1517

Asked as: second approver list empty; requester picked automatically; no second approver list; why second approver list empty

The second approver list holds the second-stage approvers. If the requester is one and second approval is required, the list is empty and the UI picks the requester automatically. If second approval is required, the requester does not qualify and nobody can final-approve the amount, Turbo says "no second-stage approvers can final-approve this amount; contact an administrator". The response includes the recurring payment count and why second approval is or is not needed.

## First approver list includes requester when second approval is needed

Lines 1519-1526

Asked as: requester added approvers; why I'm in approver list; requester first approver list; self added approver

The first approver list holds the first-stage approvers. When the PO needs second approval, the requester is not a first-stage approver, and the requester has a non-zero limit for the kind, the requester is added to the list.
