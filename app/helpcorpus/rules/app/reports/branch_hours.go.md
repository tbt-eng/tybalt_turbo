<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/branch_hours.go","sha256":"0cf740b31f45099b6f3bd56097fe1e3c46dad2cb227ea2ed1912e0bbace62e0f","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/branch_hours.go

## Branch hours CSV columns and hours format

Lines 39-49

The CSV has a header row, then one line per report row. The columns are direction, jobNumber, jobDescription, client, workBranch, surname, givenName, staffBranch and hours. The CSV leaves out the job id and the staff uid. Hours are written without trailing zeros.

## Branch hours report requires a branch

Lines 60-63

The branch hours report refuses a request with no branch, or with a branch that is only spaces. Turbo shows the message "branch is required".

## Branch hours report needs valid start and end dates

Lines 64-71

The report reads start_date and end_date from the request and refuses a request when either date is invalid. The date check itself is in parseEmployeeBranchHoursDate, which is in another file, so its exact rules and messages are not shown here.

## Start date may equal end date in branch hours report

Lines 72-74

Turbo refuses the report when the start date is after the end date. It shows "start_date must be on or before end_date". A start date equal to the end date is accepted, so a one-day report works.

## Branch hours report format defaults to JSON

Lines 76-79

The report is returned as JSON when no format is given, or when the format is blank. The format is not case-sensitive and spaces around it are ignored.

## Branch hours report format must be json or csv

Lines 80-82

Turbo refuses any format other than json or csv and shows "format must be json or csv".

## Branch hours report refuses an unknown branch

Lines 84-87

When the branch id matches no branch record, Turbo shows "unknown branch".

## Branch hours report covers staff in the branch or work for it

Lines 89-96

The report lists hours by job and staff member for the chosen branch between the start and end dates. Each row has a direction, the job, client, work branch, staff name, staff branch and hours. The report includes time worked by staff who belong to the branch and time worked for the branch. The query that decides which rows qualify, and what direction means, is in branch_hours.sql, which is not shown here. The report heading shows the branch name. When no rows match, the rows list is empty.

## Branch hours report with no matching time returns no rows

Lines 89-89

When no time matches the branch and dates, the branch hours report returns an empty list of rows. The JSON report shows an empty rows list. The CSV shows only the header row.

## Branch hours report JSON output

Lines 98-100

With the json format, Turbo returns the branch name and the list of rows.

## Branch hours CSV download filename

Lines 106-114

With the csv format, Turbo sends the report as a file download. The file name is branch_hours_<branch>_<start>_<end>.csv. Each run of characters other than letters and digits in the branch name becomes one underscore, and underscores at the start and end are trimmed. The dates are written as YYYY-MM-DD.

## Branch hours CSV leaves out job ID and staff ID

Lines 118-137

The branch hours CSV has a header row, then one row per report row. Its columns are direction, jobNumber, jobDescription, client, workBranch, surname, givenName, staffBranch and hours. The CSV leaves out the job ID (jobId) and the staff ID (uid), which the JSON report includes.

## Branch hours CSV writes hours without trailing zeros

Lines 134-134

The branch hours CSV writes hours as a plain decimal number with no trailing zeros and no exponent, for example 7.5 and not 7.50.
