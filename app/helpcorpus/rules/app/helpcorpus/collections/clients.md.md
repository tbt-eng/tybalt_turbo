<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/clients.md","sha256":"cf98e16def50f7e7f7decf2800758a986230c44e15c41997eae77b722678020b","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/clients.md

## Any signed-in person can list and view clients

Lines 7-17

Asked as: view client list; see clients; can't open client; not signed in clients

Any signed-in person can list clients and view a single client. A person who is not signed in cannot do either.

## Who can create or update a client

Lines 19-31

Asked as: edit client details; who can change client; client form locked; add new client

Only procurement (the job claim) whose admin profile is active (active = true) can create or update a client. A signed-in person without the job claim, or without an active admin profile, is refused.

## Who can delete a client and when deletion is blocked

Lines 33-40

Asked as: remove a client; can't delete client; delete client blocked; client still has jobs

Only procurement (the job claim) with an active admin profile can delete a client. Turbo refuses the deletion when a job (jobs), a project request (job_requests) or an invoicing profile record (client_invoicing_information) still refers to the client. Turbo's own server code is not limited by this rule.

## Client name must be unique and at least 2 characters

Lines 47-64

Asked as: duplicate client name; client name taken; name already exists; client name too short

A client's name is required and must have at least 2 characters. Two clients cannot have the same name (unique index on name).

## Client business development lead is required

Lines 51-51

Asked as: client account owner; who manages client; missing sales contact; bd lead required

A client must have one business development lead (business_development_lead), who is a user. The field is required and allows only one person.

## Client alias has a 120-character limit

Lines 52-58

Asked as: client alias too long; alias character limit; client nickname length; optional client address

A client's alias is optional and can have at most 120 characters. The address, city, province or state, postal code, country and phone fields are optional.
