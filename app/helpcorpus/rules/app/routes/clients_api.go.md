<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/clients_api.go","sha256":"9d89deebc9ebe3a6d7fb2fb208b69b2b39412591903161431e9bfd5201fd9f16","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/clients_api.go

## Client list shows how many jobs use each client

Lines 49-55

Each client in the client list carries a count of the jobs that reference it (referencing_jobs_count). The same count appears in the client details (lines 268-269), next to the referenced flag.

## Referenced client cannot be deleted

Lines 56-58

A client is referenced (referenced) when any job, job request, invoicing profile or note uses it. Turbo refuses to delete a referenced client.

## Contact usage counts each job or request once

Lines 81-96

A client contact's job and request counts include use as the project contact and use through an invoicing profile. A job or request that uses the contact in both roles counts once. Turbo also reports the number of invoicing profiles (profile_count) that use the contact.

## Client contacts on details include address and email preference

Lines 83-96

Each contact on client details shows the contact's address, phone and a no-email flag (no_email). It also shows counts of jobs, project requests and invoicing profiles that use the contact.

## Invoicing profiles on client details show usage and creator

Lines 98-111

Each invoicing profile listed on a client's details shows its name, contact, billing name, fax and invoicing instructions. It also shows how many jobs and project requests use it, and who created it and when.

## Client details show lead, balance date and address

Lines 113-135

The client details response includes the client's business development lead (id, given name, surname, email). It also includes the outstanding balance with its date (outstanding_balance_date), and the address, city, province or state, postal code, country and phone. The referencing jobs count and the referenced flag come with it.

## Client time divisions list committed time only

Lines 137-143

The time divisions list for a client shows the divisions that have committed time on the client's jobs. Each division has its hours. Turbo loads it apart from the client details because summing the hours reads every entry. Only the details tab shows it.

## Business development leads list

Lines 145-160

Turbo lists business development leads with their id, given name, surname and email.

## Client details need an active account

Lines 165-174

Turbo refuses to show one client's details unless the user's account is active. The message is "an active account is required to view client details". The details include invoicing profiles, and Turbo applies the same active-account rule that invoicing profiles have.

## Client details for an unknown client

Lines 175-178

When the client does not exist, the client details request fails with the message "client not found".

## Client details combine contacts and invoicing profiles

Lines 175-182

A client's details come back in one response. The response holds the client's address, phone, business development lead, outstanding balance and its date, contacts, invoicing profiles, the number of referencing jobs and the referenced flag.

## Client list returns all clients with contacts and balance

Lines 185-208

Without a client id, Turbo returns the list of clients. Each client has its contacts, referencing jobs count, outstanding balance and referenced flag. The list needs no active-account check in this file.

## Client time divisions endpoint

Lines 215-225

Turbo returns an empty list when the client has no divisions with committed time. A query failure gives "failed to load client time divisions".

## Empty contacts show as an empty list

Lines 238-254

A client with no contacts or no invoicing profiles returns empty lists in its details, not null.
