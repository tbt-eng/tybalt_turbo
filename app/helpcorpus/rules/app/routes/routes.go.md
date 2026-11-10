<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"fdb7535177465c3cb950c58959b5bb926a5ee3b46536493e786e84544e5ff6a3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expenses editing is off

Lines 14-31

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message held in ErrExpensesEditingDisabled. Other collections are not checked by this rule.

## Time editing blocked when the time feature flag is off

Lines 33-43

When the time editing feature flag is off, Turbo refuses the action with the message held in ErrTimeEditingDisabled (requireTimeEditing).

## Reports need the report claim

Lines 45-54

A person without the report claim cannot view a report. Turbo answers "you are not authorized to view this report".

## KPI reports need the KPI claim

Lines 56-65

A person without the kpi claim cannot view KPI reports. Turbo answers "you are not authorized to view KPI reports".

## Startup fails if expenditure kinds config is invalid

Lines 88-91

When the server starts, Turbo validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If it is invalid, the routes are not registered and startup returns the error.

## Version endpoint is public

Lines 93-105

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Time sheet routes require sign-in

Lines 107-131

All /api/time_sheets routes require a signed-in user. They include bundle, unbundle (recall), copy to next week, approve, reject, commit, uncommit, and tracking lists. The commit route is open to committers (the commit claim) through the shared commit handler, which checks the claim.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

The /api/export_legacy routes for time sheets, jobs and expenses are not behind user sign-in. Their handlers do their own authentication with machine_secrets.

## Expense routes require sign-in

Lines 138-168

All /api/expenses routes require a signed-in user. They cover create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, uncommit requests, attachment audit, settle, tracking and reimbursements. The file does not gate them by the expenses editing flag.

## Expense uncommit request routes

Lines 170-176

The /api/expense_uncommit_requests routes need sign-in. They list requests and reasons, approve, decline and cancel.

## Time amendment commit route

Lines 178-180

A signed-in user can call commit on a time amendment. The shared commit handler decides who may commit.

## Other signed-in user routes

Lines 182-193

Copying a time entry to tomorrow, reading user defaults and reading nav badges each need sign-in.

## Admin profile routes require sign-in

Lines 195-205

The /api/admin_profiles routes need sign-in. They list identities, save with claims, save limited, set active and clear authorized providers. Their handlers decide who may use each.

## Purchase order read routes

Lines 207-219

Reading purchase orders (pending, visible, search, approvers, approval limits, second approvers, groups) needs sign-in only. The expenses editing flag is not checked here.

## Purchase order changes blocked when expenses editing is off

Lines 221-237

Approve, reject, cancel, close, make cumulative, revise and group create or preview on purchase orders need sign-in. They also fail with the ErrExpensesEditingDisabled message when expenses editing is off.

## Legacy purchase order routes are not gated by the flag

Lines 239-243

Legacy purchase order edit, create and update routes need sign-in. This file does not check the expenses editing flag for them.

## Client, contact and vendor routes

Lines 245-266

Client, client contact and vendor routes need sign-in. They include absorbing one record into another and undoing an absorb. Vendor routes are not gated by the expenses editing flag here.

## Job routes require sign-in

Lines 270-301

All /api/jobs routes need sign-in. They cover create, upsert, details, time, expense and PO summaries, WIP, stale and unused jobs, set status, set number, change client, adopt subjobs, close and validate proposal. Handlers decide who may do each.

## WIP report modes

Lines 303-307

The WIP report has three modes: my, branch and division. Each needs sign-in.

## Reports routes need the report claim

Lines 309-328

All /api/reports routes need sign-in and the report claim. They include payroll and weekly time, expense and receipts reports, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI report route needs the KPI claim

Lines 330-338

The employee branch hours KPI report needs sign-in and the kpi claim.

## Branch hours report access depends on the requested branch

Lines 340-348

The /api/branch_hours report needs sign-in. Before running, Turbo calls requireBranchHoursAccess with the branch from the query string. The rule for who may view is in that function, not in this file.

## Stats, machine secrets and attachment audit routes

Lines 350-368

Stats, machine secrets (list, create) and attachment audit routes need sign-in. Handlers enforce admin-only access where needed.

## Currency, rate sheet and claims routes

Lines 370-401

Currency routes (list, reload rates, initialize backfill, delete), rate sheet routes (create, activate, deactivate), rate sheet entry update and claims routes need sign-in. Per the comments, rate sheets need the job claim, entry updates and claim details need the admin claim, and handlers enforce this.

## Work records routes

Lines 393-396

Listing work records and reading a work record's details need sign-in.
