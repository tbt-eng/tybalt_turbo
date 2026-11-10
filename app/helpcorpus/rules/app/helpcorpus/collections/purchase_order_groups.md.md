<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/purchase_order_groups.md","sha256":"0b25f32ca1a77a40a79606073bf345c82a96336675f5f23b912125f4a4d6c811","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/purchase_order_groups.md

## Nobody can list purchase order group records through the API

Lines 7-9

Asked as: can't see purchase order groups; list PO groups api; who can read PO groups; purchase order group list empty

The purchase_order_groups collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view a purchase order group record through the API

Lines 11-13

Asked as: open one PO group record; view single PO group api; can't open purchase order group; fetch PO group by id

The purchase_order_groups collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create a purchase order group through the API

Lines 15-17

Asked as: add new PO group; create purchase order group api; can't make PO group; new PO group not saving

The purchase_order_groups collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update a purchase order group through the API

Lines 19-21

Asked as: edit PO group; change purchase order group api; update PO group forbidden; PO group changes not saved

The purchase_order_groups collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete a purchase order group through the API

Lines 23-25

Asked as: remove PO group; delete purchase order group api; can't delete PO group; get rid of PO group

The purchase_order_groups collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Purchase order group requires a creator and a 15-character id

Lines 27-33

Asked as: PO group creator field; invalid PO group id; purchase order group id format; PO group description optional

Every purchase order group needs a creator (creator), which is one person (users) and is required. The id must be exactly 15 characters, using only lowercase letters and digits. The description is optional.
