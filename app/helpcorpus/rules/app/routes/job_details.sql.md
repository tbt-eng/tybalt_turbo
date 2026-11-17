<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_details.sql","sha256":"2974352e7f165a15336bff8c4a748b8e24d3d61b8d86bd91391c228cdf5e39c7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_details.sql

## Job details return one job selected by its ID

Lines 1-3

Asked as: get one job by id; job details for a single job; fetch job record by id; why does job lookup return only one job

The job details query returns a single job, chosen by the id given to it. Related records come back in the same result, so the app does not fetch them one by one.

## Job details show branch and rate sheet

Lines 8-8

Asked as: branch name on job; rate sheet revision; which rate sheet is job using; imported flag on job

The job details give the branch code and name, and the rate sheet's code, name and revision. They also give the imported flag (_imported) as imported.

## Job details show the parent job's number

Lines 9-10

Asked as: parent job number; subjob parent number; which job is my subjob under; missing parent number on job

For a subjob, the job details give the parent job's number (parent_number) beside the parent's ID. A job with no parent has no parent number.

## Job details include proposal and project values and balance

Lines 12-50

Asked as: proposal value on job; outstanding balance on job; client po number on job; award and completion dates

The job details give the proposal dates, the proposal value, the project value, the time and materials setting, the outstanding balance and its date. They also give the award and completion dates, the location, the authorizing document, the client PO, the client reference number and the FN agreement.

## Job details show client, contact and job owner names

Lines 15-29

Asked as: client name missing on job; job owner name blank; contact name not showing on job; client alias on job

The job details give the client's name and alias, the client contact's given name and surname, and the job owner's name and alias. The job owner is looked up in the clients records. A missing link leaves these names empty.

## Job manager and alternate manager names come from profiles

Lines 21-26

Asked as: job manager name blank; alternate manager not showing; who manages this job; manager name from profile

The job details give the manager's and the alternate manager's given name and surname. Turbo matches them to profiles by user ID (uid). A job with no alternate manager has no alternate manager name.

## Job details link a proposal and its projects

Lines 30-62

Asked as: which proposal is job linked to; projects under this proposal; proposal number for job; related projects list empty

The job details give the job's proposal ID and number (proposal_number). They also list the projects whose proposal is this job (projects_json) as ID and number pairs. The list is an empty array when there are none.

## Job details show invoicing profile, contact and instructions

Lines 35-39

Asked as: invoice contact on job; invoicing profile name; invoicing instructions wrong contact; billing contact for job

The job details show the name of the job's invoicing profile, its invoice contact's full name and email, and its invoicing instructions. The contact is the one set on the invoicing profile, not the job's own contact. The full name joins given name and surname with a space and trims it.

## Job details show who requested and who created the job

Lines 40-43

Asked as: who created this job; job creator name; who requested the project; request submitted date

The job details give the project request creator, the request created and submitted dates, the creator of the job and the creation completed value (creation_completed). They also give the full names of the request creator and of the creator, taken from profiles.

## Job details list subjobs of the job

Lines 63-67

Asked as: child jobs of job; list of subjobs; job has no subjobs showing; subjob list on job

The job details list the child jobs (subjobs) of this job as ID and number pairs (children_json). The list is an empty array when the job has no subjobs.

## Job details list the job's categories

Lines 68-72

Asked as: job categories not showing; categories on job; job category list empty; tags for job

The job details list the job's categories as ID and name pairs (categories_json). The list is an empty array when the job has no categories.

## Job details list time allocations by division

Lines 73-87

Asked as: time allocation by division; hours per division on job; job hours missing from division; allocations for job blank

The job details list the job's time allocations (allocations_json), each with its division ID, code and name and its hours. Hours count as 0 when not set. A time allocation whose division is missing is left out. A job with none gets an empty array.
