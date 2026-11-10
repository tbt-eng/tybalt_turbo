<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/expenseApprovalQueue.ts","sha256":"8e53ee63676e7a920eca490702c119e135b38377df71e441c8821b41d71ec7c6","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/expenseApprovalQueue.ts

## Approval queue skips rejected expenses

Lines 9-18

The list of expenses left to review leaves out rejected expenses (rejected is not an empty string). A rejected expense stays in the pending list until its owner recalls it, but it cannot be approved.

## Approval queue reads at most 200 pending expenses

Lines 9-17

Turbo fetches the viewer's pending expenses with a limit of 200, the endpoint's page cap. Only those 200 are checked for the next expense to review.

## After a decision, Turbo opens the next pending expense

Lines 20-27

After an expense is approved or rejected, Turbo opens the details of the first expense left to review that is not the one just decided. Turbo adds a history entry, so Back goes through the expenses just decided.

## When no expenses are left, go to Approved By Me

Lines 27-36

When no other expense is left to review, Turbo opens the approved expenses list (/expenses/approved) and adds " · nothing left to review" to the flash message.

## Flash message after approving or rejecting an expense

Lines 29-34

The flash reads "Approved" or "Rejected" followed by the expense title. An approval gets the success tone, and a rejection gets the info tone. The flash has a "View" link to the details of the expense just decided.
