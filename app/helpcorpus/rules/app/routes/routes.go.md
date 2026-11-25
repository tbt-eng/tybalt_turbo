<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"69d3a686729ec9dc0f2ca6390b6473e98511f4264e282f1e067320bc6a8a03e1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, POs and vendors refuse edits when the expenses feature flag is off

Lines 14-31

Asked as: can't edit expense; purchase order edit blocked; vendor edits disabled; expense editing turned off

When the expenses editing flag is off, Turbo refuses requests on expenses, purchase orders and vendors that go through requireExpensesEditing. It shows the message in ErrExpensesEditingDisabled. Other collections pass this check.

## Time editing flag off blocks time editing

Lines 33-43

Asked as: can't edit timesheet; time entry edit blocked; timesheet editing disabled; time changes refused

When the time editing flag is off, a check that uses requireTimeEditing refuses the request. It shows the message in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: can't open reports; report access denied; no permission for reports; reports page forbidden

A person without the report claim cannot open reports under /api/reports. Turbo says: "you are not authorized to view this report".

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see kpi reports; kpi dashboard not allowed; kpi page forbidden

A person without the kpi claim cannot open KPI reports under /api/kpi/reports. Turbo says: "you are not authorized to view KPI reports".

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds config error; app fails to boot; startup config error

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start.

## Custom endpoints require a signed-in user, except version and machine routes

Lines 93-139

Asked as: api returns unauthorized; need to log in to api; 401 on api call; api asks for sign in

The time sheet, expense, purchase order, client, vendor, job, report and other API groups require a signed-in user (users). The version endpoint /api/version is public. The legacy export endpoints use machine secrets instead of a user sign-in.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which build am I on; git commit version; version info endpoint

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Time sheet routes need sign-in; commit is open to committers

Lines 107-131

Asked as: can't commit timesheet; commit button not working; timesheet approve reject access; uncommit timesheet permission

All time sheet routes need a signed-in user. They cover bundle, unbundle, copy to next week, details, approve, reject, commit, uncommit, reviewers and approver. The commit route is meant for committers (the commit claim), and the handler enforces it. The tracking routes are for report holders.

## Time sheet tallies come in four views

Lines 120-123

Asked as: timesheet totals; pending timesheets for approver; tally counts; timesheets shared with me

Tallies are available as the user's own (by uid), pending for the approver, approved for the approver, and shared with the user as reviewer.

## Missing time sheet exceptions can be added and removed

Lines 125-131

Asked as: exclude user from missing timesheets; missing timesheet exception; ignore missing timesheet; list missing timesheets for week

For a week ending, a missing time sheet exception can be created with POST and deleted with DELETE for a given user. The tracking views for a week list time sheets, missing, ignored and not expected users.

## Legacy export endpoints use machine secrets, not user sign-in

Lines 133-136

Asked as: export api no login; legacy export token; csv export secret; export endpoint without sign in

The legacy export routes for time sheets, jobs and expenses are not behind user sign-in. Their handlers use their own authentication through machine secrets. Their results are not stated in this file.

## Other routes

Lines 138-218

Asked as: copy time entry to tomorrow; navigation badge counts; computer access denied; radiator login report

Signed-in users reach routes for time entry copy to tomorrow and own CSV, user defaults, navigation badges and admin stats. Computer routes need the computers or admin claim. Radiator login reports use machine secrets with role radiator.

## Expense routes need sign-in and cover the full lifecycle

Lines 145-175

Asked as: submit expense; recall expense; settle expense; expense reimbursements list

All /api/expenses routes need a signed-in user. They cover create, update, submit, recall, copy to tomorrow, approve, reject, commit and uncommit. They also cover uncommit requests, attachment hash audit and replace, marking an attachment missing, settle and clear settlement, plus tracking, the commit queue and reimbursements. The handlers decide who may do each.

## Expense and PO edits are not gated by the feature flag in these routes

Lines 145-268

Asked as: expense edit works with flag off; edit PO when editing disabled; create purchase order still allowed; expense edits not blocked

The expense routes and the legacy purchase order routes (edit, create, update) are not behind the expenses editing middleware in this file. Only the PO approve, reject, cancel, close, make cumulative, revise and group routes are.

## Expense uncommit requests have approve, decline and cancel actions

Lines 177-183

Asked as: approve uncommit request; decline expense unlock; cancel uncommit request; uncommit reasons list

Under /api/expense_uncommit_requests, a signed-in user can list requests and reasons, approve requests in bulk, and decline or cancel one request by id. The handlers set who may do each.

## Anyone proposes absorbs; absorb holders review

Lines 185-193

Asked as: propose absorb; withdraw absorb proposal; decline absorb proposal; merge proposal review

Absorb proposals can be listed, created, viewed, deleted, withdrawn and declined by signed-in users. The route comment says anyone proposes and absorb holders review. The handlers enforce the details.

## Absorb suggestions are for absorb holders

Lines 195-200

Asked as: duplicate clients suggestions; dismiss duplicate vendor; possible duplicate clients; suggested duplicate records

Absorb suggestions of probable duplicate clients and vendors can be listed and dismissed per collection. The route comment says they are for absorb holders.

## Time amendments are committed through their own route

Lines 202-204

Asked as: commit time amendment; amendment not committed; approve timesheet amendment; commit amended time

A signed-in user can POST to commit a time amendment. The shared commit handler applies the rules.

## Admin profile routes

Lines 220-230

Asked as: save user claims; set identity active; clear authorized provider; manage user profiles

Signed-in users reach routes to list identities, save a profile with claims, save a limited profile, set active, edit identity, and clear an authorized provider. The handlers decide who may use each.

## Purchase order read routes are not behind the expenses flag

Lines 232-244

Asked as: view purchase orders with flag off; pending POs list; PO approvers list; search purchase orders

The PO routes for pending, visible, search, approvers, approval limits, second approvers and groups need sign-in only. The expenses editing flag does not gate them.

## Purchase order approve, reject, cancel, close, convert, revise and group actions need expenses editing

Lines 246-262

Asked as: can't approve purchase order; reject PO blocked; close purchase order disabled; group PO preview

Turbo applies the expenses editing check to these purchase order actions: approve, reject, cancel, close, make cumulative, revise, and group preview and create. When the flag is off, it refuses them. Listing and viewing purchase orders is not gated by this check.

## Client routes

Lines 270-279

Asked as: client rate sheet; change invoicing profile; client notes and divisions; client contact details

Signed-in users reach client routes for list, detail, notes, time divisions, rate sheets, contact and invoicing profile usage, and replacing an invoicing profile. Procurement rules sit in the handlers.

## Absorb and undo-absorb routes exist for several record types

Lines 280-294

Asked as: undo absorb client; absorb vendor; rate sheet absorb admin; undo vendor merge

Absorb and undo absorb are available for clients, client contacts, vendors and rate sheets. Rate sheet absorb is for admins, or for absorb holders with rate_sheet_revise.

## Job routes

Lines 298-333

Asked as: job wip report; adopt subjob; close job; job proposal validation

Signed-in users reach job routes for create, upsert, details, time, staff, expense and PO summaries, WIP, status, number, client change, adopt subjobs, rate sheet periods, close and proposal validation. Handlers enforce who may act. Time summary filters are given by query parameters (division, time_type, user, category).

## WIP reports come in three modes

Lines 335-339

Asked as: wip report my branch; division wip; work in progress report; wip by mode

The WIP report is available in modes my, branch and division for signed-in users.

## Reports API requires the report claim

Lines 341-360

Asked as: payroll report access; weekly time spreadsheet denied; active jobs report not allowed; receipts spreadsheet blocked

Every route under /api/reports requires a signed-in user, and then the report claim. A user without it gets the message "you are not authorized to view this report". Payroll, weekly time, expense, receipts and payables spreadsheet reports, the time entry branch mismatches report and the active jobs report all sit behind this check.

## KPI report route requires the kpi claim

Lines 362-370

Asked as: employee branch hours kpi; kpi hours report denied; kpi branch hours access; employee hours kpi report

The employee branch hours KPI report under /api/kpi/reports needs a signed-in user with the kpi claim. A user without it gets the message "you are not authorized to view KPI reports".

## Branch hours report checks access by branch

Lines 372-380

Asked as: branch hours report access; can't see branch hours; branch query parameter error; branch hours forbidden

The branch hours report runs requireBranchHoursAccess with the branch query parameter before it returns data. This file does not show who passes that check.

## Machine secrets management is admin only

Lines 387-391

Asked as: create machine secret; list api secrets; machine secret admin; manage machine tokens

The routes to list and create machine secrets need sign-in. The comment says the handlers limit them to admins.

## Attachment audit, currency and claim routes

Lines 393-442

Asked as: orphaned attachments csv; reload exchange rates; bulk assign claims; delete orphan attachments

Signed-in users reach routes for attachment audit (targets, refresh, delete orphans, missing and orphaned CSVs), currencies (list, reload rates, backfill, delete), claims (list, assignable users, bulk assign, details) and work records. Handlers decide who may use each. Claim details need the admin claim.

## Rate sheet routes need job claim; some are admin only

Lines 410-447

Asked as: create rate sheet; rollout rate sheet revision; delete inactive rate sheet; update rate sheet entries

Rate sheet routes (create, compare, activate, deactivate, rename, former names, jobs) need the job claim according to the comment. Deleting an inactive unused sheet and rollout to the latest revision are admin actions. Updating rate sheet entries needs the admin claim.
