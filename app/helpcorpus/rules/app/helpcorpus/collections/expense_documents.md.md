<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_documents.md","sha256":"4a6d5884ebd6c91d2c6e468396a8b082b1defa98fae4dee5014fde7c261f1113","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_documents.md

## Nobody can list expense documents through the API

Lines 7-9

Asked as: list expense documents API; can't see expense documents; fetch all receipts API; no access to expense list

The expense_documents collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one expense document through the API

Lines 11-13

Asked as: view single expense document; open one receipt via API; can't read expense record; get expense document by id

The expense_documents collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create expense documents through the API

Lines 15-17

Asked as: create expense document API; upload receipt through API; can't add expense document; post new expense record

The expense_documents collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update expense documents through the API

Lines 19-21

Asked as: edit expense document API; change receipt record; update expense via API; can't modify expense document

The expense_documents collection has no update rule. Nobody can change a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete expense documents through the API

Lines 23-25

Asked as: delete expense document API; remove receipt record; can't delete expense; erase expense document

The expense_documents collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Expense document attachment: one file, 5 MB, allowed types

Lines 32-32

Asked as: attachment too large; upload file size limit; receipt file type not allowed; more than one attachment

Each expense document requires one attachment (attachment) and accepts only one file. The file must be 5 MB or smaller (5242880 bytes). Allowed types are PDF, JPEG, PNG and HEIC.

## Attachment hash must be 64 hexadecimal characters

Lines 33-33

Asked as: attachment hash invalid; hash wrong length; checksum error receipt; hash not hexadecimal

Each expense document requires a hash of its attachment (attachment_hash). The hash must be exactly 64 characters, and each character must be a hexadecimal digit (0-9, a-f, A-F).

## Each expense document records who uploaded it

Lines 34-34

Asked as: who uploaded receipt; uploader missing; uploaded by required; no user on document

The uploader (uploaded_by) is required and points to one user.

## Vendor invoice number length limit

Lines 38-39

Asked as: invoice number too long; vendor invoice number characters; invoice reference length limit; number key too long

The vendor invoice number (vendor_invoice_number) and its key (vendor_invoice_number_key) each hold at most 100 characters.

## Vendor invoice date format

Lines 40-40

Asked as: invoice date format; vendor invoice date wrong; date not accepted invoice; invoice date entry error

The vendor invoice date (vendor_invoice_date) must be in the form YYYY-MM-DD, or be empty.

## Same attachment cannot be stored twice

Lines 50-50

Asked as: duplicate receipt file; same attachment uploaded twice; file already uploaded; receipt already exists

A unique index on attachment_hash stops two expense documents from having the same attachment hash. The index applies only when the hash is not empty.

## A vendor cannot have the same invoice number twice

Lines 51-51

Asked as: duplicate invoice number vendor; vendor invoice already entered; same invoice number twice; invoice number already used

A unique index on the vendor and the vendor invoice number key (vendor, vendor_invoice_number_key) stops two expense documents from having the same pair. The index applies only when the number key is not empty. Documents with no number key are not checked.
