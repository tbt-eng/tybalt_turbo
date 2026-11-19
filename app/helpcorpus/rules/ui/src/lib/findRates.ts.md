<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/findRates.ts","sha256":"e7929a563cbf2276b33fb902b09a1ba8d90c9eed1aa53c445a562e940a93aac5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/findRates.ts

## Find by rates: what counts as a rate entered in a row

Lines 9-43

Asked as: what counts as a rate entered; empty rate field; rate field blank; overtime rate filled in

A row counts as entered when it has a rate or an overtime rate (overtime_rate). An empty number input (null, undefined or "") counts as nothing entered.

## Find by rates: overtime rate limits

Lines 13-18

Asked as: overtime rate too low; overtime rate decimals; cents on overtime rate; overtime rate under 1

An entered overtime rate below 1 shows "Overtime rate must be at least 1." An overtime rate with more than two decimal places shows "Overtime rate must have at most two decimal places." Unlike the rate, the overtime rate may have cents. The check uses whole cents with a small tolerance, so 1.15 passes. The server checks these rules too.

## Find by rates: search body sends only rates entered

Lines 20-27

Asked as: overtime only row search; overtime rate left out; search ignores overtime; sheet matches on rate only

The search sends one entry for each row that has a rate, in row order. A row with only an overtime rate is left out. Turbo sends the overtime rate only when it was entered, so a sheet then matches on the rate alone.

## Find by rates: result message

Lines 29-38

Asked as: how many sheets matched; no sheets match rates; search results count; why no results found

With no match, Turbo says "No active sheet has these rates." With one match it says "1 sheet matches". With more it says "N sheets match". Results list active sheets only, newest code first.

## Find by rates: overtime rate without a rate is refused

Lines 53-60

Asked as: overtime rate but no rate; overtime rate refused; clear overtime rate message; rate missing for overtime

If a row has an overtime rate but no rate, Find by rates shows the message "Enter the rate, or clear the overtime rate." on the rate field. A row with neither value gives no error.

## Find by rates: rate must be at least 1 and whole

Lines 61-66

Asked as: rate must be whole number; rate below 1; rate with decimals; rate not a number

A rate below 1 shows "Rate must be at least 1." A rate of 1 or more that is not a whole number shows "Rate must be a whole number." A non-numeric rate also fails the at-least-1 check.

## Find by rates: overtime rate is checked even if the rate is wrong

Lines 61-77

Asked as: two errors on one row; overtime error with bad rate; both rate errors showing; overtime checked when rate wrong

When a row has a rate error and also an overtime rate, Turbo still checks the overtime rate and can show both errors. When a row has an overtime rate, Turbo checks it only if the row has a rate.

## Find by rates: no search when nothing to search or errors

Lines 106-113

Asked as: why search did not run; search button does nothing; results disappeared after edit; no search with errors

Turbo runs no search when no rate is entered or when any rate has an error. A result shows only while the rates are the same as the ones it was searched for.

## Find by rates: hint to propose new rates

Lines 121-124

Asked as: need new rates hint; where to propose new rates; what to do after no match; propose rates from results

Under the results, with or without matches, Turbo shows "Need new rates? Choose Propose new rates above. The rates you entered come with you."

## Switching to Propose new rates carries only rows without errors

Lines 126-137

Asked as: switch to propose rows missing; rows not carried over; rows with errors dropped; why row stayed in find

On the switch to Propose new rates, Turbo carries each row that has a rate and no error on its rate or overtime rate. A row with an error is not carried, in part or in whole, and stays in Find by rates. A row with only an overtime rate is not carried.

## Only changed carried rates replace the proposal's rates

Lines 139-148

Asked as: proposal rates overwritten; changes to proposal lost; switching again resets rates; proposal edits kept

On a switch, Turbo compares the carried rates by role with the rates carried at the last switch. Only new roles and roles whose rate or overtime rate changed replace the proposal's rates. A visit that changes no rate keeps the changes made in the proposal since.

## Carried rates fill the proposed rates rows

Lines 150-178

Asked as: carried rates fill proposal; proposed rows not updated; rows sorted by role; new role row added

A carried rate replaces its role's rate. It replaces the overtime rate only when one was entered; otherwise the row keeps its overtime rate. A role with no row gets a new row. Other rows stay unchanged and the rows are sorted by role name. With no carried rates, the rows are unchanged.

## Find by rates: status line text and priority

Lines 192-200

Asked as: status line message; searching stuck message; why status says correct rates; could not search sheets

The status line shows the first that applies: "Enter the rates you know." (no rate entered), "Correct the rates marked below." (error), "Turbo could not search the rate sheets." (search failed), "Searching…" (no answer yet), then the result message.
