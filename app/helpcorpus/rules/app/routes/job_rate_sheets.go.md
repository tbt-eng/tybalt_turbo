<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_rate_sheets.go","sha256":"9048fa2db2417ebcea8a0a762088b08f72d2eae12db9c71d71d6fcc1128ae2d8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_rate_sheets.go

## Rate sheet reason limited to 500 characters

Lines 25-25

Asked as: reason too long; too many characters rate; rate note over limit; 500 characters reason

The reason for a rate sheet change can be at most 500 characters (rateSheetReasonMaxLength). Turbo counts characters, not bytes.

## Opening period has no start date

Lines 35-46

Asked as: first rate period blank start; why no start date; rates before any date; earliest period empty start

The opening period has an empty start date (Starts). It covers all time before the first dated period. A project with no history rows and no rate sheet shows no periods.

## Rate sheet routes report errors with their own status

Lines 63-81

Asked as: rate sheet error code; wrong error status rates; server error rate sheet; rate error message shows

A rate sheet error shows its own message and, when a field is named, the error code for that field. If the error has no status, Turbo treats it as a bad request. An unexpected error shows its message as a server error.

## Who may revise or set rate sheets

Lines 83-97

Asked as: who can change rates; can't edit rate sheets; permission to revise rates; why rates locked for me

Setting a first rate sheet needs an active user with the job claim. Changing the history (adding a dated period, editing or deleting a period) also needs the rate_sheet_revise claim. An inactive user has neither permission.

## Inactive users cannot revise or set rate sheets

Lines 86-90

Asked as: deactivated user rates; inactive account change rates; disabled user rate access

Turbo first checks that the user is active. An inactive user has neither the job permission nor the revise permission for rate sheets, whatever claims they hold.

## Revise needs the job claim and rate_sheet_revise claim

Lines 91-96

Asked as: can't revise rate history; rate history permission denied; job claim missing revise; rates need both claims

A user can revise a project's rate sheet history only with the job claim and the rate_sheet_revise claim. Without the job claim, Turbo does not check rate_sheet_revise, so both permissions are false.

## Proposals have no rate sheets

Lines 99-108

Asked as: proposal has no rates; P job rate sheet; rates on proposal refused; proposal rate history

Turbo refuses every rate sheet history action on a proposal (a job whose number starts with P). It shows 'Proposals do not have rate sheets.' If the project does not exist, it shows 'Project not found.'

## Missing project gives Project not found

Lines 99-103

Asked as: project not found rates; job doesn't exist rates; rate sheet missing project

When the project does not exist, a rate sheet request fails with "Project not found."

## Latest period has no end date

Lines 121-127

Asked as: latest rate no end date; blank through date; current rate open ended; why end date empty

Only the last period in the history has an empty end date (Through), meaning it is the latest. If the next period's start date cannot be parsed, the end date is empty too.

## Implied opening period when no history exists

Lines 135-143

Asked as: rate sheet no history rows; project shows opening period; default opening rate row; rates with no history

A project that has a rate sheet but no history rows shows one opening period (id 'opening') on its current rate sheet. The period has no start date and no end date.

## Period end dates come from the next period

Lines 144-160

Asked as: when does rate period end; period end day before next; through date calculation; end date of rate period

Turbo shows each period's last day (Through) as the day before the next period starts. The latest period has no end date. The opening period has an empty start date and covers all earlier time.

## Rate sheet history view and permission flags

Lines 174-195

Asked as: can revise flag; can set first flag; rate history permission flags; what rate actions can I do

Viewing a project's rate sheet history returns its periods plus two flags. can_revise is true for an active user with both the job claim and rate_sheet_revise. can_set_first is true for an active user with the job claim when the project has no rate sheet.

## Setting first rate sheet needs job claim only

Lines 189-193

Asked as: who can set first rate sheet; first rate sheet permission; set initial rate sheet rights

The history response says a user can set the first rate sheet (can_set_first) only when they hold the job claim and the project has no rate sheet yet.

## Rate sheet changes need job editing enabled

Lines 199-207

Asked as: job editing switched off; rates can't change editing disabled; editing turned off rates

Every add, change or delete of a rate sheet period is refused when job editing is switched off (utilities.ErrJobsEditingDisabled). The check runs before the permission check.

## Rate sheet changes return refreshed history

Lines 199-239

Asked as: rate history not refreshed; rates not updated after save; rate change rolled back; reload rate history after change

After a successful change, Turbo returns the project's rate sheet history rebuilt from the saved data. If any step fails, the whole change is rolled back.

## Job claim required for rate sheet changes

Lines 212-218

Asked as: need job claim change rates; active user job claim required; can't add rate period

Adding, changing or deleting a rate sheet period needs an active user with the job claim. Otherwise Turbo shows 'An active user with the job claim is required.'

## Rate sheet change is validated and synced as a whole

Lines 219-226

Asked as: rate change not saved; project rate doesn't match history; rate history validation failed; rate change transaction failed

After each change, Turbo runs in one transaction. It validates the whole history (hooks.ValidateJobRateSheetHistory) and updates the project's rate sheet to match the history (SyncJobRateSheet). If any step fails, the change is not saved.

## Revise claim required for dated changes

Lines 241-246

Asked as: rate_sheet_revise claim required; can't add dated rate; editing rate needs revise claim

Adding a dated period, editing any period, or deleting a period needs the rate_sheet_revise claim. Otherwise Turbo shows 'Changing a project's rate sheet history requires the rate_sheet_revise claim.'

## Reason is required and limited to 500 characters

Lines 248-257

Asked as: explain why rate changing; reason required rate change; empty reason rate; reason spaces trimmed

A dated add or an edit needs a reason. An empty reason gives 'Explain why the rate sheet is changing.' A reason over 500 characters gives 'Keep the reason to 500 characters or fewer.' Turbo trims spaces from the reason first and counts characters, not bytes.

## Two periods cannot start on the same date

Lines 259-271

Asked as: two periods same date; another period starts on this date; start date already used; rate date clashes

Turbo refuses a period start date that another period of the same project already uses. It shows 'Another period already starts on this date.' When editing, the period itself is ignored in the check.

## Adding a period requires a rate sheet

Lines 283-285

Asked as: select a rate sheet; no sheet chosen adding period; adding period without sheet

Adding a period without a rate sheet, or with an empty one, gives 'Select a rate sheet.'

## Set a project's first rate sheet without a date

Lines 302-312

Asked as: first rate sheet no date; set initial rates undated; optional reason first sheet; add rate without start date

Adding a period with no start date sets the first rate sheet of a project that has none. It needs only the job claim, not rate_sheet_revise. The reason is optional, but it still cannot exceed 500 characters. The period is saved with no start date.

## Undated add refused when project already has a sheet

Lines 303-305

Asked as: enter date new rate starts; no date when sheet exists; undated rate add existing project

If the project already has a rate sheet, adding a period without a start date gives 'Enter the date the new rate sheet starts.'

## Dated add needs a first sheet and a reason

Lines 314-331

Asked as: set first sheet before dated; dated rate no sheet; add dated rate fails; opening row written first

A dated period needs rate_sheet_revise. If the project has no rate sheet, Turbo shows 'Set the project's first rate sheet before adding a dated one.' A reason is required. Turbo writes the opening period row first if the history is not yet written, then checks the date is free.

## Dated change needs the project's first rate sheet

Lines 314-319

Asked as: dated change no existing sheet; first sheet before dated change; rate change needs sheet first

A dated period is refused with "Set the project's first rate sheet before adding a dated one." when the project has no rate sheet. This check runs after the revise claim check and before the reason check.

## Adding a dated period creates the opening period first

Lines 320-331

Asked as: opening row created on add; history not written yet; dated rate creates opening

Before a dated period is saved, Turbo writes the opening period row if the history is not yet written. Then it checks that the date is free.

## Finding a period to edit or delete

Lines 336-355

Asked as: this project has no rate sheet; rate period not found; reload rate history; history may have changed

Before an edit or delete, Turbo writes the opening row if the history is not yet written. The id 'opening' selects the period with no start date. If none exists, Turbo shows 'This project has no rate sheet.' A period id from another project, or an unknown one, gives 'Rate sheet period not found. The history may have changed; reload it.'

## Period not found messages

Lines 342-354

Asked as: opening id not found; period belongs to another project; unknown rate period id

The "opening" period id finds the project's period with an empty start date. If none exists, Turbo says "This project has no rate sheet." Any other id must belong to this project, otherwise Turbo says "Rate sheet period not found. The history may have changed; reload it."

## Editing a period's rate sheet, date or reason

Lines 367-395

Asked as: edit rate date or reason; change opening rates reprices; who changed rate record; reason needed even for sheet

Editing a period needs rate_sheet_revise and a reason, and a reason is always required, even if only the sheet changes. An empty rate sheet gives 'Select a rate sheet.' A new start date must not match another period's date. The edit records who changed it. Changing the opening period's sheet reprices all time before the first dated period.

## Update rejects empty rate sheet and duplicate date

Lines 378-390

Asked as: save edit empty rate sheet; edited rate date already used; update rate duplicate date

When an update sends an empty rate sheet, Turbo refuses it with "Select a rate sheet." A new start date is refused with "Another period already starts on this date." when another period of the project uses it. The period being edited is not counted as a duplicate.

## Deleting a rate sheet period

Lines 403-416

Asked as: delete rate period; remove dated rate no reason; can't delete opening period; delete rate without reason

Deleting a dated period needs rate_sheet_revise but no reason. Its time is then priced by the period before it. The opening period cannot be deleted: Turbo shows 'The opening period cannot be deleted. Change its rate sheet instead.'
