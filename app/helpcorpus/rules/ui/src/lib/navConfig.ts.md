<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"0b04f6d3701f7335a501c3b3a1c5971abe8e642490b1f64a612907806c4e99e0","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu entries without a check are shown to everyone

Lines 17-83

Asked as: menu item showing for everyone; link visible without permission; why do I see this menu entry; empty section hidden; menu items without access check

Without Show All UI, Turbo shows an entry that has no view check to everyone. An entry with a check appears only when the check passes. Turbo hides a section that has no visible items.

## WIP menu hides when it has no visible children

Lines 33-35

Asked as: WIP menu missing; can't see WIP menu; WIP section disappeared; no WIP reports showing; branch and division missing

The WIP menu item appears only if at least one of its child links is visible. My WIP is always shown. Branch and Division appear only when the person may view WIP reports (canViewWIPReports, which uses the person's claims and whether they are a branch manager).

## Show All UI shows every menu entry

Lines 60-68

Asked as: show all UI; see the whole menu; show every menu item; display all menu options; turn on full menu

With Show All UI on, Turbo shows every section, link and button in the menu, whatever the person's claims. The server still checks access when the person uses them.

## Everyone sees Entries link and New Entry button

Lines 89-100

Asked as: where is new entry button; time entries link everyone; add time entry button visible; entries button always there; green new entry button

The Time Management menu shows the Entries link (/time/entries/list) to every signed-in user. The green New Entry button (/time/entries/add) beside it has no claim check, so everyone sees it too.

## Time sheet approval links need the approver claim

Lines 101-114

Asked as: no pending approvals link; can't see approve time sheets; missing approved by me; time approval menu missing; approver links not showing

Under Time Management, Pending My Approval and Approved By Me appear only to approvers (the tapr claim). Entries, Sheets, Shared with Me and Time Off appear to everyone.

## Time tracking and work records menu access

Lines 115-126

Asked as: no time tracking link; can't see work records; tracking menu missing; work records not in menu; time report links hidden

The Time Management Tracking link appears to holders of the report or admin claim. Work Records appears to holders of the report or work_record claim.

## Time amendments menu access

Lines 127-145

Asked as: amendments menu missing; no list all amendments; can't see new amendment; time amendment button gone; amendments link not showing

The Amendments link, with its List All Amendments and New Amendment buttons, appears to holders of the tame (time amendment) or report claim.

## Purchase order menu links

Lines 152-189

Asked as: all active link gone; purchase order approval link missing; no pending PO approvals; PO menu links missing; stale POs link

Under Purchase Orders, My Purchase Orders (with New PO), Rejected (Mine), Search, Stale and Expiring (Recurring) appear to everyone. Pending My Approval appears only to approvers (the tapr claim). The All Active link is removed.

## PO Groups menu access

Lines 164-177

Asked as: PO groups missing; can't add PO group; new group button not showing; no PO group menu; PO group link hidden

The PO Groups link appears when canViewPOGroups passes for the person's claims. The New PO Group button appears only to holders of the super_po claim.

## Expense menu links for everyone and approvers

Lines 195-219

Asked as: expense approvals missing; can't see pending expenses; no new expense button; reimbursements link missing; approved expenses link gone

My Expenses (with New Expense) and My Reimbursements appear to everyone. Pending My Approval and Approved By Me appear only to approvers (the tapr claim).

## Expense tracking, settlement, commit and uncommit menu access

Lines 220-243

Asked as: settlement menu missing; commit queue not showing; uncommit requests gone; expense tracking link hidden; no commit menu

Expense Tracking appears to holders of the report or admin claim. Settlement appears to payables admins (payables_admin). Commit Queue appears to committers (the commit claim). Uncommit Requests appears to holders of the commit or admin claim.

## Job list extra buttons need the job claim

Lines 249-279

Asked as: jobs buttons missing; stale jobs button gone; unused active jobs not showing; latest jobs button missing; jobs map button

On the Jobs menu item, the Unused Active Jobs, Stale Jobs and Latest Jobs buttons appear only to procurement (the job claim). The Jobs Map button appears to everyone.

## New Project or Proposal button visibility

Lines 280-290

Asked as: can't create project; no new project button; new proposal button missing; add project request button; how to start a project request

The New Project or Proposal button appears to procurement (the job claim). It also appears to holders of the time claim while the project requests feature is on (jobRequestsEnabled). Then /jobs/add sends them to a new request.

## My WIP menu link is shown to everyone

Lines 293-298

Asked as: my wip link missing; my wip not showing; personal wip hidden; wip link for everyone

Under the WIP menu, the My WIP link (/jobs/wip/my) has no claim check, so every user sees it. The Branch and Division WIP links under it are checked separately.

## Project request menu links depend on the feature

Lines 303-314

Asked as: project requests menu missing; request approvals not showing; no project requests link; project request menu gone; approvals link missing

Project Requests appears to holders of the time, job or project_authorization claim, only while the project requests feature is on. Request Approvals appears to holders of project_authorization, only while the feature is on.

## Job Create Queue menu access

Lines 315-321

Asked as: job create queue missing; can't see create queue; create queue link hidden; job queue not showing; queue menu item gone

The Job Create Queue link appears when canViewJobCreateQueue passes. It looks at the person's claims, whether the project requests feature is on, and the attention counts. This file does not state its result.

## Value approvals and abandoned requests menu access

Lines 322-333

Asked as: value approvals missing; abandoned requests not showing; can't see value approvals; abandoned requests link; value approval menu hidden

Value Approvals appears to holders of project_authorization, even when the project requests feature is off. Abandoned Requests appears only to admins.

## Rate sheets and client menu access

Lines 334-352

Asked as: rate sheets missing; new client button gone; can't add client; clients link missing; rate sheet link hidden

The Rate Sheets link appears only to procurement (the job claim). The Clients link appears to everyone, but its New Client button appears only to procurement.

## Vendor menu access

Lines 353-365

Asked as: new vendor button missing; vendors link not showing; can't add vendor; no vendor button; vendor menu hidden

The Vendors link appears to everyone. The New Vendor button appears only to payables admins (payables_admin).

## Invoices, absorb and PO approver menu access

Lines 366-383

Asked as: invoices menu missing; absorb actions not showing; PO approvers link missing; can't see invoices; absorb link gone

Invoices appears when canViewInvoices passes for the person's claims. Absorb Actions appears to holders of the absorb claim. PO Approvers appears when canViewPOApprovers passes for the person's claims.

## Admin profiles menu access

Lines 384-389

Asked as: admin profiles missing; no admin profiles link; profiles menu hidden; can't see admin profiles; hr profiles link

The Admin Profiles link appears to holders of the admin, hr, time_off_manager or it claim.

## Claims link is shown to everyone

Lines 390-390

Asked as: claims link missing; can't see claims; claims menu hidden; where are my claims; claims not in menu

The Claims link appears to everyone.

## App Config and UI Components are admin only

Lines 391-392

Asked as: app config missing; ui components link gone; can't see app config; admin menu not showing; ui components hidden

The App Config and UI Components links appear only to admins.

## Report menu access by claim

Lines 396-439

Asked as: payroll report missing; weekly report not showing; KPI link missing; can't download CSV; reports menu hidden

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs appear only to holders of the report claim. The last two are CSV download buttons. KPI appears to holders of the kpi claim.

## Branch Hours report menu access

Lines 433-438

Asked as: branch hours missing; can't see branch hours; branch hours link gone; branch hours report hidden; no branch hours menu

The Branch Hours link appears when canViewBranchHours passes. It looks at the person's claims and whether they are a branch manager. This file does not state its result.

## Settings menu is admin only

Lines 442-461

Asked as: settings menu missing; no currencies link; time types not showing; divisions link gone; can't see settings; machine secrets missing

The Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit links appear only to admins. Other people do not see the Settings section.
