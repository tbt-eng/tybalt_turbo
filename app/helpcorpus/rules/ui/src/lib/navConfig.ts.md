<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"7427b7e3f505dcbe0db57d96114b23751bd2cba131f4c86fc0353128b76ce6c5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu entries without a check show to everyone

Lines 18-84

Asked as: menu item visible to everyone; missing menu link for some users; why do I see this link; hidden menu items not showing

A menu link or button with no view check (canView) is shown to every signed-in user. Without Show All UI, Turbo hides the entries whose check fails, and hides a section that has no entries left.

## Show All UI shows every menu link and button

Lines 61-69

Asked as: show all menu items; see every menu button; turn on show all; menu items hidden from me

With Show All UI on, Turbo shows every menu section, link and button whatever the user's claims. The server still checks access when the user opens a page.

## Flyout menu hides when it has no visible links

Lines 81-81

Asked as: empty menu dropdown disappeared; job reports menu missing; flyout menu not showing; submenu gone for me

A menu item with sub-links (such as Job Reports or WIP) disappears when none of its sub-links is visible to the user.

## Time entries and sheets menu items show to everyone

Lines 90-102

Asked as: cannot see entries menu; no new entry button; where is sheets menu; time entries menu for everyone

In the Time Management menu, Entries (with its New Entry button) and Sheets have no claim check. Every signed-in user sees them.

## Time sheet approval lists need the approver claim

Lines 103-114

Asked as: no pending approvals in time; cannot see time approvals; approved by me missing; timesheet approvals not showing

Pending My Approval and Approved By Me under Time Management show only to approvers (the tapr claim).

## Time tracking menu for reporters and admins

Lines 116-121

Asked as: time management tracking missing; cannot find time tracking; tracking link not showing; timesheet tracking access

The Time Management Tracking link shows to holders of the report claim or the admin claim.

## Work Records menu claims

Lines 122-127

Asked as: work records link missing; cannot open work records; no work records menu; work record access

The Work Records link shows to holders of the report claim or the work_record claim.

## Amendments menu and buttons share one claim check

Lines 128-146

Asked as: amendments menu not showing; cannot see amendments; new amendment button missing; list all amendments missing

The Amendments link, with its List All Amendments and New Amendment buttons, shows to holders of the tame claim or the report claim. The buttons have no check of their own.

## Time Off menu link shows to everyone

Lines 147-147

Asked as: time off link missing; cannot find time off; where is time off menu; time off menu visibility

The Time Off link in the Time Management menu has no claim check. Every user sees it.

## My Purchase Orders and New PO button show to everyone

Lines 153-164

Asked as: my purchase orders missing; new po button not showing; cannot see my POs; po menu link hidden

The My Purchase Orders link and its New PO button have no claim check. Every user sees both.

## PO Groups link depends on canViewPOGroups

Lines 165-178

Asked as: po groups link missing; cannot see po groups; new po group button missing; no PO group access

The PO Groups link shows when canViewPOGroups (in another file) accepts the user's claims. The New PO Group button shows only to holders of the super_po claim.

## PO pending approval needs the approver claim

Lines 180-185

Asked as: po approvals not showing; pending po approval missing; cannot approve purchase orders; no PO approval list

The Purchase Orders Pending My Approval link shows only to approvers (the tapr claim).

## Stale and Expiring PO lists are for payables admins

Lines 189-197

Asked as: stale po list missing; expiring purchase orders not showing; recurring po list hidden; no stale POs menu

The Stale and Expiring (Recurring) purchase order lists show only to payables admins (the payables_admin claim). Other staff see a label on the purchase order in their own list.

## My Expenses, New Expense and My Reimbursements show to everyone

Lines 203-215

Asked as: my expenses link missing; no new expense button; cannot see reimbursements; reimbursements menu not showing

The My Expenses link, its New Expense button and the My Reimbursements link have no claim check. Every user sees them.

## Expense approval lists need the approver claim

Lines 216-227

Asked as: expense approvals missing; no pending expenses to approve; approved by me expenses gone; cannot see expense approvals

Pending My Approval and Approved By Me under Expenses show only to approvers (the tapr claim).

## Expense Tracking menu for reporters and admins

Lines 228-233

Asked as: expense tracking missing; cannot find expenses tracking; expenses tracking link hidden; expense tracking access

The Expenses Tracking link shows to holders of the report claim or the admin claim.

## Expense Settlement link is for payables admins

Lines 234-239

Asked as: expense settlement missing; cannot see settlement link; settlement menu hidden; no settlement access

The Settlement link shows only to holders of the payables_admin claim.

## Commit Queue and Uncommit Requests menu claims

Lines 240-251

Asked as: commit queue missing; cannot see uncommit requests; no commit menu; uncommit link not showing

The Commit Queue link shows only to committers (the commit claim). The Uncommit Requests link shows to committers and to admins (the admin claim).

## Jobs list link and Jobs Map button show to everyone

Lines 257-266

Asked as: jobs map button missing; jobs link not showing; cannot see jobs menu; no jobs menu for me

The Jobs link and its Jobs Map button have no claim check. Every user sees them. Only the New Project or Proposal button has its own check.

## New Project or Proposal button visibility

Lines 267-277

Asked as: new project button missing; cannot create project or proposal; no new proposal button; project request button hidden

The New Project or Proposal button shows to procurement (the job claim). It also shows to staff with the time claim, but only while project requests are on (jobRequestsEnabled). Other users do not see it.

## Job Reports links are for procurement

Lines 280-292

Asked as: job reports menu missing; no unused active jobs report; stale jobs link hidden; latest jobs report not showing

The Job Reports flyout (Unused Active Jobs, Stale Jobs, No Invoicing Profile, Latest Jobs) shows only to procurement (the job claim).

## WIP menu: My WIP for all, Branch and Division restricted

Lines 293-302

Asked as: wip menu missing; branch wip not showing; division wip hidden; cannot see my wip

My WIP shows to everyone. Branch and Division WIP show only when canViewWIPReports (in another file) accepts the user's claims and branch-manager status.

## Project Requests link needs the feature on

Lines 303-308

Asked as: project requests link missing; cannot find project requests; requests menu not showing; project requests hidden

The Project Requests link shows only while project requests are on. It then shows to holders of the time, job or project_authorization claim.

## Request Approvals link needs the feature and claim

Lines 309-314

Asked as: request approvals missing; cannot approve project requests; no approvals link; request approval menu hidden

The Request Approvals link shows only while project requests are on, and only to holders of the project_authorization claim.

## Job Create Queue link depends on canViewJobCreateQueue

Lines 315-321

Asked as: job create queue missing; cannot see create queue; job queue link hidden; no job creation queue

The Job Create Queue link shows when canViewJobCreateQueue (in another file) accepts the user's claims, the project requests setting and the attention counts.

## Value Approvals link needs project authorization claim

Lines 322-327

Asked as: value approvals missing; cannot see value approvals; value approval link hidden; no value approval access

The Value Approvals link shows only to holders of the project_authorization claim. It does not depend on the project requests setting.

## Abandoned Requests link is for admins

Lines 328-333

Asked as: abandoned requests missing; cannot see abandoned requests; abandoned link hidden; no abandoned requests menu

The Abandoned Requests link shows only to holders of the admin claim.

## Rate Sheets link visibility

Lines 334-344

Asked as: rate sheets link missing; cannot find rate sheets; rate sheet menu hidden; no rate sheet access

The Rate Sheets link shows to holders of the job, admin or rate_sheet_revise claim. While project requests are on, it also shows to holders of the time or project_authorization claim.

## New Client and New Vendor buttons

Lines 345-370

Asked as: new client button missing; new vendor button not showing; cannot add vendor; cannot add client

The New Client button shows only to procurement (the job claim). The New Vendor button shows only to payables admins (the payables_admin claim). The Clients and Vendors links show to everyone.

## Invoices link depends on canViewInvoices

Lines 371-376

Asked as: invoices link missing; cannot see invoices menu; no invoices access; invoices menu hidden

The Invoices link shows when canViewInvoices (in another file) accepts the user's claims.

## Absorb Actions need absorb claim; Proposals show to all

Lines 377-387

Asked as: absorb actions missing; cannot see absorb actions; absorb proposals link; absorb menu hidden

The Absorb Actions link shows only to holders of the absorb claim. The Absorb Proposals link shows to everyone.

## PO Approvers link depends on canViewPOApprovers

Lines 388-393

Asked as: po approvers link missing; cannot see po approvers; approvers menu hidden; no po approvers access

The PO Approvers link shows when canViewPOApprovers (in another file) accepts the user's claims.

## Admin Profiles link claims

Lines 394-399

Asked as: admin profiles link missing; cannot see admin profiles; profiles menu hidden; no admin profiles access

The Admin Profiles link shows to holders of the admin, hr, time_off_manager or it claim.

## Computers link and buttons depend on canViewComputers

Lines 400-424

Asked as: computers link missing; retired computers button missing; logins button not showing; cannot see computers menu

The Computers link, with its Retired Computers, Logins and Raw Logins buttons, shows when canViewComputers (in another file) accepts the user's claims. The buttons have no check of their own.

## Claims link shows to everyone; App Config and UI Components to admins

Lines 425-427

Asked as: app config link missing; ui components menu hidden; cannot see claims link; app config not showing

The Claims link shows to everyone. The App Config and UI Components links show only to holders of the admin claim.

## Report menu items need the report claim

Lines 432-467

Asked as: payroll report missing; cannot download payroll; reports menu items hidden; active jobs csv missing

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV download show only to holders of the report claim.

## Branch Hours link depends on canViewBranchHours

Lines 468-473

Asked as: branch hours link missing; cannot see branch hours; branch hours report hidden; no branch hours access

The Branch Hours report link shows when canViewBranchHours (in another file) accepts the user's claims and branch-manager status.

## KPI report needs the kpi claim

Lines 474-474

Asked as: kpi report missing; cannot open kpi report; kpi link hidden; no kpi access

The KPI report link shows only to holders of the kpi claim.

## Settings section is for admins

Lines 477-496

Asked as: settings menu missing; cannot see settings section; currencies link hidden; no settings for me

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit show only to holders of the admin claim. Other users do not see the Settings section.
