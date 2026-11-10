<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/mileage_reset_dates.md","sha256":"964bd8d13cca70dd3b49ef1db3afee83954a23c21f0be7cc8f4c97c481b15eb7","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/mileage_reset_dates.md

## Nobody can list or view mileage reset dates through the API

Lines 7-13

Asked as: can't see mileage reset dates; list mileage reset dates API; view reset dates; mileage reset API access denied

The mileage_reset_dates collection has no list rule and no view rule. Nobody can list or view its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create, update or delete mileage reset dates through the API

Lines 15-25

Asked as: can't edit mileage reset dates; change mileage reset date; delete reset date record; add mileage reset date API

The mileage_reset_dates collection has no create rule, no update rule and no delete rule. Nobody can create, change or delete its records through the API, except a PocketBase superuser. Turbo's own server code can still do these actions.

## Mileage reset date must be written as YYYY-MM-DD

Lines 32-32

Asked as: mileage reset date format; reset date won't save; invalid reset date; date wrong format mileage

Each mileage reset date record needs a date (date). The date must match the pattern four digits, a dash, two digits, a dash, two digits, such as 2025-01-01. Turbo refuses a record without a date or with a date in another format.

## Two mileage reset date records cannot share a date

Lines 41-41

Asked as: duplicate mileage reset date; same reset date twice; reset date already exists; two reset dates same day

A unique index on the date field (idx_0o9BK5jvdD) stops two mileage reset date records from having the same date.
