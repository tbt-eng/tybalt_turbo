<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/weekly_time_by_employee.sql","sha256":"f923a8fcee0d4be7c3da8a8c290656c18a4f614e24a1d1a5f78e3ce9ae7ca6e8","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/weekly_time_by_employee.sql

## Report groups hours per employee and lists name and divisions

Lines 5-19

Asked as: hours per employee; employee name column; weekly time by person; divisions column empty

The weekly time by employee report gives one row per employee, grouped by the time entry's user (TE.uid). The Employee column shows 'surname, given name'. The divisions column lists the distinct division codes on that employee's entries, or an empty string when there are none. Rows are sorted by surname.

## Chargeable and non-chargeable hours split by job

Lines 7-8

Asked as: chargeable vs non-chargeable hours; hours without a job; billable hours split; job hours total

The report counts an entry's hours as chargeable (c) when the entry has a job. It counts them as non-chargeable (nc) when the job is empty. Total is c + nc.

## Non-chargeable worked hours are R and RT entries without a job

Lines 9-9

Asked as: non-chargeable worked hours; R and RT hours; worked time without job; nc worked count

The report counts nc_worked as hours with time type code 'R' or 'RT' and no job.

## Unworked hours are OP, OV, OH, OB and OS types

Lines 10-10

Asked as: unworked hours count; OP OV OH OB OS; overtime and holiday hours; nc unworked column

The report counts nc_unworked as hours with time type code 'OP', 'OV', 'OH', 'OB' or 'OS'. This count does not check whether the entry has a job.

## Weekly time report excludes RB time entries

Lines 16-16

Asked as: RB time missing; why is RB not shown; rb hours left out; report missing some time

The weekly time by employee report leaves out every time entry whose time type code (TT.code) is 'RB'. These hours are not in any column of the report.

## Report picks entries by a chosen date column and value

Lines 17-17

Asked as: pick week for report; filter report by date; date column value; which entries are included

The weekly time by employee report keeps only time entries where the chosen date column (date_column) equals the given value (date_column_value). The report does not filter by week in any other way.

## Percent chargeable excludes unworked hours from the base

Lines 26-26

Asked as: percent chargeable wrong; percentChargeable formula; chargeable percentage divide by zero; percent excluding unworked

The report calculates percentChargeable as 100 × c ÷ (c + nc − nc_unworked), rounded to 1 decimal place. Unworked hours are taken out of the divisor. When that divisor is 0, the result is not a normal number, and the SQL code does not guard against it.

## Percent non-chargeable worked uses total hours as the base

Lines 27-27

Asked as: percent non-chargeable worked; percentNCWorked calculation; nc worked percentage; percent of total hours

The report calculates percentNCWorked as 100 × nc_worked ÷ (c + nc), rounded to 1 decimal place. Unworked hours stay in the divisor here, unlike percentChargeable. The code does not guard against a total of 0.
