<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"d79213bbc3b1e117d44fdae2501c4a34bd7e58917fc1960dd48a6a5a598c66a5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu links without a restriction in this file

Lines 32-314

Asked as: menu link no restriction; which menu links need permission; why can I see this menu item; menu items without access check

Other menu links, such as Entries, Sheets, My Purchase Orders, My Expenses, Commit Queue, Jobs, Clients, Vendors, Admin Profiles, Claims and the Settings items, have no claim check in this file. Pages may still apply their own checks.

## Quick-action buttons beside menu items

Lines 36-243

Asked as: quick action buttons; buttons beside menu items; New Entry button missing; shortcut buttons under menu

Some menu items carry buttons: New Entry, New Amendment, New PO, New Expense, New Job, New Project Request, New Client and New Vendor. The Jobs item also has buttons for Unused Active Jobs, Stale Jobs, Latest Jobs and Jobs Map. The Amendments item has a "List All Amendments" button, and its link opens pending amendments.

## PO Groups menu link depends on canViewPOGroups

Lines 90-102

Asked as: PO Groups missing; can't see PO Groups; New PO Group button; purchase order groups hidden

The Purchase Orders menu shows "PO Groups" and its "New PO Group" button only when canViewPOGroups (in $lib/poGroups) allows it. The rule itself is not in this file.

## Uncommit Requests menu link needs commit or admin claim

Lines 133-138

Asked as: Uncommit Requests missing; can't see uncommit link; uncommit menu not showing; who can uncommit

The Expenses menu shows "Uncommit Requests" only to a committer (the commit claim) or an admin (the admin claim). Other staff do not see this link.

## WIP Branch and Division links depend on canViewWIPReports

Lines 180-189

Asked as: WIP Branch missing; WIP Division not showing; can't see WIP reports; work in progress report hidden

Under the WIP menu, "My WIP" has no restriction. "Branch" and "Division" show only when canViewWIPReports (in $lib/reports/wipReports) allows it.

## Value Approvals menu link needs project authorization claim

Lines 204-209

Asked as: Value Approvals missing; can't see value approvals; project authorization access; value approval link hidden

The Business menu shows "Value Approvals" (/jobs/value-approvals) only to staff with the project_authorization claim (requiredClaim).

## Abandoned Requests menu link is admin only

Lines 210-215

Asked as: Abandoned Requests missing; can't see abandoned project requests; abandoned queue hidden; admin only abandoned requests

The Business menu shows "Abandoned Requests" (project requests, queue=abandoned) only to an admin (the admin claim).

## Invoices menu link depends on canViewInvoices

Lines 245-245

Asked as: Invoices missing from menu; can't see invoices link; invoice menu hidden; no access to invoices

The Business menu shows "Invoices" only when canViewInvoices (in $lib/invoices) allows it.

## PO Approvers menu link depends on canViewPOApprovers

Lines 247-252

Asked as: PO Approvers missing; can't see PO approval limits; approval limits link hidden; who approves POs menu

The Business menu shows "PO Approvers" (/pos/approval-limits) only when canViewPOApprovers (in $lib/poApprovalLimits) allows it.

## UI Components menu link is admin only

Lines 255-255

Asked as: UI Components missing; can't see dev components; components menu hidden; developer components page

The Business menu shows "UI Components" (/dev/components) only to an admin (the admin claim).

## Report menu items need the report claim

Lines 261-295

Asked as: Payroll report missing; can't see reports menu items; CSV download not showing; branch mismatches download hidden

In the Reports menu, the items Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs require the report claim. The last two download CSV files (downloadTimeEntryBranchMismatchesCsv, downloadActiveJobsCsv).

## Branch Hours report link depends on canViewBranchHours

Lines 296-301

Asked as: Branch Hours missing; can't see branch hours report; branch hours link hidden; hours by branch report

The Reports menu shows "Branch Hours" only when canViewBranchHours (in $lib/reports/branchHours) allows it. The rule is not in this file.

## KPI report link needs the kpi claim

Lines 302-302

Asked as: KPI report missing; can't see KPI report; KPI link hidden; key performance report access

The Reports menu shows "KPI" only to staff with the kpi claim.
