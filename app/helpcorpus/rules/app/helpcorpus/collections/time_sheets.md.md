<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_sheets.md","sha256":"47d57edfc1e5903f6b3a52206f7523cc03ba641e10cfaeca9ec3a09e1e251e69","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_sheets.md

## Who can list time sheets through the API

Lines 9-15

Asked as: see other people's timesheets; list timesheets API; why can't I see timesheet; timesheet visibility in list

A person sees a time sheet in a list when it is their own (uid). An approver (approver) sees it once it is submitted. A reviewer on the time sheet (time_sheet_reviewers) sees it once it is submitted. A committer (the commit claim) sees it when it is submitted and approved (approved is not empty). A reporter (the report claim) sees it once it is committed (committed is not empty).

## Who can view one time sheet through the API

Lines 19-25

Asked as: open a timesheet via API; can't view timesheet; 403 on timesheet; reviewer can't open timesheet

The view rule is the same as the list rule. The owner (uid) can view the time sheet. The approver and any reviewer on the time sheet can view it once it is submitted. A committer (the commit claim) can view it when it is submitted and approved. A reporter (the report claim) can view it once it is committed.

## Nobody creates, updates or deletes time sheets through the API

Lines 27-37

Asked as: edit timesheet through API; cannot delete timesheet; create timesheet by API; API write access timesheets

The create, update and delete rules are empty. Nobody can create, update or delete a time sheet through the API, except a PocketBase superuser. Turbo's own server code can still do these actions.

## Work week hours must be between 8 and 40

Lines 45-45

Asked as: hours rejected on timesheet; too few work hours; work week hours limit; timesheet hours out of range

A time sheet needs work_week_hours, and the value must be at least 8 and at most 40.

## Week ending must be a date in YYYY-MM-DD form

Lines 47-47

Asked as: week ending date format; wrong date on timesheet; week ending error; date won't save timesheet

A time sheet needs week_ending, a text value that matches the pattern four digits, two digits, two digits (YYYY-MM-DD).

## Payroll ID format on a time sheet

Lines 56-56

Asked as: payroll number invalid; payroll ID rejected; CMS code format; timesheet payroll number error

A time sheet needs a payroll_id. The value is a whole number that does not start with 0, or the letters CMS followed by 1 or 2 digits (such as CMS5).

## One time sheet per person per week

Lines 65-65

Asked as: duplicate timesheet same week; timesheet already exists for week; two timesheets one week; submitted timesheet twice

A unique index on the person (uid) and the week ending (week_ending) stops two time sheets for the same person and the same week.
