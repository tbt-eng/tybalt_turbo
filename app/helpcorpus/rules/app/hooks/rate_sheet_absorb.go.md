<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_absorb.go","sha256":"0704d3d0ab52b246485ca522846a499662fdbfa640ebd842d8d79804eef0eb6a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_absorb.go

## What an absorbed rate sheet moves to the target

Lines 43-61

Asked as: jobs moved after merge; where did my jobs go; project requests moved to other sheet; move jobs to another rate sheet

When a rate sheet is absorbed, its jobs, project rate sheet periods (job_rate_sheets) and project requests (job_requests) move to the target. Sheets that were based on it (based_on) also point to the target. Its rate entries are owned rows and are handled with the sheet. An absorb never reprices time, so the rates must match.

## Claims needed to absorb rate sheets

Lines 62-71

Asked as: can't absorb rate sheets; no permission to merge sheets; absorb button missing; need admin to merge rates

A person needs the admin claim, or both the absorb claim and the rate sheet revise claim (rate_sheet_revise), to absorb rate sheets. The absorb claim alone is not enough, because an absorb moves jobs between contract rate sheets.

## A sheet can't absorb itself

Lines 76-76

Asked as: merge sheet into itself; same sheet as target; absorb sheet into itself error

Turbo refuses to absorb a rate sheet into itself. The message is "<sheet> can't be absorbed into <target>. A sheet can't absorb itself."

## A revision of the target can't be absorbed into it

Lines 77-77

Asked as: merge revision into its parent; same code merge refused; merge revision into target

Turbo refuses when the sheet has the same code as the target. The reason is "It is a revision of the target sheet."

## Target with no rates refuses absorb

Lines 78-78

Asked as: target has no rates; merge into empty rate sheet; target sheet is empty

Turbo refuses an absorb when the target has no rate entries. The reason is "The target has no rates."

## Refusal message when the target's rates differ from a single-revision sheet

Lines 79-79

Asked as: rates don't match merge; different rates can't merge; merge refused rates differ

A rate sheet with only one revision can be absorbed only when its rates are identical to the target's. If they differ, Turbo refuses with "Its rates are not identical to the target's."

## Refusal message for a multi-revision sheet whose latest rates differ

Lines 80-80

Asked as: latest revision rates different; last revision doesn't match target; merge revisions rates mismatch

When a sheet with several revisions is absorbed, Turbo compares only its latest revision with the target. If they differ, Turbo refuses with "Its latest revision's rates are not identical to the target's."

## Refusal message when a revised sheet is in use

Lines 81-81

Asked as: revision in use can't merge; job uses old revision; sheet revision being used

A sheet with revisions cannot be absorbed if a job or project request uses any of its revisions. Turbo refuses with "It has revisions, and a job or project request uses one of them."

## Refusal message for an inactive target and an active sheet

Lines 82-82

Asked as: merge into inactive sheet; inactive target active sheet; disabled sheet can't absorb

An inactive target cannot absorb an active sheet. Turbo refuses with "An inactive target can't absorb an active sheet."

## Job count includes old periods

Lines 93-96

Asked as: job count includes old periods; old rate period counted; job still counted on sheet

A job counts as using a rate sheet when the sheet is its current rate sheet or when any period in its rate sheet history uses it. Project requests using the sheet are counted separately.

## Rate sheet job count includes jobs using it only in old periods

Lines 125-129

Asked as: why job counted on sheet; jobs counted in old periods; job count on rate sheet wrong

Turbo counts a job as using a rate sheet when the job's current rate sheet (jobs.rate_sheet) is that sheet. It also counts a job when any period in the job's rate sheet history (job_rate_sheets) is on that sheet. Each job counts once.

## Request count counts project requests on the sheet

Lines 130-130

Asked as: request count on rate sheet; how many requests use sheet; project requests counted on sheet

Turbo counts the project requests (job_requests) whose rate sheet is the sheet. This count (request_count) decides whether a sheet with revisions can be absorbed.

## Request proposed rates are never absorbed or targets

Lines 132-133

Asked as: proposed rates sheet; can't merge proposed rates sheet; request proposed rate sheet target

Rate sheets that hold a project request's proposed rates are left out of absorb checks. They can be neither absorbed nor a target.

## Rates are compared as a whole set of role entries

Lines 139-168

Asked as: how rates are compared; identical rates check; same role rates match

Turbo compares each sheet's rate sheet entries (role, rate, overtime rate) as one rates key built by RateSheetRatesKey. Two sheets have identical rates when the keys are equal. A sheet with no entries has an empty key.

## Revisions of a sheet are ordered oldest first

Lines 164-174

Asked as: revision order; which revision is latest; sort revisions by number

Turbo groups rate sheet revisions by code and sorts each group by revision number, oldest first. Each sheet learns how many revisions its code has. The absorb checks use the latest revision as the last one in the list.

## Absorb checks run in a fixed order

Lines 180-189

Asked as: which merge error shows first; several merge errors; why this refusal reason first

Turbo checks in this order: a sheet absorbing itself, a revision of the same code as the target, a target with no rates. Only then does it check the rates, effective date and active status. The first failing check gives the refusal reason.

## Which sheet counts as the same code

Lines 184-186

Asked as: same code as target; code matches target refused; sheets with same code

Sheets are compared by code, not by the selected revision. A sheet whose code matches the target's code is refused as a revision of the target.

## Absorbing a sheet that has a single revision

Lines 191-204

Asked as: merge single revision sheet; only that sheet merged; one revision merge rules

A sheet with only one revision can be absorbed only if its rates are identical to the target's. Otherwise Turbo refuses with "Its rates are not identical to the target's." Only that sheet is absorbed.

## Absorbed sheets are deleted: one sheet or all revisions

Lines 191-224

Asked as: what gets deleted on merge; sheet deleted after absorb; all revisions deleted together

Absorbing a single-revision sheet deletes only that sheet. Absorbing a sheet with revisions deletes every revision of its code together.

## Target effective date must not be later than the sheet's

Lines 196-200

Asked as: target start date later; effective date merge error; target date later than sheet

When the sheet has one revision, Turbo refuses if the target's effective date is later than the sheet's. The message is "The target's effective date, X, is later than this sheet's, Y." Equal dates pass. Dates compare as YYYY-MM-DD text.

## Effective date refusal message names both dates

Lines 198-200

Asked as: effective date refused message; dates in merge message; target effective date too late

When the target's effective date is later than a single-revision sheet's, Turbo refuses with "The target's effective date, X, is later than this sheet's, Y." Equal dates are allowed.

## Inactive target can't absorb an active sheet

Lines 201-203

Asked as: inactive target active revision; merge revisions into inactive sheet; inactive target blocks merge

Turbo refuses when the target is inactive and the sheet is active, because absorb does not change active and the sheet would leave the Rate Sheet selector. The reason is "An inactive target can't absorb an active sheet." For a sheet with revisions, the refusal happens if any revision is active.

## Absorbing a sheet that has revisions

Lines 207-215

Asked as: merge all revisions together; absorb sheet with several revisions; revisions moved together

A sheet with several revisions can be absorbed only if no job, period or project request uses any revision. Otherwise Turbo refuses with "It has revisions, and a job or project request uses one of them." All revisions are absorbed together.

## Only the latest revision's rates must match

Lines 216-219

Asked as: only latest revision rates; last revision must match target; merge revisions rates check

When all revisions of a sheet are absorbed, only the latest revision must match the target's rates. Otherwise Turbo refuses with "Its latest revision's rates are not identical to the target's." The effective date rule does not apply because nothing moves.

## Inactive target with an active revision

Lines 220-222

Asked as: inactive target any revision active; merge with inactive target and active revision

When absorbing a sheet with revisions, Turbo refuses if the target is inactive and any revision is active. The reason is "An inactive target can't absorb an active sheet."

## Absorb is refused as a whole on the first bad sheet

Lines 229-259

Asked as: one sheet blocks whole merge; multiple sheets merge failed; why whole absorb failed

When several sheets are absorbed, Turbo checks each one. The first sheet that fails refuses the whole absorb with "<source> can't be absorbed into <target>. <reason>". Turbo also refuses if the target or a sheet is not found. A sheet is deleted only once, even if several selections lead to it.

## Absorb is refused when a sheet is not found

Lines 234-245

Asked as: rate sheet not found; target sheet missing; sheet id not found

If the target sheet is not found, Turbo refuses with "the target rate sheet was not found". If a sheet to absorb is not found, it refuses with "rate sheet <id> was not found".

## Absorbing several sheets of one code does not repeat ids

Lines 239-258

Asked as: same sheet selected twice; duplicate selection in merge; sheet listed twice

When several selected sheets lead to the same sheet ids, each id is absorbed only once.

## Absorb refusal message names the sheets and the reason

Lines 246-250

Asked as: can't be absorbed error; which sheet failed merge; merge error message sheet names

When a sheet can't be absorbed, Turbo refuses with "<source> can't be absorbed into <target>. <reason>". Each sheet is named by its code, revision and name.

## Absorbed names go into the target's name history on commit

Lines 267-293

Asked as: formerly name on sheet; old name shown after merge; previous names of rate sheet

When an absorb is committed, each distinct name of an absorbed sheet is added once to the target's name history, shown after "Formerly". Empty names and the target's own name add nothing. The history is written on commit, so an undo has nothing to remove.

## Only a user can commit a rate sheet absorb

Lines 268-270

Asked as: absorb committed by system; must be committed by user; absorb commit error

Turbo refuses to commit a rate sheet absorb unless a user commits it. The message is "a rate sheet absorb must be committed by a user."

## Absorbed sheet names are not duplicated in the history

Lines 283-293

Asked as: duplicate names formerly; same name listed twice; old name repeated in history

The name history skips the target's own name, blank names and repeated names. Each other absorbed name is added once.

## Absorbed sheets' name and code history move to the target

Lines 295-312

Asked as: old codes after merge; earlier codes lead to target; history moves after absorb

On commit, the name history and code history of absorbed sheets move to the target. Their earlier names and codes then lead to the target.

## Name history entries record who committed

Lines 318-327

Asked as: who renamed sheet; name change user record; who merged the sheet history

Each name history entry from an absorb records the old name, the target's name and the committing user (renamed_by).

## Duplicate groups need at least two codes

Lines 360-381

Asked as: duplicate groups list; same rates in different codes; revisions showing as duplicates

Duplicate groups list rate sheets with identical rates from at least two different codes. Revisions of one sheet are not duplicates of each other. Sheets without entries are in no group.

## Duplicate groups ignore sheets without rates

Lines 366-371

Asked as: empty sheet in duplicate group; sheets without rates duplicates; no rates duplicate group

Sheets that have no rate entries are in no duplicate group. Only sheets with identical rates can group.

## How the target of a duplicate group is chosen

Lines 383-400

Asked as: which sheet is group target; earliest effective date target; group target chosen how

The group's target is the sheet with the earliest effective date. Ties go to an active sheet, then the sheet with more jobs, then the lower code, then the lower revision.

## Duplicate group tie-break order

Lines 383-398

Asked as: duplicate group tie break; which sheet wins a tie; tie between sheets in group

The group's target is the sheet with the earliest effective date. On a tie, an active sheet comes first, then the sheet with more jobs, then the lower code, then the lower revision.

## Duplicate group shows who can absorb and why not

Lines 404-422

Asked as: why can't absorb in group; sheets deleted count in group; who can absorb in group

For each sheet in a group, Turbo shows whether it can absorb another sheet of the group. For a non-target sheet it shows the target's refusal reason and the number of sheets the absorb deletes.

## Reason shown when no sheet in a group can be absorbed

Lines 424-434

Asked as: no sheet in group can be absorbed; group says activate it first; earliest sheet can't absorb

When the target can absorb no sheet, the group says "No sheet in this group can be absorbed into the earliest sheet." If the cause is an inactive target with a newer revision, it says the sheet "can't absorb an active sheet, and it can't be activated." If the target is inactive with no newer revision, it says "Activate it first."

## Reason when an inactive earliest sheet blocks the group

Lines 426-433

Asked as: earliest sheet inactive; can't be activated message; inactive earliest sheet blocks group

If no sheet can be absorbed and the earliest sheet is inactive, Turbo says it can't absorb an active sheet. When it has a newer revision, the message adds that it can't be activated. Otherwise it says "Activate it first."

## Duplicate groups are listed by code and revision

Lines 438-444

Asked as: duplicate group order; duplicate groups sorted; list of duplicate groups

Duplicate groups are sorted by the code and then the revision of their first (target) sheet.

## Absorb candidates for a target

Lines 461-492

Asked as: merge candidates for target; which sheets can merge into this; absorb candidates list

For a target sheet, Turbo lists every sheet of another code with identical rates, with the refusal reason or the number of sheets an absorb deletes. A target with no rates has no candidates. Turbo returns nothing if the target does not exist.

## Absorb candidates exclude the target's own code

Lines 472-484

Asked as: candidates with same code; only other codes merge; candidates exclude own code

The candidate list for a target holds sheets with identical rates and a different code. A target with no rates has no candidates.

## Absorb candidates sorted by code and revision

Lines 477-491

Asked as: candidate list order; merge candidates sorted; sort absorb candidates

Candidates are listed by code, then revision. Each shows a refusal reason or how many sheets the absorb deletes.
