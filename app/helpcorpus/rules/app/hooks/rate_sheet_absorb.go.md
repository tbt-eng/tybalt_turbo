<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_absorb.go","sha256":"ad5b6ed1a48d1840c9dc52bcac99d8c4c4922415e15e1e578d00ad1247afa640","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_absorb.go

## What an absorb moves and deletes

Lines 13-16

Asked as: what does absorb move; what gets deleted on merge; does merging reprice time; jobs moved when merging sheets

Absorbing a rate sheet moves its jobs, project rate sheet periods (job_rate_sheets) and project requests (job_requests) to the target. Its rate sheet entries are deleted with it. An absorb never reprices time, so the rates must be identical.

## Rate sheet absorb marks jobs as changed

Lines 41-64

Asked as: jobs changed after merge; why did my jobs move; jobs moved to another rate sheet; merge changed job rate sheet

Absorbing rate sheets moves jobs between contract rate sheets (ChangesJobs is true). Turbo also moves the jobs' rate sheet periods (job_rate_sheets) and project requests (job_requests) of an absorbed sheet to the target. The rate sheet entries of an absorbed sheet are owned rows and go with the sheet.

## Name history rows move on commit, not on absorb

Lines 50-52

Asked as: undo absorb name history; undo merge name changes; old names after undo

Rate sheet name change rows (rate_sheet_name_changes) of absorbed sheets move to the target when the absorb is committed, not when it is made. An undo therefore has nothing to remove.

## Rate sheet absorb needs admin or two claims

Lines 55-60

Asked as: no permission to absorb; can't merge rate sheets; missing claim for absorb; need admin to merge sheets

Absorbing rate sheets requires the admin claim, or both the absorb claim and the rate sheet revise claim (rate_sheet_revise). The absorb claim alone is not enough, because an absorb moves jobs between contract rate sheets.

## Rate sheet absorb refusal reasons shown to the user

Lines 67-76

Asked as: why was absorb refused; absorb error message; rate sheet refused reason; merge blocked why

When Turbo refuses a rate sheet absorb, it shows one of these fixed reasons: "A sheet can't absorb itself.", "It is a revision of the target sheet.", "The target has no rates.", "Its rates are not identical to the target's.", "Its latest revision's rates are not identical to the target's.", "Its name has revisions, and a job or project request uses one of them.", or "An inactive target can't absorb an active sheet."

## A job counts as using a sheet through any period

Lines 85-116

Asked as: job still counts old sheet; sheet shows as in use; old rate period counts as use; why sheet is in use

A job uses a rate sheet when its rate sheet (jobs.rate_sheet) is the sheet, or when any period in its rate sheet history is on the sheet. Turbo counts each such job once. A sheet that only an old period uses still counts as in use. Project requests that use the sheet are counted separately.

## Rates compare as a whole set of entries

Lines 123-152

Asked as: rates look the same but differ; why rates don't match; compare rate entries; same rates different result

Turbo compares two rate sheets by the whole set of their entries (role, rate and overtime rate), through RateSheetRatesKey. A sheet with no entries has an empty rates key. Turbo treats such a sheet as having no rates.

## Revisions of a name are ordered oldest first

Lines 153-158

Asked as: which revision is first; revision order of sheets; revision number ordering

Turbo groups rate sheets by name and orders each name's revisions by revision number, oldest first. Each sheet records how many revisions its name has. Turbo uses this to tell a single-revision sheet from a name with revisions.

## A rate sheet cannot absorb itself

Lines 165-167

Asked as: merge sheet into itself; absorb same sheet; sheet as its own target

Turbo refuses with "A sheet can't absorb itself." when the sheet to absorb is the target.

## Revisions of the target cannot be absorbed

Lines 168-170

Asked as: absorb a revision of target; same name as target; merge revision into its sheet

A sheet with the same name as the target is refused with "It is a revision of the target sheet."

## Target without rates cannot absorb

Lines 171-173

Asked as: target has no rates; merge into empty sheet; empty target sheet absorb

A target with no rate entries refuses every absorb with "The target has no rates."

## Single-revision sheet needs identical rates

Lines 175-179

Asked as: rates not identical to merge; different rates can't merge; merge sheets with different rates

A sheet that is the only revision of its name can be absorbed only when its rates are identical to the target's. Otherwise Turbo refuses with "Its rates are not identical to the target's."

## Target effective date must not be later

Lines 180-184

Asked as: target effective date later; target date after sheet date; effective date blocks absorb

For a single-revision sheet, the target's effective date must be the same as or earlier than the sheet's, because its jobs and periods move. A later target date is refused with "The target's effective date, X, is later than this sheet's, Y." Equal dates pass.

## Inactive target cannot absorb an active sheet

Lines 185-187

Asked as: inactive sheet absorb active; inactive target merge; active sheet vanishes from selector

An inactive target cannot absorb an active sheet, and Turbo refuses with "An inactive target can't absorb an active sheet." Absorb does not change the active flag, and the active sheet would leave the Rate Sheet selector. An active target can absorb inactive sheets.

## Name with revisions: blocked if any revision is in use

Lines 191-196

Asked as: revision in use blocks absorb; job uses old revision; project request uses revision

A name with several revisions can be absorbed only if no job, no rate sheet period and no project request uses any of its revisions. Otherwise Turbo refuses with "Its name has revisions, and a job or project request uses one of them."

## Name with revisions absorbs as a whole; only latest rates compared

Lines 191-207

Asked as: absorb all revisions together; latest revision rates differ; merge every revision of name

When a name with revisions can be absorbed, Turbo absorbs all its revisions together. Only the latest revision's rates must match the target's. Otherwise it refuses with "Its latest revision's rates are not identical to the target's." The effective date rule does not apply, because nothing moves.

## Inactive target rule for revision groups

Lines 204-206

Asked as: inactive target with active revisions; inactive sheet with revisions merge

An inactive target cannot absorb a name with revisions if any of its revisions is active. Turbo refuses with "An inactive target can't absorb an active sheet."

## Absorb is refused with a named reason

Lines 213-234

Asked as: can't be absorbed into; rate sheet id not found; target rate sheet missing

When a target cannot absorb a sheet, Turbo refuses the whole absorb. The message reads: "<name> (rev. N) can't be absorbed into <target name> (rev. M). <reason>". If a selected sheet or the target does not exist, Turbo refuses with "rate sheet <id> was not found" or "the target rate sheet was not found".

## Absorb refusal message format

Lines 218-234

Asked as: absorb stopped by one sheet; absorb message format; target missing on absorb

A refused sheet stops the whole absorb. The message reads "<name> (rev. N) can't be absorbed into <target> (rev. N). <reason>". A missing target gives "the target rate sheet was not found", and a missing source gives "rate sheet <id> was not found".

## Absorb list has no duplicate ids

Lines 223-242

Asked as: same sheet listed twice; duplicate ids in absorb list; sheet appears twice to absorb

When several selected sheets share a name, each id to absorb is listed once.

## Committing a rate sheet absorb needs a user

Lines 250-253

Asked as: commit absorb with service account; api key commit absorb; bot can't commit merge

Turbo commits a rate sheet absorb only when a user (not another kind of account) commits it. Otherwise it refuses with "a rate sheet absorb must be committed by a user".

## Commit requires a user

Lines 251-253

Asked as: who can commit absorb; commit absorb as bot; must be committed by user

A rate sheet absorb can be committed only by a user. Otherwise Turbo refuses with "a rate sheet absorb must be committed by a user".

## Absorbed names recorded as Formerly names

Lines 258-302

Asked as: formerly names; old names after merge; previous names of sheet

On commit, each distinct name of an absorbed sheet is saved once in the target's name history (old name, the target's name, renamed by the committing user). Turbo skips the target's own name and empty names. The target shows these names after "Formerly".

## Formerly name skips the target's own name

Lines 264-302

Asked as: formerly shows target's own name; duplicate name history row; who renamed sheet on merge

On commit, Turbo adds one name history row for each distinct absorbed name. It skips a name equal to the target's name. It records the user who commits as the one who renamed it.

## Name history moves to the target on commit

Lines 278-287

Asked as: name history after merge; earlier names kept after merge; rename history moved to target

On commit, Turbo moves the name history rows (rate_sheet_name_changes) of the absorbed sheets to the target. Their earlier names are kept.

## Duplicate groups: same rates, at least two names

Lines 341-357

Asked as: duplicate rate sheets; sheets with same rates; duplicate finder groups

Duplicate groups hold sheets with identical rates from at least two different names. Sheets without entries are in no group, and revisions of one name are not duplicates of each other.

## Duplicate groups ignore sheets without entries

Lines 341-346

Asked as: empty sheets in duplicates; sheet with no rates duplicate; blank sheet not grouped

The duplicate finder groups sheets by identical rates. A sheet without entries is in no group.

## Duplicate group target choice

Lines 358-375

Asked as: which sheet is the target; earliest effective date target; target choice tie

A group's target is the sheet with the earliest effective date. On a tie, an active sheet goes first, then the sheet with more jobs, then the lower name, then the lower revision.

## Each duplicate sheet shows whether it can absorb

Lines 379-397

Asked as: can this sheet absorb; how many sheets deleted; absorb count for duplicate

For each sheet in a group, Turbo marks it as able to absorb when it can absorb at least one other sheet of the group. For a sheet other than the target, Turbo shows the reason the target cannot absorb it. It also shows how many sheets the absorb deletes: 1, or every revision of the name.

## Duplicate group lists sheets with the target first and a Reason

Lines 387-395

Asked as: reason column in duplicates; absorbs column meaning; why sheet not absorbed in group

In a group of sheets with identical rates, each sheet other than the target gets a Reason from the target's absorb check (empty when the target can absorb it). Each sheet also gets Absorbs, the number of sheets that absorbing it deletes: 1, or every revision of its name.

## Duplicate group reasons when nothing can be absorbed

Lines 399-409

Asked as: no sheet can be absorbed; earliest sheet inactive; activate it first

If no sheet can be absorbed into the target, the group says "No sheet in this group can be absorbed into the earliest sheet." If an inactive refusal caused it and the target's name has a newer revision, it says the earliest sheet is inactive, has a newer revision, can't absorb an active sheet and can't be activated. If it has no newer revision, it says "Activate it first."

## Duplicate groups are sorted by name

Lines 413-419

Asked as: order of duplicate groups; sorted duplicate list; duplicate groups sort order

Turbo lists duplicate groups by the name of each group's first sheet, then by its revision.

## Candidate list for a target

Lines 436-467

Asked as: candidate sheets to absorb; which sheets can merge into target; list of merge candidates

For a target, Turbo lists the sheets of other names with identical rates, each with the reason it cannot be absorbed (empty when it can) and how many sheets the absorb deletes. A target without rates has no candidates. A missing target gives no result.

## Candidate lookup reports a missing target

Lines 441-444

Asked as: candidates target not found; no sheet with that id; candidate lookup empty

If no rate sheet has the target's id, the candidate lookup returns nothing and reports that the sheet was not found.

## Absorb candidates for a target

Lines 446-459

Asked as: candidates for a sheet; sheets with same rates; merge candidates list

For a target with rates, the candidates are sheets of other names with identical rates. Each shows the reason it can't be absorbed, or the number of sheets the absorb deletes. A target with no rates has no candidates.

## Absorb candidates exclude the target's own name and other rates

Lines 447-458

Asked as: same name not a candidate; candidates with other rates; why my revision not listed

The candidate list for a target holds only sheets of another name whose rates are identical to the target's. A target with no rates has no candidates. Each candidate shows the refusal reason, or none, and the number of sheets the absorb deletes.
