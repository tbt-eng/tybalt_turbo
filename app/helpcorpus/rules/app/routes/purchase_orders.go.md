<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders.go","sha256":"3dc00a6d8145e2a534009ac2e6c4986024c07c5370424fd296e854366bb4c53d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders.go

## Revised purchase orders are detected by status and number

Lines 23-29

Asked as: revision not showing; PO under revision; edited PO back in list; revised PO pending list

Turbo treats a purchase order as under revision when its status is Unapproved, it already has a PO number, and it is not a legacy manual entry (legacy_manual_entry). The pending and visible PO lists and the job PO queries use this test.

## Approval of a missing purchase order is refused

Lines 272-279

Asked as: cannot approve PO; approve fails error fetching; PO missing when approving; approve button error

If the purchase order cannot be found, Turbo refuses to approve it with the message "error fetching purchase order" followed by the error.

## Only unapproved purchase orders can be approved

Lines 281-297

Asked as: can't approve my PO; approve button not working; approve rejected PO; PO already rejected approve

Approving a purchase order (PO) that is not Unapproved is refused with "only unapproved purchase orders can be approved". A PO that was already rejected is refused with "rejected purchase orders cannot be approved".

## Approving a PO with an unknown kind is refused

Lines 299-323

Asked as: expenditure kind deleted; approve fails kind gone; kind no longer exists; approval kind invalid

When a purchase order's kind (expenditure kind) no longer exists, approval is refused with 'purchase order kind is invalid or no longer exists'. Turbo uses the PO's home-currency approval total, and a PO with no job gets a normalized kind.

## Approval total uses home-currency amount

Lines 299-309

Asked as: which approvers for amount; approvers chosen by total; home currency approvers; approver depends on value

When a purchase order is approved or rejected, Turbo finds the approvers from the home-currency total (EffectiveApprovalTotalHome), the division, the kind and whether the PO has a job. A PO with no job uses its kind normalized by NormalizeExpenditureKindID.

## Invalid PO kind blocks approval or rejection

Lines 303-323

Asked as: reject fails bad kind; can't approve unknown kind; kind removed blocks PO; rejection kind invalid

If the PO's expenditure kind is unknown, approving or rejecting it is refused with "purchase order kind is invalid or no longer exists".

## Assigned approver must be valid for first stage

Lines 333-333

Asked as: approver not valid first stage; wrong approver assigned; approver can't approve amount; assigned approver rejected

If the assigned approver is empty or not a valid first-stage approver for this amount, approval is refused with "assigned approver is not valid for first-stage approval". A bypass route (policy.IsValidFirstStageViaBypass) can make the assigned approver valid.

## PO number assigned on activation

Lines 335-350

Asked as: PO number missing when active; no PO number on approval; number not generated; PO has no number

When a PO becomes Active without a PO number, Turbo generates one. A revised PO that already has a number keeps it.

## Approved PO number is kept if it already exists

Lines 335-339

Asked as: revised PO number changed; PO number reset on approval; number lost after revision; keep PO number

When a PO becomes Active, Turbo sets its status to Active. It generates a new PO number only if the PO has none, so a revised PO keeps its number.

## Approval failure messages for saving and number generation

Lines 340-347

Asked as: error generating PO number; error updating purchase order; approval save failed; PO number error approving

If Turbo cannot generate the PO number on activation, it refuses with "error generating PO number" and the cause. If saving fails, it refuses with "error updating purchase order".

## Second-stage approver can approve both stages at once

Lines 355-371

Asked as: approve both stages; second approver approves once; final approver skips first; approve both at once

If a PO needs second approval and the person approving is a second-stage approver, one action records both approvals and activates the PO. This person need not be the assigned approver. They need a final approval limit that covers the amount, otherwise they get "you do not have sufficient approval limit to complete final approval".

## Final approver's limit must cover the amount

Lines 355-362

Asked as: approval limit too low; not enough limit final approval; my limit insufficient; limit doesn't cover amount

A second-stage approver who completes final approval must have a sufficient approval limit. Otherwise Turbo refuses with 'you do not have sufficient approval limit to complete final approval'. This applies both to approving both stages at once and to second approval of a first-approved PO.

## Dual-approval PO states set approver fields and activation

Lines 363-432

Asked as: approver fields after approval; PO still unapproved after first; two stage approval status; second approver fields set

When a second-stage approver approves both stages at once, Turbo sets approver and second_approver to that person, stamps approved and second_approval, and activates the PO. When a first approver approves a PO that needs second approval, the PO stays Unapproved until second approval.

## Single-approval PO: assigned approver approves and PO becomes Active

Lines 372-432

Asked as: only assigned approver approves; can't approve PO assigned to me; not authorized to approve; single approval PO

When no second approval is required, only the assigned approver (approver) can give the first approval. Approval sets the status to Active and records the approval time and approver. Anyone else gets "you are not authorized to approve this purchase order".

## Only the assigned approver can give first approval

Lines 373-379

Asked as: who can give first approval; not authorized to approve; wrong person approving

On first approval, a caller who is not the assigned approver (approver) and not a second-stage approver of a dual-approval PO is refused with 'you are not authorized to approve this purchase order'.

## Dual approval needs priority second approver and non-empty pools

Lines 381-421

Asked as: second approval pool empty; no priority second approver; first approval pool empty; priority approver not valid

For a PO that needs second approval, the first approval is refused when the second approval pool is empty ("second approval pool is empty for this purchase order; contact an administrator"). It is also refused when no priority second approver is set ("priority second approver is required when second approval is required"). It is refused when the priority second approver is not a valid second-stage approver ("priority second approver is not valid for second-stage approval"). The first approval pool must not be empty unless the assigned approver is valid by bypass ("first approval pool is empty for this purchase order; contact an administrator"). The PO stays Unapproved after first approval.

## Second approval of a first-approved PO

Lines 436-464

Asked as: give final approval; not authorized second approval; second approval after first; final approval of PO

On a first-approved PO that needs second approval, only a second-stage approver can give the final approval. Others get "you are not authorized to perform second approval on this purchase order". The approver's limit must cover the amount, or they get "you do not have sufficient approval limit to complete final approval". Success records the second approver and activates the PO. A PO that already has a second approval is refused with "only unapproved purchase orders can be approved".

## Dual-approval PO: second approver already approved

Lines 436-443

Asked as: second approval already recorded; PO already approved twice; second stage already approved

On the second stage of a PO that needs two approvals, if a second approval is already recorded, Turbo refuses with "only unapproved purchase orders can be approved".

## Second-stage approver needs authorization for second approval

Lines 444-450

Asked as: can't do second approval; not authorized second stage; second approval refused

On the second stage of a PO that needs two approvals, a person who is not a second-stage approver is refused with "you are not authorized to perform second approval on this purchase order".

## First-approved single-stage PO can still be activated

Lines 465-485

Asked as: activate PO after first approval; PO stuck first approved; make first approved PO active; activate single stage PO

A PO that is first-approved but not Active and needs no second approval can be activated by the assigned approver or by any first-stage approver. The assigned approver must still be a valid first-stage approver. Otherwise Turbo refuses with "assigned approver is not valid for first-stage approval" or "you are not authorized to approve this purchase order".

## Approval notifications

Lines 527-580

Asked as: approval email not received; email for second approval; PO active email; approval notification missing

After a first approval that leaves second approval pending and a priority second approver set, that person gets an email asking for approval. POs in a group send one email for the whole group, unless the PO already has a number (a revision), which gets its own email. When a PO becomes Active and the approver is not its creator, the creator gets a "PO active" notification.

## Priority second approver notified after first approval

Lines 527-556

Asked as: priority approver no email; second approver not notified; priority approver notification

After first approval of a PO that needs second approval, the priority second approver gets the 'po_priority_second_approval_required' notification. For a numberless PO in a group, they get one group email instead. A revised PO, which already has a number, gets the single-PO email.

## Group PO sends one email to priority second approver

Lines 527-535

Asked as: one email for group; group PO email; grouped PO notification

When a first approval leaves a grouped PO that has no number waiting for second approval, the priority second approver gets one email for the whole group. A revised PO already has a number, so it gets the single-PO email instead.

## Creator notified when PO becomes active

Lines 557-580

Asked as: creator not notified; PO active notice; didn't get PO email; creator notification missing

When approval makes a PO Active and the approver is not the PO's creator, Turbo sends the creator the 'po_active' notification with the PO number and approver name. The creator gets no notice for their own approval.

## Pending approvals list

Lines 587-608

Asked as: what needs my approval; approvals to do; my pending list empty; items waiting for me

The pending list shows POs the caller can approve now. These are first-stage POs assigned to the caller, second-stage POs where the caller is the priority second approver, and second-stage POs open to eligible second approvers after the timeout.

## Pending and single PO lookups

Lines 610-636

Asked as: pending PO not found; can't open pending PO; PO not visible; PO lookup error

Fetching a pending PO that is not pending for the caller gives "pending purchase order not found". Fetching a PO the caller cannot see gives "purchase order not found or not visible". A missing PO id gives "purchase order id is required".

## Pending PO lookup needs an id

Lines 613-636

Asked as: pending lookup no id; blank PO id; pending PO needs id

Looking up one pending purchase order without an id gives 'purchase order id is required'. If the PO is not pending for the caller, Turbo says 'pending purchase order not found'.

## Single pending PO not found

Lines 631-636

Asked as: PO not in my approvals; missing from pending; approval not showing for me

If the purchase order is not among the caller's pending approvals, Turbo answers "pending purchase order not found".

## Visible PO list scopes and limit

Lines 646-712

Asked as: invalid scope error; PO list filter options; limit error PO list; list filter rejected

The PO list accepts scope all (the default), mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. Any other scope is refused with "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". A limit must be a positive integer, or Turbo says "limit must be a positive integer".

## Visible PO list: default scope and allowed scopes

Lines 646-658

Asked as: default PO list filter; scope not accepted; PO list filter wrong

The visible PO list defaults to scope 'all'. The allowed scopes are all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second and group. Any other value is refused with 'scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group'.

## Stale and expiring lists are for payables staff

Lines 660-700

Asked as: stale POs access; expiring POs can't see; payables only list; who sees stale POs

Only payables staff can list stale or expiring POs, and they must give stale_before or expiring_before. Others get "only payables staff can list stale or expiring purchase orders".

## Stale and expiring lists need payables staff

Lines 660-677

Asked as: stale list permission denied; expiring list not allowed; why can't I see stale

Only payables staff (payables_admin) can list the stale or expiring scope. Others are refused with "only payables staff can list stale or expiring purchase orders".

## Group scope needs a group

Lines 679-685

Asked as: group filter error; scope group no group; group PO list error

Listing POs with scope=group without a group is refused with "group is required when scope=group".

## Scope parameters are required

Lines 679-700

Asked as: stale_before required; expiring_before missing; group required error

Scope group needs group ("group is required when scope=group"). Scope stale needs stale_before and scope expiring needs expiring_before, each with its own "is required" message.

## Stale and expiring scopes need their date parameter

Lines 687-700

Asked as: stale date missing; expiring date required; date needed for stale list

The stale scope requires stale_before ('stale_before is required when scope=stale'). The expiring scope requires expiring_before ('expiring_before is required when scope=expiring').

## Visible PO list limit must be positive

Lines 701-712

Asked as: limit error PO list; list limit number invalid; limit must be number

The optional limit on the visible PO list must be a positive whole number. Otherwise Turbo says 'limit must be a positive integer'. With no limit, the list is not capped.

## Single visible PO lookup

Lines 731-751

Asked as: PO not visible error; open PO by id fails; single PO lookup fails

Looking up one visible purchase order needs an id ('purchase order id is required'). A PO that does not exist or that the caller may not see gives 'purchase order not found or not visible'.

## Searchable purchase orders list

Lines 757-773

Asked as: search purchase orders; find POs I can see; PO search list

The search list returns the purchase orders that the caller can see under the PO visibility rules, with scope "all".

## Rejecting a PO requires a reason of at least 5 characters

Lines 790-822

Asked as: reject reason too short; rejection note length; reason needs characters

Rejecting a PO needs a rejection reason. A missing reason gives "you must provide a rejection reason". A reason shorter than 5 characters after trimming spaces gives "rejection reason must be at least 5 characters".

## Rejection reason must be provided

Lines 790-804

Asked as: reject without reason; no rejection note; rejection reason missing

A rejection request with no rejection_reason field is refused with "you must provide a rejection reason".

## Rejection reason needs 5 characters

Lines 815-822

Asked as: reason with spaces rejected; rejection text too short; trimmed reason too short

After trimming spaces, a rejection reason shorter than 5 characters is refused with "rejection reason must be at least 5 characters".

## Rejecting a missing purchase order

Lines 831-838

Asked as: reject fails not found; error fetching PO reject; PO to reject missing

If the purchase order to reject cannot be found, Turbo refuses with "error fetching purchase order" followed by the error.

## Only unapproved, not-yet-rejected POs can be rejected

Lines 840-856

Asked as: can't reject again; reject approved PO; reject already rejected PO; only unapproved reject

A PO already rejected is refused with "rejected purchase orders cannot be rejected again". A PO that is not Unapproved is refused with "only unapproved purchase orders can be rejected".

## Rejecting a PO with an unknown kind is refused

Lines 858-880

Asked as: reject fails kind; rejection kind deleted; reject unknown kind

Rejection computes the approval policy from the PO's kind. When the kind no longer exists, Turbo refuses with 'purchase order kind is invalid or no longer exists'.

## Who may reject a purchase order

Lines 882-903

Asked as: who can reject PO; reject not authorized; reject first approved PO; who may reject

A first-stage approver, the assigned approver valid by bypass, or a qualified second-stage approver may reject. Others get "you are not authorized to reject this purchase order". A PO already first-approved and awaiting second approval can still be rejected by any of them, because it is not yet Active.

## Rejection records and notifies the creator

Lines 905-962

Asked as: rejection notifies creator; who rejected PO recorded; rejection email creator

Rejecting sets the rejection time, reason and rejector, and notifies the PO creator. A failed notification does not undo the rejection.

## Creator can cancel a rejected revision

Lines 987-999

Asked as: cancel rejected revision; can't cancel rejected PO; remove rejected revision; only creator cancels revision

A rejected revision (Unapproved, has a PO number, not a legacy entry, rejected) can be cancelled only by its creator, whatever its status. Others get "only the creator can cancel a rejected revision".

## Cancelling a PO: who and which status

Lines 1000-1028

Asked as: cancel PO not allowed; who can cancel PO; cancel button missing

Only Active POs can be cancelled ("only active purchase orders can be cancelled"). The PO's creator or payables staff (payables_admin) may cancel. Others get "you are not authorized to cancel this purchase order".

## Cancel: Active POs, owner or payables staff

Lines 1000-1029

Asked as: cancel non-active PO; why cancel fails status; cancel needs active status

Any other PO must be Active to cancel, or Turbo refuses with "only active purchase orders can be cancelled". The creator or payables staff (payables_admin) may cancel. Others are refused with "you are not authorized to cancel this purchase order".

## PO with expenses cannot be cancelled

Lines 1031-1054

Asked as: cancel PO with expenses; PO with expense can't cancel; expenses block cancellation

A PO that has any expense cannot be cancelled. Turbo says "this purchase order has associated expenses and cannot be cancelled". Success sets status Cancelled with the time and canceller.

## Cancellation sets status and canceller

Lines 1051-1054

Asked as: cancelled by who; cancel time recorded; status after cancel

When cancellation succeeds, the PO status becomes Cancelled, and Turbo records the cancel time and the person who cancelled it.

## Convert to Cumulative: eligible status and type

Lines 1118-1135

Asked as: convert to cumulative; one-time to cumulative fails; closed PO convert; cumulative conversion not allowed

Only an Active PO, or one closed automatically by the system, can be converted to Cumulative. Otherwise Turbo says "only active or automatically closed purchase orders can be converted to Cumulative". Imported POs closed by the system are refused. Only One-Time POs convert ("only One-Time purchase orders can be converted to Cumulative").

## Imported closed POs cannot be converted

Lines 1119-1126

Asked as: imported PO convert; imported closed cumulative; closed imported PO fails

A Closed PO counts as automatically closed only when closed_by_system is set and the PO was not imported. Imported POs fail the status check and get the error 'only active or automatically closed purchase orders can be converted to Cumulative'.

## Only payables staff can convert to Cumulative

Lines 1137-1152

Asked as: convert to cumulative not allowed; who can make cumulative; cumulative needs payables

Converting a PO to Cumulative needs the payables staff claim (payables_admin). Others get "you are not authorized to convert this purchase order to Cumulative".

## Conversion blocked if expenses exceed PO total

Lines 1154-1173

Asked as: expenses over PO total; cumulative blocked by expenses; expenses exceed PO total

Conversion is refused when the PO's expenses, committed or not, total more than the PO total, compared in whole cents. A total equal to the PO total is allowed. The message is "the expenses of this purchase order total more than the purchase order total, so it can't be Cumulative. Reduce or remove an expense first".

## Converting an auto-closed PO reopens it

Lines 1175-1196

Asked as: reopen closed PO; auto closed PO convert; no room for invoice; closed PO no room

For a PO closed by the system, committed expenses must be below the PO total, or conversion is refused with "committed expenses already use the full purchase order total, so there is no room for another invoice". Equal counts as used up. On success the PO becomes Active again, and the closed time, closer and closed-by-system flag are cleared.

## Conversion to Cumulative: reopened PO fields

Lines 1175-1199

Asked as: closed PO back active; cumulative reopen clears closed; reopened PO fields

Converting an automatically closed One-Time PO to Cumulative sets it back to Active and clears closed, closer and closed_by_system. The conversion is refused unless committed expenses total less than the PO total ('committed expenses already use the full purchase order total, so there is no room for another invoice'). Amounts are compared in whole cents.

## Child PO number format

Lines 1249-1314

Asked as: child PO number format; PO suffix numbering; parent PO children numbers; child PO number wrong

A child PO's number is its parent's number plus a two-digit suffix, starting at 01 and rising by one after the highest existing child. The parent must exist and have a number. A parent can have at most 99 children.

## Parent PO number format

Lines 1316-1377

Asked as: PO number format; unable to generate PO number; new PO number format; PO numbering rules

Turbo numbers a new parent PO YYMM-NNNN, using the current year and month, and takes the next number after the highest in that month. It auto-generates only numbers below 5000. Numbers from 5000 up are reserved for manual or imported POs. If none is free, Turbo fails with "unable to generate a unique PO number".

## Approver lookup input checks

Lines 1380-1449

Asked as: approver list error; amount must be number; division required approver; approver lookup input error

The approver list needs a numeric amount ("Amount must be a valid number"), a division ("division is required") and a kind ("kind is required"). has_job must be a boolean ("has_job must be a valid boolean"). An unknown kind gives "kind is invalid or no longer exists". The default kind is the capital expenditure kind.

## Approver lookup amount and currency

Lines 1451-1492

Asked as: currency exchange rate missing; recurring PO approver amount; approver amount wrong currency

For recurring POs, Turbo works out the total from the dates and frequency and shows recurrence errors beside those fields. An invalid currency gives "currency is invalid or no longer exists". A foreign currency without a rate gives "selected foreign currency is missing an exchange rate". The amount is converted to home currency before choosing approvers.

## Recurring PO approver lookup uses total value

Lines 1451-1477

Asked as: recurring PO approvers wrong; recurrence error approvers; approvers for recurring total

For a Recurring PO, the approver lookup replaces the amount with the total value over all payments, calculated from start date, end date and frequency. Bad dates or frequency return 'invalid_recurrence' field errors.

## Approver lookup converts foreign currency to home

Lines 1479-1492

Asked as: currency error approvers; exchange rate missing approver; approver foreign currency error

The approver lookup refuses an unknown currency ('currency is invalid or no longer exists'). It refuses a foreign currency without an exchange rate ('selected foreign currency is missing an exchange rate'). Otherwise it multiplies the amount by the currency rate to get the home-currency amount.

## Approver lookup rejects an unknown kind

Lines 1494-1507

Asked as: approver lookup kind error; kind missing for approvers; approvers kind invalid

If the kind does not exist, the approver lookup says 'kind is invalid or no longer exists'.

## Second approver lookup

Lines 1514-1536

Asked as: no second approvers; who can final approve; second approver list empty; no one can final approve

The second approver list returns second-stage approvers. If the requester is one and second approval is required, the list is empty and the UI picks the requester. If second approval is required, the requester does not qualify and no one qualifies, Turbo says "no second-stage approvers can final-approve this amount; contact an administrator".

## First approver list may add the requester

Lines 1538-1545

Asked as: requester added first approvers; I can approve myself; add me as approver; requester first approver list

If a PO needs second approval, the requester is not a first-stage approver but has a non-zero limit for the kind, Turbo adds the requester to the first-approver list.
