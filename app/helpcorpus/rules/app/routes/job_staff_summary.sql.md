<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_staff_summary.sql","sha256":"8d6d804ac7ffbcda566afab5548be49d020407174f6999e83e2ed808794dffa5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_staff_summary.sql

## Job staff summary groups time by person, not by role

Lines 1-3

Asked as: same person two roles; employee twice on job; staff listed twice; worked two roles one job

The job staff summary prices each time entry first and then groups the entries by employee (uid). An employee who worked in more than one role on the same job gets one row.

## Rate sheets column lists each sheet and revision used

Lines 4-13

Asked as: which rate sheet used; rate sheet revision column; multiple rate sheets employee; rate sheets listed per person

For each employee, the job staff summary lists the distinct rate sheets that priced their entries, as 'name (rev. N)', sorted and separated by '; '. Entries with a blank rate sheet name are left out. The column is blank when no sheet priced any entry. Time in different periods can show more than one sheet.

## Staff with no profile show blank names

Lines 12-12

Asked as: blank name in staff list; missing employee name; no profile name empty; staff name missing

When an employee has no profile, the job staff summary shows a blank given name and a blank surname (empty strings, not NULL). The row is still listed.

## Staff value and job total are sums of priced values

Lines 14-21

Asked as: staff value total wrong; job total hours sum; value per employee calculation; summary totals add up

Each row's value is the sum of the employee's entry values, or 0 when there are none. The total is the sum of the values of all rows. Hours, meals hours, estimated hours and unpriced hours are also summed for each employee.

## Row value status: Incomplete, Estimated or Rate sheet

Lines 28-29

Asked as: why value says incomplete; estimated status on employee; unpriced hours status; value status meaning

A row's value status (value_status) is 'Incomplete' when the employee has any unpriced hours. Otherwise it is 'Estimated' when the employee has any estimated hours. Otherwise it is 'Rate sheet'.

## Total status uses the whole job's hours

Lines 30-31

Asked as: job total status incomplete; whole job status estimated; total status meaning; job status overall

The total status (total_status) is 'Incomplete' when any employee on the job has unpriced hours. Otherwise it is 'Estimated' when any employee has estimated hours. Otherwise it is 'Rate sheet'.

## Percent of job value shown only when nothing is unpriced

Lines 32-32

Asked as: percent of job missing; why percent blank; percentage not showing; share of job empty

Each person's percent is their value divided by the job total, times 100, rounded to 1 decimal. Turbo leaves the percent empty when any hours on the job are unpriced. It is also empty when the total is 0.

## Job staff summary sort order

Lines 34-34

Asked as: staff list order; sort staff summary; alphabetical staff summary; order of people on job

The job staff summary sorts rows by surname, then given name, then employee ID (uid).
