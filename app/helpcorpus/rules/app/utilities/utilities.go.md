<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/utilities.go","sha256":"d1eabf4eef718cab27f32473cd9b98ba139699cced852195786f8deed41e12fa","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/utilities.go

## Inactive users cannot be assigned as manager or approver

Lines 26-44

Asked as: inactive user as manager; user not active approver; deactivated user assigned approver; admin profile inactive

Turbo counts a user as active only when their admin profile (admin_profiles) has active=true. A user with no admin profile counts as inactive. An empty user ID counts as inactive. Turbo uses this check to stop inactive users being assigned to manager or approver fields.

## Dates must be valid YYYY-MM-DD

Lines 46-52

Asked as: invalid date format; date rejected wrong format; not a valid date; date must be YYYY-MM-DD

Turbo rejects a date that is not a real date in YYYY-MM-DD form. The message is "<value> is not a valid date".

## Week ending is the next Saturday

Lines 54-72

Asked as: week ending date wrong; week end Saturday; date rolls to Saturday; which Saturday is week ending

Turbo sets a week ending date to the first Saturday on or after the date. A date that is already a Saturday stays the same.

## Pay period ending is every second Saturday from payroll epoch

Lines 74-104

Asked as: pay period end wrong; fortnight ending date; payroll Saturday every two weeks; pay period ending calculation

Turbo sets the pay period ending to the week ending date when that date is a whole multiple of 14 days from the payroll epoch (constants.PAYROLL_EPOCH). Otherwise it uses the week ending plus 7 days.

## Time off opening date must follow a pay period end

Lines 106-132

Asked as: time off opening date rejected; opening date not Sunday; time off start after pay period; opening date error

A time off opening date (opening_date) may be empty. Otherwise it must be a Sunday, and the day before it must be a pay period ending date. The errors are "opening_date must be a Sunday" and "opening_date must follow a pay period ending date".

## Committed expense pay period depends on commit week

Lines 134-191

Asked as: committed expense wrong pay period; expense landed in wrong payroll; commit week pay period mismatch; expense goes to previous payroll

When a committed expense's commit week ending is a payroll week-2 ending, Turbo uses that date as the pay period ending. When the commit falls in payroll week 1, Turbo compares the expense date with the previous payroll ending (one week before the commit week ending). An expense dated on or before that day goes to the previous payroll ending. A later expense goes to the week ending one week after the commit week ending.

## Date limit validation includes the limit date

Lines 193-214

Asked as: date too early; date too late; date limit error; date before cutoff

A date limit check accepts a date equal to the limit. As a minimum, an earlier date fails with "<date> is too early". As a maximum, a later date fails with "<date> is too late".

## PO approver must hold claim and cover the division

Lines 216-254

Asked as: cannot pick PO approver; approver missing claim; PO approver not allowed division; approver not in division

Turbo accepts a PO approver only when the user holds the required claim and has a po_approver_props record. That record must have an empty divisions list or include the PO's division. A missing or empty user fails with "value is required". A user who fails the claim or division check gets "user does not have the required claim".

## Time amendment hours must be multiples of 0.5

Lines 256-271

Asked as: time amendment hours rejected; half hour amendment; hours not half hour; negative amendment hours

Hours must be a multiple of 0.5. Negative values are allowed, because amendments can correct committed hours down. Zero passes. Otherwise the message is "must be a multiple of 0.5".

## Positive hours must be non-negative multiples of 0.5

Lines 273-288

Asked as: hours must be half hours; negative hours error; positive hours check; hours increments of half

This check accepts zero. A negative value fails with "must be a positive number". A value that is not a multiple of 0.5 fails with "must be a multiple of 0.5".

## Amounts must be non-negative multiples of 0.01

Lines 290-305

Asked as: amount cents rejected; money value decimal error; amount not whole cents; negative amount not allowed

This check accepts zero. A negative value fails with "must be a positive number". A value that is not a multiple of 0.01 fails with "must be a multiple of 0.01".

## Mileage distance must be a whole number

Lines 310-315

Asked as: mileage distance decimal; mileage km fraction error; distance must be whole; fractional kilometres mileage

When Turbo calculates a mileage expense total, a fractional distance fails with "distance must be an integer for mileage expenses".

## Mileage counts earlier committed mileage since last reset date

Lines 317-367

Asked as: mileage rate reset; mileage counter since reset; previous mileage claims rate; mileage year to date count

The mileage rate depends on the user's committed mileage expenses since the latest reset date (mileage_reset_dates) on or before the expense date. If no reset date exists, Turbo counts from the beginning. It counts earlier-dated expenses and same-day expenses with a smaller id. It ignores uncommitted expenses.

## Mileage total uses tiered rates and rounds to cents

Lines 369-512

Asked as: mileage total wrong; tiered mileage rate; no mileage rates; mileage rounding cents

Turbo reads the rate bands from the expense rate record (mileage). Each key is a starting kilometre and the value is the rate. Distance that crosses bands is split and each part is paid at its own band's rate. The total is rounded to 2 decimals. An error occurs when the rate data is missing or unreadable, or no bands exist ("no mileage rates found").

## Annual payroll period starts after previous year end date

Lines 514-540

Asked as: annual payroll start date; payroll year start wrong; no year end date found; year end payroll period

The annual payroll period starts the day after the latest payroll year end date that is before the given date. A date equal to a year end date uses the previous year end. When no earlier year end date exists, Turbo returns the error "no payroll year end date record found for the given date".

## Claim check by user

Lines 542-573

Asked as: user does not have claim; claim missing from user; check user permission; user has no claim

A user holds a claim when they have a user_claims record for a claim with that name. No user gives no claim.

## Expense rates use latest effective date on or before expense

Lines 575-599

Asked as: expense rate not found; no rate for expense date; which expense rate applies; missing expense rate record

Turbo uses the expense rate record with the latest effective_date on or before the expense date. When none exists, it fails with "no expense rate record found for the given date".

## Recurring PO requires end date and frequency

Lines 601-625

Asked as: recurring PO missing end date; recurring order no frequency; repeating purchase order needs; recurring order required fields

A recurring purchase order without an end date fails with "end_date is required for recurring purchase orders". One without a frequency fails with "frequency is required for recurring purchase orders".

## Recurring PO schedule runs from its date to its end date

Lines 627-638

Asked as: recurring PO schedule dates; payment schedule for recurring order; how recurring order payments run; recurring PO total across dates

Turbo reads the first payment date from the purchase order's date (date) and the last boundary from its end date (end_date). It multiplies the order's total (total) by the number of payments. The frequency (frequency) sets how the payments are spaced.

## Recurring PO end date must be after start date

Lines 640-649

Asked as: end date before start; recurring order dates reversed; end date earlier than start; recurring PO date order

A recurring purchase order whose end date is on or before its start date (date) fails with "end_date is before start_date".

## Recurring PO needs valid frequency and two occurrences

Lines 651-671

Asked as: invalid recurring frequency; recurring PO only once; recurring order occurs once; frequency not allowed recurring

The frequency must be Weekly, Biweekly or Monthly, otherwise Turbo says "invalid frequency". A recurring PO with fewer than 2 occurrences fails with "recurring purchase order must occur at least twice". The hint is to adjust the end date or the frequency.

## Recurring PO total is total times occurrences

Lines 673-673

Asked as: recurring PO total amount; recurring order overall value; total for repeating order; multiply payment total

Turbo sets a recurring purchase order's total value as its per-payment total multiplied by the number of occurrences.

## Recurring occurrences exclude the end date

Lines 676-705

Asked as: how many recurring payments; does last payment count; recurring payment count end date; monthly recurring payment dates

Turbo counts payments from the start date up to, but not including, the end date. The first payment falls on the start date. Weekly payments are 7 days apart and Biweekly payments are 14 days apart. Monthly payments fall on the start day of the month, or the last day of shorter months, always measured from the start date.

## Recurring PO exhausted when commits reach occurrences

Lines 707-733

Asked as: recurring PO used up; recurring order fully committed; no more recurring payments; occurrences reached commits

A recurring purchase order is exhausted when its committed expenses plus commits in progress are at least its occurrences count. Drafts and uncommitted expenses are not counted. The occurrences count is fixed when the PO is saved.

## Cumulative PO expense totals

Lines 735-772

Asked as: cumulative PO total wrong; PO expense total sum; PO running total expenses; total excluding edited expense

The expense total for a cumulative purchase order sums the totals of its expenses, or only committed ones when requested. For an update, the total can leave out the expense being edited so its old total is not counted twice. These functions do not check that the PO is cumulative.

## PO approval tier chosen by amount

Lines 774-804

Asked as: which approval tier; PO approver level by amount; approval tier for amount; amount over all tiers

Turbo picks the po_approval_tiers record with the smallest max_amount that is greater than or equal to the PO amount, so an amount equal to a tier's limit stays in that tier. It returns no claim when the amount exceeds every tier or no tiers exist.

## Lowest and highest approval tiers

Lines 806-832

Asked as: lowest approval tier; highest approval tier; no approval tiers; top tier approval limit

Turbo can look up the lowest or highest po_approval_tiers record and return its claim and max_amount. It fails with "no approval tiers found" when there are none.

## Local edits mark records as not imported

Lines 834-870

Asked as: record not sent to legacy; edit not imported; imported flag reset; record resent legacy system

When an existing record changes in any field except created, updated and the imported flag, Turbo sets _imported to false. This makes the record written back to the legacy system. A change to only those fields does not.

## Job changes mark jobs as not imported

Lines 872-911

Asked as: job not imported; job marked unimported; job reference change import; job legacy sync flag

Turbo can set _imported to false and refresh updated on a job. It can do the same for all jobs that reference a client, job_owner or contact. Any other reference column is rejected.

## Machine token validation

Lines 932-960

Asked as: machine token invalid; expired machine token; machine secret not working; token check failing

A machine token is valid only when an unexpired machine_secrets record with the given role has a matching salted SHA-256 hash. An expired record never matches.
