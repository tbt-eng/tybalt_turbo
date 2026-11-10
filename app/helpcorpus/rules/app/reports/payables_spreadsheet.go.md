<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/payables_spreadsheet.go","sha256":"3057ed3e903b6c8020c505e8a88567dfb573d6a59caa3b060d1e0e070277382f","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/payables_spreadsheet.go

## Payables spreadsheet columns and fixed values

Lines 39-82

Asked as: payables spreadsheet columns; blank columns in export; entered by TURBO; which fields are empty

Each payables row fills the type, job, division, branch, PO type, total, currency, PO number, description, supplier, employee, approved by and status columns. Entered By is always 'TURBO'. Subtotal, HST, Category, Vendor Inv #, Inv Date, Notes, Pd By and TBTE # are always blank.

## Payables date split into day, month and year

Lines 65-73

Asked as: date split into day month year; payables date columns blank; month column empty; wrong date in spreadsheet

The Date, Mon and Year columns come from the record date. Date is the day of the month, Mon is the three-letter month and Year is the full year. If the record date is empty or not in YYYY-MM-DD format, all three are blank.

## Spreadsheet rows are sorted by approval date, oldest first

Lines 84-94

Asked as: payables rows order; sort spreadsheet by approval date; oldest rows first in export; spreadsheet rows out of order

Every payables spreadsheet, daily or monthly, lists its rows in ascending order of approval date (approval_date). Any extra filter is added to the base payables query with AND.

## CSV file has a header row; TSV file has none

Lines 96-121

Asked as: csv header row missing; tsv no headers; export file column names; tab separated payables file

The CSV payables file starts with the header row of column names, then one line per PO row. The TSV file has no header row and only the data rows, separated by tabs, one row per line.

## Payables report needs the report claim

Lines 123-132

Asked as: cannot download payables report; not authorized payables spreadsheet; no access to payables export; report permission denied

Every payables spreadsheet endpoint (list of dates, daily file, monthly file) requires the report claim. A person without it sees "you are not authorized to view this report".

## Which POs appear in the payables dates list

Lines 138-155

Asked as: PO missing from payables dates; why is PO not in dates list; dash PO number not showing; unapproved PO in payables

The dates list includes only POs whose status is not 'Unapproved', that have a PO number, and that have an approval or a second approval. A PO number with a dash is eligible only if the part after the dash is below 5000. A number with no dash is eligible.

## Payables dates list covers last 28 days, excluding today

Lines 146-176

Asked as: payables dates last 28 days; why today missing from dates; which dates can I export; date range for payables

The payables dates list returns each distinct approval date from 28 days ago up to yesterday (UTC), oldest first. Today's UTC date is left out so only complete days appear. Both the 28-day start and yesterday are included.

## PO approval date uses the later second approval

Lines 150-153

Asked as: which approval date is used; second approval date later; PO date wrong in payables; approval date after second approval

A PO's approval date is its second approval date when it has one and it is later than the first approval. Otherwise it is the first approval date. The dates list and the daily file both use this date.

## Daily payables file date validation

Lines 187-197

Asked as: daily payables file error; date is required error; payables date wrong format; cannot pick today for export

The daily payables file requires a date. A missing date gives "date is required". A date not in YYYY-MM-DD format gives "date must be in YYYY-MM-DD format". A date that is today (UTC) or later gives "date must be at least one UTC day old".

## Daily payables file selects POs by approval date

Lines 199-205

Asked as: daily payables file empty; POs for one day export; download POs for a day; no rows for date

The daily payables file returns the PO rows whose approval date matches the requested date, ordered by approval date ascending. Rows come from the payables query, which the file does not show in full.

## Payables file format is CSV unless TSV is requested

Lines 207-218

Asked as: export as tsv; payables file wrong format; csv or tsv download; tab delimited export

The daily and monthly payables files return tab-separated text when the format is 'tsv'. For any other or missing format they return CSV. Only the CSV includes a header row.

## Monthly payables file needs a 4-digit month code

Lines 222-234

Asked as: monthly file yymm error; 4 digit month code; monthly payables invalid month; month code format wrong

The monthly payables file requires yymm to be exactly 4 digits. Otherwise it shows "yymm must be a 4-digit string (e.g. 2603)".

## Monthly payables file selects POs by number prefix

Lines 236-243

Asked as: monthly payables file rows; POs for a month export; which POs in monthly file; month export missing POs

The monthly payables file returns POs whose PO number starts with the yymm code followed by a dash. It does not filter by approval date or by the date range limits of the daily file. Rows are ordered by approval date ascending.

## Monthly payables file also returns the format choice

Lines 245-256

Asked as: monthly file as tsv; monthly export content type; monthly download format; monthly file csv instead of tsv

The monthly payables file returns TSV when the format is 'tsv'. For any other value it returns CSV. It sets the matching content type (text/tab-separated-values or text/csv).
