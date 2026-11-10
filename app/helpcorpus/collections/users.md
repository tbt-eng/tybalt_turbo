<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: users

`users` holds the accounts that people sign in with. Turbo's screens reach it through PocketBase's API (`/api/collections/users/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): users

```
id = @request.auth.id
```

## Who can view one record (view rule): users

```
// The caller is logged in and either
@request.auth.id != "" &&
(
  // 1. The caller id is equal to the record id or
  @request.auth.id = id ||
  // 2. The record has the po_approver claim or
  user_claims_via_uid.cid.name ?= 'po_approver' ||
  // 3. The record has the tapr claim or
  user_claims_via_uid.cid.name ?= 'tapr' ||
  // 4. The caller is the approver of at least one of users committed expenses or
  user_view_grants_via_uid.viewer ?= @request.auth.id ||
  // 5. The record is the uid, approver, or second_approver of 1+ active purchase order or
  user_view_grants_via_uid.reason ?= 'active_po' ||
  // 6. The caller has the tapr claim. Keep this last: PocketBase joins tables
  // in the order the rule names them, and each claim the caller holds repeats
  // every later join.
  @request.auth.user_claims_via_uid.cid.name ?= 'tapr'
)
```

## Who can create a record (create rule): users

An empty rule. Everybody can do this, also a person who is not signed in.

## Who can update a record (update rule): users

```
id = @request.auth.id
```

## Who can delete a record (delete rule): users

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: users

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `password` | password | required; min 8; hidden |
| `tokenKey` | text | required; min 30; max 60; hidden |
| `email` | email |  |
| `emailVisibility` | bool |  |
| `verified` | bool |  |
| `username` | text | required; min 3; max 150; pattern "^[\\w][\\w\\.\\-]*$" |
| `name` | text | max 255 |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `avatar` | file | maxSelect 1; mimeTypes ["image/jpeg","image/png","image/svg+xml","image/gif","image/webp"] |

## Indexes: users

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `__pb_users_auth__username_idx` ON `users` (username COLLATE NOCASE)`
- `CREATE UNIQUE INDEX `__pb_users_auth__email_idx` ON `users` (`email`) WHERE `email` != ''`
- `CREATE UNIQUE INDEX `__pb_users_auth__tokenKey_idx` ON `users` (`tokenKey`)`
