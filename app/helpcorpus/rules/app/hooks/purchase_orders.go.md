<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/purchase_orders.go","sha256":"5e22cadc06e745f7194678d83b89113e05288a514684b8676a4dd119fbd51a3b","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/purchase_orders.go

## Approval total equals total unless the PO is recurring

Lines 30-70

Asked as: approval total wrong; approval total differs from total; occurrences showing zero; approval amount not matching

Turbo sets a purchase order's approval total (approval_total) to its total and its occurrences to 0. For a Recurring PO, Turbo replaces both with the values it calculates from the PO (CalculateRecurringPurchaseOrderTotalValue).

## Approval total and occurrences reset before recurring calculation

Lines 30-34

Asked as: recurring po approval total reset; approval total changed on save; occurrences zero after saving; recurring po numbers overwritten

When a PO is saved, Turbo first sets its approval total (approval_total) to the total and its occurrences to 0. A recurring PO then replaces both with values from the recurring calculation.

## One-Time and Cumulative POs lose end date and frequency

Lines 42-47

Asked as: end date disappeared; frequency cleared on one-time po; cumulative po lost end date; frequency blank after save

When a purchase order is One-Time or Cumulative, Turbo clears its end date (end_date) and frequency on save.

## PO currency falls back to parent, then home currency

Lines 72-84

Asked as: po currency blank; currency empty on purchase order; child po currency missing; currency filled from parent

When a purchase order has no currency, Turbo takes its parent PO's currency if the parent has one. If the currency is still blank, Turbo uses the home currency.

## Unknown currency or foreign currency without rate is refused

Lines 86-107

Asked as: currency not found error; foreign currency no exchange rate; currency rejected on save; exchange rate missing error

Turbo refuses a purchase order whose currency does not exist, with 'referenced currency not found'. It refuses a foreign currency that has no positive exchange rate, with 'selected foreign currency is missing an exchange rate'.

## Approval total in home currency is calculated

Lines 108-108

Asked as: home currency approval total; approval total in home currency; converted approval amount wrong; approval total home zero

Turbo sets approval_total_home to the approval total multiplied by the currency's exchange rate. It multiplies by one when the currency has no rate to apply (CurrencyRateOrOne).

## Saving a PO clears its rejection

Lines 110-118

Asked as: rejection disappeared after edit; rejected po cleared on resubmit; rejection reason gone; rejector removed after saving

Each create or update of a purchase order clears its rejected date, rejection reason and rejector. Turbo assumes the creator is resubmitting after changes.

## PO branch comes from its job

Lines 120-147

Asked as: po branch changed to job branch; job has no branch; branch wrong on purchase order; referenced job missing branch

When a PO has a job, Turbo sets its branch to the job's branch. It refuses an unknown job with 'referenced job not found', and a job with no branch with 'referenced job is missing a branch'.

## Without a job, a blank branch uses the owner's default branch

Lines 148-171

Asked as: no branch chosen on po; default branch used for po; admin profile default branch missing; error looking up admin profile

When a PO has no job and no branch, Turbo uses the default branch on the owner's admin profile. It refuses if the profile is missing ('error looking up admin profile') or has no default branch ('your admin_profiles record is missing a default_branch'). A branch the owner chose is kept.

## PO requires a vendor

Lines 199-210

Asked as: no vendor selected; vendor is required error; po without supplier; supplier missing on purchase order

Turbo refuses a purchase order with no vendor: 'vendor is required'.

## PO requires a valid expenditure kind

Lines 212-238

Asked as: kind is required; invalid expenditure kind; missing expense type; po type not recognised

Turbo refuses a purchase order with no kind ('kind is required') or an unknown kind ('invalid expenditure kind').

## Kind must allow a job, and project kind needs one

Lines 239-264

Asked as: selected kind does not allow job; project kind needs job; job not allowed for kind; kind refuses job on po

Turbo refuses a PO with a job when its kind does not allow jobs (allow_job): 'selected kind does not allow job'. It refuses the project expenditure kind without a job: 'a job is required for the project expenditure kind'.

## Vendor must exist and be Active

Lines 266-293

Asked as: vendor not found; inactive vendor on po; provided vendor is not active; vendor status blocks po

Turbo refuses a purchase order whose vendor is not found ('vendor not found') or whose status is not Active ('provided vendor is not active').

## PO total must be above zero

Lines 295-307

Asked as: zero total po; total must be greater than 0; negative amount on po; po amount zero refused

Turbo refuses a purchase order whose total is 0 or less: 'total must be greater than 0'.

## Parent PO must be an active cumulative top-level PO

Lines 309-363

Asked as: child po parent not active; parent must be cumulative; parent po not found; parent already a child

For a child PO, Turbo refuses a parent that is not found ('parent PO not found'), is itself a child ('parent PO cannot be itself a child'), is not Active ('parent PO must be active') or is not Cumulative ('parent PO must be cumulative').

## Child PO must match parent's fields and currency

Lines 365-433

Asked as: child po fields differ from parent; child currency mismatch; must match parent po; child po different vendor

A child PO must have the same job, payment type, category, description, vendor and kind as its parent. Turbo refuses a difference with 'field X must match parent PO's X'. A blank child currency takes the parent's currency. Turbo refuses a different currency: 'child PO currency must match parent PO currency'.

## Only one open child PO per parent

Lines 369-399

Asked as: second child po refused; another open child po exists; child po already open; parent already has child

Turbo refuses a child PO when another child of the same parent is not Closed or Cancelled: "other child POs that are not 'Closed' or 'Cancelled' already exist". The record being edited does not count.

## PO date must be valid

Lines 436-448

Asked as: invalid po date; date must be valid; po date format wrong; purchase order date rejected

Turbo refuses a purchase order whose date is not a valid date (YYYY-MM-DD): 'date must be a valid date'. The date is required.

## Owner, approver and second approver must be active users

Lines 450-519

Asked as: owner not active user; approver inactive user; second approver not active; staff member not active

Turbo refuses a PO when the owner, approver or priority second approver is not an active user. Messages: 'the selected staff member is not an active user', 'the selected approver is not an active user', 'The selected priority second approver is not an active user'.

## Recurring PO end date rules

Lines 501-516

Asked as: recurring po end date; end date too far away; frequency not permitted non-recurring; end date before start

A Recurring PO requires an end date and a frequency. The end date must not be before the start date and must be within RECURRING_MAX_DAYS days of it. A non-recurring PO refuses an end date or frequency: 'is not permitted for non-recurring purchase orders'.

## Description needs at least 4 characters

Lines 517-517

Asked as: description too short; description needs 4 characters; po description minimum length; short po description refused

Turbo refuses a purchase order description shorter than 4 characters: 'must be at least 4 characters'.

## PO total has a maximum

Lines 520-520

Asked as: po total too big; max approval total exceeded; amount over po limit; total above maximum allowed

Turbo refuses a purchase order whose total is above MAX_APPROVAL_TOTAL.

## Child POs must be One-Time

Lines 521-521

Asked as: child po type refused; child po recurring refused; child must be one time; child po not one-time

Turbo refuses a child PO whose type is not One-Time: 'child POs must be of type One-Time'.

## Job must be Active, branch must match, division allocated

Lines 524-547

Asked as: job not active on po; branch does not match job; division not allocated to job; invalid job reference

For a PO with a job, Turbo refuses an unknown job ('invalid job reference') and a job whose status is not Active ('Job status must be Active'). It refuses a branch that differs from the job's ('branch must match the selected job'). It checks that the division is allocated to the job only on creation or when the job or division changes, so older POs can still be edited.

## Category must belong to the PO's job

Lines 549-553

Asked as: category not found; category from another job; po category wrong job; no category without job

Turbo checks the category only when the PO is created or when its job or category changes. It refuses an unknown category ('category not found') and a category from another job ('category must belong to the selected job'). A PO without a job cannot have a category.

## Budget confirmation needs a project job

Lines 555-568

Asked as: covered within project budget; budget confirmation no job; budget confirm on proposal job; project budget confirmation error

Turbo refuses 'covered within project budget' on a PO with no job, or whose job is a proposal: 'covered within project budget may only be confirmed for purchase orders linked to a project job'.

## Approver must be valid for first-stage approval

Lines 576-614

Asked as: approver not valid first stage; selected approver invalid; approver not allowed for total; first stage approver refused

Turbo refuses an approver who is not a first-stage approver for the division, total in home currency, kind and job presence, unless the owner bypass applies: 'selected approver is not valid for first-stage approval'.

## Second approver required when policy needs second approval

Lines 616-627

Asked as: second approver needed; no priority second approver; second approval missing approvers configured; second stage approver required

When second approval is required, Turbo refuses a PO with no priority second approver ('priority second approver is required when second approval is required') or one not valid for second stage. It also refuses when no first-stage or second-stage approvers are configured.

## Priority second approver cleared when not needed

Lines 628-630

Asked as: second approver disappeared; priority second approver removed; second approver cleared automatically; second approver not needed cleared

When second approval is not required, Turbo clears any priority second approver on the PO.

## Editing a first-approved PO resets its approvals

Lines 636-807

Asked as: approval reset after edit; edit cleared my approval; po approval lost when changed; approvals wiped after editing

When an Unapproved PO with a first approval but no second approval is meaningfully edited, Turbo clears the approval and second approval. The submitted approver stays as the next first approver. Changes only to status, approval, rejection, cancel, close or po_number fields do not count. A rejected revision resubmitted resets even with no change. A PO with no first approval does not reset.

## Revised PO limits which fields can change

Lines 656-710

Asked as: cannot change field on revision; revision locked fields; revised po only total editable; revision refuses field change

On a revision awaiting approval, the creator can change only the total, currency, approvers, description, division, category, job and attachment. The budget confirmation and branch change only with a new job. Turbo refuses other changes: 'only the total, currency, approvers, description, division, category, job and attachment can change on a revised purchase order, and the budget confirmation only with a new job'. A submitted attachment hash is always refused.

## Revision can switch jobs but not add or remove one

Lines 712-731

Asked as: cannot add job to revision; revision remove job; revision job change refused; child po job change

A revision can move a PO to another job. Turbo refuses adding a job ('a revision cannot add a job to a purchase order that has none') and removing one ('a revision cannot remove the job'). A child PO must keep its parent's job ('a child purchase order must keep the job of its parent').

## Approver is emailed on new PO or reset

Lines 809-989

Asked as: approver not emailed; no email on new po; revision email to approver; approver not notified after reset

When a PO is created, or an edit resets its approvals, Turbo notifies the approver after the save succeeds. A revision gets a revision email (po_revision_approval_required) that shows previous values of changed fields. Revision edits before first approval email the approver on resubmission, new approver, or a new total, currency, job, description, division, category or attachment.

## Revision email changes

Lines 903-963

Asked as: revision email previous values; email says none; attachment changed in email; email shows old values

The revision approval email lists each previous value only if it changed. A blank previous value shows as 'none'. It says whether the attachment was added, replaced or removed.

## Expense editing switch blocks saving a PO

Lines 997-1000

Asked as: expense editing disabled; cannot save po expenses locked; checkExpensesEditing failed; po saving blocked by setting

Before anything else, Turbo runs a check on whether expense editing is allowed (checkExpensesEditing). If the check fails, the PO cannot be created or updated.

## Authentication required to save a PO

Lines 1003-1015

Asked as: authentication is required; not signed in to save po; save po logged out; po save needs login

Turbo refuses to create or update a purchase order without a signed-in user: 'authentication is required'.

## Only the creator can edit an Unapproved PO

Lines 1017-1083

Asked as: someone else cannot edit po; uid cannot be changed; only creator edits po; po not unapproved cannot edit

Turbo refuses an edit by anyone but the creator: 'only the creator can edit this purchase order'. It refuses a change of owner (uid): 'uid cannot be changed'. On create, the uid must equal the signed-in user's id. It refuses a PO with status other than Unapproved: 'only unapproved purchase orders can be edited'.

## Revised PO refuses changes to locked fields

Lines 1085-1099

Asked as: revision_field_locked; locked field on revised po; revised po field refused; revised po field error

When the creator edits a revised purchase order that awaits approval, Turbo compares the submitted values with the original before cleaning. If a field outside the editable list changed, it refuses with code revision_field_locked. The message is: "only the total, currency, approvers, description, division, category, job and attachment can change on a revised purchase order, and the budget confirmation only with a new job".

## Revision job change is checked on save

Lines 1100-1108

Asked as: revision job error on save; job field refused revision; revised job rules failed; revision job change error

When a revised purchase order is saved with a different job, Turbo applies the revision job rules. A revision cannot add a job to a PO that has none. It cannot remove the job. A child PO must keep its parent's job. Turbo refuses with a code and message on the job field.

## Changing a revision's job clears category and budget confirmation

Lines 1109-1118

Asked as: category cleared after job change; budget confirmation cleared on revision; job change resets category; category gone after changing job

When a revision's job changes, Turbo clears the category and the budget confirmation unless the request sets them again.

## PO group cannot be set or changed by save

Lines 1122-1140

Asked as: po_group error; cannot set po group; change purchase order group; po group refused on save

Turbo refuses a create or update that sets or changes po_group: 'po_group can only be set when a purchase order group is created'.

## Saving a PO keeps it Unapproved

Lines 1142-1156

Asked as: status must be unapproved; cannot set po active; po status refused on save; set po closed directly

Create and update requests must leave the status as Unapproved, or Turbo refuses with 'status must be Unapproved when creating or updating purchase orders'. Active, Cancelled and Closed come only from the dedicated actions.

## Attachment hash is kept or cleared

Lines 1181-1214

Asked as: attachment hash kept; removing attachment clears hash; upload file hash; attachment hash not changed

When a new file is uploaded, Turbo sets attachment_hash from it. When the attachment is removed, Turbo clears the hash. When the attachment stays and no new file is sent, the hash is left unchanged.

## Duplicate attachment refused

Lines 1183-1213

Asked as: file already uploaded; duplicate attachment on po; same file another po; attachment already used elsewhere

Turbo hashes a new attachment. It refuses a file already on another purchase order: 'This file has already been uploaded to another purchase order'. Removing the attachment clears the hash.

## Division must be active

Lines 1216-1221

Asked as: division inactive on po; inactive division error; po division not active; division disabled blocks po

Turbo checks that the PO's division is active before validating the PO.

## Saved PO is checked against the home-currency approval policy

Lines 1222-1226

Asked as: approver checks on save; approval policy on po save; second approval check every save; home currency policy check

Every PO that goes through a normal create or update is validated with legacy mode off. So the approver checks for first-stage and second-stage approval apply to every save made through this route.

## New PO creates a deferred approval notification

Lines 1228-1245

Asked as: approval notification created; notification for new po; approver notified on create; no approver notification blank

When a PO is created, Turbo generates its id and creates an approval-required notification for the approver. The email is sent only after the save succeeds. If the approver field is blank, no notification is created.

## Update emails the approver only for resets or revision changes

Lines 1246-1254

Asked as: update emails approver; no email on po edit; when approver emailed after update; edit sends no email

On an update, Turbo creates an approval notification only if the approvals were reset or a revision changed something the approver must see. Any other update sends no email.
