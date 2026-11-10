<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: user_view_grants

`user_view_grants` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/user_view_grants/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): user_view_grants

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): user_view_grants

No rule. Nobody can do this through the API, except a PocketBase superuser.

## The query of the view: user_view_grants

```sql
SELECT
  g.id AS id,
  u.id AS uid,
  v.id AS viewer,
  CAST(g.reason AS TEXT) AS reason
FROM (
  SELECT po.uid AS id, po.uid AS uid, '' AS viewer, 'active_po' AS reason
  FROM purchase_orders po
  WHERE po.status = 'Active' AND po.uid != ''
  UNION
  SELECT po.approver, po.approver, '', 'active_po'
  FROM purchase_orders po
  WHERE po.status = 'Active' AND po.approver != ''
  UNION
  SELECT po.second_approver, po.second_approver, '', 'active_po'
  FROM purchase_orders po
  WHERE po.status = 'Active' AND po.second_approver != ''
  UNION ALL
  SELECT p.uid || p.approver, p.uid, p.approver, 'expense_approver'
  FROM (
    SELECT e.uid, e.approver
    FROM expenses e
    WHERE e.committed != '' AND e.uid != '' AND e.approver != ''
    GROUP BY e.uid, e.approver
  ) p
) g
JOIN users u ON u.id = g.uid
LEFT JOIN users v ON v.id = g.viewer
```

## Fields: user_view_grants

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `uid` | relation | maxSelect 1; to `users` |
| `viewer` | relation | maxSelect 1; to `users` |
| `reason` | text |  |
