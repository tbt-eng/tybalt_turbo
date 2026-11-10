<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/purchase_orders.go","sha256":"36ac585b3dbf3f5ce100339ea502c7949721b02381e399a1ed3745c2e0f6945a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/purchase_orders.go

## Approval total equals total except on recurring POs

Lines 33-70

Turbo sets a PO's approval total (approval_total) to its total and its occurrences to 0. On a Recurring PO it replaces both with the occurrence count and the calculated total value. The approver approves the approval total.

## One-Time and Cumulative POs clear end date and frequency

Lines 33-47

When a purchase order (PO) has type One-Time or Cumulative, Turbo clears its end date (end_date) and its frequency (frequency). Turbo also sets the occurrences count to 0 for these types.

## PO currency falls back to parent PO, then home currency

Lines 72-84

When a PO has no currency, Turbo uses its parent PO's currency if it has one. If the currency is still blank, Turbo uses the home currency.

## Unknown currency on a PO is refused

Lines 86-95

Turbo refuses to save a PO whose currency does not exist. It shows 'referenced currency not found'.

## Foreign currency without exchange rate is refused

Lines 96-107

Turbo refuses a PO in a foreign currency that has no exchange rate. It shows 'selected foreign currency is missing an exchange rate'.

## PO approval total in home currency is calculated

Lines 108-108

Turbo sets approval_total_home to the approval total times the currency's exchange rate. It uses a rate of 1 when the currency has no rate to apply (CurrencyRateOrOne).

## Editing a PO clears its rejection

Lines 110-118

Each time a PO is created or edited, Turbo clears the rejected date, the rejection reason and the rejector. Saving after a rejection prepares the PO for resubmission.

## PO with a job takes the job's branch

Lines 123-147

When a PO has a job, Turbo sets the PO's branch to the job's branch, whatever branch was submitted. Turbo refuses the PO with 'referenced job not found' if the job does not exist. It refuses with 'referenced job is missing a branch' if the job has no branch.

## PO without job and branch uses owner's default branch

Lines 148-171

When a PO has no job and a blank branch, Turbo uses the default branch (default_branch) from the owner's admin profile. Turbo refuses with 'error looking up admin profile' if the profile is missing. It refuses with 'your admin_profiles record is missing a default_branch' if the profile has no default branch. A branch the owner chose is kept.

## PO requires a vendor

Lines 199-210

Turbo refuses a PO with no vendor and shows 'vendor is required'.

## PO requires a valid expenditure kind

Lines 212-238

Turbo refuses a PO with no kind and shows 'kind is required'. It shows 'invalid expenditure kind' if the kind does not exist.

## Kind must allow a job when PO has a job

Lines 239-252

Turbo refuses a PO that has a job when its expenditure kind does not allow jobs (allow_job). It shows 'selected kind does not allow job'.

## Project kind requires a job

Lines 253-264

Turbo refuses a PO of the project expenditure kind that has no job. It shows 'a job is required for the project expenditure kind'.

## PO vendor must exist and be Active

Lines 266-293

Turbo refuses a PO whose vendor is not found ('vendor not found') or whose status is not Active ('provided vendor is not active').

## PO total must be greater than zero

Lines 295-307

Turbo refuses a PO whose total is 0 or less. It shows 'total must be greater than 0'.

## Child PO parent must be an Active Cumulative top-level PO

Lines 309-363

A child PO's parent must exist ('parent PO not found'). The parent cannot itself be a child ('parent PO cannot be itself a child'). It must be Active ('parent PO must be active') and Cumulative ('parent PO must be cumulative').

## Child PO inherits parent currency when blank

Lines 365-367

When a child PO has no currency and its parent has one, Turbo copies the parent's currency to the child.

## Only one open child PO per parent

Lines 369-399

Turbo refuses a child PO when another child of the same parent has a status other than Closed or Cancelled. It shows "other child POs that are not 'Closed' or 'Cancelled' already exist". The PO being edited is not counted, so an existing child can be edited.

## Child PO fields must match parent

Lines 401-418

A child PO must have the same job, payment type, category, description, vendor and kind as its parent. Turbo shows 'field X must match parent PO's X' for the first field that differs.

## Child PO currency must match parent

Lines 420-433

Turbo refuses a child PO whose currency differs from its parent's currency. It shows 'child PO currency must match parent PO currency'. A blank currency counts as the home currency.

## PO date must be a valid date

Lines 436-448

Turbo refuses a PO whose date is not a valid YYYY-MM-DD date. It shows 'date must be a valid date'. The date is also required.

## PO owner and approver must be active users

Lines 451-519

Turbo refuses a PO when the owner is not an active user ('the selected staff member is not an active user'). It refuses when the approver is not active ('the selected approver is not an active user'). It refuses when the priority second approver is not active ('The selected priority second approver is not an active user').

## Recurring PO end date rules

Lines 501-509

A Recurring PO requires an end date ('end_date is required for recurring purchase orders'). The end date must not be before the start date ('end date must be after start date'). It must be within the RECURRING_MAX_DAYS limit of the start date ('end date must be within N days of the start date'). A non-recurring PO cannot have an end date ('end_date is not permitted for non-recurring purchase orders').

## Recurring PO requires frequency; others forbid it

Lines 511-516

A Recurring PO requires a frequency ('frequency is required for recurring purchase orders'). A non-recurring PO cannot have one ('frequency is not permitted for non-recurring purchase orders').

## PO description minimum length

Lines 517-517

A PO description must have at least 4 characters. Turbo shows 'must be at least 4 characters'.

## PO total has a maximum

Lines 520-520

Turbo refuses a PO whose total is above MAX_APPROVAL_TOTAL.

## Child POs must be One-Time

Lines 521-521

Turbo refuses a child PO whose type is not One-Time. It shows 'child POs must be of type One-Time'.

## Division allocation checked only on create or changed job/division

Lines 524-545

When a PO has a job, Turbo checks that the division is allocated to the job only when the PO is created or the job or division is changed. Old POs can be edited for other fields despite legacy allocation mismatches.

## PO job must exist, be Active and match branch

Lines 533-540

A PO's job must exist ('invalid job reference') and must have status Active ('Job status must be Active'). The PO's branch must equal the job's branch ('branch must match the selected job').

## Covered within project budget needs a project job

Lines 549-562

Turbo refuses the 'covered within project budget' confirmation when the PO has no job or its job is a proposal. It shows 'covered within project budget may only be confirmed for purchase orders linked to a project job'.

## Approver must be valid for first stage

Lines 570-608

Turbo refuses a PO whose approver is not a first-stage approver for the division, amount and kind, unless the approver is valid via the bypass for the owner. It shows 'selected approver is not valid for first-stage approval'.

## Second approval requires a valid priority second approver

Lines 610-621

When second approval is required, a priority second approver is required ('priority second approver is required when second approval is required'). The approver must be a second-stage approver ('selected priority second approver is not valid for second-stage approval'). Turbo also reports when no first-stage or second-stage approvers are configured.

## Priority second approver cleared when not needed

Lines 622-624

When second approval is not required, Turbo clears any priority second approver on the PO.

## Revised PO editable fields list

Lines 630-679

On a revised purchase order (a revision), the creator may change only the total, currency, approver and priority second approver. Turbo also lets the computed fields (approval_total, approval_total_home, occurrences) and status and approval fields (purchaseOrderMeaningfulChangeSkipFields) differ. A change to any other field makes Turbo refuse the edit (lockedRevisionFieldChange). The comparison ignores created and updated.

## Approvals reset only for first-approved Unapproved POs

Lines 630-706

When a creator edits a purchase order, Turbo resets its approvals only if the PO is Unapproved, has a first approval (approved) and has no second approval. A draft with no first approval does not reset. Turbo skips the reset on new records and on Active, Closed or Cancelled POs. The edit must change a meaningful field, which excludes approval, rejection, cancel, close, po_number and status fields.

## Approvals reset when a first-approved PO is edited

Lines 681-706

Turbo resets approvals when the creator edits an Unapproved PO that has a first approval but no second approval. It resets only if there are meaningful changes, or if a rejected revision is resubmitted. It clears the approval and second approval, and keeps the newly submitted approver as the next first-stage approver. Edits to drafts with no first approval do not reset anything.

## Rejected revision resubmission restarts approval

Lines 699-703

Resubmitting a rejected revision restarts its approval even when nothing changed. This applies when the revision is Unapproved, has a first approval and has no second approval. Its first approver is then asked and told again.

## Revision edit emails approver on key changes

Lines 708-721

When a revision with no first approval is edited, Turbo emails the approver if it was rejected, or if the approver, total or currency changes. A blank currency counts as the home currency.

## Revision approval email shows previous total only if different

Lines 723-760

The revision approval email goes to the approver and shows the approval total with its currency. It shows the previous total only when it differs from the new one.

## Approval total text uses approval total and currency

Lines 723-727

The approval total text in notifications shows the approval total (not the single-occurrence total on recurring POs) with the PO's effective currency code.

## Revision approver email shows previous total only if changed

Lines 729-743

The email that asks a revision's approver to approve again (template po_revision_approval_required) shows the previous approval total only when it differs from the new total with its currency. Turbo blanks the previous total when the two texts are equal.

## New PO emails the approver

Lines 762-786

When a PO is created with an approver, Turbo sends that approver an approval-required email. A revision gets the revision email instead. Turbo sends it only after the save succeeds.

## No approval email without an approver

Lines 766-786

Turbo sends no approval-required notification when the purchase order has no approver. A revision gets the revision email, comparing against the original's approval total. Other POs get the standard approval-required email, deferred until the save succeeds.

## PO creation and edits need an authenticated user

Lines 795-812

Turbo refuses PO changes when no one is authenticated. It shows 'authentication is required'. It also refuses while expense editing is disabled (checkExpensesEditing).

## Expense editing can block PO changes

Lines 795-797

Before processing a purchase order create or update, Turbo runs checkExpensesEditing. If that check returns an error, Turbo refuses the action with that error. The rule it applies lives in another file.

## New PO owner must be the signed-in user

Lines 814-827

When creating a PO, the owner (uid) must be the signed-in user. Otherwise Turbo shows 'uid must be equal to the authenticated user's id'.

## Only the creator can edit a PO

Lines 828-867

Only the creator can edit a PO. Otherwise Turbo shows 'only the creator can edit this purchase order'. The owner cannot be changed ('uid cannot be changed').

## Only Unapproved POs can be edited

Lines 869-880

Turbo refuses to edit a PO whose status is not Unapproved. It shows 'only unapproved purchase orders can be edited'.

## Revised PO locks all but total, currency and approvers

Lines 882-897

On a revised PO, the creator can change only the total, currency and approvers. Other changed fields are refused with 'only the total, currency and approvers can change on a revised purchase order'.

## PO group cannot be set or changed by edit

Lines 900-918

A PO's group (po_group) cannot be set or changed by a normal create or update. Turbo shows 'po_group can only be set when a purchase order group is created'.

## PO status must stay Unapproved on create or update

Lines 920-934

Turbo refuses a PO create or update whose status is not Unapproved. It shows 'status must be Unapproved when creating or updating purchase orders'. Other status changes go through the dedicated actions.

## Owner must be allowed to use the PO branch

Lines 945-947

Turbo checks that the user can use the PO's branch (EnsureUserCanUseRecordBranch) after cleaning, and refuses the save if not.

## Reset clears approvals but keeps submitted approver

Lines 949-957

When Turbo resets approvals on an edit, it clears the first approval, second approval and second approver. It keeps the approver submitted with the edit as the next first-stage approver.

## Duplicate PO attachment is refused

Lines 961-992

Turbo refuses a new attachment whose hash matches that of another PO. It shows 'This file has already been uploaded to another purchase order'. Turbo stores the hash of a new file, and clears the hash when the attachment is removed.

## PO division must be active

Lines 994-999

Turbo checks that the PO's division is active (EnsureActiveDivision) and refuses the save with a validation error if not.

## Approval emails sent on create, reset or key revision change

Lines 1010-1032

Turbo creates the approver notification when a PO is new. On an update it creates one when approvals were reset or when a revision change needs an approver email. Otherwise it sends none. Delivery waits until the save succeeds.
