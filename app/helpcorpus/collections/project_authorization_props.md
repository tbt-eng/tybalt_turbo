<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: project_authorization_props

`project_authorization_props` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/project_authorization_props/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): project_authorization_props

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can view one record (view rule): project_authorization_props

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can create a record (create rule): project_authorization_props

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can update a record (update rule): project_authorization_props

```
@request.auth.id != '' && @request.auth.admin_profiles_via_uid.active = true && @request.auth.user_claims_via_uid.cid.name ?= 'admin'
```

## Who can delete a record (delete rule): project_authorization_props

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: project_authorization_props

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `user_claim` | relation | required; maxSelect 1; to `user_claims`; cascadeDelete |
| `max_value` | number | required; min 1; max 9007199254740991; onlyInt |

## Indexes: project_authorization_props

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX idx_pa_props_claim ON project_authorization_props (user_claim)`
