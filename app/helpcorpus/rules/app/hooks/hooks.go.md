<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"7ba6ffdfff69b9b3fe8acec47616da2660a22e287efe52547cf124f3fa5325ec","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Purchase order notification sent after save

Lines 27-72

Asked as: po email not sent; purchase order alert missing; notification not delivered after saving po; po saved but no alert

When saving a purchase order creates a notification, Turbo sends it immediately after the save. If the save fails, Turbo deletes the pending notification so the cron does not deliver it. A failure to send is only logged and does not undo the save.

## Expense, PO and vendor deletes need expenses editing enabled

Lines 81-92

Asked as: can't delete expense; delete button greyed out; cannot delete purchase order; vendor can't be edited; expenses editing disabled

When expenses editing is disabled, Turbo refuses to delete expenses and purchase orders. It also refuses to create, update or delete vendors.

## Microsoft login: new user details

Lines 127-224

Asked as: new user email blank; name missing after first microsoft sign in; username empty new account; first microsoft login profile incomplete

On a first Microsoft login (a new user), Turbo fills in an empty email from the Microsoft data, with a fallback (data.fallbackEmail). It fills in an empty name from the Microsoft full name. It builds an empty username as a unique name (buildUniqueMicrosoftUsername). Only Microsoft logins get this treatment.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: microsoft login created duplicate account; new account on microsoft sign in; relink microsoft account; failed to look up existing microsoft user

When PocketBase finds no user by Microsoft provider ID, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the external auth and logs that user in instead of creating a new user. Failures give 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## Microsoft login syncs identity and onboards

Lines 226-280

Asked as: name not updated from microsoft; microsoft sign in onboarding failed; no initial time claim; failed to sync microsoft identity

After each successful Microsoft login, Turbo syncs the user's name and email from Microsoft, and the username stays as it was. A brand-new user also gets an initial time claim. Every Microsoft login ensures the user has an admin profile, and Turbo creates no profile record because the manager is collected later in the app. All of this runs in one transaction. A failure gives 'failed to sync Microsoft identity', 'failed to assign initial time claim' or 'failed to complete Microsoft user onboarding'.

## Inactive accounts cannot log in

Lines 304-316

Asked as: can't log in deactivated account; disabled user login blocked; this account is inactive; inactive user can't sign in

When a user logs in, Turbo checks the admin profile (admin_profiles). If its active flag is off, login fails with 'this account is inactive' and no token is issued. The check runs on every login, even when the credentials are valid.

## Login without an admin profile

Lines 317-330

Asked as: account setup incomplete; can't log in no profile; missing admin profile on login; oauth login profile created

If a user has no admin profile and logs in by a method other than OAuth2, login fails with 'account setup incomplete'. If the user logged in by OAuth2, Turbo creates the missing admin profile and lets the login continue. If creating it fails, the message is 'failed to complete user onboarding'.

## Time entry saves run through ProcessTimeEntry

Lines 333-344

Asked as: time entry won't save; timesheet entry rejected; why can't save hours; time entry error on save

Before Turbo creates or updates a time entry, it runs the time entry checks in ProcessTimeEntry. A refused check returns its error to the person and the entry is not saved. This file does not show what ProcessTimeEntry checks.

## Time entry and time amendment deletes need time editing enabled

Lines 345-379

Asked as: can't delete time entry; delete time amendment greyed out; time editing disabled delete; removing hours not allowed

Deleting a time entry or a time amendment is refused when time editing is disabled (ErrTimeEditingDisabled). Creating and updating them go through their own validation (ProcessTimeEntry, ProcessTimeAmendment).

## Purchase order saves run checks and send notifications

Lines 381-404

Asked as: purchase order won't save; po save error; po checks failing; po notification sent on save

Before Turbo creates or updates a purchase order (PO), it runs ProcessPurchaseOrder, which may queue a notification. If the save then fails, Turbo deletes that pending notification. If the save succeeds, Turbo sends the notification at once.

## Job saves run through ProcessJob

Lines 409-420

Asked as: job won't save; job creation error; job checks refused; can't update job

Before Turbo creates or updates a job, it runs the job checks in ProcessJob. A refused check returns its error and the job is not saved. This file does not show what ProcessJob checks.

## Profile saves run through ProcessProfile

Lines 422-433

Asked as: profile won't save; profile update rejected; cannot edit profile

Before Turbo creates or updates a profile, it runs the checks in ProcessProfile. A refused check returns its error and the profile is not saved.

## Currency saves run through ProcessCurrency

Lines 435-446

Asked as: currency not saving; can't add currency; currency update error

Before Turbo creates or updates a currency, it runs the checks in ProcessCurrency. A refused check returns its error and the currency is not saved.

## Admin profile creation permission check

Lines 447-462

Asked as: can't create admin profile; admin profile permission denied; create admin profile forbidden; editing admin profile allowed

Creating an admin profile first runs a request permission check (EnforceAdminProfileRequestPermissions). Updating an admin profile skips that check.

## Client saves run through ProcessClient

Lines 464-475

Asked as: client won't save; cannot update client; client error on save

Before Turbo creates or updates a client, it runs the checks in ProcessClient. A refused check returns its error and the client is not saved.

## Client notes are checked on create only

Lines 477-482

Asked as: client note edit not checked; editing client note no validation; can't save client note edit

Turbo runs ProcessClientNote when a client note is created. This file adds no check when a client note is updated.

## PO approver properties saves run through ProcessPOApproverProps

Lines 484-495

Asked as: approver settings not saving; po approver properties error; approver props save failed

Before Turbo creates or updates a record of a PO approver's properties (po_approver_props), it runs the checks in ProcessPOApproverProps. A refused check returns its error and the record is not saved.

## Expense saves run in a transaction

Lines 497-528

Asked as: expense save failed receipts; attachments gone after expense failure; uploaded documents removed; expense not saved rollback

Creating or updating an expense runs its validation and save in one transaction. If it fails, Turbo cleans up the uploaded expense documents.

## New rate sheet checks effective date first

Lines 543-557

Asked as: rate sheet effective date error; new rate sheet date rejected; can't create rate sheet date

Creating a rate sheet validates its effective date (ValidateRateSheetEffectiveDate) before the general rate sheet checks. Updating a rate sheet skips the effective date check.

## Changes to categories, clients and contacts mark jobs not imported

Lines 559-596

Asked as: job marked not imported; why job needs reimport; category change job status; client edit changed job status

Creating or deleting a category marks its job as not imported. Updating a client or a client contact marks the jobs that use it as not imported, but only when something other than the _imported field changed. A client also marks jobs that it owns (job_owner).

## Notification writes update status through WriteStatusUpdated

Lines 598-612

Asked as: notification status update; notification write error; notification status field

When Turbo creates or updates a notification record, it first runs WriteStatusUpdated. Notifications are not exposed to the API, so only Turbo's own code writes them.

## Only the approver can share a time sheet

Lines 614-622

Asked as: can't share time sheet; share timesheet button missing; other manager can't approve; reviewer not allowed to share

Only the approver of a time sheet can share it with another manager, by adding or removing a reviewer row. The create and delete rules have no state condition. The hooks in this file do not change who can share. The other manager can read the shared time sheet but cannot approve or reject it.

## Time sheet reviewer must be active

Lines 625-646

Asked as: reviewer not active; selected reviewer is not an active user; can't share timesheet inactive manager; share with inactive user

When an approver shares a time sheet, the reviewer must be an active user. Otherwise Turbo refuses with 'reviewer must be an active user' (the reviewer field says 'the selected reviewer is not an active user').

## Vendor absorb commit needs expenses editing enabled

Lines 653-658

Asked as: vendor absorb won't commit; commit vendors absorb refused; absorb vendor expenses editing

Committing a vendors absorb is refused when expenses editing is disabled.

## Absorb commit is blocked by a parent absorb

Lines 660-679

Asked as: absorb commit blocked; cannot commit absorb; commit or undo absorb first; absorb waiting on another

Committing an absorb is refused while an absorb exists for a parent collection that blocks it. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'.

## Absorb commit needs claims for some collections

Lines 683-701

Asked as: no permission to commit absorb; absorb commit access denied; missing claim commit absorb

If an absorb collection has its own claim sets, the committer must hold the required claims. Otherwise Turbo refuses with 'you do not have permission to commit a X absorb'.

## Absorb commit runs its step and delete together

Lines 703-715

Asked as: absorb commit partly failed; absorb delete failed leftover; absorb commit rollback

If a collection has a commit step (OnCommit), Turbo runs it and deletes the absorb action in one transaction. A failed delete leaves nothing behind. Without a commit step, the delete just proceeds.

## New reviewer gets a notification

Lines 723-754

Asked as: reviewer got no notification; shared timesheet email not sent; share worked but no alert; notify reviewer timesheet

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to the reviewer. The sharer is the time sheet's approver. If the notification fails, the share still succeeds.
