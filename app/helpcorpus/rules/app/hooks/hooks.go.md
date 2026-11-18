<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"7685ce98cb7f66081472df953a8cb7535eb4698cc7c9c1a8736fcd0cb60aece2","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors return their message and status to the client

Lines 19-25

Asked as: error message when saving; check failed on save; why won't it save; what does this error mean

When a time entry, time amendment, purchase order, job, profile, currency, client, expense or similar save fails a Turbo check, Turbo sends back the error message and its data fields (AnnotateHookError). The person sees the message that the check gives.

## Purchase order notification is sent after a successful save

Lines 27-72

Asked as: PO email not sent; purchase order notification missing; save failed but email went out; notification sent after failed save

When a purchase order is created or updated, a notification made while processing it is sent right after the save. If the save fails, Turbo deletes that notification so it is not sent later. A failure to send does not undo the save.

## Expense and purchase order deletes and vendor changes need expenses editing enabled

Lines 81-92

Asked as: cannot delete expense; PO delete blocked; vendor edit not allowed; expenses editing turned off

Turbo refuses to delete an expense or a purchase order, and to create, update or delete a vendor, when the expenses editing setting is off. It returns the expenses editing disabled error (ErrExpensesEditingDisabled).

## Only Microsoft logins get custom onboarding

Lines 126-129

Asked as: non Microsoft login onboarding; google login setup; login with other provider; different sign in method setup

Turbo customizes onboarding only when the OAuth provider is Microsoft (compared without regard to case). Logins through other providers use PocketBase's default handling.

## Microsoft login is one transaction

Lines 135-280

Asked as: half set up user after login; partial microsoft signup; microsoft login rollback; user left incomplete after login

Turbo runs the whole Microsoft login in a single transaction. A failure in any step leaves no partly onboarded user.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: microsoft login existing account; account already exists microsoft; link microsoft account to user; microsoft login wrong account

When PocketBase finds no user by provider id for a Microsoft login, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the Microsoft account to that user and treats the login as an existing user, not a new one. If the lookup or relink fails, Turbo refuses with "failed to look up existing Microsoft user" or "failed to relink Microsoft account".

## New Microsoft user gets email, name and username from Microsoft

Lines 179-224

Asked as: missing email after microsoft signup; name blank new microsoft user; username not created; microsoft details not filled in

For a new Microsoft user, Turbo fills a blank email with the fallback email (mail, or userPrincipalName when mail is blank). It fills a blank name with the Microsoft full name when there is one. It builds a unique username when the username is blank. If the username cannot be built, Turbo refuses with "failed to prepare Microsoft username".

## New Microsoft users get the initial time claim

Lines 226-267

Asked as: initial time claim missing; new user no time claim; time claim not assigned; first login no claim

Turbo assigns the initial time claim only to a brand-new Microsoft user, not to a returning one. If assignment fails, Turbo refuses with "failed to assign initial time claim".

## Microsoft login syncs identity after auth

Lines 255-258

Asked as: microsoft name out of date; identity not updated on login; microsoft sync failed; login rolled back identity

After each successful Microsoft login, Turbo syncs the user's identity (name and email) from Microsoft and marks the record verified when needed (syncMicrosoftUserIdentity). If this fails, Turbo refuses with "failed to sync Microsoft identity" and the whole login is rolled back.

## Microsoft login ensures an admin profile but no profile

Lines 269-276

Asked as: no admin profile after microsoft; microsoft onboarding failed; why no profile created; microsoft user setup error

Every successful Microsoft login leaves the user with an admin profile. Turbo does not create a profile (profiles record), because the manager must be chosen in the app. If onboarding fails, Turbo refuses with "failed to complete Microsoft user onboarding".

## Inactive accounts cannot log in

Lines 304-316

Asked as: cannot log in account; account inactive login; locked out of account; login refused valid password

When a user logs in, Turbo checks the active flag on the user's admin profile (admin_profiles.active). If it is false, Turbo refuses the login with "this account is inactive" and issues no token, even when the credentials are valid.

## Non-Microsoft login without admin profile is refused

Lines 317-323

Asked as: account setup incomplete; login refused no profile; missing admin profile login; login fails after setup

If a user logs in by any method other than OAuth2 and has no admin profile, Turbo refuses with "account setup incomplete". Turbo does not create the missing admin profile for such a login.

## OAuth login repairs a missing admin profile

Lines 321-329

Asked as: oauth login no admin profile; fix missing admin profile; onboarding failed on login; sign in creates profile

If a user logs in with OAuth2 and has no admin profile, Turbo creates one (ensureUserAdminProfile) and then completes the login. If creation fails, Turbo refuses with "failed to complete user onboarding".

## Deleting time entries needs time editing enabled

Lines 333-355

Asked as: cannot delete time entry; time entry delete blocked; delete timesheet entry not working; time editing off delete

Turbo refuses to delete a time entry when the time editing setting is off, and returns the time editing disabled error (ErrTimeEditingDisabled). Creating and updating time entries run through ProcessTimeEntry, whose rules are in another file.

## Time entry create and update run the time entry checks

Lines 333-344

Asked as: time entry won't save; time entry check failed; timesheet entry rejected; time entry validation error

Creating or updating a time entry runs the time entry checks (ProcessTimeEntry) before the save. A failed check refuses the save with its message.

## Deleting time amendments needs time editing enabled

Lines 357-379

Asked as: cannot delete amendment; amendment delete blocked; time amendment locked; remove time amendment failed

Turbo refuses to delete a time amendment when the time editing setting is off, and returns the time editing disabled error (ErrTimeEditingDisabled). Creating and updating run through ProcessTimeAmendment, whose rules are in another file.

## Time amendment create and update run the amendment checks

Lines 357-368

Asked as: amendment won't save; time amendment rejected; amendment validation error; amendment save error

Creating or updating a time amendment runs the time amendment checks (ProcessTimeAmendment) before the save. A failed check refuses the save with its message.

## Purchase order save sends its notification

Lines 381-404

Asked as: PO email when saved; purchase order notification sent; PO save stopped; notify approver on PO

Creating or updating a purchase order runs ProcessPurchaseOrder first, and an error from it stops the save. The notification it creates is sent only after the save succeeds.

## Project, job and rate sheet hooks are registered separately

Lines 405-411

Asked as: project workflow hook; rate sheet hook; job rate sheet setup; client setup hook

Turbo registers the client setup, project workflow, job rate sheet, rate sheet draft, rate sheet code and rate sheet entry hooks from other functions. Their rules are not stated in this file.

## Job, profile and currency saves run their own checks

Lines 412-449

Asked as: job save refused; profile save error; currency save failed; cannot update job profile

Creating or updating a job runs ProcessJob. Creating or updating a profile runs ProcessProfile. Creating or updating a currency runs ProcessCurrency. A failed check refuses the save with its message.

## Admin profile creation checks request permissions first

Lines 451-465

Asked as: creating admin profile refused; admin profile permission error; cannot add admin profile; admin profile API error

When someone creates an admin profile through the API, Turbo first runs the admin profile permission check (EnforceAdminProfileRequestPermissions) and returns its error as it is. Then it runs the admin profile checks (ProcessAdminProfile). Updating an admin profile skips the permission check and runs only ProcessAdminProfile.

## Client and client note saves run their own checks

Lines 467-485

Asked as: client save error; client note not saving; add client note failed; update client note problem

Creating or updating a client runs ProcessClient. Creating a client note runs ProcessClientNote. Updating a client note has no hook here.

## PO approver properties saves run their checks

Lines 487-498

Asked as: PO approver settings not saved; approver properties error; save PO approver details; approver setup refused

Creating or updating a PO approver properties record runs ProcessPOApproverProps before the save. A failed check refuses the save with its message.

## Expense create and update run in a transaction

Lines 500-531

Asked as: expense save failed; expense receipt removed; expense not saved receipt; expense upload deleted on error

Creating or updating an expense runs ProcessExpense and the save in one transaction. If either fails, Turbo cleans up the uploaded expense documents.

## Expense deletes, PO deletes and vendor changes are gated

Lines 532-544

Asked as: delete expense blocked; vendor changes refused; PO delete not allowed; cannot edit vendors

The expenses editing gate applies to expense deletes, purchase order deletes and all vendor create, update and delete requests.

## Rate sheet creation checks effective date first

Lines 546-561

Asked as: rate sheet effective date error; rate sheet won't create; effective date rejected; rate sheet date problem

Creating a rate sheet first runs the effective date check (ValidateRateSheetEffectiveDate), then ProcessRateSheet. The sheet's code and name check come from model hooks. Updating skips the effective date check and runs ProcessRateSheet only.

## Category changes mark the job as not imported

Lines 564-577

Asked as: job not imported after category; category change job status; category added job unimported; category delete job status

Creating or deleting a category marks its job as not imported (MarkJobNotImported) once the save succeeds.

## Client and contact edits mark jobs as not imported

Lines 578-600

Asked as: job marked not imported; client change resets job import; contact edit job import; job needs reimport after client edit

When a client is updated with a meaningful change (ignoring the _imported field), Turbo marks as not imported the jobs that use it as client or as job owner. When a client contact changes meaningfully, Turbo marks the jobs that use it as contact. A change to _imported alone does nothing.

## Time sheet reviewer must be an active user

Lines 618-650

Asked as: reviewer inactive error; share timesheet reviewer refused; reviewer not active; reviewer account disabled

When an approver shares a time sheet, the reviewer must be an active user. Otherwise Turbo refuses with "reviewer must be an active user" (the reviewer field says "the selected reviewer is not an active user"). If the check itself fails, Turbo says "failed to check reviewer active status".

## Absorb commit is blocked while a parent absorb exists

Lines 652-683

Asked as: cannot commit absorb; parent absorb blocking; absorb commit blocked; absorb stuck parent

A person cannot commit an absorb for a collection while an absorb for a parent collection exists. Turbo refuses with "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first".

## Vendor absorb commit needs expenses editing enabled

Lines 657-662

Asked as: vendor absorb commit refused; commit vendors absorb blocked; vendors absorb expenses off; vendor absorb not committing

Turbo refuses to commit a vendors absorb when the expenses editing setting is off.

## Absorb commit needs the collection's claims

Lines 685-705

Asked as: no permission commit absorb; absorb claims missing; cannot commit absorb claim; absorb permission denied

When an absorb collection has its own claim sets, the person committing must hold the required claims. Otherwise Turbo refuses with "you do not have permission to commit a X absorb".

## Absorb commit runs its commit step with the delete atomically

Lines 707-720

Asked as: absorb commit failed nothing; absorb delete rolled back; absorb commit step failed; absorb not deleted after error

When a collection has an OnCommit step, Turbo runs it and the absorb action delete in one transaction. If the step fails, nothing is committed. A collection without that step is simply deleted.

## Sharing a time sheet notifies the reviewer

Lines 722-758

Asked as: reviewer not notified of share; shared timesheet notification; who sent share email; reviewer share alert missing

When a reviewer is added to a time sheet, Turbo queues a shared-time-sheet notification to that reviewer, naming the time sheet's approver as the sharer. A failure to find the time sheet or queue the notification is logged and does not block the share.
