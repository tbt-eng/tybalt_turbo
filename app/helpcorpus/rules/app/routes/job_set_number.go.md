<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_set_number.go","sha256":"33fc38f3b6cd53bd7e5ed8c6903fd31058b988a14e1efb3a9c772bdf6a7e5853","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_set_number.go

## What counts as a child-form job number

Lines 30-33

Asked as: child form job number; what does 25-123-01 mean; job number format; legacy P prefix number

Turbo treats a job number as child-form when it matches an optional leading P, 2 digits, a dash, 3 or 4 digits, a dash, then 1 or 2 digits, such as 25-123-01. The optional P is kept for old legacy rows, even though new proposal subjobs are not allowed. Turbo uses this pattern (jobChildNumberRegex) when a job is renumbered.

## Child jobs keep a child-form number

Lines 33-35

Asked as: subjob wrong number format; child job number pattern; subjob must keep number; child job renumber rejected

A child-form number looks like 25-123-01, with an optional leading P: two digits, 3 or 4 digits, then 1 or 2 digits. A job with a parent (a subjob) must keep a child-form number: "child jobs must keep a child-form number".

## Job number errors are reported on the number field

Lines 35-46

Asked as: validation failed on number; number field error; job number error code; why won't job number save

When Turbo refuses a job number change for a reason about the number, it answers "validation failed". The answer names the "number" field and gives an error code and a message for it.

## Number field rule failures keep the validator's message

Lines 48-59

Asked as: validator message job number; number field validation_error; job number rule message; number rejected by field rules

When the new job number fails the number field's own rules, Turbo shows the validator's message unchanged. It uses the validator's code, or "validation_error" when the validator gives no code.

## Changing a job number requires sign-in

Lines 71-74

Asked as: not signed in job number; change number logged out; job number unauthorized; need to log in to renumber

Turbo refuses the job number change when the person is not signed in. The message is "unauthorized".

## Job number changes blocked when job editing is off

Lines 76-82

Asked as: job editing disabled; can't edit job number; jobs editing turned off; editing locked for jobs

Turbo checks whether job editing is enabled before it changes a job number. When editing is disabled, it refuses with the error ErrJobsEditingDisabled (utilities).

## Surrounding spaces are trimmed from the new number

Lines 94-94

Asked as: extra spaces in job number; space around job number; job number whitespace; trailing space job number

Turbo trims leading and trailing spaces from the requested number. It also trims the job's current number before it compares them.

## Only active admins can change a job number

Lines 99-121

Asked as: not an admin job number; who can change job number; non admin renumber job; admin rights to renumber

Only an active account that holds the admin claim can change a job's number. Anyone else gets "you are not authorized to change job numbers".

## Job must exist to change its number

Lines 123-130

Asked as: job not found number; missing job id; can't find job to renumber; job does not exist

Turbo answers "job not found" when no job has the given id.

## Renumbering fails if the jobs number field setup is broken

Lines 132-148

Asked as: jobs collection failed to load; number field not text; jobs number setup broken; renumber configuration error

When an admin changes a job number, Turbo loads the jobs collection and its number field. If the collection cannot be loaded, it refuses with 'failed to load jobs collection'. If the number field is not a text field, it refuses with 'jobs.number field is not configured as text'.

## Setting the same number changes nothing

Lines 150-155

Asked as: same job number no change; number unchanged; job number already set; saving identical number

When the requested number equals the job's current number, Turbo makes no change and returns the current number as a success. It runs none of the other checks.

## New number must pass the number field rules

Lines 157-159

Asked as: job number pattern rejected; number field rules failed; job number format invalid; number does not match rules

Turbo validates the new number against the jobs number field's own rules (such as its pattern). When it fails, the error appears on the number field with the validation message.

## A proposal keeps a proposal-form number

Lines 161-165

Asked as: proposal number missing P; proposal must start with P; proposal number rejected; proposal renumber without P

A number that starts with P is proposal-form. A proposal cannot get a number that lacks the P: "proposals must keep a proposal-form number".

## A project cannot take a proposal-form number

Lines 161-168

Asked as: project with P number; project number starts with P; can't give project proposal number; project number not allowed

A project cannot get a number that starts with P. Turbo answers "projects must not use a proposal-form number".

## Top-level jobs cannot take a child-form number

Lines 170-177

Asked as: top level job child number; parentless job subjob number; standalone job child number; job without parent number format

A job without a parent cannot get a number in child form. Turbo answers "top-level jobs cannot use a child-form number".

## Renumbering skips the normal job save checks

Lines 179-194

Asked as: renumber skips save checks; job number bypasses hooks; forced number change cancelled proposal; renumber imported flag

Turbo writes the new number directly to the job. It does not run the usual job save hooks or validations, such as the cancelled-proposal edit lock or checks on legacy blank fields. It sets the job's _imported flag to false and updates its updated time, so writeback sees the change.

## Job numbers must be unique

Lines 195-197

Asked as: duplicate job number; number already taken; same number as another job; job number in use

When another job already has the requested number, Turbo refuses with "job number must be unique".

## Successful renumber returns job id and number

Lines 205-227

Asked as: renumber success response; what does job number change return; new job number confirmation; job id after renumber

After a successful change, Turbo returns the job id and its new number.
