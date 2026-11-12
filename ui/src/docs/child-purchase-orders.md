# Child purchase orders

A child PO covers an expense that takes a Cumulative PO over its total. It is always a One-Time PO, never Cumulative or Recurring, so it takes one expense. Its number is its parent's number plus a suffix, such as 2611-0379-01.

## Create a child PO

1. Save the expense on the Cumulative PO.
2. Select **Create Child PO**.
3. Set its total to the full expense.
4. Select an approver.
5. Select **Save**.

Turbo copies the vendor, job, category, description, kind and payment type from the parent. They must match the parent.

## Set the total to the full expense

One expense can't use two POs. It can't use the amount that remains on the parent and take the rest from the child PO. The child PO must cover the full expense.

For example, the parent has $200 left and the expense is $500. Set the child PO's total to $500, not $300. The parent keeps its $200.

## After approval

Record the expense on the child PO.

## Limits

- A PO can have one open child PO at a time.
- A child PO can't have a child PO.
- You can't [revise](/docs/revise-a-purchase-order) a PO that has child POs.
