<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/clients.sql","sha256":"f9291c04005e072647410a23f1ee5cab0e5068c47445cb90e9523c016353b756","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/clients.sql

## Referenced flag marks clients that cannot be deleted

Lines 2-4

A client is referenced (referenced) when any job, job request, invoicing profile or note uses it. A job or job request uses a client as its client or as its job owner. The comment says the client delete rule and delete hook refuse to delete referenced clients. This file does not show that refusal itself.

## Client list shows alias as empty text when none

Lines 12-12

The client list returns an empty string for a client's alias (alias) when the client has no alias.

## Client list shows contacts, or an empty list

Lines 13-31

Each client in the list carries its contacts (id, given name, surname, email) as JSON. A client with no contacts gets an empty list ('[]').

## Referencing jobs count includes jobs the client owns

Lines 14-14

The count of referencing jobs (referencing_jobs_count) counts each job once for a client when the client is the job's client or its job owner (job_owner). A client with no such jobs gets 0. Jobs with an empty client or owner value do not count for that value.

## Outstanding balance counts only jobs billed to the client

Lines 15-15

The outstanding balance (outstanding_balance) of a client adds up the outstanding balance of jobs where the client is the job's client. Jobs the client only owns (job owner) add nothing to the balance. A client with no such jobs shows 0.

## Referenced flag covers job requests, invoicing profiles and notes

Lines 16-52

Referenced is true if the client is the client or job owner of any job request, or appears in any invoicing profile (client_invoicing_information) or client note. It is also true if any job uses the client.

## Job counts and balances use jobs only, not job requests

Lines 32-43

Job counts and outstanding balances come only from the jobs table. Job requests do not add to either figure.

## Client list can be filtered by one client id

Lines 53-54

When an id is given, the list returns only that client. With no id or an empty id, it returns all clients. Clients are sorted by name.
