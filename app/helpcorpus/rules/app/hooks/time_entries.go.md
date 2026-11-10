<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/time_entries.go","sha256":"8397c98015740d495292795b0182203da13dc663ae6b453c2d74e2a84570ee0c","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/time_entries.go

## Time type decides which time entry fields are kept

Lines 30-78

Turbo clears every field on a time entry that its time type does not allow (allowed_fields). The fields id, uid, created, updated, role and branch are always kept.

## Required time entry fields come from the time type

Lines 51-62

The time type's required_fields list sets which fields a time entry must have. When required_fields is null, [] or {}, every allowed field is required.

## Time entry role is cleared when no job is set

Lines 80-83

Turbo clears the role on a time entry that has no job.

## Time entry branch follows the job's branch

Lines 88-110

When a time entry has a job, Turbo always sets its branch to the job's branch, replacing any branch given. If the job is not found, Turbo refuses with 'referenced job not found'. If the job has no branch, Turbo refuses with 'referenced job is missing a branch'.

## Time entry without a job uses the default branch

Lines 111-124

When a time entry has no job, Turbo sets its branch to the creator's default branch (default_branch) only if the branch is blank. A branch already set is kept. If the creator's admin profile has no default branch, Turbo refuses with 'your admin_profiles record is missing a default_branch'.

## Total hours limit of 18 per time entry

Lines 137-158

Hours plus meals hours on one time entry may not exceed 18. Exactly 18 is allowed. Otherwise Turbo shows 'Total hours must not exceed 18'.

## Missing required time entry fields are reported first

Lines 143-151

Turbo first checks that every required field has a value, and each missing field gets 'Value required'. Turbo runs the other checks only when no required field is missing.

## Time entry hours must be a positive multiple of 0.5

Lines 156-156

Hours on a time entry must be a positive multiple of 0.5 (utilities.IsPositiveMultipleOfPointFive).

## Time entry date must be valid

Lines 157-157

The date on a time entry must pass the date check (utilities.IsValidDate).

## Meals hours limit of 3 per time entry

Lines 159-159

Meals hours on a time entry may not exceed 3. Exactly 3 is allowed. Otherwise Turbo shows 'Meals Hours must not exceed 3'.

## Time entry description needs at least 5 characters

Lines 160-160

A time entry description must have at least 5 characters. Otherwise Turbo shows 'must be at least 5 characters'.

## Work record format depends on whether a job is set

Lines 161-161

When a time entry has a job, its work record must match ^[FKQ][0-9]{2}-[0-9]{3,4}(-[0-9]+)?$. Otherwise Turbo shows 'must be in the correct format'. When there is no job, the work record must be empty, or Turbo shows 'Work Record must be empty when job is not provided'.

## Payout request amount must be positive with whole cents

Lines 162-162

A payout request amount (payout_request_amount) must be greater than 0, or Turbo shows 'Amount must be greater than 0'. It must also be a positive multiple of 0.01 (utilities.IsPositiveMultipleOfPointZeroOne).

## Job on a time entry must exist and allow time tracking

Lines 166-174

When a time entry has a job, Turbo checks that the job exists, or it shows 'invalid job reference'. Turbo also checks that the job allows time tracking (validateJobAllowsTimeTracking). The rules of that check are in another file.

## Role is required when a time entry has a job

Lines 175-178

A time entry with a job and no role is refused with 'Role is required when a job is assigned'.

## Division must be allocated to the job

Lines 179-184

When a time entry has a job, its division must be allocated to that job (validateDivisionAllocatedToJob). The error codes and messages come from that shared helper, which is in another file.

## Time editing must be enabled to save time entries

Lines 201-207

Turbo refuses to create or update a time entry when time editing is disabled (utilities.ErrTimeEditingDisabled). The message is defined in another file.

## Time entry must belong to the signed-in user

Lines 209-213

The uid on a time entry must equal the signed-in user's id. Otherwise Turbo refuses with 'uid property must be equal to the authenticated user's id'.

## Time entry division must be active

Lines 226-231

Turbo checks that the division on a time entry is active (EnsureActiveDivision). A failure is reported as a 'Validation error'. The rules of that check are in another file.

## User must be allowed to use the time entry branch

Lines 233-235

Turbo checks that the user may use the branch of the time entry (utilities.EnsureUserCanUseRecordBranch). The rules of that check are in another file.

## Week ending is calculated from the time entry date

Lines 237-243

Turbo sets the week ending (week_ending) of a time entry from its date only. If it cannot be calculated, Turbo refuses with 'Error generating week_ending'.
