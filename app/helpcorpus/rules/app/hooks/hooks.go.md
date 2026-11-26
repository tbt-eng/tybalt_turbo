<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"0003d276d7492d6ba05a851f7e535cd62478da5617bbd37217adde9c528be1a0","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Failed notification send does not undo a saved PO

Lines 27-42

Asked as: PO saved but notification failed; notification send error after PO save; PO still saved when email fails; notification not sent after purchase order

After a purchase order saves, Turbo sends the notification that was created for it straight away. If the send fails, Turbo logs the error and keeps the saved PO. The notification is not sent immediately in that case. When no notification was created, Turbo sends nothing.

## Orphaned PO notification is deleted when the PO save fails

Lines 44-72

Asked as: PO failed to save but notice sent; delete pending notification; cron sends notice for unsaved PO; notification for PO that did not save

If a purchase order save fails after Turbo created a notification for it, Turbo deletes that pending notification (notifications record). This stops the cron from delivering a notice for a PO write that never saved. If the delete fails, Turbo logs the error and the cron may still deliver the notice.

## Expenses editing switch blocks expense-related changes when off

Lines 81-92

Asked as: expenses editing turned off; cannot change expense feature disabled; expenses locked error; feature flag expense edits blocked

Turbo checks the expenses editing feature flag before the changes that are gated on it. When the flag is off, Turbo refuses with the error ErrExpensesEditingDisabled. When the flag is on, the change continues. If Turbo cannot read the flag, it returns that error instead.

## Returning Microsoft sign-ins sync name and email only

Lines 95-122

Asked as: Microsoft login updates my name; email not updating after Microsoft sign in; sync name from Microsoft; changed my name in Microsoft

On each later Microsoft sign-in, Turbo updates the user's name when Microsoft gives a usable full name. It updates the email only when no other user already has that email. It does not rename the username once created. It creates an admin profile only if one is missing and never refreshes it from Microsoft.

## Microsoft sign-in marks user verified when email matches

Lines 98-105

Asked as: Microsoft email blank; userPrincipalName used as email; account not verified after Microsoft login; email verified Microsoft sign in

Microsoft often leaves the mail field blank, so Turbo falls back to userPrincipalName for the email. After the user is saved, Turbo stores that email and sets verified only if the stored email matches the Microsoft identity it just accepted.

## Only Microsoft sign-ins get Microsoft onboarding handling

Lines 126-129

Asked as: Google login onboarding; other OAuth provider sign in; non Microsoft social login; SSO provider different from Microsoft

Turbo applies its custom onboarding only when the OAuth provider is Microsoft (compared without regard to case). Sign-ins through any other OAuth provider skip this handling and follow PocketBase's default flow.

## Microsoft sign-in is all-or-nothing

Lines 135-280

Asked as: half onboarded user; Microsoft sign in partly failed; rollback Microsoft login; account left incomplete after login error

Turbo runs the whole Microsoft sign-in in one transaction: relinking, user creation, identity sync, initial time claim and admin profile creation. If any step fails, none of it is saved, so a user is never left partly onboarded.

## Microsoft sign-in relinks an existing user

Lines 150-177

Asked as: existing user signs in with Microsoft; duplicate account Microsoft login; link Microsoft to my old account; Microsoft account already had user

When PocketBase finds no user for the Microsoft account, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, Turbo relinks the external auth and signs that user in instead of creating a new user. If the lookup or the relink fails, Turbo refuses with 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## Failed Microsoft onboarding steps show specific errors

Lines 160-276

Asked as: Microsoft onboarding error message; Microsoft sign in failed step; failed to sync Microsoft identity; sign in error Microsoft onboarding

A failed step in Microsoft sign-in is refused with one of these messages: "failed to look up existing Microsoft user", "failed to relink Microsoft account", "failed to prepare Microsoft username", "failed to sync Microsoft identity", "failed to assign initial time claim" or "failed to complete Microsoft user onboarding".

## Microsoft first sign-in builds email, name and username

Lines 179-224

Asked as: new user blank email from Microsoft; username generated on first login; Microsoft first login name missing; empty name after Microsoft sign up

When a new user signs in with Microsoft, Turbo fills a blank email from Microsoft's mail, or from userPrincipalName when mail is blank. It fills a blank name from the full name Microsoft gives. It fills a blank username with a unique one that it builds from the Microsoft data. Existing values in the create data are kept.

## Microsoft sign-in syncs identity and creates admin profile

Lines 226-280

Asked as: Microsoft login admin profile created; identity sync after Microsoft login; profile missing after Microsoft login; Microsoft login updates profile

After each Microsoft sign-in, Turbo syncs the user's identity from Microsoft (syncMicrosoftUserIdentity) and ensures an admin profile exists. It creates no profile (profiles record), because the manager must be chosen later in the app. A new user also gets an initial time claim. Turbo runs all of this in one transaction, so a failure leaves no half-onboarded user. Only Microsoft logins get this handling.

## New Microsoft user gets an initial time claim

Lines 226-267

Asked as: new user no time claim; initial time claim missing; first time claim not assigned; claim not given to new Microsoft user

When a Microsoft sign-in creates a new user, Turbo assigns that user an initial time claim after saving the user (ensureInitialTimeClaim). A returning user, or one relinked to an existing account, does not get it.

## Microsoft sign-in does not create a profile

Lines 269-276

Asked as: no profile after Microsoft sign in; profile not created at login; manager field empty after Microsoft login; have to fill manager later

A Microsoft sign-in makes sure the user has an admin profile but does not create a profile record. The user must enter business fields such as manager in the app later.

## Inactive account is refused at sign-in

Lines 304-316

Asked as: cannot log in account disabled; deactivated user cannot sign in; account inactive login error; locked out inactive account

When the user's admin profile is not active, Turbo refuses the sign-in with "this account is inactive", even if the credentials are valid. This check runs on every sign-in method.

## Inactive accounts cannot sign in

Lines 308-316

Asked as: inactive user login blocked; active flag off sign in refused; disabled account cannot sign in; sign in refused after deactivation

After credentials check out, Turbo refuses sign-in with 'this account is inactive' when the admin profile's active flag is off. This applies to every sign-in method.

## Missing admin profile at sign-in

Lines 317-329

Asked as: account setup incomplete; no admin profile at login; OAuth2 creates profile at sign in; user onboarding failed login

If a user has no admin profile and signs in with a method other than OAuth2, Turbo refuses with 'account setup incomplete'. If the user signed in with OAuth2, Turbo creates the admin profile and lets the sign-in continue. If that fails, it refuses with 'failed to complete user onboarding'.

## Time entry saves are checked before saving

Lines 333-344

Asked as: time entry not saving; time entry validation failed; cannot save hours; timesheet entry rejected on save

Creating or updating a time entry first runs Turbo's time entry checks (ProcessTimeEntry). A failed check refuses the save and returns its error.

## Time entry and time amendment deletes need time editing on

Lines 345-379

Asked as: cannot delete time entry; delete amendment blocked; time editing turned off delete; cannot remove time amendment

Turbo refuses to delete a time entry or a time amendment when time editing is disabled. It gives the time-editing-disabled error (ErrTimeEditingDisabled).

## Time amendment and PO saves run Turbo's checks first

Lines 357-498

Asked as: PO save fails checks; time amendment validation error; job save refused; client save blocked by checks

Creating or updating a time amendment, a purchase order, a job, a profile, a currency, a client, a client note or a PO approver setting first runs that record's own checks. A failed check refuses the save. A purchase order that fails after its checks also has its pending notification deleted.

## PO notification sent after save, removed on failure

Lines 381-404

Asked as: PO notification not sent after save; notification removed when PO fails; cron not sending PO notice; send PO notice after saving

When a purchase order save creates a notification, Turbo sends it right after the save succeeds. If the save fails, Turbo deletes the pending notification so the cron does not send it for a PO that was never saved. A failed send does not undo the save.

## Only the first admin profile create checks request permissions

Lines 451-465

Asked as: cannot create admin profile via API; admin profile permission denied; edit admin profile permission check; API profile request forbidden

Creating an admin profile through the API first runs a separate permission check on the request (EnforceAdminProfileRequestPermissions). Updating an admin profile skips that check and runs only the admin profile checks.

## Client notes can be created but not edited here

Lines 480-485

Asked as: cannot edit client note; client note update blocked; change existing client note; client note changes not saved

Turbo's client note checks run only when a client note is created. This file adds no check for updating one.

## Expense saves run in a transaction

Lines 500-531

Asked as: expense upload left behind; expense save failed files; expense attachments not cleaned up; expense save rolled back

Turbo checks and saves an expense inside one transaction. When the save fails, it cleans up the uploaded expense documents.

## Expense, PO and vendor changes need expenses editing on

Lines 532-544

Asked as: cannot delete expense; vendor edits blocked; cannot delete purchase order; expenses editing off delete

When the expenses editing feature is off, Turbo refuses to delete expenses and purchase orders. It also refuses to create, update or delete vendors. It gives the expenses-editing-disabled error (ErrExpensesEditingDisabled).

## Rate sheet create checks effective date first

Lines 546-561

Asked as: rate sheet date error; effective date rejected; rate sheet create fails; rate sheet effective date wrong

Creating a rate sheet checks the effective date first, then runs the rate sheet checks. Updating skips the effective date check.

## Changes to categories, clients and contacts mark jobs not imported

Lines 563-600

Asked as: job marked not imported; job needs reimport after client change; category change resets job import; contact edit unimported job

Creating or deleting a category marks its job as not imported. Updating a client marks jobs that use it as client or job owner as not imported. Updating a client contact marks jobs that use it as contact as not imported. For clients and contacts this happens only when something other than the _imported field changes.

## Notification records get a status-updated stamp on save

Lines 602-616

Asked as: notification save refused; notification status stamp error; cannot save notification record; notification update fails

Whenever Turbo creates or updates a notification record, it first calls notifications.WriteStatusUpdated. If that call fails, the save is refused. The notifications collection is not exposed to the API, so only Turbo's own code triggers this.

## Only a time sheet's approver can share it with another manager

Lines 618-627

Asked as: share time sheet with manager; who can add reviewer; remove shared reviewer; only approver shares timesheet

A time sheet is shared with another manager through a time_sheet_reviewers row. Only the time sheet's approver can add or remove such a row. The rule has no condition on the time sheet's state. The hooks in this file do not change who can share.

## Shared reviewer can read a time sheet but not approve it

Lines 618-623

Asked as: shared timesheet cannot approve; Shared with Me approve button; reviewer read only timesheet; colleague sees my timesheet

A manager who receives a shared time sheet finds it under Shared with Me and can read it. The manager cannot approve or reject it. The person's own manager still approves it.

## Time sheet reviewer must be an active user

Lines 629-650

Asked as: reviewer not active user; cannot share with inactive manager; share timesheet with disabled user; reviewer selection error

When an approver shares a time sheet, Turbo refuses a reviewer who is not an active user with 'reviewer must be an active user'. The field message is 'the selected reviewer is not an active user'.

## Vendor absorb commit needs expenses editing on

Lines 657-662

Asked as: vendor absorb commit blocked; cannot commit vendor absorb; absorb commit expenses off; vendor absorb refused

Turbo refuses to commit a vendor absorb when expenses editing is disabled.

## Absorb commit blocked while a parent absorb exists

Lines 664-683

Asked as: cannot commit absorb; parent absorb blocking commit; absorb commit error message; commit or undo absorb first

Turbo refuses to commit an absorb while an absorb for a parent collection exists. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'. The parent collections come from absorb.GetCommitBlockers.

## Absorb commit needs the collection's claims

Lines 685-705

Asked as: no permission to commit absorb; absorb commit access denied; cannot commit rate sheet absorb; absorb claims missing

For a collection that has its own claim sets, the committer must hold the required claims. Otherwise Turbo refuses with 'you do not have permission to commit a X absorb'. The collection rule lets an absorb holder do any action and lets an admin act on rate sheet absorbs only.

## Absorb commit saves history and settles proposals

Lines 707-749

Asked as: absorb history after commit; absorb proposals settled; absorb commit notifications sent; commit absorb steps

Committing an absorb copies it to absorb history with who absorbed and who committed. It runs the collection's commit steps, then deletes the action, all in one transaction. For collections whose absorbs can be proposed, it also settles the absorb proposals. Turbo then queues an outcome notification for each settled proposal.

## New time sheet reviewer gets a notification

Lines 757-788

Asked as: notified when added as reviewer; shared timesheet notification; who shared my timesheet; reviewer notice failed share

When a reviewer is added to a time sheet, Turbo queues a shared notification for that reviewer. It names the time sheet's approver as the sharer. If the notification fails, the share still goes ahead.
