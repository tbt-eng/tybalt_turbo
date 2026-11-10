<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/expensePaymentTypes.ts","sha256":"3a668d17b424c6dc3acef335a6c82e43eb317df6276dfdac0faf6ad204c38f8a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/expensePaymentTypes.ts

## Expense payment types available in Turbo

Lines 8-26

Turbo offers seven payment types for an expense: On Account, Expense, Corporate Credit Card, Allowance, Fuel Card, Mileage and Personal Reimbursement (ExpensesPaymentTypeOptions). Any other value is not a valid payment type (isExpensePaymentType).

## Personal Reimbursement hidden unless allowed or already chosen

Lines 28-42

When a person is not allowed personal reimbursement (allowPersonalReimbursement is false), the payment type list leaves out Personal Reimbursement. The list still includes it when the expense already has Personal Reimbursement as its payment type. When personal reimbursement is allowed, the list shows all seven types.

## Default payment type list follows the personal reimbursement permission

Lines 44-48

The choices for a default payment type leave out Personal Reimbursement unless the person is allowed personal reimbursement (allowPersonalReimbursement). The list has no exception for a current value, so a saved default of Personal Reimbursement does not keep the option in the list for someone who is not allowed it.
