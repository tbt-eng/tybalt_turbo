<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"159c4fb22aabdd6d5ffe069946ea189f1831473552433b79d6f122bfaa264803","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expense editing is off

Lines 14-31

Asked as: expense editing disabled; cannot edit expenses; expenses feature turned off; purchase orders locked; vendor changes blocked

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message from ErrExpensesEditingDisabled. Other collections are not affected by this flag.

## Time editing refused when time feature flag is off

Lines 33-43

Asked as: cannot edit time; time editing disabled; time entries locked; time feature off; can't change my hours

When the time editing feature flag is off, Turbo refuses the action with the message from ErrTimeEditingDisabled (requireTimeEditing).

## Reports need the report claim

Lines 45-54

Asked as: cannot view reports; report access denied; not authorized report; no permission for reports; reports missing access

A person without the report claim cannot view reports. Turbo refuses with 'you are not authorized to view this report'. This applies to every route under /api/reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: cannot see KPI reports; KPI access denied; kpi permission missing; kpi report not allowed; no kpi access

A person without the kpi claim cannot view KPI reports. Turbo refuses with 'you are not authorized to view KPI reports'. This applies to /api/kpi/reports.

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup config error; expenditure config broken; server fails to boot

When the server starts, Turbo validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which version running; build info; git commit version; version without login

Anyone can call GET /api/version without signing in. It returns the app name, version, build, git commit, git branch and build time.

## Time sheet routes require sign-in

Lines 107-131

Asked as: time sheet not working; cannot bundle time sheet; recall timesheet; approve time sheet error; timesheet sign in needed

All /api/time_sheets routes need a signed-in user. These include bundle, unbundle (recall), copy to next week, approve, reject, commit, uncommit, tallies and tracking. Time tracking routes are for report holders.

## Legacy export routes use machine secrets, not sign-in

Lines 133-136

Asked as: export api no login; legacy export secret; machine secret export; export without signing in; automated export access

The legacy export routes for time sheets, jobs and expenses do not require a signed-in user. Each handler does its own authentication with machine_secrets.

## Expense routes require sign-in

Lines 138-170

Asked as: expense submit fails; cannot approve expense; expense recall not working; attachment upload expense; expense login required

All /api/expenses routes need a signed-in user. They cover create, update, submit, recall, copy, approve, reject, commit, uncommit, uncommit requests, attachments, settlement, tracking and reimbursements. Unlike purchase order approval routes, these routes do not check the expenses editing flag here.

## Expense uncommit request routes require sign-in

Lines 172-178

Asked as: uncommit request approval; expense uncommit decline; cancel uncommit request; uncommit reasons list; uncommit requests not showing

The routes to list, approve, decline and cancel expense uncommit requests, and to list reasons, need a signed-in user.

## Time amendments commit route

Lines 180-182

Asked as: commit time amendment; amendment commit failed; approve amendment; time change commit; amendment not committing

Committing a time amendment (POST /api/time_amendments/{id}/commit) needs a signed-in user. It uses the shared commit handler.

## Time entry copy, user defaults and nav badges routes require sign-in

Lines 184-195

Asked as: copy time entry tomorrow; user defaults missing; nav badge counts; copy entry to tomorrow; notification counts not showing

Copying a time entry to tomorrow, reading user defaults and reading nav badges need a signed-in user.

## Admin profile routes require sign-in

Lines 197-207

Asked as: edit admin profile; user profile save; clear authorized providers; deactivate user account; change user claims

Admin profile routes need a signed-in user. They cover identity, active toggle, save with claims, limited save and clearing authorized providers. Each handler decides who may do the action.

## Purchase order read routes require sign-in only

Lines 209-221

Asked as: pending purchase orders; PO search not working; list approvers; approval limits lookup; view purchase order groups

Pending, visible, search, approvers, approval limits, second approvers and group read routes for purchase orders need a signed-in user. They do not check the expenses editing flag.

## Purchase order changes blocked when expense editing is off

Lines 223-239

Asked as: cannot approve PO; purchase order rejected blocked; close purchase order; revise PO not working; cancel PO refused

Approve, reject, cancel, close, make cumulative, revise and group create or preview for purchase orders need a signed-in user. They are refused when the expenses editing flag is off.

## Legacy purchase order routes skip the expenses editing flag

Lines 241-245

Asked as: old PO edit page; legacy purchase order form; create PO old screen; update old purchase order; PO edit view

Legacy purchase order routes (edit view, create, update) need a signed-in user. This group does not use the expenses editing middleware.

## Client, contact and vendor routes require sign-in

Lines 247-268

Asked as: cannot edit client; absorb client; vendor contact update; undo client absorb; client contacts not saving

Client, client contact and vendor routes need a signed-in user. This includes absorb and undo absorb. Each handler decides who may do the action.

## Job routes require sign-in

Lines 272-303

Asked as: job not saving; change job client; close job; adopt subjobs; job status update failing

All /api/jobs routes need a signed-in user. This includes create, upsert, details, time, expense and PO summaries, set status, set number, change client, adopt subjobs, close and validate proposal.

## WIP report routes

Lines 305-309

Asked as: WIP report missing; work in progress report; my WIP figures; division WIP; branch WIP not showing

WIP reports for 'my', 'branch' and 'division' need a signed-in user.

## Report routes need the report claim

Lines 311-330

Asked as: payroll report access; payables spreadsheet download; receipts report missing; branch mismatch report; active jobs report denied

All /api/reports routes need a signed-in user who holds the report claim. They include payroll and weekly time, expense and receipts reports, payables spreadsheets, time entry branch mismatches and active jobs.

## KPI employee branch hours report needs kpi claim

Lines 332-340

Asked as: employee branch hours KPI; KPI hours report access; branch hours by employee; kpi employee report denied; employee hours kpi

The employee branch hours report needs a signed-in user who holds the kpi claim.

## Branch hours report access depends on the branch

Lines 342-350

Asked as: branch hours report denied; cannot see branch hours; branch hours access; other branch hours report; branch report permission

The branch hours report needs a signed-in user. Turbo runs requireBranchHoursAccess with the branch query parameter before it returns the report. That function decides who may see which branch.

## Stats, machine secrets and attachment audit routes require sign-in

Lines 352-370

Asked as: stats dashboard not loading; machine secrets list; attachment audit access; create machine secret; audit attachments report

The stats dashboard, machine secrets list and create, and attachment audit routes need a signed-in user. The comment says machine secrets are admin only. The handlers enforce that, not this file.

## Currency, rate sheet, claims and work record routes require sign-in

Lines 372-403

Asked as: rate sheet not saving; currency rates update; claim details access; work record missing; claims page error

Currency, rate sheet, rate sheet entry, claims and work record routes need a signed-in user. Comments say rate sheets need the job claim and claim details need the admin claim. The handlers enforce that, not this file.
