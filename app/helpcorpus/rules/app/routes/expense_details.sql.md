<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_details.sql","sha256":"2d79fb3ee10bd80f273b8d2203b02ac541e28f18575b812c626e82592b86f8b5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_details.sql

## Expense without an attachment document shows empty attachment fields

Lines 11-21

Asked as: expense has no receipt; attachment fields blank; missing receipt document; vendor invoice number empty

When an expense has no linked document (ed), the attachment, attachment hash, attachment collection ID and attachment record ID are empty strings. The vendor invoice number and date are empty, and the vendor invoice total is 0. When a document exists, the collection ID is 'pbc_2089657321'.

## Missing-attachment reason and card digits default to empty

Lines 14-37

Asked as: missing attachment reason blank; no card last 4 digits; card digits empty; why is attachment reason empty

The expense details show an empty string when the expense has no attachment missing reason (attachment_missing_reason) or no card last 4 digits (cc_last_4_digits).

## Vendor invoice details come from the expense's document

Lines 17-21

Asked as: same invoice on several expenses; vendor invoice details wrong; invoice number shared between expenses; invoice total same for expenses

The vendor invoice number, date and total come from the expense's attachment document (expense_documents), not from the expense. Expenses that share one document therefore show the same vendor invoice details.

## Unsettled expense shows settled total 0 and empty settler

Lines 39-41

Asked as: settled total showing zero; who settled expense; expense not paid yet; settler name blank

When an expense has not been settled, the details show settled total (settled_total) 0 and an empty settler (settler) and settled date (settled). The settler name is empty too.

## Expense currency shows CAD when no currency is set

Lines 45-49

Asked as: currency defaults to CAD; no currency on expense; currency symbol CAD; exchange rate shows 1

When an expense has no currency (currency), the expense details show currency code (currency_code) 'CAD', symbol (currency_symbol) 'CAD' and rate (currency_rate) 1. The currency icon and rate date are empty.

## Expense details show names as empty when a link is missing

Lines 50-64

Asked as: client name missing; category blank on expense; approver name empty; vendor name not showing

The expense details show an empty string for the client, category, kind, job number and description, division name and code, vendor name and alias, and branch when the expense lacks that link. The person names (owner, approver, rejector, settler, creator, PO owner) are also empty when no profile matches.

## Expense without a purchase order shows empty PO fields

Lines 65-90

Asked as: expense without PO fields blank; no purchase order details; PO total zero; PO number empty

When an expense has no purchase order, the PO fields in the expense details are empty strings: PO number, vendor, job, division, category, description, payment type, branch, kind and owner. The PO total (po_total) is 0.

## Purchase order currency shows CAD when none is set

Lines 84-88

Asked as: PO currency is CAD; purchase order currency missing; PO exchange rate 1; PO currency icon empty

When the linked purchase order has no currency, the expense details show the PO currency code (po_currency_code) 'CAD', symbol 'CAD' and rate 1. The PO currency icon and rate date are empty. An expense with no purchase order also gets these defaults.

## Expense details are fetched by a single expense id

Lines 91-116

Asked as: expense details by id; show one expense only; approval and settlement fields; expense lookup by id

The expense details return only the expense whose id equals the requested id. The details include the approval, rejection, commit and settlement fields of that expense.
