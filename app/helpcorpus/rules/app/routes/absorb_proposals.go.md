<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_proposals.go","sha256":"3d1de1806e30a71b260100ffc915f7ec02837f91c4037902cf1c7c9c67e8347e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_proposals.go

## Decided proposals list shows at most 200 entries

Lines 28-28

Asked as: decided absorb list limit; old absorb history missing; only 200 decided absorbs; absorb history cut off

When the absorb proposal list is filtered to decided proposals, Turbo returns at most 200 (absorbProposalDecidedListLimit). They are the most recently decided first, ordered by decision date.

## Absorb access is the absorb claim

Lines 67-70

Asked as: who can review absorbs; absorb reviewer permission; can't see absorb requests; absorb access rights

Turbo treats a person as an absorb proposal reviewer only if they hold the absorb claim. Reviewers see all proposals and may decline them.

## Each absorb proposal lists its records by name

Lines 74-102

Asked as: which records get absorbed; absorb records order; records listed in absorb; absorb record sorting

Turbo returns every absorb proposal with the records that it would absorb. Those records are sorted by record name, then by record ID.

## Who sees absorb proposals in the list

Lines 105-120

Asked as: why I see others' absorbs; absorb list shows everyone; only my absorb requests; other people's absorb proposals

An absorb holder (the absorb claim) sees every person's absorb proposals in the list. Anyone else sees only the proposals they requested.

## Absorb proposal list filters by status

Lines 121-129

Asked as: pending absorbs default list; show decided absorbs; absorb status filter; status must be pending

The list shows pending proposals by default, oldest first (ordered by created). With status=decided it shows every proposal that is not pending, newest decision first, up to 200 (absorbProposalDecidedListLimit). Any other status gets the message "status must be pending or decided".

## Viewing a single absorb proposal

Lines 137-152

Asked as: can't open absorb proposal; absorb proposal not found; view someone's absorb request; absorb proposal missing

Only an absorb holder (the absorb claim) or the proposal's requester can view one proposal. Anyone else, and any unknown id, gets "absorb proposal not found".

## Any signed-in user can propose an absorb

Lines 156-227

Asked as: request an absorb; propose merging clients; submit absorb for vendor; suggest absorbing duplicate contact

Any signed-in user can create an absorb proposal for clients, client contacts or vendors. The proposal starts as pending, with the user as its requester. Each proposed record is saved with it as open.

## Absorb proposal needs a reason

Lines 169-182

Asked as: absorb reason required; blank absorb reason error; absorb reason too long; why absorb needs explanation

A new absorb proposal needs a reason. Turbo trims spaces from the reason first. A blank reason gets "say why these records should be absorbed", and a reason longer than absorb.ProposalTextMaxLength gets "the reason is too long".

## Proposal validity is checked before the reason

Lines 172-223

Asked as: absorb blank reason and invalid; absorb failed nothing saved; absorb validation error first; absorb save partly failed

Turbo first checks the collection, target and records with absorb.ValidateProposal, then the reason. A proposal that fails validation gets that error, even if the reason is blank. The proposal and its records are saved in one transaction, so a failure saves none.

## Contact proposals record the client

Lines 189-198

Asked as: client contact absorb client; which client for contact absorb; contact absorb client name; absorb contact's client recorded

For a proposal on client contacts (client_contacts), Turbo also stores the target contact's client and the client's name on the proposal. It stores the target's and record names on the proposal as well.

## Only pending proposals can be decided

Lines 233-243

Asked as: absorb already decided; can't change absorb proposal; withdraw decided absorb; absorb proposal already finished

Withdrawing or declining works only on a pending proposal. An unknown id gets "absorb proposal not found". A proposal already decided gets "this absorb proposal has already been decided".

## Proposal being absorbed cannot be withdrawn or declined

Lines 244-250

Asked as: absorb in progress can't withdraw; commit or undo absorb; stuck absorb proposal; absorb being carried out

A pending proposal that an absorb action is carrying out cannot be withdrawn or declined. Turbo says "this absorb proposal has been absorbed; commit or undo the absorb first".

## Only the requester can withdraw a proposal

Lines 256-273

Asked as: cancel my absorb request; withdraw absorb proposal; only requester can withdraw; take back absorb request

Only the requester can withdraw a pending absorb proposal. Anyone else gets "only the requester can withdraw this absorb proposal". A withdrawal sets the status to withdrawn and records the requester as the decider. It succeeds with "absorb proposal withdrawn".

## Only the requester can delete a closed proposal

Lines 275-304

Asked as: delete absorb proposal; remove old absorb request; can't delete pending absorb; absorb proposal deletion

Only the requester can delete an absorb proposal. Anyone else gets "only the requester can delete this absorb proposal". A pending proposal cannot be deleted: Turbo says "withdraw this absorb proposal before you delete it". A proposal in any other status can be deleted, and its record rows are deleted with it.

## Only absorb holders can decline a proposal

Lines 308-316

Asked as: decline absorb not allowed; reject absorb request; who can reject absorbs; not authorized absorb decline

Only an absorb holder (the absorb claim) can decline an absorb proposal. Anyone else gets "you are not authorized to decline absorb proposals".

## Declining needs a comment and notifies the requester

Lines 324-348

Asked as: reject absorb comment required; requester notified of decline; decline absorb comment too long; absorb rejection email failed

Declining needs a comment. Turbo trims spaces from it. A blank comment gets "say why the proposal is declined", and a comment longer than absorb.ProposalTextMaxLength gets "the comment is too long". Turbo then sets the status to declined and queues a notification of the outcome to the requester. If queueing the notification fails, the decline still stands.
