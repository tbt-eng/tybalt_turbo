<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/invoice_lines.md","sha256":"7dac6a7512ae281e5f84a56f42f3aa852b1f1c40fa21ed30733d7d3cd3bcd927","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/invoice_lines.md

## Who can list or view invoice lines

Lines 7-17

Asked as: who can see invoice lines; view invoice lines access; list invoice lines permission; reporting staff can't see lines

Invoicing staff (the invoicing claim) and reporting staff (the report claim) can list and view invoice line records (invoice_lines). Nobody else can.

## Who can create an invoice line

Lines 19-23

Asked as: add invoice line; create invoice line; can't add a line to invoice; new invoice line permission

Only invoicing staff (the invoicing claim) can create an invoice line (invoice_lines). Reporting staff (the report claim) can view invoice lines but cannot create them.

## Nobody can update an invoice line through the API

Lines 25-27

Asked as: edit invoice line; change invoice line via API; update invoice line not working; fix a wrong invoice line

No update rule exists for invoice lines (invoice_lines), so nobody can update one through the API, not even invoicing staff. Only a PocketBase superuser or Turbo's own server code can change one.

## Who can delete an invoice line

Lines 29-33

Asked as: remove invoice line; delete invoice line; who can delete invoice lines; get rid of a line on invoice

Only invoicing staff (the invoicing claim) can delete an invoice line (invoice_lines).

## Invoice line kind is expense or PO

Lines 40-43

Asked as: invoice line type; invalid line kind; kind must be expense or po; line kind error

Each invoice line (invoice_lines) requires an invoice (invoice) and a kind (kind). The kind is either "expense" or "po". The expense (expense) and purchase order (purchase_order) fields are optional in the schema.

## Invoice line amount must be at least 0.01

Lines 44-44

Asked as: invoice line zero amount; negative amount on invoice; minimum line amount; amount rejected on line

An invoice line (invoice_lines) requires an amount (amount) of at least 0.01. Zero and negative amounts are not allowed.

## An expense can be on only one invoice line

Lines 54-54

Asked as: expense already invoiced; expense on two invoices; duplicate expense on invoice line; same expense added twice

A unique index (idx_invoice_lines_expense) allows an expense to appear on only one invoice line (invoice_lines). Turbo rejects a second line that uses the same expense. Lines with no expense are not limited.

## A purchase order can be on only one invoice line

Lines 55-55

Asked as: PO already on invoice; purchase order used twice; duplicate PO on invoice line; same purchase order twice

A unique index (idx_invoice_lines_purchase_order) allows a purchase order to appear on only one invoice line (invoice_lines). Turbo rejects a second line that uses the same purchase order. Lines with no purchase order are not limited.
