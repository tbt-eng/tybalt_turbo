<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_rate_sheets.go","sha256":"dd3199594614dda0e136226359ae3341e1c949c7dbab2d5d62744d9c2ca2c7a2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_rate_sheets.go

## Proposed rates need a reason of at most 500 characters

Lines 43-87

Asked as: reason required rates; reason too long; why need reason; reason field error

Every save of proposed rates needs a reason (proposed_rates.reason). Turbo trims it first. An empty reason gives "reason is required". A reason over 500 characters gives "the reason can have at most 500 characters".

## Chosen rate sheet in a project request must be active

Lines 55-69

Asked as: can't pick draft sheet; select valid rate sheet; rate sheet not accepted; draft sheet another request

A project request refuses a chosen rate sheet that is a draft, unless the draft is the request's own. An unknown sheet gives "Select a valid rate sheet." A draft that belongs to another request gives "Select an active rate sheet."

## Proposed rates need a starting rate sheet

Lines 88-94

Asked as: no starting rate sheet; based on sheet missing; choose sheet to start; which rate sheet to copy

Every save of proposed rates needs the rate sheet to start from (based_on). If it is empty, Turbo says "choose the rate sheet to start from". Turbo also validates the entries with the rate sheet entry rules before saving.

## Proposed rates must start from an active rate sheet

Lines 96-104

Asked as: start from inactive sheet; based on deactivated sheet; old rate sheet not allowed; sheet was deactivated

The starting sheet must exist and must not be a draft, or Turbo says "Start the proposed rates from an active rate sheet." An inactive sheet is accepted only when the request's own draft already started from it. So a draft keeps its source after that sheet is deactivated.

## Proposed rate with unknown role is refused

Lines 105-109

Asked as: unknown role proposed rate; role not found rates; wrong role in rates

Each entry of the proposed rates must use an existing role. Otherwise Turbo says "A proposed rate has an unknown role."

## Proposed rate sheet name must not be held by another sheet

Lines 111-123

Asked as: rate sheet name taken; duplicate sheet name; name already used sheet; proposed name exists

Turbo refuses a proposed name that another rate sheet already holds, in any letter case. The request's own draft is not counted as another sheet. The error is shown on the proposed_rates name field. The name is optional.

## Proposed rates that an active sheet already has are refused

Lines 125-135

Asked as: same rates as existing sheet; rates already exist; sheet already has these rates; duplicate rate sheet

Turbo compares the proposed rates with all active rate sheets, except for the request's own draft. If an active sheet has exactly the same rates, Turbo refuses the save and names that sheet, so the requester can choose it instead.

## Saving unchanged proposed rates changes nothing

Lines 137-148

Asked as: saved but nothing changed; no changes to draft; draft not updated; save did nothing

A new proposal always counts as changed. For an existing draft, Turbo compares the name, reason, starting sheet and rates. If none differ, the draft is not rewritten and the request is only pointed at it.

## Saving proposed rates writes a draft rate sheet

Lines 154-216

Asked as: where does my draft go; draft rate sheet created; proposed rates saved as sheet; new sheet inactive draft

A new draft rate sheet is inactive, at revision 0, and belongs to the request. Its effective date is set to today, and job creation sets the real date later. Changed rates clear the draft's rates approval, so they need a new approval. Turbo adds, updates or deletes the draft's entries to match the saved rates, and the request then points at the draft.

## Request draft is deleted when no longer used

Lines 218-225

Asked as: old draft still there; draft removed after changing sheet; leftover draft sheet; unused draft cleanup

After the request is saved with a different rate sheet, Turbo deletes the draft the request owned. If there is no draft, it does nothing.

## Procurement can rename proposed rates at job creation

Lines 227-271

Asked as: rename proposed rates at job; change sheet name when creating job; new name for proposed rates; job holder renames rates

When a job is created from a request, the job holder can give the proposed rates a new name. The rates approval stands, because the approver approved the rates and not the name. Turbo keeps the old name in the rate sheet name changes. An empty or unchanged name does nothing.

## Rename refused without proposed rates or when name is held

Lines 242-256

Asked as: can't rename rates; no proposed rates to rename; rename name already used; rename refused

If the request has no proposed rates, Turbo says "This request has no proposed rates to rename." If another rate sheet holds the new name, Turbo refuses with the message for a held name.

## Approver sees the proposal with changes marked

Lines 273-332

Asked as: what changed in proposal; approver compare rates; highlighted rate changes; show changed rates

For a request with proposed rates, Turbo adds the following to its related records: the draft's reason and the source sheet. It also adds the rates approver, each rate beside the source rate and a changed flag. A role with no source rate counts as changed. These fields appear only when the request uses a proposed draft.

## Missing roles are listed for the proposed rates

Lines 282-363

Asked as: roles missing from proposal; can't sign authorization; authorization not signable; missing role rates

Turbo lists the roles that the draft has no rate for (rate_sheet_missing_roles) and, for each one, the role name and the source sheet's rates if it has them. The list is sorted by role name. The printed authorization can be signed only when the list is empty. These fields are for display only.

## Approver must state approval of the new rates

Lines 365-390

Asked as: approval checkbox new rates; approver didn't approve rates; approve new rate sheet; rates approval needed

The approver of a request with proposed rates must state that they approve the creation of the new rate sheet and its rates. Without it, Turbo says "Approve the creation of the new rate sheet and its rates." Turbo records who approved and when on the draft. It does nothing if the request uses an active sheet.

## New submission clears the rates approval

Lines 392-404

Asked as: approval reset on resubmit; submit again approval gone; approval cleared after resubmitting; resubmit needs approval

Turbo clears the draft's approver and approval time, so that a new submission needs a new approval. If the draft has neither, it does nothing.

## Job creation needs approval of proposed rates

Lines 406-417

Asked as: can't create job rates not approved; recall request approval; proposed rates not approved; job creation blocked rates

Turbo refuses to create the job unless the draft's rates were approved by the request's own approver. The message is "The proposed rates are not approved. Recall the request for a new approval."

## Compare proposed rates with active rate sheets

Lines 419-460

Asked as: check rates before saving; compare unsaved rates; similar rate sheets; rates match existing sheet

Anyone with the time claim can compare unsaved proposed rates while project requests are enabled. Turbo returns the exact matches, which a save would refuse, and a few close sheets. The close sheets leave out the starting sheet and do not count roles that were left out of the entries. The compare changes nothing.

## Find rates: overtime rate is optional

Lines 462-469

Asked as: overtime rate search optional; search rates without overtime; find rates overtime blank; overtime not needed search

In a search for rate sheets by known rates (find_rates), an entry's overtime rate (overtime_rate) is optional. Without it, a rate sheet matches on the entry's rate alone. With it, the sheet must have the same overtime rate too.

## Find rates entry validation

Lines 483-535

Asked as: rate search error; enter at least one rate; search rate whole number; overtime rate decimal places

A rate search needs at least one entry ("enter at least one rate"). Each role must be given, known and used once. The rate must be given and be a whole number of at least 1. An overtime rate needs the role's rate, must be at least 1 and can have at most two decimal places.

## Job count of a sheet counts current and past jobs

Lines 537-573

Asked as: how many jobs use sheet; job count rate sheet; sheet usage count; jobs on rate sheet

A rate sheet's job count counts the jobs that use it now and the jobs that used it in their rate sheet history. Each job is counted once. A sheet with no jobs shows 0.

## Find rate sheets by known rates

Lines 575-638

Asked as: search rate sheets by rates; find sheet matching my rates; lookup rate sheet by rate; which sheet has these rates

Anyone with the time claim can search active rate sheets by the rates they know, while project requests are enabled. A sheet matches if it has each role with the same rate, and the same overtime rate when one is given. Roles that are left out are ignored. Each match shows its code, name, revision, effective date and the number of jobs using it, now or in history. The search does not return drafts, and it changes nothing.
