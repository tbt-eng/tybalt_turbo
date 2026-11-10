<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/jobs.md","sha256":"560997b35cc4e07f82191df5ba776457fdb0b3caeba4599f7e8ac0daefc5823f","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/jobs.md

## Who can list and view jobs

Lines 9-17

Asked as: view job list; who can see jobs; jobs visible when logged out; open a single job

Any signed-in person can list jobs and view a single job. Anyone who is not signed in cannot.

## Who can create a job

Lines 21-24

Asked as: create new job; cannot add job; job creation permission; make a new job api

Only a signed-in person with the procurement claim (the job claim) can create a job through the API.

## Who can update a job

Lines 29-29

Asked as: edit job not allowed; proposal manager edit; update job without claim; alternate manager editing

A person can update a job only if their admin profile is active (active = true). They must also have the job claim (procurement), or the job must be a proposal (number starts with P) and they must be its manager (manager) or alternate manager (alternate_manager). So a proposal's managers can edit it without the job claim.

## Job number cannot be changed

Lines 29-29

Asked as: change job number; renumber job; job number won't update; rename job code

An update request must not change the job number (number). No one can change it through the API, whatever their role.

## Client and project value are locked on non-proposals

Lines 29-29

Asked as: client locked after save; project value won't change; edit client on job; changing value of awarded job

For a job whose number does not start with P, an update must not change the client or the project value (project_value). For a proposal, these two fields can change.

## Who can delete a job

Lines 35-36

Asked as: remove a job; job delete button missing; cannot delete job; get rid of job

Only a signed-in admin (the admin claim) can delete a job through the API.

## Job with linked records cannot be deleted

Lines 38-45

Asked as: delete job blocked; job has timesheets; job still in use; delete fails linked records

A job cannot be deleted while time entries, purchase orders, expenses or job attachments still refer to it.

## Job number format

Lines 53-53

Asked as: invalid job number; job code rejected; wrong number format; duplicate job number

A job number must match this pattern: an optional P, two digits, a hyphen, 3 or 4 digits, an optional L, then an optional suffix of 1 or 2 digits and an optional suffix of 1 digit, each after a hyphen. Examples are 25-123, P25-1234 and 25-123-01. Job numbers must be unique.

## Required job fields and minimum description

Lines 54-69

Asked as: job won't save required; missing job description; job needs branch; description too short

A job needs a description of at least 3 characters, a client and a branch.

## Job field formats and limits

Lines 60-81

Asked as: job date format error; location plus code; client po too long; job status options; authorizing document type

Job dates must be written YYYY-MM-DD. The location must be a plus code. The client PO (client_po) holds at most 64 characters. The proposal value and project value must be whole numbers. Status must be one of Active, Closed, Cancelled, Awarded, Not Awarded, Submitted, In Progress or No Bid. The authorizing document must be Unauthorized, PO or PA.
