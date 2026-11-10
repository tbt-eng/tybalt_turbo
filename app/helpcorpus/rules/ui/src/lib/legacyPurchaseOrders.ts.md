<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/legacyPurchaseOrders.ts","sha256":"44a1c6125c9388640f1d03f89185bf6ecfdec4aec67491b0865e95beb07d8e30","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/legacyPurchaseOrders.ts

## Legacy purchase order save sends only a fixed set of fields

Lines 6-24

When Turbo creates or updates a legacy purchase order, it sends only these 13 fields: uid, approver, po_number, date, division, description, payment_type, total, vendor, type, kind, branch and job. Any other field on the record is not sent in the create or the update request (buildLegacyPurchaseOrderPayload).
