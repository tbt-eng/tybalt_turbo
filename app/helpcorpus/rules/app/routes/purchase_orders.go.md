<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders.go","sha256":"996973010d8c0ff97dbedf1a88d12555007cb8019f68c8bd849f2adc70f47ec5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders.go

## Revised purchase orders are identified by status, number and entry type

Lines 23-29

A purchase order is under revision when its status (status) is Unapproved, it already has a PO number (po_number), and it is not a legacy manual entry (legacy_manual_entry = 0). Turbo uses this test in the PO lists and queries.

## Second approval status messages

Lines 179-214

The second approver summary reports one of four states. "not_required" means second approval is not required. "requester_qualifies" means the requester can do it. "candidates_available" means eligible second approvers exist. "required_no_candidates" means no second-stage approver can final-approve the amount. The summary also gives the threshold, limit column, evaluated amount and second-stage timeout hours.

## Only unapproved purchase orders can be approved

Lines 281-297

Approving a purchase order that is not Unapproved is refused with "only unapproved purchase orders can be approved". A purchase order that has been rejected is refused with "rejected purchase orders cannot be approved".

## Rejected purchase orders cannot be approved

Lines 290-297

When a purchase order has a rejection date (rejected), Turbo refuses to approve it with the message "rejected purchase orders cannot be approved". Turbo checks this after the Unapproved status check.

## Approval policy uses home-currency total, kind and job

Lines 299-323

On approve, Turbo works out the approval policy from the purchase order's division, its effective approval total in home currency, its expenditure kind and whether it has a job. The kind is normalized with the has-job flag. An unknown kind is refused with "purchase order kind is invalid or no longer exists".

## Unknown PO kind blocks approval and rejection

Lines 310-317

If the purchase order's kind no longer exists, Turbo refuses approval or rejection with "purchase order kind is invalid or no longer exists".

## Activation generates a PO number if missing

Lines 335-350

When a purchase order becomes Active, it keeps any PO number it already has. If it has none, Turbo generates one. If generating fails, Turbo reports "error generating PO number" with the cause.

## Second approver can approve both stages at once

Lines 355-371

If the purchase order needs second approval and the person approving is a second-stage approver, they can complete both approvals in one action. Turbo records them as approver and second approver, sets the approved and second approval times, and makes the purchase order Active. The approver must have enough final approval limit for the total, or Turbo refuses with "you do not have sufficient approval limit to complete final approval".

## Final approval of a dual-approval PO needs sufficient limit

Lines 355-362

A second-stage approver who approves a purchase order that needs second approval must have a final approval limit that covers the total. If not, Turbo refuses with "you do not have sufficient approval limit to complete final approval". This applies both to combined first-and-second approval and to the second stage.

## First-stage approval only by the assigned approver

Lines 372-379

When a first approval is pending and the combined path does not apply, only the assigned approver (approver) may approve. Anyone else is refused with "you are not authorized to approve this purchase order".

## Second approval needs a non-empty second approver pool

Lines 381-387

On first approval of a purchase order that needs second approval, an empty second-stage pool is refused with "second approval pool is empty for this purchase order; contact an administrator".

## Priority second approver required and must be valid

Lines 389-405

On first approval of a purchase order that needs second approval, a priority second approver must be set, or Turbo refuses with "priority second approver is required when second approval is required". That person must be a valid second-stage approver, or Turbo refuses with "priority second approver is not valid for second-stage approval".

## First-stage pool and assigned approver validity checks

Lines 407-421

If second approval is required, the first-stage pool is empty and the assigned approver is not valid through the bypass, Turbo refuses with "first approval pool is empty for this purchase order; contact an administrator". An assigned approver that is missing or not valid for the first stage is refused with "assigned approver is not valid for first-stage approval".

## First approval activates single-stage purchase orders

Lines 423-432

A valid first approval sets the approved time and the approver. If no second approval is needed, the purchase order becomes Active at once. If second approval is needed, it stays Unapproved and waits for the second approver.

## Final approval of first-approved dual-approval purchase orders

Lines 436-464

For a first-approved purchase order that needs second approval, a second approval already recorded is refused with "only unapproved purchase orders can be approved". Only a second-stage approver may give it, or Turbo refuses with "you are not authorized to perform second approval on this purchase order". The approver needs enough final limit, or Turbo refuses with "you do not have sufficient approval limit to complete final approval". Success records the second approver and time and activates the purchase order.

## Second approval cannot be repeated

Lines 436-443

If a purchase order is already first-approved and already second-approved, approving it again fails with "only unapproved purchase orders can be approved".

## Only second-stage approvers can give second approval

Lines 444-450

On the second stage of a dual-approval purchase order, a caller who is not a second-stage approver is refused with "you are not authorized to perform second approval on this purchase order".

## Second approval sets second approver and activates the PO

Lines 459-464

A successful second approval sets second_approval to now and second_approver to the caller. The status becomes Active.

## Activating a first-approved single-stage purchase order

Lines 465-485

For a first-approved purchase order that needs no second approval but is not yet Active, the assigned approver must be a valid first-stage approver, or Turbo refuses with "assigned approver is not valid for first-stage approval". The caller must be the assigned approver or a first-stage approver, or Turbo refuses with "you are not authorized to approve this purchase order". Success activates it.

## Approval notifications

Lines 527-580

After a first approval that leaves the purchase order awaiting second approval with a priority second approver, Turbo notifies that person. A purchase order in a group with no PO number gets one email for the whole group. A revised purchase order, which already has a number, gets the single email. When a purchase order becomes Active and the approver is not its creator, Turbo notifies the creator that it is active.

## Priority second approver is emailed after first approval

Lines 527-556

After first approval of a PO that needs second approval and has a priority second approver, Turbo emails that person (template po_priority_second_approval_required). A PO in a group with no PO number yet sends one email for the whole group. A revised PO, which already has a number, gets the single-PO email.

## Creator is notified when PO becomes active

Lines 557-580

When an approval makes the PO Active, Turbo sends the creator the po_active email. It sends none if the creator is the person who approved.

## Pending purchase orders list

Lines 587-640

The pending list returns the purchase orders the caller can approve now, using the PO visibility rules. These are first-stage ones assigned to the caller, second-stage ones where the caller is the priority second approver, and second-stage ones open to all eligible second approvers after the second-stage timeout. A single pending purchase order that is not found for the caller is refused with "pending purchase order not found". A missing id is refused with "purchase order id is required".

## Pending PO lookup by id

Lines 610-640

Fetching one pending purchase order needs an id, otherwise Turbo says "purchase order id is required". If the PO is not in the caller's pending set, Turbo says "pending purchase order not found".

## Visible purchase orders scope and parameters

Lines 646-658

The visible purchase orders list takes a scope of all (the default), mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. Any other value is refused with "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". The scope is lowercased and trimmed.

## Scope-specific required parameters

Lines 660-681

Scope group requires group, refused with "group is required when scope=group". Scope stale requires stale_before, refused with "stale_before is required when scope=stale". Scope expiring requires expiring_before, refused with "expiring_before is required when scope=expiring".

## Limit must be a positive integer

Lines 682-693

The optional limit on the visible purchase orders list must be a positive integer. Zero, negative or non-numeric values are refused with "limit must be a positive integer". With no limit, the list is not limited.

## Single visible purchase order lookup

Lines 709-736

Fetching one purchase order requires an id ("purchase order id is required"). A purchase order that does not exist or is not visible to the caller gives "purchase order not found or not visible".

## Searchable purchase orders list

Lines 738-754

The search list returns the purchase orders the caller may see under the PO visibility rules.

## Rejection reason is required

Lines 771-803

Rejecting a purchase order needs a rejection_reason in the request. If the field is missing, Turbo refuses with "you must provide a rejection reason". If the trimmed reason is shorter than 5 characters, Turbo refuses with "rejection reason must be at least 5 characters".

## Rejection only for unapproved, not-yet-rejected purchase orders

Lines 821-837

A purchase order that is already rejected is refused with "rejected purchase orders cannot be rejected again". One that is not Unapproved is refused with "only unapproved purchase orders can be rejected".

## Who can reject a purchase order

Lines 841-884

The caller must be a first-stage approver, the assigned approver valid through the bypass, or a qualified second-stage approver. Anyone else is refused with "you are not authorized to reject this purchase order". A purchase order that is already first-approved but not yet Active can still be rejected by any such person.

## Who may reject a purchase order

Lines 863-884

The caller must be a first-stage approver, the assigned approver through the bypass rule, or a qualified second-stage approver. Otherwise Turbo says "you are not authorized to reject this purchase order". A first-approved PO that is not yet Active can still be rejected by any of these.

## Rejection records and notifies the creator

Lines 886-943

Rejecting records the time, reason and rejector, and notifies the creator. If the notification fails, the rejection still stands.

## Creator can cancel a rejected revision

Lines 968-980

A rejected revision of a purchase order can be cancelled by its creator instead of being resubmitted, and the Active and payables admin checks do not apply. Anyone else is refused with "only the creator can cancel a rejected revision".

## Who can cancel a purchase order

Lines 981-1006

Only an Active purchase order can be cancelled, or Turbo refuses with "only active purchase orders can be cancelled". The caller needs the payables admin claim (payables_admin), or Turbo refuses with "you are not authorized to cancel this purchase order".

## Purchase orders with expenses cannot be cancelled

Lines 1009-1032

A purchase order that has any expenses is refused with "this purchase order has associated expenses and cannot be cancelled". Otherwise cancelling sets the status to Cancelled and records the time and the canceller.

## Conversion to Cumulative requires Active or system-closed status

Lines 1064-1102

Only an Active purchase order, or a Closed one that the system closed, can be converted to Cumulative. Imported purchase orders closed by the system are refused. Others are refused with "only active or automatically closed purchase orders can be converted to Cumulative".

## Imported system-closed POs cannot be converted

Lines 1094-1102

A Closed PO converts only if closed_by_system is set and it is not imported (_imported). Otherwise Turbo says "only active or automatically closed purchase orders can be converted to Cumulative".

## Only One-Time purchase orders convert, by payables admin

Lines 1104-1128

Only a One-Time purchase order can be converted, or Turbo refuses with "only One-Time purchase orders can be converted to Cumulative". The caller needs the payables admin claim (payables_admin), or Turbo refuses with "you are not authorized to convert this purchase order to Cumulative".

## Reopening a system-closed purchase order on conversion

Lines 1130-1154

Converting a system-closed purchase order reopens it as Active and clears its closed time and closer. Committed expenses are compared with the total in whole cents. If they equal or exceed the total, Turbo refuses with "committed expenses already use the full purchase order total, so there is no room for another invoice". On success the type becomes Cumulative.

## Child purchase order numbers

Lines 1203-1270

A child purchase order's number is its parent's number plus a two-digit suffix, YYMM-NNNN-XX, one above the highest existing child. A missing parent fails with "parent PO not found". A parent without a number fails with "parent PO does not have a PO number". After suffix 99 it fails with "maximum number of child POs reached (99) for parent" and the parent number. A number that already exists also fails.

## Parent purchase order numbers

Lines 1271-1332

A parent purchase order number is YYMM-NNNN, using the current year and month, one above the highest existing number for that month. Numbers from 5000 up are reserved for manual or imported purchase orders, so automatic numbers stop at 4999. When none is free, it fails with "unable to generate a unique PO number".

## Approver lookup request validation

Lines 1335-1404

The approvers lookup requires a valid numeric amount ("Amount must be a valid number"), a division ("division is required") and a kind ("kind is required"). has_job, if given, must be a valid boolean ("has_job must be a valid boolean"). The kind defaults to the default capital expenditure kind and is normalized with has_job.

## Approver lookup query validation

Lines 1335-1403

The approver lookup needs a valid amount ("Amount must be a valid number"), a valid has_job if given ("has_job must be a valid boolean"), a division ("division is required") and a kind ("kind is required"). When kind is not given, the default capital expenditure kind is used.

## Recurring amount in approver lookup

Lines 1406-1432

For a Recurring purchase order, the approvers lookup replaces the amount with the total value over all payments, calculated from the start date, end date and frequency. Date or frequency errors are returned as "invalid_recurrence" field errors. Other errors give "Error calculating recurring PO total" with the cause.

## Approver lookup converts currency to home

Lines 1434-1467

The approvers lookup converts the amount to home currency using the currency's exchange rate. An invalid currency is refused with "currency is invalid or no longer exists". A foreign currency with no positive rate is refused with "selected foreign currency is missing an exchange rate". An unknown kind is refused with "kind is invalid or no longer exists".

## Second approver candidates lookup

Lines 1469-1491

The second approver lookup returns the second-stage approvers for the amount, plus a status summary. If the requester can give second approval and it is required, the list is empty so the form selects the requester. If it is required, the requester does not qualify and nobody can final-approve the amount, Turbo refuses with "no second-stage approvers can final-approve this amount; contact an administrator".

## First approver candidates include the requester

Lines 1493-1500

The first approver lookup returns the first-stage approvers. When second approval is required, the requester is not a first-stage approver but has a non-zero limit for the kind, and is eligible, they are added to the list.
