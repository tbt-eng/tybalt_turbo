<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_drafts.go","sha256":"e0ae022cb31e152ad7ec3b1adcff522bc9fe3696f2824f36a79579fe9167b641","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_drafts.go

## Proposed rates are a draft rate sheet owned by a project request

Lines 14-31

Asked as: proposed rates; propose new rates; draft rate sheet; rates for a project request

A project request can propose new rates instead of choosing an active rate sheet. The proposal is a draft: an inactive rate sheet whose request field (request) names the owning request. A rate sheet counts as a draft when its request field is not empty.

## Direct use of a draft rate sheet is refused

Lines 33-44

Asked as: cannot use draft rate sheet; rate sheet refused; rate_sheet_draft error; activated when job is created

Turbo refuses to use a draft rate sheet directly. The message is "this rate sheet is proposed in a project request; it is activated when the job is created". The error shows under the field that the UI names (code rate_sheet_draft).

## A request counts a sheet as its draft only if it owns it

Lines 55-71

Asked as: request has no draft; draft not showing; which rate sheet is mine; sheet owned by another request

A project request has a draft only when its rate sheet (rate_sheet) has the request's own ID in its request field. A request with no rate sheet, an unsaved request, or a sheet that is active or owned by another request has no draft. A missing sheet is not reported here; the normal rate sheet checks report it.

## Draft must have a reason

Lines 83-90

Asked as: reason for new rates; blank reason rates; why rates needed; invalid proposed rates

Turbo checks a draft of proposed rates when the request is submitted, approved and turned into a job. A draft with an empty or blank reason fails with "enter the reason why this project needs new rates" (invalid_proposed_rates).

## Proposed rates must differ from every active rate sheet

Lines 83-108

Asked as: duplicate rate sheet; rates same as existing sheet; identical rates; proposed rates already exist

Turbo compares the draft's rates with the active rate sheets, ignoring revisions that share the draft's name. If any active sheet has exactly the same rates, Turbo refuses with the duplicate rate sheet error under the proposed rates field. Turbo repeats the comparison at submission, approval and job creation, because an active sheet can change after the draft is saved.

## Draft must have a rate for every role

Lines 91-97

Asked as: missing rate for role; rates incomplete; no rate for role; role missing from proposal

A draft of proposed rates must have an entry for every role. If roles are missing, Turbo refuses with "the proposed rates have no rate for" followed by the missing roles, separated by commas.

## Activating a rate sheet refuses duplicate rates

Lines 112-131

Asked as: cannot activate rate sheet; activation duplicate error; same rates as active sheet; rate sheet activate fails

Before Turbo activates a rate sheet, it compares the stored entries with the other active sheets. Revisions with the same name are excluded. If an active sheet of another name has the same rates, Turbo refuses with the duplicate rate sheet error under the field active.

## Activation deactivates other revisions of the same name

Lines 112-136

Asked as: two active revisions; other revisions deactivated; make rate sheet active; new revision replaces old

When Turbo activates a rate sheet, it first deactivates the other revisions of that name. A name never has two active revisions. The caller checks permissions, completeness and newer revisions, and this code does not.

## Creating the job activates the request's draft rates

Lines 139-156

Asked as: job creation rates; rates become active after job; draft turns active; job made from request rates

When a job is created from an approved request that has a draft, Turbo validates the draft again. It then clears the draft's request field, sets its effective date (effective_date) to today's rate sheet date (RateSheetToday), and activates it. The sheet becomes an ordinary active sheet that keeps its reason. If the request uses an active sheet, Turbo does nothing.

## Deleting a project request deletes its draft rate sheet

Lines 158-183

Asked as: deleted request rate sheet; delete project request draft; draft gone after delete; request removed rates removed

When a project request is deleted by any route, Turbo also deletes its draft rate sheet, after the request itself is deleted. A sheet that job creation activated has its request field cleared, so it stays. If a draft is already gone, Turbo skips it.
