<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_entries.md","sha256":"db95093fcb400cf62662b4bf55d811986953b7c3f2151a0287f3748d4e6f5ce1","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_entries.md

## Who can list or view time entries

Lines 9-25

Asked as: see someone else's time entry; view colleague's hours; who can view time entries; approver can't see entries; list entries on submitted sheet

A person can list or view a time entry when it is their own entry (uid). The approver of the time sheet (tsid.approver) can too, once the time sheet is submitted. So can a reviewer of the time sheet (time_sheet_reviewers_via_time_sheet), once it is submitted. A committer (the commit claim) can see entries of a time sheet that is submitted, approved and not yet committed. A person with the report claim can see entries of any committed time sheet.

## Who can create a time entry

Lines 29-38

Asked as: cannot add time entry; add hours to job; create entry wrong category; category not allowed for job; time entry refused no job

A signed-in person with the time claim can create a time entry. The category must belong to the job in the request, and the request must set a job. A blank category (category) is always allowed, with or without a job. A category with no job in the request is refused.

## Who can update a time entry

Lines 42-49

Asked as: edit my hours; change entry already submitted; can't update time entry; modify hours on timesheet; change owner of entry

Only the person who owns the entry (uid), and who has the time claim, can update it. The entry must not yet be on a time sheet (tsid is empty). An update cannot change the owner (uid). Once the entry is on a time sheet, nobody can update it through the API.

## Category must match job when updating a time entry

Lines 51-58

Asked as: changing job breaks category; category error after moving job; new job category mismatch; switch job on entry; category not valid for job edit

If the update does not send a job, the new category must belong to the entry's current job. If the update sends a job, the new category must belong to that new job. A blank category is always allowed.

## Who can delete a time entry

Lines 64-66

Asked as: remove my hours; delete entry on timesheet; can't delete time entry; get rid of an entry; erase logged time

Only the person who owns the entry (uid), and who has the time claim, can delete it. The entry must not yet be on a time sheet (tsid is empty). Once it is on a time sheet, nobody can delete it through the API.

## Time entry field limits

Lines 73-84

Asked as: too many hours in a day; hours over 18; meal hours limit; work record format wrong; date rejected when saving entry

A time entry's hours (hours) cannot exceed 18. Its meals hours (meals_hours) cannot exceed 3. The work record (work_record) must look like F, K or Q, then 2 digits, a hyphen, 3 or more digits, and optionally a hyphen and digits (for example F25-123 or K25-123-4). The date (date) and week ending (week_ending) must have the form YYYY-MM-DD. The owner (uid) and time type (time_type) are required.
