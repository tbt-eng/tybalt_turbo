<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"02d708bd5d43e3c301d4fe294bf2051ab88be39566657ddcea393f61edf26599","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu checks can use claims, branch manager status, feature flag and attention counts

Lines 8-14

Asked as: menu visibility check; what decides which menu items I see; menu item hidden for me; claims and branch manager menu

A menu check can look at the person's claims, whether they are a branch manager, whether project requests are enabled, and the attention counts (attentionCounts).

## A menu entry with no check is shown to everyone

Lines 16-18

Asked as: menu link shows for everyone; no check on menu item; why can everyone see this link; menu button visible to all users

A menu link or button that has no visibility check (canView) is shown to every signed-in user. A check only narrows who sees an entry that has one.

## Claim-based menu checks pass with any one named claim

Lines 20-24

Asked as: any one of several claims; menu needs one of the roles; only need one claim; several claims menu access

A menu entry that names several claims (anyClaim) is shown to a person who holds at least one of them. The person does not need all of them.

## Project request menu entries need the feature switched on

Lines 26-30

Asked as: project requests menu missing; project request link not showing; feature switched off menu; project requests disabled no menu

A menu entry wrapped in the project requests check (jobRequests) is shown only while project requests are enabled (jobRequestsEnabled). With the feature off, nobody sees it, whatever claims they hold. When the feature is on, the entry's own check must also pass.

## WIP report links use the WIP report check

Lines 32-33

Asked as: branch WIP link missing; division WIP link not showing; who can see WIP reports; WIP report links hidden

The Branch and Division WIP links use the check canViewWIPReports, which looks at the person's claims and whether they are a branch manager (isBranchManager). This file does not say which claims pass that check.

## Show All UI shows every menu link and button

Lines 59-67

Asked as: show all UI setting; show every menu item; see all menu links; show all menu option

With Show All UI on, the navigation menu shows every section, link and button whatever the user's claims. The server still checks access when the user acts.

## Menu entries without a check are visible to everyone; empty groups are hidden

Lines 68-81

Asked as: empty menu section hidden; whole WIP group missing; menu group disappeared; section not showing in menu

A menu link or button with no canView check is shown to everyone. A section with no visible links is hidden. The WIP group is hidden when none of its child links is visible.

## A menu group with a hidden submenu link set disappears

Lines 77-79

Asked as: submenu disappeared; parent menu hidden no children; all sub-links hidden menu; dropdown menu missing

When an entry has sub-links (children) and none of them pass their checks, the entry is hidden. The entry is also hidden when it has an empty children list.

## Time Management group has links for everyone

Lines 84-100

Asked as: time management links for everyone; new entry button missing; shared with me not showing; time off link everyone

Entries, Sheets, Shared with Me and Time Off have no check, so everyone sees them. Everyone also sees the New Entry button.

## Time sheet approval links need the approver claim

Lines 100-113

Asked as: pending my approval missing timesheet; approved by me not showing; timesheet approver menu; can't see timesheet approvals

Time Management shows 'Pending My Approval' and 'Approved By Me' only to approvers (the tapr claim). 'Sheets', 'Shared with Me', 'Entries' and 'Time Off' are shown to everyone.

## Time tracking and work records menu access

Lines 114-125

Asked as: tracking link missing time; work records not showing; time tracking menu hidden; work record menu access

'Tracking' under Time Management is shown to holders of the report or admin claim. 'Work Records' is shown to holders of the report or work_record claim.

## Amendments menu access

Lines 126-144

Asked as: amendments menu missing; no amendments button; list all amendments hidden; new amendment not showing

'Amendments' and its buttons 'List All Amendments' and 'New Amendment' are shown to holders of the tame or report claim.

## Purchase order menu: approval link needs approver claim

Lines 151-188

Asked as: PO approval link missing; purchase order approvals not showing; pending PO approvals menu; can't see POs to approve

'Pending My Approval' under Purchase Orders is shown only to approvers (the tapr claim). 'My Purchase Orders' with 'New PO', 'Rejected (Mine)', 'Search', 'Stale' and 'Expiring (Recurring)' are shown to everyone.

## PO Groups menu and New PO Group button

Lines 163-176

Asked as: PO groups link missing; new PO group button missing; no PO group option; purchase order groups hidden

The 'PO Groups' link appears when canViewPOGroups (in poGroups) accepts the user's claims. The 'New PO Group' button appears only to holders of the super_po claim.

## Expense approval links need approver claim

Lines 194-218

Asked as: expense approvals not showing; approved expenses menu missing; pending expense approval link; expense approver menu

'Pending My Approval' and 'Approved By Me' under Expenses are shown only to approvers (the tapr claim). 'My Expenses', 'New Expense' and 'My Reimbursements' are shown to everyone.

## Expense tracking, settlement and commit menu access

Lines 219-242

Asked as: expense tracking missing; settlement link not showing; commit queue hidden; uncommit requests missing

Expense 'Tracking' is shown to the report or admin claim. 'Settlement' is shown to payables staff (payables_admin). 'Commit Queue' is shown to committers (commit). 'Uncommit Requests' is shown to the commit or admin claim.

## Job list shortcut buttons for procurement

Lines 251-278

Asked as: unused active jobs button missing; stale jobs button hidden; latest jobs not showing; jobs map button missing

The buttons 'Unused Active Jobs', 'Stale Jobs' and 'Latest Jobs' are shown only to procurement (the job claim). The 'Jobs Map' button is shown to everyone.

## New Project or Proposal button visibility

Lines 279-289

Asked as: new project button missing; can't create proposal; new proposal button hidden; create project button not showing

The 'New Project or Proposal' button is shown to procurement (the job claim). It is also shown to staff with the time claim while project requests are enabled (jobRequestsEnabled).

## WIP menu links

Lines 292-301

Asked as: my WIP missing; branch WIP not showing; division WIP link hidden; WIP menu missing

'My WIP' is shown to everyone. 'Branch' and 'Division' WIP links appear when canViewWIPReports (in wipReports) accepts the user's claims and branch manager status.

## Project Requests menu needs feature on and a claim

Lines 302-313

Asked as: project requests link missing; request approvals not showing; project requests hidden; no project requests menu

'Project Requests' is shown only while project requests are enabled, and only to holders of the time, job or project_authorization claim. 'Request Approvals' is shown only while project requests are enabled, and only to holders of project_authorization.

## Job Create Queue menu visibility

Lines 314-320

Asked as: job create queue missing; create queue link hidden; job queue menu not showing; no job create queue

'Job Create Queue' appears when canViewJobCreateQueue (in jobRequests) accepts the user's claims, the project requests setting and the attention counts.

## Value Approvals and Abandoned Requests menu access

Lines 321-332

Asked as: value approvals missing; abandoned requests not showing; value approval link hidden; abandoned requests menu

'Value Approvals' is shown to holders of project_authorization. 'Abandoned Requests' is shown only to admins (admin).

## Rate sheets and client creation menu access

Lines 333-351

Asked as: rate sheets link missing; new client button missing; can't add client; clients menu hidden

'Rate Sheets' is shown only to procurement (job). The 'New Client' button is also shown only to procurement. 'Clients' is shown to everyone.

## New Vendor button needs payables staff

Lines 352-364

Asked as: new vendor button missing; vendor creation hidden; can't add vendor; vendors menu button

The 'New Vendor' button is shown only to payables staff (payables_admin). 'Vendors' is shown to everyone.

## Absorb Actions, PO Approvers and Admin Profiles menu access

Lines 365-382

Asked as: absorb actions missing; PO approvers link hidden; admin profiles not showing; absorb menu missing

'Absorb Actions' is shown to holders of the absorb claim. 'PO Approvers' appears when canViewPOApprovers (in poApprovalLimits) accepts the user's claims. 'Admin Profiles' is shown to the admin, hr, time_off_manager or it claim.

## Admin-only menu entries

Lines 383-385

Asked as: app config missing; UI components menu hidden; claims menu link; admin menu entries not showing

'App Config' and 'UI Components' are shown only to admins (admin). 'Claims' is shown to everyone.

## Report menu links need the report claim

Lines 391-425

Asked as: payroll report link missing; weekly report hidden; download CSV button missing; active jobs report not showing

Under Reports, 'Payroll', 'Weekly', 'Time Off', 'Payables Spreadsheet', 'Branch Mismatches CSV' and 'Active Jobs' are shown only to holders of the report claim. The last two give a Download CSV button.

## Branch Hours and KPI report access

Lines 426-432

Asked as: branch hours link missing; KPI report not showing; KPI menu hidden; branch hours menu

'Branch Hours' appears when canViewBranchHours (in branchHours) accepts the user's claims and branch manager status. 'KPI' is shown to holders of the kpi claim.

## Settings section is admin only

Lines 435-454

Asked as: settings links missing; currencies menu hidden; time types not showing; machine secrets menu

The Settings links 'Currencies', 'Time Types', 'Divisions', 'Machine Secrets' and 'Attachment Audit' are shown only to admins (admin).
