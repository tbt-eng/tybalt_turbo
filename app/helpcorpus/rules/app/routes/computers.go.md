<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computers.go","sha256":"38a06d6cb3eb2d2a090cf52637439edd1f520b0797e6653013bef3b5aa2fdd24","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computers.go

## Computers can only be changed through the assign, retire and un-retire actions

Lines 3-7

Asked as: edit computer record; change computer details; can't update computer; modify computer via API

The computers collection has no update rule, so nobody can edit a computer record through the API. Assigning, retiring and un-retiring are server routes. A holder of the computers claim or of admin may use them.

## Computer pages show names of last user and assignee

Lines 29-57

Asked as: names missing on computer page; last user name blank; who has the computer; assignee name not showing

The Computers pages show the name of each computer's last user (last_user_name) and of the person it is assigned to (assigned_to_name). Turbo takes these names from the profiles. A holder of the computers claim can't read those profiles through the API.

## Assign is refused if the last user changed since page load

Lines 64-70

Asked as: someone else signed in; computer changed hands; stale computer page; assign refused after reload

If the assign request names a user and that user is no longer the computer's last user, Turbo refuses with "Someone else signed in to this computer since the page loaded. Reload the page and try again." The user is optional and is trimmed of spaces. When it is blank, this check is skipped and Turbo assigns the current last user.

## Who can use the Computers pages and actions

Lines 72-106

Asked as: no access to computers; can't manage computers; you do not have permission; computers permission denied

Only a holder of the computers claim or of admin can list, view, assign, retire or un-retire computers. Anyone else gets the message "you do not have permission to manage computers".

## Single computer view ignores the retired state

Lines 123-132

Asked as: retired computer not found; open retired computer; view computer by ID; retired computer page

Opening one computer by its ID finds it whether it is retired or not. The page lookup does not filter on retired (loadComputerRow lists with retired set to false but passes the ID). The loaded computer is also what every assign, retire and un-retire action returns.

## Computers list shows active computers unless retired is requested

Lines 134-143

Asked as: retired computers list; show old computers; list retired machines; computers list filter

The computers list shows computers that are not retired. It shows retired computers only when the request has retired=1 or retired=true. Any other value gives the non-retired list.

## Viewing a computer that does not exist

Lines 145-156

Asked as: computer missing; can't find computer; no such computer; blank computer page

Viewing one computer that Turbo cannot find gives the message "Computer not found."

## Assigning a computer sets the assignee to its last user

Lines 158-189

Asked as: who gets computer assigned; choose assignee; assign to someone else; computer owner set

Assigning a computer sets the assignee (assigned_to) to the computer's last user (last_user) and the assigned time (assigned_time) to now. The person cannot choose another assignee.

## Cannot assign a retired computer

Lines 171-174

Asked as: assign retired machine; retired laptop assign; can't assign old computer; assign failed retired

Assigning a retired computer is refused with "This computer is retired. Un-retire it before you assign it."

## Cannot assign a computer with no matched last user

Lines 175-179

Asked as: last sign-in no user; computer won't assign; unmatched sign-in; no matching user for computer

Assigning a computer whose last sign-in matched no Turbo user is refused with "The last sign-in on this computer matched no user in Turbo, so it can't be assigned."

## Retiring a computer

Lines 191-202

Asked as: take computer out of service; mark computer retired; decommission computer; retire laptop

Retiring a computer sets its retired time to now. If it is already retired, Turbo refuses with "This computer is already retired."

## Un-retiring a computer

Lines 204-215

Asked as: bring computer back; reactivate computer; put computer back in service; restore retired computer

Un-retiring a computer clears its retired time. If it is not retired, Turbo refuses with "This computer is not retired."

## Computer changes are saved in one transaction

Lines 217-248

Asked as: computer save failed; update conflict computer; computer changes not saved; partial computer update

Assign, retire and un-retire each load the computer, check it and save it in a single transaction. A computer that does not exist gives "Computer not found." Turbo then returns the updated computer as the pages show it.
