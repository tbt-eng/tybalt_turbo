<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_divisions_summary.sql","sha256":"a07b05a02e7e6d9c6fca3dd7ad964670ad56dc0c0f4bb4d7aafe6f06b98577dd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_divisions_summary.sql

## Job time is priced per entry before grouping by division

Lines 1-3

Asked as: price per entry; different rates same job; time entries priced before grouping; multiple roles one job; rate per role on job

Turbo prices each time entry before it groups the entries by division. An employee can use more than one role on the same job, so the entries can have different rates. This query is appended to job_priced_time_entries.sql.

## Rate sheets listed for each division with code, revision and name

Lines 4-13

Asked as: rate sheets listed per division; which rate sheet applied; rate sheet revision shown; division rate sheet labels; blank rate sheet list

Each division row lists the rate sheets that priced its entries (rate_sheets). Turbo writes a sheet as its code and '(rev. N)', and adds the name in quotes only when the sheet has a name. Each sheet is listed once, in label order, and the labels are separated by '; '. Entries with no rate sheet code add no label, and a division with no sheet shows a blank list.

## Entries with no division show blank division name and code

Lines 15-15

Asked as: division shows blank; entries without division; NULL division name; missing division code; unassigned time entries row

When a time entry has no division, or its division is not found, Turbo shows a blank division code (division_code) and a blank division name (division_name). Turbo groups these entries in one row. Turbo does not show NULL in these fields.
