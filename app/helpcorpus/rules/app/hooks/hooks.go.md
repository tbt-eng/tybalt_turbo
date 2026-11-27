<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"04fc27739e91941cfcc6eda76d5fe450ffe3bacfc55788092a00e5c55f5bc099","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors are returned with their own status and message

Lines 19-25

Asked as: hook error message shown; wrong error status returned; custom error message to user; check failed error text

When a Turbo check fails with a hook error, Turbo returns that error's message and status to the user. Other errors are returned unchanged.

## Purchase order notification is sent after a successful save

Lines 27-72

Asked as: purchase order notification sent; PO notification timing; notification after saving PO; notification when PO fails

When saving a purchase order creates a notification, Turbo sends it right after the save. If the save fails, Turbo deletes the pending notification so it is never delivered. A failure to send does not undo the save.

## Notification delivery errors do not undo a saved purchase order

Lines 27-42

Asked as: PO saved but no notification; notification failed to send; email not delivered PO; saved PO still there

If Turbo cannot send a purchase order notification after the save, it logs the error. The saved purchase order stays saved.

## Expense editing switch blocks deletes and vendor changes

Lines 83-92

Asked as: cannot delete expense; expenses editing disabled; can't edit vendors; vendor changes blocked

When expenses editing is disabled, Turbo refuses to delete an expense or a purchase order. It also refuses to create, update or delete a vendor. It returns the expenses-editing-disabled error (ErrExpensesEditingDisabled).

## Microsoft login syncs identity and onboards in one transaction

Lines 126-281

Asked as: Microsoft login setup; new user onboarding Microsoft; name email sync login; no profile after Microsoft sign in

After a Microsoft login, Turbo syncs the user's name and email (the sync runs on later logins too). It makes sure an admin profile exists. It does not create a profile (profiles record), because the manager must be chosen in the app. A first-time user also gets an initial time claim. The whole flow runs in one transaction, so a failure leaves no half-onboarded user. Only Microsoft logins get this handling.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: Microsoft account relink; existing user Microsoft login; link Microsoft to old account; duplicate user Microsoft login

When PocketBase finds no user by provider id, Turbo looks for an existing user to relink. If it finds one, it relinks the Microsoft account to that user and treats the login as a returning user, not a new one. Failures give "failed to look up existing Microsoft user" or "failed to relink Microsoft account".

## Microsoft first login builds the new user's details

Lines 179-224

Asked as: blank email Microsoft login; username from email; first login details filled; userPrincipalName used as email

On a first Microsoft login, Turbo fills a blank email from the Microsoft mail, or from userPrincipalName when mail is empty. It fills a blank name from the Microsoft full name. It builds a unique username from the email when no username is given. Details already supplied are kept.

## Inactive accounts cannot sign in

Lines 304-316

Asked as: account inactive login; can't sign in deactivated; login blocked inactive user; disabled user cannot log in

When someone signs in, Turbo checks their admin profile. If the profile is not active (active is false), Turbo refuses the login with "this account is inactive", even when the credentials are valid. This check runs on every login.

## Non-Microsoft login without admin profile is refused

Lines 317-323

Asked as: account setup incomplete; login refused no profile; no admin profile sign in; can't log in after password

If a user signs in by a method other than OAuth2 and has no admin profile, Turbo refuses with "account setup incomplete". It does not create the profile.

## OAuth2 login repairs a missing admin profile

Lines 321-329

Asked as: OAuth2 login no profile; onboarding failed login; profile created on login; missing admin profile fix

If a user signs in by OAuth2 and has no admin profile, Turbo creates one and lets the login continue. If the creation fails, Turbo refuses with "failed to complete user onboarding".

## Time entry edits are gated by the time editing setting

Lines 333-379

Asked as: cannot delete time entry; time editing disabled; can't remove timesheet entry; time amendment delete blocked

Turbo refuses to delete a time entry or a time amendment when time editing is disabled, and returns the time-editing-disabled error (ErrTimeEditingDisabled). Creating and updating run through ProcessTimeEntry and ProcessTimeAmendment, whose rules are in other files.

## Purchase order saves are processed before saving

Lines 380-404

Asked as: PO save checks; purchase order rejected; PO not saved error; notification removed after failed save

Turbo runs its purchase order checks (ProcessPurchaseOrder) on every create and update of a purchase order. If a check fails, Turbo returns the error and does not save the purchase order. If the save fails after the checks, Turbo deletes the notification it created, so no notice goes out for a PO that was not saved.

## Project, rate sheet and client setup hooks are added

Lines 405-411

Asked as: setup hooks startup; rate sheet hooks; project workflow hooks; client setup hooks added

Turbo adds the client setup, project workflow, job rate sheet, rate sheet draft, rate sheet code and rate sheet entry hooks when it starts. Their rules are in other files.

## Job saves run the job checks

Lines 412-423

Asked as: job save rejected; job checks failing; can't save job; job validation error

Turbo runs its job checks (ProcessJob) on every create and update of a job. If a check fails, Turbo returns the error and does not save the job.

## Profile saves run the profile checks

Lines 424-436

Asked as: profile save error; profile checks failing; can't save profile; profile validation fails

Turbo runs its profile checks (ProcessProfile) on every create and update of a profile. If a check fails, Turbo returns the error and does not save the profile.

## Currency saves run the currency checks

Lines 437-449

Asked as: currency save error; currency checks failing; can't save currency; currency validation fails

Turbo runs its currency checks (ProcessCurrency) on every create and update of a currency. If a check fails, Turbo returns the error and does not save the currency.

## New admin profiles need a permission check first

Lines 450-466

Asked as: new admin profile permission; can't create admin profile; admin tokens revoked; deactivated admin signed out

Creating an admin profile first runs the request permission check (EnforceAdminProfileRequestPermissions). Updating does not run it. Turbo revokes the user's tokens when an admin profile is deactivated (RevokeTokensOnDeactivation).

## Admin profile updates run checks; deactivation revokes tokens

Lines 460-466

Asked as: admin profile update checks; deactivate admin logout; sign-in tokens revoked; admin profile save error

Turbo runs its admin profile checks (ProcessAdminProfile) on every update of an admin profile. When an admin profile is saved, Turbo also runs RevokeTokensOnDeactivation. The name of that hook indicates it revokes sign-in tokens when a profile is deactivated.

## Client and client note saves run their checks

Lines 467-486

Asked as: client save rejected; client note not saved; client checks failing; note validation error

Turbo runs its client checks (ProcessClient) on every create and update of a client. It runs its note checks (ProcessClientNote) on every new client note. A failed check stops the save. Turbo has no update hook for client notes in this file.

## Purchase order approver limits run checks on save

Lines 487-499

Asked as: approver limit error; PO approver save fails; approver properties rejected; approver check failed

Turbo runs its checks (ProcessPOApproverProps) on every create and update of a purchase order approver's properties. A failed check stops the save.

## Expense saves run in one transaction with upload cleanup

Lines 501-532

Asked as: expense save failed; expense documents left behind; uploads not cleaned up; expense transaction rollback

Creating or updating an expense runs in one transaction. If it fails, Turbo cleans up the uploaded expense documents (CleanupExpenseDocumentUploads). The expense rules are in ProcessExpense.

## Deleting a purchase order is blocked when expense editing is off

Lines 533-545

Asked as: cannot delete PO; PO delete blocked; expenses editing off delete; vendor edit blocked

Turbo refuses to delete a purchase order while the expenses editing setting is off. It returns the expenses-editing-disabled error (ErrExpensesEditingDisabled). The same switch also blocks creating, updating and deleting vendors.

## New rate sheets must pass the effective date check

Lines 547-562

Asked as: rate sheet effective date error; new rate sheet rejected; effective date invalid; rate sheet date check

Creating a rate sheet first checks its effective date (ValidateRateSheetEffectiveDate). Updating a rate sheet does not run that check.

## Changing job categories marks the job as not imported

Lines 565-578

Asked as: job marked not imported; category change job reimport; job category removed import flag; job status reset

Creating or deleting a category marks its job as not imported (MarkJobNotImported), after the change succeeds.

## Client or contact edits mark related jobs not imported

Lines 579-601

Asked as: client edit reimport jobs; contact change jobs not imported; jobs need re-import after client change; _imported flag jobs

When a client changes in a meaningful way (other than the _imported field), Turbo marks jobs that use it as client or job owner as not imported. When a client contact changes the same way, Turbo marks jobs that use it as contact. A change to _imported alone does nothing.

## Notifications update their status on save

Lines 603-617

Asked as: notification status update; notification not saving; status updated on notification; notification save failed

Turbo runs WriteStatusUpdated every time a notification is created or updated, before the save. If it fails, Turbo does not save the notification.

## Time sheet reviewer must be an active user

Lines 630-651

Asked as: reviewer not active; time sheet share reviewer error; reviewer inactive user; can't share timesheet

When an approver shares a time sheet, the chosen reviewer must be an active user. Otherwise Turbo refuses with "reviewer must be an active user" (the reviewer field says "the selected reviewer is not an active user"). If the check itself fails, Turbo says "failed to check reviewer active status".

## Absorb commit blocked while a parent absorb exists

Lines 655-684

Asked as: absorb commit blocked; parent absorb exists; can't commit absorb; commit absorb error

Committing an absorb action is refused while an absorb action exists for a parent collection. The message is "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first".

## Vendor absorb commit needs expenses editing enabled

Lines 658-663

Asked as: vendor absorb commit blocked; vendor absorb expenses off; can't commit vendor absorb; vendor merge refused

Committing a vendor absorb is refused when expenses editing is disabled.

## Absorb commit needs the collection's claims

Lines 686-706

Asked as: absorb permission denied; no permission commit absorb; collection absorb claims; can't commit absorb permission

For collections with their own claim sets, committing an absorb requires the claims set for that collection. Without them Turbo refuses with "you do not have permission to commit a X absorb", where X is the collection name.

## Absorb commit saves history and settles proposals

Lines 708-750

Asked as: absorb history saved; absorb proposals settled; absorb commit outcome notice; absorb rollback failed delete

A commit runs in one transaction. It first copies the action into absorb history with who absorbed and who committed. Then it runs the collection's commit steps. For collections whose absorb can be proposed, it settles the proposals: it accepts the one that the action carries out and closes those that name an absorbed record. It then queues an outcome notification for each settled proposal. A failed delete leaves nothing behind.

## Sharing a time sheet notifies the reviewer

Lines 758-789

Asked as: reviewer notified timesheet; share timesheet notification; reviewer not notified; timesheet share email

When a reviewer is added to a time sheet, Turbo queues a notification for them, naming the time sheet's approver as the sharer. If the notification fails or the time sheet is not found, the share still succeeds.
