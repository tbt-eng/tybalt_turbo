<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_settlement.go","sha256":"5e2b9a0c11ad64ba33893a0b88ab2e49f29124a54eecb9b5ed7be3da5b12d3cf","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_settlement.go

## Only payables admins can use expense settlement

Lines 45-70

Asked as: payables admin access; missing payables_admin claim; who can use settlement; settlement permission denied; not signed in settlement

Every settlement action requires a payables admin (the payables_admin claim). It covers the unsettled list, the settled list, settling an expense and clearing a settlement. A signed-out user gets "authentication is required". A signed-in user without the claim gets "payables_admin claim required".

## Settlement list shows blank values for missing data

Lines 83-110

Asked as: blank fields in settlement list; empty vendor name; missing currency shows CAD; settled total shows zero; empty PO number column

Missing names, vendor, description, PO number, settler and card digits show as empty text. A missing settled total shows as 0. A missing currency shows code CAD and symbol CAD.

## Indicative CAD total and rate in the settlement list

Lines 92-102

Asked as: indicative CAD amount; converted CAD total; exchange rate in settlement; why CAD total is zero; foreign total in dollars

The indicative CAD total (indicative_cad_total) is the expense total times the currency rate. The rate is 1 for CAD or when no currency is set. A foreign currency with no rate gets 0, so its indicative total is 0.

## Age in days counts from approval date

Lines 104-104

Asked as: how many days waiting; expense age wrong; days since approval; negative age days; unapproved expense age

The age (age_days) is the whole days from the approval date to now. It is never negative. An expense that is not approved has an empty approval date.

## Settlement queue lists approved foreign-currency expenses

Lines 117-121

Asked as: expense missing from settlement queue; foreign expense not showing; which expenses to settle; on account card expenses list; why CAD expense hidden

The settlement list shows only expenses whose payment type is OnAccount or CorporateCreditCard and whose currency code is not CAD. Each expense must be approved, not rejected and not committed. An expense with no currency counts as CAD, so it never appears.

## Settled and unsettled lists are separate

Lines 122-126

Asked as: settled expenses list order; unsettled list sorting; newest expenses first; where settled items go; two expense lists

The unsettled list shows expenses with no settled date. The settled list shows expenses that have a settled date. Both lists are sorted by expense date, newest first, then by creation time, newest first.

## Which expenses can be settled or cleared

Lines 140-153

Asked as: cannot settle this expense; home currency settlement refused; expense not eligible to settle; wrong payment type settle; rejected expense settle error

Home-currency expenses are refused with "home-currency expenses are not part of the settlement queue". Payment types other than OnAccount or CorporateCreditCard are refused with "expense type is not eligible for settlement". An expense that is not approved, or is rejected or committed, is refused with "expense is not currently eligible for settlement".

## Settling an already settled expense is refused

Lines 154-159

Asked as: settle twice error; already paid expense; clear unsettled expense; expense already settled message; undo settled but not settled

Settling an expense that already has a settled date fails with "expense is already settled". Clearing an expense that is not settled fails with "expense is not settled".

## Settled total must be greater than zero

Lines 173-175

Asked as: zero settled amount; negative settled total; settle with zero amount; settled amount must be positive; why settle amount rejected

Settling an expense fails with "settled total must be greater than 0" when the settled total is 0 or negative. This check runs before the expense is looked up.

## Settling an expense requires an expense id

Lines 177-180

Asked as: settle with no id; missing expense id settle; settle request without id; expense id missing error

When a payables admin settles an expense without an expense id in the request path, Turbo refuses with the message "expense id is required".

## Settling re-checks eligibility inside one transaction

Lines 182-189

Asked as: settle failed nothing saved; settlement transaction failed; expense not found when settling; settle expense error generic

When a payables admin settles an expense, Turbo loads the expense and checks that it can be settled (not yet settled) in a single transaction. If the expense does not exist, or the check fails, Turbo settles nothing and refuses with "failed to settle expense".

## Settling resolves the expense's currency

Lines 190-193

Asked as: currency lookup failed settling; settle with unknown currency; currency info error settlement

Before it accepts a settlement, Turbo looks up the expense's currency (ResolveCurrencyInfo). If the lookup fails, Turbo refuses the settlement and saves nothing.

## Foreign currency rate and no-PO limit checked when settling

Lines 194-210

Asked as: no PO limit exceeded; foreign currency rate missing settle; settlement blocked by limit; rate not positive settle

Settling also checks that the foreign currency has a positive rate (validatePositiveForeignCurrencyRate). It checks the no-purchase-order limit for the settled total (validateExpenseNoPurchaseOrderLimit). A failure of either check blocks the settlement. Both rules are defined in other files.

## Settled total must be within tolerance of the expense total

Lines 197-202

Asked as: settled amount too far off; total out of tolerance; settled amount mismatch; amount differs from expense total; out of range settlement

Settling fails with code settled_total_out_of_range when the settled total is outside the tolerance for the expense total and currency. The message comes from the settled-total tolerance message (SettledTotalToleranceMessage). The tolerance itself is set in another file.

## Settling records total, settler and time

Lines 212-215

Asked as: who settled expense; record settler and time; settled date saved; settlement timestamp

A successful settlement saves the settled total, sets the settler to the payables admin who settled it, and sets the settled time to now. All of this runs in one transaction.

## Settlement errors are reported as a bad request

Lines 216-224

Asked as: settle error bad request; settlement error code; failed to settle message; settle returns error

A settlement code error returns its code and message. Any other failure shows "failed to settle expense".

## Successful settle or clear returns ok

Lines 226-226

Asked as: settle success response; clear success ok; settlement confirmed; ok true response

When settling an expense or clearing its settlement succeeds, Turbo answers {"ok": true}.

## Clearing a settlement requires an expense id

Lines 236-239

Asked as: clear settlement no id; clear request missing id; expense id needed to clear

When a payables admin clears a settlement without an expense id in the request path, Turbo refuses with the message "expense id is required".

## Clearing a settlement resets the settlement fields

Lines 241-259

Asked as: undo settlement; reset settled amount; remove settler; failed to clear settlement; unsettle expense

Clearing a settlement sets the settled total to 0 and empties the settler and settled time. The expense must still be settled and eligible. A failure shows "failed to clear settlement".

## Committed expenses cannot have settlement cleared

Lines 246-248

Asked as: unsettle committed expense; clear committed settlement; committed expense cannot change; reverse settlement committed

Clearing the settlement of a committed expense is refused with the reason "committed expenses cannot be modified".
