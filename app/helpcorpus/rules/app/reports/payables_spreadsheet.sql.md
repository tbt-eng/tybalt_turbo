<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/payables_spreadsheet.sql","sha256":"aadaf3fcac55ab261e114fae84c0195b7fd36ab0f69c3e0e1ea7b5facd9ae01a","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/payables_spreadsheet.sql

## Missing job, division, branch, type, vendor or owner show as blank

Lines 3-6

The payables spreadsheet shows an empty value for the job number, division code, branch code, PO type, vendor name and employee (the PO owner's name) when the purchase order has none. The purchase order is still listed.

## Approval date is the later of the two approvals

Lines 8-12

The approval date (approval_date) is the second approval date (second_approval) when it is set and later than the first approval date (approved). Otherwise it is the first approval date (approved).

## Currency defaults to CAD in the payables spreadsheet

Lines 14-14

When a purchase order has no matching currency, the currency column shows 'CAD'.

## Vendor and employee blank when no match

Lines 17-18

The vendor name and the employee name (the owner's given name and surname) are empty when no vendor or profile is found.

## Approved by shows the second approver when there is one

Lines 19-23

The approved-by column shows the second approver (second_approver) when the purchase order has one. Otherwise it shows the approver (approver). The name is the given name and the surname. It is empty when no profile is found.

## Payables spreadsheet leaves out unapproved purchase orders

Lines 34-34

The payables spreadsheet lists only purchase orders whose status (status) is not 'Unapproved'. Purchase orders with every other status appear if they meet the other conditions.

## Payables spreadsheet needs a PO number

Lines 35-35

The payables spreadsheet skips a purchase order whose PO number (po_number) is empty.

## Payables spreadsheet keeps only PO numbers below 5000

Lines 36-39

The payables spreadsheet reads the number after the first hyphen in the PO number as an integer. It keeps the purchase order only if that number is under 5000, so 5000 itself is excluded. This drops legacy or control numbers such as YYMM-5XXX to YYMM-9XXX. A PO number with no hyphen is always kept.

## Payables spreadsheet needs an approval date

Lines 40-40

The payables spreadsheet includes a purchase order only if it has a first approval date (approved) or a second approval date (second_approval). A purchase order with neither is left out.
