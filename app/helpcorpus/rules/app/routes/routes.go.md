<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"11df12c04d956feb31b34dd9d4bc43846af97a35778addb1dc1d719b36b9e26f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expense, PO and vendor edits blocked when expenses editing is off

Lines 14-31

Asked as: cannot edit expenses; expense edits blocked; can't change purchase order; vendor changes not saving; editing turned off

When the expenses editing feature flag is off, Turbo refuses requests on the expenses, purchase_orders and vendors collections. It gives the message of ErrExpensesEditingDisabled. Other collections are not affected by this flag.

## Time editing refused when the time feature flag is off

Lines 33-43

Asked as: can't edit my time; time entries locked; timesheet edit refused; time changes not allowed; time editing off

When the time editing feature flag is off, requireTimeEditing refuses the request. Turbo gives the message of ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: not authorized to view report; can't open payroll report; report access denied; no access to reports; why can't I see reports

Only staff with the report claim can use the report routes. Anyone else gets the message "you are not authorized to view this report". This covers payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI figures; not allowed KPI reports; employee branch hours blocked; KPI dashboard missing

Only staff with the kpi claim can use the KPI report routes. Anyone else gets the message "you are not authorized to view KPI reports". The employee branch hours report is one of them.

## Help chat is served at /api/help

Lines 84-85

Asked as: where is help chat; help chat not loading; open the assistant; help page url; chatbot slow to start

Turbo serves the Help chat at /api/help. It loads its model and corpus on first use, not at startup.

## Most Turbo API routes require a signed-in user

Lines 88-139

Asked as: API asks me to log in; unauthorized api call; need to sign in again; API returns not signed in; session required for routes

The time sheet, expense, purchase order, client, vendor, job, report, rate sheet, claim and other groups require a signed-in user (users auth). The version endpoint is public. The legacy export routes use machine secrets, not user sign-in.

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup fails config; bad expenditure config; app crashes on boot

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which build is running; git commit of server; check turbo version; build info endpoint

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit (long and short), git branch and build time.

## Time sheet bundle and unbundle have no time-editing check here

Lines 107-119

Asked as: bundle timesheet not working; unbundle time sheet; approve timesheet button; reject time sheet fails; commit timesheet problem

The time sheet routes (bundle, unbundle, copy to next week, approve, reject, commit, uncommit) need a signed-in user. This file does not call the time editing flag check for them. The check may happen in the handlers.

## Time sheet tallies come in four views

Lines 120-123

Asked as: timesheet counts dashboard; my pending timesheets; approved sheets count; shared timesheets for reviewer; timesheet totals by view

Time sheet tallies exist for the person's own sheets (uid), pending sheets for an approver, approved sheets for an approver, and shared sheets for a reviewer.

## Time sheet tracking routes for report holders

Lines 124-131

Asked as: missing timesheets list; ignored timesheets count; add timesheet exception; not expected timesheets; week list for timesheets

The time sheet tracking routes (counts, week list, missing, ignored, not expected, and add or remove a missing-timesheet exception for a user) sit in the time sheet group and need a signed-in user. This file does not check the report claim for them. Handlers do any further check.

## Legacy export endpoints use machine secrets, not user sign-in

Lines 133-139

Asked as: export timesheets script; legacy export not working; radiator login report; jobs export by date; expenses export feed

The legacy export endpoints for time sheets (by week ending), jobs (updated after a date) and expenses (updated after a date) do not use a user sign-in. They authenticate with machine secrets. Radiator login reports also use machine secrets, with the radiator role.

## Expense routes need sign-in and no feature flag

Lines 141-171

Asked as: cannot create expense; expense submit fails; attach receipt to expense; settle expense not working; expense approval button

All /api/expenses routes need a signed-in user. These routes include create, update, submit, recall, approve, reject, commit, uncommit, settle and attachment actions. This file does not apply the expenses editing flag to them. The expenses editing flag is applied only to purchase order mutation routes in this file.

## Expense uncommit request routes require sign-in

Lines 173-179

Asked as: uncommit request list; cancel uncommit request; approve uncommit request; decline uncommit reason; uncommit requests missing

All routes under /api/expense_uncommit_requests require a signed-in user. They cover listing requests, listing reasons, approving requests, declining a request and cancelling a request. This file adds no other check or feature flag; the handlers decide who may act.

## Time amendment commit route requires sign-in

Lines 181-183

Asked as: commit time amendment; amendment commit fails; time change commit button; approve amendment commit

A signed-in user can call the commit route for a time amendment (POST /api/time_amendments/{id}/commit). It uses the same commit handler as time sheets and expenses. This file adds no time-editing check; the handler decides who may commit.

## Copying a time entry to tomorrow requires sign-in

Lines 185-188

Asked as: copy time entry tomorrow; duplicate entry to next day; move entry to tomorrow; copy entry next day

Only a signed-in user can copy a time entry to tomorrow (POST /api/time_entries/{id}/copy_to_tomorrow). This file adds no time-editing check.

## User defaults and navigation badges require sign-in

Lines 190-196

Asked as: user defaults not loading; navigation badge counts; nav badges missing; my default settings; sidebar counts not showing

A signed-in user can fetch the user defaults (/api/users/defaults) and the navigation badges (/api/nav/badges). Turbo refuses both routes to anyone who is not signed in.

## Admin profile routes require sign-in; handlers check permissions

Lines 198-208

Asked as: admin profile not saving; set user active; save user claims; clear authorized provider; identity list admin

All admin profile routes require a signed-in user. They cover identity lists, active-toggle targets, saving with claims, limited saves, setting active, saving identity and clearing an authorized provider. This file does not say who may use each one; the handlers decide.

## Only purchase order reads skip the expenses editing gate

Lines 210-246

Asked as: purchase order search blocked; view pending POs; PO list not loading; approval limits lookup; second approvers list

Purchase order reads (pending, visible, search, approvers, approval limits, second approvers, groups) need only sign-in. The expenses editing flag does not block them. Legacy purchase order routes (edit, create, update) are also outside the flag gate in this file.

## Purchase order mutations need expenses editing on

Lines 224-240

Asked as: cannot approve purchase order; PO close button disabled; revise PO fails; cancel purchase order blocked; group preview PO error

Approve, reject, cancel, close, make cumulative, revise and the group preview and create actions on purchase orders check the expenses editing flag first. When the flag is off, they are refused. Reading purchase orders is not gated by this check.

## Client and contact routes require sign-in

Lines 248-262

Asked as: absorb client fails; undo client merge; merge client contacts; client absorb button; client routes denied

Client routes require a signed-in user. That includes the absorb and undo-absorb routes for clients and client contacts. This file does not say who may absorb; the absorb handler decides.

## Job, client, vendor absorb routes

Lines 256-270

Asked as: merge duplicate vendors; absorb rate sheet; undo absorb job; merge duplicate clients; absorb vendor records

Absorb and undo-absorb routes exist for clients, client contacts, vendors and rate sheets. They need a signed-in user. Handlers decide who may absorb.

## Vendor routes are not gated by the expenses editing flag here

Lines 264-270

Asked as: vendor list not loading; vendor view blocked; absorb vendor fails; undo vendor absorb; vendors still editable

Vendor routes (list, view, absorb and undo absorb) require sign-in only. This file does not apply the expenses editing flag to vendors, though requireExpensesEditing covers vendors elsewhere.

## Job routes require sign-in; handlers decide permissions

Lines 274-309

Asked as: job create fails; change job client; close job error; adopt subjobs problem; validate job proposal

All /api/jobs routes require a signed-in user. They cover create, upsert, details, notes, time, expense and PO summaries, WIP, set-status, set-number, change-client, adopt-subjobs, rate sheets, close and validate-proposal. This file adds no claim check; each handler decides who may act.

## WIP reports come in three modes

Lines 311-315

Asked as: WIP report my jobs; branch WIP report; division wip figures; work in progress report; WIP view not loading

Turbo offers the work-in-progress report in three modes: my, branch and division (/api/wip/my, /branch and /division). A signed-in user is required. This file adds no claim check.

## Reports routes need sign-in and the report claim

Lines 317-336

Asked as: payroll report download; receipts report missing; payables spreadsheet export; active jobs report access; branch mismatch report

All /api/reports routes need a signed-in user who holds the report claim. This covers payroll and weekly time and expense reports, receipts, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI reports route requires the kpi claim

Lines 338-346

Asked as: KPI reports forbidden; kpi report route access; cannot open KPI page; employee hours KPI missing

The /api/kpi/reports routes need a signed-in user who holds the kpi claim. The employee branch hours report is one of them.

## Branch hours report checks access to the requested branch

Lines 348-356

Asked as: branch hours access denied; branch hours report blocked; hours for other branch; branch report permission error

The branch hours report calls requireBranchHoursAccess with the branch query parameter before it runs. The rule for who passes is in another function that this file does not show.

## Machine secrets and claims admin routes

Lines 363-413

Asked as: machine secret admin; create machine secret; claim details admin only; list machine secrets fails

Machine secret routes (list, create) are admin only according to the comment. The claim detail endpoint needs the admin claim. Both need a signed-in user, and the handlers enforce the admin check.

## Rate sheet management claims are checked in handlers

Lines 386-405

Asked as: cannot delete rate sheet; roll projects revision; absorb duplicate rate sheets; create rate sheet permission; rate sheet admin action

Rate sheet routes need sign-in. The comments say creating needs the job claim. Deleting an inactive unused sheet and rolling projects to the latest revision are admin actions. Absorbing duplicates is for admins or holders of rate_sheet_revise. The handlers enforce these, not this file.

## Rate sheet entry update requires admin claim

Lines 420-423

Asked as: cannot update rate entry; rate sheet entry save fails; edit rate entry admin; change rate sheet line

Updating a rate sheet entry (PUT) needs a signed-in user. The comment says it needs the admin claim. The handler enforces this.
