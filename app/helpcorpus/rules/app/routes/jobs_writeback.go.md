<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_writeback.go","sha256":"0a0b379e46b0d693c96325d81c7212616ac68c296d4f65d40f02ac7236d03720","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_writeback.go

## Export shows empty lists instead of missing ones

Lines 13-31

When a job has no categories, divisions or time allocations, the export gives an empty list or map. Empty JSON text such as "", "[]" or "{}" is treated as empty, and invalid JSON also gives an empty result.

## Who may export jobs for legacy writeback

Lines 239-259

Turbo accepts either a Bearer token that is a valid unexpired machine secret for legacy_writeback, or a signed-in user with the report claim. The token is trimmed of surrounding whitespace first. A caller with neither gets the message "unauthorized".

## Export requires an updatedAfter value

Lines 261-264

The jobs export refuses a request with no updatedAfter value. Turbo shows "updatedAfter is required".

## Job export fills empty values with defaults

Lines 267-303

For each exported job, missing dates, location, client PO, client reference number, authorizing document, outstanding balance date, parent, proposal, client, contact, manager and job owner come out as empty text. A missing proposal value comes out as 0 and missing time and materials as false.

## Job export includes categories, divisions and time allocations

Lines 286-288

Each job lists its category names and the codes of divisions with time allocations. It also gives a map of division code to allocated hours. Each of these is empty when the job has none.

## Job export lists proposal number, manager and branch details

Lines 289-313

Each job shows its proposal's job number, the legacy user ID (legacy_uid) of its manager and alternate manager, and the branch code. Manager, alternate manager and contact names are given name plus surname, with surrounding spaces trimmed.

## Export maps manager and author to legacy user IDs

Lines 296-299

The export identifies a job's manager, alternate manager and client note author by their legacy ID (legacy_uid) from admin profiles, not by their Turbo user ID. A client's business development lead is converted the same way. The ID is empty when the person has no admin profile.

## Export includes jobs updated on or after the date, not imported ones

Lines 314-314

The jobs export returns only jobs whose updated time is greater than or equal to updatedAfter, so a job updated exactly at that time is included. Jobs flagged as imported (_imported) are always excluded.

## Export clients are those of matched jobs, client or job owner

Lines 324-344

The clients list holds each distinct client that is either the client or the job owner of an exported job. Each client shows its business development lead as a legacy user ID, empty if none.

## Export contacts are those referenced by matched jobs

Lines 346-359

The contacts list holds each distinct client contact that an exported job references. A missing email or client comes out as empty text.

## Client notes ignore the updated date and imported flag

Lines 361-393

For each exported client, Turbo includes all of its notes, not only recent ones, and does not filter out imported jobs' notes. Notes are ordered by client and then newest first. Each note shows the author's legacy user ID and the linked job number, empty if none.

## Client notes are listed newest first per client

Lines 361-386

The export lists a client's notes in order of newest created first. Each note carries its job ID and job number when it has a job. A note's author is the legacy ID of the user.

## Rate roles are always all exported

Lines 395-401

The export always returns every rate role, ordered by name, whatever the updatedAfter date is.

## Rate sheets and entries limited to matched jobs

Lines 403-435

The export returns only rate sheets that exported jobs use, and the entries of those sheets. Each entry has a role, rate and overtime rate.

## Export trims names of contact and managers

Lines 437-468

The export trims leading and trailing spaces from the contact name, the manager name and the alternate manager name. A missing name gives an empty string.

## Job export field mapping to legacy names

Lines 437-482

The export gives each job's client as the client name and its branch as the branch code. It gives the job owner as the owner's name, and the proposal as the proposal's number. The ID fields (clientId, jobOwnerId, clientContactId, parentId, proposalId, rateSheetId) are also included.

## Clients without notes omit the notes list

Lines 485-510

A client with no notes is exported with no notes field.

## Export groups client notes under their client

Lines 485-499

In the jobs legacy export, Turbo attaches each client note to the client it belongs to. The note carries the author's legacy user ID (uid), the job ID and number, the job-not-applicable flag, and the job status it changed to.

## Export response omits notes list for clients without notes

Lines 501-510

In the jobs legacy export, each client gets its notes list only when it has notes. A client with no notes has no notes key in the JSON (Notes is omitted when empty).

## Export response is structured into six arrays

Lines 557-565

The jobs writeback export returns one JSON object with six lists: jobs, clients, clientContacts, rateRoles, rateSheets and rateSheetEntries. Related records are linked by ID, not nested inside each job.
