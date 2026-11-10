<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/currency_rates.go","sha256":"9fb01b6c8a8045077eb154c0d5942d697cd2f37fcfde3adf806d5d2d015318f3","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/currency_rates.go

## Rates come from Bank of Canada, last 5 observations

Lines 16-16

Asked as: bank of canada rate source; exchange rate feed timeout; fx rate fetch failed; bank of canada status error

Turbo gets a currency's rate from the Bank of Canada Valet series FX{code}CAD, asking for the 5 most recent observations. The request times out after 15 seconds. A reply that is not HTTP 200 counts as an error with the message "unexpected bank of canada status" and the status.

## Rate sync reports how many rates were updated or skipped

Lines 31-34

Asked as: how many rates updated; rate sync summary counts; currencies skipped in sync; sync result numbers

A currency rate sync counts the currencies whose rate it updated (updated). It also counts the currencies it skipped because the stored rate was already as new or newer (skipped_newer).

## Plain rate sync does not return the update counts

Lines 36-39

Asked as: sync rates returns no counts; plain currency sync output; why sync shows no numbers; rate sync just error

The plain currency rate sync (SyncCurrencyRates) runs the full sync and returns only an error. It discards the counts of updated and skipped currencies.

## Rate sync skips the home currency

Lines 49-58

Asked as: home currency not fetched; why no rate for cad; skip base currency sync; home currency excluded from rates

The currency rate sync lists every currency except the home currency (HomeCurrencyCode). The comparison ignores letter case. Turbo never fetches a rate for the home currency.

## Rate sync order of currencies

Lines 49-53

Asked as: currency sync order; which currency syncs first; sort order of currencies in rates; currency without sort value

The rate sync processes currencies in order of their display sort value (ui_sort), then by code. A currency with no sort value goes last (treated as 999999).

## Rate sync continues past failures and reports them together

Lines 60-74

Asked as: one currency fails sync; sync keeps going after error; all currency rate errors together; sync failures reported with code

The currency rate sync goes through every non-home currency, even when one currency fails. It counts the currencies it updated (updated) and the currencies it skipped because the stored rate is newer or the same (skipped_newer). It returns all failures together, each tagged with its currency code.

## Blank or home currency code is not refreshed

Lines 83-86

Asked as: refresh blank currency code; refresh home currency does nothing; single currency rate refresh ignored; lowercase currency code refresh

When Turbo refreshes one currency's rate, it trims and uppercases the code first. If the code is blank or is the home currency code (HomeCurrencyCode), it does nothing and reports no update and no error.

## Missing or non-positive rate leaves currency unchanged

Lines 88-94

Asked as: no rate from bank; zero exchange rate ignored; rate missing currency unchanged; no rate date from bank

If the Bank of Canada gives no rate, a rate of zero or less, or no rate date, Turbo does not change the currency and reports no error.

## Older or same-day rate does not overwrite stored rate

Lines 100-109

Asked as: rate not updated same day; old rate date overwrite; why stored rate kept; rate already newer skipped

Turbo updates a currency's rate (rate) and rate date (rate_date) only when the fetched date is later than the stored date. If the fetched date is the same or earlier, Turbo keeps the stored rate and counts it as skipped (skipped_newer). A currency with no stored rate date always takes the fetched rate.

## Fetch requires a currency code and uses latest usable observation

Lines 112-173

Asked as: currency code required error; empty observation rate; which bank observation used; rate fetch no value found

Fetching a rate fails with "currency code is required" if the code is blank. Otherwise Turbo goes through the observations in the order the Bank gives them and takes the first one that has a non-empty value for the currency's series. It uses that value as the rate and that observation's date as the rate date. If no observation has a value, it returns no rate and no error.
