<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_proposals.go","sha256":"3a8ce33090bd3172f8de610619c8b5adbaf097ab49ffea59acdc26515e01c7ae","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_proposals.go

## Pending and decided absorb proposal lists

Lines 28-129

Asked as: list of absorb proposals; pending absorb requests; show decided proposals; filter proposals by status; old merge requests

The list shows pending proposals by default or with status=pending, oldest first. With status=decided it shows every proposal that is not pending, newest decision first, and at most 200 (absorbProposalDecidedListLimit). Any other status gives the message "status must be pending or decided".

## Who sees which absorb proposals in the list

Lines 105-118

Asked as: which absorb proposals can I see; only see my own proposals; why can't I see all merges; absorb proposals visibility; see other people's proposals

A person with the absorb claim (the absorb claim) sees every absorb proposal. Anyone else sees only the proposals they requested themselves.

## Viewing a single absorb proposal

Lines 137-153

Asked as: open an absorb proposal; view one merge proposal; absorb proposal not found; can't open proposal details; look up absorb request

A person with the absorb claim or the requester can view one absorb proposal. Anyone else, and any unknown id, gets "absorb proposal not found".

## Any signed-in user can propose an absorb

Lines 156-229

Asked as: request to merge clients; propose merging vendors; merge duplicate contacts; ask to absorb records; submit a merge request

Any signed-in user can propose an absorb of clients, client contacts or vendors. Turbo checks the target and the records with absorb.ValidateProposal and saves the proposal as pending, with the user as requester. It saves one open item for each record, and the whole save succeeds or fails together.

## Absorb proposal needs a reason

Lines 169-179

Asked as: merge request without reason; why was my absorb rejected for blank; absorb justification required; empty explanation for merge; say why records should merge

A proposal needs a reason. Turbo trims spaces from it. An empty reason gives "say why these records should be absorbed". Turbo checks the target and records first.

## Absorb proposal reason length limit

Lines 180-183

Asked as: merge explanation too long; absorb reason character limit; reason text too big; shorten merge justification; too many characters in proposal

A reason longer than absorb.ProposalTextMaxLength (the value is set in another file) gives "the reason is too long". Turbo counts the length after trimming spaces.

## Client contact proposals store the client

Lines 192-198

Asked as: client name on contact merge; which client does contact belong to; contact merge missing client; client contact proposal client id; contact proposal saved without client

When the proposal is for client contacts (client_contacts), Turbo stores the target contact's client id on the proposal. It also stores the client's name when it can find the client. The proposal is still saved if the client is not found.

## Only pending, idle proposals can be withdrawn or declined

Lines 231-252

Asked as: cancel my merge request; can't withdraw proposal; already decided merge; merge already carried out; stop pending absorb request

Withdrawing or declining works only on a pending proposal. A proposal that is not pending gives "this absorb proposal has already been decided". A proposal that an absorb action already carries out gives "this absorb proposal has been absorbed; commit or undo the absorb first". An unknown id gives "absorb proposal not found".

## Only the requester can withdraw a proposal

Lines 254-273

Asked as: someone else withdrew my request; withdraw another person's merge; only creator can cancel; who cancelled my absorb; take back merge request

Only the requester can withdraw a pending proposal. Anyone else, even a person with the absorb claim, gets "only the requester can withdraw this absorb proposal". A withdrawal sets the status to withdrawn and records who decided it.

## Deleting an absorb proposal

Lines 275-304

Asked as: remove old merge proposal; delete absorb request; can't delete pending proposal; clear finished merge requests; erase proposal history

Only the requester can delete an absorb proposal, or anyone else gets "only the requester can delete this absorb proposal". A pending proposal gives "withdraw this absorb proposal before you delete it". A withdrawn, declined, accepted or obsolete proposal can be deleted, and its record rows are deleted with it.

## Only an absorb holder can decline a proposal

Lines 306-316

Asked as: reject merge proposal; decline absorb request without permission; not allowed to turn down merge; refuse absorb proposal; can't decline merge

Only a person with the absorb claim can decline a proposal. Anyone else gets "you are not authorized to decline absorb proposals".

## Declining needs a comment

Lines 324-340

Asked as: decline without note; rejection comment required; decline reason empty; rejection note too long; why won't decline save

Declining needs a comment. Turbo trims spaces, and an empty comment gives "say why the proposal is declined". A comment longer than absorb.ProposalTextMaxLength gives "the comment is too long". Turbo checks the proposal state before the comment.

## Declining notifies the requester

Lines 339-348

Asked as: requester not told about rejection; notification after declined merge; decline email not sent; alert requester on decline; no message after rejection

After a decline, Turbo queues a notification of the outcome to the requester. If queueing fails, Turbo logs the error and the proposal still stays declined.
