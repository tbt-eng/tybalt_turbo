<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"504ab4104875e7ae0c2428849566cd92adc4a03ed698f6eec70ec498a2704bad","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, POs and vendors refuse edits when expenses editing is off

Lines 14-31

Asked as: can't edit expense; expenses editing disabled; edits blocked on expenses; vendor changes refused; PO edit not allowed

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message held in ErrExpensesEditingDisabled. Other collections skip this check. The message text is defined in another file.

## Time editing flag off blocks time actions

Lines 33-43

Asked as: can't edit timesheet; time edits blocked; time entry locked; timesheet editing turned off; time changes refused

When the time editing feature flag is off, a check that uses this helper refuses the action with the message held in ErrTimeEditingDisabled. The message text is defined in another file. This file does not show which routes call the helper.

## Reports need the report claim

Lines 45-54

Asked as: can't open report; report access denied; not authorized to view report; need permission for payroll report; report shows forbidden

Everyone who opens a report under /api/reports needs the report claim. Without it Turbo refuses with 'you are not authorized to view this report'. This covers payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI reports; KPI not authorized; KPI page forbidden; no access to kpi

Everyone who opens a KPI report under /api/kpi/reports needs the kpi claim. Without it Turbo refuses with 'you are not authorized to view KPI reports'.

## Expense config is validated at startup

Lines 88-91

Asked as: Turbo won't start; startup config error; expenditure kinds invalid; app fails to boot; server exits on start

When Turbo starts serving, it validates the expenditure kinds configuration. If it is invalid, Turbo stops serving and reports the error.

## Version endpoint is public

Lines 93-105

Asked as: app version check; which build is running; version without login; git commit version; check turbo version

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Most Turbo API routes require a signed-in user

Lines 107-132

Asked as: API returns unauthorized; need to log in to api; 401 on api call; not signed in error; api says forbidden

The time sheet, expense, PO, client, vendor, job, report, rate sheet, claim and other API groups require a signed-in user account. The exceptions are the legacy export routes, the radiator routes and the version endpoint, which use machine secrets or no auth.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: legacy export without login; export api key; machine secret export; scheduled export failing; export needs no user

The legacy exports of time sheets, jobs and expenses do not require a signed-in user. They authenticate with machine secrets handled in another file.

## Radiator and computer routes use their own auth

Lines 138-143

Asked as: radiator login report access; computer list forbidden; assign computer permission; retire computer access; radiator report denied

The radiator login report routes use machine secrets with the radiator role. The computer routes (list, details, assign, retire, un-retire) are for holders of the computers or admin claim. The checks live in other files.

## Expense routes require sign-in

Lines 145-175

Asked as: expense api unauthorized; can't submit expense; approve expense not signed in; expense endpoint 401; recall expense error

All expense routes require a signed-in user. They cover create, edit, submit, recall, approve, reject, commit, uncommit, uncommit requests, attachment checks, settle and reimbursements. Claim and state checks happen in the handlers.

## Expense uncommit requests are a separate route group

Lines 156-183

Asked as: uncommit request not found; cancel uncommit request; request to uncommit expense; decline uncommit request; uncommit requests list

Uncommit requests for expenses are created through /api/expenses/{id}/uncommit_requests. Listing, approving, declining and cancelling them use /api/expense_uncommit_requests. All need sign-in.

## Absorb proposals: anyone signed in can propose, absorb holders review

Lines 185-193

Asked as: propose absorb; who can approve absorb; absorb proposal review; withdraw absorb request; absorb request declined

Absorb proposal routes need sign-in. Per the comment, anyone proposes and absorb holders review. Create, list, view, delete, withdraw and decline are separate routes, and the handlers decide who may do each.

## Time amendments commit route

Lines 195-197

Asked as: commit time amendment; amendment won't commit; approve time change; time change commit failing; amendment stuck

Committing a time amendment goes through POST /api/time_amendments/{id}/commit and needs sign-in. The commit handler checks the rest.

## Admin profile routes require sign-in

Lines 213-223

Asked as: edit user profile; deactivate user; clear authorized provider; admin user settings; save user claims

Admin profile routes need sign-in. They cover identity, active toggle, save with claims, limited save and clearing an authorized provider. Handlers decide who may do each.

## Purchase order mutations blocked when expenses editing is off

Lines 239-261

Asked as: can't approve PO; PO approval blocked; can't close purchase order; PO groups refused; purchase order action fails

Approving, rejecting, cancelling, closing, making cumulative and revising a purchase order, and previewing or creating PO groups, all require the expenses editing flag to be on. When it is off, Turbo refuses these actions. Viewing POs, and the legacy PO create and edit routes, do not pass through this check.

## Purchase order approve, reject, cancel, close, revise and group changes need expenses editing

Lines 239-255

Asked as: reject purchase order fails; cancel PO not working; PO revise error; expenses editing disabled PO; convert PO to cumulative

Approving, rejecting, cancelling, closing, converting to cumulative and revising a PO, and previewing or creating PO groups, are refused when the expenses editing flag is off. Turbo shows the expenses-editing-disabled message. Reading PO lists, details and approvers is not behind this check.

## Legacy PO routes are not behind the expenses editing flag

Lines 257-261

Asked as: old PO edit screen; legacy purchase order form; create PO legacy; PO edit works when flag off; old purchase order page

The legacy PO routes (view for edit, create, update) require sign-in but have no expenses editing check in this file.

## Absorb and undo-absorb routes for clients, contacts, vendors and rate sheets

Lines 263-284

Asked as: absorb client; undo absorb vendor; merge rate sheet; undo merge contact; absorb contact into another

Clients, client contacts, vendors and rate sheets each have an absorb route and an undo-absorb route. The comment says rate sheet absorb is for admins, or for absorb holders with rate_sheet_revise. Handlers enforce the claims.

## Job routes require sign-in

Lines 288-323

Asked as: job api unauthorized; can't create job; close job error; job rate sheet API; subjob adoption failing

All /api/jobs routes need sign-in. They cover create, upsert, details, summaries, status, number, client change, adopt subjobs, rate sheets, close and proposal validation. Claim checks are in the handlers.

## WIP report modes

Lines 325-329

Asked as: work in progress report; WIP my branch division; WIP report not loading; wip division view; WIP branch view

The WIP report has three modes: my, branch and division. Each is its own route under /api/wip and requires sign-in.

## Reports API requires the report claim

Lines 331-347

Asked as: report not authorized message; can't view any report; reports locked for me; no reports access; all reports forbidden

Every route under /api/reports refuses a signed-in user who lacks the report claim. Turbo shows the message "you are not authorized to view this report".

## Payroll, weekly, payables and branch-mismatch reports need the report claim

Lines 339-350

Asked as: payroll report access; branch mismatch report; payables spreadsheet permission; receipts report denied; active jobs report access

The payroll time, weekly time, expense, receipts and payables spreadsheet reports, the time entry branch mismatches report and the active jobs report all sit behind the report claim check.

## Employee branch hours KPI report requires the kpi claim

Lines 352-360

Asked as: employee branch hours report; branch hours KPI denied; KPI employee hours access; employee hours kpi forbidden; branch hours kpi

The employee branch hours KPI report (/api/kpi/reports/employee_branch_hours) requires the kpi claim. A user without it sees "you are not authorized to view KPI reports".

## Branch hours report is checked per requested branch

Lines 362-370

Asked as: branch hours for one branch; can't see branch hours; branch hours access denied; branch hours report forbidden; hours by branch access

Access to the branch hours report is checked by requireBranchHoursAccess with the branch given in the request. The rule for who passes is in another file.

## Machine secrets routes exist for admins

Lines 377-381

Asked as: create machine secret; list API secrets; machine secret admin only; secrets page access; make new api secret

Listing and creating machine secrets are routes for signed-in users. The comment says admin only, but the check is in the handlers in another file.

## Currency and attachment audit routes require sign-in

Lines 383-398

Asked as: reload exchange rates; currency backfill; orphaned attachments; missing attachment CSV; delete orphan attachments

Currency routes (list, init status, reload rates, initialize backfill, delete) and attachment audit routes (targets, refresh, delete orphans, missing and orphaned CSVs) require sign-in. Further checks are in handlers.

## Rate sheet routes and their claims

Lines 400-419

Asked as: rate sheet delete inactive; roll projects to latest revision; rate sheet management job claim; rate sheet permission; rate sheet access denied

Rate sheet routes require sign-in. The comments say the job claim is required for management. Admin is required to delete an inactive unused sheet and to roll projects to the latest revision. Handlers enforce these.

## Claims detail needs the admin claim

Lines 421-427

Asked as: claim details forbidden; bulk assign claims; list assignable users; view claims admin; claims page access

Listing claims, assignable users, bulk assign and claim details need sign-in. The comment says the detail endpoint requires the admin claim, which the handler checks.

## Rate sheet entry update needs the admin claim

Lines 434-437

Asked as: edit rate sheet entry; update rate not saving; rate entry PUT forbidden; change rate sheet price; rate entry update denied

Updating a rate sheet entry (PUT /api/rate_sheet_entries/{id}) requires sign-in. The comment says the admin claim is also required, checked in the handler.
