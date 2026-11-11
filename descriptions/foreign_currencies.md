# Multi-Currency Support for POs and Expenses (Issue #121)

## Overview

This feature is built. This note started as the implementation plan and now says what
Turbo does. Turbo tags Purchase Orders and Expenses with a currency (defaulting to CAD),
displays indicative exchange rates throughout the workflow, and tracks settled CAD amounts
for non-CAD expenses. Exchange rates are sourced each weekday from the Bank of Canada Valet API.

Non-home-currency OnAccount and CorporateCreditCard expenses require explicit settlement
by a payables admin before they can be committed, managed through a dedicated Currency
Settlement queue.

Approval authority is evaluated in home currency (CAD) via `approval_total_home`, but
PO spending caps on expenses are enforced in the PO's stated currency. This ensures
that forex fluctuations after approval do not invalidate a user's ability to spend
against an approved PO.

Architectural note: the home currency is intentionally fixed to CAD in this system.
Exchange rates come from the Bank of Canada CAD-anchored feed, and settlement/reporting
semantics are defined in CAD. Supporting configurable home currency is out of scope
unless the FX provider, derived-field semantics, and migration strategy are redesigned.

---

## 1. New Collection: `currencies`

| Field       | Type   | Notes                                                      |
|-------------|--------|------------------------------------------------------------|
| `code`      | string | ISO 4217 code (e.g. "CAD", "USD"). Unique.                |
| `symbol`    | string | Display symbol (e.g. "$", "US$", "€")                     |
| `icon`      | file   | SVG icon for the currency                                  |
| `rate`      | number | 1 unit of this currency = X CAD. CAD row is always 1.     |
| `rate_date` | date   | Date of the Bank of Canada observation the rate came from  |
| `ui_sort`   | number | Sort order for currency selectors (lower = first).         |

- No migration seeds the rows. A holder of the `admin` claim adds CAD, USD and any other
  currencies on the Currencies page (`/currencies`). Any signed-in user can read them.
- The hook forces the CAD row's `rate` to 1 and refuses a foreign currency without a positive rate.
- Open-ended: an `admin` holder can add new currencies at any time. The delete route
  (`DELETE /api/currencies/{id}`) refuses a currency that a PO or expense uses.
- The `rate` and `rate_date` fields are updated in place by the weekday cron job, or when an
  `admin` holder uses Reload Rates (`POST /api/currencies/reload_rates`).
- `ui_sort` controls the display order in DSCurrencyInput selectors.

---

## 2. Daily Exchange Rate Cron Job

- **Source**: Bank of Canada Valet API (free, no API key required).
  - Example: `https://www.bankofcanada.ca/valet/observations/FXUSDCAD/json?recent=5`
    (Turbo takes the newest of the five observations that has a value)
  - Series naming: `FX{CODE}CAD` (e.g. `FXUSDCAD`, `FXEURCAD`)
- **Schedule**: Monday to Friday at 22:00 UTC (`0 22 * * 1-5`), after the Bank of Canada publishes the day's rate.
- **Behaviour**: For each non-CAD row in the `currencies` table, fetch the latest
  rate from the Bank of Canada and update `rate` and `rate_date` on that row, but only
  when the observation date is later than the stored `rate_date`.
- **Weekends/holidays**: API returns no new data; existing rate remains unchanged.
- **Implementation**: the `currency_rate_sync` job in `cron/AddCronJobs`, which calls
  `utilities.SyncCurrencyRates`.

---

## 3. Purchase Order Changes

### New fields on `purchase_orders`

| Field                 | Type                   | Notes                                                             |
|-----------------------|------------------------|-------------------------------------------------------------------|
| `currency`            | relation -> currencies | Required. Default CAD.                                            |
| `approval_total_home` | number                 | CAD equivalent of `approval_total`, computed by hook using cached `currencies.rate`. |

### Hook logic

- **Clean/Validate**:
  - `currency` must reference a valid `currencies` record. A blank currency becomes the
    parent PO's currency, or CAD when there is no parent.
  - A foreign currency must have a positive `rate`.
  - Compute `approval_total_home = approval_total * currency.rate`.
  - For CAD, `approval_total_home` = `approval_total` (rate is 1).
- **Child POs**: Must inherit parent's `currency` (validated in existing child PO rules).
- **Approval policy comparison**: Use `approval_total_home` (not `approval_total`)
  when Turbo compares the PO amount with the second-approval threshold and the
  CAD-denominated limits in `po_approver_props`.
- **Approver pool determination**: Use `approval_total_home` when resolving eligible
  first-stage and second-stage approvers.

### Currency vs. amount enforcement separation

- **Approval authority** is evaluated in CAD (`approval_total_home`).
- **Expense amount caps** are enforced in the PO's stated currency — `expense.total` is
  compared against `po.total` (+ excess) directly, both in the same currency. This means
  forex fluctuations after PO approval do not require users to obtain a new PO.
- **The fixed part of the excess** (`po_expense_allowed_excess.value`) is in CAD. It is
  converted to the PO's currency with the rate of the PO's stored home total
  (`approval_total_home / approval_total`), so it does not move with today's rate either.
  See `app_config.md`.

### UI changes

- DSCurrencyInput on PO form (default CAD, currency selectable).
- When non-CAD selected, display indicative CAD equivalent under the total
  (e.g. "CAD equivalent: $13,845.00 · rate 1.3845 · as of 2026-04-01").
- PO detail/print views show currency code/symbol alongside amounts.
- Approval screens show both the PO currency amount and the CAD equivalent.

---

## 4. Expense Changes

### New fields on `expenses`

| Field           | Type                   | Notes                                               |
|-----------------|------------------------|------------------------------------------------------|
| `currency`      | relation -> currencies | Required. Default CAD.                               |
| `settled_total` | number                 | Amount in CAD. Required for commit.                  |
| `settler`       | relation -> users      | Payables admin who entered the settlement.           |
| `settled`       | datetime               | Timestamp when settlement was set.                   |

### Currency assignment rules

- **Linked to PO**: `currency` is inherited from the PO and immutable.
- **No PO, payment type is `OnAccount`, `Expense`, or `CorporateCreditCard`**:
  User selects currency (default CAD).
- **Payment type is `Allowance`, `FuelCard`, `Mileage`, or `PersonalReimbursement`**:
  Always CAD. No currency selector shown. Hook enforces CAD.

### `settled_total` rules by payment type

| Scenario                                  | Who enters it           | When required     | Mutability                                  |
|-------------------------------------------|-------------------------|-------------------|---------------------------------------------|
| Currency = CAD (any type)                 | Hook (auto-set = total) | Always set auto   | Updates with `total`                        |
| Currency ≠ CAD, `Expense`                 | User                    | Before submit     | Immutable once submitted                    |
| Currency ≠ CAD, `OnAccount`               | Payables admin          | Before commit     | Cleared and re-entered via settlement queue |
| Currency ≠ CAD, `CorporateCreditCard`     | Payables admin          | Before commit     | Cleared and re-entered via settlement queue |

A `settled_total` entered for a non-CAD expense, by the user or by a payables admin, must
be within 20% of the current CAD equivalent (`total` times the cached rate). Turbo checks
this on save, on submit and in the settlement queue.

### Settlement fields (`settler`, `settled`) rules

- Only populated for non-home-currency `OnAccount` and `CorporateCreditCard` expenses.
- Set together with `settled_total` by the payables admin via the settlement queue.
- For `Expense` payment type (user-entered `settled_total`), `settler` and `settled`
  remain blank — the user enters `settled_total` directly on the expense form.
- For CAD expenses, all three fields (`settled_total`, `settler`, `settled`) are
  managed by the hook: `settled_total = total`, `settler` and `settled` remain blank.

### Hook logic

- **Clean**: If currency = CAD, set `settled_total = total`.
- **Validate on submit** (`Expense` type, currency ≠ CAD): `settled_total` must be > 0.
- **Validate on submit** (`Expense` type, currency ≠ CAD): `settled_total` is immutable
  (reject changes to `settled_total` if already submitted).
- **Validate on commit** (`Expense` type, currency ≠ CAD): `settled_total` must be > 0.
  A CAD expense always has `settled_total = total`, so commit does not check it.
- **Commit gate**: Non-home-currency `OnAccount`/`CorporateCreditCard` expenses must
  have non-blank `settled` timestamp to appear in the commit queue, and the commit route
  refuses them until then.
- **No-PO expense limit**: The no-PO limit check (default $100) compares against `settled_total`
  (CAD equivalent) so the limit is consistently applied in CAD. While a foreign-currency
  expense has no `settled_total` yet, the hook compares the indicative CAD amount
  (`total` times the cached rate). Submit, settle and commit check the limit again
  with the real `settled_total`.
- **PO amount caps** (One-Time, Recurring, Cumulative overflow): Compare `expense.total`
  against `po.total` in the PO's stated currency. Both are in the same currency — no
  conversion needed. Forex changes after PO approval do not affect this check.

### Rejection behaviour

- When a non-home-currency `OnAccount` or `CorporateCreditCard` expense is rejected
  (via the reject route or committer-reject path), the settlement is **automatically
  cleared**: `settled_total`, `settler`, and `settled` are all set to their zero values.
- For `Expense` payment type (user-entered `settled_total`), rejection does **not**
  clear `settled_total` — the user manages it themselves upon revision and resubmission.
- For CAD expenses, rejection has no effect on settlement fields — the hook continues
  to set `settled_total = total`.

### UI changes

- DSCurrencyInput on expense form:
  - Linked to PO: currency locked (inherited from PO), amount editable.
  - No PO, eligible type: currency selectable, amount editable.
  - Ineligible type (Allowance, etc.): currency locked to CAD, amount editable.
- `settled_total` input shown for `Expense` payment type when currency ≠ CAD
  (plain number input with CAD symbol — always settling into home currency).
- Indicative CAD equivalent shown alongside foreign currency total.
- Expense lists/detail views show currency alongside amounts.

---

## 5. DSCurrencyInput Component

A reusable Svelte 5 component for currency-aware amount entry.

### Layout

```
┌──────────────────────────────────────┐
│ [🇺🇸 USD ▾]   10,000.00          $ │
└──────────────────────────────────────┘
  ↑ currency       ↑ amount       ↑ symbol
  selector         input          (from currency)
```

- **Left side (inside field)**: Currency selector showing icon + code. Clicking opens
  a dropdown of available currencies, sorted by `ui_sort`.
- **Right side**: Currency symbol displayed as static text inside the input field.
- **Center**: Number input for the amount.

### Props

| Prop                    | Type                | Default     | Notes                                        |
|-------------------------|---------------------|-------------|----------------------------------------------|
| `amount`                | number              | —           | Bindable. The amount entered.                |
| `currency`              | string              | —           | Bindable. Currency record ID.                |
| `items`                 | CurrenciesResponse[] | `[]`       | Available currencies (pre-fetched, sorted by `ui_sort`). |
| `disabledCurrency`      | boolean             | `false`     | Locks the currency selector (e.g. expense linked to PO). |
| `disabledAmount`        | boolean             | `false`     | Locks the amount input (display-only mode).  |

The component also takes `errors`, `amountFieldName`, `currencyFieldName`, `uiName`,
`amountStep`, `amountMin`, `helperText`, `homeEquivalent`, `rate`, `rateDate`, `displayCode`
and `displaySymbol`. There is no `placeholder` prop.

### Derived state

- `symbol`: Resolved from the selected currency record's `symbol` field.
- `icon`: Resolved from the selected currency record's `icon` field (SVG).
- `code`: Resolved from the selected currency record's `code` field.

### Usage examples

```svelte
<!-- PO form: full interactive -->
<DSCurrencyInput bind:amount={total} bind:currency={currencyId}
  items={allCurrencies} />

<!-- Expense linked to PO: locked currency, editable amount -->
<DSCurrencyInput bind:amount={total} bind:currency={currencyId}
  items={allCurrencies} disabledCurrency={true} />

<!-- Read-only display -->
<DSCurrencyInput amount={total} currency={currencyId}
  items={allCurrencies} disabledCurrency={true}
  disabledAmount={true} />
```

---

## 6. Currency Settlement Queue

The Settlement page under Expenses (`/expenses/settlement`), accessible to users with `payables_admin` claim.
Manages settlement of non-home-currency `OnAccount` and `CorporateCreditCard` expenses.

### Entry criteria

An expense appears in the settlement queue when ALL of the following are true:

- Payment type is `OnAccount` or `CorporateCreditCard`
- Currency ≠ CAD (non-home currency)
- Approved (`approved` is non-blank, `rejected` is blank)
- Uncommitted (`committed` is blank)

These criteria are one SQL condition, `expenseSettlementQueueCondition` in
`app/routes/expense_settlement.go`. The queue lists, the nav badge count and the
payables rule on expense details all use it.

### Who can open a queue expense

Each row has a **View** button that opens the expense details page. A `payables_admin`
holder can load the details (`GET /api/expenses/details/{id}`) and the receipt
(`GET /api/expenses/attachment/{id}`) of an expense only while it meets the entry
criteria, settled or not. Once it is committed, or if it is rejected, they lose this
access, unless another rule lets them see it (owner, creator, approver, `commit`,
`report` or `admin`). The rule is `expensePayablesVisibilityPredicate` in
`app/routes/expenses_api.go`. It applies to those two routes only, not to the list of a
PO's expenses.

### Main Tab: Unsettled

Lists expenses matching entry criteria where `settled` is blank.

**Columns displayed:**

- Date
- Employee name (the expense's `uid`; the API field is called `creator_name`)
- Payment type, and the card's last 4 digits for a corporate card
- Vendor
- Description
- PO number (if linked)
- Currency code
- Total (in foreign currency)
- Indicative CAD equivalent (from cached rate on `currencies` row)
- Age (days since approval) is returned by the API as `age_days`, but the page does not show it

**Per-row action:** CAD amount input + **Save** button, and **View**.

**Save action** (API route, runs in a database transaction):

1. Re-read the expense record inside the transaction.
2. Verify the expense is still approved, uncommitted, and unsettled.
3. If state has changed (committed, rejected, or already settled), abort with error.
   Also abort when the amount is not above 0, is not within 20% of the current CAD
   equivalent, or (for an expense with no PO) reaches the no-PO limit.
4. Set `settled_total` to the entered amount, `settler` to the authenticated user,
   `settled` to the current timestamp.
5. Commit the transaction.

**Sort order:** Descending by date.

### Second Tab: Settled (Revision)

Lists expenses matching entry criteria where `settled` is non-blank.
This tab is **read-only** — no in-place editing of settled amounts.

**Columns displayed:** Same as unsettled tab, plus:

- Settled total (CAD)
- Settler name
- Settled timestamp

**Per-row action:** **Clear** (Clear Settlement) and **View**.

**Clear Settlement** (API route, runs in a database transaction):

1. Re-read the expense record inside the transaction.
2. Verify the expense is still uncommitted, approved, not rejected, and settled.
3. If already committed (or otherwise out of the queue), abort with error.
4. Clear `settled_total`, `settler`, and `settled`.
5. Commit the transaction.

**After clearing:** Stay on the settled tab and refresh the list. The cleared item
disappears from this tab and reappears on the unsettled tab. This allows the admin
to clear multiple items in rapid succession without navigating away.

**Revision workflow:** To revise a settled amount, the admin clears it (settled tab),
then re-enters the correct amount (unsettled tab). This two-step approach is
preferred because clearing is faster and wins the race against a concurrent commit —
the item is immediately ineligible for commit once cleared.

### Transaction safety

All settlement operations (set, clear) run in database transactions to prevent
races with concurrent commit, reject, or settlement actions:

| Race scenario          | Protection                                                       |
|------------------------|------------------------------------------------------------------|
| Settle vs. commit      | Transaction checks `committed` is blank before settling.         |
| Settle vs. reject      | Transaction checks `rejected` is blank before settling.          |
| Clear vs. commit       | Transaction checks `committed` is blank before clearing.         |
| Concurrent settlement  | Transaction checks `settled` is blank before settling (prevents double-settle). |

---

## 7. Reporting and Writeback Gap (Stakeholder Decision Required)

What Turbo does today:

- The legacy expense export (`GET /api/export_legacy/expenses/{updatedAfter}`) keeps `total`
  in the expense's own currency. When the expense has a currency set, it adds `currency`
  (the code) and `settled_total` (CAD).
- The purchase orders in that export carry no currency field.
- The payroll and weekly expense reports (`app/reports/expenses.sql`) keep Total in CAD: a
  foreign-currency expense uses `settled_total`. They add `currency` and
  `foreign_currency_total` columns for the original amount.

**Not built.** The code records no stakeholder decision, and none of the options below is
implemented as a whole. The rest of this section is the original open question.

This rollout intentionally leaves one important area unresolved: downstream reporting
and legacy writeback semantics become ambiguous once foreign-currency POs and expenses
exist.

The ambiguity is:

- `total` is stored in the PO/expense's stated currency
- `settled_total` is stored in CAD
- `approval_total_home` is stored in CAD
- existing exports, writebacks, and downstream reports may assume a single amount field
  always represents CAD or always represents the operational amount, which is no longer
  true once multi-currency records exist

This is a priority topic to review with stakeholders before the feature is considered
fully complete for downstream accounting/reporting workflows.

### Why this matters

- Aggregated reporting can accidentally mix source-currency and CAD-settled amounts.
- Legacy consumers may continue reading `total` as if it were a CAD amount.
- Writeback targets may need both the original foreign amount and the CAD-settled amount,
  with explicit currency labels, to remain auditable.
- Without a defined contract, two reports can both be "correct" internally while still
  disagreeing because they are using different money fields for different purposes.

### Options to discuss with stakeholders

#### Option A: Keep legacy writebacks unchanged for now

- Lowest implementation cost.
- Acceptable only as a temporary measure.
- Requires clear documentation that foreign-currency records are not fully represented
  in downstream reporting/export semantics yet.

#### Option B: Add explicit currency-labeled export fields

- Recommended long-term direction.
- Every exported/writeback amount field should clearly indicate whether it is:
  - source-currency amount (`total`, plus `currency_code`)
  - CAD/home-settled amount (`settled_total`)
  - CAD/home-approval amount (`approval_total_home`)
- Preserves auditability and avoids hidden unit mixing.

#### Option C: Standardize all downstream reporting on CAD/home-currency amounts

- Simplifies roll-up reporting and accounting summaries.
- Still requires carrying the original source-currency amount and code somewhere for
  audit/reference.
- Best fit if stakeholder priority is consolidated financial reporting over preserving
  source-currency operational amounts in downstream systems.

#### Option D: Split operational writebacks from accounting/reporting exports

- Operational systems receive source-currency fields.
- Accounting/reporting systems receive CAD/home-normalized fields.
- Cleanest semantically, but highest coordination and implementation cost.

### Minimum decision needed

Stakeholders should explicitly choose:

- which amount is authoritative for downstream accounting
- whether writebacks must include both source-currency and CAD/home-normalized amounts
- whether existing legacy consumers must remain backward-compatible, or can be changed
- which reports are operational views versus accounting views

Until this decision is made, multi-currency behavior inside the app can be considered
correct, but the external reporting/writeback contract should be treated as incomplete.

---

## 8. API Changes

### New routes

| Route                                      | Method | Auth             | Purpose                              |
|--------------------------------------------|--------|------------------|--------------------------------------|
| `POST /api/expenses/:id/settle`            | POST   | `payables_admin`  | Set settlement amount                |
| `POST /api/expenses/:id/clear_settlement`  | POST   | `payables_admin`  | Clear settlement (back to unsettled) |
| `GET /api/expenses/unsettled`              | GET    | `payables_admin`  | List unsettled queue                 |
| `GET /api/expenses/settled`                | GET    | `payables_admin`  | List settled (uncommitted) queue     |

Currency routes: `GET /api/currencies` (any signed-in user), and for `admin` holders
`GET /api/currencies/init_status`, `POST /api/currencies/reload_rates`,
`POST /api/currencies/{id}/initialize_backfill` and `DELETE /api/currencies/{id}`.

### Existing route changes

- **Commit route**: Gate non-home-currency `OnAccount`/`CorporateCreditCard` expenses
  on `settled` being non-blank. Exclude unsettled expenses from the commit queue query.
- **Reject route**: Auto-clear settlement (`settled_total`, `settler`, `settled`) for
  non-home-currency `OnAccount`/`CorporateCreditCard` expenses upon rejection.

---

## 9. Existing Logic Impact

| Area                          | Impact                                                                 |
|-------------------------------|------------------------------------------------------------------------|
| Approval policy resolution    | Use `approval_total_home` instead of `approval_total` to decide if two approvals are required and to compare approver limits with the PO amount. |
| Cumulative PO overflow check  | No change — sums are in PO currency, expense inherits PO currency.     |
| One-Time/Recurring PO caps    | No change — comparison in PO currency.                                 |
| No-PO expense limit ($100)    | Compare `settled_total` (CAD) against the limit.                       |
| Commit queue                  | Exclude unsettled non-home-currency OnAccount/CorporateCreditCard.     |
| Reject route                  | Auto-clear settlement for non-home-currency OnAccount/CorporateCreditCard. |
| PO cancel (block if expenses) | No change.                                                             |
| PO close                      | No change.                                                             |
| Writeback to Firebase         | The legacy expense export includes the `currency` code and `settled_total` when the expense has a currency set. The PO export has no currency field. |
| Reporting / aggregation       | Expense reports and job expense summaries use `settled_total` for foreign-currency expenses and `total` for CAD ones. |

---

## 10. Migration Plan

The migrations did steps 1, 3 and 4. Steps 2 and 5 to 7 are not migrations: an `admin`
holder does them from the Currencies page.

1. Create `currencies` collection with fields: `code`, `symbol`, `icon`, `rate`,
   `rate_date`, `ui_sort`.
2. An `admin` holder adds the CAD row and any foreign currencies on the Currencies page.
   No migration seeds them.
3. Add `currency` and `approval_total_home` fields to `purchase_orders`. The hooks
   default a blank currency to CAD.
4. Add `currency`, `settled_total`, `settler`, and `settled` fields
   to `expenses`. The hooks default a blank currency to CAD.
5. Backfill: Initialize CAD on the Currencies page
   (`POST /api/currencies/{id}/initialize_backfill` on the CAD row) sets `currency` = CAD
   for every PO and expense with a blank currency.
6. The same action sets `settled_total` = `total` for CAD expenses that have none.
7. The same action sets `approval_total_home` = `approval_total` for CAD POs and for POs
   that have none.

---

## 11. Implementation Order

All of these steps are built. Seeding and backfill happen on the Currencies page, not in
a migration (see section 10).

1. **Migration**: Create `currencies` collection, seed data, add new fields to
   POs and expenses, backfill existing records.
2. **Cron job**: Bank of Canada rate fetcher in `cron/AddCronJobs`.
3. **PO hooks**: Currency validation, `approval_total_home` computation, child PO
   inheritance, approval policy updates to use `approval_total_home`.
4. **Expense hooks**: Currency assignment, `settled_total` auto-set for CAD,
   validation for submit/commit gates, PO amount caps in PO currency.
5. **Settlement API routes**: `settle`, `clear_settlement`, `unsettled` list,
   `settled` list — all with transaction safety.
6. **Reject route update**: Auto-clear settlement on rejection for non-home-currency
   OnAccount/CorporateCreditCard.
7. **Commit route update**: Gate on settlement for non-home-currency
   OnAccount/CorporateCreditCard.
8. **DSCurrencyInput component**: Reusable currency-aware amount input with
   icon + code selector, symbol display, and lock props.
9. **UI — PO forms**: Integrate DSCurrencyInput, CAD equivalent display,
   print view updates.
10. **UI — Expense forms**: Integrate DSCurrencyInput (locked currency when PO linked),
    `settled_total` input for `Expense` payment type.
11. **UI — Currency Settlement queue**: Two-tab interface with settle and clear actions.
12. **UI — Approval screens**: Dual-amount display (foreign + CAD equivalent).
13. **Tests**: Unit tests for hooks, approval policy with currency conversion,
    settlement transaction safety, rejection auto-clear, settled_total validation,
    cron job, commit gate, PO amount caps in PO currency.
