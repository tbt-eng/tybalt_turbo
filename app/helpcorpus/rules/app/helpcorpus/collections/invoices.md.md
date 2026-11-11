<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/invoices.md","sha256":"6d48f25cab3f699e1ed8a883939a5df24cd22e870bda795ab442787f7170628c","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/invoices.md

## Who can list and view client invoices

Lines 7-17

Asked as: who can see invoices; list client invoices access; view invoice records permission; cannot see invoice list

Invoicing staff (the invoicing claim) and people with the report claim can list invoice records. They can also view one invoice record. Everyone else cannot.

## Only invoicing staff can create invoices

Lines 19-23

Asked as: create new invoice; report role make invoice; who can add invoices; invoice creation permission

Only invoicing staff (the invoicing claim) can create an invoice record through the API. Report-claim holders can read invoices but cannot create them.

## Invoicing staff can update only Draft invoices

Lines 25-29

Asked as: edit issued invoice; cannot change invoice after issue; update invoice locked; modify sent invoice

Invoicing staff (the invoicing claim) can update an invoice only while its status is "Draft". Once the status is "Issued", nobody can update it through the API. Turbo's own server code can still change it.

## Invoicing staff can delete only Draft invoices

Lines 31-35

Asked as: delete issued invoice; remove invoice permanently; cannot erase invoice; drop an invoice

Invoicing staff (the invoicing claim) can delete an invoice only while its status is "Draft". Nobody can delete an "Issued" invoice through the API. Turbo's own server code can still delete it.

## Invoice status values and required fields

Lines 42-46

Asked as: invoice status options; missing client on invoice; which status invoice can have; invoice needs a client

An invoice's status (status) is required and is either "Draft" or "Issued". Its client (client) is also required and must be one client.

## Invoice number format

Lines 43-43

Asked as: invoice number wrong; invoice number digits; invoice number rejected; invalid invoice reference

An invoice number (number) is either empty or 7 digits: 2 digits, then a month 01-12, then 3 digits. Any other value is refused.

## Invoice date format

Lines 44-44

Asked as: invoice date rejected; date format on invoice; invoice date wrong format; invoice date not saving

An invoice date (date) is either empty or written as YYYY-MM-DD. Any other value is refused.

## Invoice numbers must be unique

Lines 55-57

Asked as: duplicate invoice number; same invoice number twice; invoice number already used; reuse invoice number

Two invoices cannot have the same number. Empty numbers are exempt, so many invoices can have no number.
