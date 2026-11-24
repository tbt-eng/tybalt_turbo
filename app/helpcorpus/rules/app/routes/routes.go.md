<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"aba3172dde10780d6be40a0189da7a0f99fe425106c4876f92c2c95335a3ca62","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses editing disabled; can't edit expenses; purchase orders locked; vendor changes refused

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message in ErrExpensesEditingDisabled. Other collections skip this check. The wording of the message is defined in another file.

## Time editing needs the time feature flag

Lines 33-43

Asked as: cannot edit time; time editing off; timesheet edits blocked; time feature disabled

When the time editing feature flag is off, Turbo refuses with the message in ErrTimeEditingDisabled. The wording of the message is defined in another file.

## Reports need the report claim

Lines 45-54

Asked as: can't view report; report access denied; not authorized report; no reports permission

A person without the report claim cannot view a report. Turbo refuses with 'you are not authorized to view this report'.

## KPI reports need the KPI claim

Lines 56-65

Asked as: KPI report denied; can't see KPI; KPI access missing; kpi permission needed

A person without the kpi claim cannot view KPI reports. Turbo refuses with 'you are not authorized to view KPI reports'.

## Startup fails if expenditure kinds config is invalid

Lines 88-91

Asked as: server won't start; startup error expenditure; expenditure kinds invalid; app fails to boot

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the server does not start the routes.

## Version endpoint is public

Lines 93-105

Asked as: app version check; version without login; build number endpoint; which version running

GET /api/version needs no sign-in. It returns the app name, version, build, full version, git commit, short commit, git branch and build time.

## Most Turbo API routes need a signed-in user

Lines 107-437

Asked as: API needs login; unauthorized api call; sign in required api; 401 on timesheets

Time sheet, expense, expense uncommit request, absorb proposal, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, stats, machine secret, attachment audit, currency, rate sheet, claim and work record routes all require sign-in as a user. The legacy export routes use machine secrets instead. The radiator routes also use machine secrets (role radiator).

## Time sheet commit also covers time amendments

Lines 115-117

Asked as: who can commit; commit amendment; committing timesheets permission; commit time amendments

Committers commit time sheets through the time sheet commit route. Time amendments have their own commit route. Who may commit is decided in the handler, in another file.

## Legacy export routes use machine secrets, not sign-in

Lines 133-136

Asked as: export without login; legacy export access; machine secret export; export failing auth

The legacy export routes for time sheets, jobs and expenses do not need a signed-in user. They authenticate with machine secrets, and the check is in another file.

## Purchase order changes are blocked when expenses editing is off

Lines 238-254

Asked as: can't approve purchase order; PO approve blocked; cancel PO not working; purchase order frozen

Approve, reject, cancel, close, make cumulative and revise on a purchase order, and previewing or creating a purchase order group, need a signed-in user and the expenses editing flag on. Viewing purchase orders is not behind this flag. Legacy purchase order routes do not have this check in this file.

## Report routes require the report claim

Lines 330-349

Asked as: spreadsheet report denied; payroll report access; weekly time report blocked; active jobs report access

All routes under /api/reports require sign-in and the report claim. They cover the payroll, weekly time, expense, receipts and payables spreadsheet reports, the time entry branch mismatches report and the active jobs report.

## KPI report routes require the KPI claim

Lines 351-359

Asked as: employee hours report access; branch KPI hours denied; hours KPI not visible; KPI hours report permission

The employee branch hours KPI report requires sign-in and the kpi claim.

## Branch hours report checks branch access

Lines 361-369

Asked as: branch hours access denied; can't see branch hours; branch report permission; hours report wrong branch

The branch hours report requires sign-in. Turbo also runs requireBranchHoursAccess on the branch given in the request's 'branch' query parameter. That function is in another file, so its rule is not stated here.
