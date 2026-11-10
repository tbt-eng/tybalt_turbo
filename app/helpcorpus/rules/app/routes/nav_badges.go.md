<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/nav_badges.go","sha256":"fd17e5b9b14910c834d108e4ad3a5407e8534b1165ebe9a416edf2cb6beae84f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/nav_badges.go

## Each nav badge count belongs to a specific page link

Lines 12-19

Asked as: badge count wrong link; number on menu item; count next to menu; badge numbers navigation

Turbo returns badge counts keyed by page link: /time/sheets/pending, /expenses/pending, /pos/pending, /reports/expense/queue, /expenses/settlement and /expenses/uncommit-requests. Each count shows on the menu item for that page.

## Nav badge counts require a signed-in user

Lines 42-44

Asked as: badge not loading signed out; unauthorized nav badge; nav counts not showing login; badge request refused

Turbo refuses the navigation badge request when nobody is signed in. The message is "unauthorized".

## Pending time sheets badge counts sheets awaiting your approval

Lines 52-66

Asked as: timesheets waiting for my approval; timesheet badge number; sheets to approve count; why timesheet count wrong

The pending time sheets badge counts time sheets where you are the approver (approver) and that are not yet approved (approved is empty). Each sheet counts once. A sheet with no time entries does not count, because the count starts from time entries.

## Pending expenses badge uses the pending expenses filter

Lines 68-76

Asked as: expense badge count; pending expenses number; expenses waiting for me; expense menu count wrong

The pending expenses badge counts the expenses that match Turbo's pending filter (wherePending) for you as the signed-in user. The conditions of that filter are defined in another file.

## Pending purchase orders badge counts POs you can act on now

Lines 78-86

Asked as: PO badge number; purchase orders to approve count; POs I can act on; PO menu count

The pending purchase orders badge counts the unapproved purchase orders that the PO visibility query marks as actionable now for you (is_unapproved_actionable_now = 1). It uses the "all" visibility view with no other filters.

## Expense commit queue badge is for committers only

Lines 88-92

Asked as: commit queue badge missing; no expense count for committers; why I see no commit count; committer menu badge

Only a committer (the commit claim) gets the expense commit queue badge. Other people get no count for this link.

## Commit queue counts approved, unrejected, uncommitted submitted expenses

Lines 93-111

Asked as: expenses ready to commit count; commit badge number wrong; foreign currency expense not counted; uncommitted expenses count

The commit queue badge counts expenses that are submitted, approved, not rejected and not committed. It does not count an expense in a foreign currency (not CAD) paid by OnAccount or CorporateCreditCard that is not yet settled. Such an expense waits for settlement first. An expense with no currency is treated as CAD and is counted.

## Expense settlement badge is for payables admins only

Lines 113-117

Asked as: settlement badge missing; no settlement count for me; payables badge not showing; settlement menu number

Only a payables admin (the payables_admin claim) gets the expense settlement badge. Other people get no count for this link.

## Settlement queue counts unsettled foreign-currency card and account expenses

Lines 118-133

Asked as: foreign currency card expenses count; expenses needing settlement number; USD card expense badge; settlement count wrong

The settlement badge counts expenses paid by OnAccount or CorporateCreditCard in a currency other than CAD. The expense must be approved, not rejected, not committed and not settled. An expense with no currency is treated as CAD and is not counted. The count does not check the submitted flag.

## Expense uncommit requests badge is for admins only

Lines 135-149

Asked as: uncommit request count missing; uncommit badge not showing; admin uncommit requests number; pending uncommit menu count

Only an admin (the admin claim) gets the uncommit requests badge. It counts the expense uncommit requests with status 'pending'.
