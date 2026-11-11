<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/invoices.ts","sha256":"237d5f3d1acdb1250263e556d0f9f6b9e9a865c0f2e49c127449b7226a61ccfc","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/invoices.ts

## Issuing an invoice is a separate action

Lines 4-62

Asked as: issue invoice; draft invoice not sent; how to issue a draft; finalise invoice; invoice still draft

A draft invoice is issued through a separate issue action (POST /api/invoices/{id}/issue). Drafts and lines are written through the invoices and invoice_lines collections.

## Invoice line can be listed without a link

Lines 24-31

Asked as: invoice line no link; can't open expense from invoice; line without link; greyed out invoice line; no access to expense on invoice

Each invoice line has a can_open flag. It is false when the viewer may not open the expense or PO, and Turbo then lists the line without a link.

## Billable costs are found by job for an invoice

Lines 33-38

Asked as: billable costs by job; which jobs have billable costs; find costs to bill; jobs to invoice; billable count per job

Turbo lists the jobs that have billable costs for an invoice, with a count of billable costs for each job (billable_count). Choosing a job lists that job's billable expenses and POs.

## Cost details page shows its client invoice and billing ability

Lines 40-54

Asked as: expense billed on invoice; where was this cost invoiced; can't bill this expense; bill button missing; open the invoice from expense

The details page of an expense or PO loads the invoice it was billed on, if any (line: invoice number and date, and the amount billed). It also gets can_bill and can_open_invoice flags, which tell whether the viewer may bill the cost or open the invoice. Another file decides these values.

## Who can view client invoices

Lines 56-58

Asked as: who can see invoices; no access to client invoices; can't view invoices; invoice list hidden; invoicing permission needed

Staff with the invoicing claim or the report claim can view client invoices (canViewInvoices). Staff with neither claim cannot.

## Billing a cost adds a line to an issued invoice

Lines 91-104

Asked as: add expense to invoice; bill a PO to client; charge cost to invoice; invoice line from expense; add cost to client invoice

Billing an expense or a purchase order on a client invoice creates an invoice line (billCost) with the invoice, the kind (expense or po), the cost and an amount. The comment says the invoice must already be issued. The server rules in other files decide who may do this.

## How a cost is named in invoice lists

Lines 106-111

Asked as: invoice shows expense name; PO name on invoice; why expense has PO number; expense label without description; how costs appear in invoice

A PO is named "PO" plus its PO number. An expense is named by its description, or "Expense" when it has none. An expense that has a PO number also shows "(PO number)" after its name.
