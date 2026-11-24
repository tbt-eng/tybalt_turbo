<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_details.go","sha256":"72d3ce624e13ff91e7edf9e913a1354b503e9ef07ea8d30ed81cb6b190cc9d1c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_details.go

## Time sheet details require sign-in

Lines 23-26

Asked as: not signed in; unauthorized error; session expired; can't open timesheet details; log in required

Turbo refuses a request for time sheet details from a person who is not signed in. The message is "unauthorized".

## Time sheet details with unknown or hidden sheet

Lines 33-44

Asked as: time sheet not found; timesheet not found or not authorized; missing timesheet; can't find my timesheet; sheet does not exist

When the time sheet does not exist, Turbo answers "time sheet not found or not authorized". When the person may not view the sheet, Turbo gives the same message. A person cannot tell a missing sheet from a hidden one.

## Time sheet details response contents

Lines 46-78

Asked as: what details come back; timesheet entries sorted; entries missing job or category; time entries list; which fields show in details

Turbo returns the time sheet, its time entries, the approver info, the shared reviewer count and the owner's name. The entries are sorted by date, newest first. Each entry includes its job, time type, division, category and role.

## Shared reviewer count counts all reviewers

Lines 82-99

Asked as: reviewer count wrong; how many reviewers shown; shared reviewers number; count of reviewers; reviewers on draft timesheet

Turbo counts every shared reviewer record (time_sheet_reviewers) for the time sheet. It counts them whatever the sheet's state.

## Owner name falls back to account name

Lines 101-130

Asked as: owner name blank; name missing on timesheet; shows account name instead; empty owner name; profile name not showing

Turbo shows the owner's name as the given name and surname from their profile. When the profile has no name, Turbo uses the user's account name (users.name). When neither is set, or the user is not found, the name is empty.

## Owner can always view their own time sheet details

Lines 137-139

Asked as: view my own timesheet; can't open my draft; see my own sheet; my own timesheet details; owner cannot open timesheet

The person whose time sheet it is (uid) can always view its details. The sheet's state does not matter.

## Admin views only committed time sheets

Lines 141-147

Asked as: admin can't open timesheet; admin cannot see uncommitted sheet; admin access to someone's timesheet; manager can't view draft; admin view other's timesheet

A person with the admin claim can view the details of another person's time sheet only when it is committed (committed is set). An admin cannot view an uncommitted sheet through this claim.

## Approver views only submitted time sheets

Lines 149-151

Asked as: approver can't see timesheet; approver no access before submit; manager cannot open unsubmitted sheet; approver sees nothing; chosen approver access

The approver chosen on a time sheet can view its details only when the sheet is submitted. Before submission the approver cannot view it.

## Shared reviewer views only submitted time sheets

Lines 153-162

Asked as: reviewer can't open timesheet; shared reviewer no access; colleague can't view sheet; reviewer access before submission; reviewer sees timesheet not found

A person who is a shared reviewer of the time sheet (time_sheet_reviewers) can view its details only when the sheet is submitted.

## Committer views submitted and approved time sheets

Lines 164-170

Asked as: committer can't open timesheet; commit role cannot view; approver not yet approved sheet; committer access unapproved sheet; can't see approved sheet

A committer (the commit claim) can view the details of a time sheet that is submitted and approved (approved is set). A committer cannot view a sheet that is not yet approved.

## Report access covers committed time sheets only

Lines 172-180

Asked as: report user can't see timesheet; reporting access denied; report claim timesheet; reports cannot open draft; access denied for reports

A person with the report claim can view the details of a time sheet only when it is committed. Nobody else gets access, and Turbo denies the request.
