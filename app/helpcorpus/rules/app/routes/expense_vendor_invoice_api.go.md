<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_vendor_invoice_api.go","sha256":"a8d31084ecfd7b9031c62b9dc31efd39be37c2cf0b1d43fc03d374ca9d42c9ce","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_vendor_invoice_api.go

## Vendor invoice view shows invoice details and each expense's status

Lines 16-42

Asked as: vendor invoice details; expense status on invoice; invoice expenses list; allocated amount on invoice

The vendor invoice view shows the invoice number, date, total, currency code, who recorded it and when. It lists each expense on the invoice with its PO number, date, total, who it is for, who entered it, and whether it is submitted, approved, committed or rejected. It also shows the allocated total (allocated) and whether the viewer can open each expense (can_open).

## Vendor invoice view lists every expense sharing the invoice

Lines 44-77

Asked as: all expenses same invoice; expenses sharing invoice; invoice number lists expenses; other expenses on this invoice

Turbo shows the invoice number, date, total, currency, who recorded it and when. It lists every expense that uses the same invoice document, ordered by expense date and then by creation time. Each part shows its PO number, date, total, who the expense is for, who entered it, and its submitted, approved, committed and rejected states.

## Parts the viewer cannot open are listed without a link

Lines 58-70

Asked as: expense listed but no link; cannot open expense on invoice; greyed out expense; can't click invoice expense

A part can be opened (can_open) when the viewer passes the same visibility check as the expense list and details. A holder of the admin claim can also open any committed expense. Other parts are still listed but cannot be opened.

## Vendor invoice lookup by expense requires sign-in

Lines 83-86

Asked as: lookup invoice behind expense; not logged in invoice; unauthorized invoice lookup; invoice needs sign in

Looking up the vendor invoice behind an expense requires a signed-in user. Turbo refuses an unauthenticated request with "unauthorized".

## Who can open an expense's vendor invoice

Lines 98-106

Asked as: who can see vendor invoice; expense not found or not authorized; can't open expense invoice; admin see committed expense

Turbo returns the vendor invoice only for an expense the viewer can see, using the same visibility rule as the expense list and details (expenseVisibilityPredicate). Staff with the admin claim can also see any committed expense (committed is not empty). Otherwise Turbo answers "expense not found or not authorized".

## Expense with no vendor invoice returns an empty result

Lines 107-109

Asked as: expense has no invoice; no vendor invoice attached; empty invoice result; missing invoice file

When the expense the viewer may see has no vendor invoice document (attachment_document), Turbo returns an empty invoice with no parts. It does not give an error.

## Vendor invoice that cannot be loaded shows an error

Lines 110-114

Asked as: failed to load the invoice; invoice won't load; invoice error on expense details; invoice details error

When Turbo finds the expense's vendor invoice document but cannot load the invoice details or its expenses, it answers "failed to load the invoice". It returns no partial view. A person sees this error on the expense details lookup. The file hash lookup shows it too.

## Allocated total excludes rejected expenses

Lines 131-137

Asked as: allocated total wrong; rejected expenses in total; allocated amount math; why allocated lower

The allocated amount (allocated) on a vendor invoice adds up the totals of the expenses on it that are not rejected. Turbo sums them in cents and shows the result in currency units. Rejected expenses are listed but do not count.

## Only bookkeepers can look up a vendor invoice by file hash

Lines 148-158

Asked as: search invoice by file hash; book_keeper claim required; bookkeeper invoice lookup; file hash access denied

Looking up an invoice by file hash requires sign-in and the bookkeeper claim (book_keeper). Without sign-in Turbo answers "unauthorized". Without the claim it answers "book_keeper claim required".

## File hash lookup returns empty when the file is new

Lines 169-178

Asked as: file hash no results; new file not found; invoice hash empty result; file not on any expense

The file hash lookup matches an expense document by its SHA-256 hash (attachment_hash), and only when at least one expense uses it. If no document matches, Turbo returns an empty invoice with no parts and no error. This tells the editor that the file is new or is on no expense.

## File hash lookup shows all expenses on the invoice

Lines 179-183

Asked as: hash shows all expenses; find invoice by file checksum; bookkeeper sees hidden expenses; same file other expenses

When the hash matches, Turbo returns the invoice with every expense that shares it. This is not limited to expenses the bookkeeper can open. Parts the bookkeeper cannot open are listed without a link, by the usual visibility check.
