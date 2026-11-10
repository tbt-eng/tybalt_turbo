<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_documents.md","sha256":"f1ca91d3b72c352839deab5a8de6e86f77205cd3f79350161df80d7560a984c0","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_documents.md

## Nobody can list expense document records through the API

Lines 7-9

Asked as: list expense documents API; can't see expense docs list; api list expense records denied

The expense_documents collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one expense document record through the API

Lines 11-13

Asked as: view single expense document; open one receipt record API; can't read expense document api

The expense_documents collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create expense document records through the API

Lines 15-17

Asked as: create expense document via API; add receipt record through api; api post expense document forbidden

The expense_documents collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create records.

## Nobody can update expense document records through the API

Lines 19-21

Asked as: edit expense document API; change receipt record api; update expense doc not allowed

The expense_documents collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update records.

## Nobody can delete expense document records through the API

Lines 23-25

Asked as: delete expense document API; remove receipt record api; can't delete expense docs

The expense_documents collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete records.

## Expense document attachment: one file, 5 MB, allowed types

Lines 32-32

Asked as: receipt file too big; upload limit receipt size; which file types for receipts; multiple attachments on expense

Each expense document (expense_documents) has one required attachment (attachment), and only one file. The file must be at most 5242880 bytes (5 MB). The allowed types are PDF (application/pdf), JPEG (image/jpeg), PNG (image/png) and HEIC (image/heic).

## Expense document hash must be 64 hex characters

Lines 33-33

Asked as: attachment hash invalid; checksum wrong for receipt; hash error on upload; sha256 not accepted

Each expense document needs a file hash (attachment_hash) of exactly 64 characters. The characters must be hexadecimal digits, 0-9 and a-f in either case.

## Expense document needs one uploader

Lines 34-34

Asked as: who uploaded the receipt; missing uploader on document; expense doc no owner user

Each expense document must name one user as the uploader (uploaded_by). The field is required and holds a single user.

## Two expense documents cannot share the same file hash

Lines 40-42

Asked as: duplicate receipt upload; same file uploaded twice; receipt already exists error; repeat scan of document rejected

A unique index (idx_expense_documents_attachment_hash) stops two expense documents from having the same attachment hash. This blocks a duplicate upload of the same file. The index ignores records whose hash is empty.
