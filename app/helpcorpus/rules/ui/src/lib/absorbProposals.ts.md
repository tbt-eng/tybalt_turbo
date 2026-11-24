<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/absorbProposals.ts","sha256":"e7f37b44041ec5a8eb49f224f72a84735f1cd60cfec1dc179559836d7874aac5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/absorbProposals.ts

## Absorb proposals cover clients, contacts and vendors only

Lines 8-9

Asked as: propose absorb for job; can't propose absorb for timesheet; which records can be absorbed; absorb proposal not allowed for expense

A person can propose an absorb only for clients, client contacts (client_contacts) or vendors. Other kinds of record cannot be proposed.

## Absorb proposal statuses

Lines 11-54

Asked as: absorb proposal status; what does pending mean; proposal shows declined; withdrawn or obsolete proposal

An absorb proposal has one of five statuses: pending, accepted, declined, withdrawn or obsolete. Turbo shows the status with its first letter in capitals, such as "Pending".

## Absorb proposal keeps a copy of each record's name

Lines 13-38

Asked as: deleted client name still shows; name of absorbed record; record name after absorb; absorb in progress not committed

Each record in a proposal keeps a copy of its name (record_name), because a committed absorb deletes the record. An absorb proposal is in progress (in_progress) when an absorb holder has absorbed the records and not yet committed.

## Who can absorb directly versus only propose an absorb

Lines 40-43

Asked as: why can't I absorb directly; only propose absorb no permission; absorb claim missing; need absorb permission to merge

A user with the absorb claim absorbs records directly (canAbsorbDirectly). A user without the absorb claim can only propose an absorb. The check is only whether the claims list includes "absorb".

## Collection labels shown for absorb proposals

Lines 45-50

Asked as: proposal shows raw collection name; collection label wrong on absorb; what does Contact label mean; vendor shown as odd name

Turbo labels a proposal's collection as "Client", "Contact" or "Vendor". For any other collection name, it shows the raw name.

## Review link goes to the absorb page of the target

Lines 56-63

Asked as: where to review absorb proposal; link to absorb proposal; open proposal for contact; absorb page for vendor link

A holder reviews a proposal on the absorb page of its target record, with ?proposal=<id> added. A contact's page is /clients/<client>/contacts/<target>/absorb. A vendor's page is /vendors/<target>/absorb. A client's page is /clients/<target>/absorb.

## Absorb page loads the proposal named in the URL

Lines 65-79

Asked as: absorb page shows no proposal; proposal not found on absorb page; blank absorb page; absorb link missing proposal id

An absorb page loads a proposal only when the URL has a proposal parameter. Without it, the page shows no proposal. If the server answers with a 400 to 599 error, the page shows "Absorb proposal not found." Any other failure is passed on as it is.
