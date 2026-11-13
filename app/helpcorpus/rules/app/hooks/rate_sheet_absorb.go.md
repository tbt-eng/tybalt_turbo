<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_absorb.go","sha256":"0027063b1210c34aaeae4fd0f5fa112a553fa9dca607304c91b036d06c187073","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_absorb.go

## Rate sheet absorb never reprices time

Lines 13-52

Asked as: absorb changes my prices; merging rate sheets reprices jobs; do rates change after merge; absorb alters rates

Absorbing a rate sheet moves its jobs, project rate sheet periods (job_rate_sheets) and project requests to the target. So the two sheets must have identical rates. Turbo also repoints rate sheets that were based on the absorbed sheet (rate_sheets.based_on) to the target.

## What counts as a job using a rate sheet

Lines 18-20

Asked as: why does sheet show in use; which jobs use this rate sheet; old period still uses sheet; sheet in use by job

A job uses a rate sheet when the job's rate sheet (jobs.rate_sheet) is that sheet, or when any period in the job's rate sheet history is on that sheet. A sheet that only an old period uses still counts as in use. Project requests that use the sheet are counted separately.

## Any revision can be the target

Lines 35-35

Asked as: pick older revision as target; merge into a specific revision; target an old revision; any revision as target

Any revision of any rate sheet can be the target of an absorb.

## Absorbing rate sheets moves jobs, periods and requests to the target

Lines 41-57

Asked as: where do jobs go after merge; jobs moved to other sheet; what happens to requests when merging; move periods to target

When a rate sheet is absorbed, Turbo moves its jobs, project rate sheet periods (job_rate_sheets) and project requests to the target. Turbo also moves the sheets that proposed rates were based on (based_on). The absorbed sheet's rate entries are owned rows and are removed with it. A job whose only moved reference is an older period is not touched, so its jobs.rate_sheet field is unchanged.

## Who can absorb rate sheets

Lines 58-64

Asked as: who is allowed to absorb; absorb permission denied; need admin to merge sheets; can't merge rate sheets access

Absorbing rate sheets needs the admin claim, or both the absorb claim and the rate sheet revise claim (rate_sheet_revise). The absorb claim alone is not enough, because an absorb moves jobs between contract rate sheets.

## Rate sheet absorb needs two claim sets, not absorb alone

Lines 58-63

Asked as: absorb claim not enough; missing rate sheet revise claim; need both claims to merge; merge access claims

Turbo lets a person absorb rate sheets only with the admin claim, or with both the absorb claim and the rate sheet revise claim (rate_sheet_revise). The absorb claim alone is not enough, because an absorb moves jobs between contract rate sheets.

## Absorb reason texts shown to users

Lines 70-79

Asked as: what does absorb refusal mean; absorb error message; why can't absorb rate sheet; rate sheet merge error text

The reasons are: "A sheet can't absorb itself.", "It is a revision of the target sheet.", "The target has no rates.", "Its rates are not identical to the target's.", "Its latest revision's rates are not identical to the target's.", "Its name has revisions, and a job or project request uses one of them." and "An inactive target can't absorb an active sheet.". The date refusal reads "The target's effective date, X, is later than this sheet's, Y.".

## Job count includes jobs using the sheet only in old periods

Lines 114-119

Asked as: job count looks wrong; why job count is higher; count jobs in old periods; request count on rate sheet

For the absorb rules, a rate sheet's job count (job_count) counts each job once, whether jobs.rate_sheet is the sheet or any period of the job's rate sheet history (job_rate_sheets) is on it. The request count (request_count) counts project requests (job_requests) whose rate sheet is the sheet.

## Proposed-rate sheets are left out of absorb data

Lines 120-122

Asked as: proposed sheet missing from list; request sheets not showing in absorb; can't find request rate sheet; hidden sheets in merge

Rate sheets that belong to a project request (request is not empty) are not loaded for absorb. They are never a target, never absorbed, and never in a duplicate group or candidate list.

## Proposed rates sheets are excluded from absorb

Lines 121-122

Asked as: absorb a proposed rates sheet; merge request's proposed rates; proposed rates as target; absorb request sheet

A project request's proposed rates sheet is never absorbed and never a target. Turbo leaves out every rate sheet that belongs to a request.

## Absorb rate comparison covers every role's rate and overtime rate

Lines 128-157

Asked as: overtime rate ignored in comparison; same rates but different overtime; role rates compared when merging; identical rates check overtime

Turbo decides that two sheets have identical rates by comparing a rates key (RateSheetRatesKey) built from the role, rate and overtime rate of each rate sheet entry. A sheet with no entries has an empty key and takes part in no duplicate group.

## Revisions of a name are counted oldest first

Lines 153-163

Asked as: which revision is latest; order of revisions; latest revision checked; revision numbering order

Turbo groups rate sheets by name and orders the revisions by revision number, oldest first. The last one is the latest revision, and the check for a name with revisions compares its rates with the target's.

## A rate sheet cannot absorb itself

Lines 170-172

Asked as: merge sheet into itself; absorb same rate sheet; same sheet absorb error; sheet absorbs itself

Turbo refuses with "A sheet can't absorb itself." when the source and target are the same sheet.

## A revision of the target's own name cannot be absorbed

Lines 173-175

Asked as: merge revision into own name; absorb sibling revision; same name sheet absorb; revision of target refused

Turbo refuses with "It is a revision of the target sheet." when the source has the same name as the target.

## Target must have rates

Lines 176-178

Asked as: target has no rates; empty target rate sheet; merge into sheet without rates; target sheet blank

Turbo refuses with "The target has no rates." when the target rate sheet has no entries.

## Single-revision sheet needs identical rates

Lines 181-184

Asked as: rates not identical merge; different rates can't merge; rate mismatch absorb; single revision rates differ

A sheet that is the only revision of its name can be absorbed only when its rates are identical to the target's. Otherwise Turbo refuses with "Its rates are not identical to the target's."

## Target effective date must not be later than the sheet's

Lines 185-189

Asked as: effective date later error; target date after sheet; effective date blocks absorb; why date is refused

For a single-revision sheet, Turbo refuses when the target's effective date is later than the sheet's. The message is "The target's effective date, X, is later than this sheet's, Y." Equal dates are allowed.

## Inactive target cannot absorb an active sheet

Lines 190-192

Asked as: inactive target can't absorb; absorb active sheet into inactive; activate target before merging; inactive sheet merge error

Turbo refuses with "An inactive target can't absorb an active sheet." when the target is inactive and the sheet is active. For a name with revisions, the refusal applies when any revision is active. Absorb does not change the active flag.

## Single-revision sheet absorbs alone

Lines 193-193

Asked as: absorb only one sheet; what gets deleted on single revision; does absorb delete just this sheet; single sheet merge result

When all checks pass for a single-revision sheet, only that sheet is absorbed and deleted.

## Name with revisions absorbed only when unused

Lines 196-201

Asked as: revisions in use can't absorb; sheet with revisions used by job; name has revisions error; used revisions block merge

A name with several revisions can be absorbed only when no job, no period and no project request uses any of its revisions. Otherwise Turbo refuses with "Its name has revisions, and a job or project request uses one of them." The effective date rule does not apply, because nothing moves.

## Only the latest revision must match the target's rates

Lines 205-208

Asked as: latest revision rates mismatch; older revision rates don't matter; latest revision not identical; only newest revision compared

For a name with revisions, only the latest revision's rates must be identical to the target's. Otherwise Turbo refuses with "Its latest revision's rates are not identical to the target's." Older revisions are unreferenced, so their rates do not matter.

## Absorbing a name with revisions takes every revision

Lines 212-212

Asked as: absorb deletes all revisions; select one revision all removed; merge every revision together; all revisions go with name

When a name with revisions passes the checks, Turbo absorbs and deletes all of its revisions together. Selecting one revision is enough.

## Refusal message when absorbing rate sheets

Lines 218-239

Asked as: whole absorb refused; target not found error; one sheet blocks absorb all; absorb failed for selected sheet

When a selected sheet cannot be absorbed, Turbo refuses the whole absorb. The message is "<name> (rev. N) can't be absorbed into <target> (rev. N). <reason>". Turbo also refuses when the target or a selected sheet is not found.

## Absorbed ids are deduplicated

Lines 228-247

Asked as: same sheet absorbed twice; duplicate selection in absorb; selected revisions twice; repeated sheet in merge

When several selected sheets share a name, Turbo absorbs each sheet once.

## Absorbed names go into the target's name history

Lines 250-307

Asked as: formerly name after absorb; old name kept on target; name history rate sheet; where old sheet name goes

On commit, each absorbed name that differs from the target's name is added once to the target's name history. The target shows it after "Formerly". The history is written on commit, so an undo has nothing to remove. The name history rows of the absorbed sheets move to the target.

## Only a user can commit a rate sheet absorb

Lines 256-258

Asked as: absorb commit fails; committed by not a user; absorb by script error; commit rate sheet absorb error

Turbo refuses with "a rate sheet absorb must be committed by a user" when the committer is not a user record.

## Duplicate groups list identical rates from different names

Lines 346-361

Asked as: duplicate rate sheets list; sheets with same rates; find duplicate rate sheets; identical rates different names

The duplicates list groups rate sheets with identical rates. A group needs at least two different names, because revisions of one name are not duplicates of each other. Sheets without entries are in no group.

## Duplicate groups skip sheets that have no rates

Lines 346-351

Asked as: empty sheets in duplicates; no rates not in duplicates; blank rate sheet duplicate; sheets without rates listed

A rate sheet with no entries (empty rates key) is never in a duplicate group. In the candidate list, a target with no rates gets no candidates.

## Duplicate group target is the earliest sheet

Lines 363-380

Asked as: which sheet is duplicate target; earliest sheet chosen; why this sheet is main; target of duplicate group

The target of a duplicate group is the sheet with the earliest effective date. On a tie Turbo prefers an active sheet, then the sheet with more jobs, then name, then revision.

## Duplicate group order of sheets and groups

Lines 363-378

Asked as: order of duplicate sheets; why sheets listed this order; sort duplicate group; group sheet ordering

Within a group of identical sheets, Turbo orders them by earliest effective date. A tie goes to an active sheet, then to the sheet with more jobs, then to name, then to revision. Turbo lists the groups by the name and revision of their first sheet.

## Duplicate group shows whether each sheet can absorb

Lines 384-402

Asked as: can this sheet absorb; how many sheets deleted; absorb column duplicates; can absorb status duplicate

For each sheet in a duplicate group, Turbo shows whether it can absorb at least one other sheet. For non-targets it shows the reason the target cannot absorb it, and how many sheets the absorb deletes.

## Reason shown when no sheet in a group can be absorbed

Lines 404-414

Asked as: no sheet can be absorbed; earliest sheet inactive; why group has no absorb option; group can't absorb anything

When the group's target can absorb no sheet, the reason is "No sheet in this group can be absorbed into the earliest sheet." When an inactive refusal is the cause, the reason differs. If the target has a newer revision, it is "The earliest sheet is inactive and has a newer revision. It can't absorb an active sheet, and it can't be activated." Otherwise it is "The earliest sheet is inactive and can't absorb an active sheet. Activate it first."

## Duplicate groups are sorted by the target's name and revision

Lines 418-425

Asked as: order of duplicate groups; sort rate sheet groups; duplicate groups sorting; group list order

Duplicate groups are listed by the name of the group's first sheet (the target), then by its revision.

## Absorb candidates for a target

Lines 441-472

Asked as: candidates to absorb; which sheets can merge into this; list merge options; absorb options for sheet

The candidates for a target are the sheets of other names with identical rates. Each shows the reason the target cannot absorb it, or how many sheets an absorb deletes. A target with no rates has no candidates.

## Candidates exclude the target's own name and different rates

Lines 451-464

Asked as: candidate list different rates; why sheet not candidate; same name not listed; candidates excluded by rates

The absorb candidates for a target are the sheets of other names whose rates are identical to the target's. Each candidate shows the reason the target can't absorb it, or an empty reason, and the number of sheets the absorb deletes.
