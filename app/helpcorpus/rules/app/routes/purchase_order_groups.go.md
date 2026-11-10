<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_order_groups.go","sha256":"810a9c0486ed7e0064301042715c3c02155adc6ab445f93152ea70f6871075fe","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_order_groups.go

## Purchase order group is a template plus a CSV, one PO per staff member

Lines 1-6

Asked as: bulk purchase orders; one template many staff; make POs for whole team; upload csv of purchase orders

A super_po holder (the super_po claim) fills out one purchase order as a template and uploads a CSV with one row per staff member. Each row becomes an ordinary Unapproved purchase order owned by that staff member. All the purchase orders link to one purchase order group. Normal approval then assigns each PO its number.

## Group upload is limited to 1000 rows

Lines 29-29

Asked as: too many rows upload; max rows csv; upload has no rows; csv row maximum

A purchase order group upload must have at least one row and at most 1000 rows (maxPurchaseOrderGroupRows). Turbo shows "The upload has no rows." when empty and "The upload has N rows. The limit is 1000." when too long.

## Only four template fields can come from the CSV

Lines 31-33

Asked as: csv column not allowed; can't upload this column; which columns can I upload; upload wrong column

A CSV column can replace only total, description, division or branch. Any other upload field is refused with "\"x\" can't come from the upload". The email column is always needed and picks the owner.

## Friendlier approver error messages

Lines 37-43

Asked as: approver can't approve amount; second approver error message; confusing approver error; approver error group

For group rows, approver errors read: "The approver can't approve this amount for this kind and division."; "This amount needs a second approver. Choose one in the template."; and "The second approver can't give final approval for this amount."

## Preview offers approvers who can approve every row

Lines 93-143

Asked as: who can approve all rows; approver list preview; which approvers show up; approvers for every row

The preview lists only people who could approve every valid row, by division, amount and kind. A second approver list is offered, and second approval marked required, only when some row needs one. The owner self-approval bypass is not counted.

## Second approver offered only when a row needs one

Lines 138-141

Asked as: when is second approver needed; second approver list empty; second approval required flag; why no second approver

The preview sets SecondApprovalRequired when any row's policy needs a second approval. Only then does it list second approvers (those who could give final approval on every such row).

## Approver lists sorted by name

Lines 145-160

Asked as: approver dropdown order; approvers not sorted; sort approver names; approver list alphabetical

Offered approvers are sorted by surname, then given name, then ID.

## Create refused when any row has errors

Lines 162-170

Asked as: create button does nothing; rows have errors; can't create group with errors; create blocked by errors

Turbo creates nothing if any row has errors. It returns "N rows have errors" (code rows_invalid) with the preview so the rows can be fixed.

## Staff active status comes from admin profile

Lines 190-205

Asked as: staff marked inactive; why is staff not active; active status admin profile; staff active setting

A staff member counts as active only if their admin profile is marked active. Staff with no admin profile count as not active.

## Staff with no admin profile count as inactive

Lines 190-201

Asked as: staff member not active error; new employee not active; missing admin profile; employee shows inactive

Turbo treats a staff member as active only when their admin profile (admin_profiles) is active. A staff member with no admin profile counts as not active (active defaults to 0). The row then fails with 'This staff member is not active.'

## Row email matching ignores case and spaces

Lines 203-205

Asked as: email case sensitive; uppercase email not found; spaces in email csv; email not matching csv

Turbo looks up each CSV row's email after trimming it and lowercasing it. The staff list it matches against is also trimmed and lowercased (staffByEmail). A differently-cased email still finds the same staff member.

## Division codes and branch names are matched ignoring case

Lines 214-231

Asked as: division code case; branch name lowercase; branch not found typo case; branch name case mismatch

Turbo matches a division by its code, ignoring case. It matches a branch by its name or its code, ignoring case. This applies to template values and to uploaded values.

## Only super_po holders can preview or create PO groups

Lines 236-257

Asked as: can't create purchase orders for others; need super po to preview; no permission create group; who can make group POs

Previewing or creating a purchase order group requires the super_po claim. Without it Turbo refuses with "creating purchase orders for other staff requires the super_po claim" (shown to the user as "you do not have permission to create purchase orders for other staff").

## Group creation refused for a missing claim

Lines 236-255

Asked as: no permission create group; unable to check permissions; missing super po claim; can't make group POs

A person without the super_po claim gets 'you do not have permission to create purchase orders for other staff'. If Turbo cannot check the claim, it says 'unable to check purchase order group permissions'.

## Admins can also view PO groups

Lines 259-267

Asked as: admin view PO groups; see group list as admin; admin can't create groups; admins look at groups

Viewing the list of purchase order groups or one group is open to admin holders (the admin claim) as well as super_po holders. Anyone else gets the super_po refusal. Creating and previewing are still super_po only.

## Template errors are all returned together

Lines 272-277

Asked as: template errors all at once; only first error shown; invalid template message; see every template problem

Turbo collects all template errors and returns them together as 'purchase order group template is invalid'. When one field has several errors, only the first is kept.

## Group ID format required on create

Lines 279-281

Asked as: group id rules; group id 15 characters; invalid group id; group id format wrong

Creating a group requires a group_id of exactly 15 lowercase letters or digits. Otherwise Turbo says "group_id must be 15 lowercase letters or digits". Preview does not check it.

## Group IDs and CSV columns trimmed and lowercased

Lines 279-340

Asked as: group id uppercase error; csv headers case; spaces in csv values; group id letters digits

Turbo trims and lowercases CSV column names and trims every value before it checks the rows. On create, the group ID must be 15 lowercase letters or digits. Otherwise Turbo shows "group_id must be 15 lowercase letters or digits".

## Invalid upload field refused

Lines 283-290

Asked as: field can't come from upload; csv column refused; upload field not allowed; column rejected upload

Listing an upload field other than total, description, division or branch gives the error "\"field\" can't come from the upload".

## Template type must be One-Time or Cumulative

Lines 299-301

Asked as: type must be one-time; cumulative type error; PO type wrong; template type invalid

The template's type must be One-Time or Cumulative. Otherwise Turbo says "Type must be One-Time or Cumulative."

## Template required fields

Lines 302-315

Asked as: kind required; vendor required error; payment required group; date missing template

The template must have kind, payment type, date and vendor (messages like "Kind is required.", "Payment is required.", "Date is required.", "Vendor is required."). Approver is required only when creating, not when previewing.

## Template approver required only on create

Lines 311-315

Asked as: approver required when creating; no approver preview; approver missing create; why need approver

The template approver is required only when the group is created, with the message 'Approver is required.' The preview does not need one. The page picks the approver from the preview's list.

## Template fields replaced by the upload are not required

Lines 316-327

Asked as: branch required with csv; do I need description in template; total required with csv; template fields csv supplies

Branch, division and description are required in the template only if the CSV does not supply them. Total must be greater than 0 in the template only if the CSV does not supply total ("Total must be greater than 0.").

## Total must be greater than 0 when it comes from the template

Lines 325-327

Asked as: zero total refused; total must be above zero; template total zero; total greater than zero

When the total is not uploaded, the template total must be greater than 0. Otherwise Turbo refuses with 'Total must be greater than 0.' A zero total is refused.

## Group upload row limit checks

Lines 329-333

Asked as: no rows in upload; upload too long; over 1000 rows; empty upload file

A purchase order group upload must have at least one row and at most 1000 rows (maxPurchaseOrderGroupRows). Turbo shows "The upload has no rows." when empty and "The upload has N rows. The limit is 1000." when too long.

## CSV column names and values are trimmed and lowercased

Lines 334-340

Asked as: csv headings case sensitive; column names spaces; extra spaces csv; header case csv

Turbo trims spaces from every CSV value and trims and lowercases every column name before checking rows, so column headings are not case-sensitive.

## Uploaded total parsing

Lines 400-406

Asked as: dollar sign in csv total; total not a number; commas in amount csv; amount parsing upload

When total comes from the CSV, Turbo removes $, commas and spaces before reading the number. An empty value gives "Total is required." and a non-number gives "Total \"x\" is not a number."

## Row email checks

Lines 429-445

Asked as: email is required csv; email appears more than once; no staff with this email; duplicate email upload

Each row needs an email that matches one staff member. Errors: "Email is required.", "This email appears N times in the upload." for duplicates, "No staff member has this email.", and "This staff member is not active." Email matching ignores case and surrounding spaces.

## Row errors stop that row from creating a purchase order

Lines 429-441

Asked as: row with error no PO; duplicate email row skipped; one bad row; bad row skipped group

A row with an error produces no purchase order and counts toward the preview's error count. Rows whose email is empty, duplicated, unknown or inactive are refused before any PO validation runs. An email that appears more than once fails with 'This email appears N times in the upload.'

## Uploaded total, description, division and branch row messages

Lines 451-487

Asked as: total is required row; unknown division code; branch not found row; description missing row

A row with a missing total shows "Total is required." A total that is not a number shows 'Total "x" is not a number.' A blank description shows "Description is required." An unknown division shows 'No division has code "x".' An unknown branch shows 'No branch has the name or code "x".'

## Uploaded total errors

Lines 454-462

Asked as: csv total empty; total not a number csv; amount error upload; total blank upload

When total comes from the CSV, an empty value gives "Total is required." and a non-number gives "Total \"x\" is not a number."

## Uploaded description required

Lines 463-467

Asked as: description blank upload; csv missing description; description required row; no description csv

When description comes from the CSV, an empty value on a row gives "Description is required."

## Uploaded division matched by code

Lines 468-476

Asked as: division code csv; wrong division upload; division not found; division column upload

When division comes from the CSV, Turbo matches the division code ignoring case. Empty gives "Division is required." and an unknown code gives "No division has code \"x\"."

## Uploaded branch matched by name or code

Lines 477-485

Asked as: branch name csv; branch not found upload; branch code csv; branch column upload

When branch comes from the CSV, Turbo matches the branch name or code, ignoring case. Empty gives "Branch is required." and no match gives "No branch has the name or code \"x\"."

## Rows are validated like a normal PO

Lines 496-528

Asked as: row validation like normal PO; division inactive group; branch claim error; group PO checks

Each row is built as an Unapproved purchase order owned by the row's staff member and goes through the normal cleaning and validation. Branch is checked against the creator's branch claims, not the owner's, and is never filled from the owner's default. The division must be active.

## Created group purchase orders are Unapproved and owned by the staff member

Lines 496-510

Asked as: who owns group POs; POs status after create; po linked to group; group POs owner

Each row creates a purchase order with status Unapproved. The staff member in the row owns it. The approver and priority second approver come from the template, and the PO is linked to the group (po_group).

## Row branch is checked against the creator's branch claims

Lines 519-523

Asked as: branch access creator; branch default owner; can't use branch for staff; branch from owner default

Turbo checks each row's branch against the creator's branch access, not the owner's. It is never filled from the owner's default branch. The row's division must also be active, and the PO must pass normal PO validation.

## Row approver errors hidden until approver chosen

Lines 530-533

Asked as: approver errors on rows; why no approver error; rows no approver yet; approver error before choosing

While the template has no approver, rows are not faulted for the approver or the priority second approver.

## Preview totals and missing staff

Lines 544-598

Asked as: missing staff in upload; who is not in csv; preview totals count; staff left out upload

The preview gives row count, error count, the sum of row totals, and the currency code. It also lists active staff who do not appear in the upload (missing_active_staff).

## Template currency sets the currency code of a group

Lines 545-549

Asked as: currency code group; which currency group; currency from template; group currency wrong

The preview shows the currency code Turbo derives from the template currency (EffectiveCurrencyCode). Every purchase order in the group gets the template's currency and vendor.

## Preview does not stop on errors; create returns rows_invalid

Lines 623-644

Asked as: rows_invalid code; preview with errors; error response group; create error code

A create with errors in any row returns the code "rows_invalid" and the preview. Any other error in a group request is returned with its own status and message.

## Preview needs no approver and no group id

Lines 646-659

Asked as: preview without approver; check group before creating; dry run group; preview no group id

The preview checks the same template and rows as create, except that group_id and approver are not required. It makes no records.

## Group approver gets one email for the whole group

Lines 663-694

Asked as: one email group approval; approver email for group; how many approval emails; group approval email

On creation, the chosen approver receives one "po_group_approval_required" email for the group, with count, description, vendor and total, instead of one per PO. If the description varies by row it reads "Varies by staff member". A failure to send is logged and does not undo the creation.

## Group email description falls back to a placeholder

Lines 663-670

Asked as: varies by staff member; description blank email; email description missing; placeholder description email

When the group has no template description (it came from the upload), the approval emails show 'Varies by staff member'.

## Second approver gets one email per group

Lines 696-739

Asked as: second approver email; priority second approver notification; how often second approver emailed; second approval notice

The priority second approver gets one "po_group_second_approval_required" email per group, sent immediately when the first of its POs is first-approved and needs second approval. If that recipient already got one for the group, no more are sent.

## Second approver email is sent once per group

Lines 700-737

Asked as: duplicate second approval email; second approver gets many emails; email sent once group; repeat second approval notice

The priority second approver gets one 'po_group_second_approval_required' email per group. Turbo skips it if that recipient already received one for the group. It is sent immediately.

## Repeated group creation refused

Lines 758-766

Asked as: group already created; duplicate group id; create same group twice; repeat group creation

A group_id that already exists is refused with "this purchase order group has already been created". The check runs before and again inside the transaction, so a repeated request creates nothing.

## Duplicate group ids are rejected with a conflict

Lines 758-763

Asked as: group id conflict; already exists group; group exists error; group id taken

Turbo checks whether the group already exists both before and inside the transaction. A repeat returns 'this purchase order group has already been created'.

## Create refused when any row has errors (check)

Lines 772-774

Asked as: error count stops create; create stops preview; errors preview returned; rows error create stop

After evaluating the rows, create stops and returns the preview if the error count is above zero.

## Group creation is all or nothing

Lines 776-807

Asked as: partial group created; transaction group POs; rollback group; all or nothing group

Turbo creates the group record and all its Unapproved purchase orders in one transaction, so either every row is created or none. Each PO is linked to the group, and the group records the creator. Numbers are assigned later at approval.

## Create saves group, purchase orders and approval email together

Lines 777-807

Asked as: group saved fail; notification saved together; who created group; creator of group

Turbo saves the group, every purchase order and the approval notification in one transaction. If any save fails, nothing is created. The group's creator is the signed-in user.

## Group description stored only when from template

Lines 791-793

Asked as: group description blank; description empty upload group; where description stored; group description missing

The group stores the template description only when description is not supplied by the CSV. Otherwise the group's description is blank.

## Group approval email sending failure ignored

Lines 812-816

Asked as: email failed group still created; approval email error; email send failure; group created but no email

The approval email is sent after the transaction commits. If sending fails, Turbo only logs the error and still reports the group as created.

## Successful create returns group id, count and total

Lines 818-823

Asked as: create response group; what does create return; group id returned; success group response

A successful create returns the group id, the number of purchase orders, the total and the currency code.

## Group summary shows counts, total and defaults

Lines 840-864

Asked as: group summary; how many POs unapproved; group totals view; group details counts

A group summary shows creator, description, vendor, PO count, count still Unapproved, total, and currency. Currency defaults to CAD when none is found. The description is blank when it came from the upload.

## Group summary defaults

Lines 844-864

Asked as: currency CAD default; blank vendor group; groups newest first; group list order

A group's summary uses a blank description when none is stored. It also uses a blank vendor name when there is none, and the currency code defaults to CAD. It counts all the POs in the group and those still Unapproved, and it sums the totals. Groups are listed newest first.

## Group list is open to admins and super_po holders

Lines 866-869

Asked as: who can see groups; group list permission; can't open group list; group list access

The group list and the single-group view need the admin claim or the super_po claim. Other users get the permission error.

## Group list and single group view both need admin or super_po

Lines 866-883

Asked as: view one group access; open group denied; group page permission; single group access

Only an admin or a super_po holder can list groups or open one group.

## Missing group returns not found

Lines 884-897

Asked as: group not found; empty group id; can't find group; missing group page

Viewing one group that does not exist, or has an empty id, gives "purchase order group not found".
