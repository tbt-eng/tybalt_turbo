<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/employee_branch_hours.go","sha256":"802f5727d1c6e80fd7864e67cdb6e89410c2c0c9e557a9f777b50199723eed27","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/employee_branch_hours.go

## Report columns: employee fields, two per branch, total

Lines 21-26

The report starts with the columns payrollId, surname, givenName and defaultBranch. Each branch then adds two columns: one named for the branch, which counts hours on records with a job, and one named for the branch plus "NoJob", which counts hours on records with no job. The last column is total.

## Branch hours report needs the kpi claim

Lines 38-41

The employee branch hours report is for people who hold the kpi claim. The route group enforces this, not this file.

## Start date and end date are required

Lines 44-51

The employee branch hours report refuses a request with no start_date with the message "start_date is required". It refuses a request with no end_date with "end_date is required".

## Date error messages name the field

Lines 44-54

A missing date gives "start_date is required" or "end_date is required". An invalid date gives "<name> must be a valid date in YYYY-MM-DD format". Turbo checks start_date before end_date. A reversed range gives "start_date must be on or before end_date".

## Start date may equal end date

Lines 52-54

The employee branch hours report refuses a start date that is after the end date. The message is "start_date must be on or before end_date". When the two dates are equal, Turbo accepts them, so the range includes both dates.

## Report format defaults to JSON; only json or csv

Lines 56-62

The employee branch hours report uses json when no format is given. Turbo ignores letter case and surrounding spaces in the format. Any value other than json or csv gets the message "format must be json or csv".

## Report date range and branch hours use the dates given

Lines 64-68

The report passes the start date and end date to the query as YYYY-MM-DD text. The report date range includes both days.

## Report failure returns a generic error message

Lines 69-80

If the report cannot be built, Turbo shows "failed to create employee branch hours report". If the CSV cannot be built, Turbo shows "failed to create employee branch hours CSV".

## CSV download filename and type

Lines 77-88

With format csv, Turbo returns a file download named employee_branch_hours_<start>_<end>.csv, with the dates as YYYY-MM-DD. The content type is text/csv; charset=utf-8. The first row holds the column names.

## Dates must be real dates in YYYY-MM-DD

Lines 92-101

The employee branch hours report refuses a start_date or end_date that is not a real date in YYYY-MM-DD format. The message is "<name> must be a valid date in YYYY-MM-DD format", where <name> is start_date or end_date. A date such as 2025-02-30 is refused, and so is a date that is not zero-padded.

## Branches with a blank name get no columns

Lines 105-112

Branches whose name is empty or only spaces get no columns in the report. Turbo orders the branch columns by branch name, then by id.

## Duplicate report column names make the report fail

Lines 114-136

The report fails when a branch name gives a column name that already exists. Turbo compares names without regard to letter case. It checks against payrollId, surname, givenName, defaultBranch, total, other branch names and any branch name plus NoJob. The user sees the error "failed to create employee branch hours report".

## Report always has the same four employee columns first

Lines 114-120

The employee branch-hours report puts the four employee columns first: payrollId, surname, givenName and defaultBranch. The branch columns and the total column come after them.

## Branch columns are named by branch name and NoJob

Lines 125-152

For each branch, the report has one column named after the branch for hours on records with a job, and one named the branch name plus "NoJob" for hours on records with no job. A branch with no hours shows 0 in both columns.

## Branch hours split by job and no job

Lines 140-151

Turbo sums the hours (qt.hours) of each branch in the date range. Records whose job is not empty go to the branch column. Records with an empty job go to the branch's NoJob column. A branch with no hours shows 0.

## Total column comes last in branch hours report

Lines 153-153

The employee branch hours report always ends with a "total" column. It comes after the employee fields and the branch columns.

## Report fails if the query column count is wrong

Lines 176-182

The report fails if the query returns a different number of columns than the report builds. This is an error, not a partial report.

## Report with no matching employees has empty rows

Lines 184-187

When the query returns no rows, the employee branch hours report has its columns and an empty list of rows. It does not return null.

## Text fields blank and number fields default to 0

Lines 198-209

In the report rows, the four employee fields are text, and a missing value becomes an empty string. All other columns are numbers, and a missing value becomes 0. A value that Turbo cannot read as a number makes the report fail.

## Quotes in branch names are escaped in column names

Lines 219-221

A branch name that contains a double quote is escaped when it becomes a report column name. The report still works, and the column keeps the name with the quote.

## Numeric report values accept text; other types fail

Lines 236-251

Branch hours and total values are read as numbers, and a missing value counts as 0. A number that arrives as text is parsed as a decimal. A value that cannot be parsed, or that is of an unsupported type, makes the report fail with "failed to create employee branch hours report".

## CSV report has a header row and blank for missing values

Lines 253-275

The CSV download starts with a header row of the column names, then one row per employee. Missing values are blank, and text values are written as they are.

## CSV numbers have no trailing zeros

Lines 260-271

In the CSV, Turbo writes numbers in plain decimal form with no trailing zeros and no exponent, such as 7.5 or 8. A missing value is written as an empty cell.
