<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"37e0bd31e71ffbd9499c9be023b4f960a268a3116bf7dffc2ad3b8ad9bae9d50","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses editing check covers expenses, POs and vendors

Lines 14-31

Asked as: expenses editing off error; can't edit expenses or POs; vendor changes blocked; editing disabled message; expense editing switched off

The helper requireExpensesEditing refuses only for the collections expenses, purchase_orders and vendors, when the expenses editing flag is off. It returns the message in ErrExpensesEditingDisabled. For any other collection it allows the request. In this file only the purchase order mutation routes call it.

## Time editing check refuses when time flag is off

Lines 33-43

Asked as: time editing disabled; can't edit timesheets; time flag off; timesheet changes blocked; time editing turned off

The helper requireTimeEditing refuses with the message in ErrTimeEditingDisabled when the time editing feature flag is off. This file defines it but does not apply it to any route here.

## Reports need the report claim

Lines 45-54

Asked as: not authorized to view report; can't open reports; report access denied; no permission for reports; reports page forbidden

Every route under /api/reports needs a signed-in user who holds the report claim. A user without it is refused with "you are not authorized to view this report". This covers payroll, weekly time, expense and receipts reports, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report not authorized; can't see KPI reports; KPI access denied; no permission for KPI; kpi reports forbidden

Routes under /api/kpi/reports (such as employee branch hours) need a signed-in user who holds the kpi claim. A user without it is refused with "you are not authorized to view KPI reports".

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not working; where is the help chat; help chat slow to load; chatbot loading first time; open help assistant

Turbo adds the Help chat route at /api/help. The chat loads its model and its corpus the first time someone uses it.

## Server refuses to start if expenditure kinds config is invalid

Lines 88-91

Asked as: server won't start; expenditure kinds config error; turbo fails to start; startup config validation failed; app not starting after deploy

When Turbo starts serving, it validates the expenditure kinds configuration with ValidateExpenditureKindsConfig. If that check fails, the server does not start and the routes are not registered.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which build am I on; git commit of turbo; version without logging in; check turbo build info

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (long and short), git branch and build time.

## Routes require a signed-in user

Lines 107-139

Asked as: page asks me to log in; getting logged out of turbo; unauthorized api error; need to sign in; access denied without login

Time sheet, expense, purchase order, job, client, vendor, invoice, rate sheet, claim, currency and other custom route groups need a signed-in user from the users collection. The exceptions are /api/version and the legacy export routes under /api/export_legacy. The exports use their own authentication with machine secrets.

## Time sheet actions available

Lines 107-131

Asked as: bundle my timesheet; recall submitted timesheet; copy timesheet to next week; who can approve timesheets; missing timesheet exceptions

Time sheet routes cover bundle for a week ending, unbundle (recall), copy entries to next week, approve, reject, commit by committers, uncommit, reviewers and approver lookups, and tallies. Report holders also get tracking routes: counts, lists by week, missing, ignored and not expected. They can add or remove a missing-timesheet exception for a user and week.

## Expense actions available

Lines 138-169

Asked as: recall my expense; copy expense to tomorrow; request expense uncommit; upload expense receipt hash; mark receipt as missing

Expense routes cover create, update, submit, recall, copy to tomorrow, approve, reject, commit and uncommit. A user can also request an uncommit, upload or audit attachment hashes, mark an attachment missing, and settle or clear settlement. Lists include pending, approved, commit queue, unsettled, settled, tracking, and the caller's reimbursements by pay period.

## Expense uncommit requests can be approved, declined or cancelled

Lines 171-177

Asked as: decline uncommit request; cancel my uncommit request; approve uncommit requests; list uncommit reasons; batch approve expense requests

Under /api/expense_uncommit_requests, a user can list requests and reasons, approve requests in a batch, and decline or cancel one request by id.

## Time amendments are committed through their own route

Lines 179-181

Asked as: commit time amendment; time amendment not committing; amendment commit button; approve time change commit; commit amended hours

A time amendment is committed with POST /api/time_amendments/{id}/commit. It uses the same commit handler as time sheets and expenses.

## Time entry copy-to-tomorrow route requires sign-in

Lines 183-186

Asked as: copy time entry tomorrow; duplicate time entry next day; copy hours to tomorrow; time entry copy failing; repeat entry for tomorrow

A signed-in user can copy a time entry to tomorrow through the time entries route (copy_to_tomorrow). Requests without a signed-in user are refused.

## User defaults and navigation badges routes

Lines 188-194

Asked as: my user defaults; navigation badge counts; menu badges not showing; my default settings; sidebar notification counts

Signed-in users can fetch their own defaults (/api/users/defaults) and their navigation badges (/api/nav/badges).

## Admin profile routes available

Lines 196-206

Asked as: admin profile settings; save profile claims; clear authorized provider; set profile active; manage user identities

Signed-in users can reach the admin profile routes: identity list and identity per profile, active toggle targets, save with claims, save limited, set active, and clear an authorized provider. Who may use each route depends on the handler, which this file does not show.

## Purchase order read routes do not need expenses editing

Lines 208-221

Asked as: view purchase orders when editing off; see POs with editing disabled; read PO list; search purchase orders; view PO approvers

Signed-in users can list and view pending and visible purchase orders, their expenses, client invoice, search, approvers, second approvers, approval limits and groups. These read routes do not check the expenses editing switch, only the mutation routes do.

## Purchase order changes need expenses editing switched on

Lines 223-239

Asked as: can't approve PO; PO approve button not working; cancel purchase order blocked; close PO refused; create PO group fails

Approving, rejecting, cancelling, closing, converting to cumulative or revising a purchase order, and previewing or creating a PO group, all need the expenses editing feature flag to be on. When the flag is off, Turbo refuses these actions with the message in ErrExpensesEditingDisabled. Reading purchase orders (pending, visible, search, groups) is not gated by this flag in this file.

## Legacy purchase order routes skip the expenses editing check

Lines 241-245

Asked as: legacy PO edit; old purchase order form; create legacy purchase order; update old PO; PO works with editing off

The legacy purchase order routes (edit view, create, update) require a signed-in user. They are not in the group that checks expenses editing, so this file does not block them when expenses editing is off.

## Client invoice routes

Lines 247-252

Asked as: issue client invoice; invoice lines view; billable items for invoice; send invoice to client; invoice issue button

Signed-in users can issue an invoice (/api/invoices/{id}/issue) and read its lines and its billable items. Drafts and lines are written through their collections, not through these routes.

## Client routes

Lines 254-261

Asked as: view client list; client notes; client rate sheets; busdev leads; client time divisions

Signed-in users can list or view clients, and read client notes, time divisions, rate sheets and busdev leads (devleads).

## Absorb and undo absorb routes exist for several record types

Lines 262-276

Asked as: absorb duplicate client; undo merge of vendor; merge rate sheets; duplicate rate sheets list; undo absorb record

Clients, client contacts, vendors and rate sheets each have an absorb route (POST /{id}/absorb) and an undo absorb route (POST /undo_absorb). Rate sheet absorb also has a duplicates list and an absorb candidates list.

## Vendor list and view routes

Lines 270-274

Asked as: vendor list; look up vendor; view vendor details; vendors page not editable; find a supplier

Signed-in users can list and view vendors through /api/vendors. The vendor routes do not use the expenses editing check in this file.

## Project request routes are registered separately

Lines 278-278

Asked as: project request routes; job request registration; project requests not showing

Project request routes are added by registerJobRequestRoutes. This file does not show their rules.

## Job routes available

Lines 280-305

Asked as: create or update job; job details and notes; latest jobs list; stale or unused jobs; job time and staff summary

Signed-in users can create or upsert a job, and read details, notes, latest, unused and stale jobs. They can read time, staff, division, expense, PO and WIP summaries for a job.

## Job management actions

Lines 306-315

Asked as: change job status; set job number; change job client; close a job; validate job proposal; rate sheet periods for job

Jobs have routes to set status, set number, change client, adopt subjobs, close, and validate a proposal. Jobs also have routes to list, add, update and delete rate sheet periods. The handlers decide who may use them.

## WIP report modes

Lines 317-321

Asked as: WIP report; my work in progress; branch WIP figures; division WIP view; work in progress summary

Signed-in users can fetch the WIP report in three modes: my, branch and division.

## Report routes require report claim

Lines 323-342

Asked as: report page blocked; payroll report access; can't download receipts report; payables spreadsheet access; active jobs report missing

Every route under /api/reports first checks the report claim. These include payroll, weekly time, expense and receipts reports, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI employee branch hours report

Lines 344-352

Asked as: employee branch hours report; branch hours by employee; KPI hours report; hours per branch for staff

The employee branch hours report is at /api/kpi/reports/employee_branch_hours. It passes the kpi claim check first.

## Branch hours report checks access for the requested branch

Lines 354-362

Asked as: branch hours access denied; can't see branch hours; branch hours not allowed; branch hours for other branch

The /api/branch_hours route needs a signed-in user. Before it runs, it calls requireBranchHoursAccess with the branch query parameter. The rule for who passes is in another file.

## Stats dashboard needs sign-in

Lines 364-367

Asked as: stats dashboard blank; admin stats page; open stats page; dashboard needs login

The admin stats dashboard route (/api/stats) requires a signed-in user.

## Machine secrets routes are for admins

Lines 369-373

Asked as: machine secrets list; create machine secret; secrets page for admins; API secret creation

Listing and creating machine secrets are under /api/machine_secrets and are marked admin only in a comment. This file only requires sign-in; the admin check is in the handlers, which are not in this file.

## Attachment audit routes

Lines 375-382

Asked as: attachment audit; download missing attachments CSV; orphaned attachments report; refresh attachment audit; delete orphan attachments

Signed-in users reach the attachment audit routes: list targets, view a target's run, refresh, delete orphans, and download the missing and orphaned CSV reports.

## Currency routes

Lines 384-390

Asked as: currency list; reload exchange rates; initialize currency backfill; delete a currency; currency init status

Signed-in users reach the currency routes: list currencies, check init status, reload rates, initialize backfill for a currency, and delete a currency. The handlers decide who may use each.

## Rate sheet routes need the job claim, rollout is admin

Lines 392-409

Asked as: create rate sheet; activate rate sheet; rollout rate sheet projects; rename rate sheet; compare rate sheets

Rate sheet management routes (create, compare, activate, deactivate, rename, rollout) are commented as needing the job claim. Rollout, which moves projects to the latest revision, is marked admin. Absorbing duplicates is for an admin or a holder of rate_sheet_revise. The handlers, which are not in this file, enforce these rules.

## Claim detail and rate sheet entry edits need the admin claim

Lines 411-427

Asked as: claim details page; edit rate sheet entry; update rate entry rates; claim detail access; change rate sheet line

The claim detail route is commented as needing the admin claim, and updating a rate sheet entry (PUT /api/rate_sheet_entries/{id}) is also commented as needing the admin claim. The handlers outside this file enforce it.
