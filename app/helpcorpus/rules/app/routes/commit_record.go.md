<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/commit_record.go","sha256":"bd8f999f08d30e33e7ee1c98fc33b192b2630b5a3b4e8c94ad21ae79353bd7ec","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/commit_record.go

## Expense commit needs expense editing enabled

Lines 34-36

Before anything else, the commit action calls requireExpensesEditing for the collection. If that check returns an error, Turbo refuses the commit and returns that error. The check itself is in another file.

## Record not found when committing

Lines 44-51

If no record has the given ID in the collection, Turbo refuses the commit. The error code is record_not_found and the message starts with 'error fetching record:'.

## Already committed records cannot be committed again

Lines 53-60

If the record already has a committed timestamp, Turbo refuses the commit with 'this record is already committed'. This check runs before the permission check.

## Only committers can commit

Lines 62-82

Only a person with the commit claim (a committer) can commit a record. Anyone else gets 'you are not authorized to commit this record'. If Turbo cannot look up the claims, it shows 'error fetching user claims: ...'.

## Record must be submitted before commit

Lines 84-94

A record that is not submitted cannot be committed. Turbo says 'this record is not submitted'. Time amendments skip this check because they have no submitted, approved or rejected fields.

## Record must be approved before commit

Lines 96-103

A record with no approved timestamp cannot be committed. Turbo says 'this record is not approved'. Time amendments skip this check.

## Rejected records cannot be committed

Lines 105-112

A record with a rejected timestamp cannot be committed. Turbo says 'rejected records cannot be committed'. Time amendments skip this check.

## Expense currency must exist

Lines 115-123

When committing an expense, Turbo looks up the expense's currency. If it is not found, Turbo refuses with 'referenced currency not found'.

## Foreign-currency expense needs a settled total

Lines 125-133

A foreign-currency expense with payment type Expense cannot be committed unless its settled total (settled_total) is greater than 0. Otherwise Turbo says 'settled total must be greater than 0 before commit'. A settled total of exactly 0 is refused.

## Foreign-currency on-account and card expenses must be settled

Lines 135-143

A foreign-currency expense with payment type OnAccount or CorporateCreditCard cannot be committed until it has a settled date. Otherwise Turbo says 'foreign-currency on-account and corporate card expenses must be settled before commit'.

## Expense without a purchase order must be within the limit

Lines 145-153

When committing an expense, Turbo runs validateExpenseNoPurchaseOrderLimit with the currency and the settled total. If that check fails, the commit is refused with its error. The limit itself is defined in another file.

## Commit sets committer and committed time

Lines 156-159

A successful commit sets the committer (committer) to the person who commits and the committed timestamp (committed) to the current time.

## Committed week ending is set from the commit date

Lines 161-171

If the collection has a committed_week_ending field, Turbo sets it to the week ending date of the day of the commit. If it cannot work out the week ending, the commit fails with 'error generating week ending: ...'.

## Expense pay period ending is set on commit

Lines 178-186

Committing an expense sets its pay period ending (pay_period_ending) from the expense date and the commit week ending. If this fails, Turbo says 'error generating pay period ending: ...'.

## Mileage expense total is recalculated on commit

Lines 188-199

When a Mileage expense is committed, Turbo replaces its total with a figure calculated from the expense rate record and the mileage already committed in the fiscal year of the expense date. Errors from the rate lookup or the calculation stop the commit.

## Expense's purchase order must be Active

Lines 204-219

Committing an expense that has a purchase order is refused with 'purchase order is not active' unless the purchase order status is Active. If the purchase order cannot be found, the commit fails with 'purchase order referenced by expense not found'.

## One-time purchase order closes on commit

Lines 237-240

Committing an expense against a One-Time purchase order closes the purchase order automatically.

## Recurring purchase order closes when exhausted

Lines 241-251

Committing an expense against a Recurring purchase order counts this expense as 1 more committed. Turbo closes the purchase order only if RecurringPurchaseOrderExhausted says it is exhausted. Other uncommitted expenses are not counted. Otherwise the purchase order stays open.

## Cumulative purchase order refuses an excess total

Lines 252-269

For a Cumulative purchase order, Turbo adds this expense's total to the committed expenses total, in whole cents. If the sum is above the purchase order total plus the allowed surplus, Turbo refuses with 'the committed expenses total exceeds the total value of the purchase order beyond the allowed surplus'. The surplus comes from the PO expense excess settings.

## Cumulative purchase order closes at or above its total

Lines 270-276

If the committed total including this expense is within the limit and equals or exceeds the purchase order total, Turbo closes the purchase order. If it is below the total, the purchase order stays open.

## Automatic closure by system leaves no closer

Lines 278-287

When a commit closes a purchase order, Turbo sets its status to Closed and its closed time to the commit time. It leaves the closer blank and sets closed_by_system to true. The purchase order is saved only if it was closed. A save failure gives 'error saving purchase orders record: ...'.

## Commit success and save failure

Lines 292-314

After all checks pass, Turbo saves the record in one transaction and says 'Record committed successfully'. A save failure gives 'error saving record: ...'.

## Automatic purchase order closure helper

Lines 318-326

closePurchaseOrderBySystem sets status Closed, sets closed to the given time, clears closer and sets closed_by_system to true.
