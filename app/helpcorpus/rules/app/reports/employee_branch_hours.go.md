<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/employee_branch_hours.go","sha256":"802f5727d1c6e80fd7864e67cdb6e89410c2c0c9e557a9f777b50199723eed27","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/employee_branch_hours.go

## Report columns: employee fields, two per branch, total

Lines 21-26

Asked as: which columns in branch report; branch hours columns layout; NoJob column meaning; report column order employee

The report starts with the columns payrollId, surname, givenName and defaultBranch. Each branch then adds two columns: one named for the branch, which counts hours on records with a job, and one named for the branch plus "NoJob", which counts hours on records with no job. The last column is total.

## Branch hours report needs the kpi claim

Lines 38-41

Asked as: cannot open branch hours report; no access employee hours report; kpi permission missing; who can see branch hours

The employee branch hours report is for people who hold the kpi claim. The route group enforces this, not this file.

## Start date and end date are required

Lines 44-51

Asked as: report says start_date required; branch hours missing start date; can't run report without dates; blank date field error

The employee branch hours report refuses a request with no start_date with the message "start_date is required". It refuses a request with no end_date with "end_date is required".

## Date error messages name the field

Lines 44-54

Asked as: date error message field name; start date error on report; end date message wrong; report rejects my dates

A missing date gives "start_date is required" or "end_date is required". An invalid date gives "<name> must be a valid date in YYYY-MM-DD format". Turbo checks start_date before end_date. A reversed range gives "start_date must be on or before end_date".

## Start date may equal end date

Lines 52-54

Asked as: same start and end date; one day report range; start date after end date; dates reversed on report

The employee branch hours report refuses a start date that is after the end date. The message is "start_date must be on or before end_date". When the two dates are equal, Turbo accepts them, so the range includes both dates.

## Report format defaults to JSON; only json or csv

Lines 56-62

Asked as: report output format options; csv or json choice; format not accepted report; default report file type

The employee branch hours report uses json when no format is given. Turbo ignores letter case and surrounding spaces in the format. Any value other than json or csv gets the message "format must be json or csv".

## Report date range and branch hours use the dates given

Lines 64-68

Asked as: does report include last day; report date range inclusive; hours for chosen dates; report uses wrong dates

The report passes the start date and end date to the query as YYYY-MM-DD text. The report date range includes both days.

## Report failure returns a generic error message

Lines 69-80

Asked as: report failed generic error; branch hours report not working; csv build failed message; report keeps failing

If the report cannot be built, Turbo shows "failed to create employee branch hours report". If the CSV cannot be built, Turbo shows "failed to create employee branch hours CSV".

## CSV download filename and type

Lines 77-88

Asked as: download branch hours csv file; csv file name report; report download filename; csv header row names

With format csv, Turbo returns a file download named employee_branch_hours_<start>_<end>.csv, with the dates as YYYY-MM-DD. The content type is text/csv; charset=utf-8. The first row holds the column names.

## Dates must be real dates in YYYY-MM-DD

Lines 92-101

Asked as: invalid date 2025-02-30; date not zero padded; wrong date format report; date must be YYYY-MM-DD

The employee branch hours report refuses a start_date or end_date that is not a real date in YYYY-MM-DD format. The message is "<name> must be a valid date in YYYY-MM-DD format", where <name> is start_date or end_date. A date such as 2025-02-30 is refused, and so is a date that is not zero-padded.

## Branches with a blank name get no columns

Lines 105-112

Asked as: branch missing from report; empty branch name column; blank branch not shown; branch columns sorted order

Branches whose name is empty or only spaces get no columns in the report. Turbo orders the branch columns by branch name, then by id.

## Duplicate report column names make the report fail

Lines 114-136

Asked as: report fails duplicate branch name; branch name clashes with total; two branches same name; column name already exists

The report fails when a branch name gives a column name that already exists. Turbo compares names without regard to letter case. It checks against payrollId, surname, givenName, defaultBranch, total, other branch names and any branch name plus NoJob. The user sees the error "failed to create employee branch hours report".

## Report always has the same four employee columns first

Lines 114-120

Asked as: employee columns first order; payroll id surname columns; why employee fields lead; report starts with names

The employee branch-hours report puts the four employee columns first: payrollId, surname, givenName and defaultBranch. The branch columns and the total column come after them.

## Branch columns are named by branch name and NoJob

Lines 125-152

Asked as: branch with zero hours; columns show 0 for branch; branch name column meaning; hours column for each branch

For each branch, the report has one column named after the branch for hours on records with a job, and one named the branch name plus "NoJob" for hours on records with no job. A branch with no hours shows 0 in both columns.

## Branch hours split by job and no job

Lines 140-151

Asked as: hours with job versus without; job blank hours where; split hours by job; where do unassigned hours go

Turbo sums the hours (qt.hours) of each branch in the date range. Records whose job is not empty go to the branch column. Records with an empty job go to the branch's NoJob column. A branch with no hours shows 0.

## Total column comes last in branch hours report

Lines 153-153

Asked as: total column missing; where is total hours; total at end of report; overall hours column

The employee branch hours report always ends with a "total" column. It comes after the employee fields and the branch columns.

## Report fails if the query column count is wrong

Lines 176-182

Asked as: report column count mismatch; query columns don't match report; partial report returned; report errors unexpectedly

The report fails if the query returns a different number of columns than the report builds. This is an error, not a partial report.

## Report with no matching employees has empty rows

Lines 184-187

Asked as: no employees in report; empty branch hours results; report returns no rows; null rows in report

When the query returns no rows, the employee branch hours report has its columns and an empty list of rows. It does not return null.

## Text fields blank and number fields default to 0

Lines 198-209

Asked as: blank employee name field; missing number shows zero; text value unreadable in report; non numeric hours error

In the report rows, the four employee fields are text, and a missing value becomes an empty string. All other columns are numbers, and a missing value becomes 0. A value that Turbo cannot read as a number makes the report fail.

## Quotes in branch names are escaped in column names

Lines 219-221

Asked as: quote in branch name; double quote column name; escaped branch name report; branch name with quotes

A branch name that contains a double quote is escaped when it becomes a report column name. The report still works, and the column keeps the name with the quote.

## Numeric report values accept text; other types fail

Lines 236-251

Asked as: hours stored as text; decimal hours parsing error; unsupported value type report; number read as text

Branch hours and total values are read as numbers, and a missing value counts as 0. A number that arrives as text is parsed as a decimal. A value that cannot be parsed, or that is of an unsupported type, makes the report fail with "failed to create employee branch hours report".

## CSV report has a header row and blank for missing values

Lines 253-275

Asked as: csv empty cells; csv missing values blank; csv header and rows; text written in csv

The CSV download starts with a header row of the column names, then one row per employee. Missing values are blank, and text values are written as they are.

## CSV numbers have no trailing zeros

Lines 260-271

Asked as: csv trailing zeros decimals; csv number format 7.5; exponent in csv numbers; csv hours formatting

In the CSV, Turbo writes numbers in plain decimal form with no trailing zeros and no exponent, such as 7.5 or 8. A missing value is written as an empty cell.
