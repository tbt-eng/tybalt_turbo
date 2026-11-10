<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/uncommit_record.go","sha256":"08f07760b5c76cf142285625fd6acfcb97ac92e74bde6ccbd89a62c26baae430","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/uncommit_record.go

## Uncommit requires the editing feature to be available

Lines 14-22

Before anything else, uncommitting an expense runs requireExpensesEditing and uncommitting a time sheet runs requireTimeEditing. If either check fails, Turbo returns its error and does not uncommit. This file does not show what these checks test.

## Only admins can uncommit expenses or time sheets

Lines 27-51

Turbo checks for the admin claim before it uncommits an expense or a time sheet. A user without it sees "you are not authorized to uncommit this record". A request with no signed-in user is refused as "unauthorized".

## Uncommit refuses a record that does not exist

Lines 53-60

When the expense or time sheet id is not found, Turbo refuses the uncommit with an 'error fetching record' message (record_not_found).

## Uncommit refuses a record that is not committed

Lines 62-68

When an admin uncommits an expense or time sheet that has no committed date (committed), Turbo refuses with 'this record is not committed'.

## Expense uncommit records an uncommit request and notifies approval

Lines 70-88

When an admin uncommits an expense, Turbo finds or creates a direct uncommit request for it and applies it, with the admin's id. After the transaction succeeds, it queues an uncommit approval notice with any rejection reason. Time sheets do not use this path.

## Uncommit stops if the expense uncommit request cannot be recorded

Lines 71-78

When an admin uncommits an expense, Turbo first finds or creates a direct uncommit request for it. If that fails, Turbo refuses with 'error recording uncommit: ...' (error_recording_uncommit).

## Expense uncommit applies the request and passes on any refusal

Lines 80-88

For an expense, Turbo applies the uncommit request in place of the generic uncommit. If applying the request fails, Turbo returns that error. The expense is then treated as uncommitted and the approval notification is queued.

## Uncommit errors are returned with their message and code

Lines 99-107

When an uncommit step fails, Turbo returns the error message and its code (for example unauthorized, record_not_found, record_not_committed). Any other error returns only its message.

## Uncommit success message

Lines 109-113

After a successful uncommit, Turbo answers 'Record uncommitted successfully'. The expense approval notice is queued only after the transaction succeeds.

## Uncommit saves the record before reopening its purchase order

Lines 120-138

When a person uncommits a time sheet or expense, Turbo clears the commit fields and saves the record first. For an expense, Turbo then checks whether its purchase order should reopen. If the save fails, Turbo refuses with 'error saving record: ...' (error_saving_record).

## Uncommitting a time sheet clears its commit fields

Lines 121-126

Uncommitting a time sheet clears its committed date (committed), committer and committed week ending (committed_week_ending). The record stays approved and is no longer committed.

## Uncommitting an expense also clears pay period ending

Lines 121-131

Uncommitting an expense clears its committed date (committed), committer and pay period ending (pay_period_ending). The record stays approved and is no longer committed.

## Uncommit fails if the purchase order recalculation fails

Lines 140-147

After Turbo uncommits an expense, it recalculates its purchase order. If that step fails, Turbo refuses with 'error recalculating purchase order: ...' (error_recalculating_purchase_order). The check runs for expenses only, not time sheets.

## Uncommitted expense reopens only a Closed purchase order

Lines 152-165

After an expense is uncommitted, Turbo considers reopening its purchase order only when the expense has one and that PO's status is Closed. An expense with no PO, or a PO in any other status, leaves the PO unchanged. If the PO cannot be found, the uncommit fails with 'error recalculating purchase order'.

## Missing purchase order stops expense uncommit

Lines 153-161

If an uncommitted expense names a purchase order that does not exist, the recalculation fails with 'purchase order referenced by expense not found'. An expense with no purchase order skips the reopen check.

## One-Time PO reopens when no committed expenses remain

Lines 169-183

For a Closed One-Time purchase order, Turbo reopens it only when no expense on that PO is still committed. Any other committed expense keeps it Closed.

## Recurring PO reopens when no longer exhausted

Lines 184-191

For a Closed Recurring purchase order, Turbo reopens it when the PO is no longer exhausted (RecurringPurchaseOrderExhausted). The check counts only expenses that remain committed, and the uncommitted expense is not counted.

## Cumulative PO reopens when committed total is below PO total

Lines 192-199

For a Closed Cumulative purchase order, Turbo reopens it when the committed expense total is strictly less than the PO total, compared in whole cents. A total equal to the PO total keeps it Closed.

## Reopened purchase order becomes Active and loses closing data

Lines 201-210

When Turbo reopens a purchase order after an expense uncommit, it sets the status to Active. It clears the closed date (closed), the closer, and the closed-by-system flag (closed_by_system).
