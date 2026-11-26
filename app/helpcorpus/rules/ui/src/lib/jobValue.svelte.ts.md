<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobValue.svelte.ts","sha256":"2de873bddb8a7d6a67fca064e5bb63557bc311949731c940673a6fa3c548b8b5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobValue.svelte.ts

## Approver queue versus job history

Lines 1-3

Asked as: approval queue wrong job; old job history showing; queue page loads history; late response changes job

On an approver's queue page, Turbo loads the approval queue. On a job page, it loads the value changes of that one job. Each job page keeps its own history, so a late response for an earlier job cannot change the job now shown.

## Values show in whole Canadian dollars

Lines 9-14

Asked as: value shows cents; project value dollars format; value showing $0; non-number value amount

Project values display in Canadian dollars with no cents. A value that is not a number shows as $0.

## Approved-by-me list covers 30 days, newest first

Lines 19-36

Asked as: my approvals last 30 days; approved by me list; approvals older than a month missing; my approved items newest

The approver's own approvals list (approved_by_me) shows approvals from the last 30 days. The server limits the list to 30 days and sorts it newest first. Each entry is either a request that became a new job (new_job) or a later value change (value_change).

## Date shown for an approved-by-me entry

Lines 19-30

Asked as: date on approved entry; which date approved item; when approval date shown; new job date approved list

For a new job, the date is when the job was created. For a later value change, the date is when the change was approved.

## Approved project values ordered by approval sequence

Lines 54-59

Asked as: approved values order; approval sequence list; which project value approved first; approved values newest first

A job's approved project values appear newest first, ordered by approval sequence (approval_sequence). Only values with status Approved are in this list.

## Current project value is marked by the server

Lines 60-62

Asked as: current project value missing; which approval is current; history missing current approval; current value not marked

The approval that set the current project value is the one the server marks as current. The history can omit it when the viewer may see only some approvals, so Turbo does not work it out from the list.

## A job has at most one pending or rejected change

Lines 63-64

Asked as: second pending change blocked; change already in progress; rejected change still pending; only one change per job

A project value change that is not Approved (pending or rejected) is treated as the one change in progress. A job has at most one.

## Page refreshes after a value change action

Lines 81-122

Asked as: page not refreshing after approve; counts stale after decision; error after approve reject delete; history not reloading

After a successful approve, reject or delete, Turbo reloads the history, the attention counts and the other data on the page that shows the value. If the action fails, Turbo shows the error and does not reload.

## Approve, reject and delete send the displayed revision

Lines 93-113

Asked as: stale version approve error; document hash mismatch; approve outdated revision; rejection reason not sent

Approving, rejecting or deleting a value change sends the version and PDF hash (document_hash) that the screen shows. A decision can change the version even when the PDF hash stays the same. A rejection also sends the job data shown (review_token) and the rejection reason.

## Change PDF opens in a new tab

Lines 124-145

Asked as: PDF won't open; change PDF blank tab; cannot load PDF error; pdf tab closes

Turbo opens the PDF of a value change in a new tab. If the PDF cannot load, the tab closes and the error shows on the page. The error text is "Cannot load this PDF." when the server response is not OK.

## Value change requires job details; blockers listed

Lines 148-168

Asked as: why can't I add value change; missing job details form; completion date required value change; job missing invoicing profile

Before a value change form is filled in, Turbo lists the project details the job lacks: an invoicing profile, division hours, a location and a branch. It also lists a project completion date, but only when the job is Active. The server checks the job in full when the change is saved.

## Missing authorizing document is not a blocker

Lines 148-152

Asked as: no authorizing document warning; missing document blocks change; PA project no document; authorizing document missing gap

Turbo does not list a missing authorizing document as a gap. The server treats the project as PA, and approval records it.
