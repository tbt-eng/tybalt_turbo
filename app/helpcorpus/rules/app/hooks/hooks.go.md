<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"b86170e3a64e292093467c9a7129e23b1f8286cdbda7b59512aa2598c5197793","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Purchase order save sends its notification afterwards

Lines 27-72

Asked as: purchase order notification sent; PO notification after save; notification deleted when save fails; PO email sent too early

When a purchase order is created or updated, Turbo may create a notification first. Turbo sends it right after the save succeeds. If the save fails, Turbo deletes the pending notification so the cron does not send it for a change that never happened. A failure to send is logged and does not undo the save.

## Microsoft sign-in sets up a new user

Lines 126-224

Asked as: new Microsoft user setup; Microsoft login creates account; email filled from Microsoft; first sign in onboarding

For a new user who signs in with Microsoft, Turbo fills in an empty email from the Microsoft mail, with userPrincipalName as the fallback. It fills in an empty name from the Microsoft full name. It builds a unique username when none is given. Turbo does all of this in one transaction, so a failure leaves no half-onboarded user. Only Microsoft logins get this treatment.

## Microsoft sign-in relinks an existing user

Lines 150-177

Asked as: existing user Microsoft account link; Microsoft login matched to old user; relink Microsoft account; duplicate user Microsoft sign in

When PocketBase finds no user by provider id, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, Turbo relinks the Microsoft account to that user and treats the login as an existing user's login. Otherwise Turbo creates a new user. Failures give 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## Microsoft sign-in syncs identity and onboards

Lines 226-280

Asked as: Microsoft sync name and email; admin profile created at login; sign in identity sync error; first login time claim

After a Microsoft sign-in, Turbo syncs the user's name and email from Microsoft, and fixes the verified flag. A first-time user also gets an initial time claim. Every Microsoft user ends up with an admin profile, but Turbo creates no profile record (profiles), because the manager must be collected later in the app. A failure in any of these steps gives an error, such as 'failed to sync Microsoft identity', and undoes the sign-in work.

## Inactive accounts cannot sign in

Lines 304-316

Asked as: inactive account cannot log in; deactivated user sign in blocked; account disabled login refused; can't sign in deactivated

When a user's admin profile is not active, Turbo refuses the sign-in with 'this account is inactive'. Turbo refuses even when the credentials are valid, and it issues no token. This check runs on every sign-in.

## Sign-in without an admin profile

Lines 317-330

Asked as: account setup incomplete; no admin profile login error; login fails missing profile; OAuth sign in profile missing

A user with no admin profile who signs in by a method other than OAuth2 gets the error 'account setup incomplete'. Turbo does not create the profile for them. A user who signs in by OAuth2 gets the missing admin profile created, and the sign-in continues. If Turbo cannot create it, the error is 'failed to complete user onboarding'.

## Creating or editing a time entry runs time entry checks

Lines 333-344

Asked as: time entry not saving; timesheet entry rejected; time entry validation error; entry check fails on save

Before Turbo saves a new or edited time entry (time_entries), it runs its time entry checks (ProcessTimeEntry). If a check fails, Turbo returns that error and does not save the entry. This file does not show what the checks are.

## Time entry and amendment deletes need time editing on

Lines 345-379

Asked as: cannot delete timesheet entry; time editing disabled delete; can't remove time entry; delete time amendment blocked

Turbo refuses to delete a time entry or a time amendment when time editing is disabled. The time editing feature flag decides this (ErrTimeEditingDisabled). Creating and updating run through ProcessTimeEntry and ProcessTimeAmendment instead.

## Creating or editing a time amendment runs amendment checks

Lines 357-368

Asked as: amendment not saving; time amendment rejected; amendment validation error; time change request fails

Before Turbo saves a new or edited time amendment (time_amendments), it runs its amendment checks (ProcessTimeAmendment). If a check fails, Turbo returns that error and does not save the amendment. This file does not show what the checks are.

## Purchase order create and update run PO checks first

Lines 381-404

Asked as: purchase order won't save; PO save rejected; PO validation failed; new PO error on save

Before Turbo saves a new or edited purchase order, it runs its PO checks (ProcessPurchaseOrder). If a check fails, Turbo returns that error and does not save the PO.

## Job, profile and currency saves run their own checks

Lines 405-449

Asked as: job won't save; profile save error; currency save fails; job validation error

Before Turbo saves a new or edited job, profile or currency, it runs that record's checks (ProcessJob, ProcessProfile, ProcessCurrency). If a check fails, Turbo returns the error and does not save the record. Job saves also use the project workflow and rate sheet hooks, which this file does not detail.

## Admin profile creation needs permission checks

Lines 450-468

Asked as: who can create admin profile; admin profile permission denied; new admin profile blocked; revoked tokens after profile change

Before creating an admin profile through a request, Turbo checks the requester's permissions (EnforceAdminProfileRequestPermissions). Updates skip this check in this hook. Turbo also revokes a user's tokens when a profile is created inactive, deactivated or deleted.

## Admin profile update runs checks without the create permission check

Lines 451-465

Asked as: admin profile edit error; can't update admin profile; admin profile create permission error; profile update checks fail

Creating an admin profile first runs the request permission check (EnforceAdminProfileRequestPermissions), then the admin profile checks (ProcessAdminProfile). Updating an admin profile runs only the admin profile checks. Neither runs once the request permission check fails.

## Client, client note and invoice saves run their own checks

Lines 469-519

Asked as: client won't save; invoice save error; client note not saving; invoice line delete blocked

Before Turbo saves a client, client note, client invoice or invoice line, it runs that record's checks. These cover client create and update, client note create only, invoice create, update and delete, and invoice line create and delete. If a check fails, Turbo returns the error and does not make the change.

## PO approver properties run checks on create and update

Lines 520-532

Asked as: PO approver settings not saving; approver props error; approver setup rejected; PO approver save fails

Before Turbo saves PO approver properties (po_approver_props), on create or update, it runs the PO approver checks (ProcessPOApproverProps). If a check fails, Turbo returns the error and does not save the record.

## Expense saves run in a transaction with upload cleanup

Lines 533-565

Asked as: expense save rolled back; expense upload deleted on failure; receipt removed after error; expense transaction failed

Turbo creates and updates an expense in one transaction. If the save fails, Turbo cleans up the uploaded expense documents.

## Expense saves run expense checks

Lines 534-565

Asked as: expense not saving; expense rejected on save; expense validation error; expense check fails

Before Turbo saves a new or edited expense, it runs its expense checks (ProcessExpense). If a check fails, Turbo returns the error and does not save the expense.

## Deletes of expenses, POs and vendor changes need expenses editing on

Lines 566-578

Asked as: cannot delete expense; expenses editing off; can't delete purchase order; vendor change blocked

Turbo refuses to delete an expense or a purchase order, and refuses to create, update or delete a vendor, when expenses editing is disabled (ErrExpensesEditingDisabled).

## New rate sheet must pass the effective date check

Lines 580-595

Asked as: rate sheet date rejected; rate sheet effective date error; new rate sheet won't save; effective date invalid

When someone creates a rate sheet, Turbo first checks its effective date (ValidateRateSheetEffectiveDate), then runs the rate sheet checks. An update runs only the rate sheet checks and does not check the effective date again.

## Editing records marks jobs as not imported

Lines 597-634

Asked as: job marked not imported; jobs need reimport after edit; client change resets import; category edit import status

Creating or deleting a category marks its job as not imported. Updating a client marks the jobs that use it as client or job owner as not imported. Updating a client contact does the same for jobs that use it as contact. Client and contact updates do this only when something other than the _imported field changed.

## Notification status is stamped when a notification is saved

Lines 636-650

Asked as: notification save fails; notification status error; notification not updating; notification status stamp

Whenever Turbo creates or updates a notification record, it first updates the notification's status data (WriteStatusUpdated). If that fails, the save fails. Notifications cannot be reached through the API.

## Time sheet reviewer must be active

Lines 663-684

Asked as: reviewer not active; can't share timesheet; reviewer inactive error; share time sheet fails

When a manager shares a time sheet, the selected reviewer must be an active user. Otherwise Turbo refuses with 'reviewer must be an active user' ('the selected reviewer is not an active user').

## Vendor absorb commit needs expenses editing on

Lines 691-696

Asked as: vendor absorb commit blocked; can't commit vendor merge; vendors absorb refused; expenses editing commit

Turbo refuses to commit a vendors absorb when expenses editing is disabled.

## Absorb commit is blocked by a parent absorb

Lines 698-717

Asked as: absorb commit blocked; parent absorb exists; can't commit absorb; commit or undo first

Turbo refuses to commit an absorb while an absorb for a parent collection exists. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'. The parent collections come from absorb.GetCommitBlockers.

## Absorb commit needs the collection's claims

Lines 719-739

Asked as: no permission to commit absorb; absorb commit access denied; can't commit absorb claims; absorb permission error

For a collection with its own claim sets, a committer must hold the required claims. Otherwise Turbo refuses with 'you do not have permission to commit a X absorb'. The other collections rely on the collection rule alone.

## Absorb commit saves history and settles proposals

Lines 741-783

Asked as: absorb history missing; absorb proposal settled; absorb commit partial failure; proposal closed after absorb

A commit runs in one transaction. It copies the action into absorb history, with who absorbed and who committed. It runs the collection's own commit steps. For collections whose absorb can be proposed, it accepts the proposal the action carries out and closes the proposals that name an absorbed record. A failed step leaves nothing behind. Turbo then queues an outcome notification for each settled proposal.

## Sharing a time sheet notifies the reviewer

Lines 791-822

Asked as: reviewer not notified; time sheet share email; shared timesheet notification; share notification sent

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to the reviewer, in the name of the time sheet's approver. If the notification fails, the share still succeeds.
