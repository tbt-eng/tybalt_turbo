<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_change_client.go","sha256":"2c06dcc321d7b97396f9f3b0d615e15a415e7da2bf08ff3fe468ac3b2d841384","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_change_client.go

## Client change reason length limits

Lines 20-27

Asked as: reason too short; reason too long; character limit reason; how many characters reason

The reason for moving a project to another client must be at least 20 characters (jobClientChangeReasonMinLength) and at most 500 (jobClientChangeReasonMaxLength). Turbo trims spaces from the reason before it counts the characters.

## Client change request fields

Lines 29-36

Asked as: what to send to change client; data needed for client move; fields in client change

The change-client request takes the new client, the contact, the invoicing profile (invoicing_information), the reason and the list of confirmed subjobs.

## Approved requests keep their original client

Lines 59-59

Asked as: approved requests after moving client; approved values keep old client; old client on approved requests

Changing the client does not alter approved requests or value changes. They keep the client they were approved for.

## Changing a project's client requires a signed-in user

Lines 62-64

Asked as: not logged in client change; unauthorized error changing client; session expired move project

Turbo refuses the change-client request when the person is not signed in (e.Auth is nil), with the message "unauthorized".

## Client change reason is trimmed before checks

Lines 70-70

Asked as: extra spaces in reason; trailing spaces reason; whitespace reason saved

Turbo removes leading and trailing spaces from the reason before it checks the length and saves it in the notes.

## Client change runs in one transaction

Lines 72-90

Asked as: half done client change; partial client change saved; client change failed midway

Turbo does the whole client change in one transaction. If any check or save fails, nothing is changed. Turbo sends the notification emails only after the change succeeds.

## Client change success returns the new client

Lines 89-90

Asked as: what id returned after move; client change reply; new client id response

When the change succeeds, Turbo replies with the new client's id.

## Who may change a project's client

Lines 95-105

Asked as: who can move project to another client; permission to change client; no access change client

Only an active user with the job claim can change a project's client. Anyone else gets: "An active user with the job claim is required to change a project's client."

## Client change blocked when job editing is disabled

Lines 106-112

Asked as: cannot edit jobs right now; job editing turned off; client change greyed out

Turbo refuses to change a project's client while job editing is turned off. It returns the jobs-editing-disabled error (ErrJobsEditingDisabled).

## Project must exist to change its client

Lines 114-117

Asked as: project not found moving client; missing project client change; cannot find job to move

If the project is not found, Turbo answers "Project not found."

## Proposal client is changed in the job editor

Lines 118-120

Asked as: change client on proposal; proposal client wrong; P job client change

This action refuses a proposal (a job whose number starts with P). Turbo says: "A proposal's client can be changed in the job editor."

## Subjob cannot change client by itself

Lines 121-123

Asked as: change client of subjob; subjob client cannot change; subjob different client than parent

A subjob always uses its parent project's client. Turbo refuses with "A subjob uses its parent project's client. Change the parent project's client instead."

## Required fields and reason length for client change

Lines 125-141

Asked as: select new client error; missing contact client change; complete highlighted fields

The new client, a contact and a reason are required. Turbo says "Select the new client." and "Select a contact for the new client." when they are missing. The reason is trimmed of spaces and must be 20 to 500 characters. Shorter gives "Explain why the client is changing, in at least 20 characters." and longer gives "Keep the reason to 500 characters or fewer." Turbo reports all of these together with "Complete the highlighted fields." The invoicing profile is optional.

## Invoicing profile is optional in a client change

Lines 128-133

Asked as: no invoicing profile needed; skip invoice profile moving project; invoicing info empty after move

Turbo does not require an invoicing profile when it moves a project. The new client and the contact are required, and so is the reason. If the request sends no profile, the project and its subjobs get an empty profile (invoicing_information). Turbo checks the client, contact and profile against the new client.

## New client must differ and exist

Lines 143-150

Asked as: same client selected; pick different client; client does not exist

If the chosen client is the project's current client, Turbo says "Select a different client." If the chosen client does not exist, it says "The selected client was not found."

## Subjobs must be confirmed exactly

Lines 152-158

Asked as: subjob added while dialog open; subjobs do not match confirmation; confirm subjobs again

The request must list exactly the project's current subjobs. If the list differs, for example because a subjob was added after the dialog opened, Turbo refuses with "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again."

## Client notes follow the project to the new client

Lines 168-223

Asked as: notes after moving client; client notes move with project; old notes new client

Client notes about the project and its subjobs that belong to the old client move to the new client.

## Pending client or contact absorb blocks client change

Lines 173-175

Asked as: client merge pending; absorb waiting commit blocks; contact merge in progress

Turbo refuses if an absorb of clients or contacts, still waiting to be committed or undone, changed references on the project, its subjobs, their client notes, or the chosen contact or invoicing profile. The message is "A client or contact merge that involves this project or the chosen contact or invoicing profile is waiting to be committed or undone. Finish it before changing the client."

## Client change notes must fit 1000 characters

Lines 177-195

Asked as: note too long client change; 1000 characters notes; notes exceed limit move

Turbo records two notes of up to 1000 characters: one on the project and one under the old client. They include both client names and any subjobs. If either would be too long, Turbo refuses and tells the user how many characters to cut from the reason. The old client label is "no client" if there is none.

## Old client is named "no client" when it has none

Lines 177-180

Asked as: no client old client; project had no client; blank old client email

When the project has no old client label, Turbo writes "no client" in the notes and in the email.

## Client change notes list moved subjobs and reason

Lines 182-191

Asked as: what note is written; history note client change; list of moved subjobs note

The note on the project says "Client changed from X to Y." and lists the subjobs that moved, then gives the reason. The note on the old client says "Project N moved to Y." with the same subjob list and reason.

## Client change notes too long: shorten the reason

Lines 192-195

Asked as: reason too long for note; shorten reason by characters; note over limit

If the client names, subjob numbers and reason make either note longer than 1000 characters, Turbo refuses and says to shorten the reason by the number of characters over.

## Client change moves project and subjobs together

Lines 199-216

Asked as: subjobs move with project; all subjobs new client; legacy job zero value move

The project and all its subjobs get the new client, contact and invoicing profile. Turbo checks that the contact and profile belong to the new client. The jobs are marked as not imported, so they are written back to the legacy system. Other fields are not revalidated, so legacy jobs with a zero value or missing branch can still move.

## Moved jobs are marked as not imported

Lines 211-212

Asked as: job written back legacy; not imported flag jobs; resend job to legacy system

Turbo marks each moved job as not imported (MarkImportedFalseIfChanged), so the job is written back to the legacy system.

## Client change records notes on both clients

Lines 225-233

Asked as: note on old client; record reason both clients; history on previous client

Turbo records the reason as a note on the project. It also adds a note under the old client, with no job attached, saying where the project went. The old-client note is skipped if the project had no client.

## Subjobs confirmed must match exactly

Lines 238-250

Asked as: confirmation subjobs mismatch; no confirmation subjobs error; subjob added since dialog

The set of subjobs that the user confirms must equal the project's current subjobs. If a subjob was added since the dialog opened, or the user gave no confirmation, Turbo refuses with "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again."

## Value changes that are not approved block a client change

Lines 252-263

Asked as: value change not approved client; pending value change blocks move; rejected value change client

Turbo refuses a client change when the project or a subjob has a value change that is not Approved, that is, awaiting approval or rejected. The message says to have it approved, or ask an admin to delete it.

## Unapproved value change blocks client change

Lines 256-263

Asked as: attachment awaiting approval move; value awaiting approval client change; admin delete value change

Turbo refuses to change the client while the project or any subjob has a value change (job attachment) that is not Approved, meaning it is awaiting approval or was rejected. The message tells the user to get it approved or ask an admin to delete it.

## Open subjob request blocks client change

Lines 264-270

Asked as: subjob request open client; pending subjob request move; delete subjob request first

Turbo refuses to change the client while a project request for a subjob of this project exists. The message is "A subjob request for this project is still open. Create or delete it before changing the client."

## Pending client or contact absorb blocks a client change

Lines 273-322

Asked as: absorb not committed client; merge clients waiting undo; finish merge first

Turbo refuses a client change while an absorb of clients or client contacts, not yet committed or undone, rewrote a reference on the project, a subjob, their client notes, the chosen contact or the chosen invoicing profile. The message says to finish it first.

## Pending absorb check scope

Lines 278-322

Asked as: which absorb blocks client; merge check which records; absorb affects project client

The pending-absorb check looks at absorb actions on the clients and client contacts collections. It matches any record that the action's updated references touched and that is among the tracked records.

## Client change emails managers and approvers

Lines 324-364

Asked as: who gets email client change; managers notified moved jobs; approvers email client change

Turbo emails the managers of the moved jobs and the people who approved their values (Approved job attachments). People who are no longer active are skipped. The email gives the old and new client, subjobs and reason, and marks recipients who are approvers.
