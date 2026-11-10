<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/copy_expense.go","sha256":"9099a4610934685806aca9e4e32f21ed5d2b980459b2a13749f74582e3eb5f7b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/copy_expense.go

## Only certain fields carry over to the copy

Lines 15-31

The copy keeps these fields from the original: uid, payment_type, allowance_types, distance, description, job, division, branch, category, kind and currency. All other fields, such as status, approver, totals, attachments, settlement and pay period, are not copied. Turbo's expense processing (hooks.ProcessExpense) recomputes the derived values for the new date.

## Copy expense to tomorrow requires expense editing to be open

Lines 50-52

Before copying an expense, Turbo calls requireExpensesEditing for the expenses collection. If that check fails, its error is returned and no copy is made. This file does not show when that check fails.

## Copying an expense that does not exist is refused

Lines 62-69

If the expense id is not found, Turbo refuses the copy with the message "expense <id> not found" (code record_not_found).

## Only the owner can copy an expense

Lines 71-77

Turbo refuses the copy when the expense does not belong to the signed-in user. The message is "you are not the owner of this expense" (code unauthorized).

## Only Allowance and Mileage expenses can be copied

Lines 79-87

Turbo copies an expense to tomorrow only when its payment type (payment_type) is Allowance or Mileage. Any other type is refused with "only Allowance and Mileage expenses can be copied to tomorrow" (code unsupported_payment_type).

## Expenses linked to a purchase order cannot be copied

Lines 89-95

Turbo refuses to copy an expense that has a purchase order. The message is "cannot copy an expense that is linked to a purchase order" (code linked_to_purchase_order). The reason is that a copy would draw down the PO a second time.

## Copy fails if the original date is invalid

Lines 97-104

If the original expense's date cannot be read as YYYY-MM-DD, the copy fails with "invalid date format on original record: <error>" (code invalid_date).

## Copy is dated one day after the original

Lines 106-110

The new expense gets the original's date plus one day, not the current date. It is a new record, and the original is not changed.

## Copy ignores anything the client sends in the request

Lines 112-124

Turbo processes the copy with a request that has no body and no content type. Attachments and source-expense inputs sent to the copy endpoint therefore have no effect on the copy.

## Copy goes through normal expense processing and validation for the new date

Lines 126-136

Turbo runs the normal expense processing (hooks.ProcessExpense) on the copy before saving. Any error from it is returned and no copy is saved. The copy and its checks use the new date, so current rates and insurance expiry apply. The copy is saved in one transaction.

## Copy error responses

Lines 138-153

A copy refusal returns its message and code. A refusal without its own status is reported as a bad request. Hook errors keep their own status. Other save errors go through the expense write error handler.

## Successful copy confirmation

Lines 155-158

After a successful copy, Turbo replies "Expense copied to tomorrow" and gives the id of the new record (new_record_id).
