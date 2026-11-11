<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/commit_record.go","sha256":"2e5e9f5c694e85e5ae157e17d96c4b5732777f6cad0fb6c97f370552ce98eb7b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/commit_record.go

## Expense commit can be blocked by expenses editing setting

Lines 34-36

Asked as: expenses editing disabled; cannot commit expenses; commit blocked setting; editing turned off

Before committing, Turbo calls requireExpensesEditing for the collection. If that check returns an error, the commit stops and Turbo returns that error. This file does not show the conditions of the check.

## Commit is all-or-nothing

Lines 43-314

Asked as: partial commit failed; commit rolled back; nothing saved after error; commit all or nothing

Turbo runs the whole commit in one transaction. If any check or save fails, the record and its purchase order are not changed. On success Turbo answers 'Record committed successfully'.

## Missing record cannot be committed

Lines 44-51

Asked as: record not found; cannot find record to commit; error fetching record; wrong record id

If the record id is not found, Turbo refuses the commit with an error 'error fetching record' and the system error text.

## A committed record cannot be committed again

Lines 53-60

Asked as: already committed; commit twice; committed again; record locked after commit

If the record already has a committed timestamp, Turbo refuses the commit with 'this record is already committed'. This check runs before the permission check.

## Only a committer can commit a record

Lines 62-82

Asked as: cannot commit my timesheet; not allowed to commit; who can commit; missing committer permission

Only staff with the commit claim (a committer) can commit a time sheet, expense or other record. Anyone else gets 'you are not authorized to commit this record'.

## Record must be submitted before commit

Lines 84-94

Asked as: commit before submitting; not submitted yet; submit first; commit unsubmitted record

For every collection except time amendments, Turbo refuses to commit a record that is not submitted. The message is 'this record is not submitted'.

## Time amendments skip submitted, approved and rejected checks

Lines 84-87

Asked as: time amendment commit without approval; amendment skips approval; commit amendment directly; amendment no submit step

Time amendments have no submitted, approved or rejected fields. Turbo does not check these states when it commits one. A committer can commit a time amendment that has not been through those steps.

## Record must be approved before commit

Lines 96-103

Asked as: not approved yet; commit without approval; approval missing; waiting for approval

For every collection except time amendments, Turbo refuses to commit a record with no approved timestamp. The message is 'this record is not approved'.

## Rejected records cannot be committed

Lines 105-112

Asked as: rejected record commit; commit after rejection; rejected cannot commit; reject blocks commit

For every collection except time amendments, Turbo refuses to commit a record that has a rejected timestamp. The message is 'rejected records cannot be committed'.

## Expense currency must exist to commit

Lines 115-123

Asked as: currency missing; currency not found; unknown currency expense; bad currency code

When Turbo commits an expense, it looks up the expense's currency. If it is not found, Turbo refuses with 'referenced currency not found'.

## Foreign-currency Expense needs a settled total above 0

Lines 125-133

Asked as: zero settled total; settled amount missing; foreign expense zero total; settled total required

An expense with payment type Expense in a non-home currency needs a settled total (settled_total) greater than 0 to commit. A settled total of 0 or less is refused with 'settled total must be greater than 0 before commit'.

## Foreign-currency on-account and card expenses must be settled

Lines 135-143

Asked as: card expense not settled; on account not settled; foreign card expense commit; settle card expense first

An expense in a non-home currency with payment type OnAccount or CorporateCreditCard needs a settled timestamp to commit. Otherwise Turbo refuses with 'foreign-currency on-account and corporate card expenses must be settled before commit'.

## Expense no-PO limit is checked at commit

Lines 145-153

Asked as: no PO expense limit; expense over limit without PO; no purchase order cap; limit exceeded commit

When Turbo commits an expense, it runs validateExpenseNoPurchaseOrderLimit with the settled total. If that check fails, Turbo refuses the commit with the error that it returns. This file does not show the limit.

## Commit sets committer, time and committed week ending

Lines 156-171

Asked as: who committed record; commit date set; committed week ending wrong; committed timestamp

On commit, Turbo sets the committer to the person who commits and the committed time to now. If the collection has a committed week ending field (committed_week_ending), Turbo sets it to the week ending date of the commit date.

## Committed expense gets a pay period ending

Lines 178-186

Asked as: pay period ending blank; pay period wrong; expense pay period set; period ending on commit

When Turbo commits an expense, it sets the pay period ending (pay_period_ending) from the expense date and the committed week ending. It uses GenerateCommittedPayPeriodEnding for this.

## Mileage expense total is recalculated on commit

Lines 188-199

Asked as: mileage total wrong; mileage recalculated; mileage amount changed on commit; mileage rate commit

When Turbo commits an expense with payment type Mileage, it recalculates the total. It uses the expense rate record and the mileage already committed (CalculateMileageTotal) and replaces the total.

## Expense commit needs an Active purchase order

Lines 204-219

Asked as: PO not active; purchase order closed; cannot use closed PO; PO inactive expense

When an expense has a purchase order, the PO must have status Active to commit the expense. Otherwise Turbo refuses with 'purchase order is not active'.

## Committing against a one-time PO closes it

Lines 238-240

Asked as: one-time PO closed; PO closed after expense; single use PO closes; one off purchase order

When Turbo commits an expense against a One-Time purchase order, it closes the PO automatically.

## Recurring PO closes when all recurrences are used

Lines 241-251

Asked as: recurring PO used up; recurring purchase order ended; recurrences exhausted; PO still open recurring

When Turbo commits an expense against a Recurring PO, it counts this commit as 1 more. It closes the PO if RecurringPurchaseOrderExhausted says it is exhausted. Otherwise the PO stays Active. Other uncommitted expenses do not count.

## Cumulative PO limit and automatic close

Lines 252-277

Asked as: cumulative PO over total; PO surplus exceeded; purchase order overspent; cumulative PO closed

For a Cumulative PO, Turbo adds the committed expenses and this expense's total, in whole cents. If the sum is above the PO total plus the allowed surplus, Turbo refuses with 'the committed expenses total exceeds the total value of the purchase order beyond the allowed surplus'. If the sum is equal to or above the PO total, within the limit, Turbo closes the PO. Below the PO total, the PO stays open.

## Automatic PO close fields

Lines 278-287

Asked as: PO closed by system; who closed purchase order; closed by system flag; PO closed automatically

When a commit closes a purchase order, Turbo sets its status to Closed and closed time to the commit time. It leaves the closer blank and sets closed_by_system to true. It saves the PO only when it closed it.
