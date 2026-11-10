<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: profiles

`profiles` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/profiles/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): profiles

```
@request.auth.id != "" && (
  uid = @request.auth.id || 
  manager = @request.auth.id || 
  @request.auth.user_claims_via_uid.cid.name ?= 'tame' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'job' ||
  (
    @request.auth.user_claims_via_uid.cid.name ?= 'legacy_po_create_update' &&
    uid.admin_profiles_via_uid.active = true
  )
)
```

## Who can view one record (view rule): profiles

```
@request.auth.id != ""
```

## Who can create a record (create rule): profiles

```
@request.auth.id != "" &&
uid = @request.auth.id
```

## Who can update a record (update rule): profiles

```
@request.auth.id != "" &&
uid = @request.auth.id &&
@request.body.uid:changed = false
```

## Who can delete a record (delete rule): profiles

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: profiles

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `surname` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:[-'][a-zA-Z]+)*$" |
| `given_name` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:-[a-zA-Z]+)*$" |
| `manager` | relation | required; maxSelect 1; to `users` |
| `alternate_manager` | relation | maxSelect 1; to `users` |
| `default_division` | relation | maxSelect 1; to `divisions` |
| `uid` | relation | maxSelect 1; to `users` |
| `notification_type` | select | values ["email_text","email_html"]; maxSelect 1 |
| `do_not_accept_submissions` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `default_expense_payment_type` | select | values ["OnAccount","Expense","CorporateCreditCard","Allowance","FuelCard","Mileage","PersonalReimbursement"]; maxSelect 1 |
| `default_role` | relation | maxSelect 1; to `rate_roles` |

## Indexes: profiles

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_dvV9kj4` ON `profiles` (`uid`)`
- `CREATE INDEX `idx_iD56IM8IJg` ON `profiles` (`manager`)`
