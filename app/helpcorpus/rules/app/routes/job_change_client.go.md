<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_change_client.go","sha256":"a54ea9b19026bea2fb5adab55b8e6e4fb18849a459daf2b9583dbed0b496b503","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_change_client.go

## Reason and client note length limits

Lines 20-27

The reason for a client change must be 20 to 500 characters (jobClientChangeReasonMinLength, jobClientChangeReasonMaxLength). The notes that record the change must fit in 1000 characters (clientNoteMaxLength). Turbo trims spaces from the ends of the reason before it counts.

## Approved requests and value changes keep their original client

Lines 46-59

Changing a project's client does not alter the client that approved project requests and approved value changes were approved for. Outside this action a project's client is fixed. Only an active holder of the job claim can change it, and only with a reason.

## Who may change a project's client

Lines 62-105

Only an active user who holds the job claim may change a project's client. Anyone else gets: "An active user with the job claim is required to change a project's client." A request with no signed-in user is refused as unauthorized.

## Client change requires a signed-in user

Lines 62-64

Turbo refuses a request to change a project's client from someone who is not signed in. The response is "unauthorized".

## Reason length limits for a client change

Lines 70-141

The reason is trimmed of surrounding spaces, then must be 20 to 500 characters. A shorter reason gets "Explain why the client is changing, in at least 20 characters." A longer one gets "Keep the reason to 500 characters or fewer." Exactly 20 and exactly 500 are accepted.

## Client change is all or nothing

Lines 73-90

The whole client change runs in one transaction, so a failure part-way leaves nothing changed. Emails are sent only after it succeeds. Approved requests and value changes keep the client they were approved for.

## Client change blocked when job editing is disabled

Lines 106-112

Turbo refuses the client change when job editing is turned off. The refusal is utilities.ErrJobsEditingDisabled, and its message comes from another file. The check runs after the permission check and before the project is looked up.

## Project must exist to change its client

Lines 114-117

If the project id matches no job, Turbo refuses with "Project not found."

## Proposal client is not changed with this action

Lines 118-120

A job whose number starts with P is a proposal. Turbo refuses with "A proposal's client can be changed in the job editor."

## Subjob client cannot be changed directly

Lines 121-123

A subjob (a job with a parent) cannot have its client changed. Turbo refuses with "A subjob uses its parent project's client. Change the parent project's client instead."

## Required selections for a client change

Lines 127-144

The new client, a contact and an invoicing profile are all required. The messages are "Select the new client.", "Select a contact for the new client." and "Select an invoicing profile for the new client." Turbo reports every missing field at once, together with any reason error, under "Complete the highlighted fields."

## New client must differ from the current client

Lines 146-149

If the chosen client is the project's current client, Turbo refuses with "Select a different client."

## New client must exist

Lines 150-153

If the chosen client record is not found, Turbo refuses with "The selected client was not found."

## Subjobs must be confirmed exactly

Lines 155-161

The subjobs the user confirmed must be exactly the project's current subjobs. Otherwise Turbo refuses with "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again." This also catches a subjob added after the dialog opened.

## Subjobs move with the project

Lines 162-219

The project and all of its subjobs take the new client, contact and invoicing profile. Turbo checks that the contact and profile belong to the new client. Moved jobs are marked as not imported (via MarkImportedFalseIfChanged) so they are written back to the legacy system. They are saved without full validation of other fields, so old projects with a zero value or missing branch can still move.

## Client notes move to the new client

Lines 171-226

Client notes attached to the project or its subjobs under the old client are moved to the new client.

## Audit notes record the client change

Lines 180-236

Turbo adds a note to the project: "Client changed from X to Y." with the subjobs moved (if any) and the reason. The old client gets a note "Project N moved to Y." with the reason, unless the project had no old client. If the old client has no label, the note uses "no client".

## Combined note length limit

Lines 193-198

Each recorded note, with the client names, subjobs and reason, must fit 1000 characters. If either is longer, Turbo refuses and tells the user how many characters to cut from the reason.

## Subjob confirmation must match the current subjobs exactly

Lines 241-253

When a person changes a project's client, the list of confirmed subjobs must be exactly the project's current subjobs, no more and no fewer. If it differs, Turbo refuses with "The subjobs confirmed do not match this project's subjobs, which may have changed. Review the current list and confirm again." This stops a subjob added after the dialog opened, or a missing confirmation, from moving a job the person did not see.

## Pending value change blocks client change

Lines 259-266

Turbo refuses if the project or any subjob has a value change (job attachment) whose status is not Approved, meaning awaiting approval or rejected. The message says to have it approved, or ask an admin to delete it, before changing the client.

## Open subjob request blocks client change

Lines 267-273

Turbo refuses if any job request exists with this project as its parent. The message is "A subjob request for this project is still open. Create or delete it before changing the client."

## Pending client or contact merge blocks client change

Lines 281-309

Turbo refuses while an undoable client or contact merge involves the project, its subjobs, their client notes, or the chosen contact or invoicing profile. The message says to finish the merge (commit or undo) before changing the client.

## Client change emails managers and approvers

Lines 311-351

Turbo queues a job_client_changed email to the managers of the moved jobs and the approvers of their approved value changes. Users who are no longer active get no email. Approvers receive it flagged as approvers (ForApprover).
