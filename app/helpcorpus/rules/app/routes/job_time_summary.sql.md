<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_time_summary.sql","sha256":"3700f6c61835e20ed96a3434e81ddd1003aefa51d5cddb0d0ae3cf372a581c9b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_time_summary.sql

## Job time summary totals hours and date range

Lines 7-9

Asked as: total hours for a job; job hours total; first and last entry date; date range of job time

Turbo adds up the hours of the matching entries (total_hours). It gives the earliest and latest entry dates (earliest_entry, latest_entry).

## Job time summary lists distinct branches, divisions, time types and people

Lines 10-21

Asked as: who worked on this job; list of people on job; which branches charged to job; distinct divisions on job time

Turbo lists each distinct branch, division and time type (id and code) found in the matching entries. It also lists each distinct person (uid and full name as given name plus surname). Each value appears once, however many entries use it.

## Entries with no category show as 'No Category'

Lines 22-24

Asked as: uncategorised time entries; entries without category; blank category in job summary; no category label

In the category list, an entry with no category appears with the id 'none' and the name 'No Category'. Other categories appear with their own id and name.

## Entries still count when lookup records are missing

Lines 25-32

Asked as: missing division still counted; time entry lost after lookup deleted; hours missing from job total; orphaned time entry not showing

The summary joins divisions, time types, profiles, categories and branches with left joins. An entry whose related record is missing is not dropped for that reason. The time sheet join is also a left join, but the committed condition removes entries with no matching time sheet.

## Job time summary counts only committed time sheets

Lines 31-32

Asked as: uncommitted timesheet hours; draft time sheet not in summary; unsubmitted hours missing from job; why are my unsaved hours not counted

The job time summary includes only time entries whose time sheet has a commit date (ts.committed is not empty). Entries on time sheets that are not yet committed are left out of the totals and lists.

## Job time summary ignores entries with zero hours

Lines 33-33

Asked as: zero hour entries; why 0 hour entry not counted; empty time entry ignored; negative hours excluded

The job time summary leaves out time entries whose hours are 0 or less. Only entries with hours greater than 0 count.

## Job time summary is for one job

Lines 34-34

Asked as: summary for a single job; time report per job; job time only this job; which job does the summary cover

The job time summary covers only the time entries recorded against the job whose id is passed in (te.job = id).

## Optional filters by branch, division, time type and person

Lines 35-38

Asked as: filter job time by branch; narrow summary to person; limit job hours by division; filter time type on job report

The summary can be narrowed by branch, division, time type or person (uid). A filter that is empty or not given is ignored. A filter that is set keeps only entries that match it exactly.

## Category filter treats 'none' as uncategorized entries

Lines 39-45

Asked as: filter by no category; only uncategorized hours; category filter blank; show entries without a category

An empty or missing category filter applies no category filter. The value 'none' keeps only entries with no category. Any other value keeps only entries in that category.
