<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_adopt_subjobs.go","sha256":"a8a6d475446daa4cd587b6e753a8d39a9fbe6a1bfaa3dfbc07313b7050bf074e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_adopt_subjobs.go

## Only top-level projects with a plain number can have subjobs found

Lines 19-54

Asked as: subjobs not showing for project; why no unlinked subjobs; project number format for subjobs; subjob list empty; only top level projects

Turbo looks for unlinked subjobs only when the project number looks like 24-410 or 25-0123 (two digits, a dash, then 3 or 4 digits) and the project has no parent. For any other project the list of unlinked subjobs is empty.

## Which jobs count as unlinked subjobs of a project

Lines 55-75

Asked as: which jobs are subjobs; jobs missing parent link; find loose subjobs; nested job numbers; job number like project dash

Turbo lists jobs with no parent whose number is the project number plus a dash and one level of digits, such as 24-410-1 for 24-410. It leaves out nested numbers such as 24-410-1-2 and numbers of other projects such as 24-4100-1. The list is ordered by job number.

## Subjob blocked when the project has no client

Lines 77-79

Asked as: subjob blocked no client; project missing client adopt; set client before adopting; client needed to adopt; adoption blocked client empty

An unlinked subjob is listed with the blocker "This project has no client. Set its client first." when the project has no client. Turbo does not adopt it.

## Subjob blocked when its client differs from the project's

Lines 80-85

Asked as: subjob client mismatch; different client subjob; client doesn't match project; cannot adopt wrong client; change subjob client

An unlinked subjob whose client differs from the project's client is listed with the blocker "Its client is X. Change its client to Y first." X is the subjob's client, or "no client" when it has none. Y is the project's client. Turbo does not adopt it.

## Subjob blocked when it has subjobs of its own

Lines 86-88

Asked as: subjob has its own subjobs; subjob with children blocked; nested subjobs cannot adopt; why subjob won't adopt; subjob already has children

An unlinked subjob that already has subjobs is listed with the blocker "It has subjobs of its own." Turbo does not adopt it. Blockers are checked in this order: no project client, client mismatch, then own subjobs.

## Adopting links subjobs and marks them for legacy write-back

Lines 94-187

Asked as: adopt subjobs link parent; subjob legacy write back; imported flag subjobs; adoption skips hooks validation; subjobs linked to project

Adoption sets the project as the parent of each adoptable subjob and sets its imported flag (_imported) to false, so it is written back to the legacy system. It updates the records directly, so job hooks and record validation do not run. Blocked subjobs are left unchanged. The whole change runs in one transaction.

## Only admins can adopt subjobs

Lines 104-106

Asked as: non admin adopt subjobs; who can adopt subjobs; only admins adopt; permission to adopt subjobs; adopt button unauthorized

Adopting subjobs requires an admin (the admin claim). Anyone else gets "Only admins can adopt subjobs." A request with no signed-in user is refused as unauthorized.

## Adopting subjobs needs job editing to be enabled

Lines 142-148

Asked as: job editing disabled adopt; cannot adopt jobs editing off; jobs editing disabled error; editing turned off subjobs; adopt fails editing disabled

Turbo refuses to adopt subjobs while job editing is disabled. It returns the jobs-editing-disabled error (ErrJobsEditingDisabled). This check runs after the admin check.

## Project must exist and be top-level to adopt subjobs

Lines 150-156

Asked as: project not found adopt; project number invalid adopt; subjob cannot adopt subjobs; top level project required; adopt from subjob project

If the project is not found, Turbo says "Project not found." If the project number does not look like 24-410 or 25-0123, or the project is itself a subjob, Turbo says "Only a top-level project can adopt subjobs."

## Adoption refused when nothing can be adopted

Lines 162-172

Asked as: nothing to adopt; no subjobs can be adopted; all subjobs blocked; adoption does nothing; no adoptable subjobs

If every unlinked subjob has a blocker, or there are none, Turbo says "No unlinked subjobs of this project can be adopted."

## Confirmed list must match the adoptable subjobs exactly

Lines 173-175

Asked as: confirmed subjobs changed; subjob list changed refresh; partial adoption refused; confirm again subjobs; selection does not match subjobs

The subjobs that the user confirms must be the same set as the subjobs that can be adopted now. Order and repeats do not matter. If they differ, Turbo refuses with "The subjobs confirmed do not match the subjobs this project can adopt, which may have changed. Review the current list and confirm again." Turbo adopts all or none of the adoptable subjobs, and a partial selection is refused.

## Adoption adds a note to the project

Lines 189-193

Asked as: audit note adopted subjobs; project note adopted numbers; record of adopted subjobs; which subjobs were adopted; adoption history note

After adoption Turbo adds an audit note to the project: "Adopted unlinked subjobs: " followed by the adopted numbers, separated by commas. The response lists the adopted numbers.
