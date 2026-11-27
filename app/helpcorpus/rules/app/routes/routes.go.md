<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"19e218d143dc36daf0e3fdeec364390c2b329a1dd2624b9c585ae7bf381ef7d4","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors refuse edits when the expenses editing flag is off

Lines 14-31

Asked as: expenses editing disabled; can't edit expense; purchase order edit blocked; vendor changes refused

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors that go through this check. The message is the one in ErrExpensesEditingDisabled. Other collections pass the check.

## Time sheet editing refused when the time editing flag is off

Lines 33-43

Asked as: time sheet edit blocked; can't change timesheet; time editing turned off; timesheet locked

When the time editing feature flag is off, a check that runs before time actions refuses them. Turbo gives the message in ErrTimeEditingDisabled.

## Reports need the report claim

Lines 45-54

Asked as: can't open reports; report access denied; not allowed to view reports; reports forbidden

A person without the report claim cannot open reports under /api/reports. Turbo answers 'you are not authorized to view this report'. The claim is checked for every report route in that group.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI reports; KPI dashboard forbidden; kpi permission missing

A person without the kpi claim cannot open KPI reports under /api/kpi/reports. Turbo answers 'you are not authorized to view KPI reports'.

## Turbo checks expenditure kinds configuration at startup

Lines 88-91

Asked as: server won't start; startup config error; expenditure kinds invalid; app fails to boot

When the server starts, Turbo validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If it is invalid, the server does not finish starting up.

## Version endpoint is public

Lines 93-105

Asked as: app version number; which build is running; git commit hash; version info without login

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Time sheet routes need sign-in

Lines 107-131

Asked as: timesheet api login; time sheet sign in; timesheet route unauthorized; who can commit timesheet

All time sheet routes (bundle, unbundle, copy to next week, details, approve, reject, commit, uncommit, reviewers, approver, tallies and tracking) need a signed-in user. Tracking routes are for report holders, and the handlers apply the claim checks. Committing a time sheet is a committer action.

## Legacy export endpoints use machine secrets, not user sign-in

Lines 133-136

Asked as: export without login; legacy export api; export secret key; machine export access

The legacy export routes for time sheets, jobs and expenses do not need a user sign-in. They use their own authorization with machine secrets.

## Radiator reports use machine secrets

Lines 138-139

Asked as: radiator login report; radiator api key; radiator report access; machine report login

Radiator login report routes use machine secrets with the radiator role, not user sign-in.

## Computer routes need the computers or admin claim

Lines 141-143

Asked as: computer list forbidden; assign laptop access; retire computer permission; no computers claim

Computer routes (list, details, assign, retire, un-retire) are for people with the computers or admin claim.

## Expense routes need sign-in

Lines 145-176

Asked as: expense api login; who can approve expense; submit expense forbidden; expense copy to tomorrow

Expense actions need a signed-in user. They cover create, update, submit, recall, copy to tomorrow, approve, reject, commit, uncommit, uncommit requests, lists, attachments, settlement and reimbursements. Each handler applies its own rules. Recall is available for expenses, and commit is for committers.

## Expense uncommit requests routes

Lines 178-184

Asked as: uncommit request approve; decline uncommit; cancel uncommit request; uncommit reasons list

Signed-in users can list uncommit requests and reasons. They can approve, decline or cancel requests through these routes.

## Absorb proposal routes

Lines 186-194

Asked as: propose absorb; absorb proposal review; withdraw absorb request; who can merge clients

Any signed-in user can propose an absorb. Absorb holders review proposals. Routes cover list, create, view, delete, withdraw and decline.

## Absorb suggestions and history routes

Lines 196-208

Asked as: duplicate client suggestions; dismiss probable duplicate; absorb history list; past merges log

Absorb suggestions of probable duplicate clients and vendors can be listed or dismissed per collection, for absorb holders. Absorb history lists committed absorbs, kept for ABSORB_HISTORY_RETENTION_DAYS.

## Time amendments commit route

Lines 210-212

Asked as: commit time amendment; time amendment api; amendment commit failed; approve time change

Time amendments are committed through POST /api/time_amendments/{id}/commit, for signed-in users, with the commit handler applying its rules.

## Admin profile routes need sign-in

Lines 228-238

Asked as: edit user profile; toggle user active; save user claims; clear auth providers

Admin profile routes need sign-in. They cover identity, active toggle, saving with claims, limited saving and clearing authorized providers. Handlers apply their own checks.

## Purchase order mutations blocked when expenses editing is off

Lines 240-277

Asked as: purchase order approve blocked; PO close refused; cancel purchase order fails; group purchase orders disabled

Approve, reject, cancel, close, make cumulative, revise and group create or preview for purchase orders first check the expenses editing flag. If the flag is off, Turbo refuses them. Read routes such as pending, visible, search, approvers and groups do not check it. The legacy PO routes (edit, create, update) also do not check it.

## Client invoice routes

Lines 279-284

Asked as: issue client invoice; invoice lines api; billable items read; create invoice api

Signed-in users can issue an invoice and read its lines and billable items. Drafts and lines are written through their collections.

## Client, contact and vendor absorb routes

Lines 286-310

Asked as: merge clients; undo client absorb; vendor merge; replace invoicing profile

Absorb and undo absorb routes exist for clients, client contacts and vendors. Client routes also cover notes, time divisions, rate sheets, usage and replacing an invoicing profile.

## Job routes

Lines 314-349

Asked as: job api access; close job forbidden; adopt subjob; job proposal validation

Job routes need sign-in. They cover create, upsert, details, notes, summaries, reports, WIP, status, number, client change, adopting subjobs, rate sheet periods, close and proposal validation. Handlers apply their own claim checks.

## WIP reports by mode

Lines 351-355

Asked as: WIP report modes; my WIP report; division wip view; branch wip view

WIP reports are available in three modes: my, branch and division.

## Report routes need the report claim

Lines 357-376

Asked as: payroll report access; receipts report forbidden; payables spreadsheet download; active jobs report

Payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports need the report claim.

## Branch hours report access check

Lines 388-396

Asked as: branch hours report denied; branch hours access; hours by branch forbidden; branch query check

The branch hours report checks access using the branch query parameter, through requireBranchHoursAccess, before it runs.

## Admin routes: stats, machine secrets, attachment audit, currencies

Lines 398-424

Asked as: machine secrets list; attachment audit access; currency admin; stats page forbidden

These groups need sign-in and cover stats, machine secrets (listed as admin only), attachment audit and currencies. Handlers apply their own checks.

## Rate sheet routes and who may use them

Lines 426-445

Asked as: rate sheet delete; rollout rate sheet; absorb duplicate rate sheets; rate sheet permission

Rate sheet routes need sign-in. The comments say they need the job claim. Deleting an inactive unused sheet and rollout are admin actions. Absorbing duplicate rate sheets is for admin or absorb with rate_sheet_revise.

## Claims, work records and rate sheet entry routes

Lines 447-463

Asked as: claim details forbidden; edit rate sheet entry; work record access; rate entry admin only

Claim details need the admin claim, according to the comment. Editing a rate sheet entry needs the admin claim.
