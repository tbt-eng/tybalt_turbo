<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_write_api.go","sha256":"f5088733e28d6d06a3878fce86bdcaff5b2cf31df999d9d5f7c21d8105a064c2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_write_api.go

## Only listed expense fields can be written

Lines 15-34

Turbo accepts these expense fields: uid, date, division, description, total, payment_type, attachment, job, category, kind, allowance_types, distance, cc_last_4_digits, currency, settled_total, purchase_order, vendor and source_expense. Any other field in the request is refused with "expense field is not editable through this endpoint". The error names the field with the code not_editable and the message "field is not editable through this endpoint".

## Saving an expense requires a signed-in user

Lines 38-40

Creating or updating an expense through this endpoint is refused when the request has no signed-in user. Turbo answers "unauthorized".

## Updating an expense requires a signed-in user

Lines 54-56

Turbo refuses an update to an expense (the update endpoint) when the user is not signed in. The response is "unauthorized".

## Updating an expense that does not exist is refused

Lines 58-66

When a person updates an expense and the id is missing, Turbo answers "missing expense id". When no expense has that id, Turbo answers "expense not found".

## Expense update without an id is refused

Lines 58-61

Turbo refuses an expense update that has no expense id. It answers "missing expense id".

## Fields outside the allowed list are refused with an error

Lines 82-84

When a create or update request for an expense includes a field that is not allowed, Turbo refuses the whole request. The message is "expense field is not editable through this endpoint". The field gets the code "not_editable" and the message "field is not editable through this endpoint". Turbo checks this before it processes or saves the expense.

## Expense save runs as one all-or-nothing transaction

Lines 96-112

Turbo processes the expense (hooks.ProcessExpense) and then saves it in one transaction. If either step fails, Turbo saves nothing. The same rules apply when staff create or update an expense.

## Failed expense save removes uploaded documents

Lines 113-116

When an expense save fails, Turbo cleans up the documents that the request uploaded (CleanupExpenseDocumentUploads).

## Successful expense save returns the saved expense

Lines 118-118

When the save succeeds, Turbo returns the saved expense record.

## Uploaded attachments replace the attachment field in multipart requests

Lines 129-141

When a request is multipart/form-data and includes uploaded files under "attachment", Turbo uses those files as the expense's attachment value. If no files are uploaded, Turbo does not set the attachment from the upload.

## Rejected field check covers every submitted field

Lines 146-163

Turbo checks every field in the request data, and it stops at the first field that is not allowed. The message names only one such field, even if several are wrong.

## Expense save errors are reported as specific or generic failures

Lines 165-177

When a processing hook refuses an expense save, Turbo returns the hook's own status and message. Validation errors return "failed to save expense" with the filtered validation details. Any other error also returns "failed to save expense".
