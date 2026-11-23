<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/clients_api.go","sha256":"e28899ecb59683337e93b646edd7ec3c8f6e40456b14cc4ca064c346b55d717a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/clients_api.go

## Client list shows each client with its contacts

Lines 16-52

Asked as: client list with contacts; list of all clients; clients outstanding balance; how many jobs per client

The client list gives each client's name (name), alias (alias), contacts and the number of jobs that reference it (referencing_jobs_count). It also gives the outstanding balance (outstanding_balance). Turbo loads clients and contacts in one query, not one query per client.

## Client list contact fields

Lines 28-34

Asked as: contact details in client list; contact phone missing from list; what fields show for contacts; contact email in list

Each contact in the client list shows only an id, a given name, a surname and an email. The client list does not show other contact details such as address or phone.

## Referenced flag blocks deleting a client

Lines 53-55

Asked as: cannot delete client; delete client fails; client is in use; why can't I remove client

A client is marked referenced (referenced) when any job, project request, invoicing profile or note uses it. Turbo refuses to delete a referenced client.

## Client details show business development lead and balance date

Lines 58-67

Asked as: client business development lead; who owns the client; balance as of date; client account manager

The client details show the client's business development lead (business_development_lead) with that person's given name, surname and email. They also show the date of the outstanding balance (outstanding_balance_date).

## Client details include address and phone

Lines 70-76

Asked as: client address and phone; where is client phone number; client postal code; client full address

The client details include the client's address, city, province or state, postal code, country and phone. The client list does not include these fields.

## Contact usage counts each job or request once

Lines 78-93

Asked as: contact job count; contact counted twice; how many jobs for contact; contact invoicing profile count

A client contact's job and request counts include use as the project contact and use through an invoicing profile. A job or request that uses the contact in both roles counts once. The contact also shows a count of invoicing profiles (profile_count).

## Client details contact fields and no-email flag

Lines 80-93

Asked as: contact no email checkbox; contact address in client page; contact usage numbers; client contact details page

Each contact in the client details has a full address, a phone, and a flag for no email (no_email). It also has counts of the jobs, project requests and invoicing profiles that use the contact (job_count, request_count, profile_count).

## Client details list invoicing profiles

Lines 95-108

Asked as: client billing profiles; client invoice profile list; billing name and fax; client invoice instructions

The client details list the client's invoicing profiles. Each profile shows its name, contact, billing name, fax, invoicing instructions, the number of jobs and project requests that use it, and its creator and creation date.

## Client time divisions list committed time

Lines 134-205

Asked as: divisions hours on client; client time by division; committed hours per division; client hours breakdown

A separate list shows the divisions with committed time on the client's jobs, with the hours of each division. It is separate from the client details because summing the hours reads every entry, and only the details tab shows it. This code does not check the account status for it.

## Client details need an active account

Lines 146-154

Asked as: cannot open client details; client details access denied; inactive account client view; account not active viewing client

Turbo shows one client's details only to a person with an active account. Otherwise it refuses with "an active account is required to view client details". The details include invoicing profiles, so this follows the read rule of the invoicing profiles collection.

## Client details for an unknown client

Lines 155-158

Asked as: client id not found; client page is empty; no such client; wrong client number

When the client id matches no client, Turbo answers "client not found" instead of showing details.

## Client list without an id

Lines 165-188

Asked as: show all clients; clients without id; full client list; client overview all

With no client id, Turbo returns the list of all clients. Each client has its alias, contacts, referencing jobs count, outstanding balance and referenced flag. This list does not check for an active account in this code.

## Client details contents

Lines 207-256

Asked as: what client page shows; client with no contacts; client summary page; client details empty contacts

The details of one client include its name, alias, address, phone and business development lead (name and email). They also include the outstanding balance with its date, the contacts, the invoicing profiles and the count of referencing jobs. A client with no contacts shows an empty contact list.
