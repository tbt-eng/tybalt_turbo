<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/navConfig.ts","sha256":"b5ee438506fd1d11d1b73e41b981d9b598a4f4bcbd5fcfb0d6e4f5e3ccc01f91","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/navConfig.ts

## Menu checks can use claims, branch manager status, feature flag and attention counts

Lines 9-15

Asked as: what decides menu visibility; menu link hidden why; which checks hide links; menu check rules

A menu check can look at the person's claims, whether they are a branch manager, whether project requests are enabled, and the attention counts. Nothing else decides whether a link or button shows.

## Menu entries without a check are shown to everyone

Lines 17-19

Asked as: link visible to everyone; menu item no permission; everyone sees this link; empty menu section hidden

Without Show All UI, Turbo shows a link or button that has no canView check to every user. Turbo hides an entry whose check fails. Turbo hides a section that has no visible items.

## Menu link shows to holders of any named claim

Lines 21-25

Asked as: shown if I have one role; menu needs either permission; any one claim menu; link missing despite role

A menu entry that uses the anyClaim check is shown when the signed-in person holds at least one of the named claims. A person with none of those claims does not see it.

## Project request menus need the feature switched on

Lines 27-31

Asked as: project request menu missing; requests menu hidden; project requests feature flag; request menu off

A menu entry wrapped in the jobRequests check is shown only while the project requests feature is on (jobRequestsEnabled). It is hidden when the feature is off, even for a person who holds the claim that the entry needs.

## WIP reports menu depends on WIP report access

Lines 33-34

Asked as: WIP report menu missing; cannot see WIP reports; who can view WIP; WIP report access denied

The WIP report menu check calls canViewWIPReports with the person's claims and whether they are a branch manager. This file does not state which claims or roles pass.

## Show All UI shows every menu entry

Lines 60-68

Asked as: show everything in menu; see all menu items; show all UI setting; menu showing all links

With Show All UI on, Turbo shows every section, link and button in the menu whatever the user's claims. The server still checks access when the user acts.

## WIP menu hides when it has no visible children

Lines 69-83

Asked as: WIP vanished from menu; WIP menu gone; empty WIP section disappeared; menu section disappeared

An item with children, such as WIP, disappears from the menu when none of its children passes its check. Turbo also hides a section left with no items.

## Time entries link and New Entry button show to everyone

Lines 85-100

Asked as: New Entry button everyone; green add entry button; time entries link visible; add time entry button

The Time Management menu starts with an Entries link to the time entries list. It has a green New Entry button that opens the add-entry page. Neither has a claim check, so every signed-in user sees both.

## Time sheet approval lists need the approver claim

Lines 101-114

Asked as: pending my approval missing timesheets; approved by me not showing; timesheet approval list hidden; cannot see timesheet approvals

Staff see Pending My Approval and Approved By Me for time sheets only with the approver claim (tapr). Everyone sees Entries, Sheets, Shared with Me and Time Off.

## Time tracking and work records menu access

Lines 115-126

Asked as: time tracking menu missing; work records not showing; work records hidden; time tracking not visible

Time Tracking appears to holders of the report claim or the admin claim. Work Records appears to holders of the report claim or the work record claim (work_record).

## Amendments menu access

Lines 127-145

Asked as: amendments menu missing; list all button missing; new amendment button hidden; amendments not showing

Amendments, with its List All and New Amendment buttons, appears to holders of the time amendment claim (tame) or the report claim. The buttons share the item's check.

## Purchase order approval list needs approver claim

Lines 152-187

Asked as: PO approvals not visible; pending PO approvals missing; purchase order approval list hidden; PO approver list missing

Pending My Approval for purchase orders appears only to holders of the approver claim (tapr). My Purchase Orders, Rejected (Mine) and Search appear to everyone.

## PO Groups menu and New PO Group button

Lines 164-177

Asked as: PO groups menu missing; new PO group button missing; cannot create PO group; PO group button hidden

PO Groups appears when canViewPOGroups accepts the user's claims. The New PO Group button needs the super PO claim (super_po).

## Stale and expiring PO lists for payables staff

Lines 188-196

Asked as: stale PO list missing; expiring recurring POs hidden; payables purchase order lists; stale POs not showing

Stale and Expiring (Recurring) purchase order lists appear only to payables staff (payables_admin claim). Other staff see a label on the purchase order in their own list.

## Expense approval lists need approver claim

Lines 202-226

Asked as: expense approval list missing; pending expenses not showing; approved expenses by me hidden; expense approvals not visible

Pending My Approval and Approved By Me for expenses appear only to holders of the approver claim (tapr). My Expenses and My Reimbursements appear to everyone.

## Expense tracking, settlement and commit menus

Lines 227-250

Asked as: commit queue missing; settlement menu not showing; uncommit requests hidden; expense tracking menu missing

Expense Tracking appears to holders of the report or admin claim. Settlement appears to payables staff (payables_admin). Commit Queue appears to committers (commit claim). Uncommit Requests appears to holders of the commit or admin claim.

## Job list report buttons need the job claim

Lines 259-286

Asked as: stale jobs button missing; unused active jobs not showing; jobs map button; latest jobs button hidden

The Unused Active Jobs, Stale Jobs and Latest Jobs buttons appear only to procurement (job claim). The Jobs Map button appears to everyone.

## New Project or Proposal button visibility

Lines 287-297

Asked as: new project button missing; cannot create proposal; proposal button hidden; new project not showing

The New Project or Proposal button appears to procurement (job claim). It also appears to holders of the time claim while project requests are enabled (jobRequestsEnabled).

## WIP menu branch and division reports

Lines 300-309

Asked as: branch WIP missing; division WIP not showing; my WIP menu; branch and division WIP hidden

My WIP appears to everyone. Branch and Division WIP appear when canViewWIPReports accepts the user's claims and branch-manager status.

## Project request menus need the feature on

Lines 310-321

Asked as: request approvals missing; project requests menu hidden; project authorization menu; project requests not showing

Project Requests appears only while the project requests feature is on, to holders of the time, job or project authorization claim (project_authorization). Request Approvals appears only while the feature is on, to holders of the project authorization claim.

## Job Create Queue menu

Lines 322-328

Asked as: job create queue missing; create queue not showing; job queue menu hidden; cannot see create queue

Job Create Queue appears when canViewJobCreateQueue accepts the user's claims, the project requests setting and the attention counts.

## Value Approvals and Abandoned Requests menus

Lines 329-340

Asked as: value approvals menu missing; abandoned requests missing; abandoned requests hidden; value approvals not showing

Value Approvals appears to holders of the project authorization claim, even when project requests are off. Abandoned Requests appears only to admins (admin claim).

## Rate Sheets menu access

Lines 341-351

Asked as: rate sheets menu missing; cannot see rate sheets; rate sheet menu hidden; rate sheets not showing

Rate Sheets appears to holders of the job, admin or rate sheet revise claim (rate_sheet_revise). It also appears to holders of the time or project authorization claim while project requests are on.

## Client and vendor creation buttons

Lines 352-377

Asked as: new client button missing; add vendor button hidden; new vendor button missing; client vendor buttons not showing

The New Client button needs the job claim (procurement). The New Vendor button needs the payables claim (payables_admin). The Clients and Vendors lists appear to everyone.

## Absorb menus

Lines 378-388

Asked as: absorb actions missing; absorb proposals menu; absorb menu hidden; cannot see absorb

Absorb Actions appears only to holders of the absorb claim. Absorb Proposals appears to everyone.

## PO Approvers menu

Lines 389-394

Asked as: PO approvers menu missing; approvers menu hidden; PO approvers not showing; cannot see PO approvers

PO Approvers appears when canViewPOApprovers accepts the user's claims.

## Admin Profiles menu access

Lines 395-400

Asked as: admin profiles missing; profiles menu hidden; HR profiles menu; admin profiles not showing

Admin Profiles appears to holders of the admin, hr, time off manager (time_off_manager) or it claim.

## Computers menu and its buttons

Lines 401-425

Asked as: computers menu missing; retired computers button; logins button missing; raw logins not showing

Computers appears when canViewComputers accepts the user's claims. The Retired Computers, Logins and Raw Logins buttons have no check of their own and follow the item.

## Claims menu shown to everyone

Lines 426-428

Asked as: claims menu missing; app config not showing; UI components menu; claims menu hidden

Claims appears to everyone. App Config and UI Components appear only to admins (admin claim).

## Report menus need the report claim

Lines 433-468

Asked as: payroll report missing; weekly report menu; CSV download menu missing; report menus hidden

Payroll, Weekly, Time Off, Payables Spreadsheet, Branch Mismatches CSV and Active Jobs CSV appear only to holders of the report claim. The CSV items download a file.

## Branch Hours and KPI report menus

Lines 469-475

Asked as: branch hours menu missing; KPI report not showing; KPI menu hidden; branch hours not visible

Branch Hours appears when canViewBranchHours accepts the user's claims and branch-manager status. KPI appears only to holders of the kpi claim.

## Settings menu is for admins

Lines 478-497

Asked as: currencies menu missing; time types menu hidden; machine secrets not showing; settings menu missing

Currencies, Time Types, Divisions, Machine Secrets and Attachment Audit appear only to admins (admin claim).
