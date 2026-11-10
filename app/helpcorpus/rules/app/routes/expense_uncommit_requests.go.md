<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_uncommit_requests.go","sha256":"6661a8f0752b4091c06c9c46882f0a79828e7e0276213800082e0cd1588d8b95","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_uncommit_requests.go

## Uncommit request limits and outcome names

Lines 22-37

Asked as: uncommit request limit; how many uncommit requests shown; decline comment minimum length; uncommit note character limit

A decided-request list shows at most 200 requests (expenseUncommitDecidedListLimit). A request note holds at most 1000 characters (expenseUncommitNoteMaxLength). A decline comment needs at least 4 characters (expenseUncommitMinCommentLength). The two outcomes are uncommit_and_reject and uncommit_only.

## Request statuses for an uncommit request

Lines 29-32

Asked as: uncommit request status; pending or approved request; can a request still be changed; request status meaning

An uncommit request is pending, approved, declined or cancelled. Only a pending request can still be decided or cancelled.

## Uncommit request list rows carry expense details and a blocked flag

Lines 47-71

Asked as: uncommit request list details; blocked request row; why is request blocked; uncommit queue columns

Each row in the uncommit request list shows the request's status, outcome, note, requester, decision details and the expense snapshot (employee, PO number, vendor, description, currency, payment type, total, committed date, pay period ending). Pending rows whose eligibility has closed also carry a 'blocked' field (Blocked).

## Uncommit request failures return a code and message

Lines 81-103

Asked as: uncommit request error; generic failed uncommit error; uncommit failed with no reason; uncommit update failed message

When an uncommit request action fails with a known refusal, Turbo answers with the refusal's code and message. When the failure is not a known refusal, Turbo answers with the generic error 'failed to update uncommit request'.

## Caller needs any one of the listed claims

Lines 105-117

Asked as: which claim lets me uncommit; any one permission uncommit; multiple claims uncommit access; uncommit permission needed

Where an uncommit action lists several claims, holding any one of them is enough.

## Permission refusal message for uncommit actions

Lines 119-121

Asked as: not allowed to uncommit; no permission uncommit; uncommit access denied; why can't I uncommit

When Turbo refuses an uncommit action for lack of permission, it says "you are not authorized to" followed by the action, such as "request an uncommit" or "decline this uncommit request".

## One pending uncommit request per expense

Lines 123-135

Asked as: expense already has request; find pending uncommit for expense; no pending uncommit request; existing uncommit request on expense

Turbo looks up the one pending request (status pending) for an expense. An expense with no pending request has none to show or apply.

## Uncommit eligibility is checked from the expense itself

Lines 137-147

Asked as: can this expense be uncommitted; uncommit eligibility check; why uncommit blocked by pay period; uncommit allowed after pay period

Whether an expense can still be uncommitted depends on its payment type (payment_type), pay period ending (pay_period_ending) and commit time (committed), compared with the current time. The check is FindExpenseUncommitBlock in another file, so this file does not state its result.

## Request keeps a snapshot of the expense details

Lines 149-182

Asked as: request copies expense details; expense details after delete; snapshot of expense in request; keep details after uncommit

A new or approved request copies the employee name, PO number, vendor name, description, currency code, payment type, total, committed time, committer, committed week ending and pay period ending from the expense. Missing names become blank. The copy lets the queue and history still show the details after uncommit clears the commit fields or the expense is deleted. The currency code comes from EffectiveCurrencyCode.

## Expense details offer an uncommit request only to committers and admins

Lines 194-213

Asked as: who sees uncommit button; uncommit option on expense page; expense details uncommit state; no uncommit button showing

The expense details page returns uncommit request state only to a committer (commit claim) or an admin (admin claim), and only when the expense is committed. Otherwise it returns nothing.

## When a committer can request an uncommit

Lines 215-232

Asked as: can I ask to uncommit; request uncommit button missing; admin cannot request uncommit; uncommit button greyed out

The details page shows the expense's pending request, if any. It can offer a request (can_request) only to a committer, not to an admin without the commit claim. The expense must have no pending request, no block and not be imported. Imported expenses (_imported) are never checked for blocks.

## Who can see uncommit reasons

Lines 235-254

Asked as: view uncommit reasons; list of uncommit reasons; who can see reasons; uncommit reasons not showing

Only a committer (the commit claim) or an admin (the admin claim) can list uncommit reasons. Anyone else gets "you are not authorized to view uncommit reasons". Only active reasons are listed, ordered by sort and then name.

## Who can list uncommit requests

Lines 258-266

Asked as: see uncommit requests; uncommit request queue access; who views uncommit requests; can't open uncommit list

Only a committer (the commit claim) or an admin (the admin claim) can list expense uncommit requests. Anyone else gets "you are not authorized to view uncommit requests".

## Uncommit request list filters and limits

Lines 268-282

Asked as: uncommit list filter; show decided uncommit requests; oldest pending requests first; status filter uncommit

The list shows pending requests, oldest first, when status is empty or "pending". With status "decided" it shows every non-pending request, newest decision first, capped at 200 (expenseUncommitDecidedListLimit). Any other status gets "status must be pending or decided".

## Pending requests show whether they are now blocked

Lines 287-297

Asked as: request now blocked; pending request blocked message; uncommit request no longer allowed; why pending request blocked

In the request list, each pending request is re-checked against the uncommit eligibility rules (utilities.FindExpenseUncommitBlock) using the payment type, pay period ending and committed date. If a rule now blocks it, the row carries a "blocked" explanation. Decided requests are not re-checked.

## Only committers can request an uncommit

Lines 304-314

Asked as: how to request uncommit; admin cannot request uncommit; submit uncommit request; request uncommit not permitted

Creating an uncommit request needs the commit claim, and an admin alone is refused. Otherwise Turbo says "you are not authorized to request an uncommit". Turbo also refuses it when expense editing is not available (requireExpensesEditing).

## Uncommit request preconditions on the expense

Lines 328-349

Asked as: cannot request uncommit; expense not committed error; imported expense uncommit refused; why uncommit request refused

The expense must exist ("expense not found") and be committed ("this expense is not committed"). Imported expenses are refused with "imported expenses can't be requested". An expense with a pending request is refused with "this expense already has a pending uncommit request". An eligibility block on the expense (utilities.FindExpenseUncommitBlock) is refused with that block's own message.

## Uncommit request needs a valid reason, outcome and note

Lines 351-363

Asked as: choose reason uncommit; uncommit note required; note too long uncommit; pick outcome for uncommit

The reason must exist and be active, or Turbo says "choose a reason". The outcome must be uncommit_and_reject or uncommit_only, or Turbo says "choose whether the expense should also be rejected". A blank note (after trimming) is refused with "this reason needs a note" when the reason requires one. A note over 1000 characters is refused with "the note is too long".

## New uncommit request is saved pending with a snapshot

Lines 365-383

Asked as: save uncommit request; create uncommit request; new uncommit saved pending; submit uncommit with note

A new request is saved as pending with the requester, reason, trimmed note and outcome. It copies the employee name, PO number, vendor name, description, currency, payment type, total, committed fields and pay period ending from the expense. The currency code comes from EffectiveCurrencyCode. The snapshot is kept because uncommitting clears the commit fields and the expense can later be deleted.

## Only the requester can cancel a pending request

Lines 395-418

Asked as: cancel my uncommit request; withdraw uncommit request; cancel request not allowed; uncommit request already decided

Cancelling an uncommit request is refused unless the caller made it: "you are not authorized to cancel this uncommit request". A request that is not pending is refused with "this uncommit request has already been decided". A missing request gives "uncommit request not found". A successful cancel sets status cancelled and records the canceller and time.

## Only admins can decline an uncommit request

Lines 428-434

Asked as: decline uncommit request; reject uncommit request admin; who can decline uncommit; deny uncommit request

Declining needs the admin claim. Otherwise Turbo says "you are not authorized to decline this uncommit request".

## Decline needs a comment of 4 to 1000 characters

Lines 442-459

Asked as: decline comment too short; why declined comment required; decline reason length; comment too long decline

The decline comment is trimmed. Under 4 characters Turbo says "say why the request is being declined". Over 1000 it says "the comment is too long". Only a pending request can be declined, or Turbo says "this uncommit request has already been decided".

## Declining records the decision and notifies

Lines 460-471

Asked as: declined uncommit notification; after declining request; decision sent to requester; uncommit declined message

A decline sets status declined, with the decider, the time and the comment. It then queues a decision notification, and responds "uncommit request declined".

## Only admins can approve uncommit requests, in bulk

Lines 488-524

Asked as: approve many uncommit requests; bulk approve uncommit; approve selected uncommits; approve uncommit ids

Approving needs the admin claim, or Turbo says "you are not authorized to approve uncommit requests". The call needs a non-empty list of ids, or Turbo says "ids is required". Each request is approved in its own transaction. One failure does not stop the others, and the response lists ok, code and message for each.

## Approval re-checks the request and expense

Lines 534-554

Asked as: approve uncommit fails; expense gone when approving; approval blocked by expense; approve request already decided

Approval needs a pending request, or Turbo says "this uncommit request has already been decided". The expense must still exist ("the expense no longer exists") and be committed. Imported expenses are refused with "imported expenses can't be requested". An eligibility block on the expense is refused with the block's message.

## Approving uncommits the expense and may reject it

Lines 573-605

Asked as: approve uncommit and reject; approved uncommit uncommits expense; approval rejects expense too; uncommit_and_reject outcome effect

Approval refreshes the snapshot, uncommits the expense, and marks the request approved with the decider and time. For the uncommit_and_reject outcome, it also rejects the expense, acting as the original requester because the approving admin may lack reject permission. Both queue approval and direct uncommit use this same path.

## Approval notifications

Lines 609-617

Asked as: who is notified approval; requester told of approval; rejection notification after uncommit; notify on uncommit approval

Approval notifies the requester of the decision. If the expense was also rejected, it queues the expense-rejected notifications with the requester and the rejection reason.

## Rejection reason is built from reason name and note

Lines 627-640

Asked as: rejection reason text; reason shown on rejected expense; uncommit reason name and note; reason deleted rejection message

The rejection reason is "reason name: note" when both exist. It is the name alone or the note alone when only one exists. It is "Uncommit requested" when the reason was deleted and there is no note.

## Direct uncommit honours a pending request

Lines 652-666

Asked as: uncommit without request; direct uncommit outcome; uncommit expense directly; pending request used for uncommit

A direct uncommit applies the expense's pending request, so its outcome is honoured. With no pending request, it uses a new uncommit_only request by the caller.
