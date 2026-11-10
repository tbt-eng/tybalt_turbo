<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/vendors_api.go","sha256":"99681e236b154dfe8a2cb7389673b3d8c7b37f2eaa3e2e1bb83e37bedfdfacd0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/vendors_api.go

## Vendor list shows usage counts and whether a vendor is referenced

Lines 24-33

Asked as: can't delete vendor; vendor delete blocked; vendor in use; vendor usage count; vendor referenced by expenses

Each vendor in the response carries its expenses count (expenses_count), its purchase orders count (purchase_orders_count) and a referenced flag. The referenced flag is true when any expense or purchase order uses the vendor. The vendors deleteRule then refuses to delete a referenced vendor. This file does not compute the counts or the flag. The vendors.sql query does.

## Vendor lookup by ID returns an error if the vendor is missing

Lines 48-51

Asked as: vendor not found; missing vendor by id; vendor id doesn't exist; vendor page error; can't open vendor

When a vendor is requested by its ID and the vendors query returns no row, Turbo refuses with the message "vendor not found".

## Vendors request without an ID returns all vendors

Lines 48-60

Asked as: list all vendors; vendors without id; show every vendor; get all suppliers; vendor list empty

When no vendor ID is given, Turbo returns every vendor row from the vendors query as a list. A single vendor is returned only when an ID is given. This file does not state the list's order.
