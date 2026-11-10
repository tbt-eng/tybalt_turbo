<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_details_api.go","sha256":"3668c61b0e4bb5208e104ee75ed5fae3649fbe9ad67db1efe57d724f0f4dea0d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_details_api.go

## Job details list each allocation with division and hours

Lines 18-29

Asked as: allocation division hours; hours per division; job allocations breakdown; division code on allocation

Each allocation in the job details (allocations) shows its division (id, code and name) and its hours. Turbo reads the allocations from the allocations_json column of the job details query.

## Job details name the contact, manager and alternate manager

Lines 31-35

Asked as: job contact and manager; alternate manager on job; who manages the job; job manager name

Turbo shows the job's contact, manager and alternate manager (alternate_manager) each as an id, a given name and a surname.

## Job details show rate sheet, branch, manager and owner

Lines 37-72

Asked as: job branch and rate sheet; job owner details; fn agreement flag; job client and owner

Job details include the job's branch, rate sheet (name and revision), manager, alternate manager, contact, client (with alias) and job owner. They also include the First Nations agreement flag (fn_agreement) and the imported flag.

## Job details report whether the job is imported or has a first nations agreement

Lines 37-98

Asked as: is job imported; job client PO number; job authorizing document; job location and reference

Job details include the flags imported and fn_agreement. They also include the job's location, client PO, client reference number and authorizing document.

## Job details show a proposal's number and link when the job has one

Lines 43-44

Asked as: job linked proposal number; parent job of subjob; which proposal is job from; subjob parent number

Job details carry the linked proposal's ID and number (proposal_id, proposal_number). They also carry the parent job's ID and number (parent_id, parent_number) for a subjob. These fields are empty when the job has no such link.

## Job details show proposal dates and values plus project value

Lines 73-93

Asked as: proposal value and dates; job project value; job award and completion date; job outstanding balance

Job details include the proposal opening date, the proposal submission due date, the proposal value and the project value. They also include the project award and completion dates, the time and materials flag and the outstanding balance with its date.

## Job details show invoicing profile and invoice contact

Lines 75-79

Asked as: invoice contact on job; job invoicing profile; invoice instructions for job; who gets job invoices

Job details include the invoicing profile name, the invoice contact's name and email, the invoice instructions and the invoicing information. Each is empty when the job has none.

## Job details show who requested, created and completed the job

Lines 80-86

Asked as: who created the job; job request submitted date; who requested the job; job creation completed time

Job details include the project request's creator, created time and submitted time (request_creator, request_created, request_submitted). They also include who created the job (created_by) and when creation completed (creation_completed). Both creator names are included for display.

## Job details show the client with its name and alias

Lines 100-104

Asked as: job client name and alias; client shown on job; job owner name; client alias on job

Job details show the job's client as an id, a name and an alias. The job owner (job_owner) appears in the same shape: id, name and alias.

## Job details show the rate sheet name and revision

Lines 106-110

Asked as: rate sheet revision on job; job rate sheet name; which rate sheet job uses

The rate sheet in job details (rate_sheet) shows its id, name and revision number.

## Linked projects and child jobs show only id and number

Lines 112-115

Asked as: linked projects on job; child jobs list; job children numbers; related jobs in job details

The projects and children lists in job details give each linked job (JobRef) as an id and a number only.

## Missing numeric or flag values in job details default to zero or false

Lines 157-160

Asked as: job value shows zero; no outstanding balance shown; time and materials false; job missing values zero

When the job has no proposal value, project value or outstanding balance, details show 0. A missing time and materials flag shows false. A missing rate sheet revision shows 0.

## Only admins see unlinked subjobs

Lines 166-168

Asked as: unlinked subjobs not showing; adopt subjobs; who sees orphan subjobs; missing subjobs on project

Admins also get a list of unlinked subjobs (unlinked_subjobs). These are jobs numbered as this project's subjobs that are not linked to it. Admins can adopt them. Other users do not receive this list.

## Categories in job details show id and name

Lines 171-174

Asked as: job categories list; categories on job; job category names

Each category in job details (categories) shows an id and a name.

## Only active accounts can view job details

Lines 178-186

Asked as: cannot view job details; inactive account job access; account not active error; job details access denied

Turbo refuses to show job details when the user's account is not active, even if the user still holds a valid login token. The message is "an active account is required to view job details". The check exists because the response includes invoicing profile information.

## Job details for an unknown job are refused

Lines 189-195

Asked as: job id not found; job does not exist error; wrong job id; unknown job request

When no job matches the requested id, Turbo refuses the request with the message "job not found".

## Job details include allocations, linked jobs and categories

Lines 198-207

Asked as: job details empty lists; allocations missing from job; job linked jobs not loading; job details no error empty

The job details list the job's division allocations (hours per division), its projects, its child jobs and its categories. Each list comes from the job details query. If a list cannot be read, Turbo shows it as empty and gives no error.

## Missing job fields appear as empty values

Lines 209-266

Asked as: blank job fields; job status empty; why job field blank; empty contact on job

When a job detail has no value, Turbo shows an empty string. Missing numbers show as 0 and missing flags show as false. Examples are the status, the contact, the alternate manager, the proposal value, the project value, the outstanding balance and time and materials. The rate sheet revision shows 0 when there is none.

## Unlinked subjob lookup runs for admins only

Lines 268-280

Asked as: failed to load job; failed to find unlinked subjobs; admin job details error; admin job request fails

For a user with the admin claim, Turbo loads the job and looks for unlinked subjobs. If either step fails, Turbo refuses the whole request, with "failed to load job" or "failed to find unlinked subjobs". Non-admins skip this step.
