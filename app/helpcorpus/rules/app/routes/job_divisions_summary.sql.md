<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_divisions_summary.sql","sha256":"93b8bf5938c633fc28b2cee83486d27007ec0b5b6d1efa093c6d84022558fc05","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_divisions_summary.sql

## Time is priced per entry before grouping by division

Lines 1-3

Asked as: same job two roles; employee multiple roles rate; time priced per role; role rate on job

The job divisions summary prices each time entry first (priced_entries) and then groups by division. An employee who uses more than one role on the same job is therefore priced at each role's rate.

## Rate sheets column lists each sheet and revision used

Lines 4-13

Asked as: rate sheet revision list; which rate sheets used; sheet name rev number; rate sheets per division

For each division, Turbo lists the distinct rate sheets that priced its entries, in the form 'name (rev. N)', sorted by that label and separated by '; '. Entries with no rate sheet name are left out. A division with no named sheet shows a blank value.

## Missing division shows blank code and name

Lines 12-12

Asked as: entry without division; blank division code; null division summary; missing division name

An entry with no division is still grouped and shown, with a blank division code and a blank division name (COALESCE to ''). Turbo does not show NULL.

## Division value, hours and totals

Lines 14-24

Asked as: division hours total; job value total wrong; estimated hours sum; unpriced hours total

For each division, Turbo sums the hours, the estimated hours and the unpriced hours, and sums the value, which is 0 when there is no value. It also gives job-wide totals of value, estimated hours and unpriced hours across all divisions.

## Value status of a division and of the job

Lines 26-30

Asked as: incomplete status meaning; estimated status shown; value status job; why status incomplete

A division's value status is 'Incomplete' if it has any unpriced hours, otherwise 'Estimated' if it has any estimated hours, otherwise 'Rate sheet'. The total status follows the same order using the job-wide totals.

## Percent of job value shown only when fully priced

Lines 31-31

Asked as: percent of job blank; division percentage missing; percent value empty; why no percent

A division's percent is its value times 100 divided by the job total, rounded to 1 decimal. Turbo leaves the percent empty when any hours on the job are unpriced. It also leaves it empty when the total value is 0.

## Divisions are sorted by code, then id

Lines 33-33

Asked as: division order summary; sort divisions code; divisions listed order; blank code first

The summary lists divisions in order of division code and then division id. Entries with no division have a blank code and come first.
