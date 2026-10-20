# Client invoices

Tracking issue: [Client invoices (AR) #152](https://github.com/tbt-eng/tybalt_turbo/issues/152)

This is the first, small piece of accounts receivable. It records the
invoices we send clients, so that a cost billed back to a client (#148's
"TBTE #") is linked to the invoice that recovered it instead of being typed
onto a spreadsheet row. Full AR comes later; see "Later" at the end.

## How invoicing works today

- The admin team invoices clients: one person in each of Ottawa, Kenora,
  Kitchener/Waterloo, Toronto and Collingwood, and four in Thunder Bay.
- Invoice numbers are `YYMMNNN`, for example `2609001`: the year, the month,
  then a three-digit number that starts again at `001` every month. There is one
  sequence for the whole company.
- Numbers are taken by hand, first come, first served. The team shares a Teams
  group, and whoever takes the next number posts it there.
- Once an invoice is sent it is not cancelled. A change in amount is made with
  a credit invoice, which takes a new number. A change of address or contact
  email is a revision under the same number.
- An invoice that has not been sent yet is simply cancelled, and its number is
  used for another invoice.

## The model

### `invoices`

| Field | Notes |
| --- | --- |
| `client` | Required relation to `clients`. An invoice belongs to a client, not a job, so one invoice can cover several jobs. |
| `number` | `YYMMNNN`, chosen by the admin team. Unique across all invoices. Empty is allowed on a draft. |
| `date` | `YYYY-MM-DD`. Required to issue. |
| `total` | The invoice total in Canadian dollars, as on the invoice. Required to issue, and greater than zero. |
| `status` | `Draft` or `Issued`. Issued means sent to the client. |
| `creator` | Set by the server on create. |
| `issued_by`, `issued_at` | Set by the server when the invoice is issued. |

There is no void status. An issued invoice is never cancelled, so it keeps its
number for good.

### `invoice_lines`

One line records that a cost was billed back to the client on an invoice.

| Field | Notes |
| --- | --- |
| `invoice` | Required relation to `invoices`. |
| `kind` | `expense` or `po`. |
| `expense` or `purchase_order` | The cost billed back. Exactly one, matching `kind`. |
| `amount` | Required, greater than zero. What was billed for this cost, net of tax, in Canadian dollars. |
| `creator` | Set by the server on create. |

Lines cannot be edited. A wrong line is removed and added again.

## Rules

### Who

- Only active users with the `invoicing` claim can create, edit, issue or
  delete invoices, and add or remove lines. The claim is for the admin team
  above. `accounting` (which can cancel POs) and `payables_admin` (accounts
  payable) are not reused.
- Holders of `invoicing` or `report` can list and open invoices and lines.
- Anyone who can open an expense or PO sees the invoice it was billed on, on
  its details page.

### Invoice numbers

- **Format:** seven digits, `YYMMNNN`, where `MM` is `01` to `12`. This catches
  most typos.
- **Unique** across all invoices. If two people take the same number, the
  second is refused when they save and takes the next one from the Teams group.
- **Required to issue.** A draft may have no number yet.
- **Deleting a draft frees its number** for another invoice, the same as
  cancelling an unsent invoice today.

Turbo does not assign numbers yet. The Teams group stays the way numbers are
taken until the full AR pass.

### Drafts and issuing

- A draft can be edited or deleted.
- Issuing needs a number, a date and a total greater than zero. Issuing sets
  `status`, `issued_by` and `issued_at` in one transaction, through
  `POST /api/invoices/{id}/issue`.
- An issued invoice cannot be edited or deleted. A mistake made on an issued
  invoice is fixed directly in the database until the full AR pass adds a way
  for an admin to return it to draft.
- `status`, `creator`, `issued_by` and `issued_at` cannot be set through the
  collection API; the server ignores or refuses them.

### Lines

- A line can be added only to an issued invoice.
- The expense's or PO's job must belong to the invoice's client.
- An expense must be committed.
- A PO must be `Active` or `Closed`. An `Unapproved` PO is not yet a cost, and a
  `Cancelled` one never will be.
- **No double billing:** an expense or PO can be on at most one line. Unique
  indexes on `invoice_lines.expense` and `invoice_lines.purchase_order` enforce
  it.
- Turbo does not stop a PO and one of that PO's expenses from both being
  billed. Billing the PO or its expenses is the admin team's choice.
- Lines are checked when they are added. Uncommitting a billed expense, or
  moving a job to another client, does not remove the line.

### Clients

- A client with invoices counts as referenced: it cannot be deleted, and the
  clients list hides its Delete action.
- Absorbing a client into another moves its invoices with it
  (`invoices.client` is in the clients absorb config).

## Screens

- **Invoices** (`/invoices/list`), for `invoicing` and `report` holders: every
  invoice, newest first, with its client, number, date, total and status.
  `invoicing` holders get a New Invoice button.
- **Invoice details** (`/invoices/{id}/details`): the invoice and its lines.
  On a draft, `invoicing` holders can edit, delete and issue it. On an issued
  invoice they can add a billed-back cost: pick one of the client's jobs, then
  one of its committed expenses or Active or Closed POs that is not billed
  yet, and enter the amount billed. Lines can be removed.
- **Expense and PO details:** a "Billed to client" card shows the invoice
  number, date and amount billed for anyone who can open the record. An
  `invoicing` holder sees "Bill on invoice" on an unbilled, eligible record, to
  pick one of the client's issued invoices and enter the amount.

Adding lines from the invoice page is the main path: admin team members who do
not hold `report` cannot open committed expenses, but they can still bill them
from the invoice.

## API

| Route | Who | What |
| --- | --- | --- |
| `invoices` collection | see "Who" | List, view, create, update and delete drafts. |
| `invoice_lines` collection | see "Who" | List, view, create and delete lines. |
| `POST /api/invoices/{id}/issue` | `invoicing` | Issue a draft. |
| `GET /api/invoices/{id}/lines` | `invoicing`, `report` | An invoice's lines, with each cost's job, PO number, description, vendor and cost. |
| `GET /api/invoices/{id}/billable` | `invoicing` | The client's jobs with unbilled costs; with `?job=`, that job's unbilled costs. |
| `GET /api/expenses/client_invoice/{id}` | anyone who can open the expense | The line billing the expense, if any, and whether the caller may bill it. |
| `GET /api/purchase_orders/client_invoice/{id}` | anyone who can open the PO | The same for a PO. |

## Later: full AR

A separate pass once #148 ships:

- Turbo assigning invoice numbers on issue, replacing the Teams group.
- An admin returning an issued invoice to draft for edits.
- Credit invoices, numbered from the same sequence.
- Producing or storing the invoice PDF.
- Invoiced to date and over/under billing in WIP (#103).
- Time billing through a `billed_through` date, unbilled-time reports, and
  billing time amendments.
- #150's optional rule blocking rate sheet changes dated before the last
  invoice.
