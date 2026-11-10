<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/currencies.go","sha256":"ea11a84666c3bf00b620df0a16d0f4c31a8c7523d8635f9d5cda661e8559f0eb","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/currencies.go

## Currency code is trimmed and uppercased on save

Lines 15-17

Asked as: currency code lowercase; save currency code uppercase; extra spaces in currency code; currency symbol spaces

When a currency is saved, Turbo trims spaces from the code and converts it to uppercase. It also trims spaces from the symbol.

## Home currency rate is always set to 1

Lines 19-22

Asked as: home currency rate wrong; can't change home currency rate; rate stays at 1; base currency exchange rate

If the currency code matches the home currency code (utilities.HomeCurrencyCode), Turbo sets the rate to 1. The comparison ignores case, and any rate entered is overwritten. Turbo skips the positive-rate check for the home currency.

## Foreign currency needs a positive rate

Lines 24-35

Asked as: foreign currency rate zero; currency rate must be positive; negative exchange rate error; can't save currency rate

Turbo refuses to save a foreign currency whose rate is zero or negative. The error on the rate field is: "foreign currencies must have a positive rate" (code must_be_positive).

## Empty currency code skips the rate check

Lines 24-37

Asked as: blank currency code rate; no code currency saves; empty currency code error; currency without code

If the currency code is empty after trimming, Turbo does not check the rate, so a rate of zero or less does not cause this error. The currency is not refused by this check.
