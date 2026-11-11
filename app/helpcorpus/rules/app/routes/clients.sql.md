<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/clients.sql","sha256":"de00d32dee2429cbc98cfc476331ce59cff2c4ff8b88bf0c91bfbcdf6498dc16","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/clients.sql

## Referenced flag shows when a client cannot be deleted

Lines 2-54

Asked as: cannot delete client; client delete refused; why can't I delete client; client is in use; greyed out delete client

A client is referenced (referenced) when any job, project request (job request), invoicing profile, client note or client invoice uses it. A job or project request uses it as its client or as its job owner. The clients delete rule and delete hook refuse to delete a referenced client, so the UI uses this flag to show that a delete will be refused.

## Client list shows contacts for each client

Lines 9-31

Asked as: client contacts list; see client contacts; contacts missing for client; client's people list; show contacts per client

The client list returns each client with its id, name and alias. It adds a list of the client's contacts (contacts_json) with id, given name, surname and email. A client with no contacts gets an empty list.

## Client alias is blank when not set

Lines 12-12

Asked as: client has no alias; blank client alias; alias empty; client short name missing; no alias showing

The client list shows an empty text for a client that has no alias (alias).

## Referencing jobs count includes jobs the client owns

Lines 32-43

Asked as: jobs count wrong for client; how many jobs client owns; client job total zero; count of client's jobs; job owner not counted

The count of referencing jobs (referencing_jobs_count) counts each job once when the client is the job's client or the job's owner (job_owner). A client with no such jobs shows 0. Jobs with an empty client or owner value are ignored for that side.

## Outstanding balance counts only jobs billed to the client

Lines 33-43

Asked as: client balance wrong; client owes money total; owned jobs in balance; outstanding amount for client; client balance zero

The outstanding balance of a client adds up the outstanding balance of the jobs where the client is the job's client. Jobs that the client only owns add nothing to the balance. A client with no such jobs shows 0.

## Client list can be limited to one client

Lines 55-56

Asked as: find one client; filter client list by id; show a single client; client list sorted by name; list all clients

The client list returns one client when it gets a client id. It returns all clients when the id is missing or empty. Turbo sorts the list by client name.
