<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/purchase_orders.go","sha256":"36ac585b3dbf3f5ce100339ea502c7949721b02381e399a1ed3745c2e0f6945a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/purchase_orders.go

## Approval total equals total except on recurring POs

Lines 33-70

Asked as: recurring PO approval amount; approval total wrong on recurring; approver sees occurrences count; recurring total approved value

Turbo sets a PO's approval total (approval_total) to its total and its occurrences to 0. On a Recurring PO it replaces both with the occurrence count and the calculated total value. The approver approves the approval total.

## One-Time and Cumulative POs clear end date and frequency

Lines 33-47

Asked as: end date disappeared one-time; frequency cleared on cumulative; why occurrences is zero; one-time PO lost end date

When a purchase order (PO) has type One-Time or Cumulative, Turbo clears its end date (end_date) and its frequency (frequency). Turbo also sets the occurrences count to 0 for these types.

## PO currency falls back to parent PO, then home currency

Lines 72-84

Asked as: PO currency blank; currency filled in automatically; currency from parent PO; default currency on PO

When a PO has no currency, Turbo uses its parent PO's currency if it has one. If the currency is still blank, Turbo uses the home currency.

## Unknown currency on a PO is refused

Lines 86-95

Asked as: currency won't save; unknown currency error; currency not found on PO; invalid currency code

Turbo refuses to save a PO whose currency does not exist. It shows 'referenced currency not found'.

## Foreign currency without exchange rate is refused

Lines 96-107

Asked as: foreign currency won't save; no exchange rate error; exchange rate missing PO; can't use foreign currency

Turbo refuses a PO in a foreign currency that has no exchange rate. It shows 'selected foreign currency is missing an exchange rate'.

## PO approval total in home currency is calculated

Lines 108-108

Asked as: home currency approval amount; approval total in home currency; converted approval total; exchange rate conversion PO

Turbo sets approval_total_home to the approval total times the currency's exchange rate. It uses a rate of 1 when the currency has no rate to apply (CurrencyRateOrOne).

## Editing a PO clears its rejection

Lines 110-118

Asked as: resubmit rejected PO; rejection still showing after edit; clear rejection reason; rejected PO edited again

Each time a PO is created or edited, Turbo clears the rejected date, the rejection reason and the rejector. Saving after a rejection prepares the PO for resubmission.

## PO with a job takes the job's branch

Lines 123-147

Asked as: PO branch changed to job's; branch ignored on PO; job branch overrides branch; job missing branch error

When a PO has a job, Turbo sets the PO's branch to the job's branch, whatever branch was submitted. Turbo refuses the PO with 'referenced job not found' if the job does not exist. It refuses with 'referenced job is missing a branch' if the job has no branch.

## PO without job and branch uses owner's default branch

Lines 148-171

Asked as: blank branch on PO; default branch not set; admin profile default branch; branch filled from my profile

When a PO has no job and a blank branch, Turbo uses the default branch (default_branch) from the owner's admin profile. Turbo refuses with 'error looking up admin profile' if the profile is missing. It refuses with 'your admin_profiles record is missing a default_branch' if the profile has no default branch. A branch the owner chose is kept.

## PO requires a vendor

Lines 199-210

Asked as: no vendor on PO; vendor missing error; choose a supplier; PO vendor blank

Turbo refuses a PO with no vendor and shows 'vendor is required'.

## PO requires a valid expenditure kind

Lines 212-238

Asked as: expenditure kind missing; kind not chosen on PO; invalid kind error; which kind to pick

Turbo refuses a PO with no kind and shows 'kind is required'. It shows 'invalid expenditure kind' if the kind does not exist.

## Kind must allow a job when PO has a job

Lines 239-252

Asked as: job not allowed on this kind; kind doesn't allow job; can't attach job to PO; job rejected for kind

Turbo refuses a PO that has a job when its expenditure kind does not allow jobs (allow_job). It shows 'selected kind does not allow job'.

## Project kind requires a job

Lines 253-264

Asked as: project PO without job; project needs a job; job required for project; project kind error

Turbo refuses a PO of the project expenditure kind that has no job. It shows 'a job is required for the project expenditure kind'.

## PO vendor must exist and be Active

Lines 266-293

Asked as: vendor inactive on PO; supplier not found; vendor status not active; can't use this vendor

Turbo refuses a PO whose vendor is not found ('vendor not found') or whose status is not Active ('provided vendor is not active').

## PO total must be greater than zero

Lines 295-307

Asked as: zero total PO; PO amount zero; negative PO total; total too low

Turbo refuses a PO whose total is 0 or less. It shows 'total must be greater than 0'.

## Child PO parent must be an Active Cumulative top-level PO

Lines 309-363

Asked as: child PO parent error; parent PO not active; parent not cumulative; can't find parent PO

A child PO's parent must exist ('parent PO not found'). The parent cannot itself be a child ('parent PO cannot be itself a child'). It must be Active ('parent PO must be active') and Cumulative ('parent PO must be cumulative').

## Child PO inherits parent currency when blank

Lines 365-367

Asked as: child PO currency blank; child takes parent currency; currency copied to child; child PO missing currency

When a child PO has no currency and its parent has one, Turbo copies the parent's currency to the child.

## Only one open child PO per parent

Lines 369-399

Asked as: second open child PO; already have child PO; child PO already open; only one child allowed

Turbo refuses a child PO when another child of the same parent has a status other than Closed or Cancelled. It shows "other child POs that are not 'Closed' or 'Cancelled' already exist". The PO being edited is not counted, so an existing child can be edited.

## Child PO fields must match parent

Lines 401-418

Asked as: child PO field mismatch; child must match parent; parent details differ; child job description different

A child PO must have the same job, payment type, category, description, vendor and kind as its parent. Turbo shows 'field X must match parent PO's X' for the first field that differs.

## Child PO currency must match parent

Lines 420-433

Asked as: child currency different from parent; child PO currency mismatch; child and parent currency; currency doesn't match parent

Turbo refuses a child PO whose currency differs from its parent's currency. It shows 'child PO currency must match parent PO currency'. A blank currency counts as the home currency.

## PO date must be a valid date

Lines 436-448

Asked as: PO date invalid; wrong date format PO; date won't save; PO date missing

Turbo refuses a PO whose date is not a valid YYYY-MM-DD date. It shows 'date must be a valid date'. The date is also required.

## PO owner and approver must be active users

Lines 451-519

Asked as: inactive approver on PO; owner not active user; second approver inactive; staff member not active

Turbo refuses a PO when the owner is not an active user ('the selected staff member is not an active user'). It refuses when the approver is not active ('the selected approver is not an active user'). It refuses when the priority second approver is not active ('The selected priority second approver is not an active user').

## Recurring PO end date rules

Lines 501-509

Asked as: recurring PO end date; end date before start; end date too far; end date not allowed

A Recurring PO requires an end date ('end_date is required for recurring purchase orders'). The end date must not be before the start date ('end date must be after start date'). It must be within the RECURRING_MAX_DAYS limit of the start date ('end date must be within N days of the start date'). A non-recurring PO cannot have an end date ('end_date is not permitted for non-recurring purchase orders').

## Recurring PO requires frequency; others forbid it

Lines 511-516

Asked as: recurring PO frequency; frequency required; frequency not allowed; how often recurring

A Recurring PO requires a frequency ('frequency is required for recurring purchase orders'). A non-recurring PO cannot have one ('frequency is not permitted for non-recurring purchase orders').

## PO description minimum length

Lines 517-517

Asked as: PO description too short; description minimum length; short description error; description needs characters

A PO description must have at least 4 characters. Turbo shows 'must be at least 4 characters'.

## PO total has a maximum

Lines 520-520

Asked as: PO total too high; maximum PO amount; total above limit; PO amount limit

Turbo refuses a PO whose total is above MAX_APPROVAL_TOTAL.

## Child POs must be One-Time

Lines 521-521

Asked as: child PO wrong type; child must be one-time; child type not allowed; child PO type error

Turbo refuses a child PO whose type is not One-Time. It shows 'child POs must be of type One-Time'.

## Division allocation checked only on create or changed job/division

Lines 524-545

Asked as: division not allocated to job; division mismatch on edit; can't edit old PO; legacy allocation problem

When a PO has a job, Turbo checks that the division is allocated to the job only when the PO is created or the job or division is changed. Old POs can be edited for other fields despite legacy allocation mismatches.

## PO job must exist, be Active and match branch

Lines 533-540

Asked as: job not active on PO; invalid job reference; PO branch differs from job; job status error

A PO's job must exist ('invalid job reference') and must have status Active ('Job status must be Active'). The PO's branch must equal the job's branch ('branch must match the selected job').

## Covered within project budget needs a project job

Lines 549-562

Asked as: covered within project budget; project budget confirmation; confirm covered budget; budget confirmation not allowed

Turbo refuses the 'covered within project budget' confirmation when the PO has no job or its job is a proposal. It shows 'covered within project budget may only be confirmed for purchase orders linked to a project job'.

## Approver must be valid for first stage

Lines 570-608

Asked as: approver not valid first stage; wrong first approver; approver not allowed for amount; first-stage approver error

Turbo refuses a PO whose approver is not a first-stage approver for the division, amount and kind, unless the approver is valid via the bypass for the owner. It shows 'selected approver is not valid for first-stage approval'.

## Second approval requires a valid priority second approver

Lines 610-621

Asked as: second approver required; no second approver; priority second approver missing; second approval setup error

When second approval is required, a priority second approver is required ('priority second approver is required when second approval is required'). The approver must be a second-stage approver ('selected priority second approver is not valid for second-stage approval'). Turbo also reports when no first-stage or second-stage approvers are configured.

## Priority second approver cleared when not needed

Lines 622-624

Asked as: second approver still set; remove priority second approver; second approver not needed; clear second approver

When second approval is not required, Turbo clears any priority second approver on the PO.

## Revised PO editable fields list

Lines 630-679

Asked as: can't edit revision fields; revision locked fields; which fields on revision; revised PO field locked

On a revised purchase order (a revision), the creator may change only the total, currency, approver and priority second approver. Turbo also lets the computed fields (approval_total, approval_total_home, occurrences) and status and approval fields (purchaseOrderMeaningfulChangeSkipFields) differ. A change to any other field makes Turbo refuse the edit (lockedRevisionFieldChange). The comparison ignores created and updated.

## Approvals reset only for first-approved Unapproved POs

Lines 630-706

Asked as: approvals reset on edit; edit cleared approval; draft edit resets approval; approval gone after changing PO

When a creator edits a purchase order, Turbo resets its approvals only if the PO is Unapproved, has a first approval (approved) and has no second approval. A draft with no first approval does not reset. Turbo skips the reset on new records and on Active, Closed or Cancelled POs. The edit must change a meaningful field, which excludes approval, rejection, cancel, close, po_number and status fields.

## Approvals reset when a first-approved PO is edited

Lines 681-706

Asked as: approved PO edited again; edit resets first approval; approval restarts after edit; unapproved PO changed

Turbo resets approvals when the creator edits an Unapproved PO that has a first approval but no second approval. It resets only if there are meaningful changes, or if a rejected revision is resubmitted. It clears the approval and second approval, and keeps the newly submitted approver as the next first-stage approver. Edits to drafts with no first approval do not reset anything.

## Rejected revision resubmission restarts approval

Lines 699-703

Asked as: resubmit rejected revision; resubmitted but no change; approver asked again; restart approval on resubmit

Resubmitting a rejected revision restarts its approval even when nothing changed. This applies when the revision is Unapproved, has a first approval and has no second approval. Its first approver is then asked and told again.

## Revision edit emails approver on key changes

Lines 708-721

Asked as: revision edit email approver; approver not emailed on change; email after revision edited; notify approver of change

When a revision with no first approval is edited, Turbo emails the approver if it was rejected, or if the approver, total or currency changes. A blank currency counts as the home currency.

## Revision approval email shows previous total only if different

Lines 723-760

Asked as: revision approval email total; previous total in email; approval email shows old amount; email total differs

The revision approval email goes to the approver and shows the approval total with its currency. It shows the previous total only when it differs from the new one.

## Approval total text uses approval total and currency

Lines 723-727

Asked as: notification approval total text; wrong total in email; email shows occurrence total; currency code in notification

The approval total text in notifications shows the approval total (not the single-occurrence total on recurring POs) with the PO's effective currency code.

## Revision approver email shows previous total only if changed

Lines 729-743

Asked as: approve again email; revision approver email total; previous total blank email; re-approval email amount

The email that asks a revision's approver to approve again (template po_revision_approval_required) shows the previous approval total only when it differs from the new total with its currency. Turbo blanks the previous total when the two texts are equal.

## New PO emails the approver

Lines 762-786

Asked as: approver not emailed new PO; approval email on create; new PO notification; email sent when submitted

When a PO is created with an approver, Turbo sends that approver an approval-required email. A revision gets the revision email instead. Turbo sends it only after the save succeeds.

## No approval email without an approver

Lines 766-786

Asked as: no approver no email; PO without approver email; why no approval email; approval email missing

Turbo sends no approval-required notification when the purchase order has no approver. A revision gets the revision email, comparing against the original's approval total. Other POs get the standard approval-required email, deferred until the save succeeds.

## PO creation and edits need an authenticated user

Lines 795-812

Asked as: authentication required PO; can't save PO logged out; expense editing disabled; not signed in PO

Turbo refuses PO changes when no one is authenticated. It shows 'authentication is required'. It also refuses while expense editing is disabled (checkExpensesEditing).

## Expense editing can block PO changes

Lines 795-797

Asked as: expenses editing blocks PO; PO blocked by expenses; can't change PO expenses locked; expense lock stops PO

Before processing a purchase order create or update, Turbo runs checkExpensesEditing. If that check returns an error, Turbo refuses the action with that error. The rule it applies lives in another file.

## New PO owner must be the signed-in user

Lines 814-827

Asked as: create PO for someone else; owner must be me; uid not my id; PO owner wrong user

When creating a PO, the owner (uid) must be the signed-in user. Otherwise Turbo shows 'uid must be equal to the authenticated user's id'.

## Only the creator can edit a PO

Lines 828-867

Asked as: can't edit someone's PO; only creator edits PO; change PO owner; not the PO creator

Only the creator can edit a PO. Otherwise Turbo shows 'only the creator can edit this purchase order'. The owner cannot be changed ('uid cannot be changed').

## Only Unapproved POs can be edited

Lines 869-880

Asked as: can't edit approved PO; PO not unapproved; edit locked after approval; why PO can't be edited

Turbo refuses to edit a PO whose status is not Unapproved. It shows 'only unapproved purchase orders can be edited'.

## Revised PO locks all but total, currency and approvers

Lines 882-897

Asked as: revised PO other fields; only total can change; revision edit refused; changed fields on revision

On a revised PO, the creator can change only the total, currency and approvers. Other changed fields are refused with 'only the total, currency and approvers can change on a revised purchase order'.

## PO group cannot be set or changed by edit

Lines 900-918

Asked as: set PO group; change PO group; po group not editable; group missing on PO

A PO's group (po_group) cannot be set or changed by a normal create or update. Turbo shows 'po_group can only be set when a purchase order group is created'.

## PO status must stay Unapproved on create or update

Lines 920-934

Asked as: PO status wrong on save; status must be unapproved; set PO active directly; status change refused

Turbo refuses a PO create or update whose status is not Unapproved. It shows 'status must be Unapproved when creating or updating purchase orders'. Other status changes go through the dedicated actions.

## Owner must be allowed to use the PO branch

Lines 945-947

Asked as: can't use PO branch; branch access denied; no permission for branch; user branch not allowed

Turbo checks that the user can use the PO's branch (EnsureUserCanUseRecordBranch) after cleaning, and refuses the save if not.

## Reset clears approvals but keeps submitted approver

Lines 949-957

Asked as: approver kept after reset; approvals cleared on edit; second approver removed on edit; new approver after reset

When Turbo resets approvals on an edit, it clears the first approval, second approval and second approver. It keeps the approver submitted with the edit as the next first-stage approver.

## Duplicate PO attachment is refused

Lines 961-992

Asked as: duplicate attachment PO; file already uploaded; same receipt another PO; attachment already used

Turbo refuses a new attachment whose hash matches that of another PO. It shows 'This file has already been uploaded to another purchase order'. Turbo stores the hash of a new file, and clears the hash when the attachment is removed.

## PO division must be active

Lines 994-999

Asked as: inactive division PO; division not active; division error on PO; can't pick division

Turbo checks that the PO's division is active (EnsureActiveDivision) and refuses the save with a validation error if not.

## Approval emails sent on create, reset or key revision change

Lines 1010-1032

Asked as: when approver gets email; approval email timing; no email after save; approver notification on update

Turbo creates the approver notification when a PO is new. On an update it creates one when approvals were reset or when a revision change needs an approver email. Otherwise it sends none. Delivery waits until the save succeeds.
