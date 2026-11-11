<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/expense_vendor_invoice.go","sha256":"86279e556047a7cb3e879230f156d908c1ebf11e03b2a49fd4c4cb5d673312db","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/expense_vendor_invoice.go

## Vendor invoice details are stored on the expense document

Lines 20-26

Asked as: where invoice number saved; invoice details on the file; same invoice several expenses; invoice date stored where

Turbo keeps the vendor invoice number, date and total on the expense document (the attachment), not on the expense. Every expense that uses the same document shares these details. The expense forms accept them as extra inputs (vendor_invoice_number, vendor_invoice_date, vendor_invoice_total) next to the attachment.

## Invoice number must contain letters or digits

Lines 51-61

Asked as: invoice number only symbols; invoice number with dashes; invoice number rejected spaces; invoice number wont accept

Turbo matches invoice numbers by a key of lowercase letters and digits only, so "INV-0042", "inv 0042" and "INV0042" match. Leading zeros stay. A number with no letters or digits is refused with "The invoice number must contain letters or digits."

## Invoice details are optional inputs that keep current values

Lines 76-96

Asked as: blank invoice field erased; invoice date cleared by mistake; leave invoice number empty; invoice details overwritten

Vendor invoice number, date and total are read from the expense request. An input left out keeps the document's current value. Number and date are trimmed of spaces.

## Invoice number is trimmed of spaces

Lines 79-86

Asked as: extra spaces invoice number; space before invoice number; invoice number trailing space

Turbo trims spaces from the start and end of the invoice number and the invoice date that a request sends.

## Invoice total must be a number

Lines 87-131

Asked as: invoice total has letters; invoice amount text rejected; invoice total wont save; invoice total not numeric

If the invoice total (vendor_invoice_total) is not a number, Turbo refuses with "Enter the invoice total as a number." An empty total counts as 0.

## Invoice details need an attached invoice first

Lines 139-148

Asked as: invoice details without file; no file invoice number; need receipt before invoice; cant enter invoice without attachment

When an expense has no attached document (attachment_document), Turbo refuses invoice details (number, date or a non-zero total). It shows "Attach the invoice before entering its details." Empty inputs pass.

## Expense save fails if its attached document is missing

Lines 149-153

Asked as: expense wont save attachment; attached file missing save; save fails lost attachment; attachment gone expense error

When an expense has an attached document (attachment_document), Turbo loads that document before it checks the invoice details. If the document cannot be found, the expense save fails and nothing is written. This failure has no invoice-specific message.

## Other expenses on the invoice exclude the one being saved

Lines 154-163

Asked as: invoice check counts itself; my expense conflicts itself; invoice total includes current expense

Turbo checks the invoice rules (editing, vendor, currency, total) against the other expenses that use the same document. The expense being saved is left out of that list. The edit lock, the vendor and currency match, and the total check all depend on it.

## Who may change shared invoice details

Lines 165-171

Asked as: cant edit invoice uploaded by someone; only uploader change invoice; invoice locked to uploader; someone else's invoice details

When other expenses use the same invoice document, only the person who uploaded it or a bookkeeper (the bookkeeper claim) can change its details. Anyone else sees "Only the person who uploaded this invoice or a bookkeeper can change its details." When no other expense uses the document, anyone may change them.

## Invoice checks run in a fixed order

Lines 173-190

Asked as: which invoice error first; different invoice error message; order of invoice checks

For an expense with an invoice number, Turbo checks in this order: the vendor and currency, then that the number has a letter or digit, then that the vendor has not already entered it, then that the expenses do not exceed the total. Turbo reports the first check that fails.

## Invoice number must contain a letter or digit

Lines 177-184

Asked as: invoice number only punctuation; invoice number just hyphens; invoice number no letters

When an expense has an invoice number, Turbo builds its key from the letters and digits in it. If none remain, for example "---", Turbo refuses with "The invoice number must contain letters or digits."

## Expense without invoice number is a receipt

Lines 191-194

Asked as: expense without invoice number; remove invoice number make receipt; clear invoice details; no invoice number receipt

If the invoice number is empty after validation, Turbo treats the document as a receipt and clears all its vendor invoice details.

## Invoice records who set details and when

Lines 196-216

Asked as: who changed invoice details; invoice last edited by; invoice change history; when invoice was updated

When invoice details change, Turbo saves them on the document with the user who last set them and the time. It keeps no history, and it does not write when nothing changed.

## Invoice document is not updated when nothing changes

Lines 196-198

Asked as: saved but date unchanged; invoice updated date not changing; no change invoice saved

If the invoice details are the same after the request, Turbo does not write the document. The recorded-by person and the date stay as they were.

## Details with an attached document keep recorded-by even without a signed-in user

Lines 199-215

Asked as: edited by blank; no name on invoice change; invoice changed by nobody

Turbo saves the signed-in user's ID as the person who last set the details. If the request has no signed-in user, Turbo leaves it blank.

## Another expense on the document blocks other staff from changing details

Lines 243-247

Asked as: other staff cant change invoice; invoice edit blocked by others; invoice used by another expense

When another expense uses the same invoice document, only the person who uploaded it or a bookkeeper can change the details. Anyone else gets "Only the person who uploaded this invoice or a bookkeeper can change its details."

## Bookkeeper may fill blank invoice details despite lock

Lines 248-251

Asked as: bookkeeper add invoice number; fill empty invoice details; bookkeeper blank invoice locked

A bookkeeper can fill in an invoice that has no number, no date and a zero total, even if another expense on it is submitted. Existing details stay locked for everyone. The change must still be allowed by the uploader/bookkeeper check.

## Invoice details lock after another expense is submitted

Lines 252-258

Asked as: invoice locked after submit; cant change invoice submitted; invoice frozen for approval; rejected expense unlocks invoice

Invoice details cannot change while another expense on the same invoice is submitted and not rejected. Turbo shows "These invoice details cannot change: another expense on this invoice has been submitted for approval." Rejected expenses do not lock the details.

## Invoice number length limit

Lines 264-266

Asked as: invoice number too long; invoice number character limit; invoice number max length

The invoice number may have at most 100 characters. Otherwise Turbo shows "Keep the invoice number to 100 characters or fewer."

## Invoice date format

Lines 267-271

Asked as: invoice date wrong format; invoice date not accepted; invoice date format error

The invoice date must be a valid date in YYYY-MM-DD form. Otherwise Turbo shows "Enter a valid invoice date."

## Invoice total cannot be negative

Lines 272-274

Asked as: negative invoice amount; minus invoice total; invoice total below zero

A negative invoice total is refused with "The invoice total cannot be negative."

## Invoice details are all or nothing; required on account

Lines 276-298

Asked as: on account expense invoice needed; all invoice fields required; on account needs invoice details; invoice details missing on account

An on-account expense (payment_type OnAccount) needs invoice number, date and a total above zero. Any other expense needs all three once any one is entered. Missing fields get "Enter the vendor's invoice number.", "Enter the invoice date." or "Enter the invoice total.". The overall message is "An on-account expense needs its vendor invoice details." for on-account, else "Enter the vendor invoice details."

## Invoice vendor and currency come from the expense

Lines 305-309

Asked as: invoice currency set from expense; invoice vendor changed automatically; where invoice vendor comes from

While the expense is the only one on the document (or the document has no vendor), Turbo sets the invoice's vendor and currency from the expense. Once the invoice is shared, every expense on it must use the document's vendor and currency.

## Invoice needs a vendor

Lines 310-312

Asked as: no vendor on invoice; select vendor invoice error; invoice vendor missing

If no vendor is known for the invoice, Turbo shows "Select the vendor that issued this invoice."

## Vendor mismatch on shared invoice

Lines 313-321

Asked as: different vendor same invoice; wrong supplier on invoice; invoice from another vendor

If an expense on the invoice has a different vendor, Turbo shows "This invoice is from <vendor name>. Every expense on it must use that vendor."

## Currency mismatch on shared invoice

Lines 322-325

Asked as: wrong currency on invoice; invoice currency differs; expense currency mismatch invoice

If an expense on the invoice has a different currency, Turbo shows "Every expense on this invoice must use the invoice's currency."

## Duplicate invoice number for a vendor is refused

Lines 330-351

Asked as: invoice already entered; duplicate invoice number vendor; invoice number used before; vendor invoice already exists

Two documents cannot have the same vendor and invoice number key. Turbo shows "Invoice <number> from this vendor has already been entered", adding "against PO <po number>" and "by <name>" of the first expense's creator when known.

## Expenses cannot exceed invoice total

Lines 364-383

Asked as: expenses exceed invoice total; claims over invoice amount; expenses more than invoice

The expenses on one invoice cannot total more than the invoice total. Rejected expenses are not counted, and amounts are compared in whole cents, so an equal total passes. Otherwise Turbo shows "The expenses on this invoice would total X, more than the invoice total of Y."
