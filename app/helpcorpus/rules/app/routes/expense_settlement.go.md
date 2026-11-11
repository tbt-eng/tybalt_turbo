<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_settlement.go","sha256":"f64b7bd3f4785ace88d3faaf4e66bd0d23d922089b0a9de8a7c976f870bc7632","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_settlement.go

## Which expenses are in the settlement queue

Lines 41-53

Asked as: which expenses need settling; expense in settlement queue; foreign currency expense not paid; on account expense to settle

An expense is in the settlement queue when its payment type is On Account (OnAccount) or Corporate Credit Card (CorporateCreditCard) and its currency is not the home currency. It must also be approved, not rejected and not committed. An expense with no currency counts as home currency, so it is not in the queue.

## Only payables staff can use settlement actions

Lines 59-84

Asked as: payables access denied; can't open settlement page; missing payables admin permission; authentication is required settlement

Listing the settlement queue, settling an expense and clearing a settlement all need the payables admin claim (payables_admin). A person who is not signed in gets 'authentication is required'. A signed-in person without the claim gets 'payables_admin claim required'.

## Settlement queue lists open or settled expenses, newest first

Lines 86-148

Asked as: settlement list order; unsettled expenses list; show settled expenses; settlement queue newest first

The settlement list shows either the unsettled or the settled expenses that meet the queue condition, depending on which list is requested. Turbo sorts them by expense date, newest first, then by creation time, newest first. Each row shows the creator, vendor, PO number and currency. It also shows the total, an indicative CAD total, the rate and the settler.

## Indicative CAD total uses the currency rate

Lines 105-115

Asked as: CAD equivalent amount; converted total zero; exchange rate missing total; indicative CAD value

Turbo calculates the indicative CAD total as the expense total times the currency's rate. For CAD the multiplier is 1. If a foreign currency has no rate, the multiplier is 0.

## Age in days counts from approval date

Lines 118-118

Asked as: how long expense waiting; days since approval; expense age in queue; days outstanding settlement

The age of a queued expense (age_days) is the number of days since the approval date, counted from the date part only. It is never below 0.

## Settlement eligibility checks

Lines 150-163

Asked as: cannot settle home currency expense; why can't I settle this; expense not eligible settlement; wrong payment type settle

Turbo refuses to settle or clear an expense that is outside the queue. Home-currency expenses get 'home-currency expenses are not part of the settlement queue'. Other payment types get 'expense type is not eligible for settlement'. An expense that is not approved, or is rejected or committed, gets 'expense is not currently eligible for settlement'.

## Settling needs an unsettled expense; clearing needs a settled one

Lines 164-169

Asked as: already settled expense; settle twice error; undo settled not settled; expense not settled to clear

Settling an expense that is already settled fails with 'expense is already settled'. Clearing an expense that is not settled fails with 'expense is not settled'.

## Settled total must be greater than zero

Lines 183-185

Asked as: zero settlement amount; settled amount zero error; enter settled total; settle with no amount

When payables staff settle an expense, the settled total must be greater than 0. Otherwise Turbo refuses with 'settled total must be greater than 0'.

## Settling an expense needs an expense id

Lines 187-190

Asked as: settle without expense id; missing expense number settle; no id when settling

When payables staff settle an expense and the request has no expense id, Turbo refuses with "expense id is required".

## Settling re-checks eligibility before saving

Lines 192-199

Asked as: failed to settle expense; settle expense error changed; expense changed before settling

When payables staff settle an expense, Turbo loads it and checks it is still eligible and not yet settled (validateExpenseSettlementCandidate with requireSettled false). If the expense is not eligible or is already settled, Turbo refuses with "failed to settle expense". The check and the save run in one transaction.

## Settled total must be within tolerance of expense total

Lines 200-212

Asked as: settled amount too far off; amount outside tolerance; settlement amount mismatch; rate missing foreign currency settle

Turbo refuses to settle when the settled total is outside the tolerance for the expense total and currency. The error code is settled_total_out_of_range, with a message from the tolerance check. Turbo also refuses when the currency does not have a positive foreign-currency rate, which is checked by validatePositiveForeignCurrencyRate.

## Settling checks the no-PO expense limit

Lines 213-220

Asked as: no PO expense limit exceeded; purchase order missing limit; settle over limit without PO

Before settling, Turbo applies the limit for expenses without a purchase order (validateExpenseNoPurchaseOrderLimit) to the settled total. If that check fails, the settlement is refused with its error.

## Settling records the amount, settler and time

Lines 222-236

Asked as: record settlement payment; who settled expense; settlement time saved; mark expense as paid

A successful settlement saves the settled total, sets the settler to the signed-in payables user and sets the settled time to now. All of this happens in one transaction.

## Clearing a settlement resets it

Lines 251-271

Asked as: undo settlement; reverse settle; clear settled expense; committed expense cannot change

Payables staff can clear a settlement on an expense. Turbo sets the settled total to 0 and empties the settler and settled time. A committed expense cannot be cleared and gets 'committed expenses cannot be modified'.
