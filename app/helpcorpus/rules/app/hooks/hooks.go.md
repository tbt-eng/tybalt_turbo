<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"429743d7adb672c337e94740a214877b7ba1c7a2c2e7ae644fb12fcd0d22ffbe","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors return their message to the client

Lines 19-25

Asked as: error message hidden; vague error on save; see why it failed; error details not shown

When a hook fails with a Turbo hook error, Turbo sends that error's message and details back to the person who made the request.

## Purchase order notification is sent only after the save

Lines 27-72

Asked as: po email sent before save; notification sent though save failed; purchase order alert timing; email sent but save failed

When a purchase order is created or updated, a notification made while processing is sent after the save succeeds. If the save fails, Turbo deletes the pending notification so the cron never sends it. A failed send does not undo the save.

## Failed purchase order save deletes its pending notification

Lines 44-72

Asked as: pending po email after failed save; po notification still queued; cron sends po email for failed save; stuck purchase order notification

If a purchase order save fails after the pre-save checks, Turbo deletes the notification created for it (deleteOrphanedNotification). The notification does not stay pending for the cron to deliver. A failure to send after a successful save is only logged, and the save stays.

## Deleting expenses or purchase orders needs expenses editing enabled

Lines 81-92

Asked as: cannot delete expense; can't delete po; delete purchase order blocked; delete button does nothing expense

Turbo refuses to delete an expense or a purchase order when the expenses editing setting is off. It returns ErrExpensesEditingDisabled.

## Only Microsoft logins get custom onboarding

Lines 127-129

Asked as: non microsoft login onboarding; google login setup steps; why no onboarding for my login; other provider signup steps

The Microsoft onboarding steps run only when the OAuth provider is Microsoft. Other providers use PocketBase's default flow.

## Microsoft login runs as one all-or-nothing transaction

Lines 130-149

Asked as: half created user microsoft; partial signup after microsoft login; failed microsoft sign in leftovers; microsoft login rollback

Turbo runs the whole Microsoft sign-in in one transaction: relinking, creating the user, the identity sync, the initial time claim and the admin profile. If any step fails, the transaction rolls back. A failed sign-in leaves no partly onboarded user.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: existing user microsoft duplicate; account already exists microsoft sign in; link microsoft to old account; microsoft login creates second account

If PocketBase finds no user by provider id, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the Microsoft account and treats the login as an existing user, not a new one. A failed lookup gives "failed to look up existing Microsoft user". A failed relink gives "failed to relink Microsoft account".

## Microsoft first login builds name, email and username

Lines 179-224

Asked as: username like users123456; microsoft login missing email; blank name after microsoft sign in; auto generated username microsoft

For a first-time Microsoft user, Turbo fills any blank email from the Microsoft mail, falling back to userPrincipalName. It fills a blank name from the Microsoft full name. It builds a unique username from the data, instead of PocketBase's generated usersNNNNNN. Values already present in the create data are kept.

## New Microsoft users get the initial time claim

Lines 226-267

Asked as: new user no time claim; failed to assign time claim; first login time claim error; time claim missing new account

Only on a first-time Microsoft login, Turbo assigns the initial time claim to the new user. If this fails, Turbo refuses with "failed to assign initial time claim".

## Microsoft login syncs identity and onboards the user

Lines 255-280

Asked as: microsoft name not updated; profile missing after login; sync user details microsoft; microsoft details out of date

After each Microsoft login, Turbo syncs the user's name and email from Microsoft. Turbo makes sure the user has an admin profile, and creates no profile record because the manager must be chosen in the app. All steps run in one transaction, so a failure leaves no half-onboarded user.

## Inactive accounts cannot log in

Lines 304-316

Asked as: cannot log in inactive; account disabled login refused; this account is inactive; login blocked deactivated user

When a user logs in with valid credentials, Turbo checks the admin profile (admin_profiles) active flag. If the profile is not active, Turbo refuses the login with "this account is inactive" and issues no token.

## Non-OAuth login without an admin profile is refused

Lines 317-323

Asked as: account setup incomplete; password login refused new user; cannot sign in with password; no profile password login

If a user logs in with a password (not OAuth2) and has no admin profile, Turbo refuses with "account setup incomplete". Turbo does not create the admin profile for them.

## OAuth login with no admin profile repairs it

Lines 321-329

Asked as: admin profile missing oauth login; login fails no profile microsoft; failed to complete user onboarding; repair profile on login

If a user logs in with OAuth2 and has no admin profile, Turbo creates one (ensureUserAdminProfile) and completes the login. If the creation fails, Turbo refuses with "failed to complete user onboarding".

## Time entry saves run Turbo's time entry checks

Lines 332-344

Asked as: time entry save refused; cannot save time entry; time entry validation error; timesheet entry rejected

Before Turbo creates or updates a time entry, it runs its time entry processing (ProcessTimeEntry). If that processing fails, Turbo refuses the save and returns the error.

## Deleting time entries or time amendments needs time editing enabled

Lines 345-379

Asked as: cannot delete time entry; time amendment delete blocked; time editing off delete; can't remove time amendment

Turbo refuses to delete a time entry or a time amendment when the time editing setting is off. It returns ErrTimeEditingDisabled, and the delete runs only when time editing is enabled.

## Time amendment saves run Turbo's amendment checks

Lines 356-368

Asked as: amendment save error; cannot save time amendment; amendment validation failed; time amendment rejected

Before Turbo creates or updates a time amendment, it runs its time amendment processing (ProcessTimeAmendment). If that processing fails, Turbo refuses the save and returns the error.

## Purchase order saves run Turbo's purchase order checks

Lines 380-404

Asked as: purchase order save refused; po validation error; cannot save po; purchase order check failed

Before Turbo creates or updates a purchase order, it runs its purchase order processing (ProcessPurchaseOrder). If that processing fails, Turbo refuses the save and returns the error. The processing may queue a notification, and Turbo sends it only after the save succeeds.

## Job, profile, currency, client and client note saves are checked

Lines 406-434

Asked as: job save error; client save refused; profile save fails; cannot save job or profile

Turbo runs its own processing before it creates or updates a job, a profile, a currency or a client, and before it creates a client note. If the processing fails, Turbo refuses the save and returns the error. Client notes have no update hook here.

## Currency saves run Turbo's currency checks

Lines 435-447

Asked as: currency save error; cannot save currency; currency rejected; currency check failed

Turbo runs its currency checks (ProcessCurrency) before it creates or updates a currency record. If a check fails, the save stops and the client gets the error. This file does not show what the checks are.

## Admin profile creation checks request permissions first

Lines 448-463

Asked as: cannot create admin profile; admin profile permission denied; create profile not allowed; admin profile request rejected

Before Turbo creates an admin profile, it checks the request's permissions (EnforceAdminProfileRequestPermissions), then runs admin profile processing. An update runs only the admin profile processing, not the permission check.

## Client saves run Turbo's client checks

Lines 464-476

Asked as: client save refused; cannot save client; client validation error; client update failed

Turbo runs its client checks (ProcessClient) before it creates or updates a client record. If a check fails, the save stops and the client gets the error. This file does not show what the checks are.

## Client notes are checked on create only

Lines 477-483

Asked as: client note not saved; cannot add client note; client note edit; client note check failed

Turbo runs its client note checks (ProcessClientNote) when a client note is created. If a check fails, the note is not saved. This file sets no check on client note updates.

## PO approver properties saves are checked

Lines 484-496

Asked as: approver settings not saving; po approver properties error; approver props save fails; approver settings rejected

Turbo runs its own processing (ProcessPOApproverProps) before it creates or updates a PO approver properties record. If it fails, Turbo refuses the save.

## Failed expense save cleans up uploaded documents

Lines 497-529

Asked as: receipt uploads left behind; expense upload deleted after failure; attachments gone after expense fails; expense files cleaned up

When an expense create or update fails, Turbo cleans up the documents uploaded with that request (CleanupExpenseDocumentUploads).

## Expense saves run in one transaction

Lines 498-529

Asked as: expense half saved; expense save rollback; expense partially saved; expense save fails midway

Turbo processes an expense create or update and saves it in one transaction. If the save fails, Turbo cleans up the uploaded expense documents.

## Vendor changes need expenses editing enabled

Lines 530-542

Asked as: cannot add vendor; vendor edit blocked; vendor delete not allowed; vendor changes refused

Turbo refuses to create, update or delete a vendor when the expenses editing setting is off.

## New rate sheets must pass an effective date check

Lines 544-558

Asked as: rate sheet date error; effective date rejected; new rate sheet date check; rate sheet start date invalid

When a rate sheet is created, Turbo first checks its effective date (ValidateRateSheetEffectiveDate), then processes it. Updates skip the effective date check.

## Rate sheet update skips the effective date check

Lines 553-558

Asked as: rate sheet edit date; update rate sheet effective date; rate sheet update error; changing rate sheet date

Turbo checks the effective date only when it creates a rate sheet (ValidateRateSheetEffectiveDate). A rate sheet update runs only the rate sheet processing.

## Category and client changes mark jobs as not imported

Lines 561-597

Asked as: jobs not imported after change; job marked not imported; client change reimport jobs; category change jobs reset

Creating or deleting a category marks its job as not imported. Updating a client or client contact marks the jobs that use it as not imported. This applies only when something other than _imported changed. A client change marks jobs by client and by job owner, and a contact change marks jobs by contact.

## Notification saves update the notification status

Lines 599-613

Asked as: notification save fails; notification status error; notification record save error; notification status update failed

Turbo calls notifications.WriteStatusUpdated whenever a notification record is created or updated, before the save continues. If that call fails, the save fails. These hooks run on internal record events because notifications are not exposed to the API.

## Only the approver shares a time sheet; hooks don't change this

Lines 615-624

Asked as: reviewer cannot approve timesheet; who can share timesheet; reviewer cannot remove sharing; shared timesheet approve permission

A time sheet reviewer row shares a time sheet with another manager. That manager can read it under Shared with Me but cannot approve or reject it. Only the approver of the time sheet can add or remove a row, and the collection's rules enforce this with no state condition. The reviewer hooks in this file don't change who can share.

## Time sheet reviewer must be an active user

Lines 626-647

Asked as: reviewer not active user; cannot share timesheet inactive; selected reviewer is not an active user; share timesheet with inactive

When an approver shares a time sheet, Turbo refuses a reviewer who is not an active user. The message is "reviewer must be an active user", with "the selected reviewer is not an active user" on the reviewer field.

## Vendor absorb commit needs expenses editing enabled

Lines 654-659

Asked as: cannot commit vendor absorb; vendor absorb blocked; absorb commit editing off; vendor absorb refused

Turbo refuses to commit a vendor absorb when the expenses editing setting is off.

## Absorb commit is blocked by a parent absorb

Lines 661-680

Asked as: cannot commit absorb parent; absorb blocked by another absorb; commit or undo first; absorb waiting on parent

Turbo refuses to commit an absorb while an absorb action exists for a parent collection that blocks it. The message is "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first".

## Absorb commit needs the collection's required claims

Lines 684-702

Asked as: no permission to commit absorb; absorb commit denied; missing claims absorb; cannot commit absorb claims

For a collection with its own claim sets, Turbo checks the committer's claims before the commit. If they lack them, Turbo refuses with "you do not have permission to commit a X absorb".

## Absorb commit runs its commit step with the delete atomically

Lines 704-717

Asked as: absorb commit failed nothing kept; absorb delete rollback; commit absorb partly applied; absorb commit all or nothing

If the collection has a commit step (OnCommit), Turbo runs it and deletes the absorb action in one transaction. If the commit step or the delete fails, nothing is kept.

## Sharing a time sheet notifies the new reviewer

Lines 724-755

Asked as: reviewer not notified shared; share timesheet notification; who sent shared notice; shared notice sent to reviewer

When a reviewer is added to a time sheet, Turbo queues a shared notification for the reviewer, naming the time sheet's approver as the sharer. If the notification fails, the share still succeeds.
