<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"e54a5dc1496a635644381d72147432bf260d24420fd6eac410d5b7541ee578a5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors refuse changes when expenses editing is off

Lines 14-31

Asked as: can't edit expenses; expenses locked; purchase orders read only; vendor changes blocked; expense edits not saving

When the expenses editing feature flag is off, Turbo refuses requests for the expenses, purchase_orders and vendors collections. The refusal message is the one for ErrExpensesEditingDisabled. Other collections skip this check.

## Time editing refused when the time feature flag is off

Lines 33-43

Asked as: can't edit timesheet; time entries locked; time changes not saving; timesheet editing off; time edits refused

When the time editing feature flag is off, Turbo refuses with the message for ErrTimeEditingDisabled. The check runs only where a handler calls requireTimeEditing.

## Reports need the report claim

Lines 45-54

Asked as: can't open report; report access denied; no permission for reports; reports page forbidden; why can't I see reports

A person without the report claim cannot open any report under /api/reports. Turbo refuses with "you are not authorized to view this report". The same check covers payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: kpi report denied; can't see kpi dashboard; hours kpi not authorized; kpi access missing; kpi page blocked

A person without the kpi claim cannot open KPI reports such as the employee branch hours report. Turbo refuses with "you are not authorized to view KPI reports".

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat slow; assistant first answer slow; where is help chat; chatbot endpoint; help bot loading

Turbo registers the Help chat route at /api/help. The chat loads its model and its corpus of rule summaries on first use, so the first question may take longer to answer.

## Turbo checks expenditure kinds configuration when it starts serving

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup error expenditure kinds; config error on boot; expenditure config check fails

When Turbo starts to serve requests, it validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If the configuration is invalid, the error stops the server from starting. The result of the check is defined in another file.

## Turbo checks expenditure kinds configuration at startup

Lines 89-91

Asked as: server not starting; expenditure config broken; boot fails expenditure kinds; server refuses to serve; expenditure kinds validation error

When the server starts serving, Turbo validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If the check fails, the server does not start serving.

## Most API routes need a signed-in user

Lines 93-136

Asked as: api returns unauthorized; 401 on api routes; need to log in for api; routes require login; api says not signed in

Time sheet, expense, purchase order, client, vendor, job, report, rate sheet, claim and other routes require a signed-in user. The version endpoint (/api/version) and the legacy export endpoints are not under this check. The legacy exports use their own authentication with machine secrets.

## Version endpoint is public

Lines 93-105

Asked as: app version without login; which version is running; build info endpoint; version check no sign in; git commit of build

The version endpoint (/api/version) answers without sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Time sheet recall is the unbundle route

Lines 107-117

Asked as: recall timesheet; pull back timesheet; unbundle timesheet; withdraw submitted timesheet; undo submitted time sheet

Recalling a time sheet is done through the unbundle route. Committers commit and uncommit time sheets through the commit and uncommit routes. Time amendments have a commit route but no uncommit route.

## Time sheet tracking routes are for report holders

Lines 124-131

Asked as: missing timesheets list; timesheet exceptions; who hasn't submitted time; ignored timesheets endpoint; not expected timesheets

The time sheet tracking routes cover counts, weekly lists, missing, ignored and not-expected lists, and adding or removing a missing-timesheet exception for a person. They are described as endpoints for report holders. This file shows no report claim check on them; the handlers do any check.

## Legacy export routes use machine secrets, not sign-in

Lines 133-136

Asked as: legacy export auth; export api key; machine secret export; export without login; export calls fail auth

The legacy export routes for time sheets, jobs and expenses are outside the signed-in groups. They use their own authentication with machine secrets (machine_secrets).

## Expense routes include recall, settlement and uncommit requests

Lines 138-168

Asked as: recall expense; settle expense; copy expense to tomorrow; replace expense attachment; mark expense missing

Expenses can be submitted, recalled, approved, rejected, committed and uncommitted. They also have uncommit requests, attachment hash audit, replace and mark-missing actions, and settle and clear-settlement actions. Expenses can be copied to tomorrow.

## Expense uncommit requests have their own routes

Lines 170-176

Asked as: request expense uncommit; approve uncommit request; decline uncommit request; list uncommit reasons; cancel my uncommit

Expense uncommit requests have their own routes. A user can list requests and reasons, approve requests, and decline or cancel one request. Any signed-in user can reach these routes; the handlers decide who may act.

## Time amendments can be committed through the API

Lines 178-180

Asked as: commit time amendment; amendment not committing; finalize amendment via api; time amendment commit button; amendment commit endpoint

A time amendment has a commit route (/api/time_amendments/{id}/commit). It uses the same commit handler as time sheets and expenses.

## Admin profile routes include limited save and active toggle

Lines 195-205

Asked as: edit user profile admin; deactivate user account; clear authorized provider; admin save profile; set user inactive

Admin profile routes save with claims, save in a limited form, set active state, edit identity and clear an authorized provider. They need sign-in; the handlers decide who may act.

## Purchase order read routes work when expenses editing is off

Lines 207-219

Asked as: view purchase orders read only; search pos while locked; list approval limits; see po groups; view pending pos

The purchase order routes to view pending and visible POs, search, list approvers and approval limits, and view groups are not gated on expenses editing. Only the mutation routes are gated.

## Purchase order mutations need expenses editing enabled

Lines 221-237

Asked as: can't approve purchase order; reject po disabled; close po refused; revise po blocked; create po group fails

Approving, rejecting, cancelling, closing, making cumulative and revising a purchase order, and previewing or creating PO groups, all need the expenses editing flag on. When it is off, Turbo refuses these actions. Listing and viewing POs do not have this check.

## Legacy purchase order routes have no expenses editing gate

Lines 239-243

Asked as: legacy po edit; old purchase order form; create po legacy screen; update legacy purchase order; legacy po page

The legacy purchase order routes (edit view, create, update) need sign-in. This file does not apply the expenses editing check to them.

## Absorb and undo-absorb routes exist for clients, contacts and vendors

Lines 252-266

Asked as: merge client; absorb contact; undo merge vendor; combine duplicate vendors; unmerge client contact

Clients, client contacts and vendors each have an absorb route and an undo-absorb route. The routes need sign-in; the handlers decide who may use them.

## Job routes need a signed-in user and cover details, reports and lists

Lines 270-295

Asked as: job details api; job wip endpoint; stale jobs list; job time summary; create job api

Every /api/jobs route needs a signed-in user. Job routes give a job's details, notes, time summary, full time report, time entries, staff summary, divisions summary, expense summary and list, PO summary and list, and WIP. They also list latest, unused and stale jobs. Creating a job (POST) and upserting a job (PUT /{id}) are also routes here. Who may do each action is decided in the handlers, which are in other files.

## Job routes include status, number, client change and close

Lines 296-301

Asked as: change job client; close job; set job number; adopt subjobs; validate job proposal

Jobs have routes to set status, set number, change client, adopt subjobs, close, and validate a proposal. They need sign-in; the handlers decide who may act.

## WIP reports come in three modes

Lines 303-307

Asked as: wip my report; wip by division; branch wip; work in progress report; wip report modes

The WIP report is available in three modes: my, branch and division. All three need only sign-in at the route level.

## Report routes need the report claim

Lines 309-328

Asked as: payroll report access; receipts report denied; payables spreadsheet download; branch mismatch report; active jobs report access

All routes under /api/reports require the report claim after sign-in. They include payroll and weekly time and expense reports, receipts, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI report routes need the kpi claim

Lines 330-338

Asked as: kpi access denied; employee hours kpi; kpi not authorized message; employee branch hours report; kpi report permission

The employee branch hours KPI report needs the kpi claim after sign-in. Without it, Turbo says "you are not authorized to view KPI reports".

## Branch hours report checks access by branch

Lines 340-348

Asked as: branch hours report blocked; can't see other branch hours; branch query access; hours by branch denied; branch hours permission

The branch hours report calls requireBranchHoursAccess with the branch query parameter before it runs. The rule for who passes is in another file.

## Machine secrets routes are admin only

Lines 355-359

Asked as: create machine secret; api key admin; list machine secrets; machine secret admin only; who can make api keys

The machine secrets routes (list, create) are commented as admin only. This file only requires sign-in; the handlers enforce the admin check.

## Attachment audit routes cover refresh, orphan deletion and CSV reports

Lines 361-368

Asked as: orphaned attachments delete; missing attachments csv; refresh attachment audit; attachment audit run; orphan attachments report

Attachment audit routes list targets, show a run, refresh it and delete orphans. They also download the missing and orphaned CSV reports. They need sign-in; the handlers decide who may act.

## Currency routes include reloading rates, backfill and deletion

Lines 370-376

Asked as: reload exchange rates; currency backfill; delete currency; currency init status; refresh currency rates

Currency routes list currencies, show init status, reload rates, initialize a backfill and delete a currency. They need sign-in; the handlers decide who may act.

## Rate sheet routes need the job claim

Lines 378-386

Asked as: create rate sheet; rename rate sheet; activate rate sheet; compare rate sheets; rate sheet former names

Rate sheet routes create, compare, activate, deactivate and rename rate sheets and list former names. A comment says the job claim is required; the handlers enforce it.

## Claim detail and bulk assignment routes

Lines 388-394

Asked as: bulk assign claim; list assignable users; claim details admin; give user a claim; list claims

A comment says the claim detail route requires the admin claim. Other claim routes list claims, list assignable users, and bulk assign a claim.

## Rate sheet entry updates need the admin claim

Lines 401-404

Asked as: edit rate sheet entry; update rate entry admin; change rate sheet rate; rate entry permission; modify rate sheet line

Updating a rate sheet entry (PUT /api/rate_sheet_entries/{id}) is commented as needing the admin claim. The handler enforces it.
