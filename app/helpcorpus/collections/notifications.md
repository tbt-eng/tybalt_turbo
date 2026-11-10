<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: notifications

`notifications` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/notifications/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): notifications

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): notifications

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): notifications

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): notifications

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): notifications

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: notifications

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `recipient` | relation | required; maxSelect 1; to `users` |
| `template` | relation | required; maxSelect 1; to `notification_templates` |
| `status` | select | required; values ["pending","inflight","sent","error"]; maxSelect 1 |
| `status_updated` | date | required |
| `error` | text |  |
| `user` | relation | maxSelect 1; to `users` |
| `system_notification` | bool |  |
| `data` | json |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
