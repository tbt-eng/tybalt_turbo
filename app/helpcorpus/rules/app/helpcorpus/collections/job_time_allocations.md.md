<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/job_time_allocations.md","sha256":"69e89526af0ee3655c407a8df186d75cdfd28ad4cc2101b4162b4f269fae0005","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/job_time_allocations.md

## Who can list job time allocations

Lines 7-11

Asked as: see all job hours; list job time allocations; who can see allocations; view job hours list

Any signed-in person can list job time allocation records (job_time_allocations). A request without a signed-in user is refused.

## Who can view one job time allocation

Lines 13-17

Asked as: open one allocation; see a single job's hours; view allocation record; can't open allocation

Any signed-in person can view a single job time allocation record (job_time_allocations). A request without a signed-in user is refused.

## Nobody can create job time allocations through the API

Lines 19-21

Asked as: add job hours via api; create allocation through api; api can't insert allocations; make new allocation record

Nobody can create a job time allocation record through the API, except a PocketBase superuser. Turbo's own server code can still create these records.

## Nobody can update job time allocations through the API

Lines 23-25

Asked as: edit allocation hours via api; api update job hours; change allocation through api; allocation edit refused

Nobody can update a job time allocation record through the API, except a PocketBase superuser. Turbo's own server code can still update these records.

## Nobody can delete job time allocations through the API

Lines 27-29

Asked as: remove job allocation api; delete hours record api; api won't delete allocations; erase job time entry

Nobody can delete a job time allocation record through the API, except a PocketBase superuser. Turbo's own server code can still delete these records.

## Job time allocation fields: required values and whole hours

Lines 31-40

Asked as: allocation missing job or division; hours decimal not accepted; fractional hours error; job deleted hours gone

Each allocation needs a job (job) and a division (division), and hours (hours) must be a whole number. Hours is not required. Deleting the job also deletes its time allocations (cascadeDelete).

## Only one time allocation per job and division

Lines 42-46

Asked as: duplicate job division hours; same job twice allocation; allocation already exists for division; two hour rows same job

Two job time allocation records cannot have the same job and the same division. A unique index on job and division enforces this.
