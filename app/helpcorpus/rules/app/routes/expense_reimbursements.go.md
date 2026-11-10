<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_reimbursements.go","sha256":"691592a5bbb1f55d236b52a3a84b91ad6e359016d0d3072f4f9b0f2ef806e1de","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_reimbursements.go

## Pay period list shows your own expenses with job and amount details

Lines 22-36

Asked as: my reimbursed expenses for a period; pay period expense list; what did I get paid for; reimbursement details per period

The list for one pay period returns only the signed-in user's reimbursed expenses for that pay period ending. Each row has the date, pay date, payment type, description, distance, job number, job description, client name, currency code, foreign currency total and amount.

## Reimbursement summary lists only your own expenses by pay period

Lines 40-62

Asked as: my reimbursement summary; how much was reimbursed to me; expense totals by pay period; my expense payouts overview

The reimbursement summary returns the signed-in user's own reimbursed expenses (filtered by the user's id, uid), grouped by pay period ending. Each period shows its pay date, the number of expenses (expense_count) and the total amount (total). The reimbursed set comes from the embedded expense_reimbursements.sql query, which this file does not show.

## Pay period summary is newest first and rounds totals to cents

Lines 42-51

Asked as: pay period totals order; reimbursement totals rounding; latest pay period first; expense total cents

The reimbursement summary sorts pay periods by pay period ending, newest first. Each period total is the sum of its expense amounts rounded to 2 decimal places.

## Reimbursement summary shows all pay periods, not one

Lines 54-57

Asked as: summary for all pay periods; reimbursements across every period; show every pay period; no period filter reimbursements

The summary query passes an empty pay_period_ending value. Whether an empty value means 'all periods' depends on expense_reimbursements.sql, which this file does not show.

## Reimbursement summary load failure message

Lines 57-59

Asked as: reimbursement summary not loading; cannot see reimbursements total; reimbursement page error; summary fails to load

If Turbo cannot load the reimbursement summary, it responds with the error 'failed to load reimbursements'.

## Pay period reimbursement list requires a valid YYYY-MM-DD date

Lines 69-72

Asked as: wrong date for pay period; pay period date rejected; date format error reimbursements; pay period ending not accepted

When listing reimbursed expenses for one pay period, Turbo checks that the pay period ending is a valid date in YYYY-MM-DD format. If it is not, Turbo refuses with 'invalid date format (YYYY-MM-DD)'.

## Pay period expenses are ordered by date, payment type, then id

Lines 74-74

Asked as: order of reimbursed expenses; sort reimbursed expenses by date; reimbursement list order; expenses listed in what order

The list of reimbursed expenses for a pay period is sorted by expense date, then payment type, then id.

## Pay period list load failure message

Lines 80-82

Asked as: pay period list not loading; reimbursed expenses will not show; blank reimbursement list for period; cannot load period expenses

If Turbo cannot load the reimbursed expenses for a pay period, it responds with the error 'failed to load reimbursements'.
