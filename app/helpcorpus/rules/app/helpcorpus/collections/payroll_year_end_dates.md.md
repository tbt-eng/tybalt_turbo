<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/payroll_year_end_dates.md","sha256":"fa8a7b4d987c17bd5e1f004f7f3018c731deb63285f1372ba1fe548fea564be2","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/payroll_year_end_dates.md

## Nobody can list payroll year-end date records through the API

Lines 7-9

Asked as: can't list year-end dates; payroll year end list empty; see all year end dates api; year end dates not showing

The payroll_year_end_dates collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view a payroll year-end date record through the API

Lines 11-13

Asked as: can't open year-end date record; view one payroll year end; year end date record forbidden; api read year end date

The payroll_year_end_dates collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create payroll year-end date records through the API

Lines 15-17

Asked as: add new year end date; can't create payroll year end; make year end date api; api says not allowed to add year end

The payroll_year_end_dates collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update payroll year-end date records through the API

Lines 19-21

Asked as: can't edit year end date; change payroll year end via api; update year end date forbidden; year end date won't save changes

The payroll_year_end_dates collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete payroll year-end date records through the API

Lines 23-25

Asked as: remove year end date; can't delete payroll year end; year end record won't delete; api delete year end date denied

The payroll_year_end_dates collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Payroll year-end date must be written as YYYY-MM-DD

Lines 32-32

Asked as: year end date format wrong; payroll date invalid format; year end date rejected; how to type year end date; year end date validation error

The date field of a payroll year-end date record is required. It must match the pattern ^\d{4}-\d{2}-\d{2}$, which is four digits, two digits and two digits separated by hyphens (for example 2025-12-31). The pattern checks the format only, not whether the date is real.
