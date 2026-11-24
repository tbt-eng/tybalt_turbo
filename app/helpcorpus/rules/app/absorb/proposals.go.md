<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/proposals.go","sha256":"7c1fa464a5f3a736db316693c70ff6b6c4ee1f7aa63fc7f26cab4872ea67bc92","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/proposals.go

## Who proposes and who carries out an absorb

Lines 18-20

Asked as: request to merge records; who can propose an absorb; ask to combine duplicates; propose merging clients

Any signed-in user can propose that records be absorbed into a target. A person with absorb access carries the proposal out with the normal absorb.

## Absorb proposal statuses

Lines 25-29

Asked as: proposal status meaning; pending or accepted proposal; what states can a proposal have; proposal stuck pending

An absorb proposal has one of five statuses: pending, accepted, declined, withdrawn or obsolete. A proposal starts as pending. A decision closes it with one of the other four statuses.

## Deciding a proposal frees its records

Lines 31-33

Asked as: release records after decision; records still locked after declining; free up locked duplicates; withdrawn proposal records reusable

When a pending proposal is closed with accepted, declined, withdrawn or obsolete, Turbo marks its records as no longer open, so other proposals can use them. It records who decided (which may be blank) and when. A decision comment is cut to 1000 characters.

## Which records can have an absorb proposed

Lines 35-48

Asked as: can't propose absorb for rates; merge invoices or jobs; absorb not allowed for this type; can't propose absorb of

Staff can propose an absorb only for clients (clients), contacts (client_contacts) and vendors (vendors). For any other collection Turbo refuses with "an absorb of <collection> can't be proposed". Rate sheets are left out because their absorb has its own rules.

## Names of contacts, clients and vendors in proposals

Lines 66-88

Asked as: proposal shows blank contact name; contact shown only as email; vendor alias in brackets; client label in proposal list

A proposal stores each record's name, because a committed absorb deletes the record. A contact shows as given name and surname, followed by the email in parentheses when it has one. A contact with no name shows only the email. A vendor shows its name, followed by its alias in parentheses when it has one. A client uses the client label (ClientLabel).

## Proposal needs at least one record to absorb

Lines 99-101

Asked as: absorb with nothing selected; empty merge list; no duplicates chosen; submit proposal without records

A proposal with no records to absorb is refused with "choose at least one record to absorb".

## Kept record cannot also be absorbed

Lines 103-106

Asked as: keep record also merged; target listed as duplicate; kept client in absorb list; merge a client into itself

If the record to keep (the target) is also in the list of records to absorb, Turbo refuses with "the record to keep can't also be absorbed".

## Same record cannot be listed twice

Lines 107-109

Asked as: same client twice in merge; duplicate entry in absorb list; record added twice; double selected duplicate

If a record appears twice in the list of records to absorb, Turbo refuses with "a record is listed twice".

## Missing records stop the proposal

Lines 113-130

Asked as: deleted client in proposal; record already removed error; target no longer exists; merge with deleted record

If the record to keep no longer exists, Turbo refuses with "the record to keep no longer exists". If a record to absorb no longer exists, Turbo refuses with "a record to absorb no longer exists".

## Records must share the same parent

Lines 132-138

Asked as: contacts from different clients merge; merge records different parent; absorb across different clients; same parent required for merge

For a collection with a parent constraint (ParentConstraint in Configs), every record to absorb must have the same parent field value as the record to keep. Otherwise Turbo refuses with "every record must belong to the same <parent>", where the parent is the parent collection name without its final "s". Which collections have a constraint is set in Configs, outside this file.

## Record already in another open proposal

Lines 140-160

Asked as: client already in a proposal; duplicate proposal for record; record busy in another merge; already proposed for absorb

A record can be in only one open proposal. This check covers both the record to keep and the records to absorb. If any of them is already in another open proposal, Turbo refuses with "<name> is already in another open absorb proposal". When a proposal is being carried out, its own records are excluded from this check.

## Cannot absorb a record that another pending proposal keeps

Lines 161-174

Asked as: record is kept by another proposal; cannot merge a kept client; target of another open merge; blocked by other pending absorb

If a record to absorb is the record to keep of another pending proposal, Turbo refuses with "another open absorb proposal keeps <name>". When a proposal is being carried out, that proposal is excluded from the check.

## Decision comment is cut to 1000 characters

Lines 193-212

Asked as: decision note too long; comment cut off; long reason truncated; character limit on comment

Turbo saves a decision comment only when one is given. It cuts a longer comment to 1000 characters, counting characters and not bytes.

## Committing an absorb accepts its proposal

Lines 231-254

Asked as: committing absorb closes proposal; absorb auto accepts proposal; who gets recorded as decider; finishing merge marks accepted

When an absorb action that carries out a proposal is committed, Turbo marks that proposal accepted if it is still pending. The person who commits is recorded as the decider. This applies only to clients, contacts and vendors.

## Committed absorb makes overlapping proposals obsolete

Lines 256-316

Asked as: other proposals went obsolete; absorb cancelled my proposal; why did my proposal become obsolete; merge killed other requests

When an absorb is committed, every other pending proposal that names an absorbed record, as its record to keep or as a record to absorb, becomes obsolete. Its comment reads "Another absorb removed <names>, so this proposal can't be carried out." Turbo returns the closed proposals so their requesters can be told.

## Client absorb moves pending contact proposals to kept client

Lines 318-339

Asked as: contact proposals changed client; contact moved to new client after merge; pending contact now under kept client; client absorb updates contact proposals

When a client absorb is committed, Turbo changes every pending contact proposal that belonged to an absorbed client. Those proposals then take the kept client's id and name.
