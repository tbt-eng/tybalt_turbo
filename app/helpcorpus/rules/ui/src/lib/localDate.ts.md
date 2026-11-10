<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/localDate.ts","sha256":"d9a09c5343f28c17e0781b6c7ca4b8c28b56106259328148bd23c011e703245d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/localDate.ts

## Today's date comes from the user's device, not UTC

Lines 1-9

Asked as: wrong date in the evening; date shows tomorrow; today's date is off; timezone date wrong; date rolls over too early

Turbo works out "today" (localIsoDate) from the date on the user's own device, in the device's local timezone. It does not use UTC. In the evening, anywhere west of UTC, the date stays on the local day and does not roll over to the next day.

## Calendar dates are written as YYYY-MM-DD with zero padding

Lines 4-8

Asked as: date format wrong; date shows single digit month; leading zero in date; YYYY-MM-DD format; date has no timezone

Turbo writes a date as year-month-day (YYYY-MM-DD) and pads the month and the day to two digits. The date carries no timezone. When no date is given, Turbo uses the current date on the device.
