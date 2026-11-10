<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: purchase_order_groups

`purchase_order_groups` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/purchase_order_groups/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): purchase_order_groups

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): purchase_order_groups

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): purchase_order_groups

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): purchase_order_groups

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): purchase_order_groups

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: purchase_order_groups

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `creator` | relation | required; maxSelect 1; to `users` |
| `description` | text |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
