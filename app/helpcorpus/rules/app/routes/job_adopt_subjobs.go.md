<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_adopt_subjobs.go","sha256":"33174865d1ef4827e9f44cdb2f1709d28200738a5206a1c63d4290c281db78c3","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_adopt_subjobs.go

## Only top-level project numbers can have subjobs adopted

Lines 19-54

Asked as: adopt subjobs on a sub project; why no subjobs listed; project number format for subjobs; only top level projects subjobs

Turbo looks for unlinked subjobs only when the project number has the form 2 digits, a hyphen, then 3 or 4 digits (such as 24-410 or 25-0123). The project must also have no parent. For any other project the list of unlinked subjobs is empty.

## Which jobs count as unlinked subjobs

Lines 55-76

Asked as: unlinked subjobs list; jobs not linked to project; orphan jobs under project number; which jobs can be attached

An unlinked subjob is a job with no parent whose number is the project's number plus a hyphen and a first-level suffix, such as 24-410-1 for 24-410. Turbo excludes numbers like 24-4100-1 for 24-410. It also excludes nested numbers like 24-410-1-2 (checked with jobChildNumberRegex, defined in another file). Turbo lists them in number order.

## Blockers that stop an unlinked subjob being adopted

Lines 77-90

Asked as: subjob cannot be adopted; why subjob is blocked; subjob has no client; client mismatch subjob adopt

Turbo lists a subjob with a blocker when it cannot be adopted. If the project has no client, the blocker is "This project has no client. Set its client first." If the subjob's client differs from the project's, the blocker is "Its client is X. Change its client to Y first." (X is "no client" when empty). If the subjob has subjobs of its own, the blocker is "It has subjobs of its own." Turbo checks these in that order and shows only the first that applies.

## Only active admins can adopt subjobs

Lines 104-145

Asked as: cannot adopt subjobs; not allowed to adopt; adopt button not working; admin rights needed adopt

Adopting subjobs requires a signed-in user who is active and holds the admin claim. Anyone else gets "Only active admins can adopt subjobs." A request with no sign-in is refused as unauthorized.

## Adopting subjobs is refused when job editing is off

Lines 146-152

Asked as: adopt subjobs disabled; job editing turned off; cannot edit jobs now; adoption blocked read only

If job editing is disabled, Turbo refuses to adopt subjobs and returns utilities.ErrJobsEditingDisabled (defined in another file). This check runs after the admin check.

## Adopt subjobs needs an existing top-level project

Lines 154-160

Asked as: project not found adopt; adopt into wrong project; subjobs adopt parent project error; project cannot take subjobs

If the project does not exist, Turbo says "Project not found." If its number is not of the form 24-410 or 25-0123, or it has a parent, Turbo says "Only a top-level project can adopt subjobs."

## Nothing to adopt is an error

Lines 162-176

Asked as: nothing to adopt; no subjobs to attach; adopt returns error empty; all subjobs blocked

If none of the project's unlinked subjobs is free of blockers, Turbo refuses with "No unlinked subjobs of this project can be adopted."

## Admin must confirm the exact list of adoptable subjobs

Lines 177-179

Asked as: subjob list changed; confirm adoption again; adoption mismatch error; subjobs changed before confirming

The confirmed subjobs in the request must be exactly the set of adoptable subjobs (order and repeats do not matter). If the set differs, Turbo refuses with "The subjobs confirmed do not match the subjobs this project can adopt, which may have changed. Review the current list and confirm again." Blocked subjobs are never adopted, and adoption is all-or-nothing for the adoptable ones.

## What adoption changes on the subjobs and project

Lines 181-198

Asked as: what adopting changes; subjob parent after adoption; adopted subjobs note project; subjobs marked for write back

Adoption sets each subjob's parent to the project and sets _imported to false, so the subjobs are marked as locally changed for write-back to the legacy system. It also updates their updated time. It writes directly, so the job hooks and record validation do not run. Turbo adds a note to the project, "Adopted unlinked subjobs: " followed by the subjob numbers. It returns the list of adopted numbers.
