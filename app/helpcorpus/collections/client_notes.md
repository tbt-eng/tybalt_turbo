<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: client_notes

`client_notes` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/client_notes/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): client_notes

```
@request.auth.id != ''
```

## Who can view one record (view rule): client_notes

```
@request.auth.id != ''
```

## Who can create a record (create rule): client_notes

```
@request.auth.id != ''
```

## Who can update a record (update rule): client_notes

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): client_notes

```
@request.auth.id != "" &&
@request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Fields: client_notes

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `note` | text | required; min 10; max 1000 |
| `client` | relation | required; maxSelect 1; to `clients` |
| `job` | relation | maxSelect 1; to `jobs` |
| `job_not_applicable` | bool |  |
| `uid` | relation | required; maxSelect 1; to `users` |
| `job_status_changed_to` | select | values ["Cancelled","No Bid"]; maxSelect 1 |
| `_imported` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
