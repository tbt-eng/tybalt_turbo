<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: managers

`managers` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/managers/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): managers

```
@request.auth.id != ""
```

## Who can view one record (view rule): managers

```
@request.auth.id != ""
```

## The query of the view: managers

```sql
SELECT p.uid AS id, p.surname AS surname, p.given_name AS given_name 
FROM profiles p
INNER JOIN user_claims u ON p.uid = u.uid
INNER JOIN claims c ON u.cid = c.id
WHERE c.name = 'tapr'
```

## Fields: managers

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `surname` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:[-'][a-zA-Z]+)*$" |
| `given_name` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:-[a-zA-Z]+)*$" |
