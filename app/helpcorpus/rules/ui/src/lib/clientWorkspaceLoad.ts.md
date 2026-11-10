<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspaceLoad.ts","sha256":"8e8548be953eaf44ebfc074de473926fd1762644aa3264d8c198952759d4520d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspaceLoad.ts

## Only the details tab loads time divisions

Lines 22-30

Asked as: time divisions empty; divisions not loading; jobs tab slow; details tab divisions

Turbo fetches the client's time divisions only when the details tab is open. On the jobs and notes tabs the list of time divisions is empty. The request starts together with the client request, so the tab waits for the slower of the two.

## Client workspace shows 'Client not found.' for a missing client

Lines 33-45

Asked as: client not found; missing client error; client page error; can't open client

When Turbo loads a client's workspace and the client lookup returns not found, it shows the message "Client not found.". Any other failure while loading the client is passed on unchanged.

## Each client workspace tab loads only its own data

Lines 46-60

Asked as: notes tab empty; jobs tab no notes; tab shows no data; client tab empty lists

The jobs tab loads job lists and no notes. The notes tab loads notes and no job lists. Other tabs load neither, so their jobs and notes lists are empty.

## Job lists show 20 per page, newest number first

Lines 53-96

Asked as: how many jobs per page; job list order; jobs sorted newest; page of jobs

The open view (projects, proposals or owner) shows 20 jobs per page, sorted by number in descending order. Each view has its own page parameter: projectsPage, proposalsPage or ownerPage.

## Client notes tab lists notes and the client's jobs

Lines 62-73

Asked as: client notes and jobs; notes tab job list; jobs listed under notes; client's jobs in notes

The notes tab loads the client's notes. It also loads every job where the client is the client or the job owner (job_owner), newest first, with only id, number and description.

## Job lists split into projects, proposals and owner jobs

Lines 81-85

Asked as: projects vs proposals; proposal jobs P number; owner jobs missing; jobs where client is owner

Proposals are the client's jobs whose number starts with P. Projects are the client's jobs whose number does not start with P. Owner jobs are jobs where the client is the job owner (job_owner), whatever the number.

## Job tab counts cover all three lists

Lines 86-118

Asked as: job counts wrong; number of projects shown; job tab totals; count of proposals

Turbo shows the total number of projects, proposals and owner jobs for the client. It gets each count by loading one item for each view that is not open.

## Out-of-range job page jumps to the last page

Lines 98-108

Asked as: jobs page blank; page past end; deleted job page error; old bookmark jobs page

If the requested page is beyond the last page, for example after a job is removed or from an old bookmark, Turbo shows the last page instead. The number of pages is the total jobs divided by 20, rounded up, and never less than 1.
