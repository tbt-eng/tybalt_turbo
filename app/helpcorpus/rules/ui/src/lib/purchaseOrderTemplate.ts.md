<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/purchaseOrderTemplate.ts","sha256":"4ba8fc274bc0854dc044590ed9d272b0355ad5ab0a82c5206efb06c5837e0520","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/purchaseOrderTemplate.ts

## New blank purchase order starts with these defaults

Lines 12-33

Asked as: new blank purchase order; default values for new PO; what does a new PO start with; blank PO fields; unapproved status by default

A blank purchase order starts with status Unapproved (Unapproved) and type One-Time. Its date is today's local date and its frequency is Weekly. Its payment type is On Account (OnAccount) and its total is 0. Covered within project budget (covered_within_project_budget) is false. The PO number, requester, end date, division, description, currency, vendor, job, category and kind start empty.

## Copying a PO template keeps these fields

Lines 35-57

Asked as: copy PO from template; duplicate a purchase order; template fields carried over; clone existing PO; reuse PO template

A new purchase order made from a template copies the type, date, end date, frequency, division, description, total, currency, covered within project budget, payment type, vendor, job, category and kind from the template. It also copies the branch. When the template has no value for a field, Turbo uses the blank default instead. The new PO starts with the blank status Unapproved and the blank PO number.

## Copied template keeps its date, not today's date

Lines 43-43

Asked as: template date wrong; PO date from template; date not today; copied PO has old date; template keeps its date

A purchase order made from a template takes the template's date (date). Turbo uses today's date only when the template has no date.

## Template copy clears approvers, parent PO and attachment

Lines 57-61

Asked as: approver missing after copy; template lost approvers; copied PO no attachment; parent PO not copied; choose approvers again

A new purchase order made from a template does not copy the parent PO (parent_po), the approver (approver), the second approver (priority_second_approver) or the attachment. Turbo sets all four to empty, so the requester must choose approvers and attach files again. The branch is empty when the template has no branch.
