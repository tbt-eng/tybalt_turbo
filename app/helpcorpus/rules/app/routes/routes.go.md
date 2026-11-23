<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"d22a0576e4e9156a412aedce75da7b93d03c29a09e0300603486538aa86508cb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: can't edit expenses; expense changes blocked; purchase order edit refused; vendor edit not allowed; expenses editing disabled

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message of ErrExpensesEditingDisabled. Other collections are not blocked by this check.

## Time editing flag off blocks time actions

Lines 33-43

Asked as: can't change time entries; time sheet edits blocked; time editing turned off; cannot update timesheet; time changes refused

When the time editing feature flag is off, Turbo refuses with the message of ErrTimeEditingDisabled. This is the check in requireTimeEditing, which a handler must call.

## Reports need the report claim

Lines 45-54

Asked as: can't open reports; report access denied; not authorized to view report; no permission for reports; reports page forbidden

Staff without the report claim cannot open the report routes. Turbo refuses them with 'you are not authorized to view this report'. The check applies to every route under /api/reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: can't see KPI report; KPI access denied; employee branch hours blocked; kpi report not authorized; no KPI permission

Staff without the kpi claim cannot open the KPI report routes. Turbo refuses them with 'you are not authorized to view KPI reports'. The employee branch hours report is under this check.

## Help chat is served at /api/help

Lines 84-85

Asked as: help chat not loading; where is the help bot; chatbot slow to start; help model loading; assistant help endpoint

Turbo serves the Help chat at /api/help. The chat loads its model and its corpus on first use, not at start-up.

## Turbo checks expenditure kinds config at start-up

Lines 88-91

Asked as: server won't start; startup fails expenditure; expenditure kinds invalid; app not starting config; bad expenditure config

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start serving.

## Version endpoint is public

Lines 93-105

Asked as: what version is turbo; app build number; git commit of turbo; version check without login; which build is running

Anyone can call GET /api/version without signing in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Most Turbo API routes require sign-in

Lines 107-126

Asked as: api needs login; unauthenticated api call; 401 on turbo api; must sign in to call; who can use this endpoint

Time sheet, expense, PO, job, client, vendor, rate sheet, claim and other custom routes require a signed-in user. Each handler decides who may do the action.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-139

Asked as: export without login; legacy export access; machine secret export; radiator login report; export key authorization

The legacy export routes for time sheets, jobs and expenses do not require a user. They use their own authorization through machine_secrets. Radiator login reports use machine secrets with the radiator role.

## Computers routes are for computers or admin claim holders

Lines 141-143

Asked as: computers page forbidden; who can retire computer; assign computer permission; computer claim needed; can't unretire computer

Turbo registers list, details, assign, retire and un-retire routes for computers. The code comment says they need the computers claim or the admin claim. The check is in another file, so this file does not show the refusal.

## Expense routes require sign-in

Lines 145-175

Asked as: expense api login; submit expense not signed in; expense route unauthorized; recall expense access; reimbursement endpoint sign in

Every route under /api/expenses requires a signed-in user. The routes cover create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, uncommit requests, details, attachments, settlement and reimbursements.

## Expense uncommit request routes require sign-in

Lines 177-183

Asked as: uncommit request access; approve uncommit request; cancel uncommit request; uncommit reasons list; who can decline uncommit

The routes under /api/expense_uncommit_requests need a signed-in user (users collection). They list requests and reasons, approve requests, decline a request and cancel a request. Turbo only registers the routes here; each handler decides who may act.

## Time amendment commit route requires sign-in

Lines 185-187

Asked as: commit time amendment; amendment commit who can; approve time amendment commit; time amendment locked in; commit amendment permission

A signed-in user can call the commit action for a time amendment (POST /api/time_amendments/{id}/commit). It uses the same commit handler as time sheets and expenses. This file does not say who may commit; the handler decides.

## Time entry copy-to-tomorrow requires sign-in

Lines 189-192

Asked as: copy time entry tomorrow; copy entry next day; duplicate timesheet line tomorrow; copy to tomorrow fails; time entry copy access

Copying a time entry to tomorrow (POST /api/time_entries/{id}/copy_to_tomorrow) needs a signed-in user.

## User defaults and navigation badges require sign-in

Lines 194-200

Asked as: user defaults not loading; nav badges missing; navigation counts not showing; badge numbers unavailable; my defaults endpoint

The user defaults route (/api/users/defaults) and the navigation badges route (/api/nav/badges) need a signed-in user.

## Admin profile routes require sign-in

Lines 202-212

Asked as: admin profile access; who can save profile; claims save permission; deactivate profile action; clear authorized provider

The admin profile routes need a signed-in user. They cover identity lists, active toggle targets, save with claims, limited save, setting a profile active, saving identity, and clearing an authorized provider. This file adds no claim check, so the handlers decide who may use each one.

## Purchase order read routes have no editing flag check

Lines 214-226

Asked as: view purchase orders when locked; PO list with flag off; search purchase orders read; see POs while editing off; approval limits lookup

The purchase order routes that read data need only a signed-in user. They cover pending, visible, search, approvers, approval limits, second approvers and groups. They do not check the expenses editing flag, so people can still view purchase orders when that flag is off.

## Purchase order mutations need the expenses editing flag

Lines 228-244

Asked as: can't approve purchase order; PO approve blocked; reject PO refused; close purchase order fails; revise PO not allowed

Approve, reject, cancel, close, make cumulative and revise on a purchase order, and previewing or creating a PO group, are refused while the expenses editing flag is off. Reading POs and PO groups does not check the flag. Legacy PO routes also do not check it here.

## Legacy purchase order routes skip the expenses editing flag

Lines 246-250

Asked as: legacy PO edit access; create legacy purchase order; legacy purchase order update; fetch legacy PO for edit; old PO routes

Turbo's legacy purchase order routes (/api/purchase_orders/legacy) need a signed-in user: fetch for edit, create and update. They are in a separate group from the mutation routes, so this file does not apply the expenses editing flag check to them.

## Client routes require sign-in

Lines 252-260

Asked as: client list not loading; view client notes; client rate sheets access; time divisions for client; client page sign in

The client routes need a signed-in user. They list or view clients, notes, time divisions and rate sheets. Absorb and undo absorb also sit in this group.

## Absorb and undo absorb routes

Lines 259-273

Asked as: absorb client; undo absorb vendor; rate sheet absorb permission; absorb contacts into client; who can absorb

Absorb and undo absorb exist for clients, client contacts, vendors and rate sheets. Rate sheet absorb is for admins, or for staff with the rate_sheet_revise claim.

## Project request routes are registered separately

Lines 275-275

Asked as: project request access; job request routes; project request permissions; create project request; project request rules

Turbo registers the project request routes through registerJobRequestRoutes. This file does not show their rules.

## Job routes require sign-in

Lines 277-312

Asked as: job api not working; who can close job; job status change access; adopt subjobs; job proposal validation

All job routes need a signed-in user. They cover create, upsert, details, notes, time, staff, division, expense and PO summaries, WIP, status, number, client change, adopt subjobs, rate sheet periods, close and proposal validation. Who may do each action is decided by the handlers.

## WIP report routes: my, branch and division

Lines 314-318

Asked as: my WIP report; branch WIP view; division WIP access; wip report modes; WIP not showing for me

A signed-in user can reach three WIP report routes: /api/wip/my, /api/wip/branch and /api/wip/division. This file adds no claim check, so the handler decides what each mode shows.

## Reports routes require the report claim

Lines 320-339

Asked as: payroll report access; weekly time report blocked; payables spreadsheet access; active jobs report forbidden; receipts report permission

Every route under /api/reports requires sign-in and then the report claim. This covers the payroll, weekly time, expense and receipts reports, the payables spreadsheets, the time entry branch mismatches report and the active jobs report.

## KPI report routes require sign-in and the kpi claim

Lines 341-349

Asked as: KPI report login; kpi reports forbidden; employee hours kpi access; kpi claim needed; kpi route sign in

The employee branch hours KPI report at /api/kpi/reports requires sign-in and the kpi claim.

## Branch hours report has its own access check

Lines 351-359

Asked as: branch hours report denied; who sees branch hours; branch query access check; branch hours permission; branch report forbidden

The branch hours report requires sign-in. It then calls requireBranchHoursAccess with the branch query parameter. That function is in another file, so the rule it applies is not stated here.

## Machine secrets routes are for admins

Lines 366-370

Asked as: create machine secret; list machine secrets; who manages machine keys; machine secret admin only; secrets page access

Listing and creating machine secrets is meant for admins. The route group only requires sign-in, so the handlers enforce the admin claim.

## Attachment audit routes require sign-in

Lines 372-379

Asked as: attachment audit access; download orphaned attachments CSV; missing attachments report; delete orphan files; attachment check permission

The attachment audit routes need a signed-in user. They list targets, show a run, refresh, delete orphans, and download the missing and orphaned CSV reports. This file adds no claim check, so the handlers decide who may use them.

## Currency routes require sign-in

Lines 381-387

Asked as: currency reload rates; currency backfill access; delete currency permission; currency list not loading; who can reset currencies

The currency routes need a signed-in user. They list currencies, show initialization status, reload rates, start a backfill and delete a currency. This file adds no claim check, so the handlers decide who may do each.

## Rate sheet delete and rollout are admin actions

Lines 389-408

Asked as: delete rate sheet; rate sheet rollout; move projects to latest revision; rate sheet admin only; manage rate sheets access

Deleting a rate sheet is for an admin and works only on an inactive sheet that nothing uses. Moving projects to the latest revision (rollout) is also for an admin. Managing rate sheets needs the job claim.

## Claim details and rate sheet entry updates need admin claim

Lines 410-426

Asked as: claim details forbidden; update rate sheet entry; claim info admin only; edit rate entry permission; claim details access

The claim details endpoint requires the admin claim. Updating a rate sheet entry also requires it.
