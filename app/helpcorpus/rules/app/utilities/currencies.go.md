<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/currencies.go","sha256":"72fd558a52242f694d5482ebdae91f6051e8f23afff0e5b95c288d8379b85e9b","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/currencies.go

## Home currency is fixed to CAD

Lines 14-17

Asked as: change home currency; home currency setting; can't change base currency; why is it always CAD

Turbo's home currency is always CAD (HomeCurrencyCode). It cannot be configured, because the exchange-rate provider and the reporting and settlement logic are CAD-based.

## Foreign settlement tolerance is 20 percent

Lines 18-18

Asked as: how much can settled total differ; currency difference allowed percent; foreign amount tolerance; settlement variance limit

A settled total in a foreign currency may differ from the current CAD equivalent by at most 20% (ForeignSettlementToleranceRatio = 0.20).

## Missing currency error message

Lines 20-20

Asked as: currency not found error; no such currency; unknown currency code; currency missing from list

When a currency cannot be found, Turbo reports the error "currency not found". The blank-ID, blank-code and unknown-value cases all use this error.

## Foreign currency without rate error message

Lines 21-21

Asked as: no exchange rate for currency; foreign currency rate missing; can't use foreign currency; rate missing error

When a foreign currency has no exchange rate, Turbo reports the error "foreign currency is missing an exchange rate" (ErrForeignCurrencyRateMissing).

## Missing currency record falls back to CAD at rate 1

Lines 32-44

Asked as: currency record missing; no currency record shows CAD; default rate when record gone; blank currency symbol shows CAD

When there is no currency record, Turbo uses CAD for the code and the symbol. It sets the rate to 1 and leaves the rate date empty.

## Currency record text fields are trimmed

Lines 46-53

Asked as: extra spaces in currency code; currency code whitespace; trailing spaces currency; rate date spaces problem

When Turbo reads a currency record, it trims spaces from the code, symbol, icon and rate date (rate_date). The rate is read as a number from the record's rate field.

## Currency lookup by ID or code fails on blank or unknown value

Lines 56-84

Asked as: currency lookup failed; empty currency code search; search currency by code wrong; currency id blank error

Turbo refuses a currency lookup when the ID or code is blank or no currency matches. The error is 'currency not found', followed by 'blank id', 'blank code' or the value searched. A code lookup trims spaces and ignores case.

## Blank currency falls back to home currency CAD

Lines 90-108

Asked as: no currency on record; empty currency defaults to CAD; missing currency uses home; blank currency becomes CAD

When a record has no currency, Turbo uses the CAD currency record. If that record is missing, it uses CAD with symbol CAD and rate 1. A blank code or symbol becomes CAD, and a zero rate becomes 1.

## Named currency with missing details gets defaults

Lines 110-126

Asked as: currency with no code; empty symbol on currency; currency rate zero fix; unknown currency id error

For a currency that is set, a blank code becomes CAD and a blank symbol becomes the code. A CAD currency with rate 0 gets rate 1. A foreign currency with rate 0 keeps rate 0. An unknown currency ID returns the 'currency not found' error.

## Currency rate defaults to 1 only for CAD

Lines 128-136

Asked as: zero rate currency; rate of 0 for foreign currency; default exchange rate value; why rate shows 1

Turbo uses the currency's rate when it is above 0. With no positive rate, CAD uses 1 and a foreign currency gets 0.

## Foreign currency needs a positive exchange rate

Lines 138-146

Asked as: zero exchange rate rejected; negative rate not accepted; can't save foreign currency rate; foreign rate must be positive

Turbo refuses a foreign currency whose rate is 0 or lower. The error is 'foreign currency is missing an exchange rate'. CAD is always accepted.

## Amounts round to whole cents

Lines 148-157

Asked as: amount rounding cents; rounding decimal places; penny difference amounts; how many decimals amounts

Turbo rounds amounts to 2 decimals. It compares amounts in whole cents, rounded to the nearest cent, to avoid floating-point errors.

## Computed limits floor to whole cents

Lines 159-166

Asked as: limit rounded down; cap amount cents rounding; why limit is lower by a cent; computed limit truncated

A computed limit that falls between cents is rounded down to the cent below, for example 17.0835 becomes 17.08. Turbo first rounds away float error, so a limit that is exactly a whole cent stays at that cent.

## Indicative CAD amount uses the exchange rate

Lines 168-170

Asked as: CAD equivalent amount; convert to CAD indicative; amount in home currency; how is CAD value calculated

Turbo calculates the indicative CAD amount as the amount times the currency rate, rounded to cents. The rate is 1 for CAD.

## Settlement tolerance range is plus or minus 20 percent

Lines 172-177

Asked as: allowed range for settlement; 80 to 120 percent; minimum and maximum settlement; expected CAD range

The expected CAD amount is the amount times the rate. The allowed minimum is 80% of it and the maximum is 120%, each rounded to cents.

## Settled total check includes the range ends

Lines 179-183

Asked as: settled total accepted at limit; boundary settlement amount; edge of tolerance range; settlement exactly at max

A settled total passes when it is at or between the minimum and maximum of the tolerance range. The ends count as inside, with a tiny allowance of 0.000001 for float error.

## Message when settled total is out of tolerance

Lines 185-193

Asked as: settlement amount rejected; settled total too far; settlement warning message; total outside allowed range

Turbo shows: 'settled total must be within 20% of the current CAD equivalent ($expected; allowed range $min to $max)'. The amounts have 2 decimals.

## Blank currency code is treated as CAD

Lines 195-213

Asked as: no currency code on expense; lowercase cad code; code empty treated as CAD; unresolved currency code

When a currency has no code, Turbo treats it as CAD. Codes are trimmed and upper-cased, and CAD matches in any letter case. If the currency cannot be resolved, the effective code is CAD.

## Approval total in CAD falls back to approval total

Lines 215-225

Asked as: approval total in CAD; approval amount home currency; which approval total used; approval total zero

Turbo uses the record's approval total in CAD (approval_total_home) when it is above 0. Otherwise it uses approval_total. A missing record gives 0.

## Email amounts show commas, 2 decimals and currency code

Lines 227-236

Asked as: amount format in email; email currency display; commas in email amounts; how amounts appear in emails

Amounts in emails show thousands separators, 2 decimals and the currency code, for example '12,345.60 CAD'.
