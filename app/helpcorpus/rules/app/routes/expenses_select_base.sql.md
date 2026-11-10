<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expenses_select_base.sql","sha256":"fd3af9c32d54a4537c3ed5b753c9fd6f7a44801a043730812551aa8d210c4a88","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expenses_select_base.sql

## Expense total and distance are read as numbers

Lines 9-9

Asked as: expense total as number; distance decimal value; expense amount type; distance shows as text

Turbo reads the expense total and the distance as decimal numbers.

## Expense attachment comes from the linked expense document

Lines 11-16

Asked as: receipt file missing on expense; attachment not showing on expense; linked document attachment name; expense receipt file name

Turbo takes the attachment file name and hash from the expense document that the expense links to (attachment_document). When no document is linked, both are empty. The attachment collection ID is 'pbc_2089657321' when a document is linked, and empty otherwise. The attachment record ID is the document's ID, or empty.

## Missing-attachment reason and card digits default to empty

Lines 14-14

Asked as: blank missing attachment reason; empty card digits on expense; no reason for missing receipt; card last four blank

When an expense has no reason for a missing attachment (attachment_missing_reason), no last 4 card digits (cc_last_4_digits) or no currency value, Turbo shows empty text.

## Expense workflow fields are shown as stored

Lines 17-30

Asked as: expense rejection reason stored; who rejected my expense; expense approval details; expense commit week ending

Turbo returns the expense's rejection (rejector, rejected, rejection_reason), approval (approver, approved), submission (submitted) and commit (committer, committed, committed_week_ending) fields as stored. It also returns pay_period_ending, allowance_types, job, category and kind as stored. It does not default or change them.

## Expense distance is read as a number

Lines 31-31

Asked as: expense distance as number; distance wrong type on expense; kilometres value on expense

Turbo casts the expense distance to a real number (distance) when it reads the expense.

## Missing card digits and currency show as empty text

Lines 32-33

Asked as: expense currency blank; no card digits on expense; empty currency on expense; missing currency for expense

When an expense has no last 4 card digits (cc_last_4_digits), Turbo shows empty text. When an expense has no currency (currency), Turbo shows an empty currency value.

## Unsettled expense shows settled total of 0

Lines 34-36

Asked as: unsettled expense total zero; expense settled total 0; who settled expense blank; expense not settled amount

When an expense has no settled total (settled_total), Turbo shows 0. When it has no settler or settled date, Turbo shows them as empty text.

## Expense lookups show empty text when the record is missing

Lines 39-59

Asked as: expense client blank; expense job missing; expense vendor not showing; expense shows no PO number; broken link on expense

For an expense, Turbo shows empty text for the client, category, kind, job number, job description, division, vendor, PO number and branch names when the linked record is missing. The client name comes through the expense's job. The PO number (purchase_order_number) comes from the linked purchase order.

## Expense with no currency shows as CAD at rate 1

Lines 40-44

Asked as: expense currency defaults to CAD; expense exchange rate 1; unknown currency on expense; currency symbol CAD expense

When an expense has no matching currency record (currency), Turbo shows the currency code (currency_code) as 'CAD', the symbol (currency_symbol) as 'CAD' and the rate (currency_rate) as 1. The icon and the rate date are empty.

## Person names on an expense are given name plus surname

Lines 54-58

Asked as: expense owner name wrong; approver name on expense; expense name shows blank; who created expense name

Turbo shows the expense owner (uid_name), approver, rejector, settler and creator as given name, a space, then surname. When there is no matching profile, the name is empty.

## Expense list keeps expenses with missing linked records

Lines 60-75

Asked as: expense missing from list; expense disappears from list; expense gone after job deleted; orphan expense not listed

Turbo reads an expense with left joins to its job, client, vendor, division, category, kind, people, purchase order, currency, branch and attachment document. An expense still appears when any of these links is empty or points to a record that is gone.

## Expense client comes only through its job

Lines 61-62

Asked as: expense client not shown; client missing on expense; expense shows no client; which client is expense for

Turbo finds an expense's client from the client of the expense's job (j.client). An expense with no job shows an empty client name.

## Expense people are matched by user ID

Lines 67-74

Asked as: expense wrong person name; expense owner profile mismatch; expense user name incorrect; people on expense wrong

Turbo matches the expense owner (uid), approver, rejector, settler and creator to profiles by user ID (uid). The names it shows come from those profiles.
