<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobValue.svelte.ts","sha256":"1e1c62f0d256705e495597116512c8ac360ca501eb7e63ee4cb60ffbd929f1ee","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobValue.svelte.ts

## Project value amounts display in whole Canadian dollars

Lines 9-14

Turbo shows a project value as Canadian dollars (CAD) with no cents, rounded to the whole dollar (valueCAD). A value that is not a number shows as $0.

## Value change dates show only the date part

Lines 16-17

Turbo shows a stored timestamp as its first 10 characters, such as 2026-10-14 (valueDate). An empty timestamp shows as blank.

## Approved value changes are listed newest first

Lines 35-40

Turbo lists a job's approved value changes by approval sequence (approval_sequence), highest first. Only changes with status Approved appear in this list.

## Current project value comes from the change the server marks

Lines 41-43

Turbo takes the current value from the approved change that the server marks as current. It does not work this out from the history, because the history can leave out approvals when the viewer may see only some of them.

## A job has at most one pending or rejected value change

Lines 44-45

Turbo treats the first change that is not Approved as the job's change in progress. A job has at most one such pending or rejected change.

## Value change history loads per job or as approver queue

Lines 47-60

A job page loads the value changes for that job only. The approver's queue instead loads the changes waiting for the signed-in approver's approval (queue=approval).

## Value change actions refresh the page's value and attention counts

Lines 62-66

After a successful action, Turbo reloads the value history, the attention counts and other page data that shows the value.

## Opening the proposal form clears the error

Lines 68-72

Opening the proposal form clears any earlier error. It opens for a new value change proposal when no change is passed, or for a correction when a change is passed.

## Approve, reject and delete send the displayed version of the change

Lines 74-103

To approve, reject or delete a value change, Turbo sends the version, the document hash and the review token of the change as displayed. A rejection also sends the reason. Turbo shows the error and reports failure if the server refuses. Delete uses a DELETE request; approve and reject use POST.

## Value change PDF opens in a new tab or downloads

Lines 105-124

Turbo opens a value change's PDF in a new tab, or saves it as a file when download is requested. The file name is the job number (or "project" if none), then "value", then the approval date (or the upload date if not approved). If the PDF cannot load, Turbo shows "Cannot load this PDF."

## Missing job details block a value change

Lines 127-147

Before a value change form is filled in, Turbo lists the job details that are missing: an invoicing profile, division hours (allocations), a location and a branch. It lists a project completion date only when the job is Active. The server checks the job in full when the change is saved.

## A missing authorizing document is not a blocker

Lines 127-131

Turbo does not list a missing authorizing document as a blocker. The server treats the project as PA, and approval records it.
