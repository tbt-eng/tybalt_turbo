# Expense Uncommit Requests

## Summary

Holders of the `commit` claim can request an uncommit from the expense itself, choosing a standard reason and saying whether the expense should also be rejected. Holders of the `admin` claim work through one queue and approve requests one at a time or in bulk. Approval uncommits, and rejects if requested, in a single transaction.

Payroll-paid expenses can't be requested once their payroll has closed, and other expenses can't be requested more than 90 days after commit. Those still go to Teams and use the existing direct `Uncommit` button.

This builds on [uncommit.md](uncommit.md). The direct uncommit stays.

## What The Teams Thread Shows

Source: the `Expenses - Turbo` Teams thread, 2026-05-29 to 2026-11-05. The last 50 messages were read in full and older ones were found by searching for "uncommit", so a request that never used the word could be missing. Payment types and claims come from the local database snapshot, whose latest commit is 2026-10-30.

- There were at least 53 uncommits in 23 weeks, and 47 have a stated reason. September had 12, October had 9, and the first 5 days of November had 5.
- Since 2026-09-11, 23 of 25 requests ended in a rejection. Today that takes a second person after the admin holder, because `admin` can't reject.
  - On 2026-09-25 the expense on PO 2609-0364 was accidentally recommitted while waiting for that rejection.
  - On 2026-10-27 the usual rejector was out of office.
- Requests name a PO, but only expenses can be uncommitted. This caused repeated clarifying questions, and one no-PO expense needed its id pasted from the URL bar.
- Requests got lost or duplicated.
  - One waited a week (2026-10-16 to 2026-10-23).
  - One was for an expense already rejected (2026-10-19).
  - One PO was uncommitted three times (2604-0328).
- Age checks were ad hoc: "it is older, please confirm with Andrew", and "check with Andrew to ensure nothing was paid out".
- Of the 22 POs in the thread that still have expenses, 20 are `OnAccount` or `CorporateCreditCard`. Only 2 are `Expense`. The expenses on the other POs were deleted when those POs were cancelled, so their types are unknown.
- Uncommitting clears `committed`, `committer` and `pay_period_ending` and logs nothing, so the original commit can't be recovered later.

## Standard Reasons

Counts are the 47 requests with stated reasons, each assigned to its main reason. The other 11 asked for a PO to become Cumulative, which the companion fix below handles without an uncommit, so that isn't a reason. If such an expense is also wrong, the requester picks what's wrong with it.

| Reason                                | Count | Default outcome | Note prompt                  | Note required |
| ------------------------------------- | ----- | --------------- | ---------------------------- | ------------- |
| Duplicate PO, invoice is on another PO | 8    | Uncommit and reject | Which PO has the invoice? | Yes           |
| Amount doesn't match the invoice      | 7     | Uncommit and reject | Correct amount (tax, tip, shipping, over the 5%) | Yes |
| Vendor revised the invoice or price   | 4     | Uncommit and reject | What changed?            | No            |
| Wrong currency                        | 4     | Uncommit and reject | Correct currency         | No            |
| Wrong payment type                    | 4     | Uncommit and reject | Correct payment type     | Yes           |
| Expense is on the wrong PO            | 3     | Uncommit and reject | Which PO should it be on? | Yes          |
| Wrong document attached               | 2     | Uncommit and reject | For example, order confirmation instead of invoice | No |
| Order cancelled or billed in error    | 2     | Uncommit and reject | What happened?           | No            |
| Receipt doesn't meet policy           | 1     | Uncommit and reject | For example, not itemized | No           |
| Committed by mistake                  | 1     | Uncommit only   | What should have happened? | No          |
| Wrong vendor, job, branch or coding   | 0     | Uncommit and reject | What should it be?       | Yes           |
| Other                                 | 0     | Uncommit and reject | Explain                  | Yes           |

- `Wrong vendor, job, branch or coding` has no requests in this thread, but it came up in other chats and by email.
- 3 more requests also had the wrong document attached. They are counted under their main reason.
- The default outcome only pre-fills the form. The requester can change it.

## Who Can Do What

| Action                         | Who                                                 |
| ------------------------------ | --------------------------------------------------- |
| Create a request               | `commit` claim                                      |
| Cancel a pending request       | The requester                                       |
| View the queue and history     | `commit` or `admin` claim                           |
| Approve or decline             | `admin` claim                                       |

Approval deliberately stays limited to the `admin` claim. Requests are a queue in front of the existing uncommit, not a wider grant of it.

`admin` holders can normally open only committed expenses. They can also open any expense with a row in this history, so links on the Decided tab keep working after an uncommit.

## Eligibility

An expense can be requested when all of these hold, and each is checked again at approval:

1. Expense editing is enabled.
2. The expense is committed.
3. It isn't `_imported`.
4. It has no pending request.
5. It passes the age gate below for its payment type.

### Payroll Gate

Applies to the types repaid on payroll, per `app/routes/expense_reimbursements.sql`: `Expense`, `PersonalReimbursement`, `Allowance`, `Mileage` and legacy `Meals`.

The rule:

```
blocked when today > pay_period_ending + uncommit_payroll_lock_days
```

`uncommit_payroll_lock_days` lives in `app_config` under the `expenses` domain and defaults to 7. A payroll type with a blank `pay_period_ending` is blocked, since there's no way to tell whether it was paid.

Why 7: `GenerateCommittedPayPeriodEnding` still assigns commits to a pay period until the Saturday one week after it ends, and pay day is the following Friday (`pay_period_ending` + 13). So payroll can't be final before `pay_period_ending` + 7. The real run date isn't in the code, so if it turns out to be earlier, lower the setting.

Example: for the period ending 2026-10-24, requests are allowed through Saturday 2026-10-31, and pay day is 2026-11-06.

Why not "one week after submission":

- `submitted` is a boolean, so there is no submission date to measure from.
- It would block expenses that haven't been paid. An expense can wait weeks for approval: on PO 2607-0293 the expense dated 2026-07-24 was committed on 2026-08-25, in the period ending 2026-08-29.
- It would allow expenses that have been paid. An expense dated 2026-10-20, submitted and committed on 2026-10-30, lands in the period ending 2026-10-24. One week after submission is pay day itself.

### Non-Payroll Gate

Applies to `OnAccount`, `CorporateCreditCard` and `FuelCard`.

```
blocked when today > committed date + uncommit_non_payroll_lock_days
```

`uncommit_non_payroll_lock_days` lives in the same `app_config` domain and defaults to 90. Inside that window the queue shows how long ago each expense was committed and turns amber past 30 days (display only).

### When A Gate Closes While Pending

The request stays in the queue, flagged `Payroll closed` or `Too old`, and can't be approved from the queue. The admin holder confirms with payroll or accounting and uses direct `Uncommit`, or declines it.

## Approval

Approving from the queue, or using direct `Uncommit` on an expense with a pending request, does the following in one transaction per request:

1. Re-check eligibility.
2. Copy the commit state onto the request.
3. Clear the commit fields and reopen the PO if needed, using the existing uncommit logic.
4. If the outcome is `uncommit_and_reject`, set `rejected`, `rejector` to the requester (who holds `commit`), and `rejection_reason` to `<reason>: <note>`.
5. Mark the request `approved`.

The rejection notification goes out after the transaction, as the reject route does today, and the employee then recalls and corrects as usual.

Bulk approve runs each request in its own transaction and returns a result for each, so one failure doesn't block the rest.

Declining requires a comment and notifies the requester. There is no bulk decline.

## Data Model

Collection `expense_uncommit_reasons`, seeded by a migration and editable in the PocketBase admin:

- `name`
- `default_outcome`
- `note_prompt`
- `note_required`
- `sort`
- `active`

Collection `expense_uncommit_requests`, with no direct API access (all rules null, custom routes only, as in `uncommit.md`):

- `expense` (relation, optional): PocketBase refuses to delete a record that a required relation points at, and a rejected expense on a duplicate PO has to be deleted before the PO can be cancelled. Deleting the expense blanks this field and the snapshot keeps the history readable.
- `requester` (relation to users, required)
- `reason` (relation): required by the route, and blank only on rows written by direct `Uncommit`
- `note` (text)
- `outcome` (`uncommit_and_reject` or `uncommit_only`)
- `status` (`pending`, `approved`, `declined`, `cancelled`)
- `decided_by`, `decided_at`, `decision_comment`
- Snapshot taken at request and refreshed at approval: `employee_name`, `po_number`, `vendor_name`, `description`, `currency_code`, `payment_type`, `total`, `committed`, `committer`, `committed_week_ending`, `pay_period_ending`

A partial unique index on `expense` where `status = 'pending'` prevents duplicate requests.

Direct `Uncommit` also writes to this history. With a pending request it does what the request asked for, including rejecting, and marks it approved, and its dialog says so. It skips the eligibility checks, since it's how blocked requests get handled. With no pending request it writes a new `approved` row with the outcome `uncommit_only`.

## Routes

All mutating routes call `requireExpensesEditing`.

- `POST /api/expenses/{id}/uncommit_requests`
- `POST /api/expense_uncommit_requests/{id}/cancel`
- `GET /api/expense_uncommit_requests?status=`
- `GET /api/expense_uncommit_requests/reasons`
- `POST /api/expense_uncommit_requests/approve` with `{ ids: [] }`, also used for a single approval
- `POST /api/expense_uncommit_requests/{id}/decline`

`GET /api/expenses/details/{id}` gains an `uncommit_request` block for `commit` holders: eligibility, the blocking reason and pay day, and any pending request.

The uncommit and reject bodies move out of their handlers into shared functions that take a `txApp`, so the existing routes and this one run the same code.

## UI

Follow `ui/CLAUDE.md`.

- **Expense details:** `Request uncommit` opens a `DSPopover` with reason, outcome and note. It is hidden from `admin` holders, who have the direct `Uncommit` button and would only be asking themselves. When blocked, a `DSCallout` gives the period and pay day and says to message Dean on Teams. When pending, a `DSCallout` shows the request and a cancel action.
- **PO details:** must reach each committed expense in one click, since people think in PO numbers. Check this during the build.
- **New page `/expenses/uncommit-requests`:** `DSTabBar` for Pending and Decided, listed with `DSList`, with a pending count in the nav for `admin` holders, like the Commit Queue.
  - Each row shows PO and vendor, employee, payment type, total, committed date and age, reason and note, outcome, requester and time waiting.
  - Rows have a `DsCheck`, with `Approve selected (N)` and `Approve all`.
  - The bulk confirmation states how many will be rejected.
  - `Approve all` is disabled while a search is active, because it acts on the whole queue, not only the rows showing.

## Notifications

- **Requester:** notified on approval or decline.
- **Employee:** the existing rejection notification.
- **Admin holders:** the nav count, plus one daily reminder email while anything is pending, following `notifications/queue_reminders.go`. No email per request.

## Companion Fix: Convert Closed One-Time POs

11 of the 47 requests were about converting a PO to Cumulative. A One-Time PO closes on commit, and `make_cumulative` requires `Active`, so uncommitting is the only way to reopen it. At least 4 of the 11 said the expense itself was fine, and 2 more look that way from the data. Those uncommits move a correct expense into a new commit week and pay period for no reason.

`make_cumulative` also accepts a One-Time PO closed by the system (`closed_by_system = true`), still for `payables_admin` holders:

- If the committed total is below the PO total, it converts the PO and reopens it.
- Otherwise it refuses, because there is no room for another invoice.
- It refuses imported POs, because the import tool sets `closed_by_system` on all of them, so their real cause of closure is unknown.

## Out Of Scope

- **Time sheets:** direct uncommit is unchanged.
- **Correcting a committed expense in place:** currency and payment type fixes account for 8 of the 47 requests and could justify this later.
- **Legacy writeback:** it selects `committed != ''`, so an uncommit never reaches Tybalt and stale rows remain there. This already happens today, and the request history makes those rows easier to find.
- **Deleting rejected expenses on duplicate POs:** someone still has to delete the rejected expense before the PO can be cancelled.

## Tests

- **Permissions:** only `commit` holders can request, only `admin` holders can approve or decline, and only the requester can cancel.
- **Age gates:** every payroll type on both sides of the boundary, non-payroll types at 90 days, and the config overrides.
- **Gate closing while pending:** approval is refused.
- **Duplicate pending request:** refused, and a new request is allowed after a decision.
- **`uncommit_and_reject`:** commit fields cleared, rejection fields set, PO reopened, employee notified, and the snapshot saved.
- **`uncommit_only`:** the expense returns to the commit queue.
- **Bulk approve:** a mix of valid and invalid requests applies the valid ones and reports the rest.
- **Direct uncommit:** applies a pending request's outcome, including after its gate closed, and writes a row when there is none.
- **Editing disabled:** all mutating routes are blocked.
- **Companion fix:** convert and reopen, refuse when the PO is used up, and refuse when the PO was closed by a person.

## Decisions

- **Payroll cutoff:** `pay_period_ending` + 7 days. Andrew hasn't confirmed when payroll is actually pulled, so this is the default until someone says otherwise.
- **Old non-payroll expenses:** blocked 90 days after commit.
