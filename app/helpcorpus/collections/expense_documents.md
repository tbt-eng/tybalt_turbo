<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_documents

`expense_documents` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_documents/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_documents

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_documents

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): expense_documents

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): expense_documents

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): expense_documents

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: expense_documents

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `attachment` | file | required; maxSelect 1; maxSize 5242880; mimeTypes ["application/pdf","image/jpeg","image/png","image/heic"] |
| `attachment_hash` | text | required; min 64; max 64; pattern "^[a-fA-F0-9]{64}$" |
| `uploaded_by` | relation | required; maxSelect 1; to `users` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: expense_documents

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_expense_documents_attachment_hash` ON `expense_documents` (`attachment_hash`) WHERE `attachment_hash` != ''`
