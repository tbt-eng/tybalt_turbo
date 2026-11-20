<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computers.go","sha256":"8defcea9f1d487f1f1e976b9106300a7dfc8c0fe96b4b6b2cf003fad7562554e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computers.go

## Computers change only through assign, retire and un-retire

Lines 3-6

Asked as: edit computer record; change computer directly; update computer through api; cannot edit computer

The computers collection has no update rule, so nobody edits a computer record directly through the API. A holder of the computers claim or of admin changes a computer only with the assign, retire and un-retire actions.

## Computer pages show last user and assignee names

Lines 28-61

Asked as: computer page no name; last user name missing; assignee name blank; who has this laptop

The Computers pages show the name of the last user (last_user_name) and of the assignee (assigned_to_name), taken from their profiles. A holder of the computers claim can't read those profiles through the API, so Turbo supplies the names itself.

## Assign request body is optional

Lines 63-69

Asked as: assign without user; assign with no body; request body is not valid; assign button no user picked

An assign action may send a body with a user (user), or no body at all. With no user, Turbo assigns the computer to its last user without the page-load check. With a user, Turbo refuses the assign if it is no longer the last user. If the body can't be read, Turbo answers "The request body is not valid."

## Who can use the Computers pages and actions

Lines 71-99

Asked as: cannot see computers page; no access to computers; permission denied computers; not allowed to manage computers

Only signed-in staff who hold the computers claim or the admin claim can list, view, assign, retire or un-retire computers. Anyone else sees "you do not have permission to manage computers".

## Viewing a computer that does not exist

Lines 116-149

Asked as: computer not found; cannot open computer; computer page empty; missing computer error

Opening one computer that Turbo cannot find gives "Computer not found.". The single-computer view loads with the non-retired filter (retired=0) in loadComputerRow.

## Computers list shows retired or active computers

Lines 127-136

Asked as: see retired computers; show retired list; retired laptops missing from list; hide retired computers

The computers list shows the computers that are not retired. It shows the retired computers only when the request has retired=1 or retired=true. Any other value of retired gives the non-retired list.

## Assign sets the computer to its last user

Lines 151-182

Asked as: who gets assigned computer; assign to me; computer assigned to wrong person; how assignee is picked

Assigning a computer sets who it is assigned to (assigned_to) to the computer's last user (last_user) and sets the assigned time (assigned_time) to now. The assign request does not choose the person.

## Assign fails if the last user changed since page load

Lines 161-176

Asked as: someone else signed in; computer changed while I looked; stale page assign error; reload and try again

The assign request can name the user the person saw. If that user (trimmed of spaces) differs from the computer's current last user, Turbo refuses with "Someone else signed in to this computer since the page loaded. Reload the page and try again." If no user is sent, this check is skipped.

## Cannot assign a retired computer

Lines 164-167

Asked as: assign retired laptop; cannot assign retired computer; retired computer will not assign; assign fails retired

Turbo refuses to assign a retired computer and says "This computer is retired. Un-retire it before you assign it."

## Cannot assign when no last user matched

Lines 168-172

Asked as: last sign-in matched no user; cannot assign unknown sign in; computer has no user to assign; sign-in not matched to anyone

Turbo refuses to assign a computer whose last sign-in matched no Turbo user. The message is "The last sign-in on this computer matched no user in Turbo, so it can't be assigned."

## Retire a computer

Lines 184-195

Asked as: take computer out of service; mark computer retired; already retired error; decommission laptop

Retiring a computer sets its retired time to now. If the computer is already retired, Turbo says "This computer is already retired."

## Un-retire a computer

Lines 197-208

Asked as: bring back retired computer; restore computer; not retired error; reactivate laptop

Un-retiring clears the retired time of a computer. If the computer is not retired, Turbo says "This computer is not retired."

## Actions on an unknown computer are refused

Lines 210-253

Asked as: retire missing computer; assign deleted computer; computer transaction failed; unknown computer action error

Assign, retire and un-retire on a computer that does not exist give "Computer not found.". Each action loads, checks and saves the computer in one transaction.

## Retired computers can be viewed after un-retire or retire actions

Lines 236-240

Asked as: computer disappears after retire; page after assign shows; computer vanished after un-retire; reload after retiring

After an assign, retire or un-retire, Turbo returns the computer as the pages show it, with the last user's and assignee's names. This reload uses the non-retired filter (retired=0).
