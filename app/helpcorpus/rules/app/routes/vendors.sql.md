<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/vendors.sql","sha256":"b70acc120170020d395e975b3689146e6215915eaff085765da13b5ad5c3e88b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/vendors.sql

## Vendor referenced flag counts any expense or purchase order

Lines 1-29

Asked as: vendor marked referenced; vendor shows referenced but counts zero; can't delete vendor; vendor has uncommitted expenses; vendor flag referenced

A vendor is marked referenced when any expense or any purchase order uses it, whatever its status or commit state. The flag does not depend on the counts, so a vendor with only uncommitted expenses or non-Active purchase orders is still referenced while both counts show 0. This is the check that the vendor delete rule uses, so the UI can tell when a delete will be refused.

## Vendor expense count includes only committed expenses

Lines 10-21

Asked as: vendor expense count wrong; uncommitted expenses in vendor count; vendor shows 0 expenses; expenses_count not right; vendor expenses not counted

The vendor list counts a vendor's expenses (expenses_count) only when the expense is committed (committed is not empty). Uncommitted expenses add nothing to the count. A vendor with no expenses shows 0.

## Vendor purchase order count includes only Active purchase orders

Lines 11-29

Asked as: vendor purchase order count; vendor PO count zero; draft POs not counted vendor; purchase_orders_count wrong; vendor orders not showing

The vendor list counts a vendor's purchase orders (purchase_orders_count) only when the status is Active. Purchase orders in any other status add nothing to the count. A vendor with no purchase orders shows 0.

## Records with no vendor are ignored in vendor counts

Lines 19-27

Asked as: expense with no vendor; PO without vendor counted; blank vendor in counts; missing vendor on expense; empty vendor totals

Expenses and purchase orders with an empty vendor are left out of the vendor counts and the referenced flag.

## Vendor list can be filtered to one vendor

Lines 30-31

Asked as: filter vendor list; show one vendor; vendor list by id; list vendors sorted by name; find a vendor

The vendor list returns one vendor when an id is given. When the id is missing or empty, it returns all vendors. Results are sorted by vendor name.
