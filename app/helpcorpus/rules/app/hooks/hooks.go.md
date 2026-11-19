<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"a2444ef7a6406238df5c564492ff4e6b984a5cb6ac98dd7d61c37cf69fdb7c05","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors return their own status and message

Lines 19-25

Asked as: hook error message; generic error instead of real message; check failed but error is vague; error details not shown

When a Turbo check fails with a hook error, Turbo shows that error's message and data to the person, not a generic error.

## Failed purchase order save removes its pending notification

Lines 27-72

Asked as: notification sent for unsaved PO; PO save failed but email went out; purchase order notification after failed save; cron sent PO alert wrongly

When a purchase order create or update fails after Turbo prepares its notification, Turbo deletes that pending notification. The cron job then does not send a notification for a purchase order change that was not saved. Turbo sends a notification only after a successful save.

## Expenses editing off blocks the expense gate and its error

Lines 81-92

Asked as: cannot edit expenses; expenses locked; expense edits blocked; expenses editing turned off

Turbo checks the expenses editing setting for expense, purchase order and vendor changes. When the setting is off, Turbo returns the error ErrExpensesEditingDisabled and refuses the action. If Turbo cannot read the setting, it returns that error instead.

## Microsoft login creates an admin profile but no profile

Lines 126-133

Asked as: microsoft login no profile; manager missing after microsoft sign in; onboarding failed microsoft login; admin profile after microsoft login

Every successful Microsoft login makes sure the user has an admin profile (ensureMicrosoftUserOnboarded). It does not create a profile, because the manager is not available from Microsoft and the user must enter it later in the app. A failure gives 'failed to complete Microsoft user onboarding'. The whole Microsoft login runs in one transaction, and only Microsoft logins get this treatment.

## Microsoft login relinks an existing user instead of creating one

Lines 150-177

Asked as: microsoft account already exists; duplicate user microsoft login; link microsoft to existing account; microsoft login creates second user

When PocketBase finds no user by Microsoft provider id, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the Microsoft account to that user and treats the login as an existing user. Otherwise it creates a new user. A failed lookup or relink gives 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## Microsoft first login builds name, email and username

Lines 179-224

Asked as: name missing after microsoft signup; username from microsoft email; new microsoft user email blank; microsoft sign in fills details

When a new user signs in with Microsoft, Turbo fills the email from Microsoft's mail, or from userPrincipalName when mail is empty. It sets the name from givenName plus surname, then displayName. It derives the username from the email local-part and makes it unique. It only fills these when the sign-in did not already supply them.

## New Microsoft users get an initial time claim

Lines 226-267

Asked as: no time claim new user; initial claim missing microsoft; first microsoft login time claim; time claim not assigned

Only on a user's first Microsoft login does Turbo assign an initial time claim (ensureInitialTimeClaim). A failure gives 'failed to assign initial time claim'.

## Microsoft login syncs name and email each time

Lines 234-259

Asked as: microsoft name not updating; name changed in microsoft not showing; email not synced microsoft; username stays same after microsoft login

After every Microsoft login, Turbo syncs the user's identity from Microsoft: the name, and the email when no other user already has it (syncMicrosoftUserIdentity). It marks the record verified when the stored email matches the Microsoft identity. It never renames the username on later logins. A failure gives 'failed to sync Microsoft identity'.

## Inactive accounts cannot log in

Lines 308-316

Asked as: cannot log in disabled user; account locked out; login rejected valid password; user deactivated cannot sign in

When a user's admin profile is not active, Turbo refuses every login method with 'this account is inactive', even if the credentials are valid.

## Missing admin profile at login

Lines 317-329

Asked as: account setup incomplete; no profile at login; oauth login profile missing; user onboarding failed login

A user with no admin profile who logs in without OAuth2 gets 'account setup incomplete'. A user who logs in with OAuth2 gets an admin profile created and then logs in. If that creation fails, Turbo says 'failed to complete user onboarding'.

## Time entry and time amendment edits are gated

Lines 333-379

Asked as: cannot delete time entry; time amendment delete blocked; time editing turned off; cannot remove timesheet entry

Turbo checks time entries and time amendments with its own rules on create and update (ProcessTimeEntry, ProcessTimeAmendment). Deleting either record is refused with the time editing disabled error (ErrTimeEditingDisabled) when the time editing flag is off.

## Purchase order notification sent only after a successful save

Lines 381-404

Asked as: po notification after failed save; purchase order email not sent; send failed po save fine; pending po alert removed

Turbo checks purchase orders on create and update (ProcessPurchaseOrder). A notification made during that check is sent right after the save succeeds. When the save fails, Turbo deletes the pending notification so the cron never sends it. A failed send does not undo the save.

## Client setup, project workflow and rate sheet hooks are registered

Lines 405-411

Asked as: client hooks; project request workflow rules; rate sheet hooks; rate sheet entry checks

Turbo registers separate hooks for clients, project requests and workflow, job rate sheets, rate sheet drafts, rate sheet codes and rate sheet entries. This file does not show their rules.

## Job, profile and currency saves run their own checks

Lines 412-449

Asked as: job save refused; profile save fails; currency save rejected; job update blocked check

Creating or updating a job, a profile or a currency runs ProcessJob, ProcessProfile or ProcessCurrency first. Turbo refuses the save when that check fails. This file does not show those checks.

## Admin profile updates run ProcessAdminProfile

Lines 450-465

Asked as: cannot update admin profile; admin profile permission error; admin profile save refused; permission error admin profile

Updating an admin profile runs ProcessAdminProfile before the save. Creating one first runs a permission check. Turbo returns the permission error as it is, without the hook error format.

## Admin profile creation has an extra permission check

Lines 451-465

Asked as: cannot create admin profile; admin profile creation denied; new admin profile permission; admin profile request not allowed

Creating an admin profile first passes a request permission check (EnforceAdminProfileRequestPermissions) before the usual admin profile validation. Updates skip that check and run only the validation (ProcessAdminProfile).

## Client, note, invoice and PO approver saves run their own checks

Lines 466-529

Asked as: client save refused; invoice update blocked; PO approver settings rejected; client note not saved

Turbo runs ProcessClient for client create and update, and ProcessClientNote for client note create. It runs ProcessInvoiceCreate, ProcessInvoiceUpdate and ProcessInvoiceDelete for invoices, and the invoice line create and delete checks for invoice lines. It runs ProcessPOApproverProps for PO approver settings create and update. Turbo refuses the action when a check fails. This file does not show those checks.

## Expense and vendor changes blocked when expenses editing is off

Lines 530-575

Asked as: cannot delete vendor; vendor edits blocked; cannot delete expense; expenses locked for vendors

When the expenses editing flag is off, Turbo refuses with the expenses editing disabled error (ErrExpensesEditingDisabled) to delete expenses or purchase orders, and to create, update or delete vendors. Expense create and update run their own checks (ProcessExpense) in a transaction. Turbo removes uploaded documents if the save fails.

## Expense create and update run in one transaction

Lines 530-562

Asked as: expense files left behind; receipt upload after failed save; expense save failed uploads; expense attachments not removed

Creating or updating an expense runs its checks and the save in one transaction. When the save fails, Turbo cleans up the document files uploaded with the expense.

## Rate sheet creation checks effective date first

Lines 577-592

Asked as: rate sheet effective date error; rate sheet creation rejected; effective date wrong rate sheet; cannot create rate sheet

Creating a rate sheet checks the effective date first, then runs the rate sheet checks. Updating skips the effective-date check and runs only the rate sheet checks.

## Changes to categories, clients and contacts mark jobs not imported

Lines 595-631

Asked as: job marked not imported; job needs reimport; client change unmarks job; category change reimport jobs

Creating or deleting a category marks its job as not imported. Updating a client with a meaningful change (other than the _imported field) marks every job that uses it as client or owner as not imported. Updating a client contact does the same for jobs that use it as contact. A change to _imported alone does nothing.

## Notifications and Microsoft login set status updated time

Lines 633-647

Asked as: status updated time missing; notification timestamp wrong; last updated not set; microsoft login updated time

When a notification is created or updated, Turbo runs WriteStatusUpdated before the save.

## Time sheet reviewer must be active

Lines 660-681

Asked as: cannot share timesheet; reviewer not active; timesheet reviewer rejected; share timesheet with inactive user

When an approver shares a time sheet, the reviewer must be an active user. Otherwise Turbo refuses with 'reviewer must be an active user' (the reviewer field says 'the selected reviewer is not an active user'). These hooks do not change who may share.

## Vendor absorb commit blocked when expenses editing is off

Lines 688-693

Asked as: vendor absorb refused; cannot commit vendor absorb; absorb blocked expenses off; vendor absorb locked

Committing a vendor absorb is refused with the expenses editing disabled error (ErrExpensesEditingDisabled) when the expenses editing flag is off.

## Absorb commit blocked while a parent absorb exists

Lines 695-714

Asked as: cannot commit absorb; parent absorb blocking; absorb needs undo first; absorb already pending

Committing an absorb is refused with 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first' when an absorb action for a parent collection still exists.

## Absorb commit needs the collection's claims

Lines 716-736

Asked as: no permission to commit absorb; absorb permission denied; cannot commit collection absorb; absorb claims missing

When a collection defines its own claim sets, committing its absorb needs those claims. Otherwise Turbo refuses with 'you do not have permission to commit a <collection> absorb'.

## Absorb commit runs its commit step with the delete atomically

Lines 738-750

Asked as: absorb commit partly failed; absorb delete not rolled back; commit absorb transaction; absorb action left behind

When the collection has a commit step (OnCommit), Turbo runs it and the delete of the absorb action in one transaction. If either fails, nothing is kept.

## Sharing a time sheet notifies the reviewer

Lines 758-789

Asked as: reviewer not notified; timesheet share no notification; shared timesheet email missing; notify reviewer on share

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to the reviewer, with the time sheet's approver as the sharer. A failure to queue the notification is logged and does not block the share.
