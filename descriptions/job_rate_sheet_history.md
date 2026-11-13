# Job rate sheet history

Status: Implemented and deployed (merged to `main` on 2026-11-13)
Last updated: 2026-11-13

## Behaviour before this change

A project has one rate sheet, `jobs.rate_sheet`. It is required when the
project is created and must be active at that time. Proposals have no rate
sheet.

Changing the rate sheet of an existing project is possible, but strongly
discouraged. `validateRateSheetIsActive` in `app/hooks/jobs.go` allows it only
when all of these hold:

- The user is active and holds the `job` claim, as for any project edit.
- The user holds the `rate_sheet_revise` claim. Without it, the save fails
  with `missing_claim`.
- The new sheet is active.
- The sheet is replaced, not cleared. A rate sheet cannot be removed once set.

As of 2026-10-15, no users hold the `rate_sheet_revise` claim. In practice,
nobody can change the rate sheet of an existing project.

The job editor shows the rate sheet selector to every user who can edit the
job. Users without `rate_sheet_revise` learn that they cannot change it only
when the save fails.

### Why changing it is discouraged

Time entries store hours and role, but no rate. `job_priced_time_entries.sql`
prices every entry with the job's current sheet when a report runs. A change
therefore reprices all past time on the job, not only time from the change
onward. WIP, WIP reports, the Staff and Divisions summaries, and the
missing-rates breakdown all show the new values. The app does not keep the
previous values, so they cannot be recovered or compared.

Until this proposal is implemented, `rate_sheet_revise` holders should change a
project's rate sheet only to correct a sheet that was wrong from the start. In
that case, repricing all past time is the intended result. A rate change that
applies only from a date, such as a client accepting a new revision, cannot be
recorded correctly today.

## Problem

Rate sheet revisions already have `rate_sheets.effective_date`, and activating a
revision deactivates the others with the same name. That date says when a
revision became available. It does not say when a job moved to it. A client may
accept a new revision months after it was issued, or never.

So a project cannot move to a new sheet from a given date. The only choices are
to keep the old sheet for all time, or to switch and reprice all past time with
the new one.

## Goal

Record which rate sheet applies to a project from which date. Price each time
entry with the sheet in effect on the entry's date. Keep the history editable
by `rate_sheet_revise` holders so a wrong assignment can be corrected.

Expense allowances already work this way: the `expense_allowance_totals` view joins
the `expense_rates` row with the latest `effective_date` on or before the
expense date.

## Data model

New collection `job_rate_sheets`:

| Field            | Type                                | Notes                                                              |
| ---------------- | ----------------------------------- | ------------------------------------------------------------------ |
| `job`            | relation to `jobs`, required        | Projects only. Cascade delete.                                     |
| `rate_sheet`     | relation to `rate_sheets`, required |                                                                    |
| `effective_date` | text `YYYY-MM-DD`                   | Empty only on the opening row.                                     |
| `reason`         | text                                | Required on every row except the opening row created with the job. |
| `created_by`     | relation to `users`                 | Set by the server.                                                 |
| `updated_by`     | relation to `users`                 | Set by the server.                                                 |

`effective_date` uses the same text format as `time_entries.date` and
`rate_sheets.effective_date`, so string comparison gives date order.

Indexes:

- unique `(job, effective_date)`
- unique partial index on `job` where `effective_date = ''`, so a project has
  at most one opening row

### Row rules

- A project with a rate sheet has exactly one opening row. It covers all time
  before the first dated row, so no entry falls outside the history.
- A dated row applies from its `effective_date`, inclusive, until the next
  dated row.
- Effective dates cannot be in the future. A future date would make
  `jobs.rate_sheet` stop meaning "current" without any write to the job.
- A dated row's date must be on or after its rate sheet revision's
  `effective_date`. A revision cannot apply before it existed. The opening row
  is exempt because it covers all earlier time.
- The rate sheet must be active when a row is created or its `rate_sheet`
  changes. Rows keep inactive sheets after deactivation, as `jobs.rate_sheet`
  does today. Changing only the date or reason does not require an active sheet.
- Two consecutive rows must not use the same sheet. Such a row adds nothing,
  and rejecting it keeps the history readable. The rule is checked after
  every create, edit, and delete. Deleting a dated row, or changing its date,
  must not leave its neighbours on the same sheet: if it would, the change is
  refused, and the user deletes the later duplicate instead. Only a pair that
  the change makes is refused. A rate sheet absorb can leave a pair (see
  Rate sheet absorb below), and that pair does not block other changes.
- The opening row cannot be deleted. Changing its sheet replaces today's
  "change the job's rate sheet" and restates all time up to the first dated row.
- Proposals have no rows. `cleanJob` already clears their rate sheet.

### `jobs.rate_sheet`

Keep the field. A hook keeps it equal to the sheet on the latest row (the
latest `effective_date`, or the opening row if there are no dated rows). It is
written in the same transaction as the history change.

This keeps these working unchanged:

- the legacy writeback in `jobs_writeback.go`, which sends one sheet per job
- the `no_rate_sheet` flag in `job_wip.sql`
- job lists, details, and the job request views that show the current sheet

Once set, `jobs.rate_sheet` changes only through the history. The project
workflow guard refuses any other change unless the save carries
`ProjectWorkflowMutationRateSheet`, in the same way that
`ProjectWorkflowMutationChangeClient` guards `client`. A job save that tries it
fails in `validateRateSheetIsActive` with `use_rate_sheet_history`; clearing
still fails with `cannot_clear`. Create-time checks are unchanged.

A project's first rate sheet can still be set by any job save. A `jobs` model
hook (`AddJobRateSheetHooks` in `app/hooks/job_rate_sheets.go`) then writes the
opening row in the same transaction, recording the user who saved the job. This
covers direct creation, job request approval, and the editor.

### Projects with a rate sheet but no history

The import tool (`import_data/`) writes `jobs.rate_sheet` directly and skips
hooks, so an imported project can have a rate sheet and no rows. Pricing treats
such a project as one opening period on `jobs.rate_sheet`, and the history API
shows it as an opening period with id `opening`. Its row is written on the
first change to its history. That row keeps the sheet even if it has since been
deactivated, since it was valid when set.

## Permissions

- **View:** anyone who can view the job.
- **First rate sheet on an existing project:** a project with no rows gets its
  opening row from any active user with the `job` claim. `rate_sheet_revise` is
  not needed. This keeps today's rule, where setting a project's rate sheet for
  the first time needs only the `job` claim. It matters because most imported
  projects have no rate sheet (see Migration). As today, it reprices that
  project's past time from employee default rates to the new sheet.
- **Every other change to a row:** an active user with both the `job` and
  `rate_sheet_revise` claims, matching today's rule for changing a project's
  rate sheet. This covers dated rows, changing the opening row's sheet, and
  edits and deletes.
- **Opening row on a new project:** created automatically by the `jobs` hook,
  with no `rate_sheet_revise` check, whether the project comes from job request
  approval or direct creation.

Writes go through these routes, each in one transaction that checks the whole
history and then updates `jobs.rate_sheet`:

| Route                                        | Does                                                                        |
| -------------------------------------------- | --------------------------------------------------------------------------- |
| `GET /api/jobs/{id}/rate_sheets`             | The history, plus `can_revise` and `can_set_first` for the user             |
| `POST /api/jobs/{id}/rate_sheets`            | With no `effective_date`, sets a first sheet; with one, starts a new period |
| `PATCH /api/jobs/{id}/rate_sheets/{period}`  | Changes a period's sheet or start date                                      |
| `DELETE /api/jobs/{id}/rate_sheets/{period}` | Deletes a dated period                                                      |

`{period}` is a row id, or `opening` for a project whose history has no rows
yet. Every write except a first sheet needs a reason of up to 500 characters.
The collection's create, update, and delete API rules are null, and a model
hook refuses any row write that does not carry
`ProjectWorkflowMutationRateSheet`.

## Pricing

Only `job_priced_time_entries.sql` changes. Every time valuation already uses
it: `job_wip_api.go`, `wip_reports.go`, `job_staff_summary_api.go`,
`job_divisions_summary_api.go`, and `job_wip_missing_rates.go`.

Turn the history into date ranges for the selected jobs, then join each entry
to the range that contains its date:

```sql
, rate_history AS (
  SELECT jrs.job, jrs.rate_sheet, jrs.effective_date AS starts
  FROM selected_jobs j
  JOIN job_rate_sheets jrs ON jrs.job = j.id
  UNION ALL
  -- A project with a rate sheet but no rows: one opening period.
  SELECT j.id, j.rate_sheet, ''
  FROM selected_jobs j
  WHERE COALESCE(j.rate_sheet, '') != ''
    AND NOT EXISTS (SELECT 1 FROM job_rate_sheets jrs WHERE jrs.job = j.id)
), rate_periods AS (
  SELECT job, rate_sheet, starts,
    LEAD(starts) OVER (PARTITION BY job ORDER BY starts) AS ends
  FROM rate_history
)
```

The opening row's empty date sorts first. Then:

```sql
LEFT JOIN rate_periods p
  ON p.job = te.job
  AND te.date >= p.starts
  AND (p.ends IS NULL OR te.date < p.ends)
LEFT JOIN rate_sheets rs ON rs.id = p.rate_sheet
LEFT JOIN rate_sheet_entries rse
  ON rse.rate_sheet = p.rate_sheet AND rse.role = te.role
```

`te.date >= ''` is always true, so the opening row needs no special case. The
employee default charge-out rate and the Estimated, Partial, and unpriced rules
stay as they are.

### Performance

A project has only a few rows, and the `(job, effective_date)` index covers
the scan. WIP reports load hundreds of jobs at once, so repeat the measurements
in `wip_performance.md` before and after, on the same production snapshot.

## API and UI

### Time valuation

A job's time can be priced by several sheets, so a summary row no longer has
one sheet.

- `JobTimeValuation` replaces `rate_sheet_name` and `rate_sheet_revision` with
  `rate_sheets`: the sheets that priced the row's entries, such as
  "2026 Standard Rates (rev. 0)", joined by "; ". It is in the staff and
  division summaries and their CSV exports.
- The periods themselves, with their dates, come from the history route.
- `TimeSummaryTable.svelte` uses `jobs.rate_sheet` for the "no rate sheet"
  notice instead of inferring it from row data.

### Missing rates

`GET /api/jobs/{id}/wip/missing-rates` returns `rate_sheets`, a list of
`{ id, name, revision, roles }`, one per sheet that priced time with a missing
role rate. A role can be priced in one period and not in the next.

### Job editor and details

- On existing projects that have a rate sheet, the selector is read-only and
  shows the current sheet.
- On existing projects with no rate sheet, `job` holders can still pick one.
  Saving creates the opening row.
- Job details has a Rate sheets section (`RateSheetHistory.svelte`) listing the
  history: sheet name and revision, dates ("From the start" for the opening
  period), reason, and who set it.
- `rate_sheet_revise` holders see New period, Change, and Delete. The new
  period form defaults the date to today and lists only active sheets.
- A project with no rate sheet shows Set rate sheet to `job` holders.
- The forms warn which time will be repriced before a change is saved.

## Rate sheet absorb

Duplicate rate sheets can be absorbed (see the Rate sheets section of
`Absorb.md`). An absorb never reprices time, so it treats the history as
it treats `jobs.rate_sheet`:

- `job_rate_sheets.rate_sheet` is a moved reference. Absorb points every
  row on an absorbed sheet at the target, and undo writes back the saved
  sheet of each row. The rows are changed by SQL, so their `updated` and
  `updated_by` stay as they were and the row hooks do not run.
- A row is a use of its sheet. A sheet that only an old row uses counts as
  used, so a rate sheet name with revisions is not absorbed while any row is
  on one of them.
- The job and its latest row move together, so `jobs.rate_sheet` stays equal
  to the latest row. A job whose only moved reference is an older row is not
  saved, and its `_imported` and `updated` don't change: the writeback sends
  only `jobs.rate_sheet`.
- A dated row is on or after its sheet's effective date, and the target's
  date is not later than the absorbed sheet's, so a moved row still meets
  the revision-date rule. A moved row can be on an inactive target, as a row
  keeps a sheet that is deactivated.
- After a move, two neighbouring rows can be on the target. Pricing reads
  each row's date range and does not care. The history routes refuse only a
  pair that the change makes (`JobRateSheetSamePairs` before the change,
  `ValidateJobRateSheetHistory` after it), and deleting the later row removes
  the pair.
- While the absorb is pending, the history routes refuse with 409 a change
  to a project whose job or rows it moved (`checkPendingRateSheetAbsorb` in
  `app/routes/job_rate_sheets.go`). Undo would otherwise overwrite the change
  or leave `jobs.rate_sheet` unequal to the latest row. An absorbed sheet is
  deleted, so no row can be pointed at it.
- A project request that becomes a job while the absorb is pending gets its
  opening row on the target. `transferRequestMergeReferences` adds that row
  to the undo data with the job.

## Revision rollout

Status: implemented on branch `205` (issue 205), not yet merged.

An `admin` claim holder moves every active project on an earlier revision of
a rate sheet to its latest revision, from one date. Each moved project gets a
dated row on the latest revision, so time before the date keeps its price and
time from the date is repriced. An absorb never reprices time; this does,
which is why it needs `admin` and a confirmation. The routes are in
`app/routes/rate_sheets_rollout.go`:

| Route                                                  | Does                                                        |
| ------------------------------------------------------ | ----------------------------------------------------------- |
| `GET /api/rate_sheets/{id}/rollout?effective_date=...` | Counts the projects that would move and those skipped       |
| `POST /api/rate_sheets/{id}/rollout`                   | Moves them; body `{effective_date, reason, expected_count}` |

Both return `{effective_date, move, closed, cancelled, later_period}`. The
GET date defaults to the revision's `effective_date`; the POST needs one.

- **Who:** an active user with the `admin` claim. `job` and
  `rate_sheet_revise` are not needed. Anyone else gets 403.
- **Target:** `{id}` must be the latest revision of its name, and active,
  because a new row needs an active sheet.
- **Projects:** projects (not proposals) whose `jobs.rate_sheet` is an earlier
  revision of the same name. Only Active projects move. Closed and Cancelled
  projects are skipped and counted. A project whose latest dated row starts on
  or after the date is skipped and counted as `later_period`, because a new
  row would not become its latest.
- **Date and reason:** the rules of any dated row: a valid date, not in the
  future, not before the revision's `effective_date`
  (`hooks.CheckRateSheetPeriodDate` and `hooks.CheckRateSheetPeriodRevision`,
  which the row hook also uses). The reason is required, up to 500
  characters, and written on every new row with `created_by` and
  `updated_by` set to the admin.
- **Each project:** the same steps as a new period on the job details page
  (`changeRateSheetHistory` and `addDatedRateSheetPeriod` in
  `app/routes/job_rate_sheets.go`): the opening row is written if the
  project has none, the dated row is saved through the row hooks, the whole
  history is checked, and `hooks.SyncJobRateSheet` saves the job, so the
  writeback sends it.
- **All or nothing:** one transaction. If the number of projects to move is
  not `expected_count`, the POST refuses with 409 and returns the new counts.
- **Absorb:** refused with 409 while any rate sheet absorb is pending, since
  undo would overwrite a moved project.
- **Jobs editing:** refused while job editing is turned off, as the history
  routes are.

The projects are found with one query. Each moved project then takes about a
dozen small queries and a job save; 400 projects took about 0.4 s in a test
run. That N+1 keeps every row on the validated path.

The rate sheet details page shows admins a notice with **Move projects**
when the sheet is the active latest revision and active projects use an
earlier one. It appears right after an admin activates a revision. The
confirmation shows the counts for the chosen date, asks for the date
(default the revision's effective date) and a reason, and says which time is
repriced (`RateSheetRollout.svelte`).

## Proposed rates

Status: implemented on branch `187` (issue 187), not yet merged.

A project request can propose new rates instead of choosing an active sheet.
The request approver approves them with the request, and the sheet becomes
active in the transaction that creates the job. The decisions of 2026-11-13
in issue 187 are binding: the request approver approves (no new claim, no
second sign-off), sheets made outside a request need no approval, Turbo makes
the duplicate check, the reason is required and kept, and the request form
says that the sheet is fixed for the life of the project and must match the
contract. New rates on an existing project (issue 200) and renaming are out of
scope.

### Data model

Migration `1794610269_proposed_rate_sheets.go` adds three fields to
`rate_sheets`:

| Field      | Type                                 | Notes                                                                                     |
| ---------- | ------------------------------------ | ----------------------------------------------------------------------------------------- |
| `request`  | relation to `job_requests`, optional | Set while the sheet is a draft of that request. Cleared when job creation activates it.   |
| `reason`   | text, max 500                        | Why the client needs the rates. Kept after the request is gone.                           |
| `based_on` | relation to `rate_sheets`, optional  | The active sheet the requester started from; the approver compares with it. Kept as well. |

A sheet with `request` set is a **draft**. A draft is always inactive and
revision 0. The request's `rate_sheet` points at its draft, so job creation
copies it to the job like any other sheet. A request has at most one draft.

The migration also:

- adds `request = ""` to the `rate_sheets` list and view rules, and
  `rate_sheet.request = ""` to those of `rate_sheet_entries`, so the API lists
  no draft and no draft rate;
- refuses `request`, `based_on` and `reason` in the `rate_sheets` create rule,
  and a draft as the sheet in the `rate_sheet_entries` create rule, so the
  collection API cannot make or change a draft;
- leaves drafts out of the `rate_sheets_augmented` view (the Rate Sheets list).

The test fixtures add one draft: request `rsdraftreq00001` (by
`time@test.com`) owns sheet `rsdraftsheet001`, the 2025 Standard Rates with the
Principal at 260. `rate_sheets_import_baseline.csv` has the same row without
the new columns, as for the other test sheets.

### Request routes

`PATCH /api/job_requests/{id}` (and `POST /api/job_requests`) take an optional
`proposed_rates` object: `{name, reason, based_on, entries: [{role, rate,
overtime_rate}]}` (`app/routes/job_request_rate_sheets.go`).

- Only the creator, while the request is editable (Draft, Recalled,
  Rejected), as for every request edit. A change bumps the request version
  and clears approval like any edit, so a review page goes stale.
- Entries follow the new-sheet rules (`validateRateSheetEntryInputs`, shared
  with `POST /api/rate_sheets`): each role once, a whole rate of at least 1,
  an overtime rate of at least 1. `based_on` must be an active sheet that is
  not a draft; a saved draft may keep a source that has since been
  deactivated. The name may not equal any other sheet's name in any letter
  case, as for a rename.
- Turbo refuses rates identical to an active sheet with 409
  `duplicate_rate_sheet` under `proposed_rates`, listing the matches
  (`hooks.CompareWithActiveRateSheets`, the same check as for a new sheet).
  The request form names the sheet and offers **Use** that sheet.
- A save may not both choose a sheet and propose rates. Choosing another
  sheet, or none, deletes the request's draft.
- A request cannot choose another request's draft.
- The reason may be blank in a draft save. Submission, approval and job
  creation require a reason, a rate for every role, and no identical active
  sheet (`hooks.ValidateRequestRateSheetDraft`, called from
  `hooks.ValidateProjectRequestData` in place of the active-sheet check).

The request detail (`requestRelated`) adds `rate_sheet_proposed`,
`rate_sheet_source`, and for each rate `role_id`, `source_rate`,
`source_overtime_rate` and `changed`. These are part of the review token. The
request options leave drafts out.

### Approval

The request approver approves the rates as part of the request. No claim and
no second sign-off are added. For a request with proposed rates, the detail
route replaces the text of the approval checklist question
`rate_sheet_correct` with "The proposed rates are approved for this client."
The key and the stored checklist version stay the same, so no new checklist
version is needed and approvals of other requests are unaffected. The
creation checklist keeps its text.

### Job creation

`createJobFromRequest` first calls `hooks.ActivateRequestRateSheet`: it checks
the draft again, clears `request`, sets `effective_date` to the creation date
(UTC), and activates the sheet with `hooks.ActivateRateSheet`, which the
activate route now uses too (duplicate check, deactivate other revisions,
save). The job is then saved with the active sheet, and the `jobs` hook writes
its opening `job_rate_sheets` row as for any new project. A failure rolls
everything back; the request and its draft stay as they were.

### Returned and deleted requests

Recall and rejection do not touch the draft, so the creator can edit it. A
model hook on `job_requests` delete (`addRateSheetDraftHooks`) deletes the
request's drafts after the request on every path (creator delete, admin
delete of an abandoned request, any other delete). Job creation clears
`request` before deleting the request, so the activated sheet stays.

### Everything else refuses a draft

- `POST /api/rate_sheets/{id}/activate`, `PUT /api/rate_sheet_entries/{id}`
  and `POST /api/rate_sheets/{id}/rename` return `rate_sheet_draft`.
- `POST /api/rate_sheets` refuses a revision of a draft's name.
- Jobs and history rows need an active sheet, so neither accepts a draft.
- Rate sheet absorb leaves drafts out of its data, so a draft is neither a
  source, a target nor a duplicate. `rate_sheets.based_on` is a moved
  reference of a rate sheet absorb.

### Direct creation

Unchanged. `job` holders create sheets and activate them without approval.
Activation already refused a sheet whose rates match an active sheet of
another name (issue 189); proposed sheets, once active, are compared like any
other. Existing sheets are not checked or changed.

### Choices the issue left open

- A request chooses a sheet or proposes rates; choosing a sheet discards the
  draft.
- The requester names the new sheet; it is always revision 0.
- The draft's effective date is the save date, and becomes the activation
  date at job creation.
- The reason may be blank while the request is a draft and is required from
  submission. It is at most 500 characters.
- The proposal is checked for duplicates on every save and again at
  submission, approval and creation.
- The approval statement is replaced by text in code, not by a new checklist
  version, and only in the approval checklist.
- The requester's rates are compared with the source sheet's per role; a role
  the source lacks counts as changed.
- Near matches are not shown to the requester or the approver.

## Migration

Migration `1794579748_job_rate_sheets.go`:

1. Create `job_rate_sheets` and its indexes.
2. Backfill one opening row for every project with a non-empty `rate_sheet`,
   including imported jobs. Leave `reason` empty and `created_by` unset on
   backfilled rows.
3. Projects with no rate sheet get no rows. Their time remains priced only by
   employee defaults, as today, until a `job` holder sets a first sheet.
4. Update the `rate_sheet_revise` claim description: it now governs the history
   rather than `jobs.rate_sheet`.

The migration was first named `1792425372_job_rate_sheets.go`. It was renamed
when the branch was rebased onto `main`, because
`1794352788_claim_descriptions.go` on `main` also sets the `rate_sheet_revise`
description. With the older timestamp, a new database would run that migration
second and put back the old description, while production would run this one
second. Now this one runs last in both.

The test fixtures in `app/test_seed_data` include the collection and the
opening rows the backfill creates for fixture projects: 32 from the branch,
and 3 for the rate sheet absorb fixture jobs (`jobab...`) added on `main`.

In the local database copy checked on 2026-10-19, 6,167 of 7,043 projects (88%)
have no rate sheet, and 4,463 of those have time entries. Almost all are
imported. Only 876 projects get a backfilled row.

4,228 time entries are dated before their project's rate sheet revision took
effect. The opening row covers them, which is why it is exempt from the
revision-date rule.

The backfill produces identical prices: every project has one range covering
all time, using its current sheet. Run the WIP report and a sample of job
summaries before and after the migration and compare the results.

## Tests

In `app/job_rate_sheets_test.go` (routes), `app/routes/job_rate_sheets_pricing_test.go`
(pricing), and `app/hooks/job_rate_sheets_test.go` (hooks):

- **Pricing:** entries before, on, and after a dated row; several dated rows;
  an inactive sheet in history; a role present in one sheet and missing from
  the next; a project with no rows.
- **Row rules:** second opening row; duplicate date; future date; date before
  the revision's `effective_date`; consecutive duplicate sheet; deleting the
  opening row; inactive sheet on create; date-only edit of a row whose sheet is
  now inactive; deleting or re-dating a row so that its neighbours share a
  sheet.
- **Mirror:** `jobs.rate_sheet` follows the latest row after create, edit,
  and delete. Direct writes to it are rejected on existing projects.
- **Permissions:** missing `job` or `rate_sheet_revise`; inactive user; the
  opening row created through job request approval without
  `rate_sheet_revise`; a first sheet set on an existing project with only the
  `job` claim; changing an existing opening row's sheet without
  `rate_sheet_revise` is refused.
- **Migration:** WIP values are identical before and after the backfill.

In `app/rate_sheet_history_absorb_test.go` and `app/rate_sheet_absorb_test.go`
(rate sheet absorb):

- **Absorb:** a sheet that only an older row uses is absorbed, its row
  moves, the job is not touched, and WIP is the same before, after, after
  undo and after commit; undo restores every field of the history; the
  job's current sheet and latest row move together; a name with revisions
  is refused when only a row uses one of them, with the reason in the
  Duplicates view and the candidates; a pair of neighbouring rows left by an
  absorb does not block a new row, and a change that makes a pair is still
  refused; history changes to a moved project are refused while the absorb
  is pending; a request that becomes a job while it is pending carries the
  opening row's undo data.

In `app/rate_sheet_rollout_test.go` (revision rollout):

- **Rollout:** the counts for the default date and other dates; an admin
  without `job` or `rate_sheet_revise` moves exactly the active projects,
  with the right rows, and `jobs.rate_sheet` follows; WIP keeps the old
  price before the date and uses the new revision after it; other projects
  do not change; a second run moves nothing. Refused, with nothing written:
  a user without `admin`, an earlier or inactive revision, a future date or
  one before the revision, no reason, a wrong `expected_count` (409), and a
  pending rate sheet absorb (409).

## Decisions (2026-11-13)

- **Legacy system:** accepted. The writeback sends only the current sheet
  (`jobs.rate_sheet`). If the legacy system prices time itself, it prices
  all time with that sheet. Turbo is the record of which sheet applies from
  which date. No change is planned.
- **Revision rollout:** implemented; see Revision rollout above. An `admin`
  claim holder can deliberately move every project that uses a rate sheet to
  its latest revision. Rules:
  - Only `admin` claim holders can do it. All other users move projects one
    at a time from the job details page, as today.
  - The app offers it when a new revision is activated, and it is also
    available later as its own action.
  - It runs only after a confirmation that shows how many projects it will
    change, and how many it skips and why.
  - Each moved project gets a dated row on the new revision, so earlier time
    keeps its price.
  - The admin picks the start date. The default is the revision's effective
    date.
  - Only active projects move. Closed and cancelled projects are skipped, and
    the confirmation shows the count.
- **Invoiced periods:** deferred. The app has no record of what a client was
  invoiced, so it cannot know a project's last invoice date. Blocking a dated
  row earlier than the last invoice waits for client invoices (issues 151 and
  152, `WIP_invoiced_amount_tracking.md`). Until then, the warning and the
  required reason are the only safeguards.

## Open questions

- **Overtime:** `rate_sheet_entries.overtime_rate` is not used for pricing.
  This proposal does not change that.
