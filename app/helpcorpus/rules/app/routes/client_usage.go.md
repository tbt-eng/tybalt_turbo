<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_usage.go","sha256":"501d580c551db2475464755c06a1a73a5088f11ea931026fcf62773af4d79500","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_usage.go

## Client usage lists jobs that use the contact or invoicing profile

Lines 19-29

Asked as: jobs using a contact; which jobs use this contact; invoicing profile jobs; project contact flag; invoice recipient jobs

For a contact or an invoicing profile, Turbo lists the jobs that use it. Each job shows its number, description and status. It also shows whether the contact is the project contact (project_contact) and whether it is the invoice recipient (invoice_recipient). For an invoicing profile, only the invoice recipient flag can be true.

## Client usage lists project requests that use the contact or profile

Lines 31-40

Asked as: project requests for contact; job requests using profile; requests linked to invoicing profile; who created requests for client; request status in client usage

Turbo also lists the project requests (job requests) that use the contact or invoicing profile. Each request shows its description, status, and whether the contact is the project contact or the invoice recipient. Turbo does not send the creator and approver to the caller.

## Client usage counts the project requests it hides

Lines 42-49

Asked as: hidden requests count; why counts don't match; requests I can't see; missing requests in client; hidden_requests number

Turbo leaves out the project requests that the caller may not read, and returns their number as hidden_requests. The lists then still add up to the counts in the client details.

## Every user sees all jobs in client usage

Lines 51-57

Asked as: jobs visible to everyone; why I see all jobs; job list not filtered; other users' jobs showing; jobs ignore my permissions

Turbo does not filter the jobs in the client usage list by the caller's claims. Every signed-in user sees all of the jobs that use the contact or invoicing profile.

## Client usage lookup refuses a record from another client

Lines 59-63

Asked as: record not found for this client; contact belongs to other client; wrong client contact error; invoicing profile lookup fails; client ID mismatch

Turbo looks up a contact or an invoicing profile by its ID and checks that it belongs to the client in the path. If the record is not found, or belongs to a different client, Turbo refuses with "record not found for this client".

## Who sees project requests in client usage

Lines 81-97

Asked as: who can see project requests; requests I approve; project_authorization claim; approver sees request; missing job claim access

A caller with the job claim sees every project request in the usage list. Any other caller sees a request only if they created it. A caller with the project authorization claim (project_authorization) also sees the requests where they are the approver. The caller needs that claim for the approver case, and only checks it when they lack the job claim.
