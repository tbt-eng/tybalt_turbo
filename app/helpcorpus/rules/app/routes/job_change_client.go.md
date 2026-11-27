<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_change_client.go","sha256":"7d9e8faa36b32c8ad10cbd62f841a81d0f0505d414001b81da0c12b19de61a65","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_change_client.go

## Reason must be 20 to 500 characters

Lines 20-27

Asked as: reason too short; reason character limit; client change reason length; explain why client changing

Turbo trims the reason and counts characters. Fewer than 20 gives "Explain why the client is changing, in at least 20 characters." More than 500 gives "Keep the reason to 500 characters or fewer." Exactly 20 or 500 is accepted.

## Change is all-or-nothing

Lines 46-92

Asked as: partial client change saved; changes rolled back; client change half applied; transaction failed client

Turbo makes all the changes in one transaction. If any step fails, none of them are kept. It sends the emails only after the transaction succeeds. Approved requests and value changes keep the client they were approved for.

## Only job claim holders can change a project's client

Lines 62-101

Asked as: who can change client; no permission change client; need procurement claim; can't move project to other client

Turbo refuses the client change unless the user holds the job claim (procurement). The message is "A user with the job claim is required to change a project's client." A signed-out request is unauthorized.

## Client change is blocked while job editing is disabled

Lines 102-108

Asked as: jobs editing disabled; can't edit client while locked; editing turned off; client change blocked editing

Turbo checks whether job editing is enabled before it changes a project's client. When editing is disabled, it returns utilities.ErrJobsEditingDisabled and changes nothing.

## Project must exist to change its client

Lines 110-113

Asked as: project missing when changing client; can't find project to move; project not found client

When the project is not found, Turbo answers "Project not found."

## Proposal client is not changed with this action

Lines 114-116

Asked as: change proposal client; proposal wrong client; P job client change

A job whose number starts with P is a proposal. Turbo refuses with "A proposal's client can be changed in the job editor."

## Subjob client cannot be changed directly

Lines 117-119

Asked as: subjob client wrong; change client of subjob; subjob follows parent client

A job with a parent is a subjob, and Turbo refuses with "A subjob uses its parent project's client. Change the parent project's client instead." Subjobs follow their parent project.

## New client, contact and reason are required

Lines 121-137

Asked as: required fields client change; no contact selected; missing new client details; form errors client change

Turbo needs a new client ("Select the new client.") and a contact ("Select a contact for the new client."). It also needs a reason. Turbo checks all of them together and answers "Complete the highlighted fields." with a message on each field. Invoicing profile is optional.

## New client must differ from current client

Lines 139-142

Asked as: same client selected; already this client; pick other client

When the selected client is the project's current client, Turbo answers "Select a different client."

## Selected new client must exist

Lines 143-146

Asked as: client not found when moving; selected client missing; client record gone

When the selected client record is not found, Turbo answers "The selected client was not found."

## Confirmed subjobs must match the project's subjobs exactly

Lines 148-154

Asked as: subjob list changed; new subjob added while moving; confirm subjobs again; subjob confirmation mismatch

The user must confirm exactly the project's current subjobs. Any difference, such as a subjob added after the dialog opened or no confirmation, gives "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again."

## Client notes about moved jobs follow to new client

Lines 164-219

Asked as: client notes after moving job; notes stay with old client; moved job notes

Turbo moves the client notes that belong to the project or its subjobs and name the old client to the new client.

## Notes recording the change must fit 1000 characters

Lines 173-191

Asked as: change note too long; note over 1000 characters; shorten reason client change; client change note limit

Turbo writes a change note and a departure note that include both client names and the subjobs. If either note would exceed 1000 characters, Turbo refuses. It says how many characters to cut from the reason.

## Change is recorded as notes on the project and old client

Lines 173-229

Asked as: history of client change; where client change logged; old client note; record of moved project

Turbo adds a note on the project, "Client changed from X to Y.", with the moved subjobs and the reason. The old client also gets a note, "Project N moved to Y.", with the reason. Turbo skips the old-client note when the project had no client. The old client label is "no client" when it is missing.

## Project and subjobs move together to the new client

Lines 193-212

Asked as: subjobs not moving with project; contact not updated on subjobs; legacy sync client change; invoicing details subjobs

Turbo sets the new client, contact and invoicing information on the project and all its subjobs. It validates the client references on each job. It marks each job as not imported so it is written back to the legacy system. It saves the jobs without validating their other fields.

## Client change refused unless confirmed subjobs match exactly

Lines 234-246

Asked as: subjobs mismatch refused; client change refused subjobs; subjob added after dialog

Turbo compares the subjobs the user confirmed with the project's current subjobs. If the two lists differ in any way, it refuses to change the client. This happens when a subjob was added since the dialog opened, or when no confirmation was given. Turbo shows: "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again."

## Unapproved value change blocks client change

Lines 252-259

Asked as: value change blocks moving project; unapproved value change client; rejected value change client change; pending value change stops client

If the project or any subjob has a value change that is not Approved (awaiting approval or rejected), Turbo refuses. The message says to have it approved, or ask an admin to delete it, before changing the client.

## Open subjob request blocks client change

Lines 260-266

Asked as: subjob request open; request under project blocks client; child request stops client change

If any project request names this project as its parent, Turbo refuses with "A subjob request for this project is still open. Create or delete it before changing the client."

## Pending client or contact absorb blocks client change

Lines 270-318

Asked as: absorb blocking client change; pending absorb contacts; undo absorb clients problem; finish absorb first

Turbo refuses while an absorb of clients or contacts that has not been committed or undone touched the project, its subjobs, their client notes, the chosen contact or the chosen invoicing profile. The message says to finish the absorb first. This stops an undo from overwriting the new relations.

## Email notice goes to managers and approvers who are still active

Lines 320-361

Asked as: who gets client changed email; managers notified client change; approver email moved project; inactive managers email

Turbo queues a job_client_changed email to the managers of the project and its subjobs. It also emails the approvers of their approved value changes. It skips people who are no longer active. Approvers get the ForApprover flag set in the email data.
