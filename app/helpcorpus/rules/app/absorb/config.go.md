<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"5ec0c855bd31c1b036f50308f443a795f802c7b927b9d4936d1255ac6b5f64c1","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Claim sets: any one complete set grants access

Lines 37-40

Asked as: partial claims access denied; which permission sets grant access; claim set complete; missing one permission; access refused despite some rights

If a collection defines claim sets, a user passes by holding every claim of at least one set. The absorb claim is not checked then. A user who holds only part of every set is refused.

## Rate sheet extras: owned rows, prepare, commit hook

Lines 42-62

Asked as: rate sheet entries lost; rate sheet revisions absorb; child rows deleted on absorb; rate sheet entries not moved; prepare hook extra ids

A collection can define owned child rows, such as the entries of a rate sheet. Turbo saves them with the absorbed record, deletes them with it and restores them on undo. It does not move them to the target. A prepare check can add ids to absorb, such as other revisions of a rate sheet, or refuse the absorb with a message. An on-commit hook runs in the commit transaction before the action is deleted.

## Absorbing clients or client contacts is refused while job editing is off

Lines 59-61

Asked as: cannot merge clients job editing off; client absorb blocked; job editing disabled merge; merge clients greyed out; turn on job editing to merge

Clients and client contacts are marked ChangesJobs, because their absorb moves references in jobs. Turbo refuses the absorb of these collections while job editing is turned off. Vendors are not marked, so the job editing setting does not affect them.

## A refused absorb shows the rule's own message

Lines 64-76

Asked as: absorb refused message; why was merge refused; custom refusal reason; error message when absorbing; refused error text

When a collection's absorb check (Prepare) refuses a request, it returns a refusal (RefusedError). Turbo shows the user the message that the check gives. The helper Refused builds that message from a format and its arguments.

## Absorbed records keep their owned rows under a special key

Lines 78-81

Asked as: absorbed record owned rows key; where are rate sheet entries after merge; absorbed records stored rows; at sign field name; owned rows saved

When Turbo absorbs a record that has owned rows, it saves them in absorbed_records under the key "@owned_rows", grouped by table. No field can have this name, because a field name cannot contain "@".

## Parent constraint: absorbed records must share the same parent

Lines 88-94

Asked as: merge contacts different clients; contacts from different client; same client required merge; parent mismatch absorb; contacts must share client

A parent constraint names a parent collection and a field. Records that Turbo absorbs together must have the same value in that field. Client contacts use it with the field client, in the clients collection.

## Which records can be absorbed

Lines 96-138

Asked as: cannot merge this type; merge not available for records; unknown collection absorb; which items can be merged; absorb not supported here

Turbo configures absorb for clients, client contacts and vendors. It also configures rate sheets, which are added in the hooks package. Any other collection is unknown, and its absorb configuration lookup fails with 'unknown collection'.

## Absorbing clients moves their linked records

Lines 99-111

Asked as: merge client moves jobs; client merge what changes; jobs repointed to new client; client notes moved on merge; merged client invoicing details

When clients are absorbed, Turbo repoints these records to the target client: client contacts, client notes, jobs (client and job owner), project requests (client and job owner) and client invoicing information. Absorbing clients changes jobs, so the absorb is refused while job editing is turned off (ChangesJobs). Clients have no parent and no dependency on another absorb.

## Absorbing client contacts: same client and ordering

Lines 112-122

Asked as: merge contact jobs moved; contact merge requests invoicing; same client contacts merge; contact project requests transferred; merge contacts ordering

Client contacts absorbed together must belong to the same client (parent constraint on client). Turbo repoints jobs, project requests and client invoicing information to the target contact. The absorb is refused while job editing is turned off.

## Client contact absorb depends on client absorb

Lines 123-126

Asked as: cannot undo client merge; client merge pending contact merge; contact merge blocks client undo; commit contact merge first; undo client absorb blocked

A client contact absorb cannot be committed while a clients absorb is pending. A clients absorb cannot be undone while a client contact absorb is pending. This stops an undo from restoring references to records that a committed contact absorb deleted.

## Absorbing vendors moves purchase orders and expenses

Lines 128-134

Asked as: merge vendors purchase orders; vendor expenses moved on merge; vendor merge what happens; repoint vendor orders; combine suppliers

When vendors are absorbed, Turbo repoints the vendor on purchase orders and on expenses to the target vendor. Vendors depend on no other absorb and have no parent constraint.

## Duplicate absorb registration is rejected

Lines 140-146

Asked as: absorb registered twice; panic duplicate config; absorb already registered; two absorb configs same collection; init registration error

Registering an absorb configuration for a collection that already has one causes a panic ('absorb config for ... is already registered'). Registration is allowed only in an init function.

## Absorb claim is the default permission

Lines 148-156

Asked as: default absorb permission; need absorb claim to merge; no claim sets permission; who can merge by default; merge permission missing

A collection with no claim sets, or one that is not configured, requires the absorb claim to absorb, undo and commit.

## Claim sets: a user needs every claim of one set

Lines 158-178

Asked as: which claim set counts; all claims in set required; partial permissions merge; first complete claim set; missing a claim cannot merge

A user may absorb, undo and commit in a collection when they hold all the claims of at least one claim set. A claim set that is missing one claim does not count. Turbo checks the sets in order and stops at the first complete one.

## Which collections an absorb of a table moves references in

Lines 180-194

Asked as: which collections reference table; what merges touch this table; collections moved by absorb; list references for table; absorb affected tables

Turbo lists the collections whose absorb moves references in a given table. A collection is listed once, even when it moves several columns of that table. The list is sorted by name.

## Looking up an absorb config for an unconfigured collection fails

Lines 196-212

Asked as: unknown collection error; no merge config for collection; lookup failed unconfigured; reference config not found; parent constraint missing error

When a collection has no absorb configuration, the lookup returns the error "unknown collection: <name>". The lookup of reference configs and the parent constraint returns the same error.

## Undo of a collection's absorb is blocked by dependent collections

Lines 214-224

Asked as: cannot undo client contact merge; dependent merge blocks undo; undo blocked pending merge; contacts pending block undo; undo waiting on dependent

A collection's absorb cannot be undone while a collection that depends on it (DependsOn) has a pending absorb. Client contacts depend on clients, so a pending client contact absorb blocks the undo of a clients absorb.

## Commit of an absorb is blocked by the collection it depends on

Lines 226-235

Asked as: cannot commit contact merge; commit blocked by clients; merge commit waiting; pending client merge stops commit; commit stuck dependency

A collection's absorb cannot be committed while the collection it depends on (DependsOn) has a pending absorb. A collection with no dependency has no commit blocker. Client contacts are blocked by clients.
