<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: work_records_subjects

`work_records_subjects` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/work_records_subjects/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): work_records_subjects

```
@request.auth.id != "" &&
(
  uid = @request.auth.id ||
  work_record.creator = @request.auth.id ||
  @request.auth.user_claims_via_uid.cid.name ?= 'report'
)
```

## Who can view one record (view rule): work_records_subjects

```
@request.auth.id != "" &&
(
  uid = @request.auth.id ||
  work_record.creator = @request.auth.id ||
  @request.auth.user_claims_via_uid.cid.name ?= 'report'
)
```

## Who can create a record (create rule): work_records_subjects

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): work_records_subjects

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): work_records_subjects

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: work_records_subjects

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `work_record` | relation | required; maxSelect 1; to `work_records`; cascadeDelete |
| `uid` | relation | required; maxSelect 1; to `users` |
| `hours_on_site` | number | required |
| `hours_travel_time` | number | required |
| `distance_travelled_km` | number |  |
| `vehicle_type` | select | values ["company","personal"]; maxSelect 1 |
| `is_passenger` | bool |  |
| `company_vehicle_unit_number` | text | pattern "^(00[1-9]\|0[1-9][0-9]\|[1-9][0-9]{0,2})$" |
| `approved` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: work_records_subjects

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_work_records_subjects_work_record_uid ON work_records_subjects (work_record, uid)`
- `CREATE INDEX idx_work_records_subjects_uid_created ON work_records_subjects (uid, created)`
