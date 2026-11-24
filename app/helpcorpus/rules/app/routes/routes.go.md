<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"bd260c361e199546043662755d82c68c53c2ac7cae5c49b390b3c3477b0837ed","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses editing disabled; can't edit expenses; purchase orders locked; vendor changes refused; expenses flag off

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message from ErrExpensesEditingDisabled. Other collections skip this check. If the flag cannot be read, the request fails with that error.

## Time editing needs the time feature flag

Lines 33-43

Asked as: time editing disabled; can't change timesheet; time entries locked; time feature flag; cannot edit time

When the time editing feature flag is off, Turbo refuses the action with the message from ErrTimeEditingDisabled. If the flag cannot be read, the request fails with that error.

## Reports need the report claim

Lines 45-54

Asked as: report access denied; not authorized to view report; can't open payroll report; missing report permission; reports blocked

A person without the report claim cannot use the report routes. Turbo refuses them with "you are not authorized to view this report". These routes are under /api/reports, such as payroll, weekly time, expense, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report denied; can't see KPI reports; kpi permission missing; employee branch hours access; KPI not authorized

A person without the kpi claim cannot use the KPI report routes. Turbo refuses them with "you are not authorized to view KPI reports". The employee branch hours report is under this group.

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds invalid; startup error expenditure; config fails on boot; expenditure config broken

When the server starts, Turbo validates the expenditure kinds configuration. If it is invalid, the startup fails with that error.

## Version endpoint is public

Lines 93-105

Asked as: version without login; app version number; build number lookup; git commit shown; version check no sign-in

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit, short commit, git branch and build time.

## Most custom routes require a signed-in user

Lines 107-437

Asked as: need to log in; api says unauthorized; signed in required; route needs login; session missing for endpoint

The time sheet, expense, expense uncommit request, absorb proposal, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, report, stats, machine secret, attachment audit, currency, rate sheet, claim and work record routes all need a signed-in user. The legacy export routes and radiator routes use machine secrets instead.

## Legacy export routes use machine secrets

Lines 133-136

Asked as: export without login; legacy export failing; machine secret export; timesheet export api key; jobs export secret

The legacy export routes for time sheets, jobs and expenses do not require a user sign-in. They authenticate with machine secrets.

## Absorb routes exist for several record types

Lines 185-192

Asked as: absorb duplicates; merge duplicate clients; propose merge vendor; undo merge; combine rate sheets

Absorb and undo absorb routes exist for clients, client contacts, vendors and rate sheets. Proposing an absorb is open to any signed-in person, and absorb holders review proposals.

## Purchase order changes need expenses editing enabled

Lines 237-253

Asked as: can't approve purchase order; PO approve blocked; cancel PO refused; close purchase order fails; revise PO disabled

Approve, reject, cancel, close, make cumulative, revise, and group preview and create for purchase orders are refused when the expenses editing flag is off. Read routes and the legacy purchase order routes do not run this check.

## Branch hours report checks access for the requested branch

Lines 360-368

Asked as: branch hours report blocked; can't run branch hours; branch access denied; hours report for branch; branch report permission

Before the branch hours report runs, Turbo calls requireBranchHoursAccess with the signed-in person and the branch query parameter. The result of that check is defined in another file.

## Rate sheet routes and who can use them

Lines 398-417

Asked as: can't create rate sheet; activate rate sheet denied; rename rate sheet; delete old rate sheet; roll projects to latest revision

The rate sheet routes for creating, activating, deactivating and renaming need the job claim according to the comment. Deleting an inactive unused sheet and rolling projects out to the latest revision are admin actions. Absorbing duplicate rate sheets is for admin, or for absorb with rate_sheet_revise. The checks are in other files.
