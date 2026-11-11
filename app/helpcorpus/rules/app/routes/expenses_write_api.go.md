<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_write_api.go","sha256":"2d24560a4d5faea776acadc0d2dabe256a44f3330e5ceda33b2740b9afdf461a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_write_api.go

## Only listed expense fields can be written through this endpoint

Lines 15-167

Asked as: field not editable; can't change this field; unknown field in expense request; not_editable error; expense field refused

The create and update endpoints accept only these fields: uid, date, division, description, total, payment_type, attachment, job, category, kind, allowance_types, distance, cc_last_4_digits, currency, settled_total, purchase_order, vendor, source_expense, vendor_invoice_number, vendor_invoice_date and vendor_invoice_total. A request that has any other field is refused with "expense field is not editable through this endpoint". The error names the field with the code "not_editable" and the message "field is not editable through this endpoint".

## Vendor invoice details are written with the expense

Lines 34-38

Asked as: vendor invoice number; vendor invoice date; save invoice details on expense; invoice total on expense; attach invoice info to expense

The vendor invoice number, date and total (vendor_invoice_number, vendor_invoice_date, vendor_invoice_total) are accepted on the expense request. Turbo applies them to the expense's document.

## Only signed-in users can create or update expenses through the API

Lines 42-44

Asked as: unauthorized expense update; not signed in expense api; can't create expense via api; expense api login required; unauthorized create expense

Turbo refuses a create or update request for an expense when the caller is not signed in. The message is "unauthorized".

## Updating a missing expense is refused

Lines 62-70

Asked as: missing expense id; expense not found; update expense with no id; wrong expense id; can't find expense to update

When a caller updates an expense, Turbo needs an expense id in the path and an existing expense with that id. If the id is missing, Turbo says "missing expense id". If no expense has that id, Turbo says "expense not found".

## Fields such as status and approval or commit data cannot be set through this endpoint

Lines 86-88

Asked as: set expense status via api; approve expense through api; commit expense api; workflow field refused; can't set approval status

The allowed list has no approval, commit or other workflow fields. A request that sends one is refused as not editable. This check runs before Turbo processes or saves the expense.

## Expense validation and save run together in one transaction

Lines 100-120

Asked as: expense save rolled back; expense validation fails nothing saved; uploaded files deleted after failure; expense transaction failed; expense save all or nothing

Turbo runs the expense rules (hooks.ProcessExpense) and then saves the record in a single transaction. If either step fails, nothing is saved. Turbo then cleans up the document files uploaded with the request (CleanupExpenseDocumentUploads).

## File uploads become the expense attachment

Lines 133-145

Asked as: upload receipt to expense; attachment not replaced; file upload as attachment; expense attachment from upload; receipt file missing on expense

When the request is multipart/form-data and has uploaded files in "attachment", those files replace the attachment value in the request data. When no file is uploaded, the attachment value is not replaced.

## Expense save errors are returned with their own message

Lines 169-181

Asked as: failed to save expense; expense error message; expense rule rejected; why expense save failed; field errors expense save

When an expense rule fails, Turbo returns that rule's status and message. Other validation errors are returned as "failed to save expense" with the field errors.
