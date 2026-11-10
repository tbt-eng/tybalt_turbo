<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspaceLoad.ts","sha256":"8e8548be953eaf44ebfc074de473926fd1762644aa3264d8c198952759d4520d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspaceLoad.ts

## Only the details tab loads time divisions

Lines 22-30

Turbo fetches the client's time divisions only when the details tab is open. On the jobs and notes tabs the list of time divisions is empty. The request starts together with the client request, so the tab waits for the slower of the two.

## Client workspace shows 'Client not found.' for a missing client

Lines 33-45

When Turbo loads a client's workspace and the client lookup returns not found, it shows the message "Client not found.". Any other failure while loading the client is passed on unchanged.

## Each client workspace tab loads only its own data

Lines 46-60

The jobs tab loads job lists and no notes. The notes tab loads notes and no job lists. Other tabs load neither, so their jobs and notes lists are empty.

## Job lists show 20 per page, newest number first

Lines 53-96

The open view (projects, proposals or owner) shows 20 jobs per page, sorted by number in descending order. Each view has its own page parameter: projectsPage, proposalsPage or ownerPage.

## Client notes tab lists notes and the client's jobs

Lines 62-73

The notes tab loads the client's notes. It also loads every job where the client is the client or the job owner (job_owner), newest first, with only id, number and description.

## Job lists split into projects, proposals and owner jobs

Lines 81-85

Proposals are the client's jobs whose number starts with P. Projects are the client's jobs whose number does not start with P. Owner jobs are jobs where the client is the job owner (job_owner), whatever the number.

## Job tab counts cover all three lists

Lines 86-118

Turbo shows the total number of projects, proposals and owner jobs for the client. It gets each count by loading one item for each view that is not open.

## Out-of-range job page jumps to the last page

Lines 98-108

If the requested page is beyond the last page, for example after a job is removed or from an old bookmark, Turbo shows the last page instead. The number of pages is the total jobs divided by 20, rounded up, and never less than 1.
