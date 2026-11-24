<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/my_time_entries.sql","sha256":"999bd4e81bd242535d4c7eb67c41daa52e7faff743e539b085abcbdb100b48e2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/my_time_entries.sql

## Rejected sheet shows 'rejected' even if also approved

Lines 4-29

Asked as: rejected sheet shows approved; timesheet rejected but status wrong; rejected and approved both showing; recalled rejected sheet status

A rejected time sheet shows 'rejected' in the CSV, not 'approved' or 'pending'. Turbo checks rejected before approved and pending. The sheet stays submitted until it is recalled. A committed sheet shows 'committed' even if it also has a rejected date.

## Time entries CSV shows blanks and zeros for missing values

Lines 11-24

Asked as: blank cells in time entries export; zero hours in CSV; empty fields in timesheet download; missing values in export

When an entry has no value for a field, the CSV shows an empty string for the time type code and name, the job number and description, the division code and name, the category, the branch, the role and the work record. It also shows an empty string for the description. Missing hours (hours), meals hours (meals_hours) and payout request amount (payout_request_amount) show as 0. Entries whose linked record is missing still appear in the CSV.

## Entry with no time sheet is 'not submitted' in CSV

Lines 26-31

Asked as: entry shows not submitted; no timesheet status in export; recalled sheet status not submitted; why status says not submitted

If an entry has no time sheet, the CSV shows its time sheet status (time_sheet_status) as 'not submitted'. A recalled sheet also gives 'not submitted'. This is the final fallback when no other status applies.

## Time sheet status in CSV is chosen in a fixed order

Lines 27-30

Asked as: wrong status in csv; pending instead of approved; which status wins timesheet; status order in export

The CSV status comes from the entry's time sheet. Turbo checks in this order: committed, rejected, approved, then pending (submitted). The first match wins. A sheet that is marked submitted but is not committed, rejected or approved shows as 'pending'.

## Time entries CSV lists only the caller's own entries

Lines 41-41

Asked as: someone else's entries in my download; other people's time in export; CSV shows colleagues' hours; download includes other users

The time entries CSV download on the time sheets list page contains every time entry whose user (uid) is the caller. It never includes another person's entries.

## Time entries CSV is sorted newest first

Lines 42-42

Asked as: export order of entries; timesheet csv sorting; newest entries first download; why are rows in wrong order

The CSV orders entries by date, newest first. Entries on the same date are ordered by creation time, newest first, then by record ID, descending.
