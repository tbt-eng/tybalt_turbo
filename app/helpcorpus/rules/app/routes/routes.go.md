<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"df5c6860044df6191d6c11d1a8630513ad8d45233a777cd7ba8c4cb1fe61513a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors editing can be switched off

Lines 14-31

Asked as: expenses editing off; can't edit expenses; purchase orders locked; vendors read only; editing disabled error

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors that use this check. The error message is the one in ErrExpensesEditingDisabled. Other collections are not affected by this check.

## Time editing can be switched off

Lines 33-43

Asked as: time editing disabled; can't edit timesheet; time entries locked; time editing off

When the time editing feature flag is off, Turbo refuses the action that uses this check. The error message is the one in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: not authorized to view report; report access denied; can't open reports; no report permission

A person without the report claim cannot use any report under /api/reports. Turbo answers 'you are not authorized to view this report'. These routes need a signed-in user first.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI reports denied; can't see KPI; not authorized KPI; KPI access missing

A person without the kpi claim cannot use the KPI reports under /api/kpi/reports. Turbo answers 'you are not authorized to view KPI reports'.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; help assistant; chat help page; help model loading slow

Turbo serves the Help chat at /api/help. The chat loads its model and corpus the first time someone uses it.

## Turbo checks expenditure kinds config at start

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup config error; expenditure kinds config

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: app version number; build info without login; version page; git commit build time

Anyone can read /api/version without signing in. It returns the app name, version, build, git commit, git branch and build time.

## Most Turbo routes need a signed-in user

Lines 107-136

Asked as: need to log in; login required routes; not signed in error; session required api

Time sheet, expense, purchase order, job, client, vendor, report, rate sheet and claim routes need a signed-in user. The legacy export routes use machine secrets instead.

## Time sheet commit is for committers

Lines 115-117

Asked as: who can commit timesheets; commit button missing; can't commit time sheet; timesheet committer permission

Time sheets are committed and uncommitted through routes that need sign-in. The comment says committers may commit time sheets. The claim check runs in the handler, not in this file.

## Radiator login report routes use machine secrets, not user sign-in

Lines 133-139

Asked as: radiator login report; machine secret export; legacy export access; radiator role login

The legacy export routes for time sheets, jobs and expenses do not use the signed-in user check. They use their own authentication with machine secrets. The radiator login report routes also use machine secrets, with the role radiator.

## Computer routes need computers or admin claim

Lines 141-143

Asked as: computer list access; can't assign computer; retire computer permission; computers claim needed

Turbo registers the routes that list computers, show details, assign, retire and un-retire a computer in bindComputerRoutes. The comment says a person needs the computers claim or the admin claim. This file does not show the check itself.

## Expense routes need sign-in; no flag check here

Lines 145-175

Asked as: expense actions need login; expense route access; expense flag not checked; sign in for expenses

Every route under /api/expenses requires a signed-in user. The expense routes cover create, edit, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, attachment checks, settle and clear settlement. This file does not check the expenses editing flag for them, unlike the purchase order mutation routes.

## Uncommit requests for expenses

Lines 156-183

Asked as: ask to uncommit expense; request expense uncommit; approve uncommit requests; decline uncommit request

A person creates an expense uncommit request with POST /api/expenses/{id}/uncommit_requests. The /api/expense_uncommit_requests routes list requests and reasons, approve requests in bulk, and decline or cancel one request by id. All need sign-in.

## Expense commit queue and tracking are org-wide lists

Lines 165-172

Asked as: org-wide commit queue; expense tracking list; all expenses committed week; expense tracking counts

The expense routes include an org-wide commit queue, tracking counts and a tracking list by committed week ending. The code comment describes them as expense tracking plus the org-wide commit queue. Who may see them is decided in the handlers, not in this file.

## Expense settlement routes

Lines 168-171

Asked as: settle an expense; clear expense settlement; unsettled expenses list; who can settle

Expense routes let a person list unsettled and settled expenses, settle an expense, and clear its settlement. Who may settle is decided in the handlers, not in this file.

## Reimbursements list shows only your own expenses

Lines 173-175

Asked as: my reimbursements; expenses repaid on payroll; reimbursement by pay period; my repaid expenses

The reimbursements routes list the caller's own expenses that are repaid on payroll, grouped by pay period. A second route lists the reimbursements for one pay period ending.

## Absorb proposals: anyone proposes, absorb holders review

Lines 185-193

Asked as: propose absorb; absorb proposal review; withdraw absorb proposal; decline absorb proposal

Any signed-in user can list and create absorb proposals and open one. Routes also exist to delete, withdraw and decline a proposal. The comment says absorb holders review them. This file does not show those checks.

## Time amendments commit route

Lines 195-197

Asked as: commit time amendment; time amendment commit; amendment commit fails

A time amendment is committed with POST /api/time_amendments/{id}/commit. The route needs sign-in and uses the same commit handler as time sheets and expenses.

## Admin profile routes need sign-in

Lines 213-223

Asked as: admin profile save; identity save admin; clear authorized provider; admin claims save

The admin profile routes cover identity lists, save with claims, limited save, active toggle, identity save and clearing an authorized provider. All need sign-in. This file does not show the claim checks for them.

## Purchase order read routes are not gated by the flag

Lines 225-237

Asked as: view purchase orders off; PO search still works; read purchase orders; PO approvers list

The purchase order routes for pending, visible, search, approvers, approval limits, second approvers and groups are read-only. They need sign-in and no expenses editing check. Reading still works when expenses editing is off.

## Purchase order changes are blocked when expenses editing is off

Lines 239-261

Asked as: approve PO blocked; PO editing disabled; can't close purchase order; revise PO refused

Approve, reject, cancel, close, make cumulative and revise on a purchase order, and creating or previewing a purchase order group, first check the expenses editing flag. If it is off, Turbo refuses them. Viewing purchase orders and the legacy purchase order routes do not use this check.

## Purchase order mutation routes

Lines 239-255

Asked as: approve reject purchase order; make PO cumulative; create PO group; PO group preview

The routes for approving, rejecting, cancelling, closing, making cumulative, revising and creating a PO group are gated by the expenses editing check. Group preview is gated too.

## Legacy purchase order routes skip the editing flag

Lines 257-261

Asked as: legacy purchase order edit; old PO create route; legacy PO update; PO edit view

The legacy purchase order routes (edit view, create, update) need sign-in. They sit in a separate group without the expenses editing check, so the flag does not block them in this file.

## Client usage and invoicing profile replace routes

Lines 270-272

Asked as: where client contact used; invoicing profile replace; client usage lookup; replace invoicing profile

Clients have routes to view where a contact or an invoicing profile is used. Another route replaces an invoicing profile. All need sign-in.

## Absorb routes for clients, contacts, vendors and rate sheets

Lines 273-287

Asked as: absorb client contact; undo absorb vendor; absorb rate sheet; absorb permission clients

Absorb and undo-absorb routes exist for clients, client contacts, vendors and rate sheets. All need sign-in. The absorb permission is checked in the handlers, not in this file.

## Job routes need sign-in

Lines 291-326

Asked as: job routes sign in; adopt subjobs; close job; job WIP summary

All /api/jobs routes need sign-in. They cover create, upsert, set status, set number, change client, adopt subjobs, close, validate proposal, rate sheets, and the time, expense, PO and WIP summaries. This file does not show the claim checks.

## WIP reports have three modes

Lines 328-332

Asked as: WIP my report; WIP branch report; WIP division report; work in progress report

Turbo serves WIP reports at /api/wip/my, /api/wip/branch and /api/wip/division. Each needs sign-in.

## Report claim covers all report routes

Lines 334-353

Asked as: payroll spreadsheet report; payables report access; active jobs report; receipts report denied

Payroll, weekly, expense, receipts and payables spreadsheet reports, the time entry branch mismatch report and the active jobs report all need the report claim.

## Branch hours report has its own access check

Lines 365-373

Asked as: branch hours report access; branch hours denied; branch hours parameter

The branch hours report checks access with requireBranchHoursAccess, using the branch query parameter. The result of that check is defined in another file.

## Admin stats dashboard needs sign-in

Lines 375-378

Asked as: admin stats dashboard; stats page access; who sees stats

GET /api/stats needs sign-in. The comment calls it the admin stats dashboard. This file does not show who may view it.

## Machine secrets routes are admin only

Lines 380-384

Asked as: create machine secret; list machine secrets; machine secrets admin

The routes to list and create machine secrets need sign-in. The comment says they are for admins only. The admin check runs in the handlers.

## Attachment audit routes

Lines 386-393

Asked as: attachment audit; orphaned attachments; missing attachments CSV; delete orphan files

The attachment audit routes list targets, show a run, refresh it, delete orphans and download missing and orphaned CSV reports. All need sign-in.

## Currency routes

Lines 395-401

Asked as: reload currency rates; currency backfill; delete currency; currency init status

The currency routes list currencies, show init status, reload rates, start a backfill and delete a currency. All need sign-in.

## Rate sheet routes and who may use them

Lines 403-422

Asked as: delete rate sheet; rate sheet rollout; absorb rate sheet admin; inactive rate sheet delete

The rate sheet routes need sign-in. The comments say creating needs the job claim. Deleting needs admin and works only for an inactive sheet that nothing uses. Rollout to the latest revision is admin. Absorb is for admin, or for absorb with rate_sheet_revise.

## Claim detail route needs admin claim

Lines 424-430

Asked as: claim detail admin; bulk assign claims; assignable users list; claims list access

The claims routes need sign-in. The comment says the claim detail endpoint needs the admin claim. The list, assignable users and bulk assign routes are also here.

## Work record routes need sign-in

Lines 432-435

Asked as: work records list; open work record; work record by id

Work records can be listed and opened by record through /api/work_records. Both need sign-in.

## Rate sheet entry update needs admin claim

Lines 437-440

Asked as: update rate sheet entry; edit rate entry; rate sheet entry admin

PUT /api/rate_sheet_entries/{id} updates a rate sheet entry. The comment says it needs the admin claim. The route also needs sign-in.
