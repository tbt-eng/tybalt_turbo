<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/currency_rates.go","sha256":"9fb01b6c8a8045077eb154c0d5942d697cd2f37fcfde3adf806d5d2d015318f3","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/currency_rates.go

## Rates come from Bank of Canada, last 5 observations

Lines 16-16

Turbo gets a currency's rate from the Bank of Canada Valet series FX{code}CAD, asking for the 5 most recent observations. The request times out after 15 seconds. A reply that is not HTTP 200 counts as an error with the message "unexpected bank of canada status" and the status.

## Rate sync reports how many rates were updated or skipped

Lines 31-34

A currency rate sync counts the currencies whose rate it updated (updated). It also counts the currencies it skipped because the stored rate was already as new or newer (skipped_newer).

## Plain rate sync does not return the update counts

Lines 36-39

The plain currency rate sync (SyncCurrencyRates) runs the full sync and returns only an error. It discards the counts of updated and skipped currencies.

## Rate sync skips the home currency

Lines 49-58

The currency rate sync lists every currency except the home currency (HomeCurrencyCode). The comparison ignores letter case. Turbo never fetches a rate for the home currency.

## Rate sync order of currencies

Lines 49-53

The rate sync processes currencies in order of their display sort value (ui_sort), then by code. A currency with no sort value goes last (treated as 999999).

## Rate sync continues past failures and reports them together

Lines 60-74

The currency rate sync goes through every non-home currency, even when one currency fails. It counts the currencies it updated (updated) and the currencies it skipped because the stored rate is newer or the same (skipped_newer). It returns all failures together, each tagged with its currency code.

## Blank or home currency code is not refreshed

Lines 83-86

When Turbo refreshes one currency's rate, it trims and uppercases the code first. If the code is blank or is the home currency code (HomeCurrencyCode), it does nothing and reports no update and no error.

## Missing or non-positive rate leaves currency unchanged

Lines 88-94

If the Bank of Canada gives no rate, a rate of zero or less, or no rate date, Turbo does not change the currency and reports no error.

## Older or same-day rate does not overwrite stored rate

Lines 100-109

Turbo updates a currency's rate (rate) and rate date (rate_date) only when the fetched date is later than the stored date. If the fetched date is the same or earlier, Turbo keeps the stored rate and counts it as skipped (skipped_newer). A currency with no stored rate date always takes the fetched rate.

## Fetch requires a currency code and uses latest usable observation

Lines 112-173

Fetching a rate fails with "currency code is required" if the code is blank. Otherwise Turbo goes through the observations in the order the Bank gives them and takes the first one that has a non-empty value for the currency's series. It uses that value as the rate and that observation's date as the rate date. If no observation has a value, it returns no rate and no error.
