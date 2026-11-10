<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/receipts.sql","sha256":"cab09a766dd68f25372a6bcd89e0f9c72c39bb76026421d5bc94e3a4fbf70fc3","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/receipts.sql

## Mileage, allowance and meals expenses never appear in receipts report

Lines 1-5

The receipts report does not handle Mileage, Allowance or Meals payment types separately, because they never have attachments. It only looks for committed expenses that have an attachment document.

## Receipts report names each file in the zip by type, person, date and total

Lines 12-18

Each receipt's zip file name (zip_filename) is built as payment type, then surname and given name separated by a comma, then the expense date as year_month_day, then the total, then the original attachment file name, all joined with hyphens. The month is written as a three-letter abbreviation (Jan to Dec). The year is the first 4 characters of the expense date and the day is characters 9-10.

## Receipts report lists only committed expenses with attachments for the pay period

Lines 21-27

The payroll receipts report (receipts.sql) lists an expense only when it has been committed (committed is not empty) and its attached document has a file (attachment is not empty). Expenses are chosen by matching the date column (date_column) to the given value (date_column_value) for the pay period. Uncommitted expenses and expenses without an attachment never appear.

## Receipts report still lists expenses when the staff profile is missing

Lines 22-27

The report joins the profile and the expense document with LEFT JOINs. An expense whose profile is not found stays in the results with empty surname and given name, which blanks those parts of the zip file name. The document join does not keep expenses without a document, because the attachment check then excludes them.

## Receipts report gives the source path and hash of each attachment

Lines 30-36

For each receipt the report returns the attachment's file name (filename), its source path (the expense document id, a slash, then the attachment file name) and its hash (sha256, from attachment_hash). The collection id is the fixed value 'pbc_2089657321'.
