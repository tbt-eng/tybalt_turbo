<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/job_rate_sheets.go","sha256":"a191272dbd1286492e70529806f896656b9d4766ab8433b25ae7a3a5f4c5ae0c","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/job_rate_sheets.go

## Opening period skips the active rate sheet check

Lines 26-29

Asked as: inactive rate sheet opening period; deactivated sheet copied to period; opening period rate sheet not active; why opening period ignores active check

An opening period copied from the project's current rate sheet (rateSheetFromJobKey) does not need the sheet to be active, because the sheet may have been deactivated since it was set.

## Rate sheet history is ordered with the opening row first

Lines 45-56

Asked as: rate sheet history order; opening row sorts first; history periods wrong order; empty start date at top

A project's rate sheet history lists its periods by start date (effective_date). The opening period has an empty start date and sorts first. The opening period covers all time before the first dated period.

## Opening rate sheet period is created for projects without history

Lines 58-81

Asked as: missing opening period; project has no history rows; no first rate sheet period; rate sheet history empty; create opening period automatically

When a project has a rate sheet but no history rows, Turbo creates an opening period with an empty start date and the project's current rate sheet. Turbo creates no opening period for a proposal or for a project with no rate sheet. The actor is recorded as creator and updater.

## History must keep its opening period

Lines 94-97

Asked as: cannot remove opening row; earliest period has start date; must keep opening row; history lost its first period

After a change, if the earliest period of a project has a start date, Turbo refuses with "The rate sheet history must keep its opening row."

## Neighbouring periods cannot use the same rate sheet

Lines 98-103

Asked as: same rate sheet twice in a row; duplicate neighbouring rate sheets; consecutive periods same sheet; two periods same sheet refused

After a change, if two neighbouring periods use the same rate sheet, Turbo refuses with "Two neighbouring periods would use the same rate sheet. Delete the later one instead."

## Project's rate sheet follows the latest history period

Lines 107-127

Asked as: project rate sheet not updating; job rate sheet follows latest period; project sheet wrong after history change; jobs rate sheet out of sync

Turbo sets the project's rate sheet (jobs.rate_sheet) to the sheet of the latest period. It saves the project only when the sheet differs. It marks the project as no longer imported and skips validation of other fields.

## Rate sheet period needs an existing non-proposal project

Lines 131-137

Asked as: cannot add rate sheet to proposal; project not found rate sheet; rate sheet period for proposal; period missing project error

Turbo refuses a rate sheet period for a missing project with "The project was not found." It refuses a period for a proposal with "Proposals do not have rate sheets."

## Rate sheet period cannot move to another project

Lines 139-143

Asked as: move period to other project; change project of rate sheet period; reassign rate sheet period; wrong project on period

When a person edits a period, changing its project is refused with "A rate sheet period cannot move to another project."

## Opening period cannot gain a date and dated periods cannot lose theirs

Lines 144-147

Asked as: cannot add date to opening period; clear start date refused; remove start date from period; opening period date error

When a person edits a period, Turbo refuses adding a start date to the opening period or clearing the start date of another period. The message is "The opening period has no start date, and other periods must keep one."

## Period start date must be valid and not in the future

Lines 149-156

Asked as: invalid start date; future start date refused; date in the future rate sheet; start date tomorrow not allowed; bad date on period

A period start date must be a valid date, else Turbo says "Enter a valid date." It cannot be after today's UTC date, else Turbo says "The start date cannot be in the future." Today itself is allowed.

## Rate sheet checks run only when sheet or date changes

Lines 158-162

Asked as: edit period without changing sheet; saving period checks again; why rate sheet checks on edit; rate sheet validation on small edit

On an edit that changes neither the rate sheet nor the start date, Turbo skips the rate sheet existence, active and start date checks. A new period always gets them.

## Rate sheet must exist and be active when newly chosen

Lines 163-172

Asked as: rate sheet not found; inactive rate sheet chosen; cannot pick deactivated sheet; selected sheet not active; sheet disappeared from period

Turbo refuses a period whose rate sheet is not found with "The rate sheet was not found." A newly chosen sheet must be active, else Turbo says "The selected rate sheet is not active." A period keeps its sheet after the sheet is deactivated.

## Period cannot start before the rate sheet took effect

Lines 173-178

Asked as: start date before rate sheet effective; period earlier than sheet effective date; rate sheet took effect later; start date too early for sheet

For a dated period, if the start date is before the rate sheet's own effective date, Turbo refuses with "The start date is before <sheet name> took effect on <date>." A start date equal to the effective date passes. The opening period is exempt.

## Rate sheet periods change only through history actions

Lines 184-199

Asked as: cannot edit rate sheet period directly; rate sheet period edit refused; use history actions instead; delete rate sheet period error

Creating, updating or deleting a rate sheet period outside the rate sheet history actions is refused with "use the rate sheet history actions". Row validation runs on create and update, not on delete.

## Opening period is written when a project first gets a rate sheet

Lines 201-213

Asked as: new project rate sheet opening period; first rate sheet on project; project saved with rate sheet first time; opening period created on save

When a project is created, or saved after having no rate sheet, Turbo creates its opening period after the job save, in the same transaction. The user who set the sheet is the actor.
