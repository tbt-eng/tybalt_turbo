<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_details.sql","sha256":"56ef37fdcd64213d14ee1c748ac62bfbe75ac02b34e6ec692412e54e40b9d816","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_details.sql

## Contact job use counts jobs by contact or invoicing profile

Lines 4-21

Asked as: contact job count; jobs per contact; contact has no jobs; why job count wrong for contact; jobs linked to contact

For each contact of the client, Turbo counts a job (job_count) when the job names the contact (jobs.contact) or when the job uses an invoicing profile of this client that names the contact. A job that matches both ways counts once.

## Contact request use counts project requests by contact or profile

Lines 12-22

Asked as: project request count contact; requests for a contact; contact request total; requests not showing for contact; double counted request

For each contact of the client, Turbo counts a project request (request_count) when the request names the contact or when it uses an invoicing profile of this client that names the contact. A request that matches both ways counts once.

## Contact profile count counts invoicing profiles naming the contact

Lines 23-23

Asked as: how many profiles per contact; contact profile number; profiles under contact; billing profiles for contact count

For each contact, Turbo counts the client's invoicing profiles that name the contact as their contact (profile_count).

## Contacts are sorted by first name, last name, then id

Lines 25-25

Asked as: contact list order; sort contacts alphabetically; contacts out of order; contact name ordering

Turbo lists the client's contacts by given name, then surname, ignoring case, and then by id.

## Invoicing profiles show creator name with email fallback

Lines 29-34

Asked as: profile creator name blank; who created profile; creator shows email; profile author missing name; created by shows email

Each invoicing profile shows its creator's full name (given name and surname). When the creator has no name, Turbo shows the creator's email. When there is no email either, Turbo shows an empty value.

## Invoicing profile job and request counts

Lines 30-31

Asked as: jobs using invoicing profile; profile job count; requests on billing profile; how many jobs per profile

For each invoicing profile, Turbo counts the jobs (job_count) and the project requests (request_count) that use it.

## Invoicing profiles are sorted by name then id

Lines 35-35

Asked as: invoicing profile list order; billing profiles sorted; profiles out of order; profile alphabetical order

Turbo lists the client's invoicing profiles by name, ignoring case, and then by id.

## Missing client alias, lead and billing name show as empty

Lines 40-52

Asked as: client alias blank; no lead on client; empty client alias; lead name missing; profile no billing name; client has no bd lead

When the client has no alias, no business development lead, or a profile has no billing name, Turbo shows an empty value. The lead's given name, surname and email also show empty when there is none.

## Outstanding balance sums the client's own jobs only

Lines 48-95

Asked as: client outstanding amount; unpaid balance for client; outstanding total wrong; owed by client; balance includes owner jobs; outstanding date

Turbo adds up the outstanding balance of jobs where the client is the job's client (outstanding_balance). The date shown is the latest outstanding balance date among these jobs. Both show 0 or empty when there are no such jobs. Jobs that the client only owns as job owner are not included.

## Referencing jobs count includes jobs owned by the client

Lines 67-87

Asked as: jobs where client is owner; client job total; referencing job count; jobs tied to client; job owner count

Turbo counts each job once when the client is the job's client or its job owner (referencing_jobs_count). The count is 0 when there are none.

## Client is referenced when jobs, requests, profiles or notes use it

Lines 71-75

Asked as: client cannot be deleted; client in use; why client marked referenced; referenced flag on client; client still used; clients list referenced

Turbo marks the client as referenced (referenced) when a job uses it, a project request names it as client or job owner, or a note exists for it. The profile check does not test the client of the profile row, so any invoicing profile of this client counts. The same uses make up the clients list's referenced flag.
