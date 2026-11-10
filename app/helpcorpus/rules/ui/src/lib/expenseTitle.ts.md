<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/expenseTitle.ts","sha256":"1d6b1c74157ba5c0c1a4664d5de5ba07ca9611fbe5583029c278f7c851204aed","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/expenseTitle.ts

## Mileage expense title shows distance, not amount

Lines 13-15

Asked as: mileage expense title; mileage shows km not amount; title shows distance; why no total on mileage

A mileage expense (payment type "Mileage") gets the title "Mileage · N km", where N is the distance on the expense. The title does not show the total or the currency.

## Expense title uses vendor name as the payee

Lines 16-20

Asked as: expense title wrong vendor; expense name shows vendor; title shows merchant; expense title amount currency

For an expense that is not mileage, the title is the payee, then " · ", then the total formatted in the expense's currency (formatCurrencyAmount). The payee is the vendor name (vendor_name) with leading and trailing spaces removed.

## Expense title payee falls back to payment type name

Lines 16-19

Asked as: blank vendor expense title; no vendor name title; expense title shows payment type; empty merchant name on expense

When an expense has no vendor name, or the vendor name is only spaces, the title uses the display name of the payment type from the payment type options (allExpensePaymentTypeOptions). If no option matches, the title uses the raw payment type value (payment_type).
