<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/timesheet_writeback.go","sha256":"7c2825d23ee497387ac7e40435752aa0aa3b8e938ed9e7db30c2803cf3ccbf30","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/timesheet_writeback.go

## Exported time entries leave out job owner, contact, branch, proposal and status

Lines 14-37

In the legacy export, each time entry (timeEntryExport) does not carry the job owner, client contact, branch code, proposal number or job status. These values appear only in the time sheet's jobs tally. Job, division, client, category, meals hours, payout request amount, work record and description are left out of an entry when empty.

## Hours field name in exported time entries depends on job

Lines 39-63

In the legacy time sheet export, a time entry with a job number reports its hours as "jobHours". An entry without a job reports them as "hours".

## Exported time amendments include creator and committer details

Lines 65-95

Each time amendment in the legacy export carries the amendment's creation time (created), the creator's legacy ID and name, the commit time, the committed week ending (committed_week_ending), and the committer's legacy ID and name. Turbo shows empty text when a creator or committer has no matching profile.

## Exported time amendments are always marked committed

Lines 97-121

Every time amendment in the legacy export carries committed=true. Amendments with a job report hours as "jobHours"; those without a job report them as "hours".

## Work hours tally has three separate totals

Lines 123-127

The work hours tally in an exported time sheet has three fields: jobHours (hours on entries with a job), noJobNumber (hours on entries without a job) and hours (the same no-job hours).

## Job tally records branch, client, owner, contact, status and proposal

Lines 129-137

Each job number in an exported time sheet's job tally (jobsTally) lists the branch, client, job owner, client contact, description, status and proposal number. Turbo leaves out branch, job owner, contact and proposal when they are empty.

## Export identifies the manager by legacy ID and name

Lines 139-150

Each exported time sheet gives its manager (approver) as the manager's legacy uid (managerUid) plus the manager's name (managerName). The employee's payroll ID and legacy uid are also included.

## Time sheet export row fields

Lines 139-169

Each exported time sheet carries the work week hours, the salary flag, and the staff member's legacy ID, names and payroll ID. It also holds the entries, the tallies, the divisions, the job numbers, the time types and the payout request total.

## Export response has time sheets and time amendments

Lines 171-174

The legacy export response has two lists: timeSheets and timeAmendments.

## Lists in the export are sorted alphabetically

Lines 176-184

Divisions, job numbers and time type codes on each exported time sheet are listed in alphabetical order. This keeps the output the same on every run.

## Who may call the legacy time sheet export

Lines 188-208

The export accepts a Bearer token that is a valid unexpired machine token for "legacy_writeback". Otherwise it accepts a user who holds the report claim. Anyone else gets "unauthorized".

## Week ending date is required for export

Lines 210-213

The export refuses a request with no weekEnding value and answers "weekEnding is required".

## Which time sheets the legacy export includes

Lines 215-242

The export lists time sheets for the requested week ending that are committed (committed is not empty) and not imported (_imported = 0). It leaves out sheets whose staff payroll ID matches the generated placeholder condition (utilities.GeneratedPayrollPlaceholderSQLCondition). Staff and manager UIDs are the legacy UIDs from admin profiles.

## Exported time sheets are always approved and locked

Lines 245-250

Every exported time sheet is marked locked, approved and submitted, not rejected, with an empty rejection reason. Only committed sheets are exported.

## Time entries exported with each time sheet

Lines 252-291

Each exported time sheet carries its time entries. Missing job, division, client, category or branch values come out as empty text.

## Tallies by time type for exported time sheets

Lines 303-316

Off-rotation (OR) and off-week (OW) entries count as the number of distinct dates. Overtime payout (OTO) entries sum payout request amounts into payoutRequest. Banked (RB) entries sum hours into bankedHours. Any other time type other than R and RT adds its hours to nonWorkHoursTally under its time type code.

## Work hours tally splits job and no-job hours

Lines 315-338

Only R and RT entries count as work hours. With a job, hours go to jobHours and the job is added to the jobs tally with branch, client, owner, contact, description, status and proposal number. Without a job, hours go to both noJobNumber and hours. A job's hours do not add to the "hours" total. Meals hours sum for R and RT entries only.

## Division tally and sorted lists

Lines 316-318

R and RT entries with a division fill the divisions tally (code to name). Divisions, job numbers and time types are listed sorted.

## Time sheet export totals: banked, meals, payout and non-work hours

Lines 340-352

Each exported time sheet carries banked hours (RB entries), the meals hours total, the payout request total (OTO entries), non-work hours by time type, and the work hours tally. Divisions, job numbers and time types are listed in sorted order.

## Off-rotation and off-week tallies count distinct dates

Lines 346-347

For each exported time sheet, the off-rotation days tally (offRotationDaysTally) counts the distinct dates with an OR entry. The off-week tally (offWeekTally) counts the distinct dates with an OW entry. Two entries on the same date count once.

## Which time amendments the export includes

Lines 355-419

The export lists time amendments whose committed week ending equals the requested week, that are committed and not imported. It leaves out amendments whose staff payroll ID matches the generated placeholder condition. Creator and committer names and UIDs come from profiles, and are empty if missing.
