<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/recall_record.go","sha256":"e627e5fee032e7bec0650f180b7d83a6d5bd36a670829db06c2336d3c770ba60","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/recall_record.go

## Recall is blocked when expense editing is not allowed

Lines 25-27

Before any other check, Turbo runs requireExpensesEditing for the collection. If that check returns an error, the recall is refused with that error. This file does not show the conditions of that check.

## All recall checks and changes run in one transaction

Lines 34-102

The recall looks up the record, checks it and saves it in one database transaction. If any check or the save fails, Turbo makes no change to the record.

## Recall of a missing record is refused

Lines 35-42

Turbo refuses a recall when no record with the given id exists in the collection. The error begins "error fetching record:".

## Only the owner can recall a time sheet or expense

Lines 44-56

Turbo refuses a recall unless the person recalling owns the record. For a time sheet the owner is the person in uid. For an expense the owner is the creator (creator). Otherwise Turbo shows "you are not authorized to recall this record".

## Recall refused when record is not submitted

Lines 58-65

Turbo refuses to recall a record that is not submitted. It shows "this record is not submitted".

## Approved records cannot be recalled unless rejected

Lines 67-74

Turbo refuses to recall a record that has an approval date and no rejection date. It shows "approved records cannot be recalled unless rejected". An approved record that was also rejected can be recalled.

## Committed records cannot be recalled

Lines 76-83

Turbo refuses to recall a record that has a committed date, even if it was rejected. It shows "committed records cannot be recalled".

## Recall clears approval and rejection and unsubmits the record

Lines 85-108

A successful recall sets submitted to false. It also clears the approved date, the rejected date, the rejector (rejector) and the rejection reason (rejection_reason). The owner can then edit the record. Turbo shows "Record recalled successfully".
