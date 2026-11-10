<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_off.md","sha256":"a9eae2a32686f1d413b5c03806e0dc100f579441328500ce5b07b2f7f644d9d0","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_off.md

## Time off view is read-only

Lines 5-5

Asked as: edit time off; can't change time off; time off read only; update time off record

The time off records (time_off) are a saved query that Turbo reads and never writes. No create, update or delete rule is listed for them, so nobody can change them through the API.

## Who can list time off records

Lines 9-13

Asked as: see everyone's time off; list time off records; who can see time off; time off list access

A person can list the time off record whose id is their own. A manager (manager_uid) can list the records of the staff they manage. A person with the report claim can list all records.

## Who can view one time off record

Lines 17-21

Asked as: open time off record; view my time off; see one time off; time off details access

A person can view their own time off record. A manager (manager_uid) can view the record of each person they manage. A person with the report claim can view any record. The view rule is the same as the list rule.

## Who appears in time off records

Lines 26-46

Asked as: missing time off record; no time off for me; why no time off entry; who gets time off

Turbo makes a time off record only for staff whose admin profile has untracked_time_off false and time_sheet_expected true. Staff outside this group have no record. The record also shows the person's manager, if there is one.

## Used OV and OP hours are totals

Lines 35-36

Asked as: total used holiday hours; used OV total; how much time off used; time off sum

used_ov and used_op are the sums of committed OV and OP hours since the opening date. Each is 0 when there are none. These totals include time amendments as well as time sheet entries.

## Time sheet hours are separate from amendments

Lines 37-38

Asked as: timesheet OV hours; OV from timesheets only; amendment hours excluded; time sheet time off count

timesheet_ov and timesheet_op count only the OV and OP hours that come from a time sheet (tsid is not empty). Amendments that have no time sheet are left out of these two fields.

## Last OV and OP dates

Lines 39-40

Asked as: last time off date; when did I last take OV; last OP day

last_ov and last_op are the latest dates of counted OV and OP time. They are empty when there is no such time.

## Only committed time counts as used time off

Lines 47-72

Asked as: uncommitted hours not counted; time off not showing; draft timesheet time off; unapproved hours missing

Turbo counts only OV and OP hours that are committed. Time entries count when their time sheet is committed. Time amendments count when they are committed. Hours on uncommitted time sheets are not counted.

## Only time after the opening date counts

Lines 73-74

Asked as: opening date balance; time off before opening date; hours on start date; why first week not counted

Turbo counts used time off only when the week ending is strictly after the opening date (opening_date) on the admin profile. Hours in the week ending on the opening date itself are not counted.

## Opening balance field limits

Lines 87-89

Asked as: opening balance too high; max opening OV; opening leave limit; date format opening

The opening OV balance (opening_ov) can be at most 200. The opening OP balance (opening_op) can be at most 332. The opening date (opening_date) must have the form YYYY-MM-DD.
