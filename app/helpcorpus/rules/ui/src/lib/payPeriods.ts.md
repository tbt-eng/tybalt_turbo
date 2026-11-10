<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/payPeriods.ts","sha256":"e87f4e85bff49c81f0b1525a8c905051db942186a23becdeeb7166fe2d217d1a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/payPeriods.ts

## Pay period starts 13 days before its ending Saturday

Lines 7-13

Turbo treats a pay period as two weeks that end on a Saturday. It calculates the start date as the pay period ending date minus 13 days, which is a Sunday (payPeriodStart). The ending date is not checked to be a Saturday, so another date gives a start 13 days earlier.

## Pay date counts as upcoming only if after today

Lines 15-17

Turbo marks a pay date as upcoming only when it is later than today's local date (isUpcomingPayDate). A pay date equal to today is not upcoming.
