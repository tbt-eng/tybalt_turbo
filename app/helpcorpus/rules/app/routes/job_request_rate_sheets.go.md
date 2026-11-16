<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"9aede7c1a844798308e1f3e716271024a2837bffbfd642577da6a6e1def89925","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates reason limited to 500 characters; optional in a draft

Lines 40-43

Asked as: reason too long; reason text over limit; justification too many characters; can't save long reason; description limit exceeded

The reason for proposed rates can have at most 500 characters. Turbo answers "the reason can have at most 500 characters" when it is longer. A draft request can be saved without a reason, but the code notes that submission requires one.

## Chosen rate sheet in a project request must be active

Lines 52-66

Asked as: rate sheet rejected; can't pick rate sheet; draft rate sheet selected; invalid rate sheet choice; rate sheet not accepted

When a project request chooses a rate sheet, Turbo refuses an unknown sheet with "Select a valid rate sheet." It refuses a draft rate sheet with "Select an active rate sheet." The request's own draft (its proposed rates) is allowed.

## Proposed rates require a name and a starting rate sheet

Lines 72-90

Asked as: proposed rates won't save; no name for new rates; starting sheet not chosen; blank rates name; missing base rate sheet

Saving proposed rates in a project request needs a name and a based_on rate sheet. The name and the reason are trimmed first. Turbo answers "name is required" or "choose the rate sheet to start from" when they are empty. Turbo also checks each rate entry with the rate sheet entry rules.

## Proposed rates must start from an active rate sheet

Lines 92-100

Asked as: can't base new rates on old sheet; archived sheet as starting point; source sheet not active; start rates from inactive sheet

Proposed rates must start from an existing, non-draft rate sheet, or Turbo refuses with "Start the proposed rates from an active rate sheet." An inactive source sheet is allowed only when it is already the source (based_on) of the request's own draft.

## Proposed rate with unknown role is refused

Lines 101-105

Asked as: role doesn't exist; wrong role on rate; invalid role in proposed rates; unknown role error

Turbo refuses proposed rates when an entry names a role that does not exist, with "A proposed rate has an unknown role."

## Proposed rate sheet name must be unique, ignoring letter case

Lines 107-122

Asked as: rate sheet name taken; name already used; same name as another sheet; rate sheet name exists

Turbo refuses a proposed rate sheet name that any other rate sheet revision already has, in any letter case. It answers "another rate sheet already has this name". The request's own draft is excluded from the check.

## Proposed rates that match an active rate sheet are refused

Lines 124-133

Asked as: rates same as existing sheet; duplicate rates; identical to active sheet; rates already exist elsewhere

Turbo refuses proposed rates that exactly match the rates of an active rate sheet. The error names that sheet (via DuplicateRateSheetError), so the requester can choose it instead.

## Draft is rewritten only when the proposal changed

Lines 135-167

Asked as: draft not updated; nothing changed on save; saved proposal ignored; why draft unchanged

A new proposal always counts as changed. For an existing draft, Turbo counts it as changed when the name, reason, based_on sheet or the rates differ. When nothing changed, Turbo only points the request at the draft and rewrites nothing.

## Saving proposed rates creates an inactive draft and clears approval

Lines 153-177

Asked as: proposed rates not active; approval cleared after edit; effective date changed on save; new draft is inactive

A new proposal is saved as a draft rate sheet (revision 0, active false) that belongs to the request. A change sets name, reason and based_on, clears the rates approver and approval date (a new approval is needed), and sets the effective date to today. Job creation sets the date the sheet becomes active.

## Draft rate entries are synced by role

Lines 179-213

Asked as: role removed from draft rates; rate change not saved; old role still in draft; draft rates not updating

When a proposal is saved, Turbo adds an entry for each new role and updates an entry only when its rate or overtime rate changed. It deletes draft entries whose role is not in the proposal. The request then points at the draft.

## Unused request draft is deleted

Lines 216-223

Asked as: old draft left behind; stale draft after changing sheet; unused proposal remains; delete old proposed rates

After the request is saved with a different rate sheet, Turbo deletes the draft that the request no longer uses. A nil draft does nothing.

## Job creator can rename proposed rates; approval stands

Lines 225-271

Asked as: rename rates when creating job; change proposed rate name; no proposed rates to rename; name taken while creating job

As a job is created, the job creator (job claim, checked by the caller) can give the proposed rates a new name. An empty or unchanged name does nothing. A request with no proposed rates is refused with "This request has no proposed rates to rename." A name that another rate sheet has, in any letter case, is refused with "Another rate sheet already has this name." Turbo keeps the old name in rate_sheet_name_changes, and the rates approval stands.

## Editor and approver see proposed rates beside source rates

Lines 273-314

Asked as: compare proposed and source rates; rate differences highlighted; approver view of changes; proposed rates next to source

For a request with proposed rates, Turbo shows the draft's rates approver, the source rate sheet, and each rate beside the source's rate and overtime rate. It marks a rate as changed when its rate or overtime rate differs, or when its role is not in the source. It adds nothing when the request uses an active sheet.

## Approver must state approval of the new rate sheet

Lines 326-341

Asked as: approve new rate sheet checkbox; can't approve rates; rates approval confirmation; approver must confirm rates

When a request has proposed rates, the approver must confirm "Approve the creation of the new rate sheet and its rates." Without the confirmation Turbo refuses. With it, Turbo records the approver and the time on the draft. A request that uses an active sheet needs no such confirmation. The caller checks that the person is the request's approver.

## New submission clears the proposed rates approval

Lines 343-355

Asked as: approval gone after resubmit; resubmitting removes approval; rates need approval again; submit again clears approval

Turbo clears the rates approver and approval date on the request's draft so a new submission needs a new approval. It does nothing when no approval exists or the request has no draft.

## Job creation needs approval of proposed rates by the request's approver

Lines 357-367

Asked as: can't create job rates not approved; job creation blocked by approval; rates not approved message; recall request for approval

Turbo refuses to create a job from a request whose draft rates lack approval by the request's approver. It also refuses when the approval date is missing. The message is "The proposed rates are not approved. Recall the request for a new approval." Requests without proposed rates are not checked.
