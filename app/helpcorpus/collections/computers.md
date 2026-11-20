<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: computers

`computers` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/computers/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): computers

```
@request.auth.id != '' && (
  @request.auth.user_claims_via_uid.cid.name ?= 'computers' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can view one record (view rule): computers

```
@request.auth.id != '' && (
  @request.auth.user_claims_via_uid.cid.name ?= 'computers' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
```

## Who can create a record (create rule): computers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): computers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): computers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: computers

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `slug` | text | required |
| `serial` | text | required |
| `mfg` | text | required |
| `model` | text |  |
| `computer_name` | text |  |
| `system_type` | number | onlyInt |
| `ram` | number | onlyInt |
| `os_version` | text |  |
| `os_sku` | number | onlyInt |
| `os_arch` | text |  |
| `boot_drive` | text |  |
| `boot_drive_fs` | text |  |
| `boot_drive_cap` | number | onlyInt |
| `boot_drive_free` | number | onlyInt |
| `radiator_version` | number | onlyInt |
| `network_config` | json |  |
| `user_given_name` | text |  |
| `user_surname` | text |  |
| `email` | text |  |
| `upn` | text |  |
| `user_source_anchor` | text |  |
| `last_user` | relation | maxSelect 1; to `users` |
| `last_login` | date |  |
| `assigned_to` | relation | maxSelect 1; to `users` |
| `assigned_time` | date |  |
| `retired` | date |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: computers

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_computers_slug` ON `computers` (`slug`)`
- `CREATE INDEX `idx_computers_computer_name` ON `computers` (`computer_name`)`
