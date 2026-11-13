<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"dd2aab0b486147e3f947829685ca68149ed7533eb6bb84718590cc921f002490","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Purchase order save sends its notification only after saving

Lines 27-72

Asked as: email sent after failed save; PO change email not saved; notification when purchase order fails; purchase order email still sent

Creating or updating a purchase order runs ProcessPurchaseOrder, which may queue a notification. Turbo sends the notification only after the save succeeds. If the save fails, Turbo deletes the queued notification, so no email goes out for a PO change that did not persist. A failed send does not undo the save.

## Microsoft sign-in only customizes Microsoft logins

Lines 126-135

Asked as: Microsoft onboarding only for Microsoft; Google login skips setup; other provider login onboarding; OAuth login setup steps

Turbo applies its onboarding steps only when the OAuth provider is Microsoft. Logins through other providers skip them. The whole Microsoft flow runs in one transaction, so a failed step leaves no partly onboarded user.

## Microsoft first sign-in is all-or-nothing

Lines 135-139

Asked as: Microsoft sign-in partly set up; half-created user after Microsoft login; Microsoft login error no token; rollback Microsoft first login

Turbo runs the whole Microsoft sign-in in one transaction. If a later step fails, such as creating the admin profile, the post-save user update or the initial time claim, Turbo undoes the whole sign-in. It does not leave a partly set-up user. The person sees an error and no token is issued.

## Existing user is relinked on first Microsoft login

Lines 150-177

Asked as: existing account linked to Microsoft; duplicate user after Microsoft login; Microsoft account relink; old user Microsoft login new account

When PocketBase finds no user by Microsoft provider id, Turbo looks for a relink candidate (findMicrosoftRelinkCandidate). If it finds one, Turbo relinks the Microsoft account to that user and treats the login as an existing user, not a new one. A lookup or relink failure returns 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## New Microsoft user gets email, name and username from Microsoft

Lines 179-224

Asked as: blank email from Microsoft; Microsoft username not created; Microsoft name missing on new user; fill user details from Microsoft

For a new Microsoft user, Turbo fills a blank email from the fallback email (mail, or userPrincipalName when mail is empty). It fills a blank name with the full Microsoft name. It builds a unique username when none is given. If it cannot build the username, it returns 'failed to prepare Microsoft username'.

## New Microsoft users get an initial time claim

Lines 226-267

Asked as: new user default time claim; first Microsoft login time claim; time claim missing for new user; initial claim after Microsoft sign-in

Turbo assigns an initial time claim (ensureInitialTimeClaim) only on a user's first Microsoft login. A failure returns 'failed to assign initial time claim'.

## Microsoft login syncs identity and repairs account after sign-in

Lines 255-259

Asked as: Microsoft details not updating; sync name from Microsoft; email already used by another user; Microsoft profile refresh on login

After a Microsoft login, Turbo syncs the user's identity (name and email, with verified repair) from Microsoft. A failure returns 'failed to sync Microsoft identity'. Its code comments say username stays fixed after creation. They also say email updates only when no other user already has it.

## Microsoft login ensures an admin profile but no profile

Lines 269-280

Asked as: admin profile missing after Microsoft login; onboarding failed Microsoft; no profile created on Microsoft login; manager chosen in app

Every successful Microsoft login leaves the user with an admin profile. Turbo does not create a profile, because the manager must be chosen in the app. A failure returns 'failed to complete Microsoft user onboarding'.

## Inactive accounts cannot log in

Lines 304-316

Asked as: cannot log in deactivated user; account disabled login blocked; valid password but login refused; user marked inactive cannot sign in

On every login, if the user's admin profile is not active (active is false), Turbo refuses with 'this account is inactive'. No token is issued even when the credentials are valid.

## Missing admin profile: OAuth2 repairs it, other logins are refused

Lines 317-329

Asked as: cannot log in no admin profile; account setup incomplete message; OAuth login missing profile; login refused missing setup

If a user has no admin profile and logged in by OAuth2, Turbo creates the admin profile and lets the login continue. A failure returns 'failed to complete user onboarding'. If the login was not OAuth2, Turbo refuses with 'account setup incomplete'.

## Time entry create and update run time entry checks

Lines 333-344

Asked as: time entry save rejected; time entry validation error; cannot save my hours; time entry check failed

Creating or updating a time entry runs ProcessTimeEntry first. Its errors are returned to the user. The rules inside that function are in another file.

## Deleting time entries and time amendments needs time editing enabled

Lines 345-379

Asked as: cannot delete timesheet entry; delete time amendment blocked; time editing turned off delete; locked time entry deletion

Turbo refuses to delete a time entry or a time amendment when time editing is disabled (utilities.ErrTimeEditingDisabled). It also refuses if the setting cannot be read.

## Time amendment create and update run amendment checks

Lines 357-368

Asked as: amendment won't save; time amendment validation error; cannot change time amendment; amendment check failed

Creating or updating a time amendment runs ProcessTimeAmendment first. Its errors are returned to the user. The rules inside that function are in another file.

## Purchase order create and update validation

Lines 381-404

Asked as: PO won't save; purchase order rejected on save; purchase order validation error; cannot edit purchase order

Creating or updating a purchase order runs ProcessPurchaseOrder before the save, and its errors are returned to the user. The rules inside that function are in another file.

## Other records run their own process checks on create and update

Lines 405-494

Asked as: jobs save rejected; client save validation error; rate sheet record check failed; cannot save profile or currency

Turbo runs a record-specific check (the Process function) before it creates or updates jobs, profiles, currencies, admin profiles, clients, PO approver properties and rate sheets. Client notes get a check on create only. Each function's rules are in other files.

## New admin profile checks the requester's permissions first

Lines 447-461

Asked as: cannot create admin profile permission; new admin profile denied; add admin profile not allowed; requester lacks permission admin

Creating an admin profile first runs EnforceAdminProfileRequestPermissions. If it fails, the error is returned as it is. Then ProcessAdminProfile runs. Updates run only ProcessAdminProfile.

## Expense create and update run in one transaction

Lines 496-527

Asked as: expense save failed uploads; expense attachments removed on error; receipt upload cleaned up; expense not saved receipts gone

Creating or updating an expense runs ProcessExpense and the save in one transaction. If either fails, Turbo cleans up the uploaded expense documents (CleanupExpenseDocumentUploads).

## Expenses editing switch blocks deletes and vendor changes

Lines 528-540

Asked as: cannot delete expense; expenses locked vendor edits; vendor change blocked; expenses editing turned off

When expenses editing is disabled, Turbo refuses to delete expenses and purchase orders. It refuses to create, update or delete vendors. The error is utilities.ErrExpensesEditingDisabled.

## New rate sheet checks effective date first

Lines 542-556

Asked as: rate sheet date rejected; effective date error new rate sheet; rate sheet start date check; cannot set rate sheet date

Creating a rate sheet first runs ValidateRateSheetEffectiveDate, then ProcessRateSheet. Updating a rate sheet runs only ProcessRateSheet, so the effective date check does not run on update.

## Category changes mark the job as not imported

Lines 558-572

Asked as: category change not exported; job not re-imported after category; category edit job not imported; category delete job flag

After a category is created or deleted, Turbo marks its job as not imported (MarkJobNotImported), so the change is written back.

## Client and contact edits mark related jobs not imported

Lines 573-595

Asked as: client edit job not re-imported; contact change jobs not synced; client update job flag; jobs need reimport after client change

After a client update with meaningful changes (ignoring _imported), Turbo marks jobs not imported where the client is the job's client or job owner. After a client contact update with such changes, it does the same for jobs that use that contact. An update with no meaningful change marks nothing.

## Notification saves update the status timestamp

Lines 600-611

Asked as: notification status timestamp; notification updated time; notification save timestamp; status updated not set

Whenever a notification record is created or updated, Turbo calls WriteStatusUpdated first. Notifications are not exposed to the API.

## Time sheet can be shared only with an active reviewer

Lines 623-645

Asked as: cannot share timesheet reviewer; reviewer inactive share fails; share timesheet with inactive user; reviewer not active error

Sharing a time sheet with a reviewer fails with 'reviewer must be an active user' when the chosen reviewer is not active. The field error is 'the selected reviewer is not an active user'. If the check itself fails, Turbo says 'failed to check reviewer active status'.

## Vendor absorb commit needs expenses editing enabled

Lines 652-657

Asked as: vendor absorb commit blocked; cannot commit vendor absorb; absorb commit expenses locked; vendor absorb editing off

Committing a vendor absorb is refused when expenses editing is disabled.

## Absorb commit waits for parent absorb

Lines 659-678

Asked as: cannot commit absorb parent; parent absorb blocks commit; absorb waiting on other absorb; commit or undo absorb first

Committing an absorb for a collection is refused while an absorb exists for a parent collection that blocks it. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'.

## Absorb commit needs the collection's claims

Lines 680-700

Asked as: no permission to commit absorb; absorb commit access denied; missing claims absorb; cannot commit absorb claims

If the collection has its own claim sets, the committer must hold the required claims. Otherwise Turbo refuses with 'you do not have permission to commit a X absorb', where X is the collection name. A failed claims check returns 'failed to check claims'.

## Absorb commit runs its step and delete together

Lines 702-714

Asked as: absorb commit rolled back; absorb delete failed commit; absorb commit step error; absorb not removed after commit

When the collection has an OnCommit step, Turbo runs it and deletes the absorb action in one transaction. If the delete fails, the commit step is rolled back. Collections with no OnCommit step are simply deleted.

## New reviewer gets a time sheet shared notification

Lines 717-753

Asked as: reviewer no notification sent; timesheet share email reviewer; notify reviewer timesheet shared; reviewer added no alert

When a reviewer is added to a time sheet, Turbo queues a notification to the reviewer, naming the time sheet's approver as the sharer. If the notification or the time sheet lookup fails, Turbo logs it and the share still succeeds.
