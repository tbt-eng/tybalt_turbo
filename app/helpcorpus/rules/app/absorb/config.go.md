<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/config.go","sha256":"68ff817d9a02243b365972bd95041e9ea79dd73dcae6cedff4ef91cd2f617c7e","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/config.go

## Optional absorb settings: owned rows, checks, commit step

Lines 33-62

Asked as: rate sheet entries deleted; child rows when merging; extra check before merge; hook before merge deletes

A collection can add optional settings. Owned rows, such as rate sheet entries, are saved with the absorbed record, deleted with it and restored on undo, and are not moved to the target. A Prepare check runs inside the absorb transaction before anything changes. It can add ids, and it can refuse the absorb with a message for the user. An OnCommit step runs in the commit transaction before the action is deleted.

## Absorb refusal messages come from RefusedError

Lines 64-76

Asked as: merge refused message; why was my merge refused; error message when merging; check says cannot merge

When a collection's absorb check (Prepare) returns a refusal (RefusedError), Turbo shows the user its message text. The message is formatted by the check that refuses.

## Owned child rows are stored under a special key

Lines 78-81

Asked as: where are child rows stored; owned rows key; child records of merged item

Turbo keeps the owned child rows of an absorbed record in absorbed_records under the key "@owned_rows", grouped by table. The key cannot collide with a field name because field names cannot contain "@".

## Which collections have absorb after the main list

Lines 96-139

Asked as: which records can be merged; can I merge rate sheets; what can be absorbed

The built-in absorb list covers clients, client contacts and vendors. Rate sheets are added later by the hooks package, because their checks need the rate sheet comparison there.

## Absorbing clients moves these references to the target

Lines 99-111

Asked as: merge two clients; move jobs to another client; combine duplicate clients; client merge moves contacts

When clients are absorbed, Turbo moves the references in client contacts, client notes, jobs (client and job owner), project requests (client and job owner) and client invoicing profiles to the target client (RefConfigs). Clients have no parent, so any client can be absorbed into any other. Absorbing clients is refused while job editing is turned off (ChangesJobs).

## Absorbing client contacts: references moved and same-client rule

Lines 112-127

Asked as: merge two contacts; move contact to other contact; merge contacts from different clients

When client contacts are absorbed, Turbo moves the contact references in jobs, project requests and client invoicing profiles to the target contact. All contacts in one absorb must belong to the same client (ParentConstraint on clients, field client). Absorbing contacts is refused while job editing is turned off (ChangesJobs).

## Contact absorb waits for client absorb

Lines 123-126

Asked as: cannot commit contact merge; contact merge waiting on client; client merge blocks contacts

A client contacts absorb depends on clients (DependsOn). Turbo blocks the commit of a contact absorb while a client absorb is pending. Turbo also blocks the undo of a client absorb while a contact absorb is pending, so that undo never restores references to contacts that a committed absorb deleted.

## Absorbing vendors moves these references to the target

Lines 128-135

Asked as: merge two vendors; move purchase orders to other vendor; combine duplicate vendors

When vendors are absorbed, Turbo moves the vendor references in purchase orders, expenses and expense documents to the target vendor. Vendors have no parent and no dependency on another absorb. Absorbing vendors is not blocked by the job editing setting.

## Each collection can be registered for absorb only once

Lines 141-147

Asked as: app crashes on startup; absorb config already registered; startup panic merge setup

Registering an absorb configuration for a collection that already has one makes Turbo panic at startup with "absorb config for <name> is already registered".

## Absorb claim is the default permission

Lines 149-179

Asked as: no permission to merge; why can't I merge; merge button missing; cannot undo merge permission

To absorb, undo or commit in a collection, a person needs the absorb claim. This is so for collections that set no claim sets and for collections that have no absorb configuration. A collection with claim sets (ClaimSets) instead accepts a person who holds every claim of at least one set. Rate sheets are registered separately in the hooks package.

## Collections without claim sets need the absorb claim

Lines 149-157

Asked as: missing absorb permission; absorb claim needed; need absorb access to merge

A collection with no claim sets, or one not configured, requires the absorb claim ("absorb") to absorb, undo and commit.

## Claim sets: a user needs every claim of one set

Lines 159-179

Asked as: merge permission claim set; need all claims for merge; claim lookup error merge

When a collection sets claim sets, a user may absorb, undo and commit if they hold every claim of at least one set. The user does not need the absorb claim. A claim lookup error stops the check with an error.

## Reference moves listed by collection, sorted by name

Lines 181-195

Asked as: which collections reference this table; what does merge move in table; list collections moved by merge

Turbo finds the collections whose absorb moves references in a given table, and lists them sorted by name. Each collection appears once, even if it moves several columns in that table.

## Absorb on an unknown collection fails

Lines 197-204

Asked as: unknown collection error; merge fails unknown collection; collection not found for merge

Turbo fails an absorb request for a collection that has no absorb configuration, with the error "unknown collection: <name>".

## Absorb settings of a collection are looked up by name

Lines 206-213

Asked as: merge settings missing for name; unknown collection name lookup; reference list for collection

Turbo reads the reference list and parent constraint of a collection from its absorb configuration. An unknown collection gives the error "unknown collection: <name>".

## Undo of a parent absorb is blocked by dependent collections

Lines 215-225

Asked as: cannot undo client merge; undo merge blocked by contacts; why can't I undo merge

The collections that block undoing an absorb are those whose DependsOn names it. Undoing a clients absorb is blocked by client contacts.

## Commit is blocked until the parent collection's absorb is resolved

Lines 227-236

Asked as: commit blocked by parent merge; contact merge stuck; cannot commit until client merged

A collection's commit blockers are the collection it depends on (DependsOn). A collection with no dependency, or one not configured, has no commit blockers. Contact absorbs depend on clients.
