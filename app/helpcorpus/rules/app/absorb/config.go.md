<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"6978bedd588b3a475a6c4a963d8d548b5a8c8efbe8b9328549bfc379a5f1ca0f","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Which records change when clients are absorbed

Lines 43-54

Asked as: merge two clients; absorb client into another; client records moved to other client; combine duplicate clients

When a client is absorbed into another client, Turbo repoints these records to the target client: client contacts (client), client notes (client), jobs (client and job_owner), project requests (job_requests: client and job_owner), and client invoicing information (client). Clients have no parent constraint, so any clients can be absorbed together.

## Which records change when client contacts are absorbed

Lines 55-60

Asked as: merge two contacts; client contact duplicate; move jobs to other contact; combine contact records

When a client contact is absorbed into another contact, Turbo repoints these records to the target contact: jobs (contact), project requests (job_requests: contact) and client invoicing information (contact).

## Absorbed client contacts must belong to the same client

Lines 61-64

Asked as: merge contacts from different clients; cannot merge contacts; contacts belong to different client; contact merge refused

Client contacts that are absorbed together must all have the same client (parent constraint on the client field of clients). Contacts of different clients cannot be merged.

## Client contact absorb waits on client absorb

Lines 65-68

Asked as: cannot commit contact merge; client merge pending blocks; undo client merge blocked; contact absorb stuck

Client contacts depend on clients (DependsOn). Turbo blocks committing a client contact absorb while a client absorb is pending. Turbo also blocks undoing a client absorb while a client contact absorb is pending. This keeps an undo from restoring references to contacts that a committed absorb deleted.

## Which records change when vendors are absorbed

Lines 70-77

Asked as: merge two vendors; vendor duplicate combine; move purchase orders to other vendor; vendor expenses reassigned

When a vendor is absorbed into another vendor, Turbo repoints the vendor on purchase orders, expenses and expense documents (expense_documents) to the target vendor. Vendors have no parent constraint and no dependency on other absorbs.

## Only clients, client contacts and vendors can be absorbed

Lines 80-96

Asked as: unknown collection error; merge jobs or projects; cannot absorb this collection; absorb not allowed for

Absorb is configured only for clients, client_contacts and vendors. Any other collection is refused with the error 'unknown collection: <name>'.

## Absorb dependency checks for undo and commit

Lines 98-119

Asked as: cannot undo absorb; absorb commit blocked; why can't I undo merge; absorb dependency blocking

Undoing an absorb is blocked by every collection whose DependsOn is this collection (for clients, that is client_contacts). Committing an absorb is blocked only by the collection it depends on. Collections with no dependency, or unknown collections, have no commit blockers.
