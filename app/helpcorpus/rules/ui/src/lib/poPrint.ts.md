<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poPrint.ts","sha256":"8c5fd586cf8bf4334fa65eba86a6feabb950b62c0cdd7d7554b565476cf9935e","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poPrint.ts

## Amount format accepted for printing a purchase order

Lines 16-20

An amount to print must be digits with an optional decimal part of one or two digits (for example 100 or 100.5). Turbo rejects anything else, such as negative numbers, thousands separators or three decimals. Turbo converts the amount to integer cents (amountInCents).

## Default print amount is the full print limit

Lines 27-34

When a person opens print options for a purchase order, the amount defaults to the print limit (print_max_amount) with two decimals. The amount is blank if the limit is zero, negative or invalid. Tax and shipping are both unchecked by default.

## Only Active purchase orders can be printed

Lines 37-37

Turbo refuses to print a purchase order whose status is not Active. It shows "Only Active purchase orders can be printed." This check runs before the others.

## No approved balance blocks printing

Lines 38-39

If the print limit is missing, invalid, zero or negative, Turbo refuses to print. It shows "No approved balance remains for printing."

## Print amount must be greater than zero with valid format

Lines 40-42

Turbo refuses a print amount that is invalid or zero or less. It shows "Enter an amount greater than zero with up to two decimal places."

## Print amount cannot exceed the available approved amount

Lines 43-43

Turbo refuses a print amount that is higher than the print limit. It shows "The amount to print exceeds the available approved amount." An amount equal to the limit passes this check.

## Tax or shipping needs an amount below the limit

Lines 44-46

When tax or shipping is added, the print amount must be lower than the print limit. An amount equal to the limit is refused with "Enter a lower amount to allow for the additional charges." Without tax or shipping, an amount equal to the limit is allowed.

## Money shown in print uses currency code, two decimals

Lines 50-58

Printed amounts use the en-CA format with the currency code and exactly two decimals. The currency falls back to CAD when the purchase order has no currency code.

## Authorized amount wording on the printout

Lines 60-73

The authorized amount text is the amount plus " plus tax", " plus shipping" or " plus tax and shipping", depending on the boxes checked. For a Recurring purchase order, it adds " / " and the frequency. The frequency falls back to "period" when it is blank. An invalid or zero amount gives "Enter a valid amount."

## Warning about tax and shipping allowance

Lines 75-87

When tax or shipping is checked and the options are valid, Turbo warns that the charges must not exceed the print limit minus the print amount. It says "Tax", "Shipping" or "Tax and shipping combined", and adds "per period" for Recurring purchase orders. The warning says that exceeding it requires a new PO approved before authorizing the extra cost. With no tax or shipping, or with invalid options, there is no warning.

## Acknowledgement resets when print terms change

Lines 89-100

The acknowledgement of the warning is tied to the amount, tax, shipping, print limit, currency, type and frequency. If any of them changes, the acknowledgement no longer applies.

## Print link carries the limit and acknowledgement

Lines 102-111

The print options passed to the print page include the amount, tax and shipping flags (1 or 0), and the print limit. The print page checks that limit against a fresh read of the purchase order. Acknowledged is 1 only when tax or shipping is checked.

## Malformed print options are rejected

Lines 113-129

If any print option is missing or repeated, or tax, shipping or acknowledged is not 0 or 1, Turbo shows "Invalid print options. Open Print options again." An invalid limit shows "Invalid print limit. Open Print options again."

## Print options are validated again when read

Lines 115-141

With no print options given, Turbo uses the defaults. Any options read are validated with the same rules as in the popup, and the first error is shown.

## Tax or shipping needs acknowledgement and an unchanged limit

Lines 130-137

When tax or shipping is checked, the warning must be acknowledged, or Turbo shows "Acknowledge the additional-charge warning in Print options." If the limit in the options differs from the current print limit, Turbo shows "The available amount has changed. Open Print options again and review the warning."

## Print link holds only an opaque token

Lines 144-164

The print page URL holds only a print token, not the options. Turbo refuses a link with amount, plusTax, plusShipping, limit or acknowledged in it, and shows "Invalid print link. Open Print options again." The token must be given once and be 36 characters of hex digits and hyphens, or Turbo shows the same message.

## Print options exist only in the tab that created them

Lines 160-168

Turbo reads the saved print options for a purchase order and token from browser storage. If they are missing, for example in another tab, it shows "Print options are not available in this tab. Open Print options again." Without a print token, the default options apply.
