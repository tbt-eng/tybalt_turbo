<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"4e1db403a38bf53c10f2ca2d2e7aed737cddc8885c0d206404965a5aeda32f36","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu links show only to holders of the required claim

Lines 17-25

Asked as: link missing from menu; button hidden for me; can't see menu option; menu item not showing

Each menu link or button with a visibility check (canView) shows only to people who pass it. An entry with no check shows to everyone. This hides links only; the server still checks access.

## Project request links need the project requests feature

Lines 27-31

Asked as: no project requests menu; where is project requests; can't find request approvals; project request menu gone

Project Requests and Request Approvals show only while the project requests feature is on (jobRequestsEnabled). Project Requests shows to holders of the time, job or project authorization claim. Request Approvals shows only to holders of the project authorization claim.

## WIP menu: My WIP for all, Branch and Division restricted

Lines 33-34

Asked as: branch WIP not showing; division WIP missing; WIP menu options gone; no WIP report option

My WIP shows to everyone. Branch and Division WIP show only to people who pass canViewWIPReports (claims and whether the person is a branch manager; defined in another file).

## Show All UI displays every menu entry

Lines 60-68

Asked as: show all menu switch; see every menu item; turn on all menus; display everything in menu

When Show All UI is on, Turbo skips every visibility check and shows all sections, links and buttons whatever the person's claims. The server still checks access.

## Empty menu sections and childless groups are hidden

Lines 69-83

Asked as: empty menu section; menu heading with nothing; WIP group shows nothing; hidden parent menu item

Without Show All UI, Turbo hides a section that has no visible links. It also hides a link with child links (such as WIP) when none of the children is visible.

## Menu links open to everyone

Lines 85-101

Asked as: where are my entries; sheets menu missing; jobs link not visible; new expense button missing

Everyone sees Entries, Sheets, Shared with Me, Time Off, My Purchase Orders, Rejected (Mine), Search, Stale, Expiring (Recurring), My Expenses, My Reimbursements and Jobs, plus their New Entry, New PO and New Expense buttons.

## Time sheet approval links need the approver claim

Lines 102-113

Asked as: can't see time approvals; no timesheet approval link; approved by me missing; pending timesheets not showing

Time sheets: Pending My Approval and Approved By Me show only to holders of the approver claim (tapr). The same applies to Pending My Approval and Approved By Me under Expenses, and to Pending My Approval under Purchase Orders.

## Time Off link is open to all staff

Lines 114-146

Asked as: time off link missing; can't open time off; shared with me gone; time off menu not there

The Time Off link under Time Management (/time/off) has no claim check. Every signed-in user sees it, together with Shared with Me.

## Time tracking and work records menu access

Lines 115-126

Asked as: time tracking not showing; work records menu missing; no time tracking option; where is work records

Time Tracking shows to holders of the report or admin claim. Work Records shows to holders of the report or work record claim.

## Amendments link shows to amendment and report claim holders

Lines 127-145

Asked as: amendments menu missing; can't find amendments; no new amendment button; amendment link hidden

The Amendments link shows to holders of the time amendment (tame) or report claim. The list-all and New Amendment buttons sit on that link and show with it.

## My Purchase Orders link and New PO button show to everyone

Lines 152-163

Asked as: my POs not showing; new PO button missing; can't see purchase orders; my purchase orders gone

The My Purchase Orders link (/pos/list) and its New PO button have no claim check. Every signed-in user sees both.

## PO Groups link and New PO Group button

Lines 164-177

Asked as: PO groups menu missing; no new PO group button; can't create PO group; purchase group link hidden

PO Groups shows to people who pass canViewPOGroups (defined in another file). The New PO Group button shows only to holders of the super PO claim (super_po).

## Rejected, Search, Stale and Expiring PO links open to all

Lines 178-189

Asked as: rejected POs link gone; PO search menu missing; stale POs not showing; expiring recurring POs missing

The purchase order links Rejected (Mine), Search, Stale and Expiring (Recurring) have no claim check. The All Active link is commented out and does not appear.

## PO Pending My Approval link needs the approver claim

Lines 179-184

Asked as: PO approvals link missing; pending POs not showing; can't see POs to approve; PO pending list hidden

The Pending My Approval link under Purchase Orders (/pos/pending) shows only to holders of the approver claim (tapr).

## My Expenses, New Expense and My Reimbursements show to everyone

Lines 195-207

Asked as: expenses menu missing; reimbursements link gone; new expense not showing; can't find my expenses

The My Expenses link, its New Expense button and the My Reimbursements link have no claim check. Every signed-in user sees them.

## Expense Pending My Approval and Approved By Me need approver claim

Lines 208-219

Asked as: expense approvals not showing; expenses to approve missing; approved expenses link gone; can't see expense approvals

The expense links Pending My Approval and Approved By Me show only to holders of the approver claim (tapr).

## Expense menu links by claim

Lines 220-243

Asked as: expenses tracking missing; settlement menu hidden; commit queue not showing; uncommit requests missing

Expenses Tracking shows to holders of the report or admin claim. Settlement shows to payables admin (payables_admin). Commit Queue shows to committers (commit). Uncommit Requests shows to committers or admins.

## Job list shortcut buttons need the job claim

Lines 252-279

Asked as: unused jobs button missing; stale jobs not showing; latest jobs button gone; jobs map button missing

The Unused Active Jobs, Stale Jobs and Latest Jobs buttons show only to procurement (job claim). The Jobs Map button shows to everyone.

## New Project or Proposal button visibility

Lines 280-290

Asked as: can't add project; new proposal button missing; no new project button; project creation button hidden

The New Project or Proposal button shows to procurement (job claim). It also shows to holders of the time claim, but only while project requests are on.

## WIP menu children visibility

Lines 293-302

Asked as: my WIP missing; branch WIP hidden; WIP children not showing; division WIP option gone

My WIP has no check and shows to everyone. Branch and Division use the WIP report check.

## Project Requests link needs feature on and a claim

Lines 303-308

Asked as: job requests link missing; project requests menu hidden; can't open project requests; requests link not showing

The Project Requests link (/job-requests) shows only while project requests are on. The user must hold the time, job or project_authorization claim.

## Request Approvals link needs feature on and authorization claim

Lines 309-314

Asked as: approval queue link missing; request approvals not showing; can't see request approvals; approvals queue hidden

The Request Approvals link (/job-requests?queue=approval) shows only while project requests are on. The user must hold the project authorization claim (project_authorization).

## Job Create Queue visibility

Lines 315-321

Asked as: job create queue missing; create queue link hidden; can't see job queue; job queue not showing

The Job Create Queue link shows to people who pass canViewJobCreateQueue, which uses claims, the project requests setting and attention counts (defined in another file).

## Value Approvals needs project authorization claim

Lines 322-327

Asked as: value approvals missing; value approval link gone; can't see value approvals; approve values menu hidden

The Value Approvals link shows to holders of the project authorization claim. It does not depend on the project requests setting.

## Abandoned Requests link is admin only

Lines 328-333

Asked as: abandoned requests missing; can't see abandoned requests; abandoned link hidden; admin requests link gone

The Abandoned Requests link shows only to admins (admin claim). The setting for project requests does not affect it.

## Rate Sheets link visibility

Lines 334-344

Asked as: rate sheets link missing; can't find rate sheets; rate sheet menu hidden; no rate sheets option

Rate Sheets shows to holders of the job, admin or rate sheet revise (rate_sheet_revise) claim. It also shows to holders of the time or project authorization claim while project requests are on.

## New Client and New Vendor buttons

Lines 345-370

Asked as: new client button missing; new vendor button gone; clients link not showing; vendors menu hidden

The New Client button shows only to procurement (job claim). The New Vendor button shows only to payables admin. Clients and Vendors links show to everyone.

## Absorb Actions link needs absorb claim

Lines 371-376

Asked as: absorb actions missing; can't see absorb; absorb menu hidden; absorb link gone

The Absorb Actions link shows only to holders of the absorb claim.

## PO Approvers link uses canViewPOApprovers

Lines 377-382

Asked as: PO approvers link missing; can't find PO approvers; approvers menu hidden; PO approver list gone

The PO Approvers link shows to people who pass canViewPOApprovers (defined in another file).

## Admin Profiles link claims

Lines 383-388

Asked as: admin profiles missing; profiles link hidden; can't see admin profiles; admin profile menu gone

The Admin Profiles link shows to holders of the admin, hr, time off manager (time_off_manager) or it claim.

## Computers link uses canViewComputers

Lines 389-401

Asked as: computers menu missing; retired computers button gone; can't find computers; computers link hidden

The Computers link, with its Retired Computers button, shows to people who pass canViewComputers (defined in another file).

## Admin-only menu items

Lines 402-404

Asked as: app config missing; UI components link gone; claims link not showing; settings menu hidden for me

App Config and UI Components show only to admins (admin claim). The Settings section (Currencies, Time Types, Divisions, Machine Secrets, Attachment Audit) shows only to admins. The Claims link shows to everyone.

## Report menu links need the report claim

Lines 410-444

Asked as: payroll download missing; weekly report not showing; payables spreadsheet gone; active jobs CSV missing

Payroll, Weekly, Time Off (all), Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV downloads show only to holders of the report claim.

## Branch Hours and KPI report visibility

Lines 445-451

Asked as: branch hours link missing; KPI report not showing; can't open KPI; branch hours menu hidden

Branch Hours shows to people who pass canViewBranchHours (claims and whether the person is a branch manager; defined in another file). KPI shows only to holders of the kpi claim.

## Settings section is admin only

Lines 454-473

Asked as: no settings section; settings menu gone; currencies link missing; time types not showing

The Settings section links (Currencies, Time Types, Divisions, Machine Secrets, Attachment Audit) show only to admins, so others see no Settings section.
