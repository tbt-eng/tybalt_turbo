<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/vendors.sql","sha256":"b70acc120170020d395e975b3689146e6215915eaff085765da13b5ad5c3e88b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/vendors.sql

## Vendor referenced flag counts any expense or purchase order

Lines 1-29

A vendor is marked referenced when any expense or any purchase order uses it, whatever its status or commit state. The flag does not depend on the counts, so a vendor with only uncommitted expenses or non-Active purchase orders is still referenced while both counts show 0. This is the check that the vendor delete rule uses, so the UI can tell when a delete will be refused.

## Vendor expense count includes only committed expenses

Lines 10-21

The vendor list counts a vendor's expenses (expenses_count) only when the expense is committed (committed is not empty). Uncommitted expenses add nothing to the count. A vendor with no expenses shows 0.

## Vendor purchase order count includes only Active purchase orders

Lines 11-29

The vendor list counts a vendor's purchase orders (purchase_orders_count) only when the status is Active. Purchase orders in any other status add nothing to the count. A vendor with no purchase orders shows 0.

## Records with no vendor are ignored in vendor counts

Lines 19-27

Expenses and purchase orders with an empty vendor are left out of the vendor counts and the referenced flag.

## Vendor list can be filtered to one vendor

Lines 30-31

The vendor list returns one vendor when an id is given. When the id is missing or empty, it returns all vendors. Results are sorted by vendor name.
