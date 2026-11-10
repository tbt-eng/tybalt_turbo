<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expense_documents.go","sha256":"41202d40a9ca8db66a752f1370510897bb3353602c3c96eb85d6f41ba86ba8ca","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expense_documents.go

## Attachment file types allowed

Lines 17-22

An expense attachment must be a PDF, JPEG, PNG or HEIC file (application/pdf, image/jpeg, image/png, image/heic). Other types are rejected as an invalid attachment (invalid_mime_type).

## Failed expense saves delete newly uploaded attachment files

Lines 24-68

Turbo records each attachment file it stores during an expense save. After a rollback, the cleanup removes those files from storage (CleanupExpenseDocumentUploads). Turbo only logs a failure to delete a file and does not report it to the user.

## Reusing another expense's attachment needs the bookkeeper claim

Lines 72-76

When an expense names a source expense (source_expense), Turbo refuses the save unless the user has the book keeper claim. The message is "book_keeper claim required to reuse an attachment" (code book_keeper_required).

## Source expense takes priority over an uploaded file

Lines 72-85

If an expense has a source expense, Turbo uses it and ignores any uploaded attachment hash.

## Reusing an attachment needs a purchase order on the expense

Lines 77-79

When an expense names a source expense, Turbo refuses the save if the expense is not backed by a purchase order. The message is "purchase order is required to reuse an attachment". The bookkeeper claim check runs first.

## Source expense supplies the attachment document

Lines 80-80

With a source expense, the new expense takes the source's attachment document (attachment_document). If the source has no document-backed attachment, Turbo refuses with "source expense has no document-backed attachment".

## Existing expense keeps its attachment document when none is supplied

Lines 87-99

When an existing expense is saved with no new file and no source expense, Turbo keeps its original attachment document. Clients cannot point the expense at a different document. A new expense in this case gets no document.

## Some payment types have no attachment document

Lines 102-115

Expenses with payment type Allowance, Mileage or PersonalReimbursement get their attachment document cleared. Other payment types keep it.

## Uploaded file matching an unused document is reused

Lines 118-129

If the uploaded file's hash matches an existing expense document that no other expense uses, Turbo links that document to this expense. No claim or purchase order is needed.

## Duplicate upload is refused without the bookkeeper claim

Lines 130-132

If the uploaded file's hash matches a document that another expense already uses, Turbo refuses a user without the book keeper claim. The message is "This file has already been uploaded to another expense" (duplicate_file).

## Duplicate upload by a bookkeeper needs a purchase order

Lines 133-136

If a bookkeeper uploads a file already used by another expense, Turbo allows it only when the expense is backed by a purchase order. Otherwise it refuses with "purchase order is required to reuse an attachment".

## A new upload must be exactly one file

Lines 143-152

When the file is new to Turbo, the expense must carry exactly one attachment file. Otherwise Turbo refuses with "expected one attachment file" (invalid_file_count).

## Source expense must have a document-backed attachment

Lines 161-172

When an expense is saved with a source expense (source_expense) and the source has no attachment document (attachment_document), Turbo refuses the save. The message is "source expense has no document-backed attachment" with code missing_attachment.

## Source expense must be visible to the user

Lines 174-211

A user can reuse an attachment only from a source expense they can see, and must be signed in. Visible expenses are: their own (uid or creator); ones they approve that are submitted; approved ones, if they have the commit claim; and committed ones, if they have the report or admin claim. Otherwise Turbo says "source expense not found or not visible".

## Source expense lookup requires sign-in

Lines 175-177

Turbo refuses to reuse an attachment from a source expense when the user is not authenticated. The error is 'attachment reuse is not allowed' with 'authentication is required' (authentication_required).

## Who can see a source expense for attachment reuse

Lines 179-211

A user can use a source expense only if it is visible to them. They must be its owner (uid) or creator, or its approver once it is submitted. Committers (the commit claim) see approved expenses. Report and admin users see committed expenses. Otherwise the error is 'source expense not found or not visible' (not_found).

## Documents are matched by file hash

Lines 216-233

Turbo identifies an existing document by the hash of the uploaded file (attachment_hash). A blank hash never matches a document. If the lookup fails, Turbo treats it as no match and goes on to create a new document.

## A document counts as in use when another expense references it

Lines 235-254

Turbo treats a document as in use when any expense other than the current one references it. The current expense's own reference does not count.

## New attachment document records its uploader

Lines 279-287

A new expense document stores the file, its hash and the uploading user (uploaded_by) with creation and update times.

## Uploaded files are deleted when the save is rolled back

Lines 288-292

Turbo tracks each uploaded attachment file and deletes it if the expense save is rolled back. A file is also deleted at once if saving its document record fails.

## Missing attachment file is refused

Lines 307-309

If no file is given when a document is created, Turbo refuses with "attachment is required".

## Attachment size limit is 5 MB

Lines 311-313

Turbo rejects an attachment larger than 5,242,880 bytes (5 MB) as an invalid attachment (invalid_file_size).

## Duplicate upload error message

Lines 321-332

When a duplicate upload is refused, Turbo says "duplicate attachment detected". The attachment field shows "This file has already been uploaded to another expense" (duplicate_file).

## Attachment reuse errors use one general message

Lines 334-342

When Turbo refuses to reuse an attachment, the error message is "attachment reuse is not allowed". The field and the detail message depend on the cause.
