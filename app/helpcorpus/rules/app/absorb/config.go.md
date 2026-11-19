<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"a2ea4c6e05c176ba6fb931d819fbb91d788b4fffc893b3704b307d5a99a06dd2","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Collections with claim sets accept any one complete set

Lines 37-40

Asked as: enough permissions to merge; missing one permission for merge; partial rights not enough; which permission set works

If a collection defines claim sets, a user may absorb, undo and commit when they hold every claim of at least one set. Holding only part of every set is not enough.

## Owned child rows are deleted and restored, not moved

Lines 42-47

Asked as: child rows lost after merge; rate sheet entries disappeared; entries deleted when merging; undo brings back entries

For collections with owned rows, such as the entries of a rate sheet, the child rows of an absorbed record are saved in the absorb action. They are deleted with the record and restored on undo. They are not moved to the target.

## Absorb refusals show a message to the user

Lines 49-53

Asked as: why was my merge refused; merge blocked with message; other revisions added to merge; merge refused reason

A collection's Prepare check runs before anything changes. It may add ids to absorb, such as other revisions of a rate sheet, or refuse the absorb with a message for the user.

## Owned rows are stored under a reserved key

Lines 78-81

Asked as: @owned_rows key; where owned rows are stored; absorbed_records owned rows; reserved key for entries

An absorbed record keeps its owned child rows in absorbed_records under the key "@owned_rows". No field name can contain "@", so the key never collides with a field.

## Clients absorb moves references in eight places

Lines 99-112

Asked as: merging clients updates jobs; client merge repoints invoices; what changes when merging clients; merge clients moves contacts notes

When clients are absorbed, Turbo repoints references to the target in client contacts, client notes, jobs (client and job owner), project requests (client and job owner), client invoicing information and invoices. Clients have no parent, so any clients can be absorbed together.

## Client absorb is refused while job editing is off

Lines 110-110

Asked as: merge clients fails job editing off; can't merge clients jobs locked; client merge refused editing disabled; turn on job editing to merge

Absorbing clients changes references in jobs (ChangesJobs). Turbo refuses the absorb while job editing is turned off.

## Client contacts absorb must share the same client

Lines 113-123

Asked as: merge contacts from different clients; contacts belong to different client; combine contacts across clients error; contacts must share a client

Client contacts that are absorbed together must all belong to the same client (parent constraint on clients, field client). Absorbing moves the references in jobs, project requests and client invoicing information to the target contact.

## Client contact absorb is refused while job editing is off

Lines 123-123

Asked as: merging contacts won't work; contact merge blocked jobs read only; job editing off contacts merge; can't merge contacts right now

Absorbing client contacts changes references in jobs (ChangesJobs). Turbo refuses the absorb while job editing is turned off.

## Client contacts absorb ordering depends on clients absorb

Lines 124-128

Asked as: can't undo client merge contacts pending; commit contact merge while client merge; client merge waiting on contacts; undo clients stuck

A client contacts absorb cannot be committed while a clients absorb is pending. A clients absorb cannot be undone while a client contacts absorb is pending. This stops an undo from restoring references to contacts that a committed absorb deleted.

## Vendor absorb moves purchase orders and expenses

Lines 129-135

Asked as: merge vendors purchase orders; vendor merge moves expenses; combine vendors where do POs go; vendor duplicate merge what changes

When vendors are absorbed, Turbo repoints the vendor on purchase orders and expenses to the target. Vendors have no parent, no ordering dependency, and no extra claim sets.

## Rate sheet absorb rules are registered elsewhere

Lines 136-147

Asked as: rate sheet merge not configured here; rate sheet merge hooks; why rate sheet merge works differently; rate sheet merge checks location

Rate sheet absorb is not configured in this file. The hooks package registers it, and its checks and commit steps are not shown here. Registering a collection twice is a programming error.

## Registering a duplicate absorb config is a startup error

Lines 141-147

Asked as: already registered panic absorb; startup crash absorb config; turbo won't start duplicate absorb; absorb config panic message

Registering an absorb config for a collection that already has one makes Turbo panic with "absorb config for <name> is already registered".

## Absorb needs the absorb claim by default

Lines 149-157

Asked as: need absorb permission to merge; merge button missing permission; can't merge without absorb claim; who is allowed to merge

To absorb, undo or commit, a user needs the absorb claim. This applies to any collection that has no claim sets of its own or is not configured.

## Unconfigured collections need the absorb claim

Lines 152-157

Asked as: merge unconfigured collection permission; absorb with no claim sets; merge needs single absorb claim; collection without merge rules permission

If a collection is not configured, or has no claim sets, Turbo requires the single absorb claim to absorb, undo or commit (RequiredClaimSets).

## Absorb claim check needs every claim in one set

Lines 159-179

Asked as: permission lookup error on merge; some permissions for merge not enough; merge permission check fails; need every permission in one set

When a collection has claim sets, Turbo allows absorb, undo and commit to a user who holds every claim of at least one set (HasRequiredClaims). Holding only some claims of each set does not pass. A claim lookup error fails the check.

## Absorb of an unknown collection is an error

Lines 197-204

Asked as: unknown collection merge error; merge collection not found; absorb wrong collection name; merge config missing for collection

Asking for the absorb configuration of a collection that is not in the list returns the error "unknown collection: <name>" (GetConfig).

## Undoing an absorb is blocked by dependent collections

Lines 215-225

Asked as: can't undo merge contacts depend; undo merge blocked by dependents; why client merge undo blocked; undo blocked by related records

Turbo blocks undoing a collection's absorb when a collection whose DependsOn points to it has an absorb (GetUndoBlockers). For clients, that collection is client contacts.

## Committing an absorb is blocked while the parent collection has one

Lines 227-236

Asked as: commit merge waiting on parent; can't finalize child merge; parent merge pending commit blocked; merge commit stuck on parent

A collection that depends on another (DependsOn) cannot commit its absorb until the parent collection's absorb is resolved. Collections with no dependency, or unknown collections, have no commit blockers (GetCommitBlockers).
