<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: admin_profiles

`admin_profiles` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/admin_profiles/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): admin_profiles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): admin_profiles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): admin_profiles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): admin_profiles

```
@request.auth.id != "" &&
@request.body.id:changed = false &&
@request.body.uid:changed = false &&
@request.body.legacy_uid:changed = false &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can delete a record (delete rule): admin_profiles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: admin_profiles

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users` |
| `work_week_hours` | number | required; min 8; max 40 |
| `salary` | bool |  |
| `default_charge_out_rate` | number | required; min 50; max 1000 |
| `off_rotation_permitted` | bool |  |
| `skip_min_time_check` | select | required; values ["no","on_next_bundle","yes"]; maxSelect 1 |
| `opening_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `opening_op` | number | max 332 |
| `opening_ov` | number | max 200 |
| `payroll_id` | text | required; pattern "^(?:[1-9]\\d*\|CMS[0-9]{1,2})$" |
| `untracked_time_off` | bool |  |
| `time_sheet_expected` | bool |  |
| `allow_personal_reimbursement` | bool |  |
| `mobile_phone` | text | pattern "^\\+1 \\(\\d{3}\\) \\d{3}-\\d{4}$" |
| `job_title` | text |  |
| `personal_vehicle_insurance_expiry` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `default_branch` | relation | required; maxSelect 1; to `branches` |
| `legacy_uid` | text |  |
| `active` | bool |  |

## Indexes: admin_profiles

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_UpEVC7E` ON `admin_profiles` (`uid`)`
- `CREATE UNIQUE INDEX `idx_XnQ4v11` ON `admin_profiles` (`payroll_id`)`
