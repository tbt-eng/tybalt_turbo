<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_settlement.go","sha256":"5e2b9a0c11ad64ba33893a0b88ab2e49f29124a54eecb9b5ed7be3da5b12d3cf","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_settlement.go

## Only payables admins can use expense settlement

Lines 45-70

Every settlement action requires a payables admin (the payables_admin claim). It covers the unsettled list, the settled list, settling an expense and clearing a settlement. A signed-out user gets "authentication is required". A signed-in user without the claim gets "payables_admin claim required".

## Settlement list shows blank values for missing data

Lines 83-110

Missing names, vendor, description, PO number, settler and card digits show as empty text. A missing settled total shows as 0. A missing currency shows code CAD and symbol CAD.

## Indicative CAD total and rate in the settlement list

Lines 92-102

The indicative CAD total (indicative_cad_total) is the expense total times the currency rate. The rate is 1 for CAD or when no currency is set. A foreign currency with no rate gets 0, so its indicative total is 0.

## Age in days counts from approval date

Lines 104-104

The age (age_days) is the whole days from the approval date to now. It is never negative. An expense that is not approved has an empty approval date.

## Settlement queue lists approved foreign-currency expenses

Lines 117-121

The settlement list shows only expenses whose payment type is OnAccount or CorporateCreditCard and whose currency code is not CAD. Each expense must be approved, not rejected and not committed. An expense with no currency counts as CAD, so it never appears.

## Settled and unsettled lists are separate

Lines 122-126

The unsettled list shows expenses with no settled date. The settled list shows expenses that have a settled date. Both lists are sorted by expense date, newest first, then by creation time, newest first.

## Which expenses can be settled or cleared

Lines 140-153

Home-currency expenses are refused with "home-currency expenses are not part of the settlement queue". Payment types other than OnAccount or CorporateCreditCard are refused with "expense type is not eligible for settlement". An expense that is not approved, or is rejected or committed, is refused with "expense is not currently eligible for settlement".

## Settling an already settled expense is refused

Lines 154-159

Settling an expense that already has a settled date fails with "expense is already settled". Clearing an expense that is not settled fails with "expense is not settled".

## Settled total must be greater than zero

Lines 173-175

Settling an expense fails with "settled total must be greater than 0" when the settled total is 0 or negative. This check runs before the expense is looked up.

## Settling an expense requires an expense id

Lines 177-180

When a payables admin settles an expense without an expense id in the request path, Turbo refuses with the message "expense id is required".

## Settling re-checks eligibility inside one transaction

Lines 182-189

When a payables admin settles an expense, Turbo loads the expense and checks that it can be settled (not yet settled) in a single transaction. If the expense does not exist, or the check fails, Turbo settles nothing and refuses with "failed to settle expense".

## Settling resolves the expense's currency

Lines 190-193

Before it accepts a settlement, Turbo looks up the expense's currency (ResolveCurrencyInfo). If the lookup fails, Turbo refuses the settlement and saves nothing.

## Foreign currency rate and no-PO limit checked when settling

Lines 194-210

Settling also checks that the foreign currency has a positive rate (validatePositiveForeignCurrencyRate). It checks the no-purchase-order limit for the settled total (validateExpenseNoPurchaseOrderLimit). A failure of either check blocks the settlement. Both rules are defined in other files.

## Settled total must be within tolerance of the expense total

Lines 197-202

Settling fails with code settled_total_out_of_range when the settled total is outside the tolerance for the expense total and currency. The message comes from the settled-total tolerance message (SettledTotalToleranceMessage). The tolerance itself is set in another file.

## Settling records total, settler and time

Lines 212-215

A successful settlement saves the settled total, sets the settler to the payables admin who settled it, and sets the settled time to now. All of this runs in one transaction.

## Settlement errors are reported as a bad request

Lines 216-224

A settlement code error returns its code and message. Any other failure shows "failed to settle expense".

## Successful settle or clear returns ok

Lines 226-226

When settling an expense or clearing its settlement succeeds, Turbo answers {"ok": true}.

## Clearing a settlement requires an expense id

Lines 236-239

When a payables admin clears a settlement without an expense id in the request path, Turbo refuses with the message "expense id is required".

## Clearing a settlement resets the settlement fields

Lines 241-259

Clearing a settlement sets the settled total to 0 and empties the settler and settled time. The expense must still be settled and eligible. A failure shows "failed to clear settlement".

## Committed expenses cannot have settlement cleared

Lines 246-248

Clearing the settlement of a committed expense is refused with the reason "committed expenses cannot be modified".
