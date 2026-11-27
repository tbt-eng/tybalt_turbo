<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_api.go","sha256":"20138999c5661ef43669fc6394462bd39920cdbdabed85487658957e4e1e6076","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_api.go

## PO participation does not grant expense visibility

Lines 23-41

Asked as: PO shows more expenses than I see; can't see expenses on my PO; PO count doesn't match my list; why missing expenses on purchase order

Being a PO creator, approver or second approver does not let a person see the linked expenses. They must also pass an expense rule. So a PO's expense count can be higher than the expenses a person sees.

## My expenses list includes expenses you created for others

Lines 43-44

Asked as: expenses I entered for someone else; where are expenses I submitted for others; my list missing expenses I approve; expenses I created for colleagues

The My expenses list (/api/expenses/list) returns expenses where the caller is the owner (uid) or the creator. It does not return expenses that the caller only approves.

## Pending expenses filter: submitted and not yet approved

Lines 45-45

Asked as: expenses waiting for my approval; to approve list empty; expenses I still need to approve; approval queue for me

The pending list returns expenses where the caller is the approver, the expense is submitted (submitted = 1) and approved is empty or null. Approved expenses leave this list.

## Approved expenses list: no check that the expense is submitted

Lines 46-46

Asked as: expenses I already approved; approved expenses gone from list; committed expenses showing in approved; list of expenses I signed off

The approved list returns expenses where the caller is the approver and approved is set. Turbo does not check the submitted flag or the committed state, so committed expenses stay in this list.

## Who can see an expense in lists and PO expense view

Lines 64-70

Asked as: who can see an expense; why can't I see this expense; expense visibility for committer; who gets to view expenses

A person sees an expense if they own it (uid) or created it (creator). Its approver sees it once it is submitted. A committer (the commit claim) sees it once it is approved. A report holder (the report claim) sees it once it is committed.

## Commit holders see approved expenses; report holders see committed ones

Lines 67-70

Asked as: committer doesn't see expense; approved expense not visible to committer; report holder can't see expense; who sees committed expenses

A commit claim holder (committer) sees an expense once it is approved (approved is not empty). A report claim holder sees an expense once it is committed. An expense's approver sees it only while it is submitted (submitted = 1).

## Extra expense details access for admin holders

Lines 77-82

Asked as: admin open committed expense; admin can't open expense; view uncommit request expense; admin access to expense details

An admin (the admin claim) can open the details of a committed expense. They can also open any expense that has an uncommit request on record. They cannot open other uncommitted expenses through this extra access.

## Payables staff see foreign-currency expenses only while in settlement queue

Lines 84-93

Asked as: payables can't open foreign currency expense; payables access to receipt; settlement queue expense access lost; foreign currency expense payables view

Payables staff (the payables_admin claim) can open an expense and its receipt only while it is in the settlement queue, settled or not. Once the expense is committed or rejected, they lose this access. The queue condition is expenseSettlementQueueCondition, defined elsewhere.

## Expense details and receipt use wider visibility than lists

Lines 95-101

Asked as: can open expense but not in list; details visible but list missing; receipt access wider than list; why details show more than list

Opening one expense or its receipt accepts the shared visibility rule plus the admin and payables extras. Lists use only the shared rule.

## Expense lists are sorted by date and excluded count of query helpers

Lines 103-109

Asked as: order of pending expenses; approved list sorting; expense list pages and offset; newest expenses first in list

The pending and approved expense lists sort by expense date, newest first (e.date DESC). Each list is cut into pages with a limit and an offset.

## Expense details and receipt checks use commit, report, admin and payables claims

Lines 127-162

Asked as: who may open receipt; which roles can view expense details; permissions for expense receipt; claims needed for expense details

To decide who may open an expense or its receipt, Turbo checks whether the caller holds the commit, report, admin and payables_admin claims. Details and receipt use all four claims. The PO expenses list uses only commit and report.

## Expense details show uncommit request state for committed expense

Lines 260-261

Asked as: uncommit request status on expense; expense shows uncommit pending; uncommit request not showing; committed expense request state

The expense details response carries the state of the uncommit request (uncommit_request) that Turbo loads for the caller. It is set for commit and admin holders viewing a committed expense.

## Expense list pages hold committed expenses in fixed batches

Lines 277-282

Asked as: how many committed expenses per page; fifty committed expenses batch; page size for my expenses; default page size committed expenses

On the my expenses list, the page size for committed expenses is fixed at 50 (expensesListCommittedOverflowLimit). The default page size of 20 (defaultPageLimit) applies only to the pending and approved lists.

## My expenses list: first page shows all uncommitted plus 50 committed

Lines 280-359

Asked as: first page more than 50 rows; my expenses page one too long; why my list has so many rows; limit ignored on my expenses

The caller's expense list (owner or creator) ignores any limit the caller sends. Page 1 returns all non-committed expenses plus up to 50 committed expenses, so it can hold more than 50 rows. Non-committed expenses sort by date, newest first, then by creation, newest first.

## My expenses list ignores a limit set by the caller

Lines 282-282

Asked as: set limit on my expenses; my expenses limit not working; page size ignored in my list; limit says 50 but more rows

The My expenses list ignores any limit the caller sends. The batch size is always 50 committed expenses. The response reports Limit as 50, even when page 1 returns more rows.

## Expense details and receipt need a signed-in user

Lines 294-297

Asked as: unauthorized expense request; expense page says unauthorized; not logged in expenses error; expense list no user

Every expense list, details, PO expenses and receipt route refuses a request with no signed-in user with 'unauthorized'.

## My expenses list later pages and page count

Lines 300-306

Asked as: second page of my expenses; my expenses page count wrong; later pages only committed expenses; invalid page number expenses

Pages after 1 hold committed expenses only, in batches of 50, ordered by committed date newest first, then date. The total pages is 1 plus one for each further 50 committed expenses beyond the first 50. An invalid or non-positive page number becomes page 1.

## My expenses list can filter by purchase order

Lines 307-320

Asked as: filter my expenses by PO; expenses for one purchase order; search my expenses by order; only my expenses on a PO

When a purchase_order is given, the list returns only the caller's own expenses (owner or creator) on that PO. The total counts all of the caller's matching expenses.

## My expenses list sorting order

Lines 336-359

Asked as: order of my expenses list; why committed expenses at bottom; sort my expenses by date; uncommitted first in my expenses

On page 1 of My expenses, uncommitted expenses come first, newest date first, then the first 50 committed expenses, newest committed first. Later pages return only committed expenses in the same order.

## My expenses page count counts only committed overflow

Lines 361-365

Asked as: my expenses number of pages; why extra page of expenses; total pages expenses calculation; page count ignores uncommitted

My expenses always has at least 1 page. It adds one page for each further 50 committed expenses beyond the first 50 (rounded up). The total field counts all matching expenses, committed or not.

## PO expenses list requires the PO to be visible

Lines 412-433

Asked as: PO expenses not visible; purchase order not found on expenses; can't open PO expense list; missing PO id expenses

The list of expenses on a PO first checks that the caller can see the PO. If not, Turbo answers "purchase order not found or not visible". A missing PO id gives "purchase order id is required".

## PO expenses route needs a PO id

Lines 412-415

Asked as: PO expenses no id error; purchase order id required; expenses on PO without id; PO expenses request missing id

The PO expenses route refuses a request with no purchase order id with 'purchase order id is required'.

## PO expenses error messages

Lines 421-433

Asked as: purchase order not visible error; error fetching purchase order; PO lookup failed expenses; PO not found message

If the purchase order is not visible to the caller or does not exist, Turbo answers 'purchase order not found or not visible'. If the lookup fails, it answers 'error fetching visible purchase order'.

## PO expenses list returns only expenses the caller may see, unpaged

Lines 435-458

Asked as: PO expenses list missing items; fewer expenses than PO shows; PO expenses not paged; expenses on PO sort order

The PO expenses list returns every linked expense that passes expense visibility, with no paging. It orders by committed date newest first, then date, then creation. It can return fewer expenses than the PO's counts show.

## Expense details refused when not found or not allowed

Lines 469-500

Asked as: expense not found or not authorized; can't open expense details; expense id required error; PO owner mismatch flag

Opening an expense that does not exist, or that the caller may not view, gives "expense not found or not authorized". A missing id gives "id is required". Details also flag whether the PO owner differs from the expense owner (po_owner_uid_mismatch) and include the uncommit request state.

## Expense details and receipt need an id

Lines 476-479

Asked as: receipt request no id; expense receipt id required; details without id error; missing expense id

The expense details and receipt routes refuse a request with no id with 'id is required'.

## Expense details flag a PO owner mismatch and show uncommit request

Lines 493-497

Asked as: expense owner differs from PO owner; PO owner mismatch warning; why expense shows owner mismatch; uncommit state on expense details

Expense details add a flag when the expense owner differs from the PO owner (po_owner_uid_mismatch, from expensePurchaseOrderOwnerUIDMismatch). They also include the uncommit request state for the caller (loadExpenseUncommitRequestState).

## Expense receipt download rules

Lines 503-547

Asked as: download expense receipt; receipt file not opening; view attachment for expense; receipt download fails

The receipt download uses the same visibility as expense details and gives "expense not found or not authorized" otherwise. If the expense has no attachment, or the file is missing, Turbo says "expense attachment not found". The file opens inline.

## Receipt refused when expense has no attachment

Lines 524-529

Asked as: expense has no receipt; attachment missing on expense; receipt not found error; no attachment on expense

If the expense has no attachment, or its attachment file is missing, the receipt route answers 'expense attachment not found'. A caller who may not see the expense gets 'expense not found or not authorized'.

## Receipt opens inline as a download stream

Lines 537-545

Asked as: receipt opens in browser; receipt file name when downloaded; receipt won't load from storage; stream expense attachment

The receipt streams the stored attachment file inline under its original file name. If the file cannot be read from storage, Turbo answers 'expense attachment not found'.

## Pending expenses list for approver

Lines 557-611

Asked as: expenses to approve list; pending expenses page size; approver's pending list max rows; who sees pending expenses

The pending list shows expenses where the caller is the approver, the expense is submitted and it is not yet approved. It orders by date, newest first. The default page size is 20 and the maximum is 200.

## Pending and approved lists default to 20 rows, maximum 200

Lines 573-582

Asked as: how many rows per page expenses; limit above 200 expenses; pending list default rows; page count expense list rounding

The pending and approved expense lists show 20 rows per page unless a limit is given. A limit above 200 becomes 200. A limit or page that is not a positive whole number is ignored and the default applies (page 1, 20 rows). The page count is the total divided by the limit, rounded up.

## Approved expenses list for approver

Lines 613-665

Asked as: approved expenses for approver; list of approved expenses by date; approver approved list page size; expenses I approved with date

The approved list shows expenses where the caller is the approver and the expense has an approved date. It orders by date, newest first. The default page size is 20 and the maximum is 200. A limit above 200 becomes 200.
