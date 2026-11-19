<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"0fd9baff14b38e1ee5b343bb29c5e5c6008c47cffed9ed76f120ff39e9129ae6","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: can't edit expenses; expense changes blocked; vendor edits not allowed; purchase order editing disabled; expenses feature turned off

When the expenses editing feature flag is off, Turbo refuses changes to expenses, purchase orders and vendors with the message in ErrExpensesEditingDisabled. Other collections are not affected by this check.

## Time editing flag off blocks time editing

Lines 33-43

Asked as: can't change my timesheet; time entries locked; cannot edit my hours; timesheet editing turned off; time edits blocked

When the time editing feature flag is off, the time check returns a refusal with the message in ErrTimeEditingDisabled (requireTimeEditing).

## Reports need the report claim

Lines 45-54

Asked as: can't open payroll report; not authorized to view report; report access denied; weekly report won't load; receipts report permission

Anyone without the report claim gets the refusal 'you are not authorized to view this report'. This covers every route under /api/reports, such as payroll, weekly and receipts reports, the payables spreadsheets, the branch mismatch report and the active jobs report.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see kpi dashboard; branch hours report forbidden; employee hours report permission

Anyone without the kpi claim gets the refusal 'you are not authorized to view KPI reports'. This covers the employee branch hours report under /api/kpi/reports.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; where is the chatbot; assistant window missing; help chat slow first time

Turbo serves the Help chat at /api/help. The chat loads its model and corpus the first time someone uses it.

## Startup stops if expenditure kinds config is invalid

Lines 88-91

Asked as: server won't start; turbo fails to boot; expenditure kinds config error; app crashes on startup

When Turbo starts to serve, it validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: what version is turbo; build number; app version check; git commit hash

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit and branch, and build time.

## Most API routes require a signed-in user

Lines 107-131

Asked as: unauthorized api error; timesheet api logged out; session required api; api says log in

The time sheet, expense, purchase order, client, vendor, job, rate sheet and other custom route groups require a signed-in user from the users collection. Finer permission checks are in the handlers, not in this file.

## Legacy export endpoints use machine secrets, not user sign-in

Lines 133-136

Asked as: export timesheets script; machine secret export; legacy export api key; automated export failing

The legacy export routes for time sheets, jobs and expenses are outside the signed-in groups. They use their own authentication with machine secrets.

## Expense routes need sign-in; this file adds no editing-flag check

Lines 138-170

Asked as: expense api unauthorized; can't call expenses endpoint; expense api login needed; expenses api access error

Every /api/expenses route needs a signed-in user (the users collection). This file binds no expenses editing flag check to these routes. Unlike the purchase order change routes, any flag check on them would happen inside the handlers.

## Expense uncommit request routes need sign-in

Lines 172-178

Asked as: uncommit request not showing; reverse expense commit request; decline uncommit request; cancel uncommit request

The routes under /api/expense_uncommit_requests (list, reasons, approve, decline, cancel) require a signed-in user (users collection). This file adds no other check; the handlers decide who may act.

## Time amendment commit route needs sign-in

Lines 180-182

Asked as: commit time amendment; amendment won't commit; approve timesheet amendment; amendment commit error

Committing a time amendment (POST /api/time_amendments/{id}/commit) requires a signed-in user. It uses the same commit handler as time sheets and expenses, so the committer rule is checked there, not in this file.

## Time entry copy and user defaults need sign-in

Lines 184-195

Asked as: copy time entry tomorrow; duplicate entry to tomorrow; user defaults not loading; nav badge missing

Copying a time entry to tomorrow and fetching user defaults require a signed-in user. Nav badges also require sign-in.

## Admin profile routes need sign-in only here

Lines 197-207

Asked as: change user admin rights; toggle user active; clear auth provider; edit user claims

The admin profile routes (identity, active toggle, save with claims, save limited, set active, clear authorized provider) require a signed-in user. This file does not check admin rights; the handlers do.

## Purchase order read routes are not blocked by the expenses flag

Lines 209-221

Asked as: can't see purchase orders; pending POs not showing; search purchase orders; PO approvers list empty

Listing and viewing purchase orders (pending, visible, search, approvers, approval limits, second approvers, groups) require sign-in only. The expenses editing flag does not block them.

## Purchase order changes are blocked when expenses editing is off

Lines 223-239

Asked as: can't approve purchase order; PO cancel blocked; close PO not working; revise purchase order failing

The purchase order approve, reject, cancel, close, make_cumulative and revise routes, and the PO group preview and create routes, run the expenses editing check first. When the flag is off, Turbo refuses them. The read routes for purchase orders do not run this check.

## Legacy purchase order routes skip the expenses flag

Lines 241-245

Asked as: old purchase order create; legacy PO edit screen; update old purchase order

Legacy purchase order routes (view for edit, create, update) require sign-in. They are outside the group that checks the expenses editing flag in this file.

## Client, contact and vendor routes need sign-in

Lines 247-269

Asked as: client api unauthorized; absorb client; undo absorb vendor; contact api login error

Client, client contact and vendor routes, including absorb and undo absorb, require a signed-in user. Absorb rights are decided by the handlers.

## Job routes need sign-in; rights are checked in handlers

Lines 273-308

Asked as: job api login error; close job api; adopt subjobs; job api unauthorized

All /api/jobs routes require a signed-in user: create, upsert, details, reports, status, number, client change, adopt subjobs, rate sheets, close and validate proposal. Who may do each action is decided by the handlers.

## WIP reports come in three modes

Lines 310-314

Asked as: work in progress report; wip by branch; my wip report; division wip report

The WIP report is served for three modes: my, branch and division, at /api/wip/{mode}. Each requires sign-in only in this file.

## Report routes need the report claim before they run

Lines 316-335

Asked as: reports api forbidden; report endpoint 403; reports need permission; report route signed out

Every route under /api/reports requires a signed-in user and then the report claim.

## KPI report routes need the kpi claim

Lines 337-345

Asked as: kpi api forbidden; kpi endpoint 403; kpi route access denied

The routes under /api/kpi/reports require a signed-in user with the kpi claim.

## Branch hours report checks access for the requested branch

Lines 347-355

Asked as: branch hours report access; can't pull hours for branch; branch query denied; branch hours forbidden

The /api/branch_hours route runs requireBranchHoursAccess with the 'branch' query parameter before it builds the report. The rule for who passes is in another file.

## Stats, machine secrets and attachment audit need sign-in

Lines 357-375

Asked as: stats dashboard login; create machine secret; attachment audit not loading; machine secrets list denied

The stats dashboard, machine secrets list and create, and attachment audit routes require sign-in. The comment says machine secrets are admin only; the handlers enforce that.

## Currency routes need sign-in

Lines 377-383

Asked as: exchange rates reload; currency rates not updating; delete currency; backfill currency rates

Currency routes (list, init status, reload rates, initialize backfill, delete) require sign-in. Rights are checked in the handlers.

## Rate sheet routes need sign-in; rights checked in handlers

Lines 385-402

Asked as: rate sheet access; absorb rate sheet; rollout rate sheet; rate sheet permission denied

Rate sheet routes require sign-in. The comments say a job claim is needed to manage them, rollout is admin, and absorb is for admin or those with rate_sheet_revise. The handlers enforce this.

## Claims, work records and rate sheet entries need sign-in

Lines 404-420

Asked as: bulk assign claims; claim details denied; work records not loading; assignable users list

Claims routes (list, assignable users, bulk assign, details), work records routes and rate sheet entry update require sign-in. The comments say claim details and entry updates need the admin claim; the handlers check it.
