<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_close_api.go","sha256":"7990be6cda7a16acd975ccc6f4dc05eec87bb55d7243e48725561cd4463c4474","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_close_api.go

## Close project requires sign-in and a job id

Lines 45-52

Asked as: not logged in on close; close button unauthorized; close fails no job id; blank job id error

Turbo refuses the Close project action when the person is not signed in ("unauthorized") or when the job id is missing or blank ("missing job id").

## Close happens in one transaction

Lines 59-228

Asked as: close half saved; some changes missing after close; nothing saved after close error; partial save when closing

The close, the proposal auto-award and the audit notes run in a single transaction. If any step fails, none of them is saved.

## Close project fails when the job does not exist

Lines 60-64

Asked as: can't find job to close; close says job missing; job does not exist error; job id doesn't exist

If no job matches the id, Turbo refuses the close with "job not found".

## Close project is checked before the job claim error

Lines 60-75

Asked as: job missing error without access; not found even without permission; error shown before permission check

Turbo looks up the job before it checks the job claim. A missing job returns "job not found" even for a user without the claim.

## Closing a project requires the job claim

Lines 66-75

Asked as: no permission to close project; can't close project access; managers cannot close projects; close project forbidden

Turbo refuses to close a project when the person lacks the job claim. The message is "the job claim is required to close projects". Managers without the job claim can change allocations only.

## Proposals cannot be closed with Close project

Lines 77-80

Asked as: can't close proposal; proposal number starts with P; close a proposal blocked; proposals cannot be closed

Turbo treats a job whose number starts with P as a proposal. It refuses to close it and says "only projects can be closed with this endpoint".

## Only Active projects can be closed

Lines 82-85

Asked as: project not active closing; close greyed out for status; why can't close planned project; close fails wrong project status

Close project works only on a project whose status is Active. For any other status Turbo says "only Active projects can be closed".

## Referenced proposal missing blocks the close

Lines 88-94

Asked as: linked proposal gone; proposal deleted before close; close fails missing proposal; project's proposal can't be found

If the project refers to a proposal that Turbo cannot find, the close fails with "referenced proposal was not found".

## Terminal proposal status blocks closing the project

Lines 96-103

Asked as: cancelled proposal blocks close; no bid proposal closing error; not awarded proposal can't close; proposal ended can't close project

If the referenced proposal is Not Awarded, Cancelled or No Bid, Turbo refuses the close. The message is "referenced proposal has a terminal status and cannot be auto-awarded for close".

## Referenced proposal must be Awarded before close

Lines 110-158

Asked as: proposal not awarded close error; submitted proposal blocks close; in progress proposal can't close; proposal must be awarded first

A project whose proposal is Awarded closes normally. If the proposal is in any other status and is not an eligible imported one, Turbo refuses with "referenced proposal must be Awarded before project close".

## Imported proposals are auto-awarded when the project closes

Lines 117-150

Asked as: imported proposal becomes awarded; auto award on project close; imported proposal status changes on close; closing sets proposal awarded

If the referenced proposal is In Progress or Submitted and was imported (_imported), closing the project sets the proposal to Awarded and clears _imported. Turbo adds a note to the proposal: "Proposal auto-awarded when project <number> (<id>) was closed". The proposal is not auto-awarded when it was created in Turbo or has any other status.

## Auto-award note leaves status-changed field empty

Lines 129-145

Asked as: status changed field empty on note; why no status on proposal note; auto award note blank status; note leaves status changed blank

When Turbo auto-awards the proposal during a project close, it writes the audit note "Proposal auto-awarded when project <number> (<id>) was closed" on the proposal. It does not set job_status_changed_to on this note. That field is for No Bid and Cancelled proposal close-out comments.

## Close skips billing checks and clears imported flag

Lines 161-172

Asked as: close without invoicing profile; close needs no project value; close project skips billing; closing ignores authorizing document

Close project sets the status to Closed and runs the same job rules as a status-only close in the editor. It needs the job claim but not an invoicing profile, project value or authorizing document. The status change clears _imported. If those rules refuse, Turbo shows their error.

## Close project sets the project status to Closed

Lines 163-174

Asked as: project status becomes closed; close sets status to closed; closing changes project status; status changes to closed

Close project sets the project's status to Closed. The save runs only after the job rules in hooks.ProcessJobCore accept it.

## Schema errors on close name the field

Lines 174-204

Asked as: validation failed on close; number invalid when closing; field error when closing project; imported number malformed close

If saving the closed project fails schema validation, such as imported data with a malformed number, Turbo reports "validation failed". Each field message is prefixed with the field name, with underscores as spaces, for example "number: Invalid value format.".

## Closing a project adds a client note

Lines 206-218

Asked as: note added when closing; project closed note client; where is the close note; close adds client note

After saving, Turbo adds a client note on the project: "Project closed with the Close project action". The note fails with an error if the job has no client.

## Close response reports result

Lines 220-226

Asked as: what does close return; close result details; response after closing project; close reply fields

On success Turbo returns the project id, its status, its imported flag, that a project note was added, and the proposal result (whether it was auto-awarded, with its from and to status).

## Close project failures show code and message, or field errors

Lines 230-242

Asked as: close error code message; close returns error message; field errors when closing; close failed what message

When Close project fails, Turbo returns a code and message such as "job not found" or "only Active projects can be closed". A validation failure returns "validation failed" with per-field errors instead.

## Successful Close project returns the closed project's result

Lines 244-244

Asked as: close success response; project closed reply fields; returned after closing project; close result shows proposal

When the close succeeds, Turbo returns the project id, its new status, its imported flag and whether notes were added. The response includes the proposal result when the project references a proposal.

## Audit note on a job needs the job's client

Lines 260-264

Asked as: job has no client; client required error; close fails no client; job without client on close

Turbo writes each audit note for a job as a client note on the job's client. If the job has no client (blank after trimming), Turbo refuses with the message "client is required to create audit note". This applies to the project close note and to the proposal auto-award note.

## Client notes without a job are marked not applicable

Lines 269-291

Asked as: client note with no job; job not applicable flag; note created without job; client note blank job

When Turbo creates a client note with no job, it sets job_not_applicable to true. The job_status_changed_to field is set only when a value is given.

## Client notes record the author and trim the note text

Lines 275-287

Asked as: who wrote the note; note author name; note spaces removed; note text trimmed whitespace

Turbo saves a client note with the signed-in user as its author (uid) and removes spaces from both ends of the note text. The note gets the field job_status_changed_to only when a non-blank status is passed. The close notes pass none, so they leave it empty.
