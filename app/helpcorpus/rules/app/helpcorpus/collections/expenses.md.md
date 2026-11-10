<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expenses.md","sha256":"d3928bf81b8f04801abfa5f821c9f8a56993db5bb9b0b3eacb7639614fffb38b","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expenses.md

## Who can list or view expenses

Lines 9-25

Asked as: who can see my expense; view expense someone else; expense not showing for approver; list expenses permissions; see submitted expense

A person can list or view an expense when it is theirs (uid) or when they created it (creator). The chosen approver (approver) can see it once it is submitted (submitted = true). A committer (the commit claim) can see it once it is approved. A reporter (the report claim) can see it once it is committed.

## Who can create an expense

Lines 29-52

Asked as: create expense for someone; bookkeeper add expense for employee; cannot add expense for colleague; new expense on account; make expense for another person

A signed-in person can create an expense for themselves (uid is their own id). A bookkeeper (the book_keeper claim) can create one for another person only if it has a purchase order and its payment type (payment_type) is OnAccount or CorporateCreditCard. The creator field must be empty or the caller's own id. The request cannot set pay_period_ending.

## New expense cannot carry rejection, approval or commit data

Lines 54-66

Asked as: expense already approved on create; set approver when adding expense; preset approval or commit on new expense; create expense with rejection; new expense refused approved field

A create request is refused if it sets any rejection field (rejector, rejected, rejection_reason), approval field (approved, approver) or commit field (committed, committer, committed_week_ending).

## Expense vendor must be active

Lines 68-69

Asked as: vendor inactive on expense; expense vendor not allowed; cannot pick this supplier; vendor disabled error expense; vendor not active

When an expense is created or updated with a vendor, the vendor's status must be Active. A blank vendor passes.

## Expense category must belong to the job on create

Lines 71-76

Asked as: category wrong for job; expense category not on job; category does not match project; new expense category error; pick category for job

On create, the category must belong to the job sent in the request. A blank category is also accepted.

## Only the creator can update an expense

Lines 82-90

Asked as: cannot edit expense; change owner of expense; edit someone else's expense; update expense denied; reassign expense to another person

Only the person who created the expense (creator) can update it. The request cannot change the owner (uid), the creator or pay_period_ending.

## Update cannot change workflow fields

Lines 92-107

Asked as: cannot change expense status; edit approved expense; change approver on expense; modify submitted expense; locked expense fields when saving

An update cannot change any rejection field (rejector, rejected, rejection_reason), submitted, any approval field (approved, approver) or any commit field (committed, committer, committed_week_ending). A request can send such a field only if the value equals the stored one.

## Update checks vendor and category against the job

Lines 109-119

Asked as: edit expense category error; change job on expense; update expense vendor refused; moving expense to another job; category not valid after job change

On update, the vendor must be Active if one is sent. The new category must belong to the job: the existing job if the request sends no job, or the new job if it does. A blank category is accepted.

## Who can delete an expense

Lines 124-128

Asked as: cannot delete expense; remove expense button missing; delete submitted expense; get rid of expense; expense delete not allowed

Only the creator (creator) can delete an expense, and only while it is not submitted (submitted = false) and not committed.

## Expense field limits

Lines 135-157

Asked as: expense date format wrong; invalid payment type; card last 4 digits error; rejection reason too short; allowance type not accepted

The date must be in YYYY-MM-DD form. Payment type (payment_type) must be one of OnAccount, Expense, CorporateCreditCard, Allowance, FuelCard, Mileage or PersonalReimbursement. A rejection reason (rejection_reason) needs at least 5 characters. The card digits (cc_last_4_digits) must be exactly 4 digits. Allowance types (allowance_types) can be any of Lodging, Breakfast, Lunch and Dinner.
