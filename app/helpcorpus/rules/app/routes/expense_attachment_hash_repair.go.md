<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_attachment_hash_repair.go","sha256":"3b34f8b3e80c35b51a5a91657fa6b1e246ed774c72e53825f48913ba2fe786a0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_attachment_hash_repair.go

## Mark-missing response reports previous document and marked or noop flags

Lines 34-41

The mark-missing response gives the expense id, the new updated value, the missing reason (attachment_missing_reason) and the id of the previous attachment document (previous_attachment_document). It also gives two flags, marked and noop, which show whether Turbo changed the expense or left it as it was.

## Hash repair messages for expense attachments

Lines 45-53

Hash audit and replace use these messages: "expense attachment not found" (no storage path), "expense attachment file not found" (file missing), "failed to hash expense attachment", "expense attachment target changed", "attachment hash owner changed; rerun audit before replacing" (updated value is stale), and "calculated hash already belongs to another attachment" (unique conflict).

## Only admins can mark an expense attachment missing

Lines 63-67

Marking an expense attachment as missing requires an admin. The check is requireAdminForStoredFileHashRepair, and its refusal message is defined in another file.

## Marking attachment missing requires updated value and reason

Lines 73-80

Turbo refuses the request with "updated is required" when the updated value is blank or only spaces. It refuses with "attachment missing reason is required" when the reason is blank or only spaces. Both values are trimmed before use.

## Whitespace is trimmed from updated value and reason

Lines 80-80

When an admin marks an expense attachment missing, Turbo trims leading and trailing spaces from the updated value and from the reason. It then compares and stores the trimmed values. The reason is saved in attachment_missing_reason.

## Successful mark-missing request returns the result

Lines 80-84

When marking an attachment missing succeeds, Turbo returns the result of markExpenseAttachmentMissing. When it fails, the error goes through the expense attachment hash error handler, which shows the specific message or a generic one.

## Attachment hash audit runs on the expense's attachment document

Lines 88-96

The hash audit for an expense looks up the expense's attachment document (attachment_document) and checks that document's attachment file. The audit response includes the expense id (expense_id). Whether the hash matches comes from auditStoredFileHash, which is in another file.

## Hash replace needs the updated value from the audit

Lines 98-106

Replacing an expense attachment's hash takes the expected updated value (expectedUpdated) and runs against the same document target as the audit. Turbo refuses the replace when the owner has changed since the audit. The checks themselves are in replaceStoredFileHash, which is in another file. The response includes the expense id.

## Expense id is required to mark attachment missing

Lines 109-112

Turbo refuses with "expense id is required" when the expense id is blank after trimming.

## Marking attachment missing fails for unknown expense

Lines 121-124

When the expense does not exist, Turbo refuses with "expense not found".

## Stale expense blocks marking attachment missing

Lines 126-132

Turbo compares the expense's current updated timestamp with the one supplied. If they differ, it refuses with "expense changed; refresh before marking attachment missing". The same message appears if the expense changes between the check and the write (the update affects no row).

## Repeating the same missing reason changes nothing

Lines 134-141

If the expense has no attachment document and is already marked missing with the same reason, Turbo makes no change. It returns the current updated timestamp and flags the result as a no-op (noop).

## Marking missing clears the attachment document

Lines 134-171

Marking an attachment missing empties the expense's attachment document (attachment_document), stores the reason (attachment_missing_reason) and sets a new updated timestamp. The response reports the previous document id and marked=true. An expense that still has a document can be marked missing, which replaces the document with the reason.

## Missing reason cannot be changed once set

Lines 136-143

If the expense has no attachment document and already has a different missing reason, Turbo refuses with "attachment is already marked missing; reason cannot be changed".

## Concurrent change blocks marking attachment missing

Lines 145-167

When an admin marks an expense attachment missing, Turbo changes the expense only if its updated timestamp still equals the one supplied. If the expense changed in the meantime, Turbo refuses with "expense changed; refresh before marking attachment missing".

## Marking missing sets a new timestamp and reports the result

Lines 145-178

A successful mark clears attachment_document, saves the reason, and sets the expense's updated timestamp to the current time. The response gives marked as true, the new updated value, the reason, and the previous attachment document id. The noop flag stays false.

## Hash audit and replace look up the expense's attachment document

Lines 181-199

The attachment hash audit and replace act on the expense's attachment document record in the expense documents collection. They use its attachment file and stored hash (attachment_hash). Turbo refuses with "expense id is required" for a blank id and "expense not found" for an unknown expense. It refuses with "expense attachment document not found" when the linked document is missing.

## Expense without document-backed attachment cannot be hash-checked

Lines 201-201

If the expense has no attachment document, the hash audit and replace refuse with "expense has no document-backed attachment".

## Hash check targets the expense document's attachment file

Lines 204-223

For an expense attachment, the hash check works on the linked expense document record. It uses that record's attachment file in storage and its stored attachment_hash value. The updated timestamp that must match comes from the document record, not from the expense.

## Document with no attachment file has no storage path

Lines 204-209

If the expense document record has no attachment filename, the hash target has an empty filename and an empty storage path. The shared hash check then handles that case. Its message for an empty path is "expense attachment not found".

## Timestamp load failure message for attachment hash owner

Lines 210-213

If Turbo cannot load the updated timestamp of the attachment's document record, it fails with "failed to load attachment hash owner timestamp".

## Unexpected hash repair errors show generic message

Lines 225-231

An error that is not a known repair error is shown as "expense attachment hash repair failed".

## Audit and replace responses include the expense id

Lines 233-244

The expense attachment hash audit and replace responses both include the expense_id. They also include the results of the shared stored-file hash audit or replacement.
