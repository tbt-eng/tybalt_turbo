<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspaceLoad.ts","sha256":"e643e9b496dd707bb86ca3d2cded14ef48fb28d3d89abb60438c3b6597086a83","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspaceLoad.ts

## Only the details tab loads divisions and rate sheets

Lines 23-38

Asked as: divisions empty on jobs tab; rate sheets missing; time divisions not showing; why are divisions blank; rate sheets not loading

Turbo loads a client's time divisions and rate sheets only on the details tab. On the jobs and notes tabs these lists are empty.

## Client page shows 'Client not found.' for a missing client

Lines 41-54

Asked as: client not found error; client page error; can't open client; missing client message; client does not exist

When a client's page loads and the client record does not exist, Turbo shows the error 'Client not found.' Any other failure loading the client is passed on unchanged.

## Jobs tab loads job lists; notes tab loads notes

Lines 57-70

Asked as: notes tab empty; jobs tab no notes; notes not loading; client notes missing; why no jobs on notes tab

Turbo loads the client's jobs only on the jobs tab and the client's notes only on the notes tab. The other tabs get empty job and note lists, so each tab loads only what it shows.

## Jobs lists show 20 per page, highest number first

Lines 63-106

Asked as: how many jobs per page; job list order; jobs sorted wrong; page parameter jobs; jobs list paging

The open job list (view) shows 20 jobs per page, sorted by number from highest to lowest. Each of the other two lists is fetched with one job only, to give its count. Each list has its own page parameter: projectsPage, proposalsPage and ownerPage.

## Notes tab also lists the client's jobs, newest first

Lines 72-83

Asked as: jobs showing on notes tab; job owner jobs on notes; why are jobs listed under notes; newest jobs first notes

On the notes tab, Turbo also loads the jobs where the client is the client or the job owner (job_owner). They are sorted by newest created first, and only id, number and description are fetched.

## Client jobs split into projects, proposals and owner jobs

Lines 91-95

Asked as: proposals vs projects; where are proposals listed; owner jobs list; jobs starting with P; client jobs categories

The jobs tab has three lists. Proposals are the client's jobs whose number starts with P. Projects are the client's jobs whose number does not start with P. Owner jobs are jobs where the client is the job owner (job_owner).

## Page beyond the last page falls back to the last page

Lines 108-122

Asked as: blank job page; page past the end; old bookmark shows wrong page; jobs page empty after deleting; last page of jobs

If the requested page of a job list is past the last page, for example from an old bookmark or after a job is removed, Turbo shows the last page instead. The number of pages is at least 1, even with no jobs. The code rounds the number of pages up from the total divided by 20.

## Job counts cover all three lists

Lines 119-128

Asked as: job count wrong; total jobs number; how many projects shown; count of proposals; job totals on tab

The jobs tab shows the total number of projects, proposals and owner jobs, whichever list is open.
