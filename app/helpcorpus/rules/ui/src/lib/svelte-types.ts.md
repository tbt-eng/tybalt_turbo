<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/svelte-types.ts","sha256":"87050f64263becf2085261c8e5659327c951d711561039c245574332190abfa0","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/svelte-types.ts

## Purchase order page can carry a parent PO

Lines 32-36

The purchase order page data can carry the parent PO number and the parent currency. It can also carry a load error message.

## Second approver status values

Lines 37-53

The second approver check on a purchase order gives one of four statuses: not_required, requester_qualifies, candidates_available or required_no_candidates. The result includes the evaluated amount, the number of occurrences, the second approval threshold, the limit column and the second stage timeout in hours (second_stage_timeout_hours).

## Purchase order details include approve/reject permission

Lines 54-59

The purchase order details page data includes a flag (canApproveOrReject) that says whether the viewer can approve or reject the PO. It also includes the second approver diagnostics, which can be null.

## Linked purchase order summary for expenses

Lines 60-75

An expense can show a linked PO summary with the PO number, currency, type, payment type and status. For recurring POs it can also show the expected and remaining occurrences. It can also show the remaining amount.

## Job page lists projects that already use a proposal

Lines 76-80

When a person creates a project from a proposal, the job page data lists any existing projects that already reference that proposal (existingReferencingProjects).

## Client details paging

Lines 84-100

The client details page has separate page numbers for projects, proposals and owner jobs. It also shows counts for each of the three.

## Expenses list pages return paged augmented rows

Lines 105-111

Expense list pages receive augmented expense rows (items) with optional paging data: total pages (totalPages) and a page size (limit). A callback (createdItemIsVisible) decides whether a newly created expense appears in the list.

## Currency list shows usage by POs and expenses

Lines 113-116

Each currency row shows whether purchase orders use the currency (used_by_purchase_orders) and whether expenses use it (used_by_expenses).

## Currency init status counts blank records

Lines 118-124

The currency initialization status reports whether the home currency exists and is ready. It also counts the purchase orders and expenses that have a blank currency (blank_purchase_orders, blank_expenses).

## Purchase order list can come from three realtime sources

Lines 126-130

The purchase order list can take its live updates from the visible POs, the pending POs, or none (realtime_source: visible, pending, none). A callback (createdItemIsVisible) decides whether a newly created PO appears in the list.

## Expense commit queue row shows committer and settlement amounts

Lines 132-160

Each row in the expense commit queue (ExpenseCommitQueueRow) carries the staff member's name, the approver, committer, rejector and settler names, and the status stamps submitted, approved, rejected and committed. It also carries the expense's total and its settled total (settled_total), the currency code, symbol, rate and rate date, and the job and PO numbers. The committed stamp is optional, so it is absent on a row that is not yet committed.

## Expense settlement row shows indicative CAD total and age

Lines 162-185

Each row in the expense settlement list (ExpenseSettlementRow) shows the total and the settled total (settled_total), with the settler and the settled stamp. It also shows an indicative CAD total (indicative_cad_total) and the age of the expense in days (age_days). Other fields are the payment type, the last 4 digits of the credit card, the vendor, the PO number and the currency details with rate and rate date.

## Claim holders and assignable users

Lines 188-226

A claim has a name, a description and a list of holders. Each holder is a user with an admin profile. The claim list shows a holder count for each claim.

## Work record prefixes

Lines 228-233

A work record search result has a prefix of K, Q or F and an entry count.

## Work record entry rows show the time sheet entry details

Lines 235-247

Each work record entry row shows the work record, the week ending date, the staff member (name and uid), the hours, the job number and id, the description and the time sheet id it belongs to.

## Expense uncommit has two possible outcomes

Lines 249-257

An expense uncommit ends in one of two outcomes: uncommit and reject (uncommit_and_reject), or uncommit only (uncommit_only). Each uncommit reason has a default outcome (default_outcome).

## Some uncommit reasons require a note

Lines 251-257

Each expense uncommit reason has a note prompt (note_prompt) and a flag (note_required) that says whether a note is required.

## Expense uncommit can be blocked for two reasons

Lines 259-265

Uncommitting a committed expense can be blocked with a code of either payroll_closed or too_old. The block carries a message. It can also carry the pay period ending (pay_period_ending), the pay date (pay_date) and a limit in days (limit_days). The file does not give the limit value.

## Uncommit request statuses

Lines 267-292

An expense uncommit request has one of four statuses: pending, approved, declined or cancelled. It records who requested it, who decided it, when it was decided, and the decision comment.

## Uncommit request loses its expense link when the expense is deleted

Lines 269-270

The expense field of an uncommit request is blank once the expense has been deleted.

## Direct uncommit leaves the reason blank

Lines 287-288

An uncommit request row written by a direct uncommit (one without a request) has a blank reason name (reason_name).

## Uncommit request state for an expense

Lines 294-298

The uncommit request state of an expense says whether a request can be made (can_request). It can also show a block and the pending request, if there is one.
