<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_report_week_endings.md","sha256":"1a27d4c0613dc3038b1ca0e3d31e330bc1d1c7321676561260f44fdcda47373e","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_report_week_endings.md

## Week endings view is read-only through the API

Lines 5-29

Asked as: edit week endings; cannot change week ending; api write week ending; delete week ending record

The time_report_week_endings view has no create, update or delete rules, and Turbo only reads it. Nobody can create, change or delete its records through the API. Its data comes from the time sheets.

## Who can list week endings for time reports

Lines 7-11

Asked as: list week endings; no week endings showing; week endings missing from report; report claim list

A person who holds the report claim (name 'report') can list the records in the time_report_week_endings view. Anyone without that claim cannot list them.

## Who can view one week ending record

Lines 13-17

Asked as: open one week ending; view single week ending; cannot see week ending; week ending access denied

A person who holds the report claim (name 'report') can view a single time_report_week_endings record. Anyone without that claim cannot view one.

## Week endings view lists only weeks with committed time sheets

Lines 21-28

Asked as: week missing from list; uncommitted week not showing; unsubmitted timesheet week; week endings order newest

The time_report_week_endings view has one record for each week ending (week_ending) that has at least one committed time sheet (committed is not empty). Weeks that have only uncommitted time sheets do not appear. The newest week ending comes first.

## Week ending record id comes from the earliest time sheet id

Lines 22-27

Asked as: week ending id; where week ending id comes from; duplicate week ending id; week id smallest timesheet

Each week ending record takes its id from the smallest time sheet id (MIN(id)) among the committed time sheets for that week. The id does not identify a separate record of its own.

## Week ending field format

Lines 36-36

Asked as: week ending date format; week_ending invalid date; week ending required field; date pattern week ending

The week_ending field is required and must match the pattern YYYY-MM-DD (four digits, two digits, two digits).
