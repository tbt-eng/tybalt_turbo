<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: work_records

`work_records` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/work_records/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): work_records

```
@request.auth.id != "" &&
(
  creator = @request.auth.id ||
  work_records_subjects_via_work_record.uid ?= @request.auth.id ||
  @request.auth.user_claims_via_uid.cid.name ?= 'report'
)
```

## Who can view one record (view rule): work_records

```
@request.auth.id != "" &&
(
  creator = @request.auth.id ||
  work_records_subjects_via_work_record.uid ?= @request.auth.id ||
  @request.auth.user_claims_via_uid.cid.name ?= 'report'
)
```

## Who can create a record (create rule): work_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): work_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): work_records

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: work_records

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `number` | text | pattern "^W[0-9]{4}-[0-9]{4}$" |
| `parent_work_record` | relation | maxSelect 1; to `work_records` |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `creator` | relation | required; maxSelect 1; to `users` |
| `job` | relation | required; maxSelect 1; to `jobs` |
| `type` | relation | required; maxSelect 1; to `work_records_types` |
| `location` | text | required |
| `work_description` | text | required |
| `report_to` | text |  |
| `field_book_number` | text | pattern "^[0-9]{2}-(00[1-9]\|0[1-9][0-9]\|[1-9][0-9]{2})$" |
| `field_book_page_number` | number | min 1; onlyInt |
| `sub_contractor` | text |  |
| `equipment` | text |  |
| `supplies` | text |  |
| `lane_kms` | number |  |
| `attachment` | file | maxSelect 1; maxSize 10485760; mimeTypes ["application/pdf"] |
| `attachment_hash` | text |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: work_records

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_work_records_number ON work_records (number)`
- `CREATE INDEX idx_work_records_job_date ON work_records (job, date)`
- `CREATE UNIQUE INDEX `idx_zdRTW9LEgq` ON `work_records` (`attachment_hash`) WHERE `attachment_hash` != ''`
