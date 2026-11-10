<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: user_claims_summary

`user_claims_summary` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/user_claims_summary/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): user_claims_summary

```
@request.auth.id = id
```

## Who can view one record (view rule): user_claims_summary

```
@request.auth.id = id
```

## The query of the view: user_claims_summary

```sql
SELECT
  u.id AS id,
  COALESCE(
    (
      SELECT json_group_array(c.name)
      FROM user_claims uc
      JOIN claims c ON c.id = uc.cid
      WHERE uc.uid = u.id
    ),
    '[]'
  ) AS claims
FROM users u
```

## Fields: user_claims_summary

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `claims` | json | maxSize 1 |
