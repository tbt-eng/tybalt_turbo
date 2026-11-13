<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/utilities.ts","sha256":"7d0aa6dac1539f77513751ef2971725255653cd75c879f91f86416fc3e12e033","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/utilities.ts

## Save error message choice

Lines 21-30

Asked as: save failed which error shows; field messages on save; error message after saving; what message appears when save fails

When a save fails, Turbo shows the field messages joined together if the server gave any. Otherwise it shows the overall message, then the error text, then a fallback message.

## Payroll opening dates must be a Sunday on a 14-day cycle

Lines 44-135

Asked as: payroll date rejected; payroll opening date not accepted; which days can payroll start; payroll cycle every two weeks

Turbo accepts a payroll opening date only if it is a real date in YYYY-MM-DD form that falls on a Sunday. It must also be a whole multiple of 14 days from the epoch Sunday 2025-03-02 (PAYROLL_OPENING_DATE_EPOCH). An empty value is accepted. Dates before the epoch count too, if they are a multiple of 14 days away.

## Earliest date allowed in date inputs

Lines 44-44

Asked as: earliest date in calendar; date picker cannot go back; how far back can I pick a date; date input minimum

Date inputs use 2024-06-01 as their minimum date (DATE_INPUT_MIN).

## Job label format

Lines 65-84

Asked as: job name looks wrong; job label missing number; how job names are shown; job title format

A job label reads "number - client: description". It leaves out whichever parts are missing. The job's id replaces the number when the number is empty. When number and id are both empty, the label uses the description, then the client.

## Listing payroll opening dates between two dates

Lines 137-158

Asked as: list payroll dates in range; payroll dates between two days; show all pay periods between; upcoming payroll dates range

Turbo lists the payroll opening dates from the 14-day cycle that fall between a minimum and a maximum date, and both ends are included. It returns nothing when the minimum is after the maximum.

## Time entry tally by time type

Lines 228-277

Asked as: timesheet hours totals; how hours are counted; why hours not added to job; non-work hours total

Turbo tallies time entries by time type code. Types R and RT count as work hours: entries with no job go to hours, entries with a job go to jobHours and to that job's tally. Both add to the total and add meals hours. OR and OW entries only record their dates. OTO entries become payout requests and RB entries become bank entries. Every other type adds its hours to a non-work tally for that type and to the non-work total.

## Tally warns when a time entry lacks expanded data

Lines 241-245

Asked as: expand field missing alert; time entry alert missing data; entry left out of tally; timesheet entry error message

If a time entry has no expand data, Turbo shows the alert "Error: expand field is missing from time entry record." and leaves that entry out of the tally.

## Division tally records divisions of work entries

Lines 248-251

Asked as: divisions summary missing entry; division totals on timesheet; which entries show division; division list for work hours

Work entries (types R and RT) that have a division add that division's code and name to the divisions tally. Other time types do not.

## Short dates drop the year unless asked

Lines 279-299

Asked as: date shows without year; short date format; why year hidden on date; month and day only

Turbo shows a short date as month abbreviation and day, such as "Mar 2". It adds the year only when the year is requested.

## Date-time display falls back to raw value

Lines 301-312

Asked as: date and time shows wrong; timestamp displayed raw; blank date time; time shown oddly

Turbo shows a date-time as a short date plus the local time. An empty value shows as blank, and a value that is not a valid date shows unchanged.

## Optional text is trimmed and blank when missing

Lines 314-317

Asked as: extra spaces in text; optional field shows dash; empty note displays nothing; stray spaces removed

Turbo trims leading and trailing spaces from optional text before it displays the text (trimmedOrEmpty). A missing or null value shows as an empty string.

## Expense attachment opens in a new tab

Lines 333-377

Asked as: receipt will not open; attachment blank tab; expense file opens empty; receipt download new tab

Turbo fetches an expense attachment with the user's login and opens it in a new tab. When the server's content type is missing or generic, Turbo picks the type from the file extension: .pdf, .png, .jpg or .jpeg, otherwise a generic binary. If the fetch fails, Turbo closes the blank tab and shows no message.

## Dollar amounts are abbreviated in short form

Lines 424-449

Asked as: big amount shortened; k or M in amount; amount shows as thousands; large number abbreviation

Short dollar and number displays abbreviate large values with k, M, G, T, P or E and keep one decimal. They drop trailing zeros, so 1500 shows as $1.5k. Values below 1, including negatives, show as "0" with the "$" sign.

## Currency formatting defaults to CAD

Lines 451-478

Asked as: amount shows CAD by default; currency code missing; amount displays as zero; unknown currency display

Turbo formats amounts as currency with 2 decimals. An empty currency code becomes CAD, and a code is trimmed and upper-cased. An empty, missing or non-numeric amount shows "$0.00". When the currency code is not recognised, Turbo shows the code followed by the amount, such as "XYZ 12.00".

## Currency equivalent line shows CAD amount, rate and date

Lines 480-498

Asked as: exchange rate line; foreign amount equivalent; rate and as of date; CAD equivalent text

The foreign-currency equivalent shows the CAD amount. It adds "rate" with 4 decimals only when the rate is above 0, and adds "as of" with the first 10 characters of the rate date when a date exists. It is blank when the amount is missing.

## Indicative CAD amount calculation

Lines 500-519

Asked as: converted CAD amount wrong; how CAD total is worked out; indicative amount zero; exchange conversion calculation

Turbo calculates the indicative CAD amount as the source amount times the rate, rounded to cents. It gives 0 when the amount or rate is missing, or when the rate is 0 or negative.

## Settlement tolerance is 20% either way

Lines 521-539

Asked as: settlement range limits; how far can settled amount differ; 80 to 120 percent; allowed variance foreign expense

The settlement range for a foreign-currency expense runs from 80% to 120% of the indicative CAD total, rounded to cents. The ratio is 0.2 by default. When the total is missing or 0 or less, both bounds are 0.

## Percent display uses two decimals

Lines 541-546

Asked as: percentage decimals; percent sign format; how percent is shown; rate shown with decimals

Turbo shows percentages with 2 decimals and a % sign.

## Categories load for a job, sorted by name

Lines 555-604

Asked as: categories not showing; job category list empty; category dropdown order; categories disappeared after changing job

Turbo loads categories for the chosen job, sorted by name. With no job it loads none. An editor ignores a late response if the job has changed since. An empty job does not clear categories already shown.

## Client contacts sorted by surname then given name

Lines 606-622

Asked as: client contacts order; contact list sorted wrong; contacts missing for client; who appears first in contacts

Turbo loads a client's contacts sorted by surname, then given name. It loads none when there is no client.

## Contact display name

Lines 624-632

Asked as: contact name shows email; name order in picker; contact shows only first name; how contact names appear

Pickers show a contact as "Surname, Given". With one name only, they show that name. With neither, they show the email, then the id.

## Invoicing profile labels

Lines 634-656

Asked as: invoicing profile name; two profiles look same; invoice recipient label; profile shows no name

An invoicing profile is labelled by its name. Without a name, it uses the recipient contact's name, or "Invoice recipient". The invoicing instructions follow after " — ". When two profiles in a client's list would read the same, Turbo adds "(created date time)" to tell them apart.

## CSV export quoting and empty downloads

Lines 658-698

Asked as: CSV file opens wrong; export file empty; quotes in exported data; download gives no file

CSV exports wrap every value in quotes and double any quotes inside. They leave null values blank. Turbo downloads nothing when there are no rows or no columns. The columns come from the first row when none are given. A template download has a header row only.

## CSV downloads from a report link fail quietly

Lines 699-747

Asked as: report download does nothing; CSV from report fails; no file when exporting report; download error no message

A CSV download from a server endpoint (downloadCSV) sends the user's login token when the user is logged in. If the server answers with an error, Turbo shows no message to the user and downloads no file. The error goes only to the browser console.

## Zip download fetches the file in two steps and reports errors to the caller

Lines 749-808

Asked as: zip download failed; bulk files not downloading; zip file error; download of attachments fails

A zip download (downloadZip) first asks the endpoint for a file URL, then fetches the file from the files API. The second request carries no login token. If either step fails, or the endpoint gives no valid URL, Turbo downloads nothing and passes the error to the screen that started the download.

## Default division fills an empty division once

Lines 810-851

Asked as: division filled automatically; default division wrong; division keeps changing; amendment picks wrong division

On a new record (not an edit), Turbo fills an empty division with the user's default division from their profile. It does this once per record. For a record about another person, such as a time amendment, it uses that person's profile default division instead. It never overwrites a division already set.

## Default role fills an empty role once

Lines 853-863

Asked as: role auto filled; default role wrong; role changed on new record; role not sticking

On a new record (not an edit), Turbo fills an empty role with the user's default role from their profile. It does this once per record. It never overwrites a role already set.

## Live lists add, update and remove records

Lines 867-925

Asked as: list not refreshing; new item not showing; item disappeared from list; record vanished after update

Lists update live. A new record is added at the top only if it passes the page's visibility check. An updated record that no longer passes is removed, and one that cannot be found is removed too. A deleted record is removed.

## Purchase order active date

Lines 927-934

Asked as: PO active date wrong; purchase order start date; when PO becomes active; second approval date shown

A purchase order's active date is the date of its second approval when it has one. Otherwise it is the date it was approved. Only the first 10 characters are used.

## Where Turbo goes after an absorb

Lines 936-956

Asked as: where after absorb; redirected after merge; page after absorbing contact; absorb sends me to

After an absorb, Turbo goes to the client's details page at the contacts section for client contacts. It looks up the client from the target contact when none is given. Rate sheets go to /rate-sheets/duplicates. Other records go to their list page.
