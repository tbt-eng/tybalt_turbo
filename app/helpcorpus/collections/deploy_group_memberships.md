<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: deploy_group_memberships

`deploy_group_memberships` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/deploy_group_memberships/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): deploy_group_memberships

```
@request.auth.id != '' && user = @request.auth.id
```

## Who can view one record (view rule): deploy_group_memberships

```
@request.auth.id != '' && user = @request.auth.id
```

## Who can create a record (create rule): deploy_group_memberships

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): deploy_group_memberships

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): deploy_group_memberships

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: deploy_group_memberships

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `deploy_group` | relation | required; maxSelect 1; to `deploy_groups`; cascadeDelete |
| `user` | relation | required; maxSelect 1; to `users`; cascadeDelete |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: deploy_group_memberships

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_deploy_group_memberships_user_group` ON `deploy_group_memberships` (`user`, `deploy_group`)`
