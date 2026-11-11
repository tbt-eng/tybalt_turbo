<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poPrint.ts","sha256":"f5d43713691c63b2a328c8601a4a1e34b276907cbeb9e95664815a24c6640144","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poPrint.ts

## Printed PO amount can be lower than the approved amount

Lines 1-26

Asked as: print lower amount; vendor quoted price; print less than approved; approved amount ceiling

The approved amount of a purchase order is a ceiling, not the price. The person who prints can enter a lower amount, usually the vendor's quoted price. The print page shows only that amount, never the print limit, the allowance or the warning. Printing does not change the PO or its approved amount.

## How Turbo sets the print limit of a PO

Lines 23-25

Asked as: where print limit comes from; print limit calculation; recurring PO print limit; cumulative PO remaining amount

The print limit (print_max_amount) comes from the server. It is the PO total, or the amount per period for a Recurring PO. For a Cumulative PO it is the total less committed expenses.

## Default print amount is the full print limit

Lines 42-60

Asked as: default print amount; amount prefilled when printing; why print amount is full; empty print amount

The print options start with the print limit as the amount, with two decimals, and with Plus tax and Plus shipping off. The amount starts empty when the limit is not a valid number or is zero or less.

## Print amount format: up to two decimals

Lines 42-46

Asked as: print amount decimals; cents format rejected; amount with two decimals; comma in print amount

A print amount must be digits with at most two decimal places, such as 12 or 12.50. Any other format is rejected. Turbo converts the amount to integer cents.

## Only Active purchase orders can be printed

Lines 63-63

Asked as: cannot print purchase order; print button not working; PO not active print; print draft PO

Turbo refuses to print a PO whose status is not Active. The message is "Only Active purchase orders can be printed."

## No print when no approved balance remains

Lines 64-65

Asked as: nothing left to print; zero balance print; print blocked no balance; PO fully used printing

Turbo refuses to print when the print limit is missing or is zero or less. The message is "No approved balance remains for printing."

## Print amount must be greater than zero with valid decimals

Lines 66-68

Asked as: print amount zero rejected; negative print amount; invalid print amount error; enter valid print amount

Turbo refuses a print amount that is invalid or zero or less. The message is "Enter an amount greater than zero with up to two decimal places."

## Print amount cannot exceed the print limit

Lines 69-69

Asked as: print more than approved; amount over limit printing; exceeds approved amount; print above max

Turbo refuses a print amount above the print limit. The message is "The amount to print exceeds the available approved amount." An amount equal to the limit passes this check.

## Plus tax or Plus shipping needs an amount below the limit

Lines 70-72

Asked as: plus tax won't let me print; plus shipping same as limit; additional charges amount error; lower amount for tax

When Plus tax or Plus shipping is on, an amount equal to the print limit is refused. The message is "Enter a lower amount to allow for the additional charges." Without a plus box, an amount equal to the limit is allowed.

## Printed authorized amount text

Lines 76-99

Asked as: printed authorized amount text; currency code on print; plus tax wording on PO; recurring frequency printed text

The printed text shows the amount in the PO currency with its code, such as CAD. The currency defaults to CAD when the code is empty. It adds " plus tax", " plus shipping" or " plus tax and shipping" for the boxes ticked. For a Recurring PO it adds " / " and the frequency, or "period" when the frequency is empty. An invalid amount gives "Enter a valid amount."

## Warning states the allowance for tax and shipping

Lines 101-113

Asked as: allowance warning tax shipping; how much extra can I spend; tax and shipping combined warning; remaining allowance on print

When a plus box is ticked and the options are valid, Turbo warns that the charges must not exceed the allowance. The allowance is the print limit less the printed amount. Tax and shipping share it, and the warning says "Tax and shipping combined" when both are ticked. A Recurring PO adds "per period". The warning says a new PO must be approved before authorizing an extra cost above the allowance. No warning shows without a plus box or when the options are invalid.

## Acknowledgement depends on amount, terms, currency and limit

Lines 115-126

Asked as: acknowledge again after change; warning acknowledgement reset; changed amount need reacknowledge; acknowledgement lost

The acknowledgement of the warning is tied to the amount, the plus boxes, the print limit, the currency, the PO type and the frequency. A change to any of them needs a new acknowledgement.

## Print options need acknowledgement when a plus box is ticked

Lines 128-158

Asked as: acknowledge additional charge warning; plus tax print refused; must accept warning to print; print options acknowledgement required

When Plus tax or Plus shipping is ticked, the print options must carry the acknowledgement. Otherwise Turbo refuses with "Acknowledge the additional-charge warning in Print options." The acknowledgement is set to 1 only when a plus box is ticked.

## Print refused when the available amount changed

Lines 133-163

Asked as: available amount changed; print limit changed since options; PO updated before printing; reopen print options changed

With a plus box ticked, Turbo compares the saved limit with the current print limit of the PO. If they differ, it refuses with "The available amount has changed. Open Print options again and review the warning." Without a plus box this comparison does not run.

## Malformed saved print options are refused

Lines 139-155

Asked as: invalid print options; corrupted print settings; print options error; invalid print limit

Each saved print option must appear exactly once. The plus and acknowledged values must be 0 or 1. Otherwise Turbo says "Invalid print options. Open Print options again." An invalid limit gives "Invalid print limit. Open Print options again."

## Print link holds only an opaque key from this tab

Lines 170-194

Asked as: print link invalid; print options not available in tab; print page blank new tab; print link expired

The print link carries one key (print), which must be 36 hexadecimal characters and dashes. A bad key gives "Invalid print link. Open Print options again." Turbo reads the options from this browser tab's storage. When nothing is saved, it says "Print options are not available in this tab. Open Print options again." The limit stays out of the URL so a browser print footer cannot show it.

## Print page falls back to default options without a print key

Lines 175-186

Asked as: print link without options; print page default options; amount in URL refused; print link manually edited

When the print link has no print key, Turbo uses the default options, which print the full print limit, and validates them. Options placed directly in the link, such as amount or limit, are refused with "Invalid print link. Open Print options again."
