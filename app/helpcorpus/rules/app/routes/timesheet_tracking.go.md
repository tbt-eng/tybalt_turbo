<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_tracking.go","sha256":"c223928c35f31c4e0e33247b3d50ff7c9794dd968892ca494add07d05628fafb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_tracking.go

## Who may view committed time sheet tracking

Lines 27-62

Asked as: can't see time tracking; no access to tracking page; who can open tracking; not authorized time tracking

Time tracking views open to a holder of the report claim, the commit claim or the admin claim. The code checks them in that order, and the first claim found sets the viewer's access level. Anyone else gets the message "you are not authorized to view time tracking".

## Weekly tracking counts: committed count

Lines 73-87

Asked as: committed count wrong; how many committed sheets; committed timesheets number; count of committed

The weekly counts list one row per week ending, newest first. A time sheet counts as committed (committed_count) when its committed value is set. Every viewer sees this count, including an admin.

## Weekly tracking counts: statuses are exclusive, rejected overlaps

Lines 76-83

Asked as: rejected count too high; numbers don't add up; sheet counted twice; rejected overlaps other counts

Approved counts only time sheets approved but not committed. Submitted counts only time sheets submitted but neither approved nor committed. Rejected counts every time sheet with a rejected value, so it can overlap the other counts.

## Admin-only viewers see only committed counts

Lines 79-92

Asked as: admin sees zero counts; approved count showing 0; admin only counts zero; why admin counts empty

A viewer who has the admin claim but neither the report nor the commit claim gets 0 for approved, submitted and rejected counts in the weekly tracking counts. Only committed counts show.

## Tracking list requires a valid week ending date

Lines 136-142

Asked as: week ending error; date format error; missing week ending message; tracking date rejected

The tracking list for a week ending refuses a missing week with "weekEnding is required". It refuses a date not in YYYY-MM-DD format with "invalid date format (YYYY-MM-DD)".

## Tracking list shows blank names when a profile is missing

Lines 148-158

Asked as: name missing in tracking; approver name blank; empty name on timesheet; no name showing

The tracking list shows an empty string for a time sheet's staff given name and surname (given_name, surname) when the profile is missing. It also shows an empty string for the approver, committer and rejector names (approver_name, committer_name, rejector_name) when that person is not set or has no profile. Each name it does show is given name, a space, then surname.

## Tracking list phase of each time sheet

Lines 159-164

Asked as: what status is shown; unsubmitted label meaning; phase column wrong; sheet status label

Each time sheet in the list gets a phase. It is Committed when committed is set. Otherwise it is Approved when approved is set, and Submitted when submitted but not approved. Anything else is Unsubmitted.

## Tracking list totals are zero for sheets without entries

Lines 165-173

Asked as: totals show zero; no hours showing; empty totals tracking; timesheet has no hours

The tracking list shows 0 for every total (hours worked, stat, PPTO, vacation, sick, to bank, bereavement, OT payout request, days off rotation) on a time sheet that has no time entries for that week.

## Tracking list totals per time sheet

Lines 179-195

Asked as: hours don't add up; wrong hours in tracking; overtime payout total wrong; days off rotation count

Each row totals the hours of its time entries for that week by time type code. R and RT give hours worked, OH stat, OP PPTO, OV vacation, OS sick, RB to bank and OB bereavement. OR counts days off rotation as one per entry. OTO sums the overtime payout request amount.

## Tracking list shows submitted time sheets only

Lines 196-197

Asked as: unsubmitted sheet not listed; sheet missing from tracking; draft timesheet not showing; why sheet not in list

The tracking list for a week returns only submitted time sheets for that week ending. Unsubmitted time sheets never appear in it.

## Tracking list visibility depends on viewer's claim

Lines 198-202

Asked as: committer sees fewer sheets; submitted sheets hidden; why I see less; unapproved sheets not visible

A report holder sees all submitted time sheets for the week. An admin-only viewer sees only committed ones. A committer (without the report claim) sees only approved or committed ones, so submitted but unapproved time sheets stay hidden.

## Tracking list sort order

Lines 203-209

Asked as: tracking list order wrong; sheets out of order; sort by name; approved listed first

The tracking list sorts approved time sheets first, then submitted ones, then committed ones. Within each group it sorts by surname and then given name.

## Week ending check accepts any valid YYYY-MM-DD date

Lines 233-236

Asked as: wrong day accepted; non week ending date accepted; date validation too loose; any date allowed

The missing and not-expected lists check only that the week ending parses as a YYYY-MM-DD date (isValidWeekEnding). Turbo does not check that the date is a real week-ending day. A date in the wrong format gets the message 'invalid date format (YYYY-MM-DD)'.

## Missing and not-expected lists need a time tracking viewer

Lines 241-251

Asked as: missing list access denied; can't open missing timesheets; missing list error; missing page permission

The missing and not-expected lists run a time tracking viewer check (requireTimeTrackingViewer, defined elsewhere) first. They refuse a missing week with "weekEnding is required", and a bad date with "invalid date format (YYYY-MM-DD)".

## Missing time sheets list: who is included

Lines 259-267

Asked as: who appears on missing list; missing sheets wrong people; expected users not submitted; exception not applied missing

The missing list for a week ending shows users with no time sheet for that week, whose admin profile marks a time sheet as expected (time_sheet_expected = 1). Users with a missing-time-sheet exception for that week are left out. A user with no admin profile counts as not expected.

## Not-expected list requires viewer check and valid week

Lines 284-294

Asked as: not expected list denied; can't open not expected page; not expected list error; not expected permission

The not-expected list runs the time tracking viewer check (requireTimeTrackingViewer, defined elsewhere). It refuses a missing week with "weekEnding is required" and a bad date with "invalid date format (YYYY-MM-DD)".

## Not-expected list: who is included

Lines 302-309

Asked as: users not expected to submit; who is not expected; no admin profile not expected; not expected people list

The not-expected list shows users with no time sheet for the week whose admin profile does not mark a time sheet as expected. Users with no admin profile are included. Missing-time-sheet exceptions are not applied here.
