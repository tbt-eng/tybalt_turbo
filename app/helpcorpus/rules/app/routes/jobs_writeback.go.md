<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_writeback.go","sha256":"dd214e64290a5493a1f064dd6347b3263786995b021b56c704309da99ea2d1e0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_writeback.go

## Export empty arrays and maps give empty values

Lines 13-31

Asked as: empty list in export; no categories exported; blank divisions in export; malformed json export

When a job has no categories, no divisions or no time allocations, or the stored JSON is empty, the export gives an empty list or map. Malformed JSON also gives an empty value, with no error.

## Export response has six arrays

Lines 33-41

Asked as: what sections does export return; export response structure; jobs export keys; ids in exported job

The jobs export returns one JSON object with jobs, clients, clientContacts, rateRoles, rateSheets and rateSheetEntries. Jobs refer to clients, contacts, parent, proposal and rate sheet by ID.

## Rate sheet name in the export is the sheet's code

Lines 49-59

Asked as: rate sheet name wrong legacy; sheet description missing legacy; rate sheet shows number

In the jobs export for the legacy system, each rate sheet's name (Name) carries the sheet's code, such as 2026-11-00004. The sheet's own optional description is not sent to the legacy system.

## Export rate sheet entries carry rate and overtime rate

Lines 61-68

Asked as: rate sheet entries missing rate; overtime rate not exported; role rate in export

Each exported rate sheet entry (rateSheetEntries) gives its role, its rate sheet, its rate and its overtime rate (overtimeRate). The entries come from the rate sheets of the exported jobs.

## Export keeps client note details and links

Lines 70-81

Asked as: client note author missing; note job number not showing; job not applicable flag export

Each exported client note gives its text, its created and updated times, and its author's legacy ID (uid). It also gives the job ID and job number when a job is linked, a flag for job not applicable (jobNotApplicable), and the job status that the note changed the job to. Turbo leaves out the job ID, job number and status change when they are empty.

## Export client carries the business development lead's legacy ID

Lines 83-89

Asked as: business development lead export; client lead missing legacy id; client owner not in export

Each exported client gives its name and its business development lead as the lead's legacy ID. Turbo leaves out the lead when the client has none or the lead has no legacy ID.

## Export contact gives name, optional email and client

Lines 91-98

Asked as: contact email missing export; client contact details export; contact surname not exported

Each exported contact gives its surname, given name and client ID. Turbo leaves out the email when it is empty.

## Job export omits empty optional fields

Lines 196-238

Asked as: missing fields in jobs export; blank optional fields removed; why field not in export

The jobs export leaves out optional fields that are empty: alternate manager, proposal, dates, location, authorizing document, client PO, client reference number, categories, divisions and the ID references. Proposal value, outstanding balance, the fn agreement flag and time and materials always appear, as does jobTimeAllocations.

## Export has no time allocations key omitted for jobs

Lines 215-229

Asked as: jobTimeAllocations missing from export; empty allocations object; time allocations not showing

The jobs export always includes the time allocations of a job (jobTimeAllocations), even when the job has none. It then shows an empty object. Turbo also always includes the proposal value, outstanding balance, fnAgreement and timeAndMaterials fields of a job.

## Who can call the jobs export for the legacy system

Lines 242-262

Asked as: export unauthorized error; bearer token rejected; machine secret legacy access; who may call export

The jobs export accepts a Bearer token that matches an unexpired legacy_writeback machine secret. Trailing whitespace in the token is trimmed. If there is no valid token, it accepts a signed-in user who holds the report claim. Anyone else gets the message "unauthorized".

## Jobs export requires an updatedAfter value

Lines 264-267

Asked as: updatedAfter is required error; export missing date parameter; export asks for date

The jobs export refuses a request with an empty updatedAfter value. Turbo shows the message "updatedAfter is required".

## Jobs export returns jobs updated since the given time, not imported ones

Lines 270-318

Asked as: imported jobs not in export; jobs changed since date; export skips imported jobs

The export lists jobs whose updated time is on or after updatedAfter (the equal value is included). It leaves out jobs that were imported (_imported is not 0).

## Jobs export fills empty job fields with defaults

Lines 271-306

Asked as: job fields show zero; empty proposal value export; blank dates in export

For each job, missing dates, location, document, client PO, reference number and names come out as empty text. A missing proposal value comes out as 0 and missing time and materials comes out as 0. Outstanding balance is passed as stored.

## Job categories and divisions are listed in the export

Lines 289-291

Asked as: job categories in export; hours per division export; allocated hours map

Each exported job lists its category names and the codes of the divisions that have time allocations. It also gives a map of division code to allocated hours (jobTimeAllocations). An empty list or map is returned when there are none.

## Job export identifies people by legacy ID

Lines 299-316

Asked as: manager legacy id export; branch code in export; proposal job number export

The manager and alternate manager are exported with the legacy_uid from their admin profile, and a display name made from given name and surname. Display names are trimmed of extra spaces. The branch is exported as its branch code, and the proposal as the proposal's job number.

## Clients in the export come from job client and job owner

Lines 329-340

Asked as: which clients are exported; duplicate clients in export; client list export

The clients list holds each distinct client that is the client or the job owner of an exported job. Each has its business development lead as a legacy_uid, or empty text if none.

## Contacts in the export are those used by exported jobs

Lines 350-355

Asked as: which contacts are exported; contacts used by jobs; contact list in export

The contacts list holds each distinct client contact that an exported job uses. It gives surname, given name, email (empty if none) and client.

## Client notes export includes all notes of exported clients

Lines 364-389

Asked as: all client notes exported; old notes missing from export; note author legacy uid

For every client in the export, Turbo sends all its notes, not only recent ones, and the imported flag is not applied. Notes are listed newest first within each client. Each note shows the author's legacy_uid, the job ID and number if any, the job-not-applicable flag, and the job status it changed to.

## Rate roles are always exported in full

Lines 398-399

Asked as: rate roles always exported; role list order export; all roles in export

The export always lists all rate roles, sorted by name, whatever updatedAfter is.

## Rate sheets and entries are exported only for exported jobs

Lines 406-438

Asked as: rate sheets missing from export; sheet code instead of name; effective date revision export

The export lists the rate sheets that exported jobs use, and all entries of those sheets (role, rate, overtime rate). The sheet's name in the export is its code, such as 2026-11-00004, and not its own name. Effective date, revision and active flag are included.

## Job export reports the branch by code

Lines 440-486

Asked as: branch shown as code; client name blank export; missing linked names export

The jobs export gives each job's branch as the branch code (Branch). The client, job owner and contact names come from linked records. A missing link gives an empty value.

## Export trims display names of contact and managers

Lines 449-457

Asked as: extra spaces in names export; contact name whitespace; trailing space manager name

The jobs export trims leading and trailing spaces from the client contact name, the manager's display name and the alternate manager's display name. The contact name is the given name plus the surname.

## Client notes are attached to their client, newest first

Lines 488-513

Asked as: note order in client; client with no notes field; notes sorted newest first

In the export, each client carries its notes, grouped by client. Notes within each client run from newest to oldest by created date. A client with no notes has no notes field.

## Export response returns the jobs and related records at once

Lines 560-570

Asked as: export paging; export filtering not applied; one response all records

A successful call returns one response with the jobs, clients, client contacts, rate roles, rate sheets and rate sheet entries. Turbo does not filter or page the result.
