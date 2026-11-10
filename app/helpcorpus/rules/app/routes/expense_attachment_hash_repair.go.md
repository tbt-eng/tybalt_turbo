<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_attachment_hash_repair.go","sha256":"3b34f8b3e80c35b51a5a91657fa6b1e246ed774c72e53825f48913ba2fe786a0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_attachment_hash_repair.go

## Mark-missing response reports previous document and marked or noop flags

Lines 34-41

Asked as: mark missing response fields; what does marked flag mean; noop flag meaning; previous document id in response

The mark-missing response gives the expense id, the new updated value, the missing reason (attachment_missing_reason) and the id of the previous attachment document (previous_attachment_document). It also gives two flags, marked and noop, which show whether Turbo changed the expense or left it as it was.

## Hash repair messages for expense attachments

Lines 45-53

Asked as: hash audit error message; attachment hash stale error; hash already used by another file; file not found during hash check

Hash audit and replace use these messages: "expense attachment not found" (no storage path), "expense attachment file not found" (file missing), "failed to hash expense attachment", "expense attachment target changed", "attachment hash owner changed; rerun audit before replacing" (updated value is stale), and "calculated hash already belongs to another attachment" (unique conflict).

## Only admins can mark an expense attachment missing

Lines 63-67

Asked as: who can mark attachment missing; non admin cannot mark missing; permission denied marking attachment; user not allowed to mark missing

Marking an expense attachment as missing requires an admin. The check is requireAdminForStoredFileHashRepair, and its refusal message is defined in another file.

## Marking attachment missing requires updated value and reason

Lines 73-80

Asked as: updated is required error; missing reason blank error; blank reason when marking missing; empty updated value refused

Turbo refuses the request with "updated is required" when the updated value is blank or only spaces. It refuses with "attachment missing reason is required" when the reason is blank or only spaces. Both values are trimmed before use.

## Whitespace is trimmed from updated value and reason

Lines 80-80

Asked as: extra spaces in missing reason; spaces around updated value; whitespace ignored in reason; trailing spaces saved on attachment

When an admin marks an expense attachment missing, Turbo trims leading and trailing spaces from the updated value and from the reason. It then compares and stores the trimmed values. The reason is saved in attachment_missing_reason.

## Successful mark-missing request returns the result

Lines 80-84

Asked as: what happens after marking missing; mark missing success response; marking attachment missing failed error

When marking an attachment missing succeeds, Turbo returns the result of markExpenseAttachmentMissing. When it fails, the error goes through the expense attachment hash error handler, which shows the specific message or a generic one.

## Attachment hash audit runs on the expense's attachment document

Lines 88-96

Asked as: audit expense attachment hash; check hash of expense attachment; verify hash of expense receipt; does expense attachment hash match

The hash audit for an expense looks up the expense's attachment document (attachment_document) and checks that document's attachment file. The audit response includes the expense id (expense_id). Whether the hash matches comes from auditStoredFileHash, which is in another file.

## Hash replace needs the updated value from the audit

Lines 98-106

Asked as: replace attachment hash; replace hash after audit; hash replace refused owner changed; fix wrong attachment hash

Replacing an expense attachment's hash takes the expected updated value (expectedUpdated) and runs against the same document target as the audit. Turbo refuses the replace when the owner has changed since the audit. The checks themselves are in replaceStoredFileHash, which is in another file. The response includes the expense id.

## Expense id is required to mark attachment missing

Lines 109-112

Asked as: expense id is blank; missing expense id error; no expense id entered; expense id needed to mark

Turbo refuses with "expense id is required" when the expense id is blank after trimming.

## Marking attachment missing fails for unknown expense

Lines 121-124

Asked as: expense does not exist; can't find expense for attachment; expense id not found; unknown expense marking missing

When the expense does not exist, Turbo refuses with "expense not found".

## Stale expense blocks marking attachment missing

Lines 126-132

Asked as: expense was edited, refresh; someone else changed expense; out of date expense mark missing; expense modified before marking

Turbo compares the expense's current updated timestamp with the one supplied. If they differ, it refuses with "expense changed; refresh before marking attachment missing". The same message appears if the expense changes between the check and the write (the update affects no row).

## Repeating the same missing reason changes nothing

Lines 134-141

Asked as: already marked missing same reason; mark missing twice; nothing changed marking missing; same missing reason again

If the expense has no attachment document and is already marked missing with the same reason, Turbo makes no change. It returns the current updated timestamp and flags the result as a no-op (noop).

## Marking missing clears the attachment document

Lines 134-171

Asked as: removes attachment from expense; mark attachment missing deletes document; attachment document removed from expense; replace document with missing reason

Marking an attachment missing empties the expense's attachment document (attachment_document), stores the reason (attachment_missing_reason) and sets a new updated timestamp. The response reports the previous document id and marked=true. An expense that still has a document can be marked missing, which replaces the document with the reason.

## Missing reason cannot be changed once set

Lines 136-143

Asked as: change missing reason; already missing different reason; edit attachment missing reason error; reason already set on expense

If the expense has no attachment document and already has a different missing reason, Turbo refuses with "attachment is already marked missing; reason cannot be changed".

## Concurrent change blocks marking attachment missing

Lines 145-167

Asked as: two admins at once; expense changed while marking; conflict marking attachment missing; saving clashed with another edit

When an admin marks an expense attachment missing, Turbo changes the expense only if its updated timestamp still equals the one supplied. If the expense changed in the meantime, Turbo refuses with "expense changed; refresh before marking attachment missing".

## Marking missing sets a new timestamp and reports the result

Lines 145-178

Asked as: updated time after marking missing; timestamp changes when marked missing; noop false meaning; new updated value after mark

A successful mark clears attachment_document, saves the reason, and sets the expense's updated timestamp to the current time. The response gives marked as true, the new updated value, the reason, and the previous attachment document id. The noop flag stays false.

## Hash audit and replace look up the expense's attachment document

Lines 181-199

Asked as: expense documents attachment lookup; attachment document not found; linked document missing for expense; expense attachment document not found

The attachment hash audit and replace act on the expense's attachment document record in the expense documents collection. They use its attachment file and stored hash (attachment_hash). Turbo refuses with "expense id is required" for a blank id and "expense not found" for an unknown expense. It refuses with "expense attachment document not found" when the linked document is missing.

## Expense without document-backed attachment cannot be hash-checked

Lines 201-201

Asked as: expense has no attachment; no attachment to check hash; expense without document attachment; nothing to hash for expense

If the expense has no attachment document, the hash audit and replace refuse with "expense has no document-backed attachment".

## Hash check targets the expense document's attachment file

Lines 204-223

Asked as: which file hash check uses; attachment_hash on document; timestamp comes from document; hash check on expense document

For an expense attachment, the hash check works on the linked expense document record. It uses that record's attachment file in storage and its stored attachment_hash value. The updated timestamp that must match comes from the document record, not from the expense.

## Document with no attachment file has no storage path

Lines 204-209

Asked as: document without attachment file; empty storage path attachment; attachment file name missing; expense attachment not found error

If the expense document record has no attachment filename, the hash target has an empty filename and an empty storage path. The shared hash check then handles that case. Its message for an empty path is "expense attachment not found".

## Timestamp load failure message for attachment hash owner

Lines 210-213

Asked as: cannot load attachment owner timestamp; failed to load hash owner timestamp; timestamp load failure hash repair; attachment timestamp could not be read

If Turbo cannot load the updated timestamp of the attachment's document record, it fails with "failed to load attachment hash owner timestamp".

## Unexpected hash repair errors show generic message

Lines 225-231

Asked as: unknown hash repair error; hash repair failed generic; unexpected error repairing hash; expense attachment hash repair failed

An error that is not a known repair error is shown as "expense attachment hash repair failed".

## Audit and replace responses include the expense id

Lines 233-244

Asked as: expense id in audit response; replace response expense id; audit result includes expense; hash audit output fields

The expense attachment hash audit and replace responses both include the expense_id. They also include the results of the shared stored-file hash audit or replacement.
