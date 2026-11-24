<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_usage.go","sha256":"fff7df7c41abba0fa74b9d27afa484bdfce1cbcabd35a8cf53161fc8e3714628","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_usage.go

## Client usage list needs an active account

Lines 60-67

Asked as: inactive account client usage; can't view client details; account not active error; client usage access denied

Turbo refuses to list what uses a client's contact or invoicing profile when the caller's account is not active. It shows the message "an active account is required to view client details".

## Contact or invoicing profile must belong to the client in the path

Lines 68-72

Asked as: record not found for this client; contact belongs to other client; wrong client contact error; invoicing profile not found

Turbo looks up the contact or invoicing profile by its record ID. If the record does not exist, or its client (client) differs from the client in the path, Turbo refuses with "record not found for this client".

## Usage lists jobs and project requests for one contact or profile

Lines 74-88

Asked as: which jobs use this contact; jobs linked to invoicing profile; project contact or invoice recipient; list jobs for a contact

For a contact (client_contacts), Turbo lists the jobs and project requests (job requests) that use that contact. For an invoicing profile (client_invoicing_information), it lists those that use that profile. Each entry says whether the contact is the project contact (project_contact) or the invoice recipient (invoice_recipient). For an invoicing profile only the invoice recipient flag is true. Which records match comes from the SQL files, which this file does not show.

## Every active user sees all jobs in client usage

Lines 81-84

Asked as: why do I see all jobs; everyone sees all jobs; jobs not filtered by my access; client usage shows every job

Turbo does not filter the jobs in the client usage list by the caller's claims. Any user with an active account sees every job that uses the contact or profile.

## Who sees project requests in client usage

Lines 90-106

Asked as: who can see project requests; requests I created only; see requests I approve; job claim sees all requests

A person with the job claim sees every project request in the client usage list. Otherwise a person sees a request only if they created it. A holder of the project authorization claim (project_authorization) also sees requests where they are the approver. This follows the read rule of job_requests.

## Approver access needs the project authorization claim

Lines 94-101

Asked as: approver can't see request; approver without permission; why approver sees nothing; project authorization missing

Turbo checks the project authorization claim (project_authorization) only when the caller lacks the job claim. An approver without that claim does not see requests just because they are the approver.

## Hidden project requests are counted, not listed

Lines 100-107

Asked as: hidden requests count; requests missing from list; counts don't match client details; why are some requests hidden

Turbo leaves out project requests the caller may not read. It counts them in hidden_requests, so the lists still add up to the counts in the client details. Creator and approver IDs are not sent to the caller.
