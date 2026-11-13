<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_absorb.go","sha256":"43df139019dd2556181f7cecfaa469c0b5abcdd622c506c291d8b35e7836fa85","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_absorb.go

## What an absorb moves for rate sheets

Lines 13-15

Asked as: jobs moved to other sheet; where did my jobs go; move project requests on merge; rate sheet merge jobs

An absorb moves the jobs and project requests (job_requests) of the absorbed rate sheet to the target. The rate sheet's entries are owned rows and are deleted with it. The absorb never reprices time, so the rates must be identical.

## Absorbing rate sheets moves jobs and project requests

Lines 33-52

Asked as: what happens to sheet entries; entries deleted after absorb; old sheet entries gone; merged sheet jobs move

A rate sheet absorb moves the jobs (jobs.rate_sheet) and project requests (job_requests.rate_sheet) of each absorbed sheet to the target. Turbo also deletes the absorbed sheets' rate sheet entries (rate_sheet_entries). Turbo flags the absorb as one that changes jobs (ChangesJobs).

## Who can absorb rate sheets

Lines 43-48

Asked as: who is allowed to absorb; absorb button missing; need admin to merge; can't merge rate sheets permission

Absorbing rate sheets needs the admin claim, or both the absorb claim and the rate sheet revise claim (rate_sheet_revise). The absorb claim alone is not enough, because an absorb moves jobs between contract rate sheets.

## Absorb refusal reasons are shown to the user

Lines 56-64

Asked as: why was absorb refused; absorb error message; rate sheet merge failed reason; why can't I merge

When Turbo refuses a rate sheet absorb, it shows one of these reasons: "A sheet can't absorb itself.", "It is a revision of the target sheet.", "The target has no rates.", "Its rates are not identical to the target's.", "Its latest revision's rates are not identical to the target's.", "Its name has revisions, and a job or project request uses one of them.", or "An inactive target can't absorb an active sheet." A refusal for a later effective date has its own message.

## Absorb counts the jobs and project requests of each revision

Lines 90-100

Asked as: how many jobs use revision; revision still in use; count of requests per sheet; sheet revision usage count

For every rate sheet revision, Turbo counts the jobs (job_count) and the project requests (request_count) that use it. The absorb check for a name with several revisions refuses the absorb when any revision has a count above 0.

## Absorb compares each role's rate and overtime rate

Lines 105-131

Asked as: rates don't match; different overtime rate; same rates but not matching; role rate mismatch between sheets

Turbo loads every entry of every rate sheet with its role, rate (rate) and overtime rate (overtime_rate). It reduces them to a rates key (RateSheetRatesKey, defined in another file). Two sheets count as having identical rates only when their keys are equal.

## Revisions of a name are ordered by revision number

Lines 135-140

Asked as: which revision is newest; latest revision number; order of sheet revisions; highest revision sheet

Turbo groups rate sheets by name and sorts each group by revision number, oldest first. It records how many revisions each name has (Revisions). The latest revision of a name is the one with the highest revision number.

## A rate sheet can't absorb itself

Lines 147-149

Asked as: merge sheet into itself; absorb same sheet; self merge rate sheet

Turbo refuses an absorb when the source is the target. The reason is "A sheet can't absorb itself."

## Revisions of the target's own name can't be absorbed

Lines 150-152

Asked as: merge older revision of same sheet; absorb revision into its own sheet; same name sheet merge

Turbo refuses a source with the same name as the target. The reason is "It is a revision of the target sheet."

## Target with no rates can't absorb

Lines 153-155

Asked as: target sheet has no rates; empty target sheet; merge into blank sheet

Turbo refuses an absorb when the target has no rate entries. The reason is "The target has no rates."

## Single-revision sheet absorb rules

Lines 158-166

Asked as: effective date blocks absorb; single sheet absorb dates; start date rule merge; why date stops merge

A sheet that is the only revision of its name can be absorbed only if its rates are identical to the target's. Otherwise Turbo says "Its rates are not identical to the target's." The target's effective date must not be later than the sheet's. The same date is allowed. Otherwise Turbo says "The target's effective date, X, is later than this sheet's, Y."

## Single-revision sheet refused when target's effective date is later

Lines 159-166

Asked as: target starts after sheet; effective date too late; later start date merge; target date after sheet date

For a sheet that is the only revision of its name, Turbo refuses the absorb when the target's effective date is later than the sheet's. The message is "The target's effective date, X, is later than this sheet's, Y." Equal dates pass. Turbo checks the rates first, then the date, then the active status.

## Inactive target can't absorb an active sheet

Lines 167-169

Asked as: inactive sheet can't take active; merge into inactive sheet; active sheet inactive target; inactive sheet absorb blocked

An inactive target can't absorb an active sheet. Turbo says "An inactive target can't absorb an active sheet." Absorb does not change the active flag, so the active sheet would leave the Rate Sheet selector. For a name with revisions, the refusal applies when any revision is active.

## Absorbing a name with several revisions

Lines 173-189

Asked as: merge all revisions of sheet; absorb whole sheet name; revisions used by job blocks; merge sheet with many revisions

A name with several revisions can be absorbed only if no job and no project request uses any of its revisions. Otherwise Turbo says "Its name has revisions, and a job or project request uses one of them." When allowed, every revision is absorbed together. Only the latest revision's rates must match the target's, or Turbo says "Its latest revision's rates are not identical to the target's." The effective date rule does not apply, because nothing moves.

## Multi-revision name absorb checks only the latest revision's rates

Lines 182-189

Asked as: only latest revision rates matter; older revisions rates different; merge sheet with old revisions; revisions rate check

When a name has several revisions, only the latest revision's rates must be identical to the target's. Older revisions' rates are not compared, because nothing uses them. The inactive-target check also applies: an inactive target can't absorb the name if any revision is active. The effective date is not checked.

## Absorb refusal message names both sheets

Lines 200-224

Asked as: absorb error names sheets; which sheet can't be absorbed; rev number absorb message; refused absorb message format

When a sheet can't be absorbed, Turbo refuses the whole absorb with a message of the form "<name> (rev. N) can't be absorbed into <name> (rev. N). <reason>". An unknown target or source sheet is also refused as not found. Turbo adds the other revisions of a name that is absorbed as a whole, once each.

## Absorb request fails if a sheet is not found

Lines 200-211

Asked as: target rate sheet not found; rate sheet was not found; missing sheet during absorb; selected sheet missing

Turbo refuses a rate sheet absorb when the target sheet is not found: "the target rate sheet was not found". It also refuses when a selected sheet is not found: "rate sheet <id> was not found".

## Absorbing a whole name absorbs each revision once

Lines 217-223

Asked as: same revision listed twice; duplicate revisions in absorb list; absorb list duplicates; revision shown twice

When a name is absorbed as a whole, Turbo adds every revision of that name to the absorb list. A revision that two selections share is listed only once.

## Absorb writes rate sheet name history on commit

Lines 232-284

Asked as: former names after absorb; rename history moved; formerly list absorb; undo absorb history

A rate sheet absorb must be committed by a user, or Turbo refuses with "a rate sheet absorb must be committed by a user". On commit, the history of renames of the absorbed sheets moves to the target. Each absorbed name that differs from the target's name is added once to the target's history (the "Formerly" list), with the committer as the person who renamed. Turbo does this on commit, not on absorb, so an undo has nothing to remove.

## Rate sheet absorb must be committed by a user

Lines 233-235

Asked as: absorb commit fails; absorb run by script; absorb committed by system; commit absorb not user

Turbo refuses to commit a rate sheet absorb unless a signed-in user commits it: "a rate sheet absorb must be committed by a user".

## Name history records the absorbed names, not the target's own

Lines 240-284

Asked as: formerly list target name; old names recorded on merge; renamed by who absorb; name history rows merge

On commit, Turbo moves the name history of the absorbed sheets to the target. It adds one history row for each distinct absorbed name that differs from the target's name. Each row records the old name, the target's name and the committing user (renamed_by).

## Duplicate list shows target, absorb ability, reason and count

Lines 288-303

Asked as: duplicate rate sheets target; can absorb column; how many sheets deleted by absorb; duplicate sheet list reason

In a group of identical sheets, Turbo marks one sheet as the target. For each sheet it shows whether it can absorb at least one other sheet of the group (CanAbsorb). For a sheet other than the target it also shows why the target can't absorb it, and how many sheets the absorb deletes (Absorbs): 1, or every revision of the name. The reason is empty for the target and for a sheet that the target can absorb.

## Duplicate groups need at least two names

Lines 305-311

Asked as: duplicate group single name; revisions show as duplicates; identical sheets same name; duplicate needs two names

A duplicate group holds sheets with identical rates from at least two names, both active and inactive. Revisions of one name are not duplicates of each other. A group can carry a reason why no sheet in it can be absorbed into its target.

## Duplicate rate sheet groups and their target

Lines 317-357

Asked as: which sheet is duplicate target; earliest effective date duplicate; group identical rates target; duplicate sheets with same rates

Turbo groups rate sheets with identical rates when the group has at least two names. Revisions of one name are not duplicates of each other, and sheets without entries are in no group. The group's target is the sheet with the earliest effective date. On a tie an active sheet goes first, then the one with the most jobs, then the name and revision.

## Duplicate group target is the earliest sheet

Lines 323-357

Asked as: tie breaker duplicate sheets; which duplicate sheet is kept; earliest sheet kept duplicates; duplicate target tie

Turbo groups rate sheets with identical rates when the group has at least two names. Sheets without entries are in no group, and revisions of one name are not duplicates of each other. The target is the sheet with the earliest effective date. On a tie, an active sheet goes first, then the one with more jobs, then the name, then the lower revision.

## Duplicate list shows who can absorb and why not

Lines 361-379

Asked as: why can't absorb duplicate; marked able to absorb; deleted sheets count duplicates; duplicate list absorb reason

For each non-target sheet in a duplicate group, Turbo shows the reason the target can't absorb it, and the number of sheets the absorb deletes. A sheet is marked as able to absorb when it can absorb at least one other sheet of its group.

## Duplicate group reasons when nothing can be absorbed

Lines 381-391

Asked as: no sheet can be absorbed; earliest sheet inactive message; activate earliest sheet; duplicate group no absorb

When no sheet in a group can be absorbed into the target, the reason is "No sheet in this group can be absorbed into the earliest sheet." If an inactive refusal caused it and the target's name has a newer revision, the reason is "The earliest sheet is inactive and has a newer revision. It can't absorb an active sheet, and it can't be activated." If it has no newer revision, the reason is "The earliest sheet is inactive and can't absorb an active sheet. Activate it first."

## Duplicate groups are sorted by the target's name and revision

Lines 395-401

Asked as: duplicate groups order; sort duplicate list; duplicate list sorting

Turbo lists duplicate groups in order of their target's name, then the target's revision.

## Absorb candidates are empty for an unknown target

Lines 415-427

Asked as: absorb candidates target not found; no candidates for sheet id; candidate list unknown sheet

The list of absorb candidates for a target sheet reports that the target was not found when no rate sheet has the given id. Each candidate comes with the reason when the target can't absorb it. The reason is empty when it can.

## Absorb candidates for a target sheet

Lines 428-448

Asked as: sheets that can be merged into; list of sheets to absorb; candidates sorted by name; which sheets can merge

For a target with rates, Turbo lists every sheet of another name with identical rates. Each has the reason the target can't absorb it, or none if it can. It also shows how many sheets the absorb deletes. The list is sorted by name, then revision.

## Absorb candidates need a different name and identical rates

Lines 429-448

Asked as: candidate list empty; same name candidates; sheets with same rates list; candidates no rates target

The candidate list for a target holds only sheets of another name with identical rates. It is empty when the target has no rates. Each candidate shows the refusal reason, if any, and how many sheets the absorb deletes. Turbo sorts candidates by name, then revision.
