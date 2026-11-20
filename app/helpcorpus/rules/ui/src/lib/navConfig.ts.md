<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"ba582e033ca19373339f0339ef5f5a072d6787cace4213aeb86c0d1781459b18","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu links with no check are shown to everyone

Lines 17-19

Asked as: menu link shows for everyone; no permission check on link; why can all users see this; menu item visible to all

A menu link or button with no visibility check (canView) is shown to every signed-in user. Turbo only hides an entry when it has a check that fails.

## Claim-based menu checks pass with any one listed claim

Lines 21-25

Asked as: need every role to see menu; menu hidden with only one role; several permissions needed for link; which role unlocks this menu

A menu entry that names several claims (anyClaim) is shown when the user holds at least one of them. The user does not need all of them.

## Project request menu entries need the feature switched on

Lines 27-31

Asked as: project request links missing; can't see project requests menu; project request option gone; requests menu turned off

Menu entries wrapped in the project requests check (jobRequests) show only while project requests are enabled (jobRequestsEnabled). When the feature is off, they are hidden even for users with the right claim.

## WIP report menu check depends on claims and branch manager status

Lines 33-34

Asked as: branch WIP link missing; division WIP not showing; can't see WIP branch menu; WIP menu hidden for me

The WIP Branch and Division links use the WIP reports check (canViewWIPReports). It looks at the user's claims and whether the user is a branch manager. The result of that check is defined in another file.

## Show All UI skips every menu visibility check

Lines 60-68

Asked as: show all ui toggle; see every menu item; menu shows everything to me; why do I see all menus

With Show All UI on, the menu shows every section, link and button whatever the user's claims. The server still checks access when the user opens a page or acts.

## Menu entries without a check show to everyone; empty sections hide

Lines 69-83

Asked as: section disappeared from menu; parent menu item missing; empty menu group hidden; menu group vanished

Without Show All UI, an entry with no visibility check (canView) shows to everyone. Turbo hides a section that has no visible items. It hides a parent item whose children are all hidden.

## Time entries, sheets, shared sheets and time off menu links show to everyone

Lines 89-101

Asked as: entries link for everyone; shared with me missing; sheets menu everyone sees; new entry button visible to all

The Time Management links Entries (with its New Entry button), Sheets, Shared with Me and Time Off have no check. Every user sees them.

## Time sheet approval menu links need the approver claim

Lines 102-113

Asked as: pending my approval timesheets; approved by me timesheet link; can't see timesheet approvals; timesheet approval menu missing

Under Time Management, Pending My Approval and Approved By Me show only to holders of the approver claim (tapr).

## Time tracking and work records menu visibility

Lines 115-126

Asked as: time tracking menu missing; work records link not showing; can't find work records; time tracking link hidden

Time Tracking shows to holders of the report or admin claim. Work Records shows to holders of the report or work_record claim.

## Time amendments menu needs tame or report claim

Lines 127-145

Asked as: time amendments menu; new amendment button missing; amendments link hidden; list all amendments button

Amendments, with its List All and New Amendment buttons, shows to holders of the tame or report claim.

## Time Off link in the time menu has no check

Lines 146-146

Asked as: time off link under time; time off menu item; link to /time/off; time off link visible

The Time Off link under Time Management (/time/off) shows to every user.

## My Purchase Orders and New PO button show to everyone

Lines 152-163

Asked as: my purchase orders link; new PO button for everyone; PO button always visible; purchase orders menu everyone sees

The My Purchase Orders link and its New PO button have no check. Every user sees them.

## PO Groups menu visibility and New PO Group button

Lines 164-177

Asked as: PO groups menu missing; new PO group button; can't see PO groups; purchase order groups link

PO Groups shows to users who pass canViewPOGroups, a check in another file. The New PO Group button shows only to holders of the super_po claim.

## Purchase order Rejected, Search, Stale and Expiring links show to everyone

Lines 178-189

Asked as: rejected purchase orders link; stale POs menu; expiring recurring PO link; search purchase orders link

The Purchase Orders links Rejected (Mine), Search, Stale and Expiring (Recurring) have no check. Every user sees them. The All Active link is commented out and does not appear.

## Purchase order Pending My Approval needs approver claim

Lines 179-184

Asked as: PO pending my approval missing; purchase order approvals menu; can't see POs to approve; approve purchase orders menu

Under Purchase Orders, Pending My Approval shows only to holders of the approver claim (tapr).

## My Expenses, New Expense and My Reimbursements show to everyone

Lines 195-207

Asked as: my expenses link; new expense button missing; my reimbursements link; expense menu for everyone

The My Expenses link, its New Expense button and the My Reimbursements link have no check. Every user sees them.

## Expense approval menu links need approver claim

Lines 208-219

Asked as: expense pending approval missing; approved expenses by me link; can't see expenses to approve; expense approvals menu

Under Expenses, Pending My Approval and Approved By Me show only to holders of the approver claim (tapr).

## Expense tracking, settlement, commit queue and uncommit menu visibility

Lines 220-243

Asked as: expense tracking menu missing; settlement link not showing; commit queue missing; uncommit requests link

Expense Tracking shows to holders of the report or admin claim. Settlement shows to payables staff (payables_admin). Commit Queue shows to committers (commit). Uncommit Requests shows to holders of the commit or admin claim.

## Job list shortcut buttons need the job claim

Lines 252-279

Asked as: unused active jobs button; stale jobs button missing; latest jobs button; jobs map button

On Jobs, the buttons Unused Active Jobs, Stale Jobs and Latest Jobs show only to procurement (the job claim). The Jobs Map button shows to everyone.

## New Project or Proposal button visibility

Lines 280-290

Asked as: new project button missing; new proposal button; can't create project button; create proposal button hidden

The button shows to procurement (job claim). It also shows to holders of the time claim while the project requests feature is on (jobRequestsEnabled).

## WIP menu: My WIP for all, Branch and Division restricted

Lines 293-302

Asked as: my WIP link; branch WIP hidden; division WIP missing; WIP menu not showing

My WIP shows to everyone. Branch and Division WIP show only to users who pass canViewWIPReports, which depends on claims and on being a branch manager.

## Project Requests menu needs feature on and a claim

Lines 303-314

Asked as: project requests menu missing; request approvals link; project authorization menu; project requests link hidden

Project Requests shows only while project requests are on, and only to holders of the time, job or project_authorization claim. Request Approvals shows only while project requests are on, and only to holders of project_authorization.

## Job Create Queue menu visibility

Lines 315-321

Asked as: job create queue missing; create queue menu; job queue link not showing; can't see create queue

Job Create Queue shows to users who pass canViewJobCreateQueue, a check in another file. It uses the claims, the project requests switch and the attention counts.

## Value Approvals and Abandoned Requests menu visibility

Lines 322-333

Asked as: value approvals link; abandoned requests missing; value approvals menu hidden; abandoned requests link

Value Approvals shows to holders of project_authorization. Abandoned Requests shows to holders of admin.

## Rate Sheets menu visibility

Lines 334-344

Asked as: rate sheets menu missing; rate sheet link hidden; can't see rate sheets; rate sheets not showing

Rate Sheets shows to holders of the job, admin or rate_sheet_revise claim. While project requests are on, it also shows to holders of time or project_authorization.

## Clients and Vendors add buttons need claims

Lines 345-370

Asked as: new client button missing; new vendor button; clients link; add vendor button hidden

Clients and Vendors list links show to everyone. New Client shows to procurement (job claim). New Vendor shows to payables staff (payables_admin).

## Absorb Actions menu needs absorb claim

Lines 371-376

Asked as: absorb actions menu missing; can't see absorb menu; absorb link hidden; absorb actions not showing

Absorb Actions shows only to holders of the absorb claim.

## PO Approvers and Computers menu use other checks

Lines 377-413

Asked as: PO approvers menu; computers menu missing; retired computers button; logins button hidden

PO Approvers shows to users who pass canViewPOApprovers. Computers, with its Retired Computers, Logins and Raw Logins buttons, shows to users who pass canViewComputers. Both checks are in other files.

## Admin Profiles menu visibility

Lines 383-388

Asked as: admin profiles menu; profiles link missing; admin profiles not showing; can't see admin profiles

Admin Profiles shows to holders of the admin, hr, time_off_manager or it claim.

## Claims visible to all; App Config and UI Components admin-only

Lines 414-416

Asked as: claims menu for everyone; app config missing; UI components link hidden; app config admin only

Claims shows to everyone. App Config and UI Components show only to holders of the admin claim.

## Report menu items need report claim

Lines 422-456

Asked as: payroll report missing; weekly report menu; download CSV button missing; active jobs report link

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs reports show only to holders of the report claim. The last two have a Download CSV button.

## Branch Hours and KPI report visibility

Lines 457-463

Asked as: branch hours link; KPI report menu; branch hours not showing; KPI link missing

Branch Hours shows to users who pass canViewBranchHours, which depends on claims and on being a branch manager. KPI shows to holders of the kpi claim.

## Settings menu is admin only

Lines 466-485

Asked as: currencies menu missing; time types settings; machine secrets link; attachment audit menu

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only to holders of the admin claim.
