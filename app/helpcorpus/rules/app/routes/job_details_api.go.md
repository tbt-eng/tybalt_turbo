<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_details_api.go","sha256":"4f2971f880cf02042dea9db63a9e55ac42a325fa67b56d7912f95cc7bec58780","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_details_api.go

## Job details need an active account

Lines 180-188

Asked as: inactive account can't view job; job details blocked account; account not active job access; locked out of job details

Turbo refuses to show job details to a user whose account is not active, even if the user still holds a login token. The message is "an active account is required to view job details". The check matches the read rule for invoicing profiles, because the response includes them.

## Job details for an unknown job

Lines 195-197

Asked as: job id not found; job doesn't exist; wrong job number; no such job

When no job matches the requested id, Turbo answers "job not found".

## Job details content returned

Lines 200-268

Asked as: empty fields in job details; missing values job info; blank job details fields; what fields job details show

The job details response combines the job's client, contact, manager, alternate manager, job owner, branch, rate sheet (with its revision), invoicing profile, request and creation details, proposal and project values, outstanding balance, allocations by division, related projects, child jobs and categories. Any value that is missing comes back as an empty string, 0 or false. A list that cannot be read comes back empty.

## Only admins see unlinked subjobs

Lines 270-282

Asked as: extra subjobs list admins; jobs not linked to project; orphan subjobs visible; who sees unlinked subjobs

Admins get an extra list of unlinked subjobs (unlinked_subjobs) in the job details. These are jobs numbered as this project's subjobs but not linked to it, and admins can adopt them. Other users do not receive this list.
