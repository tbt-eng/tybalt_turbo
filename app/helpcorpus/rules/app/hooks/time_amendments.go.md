<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/time_amendments.go","sha256":"968ff306518afab58fedc338ece71ff6c57f93b37106adce45ea0100d5623975","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/time_amendments.go

## Time type decides which time amendment fields are kept

Lines 27-77

Before saving, Turbo looks up the time type (time_type) of the time amendment and clears every field that the type's allowed fields (allowed_fields) do not list. Turbo always keeps id, uid, created, creator, updated, skip_tsid_check and branch. If the time type does not exist, the save fails.

## Empty required fields list means all allowed fields required

Lines 52-59

If the time type's required fields (required_fields) are null, [] or {}, Turbo requires every allowed field. Otherwise Turbo requires only the fields listed in required_fields.

## Work record format depends on whether a job is set

Lines 86-86

When a time amendment has a job, its work record (work_record) must match ^[FKQ][0-9]{2}-[0-9]{3,4}(-[0-9]+)?$. Otherwise Turbo says "must be in the correct format". When there is no job, the work record must be empty. Otherwise Turbo says "Work Record must be empty when job is not provided".

## Total hours capped at 18

Lines 87-87

Hours plus meals hours on one time amendment must not exceed 18. Exactly 18 is allowed. Otherwise Turbo says "Total hours must not exceed 18".

## Required time amendment fields must have a value

Lines 93-101

Turbo checks the required fields first. Each required field that is empty gets the error "Value required". Turbo reports these errors and does not run the other checks until they are fixed.

## Hours must be a multiple of 0.5

Lines 106-106

The hours on a time amendment must be a multiple of 0.5. Turbo checks this with IsMultipleOfPointFive.

## Time amendment date must be valid

Lines 107-107

The date of a time amendment must pass the IsValidDate check.

## Meals hours capped at 3

Lines 109-109

Meals hours on a time amendment must not exceed 3. Exactly 3 is allowed. Otherwise Turbo says "Meals Hours must not exceed 3".

## Description needs at least 5 characters

Lines 110-110

The description of a time amendment must have at least 5 characters. Otherwise Turbo says "must be at least 5 characters".

## Work record format depends on whether a job is set (rule check)

Lines 111-111

The work record pattern is checked only when the job is present. Without a job the work record must be empty.

## Payout request amount must be positive with cents only

Lines 112-112

The payout request amount (payout_request_amount) must be greater than 0. Zero fails with "Amount must be greater than 0". The amount must also be a positive multiple of 0.01 (IsPositiveMultipleOfPointZeroOne).

## Division must be allocated to the job on create or change

Lines 115-131

For a time amendment with a job, Turbo checks that the division is allocated to the job (validateDivisionAllocatedToJob). It does this only when shouldValidateJobDivisionAllocationOnRecord says so. According to the code comment, that means on create and when the job or division changes. Unrelated edits to older records with historical mismatches are allowed.

## Job must exist and have a status

Lines 119-125

If a time amendment has a job, the job must exist. Otherwise Turbo says "invalid job reference". The job must also have a status. Otherwise Turbo says "job status must be set".

## Authentication required for time amendments

Lines 144-151

Saving a time amendment without a signed-in user fails with "authentication is required".

## Time amendments blocked when time editing is disabled

Lines 153-159

When time editing is disabled in Turbo (IsTimeEditingEnabled is false), creating or editing a time amendment is refused with ErrTimeEditingDisabled. The message text is defined in another file.

## Time amendment owner cannot be changed

Lines 161-188

On an update, changing the uid (the staff member the amendment belongs to) is refused with "uid cannot be changed". An update with no original record to compare is refused with "cannot validate time amendment update".

## Committed time amendments cannot be edited

Lines 172-179

A time amendment that has a committed date cannot be edited. Turbo says "cannot edit committed time amendment".

## Creator is set to the latest editor

Lines 191-193

On every create or update, Turbo sets the creator of the time amendment to the signed-in user. The creator is the latest person to edit it before commit.

## Week ending is derived from the date

Lines 202-208

Turbo sets the week ending (week_ending) of a time amendment from its date alone. If it cannot be generated, the save fails with "Error generating week_ending".

## Amendment must match a committed time sheet

Lines 214-246

Unless skip_tsid_check is set, Turbo looks for the time sheet of the same user and week ending. If none exists, Turbo says "No time_sheets record found, create one instead". If more than one exists, Turbo says "More than one time_sheets record found".

## Time sheet id is linked automatically

Lines 227-261

When the time sheet check passes, Turbo sets the tsid of the time amendment to the matching time sheet's id. With skip_tsid_check set, Turbo skips the whole time sheet check and does not set tsid.

## Time sheet must be committed to amend it

Lines 248-257

If the matching time sheet is not committed, Turbo refuses the amendment. It says "Found time_sheets record has not been committed. Create a time_entries record instead."
