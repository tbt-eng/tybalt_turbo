<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expenses.go","sha256":"b572134660d86ce5fcd5fe049a24562c03d70a18fe4c63f6ad1b12bfd27d28ed","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expenses.go

## Expense needs a user

Lines 27-38

Asked as: expense has no owner; expense missing staff member; cannot clean expense; no uid error

Turbo refuses to clean an expense that has no uid (the staff member it belongs to). The error is "no uid provided in for expense record".

## Expense without a PO gets a default kind

Lines 53-57

Asked as: expense kind missing; kind set automatically; no PO kind; which expenditure kind

When an expense has no purchase order, Turbo sets its kind itself. It uses the default for a job when the expense has a job, and the default for no job (capital) when it has none (DefaultExpenditureKindIDForJob).

## Bookkeeper-created expense uses creator as approver

Lines 59-64

Asked as: bookkeeper submitting for someone; approver is me; creator not approver; book keeper claim error

When a bookkeeper creates an expense on behalf of another person, the creator must hold the book_keeper claim. Turbo sets the approver (approver) to the creator. Without the claim Turbo refuses with 'the approver must have the book_keeper claim'. A missing creator gives 'creator is required'.

## Approver is the submitter's manager

Lines 65-67

Asked as: who approves my expense; wrong approver; no profile found; my manager not set

On a normal expense, Turbo sets the approver to the manager on the submitter's profile. It refuses if there is no profile ("no profile found for uid").

## Expense with a PO takes branch and currency from the PO

Lines 69-104

Asked as: branch from purchase order; currency taken from PO; PO not found; PO missing branch

When an expense links a purchase order, Turbo sets the expense branch and currency from the PO. If the PO is not found, Turbo refuses with 'referenced purchase order not found'. If the PO has no branch, it refuses with 'referenced purchase order is missing a branch'.

## Expense without a PO takes branch from job or default branch

Lines 105-154

Asked as: which branch for expense; job branch missing; no default branch; admin profile error

When an expense has no PO but has a job, Turbo sets the branch from the job. It refuses with 'referenced job not found' or 'referenced job is missing a branch'. With no job, Turbo uses the default branch of the person's admin profile. It refuses with 'error looking up admin profile' if there is no profile. It refuses with 'your admin_profiles record is missing a default_branch' if the default branch is empty.

## Mileage expense needs valid vehicle insurance

Lines 162-217

Asked as: mileage insurance expired; car insurance date; vehicle insurance not valid; cannot claim mileage

A Mileage expense needs a personal vehicle insurance expiry date on the admin profile. If it is empty, Turbo refuses with 'cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date'. If the expense date is after the expiry date, Turbo refuses with 'cannot submit mileage expense: personal vehicle insurance expired on <date>'. An expense dated on the expiry day passes.

## Mileage expense fails if expense rate record cannot load

Lines 219-231

Asked as: mileage rate record missing; cannot load rate; mileage save failing; rate record error

When a Mileage expense is saved and Turbo cannot load the expense rate record for it, the save fails with 'error loading expense rate record'. This check runs after the insurance check and before the distance check.

## Mileage distance must be a whole number

Lines 233-248

Asked as: mileage decimal distance; distance not whole number; fractional miles error; distance must be whole

A Mileage expense refuses a distance that is not an integer, with 'distance must be an integer for mileage expenses'.

## Mileage total is calculated and vendor cleared

Lines 250-278

Asked as: mileage total wrong; how mileage is calculated; vendor removed mileage; mileage recalculated at commit

For Mileage expenses, Turbo calculates the total (total) with CalculateMileageTotal from the distance and the expense rate record, and clears the vendor. At commit, Turbo recalculates the mileage with the year's committed mileage.

## Personal reimbursement needs permission on admin profile

Lines 280-304

Asked as: cannot claim personal expense; personal reimbursement blocked; reimbursement not enabled; personal reimbursement permission

A PersonalReimbursement expense is refused unless the admin profile has allow_personal_reimbursement turned on. The message is 'cannot submit personal reimbursement expense: personal reimbursement is not enabled for your profile'.

## Allowance total and description are set automatically

Lines 306-342

Asked as: allowance amount calculated; breakfast lunch dinner total; allowance description; allowance vendor cleared

For an Allowance expense, Turbo adds the rates of the chosen allowance types (Breakfast, Lunch, Dinner, Lodging) from the expense rate record to get the total. It sets the description to 'Allowance for ' followed by the types joined with commas, and clears the vendor.

## Some payment types use home currency when there is no PO

Lines 353-362

Asked as: wrong currency without PO; expense defaults to home currency; empty currency expense; home currency default

Without a PO, Allowance, FuelCard, Mileage and PersonalReimbursement expenses get the home currency. Any expense with an empty currency also gets the home currency, if one exists.

## Expense currency must exist and have a rate

Lines 364-385

Asked as: currency not found; foreign currency no rate; no exchange rate; currency error on expense

Turbo refuses an expense whose currency is not found, with 'referenced currency not found'. It refuses a foreign currency with no exchange rate, with 'selected foreign currency is missing an exchange rate'.

## Settlement fields reset on every expense save

Lines 387-398

Asked as: settled total changed; settler cleared on save; settled amount zero; settlement reset

On every save Turbo clears the settler and settled date. For a home-currency expense, it sets the settled total (settled_total) equal to the total. For a foreign-currency expense with payment type Expense, it leaves the settled total as it is. For other foreign-currency expenses, it sets the settled total to 0.

## Manager must be active and hold the tapr claim

Lines 402-469

Asked as: manager not active; manager lacks tapr; who is my approver; approver cannot be found

For a regular expense, the approver is the manager on the submitter's profile. Turbo refuses if no profile is found ('no profile found for uid'). It refuses if the manager is not active ('your manager must be an active user to create expenses'). It refuses if the manager lacks the tapr claim ('the approver must have the tapr claim').

## Bookkeeper-created expense must have an approver with bookkeeper claim

Lines 472-503

Asked as: bookkeeper approver missing; creator required error; approver needs book keeper; bookkeeper expense approver

When a bookkeeper creates an expense for another person, Turbo checks that the creator is set and holds the bookkeeper claim (book_keeper). It refuses with 'creator is required' or 'the approver must have the book_keeper claim'.

## Bookkeeper creating for others needs a PO

Lines 505-534

Asked as: bookkeeper must link PO; expense for someone else; uid not matching user; bookkeeper PO needed

A bookkeeper who creates an expense for another person must link a PO, or Turbo refuses with 'purchase order is required for bookkeeper-created expenses'. The creator must hold the book_keeper claim. Otherwise Turbo says 'uid must be equal to the authenticated user's id'.

## Bookkeeper expense PO and payment rules

Lines 536-572

Asked as: bookkeeper PO not active; PO owner mismatch; payment type not allowed; own PO bookkeeper

For a bookkeeper-created expense, the PO must be Active ('purchase order must be active'), must have a po number ('purchase order must have a po number'), and must have an owner. The expense uid must match the PO owner ('uid must match the purchase order owner'). The payment type must be OnAccount or CorporateCreditCard ('bookkeeper-created expenses must be OnAccount or CorporateCreditCard'). If the PO belongs to the bookkeeper, Turbo refuses with 'same-user purchase orders must use the regular expense flow'.

## Expense saving needs sign-in and the editing switch

Lines 579-596

Asked as: cannot save expense; not signed in expense; expense editing disabled; authentication required expense

Turbo checks whether expense editing is allowed before it processes an expense (checkExpensesEditing). It refuses an unauthenticated request with 'authentication is required'.

## Referenced PO must exist

Lines 602-617

Asked as: PO does not exist; wrong purchase order number; missing PO reference; purchase order not found

When an expense names a purchase order that does not exist, Turbo refuses with 'purchase order not found'.

## Creating an expense for another user is bookkeeper-only

Lines 622-632

Asked as: submit expense for colleague; create expense for someone; bookkeeper on behalf; who is the creator

On create, Turbo sets the creator (creator) to the signed-in user. If the expense uid is not the signed-in user, Turbo runs the bookkeeper on-behalf checks, and the creator becomes the approver. If the uid is the user's own, the expense follows the normal manager approval, even for a bookkeeper using their own PO.

## Only the creator can edit an expense

Lines 634-691

Asked as: cannot edit someone else's expense; change expense owner; edit expense as bookkeeper; creator changed error

On update, only the creator can edit an expense. Anyone else is refused with 'only the creator can edit this expense'. The uid cannot change ('uid cannot be changed'). The creator cannot change ('creator cannot be changed'). If the creator differs from the uid, the bookkeeper checks run again.

## Submitted expenses cannot be edited

Lines 693-719

Asked as: edit submitted expense; locked expense cannot change; submitted expense locked; change after submit

Turbo refuses to edit a submitted expense, with 'cannot edit submitted expense'. The check covers both the stored record and the request.

## Expense division must be active and branch usable

Lines 726-733

Asked as: expense division inactive; branch not allowed expense; cannot use branch; division not active

After cleaning, Turbo checks that the expense division is active. It also checks that the user may use the expense's branch (EnsureUserCanUseRecordBranch).

## Attachment upload is read and its file name normalized

Lines 749-768

Asked as: upload receipt; attach file to expense; attachment not uploading; file name problem

Turbo reads the upload from the 'attachment' field of the request, from the body or a multipart form. It then normalizes the pending file name. An expense keeps its files on a linked document, not on the expense itself.

## Source expense can be set from request to reuse an attachment

Lines 769-773

Asked as: reuse receipt from expense; copy attachment from expense; source expense blank; link earlier expense

When the request has a non-blank 'source_expense', Turbo sets it on the expense. A blank value is ignored.

## Empty attachment field removes the attachment

Lines 774-780

Asked as: remove receipt; delete attachment; clear attachment field; take off receipt

When the request sends a blank 'attachment' with no uploaded file and no source expense, Turbo treats it as an explicit removal. It clears the linked attachment document (attachment_document).

## Edits keep the original attachment document and missing reason

Lines 781-788

Asked as: attachment lost on edit; receipt disappears when editing; missing reason reset; keep receipt when updating

On an update, Turbo remembers the stored attachment document and the attachment missing reason. It restores them unless the attachment is replaced or removed.

## Expense kind comes from its PO

Lines 798-803

Asked as: cannot change expense kind; kind taken from PO; expense kind wrong; kind locked by PO

When an expense has a PO, it takes its kind from the PO's kind (NormalizeExpenditureKindID, with the PO's job). The user cannot choose a different one.

## Expense PO must be active

Lines 805-816

Asked as: PO is closed; inactive purchase order; cannot use closed PO; PO not active expense

Turbo refuses an expense against a PO whose status is not Active, with 'purchase order must be active'. Anyone may submit an expense against any active PO, if the other rules pass.

## Cumulative PO total excludes the expense being edited

Lines 818-842

Asked as: cumulative PO total wrong; PO running total; cumulative total error; double counted expense total

For a Cumulative PO, Turbo adds this expense's total to the sum of the other expenses on the PO. When it updates an expense, the stored copy of that expense is left out of the sum. If the sum fails, Turbo says 'error calculating cumulative total'.

## Payables admin claim is checked before validation

Lines 844-857

Asked as: payables admin check; payables admin error; payables claim failed; validation claim error

Before validating an expense, Turbo checks whether the user holds the payables admin claim (payables_admin). It passes the result to the validation. If the check fails, the save fails with 'error checking payables_admin claim'.

## Attachment document is managed by Turbo

Lines 859-867

Asked as: attachment document client value; receipt flag wrong; new expense receipt reset; attachment checks failing

A new expense, or one with the attachment explicitly removed, starts with no attachment document and no missing-receipt reason. On an update without removal, these keep their stored values. A document value sent by the client cannot satisfy the attachment checks.

## Bookkeeper claim is checked after validation

Lines 874-886

Asked as: book keeper check error; bookkeeper claim failing; attachment changed by bookkeeper; bookkeeper claim after validation

After the expense passes validation, Turbo checks whether the user holds the bookkeeper claim (book_keeper). If the check fails, the save fails with 'error checking book_keeper claim'. The result affects how the attachment document is resolved.

## Attachment handling by payment type and intent

Lines 888-909

Asked as: receipt rules by payment type; no receipt needed; missing receipt reason cleared; receipt kept or removed

Payment types that skip attachments never keep a document or a missing reason. An explicit empty attachment with no new file and no source expense removes the document and the reason. A new file or a source expense links a document and clears the missing reason. Otherwise an update keeps the original document, and a new expense has none.

## Expense of a type with no attachment never keeps a document

Lines 889-894

Asked as: receipt on card expense; no attachment payment type; document removed automatically; skip receipt payment

For payment types that skip attachments (expensePaymentTypeSkipsAttachment), Turbo clears the attachment document and the attachment missing reason. The expense accepts no document relation.
