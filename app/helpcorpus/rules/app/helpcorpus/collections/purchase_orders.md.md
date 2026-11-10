<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/purchase_orders.md","sha256":"4144fa1240beb2443c0202f7e0639acf413e1be023c90539898615f1a9e40b49","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/purchase_orders.md

## Who can list or view purchase orders

Lines 10-38

Asked as: can't see purchase order; PO missing from list; why can't I see this PO; PO not showing up

A signed-in person sees an Active purchase order. They see an Unapproved one only when it has a PO number, is not a legacy manual entry (legacy_manual_entry) and is not rejected. A person also sees a Cancelled or Closed one only if they are its creator, approver, second approver, the job's manager or alternate manager, or hold the report claim. The same filter applies to listing and to viewing one record.

## Unapproved purchase orders visible to people involved

Lines 23-31

Asked as: rejected PO still showing; pending PO visible to me; why do I see unapproved PO; unapproved order showing for me

An Unapproved purchase order is visible to its creator (uid), its approver and its rejector. It is also visible to the priority second approver, but only after the first approval (approved is set) and before the second approval is given. This lets a rejected or pending PO show to those involved. The same rule applies to listing and viewing.

## Super PO and admin see grouped purchase orders

Lines 32-38

Asked as: see grouped POs; PO group not showing; super PO access; admin can't see group orders

A person with the super_po or admin claim can list and view any purchase order that belongs to a PO group (po_group is not empty), whatever its status.

## Who can view one purchase order

Lines 44-72

Asked as: cannot open purchase order; PO page blank; access denied on PO; can't view this order

The view rule is the same filter as the list rule. Any signed-in person can view an Active purchase order. Any signed-in person can view an Unapproved purchase order that has a po_number, is not a legacy manual entry (legacy_manual_entry) and has not been rejected. The creator, the approver, the rejector and others involved can view the other Unapproved ones. A Cancelled or Closed purchase order is visible to the creator, the approver, the second approver, the job's manager or alternate manager, and a reporter (the report claim). Super PO (super_po) and admin claim holders can also view any purchase order that has a po_group.

## Cancelled and closed purchase orders have restricted visibility

Lines 46-56

Asked as: closed PO hidden; cancelled PO not visible; can't find closed order; old cancelled order missing

A Cancelled or Closed purchase order is not visible to every signed-in person. Only the creator, the approver, the second approver, the job's manager or alternate manager, or a person with the report claim can see it.

## Second approver can see a half-approved unapproved PO

Lines 57-65

Asked as: second approver can't see PO; waiting for second sign off; PO awaiting second approval visibility; half approved PO not showing

A priority second approver (priority_second_approver) can see an Unapproved purchase order only after the first approval (approved is set) and before the second approval (second_approval is empty). This holds for list and view.

## Who can create a purchase order

Lines 78-116

Asked as: make new PO; raise purchase order; who can create a PO; start a new order

A signed-in person creates a purchase order only as themselves: uid must be missing or be their own id. The status must be "Unapproved". The request must not set a PO number, legacy_manual_entry as true, rejection fields, approved, second approver fields, cancellation fields or closing fields (closed, closer, closed_by_system). The hooks check that the approver has the right claim and divisions; vendor status is also checked in hooks, not in this rule.

## Created purchase order must start Unapproved with no PO number

Lines 79-88

Asked as: new PO already active; can't set PO number on create; create PO with number; new order starts wrong status

A signed-in person creates a purchase order only if status is 'Unapproved'. The request must not send a po_number, or must send it empty. It may not set legacy_manual_entry to true.

## Creator is the signed-in user on create

Lines 90-91

Asked as: create PO for colleague; raise order on behalf of; make PO under someone else; PO under another name

On create, uid may be left out or must equal the signed-in person's id. A person cannot create a purchase order for someone else.

## Create request cannot set approval, rejection, cancel or close fields

Lines 93-116

Asked as: set approver on new PO; pre-approve PO; mark PO approved when creating; new PO rejected field error

A create request must not send rejector, rejected, rejection_reason, approved, second_approver, second_approval, cancelled, canceller, closed, closer or closed_by_system. Hooks check that the approver has the right claim and divisions, and the vendor status check is also done in the hook, not in the rule.

## Category must belong to the job on create

Lines 118-123

Asked as: category wrong for job; category not showing for job; job category error; category doesn't match job

When a purchase order is created, a category must belong to the job in the same request, or the category must be blank. A category without a job is refused.

## Who can update a purchase order

Lines 129-140

Asked as: edit purchase order; change PO after approval; can't edit my PO; update order details

Only the creator (uid) can update a purchase order through the API. The PO must be Unapproved, have no second approval yet, and not be a legacy manual entry. The uid cannot change.

## Fields a creator cannot change on update

Lines 142-164

Asked as: change PO number; edit approver on PO; update locked PO fields; PO total won't change

An update must not change the PO number, the rejection fields, approved, the second approver fields, the cancellation fields or the closing fields (closed, closer, closed_by_system). A request may include these fields only with their current values. legacy_manual_entry cannot be set to true.

## Category must match job on update

Lines 169-176

Asked as: change job on PO; category error when editing; moving PO to other job; edit category on order

On update, the new category must belong to the PO's current job when the job is not sent. If a new job is sent, the category must belong to that job. A blank category is always allowed.

## Who can delete a purchase order

Lines 182-182

Asked as: remove purchase order; delete PO; get rid of draft order; can't delete my PO

Only the creator can delete a purchase order through the API. The PO must be Unapproved and have no PO number yet.

## PO number format

Lines 190-190

Asked as: PO number invalid; wrong PO number format; PO number suffix; PO number won't save

A PO number, when present, looks like 2025-0123, with a four-digit year not starting with 0, a dash and four digits. It can end with a dash and a two-digit suffix from 01 to 99, such as 2025-0123-01.

## Rejection reason length and allowed values

Lines 191-204

Asked as: rejection note too short; reason won't save; reject comment minimum length; reject message error

A rejection reason (rejection_reason) needs at least 5 characters. Type is One-Time, Cumulative or Recurring. Payment type is OnAccount, Expense or CorporateCreditCard. Frequency is Weekly, Biweekly or Monthly. Status is Unapproved, Active, Cancelled or Closed.

## Allowed values for PO status, type, frequency and payment type

Lines 191-200

Asked as: dropdown options PO; invalid status on PO; payment type options; allowed PO types

Status is Unapproved, Active, Cancelled or Closed. Type is One-Time, Cumulative or Recurring. Frequency, when given, is Weekly, Biweekly or Monthly. Payment type is OnAccount, Expense or CorporateCreditCard.

## Required fields on a purchase order

Lines 191-221

Asked as: PO won't save missing field; required fields purchase order; mandatory PO fields; which fields are needed on PO

A purchase order requires status, uid, type, date, division, payment_type, approver and approval_total.

## Date format on purchase orders

Lines 194-195

Asked as: date format PO error; end date wrong format; PO date invalid; how to enter PO date

The date and the end date (end_date) must be in the form YYYY-MM-DD. The date is required; the end date is optional.

## Attachment limits

Lines 201-201

Asked as: upload PO file too big; attach PDF to PO; how many attachments on PO; photo upload purchase order

A purchase order has at most one attachment, up to 5242880 bytes (5 MB). It must be a PDF, JPEG, PNG or HEIC file.

## Rejection reason needs at least 5 characters

Lines 204-204

Asked as: reject reason too short; short rejection note; rejection text minimum; reason needs more letters

A rejection_reason, when given, must have at least 5 characters.

## Occurrences must be a whole number

Lines 230-230

Asked as: occurrences decimal error; occurrences not accepting number; occurrences must be whole; occurrence count invalid

The occurrences field of a purchase order accepts whole numbers only.

## PO numbers and attachments must be unique

Lines 237-239

Asked as: duplicate PO number; same PO twice; file already attached elsewhere; attachment already used on another PO

Two purchase orders cannot have the same PO number, unless it is blank. Two cannot have the same attachment hash (attachment_hash), unless it is blank, so one attachment file cannot be used on two POs.
