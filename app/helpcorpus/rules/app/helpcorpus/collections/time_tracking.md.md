<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_tracking.md","sha256":"83bc7cc4d0900fa25ee4f62c411014913508fa37f70f5267f22694bed67de66f","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_tracking.md

## Time tracking cannot be created, changed or deleted through the API

Lines 5-5

Asked as: cannot edit time tracking; change time tracking record; delete time tracking via API; create time tracking api

The time_tracking collection is a view with no create, update or delete rules. Nobody can create, update or delete its records through the API. Turbo reads it and never writes it.

## Who can list time tracking records

Lines 9-11

Asked as: who can list time tracking; can't see time tracking list; access report claim list; time tracking list refused

Only a person who holds the report claim can list the time tracking records. Everyone else is refused. The filter is `@request.auth.user_claims_via_uid.cid.name ?= 'report'`.

## Who can view one time tracking record

Lines 15-17

Asked as: cannot open time tracking record; view single time tracking; access denied time tracking record; who can see one record

Only a person who holds the report claim can view one time tracking record. Everyone else is refused.

## Placeholder payroll ID expense count for each week

Lines 22-42

Asked as: placeholder payroll id count; expenses with fake payroll number; nine digit payroll id 9; weekly placeholder expense count

For each week, Turbo counts committed expenses whose committed week ending is set and whose owner's payroll ID (payroll_id) is a placeholder. A placeholder is exactly 9 digits and starts with 9. The count matches an expense to the time sheet week by committed week ending (committed_week_ending). A week with no such expenses shows 0.

## Expense placeholder count includes expenses of staff with no admin profile

Lines 27-43

Asked as: expense no admin profile payroll; missing week row expenses only; staff without admin profile count; placeholder count missing staff

An expense whose owner has no admin profile is treated as having an empty payroll ID. It is not counted as a placeholder. Only weeks that have time sheets appear in the view, so a week with only expenses has no row.

## Time tracking has one row per week ending

Lines 33-44

Asked as: one row per week; time tracking week duplicated; order of weeks newest first; why weeks grouped together

The view groups time sheets by week ending (week_ending). It gives one row for each week and lists the newest week first. The row id is the smallest time sheet id in that week.

## Committed count counts committed time sheets

Lines 36-36

Asked as: committed count wrong; committed sheets counted rejected; how many committed timesheets; count of committed time sheets

The committed count (committed_count) for a week is the number of time sheets in that week that have a committed value. A committed sheet counts even if it also has a rejected value.

## Approved count excludes committed and rejected time sheets

Lines 37-37

Asked as: approved count not matching; approved timesheets excluded committed; rejected sheets approved count; why approved number is lower

The approved count (approved_count) for a week counts time sheets that are approved, not committed and not rejected. A sheet that is committed or rejected is not counted.

## Submitted count counts only time sheets waiting for approval

Lines 38-38

Asked as: submitted count waiting approval; pending timesheets count; submitted number wrong; timesheets awaiting approval total

The submitted count (submitted_count) for a week counts time sheets that are submitted and not yet approved, committed or rejected.
