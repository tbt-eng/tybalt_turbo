# Purchase orders

A purchase order (PO) is permission for the company to make a purchase. When the approver approves a PO, its owner can make the purchase and bind the company to it. Another person who acts for the company can also use the PO. A PO isn't an accounting record. An [expense](/docs/expenses) records the actual spending.

```flow
You | Create the PO | Enter the details and select an approver.
Approver | Approve or reject | Check the need for the purchase and the details.
Second approver | Approve or reject | Only some POs need this. This person's limit covers the total.
You | Make the purchase | The PO is now Active and has a PO number.
You | Record the expense | Select **Create Expense** on the PO.
Note: Rejected? Select **Edit and Resubmit**.
```

## Create a PO

1. Select the plus button beside **My Purchase Orders**.
2. Select the type, kind and payment. Enter the details and total.
3. Select the **Approver** and, if shown, the **Priority Second Approver**. Turbo lists only people who can approve this PO.
4. Select **Save**. Turbo sends the PO for approval.

[PO Approvers](/pos/approval-limits) shows the current limits, if you have access.

## PO types

| Type       | Use it for                          | Turbo closes it when                       |
| ---------- | ----------------------------------- | ------------------------------------------ |
| One-Time   | One expense                         | A committer commits its expense.           |
| Recurring  | Repeated expenses between two dates | A committer commits all expected expenses. |
| Cumulative | Several expenses up to a total      | Committed expenses reach the total.        |

## Branch and category

- **With a job:** the PO uses the job's branch. You can select one of the job's categories.
- **Without a job:** select the branch. Turbo starts with your default branch. Some branches need permission.

Expenses against the PO use its branch.

## Approve a PO

1. Under **Purchase Orders**, open **Pending My Approval**.
2. Open the PO and check it.
3. Select **Approve**, or select **Reject** and enter a reason.

## Change a PO

You can edit or delete your PO until it has final approval. If you make a meaningful change, the approvals reset and the PO goes back for approval.

After final approval, you can't delete your PO. If it has no expenses, you can [revise its total](/docs/revise-a-purchase-order).

## Spend more than the PO total

- **One-Time or Recurring:** an expense can go over by a small amount. The default is the lesser of 5% or $100.
- **Cumulative:** the expenses together can't go over the total.

## Convert a One-Time PO

Payables staff can change a One-Time PO to Cumulative. It then accepts several expenses, up to its total.

1. Under **Purchase Orders**, open **Search**.
2. Find the PO and open it.
3. Select **Convert to Cumulative**.
4. Review the changes.
5. Select **Convert**.

You can also convert a One-Time PO that Turbo closed since March 12, 2026. Its committed expenses must total less than the PO total.

## PO statuses

| Status         | What it means                                                               |
| -------------- | --------------------------------------------------------------------------- |
| Unapproved     | It waits for approval, or the approver rejected it.                         |
| Under revision | Its creator revised it. It waits for approval, or the approver rejected it. |
| Active         | You can buy and record expenses against it.                                 |
| Closed         | It is complete. Turbo or payables staff closed it.                          |
| Cancelled      | Payables staff cancelled it before it had any expenses.                     |
