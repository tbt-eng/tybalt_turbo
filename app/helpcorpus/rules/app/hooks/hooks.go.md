<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"b44641f9c443eb5aa9f431ee080495a097296857ab31671d66b566acd17154eb","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Notifications and hook errors on PO and other saves

Lines 19-25

Asked as: hook refused save; error message on save; save rejected with status; hook error shown to user

When a hook refuses a save with a structured error, Turbo returns that error's status and message to the person who made the request. Other errors pass through unchanged.

## Purchase order save sends its notification after it persists

Lines 27-72

Asked as: PO notification sent; purchase order notification when saved; notification sent for unsaved PO; PO save failed but notified

Turbo runs ProcessPurchaseOrder before it creates or updates a purchase order. When the save succeeds, Turbo sends the notification that the check created. When the save fails, Turbo deletes that notification, so cron never sends it for a PO that was not saved. An error in sending is logged and does not undo the save.

## Expenses editing switch blocks deletes and vendor changes

Lines 81-92

Asked as: cannot delete expense; expenses editing disabled; can't edit vendors; vendor changes blocked

When expenses editing is disabled, Turbo refuses to delete an expense or a purchase order. It also refuses to create, update or delete a vendor. The error is ErrExpensesEditingDisabled.

## Only Microsoft sign-ins get the custom onboarding

Lines 126-129

Asked as: Microsoft sign-in custom onboarding; login with Google or other provider; non-Microsoft sign-in behaviour; OAuth provider not Microsoft

Turbo applies its custom sign-in handling to the users collection only when the OAuth2 provider is Microsoft, and the name comparison ignores case. Sign-ins through any other provider skip this handling and use PocketBase's default flow.

## Microsoft sign-in changes are all-or-nothing

Lines 135-135

Asked as: Microsoft login half onboarded; partial sign-in setup rollback; Microsoft sign-in failed midway; sign-in transaction rollback

Turbo runs the whole Microsoft sign-in flow in one transaction. It covers relinking, user creation, the identity sync, the initial time claim and the admin profile. If any step fails, Turbo undoes all of them and does not leave a half-onboarded user.

## Microsoft login relinks an existing user

Lines 150-177

Asked as: existing user Microsoft login; link Microsoft account to old user; failed to relink Microsoft account; failed to look up existing Microsoft user

When PocketBase finds no user by provider ID, Turbo looks for an existing user to relink (findMicrosoftRelinkCandidate). If it finds one, it relinks the Microsoft account to that user and treats the login as a returning user, not a new one. A failed lookup returns "failed to look up existing Microsoft user". A failed relink returns "failed to relink Microsoft account".

## Microsoft first login derives email, name and username

Lines 179-224

Asked as: new Microsoft user email blank; username usersNNNNNN; name missing after Microsoft login; userPrincipalName fallback email

For a new Microsoft user, Turbo fills in an empty email from the Microsoft email, falling back to userPrincipalName. It fills in an empty name from the Microsoft full name. It builds a unique username when none is set, instead of PocketBase's generated usersNNNNNN. Turbo does not overwrite values that are already present.

## Microsoft login syncs identity and onboards the user

Lines 226-280

Asked as: sync Microsoft name and email; why no profile created on login; initial time claim new user; username not changed after login

After each Microsoft login, Turbo syncs the user's name and email from Microsoft (syncMicrosoftUserIdentity). It never renames the username after creation. Turbo creates an admin profile if one is missing, but creates no profile record, because the manager must be chosen in the app. A new user also gets an initial time claim. Only Microsoft logins get this treatment, and one failed step rolls back the whole login.

## Inactive accounts cannot sign in

Lines 304-316

Asked as: cannot log in inactive; account deactivated login fails; this account is inactive; login blocked active flag

When a user signs in with valid credentials, Turbo checks the active flag on their admin profile. If the profile is not active, Turbo refuses the login with "this account is inactive" and issues no token.

## Missing admin profile blocks password login

Lines 317-323

Asked as: account setup incomplete; no admin profile password login; can't sign in password no profile; login refused missing profile

If a user who signs in has no admin profile and did not use OAuth2, Turbo refuses the login with "account setup incomplete". Turbo does not create the admin profile in this case.

## OAuth2 login repairs a missing admin profile

Lines 321-329

Asked as: failed to complete user onboarding; OAuth login missing profile; create admin profile at login; fix profile on sign in

If a user who signs in with OAuth2 has no admin profile, Turbo creates one (ensureUserAdminProfile) and then completes the login. If the creation fails, Turbo refuses with "failed to complete user onboarding".

## Time entry and time amendment changes are checked first

Lines 333-368

Asked as: time entry save blocked; time amendment rejected; timesheet entry validation fails; can't save time entry

Turbo runs ProcessTimeEntry before it creates or updates a time entry. It runs ProcessTimeAmendment before it creates or updates a time amendment. A failure from either stops the save and returns the error.

## Deleting time entries and amendments needs time editing enabled

Lines 345-379

Asked as: cannot delete time entry; time editing disabled; can't remove timesheet line; delete amendment blocked

Turbo refuses to delete a time entry or a time amendment when time editing is disabled (ErrTimeEditingDisabled).

## Purchase order create and update run PO processing

Lines 381-404

Asked as: PO save blocked by check; purchase order not saving; PO create fails; PO update refused

Every create or update of a purchase order goes through ProcessPurchaseOrder first. If it fails, the save does not happen.

## Job, profile, currency, client, note and PO approver checks

Lines 412-501

Asked as: job save check failed; client note not saving; currency or profile save error; PO approver properties rejected

Turbo runs a check function before it saves each of these records. ProcessJob covers jobs, ProcessProfile covers profiles, and ProcessCurrency covers currencies. ProcessClient covers clients. ProcessClientNote runs on create only for client notes. ProcessPOApproverProps covers PO approver properties.

## Admin profile creation needs permission checks

Lines 451-465

Asked as: create admin profile permission denied; who can add admin profile; admin profile update rejected; profile creation not allowed

Turbo checks the request permissions (EnforceAdminProfileRequestPermissions) and then ProcessAdminProfile when it creates an admin profile. An update runs only ProcessAdminProfile.

## Admin profile changes revoke sessions

Lines 466-468

Asked as: user logged out after profile change; sessions revoked on deactivate; admin profile deleted logout; tokens revoked inactive profile

Turbo revokes the user's tokens when an admin profile is created inactive, when it is deactivated, and when it is deleted.

## Expense saves run in one transaction with upload cleanup

Lines 503-534

Asked as: expense upload left behind; expense save failed receipts; expense documents deleted on error; expense transaction failure

Turbo checks and saves an expense (ProcessExpense) in a single transaction. If it fails, Turbo cleans up the uploaded expense documents.

## Rate sheet creation validates effective date

Lines 549-564

Asked as: rate sheet effective date error; can't create rate sheet; rate sheet date invalid; rate sheet save rejected

When a rate sheet is created, Turbo first validates its effective date and then runs ProcessRateSheet. An update runs only ProcessRateSheet. The sheet's code and name check run in the model hooks.

## Rate sheet code and name checks run after the hook

Lines 556-557

Asked as: rate sheet code error; rate sheet name duplicate; rate sheet request error; rate sheet validation message

On a rate sheet save, the sheet's code and name checks run in the model hooks after the request hook, and their errors come back as rate sheet request errors.

## Category and client changes mark jobs as not imported

Lines 567-603

Asked as: job marked not imported; jobs need re-import after category change; client change re-import jobs; contact change jobs unimported

Creating or deleting a category marks its job as not imported. A client or client contact update marks the jobs that use it as not imported, but only when something other than _imported changed. For a client, Turbo marks jobs where it is the client or the job owner. For a contact, Turbo marks jobs that name it as the contact.

## Notification updates set a status-updated value

Lines 605-619

Asked as: notification status updated; notification timestamp not set; why can't edit notifications; notifications collection API access

Whenever a notification record is created or updated, Turbo updates its status-updated value through notifications.WriteStatusUpdated before it saves. Only Turbo's own code writes notifications, because the collection is not exposed to the API.

## Time sheet reviewer must be an active user

Lines 621-653

Asked as: reviewer not active; can't share timesheet with user; inactive reviewer selected; share timesheet fails reviewer

When an approver shares a time sheet, the reviewer must be an active user. Otherwise Turbo refuses with "reviewer must be an active user" (the reviewer field says "the selected reviewer is not an active user"). These hooks do not change who can share.

## Absorb commit is blocked by a parent absorb

Lines 657-686

Asked as: can't commit absorb; absorb blocked by parent; commit parent absorb first; absorb commit refused

Turbo refuses to commit an absorb while an absorb action exists for a parent collection. The message is "cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first".

## Vendor absorb commit needs expenses editing enabled

Lines 660-665

Asked as: vendor absorb commit blocked; can't commit vendor absorb; vendor merge disabled expenses; vendor commit expenses editing

Turbo refuses to commit a vendor absorb when expenses editing is disabled.

## Absorb commit needs the collection's claims

Lines 690-708

Asked as: no permission to commit absorb; absorb commit not allowed; who can commit absorb; missing claims absorb

For a collection that has its own claim sets, the person who commits must hold the required claims. Otherwise Turbo refuses with "you do not have permission to commit a X absorb".

## Absorb commit saves history and settles proposals

Lines 710-752

Asked as: absorb history record; absorb proposal accepted; absorb commit outcome notification; absorb failed partway

A commit runs in one transaction. Turbo first copies the action into absorb history, with who absorbed and who committed. It then runs the collection's commit step. For collections whose absorb can be proposed, it accepts the proposal that the action carries out and closes the proposals that name an absorbed record. A failed step leaves nothing behind. After the commit, Turbo queues an outcome notification for each settled proposal.

## New time sheet reviewer gets notified

Lines 755-791

Asked as: reviewer notification not received; shared timesheet notify reviewer; who shared my timesheet email; reviewer added alert

When a reviewer is added to a time sheet, Turbo queues a shared notification for that reviewer, naming the time sheet's approver as the sharer. If the notification fails, the share still succeeds.
