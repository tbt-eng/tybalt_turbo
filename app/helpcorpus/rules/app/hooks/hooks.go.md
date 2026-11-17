<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"0f451deafa9e0ad6ddd0103a2b1f77f09864bd77ff90233dfda33edaeb906eef","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors return their message to the requester

Lines 19-25

Asked as: hook error message; error message not shown; save hook error response; why does save show wrong error

When a save hook for a record fails with a hook error (HookError), Turbo sends the error's own status and message back to the person making the request. Other errors pass through unchanged.

## Purchase order notification sent only after save

Lines 27-72

Asked as: purchase order notification sent; PO notification after failed save; notification sent before save; PO email went out

When a purchase order save creates a notification, Turbo sends it right after the save succeeds. If the save fails, Turbo deletes the pending notification so the cron does not deliver it for a change that never happened. A failed send does not undo the save.

## Expenses, purchase orders and vendors need expenses editing on

Lines 81-92

Asked as: cannot delete expense; can't delete purchase order; vendor edit refused; expenses editing turned off

Deleting an expense or a purchase order, and creating, updating or deleting a vendor, is refused when the expenses editing feature is off (ErrExpensesEditingDisabled).

## Microsoft login ensures an admin profile but no profile

Lines 124-135

Asked as: microsoft login onboarding failed; cannot log in with microsoft; microsoft account setup error; new microsoft user no profile

Every successful Microsoft login ensures the user has an admin profile. Turbo does not create a profile, because the manager must be chosen in the app. If onboarding fails, the login is refused with 'failed to complete Microsoft user onboarding'. All the Microsoft login steps run in one transaction, so a failure leaves no partly onboarded user.

## Only Microsoft logins get the Microsoft onboarding

Lines 126-129

Asked as: google login onboarding; other oauth provider login; non microsoft oauth login; login via other provider

Turbo applies its custom Microsoft onboarding only when the OAuth provider is Microsoft (compared without regard to case). Logins through other OAuth providers skip it and continue normally.

## Microsoft login relinks an existing user instead of creating one

Lines 150-177

Asked as: microsoft account already exists; duplicate microsoft user; link existing account microsoft; failed to relink microsoft

When PocketBase finds no user by provider ID, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the Microsoft account to that user and treats the login as an existing user, not a new one. If the lookup fails, Turbo refuses with 'failed to look up existing Microsoft user'. If the relink fails, it refuses with 'failed to relink Microsoft account'.

## Microsoft first login creates user with derived email, name and username

Lines 179-224

Asked as: microsoft new user email blank; empty name after microsoft signup; userPrincipalName used as email; microsoft username generated

When a new user signs in with Microsoft, Turbo fills any blank email, name and username before the account is created. The email is the Microsoft email (mail), or the userPrincipalName when mail is blank. The name is the full name Turbo derives from the Microsoft data. The username is built as a unique name from the Microsoft data.

## Microsoft login failure for username creation

Lines 207-214

Asked as: cannot create microsoft username; username already taken microsoft; prepare microsoft username failed; microsoft signup username error

If Turbo cannot build a unique username for a new Microsoft user, the login is refused with 'failed to prepare Microsoft username'.

## Microsoft login syncs name and email on every login

Lines 226-259

Asked as: name changed in microsoft; email not updated from microsoft; microsoft details out of date; sync microsoft name

After each successful Microsoft login, Turbo syncs the user's name and email from Microsoft (syncMicrosoftUserIdentity). If the sync fails, the login is refused with 'failed to sync Microsoft identity'. The username stays as created and does not change.

## New Microsoft users get an initial time claim

Lines 226-267

Asked as: new user no time claim; initial claim missing microsoft; time claim not assigned; first login time claim

When a Microsoft login creates a new user, Turbo assigns the user an initial time claim (ensureInitialTimeClaim). Returning users do not get it again. If the assignment fails, the login is refused with 'failed to assign initial time claim'.

## Inactive accounts cannot log in

Lines 304-316

Asked as: account locked out; cannot sign in disabled user; user deactivated login fails; inactive user cannot login

A login is refused with 'this account is inactive' when the user's admin profile (active) is not set. This applies to every login method, even with valid credentials.

## Missing admin profile on login

Lines 317-329

Asked as: no profile after login; login says account setup incomplete; missing admin profile; cannot finish setup on login

If a user has no admin profile and logs in by OAuth2, Turbo creates the admin profile and lets the login continue. If creation fails, it refuses with 'failed to complete user onboarding'. If the user logs in by any other method, Turbo refuses with 'account setup incomplete'.

## Time entry saves run the time entry checks

Lines 333-344

Asked as: time entry not saving; time entry validation error; timesheet entry rejected; time entry checks failed

Every create or update of a time entry goes through the time entry processing step (ProcessTimeEntry) before it is saved. If that step fails, Turbo shows its error and does not save.

## Time entry and time amendment deletes need time editing on

Lines 345-379

Asked as: cannot delete time entry; can't remove timesheet line; delete time amendment refused; time editing disabled delete

Deleting a time entry or a time amendment is refused when the time editing feature is off (ErrTimeEditingDisabled). Creating or updating them goes through their own checks (ProcessTimeEntry, ProcessTimeAmendment).

## Time amendment saves run the time amendment checks

Lines 357-368

Asked as: amendment not saving; time amendment rejected; amendment checks failed; time correction not saved

Every create or update of a time amendment goes through the time amendment processing step (ProcessTimeAmendment) before it is saved. If that step fails, Turbo shows its error and does not save.

## Purchase order saves run PO checks and create a notification

Lines 381-404

Asked as: purchase order not saving; PO validation error; PO save rejected; purchase order checks failed

Every create or update of a purchase order goes through the purchase order processing step (ProcessPurchaseOrder) before it is saved. That step may create a notification. If it fails, Turbo shows its error and does not save.

## Failed purchase order save deletes its pending notification

Lines 386-389

Asked as: stray PO notification; notification for unsaved PO; PO change notification still sent; cron sends failed PO notice

When a purchase order save fails after the processing step has created a notification, Turbo deletes that notification so the cron does not send it for a PO change that was never saved.

## Job, profile, currency, client and client note saves run their checks

Lines 411-484

Asked as: job not saving; client save error; currency or profile rejected; client note not saved

Creates and updates of jobs, profiles, currencies and clients, and creates of client notes, go through their own processing steps (ProcessJob, ProcessProfile, ProcessCurrency, ProcessClient, ProcessClientNote) before saving. Client notes have no update hook here. A failure in a step stops the save and shows its error.

## Admin profile create checks request permissions first

Lines 450-464

Asked as: cannot create admin profile; admin profile permission denied; new admin profile error; admin profile update fails

Before an admin profile is created, Turbo checks the request's permissions (EnforceAdminProfileRequestPermissions). If that fails, Turbo returns its error as is. Then it runs the admin profile checks (ProcessAdminProfile). An update runs only the admin profile checks.

## PO approver properties saves run their checks

Lines 486-497

Asked as: PO approver settings not saving; approver props error; approver properties rejected; po approver save fails

Every create or update of a PO approver properties record (po_approver_props) goes through its processing step (ProcessPOApproverProps) before saving. A failure stops the save.

## Expense saves run in a transaction and clean up uploads

Lines 499-530

Asked as: expense save failed uploads; expense receipts left behind; expense upload cleanup; expense not saved attachments

Creating or updating an expense runs in one transaction. If it fails, Turbo cleans up the uploaded expense documents (CleanupExpenseDocumentUploads).

## Rate sheet effective date is checked on create only

Lines 545-560

Asked as: rate sheet date rejected; rate sheet effective date error; cannot create rate sheet date; rate sheet start date wrong

Creating a rate sheet checks its effective date (ValidateRateSheetEffectiveDate) before the other checks. Updating a rate sheet does not run this check.

## Changes to categories, clients and contacts mark jobs not imported

Lines 562-599

Asked as: job marked not imported; jobs need reimport after edit; category change jobs unimported; client change resets import

Creating or deleting a category marks its job as not imported. Updating a client with a change other than to _imported marks the jobs that use it as client or job owner. Updating a client contact the same way marks the jobs that use it as contact.

## Time sheet reviewer must be an active user

Lines 628-647

Asked as: cannot share timesheet; reviewer not allowed; share timesheet inactive user; reviewer cannot be added

Sharing a time sheet by adding a reviewer is refused with 'reviewer must be an active user' (reviewer_not_active: 'the selected reviewer is not an active user') when the reviewer is not active.

## Vendor absorb commit needs expenses editing on

Lines 656-661

Asked as: vendor absorb commit refused; cannot commit vendor absorb; vendors absorb blocked; absorb vendors editing off

Committing a vendors absorb is refused when the expenses editing feature is off.

## Absorb commit blocked while a parent absorb exists

Lines 663-682

Asked as: absorb commit blocked; cannot commit absorb parent; parent absorb blocking commit; commit absorb error message

Committing an absorb is refused while an absorb for a parent collection exists. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'. Turbo takes the parents from absorb.GetCommitBlockers.

## Absorb commit needs the collection's claims

Lines 684-704

Asked as: no permission to commit absorb; absorb commit forbidden; cannot commit absorb claims; absorb access denied

For a collection with its own claim sets, committing an absorb needs the required claims (absorb.HasRequiredClaims). Without them Turbo refuses with 'you do not have permission to commit a X absorb'.

## Absorb commit runs its commit step in a transaction with the delete

Lines 706-718

Asked as: absorb deleted despite error; commit step failed absorb; absorb not removed on failure; absorb commit rollback

If the collection has a commit step (OnCommit), Turbo runs it and the delete of the absorb action in one transaction. If the step fails, the absorb action is not deleted.

## New reviewer is notified that a time sheet is shared

Lines 726-757

Asked as: reviewer not notified; shared timesheet email reviewer; reviewer notification missing; who shared my timesheet

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to that reviewer, naming the time sheet's approver as the sharer. A failure to queue it is logged and does not block the share.
