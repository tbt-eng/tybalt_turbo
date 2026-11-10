<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_rates.md","sha256":"f6d841a32fe3fecc418374b8123f646c68dfde0d38024eeebe92ec685ce96d66","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_rates.md

## Nobody can list expense rates through the API

Lines 7-9

Asked as: list expense rates; can't see expense rates; expense rates API empty; read all rates

The expense_rates collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one expense rate through the API

Lines 11-13

Asked as: view one expense rate; open a single rate; can't open rate record; fetch rate by id

The expense_rates collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create expense rates through the API

Lines 15-17

Asked as: add new expense rate; create rate via API; can't add rate; insert rate record

The expense_rates collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update expense rates through the API

Lines 19-21

Asked as: edit expense rate; change rate amounts; update rate API forbidden; rate update denied

The expense_rates collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete expense rates through the API

Lines 23-25

Asked as: remove expense rate; delete rate record; can't delete rate; drop old rates

The expense_rates collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Expense rate effective date format

Lines 32-45

Asked as: rate date wrong format; effective date invalid; duplicate effective date; date not accepted for rate

Each expense rate needs an effective date (effective_date) in the form YYYY-MM-DD, four digits, two digits, two digits. The field is required. The index on this date is not unique, so two rates can have the same effective date.

## Expense rate required amounts and mileage

Lines 33-37

Asked as: mileage value too large; rate missing lunch amount; breakfast dinner lodging required; mileage json error

Each expense rate requires amounts for breakfast, lunch, dinner and lodging, and a mileage value stored as JSON. The mileage JSON is limited to 2,000,000 in size (maxSize).
