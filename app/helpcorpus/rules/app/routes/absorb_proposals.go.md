<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_proposals.go","sha256":"1ab4617a83b0147764b3a9b709b106d4e28a6b0d2e077f234c59239bc0a3eee2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_proposals.go

## Decided absorb proposals list shows at most 200

Lines 27-27

Asked as: old absorb proposals missing; decided proposals list cut off; only 200 decided proposals; why are older decided proposals gone

The decided list of absorb proposals returns at most 200 proposals (absorbProposalDecidedListLimit). The list is sorted newest decision first, so older decided proposals do not appear.

## Absorb holders review absorb proposals

Lines 66-69

Asked as: who can review absorb proposals; absorb reviewer access; see all absorb requests; absorb approver permission

Turbo treats a person as a reviewer of absorb proposals only when they hold the absorb claim. A reviewer sees every proposal and can decline one. Anyone else sees only their own proposals.

## Each absorb proposal lists its records by name

Lines 89-100

Asked as: records missing from absorb proposal; which records will be absorbed; record list order on proposal; empty absorb proposal records

Each absorb proposal in a list or in a single view comes with the records it proposes to absorb. Turbo sorts these records by record name, then by record ID. A proposal with no matching records still comes back with an empty record list.

## Who sees which absorb proposals in the list

Lines 104-118

Asked as: why can't I see other absorb proposals; see only my absorb requests; absorb proposal list visibility; who sees my absorb proposals

An absorb holder (the absorb claim) sees every absorb proposal. Anyone else sees only the proposals they requested.

## Absorb proposal list filters: pending, decided and invalid status

Lines 119-128

Asked as: absorb proposals status filter; show decided absorb proposals; pending absorb list default; status must be pending or decided

The list shows pending proposals by default, oldest first, or when status is pending. With status=decided it shows every proposal that is no longer pending, newest decision first, up to 200 (absorbProposalDecidedListLimit). Any other status gets the message "status must be pending or decided".

## Viewing one absorb proposal

Lines 136-152

Asked as: open one absorb proposal; absorb proposal not found; view absorb request details; can't open absorb proposal

An absorb holder or the proposal's requester can open a single absorb proposal. For anyone else, and for an unknown id, Turbo answers "absorb proposal not found".

## Any signed-in user can propose an absorb

Lines 155-226

Asked as: request to merge clients; propose absorbing vendors; who can submit absorb request; start absorb proposal

Any signed-in user can propose that clients, client contacts or vendors be absorbed into a target. The proposal starts as pending and records the requester. Each proposed record is saved with its name and marked open.

## Absorb proposal needs a reason

Lines 168-178

Asked as: absorb request missing reason; why is my absorb request rejected; say why these records should be absorbed; empty justification absorb

A proposal must have a reason. Turbo trims the reason and refuses an empty one with "say why these records should be absorbed". The reason check runs after absorb.ValidateProposal checks the target and records.

## Absorb proposal reason length limit

Lines 179-182

Asked as: absorb reason too long; absorb text character limit; reason too long for absorb; shorten absorb justification

A reason longer than absorb.ProposalTextMaxLength is refused with "the reason is too long". The limit is measured after trimming.

## Client contact proposals record the client

Lines 191-197

Asked as: client contact absorb client field; contact proposal client name; why does absorb show client; vendor absorb client blank

When the proposal is for client contacts, Turbo copies the target contact's client id and name onto the proposal. The proposal for clients or vendors has no client fields.

## Only pending proposals can be withdrawn or declined

Lines 232-242

Asked as: cancel absorb proposal already decided; already decided absorb request; change absorb proposal after decision; withdraw decided absorb

Withdrawing or declining a proposal that is no longer pending is refused with "this absorb proposal has already been decided". An unknown id gets "absorb proposal not found".

## Proposal already absorbed blocks withdraw and decline

Lines 243-249

Asked as: absorb already running cannot withdraw; cannot cancel absorb in progress; commit or undo the absorb; absorb request stuck being carried out

If an absorb action is already carrying out the proposal, it cannot be withdrawn or declined. Turbo says "this absorb proposal has been absorbed; commit or undo the absorb first".

## Only the requester can withdraw a proposal

Lines 253-271

Asked as: withdraw my absorb request; cancel my absorb request; only requester can withdraw; someone else's absorb request cancel

Only the person who made a pending proposal can withdraw it. Anyone else gets "only the requester can withdraw this absorb proposal". Withdrawing sets the status to withdrawn, with no comment.

## Only absorb holders can decline a proposal

Lines 276-284

Asked as: reject absorb request not authorized; not authorized to decline absorb; who can reject absorb proposals; deny absorb request permission

Only an absorb holder (the absorb claim) can decline an absorb proposal. Anyone else gets "you are not authorized to decline absorb proposals".

## Declining needs a comment

Lines 292-306

Asked as: decline absorb without comment; say why proposal is declined; decline comment too long; reject absorb needs note

A decline needs a comment. Turbo trims it and refuses an empty one with "say why the proposal is declined". A comment longer than absorb.ProposalTextMaxLength gets "the comment is too long".

## Declining notifies the requester

Lines 307-316

Asked as: requester not notified of decline; absorb rejection email; decline notification failed; notify requester absorb outcome

A decline sets the proposal to declined, with the decider and the comment. Turbo then queues a notification of the outcome to the requester. If queueing fails, the decline still stands and the failure is only logged.
