<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"7fe2eda97a3d497beee71c02b8cc731eaa2d0006ba1560db9beb7740a005a8fa","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu checks use the claims of the signed-in user

Lines 9-25

Asked as: menu depends on my permissions; why menu items missing; links change per user; menu based on my role

A menu check (canView) reads the claims of the signed-in user, whether the user is a branch manager (isBranchManager), whether project requests are on (jobRequestsEnabled) and the attention counts. The anyClaim helper shows an entry when the user holds at least one of the named claims.

## Menu entries without a check show to everyone

Lines 17-19

Asked as: why can everyone see this link; link visible to all users; no permission needed for link; menu link open to everyone

A menu link or button with no view check (canView) appears for every signed-in user. Entries such as Entries, Sheets, Shared with Me, Time Off, My Purchase Orders, Search and My Expenses are open to all in the menu.

## Project request menus need the feature on

Lines 28-31

Asked as: project requests menu missing; can't see project requests; project requests link hidden; request approvals not showing

Project Requests shows only while project requests are on, to holders of the time, job or project_authorization claim. Request Approvals shows only while the feature is on, to holders of project_authorization.

## Branch and Division WIP menus use the WIP reports check

Lines 33-34

Asked as: branch WIP link missing; division WIP not showing; WIP reports menu hidden; can't open branch WIP

The Branch and Division WIP links show only when the WIP reports check passes. That check is canViewWIPReports, which takes the user's claims and whether the user is a branch manager. This file does not show its result.

## Show All UI displays every menu entry

Lines 60-68

Asked as: show all menu items; see every menu option; display all UI; show every link

With Show All UI on, Turbo shows every section, link and button, whatever the user's claims. The server still checks access to the actions.

## Menu hides links, buttons and empty sections

Lines 69-82

Asked as: section disappeared from menu; menu link hidden; empty menu group missing; why is the menu missing

Without Show All UI, Turbo hides each link or button whose check fails. A section left with no links is hidden. The WIP group is hidden when none of its child links is visible.

## Time Off menu link is open to everyone

Lines 89-101

Asked as: shared with me link; sheets menu for everyone; entries new entry button; everyone sees entries

The Time Off link in the Time Management menu has no check (canView), so every signed-in user sees it. The Sheets, Shared with Me and Entries links are open in the same way. The Entries link also has a New Entry button that everyone sees.

## Time sheet approval menus need the approver claim

Lines 102-113

Asked as: pending my approval timesheets; approved by me timesheets; timesheet approval menu missing; approve time sheets link

Pending My Approval and Approved By Me under Time Management show only to an approver (the tapr claim).

## Time tracking and work records menu claims

Lines 115-126

Asked as: time tracking menu missing; work records link hidden; work records not showing; time tracking link gone

Time Tracking shows to holders of the report or admin claim. Work Records shows to holders of the report or work_record claim.

## Amendments menu needs time amendment or report claim

Lines 127-145

Asked as: amendments menu missing; new amendment button; list all amendments; amendments link hidden

The Amendments link, with its List All Amendments and New Amendment buttons, shows to holders of the tame or report claim.

## Time Off menu link shows to all staff

Lines 146-146

Asked as: time off menu link; who can see time off; time off link showing; time off menu for staff

The Time Off link in the Time Management menu (/time/off) has no claim check. Every signed-in user sees it.

## My Purchase Orders menu and New PO button are open

Lines 152-163

Asked as: my purchase orders button; new PO button visible; my POs link; new PO button everyone

The My Purchase Orders link and its New PO button have no claim check. Every signed-in user sees both.

## PO Groups menu and New PO Group button

Lines 164-177

Asked as: PO groups menu missing; new PO group button; can't create PO group; PO group link hidden

PO Groups shows to users who pass canViewPOGroups (its rule is in another file). The New PO Group button shows only to holders of the super_po claim.

## Rejected, Search and My Expenses menus are open

Lines 178-214

Asked as: rejected purchase orders link; my expenses link; my reimbursements menu; expenses menu for all

The Rejected (Mine) and Search links in the Purchase Orders menu have no claim check, so everyone sees them. The My Expenses link, its New Expense button and the My Reimbursements link are also open to everyone.

## Purchase order approval menu needs approver claim

Lines 179-184

Asked as: pending PO approval menu; purchase order approvals missing; PO approval link hidden; approve purchase orders menu

Pending My Approval under Purchase Orders shows only to an approver (the tapr claim).

## Stale and Expiring PO lists are for payables staff

Lines 188-196

Asked as: stale PO list; expiring recurring POs; payables stale purchase orders; stale orders list missing

The company-wide Stale and Expiring (Recurring) purchase order lists show only to payables staff (the payables_admin claim). Other users see a label on the purchase order in their own list.

## Expense approval menus need approver claim

Lines 215-226

Asked as: pending expense approval; approved expenses by me; expense approval menu missing; expense approvals link hidden

Pending My Approval and Approved By Me under Expenses show only to an approver (the tapr claim).

## Expense tracking, settlement and commit menus

Lines 227-250

Asked as: expense tracking menu missing; settlement menu hidden; commit queue not showing; uncommit requests link

Expense Tracking shows to holders of the report or admin claim. Settlement shows to payables staff (payables_admin). Commit Queue shows to a committer (the commit claim). Uncommit Requests shows to holders of the commit or admin claim.

## Job list buttons for procurement

Lines 256-286

Asked as: unused active jobs button; stale jobs button; latest jobs missing; jobs map button

The Unused Active Jobs, Stale Jobs and Latest Jobs buttons show only to procurement (the job claim). The Jobs Map button shows to everyone.

## New Project or Proposal button visibility

Lines 287-297

Asked as: new project button missing; new proposal button; create project or proposal; can't add new project

The button shows to procurement (the job claim). It also shows to staff with the time claim while project requests are on (jobRequestsEnabled).

## WIP menu: My WIP open, Branch and Division restricted

Lines 300-309

Asked as: my WIP link; WIP menu for everyone; branch WIP hidden; division WIP restricted

My WIP shows to everyone. Branch and Division WIP show only to users who pass canViewWIPReports (claims and branch manager status; its rule is in another file).

## Request Approvals menu

Lines 316-321

Asked as: request approvals menu; request approvals missing; approve project requests link; request approvals hidden

Request Approvals shows only while project requests are on, and only to holders of the project_authorization claim.

## Job Create Queue menu

Lines 322-328

Asked as: job create queue menu; job creation queue missing; create queue link hidden; job queue not showing

Job Create Queue shows to users who pass canViewJobCreateQueue, which takes their claims, the project requests setting and the attention counts. Its result is set in another file.

## Value Approvals and Abandoned Requests menus

Lines 329-340

Asked as: value approvals menu; abandoned requests link; value approvals missing; abandoned requests admin only

Value Approvals shows to holders of the project_authorization claim. Abandoned Requests shows only to admins (the admin claim). Value Approvals does not depend on the project requests setting.

## Rate Sheets menu visibility

Lines 341-351

Asked as: rate sheets menu missing; rate sheets link hidden; can't see rate sheets; rate sheet menu not showing

Rate Sheets shows to holders of the job, admin or rate_sheet_revise claim. It also shows to holders of the time or project_authorization claim while project requests are on.

## Clients and Vendors add buttons

Lines 352-377

Asked as: new client button; new vendor button; add vendor missing; add client not showing

The Clients and Vendors lists show to everyone. New Client shows only to procurement (the job claim). New Vendor shows only to payables staff (payables_admin).

## Absorb Actions menu needs absorb claim

Lines 378-383

Asked as: absorb actions menu; absorb menu missing; absorb link hidden; can't see absorb

Absorb Actions shows only to holders of the absorb claim.

## PO Approvers menu

Lines 384-389

Asked as: PO approvers menu; approvers link missing; PO approvers hidden; purchase order approvers menu

PO Approvers shows to users who pass canViewPOApprovers (its rule is in another file).

## Admin Profiles menu claims

Lines 390-395

Asked as: admin profiles menu; profiles link hidden; admin profiles missing; can't see admin profiles

Admin Profiles shows to holders of the admin, hr, time_off_manager or it claim.

## Computers menu and its buttons

Lines 396-420

Asked as: computers menu missing; retired computers button; logins button missing; raw logins not showing

Computers shows to users who pass canViewComputers (its rule is in another file). Its buttons Retired Computers, Logins and Raw Logins have no check of their own, so they show whenever the Computers link shows.

## Claims link is open; App Config and UI Components are admin only

Lines 421-423

Asked as: claims menu link; app config menu; UI components menu; app config hidden

Claims shows to everyone. App Config and UI Components show only to admins (the admin claim).

## Report menus need the report claim

Lines 429-463

Asked as: payroll report menu; weekly report link; active jobs CSV menu; branch mismatches CSV

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV show only to holders of the report claim.

## Branch Hours and KPI report menus

Lines 464-470

Asked as: branch hours menu; KPI report menu; branch hours link missing; KPI link hidden

Branch Hours shows to users who pass canViewBranchHours (claims and branch manager status; its rule is in another file). KPI shows only to holders of the kpi claim.

## Settings menu is admin only

Lines 473-492

Asked as: settings menu missing; currencies menu; time types menu; machine secrets menu

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only to admins (the admin claim).
