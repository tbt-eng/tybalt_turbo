<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_change_client.go","sha256":"0587968120e4de0026b12da5a1b6b73ab0b72e3814585c926519c78b831f880b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_change_client.go

## Reason length limits for a client change

Lines 20-27

Asked as: reason too short; reason character limit; how long must the reason be; reason too long for client change

The reason for moving a project to another client must have at least 20 characters and at most 500 characters. Turbo trims the spaces at both ends of the reason before it counts them.

## Approved requests and value changes keep their old client

Lines 59-59

Asked as: approved requests after client change; approved value changes keep old client; old client on approved request; client change approved items

A client change moves the project, its subjobs, their client notes and the contact and invoicing profile. It does not change approved project requests or approved value changes. They keep the client they were approved for.

## Client change requires a signed-in user

Lines 62-64

Asked as: change client not logged in; client change unauthorized; no user signed in client change; change client without login

The change-client action on a project refuses a request that has no signed-in user and answers 'unauthorized'.

## Client change runs as one all-or-nothing step

Lines 72-90

Asked as: client change partly failed; all or nothing client change; client change rollback; moved some jobs but not others

Turbo changes the client of the project, its subjobs, the client notes and the new notes in one transaction. If any step fails, nothing changes. Turbo queues the email notice only after the transaction succeeds.

## Client change returns the new client on success

Lines 89-90

Asked as: what does client change return; client change success response; new client id after change; confirmation after moving client

When the change succeeds, Turbo answers with the id of the new client.

## Who can change a project's client

Lines 95-105

Asked as: who can move project to another client; cannot change client permission; need procurement role to change client; job claim client change

Only an active user who holds the job claim (procurement) can move a project to another client. Anyone else gets: "An active user with the job claim is required to change a project's client."

## Client change blocked when job editing is off

Lines 106-112

Asked as: cannot change client job editing off; client change disabled; job editing disabled error; editing turned off client

Turbo refuses to change a project's client while job editing is disabled. It returns the error utilities.ErrJobsEditingDisabled; its message text is defined in another file.

## Project must exist

Lines 114-117

Asked as: project not found client change; missing project when changing client; cannot find project to move

If the project cannot be found, Turbo says "Project not found."

## Proposal client is changed in the job editor

Lines 118-120

Asked as: change client of proposal; proposal client cannot change; P job client change; proposal client job editor

A job whose number starts with P (a proposal) cannot use this action. Turbo says "A proposal's client can be changed in the job editor."

## Subjob client cannot be changed directly

Lines 121-123

Asked as: change subjob client; subjob client wrong; why can't I move subjob; subjob different client than parent

A subjob cannot have its client changed by itself. Turbo says "A subjob uses its parent project's client. Change the parent project's client instead."

## Required fields and reason length

Lines 125-144

Asked as: missing fields client change; select contact required; which fields needed to move client; highlighted fields client form

Turbo requires a new client ("Select the new client."), a contact ("Select a contact for the new client.") and an invoicing profile ("Select an invoicing profile for the new client."). The reason is trimmed of surrounding spaces and must be 20 to 500 characters. Shorter gives "Explain why the client is changing, in at least 20 characters." and longer gives "Keep the reason to 500 characters or fewer." Turbo reports all missing fields together as "Complete the highlighted fields."

## New client must differ and exist

Lines 146-153

Asked as: same client selected; pick another client; selected client not found; client does not exist when moving

The new client must differ from the project's current client, or Turbo says "Select a different client." If the selected client does not exist, Turbo says "The selected client was not found."

## Subjob confirmation must match exactly

Lines 155-161

Asked as: subjob list changed; confirm subjobs again; subjob confirmation wrong; subjobs do not match when moving client

The user must confirm exactly the project's current subjobs. If the confirmed list differs, Turbo refuses with "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again."

## Client notes follow the project

Lines 171-226

Asked as: client notes after moving project; notes stay with old client; client notes move with project; old client notes transfer

Client notes about the moved jobs that belong to the old client move to the new client.

## Combined notes must fit 1000 characters

Lines 180-198

Asked as: note too long client change; shorten reason notes too long; client change note limit; reason too long for notes

Turbo builds two notes from the old and new client names, any subjobs and the reason. Each must fit 1000 characters. If either is too long, Turbo asks the user to shorten the reason by the excess.

## Change is recorded in notes

Lines 180-236

Asked as: history of client change; where is client change recorded; client moved note on project; departure note old client

Turbo adds a note to the project: "Client changed from X to Y.", with subjobs moved and the reason. If the project had an old client, that client also gets a note "Project N moved to Y." with the reason. A project with no old client shows "no client" and gets no departure note.

## Project and subjobs move together

Lines 200-219

Asked as: subjobs contact not updated; invoicing profile moved with project; jobs imported flag after client move; contact must belong to new client

The project and all its subjobs take the new client, contact and invoicing profile in one transaction. Turbo checks that the contact and profile belong to the new client. It saves the jobs without checking their other fields, and marks them as changed from the legacy import (imported false).

## Subjob confirmation mismatch message

Lines 241-253

Asked as: subjob confirm error message; extra subjob confirmation error; missing subjob confirmation; confirmed list mismatch

If the confirmed subjobs are not exactly the project's current subjobs, Turbo refuses with 'The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again.' A missing or extra subjob causes the refusal.

## Unapproved value change blocks client change

Lines 259-266

Asked as: value change blocks client move; unapproved job attachment client change; rejected value change stops move; awaiting approval blocks client change

Turbo refuses if the project or any subjob has a value change (job attachment) whose status is not Approved, whether awaiting approval or rejected. The message says to have it approved, or ask an admin to delete it, first.

## Open subjob request blocks client change

Lines 267-273

Asked as: open subjob request blocks; subjob request stops client change; delete subjob request to move client; request still open client change

Turbo refuses while any project request names this project as its parent. It says "A subjob request for this project is still open. Create or delete it before changing the client."

## Pending client or contact merge blocks client change

Lines 277-325

Asked as: merge pending client change; absorb not committed; client merge blocks move; undo pending merge contact

Turbo refuses while a client or contact absorb that rewrote references on the project, its subjobs, their client notes, or the chosen contact or invoicing profile is still waiting to be committed or undone. The message says to finish it first.

## Email notice after client change

Lines 327-367

Asked as: who gets client changed email; notify managers client change; approvers email job moved; client change notification

Turbo queues a job_client_changed email to the managers of the moved jobs and the approvers of their approved value changes. People who are no longer active get none. The email marks which recipients are approvers.
