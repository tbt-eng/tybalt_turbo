<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"396fcaf3aec4a86e84a51affcfb998ab6c7cece920e5f34dab9fcf3ff3b08adb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expenses editing is off

Lines 14-31

Asked as: expenses locked; can't edit expenses; purchase orders read only; expense changes not allowed; vendor edits blocked

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors. It refuses with the message in ErrExpensesEditingDisabled. Other collections are not affected by this check.

## Time editing blocked when the time feature flag is off

Lines 33-43

Asked as: can't edit timesheet; time entries locked; timesheet changes not saving; time editing turned off

When the time editing feature flag is off, Turbo refuses the action with the message in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: can't open report; report access denied; not allowed to see reports; report permission missing

A person without the report claim cannot view a report. Turbo refuses with the message "you are not authorized to view this report".

## KPI reports need the KPI claim

Lines 56-65

Asked as: KPI report blocked; can't see KPI reports; KPI access denied; no KPI permission

A person without the KPI claim cannot view KPI reports. Turbo refuses with the message "you are not authorized to view KPI reports".

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; assistant slow first time; chat help bot; ask turbo a question

Turbo's Help chat answers questions at the route /api/help. It loads its model and its corpus the first time someone uses it.

## Turbo checks the expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup config error; server fails on boot

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which build is this; git commit hash; check turbo version

Anyone can read /api/version without signing in. It returns the app name, version, build, git commit and branch, and build time.

## Most Turbo API routes need a signed-in user

Lines 107-139

Asked as: api returns unauthorized; not logged in error; 401 from api; session required api

The time sheet, expense, purchase order, client, vendor, job, rate sheet, claim, currency and other API groups require a signed-in user. The legacy export routes are the exception. They use machine secrets for their own authorization.

## Legacy export endpoints use machine secrets, not user sign-in

Lines 133-136

Asked as: export link without login; legacy export secret; automated export access; export without sign in

The legacy exports of time sheets, jobs and expenses are not behind user sign-in. They use their own authorization with machine secrets.

## Expense actions available through the expenses routes

Lines 138-149

Asked as: submit expense button; who can approve expense; copy expense to tomorrow; recall my expense; commit expense

Signed-in users get routes to create and update an expense, and to submit, recall, approve, reject, commit and uncommit one. They also get a route to copy an expense to tomorrow. This file shows no check on who may do each action, so the handlers decide that.

## Expense attachment routes: view, audit hash, replace hash, mark missing

Lines 152-155

Asked as: download receipt; receipt hash audit; replace attachment hash; mark receipt missing

Expenses have routes to download an attachment, audit its hash, replace its hash and mark it missing. All of them need a signed-in user.

## Expense list, details and attachment download need sign-in

Lines 156-157

Asked as: expense list not loading; pending expenses api; approved expenses list; expense details unauthorized

The expense list, details, attachment download, pending and approved routes (/api/expenses/list, /details/{id}, /attachment/{id}, /pending, /approved) require a signed-in user. Lines 156-157 show no feature-flag check on these routes.

## Expense tracking, commit queue and settlement routes

Lines 158-165

Asked as: settle expense; clear settlement; commit queue; unsettled expenses list

Turbo offers expense tracking counts, an org-wide commit queue, unsettled and settled expense lists, and a tracking list by committed week ending. Settling an expense and clearing its settlement are separate POST actions (/settle, /clear_settlement). All need a signed-in user. This file does not show who may settle; the handlers decide that.

## Reimbursements list the caller's own expenses by pay period

Lines 166-168

Asked as: my reimbursements; expenses repaid on payroll; pay period expenses; reimbursed expenses by period

The reimbursement routes show the signed-in user's own expenses that are repaid on payroll, grouped by pay period ending. One route lists the pay periods and another lists the expenses of one pay period (/reimbursements/{payPeriodEnding}).

## Expense uncommit request routes

Lines 170-176

Asked as: approve uncommit request; decline uncommit; uncommit reasons; cancel uncommit request

Signed-in users can list uncommit requests and the uncommit reasons. They can also approve requests (POST /approve) and decline or cancel one request by id. This file does not show who may do each action; the handlers decide that.

## Time amendments can be committed through the API

Lines 178-180

Asked as: commit time amendment; amendment commit api; approve time change commit; commit amended timesheet

Turbo has a commit route for a time amendment (POST /api/time_amendments/{id}/commit). It needs a signed-in user and uses the same commit handler as time sheets and expenses.

## Copy a time entry to tomorrow

Lines 182-185

Asked as: copy time to tomorrow; duplicate timesheet entry; repeat time entry next day; copy entry tomorrow

A signed-in user can copy a time entry to tomorrow (POST /api/time_entries/{id}/copy_to_tomorrow).

## User defaults and navigation badges routes

Lines 187-193

Asked as: my user defaults; navigation badge counts; menu notification numbers; default settings for me

Signed-in users can fetch their defaults (/api/users/defaults) and their navigation badge counts (/api/nav/badges).

## Admin profile routes

Lines 195-205

Asked as: edit user profile admin; set user active; clear authorized provider; save user claims

Admin profile routes cover identity lists and details, active toggle targets, saving with claims, saving limited fields, setting active, saving identity, and clearing an authorized provider. All need a signed-in user. This file does not show who may use each one; the handlers decide that.

## Purchase order read routes ignore the expenses editing flag

Lines 207-219

Asked as: view purchase orders while locked; PO search when expenses off; list PO approvers; approval limits lookup

The purchase order routes for pending, visible, search, approvers, approval limits, second approvers and groups (GET) need only a signed-in user. They sit outside the group that checks the expenses editing flag, so the flag does not block them.

## Purchase order changes follow the expenses editing flag

Lines 221-243

Asked as: can't approve PO; PO actions blocked; revise purchase order locked; close PO not working

The purchase order actions approve, reject, cancel, close, make cumulative, revise, and group preview and create all check the expenses editing flag first. When the flag is off, Turbo refuses them. Listing and viewing purchase orders are not gated by this check. The legacy purchase order routes (edit, create, update) do not run this check in this file.

## Purchase order actions are blocked when expenses editing is off

Lines 221-237

Asked as: approve PO refused; PO expenses editing disabled; reject purchase order error; cumulative PO blocked

Approve, reject, cancel, close, make cumulative, revise, group preview and group create on purchase orders first check the expenses editing flag. If the flag is off, Turbo refuses with the message of ErrExpensesEditingDisabled.

## Legacy purchase order routes skip the expenses editing flag check

Lines 239-243

Asked as: edit old purchase order; legacy PO create; PO update works when locked; old purchase order form

The legacy purchase order routes (GET edit, POST create, PATCH update) need a signed-in user. They sit outside the group that checks the expenses editing flag.

## Client and vendor read routes

Lines 245-251

Asked as: view client list; vendor details; client notes; business development leads

Signed-in users can list and view clients and vendors. They can also view client notes and time divisions, and the list of business development leads.

## Absorb and undo absorb routes

Lines 252-266

Asked as: absorb client; undo absorb; merge vendor back; rate sheet absorb admin

Turbo offers absorb and undo absorb for clients, client contacts, vendors and rate sheets. Rate sheet absorb is for admins, or for people with the rate sheet revise claim, according to the comment in the code.

## Job routes

Lines 270-305

Asked as: create job; job status change; change job client; close job; adopt subjobs

Job routes need a signed-in user. They cover create, upsert, details, notes, time, staff, divisions, expense and PO summaries, WIP, set status, set number, change client, adopt subjobs, rate sheets, close and proposal validation. This file does not show who may do each action; the handlers decide that.

## WIP reports for my, branch and division

Lines 307-311

Asked as: WIP report my; branch WIP report; division WIP; work in progress report

Turbo offers three WIP report routes, /api/wip/my, /api/wip/branch and /api/wip/division. They need a signed-in user and no report claim in this file.

## Report routes require the report claim

Lines 313-332

Asked as: payroll report access; receipts report denied; payables spreadsheet; active jobs report

Every route under /api/reports requires a signed-in user with the report claim. This covers the payroll, weekly and receipts reports, the payables spreadsheets, the time entry branch mismatches report and the active jobs report.

## KPI report routes require the KPI claim

Lines 334-342

Asked as: branch hours KPI; employee hours KPI; KPI employee report; employee branch KPI

The employee branch hours KPI report needs a signed-in user with the KPI claim.

## Branch hours report has its own access check

Lines 344-352

Asked as: branch hours report denied; can't see branch hours; branch access check; branch hours permission

The branch hours report runs requireBranchHoursAccess for the branch in the query string before it returns data. This file does not show who passes that check.

## Machine secrets and claims detail are admin-only by comment

Lines 359-413

Asked as: machine secrets admin; claim detail admin only; who can manage secrets; rate sheet job claim

The comments say that machine secrets management is admin only, and that the claim detail endpoint requires the admin claim. This file does not itself enforce either check. Rate sheet and rate sheet entry routes carry comments about the job claim or admin claim, and this file does not enforce those either.
