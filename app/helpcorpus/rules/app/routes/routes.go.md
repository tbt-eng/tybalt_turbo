<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"cbc12f034ccaac055744513ef480895e3c42c965df9ff77d40b6f50a8d7a93df","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses editing flag covers expenses, purchase orders and vendors

Lines 14-31

Asked as: expenses editing off; can't edit expenses; vendor changes refused; purchase orders edit blocked

The expenses editing flag applies only to the collections expenses, purchase_orders and vendors. When the flag is off, Turbo refuses with the message from ErrExpensesEditingDisabled. Other collections are never blocked by this check.

## Time editing refused when time flag is off

Lines 33-43

Asked as: can't edit timesheet; time editing blocked; time entries locked; time edits refused

When the time editing feature flag is off, the time editing check refuses with the message from ErrTimeEditingDisabled. This code defines the check, but this file does not show which routes use it.

## Reports need the report claim

Lines 45-54

Asked as: not authorized to view report; can't open reports; report access denied; payroll report permission

Every route under /api/reports needs a signed-in user who holds the report claim. Anyone else gets: "you are not authorized to view this report". These reports include payroll time, weekly time, expense, receipts, payables spreadsheets, branch mismatches and active jobs.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI reports; kpi report permission; KPI reports locked

Routes under /api/kpi/reports (such as employee branch hours) need a signed-in user who holds the kpi claim. Anyone else gets: "you are not authorized to view KPI reports".

## Help chat is served at /api/help

Lines 84-85

Asked as: where is help chat; help chat slow to load; chatbot endpoint; assistant not responding

Turbo registers the Help chat route at /api/help. The Help chat loads its model and corpus on first use, not at startup.

## Turbo will not start with invalid expenditure kinds config

Lines 88-91

Asked as: server won't start; expenditure kinds config invalid; startup fails on config; turbo crashes on boot

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: app version check; what version is turbo; build number endpoint; git commit hash

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit, short commit, git branch and build time.

## Most Turbo API routes need a signed-in user

Lines 107-139

Asked as: api needs login; unauthorized api call; 401 on api; api requires signed in user

The time sheet, expense, purchase order, client, vendor, job, rate sheet, claims, currency, admin profile and similar route groups require a signed-in user from the users collection. The legacy export routes are the exception. They use machine secrets for authentication instead.

## Time sheet actions available

Lines 107-117

Asked as: bundle timesheet for week; recall timesheet; copy timesheet next week; approve or reject timesheet

A time sheet can be bundled for a week, unbundled (recalled), copied to next week, approved, rejected, committed and uncommitted. Time amendments can also be committed.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: export api no login; legacy export secret key; machine secret for export; export routes without sign in

The legacy export routes for time sheets, jobs and expenses are not behind user sign-in. Their handlers authenticate with machine secrets.

## Expense uncommit requests have their own routes

Lines 147-176

Asked as: uncommit expense request; bulk approve uncommit; decline uncommit request; cancel uncommit request

A committed expense can be uncommitted directly, or by creating an uncommit request. Requests can be listed, approved in bulk, declined or cancelled. The request reasons have their own list.

## Time amendments can be committed through the API

Lines 178-180

Asked as: commit time amendment; amendment commit api; approve time amendment; commit amendment through api

A signed-in user can call the commit action on a time amendment (/api/time_amendments/{id}/commit). It uses the same commit handler as time sheets and expenses; the handler decides who may commit.

## Time entries can be copied to tomorrow

Lines 182-185

Asked as: copy entry to tomorrow; duplicate time entry next day; time entry tomorrow copy; copy hours to tomorrow

A signed-in user can copy a time entry to tomorrow through /api/time_entries/{id}/copy_to_tomorrow. Turbo handles this copy in its own handler.

## Admin profile routes need sign-in

Lines 195-205

Asked as: admin profile access; save admin profile claims; clear authorized provider; set profile active

All admin profile routes need a signed-in user. They cover identity, active toggle targets, saving with claims, limited saves, setting active, and clearing an authorized provider. The handlers decide who may use each one.

## Purchase order read routes do not need expenses editing

Lines 207-237

Asked as: view purchase orders expenses off; PO list when editing disabled; purchase order search read only; see pending POs

Purchase order lists, pending and visible views, search, approvers, approval limits, second approvers and groups need only sign-in. They work when expenses editing is off. Only approve, reject, cancel, close, make cumulative, revise and group preview or create are blocked by the flag.

## Purchase order changes need expenses editing switched on

Lines 221-237

Asked as: can't approve purchase order; reject PO blocked; close purchase order refused; group purchase orders fails

Approving, rejecting, cancelling, closing, converting to cumulative, revising and creating groups of purchase orders go through the purchase order change routes. Turbo refuses all of them when the expenses editing feature flag is off. The refusal message is the one in ErrExpensesEditingDisabled. The check runs before the action starts.

## Legacy purchase order routes skip the expenses editing flag

Lines 239-243

Asked as: old purchase order edit form; legacy PO create update; PO edit view legacy; legacy purchase order save

The legacy purchase order routes (edit view, create, update) need sign-in only. The expenses editing middleware does not apply to them in this file.

## Client absorb and undo absorb routes

Lines 245-258

Asked as: absorb client; undo client absorb; merge client contact; undo client merge

A signed-in user can call absorb on a client or client contact, and undo absorb for each. The handler decides who may absorb.

## Vendor absorb and undo routes

Lines 260-266

Asked as: absorb vendor; undo vendor merge; merge vendors undo; vendor absorb action

A signed-in user can call absorb and undo absorb on vendors. The vendor routes are not behind the expenses editing flag in this file.

## Job rate sheet routes and job actions need sign-in

Lines 270-305

Asked as: close job; adopt subjobs; change job client; set job status

Jobs routes cover creating, upserting, set status, set number, change client, adopt subjobs, close, validate proposal, and rate sheet periods. They need sign-in here; handlers set who may use them.

## WIP report modes are my, branch and division

Lines 307-311

Asked as: WIP report my branch division; work in progress report; wip branch view; wip division report

The WIP report has three routes: /api/wip/my, /api/wip/branch and /api/wip/division. Each needs sign-in.

## Report routes need the report claim

Lines 313-332

Asked as: payables spreadsheet access; receipts report permission; branch mismatch report denied; active jobs report access

Every route under /api/reports needs sign-in and the report claim. This covers payroll, weekly, receipts, payables spreadsheets, branch mismatch and active jobs reports. Others get "you are not authorized to view this report".

## KPI report route needs the kpi claim

Lines 334-342

Asked as: employee branch hours report; branch hours KPI access; kpi employee hours denied; KPI branch report permission

The employee branch hours KPI report needs sign-in and the kpi claim. Others get "you are not authorized to view KPI reports".

## Branch hours report access depends on the branch requested

Lines 344-352

Asked as: branch hours report denied; can't see branch hours; branch hours access check; hours for other branch

The branch hours report checks access with requireBranchHoursAccess, using the branch query parameter and the signed-in user. The rule itself is in another file, so this file does not show who passes.

## Machine secrets and other admin routes check access in handlers

Lines 354-406

Asked as: machine secrets page; attachment audit access; stats page permission; claims details access

The routes for machine secrets, stats, attachment audit, currencies, rate sheet rollout and claims details require only sign-in here. Each handler applies its own admin or claim check, which this file does not show.

## Rate sheet absorb and rollout routes

Lines 391-398

Asked as: rate sheet rollout; absorb rate sheet undo; rate sheet revise claim; move projects latest revision

Rate sheets can be absorbed and the absorb undone, using the shared absorb handlers. Rollout moves projects to the latest revision. The comments say rollout and absorb are for admins, or for absorb with the rate_sheet_revise claim, but the handlers enforce this.
