<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: user_po_approver_profile

`user_po_approver_profile` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/user_po_approver_profile/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): user_po_approver_profile

```
@request.auth.id = id
```

## Who can view one record (view rule): user_po_approver_profile

```
@request.auth.id = id
```

## The query of the view: user_po_approver_profile

```sql
WITH po_props AS (
  SELECT
    uc.uid,
    pap.max_amount,
    pap.project_max,
    pap.sponsorship_max,
    pap.staff_and_social_max,
    pap.media_and_event_max,
    pap.computer_max,
    pap.divisions
  FROM user_claims uc
  JOIN claims c ON c.id = uc.cid AND c.name = 'po_approver'
  LEFT JOIN po_approver_props pap ON pap.user_claim = uc.id
)
SELECT
  u.id AS id,
  COALESCE(pp.max_amount, 0) AS max_amount,
  COALESCE(pp.project_max, 0) AS project_max,
  COALESCE(pp.sponsorship_max, 0) AS sponsorship_max,
  COALESCE(pp.staff_and_social_max, 0) AS staff_and_social_max,
  COALESCE(pp.media_and_event_max, 0) AS media_and_event_max,
  COALESCE(pp.computer_max, 0) AS computer_max,
  COALESCE(pp.divisions, '[]') AS divisions,
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
LEFT JOIN po_props pp ON pp.uid = u.id
```

## Fields: user_po_approver_profile

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `max_amount` | json | maxSize 1 |
| `project_max` | json | maxSize 1 |
| `sponsorship_max` | json | maxSize 1 |
| `staff_and_social_max` | json | maxSize 1 |
| `media_and_event_max` | json | maxSize 1 |
| `computer_max` | json | maxSize 1 |
| `divisions` | json | maxSize 1 |
| `claims` | json | maxSize 1 |
