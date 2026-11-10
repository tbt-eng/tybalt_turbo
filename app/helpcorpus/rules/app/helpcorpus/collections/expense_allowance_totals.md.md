<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_allowance_totals.md","sha256":"39930d40b519e664810ca962ddac6721290fd1e81331be9e9e1ae3fb76644ab7","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_allowance_totals.md

## Nobody can list allowance totals through the API

Lines 7-9

Asked as: list allowance totals api; api access allowance totals; get all allowance totals; allowance totals not showing via api

The expense_allowance_totals view has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code is not limited by this rule.

## Nobody can view one allowance total through the API

Lines 11-13

Asked as: view single allowance total; fetch one allowance total api; api forbidden allowance total; cannot read allowance record

The expense_allowance_totals view has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code is not limited by this rule.

## Expense with no rates in effect is left out of the view

Lines 20-40

Asked as: expense missing from allowance view; allowance total blank; expense before rates start; empty rate fields expense

The view joins expense_rates to expenses with a LEFT JOIN. An expense dated before the earliest rate effective date has no matching rate, so its rate fields are empty and its allowance total is empty. The rate fields are marked required, so such a row may fail.

## Allowance total adds the rate of each selected allowance type

Lines 28-32

Asked as: how allowance total is calculated; wrong allowance amount; sum of meal rates; mileage not in allowance total

Turbo calculates allowance_total as the sum of the rates for the allowance types (allowance_types) selected on the expense. Breakfast adds the breakfast rate, Lunch the lunch rate, Dinner the dinner rate and Lodging the lodging rate. A type that is not selected adds 0. Mileage is not part of this total.

## Allowance description lists the selected allowance types

Lines 33-38

Asked as: allowance description wrong; what allowance text shows; allowance label order; trailing space allowance description

Turbo builds allowance_description as the names of the selected allowance types, in the order Breakfast, Lunch, Dinner, Lodging. The names are separated by a space, and the trailing space is trimmed.

## Allowance totals cover only Allowance and Meals expenses

Lines 39-41

Asked as: why expense not in allowance report; expense type missing from allowance totals; only meals and allowance payments; card payment allowance totals

The expense_allowance_totals view includes only expenses whose payment type (payment_type) is 'Allowance' or 'Meals'. Expenses with any other payment type do not appear in it.

## Expense uses the newest rates effective on or before its date

Lines 40-40

Asked as: which rate applies to expense; rate change date expense; old rates used for expense; effective date of meal rates

Turbo takes the rates for an expense from the expense_rates row with the latest effective date that is on or before the expense date (date). A rate that takes effect on the expense date itself applies. A later rate does not apply.
