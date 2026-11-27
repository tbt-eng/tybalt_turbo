<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"3b05a3c5c2196662385adaa688baeb91b8c6ac36d7e0dddd9f17c6b22e3a57c9","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu section hides entries by claim unless Show All UI

Lines 17-19

Asked as: menu item missing; menu options hidden by permissions; can't see menu link; missing menu entries

Each menu entry has an optional check (canView). An entry without a check shows to everyone. Show All UI skips every check, and the server still checks access.

## anyClaim shows an entry to holders of any named claim

Lines 21-25

Asked as: any of several permissions; needs one of these roles; menu item needs either permission; only one permission required

A menu entry that uses anyClaim shows to a person who holds at least one of the named claims. A person with none of them does not see it.

## Project request menu entries need the feature switched on

Lines 27-31

Asked as: project requests menu missing; can't see project requests; requests menu gone; requests link hidden

A menu entry that uses the jobRequests wrapper shows only while the project requests feature is on (jobRequestsEnabled). With the feature off, it hides even from people who hold the right claim.

## WIP menu links

Lines 33-34

Asked as: branch WIP not showing; division WIP link missing; WIP menu links gone; WIP report access

My WIP shows to everyone. The Branch and Division WIP links show to users who pass canViewWIPReports (defined in another file), which uses their claims and whether they are a branch manager.

## Show All UI shows every menu entry

Lines 60-68

Asked as: show all toggle; see every menu item; show everything in menu; show all button

When Show All UI is on, Turbo shows every section, link and button in the menu, whatever the user's claims. The server still checks access when the user acts.

## Menu entries without a check are shown to everyone

Lines 69-83

Asked as: why does everyone see this; link visible to all users; unrestricted menu button; no permission check link

A link or button with no canView check is shown to every user. Turbo hides a section that has no visible items.

## Job Reports and WIP menus hide when they have no visible links

Lines 78-80

Asked as: empty job reports menu; job reports heading missing; WIP menu hidden; menu with no links

An item with sub-links (Job Reports, WIP) is hidden when none of its sub-links is visible to the user. Turbo hides the sub-links that the user's claims do not allow.

## Menu entries open to every user

Lines 89-101

Asked as: sheets visible to everyone; PO search for all staff; everyone sees my expenses; shared with me link

Entries, Sheets, Shared with Me, Time Off, My Purchase Orders, Rejected (Mine), PO Search, My Expenses, My Reimbursements and Jobs show to everyone.

## Time sheet approval lists need the approver claim

Lines 102-113

Asked as: time pending approval missing; approved by me timesheets; timesheet approvals not showing; tapr time approvals

Pending My Approval and Approved By Me under Time Management show only to holders of the approver claim (tapr).

## Time tracking and work records menu access

Lines 115-126

Asked as: time tracking menu missing; work records not showing; work record menu access; can't see time tracking

Time Tracking shows to holders of the report or admin claim. Work Records shows to holders of the report or work_record claim.

## Time amendments menu access

Lines 127-145

Asked as: amendments link missing; time amendment button; new amendment not showing; list all amendments

The Amendments link, with its List All Amendments and New Amendment buttons, shows to holders of the time amendment (tame) or report claim.

## Time Off menu entry shows to every user

Lines 146-146

Asked as: time off link for everyone; time off menu visible to all; time off link missing

The Time Off link under Time Management has no check, so every user sees it. The Sheets and Shared with Me links also have no check.

## My Purchase Orders menu and New PO button for everyone

Lines 152-163

Asked as: new PO button visible; my purchase orders for everyone; new PO button missing; purchase orders menu everyone

My Purchase Orders and its New PO button have no check, so every user sees them.

## PO Groups menu and New PO Group button

Lines 164-177

Asked as: PO groups menu missing; new PO group button; super PO group button; can't create PO group

PO Groups shows to users who pass canViewPOGroups (defined in another file). The New PO Group button shows only to holders of the super_po claim.

## Rejected (Mine) PO link shown to everyone

Lines 178-178

Asked as: rejected purchase orders link; rejected POs menu; rejected by me POs; rejected PO link shown

The Rejected (Mine) link under Purchase Orders has no check, so every user sees it.

## Purchase order approval list needs approver claim

Lines 179-184

Asked as: pending PO approvals missing; POs awaiting my approval; PO approver list; PO approval menu gone

Pending My Approval under Purchase Orders shows only to holders of the approver claim (tapr).

## PO Search link shown to everyone

Lines 185-187

Asked as: PO search link; search purchase orders menu; all active link missing; PO search for everyone

The Search link under Purchase Orders has no check, so every user sees it. The All Active link is commented out and does not show.

## Stale and expiring PO lists are for payables admins

Lines 188-196

Asked as: stale purchase orders list; expiring recurring POs; stale PO list missing; payables admin PO lists

The Stale and Expiring (Recurring) purchase order lists show only to holders of the payables_admin claim. Other users see a label on the purchase order in their own list.

## My Expenses and My Reimbursements shown to everyone

Lines 202-214

Asked as: my expenses button missing; new expense button; reimbursements menu; my reimbursements link

My Expenses, its New Expense button and My Reimbursements have no check, so every user sees them.

## Expense approval and tracking menu access

Lines 215-232

Asked as: pending expense approvals; expense tracking missing; approved expenses by me; expense approval menu

Pending My Approval and Approved By Me under Expenses show to approvers (tapr). Expense Tracking shows to holders of the report or admin claim.

## Expense settlement needs payables admin

Lines 233-238

Asked as: expense settlement link; settlement menu missing; settle expenses; settlement not showing

The Settlement link shows only to holders of the payables_admin claim.

## Commit queue and uncommit requests access

Lines 239-250

Asked as: commit queue missing; uncommit requests menu; committer menu access; can't see commit queue

Commit Queue shows only to committers (commit claim). Uncommit Requests shows to holders of the commit or admin claim.

## Menu section visibility for Time, Expenses and Business

Lines 251-252

Asked as: whole section missing; time management section gone; expenses section hidden; empty menu section

A menu section with no visible entries is hidden. The Time Management, Purchase Orders and Expenses sections always show Entries, My Purchase Orders and My Expenses, so they never hide.

## Jobs and Jobs Map shown to everyone

Lines 256-265

Asked as: jobs map button; jobs link visible to all; jobs menu for everyone; jobs section missing

The Jobs link and its Jobs Map button have no check, so every user sees them. Only the New Project or Proposal button has a check.

## New Project or Proposal button visibility

Lines 266-276

Asked as: new project button missing; new proposal button; create project button; new project or proposal

The New Project or Proposal button shows to procurement (job claim). It also shows to staff with the time claim, but only while project requests are on.

## Job Reports links are for procurement

Lines 279-291

Asked as: unused active jobs report; stale jobs report; no invoicing profile report; latest jobs link

Unused Active Jobs, Stale Jobs, No Invoicing Profile and Latest Jobs show only to procurement (job claim).

## WIP report links visibility

Lines 292-301

Asked as: my WIP link; branch WIP link missing; division WIP not showing; WIP links visibility

My WIP has no check and shows to all. Branch and Division WIP show only when canViewWIPReports allows.

## Project Requests menu needs feature on

Lines 302-307

Asked as: project requests menu gone; project requests link missing; requests menu hidden; project requests access

Project Requests shows only while project requests are on, and only to holders of the time, job or project_authorization claim.

## Request Approvals needs feature on and authorization claim

Lines 308-313

Asked as: request approvals menu; approve project requests link; authorization approvals missing; request approvals missing

Request Approvals shows only while project requests are on, and only to holders of the project_authorization claim.

## Job Create Queue visibility

Lines 314-320

Asked as: job create queue missing; create queue link; job queue menu; can't see create queue

Job Create Queue shows to users who pass canViewJobCreateQueue (defined in another file), which looks at claims, whether project requests are on, and attention counts.

## Value Approvals and Abandoned Requests access

Lines 321-332

Asked as: value approvals link; abandoned requests menu; project value approvals missing; abandoned requests link

Value Approvals shows to holders of the project_authorization claim, and it does not depend on the project requests setting. Abandoned Requests shows only to admins.

## Rate Sheets menu access

Lines 333-343

Asked as: rate sheets menu missing; rate sheet link not showing; rate sheet revise access; rate sheets for project requests

Rate Sheets shows to holders of the job, admin or rate_sheet_revise claim. It also shows to holders of the time or project_authorization claim, but only while project requests are on.

## New Client and New Vendor buttons

Lines 344-369

Asked as: new client button; new vendor button missing; add vendor button; add client button

The New Client button shows only to procurement (job claim). The New Vendor button shows only to payables_admin holders. The Clients and Vendors lists show to everyone.

## Absorb menu access

Lines 370-380

Asked as: absorb actions menu; absorb proposals link; absorb menu missing; absorb section access

Absorb Actions shows only to holders of the absorb claim. Absorb Proposals shows to everyone.

## PO Approvers menu access

Lines 381-386

Asked as: PO approvers menu; approvers list missing; who approves POs; PO approvers link

PO Approvers shows to users who pass canViewPOApprovers (defined in another file).

## Admin Profiles menu access

Lines 387-392

Asked as: admin profiles link; admin profiles missing; HR profiles menu; admin profiles access

Admin Profiles shows to holders of the admin, hr, time_off_manager or it claim.

## Computers menu and its buttons

Lines 393-417

Asked as: computers menu missing; retired computers button; logins button missing; raw logins button

Computers shows to users who pass canViewComputers (defined in another file). Its Retired Computers, Logins and Raw Logins buttons show with the link.

## Claims visible to all; App Config and UI Components admin only

Lines 418-420

Asked as: claims menu link; app config missing; UI components menu; admin menu links

The Claims link shows to everyone. App Config and UI Components show only to admins.

## Reports menu needs report claim

Lines 426-460

Asked as: payroll report missing; weekly report link; payables spreadsheet download; active jobs CSV download

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV download show only to holders of the report claim.

## Branch Hours and KPI report access

Lines 461-467

Asked as: branch hours report; KPI report missing; KPI menu access; branch hours link missing

Branch Hours shows to users who pass canViewBranchHours (claims and whether they are a branch manager). KPI shows only to holders of the kpi claim.

## Settings section is admin only

Lines 470-489

Asked as: currencies menu; time types missing; machine secrets link; attachment audit missing

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only to admins.
