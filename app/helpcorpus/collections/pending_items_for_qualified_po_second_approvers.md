<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: pending_items_for_qualified_po_second_approvers

`pending_items_for_qualified_po_second_approvers` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/pending_items_for_qualified_po_second_approvers/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): pending_items_for_qualified_po_second_approvers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): pending_items_for_qualified_po_second_approvers

No rule. Nobody can do this through the API, except a PocketBase superuser.

## The query of the view: pending_items_for_qualified_po_second_approvers

```sql
WITH timeout_config AS (
  SELECT
    CASE
      WHEN json_valid(value) = 1
       AND json_type(value, '$.second_stage_timeout_hours') IN ('real', 'integer')
       AND json_extract(value, '$.second_stage_timeout_hours') > 0
      THEN json_extract(value, '$.second_stage_timeout_hours')
      ELSE 24
    END AS timeout_hours
  FROM app_config
  WHERE key = 'purchase_orders'
  LIMIT 1
),
cfg AS (
  SELECT COALESCE((SELECT timeout_hours FROM timeout_config), 24) AS timeout_hours
),
qualified_users AS (
  SELECT
    u.id AS user_id,
    pap.divisions,
    pap.max_amount,
    pap.project_max,
    pap.sponsorship_max,
    pap.staff_and_social_max,
    pap.media_and_event_max,
    pap.computer_max
  FROM users u
  JOIN admin_profiles ap ON ap.uid = u.id AND ap.active = 1
  JOIN user_claims uc ON u.id = uc.uid
  JOIN claims c ON uc.cid = c.id AND c.name = 'po_approver'
  JOIN po_approver_props pap ON uc.id = pap.user_claim
),
pos_needing_second_approval AS (
  SELECT
    po.id AS po_id,
    po.approval_total,
    po.division,
    po.job,
    COALESCE(ek.name, CASE WHEN po.job != '' THEN 'project' ELSE 'capital' END) AS kind_name,
    COALESCE(
      ek.second_approval_threshold,
      ek_fallback.second_approval_threshold,
      0
    ) AS second_approval_threshold
  FROM purchase_orders po
  LEFT JOIN expenditure_kinds ek ON po.kind = ek.id
  LEFT JOIN expenditure_kinds ek_fallback
    ON ek.id IS NULL
    AND ek_fallback.name = CASE WHEN po.job != '' THEN 'project' ELSE 'capital' END
  CROSS JOIN cfg
  WHERE
    po.approved != ''
    AND po.rejected = ''
    AND po.status = 'Unapproved'
    AND po.second_approval = ''
    AND po.approved < strftime('%Y-%m-%d %H:%M:%fZ', 'now', '-' || CAST(cfg.timeout_hours AS TEXT) || ' hours')
    AND COALESCE(
      ek.second_approval_threshold,
      ek_fallback.second_approval_threshold,
      0
    ) > 0
    AND po.approval_total > COALESCE(
      ek.second_approval_threshold,
      ek_fallback.second_approval_threshold,
      0
    )
),
qualified_candidates AS (
  SELECT
    qu.user_id,
    po.po_id,
    po.approval_total,
    po.second_approval_threshold,
    CASE po.kind_name
      WHEN 'capital' THEN COALESCE(qu.max_amount, 0)
      WHEN 'project' THEN COALESCE(qu.project_max, 0)
      WHEN 'sponsorship' THEN COALESCE(qu.sponsorship_max, 0)
      WHEN 'staff_and_social' THEN COALESCE(qu.staff_and_social_max, 0)
      WHEN 'media_and_event' THEN COALESCE(qu.media_and_event_max, 0)
      WHEN 'computer' THEN COALESCE(qu.computer_max, 0)
      ELSE 0
    END AS resolved_limit
  FROM qualified_users qu
  JOIN pos_needing_second_approval po
    ON (
      json_valid(qu.divisions)
      AND (
        json_array_length(qu.divisions) = 0
        OR EXISTS (SELECT 1 FROM json_each(qu.divisions) WHERE value = po.division)
      )
    )
),
qualified_pairs AS (
  SELECT
    qc.user_id,
    qc.po_id
  FROM qualified_candidates qc
  WHERE qc.resolved_limit > qc.second_approval_threshold
    AND qc.resolved_limit >= qc.approval_total
)
SELECT
  qp.user_id AS id,
  COUNT(qp.po_id) AS num_pos_qualified
FROM qualified_pairs qp
GROUP BY qp.user_id
```

## Fields: pending_items_for_qualified_po_second_approvers

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `num_pos_qualified` | number | onlyInt |
