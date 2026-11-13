# WIP performance

Measured on 2026-09-29 with an Apple M1 Max and an isolated copy of the local
production database dump. The source database was not changed. The baseline is
commit `37dcddb5`; the optimized copy has migration
`1790712255_wip_lookup_indexes.go` applied.

The snapshot contains 16,220 jobs, 448,608 time entries, 37,161 expenses, and
7,347 purchase orders. Of 6,060 active non-proposal jobs, only 542 have a positive
project value and appear in WIP reports.

## Database and response preparation

These are medians of five sequential runs per case, using the application's Go
SQLite driver, dbx row decoding, report construction, and JSON encoding. The
report date is 2026-09-29. Runs exclude HTTP authentication, network transfer,
and browser work. They use normal operating-system file caching; no cold-disk
claim is made. No other tests ran during these measurements.

| Report scope | Rows | Before | After | Speedup |
| --- | ---: | ---: | ---: | ---: |
| My WIP, current user (two excluded jobs) | 0 | 64.33 ms | 2.64 ms | 24.4× |
| My WIP, largest manager by valued jobs | 112 | 218.04 ms | 10.42 ms | 20.9× |
| Branch WIP, all branches | 542 | 923.92 ms | 91.37 ms | 10.1× |
| Branch WIP, largest branch by valued jobs | 310 | 863.10 ms | 54.81 ms | 15.7× |
| Division WIP, largest allocated division, recorded hours required | 130 | 1,060.12 ms | 41.99 ms | 25.2× |
| Division WIP, same division, all allocated jobs | 362 | 496.84 ms | 69.08 ms | 7.2× |

The first all-branch run in each series took 1,061.65 ms before and 127.16 ms
after. This is a first-query measurement, not a cold-disk measurement.

The main changes are:

- Count the report scope separately, then price only jobs with a positive
  project value. An empty priced scope returns without the pricing query.
- Materialize the selected job set once. Start labour lookups from those jobs
  so SQLite uses the job index instead of scanning all time entries.
- Read the selected committed expenses once and reuse them for expense totals
  and PO spend. Expenses assigned to another job still reduce the selected PO.
- Add an expense index on `(job, date)` and a partial time-entry index on
  `(job, division, date)` for `hours > 0`. These support expense lookups and the
  Division WIP recorded-hours filter.

The two indexes use approximately 20.6 MiB on this snapshot. SQLite must also
maintain them when the indexed records change. The migration uses PocketBase
collection APIs, and its rollback removes only these indexes.

No response cache is added. Each request reads current data.

## Browser rendering

Closed help popups now build their contents on first use. Once opened, their
contents remain mounted so loaded details can be reused.

These medians use three runs in Chrome against the Vite development harness,
with the same mocked report response before and after. Timing starts when the
mock response is supplied and ends after rows mount and two animation frames.
It excludes server and network latency. These results indicate render cost;
they are not production page-load timings.

| Report | Before | After | DOM elements before | DOM elements after |
| --- | ---: | ---: | ---: | ---: |
| My WIP, 112 rows | 65.7 ms | 53.0 ms | 5,276 | 2,980 |
| Branch WIP, 542 rows | 271.4 ms | 215.7 ms | 26,574 | 14,479 |
| Division WIP, 130 rows | 69.8 ms | 54.3 ms | 6,529 | 3,530 |

## Correctness and repeat checks

The old and new report results were compared across all 179 manager, branch,
and division scopes in the snapshot, including both division hours settings.
Every report matched exactly: values, flags, row order, and exclusion counts.
The six timing cases also matched after the real migration was applied to a
fresh copy. Migration up, down, and up again were checked against physical
indexes and PocketBase collection metadata.

The permanent benchmark opens its database read-only. Prepare an isolated
snapshot with the migration applied, then run from `app`:

```sh
WIP_BENCHMARK_DB=/absolute/path/to/snapshot/data.db \
WIP_BENCHMARK_AS_OF=2026-09-29 \
go test ./routes -run '^$' -bench '^BenchmarkWIPReports$' -benchtime=5x -count=1
```

It selects the largest manager, branch, and allocated division from that
snapshot. Its empty-manager case uses an absent user ID, rather than the
current-user case in the table. It does not print customer identifiers.
Go's benchmark output reports mean time per operation; the table above uses
the separate five-run comparison medians.

Start the updated backend with the migration applied to use the new queries
and indexes. The frontend change also requires the updated UI build.

## Rate sheet history

Measured on 2026-10-19 with an Apple M1 Max, after time pricing in
`job_priced_time_entries.sql` began joining each time entry to its rate sheet
period (see `job_rate_sheet_history.md`). The snapshot is the local production
database dump of 2026-10-16, copied twice with SQLite's backup command from a
read-only connection. The source database was not changed.

- **Before:** `main` at `d71196b6`, with its migrations applied.
- **After:** `rate_sheet_history` at `d2160b46`, with
  `1794579748_job_rate_sheets.go` also applied. The backfill wrote 876 opening
  rows, one for each project with a rate sheet.

The snapshot contains 16,352 jobs, 453,979 time entries, 37,802 expenses, and
7,679 purchase orders. Of 6,110 active non-proposal jobs, 601 have a positive
project value.

Each case ran `BenchmarkWIPReports` one operation at a time, five times per
round, in four rounds that alternated before and after. The table shows the
median of those 20 runs. The report date is 2026-10-16. Rows matched before and
after for every case. `My WIP, no user` returns before the pricing query, so its
difference is measurement noise.

| Report scope | Rows | Before | After | Change |
| --- | ---: | ---: | ---: | ---: |
| My WIP, no user | 0 | 2.42 ms | 2.59 ms | +7% |
| My WIP, largest manager by valued jobs | 126 | 11.16 ms | 12.53 ms | +12% |
| Branch WIP, all branches | 601 | 99.89 ms | 111.63 ms | +12% |
| Branch WIP, largest branch by valued jobs | 350 | 60.84 ms | 67.15 ms | +10% |
| Division WIP, largest allocated division, recorded hours required | 149 | 46.93 ms | 51.55 ms | +10% |
| Division WIP, same division, all allocated jobs | 395 | 75.67 ms | 83.35 ms | +10% |

Every priced scope is about 10% slower, or 5 to 12 ms. Diagnostic variants
showed where the time goes:

- Removing the date range test from the period join did not change the
  timings, so the cost is not the date comparison.
- Removing `LEAD()` did not change them either, so the cost is not the window.
- Starting the history lookup from `selected_jobs` with the
  `(job, effective_date)` index did not help.

The cost is the lookup, once per time entry, into the automatic index SQLite
builds on the materialized `rate_periods` result.

### Query change, applied

The pricing query now starts the time-entry lookup from the rate periods
instead of the jobs. Each period reads its own date range from the `(job, date)` time-entry index,
and no per-entry period lookup remains. For this to keep every entry, each
selected job needs an opening period, including jobs without a rate sheet.
The rate sheet is null for those jobs, and for a gap before the first dated
row if a history has no opening row. These were the changes to
`job_priced_time_entries.sql`:

```sql
  UNION ALL
  SELECT j.id, CASE WHEN NOT EXISTS (SELECT 1 FROM job_rate_sheets jrs WHERE jrs.job = j.id)
    THEN NULLIF(j.rate_sheet, '') END, ''
  FROM selected_jobs j
  WHERE NOT EXISTS (
    SELECT 1 FROM job_rate_sheets jrs WHERE jrs.job = j.id AND jrs.effective_date = ''
  )
...
  FROM rate_periods p
  JOIN selected_jobs j ON j.id = p.job
  CROSS JOIN time_entries te ON te.job = p.job
    AND te.date >= p.starts AND (p.ends IS NULL OR te.date < p.ends)
```

No new index is needed. Measured in the same interleaved series:

| Report scope | Before | After the change | Change |
| --- | ---: | ---: | ---: |
| My WIP, largest manager by valued jobs | 11.16 ms | 12.10 ms | +8% |
| Branch WIP, all branches | 99.89 ms | 102.32 ms | +2% |
| Branch WIP, largest branch by valued jobs | 60.84 ms | 60.81 ms | 0% |
| Division WIP, largest allocated division, recorded hours required | 46.93 ms | 48.16 ms | +3% |
| Division WIP, same division, all allocated jobs | 75.67 ms | 77.40 ms | +2% |

On the migrated snapshot, the changed query returned exactly the same report
rows as the previous one for 316 scopes: every manager, branch, and allocated
division, plus all branches, each with recorded hours required and not
required. This snapshot has only opening rows.

To cover dated periods, a second copy of the local database was given 308
dated periods on 200 projects with time: one at the median date of each
project's time on another sheet, and for about half, a later one back on the
original sheet. Over every project, both queries priced the same 288,409
entries identically: date, user, role, hours, sheet, sheet rate, default rate,
and value.
