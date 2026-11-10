<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: time_report_week_endings

`time_report_week_endings` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/time_report_week_endings/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): time_report_week_endings

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## Who can view one record (view rule): time_report_week_endings

```
@request.auth.user_claims_via_uid.cid.name ?= 'report'
```

## The query of the view: time_report_week_endings

```sql
SELECT
  MIN(id) as id,
  week_ending
FROM time_sheets
WHERE committed != ''
GROUP BY week_ending
ORDER BY week_ending DESC;
```

## Fields: time_report_week_endings

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `week_ending` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
