<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/purchase_orders_augmented.md","sha256":"f59ec30bcafd4edaa3ce353dbe5126707cd88c47c283f19c0618b960519dbbc5","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/purchase_orders_augmented.md

## Purchase orders view is read-only through the API

Lines 5-5

Asked as: edit purchase order api; can't change PO through api; delete purchase order api; write access purchase orders

purchase_orders_augmented is a view that Turbo reads and never writes. The file sets no create, update or delete rules, so nobody can create, change or delete records through the API. Turbo's own server code is not limited by these rules.

## Who can list or view Active purchase orders

Lines 10-10

Asked as: can't see open PO; anyone see active purchase order; PO not showing when signed in; purchase orders visible to staff

Any signed-in person can list and view a purchase order (PO) whose status is "Active" in purchase_orders_augmented. Someone who is not signed in cannot.

## Who can list or view Cancelled or Closed purchase orders

Lines 12-21

Asked as: can't see closed PO; cancelled purchase order missing; who can view cancelled POs; closed PO hidden from me

A PO with status "Cancelled" or "Closed" is visible only to certain people. These are the PO's creator (uid), its approver, its second approver, the job's manager, the job's alternate manager, and anyone with the report claim. Other signed-in staff cannot list or view it.

## Who can list or view Unapproved purchase orders

Lines 23-30

Asked as: pending PO not showing; unapproved PO hidden; waiting approval PO invisible; draft purchase order visibility

A PO with status "Unapproved" is visible to its creator (uid), its approver and its rejector. The priority second approver (priority_second_approver) can also see it, but only once the PO has a first approval (approved) and no second approval (second_approval) yet. Nobody else can list or view it.

## View rule matches list rule for POs

Lines 36-56

Asked as: PO opens but missing from list; single PO differs from list; purchase order not found link; list and detail PO mismatch

Viewing one record in purchase_orders_augmented uses the same filter as listing: Active POs for any signed-in person, Cancelled or Closed POs for the people involved, report claim holders and the job managers, and Unapproved POs for the creator, approver, rejector, and the priority second approver after the first approval.

## Purchase order view shows approval and rejection details

Lines 76-95

Asked as: who approved this PO; approval date purchase order; why was PO rejected; second approver name shown

The purchase orders view carries the approver and the approved date, and the second approver and the second approval date. It also carries the rejector, the rejected date and the rejection reason. It carries the priority second approver (priority_second_approver). The list and view access filters use these fields.

## Purchase order view shows who closed or cancelled it

Lines 83-94

Asked as: who cancelled the PO; who closed purchase order; closed by system flag; date PO was cancelled

The purchase orders view carries the canceller and cancelled date (canceller, cancelled). It also carries the closer, the closed date and a flag for a close by the system (closer, closed, closed_by_system). The values are copied from the purchase order.

## Purchase order view includes approval total

Lines 93-96

Asked as: PO total amount; approval total missing; how much is PO worth; PO value field blank

The purchase orders view (purchase_orders_augmented) carries the purchase order's approval total (approval_total) next to its other fields. The field is required. It also carries the committer-related fields closer, closed, closed_by_system and priority_second_approver.

## View adds committed expense count and names

Lines 97-124

Asked as: expense count on PO; creator name on purchase order; vendor shown on PO view; approver names purchase order view

The view adds to each PO the number of its expenses that are committed (committed_expenses_count, counting expenses whose committed value is not empty). It also adds the names of the creator, approver, second approver, priority second approver and rejector, plus the parent PO number, vendor name and alias, job number and description, client name, division code and name, and category name.

## Committed expenses count on purchase order view

Lines 97-97

Asked as: committed expenses number wrong; uncommitted expenses counted on PO; how many expenses committed; expense count incorrect on PO

Each purchase order in the view shows committed_expenses_count: the number of expenses linked to that PO whose committed value is not empty. Expenses that are not committed are not counted.

## Names shown for people on purchase order view

Lines 98-102

Asked as: full name on purchase order; approver name blank on PO; person name missing on PO; names shown for PO people

The view shows full names (given name, a space, surname) from profiles for the creator (uid_name), approver, second approver, priority second approver and rejector. A person with no profile or no role in the PO has no name.

## Related record details on purchase order view

Lines 103-124

Asked as: PO missing vendor name; purchase order without job shows; vendor name blank on PO; parent PO empty on view

The view adds the parent PO number, vendor name and alias, job number and description, client name and id, division code and name, and category name. Each is a left join, so a PO without a parent, vendor, job or category still appears with those values empty.

## Client of a purchase order comes through its job

Lines 107-124

Asked as: client blank on purchase order; which client is on PO; PO with no job client; client name missing on PO

The client name and id on the view come from the client of the PO's job. A PO with no job has no client name or id.

## Required fields on purchase order view

Lines 131-164

Asked as: mandatory PO fields; PO won't save missing approver; required purchase order details; what fields PO needs

A purchase order must have a status, creator (uid), type, date, division, payment type and approver. The approval total is also required.

## PO number format and rejection reason length

Lines 132-146

Asked as: PO number format wrong; why PO number rejected; rejection reason too short; child PO suffix format

A PO number (po_number) looks like 2025-0001, with an optional suffix of 01 to 99 for a child PO. A rejection reason must have at least 5 characters.

## Allowed PO status, type, frequency and payment type values

Lines 133-142

Asked as: valid PO status options; payment type choices PO; PO frequency weekly monthly; allowed PO types list

A PO's status is one of Unapproved, Active, Cancelled or Closed. Its type is One-Time, Cumulative or Recurring. Its frequency, if set, is Weekly, Biweekly or Monthly. Its payment type is OnAccount, Expense or CorporateCreditCard.

## Dates in purchase orders use YYYY-MM-DD

Lines 136-137

Asked as: date format on purchase order; end date invalid PO; wrong date style PO; PO date won't accept

The date and end_date fields must match the pattern YYYY-MM-DD (four digits, two digits, two digits). The date is required. The end date is optional.

## Attachment limits on purchase orders

Lines 143-143

Asked as: upload PO file too big; attach PDF to purchase order; file type not allowed on PO; PO attachment size limit

A PO attachment holds one file of at most 5242880 bytes (5 MB). It must be a PDF, JPEG, PNG or HEIC file.

## Parent PO number uses the PO number format

Lines 171-171

Asked as: parent PO number wrong; child PO number format; parent PO invalid format; parent purchase order number rules

The parent PO number (parent_po_number) follows the same pattern as po_number: four digits, a dash, four digits, and an optional dash with a two-digit suffix from 01 to 99.

## Job number pattern on purchase order view

Lines 174-174

Asked as: job number format invalid; job number not accepted; job number suffix format; job number pattern wrong

The job number (job_number) must match: an optional P, two digits, a dash, three or four digits, an optional L, then optional one- or two-digit and one-digit suffixes after dashes.
