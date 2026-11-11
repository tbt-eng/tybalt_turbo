<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"8fadc0e278c542378acdabcd91e459be52c6a368815cc7f258645107a8bb4dc2","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Purchase order notification sent after a successful save

Lines 27-42

Asked as: purchase order email; po alert not sent; notify when po saved; po notice after save

When saving a purchase order creates a notification, Turbo sends it right after the save succeeds. A send failure is only logged and does not undo the save.

## Failed PO save deletes its pending notification

Lines 44-72

Asked as: po failed but email sent; stop notice for unsaved po; cancel po notification on error; po email for failed save

If a purchase order save fails after the notification was created, Turbo deletes that notification. This stops the cron from delivering a notice for a PO change that never saved.

## Expense, PO and vendor changes blocked when expenses editing is off

Lines 81-92

Asked as: can't delete expense; vendor edits locked; editing turned off error; can't change vendors now

While expenses editing is disabled, Turbo refuses to delete expenses and purchase orders. It also refuses to create, update or delete vendors. The error is ErrExpensesEditingDisabled.

## Only Microsoft logins get custom onboarding

Lines 126-129

Asked as: google login onboarding; okta sign in setup; non microsoft sign in; other provider login setup

Turbo applies this onboarding only to Microsoft OAuth logins. Other OAuth providers go through the normal PocketBase flow.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: microsoft login duplicate account; existing user microsoft sign in; account already exists microsoft; microsoft account linked wrong user

When PocketBase finds no user by provider ID, Turbo looks for an existing user to relink. If it finds one, Turbo relinks the Microsoft account to that user and treats the login as an existing user, not a new one. Failures give "failed to look up existing Microsoft user" or "failed to relink Microsoft account".

## Microsoft first login builds name, email and username

Lines 179-224

Asked as: blank email after microsoft login; username missing new microsoft user; empty name microsoft signup; microsoft first login details

For a new Microsoft user, Turbo fills in an empty email from the Microsoft mail, or from userPrincipalName when mail is blank. It fills an empty name from the Microsoft full name. It builds an empty username from the email local part instead of letting PocketBase generate one. Values already supplied are kept.

## Microsoft login runs in one transaction with post-login steps

Lines 226-280

Asked as: login half set up; microsoft login partially failed; new user incomplete after login; login rolled back microsoft

After a Microsoft login Turbo syncs the user's name and email from Microsoft and makes sure an admin profile exists. It assigns the initial time claim to new users only. If any step fails, the whole login is rolled back, so no half-onboarded user remains. Turbo does not create a profile record here, because the manager must be chosen later in the app.

## Inactive accounts cannot log in

Lines 304-316

Asked as: account disabled can't sign in; deactivated user login; login blocked inactive; inactive user cannot log in

When a user's admin profile is not active (active is false), Turbo refuses the login with "this account is inactive". This happens even when the credentials are valid. No token is issued.

## Login without an admin profile

Lines 317-323

Asked as: account setup incomplete; can't log in no profile; password login missing profile; no admin profile login

If a user has no admin profile and does not sign in with OAuth2, Turbo refuses the login with "account setup incomplete". Turbo does not create the admin profile for password logins.

## OAuth2 login repairs a missing admin profile

Lines 324-329

Asked as: oauth login no admin profile; sso user profile created; onboarding failed on sign in; sign in profile missing

If a user who signs in with OAuth2 has no admin profile, Turbo creates one and then lets the login continue. If creating it fails, Turbo refuses with "failed to complete user onboarding".

## Time entries run a check on every create and update

Lines 333-344

Asked as: time entry won't save; timesheet entry validation error; hours entry rejected; time entry save refused

Turbo runs its time entry processing (ProcessTimeEntry) before it saves any new or changed time entry. If that processing returns an error, Turbo refuses the save and returns the error. The rules inside ProcessTimeEntry are in another file.

## Time entry and amendment delete blocked when time editing is off

Lines 345-379

Asked as: can't delete timesheet entry; time editing locked; delete amendment not allowed; time entry delete blocked

Turbo refuses to delete a time entry or a time amendment while time editing is disabled, with the error ErrTimeEditingDisabled. If Turbo cannot read the setting, it also refuses.

## Time amendments run a check on every create and update

Lines 357-368

Asked as: amendment won't save; time amendment rejected; amendment validation error; time change request not saving

Turbo runs its time amendment processing (ProcessTimeAmendment) before it saves any new or changed time amendment. If that processing returns an error, Turbo refuses the save and returns the error. The rules inside ProcessTimeAmendment are in another file.

## Purchase order create and update run PO checks

Lines 381-404

Asked as: purchase order won't save; po validation error; po save refused; po update rejected

Turbo runs its purchase order processing (ProcessPurchaseOrder) before it saves a new or changed purchase order. If the processing fails, Turbo refuses the save. The rules inside ProcessPurchaseOrder are in another file.

## Job create and update run job checks

Lines 405-419

Asked as: job won't save; job validation error; can't save job; job update rejected

Turbo runs its job processing (ProcessJob) before it saves a new or changed job. If it fails, Turbo refuses the save. Project workflow hooks and client setup hooks are also registered here. Their rules are in other files.

## Profile create and update run profile checks

Lines 420-432

Asked as: profile won't save; profile update rejected; profile validation error; can't save profile

Turbo runs its profile processing (ProcessProfile) before it saves a new or changed profile. If it fails, Turbo refuses the save. The rules inside ProcessProfile are in another file.

## Currency create and update run currency checks

Lines 433-445

Asked as: currency won't save; currency validation failed; currency update rejected; can't save currency

Turbo runs its currency processing (ProcessCurrency) before it saves a new or changed currency. If it fails, Turbo refuses the save. The rules inside ProcessCurrency are in another file.

## Admin profile creation checks request permissions first

Lines 446-461

Asked as: can't create admin profile; admin profile permission denied; new admin profile forbidden; create admin profile error

Before the usual admin profile processing, Turbo checks the request's permissions when an admin profile is created. Updates skip this check.

## Admin profile update runs processing but no permission check

Lines 456-461

Asked as: edit admin profile permissions; admin profile update allowed; change admin profile; update admin profile error

On update of an admin profile, Turbo runs ProcessAdminProfile. The extra request permission check (EnforceAdminProfileRequestPermissions) runs only on create, not on update in this file.

## Client create and update run client checks

Lines 462-474

Asked as: client won't save; client validation error; can't save client; client update rejected

Turbo runs its client processing (ProcessClient) before it saves a new or changed client. If it fails, Turbo refuses the save. The rules inside ProcessClient are in another file.

## Client notes can only be created

Lines 475-481

Asked as: can't edit client note; client note update; edit note not possible; change saved client note

Turbo runs client note processing on create only. The hooks set no update step for client notes.

## PO approver properties create and update run checks

Lines 482-494

Asked as: approver settings won't save; po approver error; approver properties rejected; approver setup not saving

Turbo runs its PO approver properties processing (ProcessPOApproverProps) before it saves new or changed PO approver properties. If it fails, Turbo refuses the save. The rules inside ProcessPOApproverProps are in another file.

## Expense save is atomic and cleans up uploads

Lines 495-527

Asked as: expense receipt left behind; failed expense upload files; expense save failed attachments; leftover expense documents

Turbo saves an expense and its processing in one transaction. If the save fails, Turbo cleans up the expense's uploaded documents.

## Expense and PO delete and vendor gate hooks

Lines 528-540

Asked as: can't delete po; vendor changes blocked; delete expense locked; vendor delete not allowed

The expenses editing check runs as a gate on expense delete, PO delete, and vendor create, update and delete.

## New rate sheet checks its effective date

Lines 541-556

Asked as: rate sheet date error; effective date rejected; new rate sheet invalid date; rate sheet start date

Creating a rate sheet validates its effective date (ValidateRateSheetEffectiveDate) before the other rate sheet checks. Updating a rate sheet does not run this check.

## Related changes mark the job as not imported

Lines 558-595

Asked as: job not imported after change; client update job import; category delete job reimport; job marked not imported

Creating or deleting a category marks its job as not imported. Updating a client or a client contact does the same for the jobs that reference it, but only when something meaningful besides _imported changed. A client update covers jobs that use it as client or as job owner, and a contact update covers jobs that use it as contact.

## Notification records run a status update on every create and update

Lines 597-611

Asked as: notification save refused; notification record error; can't save notification; notification status error

Notifications are not exposed to the API. Each time Turbo creates or updates a notification record, it first runs notifications.WriteStatusUpdated. If that returns an error, the save is refused.

## Only the approver can share a time sheet with another manager

Lines 613-622

Asked as: share timesheet not allowed; who can share timesheet; give manager timesheet access; can't share timesheet

A time sheet is shared through a time sheet reviewer row (time_sheet_reviewers). Only the approver of the time sheet can add or remove such a row, and the rules have no state condition. The hooks in this file do not change who can share.

## A shared reviewer can read a time sheet but not approve it

Lines 613-618

Asked as: shared timesheet can't approve; shared with me read only; manager approve shared timesheet; shared timesheet view only

A manager who gets a time sheet shared finds it under Shared with Me and can read it. That manager cannot approve or reject it. The person's own manager still approves it.

## Reviewer for a shared time sheet must be active

Lines 623-645

Asked as: inactive reviewer error; reviewer not active; share with disabled user; reviewer account inactive

When an approver shares a time sheet, the chosen reviewer must be an active user. Otherwise Turbo refuses with "reviewer must be an active user" (the selected reviewer is not an active user).

## Adding a reviewer to a time sheet runs on create only

Lines 624-643

Asked as: add reviewer fails; reviewer status check failed; reviewer must be active error; adding reviewer error

Turbo checks that the reviewer is active only when a reviewer row is created through the API. The hook does not check on update. If the check for the active status itself fails, Turbo shows 'failed to check reviewer active status'. If the reviewer is inactive, Turbo shows 'reviewer must be an active user' with the field message 'the selected reviewer is not an active user'.

## Absorb commit blocked by a parent absorb

Lines 647-681

Asked as: can't commit absorb; absorb blocked by parent; commit absorb error; parent absorb exists

A person cannot commit an absorb action while an absorb for a parent collection exists. Turbo refuses with "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first".

## Vendor absorb commit blocked when expenses editing is off

Lines 652-657

Asked as: vendor absorb commit blocked; can't commit vendor absorb; vendors absorb expenses editing; vendor absorb not allowed

A person cannot commit a vendors absorb action while expenses editing is disabled.

## Absorb commit blocker lookup failure

Lines 661-670

Asked as: failed to check absorb; absorb check failed; can't verify absorb commit; absorb lookup error

When Turbo cannot check whether a parent absorb exists, it refuses the commit with 'failed to check <parent collection> absorb action'.

## Sharing a time sheet notifies the new reviewer

Lines 683-719

Asked as: reviewer not notified; shared timesheet email; notify manager timesheet shared; no email when shared

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification for that reviewer, naming the time sheet's approver as the sharer. If the notification fails, the share still succeeds.

## Reviewer notification failure does not block sharing

Lines 688-719

Asked as: share failed notification error; sharing despite email error; notification error timesheet share; share works but no alert

If Turbo cannot load the time sheet or queue the notification when a reviewer is added, it logs the error and the sharing still goes through. The notification is sent in the name of the time sheet's approver (approver).
