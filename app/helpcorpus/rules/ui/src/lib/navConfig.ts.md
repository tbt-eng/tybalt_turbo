<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"f14bebf704caf3100c3cb83a9ab628167742e75a17418c276f132618e2bb8bd6","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu shortcut buttons for creating records

Lines 35-243

Asked as: green plus button; add new record from menu; new entry shortcut; list all amendments button; create new expense from menu

Some menu items have a green plus button that opens the add page: New Entry, New Amendment, New PO, New PO Group, New Expense, New Job, New Project Request, New Client and New Vendor. The Amendments item also has a List All Amendments button.

## Menu links decided by helper functions

Lines 89-101

Asked as: PO Groups link missing; PO Approvers link hidden; Branch Hours not showing; why can't I see PO groups; menu link hidden for me

The PO Groups link uses canViewPOGroups, the PO Approvers link uses canViewPOApprovers, and the Branch Hours link uses canViewBranchHours. Each link shows only when its function allows it. Who passes is defined in other files.

## All Active purchase orders link is removed

Lines 104-106

Asked as: all active purchase orders; no All Active link; missing active POs menu; where is All Active; purchase orders menu search

The Purchase Orders menu has no All Active link. It is commented out in favour of Search.

## Uncommit Requests menu link needs commit or admin claim

Lines 132-137

Asked as: cannot see uncommit requests; uncommit link missing; no uncommit option in expenses; need committer access; who can uncommit requests

The Uncommit Requests link in the Expenses menu shows only to a committer (the commit claim) or an admin (the admin claim). Other staff do not see the link.

## Branch and Division WIP links depend on canViewWIPReports

Lines 179-188

Asked as: WIP branch link missing; division WIP not showing; can't see branch WIP report; My WIP only showing; WIP reports menu hidden

The Branch and Division links under WIP show only when canViewWIPReports allows it. The My WIP link shows to everyone. The code of canViewWIPReports is in another file, so this file does not say who passes.

## Value Approvals menu link needs project authorization claim

Lines 203-208

Asked as: value approvals not showing; no value approvals in business menu; missing approvals link; project authorization menu item; can't find value approvals

The Value Approvals link in the Business menu requires the project_authorization claim. Staff without it do not see the link.

## Admin-only menu links

Lines 209-254

Asked as: abandoned requests missing; no app config in menu; UI components link gone; business menu admin items hidden; why can't I see admin links

The Abandoned Requests, App Config and UI Components links in the Business menu require the admin claim. Other staff do not see them.

## Report menu links need the report claim

Lines 259-294

Asked as: payroll report missing; weekly report not in reports menu; time off report hidden; payables spreadsheet not showing; reports menu empty for me

The Reports menu links Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs require the report claim. Staff without it do not see them.

## Report CSV downloads run from menu buttons

Lines 269-294

Asked as: download CSV button missing; branch mismatches download; active jobs csv; no page for csv report; export jobs csv from menu

The Branch Mismatches CSV and Active Jobs items have no page. Each has a Download CSV button that downloads the file directly. Both require the report claim.

## KPI report link needs kpi claim

Lines 301-301

Asked as: KPI report not showing; no KPI in reports; kpi menu missing; can't open KPI report; kpi link hidden

The KPI link in the Reports menu requires the kpi claim.

## Settings menu links have no claim check

Lines 304-313

Asked as: settings links visible to everyone; currencies menu shown to all; time types not hidden; machine secrets visible; anyone can see settings menu

The Settings section of the menu lists Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit. None of these links sets a required claim (requiredClaim) or a visibility check (canView) in this file. The menu config therefore does not hide them from any user.
