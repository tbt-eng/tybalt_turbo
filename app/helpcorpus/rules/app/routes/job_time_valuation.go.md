<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_valuation.go","sha256":"d731d7653e2055718ca30bde37a0c743b538455fad742fc8d24d618dd320b916","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_valuation.go

## Job time valuation percent is blank for incomplete or zero totals

Lines 11-12

Asked as: percent is empty; blank percent on valuation; percent shows null; no percentage in export; valuation percent missing

The job time valuation (JobTimeValuation) leaves the percent empty (null) when the total is incomplete or zero. Turbo does this in both the API and the CSV exports.

## Job time valuation flags estimated hours and unpriced hours

Lines 11-24

Asked as: hours marked as estimated; unpriced hours showing; estimate hours on job; hours with no price; estimated hours total

The job time valuation reports estimated hours (estimated_hours) for entries priced with employee-rate estimates. It also reports unpriced hours (unpriced_hours) for entries with no price. Each is given for the row and for the total (total_estimated_hours, total_unpriced_hours).

## Job time valuation lists the rate sheets that priced the entries

Lines 14-17

Asked as: which rate sheet priced; rate sheet name missing; rate sheets used on job; rate sheet list blank; what price list applied

The rate_sheets field names each rate sheet that priced the row's entries. Each name has the code, the revision and the sheet's name when there is one, such as `2026-11-00004 (rev. 0) "2026 Standard Rates"`. Turbo joins the names with "; ". The field is empty when no rate sheet applied.

## Job time valuation has a value status and a total status

Lines 18-26

Asked as: status of valuation value; valuation total status; value status meaning; why total shows status; valuation status column

Each valuation row has a status for its value (value_status) and a status for its total (total_status), next to the value and total amounts. This file does not say which status values exist.
