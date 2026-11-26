<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobValue.svelte.ts","sha256":"fa048f0546a896e645e0f982eb346cca0ee9f4dfb9137c9afb4742d6b1e97d8f","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobValue.svelte.ts

## Dollar values show in whole Canadian dollars

Lines 9-14

Asked as: dollar amount without cents; project value in CAD; value shows as zero; currency format of job value; whole dollars only

Turbo shows project values as Canadian dollars with no cents. An amount that is not a number shows as zero.

## Approved-by-me list covers the last 30 days

Lines 19-36

Asked as: my approvals last 30 days; approvals I made list; which jobs I approved; approval history older than a month; recent approvals newest first

The approver's list of their own approvals covers the last 30 days, newest first. Each entry is a request that became a job, or a later value change. The date is the job's creation date for a new job and the approval date for a value change.

## Value change history shows approved values newest first

Lines 54-59

Asked as: old job values order; previous values of a job; value change history order; list of past project values; approved values sorted newest

On a job, Turbo lists approved project values (status "Approved") with the highest approval sequence first. The approval sequence orders a job's values.

## Current project value comes from the approval the server marks

Lines 60-62

Asked as: wrong current job value; which value is current; job shows old value; current value missing from history; why total differs from list

Turbo shows as the current value the approved change that the server marks as current. The history may leave out some approvals when the viewer may see only some of them, so Turbo does not work the current value out from the list.

## A job has at most one pending or rejected value change

Lines 63-64

Asked as: two value changes on job; cannot add second change; change already in progress; rejected value still blocking; only one pending change

Turbo treats any value change that is not Approved (pending or rejected) as the change in progress. It expects a job to have at most one.

## Approver queue lists value changes waiting for approval

Lines 66-73

Asked as: approval queue empty; changes waiting for me; my approval inbox; value changes to approve; job page shows no changes

In the approver's queue view, Turbo loads the changes that wait for the signed-in approver. On a job page, it loads that job's value change history instead.

## Page reloads after a value change action

Lines 81-122

Asked as: page not updating after approve; value history not refreshing; screen stale after reject; counts did not change; nothing reloads after delete

After a successful approve, reject or delete, Turbo reloads the value history, the attention counts and the page data that shows the value. If the action fails, Turbo shows the error and returns failure.

## Approve, reject or delete sends the displayed revision

Lines 93-113

Asked as: approved wrong version; outdated value change approved; stale version on approve; reject sends old data; changed after I opened it

When a person approves, rejects or deletes a value change, Turbo sends the version and PDF hash (document_hash) that are on screen. A rejection also sends the job data shown (review_token) and the rejection reason. A decision can change the version even when the PDF hash stays the same.

## Value change PDF opens in a new tab

Lines 124-145

Asked as: PDF will not open; value change document blank; new tab empty for PDF; cannot view attachment; PDF error on job

Turbo opens the PDF of a value change in a new tab. If the PDF cannot load, Turbo closes the tab and shows an error on the job page, such as "Cannot load this PDF."

## Job details that block a value change

Lines 148-167

Asked as: why can't I add value change; value change form won't open; missing hours blocks change; job needs location to change value; completion date required before change

Before the value change form opens, Turbo lists missing job details: division hours (no allocations), a location, and a branch. It also asks for a project completion date when the job's status is Active. A missing authorizing document or invoicing profile is not listed. The server checks the job in full when the change is saved.
