<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/utilities.ts","sha256":"473a0d86687ed9fa06253c540de7598baf8c4690cdd37a97fab2a918e37c401b","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/utilities.ts

## Save errors show field messages first

Lines 24-30

Asked as: save failed, which error shows; form not saving error message; why save error wrong text; save error field message

A failed save shows the field messages joined by spaces. Otherwise it shows the overall message, then the error text, then the fallback text.

## Payroll opening dates fall every 14 days from 2025-03-02

Lines 44-48

Asked as: payroll date rejected; wrong payroll opening day; payroll date not accepted; which dates are payroll

A valid payroll opening date is a real calendar date (YYYY-MM-DD) that falls on a Sunday. It must also be a whole number of 14-day intervals from 2025-03-02 (PAYROLL_OPENING_DATE_EPOCH). An empty value counts as valid.

## Date inputs start at 2024-06-01

Lines 44-44

Asked as: earliest date picker; can't pick old date; date before 2024 blocked; oldest selectable date

The earliest date that date inputs allow is 2024-06-01 (DATE_INPUT_MIN).

## Job label ignores blank parts and trims spaces

Lines 60-63

Asked as: job label extra spaces; blank part in job name; job name spacing; empty part job label

When Turbo builds a job label, it trims spaces from the number, id, client and description (cleanLabelPart). A part that is missing or only spaces counts as empty and is left out of the label.

## Job label format

Lines 65-84

Asked as: job name format; how job label looks; job number client description; job title missing client

A job label shows the number (or the id if there is no number), then the client, then the description, as "number - client: description". Missing parts are left out. With no parts it is empty.

## Date pickers can look ahead a set number of months

Lines 108-115

Asked as: how far ahead date picker; future date limit; date picker max months; can't choose later date

The date input maximum is today's date moved forward by the given number of months (dateInputMaxMonthsAhead). The result is a YYYY-MM-DD date.

## Year-start date limit counts back from 1 January

Lines 117-122

Asked as: earliest date this year; how far back can I date; year start date limit; backdate limit days

The minimum date for some date inputs is 1 January of the current year minus the given number of days (currentYearStartMinusDays). Turbo returns it as YYYY-MM-DD.

## Payroll opening date validation checks format, Sunday and 14-day cycle

Lines 124-135

Asked as: payroll date invalid; why payroll date not Sunday; payroll cycle check; payroll opening wrong format

isValidPayrollOpeningDate accepts an empty value. It rejects anything not in YYYY-MM-DD format or not a real date, and any date that is not a Sunday. It then accepts only dates whose distance from 2025-03-02 is a multiple of 14 days.

## Listing payroll opening dates in a date range

Lines 137-158

Asked as: list of payroll dates; payroll dates between; upcoming payroll opening dates; payroll schedule range

payrollOpeningDatesInRange returns every opening date from min to max, stepping 14 days from 2025-03-02. Both min and max are included when they are opening dates. It returns an empty list when min is after max.

## Tally stops on a time entry with no expand data

Lines 242-245

Asked as: tally alert expand missing; time entry missing data; tally skipped entry error; expand field error timesheet

If a time entry has no expand data, the tally shows the alert "Error: expand field is missing from time entry record." and skips that entry.

## Time sheet tally: work hours split by job or no job

Lines 246-261

Asked as: timesheet hours split; work hours total wrong; hours with and without job; timesheet totals by job

For time entries of type R or RT, hours with no job add to the work hours (hours). Hours with a job add to the job hours (jobHours) and to that job's total by job number. Both also add to the work total, and meals hours are summed. The entry's division code and name are recorded.

## Tally: regular entries record their division

Lines 248-251

Asked as: division shown on entry; regular entry division tally; which division recorded; training entry division

For regular (R) and regular-with-training (RT) time entries, the tally records the entry's division code with the division name when the entry has a division.

## Tally: job hours are counted per job number

Lines 252-259

Asked as: hours per job number; job hours total missing; regular hours by job; hours not added to job

A regular entry that has a job adds its hours to the job hours and to that job's total, keyed by job number. An entry with no job adds its hours to the plain hours.

## Time sheet tally: other time types

Lines 262-273

Asked as: off rotation totals; off week hours counted; payout and bank hours; non work hours tally

OR entry dates go to off-rotation dates and OW entry dates go to off-week dates. OTO entries are payout requests and RB entries are bank entries. Any other time type adds its hours to the non-work tally for that type and to the non-work total.

## Tally: dated off-rotation, off-week, payout and bank entries

Lines 262-272

Asked as: off rotation dates list; off week dates; payout requests list; bank entries not in totals

In the time sheet tally, an off-rotation (OR) entry adds its date to the off-rotation dates. An off-week (OW) entry adds its date to the off-week dates. A payout request (OTO) entry goes to the payout requests list. A bank (RB) entry goes to the bank entries list. None of these entries adds hours to the work or non-work totals.

## Date parts shown without a year by default

Lines 279-299

Asked as: date shows without year; short date format; why no year on date; month and day only

A short date shows the month name and day, such as Mar 5. The year is added only when requested.

## Date and time display falls back to the raw text

Lines 301-312

Asked as: date time display blank; invalid date shows raw; date time formatting wrong; empty date shows nothing

A date-time shows as a short date followed by the time. An empty value shows as blank, and a value that is not a valid date shows unchanged.

## Optional text is trimmed and missing text shows blank

Lines 314-317

Asked as: extra spaces in text; blank text shows null; optional field trimmed; empty value displays nothing

Turbo trims spaces from optional text before it displays it (trimmedOrEmpty). A missing value (null or undefined) displays as an empty string.

## Expense attachment link uses the expenses attachment route

Lines 329-331

Asked as: expense receipt link; attachment url expense; link to expense file; expense attachment address

The link to an expense's attachment points to the server route /api/expenses/attachment/ followed by the expense id (expenseAttachmentHref). The expense id is URL-encoded in the link.

## Expense attachment opens with a file type fallback

Lines 333-377

Asked as: receipt opens new tab; attachment file type wrong; pdf opens as download; expense file opening issue

An expense attachment opens in a new tab. If the server gives no type or application/octet-stream, the type comes from the file extension: pdf, png, jpg or jpeg. Other extensions stay octet-stream. If the fetch fails the tab closes and no message shows.

## Unknown PDF/image types default to a generic file type

Lines 333-339

Asked as: attachment wrong file type; receipt generic file; unknown file extension attachment; file opens as octet stream

When Turbo opens an expense attachment, it takes the file type from the server's Content-Type. If that is missing or application/octet-stream, it guesses from the file name: .pdf, .png, .jpg or .jpeg. Any other extension is opened as application/octet-stream.

## Expense attachment failure closes the opened tab

Lines 341-377

Asked as: blank tab after attachment; attachment tab closes; receipt fails no message; attachment did not load

Turbo opens a blank tab first, then fetches the attachment. If the fetch fails, it closes the tab and shows no message to the user.

## Expense attachment link expires after 60 seconds

Lines 372-372

Asked as: attachment link expired; receipt link stops working; temporary link timeout; attachment link 60 seconds

The temporary link to the opened expense attachment is revoked after 60 seconds.

## Short dollar and number display uses k, M, G suffixes

Lines 424-553

Asked as: short number with k; dollar amount abbreviated; numbers shown as M; display large amounts short

formatDollars and formatNumber abbreviate numbers with k, M, G, T, P and E and at most 1 decimal, with trailing zeros removed. Values below 1 show as 0 (negatives too). formatPercent shows 2 decimals with %. Non-number values are returned unchanged.

## Short number display drops trailing zeros

Lines 424-449

Asked as: trailing zeros in numbers; percent two decimals; negative shows zero; short dollar display

Short dollar, percent and number displays pass non-number values through unchanged. Dollars show as $ plus one decimal with a k, M, G, T, P or E suffix, with trailing zeros removed. Values below 1, including negatives, show as 0. Percent shows two decimals with %.

## Default currency formatting treats missing values as zero

Lines 451-478

Asked as: missing amount shows zero; currency format en-CA; amount blank shows 0.00; default currency display

A missing or non-numeric amount is shown as $0.00. Amounts show with two decimals in the en-CA format. The default currency is CAD.

## Currency defaults to CAD and falls back to a plain format

Lines 455-478

Asked as: unknown currency code; currency code blank; CAD default currency; currency shown as code

A blank currency code is treated as CAD, and codes are trimmed and upper-cased. Amounts show with 2 decimals. A missing or non-numeric amount shows "$0.00". If the currency code is not recognised, the amount shows as "CODE 0.00".

## Currency equivalent text shows CAD amount, rate and date

Lines 480-498

Asked as: CAD equivalent text; exchange rate shown; as of rate date; conversion rate display

The equivalent text starts with the CAD amount. It adds "rate" with 4 decimals only when the rate is above 0. It adds "as of" with the first 10 characters of the rate date when a date exists. It is empty when the amount is missing.

## Indicative CAD amount is source times rate, or 0

Lines 500-519

Asked as: CAD amount calculation; converted amount zero; indicative amount wrong; rate times amount

The indicative CAD amount is the source amount times the rate, rounded to 2 decimals. It is 0 when the amount or rate is missing, or when the rate is 0 or negative.

## Settlement tolerance is plus or minus 20 percent

Lines 521-539

Asked as: settlement range expense; 20 percent tolerance; foreign expense range; allowed variance settlement

The settlement range for a foreign-currency expense is the indicative CAD total minus and plus 20 percent by default (toleranceRatio 0.2), rounded to 2 decimals. If the total is missing or 0 or less, the range is 0 to 0.

## Category lists load only when a job is chosen

Lines 556-604

Asked as: categories not loading; no categories for job; categories after job change; job category list empty

Categories load for the chosen job, sorted by name. An empty job gives no categories. When the job changes, a late response for an older job is ignored. Clearing the job does not clear the categories already shown.

## Category sync keeps old categories when the job is cleared

Lines 592-604

Asked as: categories stay after clearing job; old categories still shown; clear job categories; category list not clearing

When the job is emptied, the category list is not cleared. Only the latest job's response is applied, and older responses are ignored.

## Client contacts load by client, sorted by surname then given name

Lines 607-622

Asked as: client contacts missing; contacts sorted by name; no contacts for client; client contact list

Contacts load for the chosen client, sorted by surname then given name. An empty client gives none.

## Contact name display order and fallbacks

Lines 625-632

Asked as: contact name order; contact shows email; name missing on contact; surname given name display

A contact shows as "Surname, Given". If one is missing, the other shows alone. If both are missing, the email shows, then the id.

## Invoicing profile label and duplicate marker

Lines 635-656

Asked as: invoicing profile name; duplicate invoice profile; invoice recipient label; profile shows created date

An invoicing profile shows its name, or else its recipient contact's name, or else "Invoice recipient". Instructions follow after " — ". If another profile of the client reads the same, "(created date time)" is added.

## CSV download quoting and empty cases

Lines 658-686

Asked as: csv quotes wrong; csv empty download; csv download not working; csv template header

CSV cells are wrapped in double quotes and inner quotes are doubled. Null values become empty. No file downloads when there are no rows or no headers. The header-only template download does not need rows.

## CSV download with no rows does nothing

Lines 658-681

Asked as: csv download nothing happens; no csv file; csv headers missing; export csv blank rows

Downloading CSV rows does nothing when there are no rows or no headers. Without given headers, the first row's keys become the headers. Every value is quoted and inner quotes are doubled, and missing values are blank.

## CSV template has only a header row

Lines 684-686

Asked as: csv template only headers; blank csv template; download header row only; template file empty

The CSV template download contains only the quoted header row followed by a newline.

## CSV and zip downloads fail silently or rethrow

Lines 700-808

Asked as: zip download failed; csv download error silent; zip file error message; download fails no message

A failed CSV download is only logged and the user sees no message. A failed zip download is rethrown to the caller. The zip is fetched in two steps: first the file URL, then the file.

## Default division applies once to a new record

Lines 814-851

Asked as: default division new record; division filled automatically; profile division on new record; division not changing on edit

A new record with no division gets the user's profile default division once. It is not applied when editing or when a division is already set. When a subject user is given (time amendments), the default division of that user's profile is used, fetched once per user.

## Default division with a chosen person is applied from their profile

Lines 819-842

Asked as: division from chosen person; time amendment default division; person default division filled; division for selected user

When a person is chosen (such as on a time amendment), Turbo fills an empty division with that person's profile default division (default_division). It does this once per person for the record. It does not overwrite a division that is already set. It does nothing when editing.

## Default role applies once to a new record

Lines 853-863

Asked as: default role new record; role filled automatically; profile role on new record; role not changing on edit

A new record with no role gets the user's profile default role once. It is not applied when editing.

## Live lists treat a not-found response as a removed record

Lines 877-878

Asked as: live list record missing; list shows 404; record gone from list; updated record disappeared

When a live list reloads a changed record and the server answers 404 (not found), Turbo treats the record as gone (isNotFoundError). A 404 on a new record skips it silently. A 404 on an updated record removes it from the list.

## Live lists drop records that no longer match

Lines 880-925

Asked as: live list record removed; list not matching filter; record no longer visible; deleted record still listed

In live lists, an updated record that no longer passes the page's visibility check is removed. So is one that can no longer be found. A created record that fails the check is not added. Deleted records are removed.

## Live list adds new records only if visible

Lines 885-899

Asked as: new record not appearing; live list adds record; created record missing from list; record added to top

For a live list, a new record is added at the top only if it passes the page's visibility check and its augmented record loads. A record that fails to load, or is not found, is skipped.

## Live list removes updated records that vanish

Lines 900-921

Asked as: updated record vanished; record removed after edit; live list record disappears; edited item gone from list

An updated record is removed from a live list when it no longer passes the visibility check or its augmented record is not found. A deleted record is always removed.

## Purchase order active date

Lines 927-934

Asked as: purchase order active date; PO active date; second approval date shown; PO date wrong

A purchase order's active date is the first 10 characters of its second approval date if it has one. Otherwise it is the first 10 characters of its approval date.

## Redirect after merging records

Lines 937-953

Asked as: after merging contacts; redirect after merge; where merge takes me; merged client contacts page

After merging client contacts, Turbo goes to that client's details page, contacts section. The client comes from the given client or the target contact. After merging other record types, it goes to that collection's list page.
