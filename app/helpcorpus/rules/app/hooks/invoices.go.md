<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/invoices.go","sha256":"7bd02e267b83bf92b5efa9fb8cad1ad14e8319c8ca7c3fce6ee272ddebaee700","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/invoices.go

## Invoice number format

Lines 19-19

Asked as: invoice number rejected; wrong invoice number format; invoice number too short; seven digit invoice number; invoice number month

When a client invoice has a number, the number must be seven digits in the form YYMMNNN, with a month from 01 to 12. Otherwise Turbo refuses with "invoice number must be seven digits, YYMMNNN, with a month from 01 to 12". An invoice without a number skips this check while it is a draft.

## Only active invoicing staff change invoices

Lines 31-53

Asked as: cannot change invoice; can't edit invoices; no access to invoices; not signed in invoice; inactive user invoice error

Turbo refuses invoice and invoice line changes from a user who is not signed in ("authentication required"). It also refuses an inactive user ("inactive users cannot change invoices"). It also refuses a user without the invoicing claim ("the invoicing claim is required").

## New invoice starts as a draft

Lines 55-65

Asked as: new invoice status; why is my invoice draft; created invoice shows draft; issued by blank; new invoice creator

When invoicing staff create an invoice, Turbo sets its status to Draft and its creator to the person creating it. Turbo clears the issued-by (issued_by) and issued-at (issued_at) fields. Turbo then runs the draft checks.

## Only drafts can be edited

Lines 67-82

Asked as: can't edit invoice; invoice locked; change issued invoice; edit sent invoice; invoice cannot be changed

Turbo lets invoicing staff edit a draft's client, number, date and total. Turbo refuses to edit an invoice that is not a Draft with "an issued invoice cannot be edited". Status, creator, issued-by and issued-at keep their saved values on every edit, so the issue route is the only way to issue an invoice.

## Only drafts can be deleted

Lines 84-94

Asked as: delete invoice button; remove issued invoice; can't delete invoice; free up invoice number; get rid of invoice

Turbo lets invoicing staff delete a draft invoice, which frees its number. Turbo refuses to delete an invoice that is not a Draft with "an issued invoice cannot be deleted".

## Draft requires a client

Lines 96-99

Asked as: invoice without client; save invoice no client; client missing on invoice; invoice client blank

Turbo refuses to save a draft invoice without a client. The message is "client is required".

## Invoice date must be valid

Lines 103-105

Asked as: invoice date rejected; wrong date format invoice; date not accepted on invoice; invoice date error

If a draft invoice has a date, it must be a real date in YYYY-MM-DD form. Otherwise Turbo refuses with "date must be a valid YYYY-MM-DD date". An empty date is allowed on a draft.

## Invoice total cannot be negative

Lines 106-108

Asked as: invoice total negative; negative invoice amount; minus total invoice; credit invoice negative

Turbo refuses a draft invoice with a negative total. The message is "total cannot be negative". A total of zero is allowed on a draft.

## Invoice number must be unique

Lines 120-141

Asked as: duplicate invoice number; number already used; same invoice number twice; invoice number taken

Turbo refuses an invoice number that another invoice already uses. The message is "invoice number <number> is already used by another invoice". The check runs when a draft is saved and when it is issued.

## Requirements to issue an invoice

Lines 143-159

Asked as: cannot issue invoice; issue button not working; why won't invoice issue; invoice won't issue; issue invoice error

Turbo issues only a Draft invoice, and refuses otherwise with "the invoice is already issued". The draft needs a number ("an invoice needs a number before it is issued"). It needs a date ("an invoice needs a date before it is issued"). Its total must be greater than zero ("an invoice needs a total greater than zero before it is issued"). The draft checks run again too.

## Billing a cost records who billed it

Lines 161-169

Asked as: who billed this cost; who added invoice line; billed by whom; who put cost on invoice

When invoicing staff add an invoice line (bill a cost), Turbo sets the line's creator to that person and then checks the line rules.

## Invoicing staff can remove invoice lines

Lines 171-174

Asked as: delete invoice line; remove line from invoice; take cost off invoice; drop line from invoice

Any active user with the invoicing claim can delete an invoice line. Turbo adds no other check on a deletion.

## Costs bill only on issued invoices

Lines 179-185

Asked as: cannot bill cost to draft; bill cost before issuing; invoice not found on line; add cost to invoice fails; cost to unissued invoice

Turbo refuses a line when its invoice is missing ("invoice not found"). It also refuses a line when the invoice is not Issued, with "costs can be billed only on an issued invoice".

## Billed amount must be positive

Lines 186-188

Asked as: zero amount line; billed amount zero; line amount must be positive; invoice line zero value

Turbo refuses an invoice line whose amount is zero or less. The message is "the amount billed must be greater than zero". The amount is compared in cents.

## Line kind and one cost per line

Lines 190-205

Asked as: line kind error; expense and po same line; line needs expense or po; invoice line both expense and po

A line's kind must be expense or po, otherwise Turbo says "kind must be expense or po". An expense line needs its expense, and a po line needs its purchase order. Turbo refuses a line that names both with "a line bills one expense or one purchase order, not both". Turbo refuses a line that names neither with "<field> is required for a line of kind <kind>".

## Only committed expenses can be billed

Lines 207-213

Asked as: bill uncommitted expense; expense not committed; pending expense billing; expense cannot be billed; expense not found

Turbo refuses to bill an expense that is not committed with "only a committed expense can be billed". Turbo says "expense not found" when the expense does not exist.

## Only Active or Closed POs can be billed

Lines 214-218

Asked as: bill purchase order fails; po not active billing; draft po billing error; cannot bill po

Turbo refuses to bill a purchase order unless its status is Active or Closed. The message is "only an Active or Closed purchase order can be billed".

## Cost must be on a job of the invoice's client

Lines 219-247

Asked as: cost without job; job belongs to other client; wrong client job invoice; job client mismatch; cost job not found

A billed cost must have a job. Without one Turbo says "only a cost on a job can be billed to a client". If the job is not found, Turbo says "the cost's job was not found". If the job's client differs from the invoice's client, Turbo says "the cost's job <number> belongs to a different client than the invoice".

## A cost can be billed only once

Lines 223-233

Asked as: cost already billed; double billing cost; cost on another invoice; bill same cost twice

Turbo refuses to bill a cost that is already on another invoice line. The message is "this cost is already billed on invoice <number>". The number is blank if that invoice has none.
