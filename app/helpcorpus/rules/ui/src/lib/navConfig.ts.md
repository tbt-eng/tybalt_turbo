<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"991a1c318d7afe89a89f19e10add3bab4a0ec7d9b8dda33c40d2248bfbcea55c","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu visibility depends on four facts about the user

Lines 9-14

Asked as: why can't I see menu; menu item missing; what decides my menu; menu links disappear

Turbo decides whether a menu link or button shows from four facts: the user's claims, whether the user is a branch manager (isBranchManager), whether the project requests feature is on (jobRequestsEnabled), and the user's attention counts (attentionCounts). A menu check looks at nothing else.

## Menu entries without a check show to everyone

Lines 16-18

Asked as: everyone sees this link; no permission needed for menu; link shows for all users; why do all users see

A menu link or button with no visibility check (canView) shows to every signed-in user. Examples are Entries, Sheets, Time Off, My Purchase Orders, My Expenses and Claims.

## Any one named claim shows a claim-gated menu entry

Lines 20-24

Asked as: only need one role; any matching role menu; one permission is enough

A menu entry gated with anyClaim shows to a user who holds at least one of the named claims. The user does not need all of them.

## Project request menu entries need the feature on

Lines 26-30

Asked as: project requests menu missing; can't find job requests; where are project requests

Project Requests and Request Approvals show only while the project requests feature is on (jobRequestsEnabled). Project Requests also needs the time, job or project_authorization claim. Request Approvals needs the project_authorization claim.

## WIP menu: Branch and Division need WIP report access

Lines 32-33

Asked as: branch WIP link missing; division WIP hidden; can't see WIP reports

My WIP shows to everyone. The Branch and Division WIP links show only when canViewWIPReports accepts the user's claims and branch manager status.

## Show All UI displays every menu entry

Lines 59-67

Asked as: show everything toggle; see all menu items; turn on all menus; display every link

With Show All UI on, Turbo shows every section, link and button whatever the user's claims. The server still checks access when the user acts.

## Hidden menu entries, empty groups and empty sections

Lines 68-82

Asked as: section disappeared from menu; menu group gone; parent link vanished; empty menu section

Without Show All UI, Turbo hides each link or button that fails its check. A link with child links disappears when none of its children is visible. A section disappears when it has no visible link.

## Time entries and sheets menu links shown to everyone

Lines 88-100

Asked as: entries link for everyone; new entry button visible; sheets menu for all

The Time Management links Entries and Sheets have no view check, so everyone sees them. The New Entry button on Entries also has no check and shows to everyone.

## Time sheet approval menu links need approver claim

Lines 101-112

Asked as: timesheet approval link missing; pending timesheets not showing; approved by me missing

Pending My Approval and Approved By Me under Time Management show only to approvers (the tapr claim). The Pending My Approval link under Purchase Orders and Expenses needs the same claim, as does Approved By Me under Expenses.

## Time tracking and work records menu access

Lines 114-125

Asked as: time tracking menu hidden; work records not showing; can't find work records

Time Tracking shows to holders of the report or admin claim. Work Records shows to holders of the report or work_record claim.

## Amendments menu link claims

Lines 126-144

Asked as: amendments link missing; new amendment button hidden; list all amendments not showing

The Amendments link, with its List All Amendments and New Amendment buttons, shows to holders of the tame or report claim. The buttons have no extra check.

## Time Off menu link shown to everyone

Lines 145-145

Asked as: where is time off; time off link; time off menu for everyone

The Time Management menu link Time Off (/time/off) has no view check. Every signed-in user sees it.

## My Purchase Orders link and New PO button open to all

Lines 151-162

Asked as: my POs link; new PO button missing; purchase orders menu for all

The Purchase Orders menu link My Purchase Orders and its New PO button have no view check. Every user sees both.

## PO Groups menu needs a PO group view check

Lines 163-176

Asked as: PO groups link missing; new PO group button; can't see PO groups

The PO Groups link shows when canViewPOGroups accepts the user's claims. The New PO Group button also needs the super_po claim.

## Rejected, Search, Stale and Expiring PO links shown to everyone

Lines 177-188

Asked as: rejected POs link; stale POs menu; expiring recurring PO link

The Purchase Orders links Rejected (Mine), Search, Stale and Expiring (Recurring) have no view check. Everyone sees them. The All Active link is commented out and does not appear.

## Pending My Approval PO link needs approver claim

Lines 178-183

Asked as: PO approvals not showing; pending POs link missing; pos pending link

The Purchase Orders link Pending My Approval (/pos/pending) shows only to holders of the approver claim (tapr). Others do not see it, unless Show All UI is on.

## My Expenses and My Reimbursements links open to all

Lines 194-206

Asked as: my expenses missing; new expense button; reimbursements link

The Expenses links My Expenses and My Reimbursements, and the New Expense button, have no view check. Every user sees them.

## Pending My Approval expenses link needs approver claim

Lines 207-212

Asked as: expense approvals link missing; pending expenses menu; expenses to approve

The Expenses link Pending My Approval (/expenses/pending) shows only to holders of the approver claim (tapr).

## Expense menu links by claim

Lines 219-242

Asked as: expenses tracking missing; settlement link hidden; commit queue not showing; uncommit requests missing

Expenses Tracking shows to holders of the report or admin claim. Settlement shows to payables staff (the payables_admin claim). Commit Queue shows to committers (the commit claim). Uncommit Requests shows to holders of the commit or admin claim.

## Job list buttons need procurement claim

Lines 252-278

Asked as: unused active jobs button missing; stale jobs button; jobs map button

The Unused Active Jobs, Stale Jobs and Latest Jobs buttons on Jobs show only to procurement (the job claim). The Jobs Map button shows to everyone.

## New Project or Proposal button visibility

Lines 279-289

Asked as: new project button missing; can't create project; new proposal button

The New Project or Proposal button shows to procurement (the job claim). It also shows to holders of the time claim while project requests are on. With the feature off, staff who have only the time claim do not see it.

## WIP menu links

Lines 292-301

Asked as: WIP menu links; work in progress menu; branch wip hidden

My WIP shows to everyone. Branch and Division WIP show only when canViewWIPReports accepts the user's claims and branch manager status.

## Project Requests menu link needs feature and claim

Lines 302-307

Asked as: job requests link missing; project requests menu gone; can't open project requests

The Project Requests link (/job-requests) shows only while the project requests feature is on (jobRequestsEnabled). It also needs the time, job or project authorization claim (project_authorization).

## Request Approvals menu link needs feature and authorization claim

Lines 308-313

Asked as: approval queue link; request approvals missing; job requests approval queue

The Request Approvals link (/job-requests?queue=approval) shows only while project requests are on. It also needs the project authorization claim (project_authorization).

## Job Create Queue menu link

Lines 314-320

Asked as: job create queue missing; job creation queue link; create queue menu

The Job Create Queue link shows when canViewJobCreateQueue accepts the user's claims, the project requests setting and the attention counts. Another file decides the result.

## Value Approvals and Abandoned Requests menu claims

Lines 321-332

Asked as: value approvals menu; abandoned requests link; can't see abandoned requests

Value Approvals shows to holders of the project_authorization claim, and does not depend on the project requests feature. Abandoned Requests shows only to admins (the admin claim).

## Rate Sheets menu link visibility

Lines 333-343

Asked as: rate sheets link missing; rate sheet menu hidden; rate sheets not showing

Rate Sheets shows to holders of the job, admin or rate_sheet_revise claim. While project requests are on, it also shows to holders of the time or project_authorization claim.

## Client and vendor add buttons

Lines 344-369

Asked as: new client button missing; add vendor button; can't add client

The New Client button shows only to procurement (the job claim). The New Vendor button shows only to payables staff (the payables_admin claim). The Clients and Vendors links show to everyone.

## Absorb, PO approvers and admin profile menu links

Lines 370-387

Asked as: absorb actions missing; PO approvers link; admin profiles menu

Absorb Actions shows to holders of the absorb claim. PO Approvers shows when canViewPOApprovers accepts the user's claims. Admin Profiles shows to holders of the admin, hr, time_off_manager or it claim.

## Claims menu link shown to everyone

Lines 388-388

Asked as: claims link missing; where are my claims; claims list menu

The Business menu link Claims (/claims/list) has no view check. Every user sees it.

## Admin-only menu entries

Lines 389-390

Asked as: app config missing; UI components menu; settings section hidden

App Config and UI Components show only to admins (the admin claim). The Settings section also needs the admin claim for Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit.

## Report menu links need report claim

Lines 396-430

Asked as: payroll report link missing; reports menu hidden; CSV download missing

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs under Reports show only to holders of the report claim. The last two are CSV download buttons.

## Branch Hours and KPI report menu access

Lines 431-437

Asked as: branch hours link; KPI report missing; KPI menu

Branch Hours shows when canViewBranchHours accepts the user's claims and branch manager status. KPI shows to holders of the kpi claim.

## Settings menu entries are admin only

Lines 440-459

Asked as: currencies link missing; time types menu; machine secrets menu; attachment audit missing

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only to admins (the admin claim).
