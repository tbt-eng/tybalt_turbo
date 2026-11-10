<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expenses

`expenses` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/expenses/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expenses

```
uid = @request.auth.id ||
creator = @request.auth.id ||
(approver = @request.auth.id && submitted = true) ||
(approved != "" && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(committed != "" && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can view one record (view rule): expenses

```
uid = @request.auth.id ||
creator = @request.auth.id ||
(approver = @request.auth.id && submitted = true) ||
(approved != "" && @request.auth.user_claims_via_uid.cid.name ?= 'commit') ||
(committed != "" && @request.auth.user_claims_via_uid.cid.name ?= 'report')
```

## Who can create a record (create rule): expenses

```
// the caller is authenticated
@request.auth.id != "" &&

// the pay_period_ending is not set or changed
@request.body.pay_period_ending:changed = false &&

// creator is server-managed by the hook
(@request.body.creator:isset = false || @request.body.creator = @request.auth.id) &&

// regular owner flow, or narrow bookkeeper-on-behalf PO flow
(
  @request.body.uid = @request.auth.id ||
  (
    @request.auth.user_claims_via_uid.cid.name ?= "book_keeper" &&
    @request.body.purchase_order:isset = true &&
    @request.body.purchase_order != "" &&
    @request.body.uid != @request.auth.id &&
    (
      @request.body.payment_type = "OnAccount" ||
      @request.body.payment_type = "CorporateCreditCard"
    )
  )
) &&

// no rejection properties are submitted
@request.body.rejector:isset = false &&
@request.body.rejected:isset = false &&
@request.body.rejection_reason:isset = false &&

// no approval properties are submitted
@request.body.approved:isset = false &&
@request.body.approver:isset = false &&

// no committed properties are submitted
@request.body.committed:isset = false &&
@request.body.committer:isset = false &&
@request.body.committed_week_ending:isset = false &&

// if present, vendor is active
(@request.body.vendor = "" || @request.body.vendor.status = "Active") &&

// if present, the category belongs to the job, otherwise is blank
(
  // compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can update a record (update rule): expenses

```
// only the effective draft owner can update the record
@request.auth.id = creator &&

// the pay_period_ending is not set or changed
@request.body.pay_period_ending:changed = false &&

// ownership fields must not change
@request.body.uid:changed = false &&
@request.body.creator:changed = false &&

// no rejection properties are submitted
(@request.body.rejector:isset = false || rejector = @request.body.rejector) &&
(@request.body.rejected:isset = false || rejected = @request.body.rejected) &&
(@request.body.rejection_reason:isset = false || rejection_reason = @request.body.rejection_reason) &&

// submitted is not changed
(@request.body.submitted:isset = false || submitted = @request.body.submitted) &&

// no approval properties are submitted
(@request.body.approved:isset = false || approved = @request.body.approved) &&
(@request.body.approver:isset = false || approver = @request.body.approver) &&

// no committed properties are submitted
(@request.body.committed:isset = false || committed = @request.body.committed) &&
(@request.body.committer:isset = false || committer = @request.body.committer) &&
(@request.body.committed_week_ending:isset = false || committed_week_ending = @request.body.committed_week_ending) &&

// if present, vendor is active
(@request.body.vendor = "" || @request.body.vendor.status = "Active") &&

// if present, the category belongs to the job, otherwise is blank
(
  // the job is unchanged, compare the new category to job
  ( @request.body.job:isset = false && @request.body.category.job = job ) ||
  // the job has changed, compare the new category to the new job
  ( @request.body.job:isset = true && @request.body.category.job = @request.body.job ) ||
  @request.body.category = ""
)
```

## Who can delete a record (delete rule): expenses

```
@request.auth.id != "" &&
submitted = false &&
committed = "" &&
@request.auth.id = creator
```

## Fields: expenses

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users` |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `division` | relation | required; maxSelect 1; to `divisions` |
| `description` | text |  |
| `total` | number |  |
| `payment_type` | select | required; values ["OnAccount","Expense","CorporateCreditCard","Allowance","FuelCard","Mileage","PersonalReimbursement"]; maxSelect 1 |
| `rejector` | relation | maxSelect 1; to `users` |
| `rejected` | date |  |
| `rejection_reason` | text | min 5 |
| `approver` | relation | maxSelect 1; to `users` |
| `approved` | date |  |
| `creator` | relation | maxSelect 1; to `users` |
| `job` | relation | maxSelect 1; to `jobs` |
| `category` | relation | maxSelect 1; to `categories` |
| `pay_period_ending` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `allowance_types` | select | values ["Lodging","Breakfast","Lunch","Dinner"]; maxSelect 4 |
| `submitted` | bool |  |
| `committer` | relation | maxSelect 1; to `users` |
| `committed` | date |  |
| `committed_week_ending` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `distance` | number |  |
| `cc_last_4_digits` | text | pattern "^\\d{4}$" |
| `purchase_order` | relation | maxSelect 1; to `purchase_orders` |
| `vendor` | relation | maxSelect 1; to `vendors` |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |
| `_imported` | bool |  |
| `branch` | relation | maxSelect 1; to `branches` |
| `kind` | relation | maxSelect 1; to `expenditure_kinds` |
| `currency` | relation | maxSelect 1; to `currencies` |
| `settled_total` | number |  |
| `settler` | relation | maxSelect 1; to `users` |
| `settled` | date |  |
| `attachment_document` | relation | maxSelect 1; to `expense_documents` |
| `attachment_missing_reason` | text |  |

## Indexes: expenses

A unique index means that two records can't have the same values in its fields.

- `CREATE INDEX `idx_8LRpecUoxd` ON `expenses` ( `purchase_order`, `committed` )`
- `CREATE INDEX `idx_slBmqtw6SZ` ON `expenses` (`date`)`
- `CREATE INDEX `idx_3TRP1AbuJv` ON `expenses` ( `branch`, `job` )`
- `CREATE INDEX `idx_expenses_uid_date` ON `expenses` (`uid`, `date`)`
- `CREATE INDEX `idx_expenses_approver_submitted_date` ON `expenses` (`approver`, `submitted`, `date`)`
- `CREATE INDEX `idx_expenses_po_date` ON `expenses` (`purchase_order`, `date`)`
- `CREATE INDEX `idx_expenses_approved_nonempty` ON `expenses` (`approved`) WHERE `approved` != ''`
- `CREATE INDEX `idx_expenses_committed_nonempty` ON `expenses` (`committed`) WHERE `committed` != ''`
- `CREATE INDEX `idx_Y3uLpJvqvc` ON `expenses` (`committed_week_ending`)`
- `CREATE INDEX `idx_expenses_creator_date` ON `expenses` (`creator`, `date`)`
- `CREATE INDEX `idx_expenses_creator_submitted_date` ON `expenses` (`creator`, `submitted`, `date`)`
- `CREATE INDEX `idx_expenses_job_date` ON `expenses` (`job`, `date`)`
- `CREATE INDEX `idx_expenses_uid_approver_committed` ON `expenses` (`uid`, `approver`) WHERE `committed` != ''`
