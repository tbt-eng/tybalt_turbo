<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"319d576c9da9b62600cbead7b078764c5e3f913f560acc34efe37c54c44626e3","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Purchase order saves send a notification after saving

Lines 27-72

Asked as: purchase order notification not sent; PO saved but notification missing; notification deleted after failed PO save; cron not delivering PO notice

A purchase order create or update runs ProcessPurchaseOrder, which may queue a notification. Turbo sends it right after a successful save. If the save fails, Turbo deletes the queued notification so the cron job does not deliver it. A failure to send is logged and does not undo the save.

## Microsoft login creates a new user from Microsoft data

Lines 126-224

Asked as: new Microsoft user missing email; Microsoft login creates account; username not set on first login; Microsoft name not filled in

Only Microsoft logins are customized. For a new user, Turbo sets the email from the Microsoft mail, or userPrincipalName if mail is blank. It sets the name from the Microsoft full name when one is available. It builds a unique username when none is given. Turbo sets these only if the field is empty.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: Microsoft login treated as new user; existing user Microsoft relink; could not link Microsoft account; failed to relink Microsoft account

When PocketBase finds no user by provider ID, Turbo looks for an existing user to relink. If it finds one, Turbo relinks the Microsoft account to that user and treats the login as a returning user, not a new one. If the lookup or relink fails, Turbo refuses with "failed to look up existing Microsoft user" or "failed to relink Microsoft account".

## Microsoft login syncs identity and onboards the user

Lines 226-280

Asked as: Microsoft login name not updating; login refused after Microsoft sign in; failed to sync Microsoft identity; onboarding failed on Microsoft login

After each Microsoft login, Turbo syncs the user's name and email from Microsoft, and the username stays unchanged. A brand-new user also gets an initial time claim. Every Microsoft login ensures an admin profile exists. Turbo does not create a profile (profiles record) here, because the manager must be chosen in the app. The whole flow runs in one transaction. A failure at any step refuses the login with a message such as "failed to sync Microsoft identity", "failed to assign initial time claim" or "failed to complete Microsoft user onboarding".

## Inactive accounts cannot log in

Lines 304-316

Asked as: cannot log in account disabled; login blocked deactivated user; this account is inactive; valid password but login refused

When a user logs in, Turbo looks up the admin profile. If its active flag is off, Turbo refuses the login with "this account is inactive" and issues no token, even when the credentials are valid. This check runs on every login.

## Missing admin profile blocks non-Microsoft login

Lines 317-323

Asked as: account setup incomplete; password login refused no profile; missing admin profile login error; can't log in with password

If a user has no admin profile and does not log in through OAuth2, Turbo refuses the login with "account setup incomplete". Turbo does not create the missing profile for password logins.

## OAuth login repairs a missing admin profile

Lines 321-329

Asked as: OAuth login no admin profile; SSO login profile missing; failed to complete user onboarding; profile created on Microsoft sign in

If a user logs in through OAuth2 and has no admin profile, Turbo creates one and lets the login continue. If creating it fails, Turbo refuses with "failed to complete user onboarding".

## Time entry and amendment saves are validated

Lines 333-368

Asked as: time entry won't save; time amendment save rejected; timesheet entry validation error; can't update time entry

Creating or updating a time entry runs ProcessTimeEntry. Creating or updating a time amendment runs ProcessTimeAmendment. A rule failure refuses the save.

## Deleting time entries needs time editing enabled

Lines 345-379

Asked as: cannot delete time entry; time editing disabled; delete timesheet entry blocked; can't remove time amendment

Turbo refuses to delete a time entry or a time amendment when the time editing feature flag is off. It returns the time-editing-disabled error (ErrTimeEditingDisabled).

## Purchase order create and update flow

Lines 381-404

Asked as: PO validation before save; purchase order rejected on save; PO create fails checks; PO update blocked

Turbo sends every purchase order create and update through ProcessPurchaseOrder before saving.

## Records for jobs, profiles, currencies, clients, invoices and PO approver props are validated

Lines 408-525

Asked as: job won't save validation; profile save error; currency or client save rejected; invoice line cannot be added or removed

Turbo validates a record on create or update before saving it, through ProcessJob, ProcessProfile, ProcessCurrency, ProcessClient and ProcessPOApproverProps. Client notes are validated on create only. Invoices run ProcessInvoiceCreate, ProcessInvoiceUpdate and ProcessInvoiceDelete. Invoice lines run ProcessInvoiceLineCreate and ProcessInvoiceLineDelete, with no update hook.

## Admin profile creation has an extra permission check

Lines 447-461

Asked as: cannot create admin profile; admin profile permission denied; new admin profile forbidden; editing admin profile works but creating fails

Creating an admin profile first runs EnforceAdminProfileRequestPermissions and then ProcessAdminProfile. An update runs only ProcessAdminProfile.

## Expense saves run in a transaction with upload cleanup

Lines 527-558

Asked as: expense save failed files remain; expense upload cleanup; expense receipt uploads deleted; expense save rolled back

Creating or updating an expense runs ProcessExpense inside one transaction. If the save fails, Turbo cleans up the uploaded expense documents.

## Expense, PO and vendor changes need expenses editing enabled

Lines 559-571

Asked as: cannot delete expense; vendor edit blocked; expenses editing disabled; can't change PO or vendor

Turbo refuses to delete an expense or a purchase order, and to create, update or delete a vendor, when the expenses editing flag is off. It returns the expenses-editing-disabled error (ErrExpensesEditingDisabled).

## Rate sheet effective date is checked only on create

Lines 573-587

Asked as: rate sheet effective date error; rate sheet date rejected; new rate sheet date check; rate sheet update date not checked

Creating a rate sheet runs ValidateRateSheetEffectiveDate and then ProcessRateSheet. An update runs only ProcessRateSheet, so the effective date check does not run on update.

## Category and client changes mark jobs as not imported

Lines 589-626

Asked as: job marked not imported; job needs reimport after category change; client edit unimported jobs; jobs reset after client update

Creating or deleting a category marks its job as not imported (MarkJobNotImported). Updating a client marks the jobs that use it as client or job owner as not imported. Updating a client contact marks the jobs that use it as contact as not imported. Jobs are marked only when the record has meaningful changes, ignoring _imported.

## Notification saves run a status-updated step first

Lines 627-642

Asked as: notification not saved; notification status error; notification write fails; status updated step blocks save

Turbo calls notifications.WriteStatusUpdated before it saves a notification record, on create and on update. If that step returns an error, Turbo does not save the notification. Notifications are not exposed to the API, so only Turbo's own code writes them. What WriteStatusUpdated sets is in another file and is not stated here.

## Time sheet reviewer must be an active user

Lines 644-676

Asked as: cannot share timesheet with reviewer; reviewer inactive share error; share timesheet fails inactive; reviewer not active user

When an approver shares a time sheet, Turbo refuses if the selected reviewer is not an active user. The message is "reviewer must be an active user", and the reviewer field reads "the selected reviewer is not an active user". A failed status check returns "failed to check reviewer active status". Only the approver can share. These hooks do not change that.

## Absorb commit is blocked by a parent absorb

Lines 678-712

Asked as: cannot commit absorb; absorb blocked by parent; commit absorb error parent absorb; undo absorb first

Turbo refuses to commit an absorb action while an absorb for a parent collection exists. The message is "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first". The blocking parent collections come from absorb.GetCommitBlockers.

## Vendor absorb commit needs expenses editing enabled

Lines 683-688

Asked as: vendor absorb commit blocked; cannot commit vendor absorb; vendor import commit disabled; absorb vendors refused

Turbo refuses to commit an absorb of vendors when the expenses editing flag is off.

## Sharing a time sheet notifies the reviewer

Lines 714-750

Asked as: reviewer not notified of shared timesheet; share timesheet notification missing; reviewer gets no share alert; timesheet shared but no email

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to that reviewer. It uses the time sheet's approver as the sharer. If the notification fails, Turbo logs the error and still saves the share.
