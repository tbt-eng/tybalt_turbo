<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_writeback.go","sha256":"39c8ecfec9ac5adb44132bc5c6f4692ad7e4d93b44b0465804cb172c03f283cc","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_writeback.go

## PO approver props timestamp fallback

Lines 12-12

If a PO approver props record has a blank created or updated time, Turbo uses the other one. If both are blank, it uses 1970-01-01 00:00:00.000Z (poApproverPropsMissingTimestampFallback).

## Expense export omits the rejected flag when false

Lines 234-238

In the writeback export, each expense always has submitted, approved and committed flags. The rejected flag appears only when it is true, and a rejection reason is given only when one exists (rejected, rejectionReason).

## Allowance type match is exact and case-sensitive

Lines 256-265

The expenses writeback export sets the Breakfast, Lunch, Dinner and Lodging flags by looking for those exact words in the expense's allowance types (allowance_types). A type that is spelled or capitalised differently does not set a flag.

## Who may call the expenses writeback export

Lines 269-289

The export accepts a Bearer token that matches an unexpired machine secret for legacy_writeback (ValidateMachineToken). Otherwise it accepts a signed-in user who holds the report claim. Anyone else gets the message "unauthorized".

## Writeback export requires updatedAfter

Lines 291-294

If the request has no updatedAfter value, Turbo refuses the export with the message "updatedAfter is required".

## Which expenses the writeback exports

Lines 296-384

The export includes only expenses updated at or after updatedAfter (>=, so the equal timestamp is included). Each expense must be committed and not imported (_imported = 0). Expenses whose owner's payroll ID is a generated placeholder are left out.

## Missing related records export as blank

Lines 296-374

When a user, approver, committer, job, division, vendor or PO link is missing, its exported fields are blank or 0 instead of causing an error.

## Settled total only with explicit currency

Lines 297-309

The export includes the settled total only when the expense has a currency set. Expenses with a blank currency omit it, and their currency code is blank.

## Query failures abort the export

Lines 386-391

If a database query for expenses, vendors, purchase orders or PO approver props fails, the export returns an error. The message is "failed to query" plus the name of the data and the error text.

## Vendors exported with writeback

Lines 393-424

The vendor list covers vendors on the matched expenses, which skip blank vendors and placeholder-payroll owners. It also covers vendors on every non-imported purchase order, even when updatedAfter would exclude that PO. Each vendor appears once.

## Purchase orders ignore updatedAfter and skip Unapproved

Lines 433-486

The export returns every non-imported purchase order whatever its update time. It leaves out purchase orders whose status is Unapproved, because they have no final PO number yet.

## All PO approver limits are exported

Lines 493-510

The export returns every PO approver props record (the limits for each approver), not just recent ones, so the legacy system gets a full copy.

## Export fails when an approver has no legacy uid

Lines 512-522

If any PO approver props record belongs to a user with no legacy uid, the whole export fails. The message is "missing legacy uid for po_approver_props id" followed by that record's id.

## Exported expense id is the immutable ID

Lines 525-529

In the expenses writeback export, each expense's immutableID is the expense's own record id. The legacy system uses it to match the expense when it folds the data in.

## Expense export identifies people by legacy uid

Lines 530-535

In the writeback export, the expense owner, manager (approver), committer and rejector are each identified by their legacy uid. Each also has a name taken from their profile. Turbo's own user IDs are not used.

## Exported names are trimmed of extra spaces

Lines 534-534

The expenses writeback export trims leading and trailing spaces from the owner's display name, the manager's name, the committer's name and the rejector's name.

## Expense export gives job number, client, category name and PO number

Lines 549-556

In the writeback export, each expense carries the job number (job), the job description, the client name, the category name and the purchase order number (po) as text. It does not carry job or category IDs. The vendor and purchase order IDs are listed separately (vendorId, purchaseOrderId).

## Expense workflow fields exported as flags

Lines 557-574

Approved, committed and rejected are exported as true when their timestamp is set, and false when blank. The commit time is exported as the committed timestamp. Names are trimmed of extra spaces.

## Allowance flags only when one applies

Lines 580-590

Breakfast, lunch, dinner and lodging flags are exported only if the expense has at least one of those allowance types. Then all four flags are exported as true or false. If it has none, the flags are left out.

## Expense total hidden for Allowance, Mileage, Meals

Lines 592-595

The export leaves out the expense total when the payment type is Allowance, Mileage or Meals. For other payment types it includes the total.

## Settled total exported only when expense has currency

Lines 597-602

The expenses writeback export includes settled_total on an expense only when the exported currency code is not blank. An expense with no currency leaves settled_total out of the payload.

## Vendor export lists id, name, alias and status

Lines 606-614

Each vendor in the writeback export has an id and a name. The alias and the status are left out when they are blank.

## Purchase order export uses internal IDs for job, category, parent PO and branch

Lines 619-654

In the writeback export, a purchase order's job (job), category (category), parent PO (parentPo) and branch (branch) are internal PocketBase IDs. They are not job numbers or names, which is how expenses are exported. The division is given as its code and name. The creator, approvers, canceller, closer and rejector are given by legacy uid.

## Purchase order export gives approval events as timestamps

Lines 633-648

In the writeback export, a purchase order's approved, second approval, cancelled, closed and rejected values are timestamps, and each is left empty when that event has not happened. Expenses are exported as true/false flags instead. The export also carries the PO's status, total and approval total (approvalTotal).

## PO approver limits exported per category

Lines 660-672

Each PO approver record in the writeback export gives the approver's legacy uid. It also gives the limits max_amount, project_max, sponsorship_max, staff_and_social_max, media_and_event_max and computer_max, plus the divisions list and the created and updated timestamps.

## Writeback response has four lists

Lines 675-681

A successful expenses writeback export returns one response with four lists: expenses, vendors, purchaseOrders and poApproverProps. A list with no matching rows comes back empty and is not left out.
