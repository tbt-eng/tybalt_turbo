<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/svelte-types.ts","sha256":"902a4993d148f1677f82c2465ed97f444977b12a7974bd02d76e4bdc66d6cf05","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/svelte-types.ts

## Purchase order page can carry parent PO and load error

Lines 32-36

Asked as: parent PO number; PO load failed; purchase order won't load; child PO parent currency

A purchase order page can carry the parent PO's number (parent_po_number) and currency (parent_currency). It can also carry a load error message (loadError) when the record fails to load.

## Second approver check reports one of four statuses

Lines 37-53

Asked as: second approver needed; why does PO need second approver; second approval threshold; second approver status

The check for a purchase order's second approver returns a status (SecondApproverStatus): not_required, requester_qualifies, candidates_available or required_no_candidates. It also returns a reason code and message, the evaluated amount, the occurrence count, the second approval threshold, the limit column used, and the second-stage timeout in hours.

## Purchase order details page flags whether the viewer can approve or reject

Lines 54-59

Asked as: can I approve this PO; approve or reject button missing; viewer permission to approve; PO details approval rights

The purchase order details page loads the PO, its expenses and the second approver diagnostics (which can be empty). It also carries a flag (canApproveOrReject) that says whether the viewer can approve or reject the PO.

## Currency defaults to CAD when a purchase order's currency is blank

Lines 60-73

Asked as: blank PO currency; PO currency is empty; currency code missing on PO; default currency CAD

A linked purchase order summary carries a currency code (currency_code) for the currency its amounts are in. Turbo treats that code as CAD when the purchase order's currency is blank.

## Recurring purchase order summary shows remaining occurrences and amount

Lines 70-72

Asked as: recurring PO remaining occurrences; how many PO payments left; recurring purchase order remaining amount; expected occurrences recurring

A linked purchase order summary can show the expected and remaining occurrences of a recurring PO, and the remaining amount. Each of these three values is optional.

## Creating a project from a proposal lists existing projects that reference it

Lines 78-82

Asked as: projects already made from proposal; duplicate project from proposal; proposal already has project; see projects linked to proposal

The job page can list the existing projects that already reference the proposal when a person creates a project from a proposal.

## Client details page has separate tabs and page counters

Lines 86-102

Asked as: client page tab page number; client projects tab paging; client owner tab count; client notes pages

The client details page keeps a separate page number for the projects tab (projectsPage), the proposals tab (proposalsPage) and the owner tab (ownerPage). It also keeps a count of jobs for each of these three tabs (counts). Notes have their own page and total pages (page, totalPages).

## Expense list pages are paged

Lines 107-113

Asked as: expense list page count; how many pages of expenses; expense list paging; expenses total pages

Expense list pages return rows with the total number of pages (totalPages) and the page size (limit). Both values are optional.

## Currency list marks currencies used by purchase orders or expenses

Lines 115-118

Asked as: which currencies are in use; currency used by expenses; can I delete this currency; currency used on POs

Each currency row shows whether any purchase order uses it (used_by_purchase_orders) and whether any expense uses it (used_by_expenses).

## Currency setup status counts records with a blank currency

Lines 120-126

Asked as: currency setup check; home currency missing; records with no currency; blank currency count

The currency setup status reports whether the home currency exists and is ready. It also counts the purchase orders and expenses with a blank currency (blank_purchase_orders, blank_expenses).

## Expense commit queue rows show approval, rejection and settlement details

Lines 134-162

Asked as: expense commit queue who approved; who committed this expense; expense rejected by whom; settled expense queue details

Each row in the expense commit queue shows who approved, committed, rejected and settled the expense. It also shows the total, the settled total, the currency with its rate and rate date, and the PO number. The committed date is optional.

## Expense settlement rows show indicative CAD total and age in days

Lines 164-187

Asked as: expense settlement CAD estimate; how old is this expense; credit card last four digits; settlement list payment type

Each row in the expense settlement list shows the total and the settled total. It also shows the settler, the payment type, the last 4 digits of the credit card, the currency rate and rate date, an indicative CAD total (indicative_cad_total), and the expense's age in days (age_days).

## Claim list shows how many people hold each claim

Lines 189-204

Asked as: how many people have this claim; claim holders list; who holds the claim; claim list holder count

Each claim in the claims list has a name, a description and a holder count (holder_count). The claim details list each holder with their given name, surname and admin profile.

## Assignable users for a claim show name, username and email

Lines 206-228

Asked as: who can be assigned to claim; claim assignee list; users available for claim; assign user to claim

For a claim, Turbo returns a list of users who can be assigned to it (assignable_users). Each user has a name, username, email and admin profile ID.

## Work record search uses prefix K, Q or F

Lines 230-235

Asked as: search work record by K Q F; work record prefix search; how many time entries for record; work record lookup

A work record search result has a prefix of K, Q or F. It also shows the number of time entries (entry_count) and the search text.

## Work record entries show week, hours, job and person

Lines 237-249

Asked as: which timesheet has this entry; hours per week on work record; entries for work record; job number on time entry

Each entry found for a work record shows the week ending, the hours, the job number and the job. It also shows the description, the person's name and the time sheet it belongs to (timesheet_id).

## Uncommit outcome is either uncommit and reject, or uncommit only

Lines 251-259

Asked as: uncommit only option; reverse committed expense; undo commit expense reject; uncommit note required

Uncommitting a committed expense has one of two outcomes: uncommit_and_reject or uncommit_only. Each uncommit reason has a default outcome, a note prompt, and a flag that says whether a note is required.

## Uncommit is blocked when payroll is closed or the expense is too old

Lines 261-267

Asked as: cannot uncommit expense; uncommit blocked payroll; expense too old to uncommit; why uncommit grayed out

An uncommit can be blocked for one of two reasons: payroll_closed or too_old. The block carries a message. It can also carry the pay period ending, the pay date and the limit in days.

## Uncommit request status and blank fields

Lines 269-294

Asked as: uncommit request cancelled; uncommit request blank expense; uncommit request reason missing; uncommit request declined

An expense uncommit request has the status pending, approved, declined or cancelled. Its expense field is blank once the expense is deleted. Its reason name is blank on rows written by a direct uncommit.

## Uncommit request state tells whether a request can be made

Lines 296-300

Asked as: can I request uncommit; uncommit request button disabled; pending uncommit request; why no uncommit request allowed

The uncommit request state for an expense says whether a request can be made (can_request). It can include the block that prevents a request, and the pending request.
