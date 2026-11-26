<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"8e02f5fdb671c33734bec88289ce73af19f688854cd8c67c37488b125da23c83","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu entries without a check are shown to everyone

Lines 17-19

Asked as: everyone sees this link; menu item has no check; visible to all users; why can everyone see this

A menu link or button with no canView check is shown to everyone. A check hides the entry from users who fail it.

## A menu link with several claims needs only one

Lines 21-25

Asked as: need only one role; any one permission enough; menu needs one of several; do I need every claim

Where a menu link or button lists several claims, Turbo shows it to anyone who holds at least one of them (anyClaim). Holding all of the claims is not required.

## Project request menu links need the feature on

Lines 27-31

Asked as: project requests menu missing; request approvals link hidden; can't see project requests; project request link gone

"Project Requests" shows for holders of the time, job or project_authorization claim, only while project requests are on. "Request Approvals" shows for holders of project_authorization, only while the feature is on.

## WIP Branch and Division links depend on the WIP reports check

Lines 33-34

Asked as: branch wip link missing; division wip hidden; wip reports menu not showing; no wip branch option

The WIP menu links Branch and Division use the WIP reports check (wipReports). That check calls canViewWIPReports with the person's claims and whether they are a branch manager. The result of that function is defined in another file.

## Show All UI shows every menu entry

Lines 60-68

Asked as: show all ui toggle; see every menu item; all menu entries showing; menu shows everything

When Show All UI is on, Turbo shows every section, link and button in the menu whatever the user's claims. The server still checks access when the user acts.

## Empty menu sections and empty WIP group are hidden

Lines 69-83

Asked as: menu section disappeared; empty menu group missing; wip group gone; menu group hidden

Turbo hides a menu section that has no visible items. It also hides a link that has children when none of its children are visible.

## Time sheet approval lists need the approver claim

Lines 89-114

Asked as: time pending my approval missing; approved by me timesheets; time approver list hidden; can't see sheets to approve

Time sheet "Pending My Approval" and "Approved By Me" show only for approvers (the tapr claim). Sheets, Shared with Me, Entries and Time Off show for everyone.

## Time tracking and work records menu access

Lines 115-126

Asked as: time tracking menu missing; work records link hidden; can't find time tracking; work records not showing

Time "Tracking" shows for holders of the report or admin claim. "Work Records" shows for holders of the report or work_record claim.

## Time amendments menu access

Lines 127-145

Asked as: amendments link missing; new amendment button hidden; list all amendments missing; amendments menu not showing

The "Amendments" link, with its List All and New Amendment buttons, shows only for holders of the tame or report claim.

## Time Off menu link is shown to everyone

Lines 146-146

Asked as: time off link visible; time off menu for everyone; why do I see time off; time off link shown

The Time Off link in the Time Management menu (/time/off) has no check, so every signed-in person sees it.

## My Purchase Orders link and New PO button are shown to everyone

Lines 152-163

Asked as: new po button missing; my purchase orders link; where is new po; my pos menu shown

The My Purchase Orders link (/pos/list) and its New PO button (/pos/add) have no check. Every signed-in person sees both.

## PO Groups menu and New PO Group button

Lines 164-177

Asked as: po groups menu missing; new po group button; can't create po group; po group button hidden

"PO Groups" shows when canViewPOGroups accepts the user's claims. The "New PO Group" button needs the super_po claim.

## Purchase order menu: approver and payables lists

Lines 178-196

Asked as: po pending my approval missing; stale purchase orders hidden; expiring recurring pos missing; po label in my list

PO "Pending My Approval" needs the approver claim (tapr). The company-wide "Stale" and "Expiring (Recurring)" lists need the payables_admin claim. Other users see a label on the PO in their own list. My Purchase Orders, Rejected (Mine) and Search show for everyone.

## Expense approval lists need the approver claim

Lines 202-226

Asked as: expense pending approval missing; approved by me expenses; expenses to approve not showing; expense approver list hidden

Expense "Pending My Approval" and "Approved By Me" show only for approvers (the tapr claim). My Expenses and My Reimbursements show for everyone.

## Expense tracking, settlement and commit menu access

Lines 227-250

Asked as: expense settlement missing; commit queue not showing; uncommit requests hidden; expense tracking link gone

Expense "Tracking" needs the report or admin claim. "Settlement" needs payables_admin. "Commit Queue" needs the commit claim. "Uncommit Requests" needs the commit or admin claim.

## Job report buttons need the job claim

Lines 259-293

Asked as: unused active jobs button; stale jobs button missing; jobs map button; latest jobs hidden

The Jobs buttons "Unused Active Jobs", "Stale Jobs", "Projects Without an Invoicing Profile" and "Latest Jobs" show only for procurement (the job claim). The Jobs Map button shows for everyone.

## New Project or Proposal button visibility

Lines 294-304

Asked as: new project button missing; can't request new project; add project proposal button; new proposal button hidden

The "New Project or Proposal" button shows for procurement (the job claim). It also shows for staff with the time claim, but only while project requests are on. Then /jobs/add sends them to a new request.

## WIP menu links

Lines 307-316

Asked as: my wip link; wip menu not showing; wip links missing; which wip menus show

"My WIP" shows for everyone. "Branch" and "Division" WIP show only when canViewWIPReports accepts the user's claims and branch manager status.

## Project request links need project requests enabled

Lines 317-328

Asked as: project requests link missing; request approvals hidden; project request menu gone; project requests turned off

The Project Requests link needs the time, job or project_authorization claim and the feature on. The Request Approvals link needs project_authorization and the feature on.

## Job Create Queue menu link

Lines 329-335

Asked as: job create queue missing; create queue menu hidden; job queue link gone; job create link

"Job Create Queue" shows when canViewJobCreateQueue accepts the user's claims, the project requests setting and the attention counts.

## Value Approvals and Abandoned Requests menu access

Lines 336-347

Asked as: value approvals link missing; abandoned requests hidden; value approvals menu gone; abandoned requests not showing

"Value Approvals" shows for holders of the project_authorization claim. "Abandoned Requests" shows for holders of the admin claim.

## Rate Sheets menu access

Lines 348-358

Asked as: rate sheets missing; rate sheet menu hidden; can't open rate sheets; rate sheets link gone

"Rate Sheets" shows for holders of the job, admin or rate_sheet_revise claim. It also shows for holders of the time or project_authorization claim, only while project requests are on.

## Client and vendor creation buttons

Lines 359-384

Asked as: new client button missing; new vendor button hidden; clients link shown; vendor button not showing

The Clients and Vendors links show for everyone. The "New Client" button needs the job claim. The "New Vendor" button needs the payables_admin claim.

## Absorb menu access

Lines 385-395

Asked as: absorb actions missing; absorb proposals link; absorb menu hidden; absorb link not showing

"Absorb Actions" shows only for holders of the absorb claim. "Absorb Proposals" shows for everyone.

## PO Approvers menu link

Lines 396-401

Asked as: po approvers link missing; po approvers menu hidden; can't see po approvers; approvers menu gone

"PO Approvers" shows when canViewPOApprovers accepts the user's claims.

## Admin Profiles menu access

Lines 402-407

Asked as: admin profiles missing; profiles menu hidden; admin profiles link gone; can't open admin profiles

"Admin Profiles" shows for holders of the admin, hr, time_off_manager or it claim.

## Computers menu and its buttons

Lines 408-432

Asked as: computers menu missing; logins button hidden; retired computers not showing; raw logins missing

"Computers" shows when canViewComputers accepts the user's claims. Its buttons (Retired Computers, Logins, Raw Logins) have no separate check.

## Claims link shows for everyone; app config and UI components need admin

Lines 433-435

Asked as: app config missing; ui components menu; claims link shown; app config hidden

"Claims" shows for everyone. "App Config" and "UI Components" show only for holders of the admin claim.

## Report menu access by claim

Lines 438-483

Asked as: payroll report missing; kpi report link; reports menu hidden; payables spreadsheet missing

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs show for holders of the report claim. KPI needs the kpi claim.

## Branch Hours menu link

Lines 476-481

Asked as: branch hours link missing; branch hours menu hidden; can't see branch hours; branch hours not showing

"Branch Hours" shows when canViewBranchHours accepts the user's claims and branch manager status.

## Settings menu needs the admin claim

Lines 485-504

Asked as: currencies menu missing; time types settings hidden; machine secrets not showing; settings menu gone

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only for holders of the admin claim.
