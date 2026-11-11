<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"86c82e72ff42a293331a2d6b89cf208e84ebbf0d762fb62721612c56c2d40212","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu links with no claim requirement

Lines 32-316

Asked as: link visible to everyone; no permission needed menu; menu link for all staff; why can everyone see this

Other menu links have no claim or check in this file, so they show to all staff. Examples are time entries, sheets, My Purchase Orders, My Expenses, Jobs, Clients, Vendors, My WIP, Rate Sheets, Admin Profiles, Claims and the Settings links. Pages may still restrict access.

## Menu shortcut buttons for creating records and lists

Lines 39-243

Asked as: quick create button; plus button on menu; shortcut button missing; new entry button

Some menu items have shortcut buttons. These include New Entry, New Amendment, New PO, New PO Group, New Expense, New Job, New Project Request, New Client and New Vendor. Jobs also has Unused Active Jobs, Stale Jobs, Latest Jobs and Jobs Map buttons. Amendments has a List All Amendments button.

## Menu links with their own visibility checks

Lines 90-302

Asked as: PO Groups menu hidden; can't see WIP Branch; invoices link missing; branch hours not showing

Some links use a check function instead of a single claim. They are PO Groups (canViewPOGroups), WIP Branch and Division (canViewWIPReports), Invoices (canViewInvoices), PO Approvers (canViewPOApprovers) and Branch Hours (canViewBranchHours). The results of those functions are defined in other files.

## All Active PO list is removed from the menu

Lines 105-107

Asked as: all active PO link gone; missing all active purchase orders; where is all active PO; all active POs removed

The All Active purchase order link is commented out. Search replaces it.

## Uncommit Requests menu link needs commit or admin claim

Lines 133-138

Asked as: can't see uncommit requests; committer link missing; uncommit menu hidden; uncommit option not showing

The Uncommit Requests link in the Expenses menu shows only to a committer (the commit claim) or an admin (the admin claim). Other staff do not see this link.

## WIP menu groups three reports with no page of its own

Lines 180-189

Asked as: WIP menu has no page; clicking WIP does nothing; WIP item not a link; where is my WIP report

The WIP item has no link of its own and groups My WIP, Branch and Division. Branch and Division use the canViewWIPReports check.

## Value Approvals menu link needs project authorization claim

Lines 204-209

Asked as: value approvals not showing; approvals link missing in business; project authorization menu hidden; can't find value approvals

The Value Approvals link in the Business menu requires the project_authorization claim. Staff without this claim do not see it.

## Abandoned Requests menu link is for admins only

Lines 210-215

Asked as: abandoned requests hidden; can't find abandoned project requests; abandoned project requests missing; abandoned queue not visible

The Abandoned Requests link (project requests, queue=abandoned) shows only to an admin (the admin claim).

## App Config and UI Components menu links are for admins only

Lines 255-256

Asked as: app config link missing; UI components link gone; business menu admin links; can't see app config

The App Config and UI Components links in the Business menu require the admin claim. Other staff do not see them.

## Report menu links need the report claim

Lines 262-296

Asked as: payroll report not showing; weekly report missing; payables spreadsheet hidden; CSV report links missing

Payroll, Weekly, Time Off (all staff), Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV links in the Reports menu require the report claim. Staff without it do not see them.

## Two Reports items download CSV files directly

Lines 271-296

Asked as: download CSV button; branch mismatches download; active jobs csv download; CSV button does nothing

Branch Mismatches CSV and Active Jobs have no page. Their Download CSV buttons run downloadTimeEntryBranchMismatchesCsv and downloadActiveJobsCsv.

## KPI report link needs the kpi claim

Lines 303-303

Asked as: KPI report missing; where is the KPI link; KPI menu not visible; can't open KPI reports

The KPI link in the Reports menu requires the kpi claim.
