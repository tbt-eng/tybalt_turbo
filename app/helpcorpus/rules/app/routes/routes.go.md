<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"29a8082902019549eacccb4f34ab89d18b695097a733bd324db6290b0df3670f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need expenses editing turned on

Lines 14-31

Asked as: can't edit expense; expense changes blocked; purchase order edit not allowed; vendor edits disabled

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message from ErrExpensesEditingDisabled. Other collections are not checked by this function.

## Time editing needs the time feature flag turned on

Lines 33-43

Asked as: can't change timesheet hours; time entries locked; editing time disabled; timesheet edit blocked

When the time editing feature flag is off, Turbo refuses with the message from ErrTimeEditingDisabled. This check is a helper function; this file does not show which routes call it.

## Only holders of the report claim can view reports

Lines 45-54

Asked as: can't open report; report access denied; not authorized to view report; reports page forbidden

A person without the report claim is refused with "you are not authorized to view this report". This applies to every route under /api/reports, such as payroll, weekly time, expense and receipts reports, payables spreadsheets, branch mismatches and active jobs.

## Only holders of the kpi claim can view KPI reports

Lines 56-65

Asked as: kpi report access denied; can't see kpi numbers; kpi reports forbidden; kpi page not allowed

A person without the kpi claim is refused with "you are not authorized to view KPI reports". This applies to the KPI employee branch hours report under /api/kpi/reports.

## Turbo refuses to start if expenditure kinds config is invalid

Lines 89-91

Asked as: server won't start; expenditure config error; app crashes on startup; startup config invalid

At server start, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start serving.

## Version endpoint is public

Lines 93-105

Asked as: app version check; what build is this; git commit version; which release is running

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (long and short), git branch and build time.

## Time sheet routes require sign-in

Lines 107-131

Asked as: timesheet not loading logged out; time sheet login required; bundle timesheet; recall timesheet

All /api/time_sheets routes require a signed-in user: bundle, unbundle (recall), copy to next week, details, approve, reject, commit, uncommit, reviewers, approver, tallies and tracking. Committing is meant for committers; the claim check is in the handler, not this file.

## Time sheet tallies come in four kinds

Lines 120-123

Asked as: timesheet totals; pending hours count; approved tallies; shared timesheet counts

Tallies are available for the user's own time sheets (uid), pending for the approver, approved for the approver, and shared for the reviewer.

## Time tracking endpoints are for report holders

Lines 124-131

Asked as: missing timesheets list; who hasn't filed time; ignore missing timesheet; time tracking not showing

Tracking counts, weekly tracking lists, missing, ignored and not-expected lists, and adding or removing a missing-timesheet exception for a user all sit under /api/time_sheets. The comment says they are for report holders; the claim check is in the handlers.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: export timesheets csv; legacy export secret; jobs export api; expenses export unauthorized

The legacy exports of time sheets by week ending, jobs and expenses (both by updated-after) are not behind user sign-in. Their handlers do custom authentication with machine secrets.

## Radiator login reports use machine secrets

Lines 138-139

Asked as: radiator login report; radiator dashboard logins; radiator api key; radiator data feed

The radiator login report routes use machine-secret authentication with the radiator role instead of user sign-in.

## Expense routes require sign-in

Lines 145-175

Asked as: expense api unauthorized; submit expense; attachment missing; expense signed out error

All /api/expenses routes require a signed-in user: create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, uncommit requests, lists, details, attachments, attachment hash audit and replace, mark attachment missing, settle and reimbursements.

## Expense uncommit requests routes

Lines 156-183

Asked as: request to uncommit expense; decline uncommit request; cancel uncommit request; unlock committed expense

A signed-in user can create an uncommit request on an expense, list requests and reasons, approve, decline or cancel requests.

## Expense settlement and reimbursement routes

Lines 165-175

Asked as: expense repaid; reimbursement payroll; unsettled expenses; clear expense settlement

Expenses can be listed as unsettled or settled, settled, and have settlement cleared. A person sees their own expenses repaid on payroll, grouped by pay period, under reimbursements. There is also an org-wide expense commit queue.

## Absorb proposal routes: anyone proposes, absorb holders review

Lines 185-193

Asked as: propose merge; withdraw absorb proposal; decline absorb; suggest combining records

Any signed-in user can list, create, view, delete or withdraw absorb proposals. Declining is also routed here. The comment says absorb holders review them; the claim check is in the handlers.

## Absorb suggestions and history are for absorb holders

Lines 195-207

Asked as: duplicate clients; possible duplicate vendors; absorb history; who merged records

Absorb suggestions of probable duplicate clients and vendors can be listed and dismissed per collection. Absorb history shows committed absorbs with who absorbed and who committed, kept for ABSORB_HISTORY_RETENTION_DAYS.

## Time amendments can be committed through the API

Lines 209-211

Asked as: commit time amendment; amendment approval api; finalise amendment

POST /api/time_amendments/{id}/commit uses the shared commit handler for time amendments and needs sign-in.

## Admin profile routes

Lines 227-237

Asked as: change user identity; set user active; clear auth provider; save user claims

Signed-in users can reach routes to list identities, set active, save with claims, save limited, save identity and clear an authorized provider. The claim checks are in the handlers.

## Purchase order approve, reject, cancel, close, make cumulative, revise and group create need expenses editing

Lines 253-269

Asked as: approve purchase order fails; close po; cumulative purchase order; group purchase orders

These PO actions require sign-in and the expenses editing flag: approve, reject, cancel, close, make_cumulative, revise, group preview and group create. When the flag is off, Turbo refuses them. Reading POs (pending, visible, search, approvers, groups) does not check the flag in this file.

## Legacy purchase order routes skip the expenses editing flag check

Lines 271-275

Asked as: edit old purchase order; create legacy po; old po form

The legacy PO routes (edit view, create, update) require sign-in but are not behind the expenses editing flag middleware in this file.

## Absorb and undo-absorb routes exist for clients, contacts, vendors and rate sheets

Lines 286-301

Asked as: merge clients; merge vendor records; undo merge; rate sheet merge

Records of clients, client contacts, vendors and rate sheets can be absorbed into a target, and each has an undo absorb route. Rate sheet absorb is for admin, or absorb with rate_sheet_revise, per the comment.

## Job routes require sign-in

Lines 305-340

Asked as: create job; job wip; close job; adopt subjobs

All /api/jobs routes require sign-in. They include create, upsert, details, notes, time and expense summaries, POs, WIP, set status, set number, change client, adopt subjobs, rate sheet periods, close and validate proposal.

## WIP reports come in my, branch and division modes

Lines 342-346

Asked as: wip report; work in progress report; branch wip; division wip

Signed-in users can request a WIP report in three modes: my, branch or division.

## Reports routes require sign-in and the report claim

Lines 348-367

Asked as: report 403; payroll report forbidden; payables spreadsheet access; active jobs report access

Every route under /api/reports requires a signed-in user and then the report claim. The payroll time, weekly time, payroll and weekly expense, receipts, payables spreadsheet, branch mismatch and active jobs reports all fall under this check.

## KPI report routes require sign-in and the kpi claim

Lines 369-377

Asked as: kpi branch hours; employee hours kpi; kpi report login

The employee branch hours report under /api/kpi/reports requires a signed-in user and the kpi claim.

## Branch hours report checks access to the requested branch

Lines 379-387

Asked as: branch hours report denied; hours by branch access; see other branch hours

The /api/branch_hours report requires sign-in. It calls requireBranchHoursAccess with the branch given in the request's branch query parameter before it runs. The rules of that check are in another file.

## Machine secrets management is admin only

Lines 394-398

Asked as: api key management; create machine secret; secret keys admin; list machine secrets

Listing and creating machine secrets requires sign-in. The comment says it is admin only; the check is in the handlers.

## Attachment audit routes

Lines 400-407

Asked as: missing attachments csv; delete orphaned attachments; attachment audit run; attachment check report

Attachment audit routes list targets, show a target's run, refresh, delete orphans and download missing or orphaned CSV reports.

## Currency routes

Lines 409-415

Asked as: exchange rates reload; delete currency; currency backfill; currency setup status

Signed-in users reach routes to list currencies, check init status, reload rates, initialize backfill and delete a currency.

## Rate sheet management routes

Lines 417-436

Asked as: rate sheet rollout; delete rate sheet; duplicate rate sheet; rename rate sheet

Rate sheet routes cover create, compare, activate, deactivate, delete, rename, former names, jobs, rollout and duplicates. The comments say the job claim is required, delete is for an inactive sheet that nothing uses (admin), and rollout moves projects to the latest revision (admin).

## Claims routes: details need the admin claim

Lines 438-444

Asked as: bulk assign claims; claim details admin; list claims users; assignable users

Listing claims, assignable users and bulk assign need sign-in. The comment says the claim details endpoint requires the admin claim.

## Rate sheet entries update needs the admin claim

Lines 451-454

Asked as: edit rate entry; update rate sheet entry; rate entry forbidden; change rate sheet price

PUT /api/rate_sheet_entries/{id} requires sign-in. The comment says it requires the admin claim, which the handler checks.
