<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_api.go","sha256":"c51e0f5afde3247dfef2c501b896f58c233a379ad299448fe58f2b0166630f0b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_api.go

## PO participants do not automatically see linked expenses

Lines 24-41

Asked as: PO creator can't see expenses; expense count higher than list; PO totals but no expenses; linked expenses hidden from PO owner

Being a PO creator, approver or second approver does not by itself let a person see the expenses linked to that PO. A person sees an expense only if they own it, are its approver while it is submitted, hold the commit claim and it is approved, or hold the report claim and it is committed. The PO's expense counts can therefore be higher than the expenses a person can open.

## Who sees an expense in lists and the PO expense view

Lines 42-70

Asked as: why can't I see this expense; expense missing from my list; who can view an expense; expense hidden from approver

A person sees an expense when they own it (uid) or created it (creator). The expense's approver sees it while it is submitted. A committer (the commit claim) sees it once it is approved. A reporter (the report claim) sees it once it is committed. PO participation alone does not grant access to a PO's expenses.

## My expenses list matches owner or creator

Lines 43-44

Asked as: my expenses list missing items; expenses I created; filter my expenses by PO; expenses where I'm creator

The list of a person's own expenses (whereListMine) includes expenses where the person is the owner (uid) or the creator. When a purchase_order filter is given, it keeps only expenses on that purchase order (whereListMineByPO).

## Pending list filter: submitted and not yet approved

Lines 45-45

Asked as: expenses waiting for my approval; to approve list empty; submitted but not approved; pending expenses for me

The pending expenses list shows expenses where the caller is the approver, the expense is submitted (submitted = 1), and the approved field is empty or null.

## Approved list filter: expenses the caller approved

Lines 46-46

Asked as: expenses I already approved; approved by me list; my approval history; expenses I signed off

The approved expenses list shows expenses where the caller is the approver and the approved field is set. It does not check whether the expense is committed later.

## Approver sees submitted expenses in shared visibility

Lines 64-70

Asked as: approver can't open submitted expense; can't view expense to approve; need to see expense I approve; approver access to expense

Under the shared expense visibility rule, an expense's approver sees it once it is submitted (submitted = 1). A commit claim holder sees it when approved is not empty. A report claim holder sees it when committed is not empty. The owner and the creator always see it.

## Admin claim holders see committed or uncommit-requested expenses

Lines 72-82

Asked as: admin can't open expense; admin cannot see uncommitted expense; admin access to committed expense; admin expense details blocked

On the expense details page and the receipt, an admin (the admin claim) sees an expense only when it is committed or has an uncommit request on record. Other uncommitted expenses stay hidden from them.

## Admin visibility runs on uncommit requests

Lines 77-82

Asked as: uncommit request not visible to admin; admin sees uncommit requested expense; expense with uncommit request for admin; uncommitted expense still shown to admin

For an admin claim holder, the details and receipt routes add expenses that are committed (committed is not empty). They also add any expense that has an uncommit request record, even when it is no longer committed.

## Payables staff see expenses only in the settlement queue

Lines 84-93

Asked as: payables lost access to expense; payables can't open expense after commit; expense rejected payables no access; settlement queue expense access

On the expense details page and the receipt, payables staff (the payables_admin claim) see an expense only while it is in the settlement queue, whether settled or not. Once it is committed or rejected, they lose access. This uses expenseSettlementQueueCondition, the same condition as the queue.

## Payables visibility uses the settlement queue condition

Lines 91-93

Asked as: payables staff no access; accounts team can't view expense; payables access ended; payables visibility expense

A payables staff member (payables_admin claim) sees an expense only when it meets expenseSettlementQueueCondition, whether settled or not. This access ends when the expense is committed or rejected.

## Details and receipt use wider visibility than lists

Lines 95-101

Asked as: can open expense but not in list; details show expense list hides it; receipt visible but not in list; why details differ from list

Opening one expense and its receipt uses the shared visibility rule plus the admin and payables extras. Lists and the PO expenses view use only the shared rule.

## Details and receipt add admin and payables access

Lines 98-100

Asked as: admin can open details; payables can open receipt; admin receipt access; extra access for details page

The details and receipt routes combine the shared visibility rule with the admin rule and the payables rule. The lists use only the shared rule.

## Pending and approved lists sort by date, newest first

Lines 107-109

Asked as: pending list order; approved list sorted by date; newest approvals first; list paging for approvals

The list query behind the pending and approved lists sorts by expense date, newest first, and applies the limit and offset.

## Expense details also expose settlement and currency fields

Lines 197-206

Asked as: foreign currency expense details; exchange rate on expense; settled total and settler; expense settlement date shown

Each expense row carries the currency, exchange rate and rate date. It also carries the settled total, the settler and the date it was settled. These apply to foreign-currency expenses that payables staff settle.

## Expense details show vendor invoice and PO comparison fields

Lines 226-258

Asked as: vendor invoice number on expense; compare expense with PO fields; invoice date and total in details; expense PO vendor job

The expense details response adds the vendor invoice number, date and total. It also adds the linked PO's vendor, job, division, category, description, payment type, branch, kind, total, currency and owner. These let a viewer compare the expense with its PO.

## Expense details flag a PO owner who differs from the expense owner

Lines 259-259

Asked as: expense owner differs from PO owner; owner mismatch warning; po_owner_uid_mismatch flag; expense and PO owner different

The expense details response has a flag (po_owner_uid_mismatch). Turbo sets it by comparing the expense's owner (uid) with the PO's owner (po_uid), using a helper in another file.

## Uncommit request state appears on committed expense details

Lines 260-261

Asked as: uncommit request status on expense; committed expense uncommit request shown; uncommit request state details; uncommit status for viewer

Turbo attaches the uncommit request state (uncommit_request) to the expense details. This is for commit claim and admin claim holders who view a committed expense. Turbo loads it for the person viewing.

## My expenses list: first page layout

Lines 272-354

Asked as: my expenses first page rows; more than 50 rows returned; committed expenses at top of my list; first page of my expenses

The my-expenses list (/api/expenses/list) shows expenses the caller owns or created. A purchase_order filter narrows it to one PO. Page 1 returns all non-committed expenses, newest date first, plus up to 50 committed expenses, most recently committed first. The page size in the request (limit) is ignored, so page 1 can hold more than 50 rows.

## Pagination defaults and cap for expense lists

Lines 272-273

Asked as: limit above 200; how many rows per page; default page size 20; invalid page number fallback

The PO expenses view and the pending and approved lists use 20 rows per page by default. A limit above 200 becomes 200. An invalid or non-positive limit or page falls back to the default, and the page falls back to 1.

## My expenses list ignores the limit parameter

Lines 277-341

Asked as: limit parameter not working; page size ignored my expenses; my expenses more than 50 rows; limit ignored in expenses list

The my-expenses list ignores any limit the caller gives. It uses 50 committed expenses per batch. Page 1 returns all non-committed expenses and then up to 50 committed ones, so it can return more than 50 rows.

## Pending and approved lists need sign-in

Lines 289-293

Asked as: unauthorized expense request; expense route not signed in; missing id in expense request; expense api needs login

The pending, approved, my expenses, details and receipt routes refuse a caller who is not signed in with the message 'unauthorized'. They also require an id in the path where one is used.

## Later pages of my expenses hold committed overflow

Lines 340-360

Asked as: page 2 of my expenses; committed expenses on later pages; my expenses overflow pages; offset for expense pages

Pages after page 1 of my expenses show only committed expenses, starting at offset 50 + (page-2)*50. Total pages is 1 plus the committed expenses beyond the first 50, in batches of 50, rounded up. Committed expenses sort by committed time, then date, newest first.

## My expenses list: later pages and page count

Lines 356-371

Asked as: page count my expenses wrong; total on my expenses list; invalid page number expenses; later pages show only committed

Pages after the first return committed expenses only, in batches of 50, starting after the first 50. The page count is 1 plus the number of extra batches needed for committed expenses beyond 50. An invalid or non-positive page number becomes page 1. The response total counts all the caller's expenses, not only committed ones.

## PO expenses view can show fewer expenses than the PO counts

Lines 375-470

Asked as: PO shows fewer expenses than count; committed expenses count mismatch; PO expense list shorter; expenses missing from PO view

After the PO check, the PO's expenses are filtered again by expense visibility. A PO's creator or approvers see only the linked expenses they own, created, approve while submitted, or can see through the commit or report claim. The list can be shorter than PO totals such as committed_expenses_count. Results are sorted by committed date, then expense date, newest first.

## PO expenses list sort and missing id

Lines 407-470

Asked as: purchase order id is required; PO expenses sort order; empty PO id error; PO expenses route error

The PO expenses route refuses an empty purchase order id with 'purchase order id is required'. It sorts expenses by committed time, then date, then creation, newest first.

## PO expenses view requires a visible purchase order

Lines 412-428

Asked as: purchase order not found or not visible; cannot open PO expenses; PO not visible error; PO expense view denied

The PO expenses view first checks that the caller can see the purchase order. If not, Turbo says "purchase order not found or not visible". Without a login, the request is refused as unauthorized.

## Default page size is 20, capped at 200

Lines 437-446

Asked as: rows per page default; limit 200 maximum; invalid limit defaults; page size for approval lists

The pending, approved and PO expense lists show 20 rows per page unless a limit is given. A limit above 200 becomes 200. A limit or page that is not a positive number is ignored, so the defaults apply (page 1, 20 rows). Total pages is the total divided by the limit, rounded up.

## PO expenses list page count rounds up

Lines 479-486

Asked as: total pages zero expenses; PO expenses total pages; page count rounding; how many pages PO expenses

The PO expenses response gives the page, the limit used, the total number of visible matching expenses, and total pages. Total pages is the total divided by the limit, rounded up. With zero visible expenses, total pages is 0.

## Expense details do not inherit PO visibility

Lines 490-496

Asked as: can see PO but not expense; linked PO access doesn't open expense; PO visible expense hidden; PO access not enough for expense

Expense details are checked against expense visibility only. A person who can see the linked PO still cannot open an expense unless they pass the expense rules.

## Expense details refuse people without access

Lines 497-529

Asked as: expense not found or not authorized; cannot open expense details; access denied on expense; expense details blocked

Opening an expense the caller may not view returns "expense not found or not authorized". The details include a flag when the expense owner differs from the PO owner (po_owner_uid_mismatch), and the state of any uncommit request.

## Expense details not found message

Lines 509-521

Asked as: expense not found message; expense does not exist error; details route missing expense; expense not found for me

The details route answers 'expense not found or not authorized' both when the expense does not exist and when the caller may not see it. It also flags when the expense owner differs from the PO owner (po_owner_uid_mismatch).

## Expense receipt download rules

Lines 531-575

Asked as: download receipt error; expense attachment not found; receipt opens in browser; receipt download denied

Downloading a receipt uses the same visibility as expense details. If the caller may not view the expense, Turbo says "expense not found or not authorized". If the expense has no attachment or the file is missing, Turbo says "expense attachment not found". The file opens inline.

## Receipt route returns not found when no file

Lines 549-573

Asked as: receipt file missing; no receipt attached error; receipt not in storage; receipt route not found

The receipt route uses the same visibility as details and answers 'expense not found or not authorized' otherwise. It answers 'expense attachment not found' when the expense has no attachment or the file is missing from storage. The file opens inline under its stored name.

## Pending expenses list for approvers

Lines 585-639

Asked as: approver pending list; expenses to approve list; awaiting approval expenses; submitted expenses for approver

The pending list shows expenses where the caller is the approver, the expense is submitted and it is not yet approved. Results are sorted by expense date, newest first.

## Approved expenses list for approvers

Lines 641-693

Asked as: approved expenses for approver; expenses I approved later committed; approved date list; approved expenses sorted

The approved list shows expenses where the caller is the approver and the expense has an approved date. It includes expenses that were later committed. Results are sorted by expense date, newest first.
