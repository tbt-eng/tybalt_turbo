<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/expenseUncommitRequests.ts","sha256":"fa64e57dbf3766f482ddba38938e4b5069713605a3ba2b6a147233c6ad907de0","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/expenseUncommitRequests.ts

## Expense uncommit outcome choices

Lines 8-11

When a committer uncommits an expense, Turbo offers two outcomes: "Uncommit and reject" (uncommit_and_reject) and "Uncommit only" (uncommit_only).

## Payment types treated as payroll expenses

Lines 13-18

Turbo treats an expense as paid on payroll when its payment type is Expense, PersonalReimbursement, Allowance, Mileage or Meals. Any other payment type is not payroll. The code notes this list must match IsPayrollExpensePaymentType in the backend.

## Uncommit request reason shown with note

Lines 24-27

An expense uncommit request shows its reason name, or "Direct uncommit" when it has no reason name. When the request has a note, Turbo shows the reason followed by a colon and the note.

## Uncommit block labels

Lines 29-31

An uncommit block with the code payroll_closed is labelled "Payroll closed". Any other block code is labelled "Too old".

## Uncommit block explanation messages

Lines 33-41

A too_old block says "This expense was committed more than N days ago.", where N is the block's limit in days (limit_days). A payroll block with a pay period ending and a pay date says payroll for that period may already have run, and gives the pay day. A payroll block without both dates says "This expense is paid on payroll and has no pay period, so it may already have been paid."

## Uncommit request status label

Lines 43-45

Turbo shows an uncommit request's status with its first letter in capitals.

## Days ago label for committed date

Lines 47-51

Turbo shows "today" when the whole days since the datetime are 0 or fewer, "1 day ago" for exactly 1 day, and "N days ago" for more. Days are whole 24-hour periods, rounded down, not calendar days.

## Days since counts whole 24-hour periods

Lines 53-56

Turbo counts days since a datetime as the elapsed time divided by 86,400,000 milliseconds, rounded down. It reads PocketBase datetimes by replacing the space with a T.
