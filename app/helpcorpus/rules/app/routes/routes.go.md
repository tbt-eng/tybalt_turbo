<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"67ba29eb6f1d1f380717133cfac508295002bf812b9442fe26d6f3a1d9f00c3b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses editing off; cannot edit expenses; purchase orders read only; vendor changes refused

When the expenses editing feature flag is off, Turbo refuses writes to the expenses, purchase_orders and vendors collections. The refusal message is the one in ErrExpensesEditingDisabled. Other collections skip this check.

## Time sheet editing needs the time editing flag

Lines 33-43

Asked as: time editing disabled; cannot change time sheet; time sheet locked; time editing flag off

When the time editing feature flag is off, a time sheet action that uses this check is refused with the message in ErrTimeEditingDisabled (requireTimeEditing). Nothing in this file applies the check to a route.

## Reports need the report claim

Lines 45-54

Asked as: cannot open report; report access denied; not authorized to view report; no permission for reports

A person without the report claim cannot open any report under /api/reports. Turbo refuses with 'you are not authorized to view this report'. This covers payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the KPI claim

Lines 56-65

Asked as: KPI report access denied; cannot see KPI reports; no KPI permission; KPI reports blocked

A person without the kpi claim cannot open KPI reports under /api/kpi/reports, such as employee branch hours. Turbo refuses with 'you are not authorized to view KPI reports'.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; help chat slow first time; where is the help chat; help chat endpoint

Turbo registers the Help chat route at /api/help. The chat loads its model and corpus the first time it is used.

## Startup fails if expenditure kinds config is invalid

Lines 88-91

Asked as: server won't start; startup error expenditure kinds; invalid expenditure config; routes not registered

When Turbo starts serving, it validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). An invalid configuration stops the routes from being registered.

## Staff must sign in to use Turbo API routes

Lines 93-136

Asked as: api asks me to log in; unauthorized api call; need to sign in; version endpoint public

Nearly all Turbo API routes require a signed-in user (users collection). The version endpoint, /api/version, is public. The legacy export routes under /api/export_legacy skip sign-in and use machine secrets instead.

## Time sheet routes need sign-in and have no flag check here

Lines 107-131

Asked as: time sheet route no flag; time sheet not gated; time sheet api login

All /api/time_sheets routes require a signed-in user. They cover bundle, unbundle (recall), copy to next week, approve, reject, commit, uncommit, tallies and tracking. The time editing flag is not checked at route level in this file.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: legacy export no login; export machine secret; legacy export authentication; export without sign in

The legacy export routes for time sheets, jobs and expenses are outside the sign-in groups. They use their own authentication through machine secrets, checked in the handlers.

## Expense routes need sign-in; recall and uncommit requests included

Lines 138-169

Asked as: expense recall not working; uncommit expense request; expense api sign in; copy expense to tomorrow

All /api/expenses routes require a signed-in user. They cover create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit and uncommit requests. Unlike purchase order changes, this file adds no expenses-editing flag check to them.

## Expense attachment audit, hash replace and missing marking

Lines 153-156

Asked as: expense attachment hash; replace attachment hash; mark attachment missing; audit expense attachment

Signed-in users reach routes to audit an expense attachment hash, replace the hash, and mark the attachment as missing. This file adds no claim check; the handlers decide who may do it.

## Expense settlement routes

Lines 162-165

Asked as: settle expense; clear settlement; unsettled expenses list; who can settle expenses

Turbo has routes to list unsettled and settled expenses, to settle an expense (/settle) and to clear a settlement (/clear_settlement). This file only requires sign-in; who may settle is decided in the handlers (settlement is for payables staff).

## Expense reimbursements are the caller's own, by pay period

Lines 167-169

Asked as: my reimbursements; expenses repaid on payroll; reimbursements by pay period; my expenses paid

The /api/expenses/reimbursements routes list the caller's own expenses repaid on payroll, grouped by pay period. A pay period ending date picks one period.

## Expense uncommit request routes

Lines 171-177

Asked as: uncommit request reasons; approve uncommit request; decline uncommit request; cancel uncommit request

Signed-in users can list uncommit requests and their reasons, approve requests, and decline or cancel a request by id. This file adds no claim check; the handlers decide who may act.

## Time amendments commit route

Lines 179-181

Asked as: commit time amendment; amendment commit button; time amendment not committing

A signed-in user can call POST /api/time_amendments/{id}/commit. It uses the same commit handler as time sheets and expenses; the handler decides who may commit.

## Admin profile, user defaults and nav badge routes

Lines 188-206

Asked as: save user profile; user defaults settings; nav badge count; clear authorized provider

Signed-in users reach routes for admin profile identity, save with claims, limited save, active toggle and clearing an authorized provider. They also reach user defaults and nav badges. The handlers decide who may save or change each.

## Purchase order read routes are open to signed-in users

Lines 208-221

Asked as: purchase orders blocked read; view purchase orders off; pending purchase orders list; PO search not working

Pending, visible, search, approvers, approval limits, second approvers and group read routes need only sign-in. They are not blocked when expenses editing is off.

## Purchase order changes are blocked when expenses editing is off

Lines 223-239

Asked as: cannot approve purchase order; PO approval refused; cannot close purchase order; cannot create PO group

Approving, rejecting, cancelling, closing, converting to cumulative and revising a purchase order, and previewing or creating PO groups, all check the expenses editing flag first. When the flag is off, Turbo refuses them. Viewing and searching purchase orders is not gated this way.

## Legacy purchase order routes skip the expenses editing flag

Lines 241-245

Asked as: edit purchase order legacy; create purchase order old; update PO not blocked

The legacy purchase order routes (edit, create, update) need only sign-in. They are in a separate group from the PO mutation routes, so this file does not apply the expenses-editing flag to them.

## Client invoice issue route

Lines 247-252

Asked as: issue client invoice; invoice lines billable items; create invoice for client

A signed-in user can issue an invoice, and fetch its lines and billable items. Drafts and lines are written through their own collections.

## Client, contact and vendor absorb routes

Lines 254-275

Asked as: merge clients; absorb vendor; undo absorb client; combine contacts

Signed-in users can reach routes to absorb one client, client contact or vendor into another and to undo an absorb. This file adds no claim check; the handlers decide who may do it.

## Job routes need sign-in

Lines 279-310

Asked as: job api sign in; job summary not loading; close job route; adopt subjobs

All /api/jobs routes require a signed-in user. They cover create, upsert, details, time, expense and PO summaries, WIP, set status, set number, change client, adopt subjobs, close and validate proposal. Claims for each action are checked in the handlers.

## WIP report modes

Lines 312-316

Asked as: my WIP report; branch WIP; division WIP report; work in progress report

Turbo offers WIP reports in three modes: my, branch and division. They are served at /api/wip/my, /api/wip/branch and /api/wip/division and need sign-in.

## All report routes need the report claim

Lines 318-337

Asked as: report needs claim; payroll report blocked; receipts report access; payables spreadsheet permission

Every route under /api/reports requires sign-in and the report claim. These include payroll and weekly time, expense and receipts reports, payables spreadsheets, time entry branch mismatches and active jobs.

## Employee branch hours KPI report

Lines 339-347

Asked as: branch hours KPI; employee hours by branch; branch hours report login

The employee branch hours report at /api/kpi/reports/employee_branch_hours requires sign-in and the KPI claim.

## Branch hours report checks branch access

Lines 349-357

Asked as: branch hours not allowed; branch hours access denied; which branch hours I see

Before the branch hours report runs, Turbo calls requireBranchHoursAccess with the branch query parameter. Who passes is decided in that function, which is not in this file.

## Admin stats and machine secrets routes require sign-in only here

Lines 359-410

Asked as: machine secrets admin; stats route access; rate sheets permission

The stats, machine secrets, claims, rate sheets and currencies routes require a signed-in user at the router. Comments say machine secrets are for admins and some claims and rate sheet routes need the admin or job claim. Those checks are in the handlers, not in this file.

## Attachment audit routes

Lines 370-377

Asked as: attachment audit run; download orphaned files csv; delete orphan attachments; missing attachments report

Signed-in users reach routes to list audit targets, view a run, refresh it, delete orphans and download the missing and orphaned CSV reports. This file adds no claim check; the handlers decide who may do it.

## Currency routes

Lines 379-385

Asked as: reload exchange rates; currency backfill; delete currency; currency init status

Signed-in users reach routes to list currencies, check init status, reload rates, start a backfill and delete a currency. Claim checks happen in the handlers.

## Rate sheet routes

Lines 387-392

Asked as: activate rate sheet; deactivate rate sheet; create rate sheet; rate sheet job claim

Signed-in users reach routes to create, activate and deactivate rate sheets. A comment says the job claim is required; the handlers enforce it, not this file.

## Claims and work record routes

Lines 394-405

Asked as: bulk assign claim; assignable users list; claim details view; work records list

Signed-in users reach routes to list claims, list assignable users, bulk assign a claim, and view claim details. A comment says claim details need the admin claim. Work record list and detail routes also need only sign-in here.
