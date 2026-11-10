<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/approve_record.go","sha256":"0d26de978799fe8e65dc5be3acc186b2d150b9043a1a7ab1357af2c0d575f167","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/approve_record.go

## Approval works only for expenses and time sheets

Lines 24-35

The approve action is available only for expenses (expenses) and time sheets (time_sheets). Any other record type is refused with "unsupported approval collection".

## Expense approval depends on the expense editing check

Lines 25-28

Before approving an expense, Turbo runs the expense editing check (requireExpensesEditing). If that check returns an error, Turbo returns that error and does not approve the expense. This file does not show the conditions of the check.

## Time sheet approval depends on the time editing check

Lines 29-32

Before approving a time sheet, Turbo runs the time editing check (requireTimeEditing). If that check returns an error, Turbo returns that error and does not approve the time sheet. This file does not show the conditions of the check.

## Approval fails when the record is not found

Lines 43-50

If no record matches the ID, Turbo refuses the approval. The error code is record_not_found and the message begins "error fetching record:" followed by the details.

## Only the assigned approver can approve a record

Lines 52-59

Only the person set as the approver (approver) on the record can approve it. Anyone else gets "you are not authorized to approve this record".

## Only submitted records can be approved

Lines 61-68

Turbo refuses to approve a record that is not submitted. The message is "only submitted records can be approved".

## Committed records cannot be approved

Lines 70-77

Turbo refuses to approve a record that has a committed date (committed). The message is "committed records cannot be approved".

## A record cannot be approved twice

Lines 79-86

Turbo refuses to approve a record that already has an approved date (approved). The message is "this record is already approved".

## Purchase order must be active to approve

Lines 88-94

Before approving, Turbo checks the record's purchase order with validateExpensePurchaseOrderIsActive. If the check fails, Turbo returns its error and does not approve the record. A lookup failure (purchase_order_lookup_error) is reported as a server error. Any other failure is reported as a bad request. The conditions of the check are in another file.

## Approving sets the approval timestamp

Lines 96-115

When all checks pass, Turbo sets the approved date (approved) to the current time and saves the record. It then shows "Record approved successfully". If the save fails, the error code is error_saving_record and the message begins "error saving record:".
