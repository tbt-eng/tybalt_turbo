<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/utilities.ts","sha256":"a3bea9b7a1640ed3918c9a137988fc1b62f047bdc863953a95f2f98584053d9a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/utilities.ts

## Save error message choice

Lines 24-30

When a save fails, Turbo shows the API's field messages joined by spaces. If there are none, it shows the overall message, then the error text, then a fallback message that the screen supplies.

## Payroll opening dates fall every 14 days on Sundays

Lines 44-135

A payroll opening date is valid only if it is a real calendar date written YYYY-MM-DD, falls on a Sunday, and is a whole multiple of 14 days from the epoch 2025-03-02 (PAYROLL_OPENING_DATE_EPOCH, PAYROLL_OPENING_DATE_INTERVAL_DAYS). An empty value is accepted as valid. Dates before the epoch also count if they are a multiple of 14 days away.

## Date input limits: earliest date and months ahead

Lines 44-44

The earliest date a date input accepts is 2024-06-01 (DATE_INPUT_MIN). A maximum date can be set a given number of months after today. A minimum date can be set a given number of days before 1 January of the current year.

## Job label format and fallbacks

Lines 65-84

A job label reads 'number - client: description'. When the client is missing it reads 'number - description'. When the description is missing it reads 'number - client'. Turbo uses the job id when there is no number. With no number and no id, the label falls back to the description, then to the client. Parts are trimmed, and a missing job gives an empty label.

## Date input range limits for months ahead and year start

Lines 108-122

The latest allowed date on a date input is today plus the given number of months (dateInputMaxMonthsAhead). A helper gives the date that is a given number of days before 1 January of the current year (currentYearStartMinusDays).

## Empty payroll opening date is accepted as valid

Lines 124-135

A blank payroll opening date passes validation (isValidPayrollOpeningDate). A non-blank value must be written YYYY-MM-DD and be a real calendar date. It must also fall on a Sunday and be a multiple of 14 days from 2025-03-02.

## Listing payroll opening dates between two dates

Lines 137-158

Turbo lists payroll opening dates from the 14-day cycle that fall between a start and an end date, and includes both end dates when they match. When the start is after the end, the list is empty.

## Time entries are tallied by time type

Lines 240-273

Turbo tallies time entries by time type code. Types R and RT count as work hours. OR entries list off-rotation dates. OW entries list off-week dates. OTO entries are payout requests. RB entries are bank entries. Every other type adds its hours to non-work hours by code and to the non-work total.

## Time entry missing expand data raises an alert

Lines 241-244

If a time entry has no expanded related records, Turbo shows the alert 'Error: expand field is missing from time entry record.' and leaves that entry out of the tally.

## Work hours split into job hours and other hours

Lines 247-260

For R and RT entries with no job, the hours go to work hours (hours). With a job, they go to job hours and to that job's tally by job number. Both add to the work total. Meal hours add to the meals total, and the entry's division is recorded by code and name.

## Division tally records divisions of work entries

Lines 247-250

For work time entries (types R and RT) that have a division, Turbo records the division name under its code in the divisions tally (divisionsTally).

## Meal hours are summed only for work entries

Lines 259-260

Turbo adds meals hours (meals_hours) to the meals tally only for work entries (types R and RT). Other time types do not count toward it.

## Other time types are tallied as non-work hours

Lines 261-272

Time types other than R, RT, OR, OW, OTO and RB add their hours to a non-work tally for that type code and to the non-work total. Off rotation (OR) and off week (OW) entries add only their dates to lists. Payout request (OTO) and bank (RB) entries are collected as separate lists, and they add no hours.

## Short date and date-time display

Lines 278-311

A short date shows as 'Mon D', such as 'Mar 2', and adds ', YYYY' only when the year is requested. A date-time shows as a short US-style date plus the local time. An empty value shows as empty. An unparseable value is shown unchanged.

## Trimmed text helper returns empty for missing values

Lines 313-316

Optional text shows as an empty string when it is missing. Otherwise Turbo trims leading and trailing whitespace (trimmedOrEmpty).

## Time entry attachment and file URL helpers build links

Lines 318-330

Turbo builds a PocketBase file link from the collection id, record id and file name, and encodes each part. An expense attachment link points to /api/expenses/attachment/ followed by the expense id. A currency icon link uses the "currencies" collection.

## Expense attachment opens in a new tab with a file type fallback

Lines 332-376

Opening an expense attachment fetches the file with the user's login token and opens it in a new tab. If the server gives no type or a generic one, Turbo picks the type from the file extension: .pdf, .png, .jpg or .jpeg, otherwise a generic binary type. The temporary link lasts 60 seconds. On failure the new tab is closed and the error goes only to the console, with no message to the user.

## Dollar and number abbreviations

Lines 423-452

Dollar and number displays shorten large values with suffixes k, M, G, T, P and E, with at most one decimal place and trailing zeros removed. Dollar values get a '$' prefix. Values below 1, including negatives, show as '0'. Non-number values are shown unchanged.

## Currency formatting defaults to CAD

Lines 450-477

Currency amounts show with two decimals in en-CA format. An empty or missing currency code means CAD, and codes are trimmed and put in upper case. A missing or non-number amount shows as '$0.00'. If the currency code is not recognized, the amount shows as the code followed by the amount, such as 'XYZ 12.00'.

## Currency equivalent shows rate and rate date

Lines 479-497

A currency equivalent shows the CAD amount, then 'rate' with 4 decimals only when the rate is a positive number, then 'as of' with the first 10 characters of the rate date only when a date is given. The parts are joined with ' · '. A missing amount gives empty text.

## Indicative CAD amount and settlement tolerance

Lines 499-538

The indicative CAD amount is the source amount times the rate, rounded to cents. It is 0 when the amount or the rate is missing, or when the rate is zero or negative. The settlement tolerance range is 20% below and above the indicative CAD total by default, rounded to cents. It is 0 to 0 when the total is missing or zero or negative.

## Percent display

Lines 540-545

A percentage shows with 2 decimals and a '%' sign. Non-number values are shown unchanged.

## Categories load only for a selected job

Lines 555-603

Turbo lists the categories of a job sorted by name. When no job is selected, the list is empty. Editors keep the existing categories when the job is cleared, and ignore out-of-date responses after fast job changes.

## Client contacts list and naming

Lines 606-631

A client's contacts are listed by surname then given name, and an empty client gives an empty list. A contact shows as 'Surname, Given'. If a part is missing, the other part is shown alone. With neither, Turbo shows the email, then the contact id.

## Invoicing profile label and duplicate disambiguation

Lines 634-655

An invoicing profile is labelled by its name. With no name it uses the recipient contact's name, or 'Invoice recipient' if there is none. Invoicing instructions follow after ' — '. When another profile of the same client has identical text, Turbo adds '(created date-time)' to tell them apart.

## CSV export quoting and empty results

Lines 657-697

CSV exports put every value in double quotes, double any quote inside it, and leave null or undefined values blank. No file downloads when there are no rows or no columns. Without given headers, the first row's keys are the headers. A template download has only the header row.

## Server CSV and zip downloads

Lines 699-807

A CSV download from the server uses the user's login token and saves the file under the given name. A failure is only logged to the console. A zip download fetches a file address first, then the file, and passes any failure to the caller. The second fetch does not send the login token.

## Default division fills a new record once

Lines 813-850

For a new record (not being edited) with no division, Turbo fills in the user's default division (default_division). It does this only once per record. When a subject user is given, as in time amendments, it uses that user's profile default division instead. It repeats for a different user but not for the same one.

## Default role fills a new record once

Lines 852-862

For a new record (not being edited) with no role, Turbo fills in the user's default role (default_role) from their profile. It does this only once per record.

## Live list updates after record changes

Lines 866-924

A list updates live when a record is created, updated or deleted. A created record is added at the top only if it passes the page's visibility check. An updated record that no longer passes is removed. An updated record that is not found is removed. A deleted record is removed. A failed load of the record is only logged.

## Purchase order active date

Lines 926-933

A purchase order's active date is the first 10 characters of its second approval date (second_approval) when it has one. Otherwise it is the first 10 characters of its approval date (approved).

## Redirect after absorbing records

Lines 936-953

After an absorb, a client contact goes to that client's details page at the contacts section. Turbo uses the given client, or looks up the contact's client. Any other record type goes to its list page.
