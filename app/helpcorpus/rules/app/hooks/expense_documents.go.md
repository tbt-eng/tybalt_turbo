<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expense_documents.go","sha256":"41202d40a9ca8db66a752f1370510897bb3353602c3c96eb85d6f41ba86ba8ca","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expense_documents.go

## Attachment file types allowed

Lines 17-22

Asked as: wrong file type; upload pdf or photo; invalid attachment type; heic file rejected; can't attach this file

An expense attachment must be a PDF, JPEG, PNG or HEIC file (application/pdf, image/jpeg, image/png, image/heic). Other types are rejected as an invalid attachment (invalid_mime_type).

## Failed expense saves delete newly uploaded attachment files

Lines 24-68

Asked as: failed save leftover files; orphan attachment files; files left after error; save failed uploads remain; clean up uploaded receipts

Turbo records each attachment file it stores during an expense save. After a rollback, the cleanup removes those files from storage (CleanupExpenseDocumentUploads). Turbo only logs a failure to delete a file and does not report it to the user.

## Reusing another expense's attachment needs the bookkeeper claim

Lines 72-76

Asked as: reuse someone's receipt; copy receipt from other expense; need bookkeeper to reuse; book keeper claim needed; use another expense's attachment

When an expense names a source expense (source_expense), Turbo refuses the save unless the user has the book keeper claim. The message is "book_keeper claim required to reuse an attachment" (code book_keeper_required).

## Source expense takes priority over an uploaded file

Lines 72-85

Asked as: uploaded file ignored; source expense overrides upload; why my new file not used; attachment hash ignored; source wins over file

If an expense has a source expense, Turbo uses it and ignores any uploaded attachment hash.

## Reusing an attachment needs a purchase order on the expense

Lines 77-79

Asked as: reuse receipt without PO; attachment needs purchase order; no PO reuse attachment; link receipt needs PO; source expense no purchase order

When an expense names a source expense, Turbo refuses the save if the expense is not backed by a purchase order. The message is "purchase order is required to reuse an attachment". The bookkeeper claim check runs first.

## Source expense supplies the attachment document

Lines 80-80

Asked as: copy source receipt; inherit attachment from expense; source has no receipt; pull receipt from other expense; attachment from source expense

With a source expense, the new expense takes the source's attachment document (attachment_document). If the source has no document-backed attachment, Turbo refuses with "source expense has no document-backed attachment".

## Existing expense keeps its attachment document when none is supplied

Lines 87-99

Asked as: change receipt on saved expense; edit expense receipt; swap attachment document; expense lost its receipt; update receipt on existing expense

When an existing expense is saved with no new file and no source expense, Turbo keeps its original attachment document. Clients cannot point the expense at a different document. A new expense in this case gets no document.

## Some payment types have no attachment document

Lines 102-115

Asked as: mileage no receipt; allowance attachment removed; personal reimbursement receipt gone; receipt disappeared after payment type; why receipt cleared

Expenses with payment type Allowance, Mileage or PersonalReimbursement get their attachment document cleared. Other payment types keep it.

## Uploaded file matching an unused document is reused

Lines 118-129

Asked as: same file uploaded again; receipt already in system; reuse existing receipt upload; file matches old receipt; previously uploaded receipt

If the uploaded file's hash matches an existing expense document that no other expense uses, Turbo links that document to this expense. No claim or purchase order is needed.

## Duplicate upload is refused without the bookkeeper claim

Lines 130-132

Asked as: file already used; receipt already attached elsewhere; duplicate receipt upload; same receipt another expense; upload rejected duplicate

If the uploaded file's hash matches a document that another expense already uses, Turbo refuses a user without the book keeper claim. The message is "This file has already been uploaded to another expense" (duplicate_file).

## Duplicate upload by a bookkeeper needs a purchase order

Lines 133-136

Asked as: bookkeeper duplicate receipt; bookkeeper reuse file; duplicate upload as bookkeeper; bookkeeper receipt already used; bookkeeper needs PO for file

If a bookkeeper uploads a file already used by another expense, Turbo allows it only when the expense is backed by a purchase order. Otherwise it refuses with "purchase order is required to reuse an attachment".

## A new upload must be exactly one file

Lines 143-152

Asked as: multiple files attached; attach two receipts; only one file allowed; expected one attachment; too many attachments

When the file is new to Turbo, the expense must carry exactly one attachment file. Otherwise Turbo refuses with "expected one attachment file" (invalid_file_count).

## Source expense must have a document-backed attachment

Lines 161-172

Asked as: source expense no receipt; source missing attachment; source has no document; reuse expense without receipt; missing receipt on source

When an expense is saved with a source expense (source_expense) and the source has no attachment document (attachment_document), Turbo refuses the save. The message is "source expense has no document-backed attachment" with code missing_attachment.

## Source expense must be visible to the user

Lines 174-211

Asked as: source expense not visible; can't see other expense; cannot access expense receipt; expense not found when reusing; can't pick that expense

A user can reuse an attachment only from a source expense they can see, and must be signed in. Visible expenses are: their own (uid or creator); ones they approve that are submitted; approved ones, if they have the commit claim; and committed ones, if they have the report or admin claim. Otherwise Turbo says "source expense not found or not visible".

## Source expense lookup requires sign-in

Lines 175-177

Asked as: not logged in reuse; reuse receipt signed out; sign in to reuse attachment; authentication needed attachment; session expired attachment

Turbo refuses to reuse an attachment from a source expense when the user is not authenticated. The error is 'attachment reuse is not allowed' with 'authentication is required' (authentication_required).

## Who can see a source expense for attachment reuse

Lines 179-211

Asked as: who can see expense; approver can't see expense; committer access expense; why expense hidden; visibility of source expense

A user can use a source expense only if it is visible to them. They must be its owner (uid) or creator, or its approver once it is submitted. Committers (the commit claim) see approved expenses. Report and admin users see committed expenses. Otherwise the error is 'source expense not found or not visible' (not_found).

## Documents are matched by file hash

Lines 216-233

Asked as: how files are matched; blank hash receipt; same file different hash; document lookup failed; hash matching receipt

Turbo identifies an existing document by the hash of the uploaded file (attachment_hash). A blank hash never matches a document. If the lookup fails, Turbo treats it as no match and goes on to create a new document.

## A document counts as in use when another expense references it

Lines 235-254

Asked as: receipt still in use; which expense uses document; document counted as used; my own reference receipt; receipt linked to other expense

Turbo treats a document as in use when any expense other than the current one references it. The current expense's own reference does not count.

## New attachment document records its uploader

Lines 279-287

Asked as: who uploaded receipt; record uploader receipt; receipt upload timestamp; track receipt uploader; receipt metadata saved

A new expense document stores the file, its hash and the uploading user (uploaded_by) with creation and update times.

## Uploaded files are deleted when the save is rolled back

Lines 288-292

Asked as: rollback deletes files; file removed after error; uploaded file deleted on failure; save failed file gone; files removed when save fails

Turbo tracks each uploaded attachment file and deletes it if the expense save is rolled back. A file is also deleted at once if saving its document record fails.

## Missing attachment file is refused

Lines 307-309

Asked as: no file attached; attachment is required; save without receipt error; missing file on create; empty upload error

If no file is given when a document is created, Turbo refuses with "attachment is required".

## Attachment size limit is 5 MB

Lines 311-313

Asked as: file too big; receipt over 5mb; attachment size limit; large photo rejected; file size error

Turbo rejects an attachment larger than 5,242,880 bytes (5 MB) as an invalid attachment (invalid_file_size).

## Duplicate upload error message

Lines 321-332

Asked as: duplicate attachment detected; duplicate error message; why says already uploaded; attachment field error duplicate; duplicate warning receipt

When a duplicate upload is refused, Turbo says "duplicate attachment detected". The attachment field shows "This file has already been uploaded to another expense" (duplicate_file).

## Attachment reuse errors use one general message

Lines 334-342

Asked as: attachment reuse not allowed; generic reuse error; can't reuse attachment message; reuse error unclear; why reuse refused

When Turbo refuses to reuse an attachment, the error message is "attachment reuse is not allowed". The field and the detail message depend on the cause.
