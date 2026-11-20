<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"7a28b67f2ca1e785c9a0bbead40838735ed897b6fa550391b8c6e72d86eaa377","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expenses editing is off

Lines 14-31

Asked as: expenses editing disabled; can't edit expenses; expense changes blocked; purchase orders locked; vendor edits forbidden

When the expenses editing feature flag is off, Turbo refuses actions on the expenses, purchase_orders and vendors collections with a forbidden error. The message is the one in ErrExpensesEditingDisabled. Other collections are not checked by this gate.

## Time editing refused when the time feature flag is off

Lines 33-43

Asked as: can't edit time; time editing off; time entries locked; time changes forbidden; timesheet edits blocked

When the time editing feature flag is off, Turbo refuses with a forbidden error. The message is the one in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: can't open reports; not authorized to view report; report access denied; no report permission; payroll report forbidden

A person without the report claim cannot use any /api/reports endpoint. Turbo answers 'you are not authorized to view this report'. The endpoints cover payroll, weekly time, expense and receipts reports, payables spreadsheets, branch mismatches and active jobs.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI reports denied; can't see KPI; kpi access missing; not authorized KPI; KPI report permission

A person without the kpi claim cannot use the KPI reports under /api/kpi/reports, such as the employee branch hours report. Turbo answers 'you are not authorized to view KPI reports'.

## Help chat route

Lines 84-85

Asked as: help chat not working; help chat slow; help chat first load; help bot startup; assistant chat loading

The Help chat is served at /api/help. It loads its model and corpus on first use.

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds error; startup config check failed; expenditure config invalid; app fails to boot

When the server starts, Turbo validates the expenditure kinds configuration. If the check fails, the server does not start.

## Routes need a signed-in user

Lines 93-136

Asked as: api says not signed in; need to log in for api; unauthenticated api call; api requires login; route rejects anonymous

Nearly every custom route requires a signed-in user from the users collection. The exceptions are /api/version, which anyone can read, and the legacy export routes (/api/export_legacy), which use machine secrets for their own authentication.

## Version endpoint is public

Lines 93-105

Asked as: app version info; which build is running; git commit of server; check turbo version; version without login

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Time sheet commit and approval routes

Lines 109-117

Asked as: bundle timesheet approval; reject timesheet route; committer timesheet commit; uncommit timesheet; timesheet approve endpoint

Time sheets have routes to bundle, unbundle, approve, reject, commit and uncommit. The comment says committers commit time sheets. Time amendments also have a commit route.

## Bundle, unbundle and copy time sheet routes

Lines 109-111

Asked as: recall timesheet; copy timesheet next week; bundle by week ending; copy entry to tomorrow; unbundle timesheet

Staff bundle a time sheet by week ending. They unbundle (recall) a time sheet and copy its entries to the next week by time sheet id. A time entry can be copied to tomorrow. All need a signed-in user.

## Time sheet reports and tracking routes

Lines 120-131

Asked as: timesheet tallies; missing timesheet exception; ignored timesheets list; not expected timesheets; add missing timesheet user

Time sheet tallies (own, pending, approved, shared) and the tracking routes (counts, week lists, missing, ignored, not expected) need a signed-in user. Staff add or remove a missing time sheet exception per week and per user (uid). This file does not check a claim for these routes. The handlers decide who may use them.

## Legacy exports use machine secrets, not a user sign-in

Lines 133-136

Asked as: legacy export failing; export without login; machine secret export; legacy export access; export api key

The legacy export routes for time sheets, jobs and expenses are not behind the user sign-in check. Their handlers do their own authentication with machine_secrets.

## Expense routes need sign-in; the expenses editing flag is not checked here

Lines 138-168

Asked as: expense edit allowed while flag off; expense routes no flag; can edit expense when locked; expense api editing; expense submit recall endpoint

All /api/expenses routes require a signed-in user. They include create, update, submit, recall, approve, reject, commit, uncommit, settle and attachment actions. This file does not attach the expenses editing flag check to them.

## Expense uncommit requests routes

Lines 149-176

Asked as: request uncommit expense; cancel uncommit request; approve uncommit request; decline uncommit; uncommit request list

Expense uncommit requests are listed, approved, declined and cancelled through /api/expense_uncommit_requests. Staff create a request with POST /api/expenses/{id}/uncommit_requests. Approve takes no id in the path. Decline and cancel take the request id. All of these need a signed-in user.

## Time amendments can be committed through a route

Lines 178-180

Asked as: commit time amendment; amendment commit endpoint; who can commit amendment; time amendment approval; commit amendment route

POST /api/time_amendments/{id}/commit uses the same commit handler as time sheets and expenses. It needs a signed-in user. The handler decides who may commit.

## Admin profile routes

Lines 195-205

Asked as: admin profile edit; save user claims; set user active; clear authorized provider; admin user settings

Admin profile routes need a signed-in user. They cover identity, active toggle targets, save with claims, limited save, set active and clearing an authorized provider. The handlers decide who may do each.

## Purchase order lists and details need a signed-in user

Lines 207-219

Asked as: pending purchase orders; PO search list; purchase order approvers; approval limits lookup; second approver list

The purchase order routes for pending, visible, search, approvers, approval limits, second approvers and groups require a signed-in user (users). Turbo does not check the expenses editing flag on these read routes. Only the mutation routes check that flag.

## Purchase order actions blocked when expenses editing is off

Lines 221-237

Asked as: can't approve PO; PO cancel blocked; close purchase order refused; revise PO disabled; group PO creation blocked

The purchase order approve, reject, cancel, close, make_cumulative and revise actions, and group preview and group creation, check the expenses editing flag first. When the flag is off, Turbo refuses them with the expenses editing disabled message. Purchase order list and view routes are not gated by this flag.

## Legacy purchase order routes skip the expenses editing flag

Lines 239-243

Asked as: legacy PO edit; old purchase order form; create legacy PO; PO edit works when locked; legacy PO update

The legacy purchase order routes (edit view, create, update) require a signed-in user. They are not in the group that checks the expenses editing flag. Turbo does not refuse them when expenses editing is off in this file.

## Absorb routes exist for clients, contacts, vendors and rate sheets

Lines 253-267

Asked as: merge duplicate clients; absorb duplicate vendor; undo merge contact; combine rate sheets; absorb undo

Turbo offers an absorb action and an undo-absorb action for clients, client contacts, vendors and rate sheets. Absorb combines duplicates into a target. For rate sheets, the route comments say absorb is for an admin, or for a person with the rate_sheet_revise claim.

## Job rate sheet and job management routes

Lines 272-306

Asked as: change job client; close job; adopt subjobs; job status change; rate sheet period edit

Job routes need a signed-in user. They cover create, upsert, set status, set number, change client, adopt subjobs, close and validate proposal. Rate sheet periods are added, updated and deleted per job. The handlers decide who may do each action.

## Work-in-progress report modes

Lines 308-312

Asked as: WIP report; work in progress report; job WIP missing rates; my WIP division; branch WIP view

Work-in-progress (WIP) reports exist in three modes: my, branch and division. They need a signed-in user. A job WIP and its missing rates are also available per job.

## Report routes sit behind the report claim

Lines 314-333

Asked as: report page locked; reports permission missing; receipts report access; payables spreadsheet denied; branch mismatch report

All /api/reports routes require the report claim after sign-in. They include payroll and weekly time and expense reports, receipts, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI employee branch hours report

Lines 335-343

Asked as: employee branch hours; branch hours KPI; employee hours by branch; KPI hours report; employee hours report access

The employee branch hours report is at /api/kpi/reports/employee_branch_hours. It requires the kpi claim after sign-in.

## Branch hours report access check

Lines 345-353

Asked as: branch hours denied; can't see branch hours; branch hours access; branch hours query; branch hours forbidden

The branch hours report (/api/branch_hours) runs requireBranchHoursAccess with the branch query parameter before it returns data. The rule for who passes is in another file.

## Machine secrets routes are admin only

Lines 360-364

Asked as: create machine secret; machine secret list; secret generation admin; api secret creation; who can make secrets

The machine secrets routes list and create secrets. The route comment says they are for admins only. The check is in the handlers.

## Attachment audit routes

Lines 366-373

Asked as: attachment audit; orphaned attachments; missing attachments csv; delete orphan files; attachment audit run

Attachment audit routes need a signed-in user. They cover listing targets, viewing a run, refreshing, deleting orphans and downloading the missing and orphaned CSV reports.

## Currency routes

Lines 375-381

Asked as: exchange rates reload; currency backfill; delete currency; currency init status; update currency rates

Currency routes need a signed-in user. They cover listing, init status, reloading rates, initializing a backfill and deleting a currency. The handlers decide who may use them.

## Rate sheet routes and their admin notes

Lines 383-402

Asked as: rate sheet rollout; delete rate sheet; rename rate sheet; rate sheet former names; activate rate sheet

Rate sheet routes need a signed-in user. They cover create, compare, activate, deactivate, delete, rename, rollout and former names. The file notes that delete (an inactive sheet nothing uses) and rollout (move projects to the latest revision) are for admins. Absorbing a duplicate needs admin, or absorb with the rate_sheet_revise claim.

## Claims routes

Lines 404-410

Asked as: claim details; assign claim bulk; list claims; assignable users; claim admin detail

Claims routes need a signed-in user. They list claims, show assignable users and bulk assign a claim. The file notes the claim detail endpoint needs the admin claim.

## Rate sheet entry updates

Lines 417-420

Asked as: rate sheet entry edit; update rate entry; change rate sheet line; edit rate sheet price; rate entry admin only

PUT /api/rate_sheet_entries/{id} updates a rate sheet entry. It needs a signed-in user. The file notes the admin claim is required, and the handler enforces it.
