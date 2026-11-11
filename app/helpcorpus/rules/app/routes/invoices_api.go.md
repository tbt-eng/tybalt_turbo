<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/invoices_api.go","sha256":"239eb06fc69f6c96daf23cf3cb25dd642b0d96914cc2c9962a27c82d1d2819be","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/invoices_api.go

## Drafts are edited elsewhere; these routes only issue and list

Lines 16-19

Asked as: where to edit draft; can't change draft invoice; draft invoice edits; fix invoice before sending

Turbo creates, edits and deletes draft invoices through the invoices collection. It handles their lines through the invoice_lines collection. The routes in this file issue an invoice and return the joined cost details that the screens need.

## Expense invoice line falls back to CAD when no currency

Lines 52-75

Asked as: currency blank on invoice; no currency on invoice row; why CAD on invoice; currency missing invoice cost

An invoice cost row for an expense or PO shows the currency code (currency_code) from the cost's currency. It shows 'CAD' when the cost has no currency.

## Blank job, PO number, description and vendor on cost rows

Lines 52-64

Asked as: empty fields on invoice; missing vendor name invoice; expense shows PO number; blank job number invoice

A cost row on an invoice shows blank text, not nothing, when the job number, description or vendor name is missing. An expense row shows its linked PO number when it has one, and blank when it has none.

## Invoice lines show expense and PO details, sorted

Lines 54-100

Asked as: what is on invoice lines; invoice line order; who created invoice line; invoice line details list

Each invoice line shows the expense or PO it bills, with job number, PO number, description, vendor, date, cost total, currency, line amount and creator name. Lines sort by job number, then date, then creation time. A missing currency shows as CAD.

## Billable list needs invoicing staff; lists client jobs or one job's costs

Lines 77-268

Asked as: unbilled costs for client; what can I bill; list jobs to invoice; costs waiting to invoice

Listing billable costs for an invoice requires an active invoicing user. Without a job parameter, Turbo lists the invoice client's jobs that have at least one unbilled cost, with the count, sorted by job number descending. With a job parameter, it lists that job's unbilled costs, newest date first. The job must belong to the invoice's client.

## Which costs count as billable

Lines 77-83

Asked as: why cost not billable; expense won't show for billing; PO not billable status; cost missing from billing list

An expense is billable when it is committed and has no invoice line. A PO is billable when its status is Active or Closed and it has no invoice line. Costs already on an invoice line never appear as billable.

## Invoice lines can be viewed by invoicing or report holders

Lines 141-197

Asked as: can't see invoice lines; invoice lines access denied; no access to invoice lines; invoice lines won't load

Listing an invoice's lines requires the invoicing or report claim. Otherwise Turbo says "the invoicing or report claim is required". If the invoice does not exist, Turbo says "invoice not found".

## Only invoicing staff can issue a draft invoice

Lines 157-167

Asked as: can't issue invoice; issue button not working; who can send invoice; invoice won't go out

Issuing a draft invoice requires an active invoicing user (hooks.RequireActiveInvoicingUser). The code does not state what that check requires. If the invoice does not exist, Turbo says "invoice not found".

## Issuing an invoice sets status, issuer and time

Lines 168-183

Asked as: who issued invoice; when invoice was sent; invoice marked issued; issued date on invoice

When a draft passes the issue validation (hooks.ValidateInvoiceForIssue), Turbo sets its status to Issued. It records the issuing user (issued_by) and the current UTC time (issued_at). The validation runs in the same transaction as the save, and the code in this file does not state its checks.

## Each invoice line says whether the viewer can open its cost

Lines 209-239

Asked as: can't open expense from invoice; link to expense greyed out; can't open cost on line; invoice line link disabled

Each line has can_open, true only when the viewer may open that expense or PO. A PO line uses purchase order visibility. An expense line uses expense visibility. An admin (the admin claim) can also open any committed expense.

## Cost invoice view reports whether the viewer can open the invoice

Lines 269-281

Asked as: no link to invoice on expense; can't open invoice from expense; invoice link missing on PO; view invoice button hidden

On an expense or PO details page, the viewer can open the invoice (can_open_invoice) only if they hold the invoicing or report claim.

## Cost invoice view fields: total, currency and client

Lines 269-289

Asked as: expense total on invoice view; client on cost details; currency on cost invoice; cost invoice details fields

The cost invoice view gives the cost's total, its currency code and the client of its job (client_id). The view has the invoice line (line) only when the cost is already billed on an invoice.

## Invoice line details show the invoice number and date, blank when missing

Lines 291-298

Asked as: invoice number missing on expense; which invoice billed expense; blank invoice date on cost; expense invoice number blank

The details page of an expense or PO shows the invoice line that billed it, with the invoice number and date. If the invoice has no number or date, Turbo shows them blank. Turbo returns only one line (LIMIT 1) when it finds the cost on more than one.

## Cost invoice view: billing status and who can bill

Lines 300-322

Asked as: can I bill this cost; bill button missing; is cost already invoiced; who can bill expense

A cost's invoice view shows the invoice line that bills it, if any, and the cost total and currency (default CAD). Holders of invoicing or report can open the invoice (can_open_invoice). The cost can be billed (can_bill) only when it has no line, it is billable, its job exists, and the viewer is an active invoicing user.

## Expense invoice view limits who sees it and when billable

Lines 324-355

Asked as: expense not found; can't see expense invoice view; expense not authorized invoice; why expense invisible invoice

The expense invoice view works only for expenses the viewer can see; an admin can also see any committed expense. Otherwise Turbo says "expense not found or not authorized". An expense is billable once committed.

## Purchase order invoice view needs a visible PO

Lines 357-380

Asked as: purchase order not visible; PO not found invoice; can't open PO invoice; PO invoice page error

The PO invoice view works only for POs the viewer can see. Otherwise Turbo says "purchase order not found or not visible". A PO is billable when its status is Active or Closed.
