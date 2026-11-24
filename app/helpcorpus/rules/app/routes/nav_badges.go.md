<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/nav_badges.go","sha256":"a66034d28989d00a26af41c3dd4cc44f4a8417bd9a5d4781ce55d2229482f0d0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/nav_badges.go

## Pending purchase orders badge counts unapproved actionable POs

Lines 22-33

Asked as: pending PO count; unapproved purchase orders badge; POs waiting on me; purchase order badge number

The pending purchase orders badge counts purchase orders that are unapproved and actionable by you now (is_unapproved_actionable_now = 1). It uses the shared purchase order visibility query with the "all" view for the signed-in user.

## Nav badges require sign-in

Lines 43-45

Asked as: badge not loading logged out; nav counts without sign in; unauthorized badge error; badge request not signed in

Turbo refuses a navigation badge request from a person who is not signed in. The message is "unauthorized".

## Project workflow badges are added for every signed-in user

Lines 48-50

Asked as: project badges missing; project workflow count showing; why do I see project numbers; project counts in nav

Turbo adds project workflow badge counts for every signed-in person before the other counts. Another function (addProjectWorkflowBadges) decides what they count, and this file does not show it.

## Pending time sheets badge counts sheets awaiting you

Lines 53-67

Asked as: timesheet approval count; sheets waiting for my approval; pending timesheets number; timesheet badge showing wrong count

The pending time sheets badge counts time sheets where you are the approver and that are not yet approved (approved is empty). A time sheet counts only if it has at least one time entry, and it counts once however many entries it has.

## Pending expenses badge uses the pending filter

Lines 69-77

Asked as: pending expense count; expenses waiting on me badge; expense badge number; why expense count is high

The pending expenses badge counts the expenses that match Turbo's pending-expense filter (wherePending) for you as the signed-in user. This file does not show what that filter contains.

## Pending purchase orders badge is shown to every signed-in user

Lines 79-87

Asked as: why I see PO count; purchase order number in menu; PO badge for everyone; badge showing POs without permission

Every signed-in person gets the pending purchase orders count, with no claim needed. The count comes from the same visibility rules as the purchase order list.

## Expense commit queue badge only for committers

Lines 89-93

Asked as: no commit queue badge; can't see commit count; commit badge missing; who sees commit queue number

Only a committer (the commit claim) gets the expense commit queue badge. People without the claim get no count for it.

## Commit queue counts approved, uncommitted, submitted expenses

Lines 94-101

Asked as: expenses to commit count; commit queue total; approved expenses not committed; count of expenses to post

The commit queue badge counts expenses that are submitted, approved, not committed and not rejected. It counts expenses from all staff, not only your own.

## Commit queue excludes unsettled foreign-currency card and on-account expenses

Lines 102-106

Asked as: foreign currency expense missing from commit; USD expense not in queue; on account expense hidden; corporate card expense not counted

The commit queue badge leaves out an expense when its currency is not CAD, its payment type is OnAccount or CorporateCreditCard, and it is not yet settled. An expense with no currency counts as CAD, so it stays in the queue. Once the expense is settled, it counts.

## Expense settlement badge only for payables admins

Lines 114-130

Asked as: settlement badge missing; payables admin count; expenses to settle number; no settlement count showing

Only a person with the payables_admin claim gets the expense settlement badge. It counts expenses that meet the settlement queue condition (expenseSettlementQueueCondition) and are not yet settled.

## Uncommit requests badge only for admins

Lines 132-146

Asked as: uncommit request count; uncommit badge missing; pending uncommit requests; admin uncommit number

Only an admin (the admin claim) gets the expense uncommit requests badge. It counts the requests whose status is 'pending'.

## Absorb proposals badge only for staff with absorb access

Lines 148-162

Asked as: absorb badge missing; absorb proposal count; pending absorb requests; no absorb count in nav

Only staff with absorb access (the absorb claim) get the absorb proposals badge. It counts the absorb proposals whose status is 'pending'.
