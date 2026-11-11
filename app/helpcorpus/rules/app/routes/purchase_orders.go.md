<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders.go","sha256":"ae6df821be2a8c9d8e9d04b5fe936061c435818d6f05cc6dc4379cfc9f3bf99a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders.go

## Revised purchase order is recognised by status, number and legacy flag

Lines 23-29

Asked as: PO shows as under revision; revised PO not in list; unapproved PO with number; revision flag on purchase order

Turbo treats a purchase order as under revision when its status is Unapproved, it already has a PO number (po_number) and it is not a legacy manual entry (legacy_manual_entry = 0). Lists and job summaries use this test.

## Second approver lookup fails when nobody can final-approve

Lines 179-214

Asked as: no second approver available; can't find second approver; no one can final approve; second approver error message

If second approval is required, the requester does not qualify and no second-stage approver exists, the lookup fails with "no second-stage approvers can final-approve this amount; contact an administrator". The result status is otherwise not_required, requester_qualifies or candidates_available.

## Second approver lookup returns the timeout

Lines 185-196

Asked as: second approver timeout hours; how many recurring payments; second approval threshold lookup; approver lookup result details

The second approver lookup returns the PO's second-stage timeout in hours, the second approval threshold and the number of recurring payments. For a non-recurring PO the payment count is 0.

## Second-stage timeout is shown with second approver candidates

Lines 216-234

Asked as: when do second approvers see PO; second stage timeout visibility; PO hidden from second approvers; second approver candidates timing

The visibility lists use the second-stage timeout hours (GetPurchaseOrderSecondStageTimeoutHours) to decide when general second approvers see a first-approved PO. The value comes from another file.

## Single purchase order lookup returns nothing when not visible

Lines 236-255

Asked as: purchase order lookup returns nothing; can't open PO by id; PO id returns empty; cannot see purchase order

Looking up a purchase order by id returns nothing, both when it does not exist and when the caller may not see it under the PO visibility rules.

## Approving a purchase order that does not exist is refused

Lines 271-279

Asked as: approve PO that doesn't exist; approval fails po not found; error fetching purchase order approve; approve missing purchase order

When staff approve a purchase order and Turbo cannot find it, Turbo refuses the approval. The message starts with "error fetching purchase order:" and then gives the lookup error (code po_not_found).

## Only unapproved purchase orders can be approved

Lines 281-297

Asked as: can't approve rejected PO; approve PO already active; why won't PO approve; only unapproved can be approved

Approving a purchase order whose status is not Unapproved fails with "only unapproved purchase orders can be approved". A purchase order that has been rejected fails with "rejected purchase orders cannot be approved". The Unapproved status is checked first.

## Approval fails when the purchase order kind is invalid

Lines 299-323

Asked as: PO kind invalid approval; approval fails kind missing; kind no longer exists approve; approve PO kind error

If the purchase order's kind no longer exists, approval fails with "purchase order kind is invalid or no longer exists". Turbo normalizes the kind from whether the PO has a job. It uses the approval total in home currency (EffectiveApprovalTotalHome) to find the approvers.

## PO number is generated at activation only if missing

Lines 335-350

Asked as: PO number not changing on activation; revised PO new number; PO number when activated; keep PO number revision

When a PO becomes Active, Turbo generates a PO number only if the PO has none. A revised PO keeps its existing number.

## Second-stage approver can approve both stages at once

Lines 353-371

Asked as: approve both stages at once; second approver approves everything; combined approval single action; skip first approval as second approver

On a PO that needs second approval and has no first approval yet, a second-stage approver (the caller) approves both stages in one action. Turbo records them as approver and second approver, sets both approval times, and makes the PO Active. A PO number is generated if it has none.

## Combined approval needs a sufficient final limit

Lines 355-362

Asked as: not enough limit to approve; insufficient approval limit final; limit too low for PO; approval limit error combined

A second-stage approver who approves both stages at once needs a limit that covers the approval total. Otherwise approval fails with "you do not have sufficient approval limit to complete final approval". The same message applies to the second approval step.

## Only the assigned approver can give first approval

Lines 372-379

Asked as: not authorised to approve PO; not the assigned approver; can't give first approval; first approval permission denied

For a first approval that does not use the combined path, the caller must be the assigned approver on the PO. Anyone else gets "you are not authorized to approve this purchase order".

## First approval checks pools and the priority second approver

Lines 381-413

Asked as: second approval pool empty; no priority second approver; first approval pool empty; priority approver not valid

When the PO needs second approval, first approval runs these checks in order. An empty second-approver pool fails with "second approval pool is empty for this purchase order; contact an administrator". A missing priority second approver fails with "priority second approver is required when second approval is required". A priority second approver who is not valid for the second stage fails with "priority second approver is not valid for second-stage approval". An empty first-approver pool fails with "first approval pool is empty for this purchase order; contact an administrator", unless the assigned approver is valid for the first stage through the bypass.

## Assigned approver must be valid for first stage

Lines 415-421

Asked as: assigned approver invalid first stage; assigned approver not valid; owner bypass approver; approver not eligible first stage

If the assigned approver is empty or not valid for first-stage approval, approval fails with "assigned approver is not valid for first-stage approval". A valid assigned approver can also qualify through the bypass rule for the PO owner.

## First approval activates single-approval purchase orders

Lines 423-432

Asked as: first approval makes PO active; single approval PO activated; PO still waiting after first approval; approval time recorded first stage

A successful first approval records the approval time and the approver. If the PO does not need second approval, it becomes Active and gets a PO number if it has none. If second approval is needed, it stays Unapproved and waits for the second approval.

## Second approval completes a dual-approval purchase order

Lines 435-464

Asked as: final approval of dual PO; second approval completes PO; already second approved error; second stage approve not authorised

A PO that already has its first approval and needs second approval is completed by a second-stage approver whose limit covers the amount. The PO becomes Active. A PO that already has its second approval fails with "only unapproved purchase orders can be approved". A caller who is not a second-stage approver gets "you are not authorized to perform second approval on this purchase order". A caller whose limit is too low gets "you do not have sufficient approval limit to complete final approval".

## Approving a first-approved single-stage purchase order

Lines 465-485

Asked as: approve first approved single PO; PO stuck after first approval; reapprove single stage PO; assigned approver no longer first stage

A single-stage PO that was first-approved but never activated can be approved again. The assigned approver must still be a first-stage approver, or approval fails with "assigned approver is not valid for first-stage approval". The caller must be the assigned approver or any first-stage approver. Otherwise approval fails with "you are not authorized to approve this purchase order". Approval then makes the PO Active.

## Approval saves the updated purchase order and returns it

Lines 488-524

Asked as: what is returned after approving; approved PO saved details; save approval transaction; PO details after approve

After a successful approval, Turbo saves the purchase order in one transaction. It then returns the updated record with the creator, approver, division and job details.

## Priority second approver is notified after first approval

Lines 527-556

Asked as: who gets email after first approval; priority approver notification email; second approver emailed; group email approval notice

After a first approval that leaves the PO waiting for second approval, Turbo notifies the priority second approver. A PO in a group with no number yet sends one email for the whole group. A revised PO already has a number, so it sends a single-PO email instead.

## Creator is notified when a purchase order becomes active

Lines 557-580

Asked as: creator notified PO active; email when PO approved; who is told PO approved; no email to approver

When an approval makes the PO Active, Turbo notifies the creator. It sends no notification if the creator is the one who approved.

## Pending purchase orders list

Lines 587-640

Asked as: POs waiting for my approval; my approvals list; pending approvals for me; PO not in my approval queue

Turbo lists the purchase orders the caller can approve now. These are first-stage POs assigned to the caller, second-stage POs where the caller is the priority second approver, and second-stage POs that other eligible second approvers see after the timeout. A single pending PO that is not available to the caller gets "pending purchase order not found". A missing ID gets "purchase order id is required".

## Pending purchase order lookup by id returns not found

Lines 610-640

Asked as: pending PO not found; PO not in my pending list; empty id pending lookup; find pending approval by id

If the purchase order is not among the caller's pending approvals, Turbo answers 'pending purchase order not found'. An empty id gives 'purchase order id is required'.

## Visible purchase orders list scopes and required parameters

Lines 642-681

Asked as: invalid PO list scope; filter purchase orders by status; scope group needs group; stale or expiring PO filter

The visible purchase orders list accepts scope all (the default), mine, active, rejected, stale, expiring, approved_by_me_awaiting_second or group. Any other scope fails with "scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group". Scope group needs a group, scope stale needs stale_before, and scope expiring needs expiring_before. The messages are "group is required when scope=group", "stale_before is required when scope=stale" and "expiring_before is required when scope=expiring".

## Visible purchase orders scope defaults to all and is validated

Lines 646-658

Asked as: default PO list filter; unknown scope on PO list; list all purchase orders default; scope validation message

The visible purchase orders list uses scope 'all' when none is given. An unknown scope returns 'scope must be one of: all, mine, active, rejected, stale, expiring, approved_by_me_awaiting_second, group'.

## Visible purchase orders limit must be positive

Lines 682-693

Asked as: limit on PO list invalid; purchase order list row limit; negative limit PO list; limit must be positive

The optional limit on the visible purchase orders list must be a positive integer. Zero, a negative number or text fails with "limit must be a positive integer". Without a limit, Turbo applies none.

## Single purchase order lookup respects visibility

Lines 709-736

Asked as: PO not found or not visible; can't view purchase order; open PO by id denied; PO id required message

Turbo returns one purchase order only if it is visible to the caller. A PO that is missing or not visible gets "purchase order not found or not visible". A missing ID gets "purchase order id is required".

## Searchable purchase orders list uses visibility rules

Lines 738-754

Asked as: search purchase orders list; find POs I can see; PO search results; all visible purchase orders

The purchase order search list returns the purchase orders the caller can see under the PO visibility rules, with scope 'all'.

## Rejection reason is required and needs 5 characters

Lines 771-803

Asked as: reject PO without reason; rejection reason too short; reason needs more characters; why reject message needed

Rejecting a purchase order without a rejection reason fails with "you must provide a rejection reason". A reason shorter than 5 characters after trimming spaces fails with "rejection reason must be at least 5 characters".

## Rejecting a purchase order that does not exist is refused

Lines 811-819

Asked as: reject PO that doesn't exist; rejection fails po not found; error fetching PO reject; reject missing purchase order

When staff reject a purchase order and Turbo cannot find it, Turbo refuses the rejection. The message starts with "error fetching purchase order:" and then gives the lookup error (code po_not_found). A rejection reason of at least 5 characters is checked before the lookup.

## Only unapproved, unrejected purchase orders can be rejected

Lines 821-861

Asked as: can't reject approved PO; reject already rejected PO; cannot reject again; rejecting active purchase order

A PO that is already rejected fails with "rejected purchase orders cannot be rejected again". A PO that is not Unapproved fails with "only unapproved purchase orders can be rejected". An invalid kind fails with "purchase order kind is invalid or no longer exists".

## Who can reject a purchase order

Lines 863-884

Asked as: who is allowed to reject; not authorised to reject PO; reject after first approval; approver rejection rights

A first-stage approver, a qualified second-stage approver, or the assigned approver who is valid through the bypass can reject. Anyone else gets "you are not authorized to reject this purchase order". A PO that already has its first approval can still be rejected by these people, because it is not yet Active.

## Rejection records the reason and notifies the creator

Lines 886-943

Asked as: rejection notification to creator; record who rejected PO; rejected PO email sent; rejection reason saved

Rejecting a PO records the time, the reason and the rejector. Turbo then sends the creator a rejection notification. If the notification fails, the rejection stays.

## Rejection notification failure does not undo the rejection

Lines 926-945

Asked as: rejection email failed; notification error on reject; rejection still saved after error; reject email not sent

If the rejection notification to the creator cannot be sent, Turbo logs the error and still returns the rejected purchase order.

## Cancelling a purchase order that does not exist is refused

Lines 958-966

Asked as: cancel PO that doesn't exist; cancellation fails not found; error fetching PO cancel; cancel missing purchase order

When staff cancel a purchase order and Turbo cannot find it, Turbo refuses the cancellation. The message starts with "error fetching purchase order:" and then gives the lookup error (code po_not_found). Turbo makes this check before any check on status or claim.

## Creator can cancel a rejected revision

Lines 968-980

Asked as: cancel my rejected revision; creator cancel rejected PO; cancel revision without payables; only creator can cancel

A rejected revision keeps its PO number but has no approval. Its creator can cancel it without payables staff rights, and the Active status check does not apply. Anyone else gets "only the creator can cancel a rejected revision".

## Cancelling a purchase order needs payables staff and Active status

Lines 981-1032

Asked as: cancel active purchase order; who can cancel PO; payables staff cancel PO; only active can be cancelled

Only an Active PO can be cancelled, or it fails with "only active purchase orders can be cancelled". The caller must have the payables_admin claim (payables staff), or it fails with "you are not authorized to cancel this purchase order". Cancelling sets the status to Cancelled and records the time and the canceller.

## Purchase order with expenses cannot be cancelled

Lines 1009-1027

Asked as: PO has expenses cancel; can't cancel PO with expense; expense linked blocks cancel; cancel blocked by expenses

If any expense is linked to the PO, cancelling fails with "this purchase order has associated expenses and cannot be cancelled".

## Cancellation records time and canceller

Lines 1029-1060

Asked as: cancelled date and who cancelled; record canceller; PO cancel timestamp; status changed to cancelled

Cancelling sets the status to Cancelled, and records the cancel time (cancelled) and the person who cancelled (canceller).

## Convert to Cumulative: allowed statuses and type

Lines 1064-1113

Asked as: convert PO to cumulative; which POs can become cumulative; one-time to cumulative; closed PO cumulative conversion

Only an Active PO, or a PO that Turbo closed automatically (closed_by_system), can be converted to Cumulative. Otherwise conversion fails with "only active or automatically closed purchase orders can be converted to Cumulative". Imported POs that are closed are refused, because their cause of closure is unknown. Only a One-Time PO can be converted, or it fails with "only One-Time purchase orders can be converted to Cumulative".

## Converting a purchase order returns it as Cumulative

Lines 1096-1104

Asked as: PO type now cumulative; result after cumulative conversion; imported PO closed cannot convert; converted PO returned

After conversion, the type becomes Cumulative and Turbo returns the saved purchase order. Imported purchase orders closed by the system cannot be converted.

## Convert to Cumulative needs payables staff

Lines 1115-1130

Asked as: who can make cumulative; cumulative needs payables rights; not authorised to convert PO; payables staff convert

Converting a PO to Cumulative needs the payables_admin claim (payables staff). Others get "you are not authorized to convert this purchase order to Cumulative".

## Expenses cannot total more than the PO to convert

Lines 1132-1151

Asked as: expenses exceed PO total convert; too many expenses cumulative; reduce expense before converting; expenses total over PO

Conversion fails if all expenses of the PO, committed or not, total more than the PO total, compared in whole cents. The message is "the expenses of this purchase order total more than the purchase order total, so it can't be Cumulative. Reduce or remove an expense first". A total equal to the PO total passes.

## Converting an auto-closed PO reopens it

Lines 1153-1177

Asked as: reopen auto closed PO; auto closed PO conversion; closed PO becomes active again; clear closed date on PO

When a PO that Turbo closed automatically is converted, committed expenses must be below the PO total. If they equal or exceed it, conversion fails with "committed expenses already use the full purchase order total, so there is no room for another invoice". Otherwise the PO becomes Active again, and its closed time, closer and closed_by_system flag are cleared. The type then changes to Cumulative.

## Auto-closed conversion refused when committed expenses use the full total

Lines 1153-1169

Asked as: no room for another invoice; committed expenses use full total; auto closed convert refused; committed expenses equal total

When converting an auto-closed PO, if committed expenses equal or exceed the PO total (compared in cents), Turbo refuses with 'committed expenses already use the full purchase order total, so there is no room for another invoice'.

## Auto-closed conversion clears the closure details

Lines 1170-1174

Asked as: closed fields cleared; status back to active; reset closed PO status; closed by system flag removed

Converting an auto-closed purchase order sets status back to Active and clears closed, closer and closed_by_system.

## Parent PO number format and sequence

Lines 1211-1224

Asked as: how PO numbers are formatted; next PO number sequence; PO number 4999 limit; unable to generate PO number

A new parent PO number has the format YYMM-NNNN, from the current year and month. The next number is the highest existing one for that month plus one. Numbers from 5000 up are reserved for manual or imported POs, so automatic numbering stops at 4999. If no number is free, generation fails with "unable to generate a unique PO number".

## Child PO number format and limit

Lines 1226-1292

Asked as: child PO number format; maximum child POs reached; 99 child POs limit; child PO suffix numbering

A child PO number is the parent's PO number plus a two-digit suffix, such as 2401-0001-01. The suffix is the highest existing child suffix plus one. The parent must exist and have a PO number. At most 99 children are allowed, and the next one fails with "maximum number of child POs reached (99)" and the parent number.

## Parent PO numbers: automatic range

Lines 1296-1353

Asked as: automatic PO number range; PO number below 5000; skip existing PO number; monthly PO number counting

Automatic parent PO numbers count up from the highest number of the current month below 5000. Turbo skips a number that already exists.

## Approver lookup validates its inputs

Lines 1358-1427

Asked as: approver lookup amount invalid; division missing approver lookup; has job not true or false; kind required approver search

The approver lookup needs a valid numeric amount, or it fails with "Amount must be a valid number". A has_job value that is not a boolean fails with "has_job must be a valid boolean". An empty division fails with "division is required". An empty kind fails with "kind is required". If no kind is given, it defaults to the default capital expenditure kind.

## Approver lookup defaults and optional inputs

Lines 1358-1426

Asked as: approver lookup default kind; approver search optional fields; capital expenditure default kind; numeric amount required lookup

The approver lookup needs a numeric amount ('Amount must be a valid number') and a division ('division is required'). Kind defaults to the capital expenditure kind, and has_job must be a boolean ('has_job must be a valid boolean') when given.

## Recurring PO approver lookup uses the total value

Lines 1429-1455

Asked as: recurring PO approver total; approvers for recurring PO; recurring PO total value; recurring PO start end frequency

For a Recurring PO, the approver lookup calculates the number of payments and the total value from the start date, end date and frequency. It finds approvers for that total, not for one payment. Recurrence errors return "invalid_recurrence" with the field messages. Other errors return "Error calculating recurring PO total".

## Recurring approver lookup reports recurrence errors

Lines 1439-1454

Asked as: recurring dates invalid error; bad frequency recurring PO; invalid recurrence message; recurring calculation error

When a recurring PO's dates or frequency are invalid, the approver lookup returns 'invalid_recurrence' with field errors. Other calculation failures say 'Error calculating recurring PO total'.

## Approver lookup converts foreign currency to home currency

Lines 1457-1470

Asked as: foreign currency approver conversion; currency no longer exists; no exchange rate for currency; home currency amount approvers

The approver lookup converts the amount to home currency with the currency's exchange rate. An unknown currency fails with "currency is invalid or no longer exists". A foreign currency with no rate (zero or less) fails with "selected foreign currency is missing an exchange rate".

## Approver lookup refuses invalid or unrated currency

Lines 1457-1469

Asked as: currency missing rate error; unknown currency on PO; exchange rate zero foreign currency; currency invalid message

An unknown currency gives 'currency is invalid or no longer exists'. A foreign currency without an exchange rate gives 'selected foreign currency is missing an exchange rate'.

## Approver lookup rejects an unknown kind

Lines 1472-1485

Asked as: kind no longer exists lookup; approver search invalid kind; unknown PO kind lookup; kind invalid for approvers

If the kind no longer exists, the approver lookup fails with "kind is invalid or no longer exists".

## Second approver lookup: requester who qualifies

Lines 1492-1500

Asked as: I can second approve myself; requester qualifies second approval; no second approver needed for me; requester approves second stage

When second approval is required and the requester is a second-stage approver, the lookup returns an empty approver list. The requester then approves in the second stage. The result is marked requester_qualifies.

## Second approver list is empty for a qualifying requester

Lines 1492-1499

Asked as: second approver list empty; requester selected automatically; empty approver list qualifying; auto select me second approval

If second approval is required and the requester can second-approve, the list is returned empty so the UI selects the requester automatically.

## Second approver lookup refuses when the pool is empty

Lines 1501-1508

Asked as: no second approvers available amount; nobody can final approve; second approver pool empty lookup; contact administrator second approval

If second approval is required, the requester does not qualify and no approvers exist, Turbo says 'no second-stage approvers can final-approve this amount; contact an administrator'.

## First approver list adds a requester with a limit

Lines 1516-1523

Asked as: add myself as first approver; requester with limit first approver; I can approve first stage; requester eligible first approval

The first approver list holds the first-stage approvers. When second approval is required, the requester is not a first-stage approver and has a non-zero limit for the kind, Turbo adds the requester as an eligible approver.
