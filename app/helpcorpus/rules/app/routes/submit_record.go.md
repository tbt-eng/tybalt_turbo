<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/submit_record.go","sha256":"7ce44c00d9d8f46a055172aa3d56175dc8523d9c429cc855fa59e854a7a4b725","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/submit_record.go

## Only the owner can submit a time sheet or expense

Lines 45-57

Asked as: not authorized to submit; submit someone else's timesheet; can't submit my colleague's expense; 403 submit record; wrong owner submit

Turbo refuses to submit a record unless the caller owns it, and answers "you are not authorized to submit this record" (403). For a time sheet the owner is the person in uid. For an expense the owner is the person in creator.

## A submitted record cannot be submitted again

Lines 59-66

Asked as: submitted twice; already submitted; submit button greyed out; resubmit timesheet; sent it again

When a record is already submitted, Turbo refuses with "this record is already submitted" (400).

## Rejected records cannot be submitted

Lines 68-75

Asked as: rejected can't submit; resubmit rejected expense; rejected timesheet won't submit; fix and resubmit rejected

When a record has a rejected date, Turbo refuses to submit it and answers "rejected records cannot be submitted" (400).

## Submit checks that the linked purchase order is active

Lines 77-83

Asked as: purchase order not active; PO closed can't submit; linked PO invalid; submit fails because of PO; purchase order lookup error

Before submitting, Turbo checks the purchase order on the record with validateExpensePurchaseOrderIsActive. A failed check returns 400 with that check's own error. A failed purchase order lookup (purchase_order_lookup_error) returns 500. The exact rules for an active PO are in another file.

## Expense submit needs a valid currency

Lines 85-93

Asked as: currency not found on expense; expense currency missing; currency lookup failed submit; can't find currency to submit

When an expense is submitted, Turbo looks up its currency. If the lookup fails, Turbo refuses with "referenced currency not found" (400).

## Expense submit checks the foreign currency rate

Lines 94-97

Asked as: exchange rate invalid submit; foreign rate error; bad conversion rate expense; currency rate rejected

On submit, Turbo runs validatePositiveForeignCurrencyRate on the expense's currency. If it fails, Turbo returns that error (400). The rate rule itself is in another file.

## Foreign-currency expenses need a settled total before submit

Lines 99-106

Asked as: settled amount missing; need settled total to submit; foreign expense won't submit; enter settled amount before submit; zero settled total

For an expense with payment type "Expense" in a currency other than the home currency, the settled total (settled_total) must be above 0. Otherwise Turbo refuses with "settled total is required before submitting a foreign-currency expense" (400).

## Settled total must be within tolerance of the total

Lines 107-118

Asked as: settled total doesn't match; settled amount out of tolerance; settled total too different; amount differs from settled; settled range error

For a foreign-currency expense with payment type "Expense", Turbo checks that the settled total is within tolerance of the total for that currency. If it is not, Turbo refuses with code settled_total_out_of_range and a tolerance message (400). The tolerance value comes from another file.

## Expense without a PO is checked against a limit on submit

Lines 120-128

Asked as: expense over limit without PO; no PO expense limit; expense too large no purchase order; limit exceeded submitting expense

When an expense is submitted, Turbo runs validateExpenseNoPurchaseOrderLimit using the currency and the settled total. If it fails, Turbo returns that error (400). The limit's value is in another file.

## Submitting sets the record to submitted

Lines 131-153

Asked as: submit successful; timesheet sent to approval; record marked submitted; did submit work

A successful submit sets submitted to true. Turbo then answers "Record submitted successfully".

## Expense is auto-approved when the owner is also the approver

Lines 133-135

Asked as: expense auto approved; approve my own expense; my own manager approver; self approval expense; approved myself on submit

When the expense's uid and its approver are both the submitting user, Turbo sets approved to the current time on submit. This applies to expenses only, not time sheets.

## Expense submit is blocked when expense editing is off

Lines 26-28

Asked as: can't submit expense editing off; expenses locked; expense editing disabled; submit blocked expenses closed; expense submissions turned off

Before anything else, Turbo runs requireExpensesEditing for the collection. If it returns an error, the submit stops with that error. Its conditions are defined in another file.
