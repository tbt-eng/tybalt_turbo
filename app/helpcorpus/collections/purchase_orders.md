<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: purchase_orders

`purchase_orders` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/purchase_orders/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): purchase_orders

```
(
(status = "Unapproved" && po_number != "" && legacy_manual_entry = false && rejected = "" && @request.auth.id != "") || (status = "Active" && @request.auth.id != "") ||
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
) || (
po_group != '' &&
(
  @request.auth.user_claims_via_uid.cid.name ?= 'super_po' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
)
```

## Who can view one record (view rule): purchase_orders

```
(
(status = "Unapproved" && po_number != "" && legacy_manual_entry = false && rejected = "" && @request.auth.id != "") || (status = "Active" && @request.auth.id != "") ||
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
) || (
po_group != '' &&
(
  @request.auth.user_claims_via_uid.cid.name ?= 'super_po' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'admin'
)
)
```

## Who can create a record (create rule): purchase_orders

```
// the caller is authenticated
@request.auth.id != "" &&

// standard flow cannot set legacy_manual_entry
(@request.body.legacy_manual_entry:isset = false || @request.body.legacy_manual_entry = false) &&

// no po_number is submitted
(@request.body.po_number:isset = false || @request.body.po_number = "") &&

// status is Unapproved
@request.body.status = "Unapproved" &&

// the uid is missing or is equal to the authenticated user's id
(@request.body.uid:isset = false || @request.body.uid = @request.auth.id) &&

// no rejection properties are submitted
@request.body.rejector:isset = false &&
@request.body.rejected:isset = false &&
@request.body.rejection_reason:isset = false &&

// approved isn't set. We check that the approver has the appropriate claim and divisions in payload in hooks, commenting out the previous check for po_approver here. 
@request.body.approved:isset = false &&
//@request.body.approver.user_claims_via_uid.cid.name ?= 'po_approver' &&

// no second approver properties are submitted
@request.body.second_approver:isset = false &&
@request.body.second_approval:isset = false &&

// no cancellation properties are submitted
@request.body.cancelled:isset = false &&
@request.body.canceller:isset = false &&

// no closed properties are submitted
@request.body.closed:isset = false &&
@request.body.closer:isset = false &&
@request.body.closed_by_system:isset = false &&

// vendor is active (disabled, perform this in the hook for better error messages)
// @request.body.vendor.status = "Active" &&

// if present, the category belongs to the job, otherwise is blank
(
  // compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can update a record (update rule): purchase_orders

```
legacy_manual_entry = false &&
(@request.body.legacy_manual_entry:isset = false || @request.body.legacy_manual_entry = false) &&

// only the creator can update the record
uid = @request.auth.id &&

// the uid must not change
@request.body.uid:changed = false &&

// status is Unapproved and second approval has not been performed
status = 'Unapproved' &&
second_approval = "" &&

// no po_number is submitted
(@request.body.po_number:isset = false || po_number = @request.body.po_number) &&

// no rejection properties are submitted
(@request.body.rejector:isset = false || rejector = @request.body.rejector) &&
(@request.body.rejected:isset = false || rejected = @request.body.rejected) &&
(@request.body.rejection_reason:isset = false || rejection_reason = @request.body.rejection_reason) &&

// approved is unchanged
(@request.body.approved:isset = false || approved = @request.body.approved) &&

// no second approver properties are submitted
(@request.body.second_approver:isset = false || second_approver = @request.body.second_approver) &&
(@request.body.second_approval:isset = false || second_approval = @request.body.second_approval) &&

// no cancellation properties are submitted
(@request.body.cancelled:isset = false || cancelled = @request.body.cancelled) &&
(@request.body.canceller:isset = false || canceller = @request.body.canceller) &&

// no closed properties are submitted
(@request.body.closed:isset = false || closed = @request.body.closed) &&
(@request.body.closer:isset = false || closer = @request.body.closer) &&
(@request.body.closed_by_system:isset = false || closed_by_system = @request.body.closed_by_system) &&

// vendor is active (disabled, perform this in the hook for better error messages)
// @request.body.vendor.status = "Active" &&

// if present, the category belongs to the job, otherwise is blank
(
  // the job is unchanged, compare the new category to job
  ( @request.body.job:isset = false && @request.body.category.job = job ) ||
  // the job has changed, compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
 )
```

## Who can delete a record (delete rule): purchase_orders

```
@request.auth.id = uid && status = 'Unapproved' && po_number = ''
```

## Fields: purchase_orders

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
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
| `vendor` | relation | maxSelect 1; to `vendors` |
| `parent_po` | relation | maxSelect 1; to `purchase_orders` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `closer` | relation | maxSelect 1; to `users` |
| `closed` | date |  |
| `closed_by_system` | bool |  |
| `priority_second_approver` | relation | maxSelect 1; to `users` |
| `approval_total` | number | required |
| `_imported` | bool |  |
| `branch` | relation | maxSelect 1; to `branches` |
| `attachment_hash` | text |  |
| `kind` | relation | maxSelect 1; to `expenditure_kinds` |
| `legacy_manual_entry` | bool |  |
| `currency` | relation | maxSelect 1; to `currencies` |
| `approval_total_home` | number |  |
| `covered_within_project_budget` | bool |  |
| `occurrences` | number | onlyInt |
| `po_group` | relation | maxSelect 1; to `purchase_order_groups` |

## Indexes: purchase_orders

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_6Ao8pCT` ON `purchase_orders` (`po_number`) WHERE `po_number` != ''`
- `CREATE INDEX `idx_lVCg50dCG9` ON `purchase_orders` ( `job`, `date DESC` ) WHERE status = 'Active'`
- `CREATE UNIQUE INDEX `idx_Ml6Pmg44QP` ON `purchase_orders` (`attachment_hash`) WHERE `attachment_hash` != ''`
- `CREATE INDEX `idx_purchase_orders_po_group` ON `purchase_orders` (`po_group`) WHERE `po_group` != ''`
