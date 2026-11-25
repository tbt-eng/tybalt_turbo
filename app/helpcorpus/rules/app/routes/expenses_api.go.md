<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_api.go","sha256":"0795d8fe767959dd243ee4de07537ebe11c422446cb28c9d0932fb4d48ac1904","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_api.go

## PO participants do not automatically see the PO's expenses

Lines 21-41

Asked as: PO shows fewer expenses; can't see PO expenses; PO totals don't match expenses; missing expenses on purchase order; PO creator can't see expenses

Being a purchase order's creator, approver or second approver does not by itself let a person see the expenses linked to that PO. Expense visibility comes only from the expense rules: the expense owner, the expense approver while the expense is submitted, a committer (the commit claim) once it is approved, and a reporter (the report claim) once it is committed. A person can therefore open a PO and see fewer linked expenses than the PO's totals count.

## Who can see an expense in lists and on a purchase order

Lines 42-70

Asked as: who can see expense; expense not showing in list; why can't I see this expense; expense visibility rules; see expense on PO

An expense is visible to its owner (uid) or its creator, whoever they are. Its approver sees it once it is submitted. A committer (the commit claim) sees it once it is approved. A report holder (the report claim) sees it once it is committed. PO participation alone gives no access to the expense.

## My expenses list shows only own expenses

Lines 43-44

Asked as: my expenses list; expenses not showing for me; filter expenses by PO; not signed in expenses error; list only my expenses

The expenses list returns only expenses where the caller is the owner (uid) or the creator. An optional purchase_order filter narrows the list to that PO. A caller who is not signed in gets "unauthorized".

## Pending expenses list for approvers

Lines 45-45

Asked as: expenses waiting for my approval; pending approvals list; expenses to approve; approver queue not loading; pending expenses order

The pending list shows expenses where the caller is the approver, the expense is submitted and it is not yet approved. The default page size is 20 and the maximum is 200. A larger limit is cut to 200, and an invalid limit gives 20. The list is sorted by date, newest first.

## Approved expenses list for approvers

Lines 46-46

Asked as: approved expenses list; expenses I approved; approved expenses gone after commit; list of approved claims; approver history of expenses

The approved list shows expenses where the caller is the approver and the expense has an approval date. It stays in the list after the expense is committed. It is sorted by date, newest first. It uses a default page size of 20, capped at 200.

## Admin holders see committed or uncommit-requested expenses on details

Lines 72-82

Asked as: admin can't see committed expense; admin view committed expense; uncommit request hidden; admin access expense details; admin cannot open expense

On the expense details page and the receipt, an admin holder (the admin claim) also sees an expense that is committed. They also see any expense that has an uncommit request on record. Other uncommitted expenses stay hidden from them.

## Payables staff see expenses only while in settlement queue

Lines 84-101

Asked as: payables can't open expense; foreign currency settlement access; payables staff expense access; settlement queue expense hidden; payables lost access expense

On the expense details page and the receipt, payables staff (the payables_admin claim) see an expense only while it is in the foreign-currency settlement queue, settled or not. Once the expense is committed or rejected, they lose this access. Lists do not include this extra access.

## Pending and approved expense lists sort newest date first

Lines 107-109

Asked as: expense list sort order; expenses newest first; expense list order by date; oldest expenses showing first; sorting of pending list

The pending expenses list and the approved expenses list show expenses by expense date, newest first (ORDER BY e.date DESC). Both lists page the results with a limit and an offset.

## Expense details include PO comparison and uncommit request data

Lines 226-259

Asked as: expense details PO comparison; owner mismatch warning; PO owner differs from expense; uncommit request status on expense; expense details extra fields

The expense details response adds the linked PO's vendor, job, division, category, description, payment type, branch, kind, total, currency and owner, so a viewer can compare them with the expense. It also sets a flag (po_owner_uid_mismatch) that compares the expense's person with the PO's owner. For commit and admin holders viewing a committed expense, it includes the state of any uncommit request (uncommit_request).

## My expenses list pagination: all uncommitted first, committed in batches of 50

Lines 274-356

Asked as: expenses list pages wrong; more than 50 rows on page 1; committed expenses paging; page size ignored expenses; expenses list pagination

Page 1 of the expenses list returns all uncommitted expenses, then up to 50 committed expenses (newest committed first). Page 1 can therefore hold more than 50 rows. Later pages return 50 committed expenses each. A limit sent by the caller is ignored. A missing, invalid or non-positive page is treated as page 1.

## My expenses list total pages calculation

Lines 358-372

Asked as: expenses list total pages; wrong page count expenses; how many pages expenses; total count expenses wrong; pagination total incorrect

The expenses list reports total as the count of all matching expenses. Total pages is 1, plus one page for each further 50 committed expenses beyond the first 50.

## Purchase order expenses require the PO to be visible

Lines 409-430

Asked as: PO expenses not visible; purchase order not found error; can't load PO expenses; PO id required error; PO not visible to me

The expenses of a purchase order are returned only if the caller can see that PO. Otherwise Turbo answers "purchase order not found or not visible". A missing PO id gives "purchase order id is required".

## Purchase order expense list is filtered by expense visibility

Lines 432-455

Asked as: PO expense list shows fewer; other people's expenses on PO; PO expenses hidden; expenses missing from PO list; PO expense list sorting

The expenses list for a purchase order shows only the linked expenses that the caller may see under the expense rules. PO owners, approvers and second approvers do not see other people's expenses unless they also qualify as owner, creator, approver, committer or report holder. The list can therefore show fewer expenses than the PO's counts. It is not paged and is sorted by committed date, then date, newest first.

## Expense details access and error

Lines 466-498

Asked as: expense not found or not authorized; can't open expense details; expense details access denied; expense details warning flag; expense details error

The expense details page returns the expense only if the caller passes the details visibility rules (owner, creator, approver of a submitted expense, commit, report, admin or payables extras). Otherwise Turbo answers "expense not found or not authorized". The details include a warning flag when the expense owner differs from the PO owner, and the state of any uncommit request.

## Expense receipt download rules

Lines 500-543

Asked as: download expense receipt; receipt not found; expense attachment missing; can't open receipt file; receipt file name inline

The receipt is available to the same people who can open the expense details. Otherwise Turbo answers "expense not found or not authorized". If the expense has no attachment, or the file cannot be read, Turbo answers "expense attachment not found". The file opens inline with its own file name.

## Pending expenses paging defaults

Lines 556-608

Asked as: pending list page size; limit above 200 pending; pending expenses per page; invalid limit pending expenses; pending expenses paging

The pending expenses list uses a default page size of 20. A limit above 200 is cut to 200. An invalid or non-positive limit or page falls back to the default or page 1.

## Approved expenses paging defaults

Lines 611-662

Asked as: approved list page size; approved expenses per page; approved expenses limit; approved list unauthorized; approved expenses paging

The approved expenses list uses a default page size of 20. A limit above 200 is cut to 200. An invalid or non-positive limit or page falls back to the default or page 1. A caller who is not signed in gets "unauthorized".
