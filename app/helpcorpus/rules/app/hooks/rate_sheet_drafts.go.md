<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_drafts.go","sha256":"fd74ff5ed75d8100f631b47bed734f4cba7fe01c3417e1458ff86864767cdea9","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_drafts.go

## Proposed rate sheet is a draft owned by a project request

Lines 27-31

Asked as: rate sheet owned by request; proposed rates draft; rate sheet not yet active; which request owns this sheet

A rate sheet whose request field (request) names a project request is a draft. The project request owns the draft. Turbo treats any rate sheet with an empty request field as an ordinary sheet.

## Direct use of a draft rate sheet is refused

Lines 33-44

Asked as: cannot use proposed rate sheet; rate sheet refused directly; activated when job is created; rate_sheet_draft error

When a rate sheet draft is used directly, Turbo refuses with: "this rate sheet is proposed in a project request; it is activated when the job is created". The error shows under the field that the screen names (code rate_sheet_draft).

## A request's draft is found only if that request owns it

Lines 55-71

Asked as: request has no draft rates; draft not showing for request; sheet belongs to another request; rates from another request

Turbo treats the rate sheet of a project request as its draft only when the sheet's request field names that same request. A request with no rate sheet, an unsaved request, or a sheet that another request owns or that is an ordinary sheet has no draft. A missing sheet is not reported here; the normal rate sheet checks report it.

## Draft needs a reason

Lines 87-90

Asked as: reason field empty; why new rates needed; cannot submit project request; blank reason for rates

When a project request with proposed rates is submitted, approved or turned into a job, Turbo refuses a draft whose reason is empty or only spaces. The message is "enter the reason why this project needs new rates".

## Draft needs a rate for every role

Lines 91-97

Asked as: missing role rate; proposed rates incomplete; no rate for role; rate missing for a role

Turbo refuses proposed rates that lack a role's entry. The message is "the proposed rates have no rate for " followed by the missing roles, separated by commas. The check runs at submission, approval and job creation.

## Proposed rates must differ from every active rate sheet

Lines 98-109

Asked as: rates same as active sheet; duplicate proposed rates; rates already exist; proposed rates identical to current

Turbo refuses proposed rates when an active rate sheet has exactly the same rates. It shows the duplicate-rate-sheet error under proposed rates (proposed_rates). The draft's own code is excluded from the comparison. Turbo repeats the comparison at each step, because an active sheet can change after the draft is saved.

## Activating a rate sheet checks duplicates and replaces other revisions

Lines 117-137

Asked as: cannot activate rate sheet; active duplicate rate sheet error; old revisions still active; activate replaces previous revision

When Turbo activates a rate sheet, it compares the stored entries with the active sheets of other codes. If the rates are exactly the same, it refuses with the duplicate error under the field active. Otherwise it deactivates the other revisions of the same code, then sets the sheet active, so a code never has two active revisions.

## Job creation activates the request's draft rates

Lines 139-166

Asked as: job creation proposed rates; draft turns active on job; rates activated when job made; effective date set on job

When a job holder creates a job from an approved project request that has a draft, Turbo activates the draft in the same transaction. It sets the draft's effective date to today (RateSheetToday) and clears its request field. The sheet then keeps its reason as an ordinary active sheet. A request that uses an active sheet gets no activation.

## Job creation validates the draft again

Lines 148-152

Asked as: job fails draft checks; rates rejected when creating job; duplicate error at job creation; draft revalidated for job

Before activating the draft, Turbo runs the draft checks again: reason, a rate for every role, and no duplicate of an active sheet. A duplicate is reported under proposed rates.

## Older draft revision cannot replace a newer one

Lines 153-162

Asked as: older rate sheet revision; newer revision exists; old proposed rates won't activate; outdated draft cannot activate

If a newer revision of the draft's code exists, Turbo refuses to activate it. The message is "a newer revision of the proposed rate sheet exists". This is a safeguard, because nothing can revise a draft.

## Deleting a project request deletes its draft rate sheet

Lines 168-193

Asked as: deleted request draft gone; project request deletion removes rates; rate sheet disappeared after delete; draft removed with request

Deleting a project request also deletes the draft rate sheets that it owns, whoever deletes it. Turbo deletes the request first, then the drafts. A sheet already gone is skipped. Job creation clears the draft's request first, so an activated sheet stays after the request is deleted.
