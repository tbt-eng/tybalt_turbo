<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_purchase_order_validation.go","sha256":"0517d97a51f4f19199b51982b3f0c3342189f702116434979ccca592bf755464","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_purchase_order_validation.go

## PO check is skipped for expenses without a PO

Lines 15-22

Turbo skips the active-PO check when the record is not an expense or when the expense has no purchase order.

## Expense with a purchase order needs the PO to exist

Lines 19-31

When an expense is submitted or approved and names a purchase order (purchase_order), Turbo checks that the PO exists. If it does not exist, Turbo refuses with "purchase order referenced by expense not found".

## Other PO lookup errors are refused

Lines 32-35

If Turbo fails to fetch the purchase order for a reason other than not found, it refuses with "error fetching purchase order:" followed by the error (code purchase_order_lookup_error).

## Expense with a purchase order needs an Active PO

Lines 38-43

When an expense is submitted or approved and names a purchase order, the PO's status must be "Active". Any other status makes Turbo refuse with "purchase order is not active".

## Owner mismatch between expense and PO

Lines 50-52

Turbo flags a mismatch only when both the expense owner and the PO owner are present and they differ. If either owner is empty, there is no mismatch.

## No-PO cap does not apply in some cases

Lines 58-76

The no-PO cap does not apply to home-currency expenses, to expenses that have a purchase order, or to payment types Mileage, FuelCard, PersonalReimbursement and Allowance. It also does not apply when the configured limit is 0 or less.

## Foreign-currency expense without a PO has a cap

Lines 73-84

A foreign-currency expense with no purchase order is refused when its settled total is at or above the configured no-PO limit (GetNoPOExpenseLimit). The message is "a purchase order is required for expenses of $X or more", with X shown to 2 decimals. A total equal to the limit is refused.

## Foreign currency needs a positive exchange rate

Lines 87-93

If a selected foreign currency lacks a positive exchange rate, Turbo refuses with "selected foreign currency is missing an exchange rate".
