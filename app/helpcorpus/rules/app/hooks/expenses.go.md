<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expenses.go","sha256":"e76b69c250917d8115102200b12059b5d9b6d4d56dd3f0b209237b931c89d8b3","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expenses.go

## Expense needs a user ID

Lines 27-38

Asked as: expense save refused no uid; missing user id on expense; no uid provided; expense has no user

When Turbo cleans an expense whose user (uid) is empty, it refuses the save. The message is 'no uid provided in for expense record'.

## Expense without a purchase order gets a default kind

Lines 53-57

Asked as: expense kind empty; why is my expense capital; expense defaults to project; kind set automatically

When an expense has no purchase order, Turbo sets its kind (kind) by itself. It uses the default kind for an expense with a job when the expense has a job, and the default kind for one without a job otherwise. The comment in the code says the default is capital without a job and project with a job.

## Expense approver is the creator when a bookkeeper enters it

Lines 59-64

Asked as: bookkeeper entered expense approver; approver set to me; creator must be bookkeeper; approver must have book keeper claim

When a bookkeeper (the book_keeper claim) creates an expense on behalf of another person, Turbo sets the approver (approver) to the creator. The creator must have the book_keeper claim. Otherwise Turbo refuses with 'the approver must have the book_keeper claim'. A missing creator gives 'creator is required'.

## Expense approver is the submitter's manager

Lines 65-66

Asked as: who approves my expense; expense approver is my manager; manager not set on expense; no profile found for uid

For a regular expense, Turbo sets the approver (approver) to the manager on the submitter's profile. A missing profile gives 'no profile found for uid'.

## Expense with a purchase order takes its branch and currency from the PO

Lines 69-104

Asked as: expense branch from purchase order; currency taken from PO; referenced purchase order missing branch; PO not found on expense

When an expense has a purchase order, Turbo sets the branch (branch) and the currency (currency) from the PO. A PO that is not found gives 'referenced purchase order not found'. A PO with no branch gives 'referenced purchase order is missing a branch'.

## Expense without a PO takes its branch from the job

Lines 106-129

Asked as: branch from job; referenced job not found; job has no branch; expense job branch missing

When an expense has a job and no PO, Turbo sets the branch from the job. A job that is not found gives 'referenced job not found'. A job with no branch gives 'referenced job is missing a branch'.

## Expense with no PO and no job uses default branch

Lines 130-153

Asked as: no default branch; admin profile missing default branch; expense branch missing no job; default branch error

When an expense has neither a PO nor a job, Turbo sets the branch from the default branch (default_branch) on the submitter's admin profile. Without a default branch Turbo refuses with 'your admin_profiles record is missing a default_branch'. Without an admin profile it reports an error looking up the admin profile.

## Mileage needs a valid vehicle insurance date

Lines 163-188

Asked as: mileage refused insurance; car insurance date missing; vehicle insurance expiry needed; mileage needs insurance date

For a Mileage expense, the admin profile must have a personal vehicle insurance expiry date (personal_vehicle_insurance_expiry). Without it Turbo refuses with 'cannot submit mileage expense: personal vehicle insurance expiry must be updated with a valid date'.

## Mileage refused after insurance expiry date

Lines 198-217

Asked as: mileage after insurance expired; insurance expired mileage claim; can't claim mileage expired insurance; mileage date past insurance

Turbo refuses a Mileage expense dated after the insurance expiry date, with 'cannot submit mileage expense: personal vehicle insurance expired on <date>'. An expense dated on the expiry date itself passes.

## Mileage total is calculated and vendor cleared

Lines 219-271

Asked as: mileage total wrong; mileage amount calculated; vendor cleared mileage; mileage rate calculation

For a Mileage expense, Turbo calculates the total (total) from the distance and the expense rate record, and clears the vendor. At commit it recalculates the mileage using the year's committed mileage.

## Mileage distance must be a whole number

Lines 233-248

Asked as: mileage decimal distance; distance not whole number; mileage km must be integer; fractional mileage distance

Turbo refuses a Mileage expense whose distance is not an integer, with 'distance must be an integer for mileage expenses'.

## Personal reimbursement requires profile permission

Lines 280-304

Asked as: personal reimbursement refused; can't submit personal reimbursement; personal reimbursement not enabled; reimbursement permission profile

Turbo refuses a PersonalReimbursement expense unless the admin profile allows personal reimbursement (allow_personal_reimbursement). The message is 'cannot submit personal reimbursement expense: personal reimbursement is not enabled for your profile'.

## Allowance total and description are set automatically

Lines 306-343

Asked as: allowance total wrong; breakfast lunch dinner allowance; allowance description automatic; allowance amount calculated

For an Allowance expense, Turbo sums the rates of the selected allowance types (Breakfast, Lunch, Dinner, Lodging) from the expense rate record. It sets the total to that sum, sets the description to 'Allowance for ' plus the types joined by commas, and clears the vendor.

## Expense currency falls back to home currency

Lines 353-362

Asked as: currency defaults home currency; no currency on expense; expense currency missing; fallback home currency

For a no-PO expense of type Allowance, FuelCard, Mileage or PersonalReimbursement, Turbo sets the currency to the home currency. If any expense still has no currency and a home currency exists, Turbo uses the home currency.

## Some payment types without a PO always use home currency

Lines 353-358

Asked as: mileage always home currency; fuel card currency locked; allowance currency home only; why is currency greyed out

For an expense with no PO, Turbo sets the currency to the home currency when the payment type is Allowance, FuelCard, Mileage or PersonalReimbursement.

## Currency must exist and foreign currency needs a rate

Lines 364-385

Asked as: currency not found expense; foreign currency no exchange rate; exchange rate missing; currency refused on expense

Turbo refuses an expense whose currency is not found, with 'referenced currency not found'. It refuses a foreign currency with no exchange rate, with 'selected foreign currency is missing an exchange rate'.

## Settlement fields are reset on save

Lines 387-398

Asked as: settled total reset; settled date cleared; settler removed on save; settled amount zero foreign

On every save Turbo clears the settler (settler) and settled date (settled). For a home-currency expense it sets the settled total (settled_total) to the total. For a foreign-currency expense with payment type Expense it leaves the settled total alone. For other foreign-currency expenses it sets the settled total to 0.

## Manager must be active and hold the tapr claim

Lines 402-470

Asked as: manager not active; manager missing tapr; approver lacks tapr claim; manager can't create expenses

Turbo refuses an expense when the manager is not an active user, with 'your manager must be an active user to create expenses'. It also refuses when the manager lacks the tapr claim, with 'the approver must have the tapr claim'.

## Bookkeeper can create expenses for another person only on a PO

Lines 505-534

Asked as: bookkeeper expense without PO; create expense for someone else; uid must equal my id; purchase order required bookkeeper

A bookkeeper (the book_keeper claim) may create an expense for another person only with a purchase order. Without one Turbo says 'purchase order is required for bookkeeper-created expenses'. A creator without the claim gets 'uid must be equal to the authenticated user's id'.

## Bookkeeper on-behalf expense PO and owner rules

Lines 536-562

Asked as: bookkeeper PO owner mismatch; PO owner must match uid; same user purchase order; PO not active bookkeeper

For a bookkeeper-created expense the PO must be Active ('purchase order must be active') and have a PO number ('purchase order must have a po number'). It needs an owner ('purchase order is missing an owner'). The expense's uid must match the PO owner ('uid must match the purchase order owner'). If the PO owner is the bookkeeper, Turbo says 'same-user purchase orders must use the regular expense flow'. The uid is set to the PO owner.

## Bookkeeper expenses limited to two payment types

Lines 564-572

Asked as: bookkeeper payment type; only on account or credit card; bookkeeper payment method refused; corporate card bookkeeper limit

A bookkeeper-created expense must have payment type OnAccount or CorporateCreditCard. Otherwise Turbo refuses with 'bookkeeper-created expenses must be OnAccount or CorporateCreditCard'.

## Expense editing requires sign-in

Lines 585-596

Asked as: authentication required expense; not signed in expense; can't save expense logged out; expense needs login

Turbo refuses to create or update an expense when no user is signed in, with 'authentication is required'.

## Expense's PO must exist

Lines 602-617

Asked as: purchase order not found expense; PO missing on expense; expense PO does not exist; wrong PO number expense

Turbo refuses an expense when its purchase order is not found, with 'purchase order not found'.

## New expense records its creator

Lines 622-632

Asked as: who created expense; creator not recorded; expense creator set to me; record who entered expense

On creation Turbo sets the creator (creator) to the signed-in user. Creating for another person is allowed only on the bookkeeper path.

## Only the creator may edit an expense, and uid and creator cannot change

Lines 634-691

Asked as: can't edit expense; only creator edits expense; uid cannot be changed; change expense owner

On update only the creator can edit the expense, otherwise Turbo says 'only the creator can edit this expense'. The uid cannot be changed ('uid cannot be changed'). The creator cannot be changed ('creator cannot be changed'). If the original record is missing, Turbo says 'cannot validate ownership for this expense update'. An expense created by a bookkeeper on behalf of someone else is re-validated under the bookkeeper rules on each update.

## Submitted expenses cannot be edited

Lines 693-719

Asked as: edit submitted expense; locked expense can't change; submitted expense edit blocked; change expense after submit

Turbo refuses to edit an expense that is submitted, with 'cannot edit submitted expense'. The check covers both the stored and the incoming record.

## Expense division must be active and branch usable

Lines 726-733

Asked as: division inactive expense; can't use branch on expense; expense branch not allowed; division not active

After cleaning, Turbo checks that the division on the expense is active. It also checks that the user may use the expense's branch, through EnsureUserCanUseRecordBranch (its rules are in another file).

## Attachment input handling and removal

Lines 759-780

Asked as: remove receipt upload; attachment field name; delete attached document expense; link source expense

The upload field is named 'attachment'. An empty attachment value with no new file and no source expense (source_expense) removes the attached document and its missing-receipt reason. A source expense id in the request is stored on the expense.

## Expense kind comes from its purchase order

Lines 798-803

Asked as: expense kind from PO; kind differs from PO; why kind not my choice; PO kind applied to expense

When an expense has a PO, Turbo sets the expense kind (kind) from the PO's kind, not from the user. The PO's job decides the kind when the PO has none.

## PO must be active for an expense

Lines 805-816

Asked as: closed PO expense; PO not active submit; anyone can use PO; expense on inactive purchase order

Turbo refuses an expense against a purchase order whose status is not Active, with 'purchase order must be active'. Any person may submit an expense against any active PO if the other rules pass.

## Cumulative PO total excludes this expense's stored row

Lines 818-842

Asked as: cumulative PO total wrong; PO running total double counted; cumulative purchase order limit; PO total includes old amount

For a Cumulative PO, Turbo adds the totals of the other expenses on the PO to this expense's total. On update it leaves out the stored row of this expense so the old total is not counted twice.

## Attachment document is server-managed

Lines 859-867

Asked as: attachment document changed; can't set receipt document; receipt cleared on save; missing receipt reason removed

A client cannot set the attachment document (attachment_document). On create or explicit removal Turbo clears it and the missing-receipt reason. On a normal update it keeps the stored values before validation.

## Attachment handling by payment type and intent

Lines 888-913

Asked as: receipt not saved; attachment skipped for payment type; vendor invoice details saved; receipt kept on update

Payment types that skip attachments never keep an attachment document or a missing reason. With a file or source expense, Turbo links the resolved document and clears the missing reason. Without either, an update keeps the original document and a new expense gets none. Vendor invoice details are then saved on the document.
