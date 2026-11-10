<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/utilities.go","sha256":"0bc24d5717167ed7b2fa2962a39719eae55a8609467155a050ee60a9dced7673","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/utilities.go

## Inactive or profile-less users count as inactive

Lines 26-44

Asked as: user not active; can't assign approver; deactivated user approver; manager not showing; inactive user as manager

Turbo treats a user as active only when an admin profile (admin_profiles) exists with active=true. A user with no admin profile, an empty user ID or active=false is inactive. This check stops inactive users from being assigned as manager or approver (IsUserActive).

## Date fields must be valid YYYY-MM-DD dates

Lines 46-52

Asked as: invalid date error; date not accepted; wrong date format; is not a valid date; date field rejected

Turbo rejects a date value that is not a valid YYYY-MM-DD date. The message is "<value> is not a valid date".

## Week ending is the next Saturday, or the same day

Lines 54-72

Asked as: week ending wrong day; week ending moved; timesheet end date Saturday; why week ending changed; week end date shifted

Turbo sets the week ending to the first Saturday on or after the given date. A date that is already a Saturday stays unchanged (GenerateWeekEnding).

## Pay period ending is every second Saturday from payroll epoch

Lines 74-104

Asked as: pay period end date wrong; fortnight ending date; which Saturday is pay period; biweekly pay period end; pay period date off

Turbo takes the week ending of the date. If that Saturday is a whole multiple of 14 days from the payroll epoch (PAYROLL_EPOCH), it is the pay period ending. Otherwise the pay period ending is the Saturday 7 days later. A non-multiple of 24 hours in the interval gives an error.

## Time off opening date must follow a pay period end

Lines 106-132

Asked as: opening date not Sunday; time off start date error; opening_date must be a Sunday; time off opening rejected; time off start date

A time off opening date (opening_date) may be empty. If set, it must be a Sunday, with the error "opening_date must be a Sunday". The Saturday before it must be a pay period ending date, or the error is "opening_date must follow a pay period ending date".

## Committed expense pay period bucketing

Lines 134-191

Asked as: expense in wrong pay period; expense landed previous period; committed expense wrong period; expense period bucket; expense assigned to wrong period

Turbo sets the pay period ending of a committed expense from its date and the commit's week ending. If the commit week ending is itself a pay period ending (payroll week 2), that date is used. If the commit is in payroll week 1 and the expense date is on or before the previous Saturday (7 days earlier), the previous Saturday is used. Otherwise the Saturday 7 days after the commit week ending is used.

## Date limit validation includes the limit date

Lines 193-214

Asked as: date too early error; date too late error; limit date allowed; date limit check; is too late message

A date limit check accepts the limit date itself. With a minimum, an earlier date fails with "<date> is too early". With a maximum, a later date fails with "<date> is too late". An unparseable date also fails.

## PO approver division permission check

Lines 216-254

Asked as: can't approve PO for division; approve purchase orders division; division approver permission; no approval rights division; PO approver not allowed

A user may approve purchase orders for a division only if they hold the given claim and their approver properties (po_approver_props) list no divisions or list that division. An empty divisions list means all divisions. Otherwise the error is "user does not have the required claim". An empty value gives "value is required".

## Hours must be multiples of 0.5, negatives allowed

Lines 256-271

Asked as: hours must be half hour; half hour increments; negative hours allowed; time amendment negative hours; hours not multiple half

A value must be a multiple of 0.5 or Turbo says "must be a multiple of 0.5". Zero passes. Negative values pass, because time amendments can correct committed hours either way.

## Positive multiple of 0.5 validation

Lines 273-288

Asked as: negative hours not allowed; must be positive number; half hour rounding error; hours value rejected; positive hours only

A value must be zero or a positive multiple of 0.5. A negative value fails with "must be a positive number". A non-multiple fails with "must be a multiple of 0.5".

## Positive multiple of 0.01 validation

Lines 290-305

Asked as: amount must be cents; two decimal places error; money value rounding; negative amount not allowed; amount multiple 0.01

A value must be zero or a positive multiple of 0.01. A negative value fails with "must be a positive number". A non-multiple fails with "must be a multiple of 0.01".

## Mileage distance must be a whole number

Lines 310-315

Asked as: mileage decimal distance error; km must be whole number; distance fraction not allowed; mileage distance integer; decimal kilometres mileage

A mileage expense with a distance that is not an integer fails with "distance must be an integer for mileage expenses".

## Mileage total uses cumulative distance since last reset

Lines 317-367

Asked as: mileage rate reset; mileage total wrong; why mileage price lower; kilometres since reset; mileage count since reset

Turbo prices a mileage expense by the committed mileage the same user already has since the latest mileage reset date on or before the expense date. With no reset date, all earlier mileage counts. Only committed Mileage expenses count. Same-day expenses count when their id is smaller.

## Mileage rate bands are tiered by cumulative distance

Lines 369-512

Asked as: mileage rate tiers; mileage rate went down; no mileage rates found; mileage rate calculation wrong; mileage bands pricing

Each band's lower bound is in km and has its own rate. Distance that falls in a band is paid at that band's rate, and an expense that crosses bands is split across them. The result is rounded to 2 decimals. Errors occur if the rates are not JSON, a band key is not a number, or no rates exist ("no mileage rates found").

## Annual payroll period starts the day after year end

Lines 514-540

Asked as: annual payroll period start; year end payroll date; no payroll year end record; payroll year start date; annual payroll start day

Turbo starts the annual payroll period the day after the latest payroll year end date before the given date. A date equal to a year end date uses the earlier year end. If none exists the error is "no payroll year end date record found for the given date".

## Claim check is false without a user

Lines 542-573

Asked as: user doesn't have claim; permission not showing; claim check false; missing user claim; user lacks permission

A user has a claim only if a user claim record links them to a claim with that name. A missing user or empty user ID has no claims.

## Expense rate is the latest effective on or before the expense date

Lines 575-599

Asked as: no expense rate found; rate for expense date; expense rate missing; which expense rate applies; expense rate record error

Turbo picks the expense rate record with the latest effective date on or before the expense date. If none exists the error is "no expense rate record found for the given date".

## Recurring PO requires end date and frequency

Lines 601-625

Asked as: recurring PO missing end date; frequency required recurring; end_date is required; recurring order no frequency; recurring PO needs fields

A recurring purchase order without an end date fails with "end_date is required for recurring purchase orders". One without a frequency fails with "frequency is required for recurring purchase orders".

## Recurring PO dates must be valid dates

Lines 626-637

Asked as: recurring PO date parsing error; recurring total value fails; recurring PO date error; recurring dates unreadable; recurring calculation date failure

When Turbo calculates the total value of a recurring purchase order, it reads the start date (date) and end date (end_date) as YYYY-MM-DD dates. A date that cannot be read stops the calculation with a date-parsing error. Turbo gives no custom message for this case.

## Recurring PO end date must be after start date

Lines 640-649

Asked as: end date before start date; recurring end before start; end date earlier than start; recurring end date wrong; end_date is before start_date

A recurring purchase order whose end date is on or before its start date fails with "end_date is before start_date" (end_date must be after start_date).

## Recurring PO frequency must be known

Lines 651-660

Asked as: invalid frequency error; frequency not accepted; weekly biweekly monthly options; unknown recurring frequency; frequency value rejected

An unknown frequency fails with "invalid frequency". Valid values are Weekly, Biweekly and Monthly.

## Recurring PO must occur at least twice

Lines 662-671

Asked as: recurring occurs once only; only one occurrence; recurring PO too short; recurring must repeat twice; single occurrence error

A recurring purchase order with fewer than 2 occurrences fails with "recurring purchase order must occur at least twice". The hint says to adjust the end date or the frequency.

## Recurring PO total value is total times occurrences

Lines 673-673

Asked as: recurring PO total wrong; total value multiplied; recurring value calculation; recurring total times occurrences; recurring order total amount

Turbo sets the total value of a recurring purchase order to its total multiplied by its number of occurrences.

## Recurring occurrences count payments before the end date

Lines 676-705

Asked as: how many recurring payments; number of occurrences; end date excluded payments; monthly payments count; recurring payment dates

The start date is the first payment and the end date is excluded. Weekly and Biweekly payments repeat every 7 or 14 days. Monthly payments fall on the start day of the month, or the last day of shorter months, always measured from the start date. A one-year monthly term of 2026-02-01 to 2027-02-01 has 12 payments.

## Recurring PO is exhausted when commits reach occurrences

Lines 707-733

Asked as: recurring PO used up; recurring PO closed; no more recurring commits; occurrences all committed; recurring order finished

A recurring purchase order is exhausted when its committed expenses plus pending commits reach or equal its occurrences count (>=). Drafts and uncommitted expenses do not count. The occurrences count was fixed when the PO was saved.

## Cumulative PO expense total, all or committed only

Lines 735-755

Asked as: PO total of expenses; sum of PO expenses; PO spent so far; total of committed expenses; expenses total on PO

Turbo sums the totals of a purchase order's expenses, or only the committed ones when requested. No expenses gives 0. The function does not check that the PO is Cumulative.

## Required PO approver tier is the smallest covering the amount

Lines 757-787

Asked as: which approval tier applies; approver level for PO; amount above all tiers; PO needs which approver; approval tier amount

Turbo picks the approval tier with the lowest maximum amount that is at least the PO amount (an amount equal to the maximum qualifies). If the amount exceeds every tier, or there are no tiers, the claim is empty.

## Lowest or highest approval tier lookup

Lines 789-815

Asked as: no approval tiers found; lowest approval tier; highest approval tier amount; top approval tier lookup; approval tier missing

Turbo can return the claim and maximum amount of the lowest or highest approval tier. If no tiers exist the error is "no approval tiers found".

## Edits mark records as not imported for writeback

Lines 817-853

Asked as: record sent back to legacy; imported flag reset; why record not imported; writeback after edit; edit resets imported flag

When an existing record changes in any field other than created, updated or _imported, Turbo sets _imported to false. The record is then written back to the legacy system. Unchanged records keep their flag.

## Job changes mark jobs not imported

Lines 855-894

Asked as: job not imported; invalid jobs reference column; job update legacy sync; job change writeback; jobs client owner contact update

Turbo marks a job as not imported and refreshes its updated time. It can also do this for all jobs that reference a given client, job owner or contact. Any other reference column gives the error "invalid jobs reference column".

## Only tables with an _imported column are flagged for writeback

Lines 896-913

Asked as: table has imported column; direct SQL imported flag; writeback table check; imported column missing; SQL update imported flag

Turbo checks whether a table has an `_imported` column (TableHasImportedColumn) by reading the table's column list. It reports true only when a column named `_imported` exists. Direct SQL updates use this check to decide whether to set `_imported = false`, because those updates bypass the record hooks.

## Machine token valid only if unexpired and hash matches

Lines 915-943

Asked as: machine token rejected; expired machine secret; token hash mismatch; invalid machine token; machine secret not working

A machine token is valid only if an unexpired machine secret with the same role has a SHA-256 hash equal to salt plus token. An expired secret or a mismatch is rejected.
