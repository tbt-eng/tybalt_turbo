<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"2cbf72d63747b52bbce3f6e797091f5c11c652a979f18d49fdf7a9e56c338c03","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expense editing is off

Lines 14-31

Asked as: expenses locked; can't edit expenses; expense editing disabled; purchase orders forbidden; vendors read only

When the expenses editing feature flag is off, Turbo refuses requests on the expenses, purchase_orders and vendors collections with a forbidden error. The error message is the one for ErrExpensesEditingDisabled. Other collections are not checked.

## Time editing blocked when the time feature flag is off

Lines 33-43

Asked as: can't edit my time; time editing off; timesheet locked; time entries forbidden; time feature disabled

When the time editing feature flag is off, Turbo refuses the request with a forbidden error. The error message is the one for ErrTimeEditingDisabled (utilities.ErrTimeEditingDisabled).

## Reports need the report claim

Lines 45-54

Asked as: can't see reports; report access denied; no permission for reports; reports not authorized; why can't I open report

Only staff with the report claim can view a report. Anyone else gets the message "you are not authorized to view this report".

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report denied; can't see KPI; kpi access missing; KPI not authorized; dashboard metrics forbidden

Only staff with the kpi claim can view KPI reports. Anyone else gets the message "you are not authorized to view KPI reports".

## Turbo checks expenditure kinds configuration at startup

Lines 89-91

Asked as: server won't start; startup fails expenditure; expenditure kinds invalid; app crashes on boot; config error at startup

When the server starts, Turbo validates the expenditure kinds configuration (ValidateExpenditureKindsConfig). If it is invalid, startup fails with that error.

## Version endpoint is public

Lines 93-105

Asked as: app version number; what build is this; version without login; git commit info; check build version

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit (full and short), git branch and build time.

## Most Turbo API routes need a signed-in user

Lines 107-417

Asked as: API needs login; unauthorized API call; not signed in error; api returns 401; need to log in first

The time sheet, expense, expense uncommit request, time amendment, time entry, purchase order, client, vendor, job, report, rate sheet, claim and other custom API routes require a signed-in user (the users collection). Only the version endpoint and the legacy export routes skip this check.

## Legacy export routes use machine secrets, not user sign-in

Lines 133-136

Asked as: export without login; legacy export secret; machine secret export; export token access; scripted export access

The legacy export routes for time sheets, jobs and expenses are outside the user sign-in group. They use their own authorization through machine_secrets.

## Purchase order changes need expense editing enabled

Lines 221-237

Asked as: can't approve PO; purchase order frozen; PO approval blocked; can't close purchase order; cancel PO not working

Approving, rejecting, cancelling, closing, making cumulative, revising and creating groups of purchase orders go through a check of the expenses editing flag. When the flag is off, Turbo refuses them. Viewing routes for purchase orders do not have this check.

## Absorb and undo absorb are available for several record types

Lines 245-266

Asked as: absorb a client; merge vendor records; absorb rate sheet; undo absorb; combine duplicate records

Absorb and undo absorb routes exist for clients, client contacts, vendors and rate sheets. The rate sheet comments say an admin can absorb rate sheets, or a person with rate_sheet_revise.

## Report routes require the report claim

Lines 313-332

Asked as: payroll report denied; weekly report blocked; receipts report access; payables spreadsheet no access; active jobs report

Every route under /api/reports requires the report claim, including payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI report routes require the kpi claim

Lines 334-342

Asked as: employee hours report kpi; branch hours by employee; kpi employee report blocked; employee hours denied

The employee branch hours report under /api/kpi/reports requires the kpi claim.

## Branch hours report checks access for the requested branch

Lines 344-352

Asked as: branch hours not showing; branch report denied; hours for my branch; branch access error; can't see branch hours

Before the branch hours report runs, Turbo checks the caller's access for the branch given in the branch query parameter (requireBranchHoursAccess). The rules of that check are not in this file.

## Machine secrets and claims routes have admin limits

Lines 359-363

Asked as: machine secrets admin; who can edit claims; claim details restricted; rate sheet entry admin only; admin only settings

The machine secrets routes are labelled admin only. The claim detail route requires the admin claim. The rate sheet entry update route requires the admin claim. The checks are in the handlers, not in this file.
