<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"82b8b95945fd866df22add49fce25eac23646019775ed616df6713bfea145b92","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Purchase order changes blocked when expenses editing is off

Lines 14-31

Asked as: can't approve purchase order; PO buttons greyed out; expenses editing turned off; can't cancel PO; purchase order locked

Approving, rejecting, cancelling, closing, converting to cumulative, revising a purchase order and creating or previewing purchase order groups all need the expenses editing feature switched on. When it is off, Turbo refuses with the message of ErrExpensesEditingDisabled. The check covers the "expenses", "purchase_orders" and "vendors" collections in requireExpensesEditing.

## Time editing refused when the time feature is off

Lines 33-43

Asked as: can't edit time; time entry locked; time editing switched off; timesheet changes refused

The helper requireTimeEditing refuses with the message of ErrTimeEditingDisabled when the time editing feature flag is off. The helper is defined here, but this file does not call it on any route.

## Reports need the report claim

Lines 45-54

Asked as: can't open payroll report; not authorized to view report; reports access denied; no access to reports

Every request under /api/reports needs the report claim. Without it Turbo refuses with "you are not authorized to view this report". This covers payroll, weekly time, expense and receipts reports, payables spreadsheets, branch mismatches and active jobs.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI reports; not authorized KPI; employee branch hours missing

Every request under /api/kpi/reports needs the kpi claim. Without it Turbo refuses with "you are not authorized to view KPI reports". The employee branch hours report is under this group.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; where is help chat; chat assistant url; help page missing

Turbo registers the Help chat at /api/help. The chat loads its model and its corpus on first use, not at start-up.

## Turbo refuses to start with a bad expenditure kinds setting

Lines 88-91

Asked as: server won't start; expenditure kinds error; Turbo fails to start; startup config error

When Turbo starts to serve, it checks the expenditure kinds configuration (ValidateExpenditureKindsConfig). If the check fails, the server does not start serving.

## Version endpoint is public

Lines 93-105

Asked as: what version is this; app build number; git commit of turbo; version without login

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Time sheet routes require sign-in

Lines 107-131

Asked as: can't submit my week; recall timesheet; timesheet sign in error; can't approve timesheet; bundle my time

All /api/time_sheets routes need a signed-in user. They cover bundling (submitting) a week, unbundling (recalling), copying entries to next week, details, approve, reject, commit, uncommit, reviewers, approver, tallies and tracking. The commit comment says committers commit time sheets; the claim check happens in the handler, not in this file.

## Time sheet tallies come in four kinds

Lines 120-123

Asked as: timesheet counts; pending timesheets total; how many approved sheets; my time sheet tally

Tallies are available for the user's own time sheets (uid), pending approval (approver), approved (approver), and shared for review (reviewer). Handler parameters 1,0 mean pending and 0,1 mean approved.

## Time tracking and missing-sheet exceptions routes

Lines 124-131

Asked as: add missing timesheet exception; ignore missing timesheet; people who haven't submitted time; time tracking counts

Time tracking routes show counts, a week list, missing, ignored and not-expected time sheets. A user can add or remove a missing time sheet exception for a user and week. The comment says these are for report holders; the claim check is in the handlers.

## Legacy export endpoints use machine secrets

Lines 133-136

Asked as: legacy export without login; machine secret export; export jobs legacy; old export api

The export_legacy routes for time sheets, jobs and expenses do not need a user sign-in. They use custom authentication through machine_secrets in the handlers.

## Expense routes and their actions

Lines 138-157

Asked as: can't edit expense; submit expense claim; attach receipt to expense; approve expense route; copy expense to tomorrow

Signed-in users reach expense create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, uncommit requests, list, details, attachment, pending and approved. Attachment hash audit, replace and mark-missing are also available. Permissions are checked in the handlers.

## Expense tracking, commit queue and settlement routes

Lines 158-165

Asked as: settle expenses; unsettled expenses list; commit queue expenses; clear settlement

Signed-in users reach expense tracking counts, the org-wide commit queue, unsettled and settled lists, settle and clear settlement. Tracking by committed week is also available.

## Expense reimbursements by pay period

Lines 166-168

Asked as: expenses repaid on payroll; reimbursed expenses by period; pay period expenses; what was paid back

The reimbursement routes list the caller's own expenses repaid on payroll, grouped by pay period, and each period's expenses.

## Expense uncommit request routes

Lines 170-176

Asked as: expense uncommit requests; decline uncommit request; cancel uncommit request; uncommit reasons

Signed-in users can list expense uncommit requests and reasons, approve requests, and decline or cancel one request by id.

## Time amendments commit route

Lines 178-180

Asked as: commit time amendment; amendment commit button; time amendment not committing

POST /api/time_amendments/{id}/commit uses the shared commit handler for time_amendments and needs a signed-in user.

## Copy a time entry to tomorrow

Lines 182-185

Asked as: copy time entry tomorrow; duplicate time entry next day; repeat entry tomorrow

A signed-in user can copy a time entry to tomorrow with POST /api/time_entries/{id}/copy_to_tomorrow.

## User defaults and nav badges routes

Lines 187-193

Asked as: my default settings; nav badge counts; navigation notification numbers; user defaults not loading

Signed-in users can get their defaults from /api/users/defaults and navigation badges from /api/nav/badges.

## Admin profile routes

Lines 195-205

Asked as: edit admin profile; set user active; clear authorized provider; identity save error; profile claims save

Signed-in users reach admin profile identity lists, active toggle targets, save with claims, save limited, set active, identity save and clearing an authorized provider. Who may use each is decided in the handlers.

## Purchase order read routes are not gated by the expenses flag

Lines 207-219

Asked as: view purchase orders; pending POs list; search purchase orders; see approvers list; PO groups listing

Signed-in users can read pending, visible, searchable purchase orders, approvers, approval limits, second approvers and groups. These routes do not check the expenses editing flag.

## Purchase order changes need expenses editing on

Lines 221-237

Asked as: PO approve button disabled; can't make PO cumulative; revise purchase order blocked; group preview refused; expenses flag off PO

The PO approve, reject, cancel, close, make_cumulative, revise and group create/preview routes need the expenses editing flag. When it is off they are refused before the handler runs.

## Legacy purchase order routes

Lines 239-243

Asked as: legacy purchase order edit; old PO form; create legacy PO; update old purchase order

Signed-in users can fetch a legacy purchase order for edit, create one and update one. These routes do not pass the expenses editing check in this file.

## Client, contact and vendor routes

Lines 245-266

Asked as: absorb client; undo absorb vendor; client contacts list; vendor details; business development leads

Signed-in users reach client list, details, notes, time divisions and business development leads. Absorb and undo absorb exist for clients, client contacts and vendors, and vendor list and details are available. The handlers decide who may absorb.

## Job routes

Lines 270-301

Asked as: create job; job status change; change job client; close job; stale jobs list; job WIP summary

Signed-in users reach job create, upsert, details, notes, time, staff, division, expense and PO summaries, WIP, lists, unused and stale jobs. Set status, set number, change client, adopt subjobs, close and validate proposal are also available. Query parameters filter the time summary and entries.

## WIP report modes

Lines 303-307

Asked as: WIP by division; my WIP report; branch WIP view; work in progress report modes

WIP reports are available in three modes, my, branch and division, for signed-in users.

## Report routes need sign-in and the report claim

Lines 309-328

Asked as: payroll report access; can't see payables spreadsheet; branch mismatch report; active jobs report denied

Every route under /api/reports needs a signed-in user. It then needs the report claim. A user without the claim gets 'you are not authorized to view this report'. The check covers payroll, weekly time, expense and receipts reports, the payables spreadsheets, the time entry branch mismatches report and the active jobs report.

## KPI report routes need sign-in and the kpi claim

Lines 330-338

Asked as: KPI reports locked; kpi access needed; employee hours KPI blocked

Routes under /api/kpi/reports need a signed-in user. They also need the kpi claim. A user without the claim gets 'you are not authorized to view KPI reports'. The employee branch hours report is under this group.

## Branch hours report checks branch access

Lines 340-348

Asked as: branch hours access denied; branch hours report blocked; can't view branch hours

GET /api/branch_hours checks access for the branch in the "branch" query parameter using requireBranchHoursAccess. The rule is in another function.

## Stats and machine secrets routes

Lines 350-359

Asked as: system stats page; create machine secret; list machine secrets

Signed-in users reach /api/stats and machine secret list and create. The comment says machine secrets are admin only; the check is in the handlers.

## Attachment audit routes

Lines 361-368

Asked as: attachment audit run; orphaned attachments delete; missing attachments csv; refresh attachment audit

Signed-in users reach attachment audit targets, a target's run, refresh, delete orphans, and CSV downloads of the missing and orphaned reports.

## Currency routes

Lines 370-376

Asked as: currency rates reload; initialize currencies; delete currency; exchange rates backfill

Signed-in users reach currency list, initialization status, reload rates, initialize backfill and delete.

## Rate sheet routes

Lines 378-391

Asked as: rate sheet duplicates; compare rate sheets; activate rate sheet; rename rate sheet; former rate sheet names

Signed-in users reach rate sheet create, compare, activate, deactivate, rename and former names. The comment says the job claim is required. Duplicates, absorb candidates, absorb and undo absorb are also there, for admin or absorb with rate_sheet_revise.

## Claims and work record routes

Lines 393-404

Asked as: bulk assign claims; assignable users list; work records list; claim details access

Signed-in users reach the claims list, assignable users, bulk assign and claim details. The comment says the detail endpoint requires the admin claim. Work records list and details are also available.

## Rate sheet entry update route

Lines 406-409

Asked as: edit rate sheet entry; update rate entry; change rate sheet line

PUT /api/rate_sheet_entries/{id} updates a rate sheet entry for a signed-in user. The comment says the admin claim is required.
