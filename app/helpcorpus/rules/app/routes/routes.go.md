<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"00ae1e300eb76f6ec7b35a219f38f4690a4147d2b40fa031014fb8ceb6aca106","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors need the expenses editing flag

Lines 14-31

Asked as: expenses editing disabled; can't edit expenses; purchase orders locked; vendors read only

When the expenses editing feature flag is off, Turbo refuses actions on the expenses, purchase_orders and vendors collections. The refusal message is the one in utilities.ErrExpensesEditingDisabled. Other collections skip this check.

## Time editing needs the time editing flag

Lines 33-43

Asked as: time editing disabled; can't change timesheet; time entries locked; timesheet edits refused

When the time editing feature flag is off, a time editing check refuses with the message in utilities.ErrTimeEditingDisabled. The helper requireTimeEditing runs this check. This file does not show which handlers call it.

## Reports need the report claim

Lines 45-54

Asked as: not authorized to view report; no access to reports; report permission denied; can't open payroll report

A person without the report claim cannot use any report under /api/reports. Turbo refuses with "you are not authorized to view this report". This covers the payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports need the kpi claim

Lines 56-65

Asked as: KPI report access denied; can't see KPI reports; no kpi permission; kpi reports forbidden

A person without the kpi claim cannot use the KPI reports under /api/kpi/reports, such as the employee branch hours report. Turbo refuses with "you are not authorized to view KPI reports".

## Turbo checks expenditure kinds config at startup

Lines 88-91

Asked as: server won't start; expenditure kinds config error; startup validation failed; app crashes on launch

When the server starts, Turbo validates the expenditure kinds configuration. If the validation fails, the server does not start.

## Version endpoint is public

Lines 93-105

Asked as: version without login; app build number; which version is running; git commit of build

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Signed-in users only for most custom routes

Lines 107-427

Asked as: not signed in; request refused logged out; login required for api; need to log in

The time sheet, expense, expense uncommit request, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, report, stats, machine secret, attachment audit, currency, rate sheet, claim and work record routes require a signed-in user. Without one, Turbo refuses the request.

## Legacy export and radiator routes use machine secrets

Lines 133-139

Asked as: export without login; radiator report access; machine secret export; legacy export link

The legacy export routes for time sheets, jobs and expenses do not need a user sign-in. They use their own authentication through machine_secrets. The radiator login report routes do the same, with the role radiator.

## Purchase order changes need the expenses editing flag

Lines 228-244

Asked as: can't approve purchase order; PO approval blocked; cannot close PO; PO group creation refused

Approving, rejecting, cancelling, closing, converting to cumulative and revising a purchase order, and previewing or creating PO groups, all need the expenses editing flag on. When it is off, Turbo refuses them. Reading POs and the legacy PO routes do not use this check in this file.

## Branch hours report checks access for the requested branch

Lines 352-360

Asked as: branch hours report denied; can't see branch hours; branch hours access; hours for other branch

The /api/branch_hours report runs an access check (requireBranchHoursAccess) with the branch query parameter before it returns data. This file does not show who passes the check.

## Machine secrets are managed by admins

Lines 367-371

Asked as: who manages machine secrets; add machine secret; machine secret permission; secret admin only

Only admins manage machine secrets, by the route comment. The list and create routes need a signed-in user. This file does not show the admin check.
