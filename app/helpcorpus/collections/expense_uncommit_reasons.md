<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_uncommit_reasons

`expense_uncommit_reasons` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_uncommit_reasons/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_uncommit_reasons

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_uncommit_reasons

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can create a record (create rule): expense_uncommit_reasons

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): expense_uncommit_reasons

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): expense_uncommit_reasons

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: expense_uncommit_reasons

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `name` | text | required |
| `default_outcome` | select | required; values ["uncommit_and_reject","uncommit_only"]; maxSelect 1 |
| `note_prompt` | text |  |
| `note_required` | bool |  |
| `sort` | number |  |
| `active` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
