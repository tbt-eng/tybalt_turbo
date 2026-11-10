<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/hooks.go","sha256":"4838c4c7dd09f63bfa69a22edd5c7adacad36d8ad8505b055d19da8c412c8768","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/hooks.go

## Hook errors return the error body to the user

Lines 19-25

When a record hook fails with a HookError, Turbo sends that error as JSON, with the error's own status, back to the person who made the request. Other errors are returned unchanged.

## Purchase order notification sent only after save

Lines 27-72

When saving a purchase order creates a notification, Turbo sends it right after the save succeeds. If the save fails, Turbo deletes the pending notification so it is never delivered. A send failure is logged and does not undo the save.

## Expense, PO and vendor changes blocked when expenses editing is off

Lines 81-92

When expenses editing is disabled, Turbo refuses to delete expenses and purchase orders. It also refuses to create, update or delete vendors. It returns the expenses-editing-disabled error (ErrExpensesEditingDisabled).

## Only Microsoft OAuth logins get custom onboarding

Lines 127-129

Turbo applies its custom onboarding only to Microsoft OAuth logins. Other providers use the default flow.

## Microsoft login syncs name and email each time

Lines 135-136

After a Microsoft login, Turbo syncs the user's identity (name and email) from Microsoft. If this fails, login is refused with 'failed to sync Microsoft identity'. The whole Microsoft login runs in one transaction, so a failure leaves no half-onboarded user.

## Microsoft login relinks an existing user

Lines 150-177

When Microsoft login finds no provider-id match, Turbo looks for an existing user to relink. If one is found, Turbo relinks the Microsoft account and treats the login as an existing user, not a new one. Failures give 'failed to look up existing Microsoft user' or 'failed to relink Microsoft account'.

## Microsoft login creates new user with derived email

Lines 179-200

When a new user signs in with Microsoft, Turbo sets the email from Microsoft data, falling back to userPrincipalName when `mail` is blank. It only does this when no email is already set.

## New Microsoft user name taken from Microsoft data

Lines 201-205

For a new Microsoft user, Turbo sets the name from the Microsoft full name when it is not blank and no name is already set.

## New Microsoft user gets generated unique username

Lines 207-214

For a new Microsoft user with no username, Turbo derives a unique username from the email. If this fails, login is refused with 'failed to prepare Microsoft username'.

## New Microsoft user gets an initial time claim

Lines 226-267

A user created by their first Microsoft login is given an initial time claim. If this fails, login is refused with 'failed to assign initial time claim'.

## Syncing Microsoft identity on login

Lines 255-259

After Microsoft authentication, Turbo syncs the user's name and email from Microsoft. The email changes only when no other user has it, and the user is verified when the stored email matches the Microsoft identity. The username is not changed on later logins.

## Microsoft login ensures an admin profile exists

Lines 269-276

Every successful Microsoft login ensures the user has an admin profile (admin_profiles). No business profile is created, because fields such as manager must be collected in the app. If this fails, login is refused with 'failed to complete Microsoft user onboarding'.

## Inactive accounts cannot log in

Lines 304-316

When a user logs in successfully, Turbo checks the admin profile (admin_profiles) for the user. If the profile's active flag is off, Turbo refuses the login with "this account is inactive" and issues no token.

## Login without admin profile is refused unless OAuth

Lines 317-330

A user with no admin profile who logs in by a non-OAuth method is refused with 'account setup incomplete'. An OAuth2 user missing one gets it created at login. If that fails, the message is 'failed to complete user onboarding'.

## Time entry create and update run through time entry processing

Lines 333-344

Turbo runs ProcessTimeEntry on every time entry create and every time entry update before it saves the record. If ProcessTimeEntry returns an error, the save stops and the error goes back to the user. The rules inside ProcessTimeEntry are in another file.

## Time entry deletion needs time editing enabled

Lines 345-355

Turbo refuses to delete a time entry when time editing is disabled, and returns the time-editing-disabled error (ErrTimeEditingDisabled).

## Time amendment create and update run through amendment processing

Lines 357-368

Turbo runs ProcessTimeAmendment on every time amendment create and every time amendment update before it saves the record. If ProcessTimeAmendment returns an error, the save stops and the error goes back to the user. The rules inside ProcessTimeAmendment are in another file.

## Time amendment deletion needs time editing enabled

Lines 369-379

Turbo refuses to delete a time amendment when time editing is disabled, and returns the time-editing-disabled error (ErrTimeEditingDisabled).

## Purchase order create and update notifications

Lines 381-404

Creating or updating a purchase order runs the purchase order rules first. Turbo sends any notification created only after the save succeeds, and deletes it if the save fails.

## Client setup and project workflow hooks are registered

Lines 405-407

Turbo registers the client setup hooks (AddClientSetupHooks) and the project workflow hooks (AddProjectWorkflowHooks) alongside the other record hooks. Their rules are in other files.

## Job create and update run through job processing

Lines 408-419

Turbo runs ProcessJob on every job create and every job update before it saves the record. If ProcessJob returns an error, the save stops and the error goes back to the user. The rules inside ProcessJob are in another file.

## Profile create and update run through profile processing

Lines 421-432

Turbo runs ProcessProfile on every profile create and every profile update before it saves the record. If ProcessProfile returns an error, the save stops and the error goes back to the user. The rules inside ProcessProfile are in another file.

## Currency create and update run through currency processing

Lines 434-445

Turbo runs ProcessCurrency on every currency create and every currency update before it saves the record. If ProcessCurrency returns an error, the save stops and the error goes back to the user. The rules inside ProcessCurrency are in another file.

## Admin profile creation needs permission check

Lines 447-461

Creating an admin profile first checks the requester's permissions. Updates skip that check and only run the admin profile processing.

## Admin profile update has no separate permission check

Lines 456-461

On admin profile update, Turbo runs only ProcessAdminProfile. The EnforceAdminProfileRequestPermissions check runs on create only, in this file.

## Client create and update run through client processing

Lines 463-474

Turbo runs ProcessClient on every client create and every client update before it saves the record. If ProcessClient returns an error, the save stops and the error goes back to the user. The rules inside ProcessClient are in another file.

## Client note creation runs through note processing

Lines 476-481

Turbo runs ProcessClientNote on every client note create before it saves the record. Client notes have no update hook in this file. If ProcessClientNote returns an error, the save stops and the error goes back to the user.

## PO approver props create and update run through processing

Lines 483-494

Turbo runs ProcessPOApproverProps on every create and every update of po_approver_props (purchase order approver settings) before it saves the record. If ProcessPOApproverProps returns an error, the save stops and the error goes back to the user. The rules inside ProcessPOApproverProps are in another file.

## Expense saves are transactional with upload cleanup

Lines 496-527

Creating or updating an expense runs in one transaction. If it fails, Turbo cleans up the uploaded expense documents.

## Delete gates for expenses, purchase orders and vendors

Lines 528-540

Turbo applies the expenses-editing check to deleting expenses and purchase orders, and to creating, updating and deleting vendors.

## Rate sheet effective date checked only on create

Lines 542-556

Turbo validates the effective date of a rate sheet only when it is created. Updates run the rate sheet processing without that check.

## Related changes mark jobs as not imported

Lines 559-595

Creating or deleting a category marks its job as not imported. Updating a client marks its jobs (by client or job owner) as not imported. Updating a client contact does the same for jobs that use it as contact. Clients and contacts only trigger this when something other than `_imported` changed.

## Notification records run a status-updated step on create and update

Lines 600-611

Turbo runs notifications.WriteStatusUpdated whenever a notification record is created or updated, and stops the save if it returns an error. What it sets is defined in another file.

## Time sheet reviewer must be an active user

Lines 613-635

Turbo refuses to add a time sheet reviewer who is not an active user. The message is 'reviewer must be an active user', with 'the selected reviewer is not an active user' on the reviewer field.

## Absorb commit blocked by parent absorb

Lines 637-671

Turbo refuses to commit an absorb action while an absorb for a parent collection exists. The message is 'cannot commit X absorb while Y absorb exists; commit or undo the Y absorb first'.

## Vendor absorb commit needs expenses editing enabled

Lines 642-647

Turbo refuses to commit a vendors absorb action when expenses editing is disabled.

## New reviewer gets a time sheet shared notification

Lines 678-709

When a reviewer is added to a time sheet, Turbo queues a shared notification for the reviewer. The sharer is the time sheet's approver. A failure to queue is logged and does not block adding the reviewer.
