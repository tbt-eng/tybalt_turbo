<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/nav_badges.go","sha256":"ae7e0c740f9f8563dc9f2a2e932548550a4575d604de0e86f394e55b1d8c18da","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/nav_badges.go

## Pending purchase orders badge counts actionable unapproved POs

Lines 21-31

Asked as: pending PO badge; unapproved purchase orders count; POs waiting for my approval; purchase order badge number wrong

The pending purchase orders badge (/pos/pending) counts purchase orders that the PO visibility query marks as unapproved and actionable now (is_unapproved_actionable_now = 1). It uses the "all" visibility view for you. The visibility rules themselves live in another file.

## Nav badge counts require a signed-in user

Lines 42-44

Asked as: badge request unauthorized; nav counts not showing; not signed in badges; no badge numbers appear

Turbo refuses the navigation badge request when the person is not signed in. It returns the message "unauthorized" and gives no counts.

## Project workflow badges are added for every user

Lines 47-49

Asked as: project workflow badges missing; badges fail to load; whole badge menu broken; project badge error

Turbo adds the project workflow badges (addProjectWorkflowBadges) to the badge counts before the other counts. If that step fails, the whole badge request fails. This file does not show which badges it adds.

## Pending time sheets badge counts sheets awaiting your approval

Lines 52-66

Asked as: timesheets waiting on me; time sheet approval count; timesheet badge wrong number; sheets to approve badge

The pending time sheets badge (/time/sheets/pending) counts time sheets where you are the approver and the sheet is not yet approved (approved is empty). A sheet that has no time entries is not counted. Each sheet counts once, however many entries it has.

## Pending expenses badge uses the pending-expenses filter

Lines 68-76

Asked as: expense approval badge count; expenses waiting for my approval; pending expense number off; expense badge wrong

The pending expenses badge (/expenses/pending) counts expenses with the shared pending filter (wherePending) and your user ID as the approver (auth). This file does not show what that filter contains.

## Pending purchase orders badge is built for every user

Lines 78-86

Asked as: why do I see PO count; everyone sees purchase orders badge; PO badge for all users; purchase order count showing for me

Every signed-in user gets the pending purchase orders count, whatever their claims.

## Commit queue badge is shown only to committers

Lines 88-93

Asked as: commit queue badge missing; no badge on expense queue; can't see commit count; expense queue number hidden

Only a committer (the commit claim) gets the expense commit queue badge (/reports/expense/queue). Other users get no count for this link.

## Commit queue counts approved, uncommitted, submitted expenses

Lines 94-100

Asked as: expenses to commit count; approved expenses not committed; commit queue number wrong; what counts in commit queue

The commit queue badge counts expenses that are submitted, approved, not rejected and not committed.

## Commit queue skips unsettled foreign-currency card and account expenses

Lines 101-105

Asked as: foreign currency expense missing from queue; USD card expense not in commit count; unsettled account expense hidden; commit count skips my expense

The commit queue badge leaves out an expense when its currency is not CAD, its payment type is OnAccount or CorporateCreditCard, and it is not yet settled. An expense with no currency record counts as CAD, so Turbo does not leave it out. Once such an expense is settled, it is counted.

## Settlement badge is shown only to payables staff

Lines 113-129

Asked as: settlement badge missing; payables badge not showing; expenses to settle count; settlement queue number

Only payables staff (the payables_admin claim) get the expense settlement badge (/expenses/settlement). It counts expenses that match the settlement queue condition (expenseSettlementQueueCondition, defined elsewhere) and are not yet settled.

## Uncommit requests badge is shown only to admins

Lines 131-145

Asked as: uncommit request badge missing; no uncommit count; admin uncommit requests number; uncommit badge not showing

Only an admin (the admin claim) gets the uncommit requests badge (/expenses/uncommit-requests). It counts expense uncommit requests whose status is 'pending'.
