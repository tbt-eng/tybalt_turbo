<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/utilities.ts","sha256":"5cbd8df33362b1421607c74a37e322e2579df5e458e8bb5b37c9aa33a2701ec1","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/utilities.ts

## Save error message preference

Lines 21-30

Asked as: save failed message; error when saving; why won't it save; several field errors together

A failed save shows the API's field messages joined together when it sent any. Otherwise it shows the overall message, then the error text, then a fallback message.

## Payroll opening dates fall every 14 days from 2025-03-02

Lines 44-135

Asked as: payroll date rejected; invalid payroll opening date; which dates are payroll Sundays; payroll date every two weeks

A payroll opening date is valid only if it is a real calendar date written YYYY-MM-DD and falls on a Sunday. It must also be a whole multiple of 14 days from 2025-03-02 (PAYROLL_OPENING_DATE_EPOCH). An empty value is accepted as valid.

## Date pickers: earliest date and look-ahead limits

Lines 44-122

Asked as: date picker earliest date; how far back dates go; date picker too far ahead; can't pick older date

The earliest date that date inputs accept is 2024-06-01 (DATE_INPUT_MIN). The helper dateInputMaxMonthsAhead gives today's date plus the given number of months as the latest date. The helper currentYearStartMinusDays gives 1 January of this year minus the given number of days.

## Job label format

Lines 60-84

Asked as: job name format; job title looks wrong; how job shows in list; job label missing client

A job label reads "number - client: description" when all three parts exist. It drops the missing parts, so it can read "number - description" or "number - client". It uses the job id when the number is missing. If no number or id exists, it shows the description, then the client. Each part is trimmed first.

## Listing payroll opening dates between two dates

Lines 137-158

Asked as: payroll dates in range; list pay periods between dates; payroll dates between two days

Turbo lists every 14-day payroll opening date from 2025-03-02 that falls between the start and end dates, and both ends are included. It returns an empty list when the start date is after the end date.

## How a time sheet's entries are tallied by time type

Lines 228-277

Asked as: timesheet totals wrong; hours summary by job; why hours not adding up; timesheet tally by type

Time entries of type R or RT count as work hours. Entries without a job go to the plain hours tally. Entries with a job go to the job hours tally and to that job's tally, by job number. All work entries add to the work total and meals hours, and record the division by its code and name. OR entries are listed as off-rotation dates and OW entries as off-week dates. OTO entries are listed as payout requests and RB entries as bank entries. Every other time type adds its hours to a non-work tally for that type and to the non-work total.

## Tally skips entries missing expanded data

Lines 241-245

Asked as: entry missing expand error; timesheet expand field alert; time entry left out

When a time entry lacks its expanded related data, Turbo shows the alert "Error: expand field is missing from time entry record." It leaves that entry out of the tally.

## Short date format

Lines 279-299

Asked as: date shows as Mar 2; short date format; date without year

A short date shows the month abbreviation and day, such as "Mar 2". It adds ", year" only when the year is requested.

## Date and time display

Lines 301-312

Asked as: date time shows wrong; timestamp display format; date shows invalid text

A date-time displays as a short month, a 2-digit day and the year, followed by the local time. An empty value shows as blank. A value that is not a valid date shows as it was given.

## Expense attachment opens in a new tab with a fallback file type

Lines 333-377

Asked as: receipt won't open; attachment opens blank tab; expense file opens new tab; attachment download fails silently

Opening an expense attachment fetches the file with the user's sign-in token and shows it in a new browser tab. If the server sends no file type or a generic one, Turbo picks the type from the file name: .pdf, .png, .jpg or .jpeg. Any other file name is treated as a generic binary file. If the fetch fails, the new tab closes and the error is only logged.

## Dollar and number abbreviations

Lines 424-453

Asked as: big numbers shortened; k M suffix numbers; amount shows 1.2k; negative number shows zero

Short dollar amounts and numbers use k, M, G, T, P and E suffixes with at most 1 decimal, and trailing zeros are dropped. A value below 1, including a negative value, shows as "0". Values that are not numbers pass through unchanged.

## Currency formatting defaults to CAD

Lines 451-478

Asked as: amount shows $0.00; currency defaults to CAD; unknown currency code display; currency code missing

Currency amounts show with 2 decimals in en-CA format. A blank currency code means CAD, and codes are trimmed and upper-cased. An empty or invalid amount shows as "$0.00". If the currency code is not recognized, the amount shows as the code followed by the amount, such as "XYZ 12.50".

## CAD equivalent text shows rate and date only when present

Lines 480-498

Asked as: CAD equivalent text; exchange rate shown with amount; as of date missing; rate text under amount

The CAD equivalent shows the CAD amount, then "rate" with 4 decimals only if the rate is above 0, then "as of" with the first 10 characters of the rate date if one is given. The parts are joined with " · ". It shows nothing when the amount is missing.

## Indicative CAD amount is rounded; zero without a valid rate

Lines 500-519

Asked as: indicative amount is zero; CAD conversion wrong; exchange rate not applied; converted amount rounding

The indicative CAD amount is the source amount times the exchange rate, rounded to cents. It is 0 when the amount is missing or the rate is missing or not above 0.

## Foreign currency settlement tolerance is 20 percent

Lines 521-539

Asked as: 20 percent tolerance; settlement range min max; foreign currency variance allowed; settlement range limits

The settlement tolerance range is the indicative CAD total minus 20% up to plus 20% by default, rounded to cents. If the total is missing or 0 or less, both the minimum and the maximum are 0.

## Percent shows two decimals

Lines 541-546

Asked as: percent two decimals; percentage shows wrong decimals; percent sign missing

A numeric percent shows with 2 decimals and a % sign. A value that is not a number passes through unchanged.

## Categories load only when a job is chosen

Lines 555-604

Asked as: categories empty after job; job categories dropdown; category list not updating; categories don't change with job

Category options come from the categories of the chosen job, sorted by name. An empty job gives no categories. When the job changes quickly, only the latest response is applied. An empty job leaves the earlier category options as they are.

## Client contacts load by client and sort by surname

Lines 606-622

Asked as: contacts order by surname; client contacts list sorted; contacts not showing for client

A client's contacts are listed sorted by surname, then given name. A blank client gives no contacts.

## Contact display name fallbacks

Lines 624-632

Asked as: contact shows email; name blank in contact; contact name format; contact shows record id

A contact shows as "Surname, Given" when both exist. With only one name, it shows that one. With neither, it shows the email, then the record id.

## Invoicing profile label and duplicate disambiguation

Lines 634-656

Asked as: invoice profile name duplicate; two invoicing profiles same name; invoice recipient label; invoicing profile shows created time

An invoicing profile is labelled by its name. If it has no name, it uses its recipient contact's name, or "Invoice recipient" if there is none. Its invoice instructions follow after " — ". When another profile of the same client has the same label, Turbo adds " (created date and time)" to tell them apart.

## CSV download from rows

Lines 658-681

Asked as: export csv empty; csv quotes values; download rows as csv; csv blank cells

A CSV download from rows does nothing when there are no rows or no columns. Without given headers, the columns come from the first row's fields. Every value is wrapped in quotes with inner quotes doubled, and empty values become blank.

## CSV template has only headers

Lines 683-686

Asked as: blank csv template; import template headers only; csv file with only headers

A template download is a CSV with only the header row and a trailing newline.

## CSV report download failure handling

Lines 700-753

Asked as: csv download does nothing; report download no error; csv export silent failure

If a CSV download from the server fails, Turbo only logs the error and the user sees nothing. It shows the error to the caller only when the caller asks for it.

## Zip download needs a file URL from the server

Lines 755-814

Asked as: zip download error; download zip file link; zip won't download; zip file link invalid

A zip download first asks the server for a file URL, then downloads the file from it. Turbo raises an error if the server fails or sends no valid URL. The second download request carries no sign-in token.

## Default division is applied once to a new record

Lines 816-857

Asked as: division filled in automatically; new record default division; division not set on time amendment; division blank on new entry

A new record without a division gets the signed-in user's default division (default_division), once for each record. Nothing is applied when editing. When a subject user is given, as in time amendments, the default division of that person's profile applies instead. It applies only if the division is still empty, and it is looked up once per person.

## Default role is applied once to a new record

Lines 859-869

Asked as: role filled automatically; new record default role; role blank on new record

A new record without a role gets the signed-in user's default role (default_role), once for each record. Nothing is applied when editing.

## Live lists drop records that stop being visible

Lines 871-931

Asked as: record disappeared from list; list not updating after edit; new item missing from list; deleted item still showing

Live lists add a created record at the top only if it passes the page's visibility check. When a record is updated and no longer passes the check, it is removed. If its details are not found, it is also removed. A deleted record is removed.

## Purchase order active date

Lines 933-940

Asked as: purchase order active date; PO active date wrong; which approval date for PO

A purchase order's active date is the first 10 characters of the second approval time when it has one. Otherwise it is the first 10 characters of the approval time.

## Where Turbo goes after an absorb

Lines 942-962

Asked as: where after absorb contacts; absorb sends me to; page opens after absorbing; page after absorbing rate sheets

After absorbing contacts, Turbo opens the client's details page at the contacts section. It uses the given client, or the client of the target contact. After absorbing rate sheets, it opens the rate sheet duplicates page. For any other collection, it opens that collection's list page.
