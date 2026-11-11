<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_api.go","sha256":"9402f400f1938b6d219863960fddcaec0b4ea876da4d09f6bccc828041e77809","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_api.go

## Expense visibility by role: owner, approver, commit, report

Lines 33-37

Asked as: who can see my expense; approver cannot see expense; who has access to expense; why manager cannot view claim

A caller sees an expense if they own it or created it. Its approver sees it while it is submitted. A commit claim holder (committer) sees it once it is approved. A report claim holder sees it once it is committed.

## PO participants do not automatically see linked expenses

Lines 38-41

Asked as: PO creator cannot see expenses; some expenses missing on PO page; purchase order hides linked expenses; PO approver no expense access

Being a PO creator, approver or second approver gives no access to the expenses linked to that PO. A person can open a PO details page and still see only some of its linked expenses. They see the rest only if an expense rule also applies to them.

## My expenses list covers own and created expenses

Lines 43-44

Asked as: my expenses list; list of claims I created; filter expenses by purchase order; show expenses I submitted

The expenses list returns expenses the caller owns (uid) or created (creator). An optional purchase_order filter narrows it to one PO.

## Pending expenses list for approvers

Lines 45-45

Asked as: expenses waiting for my approval; pending approvals list; things I need to approve; list of claims to approve

The pending list shows expenses where the caller is the approver, the expense is submitted and it is not yet approved.

## Pending expenses list includes only submitted, unapproved expenses

Lines 45-46

Asked as: rejected expense in pending list; pending approvals wrong items; submitted but not in approvals; why pending list shows this

A pending expense is one where the caller is the approver, it is submitted, and it has no approval date. The approved list shows expenses where the caller is the approver and an approval date is set. It does not check submitted, committed or rejected.

## Approved expenses list for approvers

Lines 46-46

Asked as: expenses I approved; approved expenses history; list of approved claims; claims I signed off

The approved list shows expenses where the caller is the approver and the expense has an approved date.

## Who can open an expense's details or receipt

Lines 64-70

Asked as: who can open expense; cannot open expense details; who can view receipt; can approver open my claim

A person can open an expense (and its receipt) if they own it (uid) or created it (creator). Its approver can open it once it is submitted. A committer (the commit claim) can open it once it is approved. A reporter (the report claim) can open it once it is committed.

## Admin holders see committed or uncommit-request expenses

Lines 77-82

Asked as: admin cannot open expense; admin access to uncommitted expense; why admin cannot see claim; admin view committed expense

An admin (the admin claim) can open an expense's details only if it is committed or has an uncommit request on record. Other uncommitted expenses stay hidden from them unless another rule gives access.

## Payables staff see expenses only in settlement queue

Lines 91-93

Asked as: payables cannot see expense; finance staff expense access; settlement queue access; payables lost access to claim

Payables staff (the payables_admin claim) can open an expense and its receipt only while it is in the settlement queue, settled or not. Once the expense is committed or rejected, they lose this access. The queue condition is expenseSettlementQueueCondition, defined elsewhere.

## Details and receipt combine all three visibility rules

Lines 98-100

Asked as: receipt access rules; why details open for payables; who gets receipt access; combined expense access rules

The expense details and receipt routes allow anyone who passes the shared rule, the admin rule or the payables rule. The expense lists use only the shared rule.

## Commit and report claims decide expense visibility

Lines 127-142

Asked as: commit claim visibility; report permission on expense; has commit has report; claim holder expense access

Turbo checks whether the caller holds the commit claim (has_commit) and the report claim (has_report). The shared expense visibility rule uses these two checks. It does not look at purchase order participation.

## Admin and payables claims widen expense details access

Lines 146-162

Asked as: admin payables permission details; payables admin expense access; admin claim widens access; payables can open receipt

For expense details and receipts, Turbo also checks the admin claim (has_admin) and the payables_admin claim (has_payables_admin). These add the admin and payables visibility rules to the shared expense rule.

## Uncommit request state shown on expense details

Lines 257-258

Asked as: uncommit request status; uncommit request not showing; uncommit request on expense; why uncommit state missing

Expense details can include the state of an uncommit request (uncommit_request). Turbo sets it for commit and admin holders who view a committed expense. It is omitted when there is none.

## Default page size is 20

Lines 269-270

Asked as: default page size; how many rows per page; list shows only 20 items; default number of results

List endpoints that accept a limit use a default page size of 20 (defaultPageLimit) when the caller gives none.

## My expenses list loads committed expenses 50 at a time

Lines 272-274

Asked as: committed expenses paging; why only 50 committed expenses; more committed expenses next page; committed expenses in batches

The first page of the caller's expenses shows all non-committed expenses plus up to 50 committed ones. Later pages show committed expenses in batches of 50 (expensesListCommittedOverflowLimit). The caller's own limit is ignored.

## My expenses list ignores the limit parameter

Lines 274-351

Asked as: limit parameter ignored; page size not changing; set limit to 100; expenses list more than 50 rows

The expenses list ignores a caller-supplied limit. Page 1 returns all uncommitted expenses, then up to 50 committed ones, so it can return more than 50 rows. Later pages return committed expenses in batches of 50. The response limit is 50.

## Expense list requires sign-in

Lines 286-289

Asked as: unauthorized expense list; not logged in expenses; expenses list 401 error; need sign in for list

The expenses list, PO expenses list, pending list and approved list refuse a caller who is not signed in with "unauthorized".

## My expenses list ordering and page count

Lines 293-298

Asked as: order of my expenses; why committed shown after uncommitted; page count my expenses; invalid page number expenses

On page 1, uncommitted expenses come first, newest date first. Committed expenses follow, ordered by committed date, then expense date, newest first. Total pages is 1, plus one for each further batch of 50 committed expenses beyond the first 50. An invalid or non-positive page number becomes page 1.

## My expenses list page count and committed overflow

Lines 353-357

Asked as: total pages my expenses; committed expenses on page two; how many pages expenses; newest committed first page

On page 1 of the my-expenses list, Turbo shows all non-committed expenses plus up to 50 committed ones. Later pages show 50 committed expenses each, newest committed first. Total pages is 1 plus the committed expenses beyond the first 50, divided by 50 and rounded up.

## PO expenses list shows only expenses the caller may see

Lines 397-467

Asked as: PO expenses fewer than committed count; missing expenses on PO list; PO list not showing all expenses; PO expenses only some visible

The PO expenses list applies expense visibility, not PO visibility. A PO creator, approver or second approver sees only the linked expenses they own, created, approve (once submitted), or may see through the commit or report claim. The list can show fewer expenses than the PO's committed expense count.

## Missing PO id on the PO expenses list is refused

Lines 399-407

Asked as: purchase order id is required; PO id missing error; PO expenses unauthorized; no PO id given

The PO expenses list returns the error 'purchase order id is required' when no PO id is given. Without sign-in it returns 'unauthorized'.

## PO expenses list requires a visible PO

Lines 404-425

Asked as: purchase order not found; PO not visible error; hidden purchase order expenses; cannot list PO expenses

Listing the expenses of a PO first checks that the caller can see that PO. If the PO is missing or hidden, Turbo answers "purchase order not found or not visible". A missing PO id gives "purchase order id is required".

## Page size default and cap

Lines 434-443

Asked as: maximum 200 rows; limit above 200; invalid limit or page; default rows per list

The PO expenses, pending and approved lists default to 20 rows per page. A limit above 200 becomes 200. An invalid or non-positive limit or page falls back to the default. These lists are ordered by newest expense date; the PO list puts committed date first.

## PO expenses list page size and page count

Lines 467-482

Asked as: PO expenses page count; sort order PO expenses; number of pages on PO list; PO list paging

The expenses list of a purchase order uses the page and limit from the request. The page count is the number of visible expenses divided by the page size, rounded up. The list is sorted by committed date, then expense date, then creation date, newest first.

## Expense details and receipt require sign-in and an id

Lines 498-504

Asked as: id is required expense; expense details unauthorized; missing expense id; expense details need login

The expense details and receipt routes refuse a caller who is not signed in with "unauthorized". They refuse a missing id with "id is required".

## Expense details hidden from unauthorized callers

Lines 512-517

Asked as: expense not found or not authorized; cannot find my expense; expense hidden from me; not authorized to view expense

When the expense does not exist or the caller may not see it, Turbo answers "expense not found or not authorized". The message is the same in both cases.

## Expense details add PO owner mismatch and uncommit state

Lines 518-522

Asked as: PO owner mismatch; owner differs from PO owner; uncommit state on details; expense details flag mismatch

Expense details include a flag (po_owner_uid_mismatch) computed from the expense owner and the PO owner by expensePurchaseOrderOwnerUIDMismatch, defined elsewhere. They also include the uncommit request state, loaded by loadExpenseUncommitRequestState, defined elsewhere.

## Receipt download uses the expense details visibility rule

Lines 540-554

Asked as: cannot download receipt; receipt not authorized; who can download receipt; receipt access denied

Turbo opens a receipt only for a caller who may view the expense details. Otherwise it answers "expense not found or not authorized". If the expense has no stored attachment, it answers "expense attachment not found".

## Receipt download refuses when no attachment exists

Lines 552-565

Asked as: no receipt attached; attachment not found; receipt missing on expense; expense has no attachment

If the expense has no attachment, or its attachment location is missing, Turbo answers "expense attachment not found". The same message appears when the file cannot be read from storage.

## Missing receipt is reported with an error message

Lines 552-567

Asked as: receipt file cannot be read; attachment missing from storage; receipt file missing; receipt error message

If the attachment file cannot be read from storage, Turbo returns 'expense attachment not found'. A missing attachment name, collection or record id gives the same message.

## Receipt is shown inline

Lines 569-570

Asked as: receipt opens in browser; receipt download inline; receipt file name shown; view receipt in tab

A permitted caller gets the receipt file streamed inline under its stored file name.

## Pending expenses list needs sign-in and uses page size

Lines 584-636

Asked as: approval list paging; approval list max 200; pending list needs login; waiting approvals page size

The list of expenses waiting for the caller's approval requires sign-in. It uses a default page size of 20 and a maximum of 200. Invalid page or limit values fall back to page 1 and 20 per page. Results are ordered by expense date, newest first.

## Approved expenses list needs sign-in and uses page size

Lines 639-691

Asked as: approved list paging; approved list needs login; approved expenses max rows; approved list page size

The list of expenses the caller approved requires sign-in. It uses a default page size of 20 and a maximum of 200. Invalid page or limit values fall back to page 1 and 20 per page. Results are ordered by expense date, newest first.
