<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/job_rate_sheets.md","sha256":"6cff685784f658202b5e059a3877667fe02e66d2ce8af4a4465aa67eb03d2874","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/job_rate_sheets.md

## Who can list or view job rate sheet records

Lines 7-17

Asked as: who can see job rate sheets; view rate sheets list; not signed in rate sheets; access job rate sheet records

Any signed-in person can list job rate sheet records (job_rate_sheets) and view one record. A person who is not signed in cannot.

## Nobody can create a job rate sheet record through the API

Lines 19-21

Asked as: cannot create rate sheet via API; add new job rate sheet; API create rate sheet fails; only superuser creates rate sheets

No rule allows creating a job rate sheet record (job_rate_sheets) through the API. Only a PocketBase superuser can, and Turbo's own server code can still create one.

## Nobody can update or delete job rate sheet records through the API

Lines 23-29

Asked as: cannot edit rate sheet; delete job rate sheet API; change rate sheet record; update rate sheet through API

No rule allows updating or deleting a job rate sheet record (job_rate_sheets) through the API. Only a PocketBase superuser can, and Turbo's own server code can still do both.

## Job and rate sheet are required on a job rate sheet record

Lines 36-37

Asked as: rate sheet missing job; job rate sheet required fields; deleting job removes rate sheets; rate sheet without a job

Each job rate sheet record (job_rate_sheets) must name one job and one rate sheet (rate_sheet). Turbo deletes a job's rate sheet records when it deletes the job (cascadeDelete).

## Effective date must be a YYYY-MM-DD date or empty

Lines 38-38

Asked as: effective date wrong format; rate sheet date invalid; date format for effective date; empty effective date opening

The effective date (effective_date) of a job rate sheet record must be a date written as YYYY-MM-DD, or it must be empty. An empty date marks the opening rate sheet of the job.

## Reason is limited to 500 characters

Lines 39-39

Asked as: reason too long; rate sheet note character limit; reason field length; how long can reason be

The reason on a job rate sheet record (reason) holds at most 500 characters. It is optional.

## A job cannot have two rate sheets with the same effective date

Lines 49-49

Asked as: two rate sheets same date; duplicate effective date job; job already has rate sheet that date; rate sheet date clash

Two job rate sheet records of one job cannot have the same effective date (unique index on job and effective_date).

## A job can have only one opening rate sheet

Lines 50-50

Asked as: second opening rate sheet; job already has opening rate; two rate sheets no date; blank date rate sheet refused

A job can have only one job rate sheet record with an empty effective date (the opening rate sheet). A second one with an empty date is refused (unique index idx_job_rate_sheets_opening).
