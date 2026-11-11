<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expenses.go","sha256":"a8ac88b61eb70314a26349b8e39cd9ebed5a6bfa07bd80ab1c395262c7cce367","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expenses.go

## Expense requires a uid

Lines 27-38

Asked as: expense has no uid; cannot save expense missing user id; no uid error

Turbo refuses to clean an expense with no uid ('no uid provided in for expense record').

## Expense without a PO gets a default kind

Lines 53-57

Asked as: expense kind not set; what kind without PO; default capital or project

When an expense has no purchase order, Turbo sets its kind from whether it has a job. The helper utilities.DefaultExpenditureKindIDForJob(hasJob) picks the kind. The comments say a no-job expense defaults to capital and a job expense to project.

## Approver is the manager, or the bookkeeper creator

Lines 59-67

Asked as: who approves my expense; expense approver is my manager; who gets approval request

A normal expense takes the manager on the submitter's profile as its approver. An expense that a bookkeeper creates on behalf of another person takes the bookkeeper (the creator) as its approver.

## PO expense takes branch and currency from the PO

Lines 69-104

Asked as: expense branch from PO; currency changed from PO; PO branch expense error

An expense linked to a purchase order takes its branch and currency from that PO. Turbo refuses if the PO is not found ('referenced purchase order not found'). It also refuses if the PO has no branch ('referenced purchase order is missing a branch').

## Non-PO expense branch comes from job or default branch

Lines 105-154

Asked as: wrong branch on expense; job missing branch error; default branch admin profile

Without a PO, an expense takes its branch from its job. Turbo refuses if the job is missing ('referenced job not found') or has no branch ('referenced job is missing a branch'). With no job, Turbo uses the submitter's default branch from the admin profile. If none is set, it says 'your admin_profiles record is missing a default_branch'. If the admin profile cannot be found, it says 'error looking up admin profile'.

## Mileage needs valid vehicle insurance

Lines 159-217

Asked as: mileage insurance expired; vehicle insurance expiry error; cannot claim mileage insurance

A mileage expense needs the admin profile's personal vehicle insurance expiry to be set. If it is empty, Turbo says 'cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date'. If the expense date is after the expiry date, Turbo says 'cannot submit mileage expense: personal vehicle insurance expired on <date>'. An expense dated on the expiry day itself passes.

## Mileage expense needs an expense rate record

Lines 218-231

Asked as: mileage rate record error; error loading expense rate; cannot save mileage rate

When a mileage expense is saved, Turbo loads the expense rate record (utilities.GetExpenseRateRecord) after the insurance check. If it cannot load, Turbo refuses the save with "hook error when cleaning expense" (error_loading_expense_rate_record: "error loading expense rate record").

## Mileage distance must be a whole number

Lines 233-248

Asked as: mileage distance decimal; distance must be whole number; fractional miles not accepted

Mileage distance must be an integer, or Turbo says 'distance must be an integer for mileage expenses'.

## Mileage total is calculated and vendor cleared

Lines 250-278

Asked as: mileage total wrong; how mileage amount calculated; mileage vendor cleared

For a mileage expense, Turbo calculates the total from the distance and the expense rate (utilities.CalculateMileageTotal). It clears the vendor. The commit step recalculates the total using the whole year's committed mileage.

## Personal reimbursement needs permission on profile

Lines 280-304

Asked as: cannot submit personal reimbursement; personal reimbursement disabled; reimbursement not enabled for me

A personal reimbursement expense is refused unless the admin profile allows it (allow_personal_reimbursement). Turbo says 'cannot submit personal reimbursement expense: personal reimbursement is not enabled for your profile'.

## Allowance total and description are calculated

Lines 306-342

Asked as: allowance total calculation; allowance description meals; breakfast lunch dinner amount

For an allowance expense, Turbo sums the expense rate for each selected allowance type (Breakfast, Lunch, Dinner, Lodging). It sets that sum as the total. It sets the description to 'Allowance for ' followed by the types joined with commas, and clears the vendor.

## Allowance expense needs an expense rate record

Lines 306-319

Asked as: allowance save fails; allowance rate error; allowance error loading rate

When an allowance expense is saved, Turbo loads the expense rate record. If it cannot load, Turbo refuses the save with "hook error when cleaning expense" (error_loading_expense_rate_record: "error loading expense rate record").

## Home currency forced for some payment types without PO

Lines 353-362

Asked as: currency switched to home; fuel card currency wrong; mileage home currency forced

Without a PO, an expense paid as Allowance, FuelCard, Mileage or PersonalReimbursement uses the home currency. An empty currency also falls back to the home currency.

## Foreign currency must exist and have a rate

Lines 364-385

Asked as: currency has no exchange rate; foreign currency not found; exchange rate missing expense

Turbo refuses an expense whose currency is not found ('referenced currency not found'). It refuses a foreign currency with no positive exchange rate ('selected foreign currency is missing an exchange rate').

## Settlement fields are reset on save

Lines 387-398

Asked as: settled total reset; settled date cleared on save; settler removed after edit

On every save, Turbo clears the settler and the settled date. In the home currency, it sets the settled total to the expense total. In a foreign currency with payment type Expense, it leaves the settled total alone. In a foreign currency with any other payment type, it sets the settled total to 0.

## Expense manager approver must be active and hold tapr

Lines 402-470

Asked as: manager not active expense; approver lacks tapr; no profile for my user

Turbo refuses an expense if the profile has no match for the uid ('no profile found for uid'). It refuses if the manager is not an active user: 'your manager must be an active user to create expenses'. It refuses if the manager lacks the approver claim (tapr): 'the approver must have the tapr claim'.

## Manager approver lookup failures

Lines 402-429

Asked as: cannot find my manager; no manager on profile; approver lookup failed

Turbo needs the submitter's profile to find the manager. If no profile exists, it refuses with "no profile found for uid". The approver (manager) is checked as active and holding tapr.

## Bookkeeper approver must be active with book_keeper claim

Lines 472-520

Asked as: bookkeeper creating expense for someone; book_keeper claim missing; creator not active bookkeeper

When a bookkeeper creates an expense for someone else, the creator must be provided. The creator must be an active user ('bookkeeper approver must be an active user'). The creator must hold the bookkeeper claim (book_keeper), or Turbo says 'the approver must have the book_keeper claim'.

## Bookkeeper creating on behalf must follow strict rules

Lines 522-590

Asked as: expense for another person; bookkeeper expense needs PO; on behalf of owner mismatch

A person who creates an expense for someone else must hold the bookkeeper claim (book_keeper). Otherwise Turbo says 'uid must be equal to the authenticated user's id'. The expense needs a PO ('purchase order is required for bookkeeper-created expenses'). The PO must be Active, have a po number, and have an owner. The expense uid must match the PO owner ('uid must match the purchase order owner'). A PO that belongs to the bookkeeper must use the regular flow ('same-user purchase orders must use the regular expense flow'). The payment type must be OnAccount or CorporateCreditCard.

## Expense saving needs authentication and open editing

Lines 596-613

Asked as: not signed in expense; expense editing closed authentication; authentication is required expense

Turbo checks whether expense editing is open (checkExpensesEditing). It refuses unauthenticated requests ('authentication is required').

## Expense editing can be switched off

Lines 596-598

Asked as: expense editing disabled; cannot edit expenses now; expenses editing closed

Before anything else, Turbo checks whether expense editing is open (checkExpensesEditing). If editing is closed, the save is refused with that check's error.

## Expense needs a signed-in user

Lines 600-613

Asked as: not logged in save expense; validating uid error; session expired expense

Turbo refuses to save an expense when nobody is authenticated. The error is "hook error when validating uid" (authentication_required: "authentication is required").

## Linked purchase order must exist

Lines 615-634

Asked as: purchase order not found; PO number wrong; linked PO missing

When an expense names a purchase order that cannot be found, Turbo says 'purchase order not found'.

## On create, uid must be the user unless bookkeeper

Lines 639-649

Asked as: expense created for someone else; submit expense on behalf; creator set to me

On create, Turbo sets the creator to the signed-in user. If the uid differs from the user, the bookkeeper on-behalf rules apply and the creator becomes the approver.

## Only the creator can edit; uid and creator are fixed

Lines 650-708

Asked as: cannot edit someone else's expense; change uid of expense; creator changed error

Only the creator can edit an expense ('only the creator can edit this expense'). The uid cannot be changed ('uid cannot be changed'). The creator cannot be changed ('creator cannot be changed'). On an expense a bookkeeper created for someone else, the bookkeeper rules are checked again.

## Submitted expenses cannot be edited

Lines 710-736

Asked as: edit submitted expense; submitted expense locked; recall expense to edit

Turbo refuses to edit a submitted expense, whether the stored or the new record is submitted: 'cannot edit submitted expense'. A person must recall it first.

## Division must be active and branch usable

Lines 743-750

Asked as: division inactive expense; branch not usable expense; cannot use expense branch

The expense's division must be active. Turbo also checks that the user may use the record's branch (EnsureUserCanUseRecordBranch).

## Attachment upload is read from the request's attachment field

Lines 766-785

Asked as: receipt upload not working; attachment field request; file upload expense API

Turbo reads an uploaded file from the request field "attachment" (JSON body or multipart form) and treats it as a pending attachment. Turbo sanitizes the pending file names. Expense files are stored as expense documents, not on the expense itself.

## Attachment document is server-managed

Lines 786-797

Asked as: attachment document cleared; missing receipt reason cleared; cannot set attachment

A client cannot set the attachment document. On create, or when the attachment is explicitly removed, Turbo clears it and the missing-receipt reason. On other updates it keeps the stored values. An empty attachment with no new file and no source expense counts as an explicit removal.

## Source expense can be linked to reuse an attachment

Lines 786-790

Asked as: reuse receipt from another expense; source expense attachment; copy receipt previous expense

If the request has a non-blank source_expense, Turbo sets it on the expense so that the attachment of that expense can be reused. A blank value is ignored.

## Sending an empty attachment removes it only when nothing replaces it

Lines 791-805

Asked as: remove receipt from expense; clear attachment value; delete attachment expense

If the request sends an empty "attachment" value, Turbo treats it as an explicit removal only when there is no uploaded file and no source_expense. Then it clears attachment_document. On a new expense or on explicit removal, Turbo also clears attachment_missing_reason. Otherwise it restores both from the stored expense before validation.

## Expense kind is inherited from PO

Lines 815-820

Asked as: expense kind cannot change; kind taken from PO; why kind locked

For an expense with a PO, the kind comes from the PO's kind. A user cannot edit it. Turbo normalizes it with the PO's job status.

## PO must be Active for expense

Lines 822-833

Asked as: purchase order not active; PO closed expense; submit against inactive PO

Turbo refuses an expense on a purchase order that is not Active: 'purchase order must be active'. Any user can submit against any Active PO.

## Cumulative PO total excludes this expense's old total

Lines 835-859

Asked as: cumulative PO total double counted; PO total wrong on edit; cumulative purchase order limit

For a Cumulative PO, validation adds the totals of the other expenses on the PO to this one. On an update, Turbo leaves out this expense's stored total so it is not counted twice.

## Expense total is not set from the PO's expenses for non-cumulative POs

Lines 840-859

Asked as: PO total ignores other expenses; non cumulative PO total; PO total zero validation

Turbo sums the other expenses on a purchase order only when the PO's type (type) is "Cumulative". For other PO types the existing total passed to validation is 0.

## Payables admin claim is checked before validation

Lines 861-889

Asked as: error checking claim; payables admin permission; payables claim failed

Turbo checks whether the user holds the payables admin claim (payables_admin) and passes the result to expense validation. If the check fails, Turbo refuses with "error checking claim".

## Attachment rules by payment type on save

Lines 891-926

Asked as: payment type no receipt needed; receipt skipped for payment type; missing receipt reason saved

For payment types that skip attachments, Turbo clears the document and the missing-receipt reason. On explicit removal, Turbo clears both. With a new file or source expense, Turbo creates or reuses an expense document and clears the missing reason. With none of these, an existing expense keeps its stored document and a new expense has none.

## Bookkeeper claim check can fail

Lines 891-903

Asked as: error checking bookkeeper claim; book keeper check failed; claim check error save

Turbo checks whether the user holds the bookkeeper claim (book_keeper) when saving an expense. If the check fails, it refuses with "error checking claim".

## Attachment handling on save

Lines 905-926

Asked as: receipt not saved on expense; attachment lost on update; new receipt not linked

Payment types that skip attachments keep no document and no missing reason. Explicit removal clears both. A new file or a source expense links a document and clears the missing reason. An update with no attachment input keeps the stored document. A create with none leaves the document empty.

## Request attachment files accepted forms

Lines 931-970

Asked as: attach multiple files expense; attachment list accepted; multipart attachment files

Turbo accepts the attachment as a single file or a list of files, and ignores other values. If the body has none and the request is multipart, it reads the uploaded files named "attachment".
