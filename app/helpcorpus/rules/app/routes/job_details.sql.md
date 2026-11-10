<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_details.sql","sha256":"ae59f34c648412fce89adddae9101d1064e34a54ed7f0f92ba86f3e6ada96644","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_details.sql

## Job details returns core job fields and the imported flag

Lines 4-14

Asked as: job number missing; job status and location; imported flag; client PO not showing; authorizing document on job

Job details returns the job's number, description, status, location, authorizing document, client PO and client reference number. It returns the internal import flag (_imported) as imported.

## Job details shows the parent job number for a subjob

Lines 9-10

Asked as: subjob parent number; which job is this under; parent job missing; job's parent job

Job details returns the job's parent (parent_id) and the parent job's number (parent_number). Both are empty when the job has no parent.

## Job details shows client name and alias

Lines 15-17

Asked as: client alias not showing; client name on job; which client is this job for

Job details returns the job's client (client_id) with the client's name and alias (client_name, client_alias).

## Job details shows the job contact's name

Lines 18-20

Asked as: job contact person; client contact missing; contact name on job

Job details returns the job's client contact (contact_id) with the contact's given name and surname.

## Job details shows manager and alternate manager names

Lines 21-26

Asked as: who manages this job; alternate manager blank; job manager name

Job details returns the job's manager and alternate manager (manager_id, alternate_manager_id). Each comes with the given name and surname from the staff profile.

## Job details shows the job owner from the clients table

Lines 27-29

Asked as: owner of the job; job owner name wrong; owning client on job

Job details returns the job owner (job_owner_id) with the name and alias of the client record that it points to (job_owner_name, job_owner_alias).

## Job details shows the proposal behind a job

Lines 30-31

Asked as: which proposal is this from; proposal number on job; job's source proposal

Job details returns the job's proposal (proposal_id) and the proposal's number (proposal_number).

## Job details returns agreement and project dates

Lines 32-34

Asked as: FN agreement date; project award date; project completion date; agreement date missing

Job details returns the job's FN agreement (fn_agreement), project award date and project completion date.

## Job details shows invoice contact name and instructions from the invoicing profile

Lines 35-39

Asked as: who gets the invoice; invoice instructions missing; invoicing profile contact; billing contact on job

The invoice contact name, email, profile name and invoice instructions come from the job's invoicing profile (invoicing_information) and that profile's contact. The contact name joins given name and surname with a space and trims the ends.

## Job details shows request creator and job creator names

Lines 40-43

Asked as: who created the request; request submitted date; who made this job; creator name blank

The request creator name and the created-by name each join the person's given name and surname with a space and trim the ends. The query also returns the request dates (request_created, request_submitted) and creation_completed.

## Job details returns proposal dates and values

Lines 44-46

Asked as: proposal due date; proposal opening date; proposal value amount; proposal submission deadline

Job details returns the proposal opening date, the proposal submission due date and the proposal value. Lines 44-46 hold these fields.

## Job details returns project value, time and materials, and outstanding balance

Lines 47-50

Asked as: project total value; time and materials setting; outstanding balance amount; amount still owed on job

Job details returns the project value, the time-and-materials setting, the outstanding balance and the date of the outstanding balance.

## Job details shows branch code and name

Lines 51-53

Asked as: which branch the job is in; branch code on job; office branch for job

Job details returns the job's branch (branch_id) with its code and name.

## Job details shows rate sheet name and revision

Lines 54-56

Asked as: rate sheet version; which rate card; rate sheet name blank

Job details returns the job's rate sheet (rate_sheet_id) with the rate sheet's name and revision.

## Job details lists projects created from a proposal

Lines 57-61

Asked as: jobs created from proposal; projects from this proposal; jobs spawned by proposal

The job details query returns the id and number of every job whose proposal is this job (projects_json). The list is an empty array when no job points to this job as its proposal.

## Job details lists subjobs of the job

Lines 62-66

Asked as: child jobs list; jobs under this job; list of subjobs

The job details query returns the id and number of every job whose parent is this job (children_json). The list is an empty array when the job has no subjobs.

## Job details lists the job's categories

Lines 67-71

Asked as: job categories list; category missing on job; which categories job has

The job details query returns the id and name of every category that belongs to this job (categories_json). The list is an empty array when the job has no categories.

## Job details lists time allocations by division

Lines 72-86

Asked as: hours by division; time allocated per division; division hours blank; time split across divisions

The job details query returns each time allocation (job_time_allocations) of the job with its division id, code and name, and its hours. Hours show as 0 when the allocation has no hours value. The list is an empty array when the job has no allocations. An allocation whose division record is missing is left out.

## Job details returns the job even when related records are missing

Lines 88-100

Asked as: job blank when contact missing; job not loading missing data; empty fields on job; job shows nothing for client

The job details query uses left joins for the proposal, parent, invoicing profile, client, contact, managers, job owner, branch, rate sheet and creator profiles. The job still returns when any of them is empty, and their fields are then null.

## Job details returns exactly one job selected by id

Lines 101-101

Asked as: open job by id; job lookup by id; why only one job returns; fetch single job

The job details query returns the single job whose id equals the supplied id. The id must be provided. It does not filter by status or by who is asking.
