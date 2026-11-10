<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expenses.go","sha256":"80ed5c7d33eaad7ad2617f1e8605a093ff6ea3e7ec84cb5adcc27be6cf0c4454","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expenses.go

## Expense without a user is refused when cleaned

Lines 27-38

When Turbo cleans an expense that has no user (uid), it refuses the save. The message is "no uid provided in for expense record".

## Home currency lookup and fallback to blank currency

Lines 40-51

Turbo looks up the home currency before it cleans an expense. For some payment types, it sets the expense currency to the home currency. If Turbo finds no home currency, it sets the currency to blank instead.

## Expense without a PO gets a default kind

Lines 53-57

When an expense has no purchase order, Turbo sets its kind by default. The kind depends on whether the expense has a job. The code's default is utilities.DefaultExpenditureKindIDForJob; the file does not say which kind it returns for each case (capital for no job, project for a job, per the docs).

## Bookkeeper-created expense is approved by its creator

Lines 59-67

When a bookkeeper creates the expense for another user, Turbo sets the approver to the creator. The creator must have an active user account, or Turbo refuses with 'bookkeeper approver must be an active user'. The creator must also hold the book_keeper claim, or Turbo refuses with 'the approver must have the book_keeper claim'. The creator is required.

## PO-linked expense takes branch and currency from the PO

Lines 69-104

When an expense links a purchase order, Turbo sets the expense's branch and currency from the PO. It refuses with 'referenced purchase order not found' if the PO is missing. It refuses with 'referenced purchase order is missing a branch' if the PO has no branch.

## Expense without PO takes branch from job or default branch

Lines 105-154

With no PO, Turbo takes the expense's branch from its job. It refuses with 'referenced job not found' or 'referenced job is missing a branch' if the job is missing or has no branch. With no job, Turbo uses the default branch (default_branch) on the user's admin profile. It refuses with 'your admin_profiles record is missing a default_branch' if none is set, or 'error looking up admin profile' if the profile is missing.

## Mileage expense needs valid personal vehicle insurance

Lines 162-217

A mileage expense needs the user's personal vehicle insurance expiry date on their admin profile. Turbo refuses with 'cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date' if the date is empty. It refuses with 'cannot submit mileage expense: personal vehicle insurance expired on <date>' if the expense date is after the expiry. An expense dated on the expiry day is accepted.

## Mileage total is calculated from the rate

Lines 219-271

For a mileage expense, Turbo calculates the total from the distance and the expense rate record (utilities.CalculateMileageTotal). It clears the vendor. Turbo shows 'error loading expense rate record' or 'error calculating mileage total' if these fail. At commit, Turbo recalculates the total with the year's committed mileage.

## Mileage distance must be a whole number

Lines 233-248

A mileage expense must have a distance that is an integer, or Turbo refuses with 'distance must be an integer for mileage expenses'.

## Personal reimbursement needs profile permission

Lines 280-304

A personal reimbursement expense is refused with 'cannot submit personal reimbursement expense: personal reimbursement is not enabled for your profile' unless the admin profile has allow_personal_reimbursement turned on.

## Allowance total is the sum of the meal and lodging rates

Lines 306-342

For an allowance expense, Turbo adds up the rates of the selected allowance types (Breakfast, Lunch, Dinner, Lodging) from the expense rate record. It sets the total to that sum. It sets the description to 'Allowance for ' followed by the types joined by commas. It clears the vendor.

## Currency defaults to home currency

Lines 353-373

With no PO, allowance, fuel card, mileage and personal reimbursement expenses use the home currency. Any expense with an empty currency also gets the home currency when one exists. Turbo refuses with 'referenced currency not found' if the currency cannot be found.

## Foreign currency needs an exchange rate

Lines 374-385

Turbo refuses an expense in a foreign currency with 'selected foreign currency is missing an exchange rate' when that currency has no positive rate.

## Settled total follows currency

Lines 387-398

For an expense in the home currency, Turbo sets settled_total to the total. For a foreign-currency expense with payment type Expense, Turbo leaves settled_total unchanged. For any other foreign-currency payment type, Turbo sets settled_total to 0. In every case Turbo clears the settler and the settled date.

## Approver is the manager, who must be active and hold tapr

Lines 402-470

For a normal expense, Turbo sets the approver to the manager on the submitter's profile. It refuses with 'your manager must be an active user to create expenses' if the manager is not active. It refuses with 'the approver must have the tapr claim' if the manager lacks the tapr claim. It refuses with 'no profile found for uid' if the submitter has no profile.

## Bookkeeper-created expense is approved by its creator (checks)

Lines 472-520

The creator must be provided ('creator is required'), be an active user, and hold the book_keeper claim. Otherwise the expense is refused.

## Bookkeeper creating for another user has strict conditions

Lines 522-590

A user may create an expense for another user only as a bookkeeper (the book_keeper claim) and with a purchase order. The PO must be Active and have a po_number. The PO must have an owner who is not the bookkeeper (same-user POs use the regular flow). The expense uid must equal the PO owner. The payment type must be OnAccount or CorporateCreditCard. A caller without the claim gets 'uid must be equal to the authenticated user's id'.

## Expense saves require sign-in and expense editing

Lines 596-634

Turbo first checks whether expense editing is allowed (checkExpensesEditing; the file does not say when it refuses). Turbo refuses an unauthenticated user with 'authentication is required'. It refuses a purchase order that does not exist with 'purchase order not found'.

## New expense records the creator

Lines 639-649

When a new expense is created, Turbo sets the creator to the signed-in user. If the expense uid is the creator's own id, it follows the normal manager approval flow. If it is another user's id, the bookkeeper on-behalf rules apply.

## Only the creator can edit an expense; uid and creator are fixed

Lines 650-708

Only the creator may edit an expense, or Turbo refuses with 'only the creator can edit this expense'. The uid cannot change ('uid cannot be changed'). The creator cannot change ('creator cannot be changed'). An expense that was created by someone other than its owner keeps bookkeeper approval and must again meet the bookkeeper on-behalf rules.

## Submitted expense cannot be edited

Lines 710-736

Turbo refuses to edit an expense that is already submitted, or to save one marked submitted, with 'cannot edit submitted expense'. Recall the expense first.

## Division must be active and user must be allowed the branch

Lines 743-750

After cleaning, Turbo checks that the expense's division is active (EnsureActiveDivision). Turbo also checks that the user may use the expense's branch (EnsureUserCanUseRecordBranch). The file does not state the messages.

## Attachment upload and source expense handling

Lines 766-797

The expense form's attachment field is a virtual input; Turbo stores the file on an expense document. A source_expense in the request is copied to the expense. If the request sends an empty attachment with no new file and no source_expense, Turbo treats it as an explicit removal and clears the linked document.

## Original attachment values are kept for edits

Lines 798-805

When a person edits an existing expense, Turbo remembers its original attachment document (attachment_document) and missing-attachment reason (attachment_missing_reason). A new expense has no original values to remember.

## PO-linked expense inherits kind and PO must be active

Lines 815-834

An expense with a PO takes its kind from the PO, and the user cannot edit it. Turbo refuses with 'purchase order must be active' unless the PO status is Active. Any user can submit against any active PO.

## Cumulative PO sums existing expenses

Lines 835-854

For a Cumulative PO, Turbo totals the PO's existing expenses to check the new expense against the PO. It refuses with 'error calculating cumulative total' if that total fails.

## Payables admin claim is checked before validating an expense

Lines 856-869

Before validating an expense, Turbo checks whether the person saving it holds the payables admin claim (payables_admin). If the check fails, Turbo shows "error checking payables_admin claim". Turbo passes the result, with the PO and the cumulative total of existing expenses, to validateExpense. That function is in another file, so its effect is not stated here.

## Attachment document is server-managed

Lines 871-879

A client cannot set the attachment document. For a new expense, or when the attachment is removed, Turbo clears the attachment document and the missing-attachment reason before validation. Otherwise Turbo restores the saved values.

## Bookkeeper claim is checked after validation

Lines 881-898

After the expense passes validateExpense, Turbo checks whether the person saving it holds the bookkeeper claim (book_keeper). If the check fails, Turbo shows "error checking book_keeper claim". Turbo then uses the result when it chooses the attachment document for the save.

## Attachment handling after validation

Lines 900-923

Some payment types skip attachments, and Turbo clears the document and reason for them. An explicit removal clears both. A new file or source expense makes Turbo create or reuse an expense document and clear the missing reason. Otherwise an existing expense keeps its document, and a new one has none.

## Attachment file is read from request body or multipart upload

Lines 926-946

Turbo takes the expense attachment from the request body's "attachment" field. If the body gives no files and the request is a multipart form, Turbo reads the uploaded file named "attachment" instead. A missing uploaded file is not an error. Turbo keeps any files it finds as the expense's pending attachment.

## Only real file values count as attachments

Lines 948-965

When Turbo reads the "attachment" value, it accepts one file or a list of files. In a mixed list, Turbo keeps the file items and ignores the others. Any other kind of value gives no attachment.
