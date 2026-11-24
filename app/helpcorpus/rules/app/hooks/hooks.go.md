<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"6e85af51d327a044c4d00cdfebbb7b0fa0c24d3ecb450ad14bb11fb2549db282","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## PO notification send failure does not undo the save

Lines 27-42

Asked as: PO saved but email failed; notification failed after purchase order; PO kept though send failed; approval email not sent; retry failed notification later

After a purchase order saves, Turbo sends the notification it created at once. If sending fails, Turbo only logs the error and keeps the saved PO. The notification stays for later delivery (sendNotificationAfterSave).

## Failed purchase order save deletes its pending notification

Lines 44-72

Asked as: PO failed but notification still sent; stop notification for unsaved PO; PO error got email anyway; cron sends notice for failed PO; remove pending PO notification

When a purchase order create or update fails after its checks pass, Turbo deletes the notification that the checks created. This stops the cron from later sending a notification for a PO write that never saved (deleteOrphanedNotification).

## Expenses editing flag off refuses the action

Lines 81-92

Asked as: cannot edit expenses; expenses locked; expense editing turned off; can't change expense; expenses feature disabled error

When the expenses editing feature flag is off, Turbo refuses the action with ErrExpensesEditingDisabled. If Turbo cannot read the flag, it returns that error instead (checkExpensesEditing).

## Microsoft returning-user sync keeps username and only updates name and email

Lines 95-122

Asked as: Microsoft login changed my username; name not updating from Microsoft; email not updated after Microsoft sign-in; Microsoft sync overwrote username; returning Microsoft user details

On later Microsoft sign-ins Turbo syncs name and email from Microsoft, but it does not rename the username after creation. It changes the email only when no other user has that email. It creates an admin profile only if one is missing and never refreshes it from Microsoft.

## Microsoft onboarding is all-or-nothing

Lines 124-125

Asked as: half created Microsoft account; Microsoft signup left partial user; sign-in failed midway user created; no manager profile created at signin; Microsoft onboarding rolled back

Turbo runs the whole Microsoft sign-in in one transaction. If creating the admin profile or updating the user after save fails, no partly onboarded user remains. Turbo creates no profile record at sign-in, because the manager must be chosen later in the app.

## Only Microsoft OAuth logins get the custom onboarding

Lines 127-129

Asked as: custom onboarding for Google login; non Microsoft provider signup; Google sign-in skips setup; login provider case onboarding; other OAuth provider onboarding

The Microsoft onboarding steps run only when the OAuth provider is Microsoft, compared without regard to case. Other providers go straight through.

## Microsoft sign-in can relink an existing user

Lines 150-177

Asked as: Microsoft account linked to wrong user; existing account Microsoft sign in; same person new Microsoft login; reconnect Microsoft to old account; duplicate account Microsoft login

When PocketBase finds no user by provider id on a Microsoft login, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If one is found, Turbo relinks the Microsoft account to that user and treats the login as an existing user, not a new one. If the lookup or relink fails, the login is refused.

## Microsoft new-user onboarding failures give specific messages

Lines 160-276

Asked as: Microsoft sign-in error message; failed to relink Microsoft; Microsoft username error; failed to sync Microsoft identity; initial time claim error

A failed Microsoft sign-in step gives one of these messages: "failed to look up existing Microsoft user", "failed to relink Microsoft account", "failed to prepare Microsoft username", "failed to sync Microsoft identity", "failed to assign initial time claim" or "failed to complete Microsoft user onboarding".

## Microsoft first sign-in builds email, name and username

Lines 179-224

Asked as: username from Microsoft email; new Microsoft user name blank; userPrincipalName used as email; Microsoft first login fields; username already taken Microsoft

A brand-new user who signs in with Microsoft gets email from the mail field, or from userPrincipalName when mail is blank. Name comes from the full name (given name plus surname, then display name). Username is built from the email's local part and made unique. Only blank values are filled.

## Microsoft login syncs identity and onboards the user

Lines 226-280

Asked as: Microsoft login creates admin profile; new user missing manager profile; Microsoft login user setup; half onboarded after Microsoft; Microsoft signin sets up account

After a successful Microsoft login Turbo syncs the user's name and email from Microsoft, and ensures the user has an admin profile. For a new user it also assigns the initial time claim. Turbo does not create a profile, because the manager must be chosen later in the app. All steps run in one transaction, so a failure leaves no half-onboarded user.

## New Microsoft user gets an initial time claim

Lines 226-267

Asked as: new user no initial time claim; time claim missing after signup; relinked user time claim; first time claim Microsoft; new account no time access

Turbo assigns an initial time claim only to a user that the Microsoft sign-in just created as new. A user relinked to an existing account does not get one (ensureInitialTimeClaim).

## Login backstop gives specific errors

Lines 304-330

Asked as: account setup incomplete message; no profile login error; can't log in no admin profile; user onboarding failed login; login says account inactive

A user with no admin profile who signs in by a method other than OAuth2 gets the error "account setup incomplete". If Turbo cannot create the missing admin profile for an OAuth2 user, the error is "failed to complete user onboarding". The inactive-account refusal reads "this account is inactive".

## Inactive accounts cannot sign in

Lines 308-316

Asked as: inactive user cannot log in; deactivated account login blocked; correct password still refused; account disabled sign in; this account is inactive

A user whose admin profile is not active (active is false) cannot sign in even with valid credentials. Turbo refuses with 'this account is inactive'.

## Missing admin profile at login

Lines 317-329

Asked as: login without admin profile; missing profile after OAuth login; create admin profile at login; account setup incomplete; non OAuth login refused

If a user has no admin profile and signs in with OAuth2, Turbo creates the admin profile and lets the login continue. If creation fails, it refuses with 'failed to complete user onboarding'. A non-OAuth login with no admin profile is refused with 'account setup incomplete'.

## Time entry and time amendment saves are checked

Lines 333-368

Asked as: time entry checks on save; time amendment validation; timesheet entry not saving; amendment rejected on save; time entry validation errors

Turbo runs its time entry checks (ProcessTimeEntry) on every create and update of a time entry. It runs the time amendment checks (ProcessTimeAmendment) on every create and update of a time amendment.

## Deleting time entries or amendments needs time editing enabled

Lines 345-379

Asked as: cannot delete time entry; delete timesheet line blocked; time editing turned off; can't remove time amendment; time editing disabled error

Turbo refuses to delete a time entry or a time amendment when the time editing setting is off, and returns the time-editing-disabled error (ErrTimeEditingDisabled).

## Purchase order saves run checks and send a notification

Lines 381-404

Asked as: PO save email sent; purchase order notification on save; PO checks run on save; send email when PO saved; PO failed email removed

Turbo runs its purchase order checks (ProcessPurchaseOrder) on every create and update. A notification made by those checks is sent right after the save. If the save fails, Turbo deletes that notification so it is not sent later for a PO that was not saved. A failure to send is only logged and does not undo the save.

## Job, profile, currency, client and PO approver saves are checked

Lines 412-498

Asked as: job save validation; profile save checks; client notes on create; admin profile permission error; PO approver save rules

Turbo runs its own checks on every create and update of jobs, profiles, currencies, clients and PO approver properties. Client notes are checked on create. Admin profiles are checked on create and update, and creating one also first enforces admin profile request permissions.

## Expense saves run in one transaction

Lines 500-531

Asked as: expense save failed attachments; expense upload deleted on error; receipt removed after failed save; expense transaction rollback; expense create update fails

Turbo checks and saves an expense (create or update) in one transaction. If it fails, Turbo cleans up the uploaded expense documents.

## Expenses editing flag blocks several actions

Lines 532-544

Asked as: cannot delete expense; cannot delete purchase order; vendor cannot be added; vendor edit blocked; expenses editing off vendors

When the expenses editing setting is off, Turbo refuses to delete expenses and purchase orders. It also refuses to create, update or delete vendors.

## Rate sheet create checks the effective date first

Lines 546-561

Asked as: rate sheet effective date error; rate sheet create date check; rate sheet date validation; new rate sheet rejected; rate sheet update checks

When a rate sheet is created, Turbo validates its effective date before the other rate sheet checks. Updates run only the rate sheet checks.

## Changes to categories, clients and contacts mark jobs not imported

Lines 564-600

Asked as: job marked not imported; jobs need reimport after client change; category change job import; contact update jobs unimported; client edit affects jobs

Creating or deleting a category marks its job as not imported. Updating a client or client contact marks the jobs that use it as not imported, but only when something meaningful changed (ignoring _imported). For a client, this covers jobs where it is the client or the job owner.

## Notification records are written only by Turbo's own code

Lines 602-616

Asked as: cannot create notification via API; edit notifications directly; notification write refused; notifications read only; WriteStatusUpdated error

The notifications collection is not exposed to the API, so staff cannot create or update a notification through a request. Every time Turbo's own code creates or updates a notification, it calls notifications.WriteStatusUpdated first. If that call returns an error, the save fails. This file does not show what WriteStatusUpdated sets.

## Only the approver can share a time sheet with another manager

Lines 618-627

Asked as: share time sheet with manager; other manager cannot approve; who can share timesheet; shared with me approve; reviewer cannot reject timesheet

Only the approver of a time sheet can share it, by adding a reviewer (a time_sheet_reviewers row). The other manager finds it under Shared with Me and can read it. The other manager cannot approve or reject it. The hooks in this file do not change who can share, and the sharing rules have no condition on the time sheet's state.

## Time sheet reviewer must be an active user

Lines 629-650

Asked as: reviewer not active error; share timesheet inactive user; can't share with user; reviewer must be active; sharing blocked inactive manager

When a time sheet is shared with a reviewer, Turbo refuses if the reviewer is not active. The message is 'reviewer must be an active user'.

## Vendor absorb commit needs expenses editing enabled

Lines 657-662

Asked as: vendor absorb commit blocked; cannot commit vendor absorb; absorb commit expenses off; vendor absorb refused; vendor merge not allowed

Turbo refuses to commit a vendor absorb when the expenses editing setting is off.

## Absorb commit blocked while a parent absorb exists

Lines 664-683

Asked as: cannot commit absorb parent; absorb blocked by other absorb; commit or undo first; parent absorb exists error; absorb order conflict

Turbo refuses to commit an absorb while an absorb exists for a parent collection that blocks it. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'.

## Absorb commit needs the collection's claims

Lines 687-705

Asked as: no permission to commit absorb; absorb commit access denied; missing claims for absorb; can't commit absorb permission; absorb claims required

For a collection that defines its own claim sets, the person must hold the required claims to commit an absorb. Otherwise Turbo refuses with 'you do not have permission to commit a X absorb'.

## Absorb commit settles proposals

Lines 707-746

Asked as: absorb accepts proposals; proposal closed after commit; absorb settles proposals; proposal outcome notification; commit absorb proposal update

Committing an absorb for a collection whose absorb can be proposed accepts the proposal the action carries out and closes the proposals that name an absorbed record. The commit steps, proposal updates and delete run in one transaction. Turbo then queues outcome notifications for the settled proposals. A failed notification is only logged.

## Sharing a time sheet notifies the reviewer

Lines 753-784

Asked as: reviewer notified when shared; shared timesheet email; notify manager share sheet; share notification failed; who shared the timesheet

After a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to that reviewer, naming the time sheet's approver as the sharer. A failure to notify is only logged and does not block the share.
