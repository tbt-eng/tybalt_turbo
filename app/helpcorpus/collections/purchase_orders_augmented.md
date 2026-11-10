<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: purchase_orders_augmented

`purchase_orders_augmented` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/purchase_orders_augmented/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): purchase_orders_augmented

```
(status = "Active" && @request.auth.id != "") ||
(
  (status = "Cancelled" || status = "Closed") &&
  (
    @request.auth.id = uid ||
    @request.auth.id = approver ||
    @request.auth.id = second_approver ||
    job.manager = @request.auth.id ||
    job.alternate_manager = @request.auth.id ||
    @request.auth.user_claims_via_uid.cid.name ?= "report"
  )
) ||
(
  status = "Unapproved" &&
  (
    @request.auth.id = uid ||
    @request.auth.id = approver ||
    @request.auth.id = rejector ||
    (approved != "" && second_approval = "" && @request.auth.id = priority_second_approver)
  )
)
```

## Who can view one record (view rule): purchase_orders_augmented

```
(status = "Active" && @request.auth.id != "") ||
(
  (status = "Cancelled" || status = "Closed") &&
  (
    @request.auth.id = uid ||
    @request.auth.id = approver ||
    @request.auth.id = second_approver ||
    job.manager = @request.auth.id ||
    job.alternate_manager = @request.auth.id ||
    @request.auth.user_claims_via_uid.cid.name ?= "report"
  )
) ||
(
  status = "Unapproved" &&
  (
    @request.auth.id = uid ||
    @request.auth.id = approver ||
    @request.auth.id = rejector ||
    (approved != "" && second_approval = "" && @request.auth.id = priority_second_approver)
  )
)
```

## The query of the view: purchase_orders_augmented

```sql
SELECT
  po.id,
  po.po_number,
  po.status,
  po.uid,
  po.type,
  po.date,
  po.end_date,
  po.frequency,
  po.division,
  po.description,
  po.total,
  po.payment_type,
  po.attachment,
  po.rejector,
  po.rejected,
  po.rejection_reason,
  po.approver,
  po.approved,
  po.second_approver,
  po.second_approval,
  po.canceller,
  po.cancelled,
  po.job,
  po.category,
  po.kind,
  po.vendor,
  po.parent_po,
  po.created,
  po.updated,
  po.closer,
  po.closed,
  po.closed_by_system,
  po.priority_second_approver,
  po.approval_total,
  (SELECT COUNT(*) FROM expenses WHERE expenses.purchase_order = po.id AND expenses.committed != "") AS committed_expenses_count,
  (p0.given_name || " " || p0.surname) AS uid_name,
  (p1.given_name || " " || p1.surname) AS approver_name,
  (p2.given_name || " " || p2.surname) AS second_approver_name,
  (p3.given_name || " " || p3.surname) AS priority_second_approver_name,
  (p4.given_name || " " || p4.surname) AS rejector_name,
  po2.po_number AS parent_po_number,
  v.name AS vendor_name,
  v.alias AS vendor_alias,
  j.number AS job_number,
  cl.name AS client_name,
  cl.id AS client_id,
  j.description AS job_description,
  d.code AS division_code,
  d.name AS division_name,
  c.name AS category_name
FROM purchase_orders AS po
LEFT JOIN profiles AS p0 ON po.uid = p0.uid
LEFT JOIN profiles AS p1 ON po.approver = p1.uid
LEFT JOIN profiles AS p2 ON po.second_approver = p2.uid
LEFT JOIN profiles AS p3 ON po.priority_second_approver = p3.uid
LEFT JOIN profiles AS p4 ON po.rejector = p4.uid
LEFT JOIN purchase_orders AS po2 ON po.parent_po = po2.id
LEFT JOIN vendors AS v ON po.vendor = v.id
LEFT JOIN jobs AS j ON po.job = j.id
LEFT JOIN divisions AS d ON po.division = d.id
LEFT JOIN categories AS c ON po.category = c.id
LEFT JOIN clients AS cl ON j.client = cl.id
```

## Fields: purchase_orders_augmented

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `po_number` | text | pattern "^([1-9]\\d{3})-(\\d{4})(?:-(0[1-9]\|[1-9]\\d))?$" |
| `status` | select | required; values ["Unapproved","Active","Cancelled","Closed"]; maxSelect 1 |
| `uid` | relation | required; maxSelect 1; to `users` |
| `type` | select | required; values ["One-Time","Cumulative","Recurring"]; maxSelect 1 |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `end_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `frequency` | select | values ["Weekly","Biweekly","Monthly"]; maxSelect 1 |
| `division` | relation | required; maxSelect 1; to `divisions` |
| `description` | text |  |
| `total` | number |  |
| `payment_type` | select | required; values ["OnAccount","Expense","CorporateCreditCard"]; maxSelect 1 |
| `attachment` | file | maxSelect 1; maxSize 5242880; mimeTypes ["application/pdf","image/jpeg","image/png","image/heic"] |
| `rejector` | relation | maxSelect 1; to `users` |
| `rejected` | date |  |
| `rejection_reason` | text | min 5 |
| `approver` | relation | required; maxSelect 1; to `users` |
| `approved` | date |  |
| `second_approver` | relation | maxSelect 1; to `users` |
| `second_approval` | date |  |
| `canceller` | relation | maxSelect 1; to `users` |
| `cancelled` | date |  |
| `job` | relation | maxSelect 1; to `jobs` |
| `category` | relation | maxSelect 1; to `categories` |
| `kind` | relation | maxSelect 1; to `expenditure_kinds` |
| `vendor` | relation | maxSelect 1; to `vendors` |
| `parent_po` | relation | maxSelect 1; to `purchase_orders` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `closer` | relation | maxSelect 1; to `users` |
| `closed` | date |  |
| `closed_by_system` | bool |  |
| `priority_second_approver` | relation | maxSelect 1; to `users` |
| `approval_total` | number | required |
| `committed_expenses_count` | json | maxSize 1 |
| `uid_name` | json | maxSize 1 |
| `approver_name` | json | maxSize 1 |
| `second_approver_name` | json | maxSize 1 |
| `priority_second_approver_name` | json | maxSize 1 |
| `rejector_name` | json | maxSize 1 |
| `parent_po_number` | text | pattern "^([1-9]\\d{3})-(\\d{4})(?:-(0[1-9]\|[1-9]\\d))?$" |
| `vendor_name` | text | required; min 3 |
| `vendor_alias` | text | min 3 |
| `job_number` | text | pattern "^(P)?[0-9]{2}-[0-9]{3,4}L?(-[0-9]{1,2})?(-[0-9])?$" |
| `client_name` | text | required; min 2 |
| `client_id` | relation | maxSelect 1; to `clients` |
| `job_description` | text | required; min 3 |
| `division_code` | text | required; min 1 |
| `division_name` | text | required; min 2 |
| `category_name` | text | required; min 3 |
