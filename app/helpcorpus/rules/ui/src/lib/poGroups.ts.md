<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poGroups.ts","sha256":"7bc90e4ce0c5ed372b9fb9f191c0aaf5e54a2387d082346d214503c81830b7f3","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poGroups.ts

## Who can view purchase order groups

Lines 8-11

A holder of the super_po claim can view purchase order groups. A holder of the admin claim can also view them. Staff with neither claim cannot (canViewPOGroups).

## Group preview rows carry per-field errors

Lines 35-45

Each row in a purchase order group preview can carry errors by field (errors). The preview counts these in error_count.

## Group preview lists active staff missing from the upload

Lines 49-55

The preview of a purchase order group returns a list of active staff (missing_active_staff) with their uid, email and name. It also gives the row count, the number of rows with errors (error_count), the total and the currency code.

## Second approver on a PO group needs a row over the threshold

Lines 56-60

The PO group preview lists as approvers (approvers) only the people who could approve every row. The preview marks a second approver as required (second_approval_required) only when some row is over the second approval threshold. The code in this file does not state the threshold's value.

## Group description is blank when each PO's description comes from the upload

Lines 68-70

A purchase order group's summary description holds the template's description. It is blank when each purchase order's description came from the CSV upload.

## Purchase order group title falls back to vendor name

Lines 78-83

A purchase order group is titled by its template description when it has one. If the description is blank, the title is the vendor name followed by "purchase orders". If the vendor name is also blank, the title is "PO Group" (poGroupTitle).

## Purchase order group with no description or vendor is titled PO Group

Lines 80-83

A purchase order group is titled "PO Group" when it has no template description and no vendor name. This is the last fallback after the description and the vendor name (poGroupTitle).
