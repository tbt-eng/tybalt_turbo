<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_rate_sheets.go","sha256":"4906045336023b361f45b07993dfc9c3da90b128a1fe7b6366c460ee830c9fcc","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_rate_sheets.go

## Opening period id stands for an unwritten history

Lines 21-23

Asked as: why is the id opening; unwritten history id; first period id opening; history not saved yet

A project whose rate sheet history has not been written yet has an opening period with the id "opening" (openingRateSheetPeriodID). Turbo creates its row on the first change to the history.

## History view shows empty start date for opening period

Lines 35-46

Asked as: blank start date history; no start date first period; latest period no end; why start date empty

In a project's rate sheet history, the opening period (the first one) has an empty start (starts) and covers all earlier time. The latest period has an empty end (through).

## Who can see history and what the flags mean

Lines 83-97

Asked as: no revise option; can't see revise button; rate sheet permissions missing; can I edit rate history

Viewing a project's rate sheet history returns its periods. It also returns can_revise, which is true for an active user with both the job claim and the rate_sheet_revise claim. It returns can_set_first, which is true for an active user with the job claim when the project has no rate sheet yet. A user without the job claim gets both flags false.

## Proposals have no rate sheet history

Lines 99-108

Asked as: proposal rate sheet error; P job rate sheet; can't add rate sheet proposal; project not found rate

Every rate sheet history action (view, add, change, delete) on a proposal (a job whose number starts with P) is refused with "Proposals do not have rate sheets." An unknown project gives "Project not found."

## History shows names of who created and updated a period

Lines 110-119

Asked as: who created the period; name missing on period; last updated by blank; who changed this rate

Each period in the history shows the name of its creator and last updater, taken from their profile (given name and surname). The name is empty when there is no profile.

## History view for a project with a sheet but no history

Lines 134-143

Asked as: project has sheet but no history; existing rate sheet no periods; old project no history; where is the history

A project that has a rate sheet but no saved history shows one implied opening period on its current rate sheet. The opening period has the ID "opening" and an empty start date.

## How period end dates are shown

Lines 144-160

Asked as: when does period end; through date wrong; rate period last day; end date of period

Periods are listed in order. Each period's last day (through) is the day before the next period starts. The latest period has no end date. The opening period starts empty and covers all earlier time.

## Rate sheet history view lists periods and permission flags

Lines 174-195

Asked as: show rate sheet history; list rate periods; view rate history page; see project rate periods

The history view for a project returns its periods plus two flags: can_revise and can_set_first. The view itself checks no permission beyond loading the project, and returns 'Project not found.' when the project does not exist.

## Rate sheet history change runs as one transaction

Lines 197-246

Asked as: change failed nothing saved; partial rate sheet update; rate change rolled back; save half done rates

Each add, change or delete of a period runs in one transaction, so a refused change leaves the history unchanged. After success Turbo returns the project's updated periods.

## Rate sheet changes need job editing enabled

Lines 199-207

Asked as: can't edit rate sheet; rate sheet locked; job editing turned off; rates read only

Any add, change or delete of a rate sheet period is refused when job editing is switched off (ErrJobsEditingDisabled). The change runs in a single transaction.

## Changing history needs an active user with job claim

Lines 212-218

Asked as: can't add rate period; not allowed change rate sheet; inactive user rate sheet; no job access rates

An add, change or delete of a rate sheet period is refused with "An active user with the job claim is required." unless the user is active and has the job claim.

## History is validated and job sheet synced after each change

Lines 226-245

Asked as: rate sheet sync failed; history validation error; change reverted after save; rate change undone automatically

After each add, change or delete, Turbo validates the whole history (ValidateJobRateSheetHistory) and updates the project's rate sheet (SyncJobRateSheet). If either fails, the whole change is rolled back. On success Turbo returns the new list of periods.

## Pending rate sheet absorb blocks history changes

Lines 248-280

Asked as: absorb waiting to commit; rate sheet absorb blocking; can't change rates after absorb; undo absorb first

Turbo refuses to change a project's rate sheet history while an uncommitted rate sheet absorb has moved that project's rate sheet or one of its periods. The message is "A rate sheet absorb that moved this project's rate sheet is waiting to be committed or undone. Finish it before changing the rate sheet history." The user commits or undoes the absorb first.

## Dated period changes need rate_sheet_revise claim

Lines 282-287

Asked as: need revise permission; can't add dated period; missing rate sheet revise; dated period permission denied

Adding a dated period, changing a period and deleting a period each need the rate_sheet_revise claim. Without it Turbo refuses with "Changing a project's rate sheet history requires the rate_sheet_revise claim."

## Reason is required and limited to 500 characters

Lines 289-298

Asked as: reason required rate change; explain why field; reason too long; reason 500 characters limit

A reason is required when Turbo adds a dated period or changes a period. It is trimmed first. An empty reason gives "Explain why the rate sheet is changing." A reason over 500 characters gives "Keep the reason to 500 characters or fewer."

## Two periods cannot start on the same date

Lines 300-312

Asked as: two periods same date; date already used rate; duplicate start date period; start date taken

Turbo refuses a period whose start date another period of the project already uses, with "Another period already starts on this date." When a period is changed, it is not compared with itself.

## Reason is trimmed before it is checked

Lines 323-323

Asked as: reason only spaces; blank reason spaces; whitespace reason error; spaces count as empty reason

When a person adds or changes a period, Turbo removes spaces from both ends of the reason before it checks it. A reason of only spaces counts as empty.

## Adding a period requires a rate sheet

Lines 324-326

Asked as: no rate sheet selected; add period without sheet; select a rate sheet error; new period no sheet

Adding a period without a rate sheet is refused with "Select a rate sheet."

## Period start date is trimmed

Lines 327-331

Asked as: date spaces error; blank effective date; start date only spaces; effective date treated first

When a person adds a period, Turbo trims spaces around the effective date. A date of only spaces counts as no date, so the period is treated as the first rate sheet.

## New period records who created and updated it

Lines 337-341

Asked as: who added period; created by updated by; who made this change; who last edited rate

A new rate sheet period stores the person who adds it as both its creator (created_by) and its last updater (updated_by). A change to a period sets updated_by to the person who changes it.

## Setting a project's first rate sheet needs only job claim

Lines 343-353

Asked as: set first rate sheet; no sheet yet add; first rate sheet permission; add sheet with no date

Adding a period with no start date sets the first rate sheet of a project that has none. It needs only the job claim, not rate_sheet_revise. The reason is optional but is limited to 500 characters. The period is saved as the opening period with an empty date.

## Undated period refused when project already has a sheet

Lines 343-346

Asked as: enter start date error; new sheet needs date; missing effective date; undated period refused

If the project already has a rate sheet, adding a period without a start date is refused with "Enter the date the new rate sheet starts."

## Dated period needs a first rate sheet and writes the opening row

Lines 355-372

Asked as: dated period no first sheet; set first sheet before dated; add dated period no sheet; opening row created

A dated period is refused with "Set the project's first rate sheet before adding a dated one." when the project has no rate sheet. Otherwise Turbo first creates the opening period row if the history is not yet written. Then it checks the date is free and saves the new period with the reason.

## Finding a period to change or delete

Lines 377-396

Asked as: period not found reload; project has no rate sheet; history changed reload; can't find period to change

Turbo creates the opening row if needed before finding a period. The ID "opening" means the period with no start date; if none exists the message is "This project has no rate sheet." A period ID that does not exist or belongs to another project gives "Rate sheet period not found. The history may have changed; reload it."

## Changing a period's rate sheet or start date

Lines 408-436

Asked as: change period rate sheet; reprice before first period; edit period start date; change start date of period

Changing a period needs rate_sheet_revise and a valid reason. A blank rate sheet is refused with "Select a rate sheet." A new start date must be free. Turbo stores the reason and records the user as updater. Changing the opening period's sheet reprices all time before the first dated period.

## Changing an existing period needs a reason and revise claim before lookup

Lines 408-431

Asked as: change period empty rate sheet; select a rate sheet error; empty start date change; change period without reason

When changing a period, Turbo checks the rate_sheet_revise claim first, then the reason, then finds the period. Sending an empty rate sheet refuses with 'Select a rate sheet.' and an empty start date is allowed only if no other period starts that day.

## Deleting a period

Lines 440-458

Asked as: delete rate sheet period; cannot delete opening period; remove dated period; delete rate period

Deleting a dated period needs rate_sheet_revise, and no reason is required. Its time is then priced by the period before it. The opening period cannot be deleted; Turbo refuses with "The opening period cannot be deleted. Change its rate sheet instead."
