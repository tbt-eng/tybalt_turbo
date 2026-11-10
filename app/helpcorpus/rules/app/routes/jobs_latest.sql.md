<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_latest.sql","sha256":"f654707c055252c579ec287e2c066ed2820c4aebf550efa7c6a2a1dfe4d4024c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_latest.sql

## Latest jobs list takes the newest jobs from each group

Lines 4-43

The latest jobs list takes at most :limit proposals and at most :limit projects, newest created first. The limit applies to each group separately, so the list can hold up to twice the limit.

## Missing outstanding balance shows as 0 and missing date as blank

Lines 13-14

In the latest jobs list, a job with no outstanding balance shows 0. A job with no outstanding balance date shows an empty text.

## Proposals are jobs whose number starts with P

Lines 17-18

The latest jobs list treats a job as a proposal when its number starts with 'P' (LIKE 'P%'). Every other job counts as a project.

## Jobs with no matching client still appear in latest jobs list

Lines 51-62

The latest jobs list joins the client, branch and manager as optional. A job stays in the list even if the lookup finds no match. The client name comes from the clients table and is empty when no client matches.

## Missing branch or manager shows blank in latest jobs list

Lines 52-53

The latest jobs list shows an empty branch when a job has no branch. It shows an empty manager when the job has no manager. The manager appears as given name, a space, then surname.

## Latest jobs list shows proposals first, then projects

Lines 56-63

The latest jobs list returns two groups: proposals (group_name 'Proposals'), then projects (group_name 'Projects'). Each group is sorted by created date, newest first.
