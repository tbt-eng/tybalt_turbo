# Branch Hours Report

## Purpose

This read-only report shows how one branch's hours moved between branches, by
project and person. Branch managers use it to see their inter-office effort:
the hours other branches' staff worked for their branch, and the hours their
staff worked for other branches.

## Access

- Branch managers (`branches.manager`), and users with the `kpi` or `admin`
  claim, can view the report. This is the same access as the branch and
  division WIP reports (`canViewBranchReports` in `app/routes/wip_reports.go`).
- Any permitted user can select any branch.
- The Branch Hours navigation item is under Reports.

## Page

- `/reports/branch-hours` displays this report.
- The branch defaults to the branch the user manages, or else the user's own
  default branch.
- The period defaults to January 1 of the current year through today. Start
  and end dates are inclusive. Changing a control reloads the report.
- The branch, dates and selected tab are kept in the URL hash, so a report can
  be shared as a link.
- Three tabs split the hours by direction, each labelled with its total:
  - **Received**: other branches' staff working for the selected branch.
  - **Supplied**: the selected branch's staff working for other branches.
  - **Internal**: the selected branch's staff working for the selected branch.
- Each tab lists one row per project and work branch, largest first. Expanding
  a row lists the people who recorded its hours. Time with no job appears as
  "No project".
- A search box filters projects by project number, description, client, branch
  or person.
- Download CSV exports every direction as flat rows.

## Data Rules

The report uses the same rules as the KPI Employee Branch Hours report
(`kpi_employee_branch_hours.md`), so its totals reconcile with it:

- The source is `time_entries`. The report includes only `R` and `RT` time
  types, does not include `time_amendments`, and does not require a committed
  time sheet.
- The staff branch is the employee's current `admin_profiles.default_branch`.
  The work branch is the branch on the time entry.
- A time entry is included when its staff branch or its work branch is the
  selected branch. Its direction is `internal` when both are the selected
  branch, `supplied` when only the staff branch is, and `received` otherwise.

For a selected branch, Internal plus Received equals the KPI report's branch
and NoJob columns for that branch. Internal plus Supplied equals the KPI totals
of the employees whose default branch is that branch.

## API

```text
GET /api/branch_hours
```

Parameters:

- `branch`: required branch ID.
- `start_date`: required date in `YYYY-MM-DD` format.
- `end_date`: required date in `YYYY-MM-DD` format.
- `format`: optional `json` or `csv`. The default is `json`.

The JSON response has the branch name and one row per job, employee and work
branch:

```json
{
  "branch": "Kitchener-Waterloo",
  "rows": [
    {
      "direction": "received",
      "jobId": "...",
      "jobNumber": "26-0412",
      "jobDescription": "...",
      "client": "...",
      "workBranch": "Kitchener-Waterloo",
      "uid": "...",
      "payrollId": "...",
      "surname": "...",
      "givenName": "...",
      "staffBranch": "Thunder Bay",
      "hours": 12.5
    }
  ]
}
```

The CSV has the columns `direction`, `jobNumber`, `jobDescription`, `client`,
`workBranch`, `payrollId`, `surname`, `givenName`, `staffBranch` and `hours`.
No-job rows have an empty `jobNumber`.
