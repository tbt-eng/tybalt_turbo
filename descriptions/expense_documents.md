# Expense Documents

## Goal

Expense attachments used to be stored directly on `expenses` records and guarded by a unique hash on `expenses.attachment_hash`. That prevented a legitimate finance workflow: a `book_keeper` receives one vendor attachment that supports multiple PO-backed expenses, then needs to enter one expense per PO while reusing the same supporting document.

Turbo now works this way:

- `expenses` remain the accounting rows.
- `expense_documents` own uploaded receipt/invoice files.
- Multiple expenses may reference the same `expense_documents` row.
- Only users with the `book_keeper` claim may reuse a document that another expense already uses.
- Everyone else keeps the same UX and duplicate-file protection they have today.

The `book_keeper` claim, `expenses.creator`, `creator != uid` semantics,
bookkeeper-on-behalf PO eligibility, and creator-as-approver rules are specified
in `book_keepers.md`. This spec relies on that actor model and defines only the
attachment/document storage, duplicate-file, source-expense reuse, legacy
fallback, and backfill behavior.

This was implemented in three phases. All three are done.

1. Phase 1 cut over the app so new expense uploads use `expense_documents`, while legacy expense attachments continued to work through fallback reads.
2. Phase 2 backfilled existing legacy expense attachments into `expense_documents`.
3. Phase 3 removed the legacy fallback code and the legacy attachment columns (`expenses.attachment`, `expenses.attachment_hash`).

Mixed storage is over. Every expense attachment is an `expense_documents` row that the expense points at through `expenses.attachment_document`, and the app has no legacy fallback. The Phase 1 and Phase 2 sections below are kept as a record of the rollout. Where they describe a fallback to `expenses.attachment` or `expenses.attachment_hash`, the text says that the fallback is gone.

## Current Storage Model

PocketBase stores files under the collection id and record id:

```text
{collection_id}/{record_id}/{filename}
```

Legacy expense attachments were stored under the `expenses` collection id:

```text
o1vpz1mm7qsfoyy/{expense_id}/{expense_attachment_filename}
```

Expense documents are stored under the `expense_documents` collection id (`pbc_2089657321`):

```text
{expense_documents_collection_id}/{expense_document_id}/{document_attachment_filename}
```

That storage-prefix change was the main reason for a phased rollout. Phase 1 read both locations. Phase 2 copied old S3 objects into the new prefixes and linked the old expenses to the new document rows.

## Phase 1: Cut Over New Writes

Phase 1 is built and merged. It shipped schema changes, document-first reads with legacy fallback, and the new write path in the same release. Phase 3 later removed the legacy fallback.

### 1. Add the `expense_documents` schema

Migration `app/migrations/1777564167_created_expense_documents.go` adds a base collection named `expense_documents` with these fields:

- `attachment`: file, required, `maxSelect: 1`, max size 5 MB, MIME types PDF, JPEG, PNG and HEIC (the same as the old `expenses.attachment` field).
- `attachment_hash`: text, required, 64-character SHA-256 hex string.
- `uploaded_by`: relation to `users`, required.
- standard PocketBase `created` and `updated` autodate fields.

It has this index:

```sql
CREATE UNIQUE INDEX `idx_expense_documents_attachment_hash`
ON `expense_documents` (`attachment_hash`)
WHERE `attachment_hash` != ''
```

Every API rule on `expense_documents` is null, so there is no direct API access:

- no general list/search UI;
- no user-created direct document records through the generic collection API;
- no direct document browsing.

The backend expense save flow creates and links documents. Users interact through expenses, not through a standalone file cabinet.

### 2. Add `expenses.attachment_document`

The same migration adds a nullable relation field to `expenses`:

```text
attachment_document -> expense_documents
```

Phase 1 kept these existing fields as the fallback source for old rows:

- `expenses.attachment`
- `expenses.attachment_hash`
- the existing unique index on `expenses.attachment_hash`

Phase 3 removed all three (migration `app/migrations/1779720278_drop_expense_legacy_attachment_fields.go`).

### 3. Introduce an effective attachment resolver

There is one backend concept for "the attachment for this expense". Phase 1 fell back to the legacy `expenses` fields. Since Phase 3 there is no fallback:

```text
if expenses.attachment_document is set:
  use expense_documents.attachment and expense_documents.attachment_hash
else:
  the expense has no attachment
```

There is no resolver function. Each query joins `expense_documents ed ON e.attachment_document = ed.id` and returns:

- filename (`attachment`);
- SHA-256 hash (`attachment_hash`);
- collection id (`attachment_collection_id`, always the `expense_documents` id);
- record id (`attachment_record_id`, the document id);
- the storage path is built from those: `{collection id}/{record id}/{filename}`.

All code that displays, downloads, zips, exports, or writes back expense attachments uses this join. No code reads `e.attachment` or `e.attachment_hash` any more.

These paths use it:

- expense details API and details page;
- expense list and tracking responses when attachment data is shown;
- receipt zip query and zip creation;
- legacy expenses writeback/export;
- any expense report SQL that selects `e.attachment` or `e.attachment_hash`;
- generated TypeScript types after the schema changes.

### 4. Add an authorized expense attachment download route

Because document files are not stored under the expense record, Turbo does not expose document files through collection rules.

The route is:

```text
GET /api/expenses/attachment/{id}
```

Behavior:

1. Require auth.
2. Load the expense through the same visibility policy used by expense details.
3. Resolve the effective attachment.
4. Stream the underlying PocketBase file only if the caller can view the expense.

The details page and list links use this expense-scoped route rather than constructing document collection file URLs directly.

Server-side report generation reads files internally through PocketBase filesystem APIs, using the storage path.

### 5. Add transactional custom expense write endpoints

The editor saves expenses through app-owned write endpoints rather than the generic PocketBase collection write endpoints:

```text
POST  /api/expenses
PATCH /api/expenses/{id}
```

This is intentional because an expense attachment upload now mutates two collections:

```text
expense_documents
expenses
```

The custom endpoints are the canonical Phase 1 write path for expense create/update from the UI. They accept the same logical form payload the editor already produces, including file uploads, duplicate-upload reuse, and `source_expense` reuse.

The endpoint implementation:

1. requires an authenticated `users` record;
2. loads or creates the target `expenses` record;
3. loads request data into the normal PocketBase `RecordUpsert` form so collection field validation and normal model save hooks still run;
4. calls the existing expense processing code for ownership, PO, branch, approver, currency, cumulative total, bookkeeper-on-behalf, attachment-required, and submitted-expense validation;
5. resolves the effective attachment document before saving;
6. saves the `expense_documents` row and the `expenses` row inside one `app.RunInTransaction` call;
7. if the transaction fails after a new document file was uploaded, rolls back the DB changes and best-effort deletes the uploaded document object from PocketBase storage/S3.

The generic PocketBase `expenses` create/update request hooks also wrap their existing `ProcessExpense` + `e.Next()` flow in a transaction and use the same cleanup tracker. This keeps direct collection API writes safe during the transition, but the UI uses only the custom endpoints.

Keep the visible form unchanged for ordinary users:

- same fields;
- same attachment control;
- same validation messages where possible;
- same redirect after save.

All existing backend validation remains centralized in the current expense processing path. The custom endpoint owns the transaction boundary and request parsing, not a second copy of PO, branch, approver, bookkeeper-on-behalf, currency, and attachment-required rules.

An upload in the `attachment` request field is not a legacy write. `attachment` is now only a request field on expenses, not a column. Whether the upload arrives through the custom endpoint or the generic collection API, the backend converts it to an `expense_documents` link before persistence.

The editor trims its save payload to editable expense fields before calling the custom endpoint. Existing document-backed filenames shown in the attachment control are not resent as `attachment` values. Only a new `File`, an explicit empty attachment removal, or `source_expense` is sent as attachment intent.

### 6. Implement document resolution during create/update

The backend helper used by the expense create/update hooks (`app/hooks/expense_documents.go`) is:

```text
resolveExpenseDocumentForSave(app, e, attachmentHash, hasBookKeeperClaim, poBacked)
```

It handles four cases.

#### Case A: New file, hash not seen before

1. Compute SHA-256 from the uploaded file.
2. Create a new `expense_documents` record.
3. Store the uploaded file on `expense_documents.attachment`.
4. Set `expense_documents.attachment_hash`.
5. Set `expense_documents.uploaded_by` to the authenticated user.
6. Set `expenses.attachment_document` to the new document id.

#### Case B: New file, hash already exists, caller is not `book_keeper`

Reject with the existing duplicate attachment semantics when the hash is already
used by another expense. "Already used" means that another expense references the
matching `expense_documents` row through `expenses.attachment_document`. (Phase 1
also checked legacy-only expenses through `expenses.attachment_hash`; that field
is gone.)

```text
field: attachment
code: duplicate_file
message: This file has already been uploaded to another expense
```

If the matching `expense_documents` row exists but no expense references it,
the upload relinks that orphaned document row. Turbo does not create a second
`expense_documents` row for the same hash. This preserves the duplicate-file UX for everyone without the
`book_keeper` claim while avoiding a stale orphaned document row blocking a
valid upload.

#### Case C: New file, hash already exists, caller has `book_keeper`

Allow reuse only when the target expense is PO-backed:

```text
expenses.purchase_order != ''
```

Then:

1. Do not create a duplicate `expense_documents` row.
2. Set `expenses.attachment_document` to the existing document id.

If the target expense has no PO, reject the reuse. The special workflow exists only for splitting a shared vendor document across multiple PO-backed expenses.

#### Case D: Existing source expense document

For the "Create another expense with this attachment" workflow, accept a `source_expense` id in the create request.

Allow this only when:

- caller has the `book_keeper` claim;
- target expense is PO-backed;
- source expense is visible to the caller (the caller is its `uid`, its `creator`, its approver once submitted, a `commit` holder once approved, or a `report` or `admin` holder once committed);
- source expense has an `attachment_document`.

Then set the target expense's `attachment_document` to the source document id.

There is no arbitrary document picking. The only reuse flows are duplicate upload and "create another expense with this attachment" from an existing visible expense.

### 7. Automatically migrate legacy attachments when edited

**Removed.** During Phases 1 and 2, saving an expense that still had a legacy `expenses.attachment` copied the file into `expense_documents`. Phase 3 removed this code along with the legacy columns, because no expense has a legacy attachment any more.

Edit behavior today:

- If the user uploads a replacement file, Turbo uses the new uploaded file and the normal document-resolution rules.
- If the user does not upload a replacement, the expense keeps its `attachment_document`. A client cannot set that relation directly.
- If the user removes an attachment and the expense type permits no attachment, Turbo clears the document relation. For a type that requires an attachment, the save fails with `attachment is required`.
- `Allowance`, `Mileage` and `PersonalReimbursement` expenses never keep a document.

### 8. Add bookkeeper-only source reuse UX

On the expense details page, callers with `book_keeper` see an action only when the source expense has an attachment. The action button is next to the attachment field (download) and its title reads:

```text
Create another expense with this attachment
```

The action opens PO search (`/pos/search?source_expense={expense_id}`). Choosing a PO there leads into the existing PO-backed create flow. The target expense still gets its PO-specific defaults from the selected PO. The source expense contributes only the attachment document.

Route shape of the create page:

```text
/expenses/add/{poid}?source_expense={expense_id}
```

The page loader keeps `source_expense` in page data, and the editor includes it in the create request. The backend remains authoritative and validates every condition listed in Case D.

### 9. Preserve editor record shape for edit pages

Expense edit pages load the normal PocketBase `expenses` record as the item to edit, not the custom details response:

```text
pb.collection("expenses").getOne(expense_id, {
  expand: "purchase_order,attachment_document"
})
```

This is intentional. The editor relies on PocketBase system fields and `expand.purchase_order` to determine ownership, PO payment type, PO currency, and recurring/cumulative PO hints. Loading the custom details response would drop that record shape and make some PO-backed or bookkeeper-created expenses behave like non-editable records.

The edit loader also calls `GET /api/expenses/details/{id}` and uses its `attachment` value as the display filename, while preserving the canonical `attachment_document` relation on the expense record. There is no fallback to `expenses.attachment`; that field is gone.

The editor's attachment link uses the expense-scoped download route:

```text
GET /api/expenses/attachment/{id}
```

The UI does not construct direct `expense_documents` file URLs. The route keeps authorization scoped to expense visibility.

Phase 3 removed the editor's fallback to `expenses.attachment`. The editor relies only on `attachment_document`.

### 10. Known issues and caveats

Phase 1 was a mixed-storage rollout. The app wrote new expense attachments to `expense_documents`, while legacy rows pointed at files under the `expenses` collection until the Phase 2 backfill. Code that read expense attachments then used this pattern, which Phase 3 removed:

```text
COALESCE(expense_documents.attachment, expenses.attachment)
COALESCE(expense_documents.attachment_hash, expenses.attachment_hash)
```

Known caveats. Mixed storage has ended, so each caveat below says whether it still holds:

- (History.) Rolling back to pre-Phase-1 code after new document-backed writes have landed is not data-neutral. Pre-Phase-1 code reads `expenses.attachment` directly, so newly created document-backed expenses with blank legacy attachment fields can lose attachment visibility until the Phase 1 code is restored or the data is repaired.
- (Still holds.) Old direct file URLs under `/api/files/expenses/...` no longer work, because `expenses` has no file field. UI links use `GET /api/expenses/attachment/{id}` so authorization stays expense-scoped.
- (Still holds.) Legacy writeback keeps the old payload shape. It exports the document's `attachment` and `attachmentHash` values, not `attachment_document`, so downstream legacy consumers do not need a schema change. A downstream process that fetches files by filename must look under the `expense_documents` storage prefix.
- (Still holds.) Receipt ZIP generation reads document-backed attachments. The ZIP cache manifest includes collection id, storage path, zip filename, and hash so a reused document or duplicate hash cannot produce a stale ZIP.
- (History.) The ZIP cache migration deleted existing `zip_cache` rows, so the first receipt ZIP requests after that deploy regenerated those archives.
- (Still holds.) A concurrent upload of the same new file can still race at the `expense_documents.attachment_hash` unique index. The expected steady-state behavior is still one document row per hash, but the losing request may surface a lower-level save error rather than the normal `duplicate_file` validation message.
- (History.) Phase 1 did not delete old legacy expense attachment objects.
- (Still holds.) Expense deletion may leave an unreferenced `expense_documents` row. A later upload of the same file by any user relinks that orphaned document, as long as no other expense references it.
- (Removed.) Editing or reusing a legacy-only expense used to lazily create or reuse an `expense_documents` row. That code is gone; a source expense must now have an `attachment_document`.
- (Still holds.) The generic PocketBase `expenses` create/update hooks remain active and convert direct uploads to `expense_documents`. The custom `/api/expenses` routes are the UI path, but direct collection API writes should not silently create new legacy-only attachments.
- (Still holds.) `image/heic` is accepted by backend validation, but browser preview/download behavior may be less polished than PDF, PNG, or JPEG. The attachment route should still stream the original file.

### 11. Phase 1 validation before deploy

This is the checklist that was written for the Phase 1 deploy. It is kept as a record. The steps about legacy attachments no longer apply, because no expense has one.

Test on `multi_use_attachment` with a dumped production DB.

Minimum manual checks:

1. Open old expenses with legacy attachments. Download links still work.
2. Create a new expense with a new attachment. Verify the S3 key uses the `expense_documents` collection id, not the `expenses` collection id.
3. Create a normal non-bookkeeper duplicate. Verify it is rejected.
4. Create a bookkeeper PO-backed duplicate. Verify it reuses the existing `expense_documents` row.
5. Use "Create another expense with this attachment" as a bookkeeper. Verify the new PO expense shares the same document.
6. Edit an old legacy expense without changing the attachment. Verify it gets an `attachment_document` and the copied document file downloads.
7. Generate receipt zips and legacy writeback/export output for a mix of legacy and document-backed expenses.

After deploying Phase 1 to production, run a short smoke test before treating the
release as stable. The goal is not to re-test every backend branch; it is to
prove that the production app, database schema, file storage, auth rules, and
integration read paths all agree on the mixed-storage model.

Production smoke test:

1. Confirm the migration ran:
   - `expense_documents` exists;
   - `expenses.attachment_document` exists;
   - `zip_cache` has the `manifest` field instead of `hashes` and `filenames`.
2. Open a known old expense that still has a legacy `expenses.attachment` value and no `attachment_document`.
   - The details page should show the attachment.
   - `GET /api/expenses/attachment/{id}` should download or preview the file.
   - The file should still resolve from the old `expenses` storage prefix.
3. Create a normal user expense with a new PDF, PNG, or JPEG attachment.
   - The expense should save through `/api/expenses`.
   - The expense row should have `attachment_document` set.
   - The expense row should have blank legacy `attachment` and `attachment_hash` fields.
   - The linked `expense_documents` row should have the filename, hash, and `uploaded_by`.
   - The stored object should live under the `expense_documents` collection prefix.
   - The details/list attachment link should download through `/api/expenses/attachment/{id}`.
4. Attempt a duplicate upload as a normal non-`book_keeper` user.
   - If another expense still references the matching document, the request should fail with the existing `attachment.duplicate_file` style error.
   - If the matching document row is orphaned and no legacy-only expense has the same hash, the request should save and reference the existing document row.
   - If the matching document row is orphaned but a legacy-only expense still has the same hash, the request should fail with the existing `attachment.duplicate_file` style error.
   - No new `expense_documents` row should be created.
5. Create a PO-backed expense as a `book_keeper` using a duplicate attachment.
   - The new expense should save.
   - It should reference the existing `expense_documents` row.
   - It should not create a second document row with the same hash.
6. From an expense details page as a `book_keeper`, use "Create another expense with this attachment".
   - The PO search flow should preserve `source_expense`.
   - The created PO-backed expense should share the same `attachment_document`.
   - The target expense should still inherit its normal PO-specific defaults.
7. Edit a legacy-only expense without changing its file.
   - The edit should save.
   - The expense should gain an `attachment_document`.
   - The old legacy filename should remain available for Phase 1 fallback.
   - The copied document-backed attachment should download successfully.
8. Generate a receipt ZIP for a date/week that includes both legacy-only and document-backed expenses.
   - The ZIP should include both files.
   - Re-running the same report should hit the regenerated `zip_cache` record.
   - Changing the attachment set later should miss the cache and regenerate the ZIP.
9. Run the legacy expense export/writeback path for rows created or edited after Phase 1.
   - The payload should still expose effective `attachment` and `attachmentHash` values.
   - Downstream consumers should not need to know about `attachment_document`.
10. Monitor the first production users for errors around:
    - `POST /api/expenses`;
    - `PATCH /api/expenses/{id}`;
    - `GET /api/expenses/attachment/{id}`;
    - weekly receipt ZIP generation;
    - legacy expense export/writeback.

After the smoke test passes, keep Phase 1 running in mixed-storage mode for
several normal business days before starting Phase 2. During that soak period,
prefer roll-forward fixes for any issue that appears. Avoid reverting to
pre-Phase-1 code unless there is a severe incident and a data repair plan has
already accounted for document-backed expenses with blank legacy attachment
fields and reused documents that cannot cleanly map back to the old unique
`expenses.attachment_hash` invariant.

Before starting Phase 2, produce a small production baseline report:

- count of legacy-only expenses with `attachment != ''` and no `attachment_document`;
- count of document-backed expenses;
- count of `expense_documents` rows referenced by more than one expense;
- count of document-backed expenses whose legacy `attachment` field is blank;
- any attachment download, receipt ZIP, or writeback errors observed during the soak period.

### 12. Phase 1 backend tests

Backend tests cover Phase 1 (for example `TestExpenseDocumentsPhase1CreateAndReuse` and `TestExpenseDocumentsDocumentOnlyIntegration` in `app/expenses_test.go`). The list below is the original plan. Its legacy-fallback and auto-migration cases no longer apply, because Phase 3 removed that code.

Expense document creation and linking:

- creating an expense with a new attachment creates exactly one `expense_documents` record;
- the created expense has `attachment_document` set;
- the created expense does not write a new legacy `expenses.attachment` value;
- the stored document hash matches the uploaded file content;
- the document `uploaded_by` field is the authenticated user.

Duplicate hash behavior:

- non-`book_keeper` uploading a file whose hash already exists in `expense_documents` receives the existing field-level duplicate attachment error when another expense references that document;
- non-`book_keeper` uploading a file whose hash already exists only on an orphaned `expense_documents` row succeeds and references that existing document;
- non-`book_keeper` uploading a file whose hash exists on an orphaned `expense_documents` row and on a legacy-only expense is still rejected;
- non-`book_keeper` uploading a file whose hash exists only on a legacy expense is still rejected;
- `book_keeper` uploading a duplicate file for a PO-backed expense links the existing document and creates no duplicate document;
- `book_keeper` uploading a duplicate file for a non-PO expense is rejected.

Source expense reuse:

- `book_keeper` can create a PO-backed expense using `source_expense` when the source has an `attachment_document`;
- `book_keeper` can create a PO-backed expense using `source_expense` when the source is legacy-only, and the source is migrated to `expense_documents` first;
- non-`book_keeper` using `source_expense` is rejected;
- `source_expense` is rejected when the source expense is not visible to the caller;
- `source_expense` is rejected when the source expense has no effective attachment;
- `source_expense` is rejected when the target expense has no PO.

Legacy fallback and auto-migration:

- expense details return a downloadable attachment for legacy-only expenses;
- receipt zip generation includes both legacy-only and document-backed expense attachments;
- legacy writeback/export includes effective attachment filename and hash for both legacy-only and document-backed expenses;
- editing a legacy-only expense without replacing its file creates or reuses an `expense_documents` row and links it;
- editing a legacy-only expense with a replacement file uses the replacement document and does not migrate the old file;
- removing an attachment from an expense type that permits no attachment clears the effective attachment relation/fields.

Regression coverage:

- existing PO-backed expense validation still enforces active PO, matching job, matching currency, and cumulative overflow behavior;
- existing bookkeeper-on-behalf validation still requires `book_keeper`, a PO-backed target, valid payment type, and PO owner matching;
- submitted expenses still cannot be edited through the generic expense save path;
- direct generic collection file upload to `expenses.attachment` is converted to `expense_documents` and does not persist a new legacy expense attachment;
- custom `POST /api/expenses` creates a document-backed expense;
- custom `PATCH /api/expenses/{id}` updates through the same document-first workflow;
- a failed final expense save after document creation rolls back the `expense_documents` row so the unique document hash cannot block a later valid upload.

## Phase 2: Backfill Existing Legacy Attachments

Phase 2 is over. The backfill command described here was built, and Phase 3 then removed it from the repo together with the legacy columns, so it can no longer be run. The code does not record how the production run went. This section is kept as a record of how the backfill was designed.

The backfill should be idempotent and split into prepare, manual S3 copy, verify, and apply steps. This prevents the app from pointing an expense at a document path before the file exists in S3.

Phase 2 is a local-first maintenance workflow:

1. Dump the production database locally.
2. Run `prepare` against the local dump to precompute the migration manifest.
3. Run the generated `copy_s3.sh` to copy legacy S3 objects into their new document-owned prefixes.
4. Run `verify` against the copied objects.
5. Stop production.
6. Re-dump or otherwise refresh the production database state locally.
7. Run `apply` from the verified manifest so the database rows match the already-copied files.

The target set is committed legacy-only expenses:

```sql
SELECT *
FROM expenses
WHERE attachment != ''
  AND (attachment_document = '' OR attachment_document IS NULL)
  AND committed != ''
```

The copy step must copy, never move or delete, legacy S3 objects. Any reconcile mismatch during `verify` or `apply` should abort rather than inventing a new mapping.

### 1. Add a backfill command

The command lived under `app/cmd/backfill_expense_documents`. It is no longer in the repo.

It should support at least these modes:

```text
prepare
verify
apply
report
```

The command uses the app database and PocketBase collection metadata so it can discover:

- `expenses` collection id;
- `expense_documents` collection id;
- file names;
- record ids;
- hashes.

### 2. Prepare mode

`prepare` scans for committed expenses that still need migration:

```sql
SELECT *
FROM expenses
WHERE attachment != ''
  AND (attachment_document = '' OR attachment_document IS NULL)
  AND committed != ''
```

For each row:

1. Determine the legacy S3 key:

   ```text
   {expenses_collection_id}/{expense_id}/{expenses.attachment}
   ```

2. Determine the attachment hash:

   - use `expenses.attachment_hash` when present;
   - if blank, read the legacy file and compute SHA-256;
   - if the file cannot be read, record an error and skip that expense.

3. Determine the document id and document attachment filename for the attachment hash:

   - if the hash already exists in `expense_documents`, use that document id and the row's current `attachment` filename;
   - if the hash already appeared earlier in the manifest, reuse the same precomputed document id and document attachment filename;
   - otherwise generate a new deterministic `expense_document_id` for the manifest and use the legacy expense attachment filename as `document_attachment`.

   Any generated `expense_document_id` must be deterministically derived from
   the attachment hash and must be a valid PocketBase record id: exactly 15
   lowercase alphanumeric characters. Use the first 15 characters of the
   lowercase base32 encoding of:

   ```text
   SHA-256("expense_document:" + attachment_hash)
   ```

   Do not use the raw attachment hash prefix directly as the document id.

4. Determine the future `expense_documents` row values for newly generated document ids:

   - `attachment` set to the same filename as the legacy attachment;
   - `attachment_hash` set to the computed hash;
   - `uploaded_by` set to `expenses.creator` when present, otherwise `expenses.uid`.

   When the manifest reuses an existing `expense_documents` row, do not rewrite
   that row's attachment filename. The manifest's `new_s3_key` must use the
   existing document id and existing document attachment filename.

5. Do not create `expense_documents` rows and do not set `expenses.attachment_document` in `prepare`. The database writes happen in `apply`, after S3 copies are complete and production is stopped.

6. Write a manifest row containing:

   ```text
   expense_id
   expense_attachment
   attachment_hash
   expense_document_id
   document_attachment
   uploaded_by
   old_s3_key
   new_s3_key
   copy_required
   status
   ```

7. Write an executable copy script derived from the manifest.

Recommended output paths:

```text
tmp/expense_document_backfill/manifest.tsv
tmp/expense_document_backfill/copy_s3.sh
tmp/expense_document_backfill/errors.tsv
```

`prepare` must be safe to rerun. For the same input database and files, it should preserve the same `expense_document_id` values so the generated S3 destination keys stay stable.

### 3. Manual S3 copy

After `prepare`, manually copy the existing S3 objects from the old expense prefixes to the new expense document prefixes.

This manual copy happens before the database is mutated. The manifest's precomputed `expense_document_id` values determine the destination S3 keys that `apply` will later reference.

For each manifest row with `copy_required = true`, copy:

```text
from: {expenses_collection_id}/{expense_id}/{expenses.attachment}
to:   {expense_documents_collection_id}/{expense_document_id}/{expense_documents.attachment}
```

Example key shape:

```text
from: o1vpz1mm7qsfoyy/b4o6xph4ngwx4nw/receipt.pdf
to:   {expense_documents_collection_id}/{expense_document_id}/receipt.pdf
```

The generated `copy_s3.sh` should contain explicit copy commands so the operator can review exactly what will move before running it. It must not delete the old objects.

Example command shape:

```bash
aws s3 cp "s3://$TYBALT_S3_BUCKET/o1vpz1mm7qsfoyy/{expense_id}/{filename}" \
  "s3://$TYBALT_S3_BUCKET/{expense_documents_collection_id}/{expense_document_id}/{filename}"
```

Do not create `expense_documents` rows or set `expenses.attachment_document` for rows whose new S3 object has not been copied and verified.

### 4. Verify mode

After the manual S3 copy, run `verify`.

`verify` reads the manifest and checks every target object:

1. New S3 object exists.
2. New S3 object can be read through PocketBase filesystem APIs.
3. SHA-256 of the copied object matches the manifest `attachment_hash`.
4. Existing legacy source object still exists.
5. The current expense row is still legacy-only or already linked to the expected document.
6. No existing `expense_documents` row maps the same hash to a different id.
7. If the manifest reuses an existing `expense_documents` row, that row's `attachment_hash` and `attachment` still match the manifest.

Rows that fail verification remain unapplied. Write failures to:

```text
tmp/expense_document_backfill/verify_errors.tsv
```

`verify` must be safe to rerun after recopying failed rows.

### 5. Apply mode

After `verify` succeeds and production is stopped, `apply` reconciles the refreshed database against the manifest and performs the database writes. It inserts any missing `expense_documents` rows using the manifest's precomputed ids, then updates each verified expense:

```text
expenses.attachment_document = expense_document_id
```

Leave these fields in place during the first backfill pass:

```text
expenses.attachment
expenses.attachment_hash
```

Keeping the legacy fields until cleanup gives a rollback path and makes it easier to audit the migration.

`apply` should run all DB changes in one transaction and be idempotent:

- skip rows already linked to the expected document;
- abort if a row is linked to a different document;
- abort if an expected expense is missing;
- abort if an `expense_documents` row already uses the same hash with a different id;
- abort if a copied S3 target is missing or fails hash verification;
- never overwrite a non-empty `attachment_document`.

### 6. Report mode

`report` summarizes migration progress:

- total expenses with legacy attachment;
- total with `attachment_document`;
- total still legacy-only committed expenses;
- total document-backed expenses whose legacy `attachment` field is blank;
- total document-backed but missing target S3 object;
- duplicate hashes sharing one document;
- rows with missing legacy files;
- rows with blank or invalid hash.

This report is the acceptance checkpoint before cleanup.

### 7. Cleanup after successful backfill

This cleanup is done. It was Phase 3:

1. Remove all remaining code paths that rely on legacy fallback.
2. Remove `expenses.attachment`.
3. Remove `expenses.attachment_hash`.
4. Remove the old unique index on `expenses.attachment_hash`.
5. Update generated PocketBase types.
6. Update imports/exports/writeback queries so `expense_documents` is the only source for expense attachments.
7. Remove any UI handling that displays legacy expense attachment links.

`attachment_missing_reason` is the permanent audit marker for historical
receipts that were supposed to exist but were already unrecoverable in the
legacy system. It is set only by an admin-only repair endpoint
(`POST /api/expenses/{id}/attachment_missing/mark`, `admin` claim) that clears
the document relation in a guarded transaction after the operator has verified
the receipt cannot be recovered. Once set, the reason cannot be changed. These
rows have no `expense_documents` row, and the expense details page displays the
missing reason instead of treating the expense as an ordinary attachmentless row.

Cleanup shipped as its own change after the backfill.

### 8. Phase 2 backend tests

These tests covered the backfill command. They were removed from the repo with it.

Prepare mode:

- `prepare` precomputes document ids for legacy-only committed expenses with valid attachments;
- `prepare` reuses an existing or already-manifested document id when another expense has the same hash;
- `prepare` uses `expenses.attachment_hash` when present;
- `prepare` computes the SHA-256 hash from the legacy file when `expenses.attachment_hash` is blank;
- `prepare` records an error and skips rows whose legacy file cannot be read;
- `prepare` writes a manifest with correct old and new S3 keys;
- rerunning `prepare` is idempotent and preserves deterministic 15-character lowercase alphanumeric manifest document ids derived from the attachment hash;
- generated manifest document ids use the namespaced SHA-256/base32 derivation, not a raw attachment hash prefix.

Manual copy support:

- generated copy commands use the legacy expense collection prefix as the source;
- generated copy commands use the `expense_documents` collection prefix and document id as the destination;
- rows that already have a verified destination object are marked as not requiring copy;
- generated scripts never delete legacy objects.

Verify mode:

- `verify` succeeds when the copied destination object exists and its SHA-256 matches the document hash;
- `verify` fails when the destination object is missing;
- `verify` fails when the copied destination object hash does not match;
- `verify` fails when the legacy source object is missing;
- rerunning `verify` after correcting a failed copy succeeds without changing linked expenses.

Apply mode:

- `apply` inserts missing `expense_documents` rows using manifest ids;
- `apply` sets `expenses.attachment_document` only for verified rows;
- `apply` leaves `expenses.attachment` and `expenses.attachment_hash` untouched during the backfill pass;
- `apply` skips rows already linked to the expected document;
- `apply` aborts when a row is linked to a different document;
- `apply` aborts when a matching hash exists under a different document id;
- `apply` is transactional for DB updates and leaves no partially applied batch after an injected failure.

Report and cleanup readiness:

- `report` counts legacy-only, document-backed, missing-copy, missing-source, and duplicate-hash rows correctly;
- cleanup tests verify the app no longer needs legacy fallback after all rows are linked;
- post-cleanup export/report tests use only `expense_documents` as the source of expense attachment data.

## Testing Requirements

The legacy-fallback, auto-migration and backfill items in these lists no longer apply, because Phase 3 removed that code.

Backend tests should cover:

- new expense with a new document;
- duplicate upload rejected for non-bookkeeper;
- duplicate upload reused for bookkeeper on PO-backed expense;
- duplicate upload rejected for bookkeeper on non-PO expense;
- source expense reuse accepted for bookkeeper;
- source expense reuse rejected for non-bookkeeper;
- source expense reuse rejected when source is not visible;
- source expense reuse rejected when source has no effective attachment;
- legacy fallback read for expense details and receipt generation;
- legacy edit auto-migrates to `expense_documents`;
- replacement attachment on a legacy expense uses the replacement document rather than migrating the old file;
- backfill prepare is idempotent;
- backfill verify refuses bad or missing copies;
- backfill apply skips already-linked rows and refuses mismatched links.
- admin repair can mark an unrecoverable historical attachment missing while
  preserving a text reason and clearing the document relation.

Frontend tests or manual checks should cover:

- ordinary users see the same expense form and attachment control;
- bookkeepers see the details-page action only when an attachment exists;
- source reuse lands on the existing PO-backed expense form;
- admins can audit/repair a bad hash or mark a confirmed unrecoverable
  historical attachment missing from the expense details dialog;
- validation errors still display next to the attachment field.

Integration/manual testing against the dumped production DB should include:

- mixed legacy and document-backed expenses in details/list/tracking;
- receipt zip generation for a mixed pay period;
- legacy writeback/export for a mixed updated-after range;
- editing old rows with legacy attachments;
- the exact generated S3 copy manifest for Phase 2.

## Operational Notes

Phases 1 to 3 are done. The notes about deploying a phase are history.

- Phase 1 requires a normal deployment but not an outage window.
- Phase 1 should be deployed only after testing locally on `multi_use_attachment` against a dumped production DB.
- Phase 2 does not need to be scheduled immediately.
- The manual S3 copy belongs to Phase 2, after `prepare` and before `verify`/`apply`.
- The old S3 objects under the `expenses` collection prefix should not be deleted during Phase 2. Deletion, if desired, should be considered only after a separate retention/rollback decision.
- The `expense_documents.attachment_hash` unique index remains the long-term duplicate prevention mechanism.
- Orphaned `expense_documents` rows may be relinked by later uploads of the same file; this is not treated as multi-expense reuse unless another expense still references the document.
- `book_keeper` reuse changes only which expense may point at an existing document. It never creates two document records with the same hash.
- Historical required-attachment expenses whose original files are permanently
  lost should be marked with `attachment_missing_reason`; do not upload
  placeholders, invent replacement files, or blank them silently.

## Phase 3: cleanup

Phase 3 is done. Migration `app/migrations/1779720278_drop_expense_legacy_attachment_fields.go`
removed `expenses.attachment`, `expenses.attachment_hash` and the unique index on
that hash, and the legacy fallback code and the backfill command are gone. A
known unrecoverable historical attachment is represented by
`attachment_missing_reason`.

## Vendor invoice details (#153)

Status: implemented on the `vendor_invoices` branch.

A vendor invoice is the document that substantiates an expense, whoever enters
the expense. Its details are stored on the `expense_documents` row, not on the
expense, and not in a separate vendor invoice record. One invoice split across
several POs is one document shared by several expenses, as bookkeepers already
do.

### Fields on `expense_documents`

| Field                       | Type                    | Notes                                                                                                                          |
| --------------------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `vendor`                    | relation → `vendors`    | Copied from the expense when the details are first written.                                                                    |
| `vendor_invoice_number`     | text, max 100           | As entered.                                                                                                                    |
| `vendor_invoice_number_key` | text                    | Server-set: the number in lowercase with every character except letters and digits removed. Used only for the duplicate check. |
| `vendor_invoice_date`       | text, `YYYY-MM-DD`      | The date on the invoice.                                                                                                       |
| `vendor_invoice_total`      | number                  | The invoice total, taxes included, in the invoice currency.                                                                    |
| `vendor_invoice_currency`   | relation → `currencies` | Copied from the expense. Empty means the home currency, as on expenses.                                                        |

A document with a `vendor_invoice_number` is a vendor invoice. One without is a
receipt. A unique partial index on `(vendor, vendor_invoice_number_key)` where
the key is not empty prevents the same invoice from being recorded twice.

The collection's API rules stay null. Like the file itself, the details are
written only by the expense write path.

### Entering details

The expense write endpoints (`POST /api/expenses`, `PATCH /api/expenses/{id}`)
accept three virtual inputs alongside `attachment`:

- `vendor_invoice_number`
- `vendor_invoice_date`
- `vendor_invoice_total`

They are not expense fields. After the document is resolved, `ProcessExpense`
applies them to the expense's document in the same transaction as the expense
save. An input left out of the request keeps the document's current value, so a
bookkeeper reusing a document carries its details.

### Rules

- **Required:** on an `OnAccount` expense that has a document, the number, date,
  and a total greater than zero must be present after the inputs are applied.
  This includes `OnAccount` expenses without a PO. Other payment types may record
  them; a `CorporateCreditCard` expense needs only its receipt. Historical rows
  with `attachment_missing_reason` and no document are exempt.
- **No document:** details sent for an expense with no document are refused:
  attach the invoice first.
- **Same vendor and currency:** once a document has a number, every expense
  linked to it must have the document's vendor and currency.
- **Who can change them:** the document's uploader (`uploaded_by`), a
  `book_keeper`, or anyone when no other expense uses the document.
- **When they lock:** once another expense on the document is submitted and not
  rejected, the details cannot change, because that expense is waiting for
  approval against them. The expense being saved is always unsubmitted, since
  submitted expenses cannot be edited.
- **Duplicates:**
  - The same file is caught by `attachment_hash`, as before.
  - A different file of the same invoice is caught by the unique key. The error
    gives the PO number and the name of whoever entered an expense on the
    existing invoice, not a link, because the user may not be allowed to view
    that expense.
  - The same invoice reused for more than it is worth is caught by the total:
    the totals of the expenses on a document, other than rejected ones, cannot
    exceed `vendor_invoice_total`.
- **Cents:** totals are compared in whole cents with
  `utilities.CurrencyAmountToCents`, so a sum like 16.26 + 0.01 matches a 16.27
  invoice exactly.
- **No backfill:** details start with expenses entered after release.

### Vendor merges

`expense_documents.vendor` is merged with the other vendor references. If both
vendors have an invoice with the same number key, the merge is refused with a
message naming the invoice number, since merging would record it twice.

### Payment (#148)

#148 adds approval to pay and the payment record to the same document, once per
invoice, after every expense on it is approved and the expenses add up to
`vendor_invoice_total`.
