# Expenses

An expense is the accounting record of a purchase. It records the actual amount spent. If a [PO](/docs/purchase-orders) covers the purchase, you must enter the expense directly against that PO. The expense then records spending against the PO. One PO can have several expenses.

```flow
You | Enter the expense | Start from an Active PO, or from **My Expenses** if you don't need a PO.
You | Submit it | Saving doesn't submit it.
Your manager | Approve or reject | Check the expense and its attachment.
Committer | Commit it | Turbo locks the expense.
Note: Rejected? Select **Recall**, edit the expense and submit it again.
```

## When you need a PO

- **With a job:** you need a PO, except for the other payment types below.
- **Without a job:** you need a PO at or above the no-PO limit, except for the other payment types below. The default limit is $100.

## Payment types for purchases

| Payment type          | Use it for                     | Reimbursed |
| --------------------- | ------------------------------ | ---------- |
| On Account            | A charge to a vendor account   | No         |
| Expense               | A purchase that you paid for   | Yes        |
| Corporate Credit Card | A purchase on a corporate card | No         |

A PO uses one of these types. An On Account expense needs its [vendor invoice](/docs/vendor-invoices) details.

## Other payment types (these never need a PO)

| Payment type           | Use it for                                           | Reimbursed |
| ---------------------- | ---------------------------------------------------- | ---------- |
| Fuel Card              | Fuel on a fuel card                                  | No         |
| Allowance              | Meals and lodging. Turbo sets the amount.            | Yes        |
| Mileage                | Travel in your vehicle. Turbo calculates the amount. | Yes        |
| Personal Reimbursement | A personal cost, if allowed for you                  | Yes        |

## Enter an expense from a PO

1. Select **Create Expense** on an Active PO. Turbo fills the form from the PO.
2. Change the payment type, vendor, description or category if necessary.
3. Enter the total and add an attachment.
4. Select **Save**.
5. Select **Submit**.

You can't change the job, currency or branch.

## Enter an expense without a PO

1. Select the plus button beside **My Expenses**.
2. Select the payment type. You can select a job only for one of the other payment types.
3. Enter the details and add an attachment. Allowance, Mileage and Personal Reimbursement expenses don't need one.
4. Select **Save**.
5. Select **Submit**.

## Category

A category groups costs within one job. It is optional, and it needs a job. Turbo shows the field only when the job has categories.

- **With a PO:** Turbo copies the PO's category. You can select another category of the PO's job.
- **Without a PO:** select the job first. Then you can select one of its categories.

## Branch

You don't select the branch of an expense. Turbo sets it when you save:

- **With a PO:** the PO's branch.
- **Without a PO, with a job:** the job's branch.
- **Without a PO or a job:** your default branch. An admin sets this on your admin profile.

## Approve an expense

1. Under **Expenses**, open **Pending My Approval**.
2. Open the expense and check it.
3. Select **Approve**, or select **Reject** and enter a reason.

## Change an expense

Select **Recall** to change a submitted expense before approval, or a rejected expense. Then edit it and submit it again. You can't change a committed expense.

## Expense statuses

| Status    | What happens next                     |
| --------- | ------------------------------------- |
| Draft     | You edit and submit it.               |
| Submitted | Your manager approves or rejects it.  |
| Approved  | A committer commits or rejects it.    |
| Rejected  | You recall, edit and submit it again. |
| Committed | It is final. It can close its PO.     |

## Before and after commit

Payables staff must settle a foreign-currency On Account or Corporate Credit Card expense before commit.

**My Reimbursements** shows committed expenses that payroll repaid to you.
