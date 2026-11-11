<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"cc734974accb24f7bf1db6c75d71e0d61e1b9874cc615cb8057be17fdc6b1511","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Absorbing clients repoints these related records

Lines 43-55

Asked as: merge two clients; absorb client into another; client duplicate merge; jobs moved to other client; client merged what happens to invoices

When a client is absorbed (merged) into another, Turbo repoints these references to the target client: client contacts (client), client notes (client), jobs (client and job_owner), project requests (job_requests: client and job_owner), client invoicing information (client) and client invoices (invoices: client).

## Absorbing client contacts repoints jobs, requests, invoicing

Lines 56-61

Asked as: merge two client contacts; contact duplicate merge; contact jobs moved to new contact; absorb contact into another; contact invoicing details moved

When a client contact is absorbed into another, Turbo repoints the contact on jobs (contact), project requests (job_requests: contact) and client invoicing information (contact) to the target contact.

## Absorbed client contacts must share one client

Lines 62-65

Asked as: merge contacts from different clients; cannot combine contacts of two clients; contacts belong to different client; merge contacts wrong client error; same client contacts merge

Client contacts that are absorbed together must all belong to the same client (parent constraint: clients, field client). Contacts of different clients cannot be merged.

## Client contact absorb ordering depends on clients absorb

Lines 66-70

Asked as: cannot commit contact merge; cannot undo client merge; undo client absorb blocked; contact merge pending error; client undo stuck pending contacts

Turbo blocks committing a client contact absorb while a clients absorb is pending. Turbo also blocks undoing a clients absorb while a client contact absorb is pending (DependsOn: clients). This stops an undo from restoring references to contacts that a committed absorb deleted.

## Absorbing vendors repoints purchase orders and expenses

Lines 71-77

Asked as: merge two vendors; vendor duplicate merge; purchase orders moved to other vendor; expenses switched to new vendor; combine vendor records

When a vendor is absorbed into another, Turbo repoints the vendor on purchase orders and expenses to the target vendor. Vendors have no parent constraint and no ordering dependency, so any vendors can be absorbed together and in any order.

## Only clients, client contacts and vendors can be absorbed

Lines 80-96

Asked as: merge projects; absorb not available for this; unknown collection error; merge jobs option missing; can't merge invoices

Absorb is configured only for clients, client_contacts and vendors. For any other collection Turbo returns the error "unknown collection: <name>".

## Commit and undo blockers follow the dependency

Lines 98-119

Asked as: why commit blocked; cannot undo absorb; undo merge not allowed; commit merge waiting on other; absorb blocked by dependency

Committing an absorb is blocked by the collection it depends on: client contacts depend on clients. Undoing an absorb is blocked by every collection that depends on it: clients are blocked by client contacts. Clients and vendors have no commit blockers.
