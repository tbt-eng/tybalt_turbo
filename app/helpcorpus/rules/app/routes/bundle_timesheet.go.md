<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/bundle_timesheet.go","sha256":"884fbfe515c2d8ca39d573a7f1fdce2445d85b9d47ebd13f7a800a5335df286e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/bundle_timesheet.go

## Bundling a time sheet requires time editing and the time claim

Lines 16-22

Asked as: can't bundle timesheet; missing time claim; no permission to bundle; bundle button not working

Before bundling, Turbo runs requireTimeEditing and requireTimeClaim for the signed-in user. If either check fails, Turbo returns that error and bundles nothing. This file does not show what the checks require.

## Week ending date must be valid and a Saturday

Lines 24-32

Asked as: invalid week ending date; week ending not saturday; date format error; wrong date for timesheet

The week ending date must be in YYYY-MM-DD format. Otherwise Turbo shows "Invalid date format. Use YYYY-MM-DD". A date that is not a Saturday gets "Week ending date must be a Saturday".

## Only one time sheet per user per week ending

Lines 59-68

Asked as: timesheet already exists; duplicate timesheet same week; two timesheets same week; already submitted this week

If the user already has a time sheet for that week ending date, Turbo refuses the bundle with "a time sheet already exists for this user and week ending date".

## Bundle collects all the user's time entries for the week

Lines 70-132

Asked as: time entries won't bundle; entries fail validation; bundle stops with error; some hours invalid

Turbo gathers every time entry (time_entries) of the user with that week ending date. It validates them as a group with validateTimeEntries. If validation fails, the bundle stops and the error is returned. The details of the validation are in another file.

## Bundling needs an admin profile

Lines 79-86

Asked as: admin profile missing; error fetching admin profile; profile not found bundling; no admin profile

Turbo reads the user's admin profile for the time sheet values. If the profile is missing, the bundle fails with "error fetching user's admin profile".

## Payroll year end date used for time off balance check

Lines 103-132

Asked as: payroll year end date error; no year end date found; pptO balance check failing; vacation balance bundle error

Turbo uses the latest payroll year end date (payroll_year_end_dates) that is on or before the week ending date, so a date equal to the week ending counts. It passes that date to the time entry validation, which checks PPTO and vacation opening balances. If no such date exists, the bundle fails with "no payroll year end dates found for the week ending" and the date. If the date is not a Saturday, it fails with "payroll year end date ... is not a Saturday, contact support".

## Manager must be an active user to approve

Lines 140-165

Asked as: manager not active; approver inactive; manager left company; can't send to manager

The approver is the manager set on the user's profile. If the manager is not an active user, Turbo refuses with "the approver (your manager) is not an active user" (code approver_not_active).

## Manager must hold the tapr claim

Lines 167-178

Asked as: manager lacks tapr; approver not qualified; manager can't approve timesheet; tapr approver error

If the manager lacks the tapr claim, Turbo refuses the bundle with "the approver must have the tapr claim" (code unqualified_approver).

## New time sheet is created submitted with manager as approver

Lines 182-187

Asked as: timesheet already submitted; who approves my timesheet; timesheet sent to manager; submitted automatically

A successful bundle creates a time sheet for the user and week ending. It sets the approver to the user's manager and marks it submitted (submitted = true).

## Self-managed user's time sheet is approved automatically

Lines 189-194

Asked as: own manager timesheet; auto approved timesheet; approve my own timesheet; self managed approval

When the user is their own manager, Turbo sets the approved time on the time sheet to now at bundle time. The tapr claim check already guarantees that such a user holds the tapr claim.

## Time sheet copies hours, salary and payroll ID from admin profile

Lines 196-199

Asked as: hours copied from profile; salary missing on timesheet; wrong work week hours; payroll id blank timesheet

The new time sheet takes work_week_hours, salary and payroll_id from the user's admin profile at bundle time.

## Skip minimum time check resets after the next bundle

Lines 205-212

Asked as: skip minimum time check; min time check keeps skipping; skip check reset; minimum hours override once

If the admin profile's skip_min_time_check is "on_next_bundle", Turbo sets it to "no" once the time sheet is saved. The skip applies to one bundle only.

## Time entries are linked to the new time sheet

Lines 214-223

Asked as: entries not linked to timesheet; time entries missing timesheet id; bundle didn't save entries; entries not attached

After creating the time sheet, Turbo writes its ID (tsid) to every time entry of that user and week. All the changes happen in one transaction.

## Bundle responses

Lines 225-237

Asked as: what does bundle message mean; bundle error message; success message timesheet; timesheet processed message

Validation failures return the validation error, and approver failures return a message and a code. Other failures return the error text. On success Turbo shows "Time sheet processed successfully".
