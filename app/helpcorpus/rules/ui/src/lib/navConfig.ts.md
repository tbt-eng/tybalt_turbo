<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"cc80d972ffe8f83d92e428a6713d3bda3858df1d0970525422759a7a0f019c32","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Quick-add buttons next to menu items

Lines 35-242

Some menu items have a green "New" button: New Entry, New Amendment, New PO, New PO Group, New Expense, New Job, New Project Request, New Client and New Vendor. The Amendments item also has a "List All Amendments" button. The Jobs item also has blue buttons for Unused Active Jobs, Stale Jobs, Latest Jobs and Jobs Map.

## Menu items limited by a view check function

Lines 89-92

Some menu items show only when a check function allows it: PO Groups (canViewPOGroups), the WIP "Branch" and "Division" reports (canViewWIPReports), PO Approvers (canViewPOApprovers) and Branch Hours (canViewBranchHours). Their rules are in other files. "My WIP" has no check.

## Menu items open queues without a claim

Lines 131-131

Queue items like Request Approvals and Job Create Queue, and the Commit Queue, have no claim or check in this menu file. The menu shows them to all staff.

## Uncommit Requests menu needs commit or admin claim

Lines 132-137

The Expenses menu shows "Uncommit Requests" only to a committer (the commit claim) or an admin (the admin claim). Other staff do not see it.

## WIP menu groups three reports

Lines 179-188

The "WIP" menu item has no page of its own. It holds three sub-items: My WIP, Branch and Division.

## Value Approvals menu needs project authorization claim

Lines 203-208

The Business menu shows "Value Approvals" only to staff with the project_authorization claim (requiredClaim).

## Abandoned Requests menu is admin only

Lines 209-214

The Business menu shows "Abandoned Requests" (project requests in the abandoned queue) only to an admin (the admin claim).

## UI Components menu is admin only

Lines 253-253

The Business menu shows "UI Components" only to an admin (the admin claim).

## Report menu items need the report claim

Lines 259-293

The Reports menu shows Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs only to staff with the report claim.

## Branch Mismatches CSV and Active Jobs download files

Lines 268-293

"Branch Mismatches CSV" and "Active Jobs" have no page. Each has a "Download CSV" button that downloads a CSV file (downloadTimeEntryBranchMismatchesCsv, downloadActiveJobsCsv).

## KPI report needs the kpi claim

Lines 300-300

The Reports menu shows "KPI" only to staff with the kpi claim.

## Settings menu items have no claim or view check

Lines 303-312

The Settings menu section lists Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit. None of these items has a required claim (requiredClaim) or a view check (canView) in the menu configuration. Any restrictions on these pages come from somewhere other than this menu file.
