<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: notification_templates

`notification_templates` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/notification_templates/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): notification_templates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): notification_templates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): notification_templates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): notification_templates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): notification_templates

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: notification_templates

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `code` | text | required; min 3; pattern "^[a-z]+(?:_[a-z]+)*$" |
| `description` | text |  |
| `subject` | text |  |
| `text_email` | text |  |
| `html_email` | text |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: notification_templates

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_BHMSwRwcDq` ON `notification_templates` (`code`)`
