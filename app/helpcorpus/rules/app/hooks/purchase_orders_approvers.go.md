<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/purchase_orders_approvers.go","sha256":"edc02e854ce90065d28a16f0cbdf22289c56de2e883667dc908b95a7e486996b","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/purchase_orders_approvers.go

## Approver divisions must be a JSON array of strings

Lines 22-63

When a purchase order approver's settings (po_approver_props) are saved, the divisions list must be an array of strings. Otherwise Turbo refuses with "division validation error" and says "divisions must be an array of strings". A value that is not an array at all is refused with "divisions must be a JSON array".

## Empty or missing divisions accepted for approver settings

Lines 25-35

A purchase order approver's settings (po_approver_props) can be saved with no divisions. Turbo treats a missing divisions value as an empty list. An empty JSON value is also treated as an empty list.

## Each approver division must be active

Lines 65-69

Turbo checks every division in a purchase order approver's settings (po_approver_props) with EnsureActiveDivision. If a check fails, Turbo returns that error and does not save the settings. The rules for an active division are in another file.

## Approver settings require a user

Lines 71-80

A purchase order approver's settings (po_approver_props) must name a user (user_claim). A blank value or one with only spaces is refused with "po_approver_props validation error" and "user_claim is required".

## One approver settings record per user

Lines 82-110

Each user (user_claim) can have only one purchase order approver settings record (po_approver_props). Saving a second record for the same user is refused with "po_approver_props validation error" and "po_approver_props already exists for this user_claim". Editing the user's own existing record is allowed.
