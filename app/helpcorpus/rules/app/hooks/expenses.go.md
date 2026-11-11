<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expenses.go","sha256":"27ecf124030a75a6e8fd1b12391b1db02fbeeead6ea7f3f82218ba8a529d29c2","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expenses.go

## Expense needs an owner (uid)

Lines 27-38

Asked as: expense has no owner; cannot clean expense; missing uid on expense; expense without user

Turbo refuses to clean an expense that has no uid. The message is 'no uid provided in for expense record'.

## Expense without a PO gets a default kind

Lines 53-57

Asked as: expense kind wrong; kind set automatically; capital or project default; which kind is my expense

When an expense has no purchase order, Turbo sets its kind by itself. The default depends on whether the expense has a job: capital with no job, project with a job. The code calls utilities.DefaultExpenditureKindIDForJob.

## Expense approver is the creator when a bookkeeper enters it

Lines 59-64

Asked as: bookkeeper entered expense for me; who approves bookkeeper expense; bookkeeper approver not active; approver is me

When a bookkeeper creates an expense for another person, Turbo sets the approver to the creator. The creator must be an active user and hold the bookkeeper claim (book_keeper). Otherwise Turbo refuses with 'bookkeeper approver must be an active user' or 'the approver must have the book_keeper claim'. A missing creator gives 'creator is required'.

## Expense approver is the submitter's manager

Lines 65-66

Asked as: who approves my expense; my manager cannot approve; no profile found; manager not set on profile

For a regular expense, Turbo sets the approver to the manager on the submitter's profile (profiles.manager). With no profile Turbo refuses with 'no profile found for uid'. If the manager is not active it refuses with 'your manager must be an active user to create expenses'. If the manager lacks the approver claim (tapr) it refuses with 'the approver must have the tapr claim'.

## PO expense takes branch and currency from the PO

Lines 69-104

Asked as: branch from purchase order; currency from PO; PO not found on expense; PO has no branch

An expense with a purchase order gets its branch and currency from that PO. Turbo refuses with 'referenced purchase order not found' if the PO is missing, and with 'referenced purchase order is missing a branch' if the PO has no branch.

## Referenced job or PO must exist and have a branch

Lines 70-128

Asked as: purchase order not found; job not found; job has no branch; missing branch on job or PO

If the PO is not found, Turbo refuses with "referenced purchase order not found". If the PO has no branch, it refuses with "referenced purchase order is missing a branch". For an expense with no PO, a missing job gives "referenced job not found" and a job with no branch gives "referenced job is missing a branch".

## Expense without PO takes branch from job or default branch

Lines 105-154

Asked as: branch for expense without PO; default branch missing; job branch missing; which branch does expense use

An expense without a PO gets its branch from its job. If the job is not found Turbo refuses with 'referenced job not found'. If the job has no branch it refuses with 'referenced job is missing a branch'. With no job, Turbo uses the default branch on the submitter's admin profile. If that is empty it refuses with 'your admin_profiles record is missing a default_branch'. If the admin profile cannot be found it refuses with 'error looking up admin profile'.

## Missing default branch or admin profile blocks no-job, no-PO expense

Lines 130-152

Asked as: no job expense branch error; admin profile error; default branch empty; expense with no PO or job

An expense with no PO and no job takes the default branch from the submitter's admin profile (default_branch). If the profile is missing, Turbo refuses with "error looking up admin profile". If the default branch is empty, it refuses with "your admin_profiles record is missing a default_branch".

## Mileage expense needs valid personal vehicle insurance

Lines 162-217

Asked as: mileage insurance expired; cannot claim mileage; vehicle insurance date; insurance expiry blocks mileage

A mileage expense is refused if the admin profile has no personal vehicle insurance expiry: 'cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date'. It is also refused if the expense date is after the expiry date: 'cannot submit mileage expense: personal vehicle insurance expired on <date>'. An expense dated on the expiry date is allowed.

## Mileage needs insurance expiry date on admin profile

Lines 163-188

Asked as: insurance expiry date missing; update vehicle insurance; mileage needs insurance date; admin profile insurance blank

A mileage expense needs the admin profile to have a personal vehicle insurance expiry date. If it is empty, Turbo refuses with "cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date". If the admin profile is not found, Turbo refuses with an error about the profile lookup.

## Mileage expense dated after insurance expiry is refused

Lines 206-217

Asked as: mileage date after insurance; insurance ran out; expense date past insurance; insurance expired on date

Turbo compares the expense date with the insurance expiry date. It refuses only when the expense date is after the expiry date, so an expense dated on the expiry day passes. The message is "cannot submit mileage expense: personal vehicle insurance expired on <date>".

## Mileage distance must be a whole number and sets total

Lines 219-278

Asked as: mileage distance decimal; miles must be whole number; mileage total wrong; vendor cleared on mileage

A mileage expense needs a whole-number distance, or Turbo refuses with 'distance must be an integer for mileage expenses'. Turbo calculates the total from the distance and the expense rate (utilities.CalculateMileageTotal). It also clears the vendor. At commit the mileage total is recalculated using the year's committed mileage.

## Mileage and allowance need an expense rate record

Lines 219-231

Asked as: expense rate not loading; mileage rate missing; allowance rate error; cannot load rate record

For mileage and allowance expenses, Turbo loads the expense rate record. If it cannot, it refuses with "error loading expense rate record".

## Mileage and allowance clear the vendor

Lines 264-342

Asked as: allowance vendor removed; allowance description changed; vendor disappears on allowance; allowance total meals

Turbo clears the vendor on a mileage or allowance expense. For an allowance, Turbo also sets the description to "Allowance for " followed by the allowance types separated by commas. The total is the sum of the rate for each type (breakfast, lunch, dinner, lodging).

## Personal reimbursement needs permission on the profile

Lines 280-304

Asked as: cannot claim personal reimbursement; reimbursement not allowed; enable personal reimbursement; reimbursement option blocked

A personal reimbursement expense is refused unless the admin profile has allow_personal_reimbursement turned on. The message is 'cannot submit personal reimbursement expense: personal reimbursement is not enabled for your profile'.

## Allowance expense total and description are calculated

Lines 306-342

Asked as: allowance total wrong; breakfast lunch dinner totals; allowance description text; allowance amount calculated

For an allowance expense Turbo adds up the rates for each selected allowance type (Breakfast, Lunch, Dinner, Lodging) from the expense rate record. The sum becomes the total. The description becomes 'Allowance for ' plus the types joined by commas. The vendor is cleared.

## Home currency for non-PO allowance, fuel card, mileage, reimbursement

Lines 353-362

Asked as: currency set to home currency; fuel card currency; mileage currency automatic; no currency on expense

Without a PO, expenses with payment type Allowance, FuelCard, Mileage or PersonalReimbursement are set to the home currency. Any expense with no currency also gets the home currency, if one exists.

## Expense currency must exist and have a rate

Lines 364-385

Asked as: currency not found; foreign currency no rate; missing exchange rate; pick another currency

Turbo refuses an expense whose currency is not found: 'referenced currency not found'. It refuses a foreign currency with no positive exchange rate: 'selected foreign currency is missing an exchange rate'.

## Settled total reset on save

Lines 387-398

Asked as: settled total changed; settled amount reset; settled total zero; settler cleared on save

On every save Turbo clears the settler and settled fields. For a home-currency expense it sets the settled total equal to the total. For a foreign-currency expense with payment type Expense it leaves the settled total alone. For other foreign-currency payment types it sets the settled total to 0.

## Manager approver checks

Lines 402-470

Asked as: manager checks failing; manager not active; manager lacks tapr; manager becomes approver

The manager approver (profile manager) is checked in three steps. The manager must be an active user ('your manager must be an active user to create expenses'). The manager must hold the tapr claim ('the approver must have the tapr claim'). Then the manager becomes the expense's approver.

## Manager must be active to approve expenses

Lines 402-447

Asked as: manager inactive for approval; approver not active; no manager to approve; submit to manager error

For a regular expense, Turbo takes the approver from the manager on the submitter's profile. If the manager is not an active user, Turbo refuses with "your manager must be an active user to create expenses". If the submitter has no profile, Turbo refuses with "no profile found for uid".

## Manager approver needs the tapr claim

Lines 449-469

Asked as: manager cannot approve; manager missing approver claim; tapr claim needed; approver permission missing

The manager on the submitter's profile must hold the time sheet and expense approver claim (tapr). Otherwise Turbo refuses with "the approver must have the tapr claim". When the check passes, Turbo sets the manager as the approver.

## Bookkeeper approver must be active and hold the book_keeper claim

Lines 472-520

Asked as: bookkeeper creator not active; book_keeper claim missing; creator required error; bookkeeper approver check

When a bookkeeper enters an expense for another person, Turbo checks the creator as approver. The creator must be given, must be an active user, and must have the bookkeeper claim (book_keeper). Otherwise Turbo refuses with "bookkeeper approver must be an active user", "creator is required" or "the approver must have the book_keeper claim".

## Bookkeeper creating on behalf must follow PO rules

Lines 522-590

Asked as: bookkeeper PO rules; PO must be active for bookkeeper; payment type not allowed for bookkeeper; same user PO refused

A person who creates an expense for another user must hold the bookkeeper claim and must name a PO. Otherwise Turbo refuses with 'purchase order is required for bookkeeper-created expenses' or 'uid must be equal to the authenticated user's id'. The PO must be Active ('purchase order must be active'). It must have a PO number ('purchase order must have a po number') and an owner ('purchase order is missing an owner'). A PO owned by the bookkeeper is refused: 'same-user purchase orders must use the regular expense flow'. The expense uid must be the PO owner ('uid must match the purchase order owner'). The payment type must be OnAccount or CorporateCreditCard ('bookkeeper-created expenses must be OnAccount or CorporateCreditCard').

## Bookkeeper on-behalf expense needs a PO

Lines 522-531

Asked as: bookkeeper expense without PO; PO needed for someone else expense; enter expense for colleague; on behalf needs PO

A bookkeeper creating an expense for someone else must link a purchase order. Otherwise Turbo refuses with "purchase order is required for bookkeeper-created expenses".

## Only bookkeepers can create expenses for others

Lines 533-551

Asked as: cannot enter expense for others; not my expense uid; only bookkeepers enter for others; wrong user on expense

If the logged-in user lacks the bookkeeper claim (book_keeper) and the uid is not their own, Turbo refuses with "uid must be equal to the authenticated user's id".

## Bookkeeper PO must be active, numbered and owned by another user

Lines 553-574

Asked as: PO owned by me; PO has no number; PO missing owner; PO not active for bookkeeper

The PO must be Active ("purchase order must be active") and have a PO number ("purchase order must have a po number"). The PO needs an owner ("purchase order is missing an owner"). A PO owned by the bookkeeper is refused with "same-user purchase orders must use the regular expense flow".

## Bookkeeper expense owner and payment type restrictions

Lines 575-590

Asked as: payment type on bookkeeper expense; expense owner must match PO; OnAccount or corporate card; bookkeeper payment method

The expense uid must match the PO owner ("uid must match the purchase order owner"). The payment type must be OnAccount or CorporateCreditCard ("bookkeeper-created expenses must be OnAccount or CorporateCreditCard").

## Expense saving needs login and an existing PO

Lines 596-634

Asked as: cannot save expense; not logged in save; PO does not exist; expense save refused

Turbo refuses an expense save without authentication ('authentication is required'). It refuses a purchase order that does not exist ('purchase order not found'). Turbo first checks whether expense editing is allowed (checkExpensesEditing).

## Expenses are not editable when editing is turned off

Lines 596-598

Asked as: expenses locked; editing turned off; cannot edit expenses now; expense editing disabled

Before processing an expense, Turbo checks whether expense editing is allowed (checkExpensesEditing). If that check fails, the save is refused.

## Expense needs a login

Lines 602-613

Asked as: not logged in; session expired expense; authentication needed; save without login

Saving an expense without a logged-in user is refused with "authentication is required".

## Creator is set on create and cannot change

Lines 636-708

Asked as: who created expense; cannot change creator; cannot change owner; creator set on create

On create Turbo sets the creator to the logged-in user. If the expense uid differs from the user, the bookkeeper on-behalf rules apply. On update only the creator can edit ('only the creator can edit this expense'). The uid cannot change ('uid cannot be changed'). The creator cannot change ('creator cannot be changed'). An expense whose creator differs from its owner is treated as bookkeeper-created and uses the on-behalf checks again.

## Creating an expense for another user needs the bookkeeper path

Lines 639-649

Asked as: expense for another user; expense for myself approval; bookkeeper expense for self; uid differs from me

Creating an expense where uid differs from the logged-in user goes through the bookkeeper checks. An expense for yourself uses the regular manager approval, even for a bookkeeper.

## Only the creator can edit an expense; uid cannot change

Lines 650-702

Asked as: only creator can edit; cannot change expense owner; cannot validate expense update; someone else edited my expense

On update, only the creator may edit the expense. Others get "only the creator can edit this expense". The owner (uid) cannot change: "uid cannot be changed". The creator cannot change: "creator cannot be changed". If the original record is missing, Turbo refuses with "cannot validate ownership for this expense update".

## Bookkeeper rules apply again on update of on-behalf expense

Lines 703-708

Asked as: editing on-behalf expense; bookkeeper expense edit rules; creator stays approver; edit expense for colleague

If an expense was created by someone other than its owner, the bookkeeper-on-behalf checks run again on every edit. The creator stays the approver.

## Submitted expenses cannot be edited

Lines 710-736

Asked as: cannot edit submitted expense; expense already submitted; recall expense first; edit after submit

Turbo refuses to edit an expense that is already submitted or that the request marks as submitted: 'cannot edit submitted expense'. A person must recall it first.

## Division must be active and branch usable

Lines 743-750

Asked as: division inactive; cannot use branch; branch not allowed; division not active on expense

After cleaning, Turbo checks that the expense's division is active. It also checks that the user may use the expense's branch (EnsureUserCanUseRecordBranch).

## Source expense in request is copied to the expense

Lines 786-790

Asked as: copy existing expense; reuse source expense; source expense attachment; duplicate expense receipt

A non-empty source_expense in the request is set on the expense. Turbo uses it to reuse that expense's document as the attachment.

## Attachment document is managed by Turbo

Lines 791-797

Asked as: attachment keeps changing; cannot set receipt; attachment removed by itself; receipt document managed

A client cannot set the attachment document. On create, or when the attachment is explicitly removed, Turbo clears the attachment document and the missing-attachment reason. On other updates it restores the stored values before validation. Sending an empty attachment with no new file and no source expense counts as removal.

## Removing attachment needs an empty attachment and no new file or source expense

Lines 792-797

Asked as: remove receipt from expense; delete attachment; clear attachment field; empty attachment removal

A request with an empty "attachment" field counts as an explicit removal only when no file is uploaded and the expense has no source expense. Then Turbo clears the attachment document.

## Stored attachment and file-hash read

Lines 798-813

Asked as: receipt upload fails; file hash error; attachment not saving; reuse same receipt file

Turbo hashes any uploaded attachment file to create or reuse its expense document. If hashing fails, the save is refused. On update, Turbo remembers the stored attachment document and missing reason.

## Expense kind comes from its PO

Lines 815-820

Asked as: kind from purchase order; cannot change expense kind; kind follows PO; expense kind locked

An expense with a PO takes its kind from the PO. The user cannot edit it. Turbo normalizes the kind using whether the PO has a job.

## PO must be active for expenses

Lines 822-827

Asked as: PO not active for expense; closed purchase order; purchase order inactive; cannot use this PO

An expense on a purchase order that is not Active is refused: 'purchase order must be active'.

## Any staff may use any active PO

Lines 829-833

Asked as: anyone use any PO; use someone else's PO; PO owner check; submit against other PO

Anybody can submit an expense against any active PO if the other rules are met. This was decided on August 26, 2025. Turbo does not check the caller against the PO.

## Cumulative PO total excludes this expense's old total

Lines 835-859

Asked as: cumulative PO total; PO running total; cumulative total error; PO over limit

For a Cumulative PO, Turbo sums the totals of the other expenses on the PO and adds this expense's total. On an update it leaves out this expense's stored total. If the calculation fails Turbo says 'error calculating cumulative total'.

## Cumulative PO sum for new expense uses all expenses

Lines 840-859

Asked as: cumulative PO new expense; cumulative total wrong; PO total includes all; cumulative calculation error

For a Cumulative PO, a new expense is checked against the cumulative total of all expenses on the PO. An update uses the total of the other expenses only. An error gives "error calculating cumulative total".

## Payables staff claim is checked during expense save

Lines 861-874

Asked as: payables admin check; payables staff claim; claim check error; payables permission error

Turbo checks whether the user has the payables_admin claim (payables staff) and passes the result to the expense validation. If the check fails, Turbo refuses with "error checking payables_admin claim".

## Attachment fields are reset before validation

Lines 876-884

Asked as: receipt cleared on save; attachment reset; client attachment ignored; missing attachment reason reset

For a new expense or an explicit attachment removal, Turbo clears the attachment document and the missing-attachment reason. For other updates, it restores both from the stored record, so a value from the client is ignored.

## Invoice details and bookkeeper claim used when saving attachment

Lines 891-930

Asked as: vendor invoice details; invoice saved with expense; book_keeper claim error; attachment transaction failed

Turbo checks the bookkeeper claim (book_keeper). If the check fails, it refuses with "error checking book_keeper claim". It then saves the vendor invoice details with the document in the same transaction as the expense.

## Attachment handling by payment type and request

Lines 905-930

Asked as: which payment types need receipt; receipt not kept; no receipt needed; new file replaces receipt

Payment types that skip attachments never keep a document or missing reason (expensePaymentTypeSkipsAttachment). An explicit removal clears both. A new file or a source expense links a document, reusing or creating it (resolveExpenseDocumentForSave), and clears the missing reason. With no attachment input, an existing expense keeps its document and a new expense has none.

## Attachment upload can come from request body or multipart form

Lines 933-953

Asked as: upload receipt; attach file to expense; multipart upload; file name attachment

Turbo reads the attachment from the request body field "attachment". For multipart/form-data requests, it also looks for uploaded files named "attachment". File names are normalized.
