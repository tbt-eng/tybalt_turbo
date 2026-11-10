<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_adopt_subjobs.go","sha256":"a8a6d475446daa4cd587b6e753a8d39a9fbe6a1bfaa3dfbc07313b7050bf074e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_adopt_subjobs.go

## Which projects can have unlinked subjobs listed

Lines 19-54

Turbo lists unlinked subjobs only for a top-level project whose number looks like 24-410 or 25-0123 (two digits, a dash, then 3 or 4 digits). A project that has a parent (parent) or any other number format gets an empty list.

## Which jobs count as unlinked subjobs

Lines 55-74

An unlinked subjob is a job with no parent whose number is the project's number plus a dash and a single level suffix, such as 24-410-1 for 24-410. Nested numbers such as 24-410-1-2 are left out. A number like 24-4100-1 does not match project 24-410. The list is sorted by number.

## Blocker: project has no client

Lines 77-79

When the project has no client, every unlinked subjob is listed with the blocker "This project has no client. Set its client first." and cannot be adopted.

## Blocker: subjob client differs from project client

Lines 80-85

A subjob whose client differs from the project's client is listed with the blocker "Its client is X. Change its client to Y first." X shows "no client" when the subjob has none. The subjob is not adopted.

## Blocker: subjob has its own subjobs

Lines 86-88

A job that already has subjobs of its own is listed with the blocker "It has subjobs of its own." and cannot be adopted. The client check comes first, so a client mismatch shows its blocker instead.

## Adoption links subjobs without running job checks

Lines 96-187

Adoption sets each adoptable subjob's parent to the project and sets _imported to false so it is written back to the legacy system. It updates the records directly, so job hooks and record validation do not run. Blocked subjobs are left unlinked. All changes happen in one transaction.

## Only admins can adopt subjobs

Lines 104-141

Adopting subjobs requires an admin (the admin claim). Anyone else gets "Only admins can adopt subjobs." Users who are not signed in are refused as unauthorized.

## Adopting requires job editing to be enabled

Lines 142-148

Turbo refuses to adopt subjobs when job editing is switched off. It returns the error ErrJobsEditingDisabled (defined elsewhere). This check runs after the admin check.

## Project must exist and be top-level

Lines 150-156

If the project is not found, Turbo says "Project not found." If its number does not look like 24-410 or 25-0123, or it has a parent, Turbo says "Only a top-level project can adopt subjobs."

## Nothing adoptable is refused

Lines 158-172

Only subjobs with no blocker are adopted. If none qualify, Turbo says "No unlinked subjobs of this project can be adopted."

## Confirmed list must match adoptable subjobs

Lines 173-175

The subjobs the admin confirms must be exactly the set that can be adopted now. Order and repeats are ignored. If the set differs, Turbo says "The subjobs confirmed do not match the subjobs this project can adopt, which may have changed. Review the current list and confirm again." Nothing is changed.

## Adoption writes a note on the project

Lines 189-193

After adoption, Turbo adds a note to the project: "Adopted unlinked subjobs: " followed by the adopted numbers separated by commas. The request returns the list of adopted numbers.
