<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_rate_sheets.go","sha256":"cde5c060302d088795b578a31b4581ea9bafeef8c0de0f682cff8db2e543ce8b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_rate_sheets.go

## Opening period has a fixed ID until its row is saved

Lines 22-24

Asked as: opening id; period id opening; first rate sheet id; opening row missing

A project whose rate sheet history has not been written is addressed as the period "opening" (openingRateSheetPeriodID). Turbo creates its row on the first change.

## Rate sheet reason is limited to 500 characters

Lines 26-26

Asked as: reason too long; reason character limit; how many characters reason; reason text too big

The reason for a rate sheet period change is capped at 500 characters (rateSheetReasonMaxLength). The limit counts characters, not bytes.

## Opening period is shown without a start date or end

Lines 37-48

Asked as: blank start date meaning; period no start date; why end date empty; latest period no end

A period with no start date (Starts empty) is the opening period; it covers all earlier time. The last period has no end date (Through empty), meaning it applies to the latest time.

## Rate sheet errors are reported with the message and field

Lines 65-83

Asked as: rate sheet error message; error next to field; refused change message; where error shows

When a rate sheet change is refused, Turbo shows the message and ties it to a field (rate_sheet, effective_date or reason) when the check concerns one. A refusal with no status counts as a bad request. An error that is not a rule refusal shows its own text.

## Who sees revise and set-first-sheet options

Lines 85-99

Asked as: cannot see revise button; revise option missing; set first sheet hidden; who can revise

The history view tells the user whether they can revise (the rate_sheet_revise claim) and whether they can set the first rate sheet. Both need an active user with the job claim. Revising also needs the rate_sheet_revise claim. Setting the first sheet is offered only when the project has no rate sheet.

## Proposals have no rate sheets

Lines 101-110

Asked as: proposal rate sheet; P job no sheets; proposals do not have rate sheets; rate sheet for proposal

Turbo refuses any rate sheet history request for a job whose number starts with P. The message is "Proposals do not have rate sheets." A project that does not exist gets "Project not found."

## Project not found when opening rate sheet history

Lines 101-105

Asked as: project not found history; job number missing; rate sheet project does not exist; history project not found

Turbo answers "Project not found." when the project in a rate sheet history request does not exist.

## Periods are listed with who created and last updated them

Lines 112-175

Asked as: who created period; who last updated period; creator name history; period editor name

Each period in the history shows the name of the person who created it and who last updated it (given name plus surname), its reason, its update time, and a summary of its rate sheet: code, name, revision, effective date and active flag.

## Opening period is implied when a project has a sheet but no history

Lines 137-145

Asked as: sheet but no history; history empty but sheet; opening shown without history; project has sheet no periods

If a project has a rate sheet but no history rows, Turbo shows one opening period (id "opening") on the project's current rate sheet. The opening period has no start date and covers all earlier time. Turbo writes its row only when the history is first changed.

## Each period ends the day before the next starts

Lines 146-162

Asked as: period end date wrong; through date; why period ends day before; end date of period

In a project's rate sheet history, each period's last day (through) is the day before the next period's start date. The latest period has no end date.

## Rate sheet history is fetched with permission flags and no access check

Lines 177-198

Asked as: can revise flag; can set first flag; history view permissions; history permission flags

Viewing a project's rate sheet history returns its periods plus two flags: can revise (can_revise) and can set first (can_set_first). The view itself does not refuse anyone; it only reports the flags. The can-set-first flag is true only for an active user with the job claim on a project that has no rate sheet.

## Rate sheet history changes need editing enabled

Lines 202-210

Asked as: cannot edit rate sheet; job editing switched off; rate sheet changes disabled; editing disabled error

Every add, change or delete of a project's rate sheet history is refused when job editing is switched off. The refusal is the error for disabled job editing (ErrJobsEditingDisabled).

## Changed history is returned after each change

Lines 202-241

Asked as: list after change; history refresh after save; failed change rollback; periods after edit

After a successful add, update or delete of a period, Turbo returns the project's rebuilt list of periods. The whole change runs in one transaction, so a refused change leaves the history unchanged.

## Changing history requires an active user with the job claim

Lines 215-221

Asked as: no job claim; not active user error; cannot add rate sheet period; job claim required

Adding, changing or deleting a rate sheet period needs an active user with the job claim. Otherwise Turbo says "An active user with the job claim is required."

## History change is validated and syncs the project's rate sheet

Lines 243-258

Asked as: project rate sheet not updated; latest period sets sheet; rate sheet did not sync; jobs rate sheet changed

Each history change runs in one transaction. Turbo then checks the whole history and sets the project's rate sheet (jobs.rate_sheet) to the latest period. If any step fails, the change is not saved.

## After any change the whole history is validated and the project's sheet synced

Lines 246-258

Asked as: history check after edit; validation after delete; project sheet mismatch; validate whole history

After each add, update or delete, Turbo checks the whole history against the pairs that existed before (via ValidateJobRateSheetHistory) and then sets the project's rate sheet (jobs.rate_sheet) to match the latest period.

## Pending rate sheet absorb blocks history changes

Lines 260-292

Asked as: absorb waiting to commit; undo absorb first; absorb blocks history; pending absorb error

Turbo refuses to change a project's rate sheet history while an uncommitted rate sheet absorb has moved that project's rate sheet or one of its periods. The message is "A rate sheet absorb that moved this project's rate sheet is waiting to be committed or undone. Finish it before changing the rate sheet history." The user must commit or undo the absorb first.

## Revise claim required to change history

Lines 294-299

Asked as: no revise permission; rate_sheet_revise missing; cannot change period claim; changing history needs claim

Changing or deleting a period, or adding a dated one, needs the rate_sheet_revise claim. Otherwise Turbo says "Changing a project's rate sheet history requires the rate_sheet_revise claim."

## Reason is required and limited to 500 characters

Lines 301-310

Asked as: empty reason error; explain why rate sheet changing; reason is required; reason blank refused

A period change needs a reason. An empty reason gets "Explain why the rate sheet is changing." A reason over 500 characters gets "Keep the reason to 500 characters or fewer." Turbo trims spaces from the reason first.

## Two periods cannot start on the same date

Lines 312-324

Asked as: two periods same date; another period starts on date; duplicate start date; date already used period

Turbo refuses a period whose start date is already used by another period of the project. The message is "Another period already starts on this date."

## Reason is trimmed before it is checked

Lines 335-335

Asked as: spaces in reason; reason only spaces; trailing space reason; reason counts as empty

Turbo trims spaces from the start and end of the reason before checking it. A reason of only spaces counts as empty.

## Adding a period requires choosing a rate sheet

Lines 336-338

Asked as: no rate sheet chosen; add period without sheet; select a rate sheet error; period sheet missing

Adding a period without a rate sheet gets "Select a rate sheet."

## Setting a project's first rate sheet needs no date or revise claim

Lines 345-353

Asked as: set first rate sheet; add sheet without date; first sheet no revise; first sheet reason optional

An add with no start date sets the first rate sheet of a project that has none. It needs only the job claim. The reason is optional here, but a reason over 500 characters is refused.

## Undated add refused when project already has a sheet

Lines 346-348

Asked as: enter start date error; date needed new sheet; undated add refused; new sheet start date required

If the project already has a rate sheet, an add without a start date gets "Enter the date the new rate sheet starts."

## Reason too long on undated first-sheet add is refused

Lines 349-351

Asked as: long reason first sheet; first sheet reason limit; undated first sheet reason

When setting a project's first rate sheet with no date, a reason is optional, but a reason over 500 characters is refused with "Keep the reason to 500 characters or fewer."

## Dated period needs a first sheet and a reason

Lines 355-364

Asked as: add dated period error; dated period needs sheet; set first sheet before dated; dated period no reason

Adding a dated period needs the rate_sheet_revise claim and a reason. If the project has no rate sheet, Turbo says "Set the project's first rate sheet before adding a dated one."

## Dated add writes the opening row first

Lines 372-380

Asked as: dated add opening row; first dated period saves opening; adding dated period no history

Adding a dated period to a project with no written history first creates the opening period row. Turbo then checks the date is free and saves the new period.

## New period records who made it

Lines 384-397

Asked as: created by user; who made period; last editor recorded; period stores user

A new period stores the user as both its creator (created_by) and its last editor (updated_by).

## Period lookup errors

Lines 401-418

Asked as: period not found reload; history changed reload; project has no rate sheet; period belongs another project

If the period id is not found or belongs to another project, Turbo says "Rate sheet period not found. The history may have changed; reload it." If the opening period has no row, Turbo says "This project has no rate sheet."

## Updating a period

Lines 423-459

Asked as: change period sheet; edit period start date; reprice opening period; update rate sheet period

Changing a period needs the rate_sheet_revise claim and a reason. The user can change its rate sheet and its start date. An empty rate sheet gets "Select a rate sheet." The change records the user as last editor. Changing the opening period's sheet reprices all time before the first dated period.

## Update requires revise claim and reason before looking up the period

Lines 430-440

Asked as: update period error order; reason checked before lookup; update needs claim first

Updating a period first requires the rate_sheet_revise claim, then a reason (1 to 500 characters), and only then looks up the period.

## Updating a period needs a rate sheet when one is sent

Lines 441-446

Asked as: empty sheet on update; period keeps current sheet; update rate sheet blank

When an update names a rate sheet but leaves it empty, Turbo refuses with "Select a rate sheet." If no rate sheet is sent, the period keeps its current one.

## Update changes start date only when one is sent, and checks duplicates

Lines 447-455

Asked as: change start date; date already used update; update start date only

An update sets a new start date only if one is sent, trimmed of spaces. It refuses a date that another period of the project already uses. The period's reason and last-updated-by are always set to the new reason and the current user.

## Deleting a period

Lines 462-480

Asked as: delete rate sheet period; remove period; previous period prices time; delete dated period

Deleting a period needs the rate_sheet_revise claim. The opening period cannot be deleted. Turbo says "The opening period cannot be deleted. Change its rate sheet instead." After a delete, the previous period prices that time.

## Opening period cannot be deleted

Lines 466-478

Asked as: delete opening period; cannot delete first period; opening period delete error

Deleting the opening period (no start date) is refused with "The opening period cannot be deleted. Change its rate sheet instead." Deleting a dated period needs the rate_sheet_revise claim; its time is then priced by the period before it.
