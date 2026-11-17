<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheet_codes.md","sha256":"bdededcae2677058c6bfac258faacb068575904f22bd1696742af03855728492","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheet_codes.md

## Nobody can list rate sheet code records through the API

Lines 7-9

Asked as: list rate sheet codes via API; cannot see rate sheet codes; fetch all rate sheet codes; API returns empty rate sheet list

The rate_sheet_codes collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one rate sheet code record through the API

Lines 11-13

Asked as: view single rate sheet code; open one rate sheet code record; API forbidden rate sheet code; cannot read a code record

The rate_sheet_codes collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create rate sheet code records through the API

Lines 15-17

Asked as: create rate sheet code API; add new rate sheet code; POST rate sheet code forbidden; API cannot insert rate code

The rate_sheet_codes collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update rate sheet code records through the API

Lines 19-21

Asked as: edit rate sheet code; change rate sheet code via API; update code record forbidden; modify rate code API

The rate_sheet_codes collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete rate sheet code records through the API

Lines 23-25

Asked as: delete rate sheet code; remove rate sheet code API; cannot delete code record; erase rate code entry

The rate_sheet_codes collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Rate sheet code format: year, month and five-digit serial

Lines 32-32

Asked as: rate code wrong format; invalid rate sheet code pattern; code looks like 2026-11-4; rate code format rejected

A rate sheet code (code) must match the pattern YYYY-MM-NNNNN, such as 2026-11-00004. The month part must be 01 to 12. The serial part has exactly five digits.

## Rate sheet code month field format

Lines 33-33

Asked as: rate sheet month format; month field rejected; month not YYYY-MM; invalid month value rate sheet

The month field of a rate sheet code must match YYYY-MM, with a month from 01 to 12.

## Rate sheet code serial must be a whole number from 1 to 99999

Lines 34-34

Asked as: rate sheet serial number invalid; serial zero rejected; serial decimal not allowed; serial over 99999

The serial of a rate sheet code must be a whole number from 1 to 99999. Zero and decimals are not allowed.

## Two rate sheet codes cannot share the same code

Lines 41-41

Asked as: duplicate rate sheet code; same code twice; code already exists rate sheet; rate code taken error

A unique index (idx_rate_sheet_codes_code) stops two records from having the same code.

## Each month allows a serial number only once

Lines 42-42

Asked as: same serial same month; serial already used this month; duplicate month serial; repeated serial number for month

A unique index (idx_rate_sheet_codes_month_serial) stops two records from having the same month and serial together. Each month can use each serial once.
