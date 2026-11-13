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

## Decisions (2026-11-13)

- **Legacy system:** accepted. The writeback sends only the current sheet
  (`jobs.rate_sheet`). If the legacy system prices time itself, it prices
  all time with that sheet. Turbo is the record of which sheet applies from
  which date. No change is planned.
- **Revision rollout:** to be built, not yet implemented. An `admin` claim
  holder can deliberately move every project that uses a rate sheet to its
  latest revision. Rules:
  - Only `admin` claim holders can do it. All other users move projects one
    at a time from the job details page, as today.
  - The app offers it when a new revision is activated, and it is also
    available later as its own action.
  - It runs only after a confirmation that shows how many projects it will
    change.
  - Each moved project gets a dated row on the new revision, so earlier time
    keeps its price.
- **Invoiced periods:** deferred. The app has no record of what a client was
  invoiced, so it cannot know a project's last invoice date. Blocking a dated
  row earlier than the last invoice waits for client invoices (issues 151 and
  152, `WIP_invoiced_amount_tracking.md`). Until then, the warning and the
  required reason are the only safeguards.

## Open questions

- **Revision rollout details:** the start date of the new rows (the
  revision's effective date, or a date the admin picks), and whether closed
  and cancelled projects are moved too.
- **Overtime:** `rate_sheet_entries.overtime_rate` is not used for pricing.
  This proposal does not change that.
