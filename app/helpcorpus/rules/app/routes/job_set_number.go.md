<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_set_number.go","sha256":"b94f89220b73535bcc8335ef2db5fd44561704d8cd2612bcc2a15a02fbf297e7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_set_number.go

## Child and top-level numbers cannot swap forms

Lines 30-169

Asked as: child job number format; top-level job wrong number form; sub job number 25-123-01; child number on parent job

A job with a parent is a child job. A child job must keep a child-form number, such as 25-123-01: "child jobs must keep a child-form number". A top-level job cannot take one: "top-level jobs cannot use a child-form number". The child form is an optional P, two digits, a dash, 3 or 4 digits, a dash, then 1 or 2 digits.

## Changing a job number requires sign-in

Lines 71-74

Asked as: change job number not signed in; renumber job logged out; unauthorized job number change; need to log in to rename job

Turbo refuses a request to change a job number when the person is not signed in. The message is "unauthorized".

## Job number change blocked when job editing is off

Lines 76-82

Asked as: cannot edit job number; job editing turned off; job number locked; renumber blocked editing disabled

Turbo refuses to change a job number while job editing is disabled. It returns the jobs-editing-disabled error (ErrJobsEditingDisabled) from the utilities package. This check runs before the admin check.

## Job must exist to change its number

Lines 84-122

Asked as: job not found renumber; missing job id; wrong job id number change; job doesnt exist to rename

Turbo answers "job not found" when no job has the id in the request. It also refuses a request with no job id, with the message "missing job id".

## Job number is trimmed; same number changes nothing

Lines 94-147

Asked as: spaces around job number; same job number no change; trailing space job number; renumber to same number

Turbo trims spaces from the start and end of the new number. If the trimmed number equals the job's current number, Turbo makes no change and returns the current number as a success. Turbo skips all other checks in this case.

## Successful renumber returns job id and number

Lines 98-219

Asked as: job number changed confirmation; what does renumber return; renamed job result; job number update success

After a successful change, Turbo returns the job id and its number. The renumber and the checks run in one transaction.

## Only admins can change job numbers

Lines 99-113

Asked as: who can change job number; not allowed to rename job; non admin job number edit; need admin to renumber

Only a person with the admin claim can change a job's number. Anyone else gets "you are not authorized to change job numbers".

## New job number must pass the number field's rules

Lines 149-151

Asked as: job number invalid format; job number validation failed; new job number rejected; job number pattern error

Turbo checks the new number against the validation of the jobs number field (the pattern and limits in the jobs collection). If it fails, Turbo refuses the change and shows the validation message on the number field.

## Proposal and project numbers cannot swap forms

Lines 153-160

Asked as: proposal number on project; P number wrong form; project cannot be proposal number; proposal must keep P number

A number that starts with P is a proposal-form number. A proposal must keep a proposal-form number: "proposals must keep a proposal-form number". A project must not take one: "projects must not use a proposal-form number".

## Renumbering skips normal job save checks

Lines 171-186

Asked as: renumber skips save checks; renumber cancelled proposal; job renumber ignores validation; why renumber bypasses lock

Renumbering writes the new number straight to the job and does not run the normal job save hooks or validation. For example, it skips the edit lock on cancelled proposals and the checks on legacy blank fields. It marks the job as locally modified (_imported set to false) and updates its timestamp, so writeback sees the change.

## Job number must be unique

Lines 187-189

Asked as: duplicate job number; number already used; job number taken; same number as another job

If another job already has the new number, Turbo refuses the change with "job number must be unique".
