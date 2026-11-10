<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/close_purchase_order.go","sha256":"8c8e3c2effd4fbf3947d0acec58624fe8283aee4b8b33cf31cb977cc2850b742","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/close_purchase_order.go

## Closing runs as one all-or-nothing transaction

Lines 22-118

Asked as: close fails halfway; partial close of PO; PO changed despite error; rollback when closing; all or nothing close

The permission check, the checks on the purchase order and its expenses, and the save all run in one transaction. If any step fails, the purchase order is not changed.

## Only payables staff can close a purchase order

Lines 23-38

Asked as: who can close PO; can't close purchase order; not authorized to close; no permission to close PO; payables only close

Turbo closes a purchase order only for a user with the payables_admin claim (payables staff). Any other user gets the message "you are not authorized to close purchase orders".

## Closing checks run in a fixed order

Lines 23-101

Asked as: which close error shows first; several reasons close refused; order of closing checks; why wrong error on close

Turbo checks in this order: payables claim, One-Time type, at least one committed expense, Active status, then pending expenses. It shows only the first refusal that applies. A PO that is not Active but has no committed expense gets the "no expenses" message, not the "not active" message.

## Closing a missing purchase order is refused

Lines 40-48

Asked as: purchase order not found; PO does not exist; error fetching purchase order; close missing PO

If the purchase order does not exist, Turbo refuses the close and shows "error fetching purchase order:" followed by the error detail.

## One-Time purchase orders cannot be closed manually

Lines 50-58

Asked as: close one-time PO; can't close one time order; manually close one-time; cancel one-time purchase order

Turbo refuses to close a purchase order whose type is One-Time. It shows "One-Time purchase orders may be cancelled but not manually closed".

## Closing needs at least one committed expense

Lines 60-77

Asked as: close PO with no expenses; no expenses PO close; PO has no committed expense; cancel instead of closing

Turbo refuses to close a purchase order that has no committed expense. It shows "only cumulative or recurring purchase orders with at least one associated expense may be closed manually. Cancel the purchase order instead."

## Only active purchase orders can be closed

Lines 79-86

Asked as: close PO not active; close draft PO; PO status blocks closing; only active purchase orders

Turbo refuses to close a purchase order whose status is not Active. It shows "only active purchase orders can be closed". This check runs after the One-Time and committed-expense checks, so those messages come first.

## Pending expenses block closing a purchase order

Lines 88-101

Asked as: expense blocking PO close; PO won't close; why can't I close PO; unresolved expenses stop closing; pending expense close error

Turbo refuses to close a purchase order while it has an expense that is submitted, not committed and not rejected. Drafts, recalled expenses and rejected expenses do not block closing. Approved expenses that are not yet committed do block it. Turbo shows a message that describes the blocking expenses.

## Closing a purchase order records closer, time and Closed status

Lines 103-106

Asked as: who closed the PO; closed date on PO; PO closed status set; record closer and time

When payables staff close a purchase order and every check passes, Turbo sets its status to Closed (status). It records the current time as the closed time (closed). It records the person who closed it as the closer (closer).

## Closing fails if the purchase order cannot be saved

Lines 109-115

Asked as: error saving purchase order; PO close won't save; save failed when closing; closing not saved

If Turbo cannot save the closed purchase order, the closing fails with the message "error saving purchase order: ..." (code error_saving_purchase_order). The purchase order is not closed.

## A refused closing returns a message and a code

Lines 120-132

Asked as: close error code; closing refused message; error text closing PO; PO unchanged after refusal

When a check fails, Turbo returns the check's message and its error code, and the purchase order stays unchanged. An unexpected error returns only the error text.

## Successful closing returns the updated purchase order

Lines 134-139

Asked as: PO after closing; response after close; updated purchase order returned; closed PO details

After a purchase order closes, Turbo reloads it from the database and returns it. The returned record shows status Closed, along with the closer and closed time.

## Drafts, recalled and rejected expenses do not block closing

Lines 143-146

Asked as: draft expense blocks close; recalled expense PO close; rejected expense stop closing; does draft stop PO close

Only submitted expenses that are not committed and not rejected block manual closing. Draft and recalled expenses do not block it. Approved expenses that are not yet committed do block it.

## Which expenses count as pending when closing a purchase order

Lines 155-164

Asked as: which expenses block PO; submitted expense stops closing; expense counted as pending; pending expenses on PO

An expense on the PO blocks closing if it is submitted and neither committed nor rejected. Drafts, recalled expenses, rejected expenses and committed expenses do not block.

## Pending expense message lists approvers by name

Lines 171-189

Asked as: who must approve expense; approver names close PO; awaiting approval names; list of approvers blocking

If submitted expenses on the PO are not yet approved, the closing refusal reads "This PO has one or more expenses awaiting approval by" followed by the approvers' names, separated by commas. Each name appears once. Names are sorted alphabetically.

## Approved but uncommitted expenses block closing a purchase order

Lines 173-192

Asked as: approved expense not committed; awaiting commitment PO; approved but not committed; commit expense before closing

If a submitted expense on the PO is approved but not yet committed, closing is refused. The message reads "This PO has one or more approved expenses awaiting commitment."

## Missing approver name still blocks closing

Lines 176-185

Asked as: approver name missing; approver has no name; unavailable approver close; nameless approver blocks PO

If an unapproved pending expense has an approver with no profile name, the block stays. The message names "an approver whose name is unavailable" in the list of approvers.

## Both pending reasons combine in one closing message

Lines 186-193

Asked as: two messages closing PO; approval and commitment message; both pending reasons shown; combined close refusal message

When a PO has expenses awaiting approval and approved expenses awaiting commitment, Turbo shows both sentences in one message. The approval sentence comes first. A single space separates the two sentences.
