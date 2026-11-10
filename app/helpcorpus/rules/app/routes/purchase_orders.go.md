<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders.go","sha256":"996973010d8c0ff97dbedf1a88d12555007cb8019f68c8bd849f2adc70f47ec5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders.go

## Revised purchase orders are identified by status, number and entry type

Lines 23-29

Asked as: revised purchase order; PO being edited again; resubmitted PO in list; unapproved PO with number

A purchase order is under revision when its status (status) is Unapproved, it already has a PO number (po_number), and it is not a legacy manual entry (legacy_manual_entry = 0). Turbo uses this test in the PO lists and queries.

## Second approval status messages

Lines 179-214

Asked as: second approver status; why second approval needed; second approval summary message; what second approval state means

The second approver summary reports one of four states. "not_required" means second approval is not required. "requester_qualifies" means the requester can do it. "candidates_available" means eligible second approvers exist. "required_no_candidates" means no second-stage approver can final-approve the amount. The summary also gives the threshold, limit column, evaluated amount and second-stage timeout hours.

## Only unapproved purchase orders can be approved

Lines 281-297

Asked as: can't approve PO; approve button not working; PO already approved error; cannot approve purchase order

Approving a purchase order that is not Unapproved is refused with "only unapproved purchase orders can be approved". A purchase order that has been rejected is refused with "rejected purchase orders cannot be approved".

## Rejected purchase orders cannot be approved

Lines 290-297

Asked as: approve rejected PO; rejected PO won't approve; bring back rejected purchase order; approving a rejected order

When a purchase order has a rejection date (rejected), Turbo refuses to approve it with the message "rejected purchase orders cannot be approved". Turbo checks this after the Unapproved status check.

## Approval policy uses home-currency total, kind and job

Lines 299-323

Asked as: which approval policy applies; approval limits by division; home currency approval total; approval rules for PO

On approve, Turbo works out the approval policy from the purchase order's division, its effective approval total in home currency, its expenditure kind and whether it has a job. The kind is normalized with the has-job flag. An unknown kind is refused with "purchase order kind is invalid or no longer exists".

## Unknown PO kind blocks approval and rejection

Lines 310-317

Asked as: PO kind deleted; invalid expenditure kind; kind no longer exists; can't approve after kind removed

If the purchase order's kind no longer exists, Turbo refuses approval or rejection with "purchase order kind is invalid or no longer exists".

## Activation generates a PO number if missing

Lines 335-350

Asked as: PO number missing on activation; no PO number after active; error generating PO number; activated PO without number

When a purchase order becomes Active, it keeps any PO number it already has. If it has none, Turbo generates one. If generating fails, Turbo reports "error generating PO number" with the cause.

## Second approver can approve both stages at once

Lines 355-371

Asked as: approve both stages at once; approve and second approve together; combined approval; I'm second approver approving alone

If the purchase order needs second approval and the person approving is a second-stage approver, they can complete both approvals in one action. Turbo records them as approver and second approver, sets the approved and second approval times, and makes the purchase order Active. The approver must have enough final approval limit for the total, or Turbo refuses with "you do not have sufficient approval limit to complete final approval".

## Final approval of a dual-approval PO needs sufficient limit

Lines 355-362

Asked as: not enough approval limit; can't final approve over limit; approval limit too low; final approval limit error

A second-stage approver who approves a purchase order that needs second approval must have a final approval limit that covers the total. If not, Turbo refuses with "you do not have sufficient approval limit to complete final approval". This applies both to combined first-and-second approval and to the second stage.

## First-stage approval only by the assigned approver

Lines 372-379

Asked as: not authorized to approve; someone else approving my PO; wrong approver can't approve; only assigned approver approves

When a first approval is pending and the combined path does not apply, only the assigned approver (approver) may approve. Anyone else is refused with "you are not authorized to approve this purchase order".

## Second approval needs a non-empty second approver pool

Lines 381-387

Asked as: no second approvers available; second approval list empty; second approver pool empty; contact admin no second approver

On first approval of a purchase order that needs second approval, an empty second-stage pool is refused with "second approval pool is empty for this purchase order; contact an administrator".

## Priority second approver required and must be valid

Lines 389-405

Asked as: priority second approver required; pick second approver; invalid second approver chosen; second approver not allowed

On first approval of a purchase order that needs second approval, a priority second approver must be set, or Turbo refuses with "priority second approver is required when second approval is required". That person must be a valid second-stage approver, or Turbo refuses with "priority second approver is not valid for second-stage approval".

## First-stage pool and assigned approver validity checks

Lines 407-421

Asked as: first approval pool empty; assigned approver not valid; approver not valid first stage; first stage approver missing

If second approval is required, the first-stage pool is empty and the assigned approver is not valid through the bypass, Turbo refuses with "first approval pool is empty for this purchase order; contact an administrator". An assigned approver that is missing or not valid for the first stage is refused with "assigned approver is not valid for first-stage approval".

## First approval activates single-stage purchase orders

Lines 423-432

Asked as: first approval sets approved; PO active after approval; waiting for second approver; approved but still unapproved

A valid first approval sets the approved time and the approver. If no second approval is needed, the purchase order becomes Active at once. If second approval is needed, it stays Unapproved and waits for the second approver.

## Final approval of first-approved dual-approval purchase orders

Lines 436-464

Asked as: give second approval; final approval after first; second approval already done; final approve dual PO

For a first-approved purchase order that needs second approval, a second approval already recorded is refused with "only unapproved purchase orders can be approved". Only a second-stage approver may give it, or Turbo refuses with "you are not authorized to perform second approval on this purchase order". The approver needs enough final limit, or Turbo refuses with "you do not have sufficient approval limit to complete final approval". Success records the second approver and time and activates the purchase order.

## Second approval cannot be repeated

Lines 436-443

Asked as: approve twice; approve again already approved; double approval error; second approval repeated

If a purchase order is already first-approved and already second-approved, approving it again fails with "only unapproved purchase orders can be approved".

## Only second-stage approvers can give second approval

Lines 444-450

Asked as: not a second approver; can't give second approval; who can second approve; second stage access denied

On the second stage of a dual-approval purchase order, a caller who is not a second-stage approver is refused with "you are not authorized to perform second approval on this purchase order".

## Second approval sets second approver and activates the PO

Lines 459-464

Asked as: second approval activates PO; second approver recorded; PO active after second approval; second approval time set

A successful second approval sets second_approval to now and second_approver to the caller. The status becomes Active.

## Activating a first-approved single-stage purchase order

Lines 465-485

Asked as: approved but not active; activate approved PO; stuck approved purchase order; single stage approved not activated

For a first-approved purchase order that needs no second approval but is not yet Active, the assigned approver must be a valid first-stage approver, or Turbo refuses with "assigned approver is not valid for first-stage approval". The caller must be the assigned approver or a first-stage approver, or Turbo refuses with "you are not authorized to approve this purchase order". Success activates it.

## Approval notifications

Lines 527-580

Asked as: approval emails; who gets notified on approval; group email for PO; approval notification sent

After a first approval that leaves the purchase order awaiting second approval with a priority second approver, Turbo notifies that person. A purchase order in a group with no PO number gets one email for the whole group. A revised purchase order, which already has a number, gets the single email. When a purchase order becomes Active and the approver is not its creator, Turbo notifies the creator that it is active.

## Priority second approver is emailed after first approval

Lines 527-556

Asked as: email to second approver; second approver didn't get email; priority approver notification; second approval email missing

After first approval of a PO that needs second approval and has a priority second approver, Turbo emails that person (template po_priority_second_approval_required). A PO in a group with no PO number yet sends one email for the whole group. A revised PO, which already has a number, gets the single-PO email.

## Creator is notified when PO becomes active

Lines 557-580

Asked as: creator email PO active; notify creator; told PO is active; creator not notified

When an approval makes the PO Active, Turbo sends the creator the po_active email. It sends none if the creator is the person who approved.

## Pending purchase orders list

Lines 587-640

Asked as: my pending approvals; POs waiting for me; to approve list; pending purchase orders list

The pending list returns the purchase orders the caller can approve now, using the PO visibility rules. These are first-stage ones assigned to the caller, second-stage ones where the caller is the priority second approver, and second-stage ones open to all eligible second approvers after the second-stage timeout. A single pending purchase order that is not found for the caller is refused with "pending purchase order not found". A missing id is refused with "purchase order id is required".

## Pending PO lookup by id

Lines 610-640

Asked as: find pending PO by id; pending PO not found; PO id required; look up one pending order

Fetching one pending purchase order needs an id, otherwise Turbo says "purchase order id is required". If the PO is not in the caller's pending set, Turbo says "pending purchase order not found".

## Visible purchase orders scope and parameters

Lines 646-658

Asked as: invalid scope; PO list filter options; scope value error; filter purchase orders by scope

The visible purchase orders list takes a scope of all (the default), mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. Any other value is refused with "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". The scope is lowercased and trimmed.

## Scope-specific required parameters

Lines 660-681

Asked as: group parameter missing; stale before date needed; expiring before date required; scope needs date

Scope group requires group, refused with "group is required when scope=group". Scope stale requires stale_before, refused with "stale_before is required when scope=stale". Scope expiring requires expiring_before, refused with "expiring_before is required when scope=expiring".

## Limit must be a positive integer

Lines 682-693

Asked as: limit positive number; PO list limit error; limit must be number; how many POs to show

The optional limit on the visible purchase orders list must be a positive integer. Zero, negative or non-numeric values are refused with "limit must be a positive integer". With no limit, the list is not limited.

## Single visible purchase order lookup

Lines 709-736

Asked as: PO not visible; can't open purchase order; purchase order not found; view single PO

Fetching one purchase order requires an id ("purchase order id is required"). A purchase order that does not exist or is not visible to the caller gives "purchase order not found or not visible".

## Searchable purchase orders list

Lines 738-754

Asked as: search purchase orders; which POs can I see; PO search results; list POs I may view

The search list returns the purchase orders the caller may see under the PO visibility rules.

## Rejection reason is required

Lines 771-803

Asked as: reject without reason; rejection note too short; reason missing when rejecting; rejection reason needed

Rejecting a purchase order needs a rejection_reason in the request. If the field is missing, Turbo refuses with "you must provide a rejection reason". If the trimmed reason is shorter than 5 characters, Turbo refuses with "rejection reason must be at least 5 characters".

## Rejection only for unapproved, not-yet-rejected purchase orders

Lines 821-837

Asked as: reject twice; cannot reject PO; already rejected order; reject approved PO error

A purchase order that is already rejected is refused with "rejected purchase orders cannot be rejected again". One that is not Unapproved is refused with "only unapproved purchase orders can be rejected".

## Who can reject a purchase order

Lines 841-884

Asked as: reject permission; can't reject PO; not allowed to reject; who is allowed to reject

The caller must be a first-stage approver, the assigned approver valid through the bypass, or a qualified second-stage approver. Anyone else is refused with "you are not authorized to reject this purchase order". A purchase order that is already first-approved but not yet Active can still be rejected by any such person.

## Who may reject a purchase order

Lines 863-884

Asked as: rejecting first approved PO; reject approved not active PO; reject access for approvers; reject after first approval

The caller must be a first-stage approver, the assigned approver through the bypass rule, or a qualified second-stage approver. Otherwise Turbo says "you are not authorized to reject this purchase order". A first-approved PO that is not yet Active can still be rejected by any of these.

## Rejection records and notifies the creator

Lines 886-943

Asked as: rejected PO email creator; rejection notification failed; who rejected PO; record of rejection

Rejecting records the time, reason and rejector, and notifies the creator. If the notification fails, the rejection still stands.

## Creator can cancel a rejected revision

Lines 968-980

Asked as: cancel rejected revision; withdraw rejected PO; delete my rejected change; creator cancel rejected edit

A rejected revision of a purchase order can be cancelled by its creator instead of being resubmitted, and the Active and payables admin checks do not apply. Anyone else is refused with "only the creator can cancel a rejected revision".

## Who can cancel a purchase order

Lines 981-1006

Asked as: cancel PO button; cancel active purchase order; not allowed to cancel; who can cancel PO

Only an Active purchase order can be cancelled, or Turbo refuses with "only active purchase orders can be cancelled". The caller needs the payables admin claim (payables_admin), or Turbo refuses with "you are not authorized to cancel this purchase order".

## Purchase orders with expenses cannot be cancelled

Lines 1009-1032

Asked as: cancel PO with expenses; PO has expenses can't cancel; cancel blocked by expense; expenses stop cancellation

A purchase order that has any expenses is refused with "this purchase order has associated expenses and cannot be cancelled". Otherwise cancelling sets the status to Cancelled and records the time and the canceller.

## Conversion to Cumulative requires Active or system-closed status

Lines 1064-1102

Asked as: convert to cumulative; make PO cumulative; cumulative conversion not allowed; change PO to cumulative type

Only an Active purchase order, or a Closed one that the system closed, can be converted to Cumulative. Imported purchase orders closed by the system are refused. Others are refused with "only active or automatically closed purchase orders can be converted to Cumulative".

## Imported system-closed POs cannot be converted

Lines 1094-1102

Asked as: imported PO convert; system closed PO cumulative; auto closed PO convert; closed PO won't convert

A Closed PO converts only if closed_by_system is set and it is not imported (_imported). Otherwise Turbo says "only active or automatically closed purchase orders can be converted to Cumulative".

## Only One-Time purchase orders convert, by payables admin

Lines 1104-1128

Asked as: one-time PO convert; payables admin convert; who can convert PO; convert one-time order

Only a One-Time purchase order can be converted, or Turbo refuses with "only One-Time purchase orders can be converted to Cumulative". The caller needs the payables admin claim (payables_admin), or Turbo refuses with "you are not authorized to convert this purchase order to Cumulative".

## Reopening a system-closed purchase order on conversion

Lines 1130-1154

Asked as: reopen closed PO; no room for invoice; committed expenses full total; reopen PO on conversion

Converting a system-closed purchase order reopens it as Active and clears its closed time and closer. Committed expenses are compared with the total in whole cents. If they equal or exceed the total, Turbo refuses with "committed expenses already use the full purchase order total, so there is no room for another invoice". On success the type becomes Cumulative.

## Child purchase order numbers

Lines 1203-1270

Asked as: child PO number; PO suffix numbering; sub PO number format; parent PO number suffix

A child purchase order's number is its parent's number plus a two-digit suffix, YYMM-NNNN-XX, one above the highest existing child. A missing parent fails with "parent PO not found". A parent without a number fails with "parent PO does not have a PO number". After suffix 99 it fails with "maximum number of child POs reached (99) for parent" and the parent number. A number that already exists also fails.

## Parent purchase order numbers

Lines 1271-1332

Asked as: next PO number; PO number sequence; no unique PO number; how PO numbers are assigned

A parent purchase order number is YYMM-NNNN, using the current year and month, one above the highest existing number for that month. Numbers from 5000 up are reserved for manual or imported purchase orders, so automatic numbers stop at 4999. When none is free, it fails with "unable to generate a unique PO number".

## Approver lookup request validation

Lines 1335-1404

Asked as: approvers lookup amount invalid; approver lookup missing division; has job boolean error; approvers request validation

The approvers lookup requires a valid numeric amount ("Amount must be a valid number"), a division ("division is required") and a kind ("kind is required"). has_job, if given, must be a valid boolean ("has_job must be a valid boolean"). The kind defaults to the default capital expenditure kind and is normalized with has_job.

## Approver lookup query validation

Lines 1335-1403

Asked as: approver search query error; find approvers missing kind; approver query amount; approver lookup parameters

The approver lookup needs a valid amount ("Amount must be a valid number"), a valid has_job if given ("has_job must be a valid boolean"), a division ("division is required") and a kind ("kind is required"). When kind is not given, the default capital expenditure kind is used.

## Recurring amount in approver lookup

Lines 1406-1432

Asked as: recurring PO total; recurring approver amount; invalid recurrence dates; recurring total for approvers

For a Recurring purchase order, the approvers lookup replaces the amount with the total value over all payments, calculated from the start date, end date and frequency. Date or frequency errors are returned as "invalid_recurrence" field errors. Other errors give "Error calculating recurring PO total" with the cause.

## Approver lookup converts currency to home

Lines 1434-1467

Asked as: exchange rate missing; foreign currency approver lookup; currency invalid approver; approver lookup currency error

The approvers lookup converts the amount to home currency using the currency's exchange rate. An invalid currency is refused with "currency is invalid or no longer exists". A foreign currency with no positive rate is refused with "selected foreign currency is missing an exchange rate". An unknown kind is refused with "kind is invalid or no longer exists".

## Second approver candidates lookup

Lines 1469-1491

Asked as: no second approvers for amount; list of second approvers empty; cannot final approve amount; second approver candidates

The second approver lookup returns the second-stage approvers for the amount, plus a status summary. If the requester can give second approval and it is required, the list is empty so the form selects the requester. If it is required, the requester does not qualify and nobody can final-approve the amount, Turbo refuses with "no second-stage approvers can final-approve this amount; contact an administrator".

## First approver candidates include the requester

Lines 1493-1500

Asked as: requester as first approver; I am not in approver list; approver list missing me; add myself as first approver

The first approver lookup returns the first-stage approvers. When second approval is required, the requester is not a first-stage approver but has a non-zero limit for the kind, and is eligible, they are added to the list.
