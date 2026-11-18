<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"f31025d7cb450cc45627ebcee35cf41f0c79562d900ada61cf2ce3f5588e7cea","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses feature disabled; cannot edit expenses; purchase orders blocked; vendor changes refused; expenses editing turned off

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors. The refusal message is the one in ErrExpensesEditingDisabled. Other collections are not affected by this check.

## Time editing is refused when the time feature flag is off

Lines 33-43

Asked as: time editing disabled; cannot change timesheet; time changes refused; time feature off; can't edit my hours

When the time editing feature flag is off, Turbo refuses time editing actions that use this check. It gives the message in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: no access to reports; report permission missing; can't open payroll report; not authorized to view report; reports locked for me

A person without the report claim cannot use any /api/reports route. Turbo refuses with "you are not authorized to view this report". The claim is checked for every report in the group, such as payroll, weekly, receipts, payables spreadsheets, branch mismatches and active jobs.

## KPI reports need the kpi claim

Lines 56-65

Asked as: kpi report access; can't see kpi; kpi permission missing; not authorized kpi reports; branch hours kpi locked

A person without the kpi claim cannot use the KPI reports (/api/kpi/reports, such as employee branch hours). Turbo refuses with "you are not authorized to view KPI reports".

## Startup stops if the expenditure kinds config is invalid

Lines 88-91

Asked as: server won't start; startup fails expenditure kinds; expenditure config invalid; app crashes on boot; expenditure kinds error

When Turbo starts, it validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Most API routes require a signed-in user

Lines 93-136

Asked as: anonymous api calls; api needs login; not signed in error; unauthenticated requests refused; sign in required api

The time sheet, expense, expense uncommit request, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, report, stats, machine secret, attachment audit, currency, rate sheet, claim, work record and rate sheet entry routes all require a signed-in user. The version endpoint (/api/version) is public. The legacy export routes do not use user sign-in and use their own authorization (machine secrets).

## Version endpoint is public

Lines 93-105

Asked as: version info without login; api version check; which build is running; app version endpoint; build and commit details

The /api/version route needs no sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Time sheet and expense routes require sign-in

Lines 107-206

Asked as: anonymous timesheet access; timesheets need login; expense api unauthorized; no login for expenses; refused without signing in

Every route under /api/time_sheets, /api/expenses, /api/expense_uncommit_requests, /api/time_amendments, /api/time_entries, /api/users, /api/nav and /api/admin_profiles needs a signed-in user (the users collection). Turbo refuses anonymous callers.

## Time sheet approve, reject, commit and uncommit routes

Lines 113-117

Asked as: approve timesheet; reject timesheet; commit timesheet; uncommit timesheet; who can approve sheets

Time sheets can be approved, rejected, committed and uncommitted through /api/time_sheets/{id}/approve, /reject, /commit and /uncommit. Time amendments can be committed through their own commit route. The code in this file does not show who passes each check.

## Time sheet tracking routes are for report holders

Lines 124-131

Asked as: missing timesheets list; tracking counts timesheets; ignored timesheets; add missing exception; weekly timesheet tracking

The time sheet tracking routes are meant for report holders. They cover tracking counts, the weekly list, missing, ignored and not expected time sheets. They also add or remove a missing time sheet exception for a user and a week.

## Legacy export routes use machine secrets, not sign-in

Lines 133-136

Asked as: export legacy without login; legacy export access; machine secret export; export script authentication; legacy data export failing

The legacy export routes for time sheets, jobs and expenses (/api/export_legacy/...) are not behind the signed-in user check. They use their own authentication through machine_secrets.

## Expense workflow routes

Lines 140-176

Asked as: submit expense; recall expense; copy expense to tomorrow; request expense uncommit; approve expense uncommit

An expense has routes to submit, recall, copy to tomorrow, approve, reject, commit and uncommit. Staff can also request an uncommit of an expense (/uncommit_requests). Expense uncommit requests can be approved, declined or cancelled. The code here does not state who may do each step.

## Expense settlement and attachment routes

Lines 153-164

Asked as: settle expense; clear settlement; unsettled expenses list; attachment hash replace; mark receipt missing

Expenses have routes to settle, to clear a settlement, and to list unsettled and settled expenses. They also have routes to audit or replace an attachment hash and to mark an attachment missing. This file does not show who may use them.

## Reimbursements list shows the caller's own expenses

Lines 166-168

Asked as: my repaid expenses; payroll reimbursements list; my reimbursements by pay period; expenses repaid on payroll; reimbursed expenses view

The reimbursements routes list the caller's own expenses that are repaid on payroll. They group them by pay period (payPeriodEnding).

## Purchase order approve, reject, cancel, close, convert, revise and group routes need expenses editing

Lines 207-243

Asked as: purchase order approval blocked; cannot close purchase order; revise po disabled; cancel po refused; po group create blocked

Approve, reject, cancel, close, make cumulative, revise and group preview/create on a purchase order are refused when the expenses editing flag is off. Listing and viewing purchase orders (pending, visible, search, approvers, groups) are not gated by this flag in this file. The legacy purchase order routes (create, edit and update) are also not gated here.

## Absorb and undo absorb routes exist for clients, contacts, vendors and rate sheets

Lines 253-267

Asked as: absorb duplicate client; undo merge contact; absorb vendor; rate sheet merge; merge rate sheets back

Turbo offers an absorb action and an undo absorb action for clients, client contacts, vendors and rate sheets. Rate sheet absorb is for admins, or for those who have the rate_sheet_revise claim, according to the comment. Rate sheet rollout to the latest revision is an admin action.

## Jobs routes require sign-in

Lines 271-306

Asked as: job create refused; jobs api access; close job; change job client; adopt subjobs

All /api/jobs routes need a signed-in user. They cover job create, upsert, details, time, expense and PO summaries, WIP, rate sheets, set-status, set-number, change-client, adopt-subjobs, close and validate-proposal. Who may do each action is not stated in this file.

## WIP report modes

Lines 308-312

Asked as: my wip report; branch wip view; division wip; work in progress report; wip by branch

The WIP report has three routes: /api/wip/my, /api/wip/branch and /api/wip/division. They need only sign-in here.

## Report routes require the report claim

Lines 314-333

Asked as: payroll time report denied; weekly time report access; receipts report permission; payables spreadsheet denied; active jobs report locked

All /api/reports routes (payroll time, weekly time, expenses, receipts, payables spreadsheets, branch mismatches, active jobs) check the report claim after sign-in. Without it Turbo says "you are not authorized to view this report".

## KPI report route requires the kpi claim

Lines 335-343

Asked as: branch hours kpi denied; employee hours kpi; kpi report not allowed; kpi access refused; employee branch kpi

The employee branch hours KPI report (/api/kpi/reports/employee_branch_hours) needs sign-in and the kpi claim. Without it Turbo says "you are not authorized to view KPI reports".

## Branch hours report access depends on the requested branch

Lines 345-353

Asked as: branch hours report denied; branch query access; hours report for branch; cannot see branch hours; branch report permission

The branch hours report runs a check (requireBranchHoursAccess) using the signed-in user and the branch query parameter before it runs. The result of that check is defined in another file.

## Admin, attachment audit, currency and claim routes

Lines 355-418

Asked as: machine secrets admin; claim details refused; rate sheet entry update denied; attachment audit access; currency routes login

The stats, machine secrets, attachment audit, currencies, claims, work records and rate sheet entry routes all need sign-in. Comments say machine secrets are admin only, claim details need the admin claim and rate sheet entry updates need the admin claim. This file does not show the checks.

## Rate sheet routes and who may use them

Lines 383-400

Asked as: create rate sheet; compare rate sheets; activate rate sheet; rollout rate sheet revision; rate sheet job claim

Rate sheet create, compare, activate, deactivate, rename and jobs routes need sign-in. Comments say the job claim is required for these. Rollout to the latest revision is for admins. Absorbing duplicate rate sheets is for admins, or for those with the rate_sheet_revise claim.
