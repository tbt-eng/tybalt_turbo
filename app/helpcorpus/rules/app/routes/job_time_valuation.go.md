<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_valuation.go","sha256":"8c92cd742a04ab169af72dac811a3751bfdba1a66d5abaa628aec0377dcbbd82","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_valuation.go

## Job time valuation percent is blank when the total is incomplete or zero

Lines 11-23

Asked as: percent is blank; job valuation percent empty; null percent export; valuation total zero percent; unpriced hours in export

In the API and CSV exports, the percent (Percent) of a job time valuation row is null when the total is incomplete or zero. The row also reports estimated hours (employee-rate estimates) and unpriced hours, each for the row and for the total.

## Job time valuation shows rate sheets used to price entries

Lines 13-16

Asked as: which rate sheet priced; rate sheet on job valuation; job valuation rate sheets empty; what rates were used; valuation rate sheet list

A job time valuation row (JobTimeValuation) lists the rate sheets that priced its time entries (rate_sheets), such as "2026 Standard Rates (rev. 0)", joined by "; ". The list is empty when no rate sheet applied.
