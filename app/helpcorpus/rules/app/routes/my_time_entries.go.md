<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/my_time_entries.go","sha256":"764c5f5a7d09ca9cc47ebd59fb2557dcdf3746a92d8e5b8833b509a7371f66c0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/my_time_entries.go

## My time entries CSV columns and file name

Lines 37-108

Asked as: download time entries csv; csv column headings; time entries export file name; what columns in my export

The CSV starts with a header row, then one row per time entry. The columns are date, week_ending, time_type, time_type_name, hours, meals_hours, job_number, job_description, division, division_name, category, branch, role, work_record, payout_request_amount, description and time_sheet_status. The browser downloads the file as my_time_entries.csv.

## My time entries CSV export covers only your own entries

Lines 57-68

Asked as: export someone else's entries; see colleague's time entries; export another person's hours; download other user's timesheet data

The my time entries CSV export returns only the time entries whose owner (uid) is the signed-in caller. It needs no claim, so every signed-in person can export their own entries. Nobody can use it to export another person's entries.

## My time entries CSV includes entries from time sheets in any state

Lines 57-68

Asked as: missing entries in export; export draft timesheets; csv shows unsubmitted hours; why are some entries not exported

The CSV export lists every time entry of the caller, whatever the status of its time sheet. Each row shows the time sheet status (time_sheet_status), so draft, submitted and committed time sheets all appear. The export lists the newest entries first.

## Payout request amount is blank when it is zero

Lines 75-78

Asked as: empty payout amount in csv; payout column shows blank; why zero payout shows nothing; payout request cell empty

In the my time entries CSV, the payout request amount (payout_request_amount) appears with 2 decimal places when it is not zero. When it is zero, or the entry has none, the cell is empty instead of 0.00.

## Hours and meals hours are written without padding

Lines 84-85

Asked as: hours show without decimals; csv hours formatting; why 8 instead of 8.00; meals hours formatted wrong

In the my time entries CSV, hours and meals hours (meals_hours) appear with no fixed decimals and no trailing zeros. A value of 8 appears as 8, and 7.5 appears as 7.5. Zero values still appear as 0.
