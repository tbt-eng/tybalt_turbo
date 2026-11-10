<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/reject_record.go","sha256":"dfe0b2ae322b4e30e6f84272271aff37aed03dd1ddf832b17eec323e7f1f7fe9","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/reject_record.go

## Who may reject a time sheet or expense

Lines 57-73

Asked as: who can reject timesheet; not authorized to reject; reject someone else's record; approver rejects expense

The assigned approver of a record may reject it. A committer (the commit claim) may also reject it. Anyone else gets 'you are not authorized to reject this record' (rejection_unauthorized).

## Committer can reject only approved records

Lines 75-83

Asked as: committer reject unapproved; commit user cannot reject; reject before approval; only approved records rejected

A committer who is not the record's approver can reject only a record that has been approved. If the record is not approved, Turbo refuses with 'only approved records can be rejected by a commit user' (record_not_approved). The approver can reject a submitted record before approval.

## Only submitted records can be rejected

Lines 139-145

Asked as: cannot reject draft; reject not submitted; reject unsubmitted timesheet; record not submitted error

Turbo refuses to reject a record that is not submitted. The message is 'only submitted records can be rejected' (record_not_submitted).

## Committed records cannot be rejected

Lines 147-153

Asked as: reject committed time; committed timesheet cannot reject; locked time reject; can't reject committed

Turbo refuses to reject a record that has a committed time. The message is 'committed records cannot be rejected' (record_committed).

## A record cannot be rejected twice

Lines 155-161

Asked as: reject twice; already rejected; rejected record again; reject same record again

Turbo refuses to reject a record that is already rejected. The message is 'this record is already rejected' (record_already_rejected).

## Rejection reason must be at least 4 characters

Lines 163-169

Asked as: rejection reason too short; reason minimum length; reason too short error; how long rejection reason

Turbo refuses a rejection whose reason is shorter than 4 characters. The message is 'rejection reason must be at least 4 characters long'. The length is counted in bytes, not characters. This check runs after the submitted, committed and already-rejected checks.

## Rejecting records the time, reason and rejector

Lines 171-174

Asked as: what happens when rejected; who rejected record; rejection time stored; record rejected successfully

A successful rejection sets the rejected time to now. It also stores the rejection reason and the rejector (the user who rejected). It returns 'record rejected successfully'.

## Rejection requires a reason in the request

Lines 34-40

Asked as: reject without reason; no reason given; rejection body invalid; reject request fails

If the request body cannot be read as a rejection, Turbo answers 'you must provide a rejection reason' (invalid_request_body).

## Missing record returns not found

Lines 48-55

Asked as: record not found reject; 404 rejecting record; record id missing; reject nonexistent record

If the record ID does not exist in the collection, Turbo answers with a 404 and the code record_not_found.

## Rejected foreign-currency account expenses lose settlement

Lines 176-191

Asked as: rejected expense settlement cleared; reject expense settled total; currency not found expense; reject foreign currency expense

When an expense is rejected, its currency must exist; otherwise Turbo refuses with 'referenced currency not found' (invalid_currency). If the currency is not the home currency and the payment type is OnAccount or CorporateCreditCard, rejection clears the settled total (set to 0), the settler and the settled date.

## Rejection notifications for time sheets and expenses

Lines 104-131

Asked as: rejection email notification; notify when rejected; rejection alert not sent; reject notification queue failed

After a successful rejection of a time sheet or an expense, Turbo queues rejection notifications that include the reason. If queuing fails, the error is only logged and the rejection still succeeds.

## Expense rejection depends on expense editing being allowed

Lines 28-30

Asked as: cannot reject expense editing; expense editing disabled reject; reject expense blocked; expenses editing not allowed

Before rejecting, Turbo calls requireExpensesEditing for the collection and returns its error if it fails. The file does not show the conditions of that check.
