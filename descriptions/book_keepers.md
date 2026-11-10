# Bookkeepers entering expenses against purchase orders

## Summary

This feature is built. Turbo used to assume the person entering an expense is also the person who incurred the expense. That assumption is not true for a meaningful share of AP work: bookkeepers often receive vendor invoices against existing purchase orders and enter the resulting expenses on behalf of the PO owner.

The old model routed those expenses to the bookkeeper's manager because `expenses.uid` was set to the person creating the record, and the expense hook derived `expenses.approver` from that user's manager. That created a large amount of approval busywork that does not reflect the real control point. The PO has already authorized the spend; the bookkeeping step is matching an invoice to that active PO.

Turbo now works this way:

- there is a `book_keeper` claim;
- `expenses.creator` is a relation to `users`;
- `expenses.uid` is the person the expense is for;
- the hook sets `expenses.creator` to the authenticated user who created the row;
- a narrow bookkeeping-on-behalf flow allows `creator != uid`;
- for that flow, the hook sets `expenses.approver = expenses.creator`;
- the bookkeeper can approve only those expense records where they are the assigned approver.

This is intentionally limited to the green plus flow: creating an expense from an existing active `purchase_orders` record.

Document reuse for bookkeeper-entered PO expenses is specified separately in
`expense_documents.md`. This spec defines the actor, ownership, visibility, and
approval model for bookkeeping-on-behalf; it does not define the attachment
storage, duplicate-file, source-expense reuse, or legacy attachment backfill
model.

## Status quo

This section describes Turbo before the bookkeeper feature, with notes on what changed.

### Expense ownership

Before the bookkeeper feature, `expenses.uid` meant all of the following:

- the person the expense belongs to;
- the person who created the expense;
- the person allowed to edit draft expenses;
- the person allowed to submit, delete, and recall the expense through the current owner flows;
- the profile used to derive the manager approver.

`ProcessExpense` used to enforce this at create time: `expenses.uid` had to equal the authenticated user's id, and on update the original `uid` had to equal the authenticated user. Now `creator` holds the person who created the expense, only the `creator` can edit, submit, delete and recall it, and `uid` may differ from the creator in the bookkeeper flow below. `uid` still cannot be changed.

### Expense approval

Expense approval is manager based for regular expenses, and this still holds:

- `cleanExpense` loads the profile for `expenses.uid`;
- it reads `profiles.manager`;
- it validates that manager is active and has the `tapr` claim;
- it writes that manager to `expenses.approver`.

The approval route is then strict and simple: the authenticated user can approve only when `record.approver == auth.id`, the record is submitted, the record is not committed, and the record is not already approved.

### Purchase-order expense creation

There are two creation paths:

- `/expenses/add` creates an expense from scratch and cannot specify a PO number.
- `/expenses/add/{poid}` creates an expense from an existing PO, pre-filling values from the linked `purchase_orders` record.

For PO-linked expenses:

- the PO must exist and be `Active`;
- the expense `kind`, branch, currency, job, and other PO-derived values are normalized against the PO;
- current policy allows any authenticated user to submit an expense against any active PO, subject to validation.

The UI used to set `item.uid` to the authenticated user before every save, so a bookkeeper using the green plus icon created the expense as themself. It now sends the PO owner's id as `uid` in the bookkeeper on-behalf case described below.

## Problem

Before this feature, when bookkeeping entered an invoice against someone else's PO:

- the expense was recorded under the bookkeeper's `uid`, not the PO owner's `uid`;
- the bookkeeper's manager became the expense approver;
- the manager received a queue of invoice-matching approvals that were really bookkeeping work;
- the expense did not clearly record who entered it versus who the purchase belonged to.

The underlying PO already carries the staff member, vendor, job/category, amount, kind, branch, and approval state. For invoice entry against an existing active PO, the bookkeeping user is not asking their own manager to approve their spending. They are recording that an approved purchase has produced an invoice.

## Desired behavior

Turbo does everything in this section today.

### Flow comparison

Both PO-linked paths use the same green plus UI and keep the same baseline PO-linked validation. The table below calls out only the differences between the regular path and the narrow bookkeeper on-behalf path.

| Area                               | Regular PO-linked flow                                                                                                    | Bookkeeper on-behalf flow                                                                                                     |
|------------------------------------|---------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| Who can use it                     | Any authenticated user.                                                                                                   | Authenticated user with `book_keeper`.                                                                                        |
| When this path is used             | The expense is being created for the authenticated user, including when the user has `book_keeper` or owns the linked PO. | The expense is being created for another user, and that user is the linked PO owner.                                          |
| Expense `uid`                      | The authenticated user.                                                                                                   | The linked PO owner.                                                                                                          |
| Expense `creator`                  | The authenticated user, so `creator = uid = auth.id`.                                                                     | The authenticated bookkeeper, so `creator != uid`.                                                                            |
| Relationship to PO owner           | The expense can be linked to someone else's PO, but it is still the caller's expense.                                     | The expense is for the PO owner, and the bookkeeper is recorded separately as the entering user.                              |
| Extra PO eligibility               | No additional owner or PO-number requirement beyond normal PO-linked expense validation.                                  | PO must have a non-empty `po_number`, and PO owner must differ from the bookkeeper.                                           |
| Payment types                      | Existing expense payment-type rules apply.                                                                                | Only `OnAccount` and `CorporateCreditCard`; the linked PO's payment type does not control bookkeeper eligibility.             |
| Approver                           | Derived from the expense owner's manager, and that manager must be active with `tapr`.                                    | Set to the bookkeeper/creator. The creator must be active with `book_keeper`; `tapr` is not required for this assignment.     |
| Draft edit, delete, submit, recall | The caller controls the draft because `creator = uid`.                                                                    | The bookkeeper controls the draft because `creator != uid`; the PO owner can view when visible but is not the draft operator. |
| Approval queue                     | Submitted expense appears for the manager-derived approver.                                                               | Submitted expense appears for the bookkeeper because `approver = creator`.                                                    |
| Same-user bookkeeper PO            | A bookkeeper creating against their own PO stays on this regular path.                                                    | Not applicable; same-user POs do not use the on-behalf path and do not create self-approval.                                  |

### Path selection

There is no separate button or toggle for the bookkeeper flow. The two PO-linked flows are selected by the `uid` that reaches the expense create hook.

Regular PO-linked flow is active when the create payload has:

- `uid = auth.id`.

This remains true even when the linked PO belongs to another user. In that case, the user is creating their own expense against another person's PO.

Bookkeeper on-behalf flow is attempted only when the create request has:

- `uid != auth.id`;
- `purchase_order` set to the linked PO.

When this path is accepted, the server sets `creator = auth.id`.

The frontend sends `uid = purchase_orders.uid` only for the narrow green plus case where all of the following are true:

- the user has `book_keeper`;
- the linked PO owner differs from the current user;
- the linked PO has a non-empty `po_number`;
- the selected expense payment type is `OnAccount` or `CorporateCreditCard`.

If any of those conditions is false, the frontend keeps the regular behavior and sends `uid = auth.id`. The server then treats the create as a regular PO-linked expense.

The server is still authoritative. If any client submits `uid != auth.id`, the hook accepts it only after validating the full bookkeeper eligibility list: `book_keeper` claim, active numbered PO, PO owner match, allowed expense payment type, different PO owner, and existing branch and PO-linked expense validation.

### Definitions

- `uid`: the person the expense is for. For the bookkeeper flow, this must be the linked PO's `uid`.
- `creator`: the authenticated user who created the expense row.
- `approver`: the user responsible for approving the expense row.
- Bookkeeper-entered expense: an expense where `creator != uid`.

### New claim

Turbo has a claim named `book_keeper`.

The claim grants the narrow ability to create an eligible PO-linked expense on behalf of the PO owner. It does not grant broad approval powers, commit powers, report access, payables admin access, or the ability to create arbitrary expenses for arbitrary users.

### New `expenses.creator` field

The `expenses` collection has a `creator` field:

- type: relation to `users`;
- max select: 1;
- set by the hook to the authenticated user on every new row, whatever the request sends;
- backfilled on existing rows as `creator = uid`;
- immutable after create.

The field is not named `created` because PocketBase already has the system `created` timestamp. The actor field is `creator`.

### Eligibility for `creator != uid`

The system may accept `creator != uid` only when all conditions are true:

- authenticated user has the `book_keeper` claim;
- the expense is being created, not updated;
- the request is using the PO-linked green plus flow, meaning `purchase_order` is present;
- the linked `purchase_orders` record exists;
- the linked PO is `Active`;
- the linked PO has a non-empty `po_number`;
- the linked PO's `uid` is not the authenticated user;
- submitted expense `uid` equals the linked PO's `uid`;
- expense `payment_type` is `CorporateCreditCard` or `OnAccount`;
- the caller is allowed to use the linked PO's branch under existing branch-claim rules.

If any of these conditions are not met, the existing rule applies: `uid` must equal the authenticated user. When the creator later edits a bookkeeper-entered draft, the hook checks the same list again, so the PO must still be `Active`.

### Same-user bookkeeper case

If a user with the `book_keeper` claim creates an expense against their own PO, regular rules apply.

That means:

- `uid = creator = auth.id`;
- approver is derived from the user's manager;
- the record is not self-approved merely because the user holds `book_keeper`.

### Approval assignment

For normal expenses:

- `creator = auth.id`;
- `uid = auth.id`;
- `approver = profiles.manager` for `uid`;
- approver must be active and hold `tapr`.

For eligible bookkeeper-entered expenses:

- `creator = auth.id`;
- `uid = purchase_orders.uid`;
- `approver = creator`;
- approver must be active and hold `book_keeper`;
- `tapr` is not required for this special creator-as-approver assignment.

Bookkeepers do not gain approval authority over all AP expenses. They can approve only rows where `expenses.approver` is their own user id. For this flow, that is true only because the create hook assigned them as the approver on an eligible `creator != uid` record.

### Draft lifecycle

When `creator != uid`, the creator is the operational owner of the draft.

The creator can:

- see the draft in their expense list;
- edit the draft while it is unsubmitted;
- submit the draft;
- delete the draft while it is unsubmitted;
- recall it when recall is otherwise allowed;
- see it in their pending approval queue after submission;
- approve it when the normal approval-route state checks pass.

The `uid` user remains visible as the person the expense is for, but is not the person responsible for editing an invoice that bookkeeping entered. The `uid` user can view the expense, but cannot edit, submit, recall or delete it, and cannot approve it unless they are also the assigned approver under existing rules.

### Queues and visibility

The pending approval queue is driven by:

```sql
expenses.approver = auth.id
AND expenses.submitted = true
AND expenses.approved is blank
```

Because bookkeeper-entered rows set `approver = creator`, submitted rows where `creator != uid` naturally appear in the creator's approval queue.

The "my expenses" list and expense details visibility account for `creator`:

- both show rows where `uid = auth.id`;
- both show rows where `creator = auth.id`;
- expense details also show rows where `approver = auth.id` and the row is submitted;
- expense details keep the existing commit/report visibility behavior.

This prevents a bookkeeping-entered draft from disappearing immediately after creation.

### Audit display

Expense details and list rows show both people when they differ:

- "Submitted for": `uid_name` (details page);
- "Entered by": `creator_name` (details page and list rows);
- the list row also shows `approver_name` once the expense is approved.

When `creator == uid`, the details page keeps the simple "Submitted By" language and list rows show only `uid_name`.

## Non-goals

This change does not:

- let bookkeepers create expenses from scratch for other users;
- let bookkeepers choose any arbitrary `uid`;
- apply to `Expense`, `Allowance`, `FuelCard`, `Mileage`, or `PersonalReimbursement`;
- apply to PO-linked expenses without a PO number;
- apply to inactive, cancelled, closed, or unapproved POs;
- let bookkeepers approve expenses unless they are the assigned approver;
- replace commit/payables/report controls;
- bypass cumulative PO total checks, recurring date checks, attachment rules, settlement rules, or currency rules.

## Implementation plan

Every step in this plan is built. The steps below say what the code does.

### 1. Data migration

Migration `app/migrations/1777381896_book_keeper_expense_creator.go`:

- creates the `book_keeper` claim if it does not already exist;
- adds `creator` to `expenses` as a relation to `users`;
- backfills existing expenses with `creator = uid`;
- adds indexes useful for list and visibility queries:
  - `expenses(creator, date)`;
  - `expenses(creator, submitted, date)`.

The migration is idempotent around the claim insert so environments with a manually created claim do not fail.

### 2. PocketBase collection rules

The `expenses` collection rules recognize the split:

- create rule permits regular creates where `uid` equals auth id;
- create rule permits the shape of a possible bookkeeper-on-behalf create only for authenticated users with `book_keeper`, with `purchase_order` present, `uid != auth.id`, and an allowed payment type;
- update/delete rules permit draft operations by `creator`;
- approval, rejection, commit, and settlement remain route-controlled where they already are.

Hook validation remains authoritative. Collection rules are an initial gate for authentication, actor shape, and basic field protection; the hook enforces full PO eligibility, PO owner matching, PO payment compatibility, branch-claim access, and approval assignment.

### 3. Expense hook changes

`ProcessExpense` and `cleanExpense` explicitly calculate the expense actor model:

1. Resolve `auth.id`.
2. On create, set `creator = auth.id` server-side, ignoring any submitted `creator`.
3. On update, require `creator` to remain unchanged.
4. Resolve the linked PO before deciding whether `uid != auth.id` is allowed.
5. If `uid != auth.id`, validate the full bookkeeper eligibility list.
6. If eligible, force `uid = po.uid` and `approver = creator`.
7. If `uid != auth.id` and the create is not eligible, refuse it. If `uid == auth.id`, use the manager-derived approver.

The manager/tapr validation stays in place for regular expenses. The bookkeeper path validates active `creator` plus `book_keeper` instead.

### 4. Routes and SQL projections

Expense SQL and response structs:

- include `creator`;
- include `creator_name`;
- include creator in details responses;
- include creator in list responses;
- both `expenses_select_base.sql` and `expense_details.sql` select them.

Route predicates:

- `/api/expenses/list`: includes `creator = auth.id` as well as `uid = auth.id`;
- `/api/expenses/list?purchase_order=...`: includes creator-owned rows for that PO;
- `/api/expenses/details/{id}`: includes `creator = auth.id` in the expense visibility predicate;
- `/api/purchase_orders/visible/{id}/expenses`: keeps the existing PO visibility gate, then allows rows visible through the expanded expense predicate.

The pending and approved approval queues use `approver = auth.id`.

### 5. Submit, recall, delete, and edit lifecycle

Every expense route or hook that used to say "only `uid` can do this" now uses `creator` as the operational owner:

- regular rows: `creator = uid`;
- bookkeeper-entered rows: `creator != uid`.

`creator` ownership applies to draft edit, draft delete, submit, and recall checks. Approved/committed locks are unchanged.

### 6. Frontend changes

The green plus expense creation flow:

- the visible PO payload includes `uid`;
- when creating from `/expenses/add/{poid}`, the editor detects:
  - current user has `book_keeper`;
  - linked PO `uid` differs from current user;
  - linked PO has a PO number;
  - the selected payment type is `OnAccount` or `CorporateCreditCard`;
- it sends the linked PO `uid` as `uid` for that narrow case;
- otherwise it sends the authenticated user as `uid`.

The editor does not allow arbitrary user selection for this feature. The only allowed alternate `uid` is the linked PO owner.

List/details displays show "Entered by" when `creator != uid`, and action visibility uses `creator` rather than `uid` for edit/delete/submit/recall.

### 7. Tests

Backend tests exist for this flow (the `TestBookKeeper...` tests in `app/expenses_test.go`). The planned coverage was:

- existing expense create still requires `uid = auth.id`;
- existing manager-derived approver behavior remains unchanged;
- `book_keeper` can create an `OnAccount` expense against another user's active PO and gets `creator = auth.id`, `uid = po.uid`, `approver = auth.id`;
- same bookkeeper creating against their own PO follows regular manager approval rules;
- non-bookkeeper cannot create with `uid != auth.id`;
- bookkeeper cannot create for another user without a PO;
- bookkeeper cannot create for a user who is not the linked PO owner;
- bookkeeper cannot use inactive/closed/cancelled/unapproved POs;
- bookkeeper cannot use payment types outside `CorporateCreditCard` and `OnAccount`;
- draft creator can edit/submit/delete their bookkeeper-entered expense;
- `uid` user cannot approve unless they are the assigned approver;
- creator/bookkeeper can approve only when `approver = creator`;
- unrelated bookkeeper cannot approve;
- pending queue includes submitted bookkeeper-entered rows for the creator;
- list/details visibility includes creator-owned rows before submission.

The plan also called for frontend tests or targeted manual QA for these. No automated frontend test covers them:

- green plus from another user's eligible PO sets the expense for the PO owner;
- green plus from own PO keeps regular behavior;
- ineligible payment types do not expose the on-behalf behavior;
- details page clearly shows the PO owner and the entering bookkeeper.

## Rollout notes

This shipped without changing historical expense semantics because the migration backfilled existing rows with `creator = uid`. Regular expenses behave as before.

The migration and hook changes are in the code. The planned operational rollout was:

1. Deploy migration and hook changes.
2. Assign `book_keeper` only to bookkeeping staff who should perform AP invoice entry.
3. Verify one eligible active PO for another user in staging.
4. Create, submit, approve, reject, recall, and commit a bookkeeper-entered `OnAccount` expense.
5. Repeat for `CorporateCreditCard`.
6. Verify same-user bookkeeper PO entry still routes to manager approval.

## Open questions

The code answers these as follows:

- Bookkeeper-entered expenses do not auto-submit after save. They stay explicit draft -> submit -> approve like other expenses.
- **Not built.** The PO owner gets no notification that bookkeeping entered an expense on their behalf.
- The audit language says "Submitted for" for `uid`.
- `book_keeper` adds no navigation items. The "Pending My Approval" item under Expenses still needs `tapr`, so a bookkeeper without `tapr` approves from the expense's details page or from the `Approve` button on their own expense list.
