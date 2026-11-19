<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/findRates.ts","sha256":"6ba4ce98a62c51f27f5b7edea4f9bc08d81ddcdb2df431e39dd7bab26f928994","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/findRates.ts

## Find by rates: what counts as a row entered

Lines 9-43

Asked as: what counts as entered; row not counted; empty rate field; blank row search

In Find by rates on a project request, a row counts as entered when it has a rate or an overtime rate. An empty input (null, undefined or "") counts as nothing entered.

## Find by rates: overtime rate limits

Lines 13-77

Asked as: overtime rate too low; overtime cents not allowed; overtime decimal places error; overtime rate under 1

An overtime rate below 1 gets the error "Overtime rate must be at least 1." An overtime rate with more than two decimal places gets "Overtime rate must have at most two decimal places." Unlike the rate, the overtime rate may have cents. Turbo compares it in whole cents with a tolerance of 1e-6, so 1.15 passes. A row with a rate and no overtime rate has no overtime check.

## Find by rates: rows without a rate are left out

Lines 20-27

Asked as: search without overtime rate; rate only search; overtime left out of search; why rows ignored in search

Turbo searches only rows that have a rate. Each rate is sent with its overtime rate only when the requester entered one. Without an overtime rate, a sheet matches on the rate alone.

## Find by rates: result line

Lines 29-38

Asked as: no sheets found; how many sheets match; search returned nothing; sheet count message

With no match, Turbo says "No active sheet has these rates." With matches it says "1 sheet matches" or "N sheets match". Only active sheets are searched, and the newest code comes first.

## Find by rates: overtime rate needs the role's rate

Lines 53-60

Asked as: overtime rate but no rate; overtime without rate error; clear overtime rate; search stuck on rate field

A row with an overtime rate and no rate gets the error "Enter the rate, or clear the overtime rate." on its rate field. The search does not run until the requester fixes it.

## Find by rates: rate must be at least 1

Lines 61-63

Asked as: rate too small; rate under 1; rate not a number; invalid rate value

A rate below 1 gets the error "Rate must be at least 1." A value that is not a number gets the same error.

## Find by rates: rate must be a whole number

Lines 64-66

Asked as: rate has decimals; rate with cents; fractional rate error; rate 25.5 rejected

A rate of 1 or more that is not a whole number, such as 25.5, gets the error "Rate must be a whole number."

## Find by rates: when no search runs

Lines 106-113

Asked as: search not running; find button does nothing; results disappeared after edit; no search happening

Turbo runs no search when no rate is entered, or when any entered rate has an error. A search result shows only while the rates are the same as the ones it was made for.

## Carrying found rates into proposed rates

Lines 121-149

Asked as: copy found rates to proposed; use found rates; apply search rates to proposed; found rates not updating rows

When found rates are carried into proposed rates, each one replaces its role's rate. Its overtime rate replaces the row's only when one was entered, so otherwise the row keeps its overtime rate. A role with no row gets a new row. Other rows stay unchanged, and rows are sorted by role name. With no carried rates, the rows do not change.

## Find by rates: status line text

Lines 151-171

Asked as: status message meaning; why status says enter rates; searching stuck message; search failed message

The status line shows, in this order of priority: "Enter the rates you know." when no rate is entered; "Correct the rates marked below." when a rate has an error; "Turbo could not search the rate sheets." when the search failed; "Searching…" until the count arrives; then the result line.
