<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: admin_profiles_augmented

`admin_profiles_augmented` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/admin_profiles_augmented/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): admin_profiles_augmented

```
@request.auth.id != "" &&
(
  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'hr' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'time_off_manager'
)
```

## Who can view one record (view rule): admin_profiles_augmented

```
@request.auth.id != "" &&
(
  @request.auth.user_claims_via_uid.cid.name ?= 'admin' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'hr' ||
  @request.auth.user_claims_via_uid.cid.name ?= 'time_off_manager'
)
```

## The query of the view: admin_profiles_augmented

```sql
SELECT ap.id,
  ap.uid,
  ap.active,
  ap.work_week_hours,
  ap.salary,
  ap.default_charge_out_rate,
  ap.off_rotation_permitted,
  ap.skip_min_time_check,
  ap.opening_date,
  ap.opening_op,
  ap.opening_ov,
  ap.payroll_id,
  ap.untracked_time_off,
  ap.time_sheet_expected,
  ap.allow_personal_reimbursement,
  ap.mobile_phone,
  ap.job_title,
  ap.personal_vehicle_insurance_expiry,
  ap.default_branch,
  p.given_name,
  p.surname,
  po.po_approver_props_id,
  po.po_approver_user_claim_id,
  po.po_approver_max_amount,
  po.po_approver_project_max,
  po.po_approver_sponsorship_max,
  po.po_approver_staff_and_social_max,
  po.po_approver_media_and_event_max,
  po.po_approver_computer_max,
  COALESCE(po.po_approver_divisions, '[]') AS po_approver_divisions
FROM admin_profiles ap
LEFT JOIN users u ON u.id = ap.uid
LEFT JOIN profiles p ON u.id = p.uid
LEFT JOIN (
  SELECT
    uc.uid,
    uc.id AS po_approver_user_claim_id,
    pap.id AS po_approver_props_id,
    pap.max_amount AS po_approver_max_amount,
    pap.project_max AS po_approver_project_max,
    pap.sponsorship_max AS po_approver_sponsorship_max,
    pap.staff_and_social_max AS po_approver_staff_and_social_max,
    pap.media_and_event_max AS po_approver_media_and_event_max,
    pap.computer_max AS po_approver_computer_max,
    pap.divisions AS po_approver_divisions
  FROM user_claims uc
  INNER JOIN claims c ON c.id = uc.cid AND c.name = 'po_approver'
  LEFT JOIN po_approver_props pap ON pap.user_claim = uc.id
) po ON po.uid = ap.uid;
```

## Fields: admin_profiles_augmented

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users` |
| `active` | bool |  |
| `work_week_hours` | number | required; min 8; max 40 |
| `salary` | bool |  |
| `default_charge_out_rate` | number | required; min 50; max 1000 |
| `off_rotation_permitted` | bool |  |
| `skip_min_time_check` | select | required; values ["no","on_next_bundle","yes"]; maxSelect 1 |
| `opening_date` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `opening_op` | number | max 332 |
| `opening_ov` | number | max 200 |
| `payroll_id` | text | required; pattern "^(?:[1-9]\\d*\|CMS[0-9]{1,2})$" |
| `untracked_time_off` | bool |  |
| `time_sheet_expected` | bool |  |
| `allow_personal_reimbursement` | bool |  |
| `mobile_phone` | text | pattern "^\\+1 \\(\\d{3}\\) \\d{3}-\\d{4}$" |
| `job_title` | text |  |
| `personal_vehicle_insurance_expiry` | text | pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `default_branch` | relation | required; maxSelect 1; to `branches` |
| `given_name` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:-[a-zA-Z]+)*$" |
| `surname` | text | required; min 2; max 48; pattern "^[a-zA-Z]+(?:[-'][a-zA-Z]+)*$" |
| `po_approver_props_id` | json | maxSize 1 |
| `po_approver_user_claim_id` | json | maxSize 1 |
| `po_approver_max_amount` | json | maxSize 1 |
| `po_approver_project_max` | json | maxSize 1 |
| `po_approver_sponsorship_max` | json | maxSize 1 |
| `po_approver_staff_and_social_max` | json | maxSize 1 |
| `po_approver_media_and_event_max` | json | maxSize 1 |
| `po_approver_computer_max` | json | maxSize 1 |
| `po_approver_divisions` | json | maxSize 1 |
