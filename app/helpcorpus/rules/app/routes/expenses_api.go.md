<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_api.go","sha256":"d6ce7f7ebf9c288d55c0c24cb97b72096555aef55503da0366b90e4ea802b291","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_api.go

## My expenses list: what it includes

Lines 43-44

The expenses list shows expenses where the caller is the owner (uid) or the creator. Turbo can narrow it to one purchase order with the purchase_order filter.

## Pending expenses list definition

Lines 45-45

The pending list shows expenses where the caller is the approver, the expense is submitted and it is not yet approved.

## Approved expenses list definition

Lines 46-46

The approved list shows expenses where the caller is the approver and the expense has an approved timestamp.

## Who can see an expense

Lines 64-70

A person can see an expense if they own it (uid) or created it (creator). Its approver can see it once it is submitted. A committer (the commit claim) can see it once it is approved. A reporter (the report claim) can see it once it is committed. Being involved in the linked purchase order does not grant access.

## Admin holders see committed or uncommit-requested expenses

Lines 77-82

Admin claim holders get extra access to expense details. They can see an expense that is committed, or one that has an uncommit request on record. They cannot see other uncommitted expenses through this admin access.

## Pending expenses list sorted by date, newest first

Lines 89-91

The pending expenses list and the approved expenses list sort by expense date, newest first. Total pages is the total divided by the limit, rounded up.

## Visibility checks look up commit and report claims

Lines 109-124

Expense visibility depends on whether the caller holds the commit claim and the report claim. If these claims cannot be checked, Turbo returns "error checking expense visibility claims".

## Default page size is 20 rows

Lines 231-232

Expense lists that take a limit use a default page size of 20 rows (defaultPageLimit) when the caller gives no valid limit.

## My expenses list ignores the limit and pages committed expenses

Lines 232-313

The expenses list ignores any limit the caller sends. Page 1 shows all non-committed expenses and then up to 50 committed expenses. Each later page shows the next 50 committed expenses. Page 1 can therefore hold more than 50 rows.

## Expense endpoints require sign-in

Lines 249-251

The expense list, purchase order expenses, details, attachment, pending and approved endpoints refuse a caller who is not signed in with the message "unauthorized".

## Invalid page value defaults to page 1

Lines 255-260

A page that is missing, not a number, or not above 0 is treated as page 1. This applies to the expense list endpoints.

## My expenses list ordering

Lines 290-299

On page 1, non-committed expenses are sorted by date, newest first, and then by creation, newest first. Committed expenses are sorted by committed time, then date, then creation, newest first.

## My expenses list page count

Lines 315-329

The expenses list has one page when there are 50 or fewer committed expenses. Each further 50 committed expenses, rounded up, adds one page. Total counts all of the caller's expenses.

## Purchase order expenses need a visible purchase order

Lines 366-387

The purchase order expenses list needs a purchase order id, otherwise it says "purchase order id is required". If the caller cannot see the purchase order, or it does not exist, Turbo says "purchase order not found or not visible".

## Expense id is required

Lines 366-369

Expense details and attachment requests without an id are refused with "id is required". A purchase order expenses request without an id is refused with "purchase order id is required".

## Page size default and cap

Lines 396-404

The purchase order expenses, pending and approved lists use 20 rows per page by default. A limit above 200 is cut to 200. A limit that is not a positive number is ignored.

## Purchase order expenses are filtered by expense visibility

Lines 415-429

The purchase order expenses list shows only the linked expenses the caller may see under the expense rules. Owners, approvers and second approvers of the purchase order get no extra access. The list can therefore be shorter than the purchase order's expense counts. Rows are sorted by committed time, then date, then creation, newest first.

## Purchase order expenses list orders committed expenses first

Lines 429-429

The list of expenses linked to a purchase order sorts by committed date (e.committed) newest first, then expense date newest first, then creation time newest first. The count and the list use the same visibility filter, so the total matches the rows a person may see.

## Purchase order expenses list page count and total

Lines 438-445

The purchase order expenses list returns the page, the limit, the total of visible expenses and total pages. Total pages is the total divided by the limit, rounded up.

## Expense details access and refusal

Lines 463-484

Expense details are returned only if the caller passes the expense visibility rules or the admin rule. Otherwise Turbo says "expense not found or not authorized". A missing id gives "id is required".

## Expense details include purchase order owner mismatch and uncommit state

Lines 485-489

Expense details carry a flag when the expense owner differs from the purchase order owner (po_owner_uid_mismatch). They also carry the expense's uncommit request state for the caller (uncommit_request).

## Expense attachment download

Lines 507-542

The attachment download uses the same access rules as expense details, and refuses with "expense not found or not authorized". If the expense has no attachment, or the file cannot be read, Turbo says "expense attachment not found". Otherwise it streams the file inline under its file name.

## Expense attachment uses the same visibility as expense details

Lines 507-523

The attachment download applies the expense visibility rules and the admin rule used for expense details. A caller who cannot see the expense gets "expense not found or not authorized".

## Expense attachment missing returns not found

Lines 524-537

If the expense has no attachment file, or its attachment collection or record id is empty, Turbo refuses with "expense attachment not found". The same message appears when the stored file cannot be read.

## Attachment opens inline with the original file name

Lines 541-542

A permitted attachment download is sent inline as a binary stream, with the stored file name as the file name.

## Boolean claim flags become 0 or 1 for visibility

Lines 546-552

Turbo turns the commit, report and admin claim checks into 1 (claim held) or 0 (not held) before it applies the visibility rules.

## Pending list shows submitted, unapproved expenses assigned to caller

Lines 554-608

The pending list shows expenses where the caller is the approver, the expense is submitted, and the approved date is empty or missing.

## Pending expenses list requires sign-in

Lines 558-561

The pending expenses list refuses a caller who is not signed in and answers 'unauthorized'.

## Pending and approved lists use the same page size rules

Lines 563-579

The pending and approved expense lists take the page and limit from the request. An invalid or non-positive value falls back to page 1 and a limit of 20. A limit above 200 becomes 200.

## Pending and approved lists are paginated

Lines 570-605

The pending and approved expense lists are paginated, sorted by expense date, newest first. The page count is the total divided by the page size, rounded up.

## Approved expenses list: only the approver's approved expenses

Lines 606-662

The approved list (whereApproved) shows expenses where the caller is the approver (approver) and the approved date (approved) is not empty. It does not check the submitted flag. Staff see only the expenses they approved.

## Approved expenses list is paginated, newest date first

Lines 610-662

The approved list uses the page and limit rules of the pending list. It defaults to 20 rows per page and caps the page size at 200. It sorts by expense date, newest first, and returns the total and total pages.

## Approved list requires sign-in

Lines 613-616

The approved expenses list refuses a caller who is not signed in and answers 'unauthorized'.
