<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"d38cbd4ecfee5783da65a69c4d02bc1ea4012ecf7e5b11d55f6a192185a0d228","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses locked; can't edit expenses; purchase orders blocked; vendor edits refused; expenses editing off

When the expenses editing feature flag is off, Turbo refuses actions on the collections expenses, purchase_orders and vendors. The refusal message is the one in ErrExpensesEditingDisabled. Other collections are not affected by this flag.

## Time editing refused when the time feature flag is off

Lines 33-43

Asked as: can't edit timesheet; time edits blocked; timesheet changes refused; time editing off

When the time editing feature flag is off, Turbo refuses the action. The refusal message is the one in ErrTimeEditingDisabled. This helper is defined in this file, but the routes shown here do not call it.

## Reports need the report claim

Lines 45-54

Asked as: can't open reports; report access denied; no permission for reports; not authorized to view report

A person without the report claim cannot use the report endpoints. Turbo refuses with 'you are not authorized to view this report'. This applies to every route under /api/reports, such as payroll, weekly time, expense and receipts reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: kpi report denied; can't see kpi reports; kpi access refused; kpi permission missing

A person without the kpi claim cannot use the KPI report endpoints. Turbo refuses with 'you are not authorized to view KPI reports'. The employee branch hours report is under this rule.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; where is help chat; help bot slow first time; assistant endpoint url

Turbo registers the Help chat route at /api/help when it adds its custom routes. The chat loads its model and corpus on first use, not at startup.

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup config error; server not serving requests

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start serving.

## Version endpoint is public

Lines 93-105

Asked as: app version info; build number endpoint; git commit version; which version is running

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Most custom routes require a signed-in user

Lines 107-118

Asked as: 401 on api routes; signed out api error; need to log in; not signed in error

Time sheet, expense, expense uncommit request, absorb, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, report, KPI, branch hours, stats, machine secret, attachment audit, currency, rate sheet, claim and work record routes all require a signed-in user. Individual handlers then apply their own claim checks.

## Time sheet tracking endpoints are for report holders

Lines 124-131

Asked as: missing timesheets list; timesheet exceptions ignored; tracking counts; add missing timesheet exception

The time sheet tracking routes are tracking counts, the weekly list, missing, ignored and not expected. The routes that add or remove a missing time sheet exception for a user are also in this group. Their comment says they are for report holders. This file does not check the claim, and the handlers do.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: legacy export failing; export without login; machine secret export; old export route error

The legacy export routes for time sheets, jobs and expenses do not use the user sign-in rule. Their handlers check machine secrets.

## Radiator and computer routes have their own access

Lines 138-143

Asked as: radiator login report; computer assign retire; computers claim access; retire computer error

Radiator login reports use machine secrets with the role radiator. Computer routes (list, details, assign, retire, un-retire) are for holders of the computers or admin claim, as the comment says. The handlers apply the checks, not this file.

## Expense workflow actions available

Lines 145-156

Asked as: submit expense; recall expense; copy expense to tomorrow; approve reject expense

Expenses can be created, edited, submitted, recalled, copied to tomorrow, approved, rejected, committed and uncommitted through the API. Staff can ask for an uncommit through an uncommit request. These actions are not gated by the expenses editing flag in this file.

## Expense settlement and attachment routes

Lines 157-177

Asked as: settle expense; clear settlement; replace attachment hash; mark attachment missing

Expenses have routes to settle, clear a settlement, and list settled and unsettled expenses. Others audit or replace an attachment hash, and mark an attachment missing. The commit queue is org-wide. Reimbursements are listed for the caller's own expenses, grouped by pay period.

## Expense uncommit request routes

Lines 179-185

Asked as: uncommit request approve; decline uncommit; cancel uncommit request; uncommit reasons list

Uncommit requests for expenses can be listed, and their reasons listed. A request can be approved (in bulk), declined or cancelled.

## Absorb proposals: anyone proposes, absorb holders review

Lines 187-195

Asked as: propose absorb; withdraw absorb proposal; delete absorb proposal; review absorb proposals

Any signed-in user can propose an absorb. A proposal can be listed, viewed, deleted, withdrawn or declined. Absorb holders review proposals.

## Absorb suggestions and history are for absorb holders

Lines 197-209

Asked as: duplicate clients suggestions; absorb history; dismiss duplicate vendor; how long absorb history kept

Absorb suggestions show probable duplicate clients and vendors, and can be dismissed. Absorb history lists committed absorbs, with who absorbed and who committed. History is kept for ABSORB_HISTORY_RETENTION_DAYS.

## Time amendments can be committed through the API

Lines 211-213

Asked as: commit time amendment; amendment commit api; time amendments commit

Turbo adds a commit route for time amendments (POST /api/time_amendments/{id}/commit) for signed-in users. It uses the same commit handler as time sheets and expenses. This file does not show who may commit.

## Time entry copy and CSV routes need sign-in

Lines 215-219

Asked as: copy time entry tomorrow; download time csv; export my time entries

A signed-in user can copy a time entry to tomorrow and download their own time entries as CSV (/mine/csv). Signed-out requests get no access.

## User defaults and navigation badges need sign-in

Lines 221-227

Asked as: user defaults; navigation badges; nav badge counts; menu badges not showing

Turbo serves a user's defaults (/api/users/defaults) and the navigation badges (/api/nav/badges) only to signed-in users.

## Admin profile routes need sign-in and have separate save routes

Lines 229-239

Asked as: save admin profile; set profile active; clear authorized provider; admin profile claims

All admin profile routes need a signed-in user. Turbo has separate routes to save an admin profile with claims, to save it in a limited form, to set it active, to save its identity, and to clear an authorized provider. This file does not show which claim each route requires.

## Purchase order changes need expenses editing enabled

Lines 241-271

Asked as: can't approve purchase order; close po blocked; revise purchase order; create po group refused

Approve, reject, cancel, close, make cumulative and revise on a purchase order are refused when the expenses editing flag is off. Creating a purchase order group, or previewing one, is refused too. Read-only PO routes (pending, visible, search, approvers, groups) are not gated by this flag.

## Legacy purchase order routes are not gated by the flag

Lines 273-277

Asked as: legacy purchase order edit; create old po; update legacy po; old po form

Legacy purchase order routes (edit view, create, update) require a signed-in user. They sit outside the expenses editing flag check in this file.

## Absorb and undo-absorb routes by record type

Lines 289-303

Asked as: undo absorb client; undo absorb vendor; absorb rate sheet; undo absorb contact

Absorb and undo-absorb routes exist for clients, client contacts, vendors and rate sheets. Rate sheets: the comment says admin, or absorb with rate_sheet_revise. This file does not check the claims.

## Job routes

Lines 307-342

Asked as: job notes; close job; adopt subjob; unused jobs list; job status change

Job routes cover create, upsert, details and notes. They also cover time, staff, division, expense and PO summaries, WIP, status, number and client changes, and rate sheet periods. Subjobs can be adopted. A job can be closed and a proposal validated. Unused and stale job lists exist.

## Work in progress (WIP) reports come in three modes

Lines 344-348

Asked as: wip my report; branch wip; division wip; work in progress report

Turbo serves the WIP report in three modes: my, branch and division. Each mode has its own route (/api/wip/my, /branch, /division) for signed-in users.

## Reports need sign-in and the report claim

Lines 350-369

Asked as: payroll report access; receipts report denied; active jobs report; payables spreadsheet

All routes under /api/reports require a signed-in user who has the report claim. Anyone else sees "you are not authorized to view this report". These routes include the payroll, weekly and receipts reports, the payables spreadsheets, the time entry branch mismatches report and the active jobs report.

## KPI reports need sign-in and the kpi claim

Lines 371-379

Asked as: kpi reports not authorized; kpi report access; employee hours kpi

The routes under /api/kpi/reports, such as the employee branch hours report, require a signed-in user who has the kpi claim. Anyone else sees "you are not authorized to view KPI reports".

## Branch hours report checks access for the requested branch

Lines 381-389

Asked as: branch hours report denied; branch query parameter; hours report wrong branch

Before the branch hours report runs, Turbo checks the user's access for the branch given in the 'branch' query parameter. It uses requireBranchHoursAccess, which is defined in another file.

## Machine secrets and stats routes need sign-in

Lines 391-401

Asked as: stats dashboard; machine secrets list; create machine secret

Signed-in users reach the stats dashboard (/api/stats) and the machine secrets list and create routes. A comment says machine secrets are admin only. The check is in the handler.

## Machine secrets, claims and stats endpoints

Lines 396-401

Asked as: claim details admin; machine secret admin only; rate sheet entry admin

The machine secrets routes are described as admin only. The claim detail endpoint requires the admin claim, according to the comment. Rate sheet entry updates need the admin claim, according to the comment. This file does not check these claims.

## Attachment audit routes need sign-in

Lines 402-409

Asked as: attachment audit run; orphaned attachments csv; missing attachments report; delete orphan attachments

Signed-in users reach the attachment audit routes. They list targets, show a run, refresh it, delete orphans, and download the missing and orphaned CSV reports. This file does not show which claim they need.

## Currency routes need sign-in

Lines 411-417

Asked as: reload currency rates; currency backfill; delete currency; currency initialization status

Signed-in users reach the currency routes. They list currencies, show initialization status, reload rates, start a backfill for one currency and delete a currency. This file does not show which claim they need.

## Rate sheet management routes

Lines 419-438

Asked as: create rate sheet; compare rate sheets; rollout rate sheet; activate rate sheet; delete inactive rate sheet

Rate sheet routes include create, compare, activate, deactivate, rename, delete and rollout. The comments say the job claim is required. They say delete applies to an inactive sheet that nothing uses (admin). They say rollout moves projects to the latest revision (admin). This file does not check the claims.

## Claims routes need sign-in; detail needs admin

Lines 439-446

Asked as: list claims; bulk assign claim; claim details; assignable users for claim

Signed-in users can list claims, list assignable users for a claim, bulk assign a claim and open a claim's details. A comment says the detail endpoint needs the admin claim. The check is in the handler.

## Work records routes need sign-in

Lines 448-451

Asked as: work records list; open work record; work record details

A signed-in user can list work records and open one record's details. Signed-out requests get no access.

## Rate sheet entry updates need sign-in

Lines 453-456

Asked as: update rate entry; change rate sheet entry; edit rate sheet price

Turbo lets a signed-in user update a rate sheet entry (PUT /api/rate_sheet_entries/{id}). The comment says this needs the admin claim. The route check itself is in the handler, not in this file.
