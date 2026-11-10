<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_valuation.go","sha256":"0798eb3d0c4a314293e44fbcc39f54089acfd2d75dc0c0eadc8705ac9283d990","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_valuation.go

## Job time valuation shows estimated and unpriced hours, with percent blank when incomplete

Lines 11-25

Asked as: job valuation percent blank; why is valuation percent empty; unpriced hours on job; estimated hours missing from export; job time valuation CSV blank

The job time valuation in the API and CSV exports identifies employee-rate estimates (estimated_hours) and unpriced hours (unpriced_hours). The percent (percent) is null when the total is incomplete or zero. The file only declares these fields. The calculation is in the query job_priced_time_entries.sql, and this file does not show its result.
