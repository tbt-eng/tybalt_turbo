<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_details_api.go","sha256":"2e7fe611d114925e5bb8a4d3d5b035760e73a35c10204715dbf86cbaaadd890f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_details_api.go

## Job details return an error when the job is not found

Lines 186-188

Asked as: job not found; job details wrong id; no such job; job id doesn't exist

When a job details request names a job ID that matches no row, Turbo refuses it with the message "job not found".

## Job details list allocations, projects, children and categories

Lines 191-200

Asked as: job details allocations empty; missing allocations on job; job children not showing; job categories blank list

Job details include the job's allocations (division and hours), its projects, its children and its categories. Turbo takes them from the job details query. If that data cannot be read, the list comes back empty and Turbo gives no error.

## Missing job details fields come back as empty values

Lines 202-259

Asked as: blank field job details; job manager shows empty; rate sheet revision 0; missing contact on job

Job details show an empty string, 0 or false for a field that has no value, such as a missing manager, contact, branch, rate sheet, dates, proposal value or project value. Turbo does not leave these fields out. A missing rate sheet revision shows as 0.

## Only admins see unlinked subjobs in job details

Lines 261-273

Asked as: unlinked subjobs not visible; can't see subjobs list; why no unlinked subjobs; subjobs missing for non admin

Turbo adds the list of unlinked subjobs (unlinked_subjobs) to job details only when the person holds the admin claim. These are jobs numbered as this project's subjobs that are not linked to it. Other people do not receive the list. Turbo also loads the job record and runs findUnlinkedSubjobs for admins, and shows an error if either step fails.
