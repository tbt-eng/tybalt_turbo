<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_amendments.md","sha256":"1538fc1d93218c6d8617158e9b06c27e71bc26ec11eac6086a0b4b45b55b0eb6","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_amendments.md

## Who can list or view time amendments

Lines 7-19

Asked as: who can see time amendments; view amendment records; list amendments permission; cannot open amendment list

Staff with the tame claim or the report claim can list time amendment records and view one record. Other people cannot.

## Who can create a time amendment

Lines 21-31

Asked as: create time amendment; add amendment to timesheet; new amendment fails; tsid set automatically

Only staff with the tame claim can create a time amendment. The request must not send a time sheet (tsid), or must send it blank, because Turbo's hook sets it. The request must not send committed, committer or committed_week_ending.

## Category must belong to the job on create

Lines 33-38

Asked as: category not on job; category does not match job; create amendment category error; job and category mismatch

When a person creates a time amendment, a category must belong to the job in the same request. A request with a category and no job fails. The category can also be blank.

## Who can update a time amendment and when

Lines 43-49

Asked as: edit committed amendment; cannot change amendment; change person on amendment; update amendment not allowed

Only a signed-in person with the tame claim can update a time amendment. The amendment must not be committed (committed is blank). The request must not change the person it is about (uid).

## Update cannot change time sheet or commit fields

Lines 51-57

Asked as: change committed week ending; edit timesheet on amendment; committer field changed; update blocked on commit fields

An update may send tsid, committed, committer or committed_week_ending only if each value equals the stored one. Turbo sets these fields itself, so a person cannot change them through an update.

## Category must match the job on update

Lines 59-66

Asked as: change job on amendment; category wrong after changing job; update category mismatch; moved amendment to new job

When a person updates a time amendment, the new category must belong to the job. If the request sends no job, the category must belong to the stored job. If the request sends a job, the category must belong to that new job. The category can also be blank.

## Who can delete a time amendment

Lines 69-74

Asked as: remove amendment; cannot delete committed amendment; erase time entry; delete amendment button missing

Only staff with the tame claim can delete a time amendment, and only while it is not committed (committed is blank). A committed amendment cannot be deleted through the API.

## Time amendment field limits

Lines 80-102

Asked as: hours out of range; meals hours too high; work record format wrong; amendment date rejected; missing required amendment field

Hours must be between -18 and 18. Meals hours (meals_hours) can be at most 3. The work record must match the pattern like F25-123 or K25-123-1, starting with F, K or Q. The date and week ending must be in the form YYYY-MM-DD. The fields uid, hours, time_type, date, week_ending, creator and branch are required.
