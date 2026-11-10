<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_details.sql","sha256":"3f5d0df5916c82a9af4db6b4d9724c6af703658b1d5b6b527725df45cb0dd462","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_details.sql

## Missing attachment reason, card digits and distance in expense details

Lines 9-32

Asked as: blank attachment reason; missing card digits; expense shows empty card number; distance shows wrong format

The reason for a missing attachment (attachment_missing_reason) and the last 4 card digits (cc_last_4_digits) show blank when not set. The total and distance are returned as numbers.

## Expense attachment details come from the linked expense document

Lines 11-16

Asked as: attachment not showing on expense; receipt file missing; no file on expense; attachment hash blank

The attachment file and its hash come from the expense's attachment document (ed) and show blank when there is none. The attachment collection id is 'pbc_2089657321' and the record id is the document's id only when a document is linked. Otherwise both are blank.

## Expense details show CAD when the expense has no currency

Lines 33-44

Asked as: expense currency blank; no currency on expense; expense shows CAD; currency symbol missing

When an expense has no currency or the currency is not found, the expense details show the currency code 'CAD', the symbol 'CAD' and a rate of 1. The currency icon and rate date are blank.

## Expense details show the settled total as 0 when unsettled

Lines 34-36

Asked as: settled amount shows zero; expense not settled total; who settled expense; settled date empty

When an expense has no settled total (settled_total), the details show 0. When it has no settler (settler) or settled date (settled), they show blank text.

## Expense details include linked purchase order information

Lines 37-85

Asked as: purchase order details missing; PO info not on expense; expense linked PO blank; PO vendor not showing

When an expense links to a purchase order, the details include the PO number, vendor, job, division, category, description, payment type, branch, kind, total, currency and owner. All of these show blank when there is no purchase order, except the total, which shows 0.

## Expense details show people by full name

Lines 54-58

Asked as: name missing on expense; owner shows blank; approver name not showing; expense shows no person name

The expense's owner (uid), approver, rejector, settler and creator, and the purchase order's owner, appear as given name plus surname. A person who is not set or not found shows as blank. This blank applies even when the name parts are missing.

## Purchase order currency falls back to CAD in expense details

Lines 78-83

Asked as: PO currency blank; purchase order no currency; PO shows CAD rate; PO currency symbol missing

When the linked purchase order has no currency, the expense details show the purchase order's currency code as 'CAD', the symbol as 'CAD' and the rate as 1. The icon and rate date are blank.

## Expense details look up one expense by its id

Lines 86-111

Asked as: expense details not loading; expense missing from details; job name blank on expense; expense hidden when link missing

The details return a single expense that matches the requested id. Related records (job, client, vendor, division, category, kind, branch) are joined optionally, so missing links show blank names instead of hiding the expense.
