<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"98c7f7541000703725607254c0fe24f4e2a7a11e932981da32ae817f4df9758a","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Absorbing clients repoints these records to the target client

Lines 43-54

When clients are absorbed (merged), Turbo repoints the client on client contacts, client notes, jobs, project requests (job_requests) and client invoicing information to the target client. It also repoints the job owner (job_owner) on jobs and project requests. Clients have no parent constraint and no dependency.

## Absorbing client contacts repoints jobs, requests and invoicing information

Lines 55-60

When client contacts are absorbed, Turbo repoints the contact on jobs, project requests (job_requests) and client invoicing information to the target contact.

## Absorbed client contacts must belong to the same client

Lines 61-64

Client contacts absorbed together must all have the same client (parent constraint on the client field of the clients collection). Turbo does not merge contacts of different clients.

## Client contact absorb waits on pending client absorb

Lines 65-69

Client contacts depend on clients (DependsOn). Turbo blocks the commit of a client contacts absorb while a clients absorb is pending. It also blocks the undo of a clients absorb while a client contacts absorb is pending.

## Absorbing vendors repoints purchase orders and expenses

Lines 70-76

When vendors are absorbed, Turbo repoints the vendor on purchase orders and expenses to the target vendor. Vendors have no parent constraint and no dependency, so other absorbs do not block them.

## Only clients, client contacts and vendors can be absorbed

Lines 79-95

Absorb configuration exists only for clients, client contacts and vendors. For any other collection, Turbo returns the error "unknown collection: <name>".

## Undo of an absorb is blocked by dependent collections

Lines 97-107

The collections whose DependsOn names this collection block undoing its absorb. Only client contacts depend on clients, so only the clients undo is blocked by them. Client contacts and vendors have no blockers.

## Commit of an absorb is blocked by the parent collection

Lines 109-118

The commit blockers of a collection are the one collection named in its DependsOn. Client contacts are blocked by clients. Clients, vendors and unknown collections have no commit blockers.
