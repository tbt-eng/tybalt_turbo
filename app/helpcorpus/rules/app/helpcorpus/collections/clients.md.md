<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/clients.md","sha256":"95763bc8f9fa4de74b2a61c1be666dbb27b484e6324c4b04ba5de30127ecbc5b","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/clients.md

## Who can list and view clients

Lines 7-17

Asked as: list clients; see all clients; view client details; who can see clients; clients visible to everyone

Any signed-in person can list the clients and view one client. No claim is needed.

## Who can create a client

Lines 19-24

Asked as: add a new client; create client; cannot add client; procurement create client; new client form fails

Only procurement (the job claim) can create a client through the API. The person must also have an active admin profile (active = true).

## Who can update a client

Lines 26-31

Asked as: edit client; change client details; update client via api; cannot edit client; client update rejected

Only procurement (the job claim) can update a client through the API. The person must also have an active admin profile (active = true). The filter does not limit which fields a request changes.

## Who can delete a client and when deletion is blocked

Lines 33-40

Asked as: delete client; remove client; cannot delete client; client has jobs; deletion blocked

Only procurement (the job claim) with an active admin profile can delete a client through the API. The filter blocks deletion while a job, a project request (job_requests) or an invoicing profile record (client_invoicing_information) still refers to the client. Turbo's own code is not limited by this rule.

## Client name is required and must be unique

Lines 47-47

Asked as: client name too short; duplicate client name; name already taken; client name exists; unique client name

A client needs a name of at least 2 characters (name). Two clients cannot have the same name, because of a unique index on name.

## Client business development lead is one optional person

Lines 51-51

Asked as: business development owner; assign bd lead; client sales lead; second bd lead; bd contact for client

A client can have at most one business development lead (business_development_lead), who must be a user. The field is optional.

## Client alias length limit

Lines 58-58

Asked as: client alias too long; alias character limit; max alias length; client nickname length; shorter client alias

A client's alias (alias) can have at most 120 characters. The alias is optional.
