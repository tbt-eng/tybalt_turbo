<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_set_number.go","sha256":"b94f89220b73535bcc8335ef2db5fd44561704d8cd2612bcc2a15a02fbf297e7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_set_number.go

## Child and top-level numbers cannot swap forms

Lines 30-33

A job with a parent (a subjob) must get a child-form number: "child jobs must keep a child-form number". A job without a parent cannot get one: "top-level jobs cannot use a child-form number". Child form is an optional P, 2 digits, a dash, 3 or 4 digits, a dash, then 1 or 2 digits (such as 25-123-01).

## Rejected job number is reported on the number field

Lines 48-59

When Turbo rejects a new job number because the number field's own validation fails, it answers "validation failed" with the error attached to the number field. The error carries the validation message and a code. The code is "validation_error" unless the validation error supplies its own code.

## Duplicate number is recognised from the database error text

Lines 61-67

Turbo treats a failed job number update as a duplicate number when the database error text contains "unique constraint failed" or "not unique", ignoring case. Any other database error is not reported as a duplicate.

## Changing a job number requires being signed in

Lines 71-74

Turbo refuses the job number change from a user who is not signed in. The message is "unauthorized".

## Job number changes are blocked when job editing is off

Lines 76-82

When job editing is disabled (utilities.IsJobsEditingEnabled is false), Turbo refuses to change a job number. It returns the error utilities.ErrJobsEditingDisabled. The message text is defined in another file.

## Number is trimmed and unchanged number is a no-op

Lines 94-147

Turbo trims spaces from the new number and from the current number before comparing. If they are equal, Turbo changes nothing and returns the job's current number as a success.

## Only admins can change a job number

Lines 99-113

Only a user with the admin claim can change a job number. Anyone else gets "you are not authorized to change job numbers".

## Job must exist to change its number

Lines 115-122

If no job has the given id, Turbo answers "job not found".

## New job number must pass the number field's own validation

Lines 149-151

Before any other check, Turbo validates the new number against the jobs number field rules (text field validation). A failure is returned as a validation error on "number" with status 400 and the validator's message. The rules themselves are defined elsewhere.

## Proposal and project numbers cannot swap forms

Lines 153-160

A job number starting with P is a proposal form. A proposal must keep a P number: "proposals must keep a proposal-form number". A project must not take a P number: "projects must not use a proposal-form number".

## Child-form check uses the number's shape, not the parent

Lines 162-169

Turbo decides if the new number is child-form by the pattern above. It decides if the job is a child by whether the job has a parent. The two must agree.

## Renumbering skips the normal job save rules

Lines 171-186

An admin renumber writes the number straight to the database and skips the normal job save hooks and validation (such as cancelled-proposal edit locks). It sets the job as locally modified (_imported = false) and updates the updated time, so writeback sees the change.

## Duplicate job number is refused

Lines 187-189

If another job already has the number, Turbo refuses with "job number must be unique" (code validation_not_unique).

## Successful renumber returns job id and number

Lines 197-219

On success Turbo returns the job id and the number now stored.
