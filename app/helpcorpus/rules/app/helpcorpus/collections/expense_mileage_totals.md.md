<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_mileage_totals.md","sha256":"a683e8f570e172050aab6fda8ad70c4b9b94d5a0cd826ff079f188cffb38c2cd","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_mileage_totals.md

## Nobody can list mileage totals through the API

Lines 7-9

Asked as: list mileage totals; mileage totals API list; fetch all mileage records; why can't I list mileage

The expense_mileage_totals view has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code is not limited by this.

## Nobody can view one mileage total through the API

Lines 11-13

Asked as: view one mileage total; mileage record forbidden; get single mileage total; mileage total access denied

The expense_mileage_totals view has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code is not limited by this.

## Mileage rates come from the latest rate table in effect

Lines 19-31

Asked as: which mileage rate applies; wrong per km rate; rate table effective date; mileage rate changed

Each mileage expense uses the rate table (expense_rates) with the latest effective date on or before the expense date. The table's mileage setting gives tiers: each key is the lower distance of a tier and each value is the rate. A tier ends where the next tier begins.

## Mileage totals cover only expenses paid by Mileage

Lines 34-62

Asked as: mileage missing from totals; cash expense not in mileage; why expense not in totals; payment type mileage filter

The mileage totals view reads only expenses whose payment type (payment_type) is 'Mileage'. Other expenses do not appear in it.

## Cumulative distance restarts at each mileage reset date

Lines 39-54

Asked as: mileage counter reset; distance total restarts; yearly mileage reset; running km total wrong

For each person (uid), Turbo adds up mileage distance in date order. The running total (cumulative) restarts at the latest reset date (mileage_reset_dates) on or before the expense date. An expense dated on a reset date belongs to the new period.

## An expense that crosses a tier boundary is split across tiers

Lines 72-109

Asked as: trip spans two rates; split mileage between tiers; km charged at two rates; tier boundary mileage

An expense covers the span from its cumulative distance minus its own distance, up to its cumulative distance. Turbo prices each part of that span at the rate of the tier it falls in. A tier counts only when the expense overlaps it. The top tier has no upper limit and is treated as 1e9 km.

## Mileage total sums distance times rate, rounded to 2 decimals

Lines 111-128

Asked as: mileage amount calculation; why mileage total is zero; km times rate amount; mileage rounding cents

Each mileage expense's total (mileage_total) is the sum of the kilometres in each tier times that tier's rate. Turbo rounds the sum to 2 decimals. It is 0 when no tier matches, for example when no rate table applies.
