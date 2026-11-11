<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/commit_record.go","sha256":"ae3393e0ab1f2e80e5336ac2f9116f00fb0533365bd3d4ed0cc0366f275efe57","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/commit_record.go

## Commit refused when expenses editing is unavailable

Lines 34-36

Asked as: can't commit expenses; commit button not working; expenses editing locked; commit blocked expenses

Before committing, Turbo calls requireExpensesEditing for the collection. If that check fails, the commit stops with its error. This file does not show what the check requires.

## Commit refuses when the record is not found

Lines 43-51

Asked as: record not found on commit; commit says missing record; can't find my timesheet to commit; error fetching record

When a committer commits a record, Turbo looks it up by its ID. If no record exists, Turbo refuses with 'error fetching record' and the lookup error (code record_not_found).

## Commit failures roll back the whole commit

Lines 43-44

Asked as: commit failed nothing saved; changes gone after commit; half committed record; partial save on commit

Turbo commits the record and any purchase order closure in one transaction. If any check or save fails, Turbo refuses the commit and shows the error message, and none of the changes are kept.

## Already committed records cannot be committed again

Lines 53-60

Asked as: commit button greyed out; already committed twice; commit it again; why can't I recommit

Committing a time sheet, expense or time amendment that already has a committed date is refused with "this record is already committed".

## Only committers can commit a record

Lines 62-82

Asked as: I can't commit this; no commit button for me; who is allowed to commit; not allowed to commit

Only a committer (the commit claim) can commit a record. Anyone else gets "you are not authorized to commit this record".

## Record must be submitted before commit

Lines 84-94

Asked as: commit says not submitted; forgot to submit; submit before committing; unsubmitted timesheet commit

A record that is not submitted is refused with "this record is not submitted". Time amendments skip this check, because they have no submitted, approved or rejected fields.

## Record must be approved before commit

Lines 96-103

Asked as: manager hasn't approved yet; commit before approval; waiting for approval to commit; not approved commit error

A submitted record with no approved date is refused with "this record is not approved". Time amendments skip this check.

## Rejected records cannot be committed

Lines 105-112

Asked as: rejected timesheet commit; manager rejected my claim; commit a rejected expense; rejected record won't commit

A record with a rejected date is refused with "rejected records cannot be committed". Time amendments skip this check.

## Expense currency must exist

Lines 115-123

Asked as: currency missing on expense; unknown currency on commit; wrong currency expense; currency not found commit

When an expense is committed, Turbo looks up its currency. If it is not found, the commit is refused with "referenced currency not found".

## Foreign-currency Expense needs a settled total above 0

Lines 125-133

Asked as: settled total zero; foreign currency expense won't commit; settled amount is 0; settled total error commit

A foreign-currency expense with payment type Expense needs a settled total (settled_total) greater than 0. Otherwise the commit is refused with "settled total must be greater than 0 before commit". A settled total of exactly 0 is refused.

## Foreign-currency on-account and card expenses must be settled

Lines 135-143

Asked as: card expense not settled; on account expense won't commit; unsettled corporate card foreign; foreign expense needs settling

A foreign-currency expense with payment type OnAccount or CorporateCreditCard needs a settled date before commit. Otherwise Turbo refuses with "foreign-currency on-account and corporate card expenses must be settled before commit".

## Expense without a purchase order limit check at commit

Lines 145-153

Asked as: over purchase order limit; expense exceeds PO limit; PO limit blocks commit; expense too big for PO

When an expense is committed, Turbo runs validateExpenseNoPurchaseOrderLimit with the expense's currency and settled total. If it returns an error, the commit is refused with that error. This file does not show the limit.

## Commit sets committer, time and week ending

Lines 156-171

Asked as: who committed this; committed date and time; committed week ending wrong; when was it committed

A successful commit sets the committer to the person committing and the committed time to now. If the collection has a committed_week_ending field, Turbo sets it to the week ending for today's date.

## Committed expense gets a pay period ending

Lines 178-186

Asked as: pay period wrong on expense; pay period changed after commit; which pay period is expense in; pay period ending missing

When an expense is committed, Turbo sets its pay period ending (pay_period_ending) from the expense date and the committed week ending (via GenerateCommittedPayPeriodEnding).

## Mileage expense total is recalculated at commit

Lines 188-199

Asked as: mileage total changed; mileage amount recalculated; mileage wrong after commit; mileage total error commit

When a Mileage expense is committed, Turbo recalculates its total from the committed mileage in the fiscal year of the expense date, using the expense rate record. Any error from that calculation stops the commit.

## Expense's purchase order must be Active to commit

Lines 204-219

Asked as: purchase order not active; PO closed can't commit expense; PO status blocks commit; PO is closed

Committing an expense that references a purchase order is refused with "purchase order is not active" unless the PO status is Active.

## Committing against a One-Time PO closes it

Lines 237-240

Asked as: one-time PO closed automatically; PO closed after my expense; single use PO closed; PO shut after expense

When an expense is committed against a One-Time purchase order, Turbo closes the PO automatically. It sets status to Closed, closed to the commit time, closer blank and closed_by_system to true.

## Recurring PO closes when all recurrences are committed

Lines 241-251

Asked as: recurring PO closed; all recurrences committed PO; recurring order still open; PO closed after last recurrence

When an expense is committed against a Recurring PO, Turbo counts this commit as one. If every recurrence now has a committed expense, it closes the PO automatically. Otherwise the PO stays Active. Uncommitted expenses are not counted.

## Cumulative PO commit limit and automatic close

Lines 252-278

Asked as: cumulative PO exceeded; total exceeds PO value; PO running total over; PO limit reached commit

For a Cumulative PO, Turbo adds the committed expenses and this expense's total, compared in whole cents. If the sum exceeds the PO total, the commit is refused with "the committed expenses total exceeds the total value of the purchase order". No excess is allowed. If the sum equals the PO total, Turbo closes the PO automatically. If it is less, the PO is unchanged.

## Automatic PO closure is saved with the commit

Lines 279-300

Asked as: PO closed but expense failed; PO save error on commit; purchase order save failed; commit failed closing PO

When a commit closes a PO, Turbo saves the PO in the same transaction as the expense. If either save fails, the commit fails with an error.

## Automatic purchase order closure leaves the closer blank

Lines 319-327

Asked as: who closed the PO; PO closed by system; closer field blank; PO closed without a person

When committing an expense closes its purchase order automatically, Turbo sets the status to Closed and sets the closed time to the commit time. It sets closed_by_system to true and leaves the closer blank, because no person closed it.
