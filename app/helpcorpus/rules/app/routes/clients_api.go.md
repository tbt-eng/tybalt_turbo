<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/clients_api.go","sha256":"fded9d9e5b1011b3de41f652a14e965460b9419dc73cef02095800707a82b710","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/clients_api.go

## Client list returns clients with contacts in one query

Lines 14-177

Asked as: client list slow; list clients with contacts; clients list many queries; client list outstanding balance; referencing jobs count

The client list returns each client with its contacts, its referencing jobs count (referencing_jobs_count), its outstanding balance and a referenced flag. One SQL query (clients.sql) fetches all of it, so the list does not make a query for each client.

## Referenced flag shows when a client cannot be deleted

Lines 51-54

Asked as: cannot delete client; delete client refused; client won't delete; why can't I delete a client; client marked referenced

A client is marked referenced (referenced) when any job, project request, invoicing profile or note uses it. Turbo refuses to delete a referenced client. The list and the client details both carry this flag.

## Contact use counts each job or request once

Lines 76-91

Asked as: contact job count; contact counted twice; contact request count wrong; contact profile count; contact has no email

A contact's job count and request count on the client details include use as a project contact and use through an invoicing profile. A job or request that uses the contact in both roles counts once. Each contact also shows its profile count and a no-email flag (no_email).

## Client time divisions are a separate request

Lines 132-138

Asked as: client time divisions; divisions with hours for client; client hours by division; client division breakdown; no time on client

A separate request lists the divisions with committed time on the client's jobs, with each division's code, name and hours. The client details do not include them because summing the hours reads every time entry, and only the details tab shows them. A client with no committed time gets an empty list.

## Client details of an unknown client are refused

Lines 142-147

Asked as: client not found; client details missing; open unknown client; client id does not exist; client page empty

When a person opens the details of a client id that does not exist, Turbo answers "client not found".

## Client time divisions load error

Lines 184-194

Asked as: time divisions error; client divisions won't load; divisions load failed; client hours error; divisions query failing

When the query for a client's time divisions fails, Turbo answers "failed to load client time divisions".

## Client details include lead, balance, address, contacts and invoicing profiles

Lines 196-245

Asked as: client business development lead; client address and phone; client details balance; client invoicing profiles shown; who is the client lead

The client details show the business development lead (name and email), the outstanding balance and its date, the address fields and phone, the contacts, the invoicing profiles and the referencing jobs count. The contacts and the profiles come from the same query as the client.

## Client without contacts or profiles returns empty lists

Lines 207-223

Asked as: client has no contacts; no invoicing profiles on client; empty contacts list; client shows null contacts; client without profiles

When a client has no contacts or no invoicing profiles, the client details return an empty list for each, not null.
