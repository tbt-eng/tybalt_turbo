<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_details.sql","sha256":"365d5071c433b02a62c066a24d3285a3d6bd9bfc3383d15e9abbea52c3e53cd0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_details.sql

## Contact usage counts jobs once per contact

Lines 4-11

Asked as: contact job count doubled; same contact counted twice jobs; job count for contact; contact jobs counted twice

On the client details page, a contact's job count (job_count) counts each job once, even when the job's contact and the job's invoicing profile contact are the same person. The count includes jobs that name the contact directly and jobs that use an invoicing profile of this client whose contact is that person.

## Contact usage counts project requests once per contact

Lines 12-18

Asked as: contact request count wrong; project requests counted twice; requests per contact; contact requests duplicate

On the client details page, a contact's request count (request_count) counts each project request once. A request counts when it names the contact directly, or when it uses an invoicing profile whose contact is that contact. A request that does both counts once. Only contacts of this client are listed.

## Contact list order and counts on client details

Lines 19-26

Asked as: contact list order; sort contacts on client page; how contacts are sorted; contact profile count

The client details page lists the client's contacts sorted by given name, then surname, ignoring case, then id. Each contact shows its job count, project request count and the number of invoicing profiles (profile_count) that use it.

## Invoicing profile creator name fallback

Lines 27-36

Asked as: profile creator shows email; who created invoicing profile; creator name missing; invoicing profile order

Each invoicing profile on the client details page shows the creator's name as given name plus surname. If both are blank, it shows the creator's login email. If that is missing too, it shows an empty string. Profiles are sorted by name, ignoring case, then id.

## Invoicing profile job and request counts

Lines 30-31

Asked as: how many jobs use profile; profile job count; billing profile usage; profile request total

Each invoicing profile on the client details page shows how many jobs and how many project requests use it (job_count, request_count).

## Client details show empty text instead of missing values

Lines 37-66

Asked as: blank fields on client page; client page shows nothing; missing alias on client; zero balance no jobs

On the client details page, a missing alias, business development lead, lead name or lead email shows as an empty string. A profile's missing billing name also shows as empty. A client with no jobs shows an outstanding balance of 0 and an empty balance date.

## Outstanding balance sums the client's own jobs only

Lines 48-49

Asked as: client outstanding balance wrong; owed amount includes other client; job owner balance included; client balance total

The client's outstanding balance (outstanding_balance) is the sum of the outstanding balances of jobs whose client is this client. The balance date is the latest outstanding balance date among those jobs. Jobs that this client only owns (job_owner) do not count.

## Business development lead details come from profile and login

Lines 50-52

Asked as: lead email missing; business development lead name; sales lead details; lead shown on client

The lead's given name and surname come from the lead's profile. The lead's email comes from the lead's login record (users).

## Referencing jobs count includes jobs owned by the client

Lines 67-88

Asked as: referencing jobs count; jobs owned by client counted; client job count double; jobs linked to client total

The referencing jobs count (referencing_jobs_count) counts each job once when its client is this client or its job owner (job_owner) is this client. A job that matches both counts once.

## When a client counts as referenced

Lines 68-76

Asked as: client marked referenced; why client cannot be deleted; client referenced status; client has links

A client counts as referenced (referenced) when any of these is true: a job uses it as client or job owner; a project request names it as client or job owner; it has any invoicing profile; it has a note; or it has an invoice. The invoicing profile check looks only at the profiles already filtered to this client.
