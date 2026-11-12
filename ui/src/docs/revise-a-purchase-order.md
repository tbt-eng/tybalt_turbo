# Revise a purchase order

Change the total, currency, description, division, job or category of an Active purchase order that has no expenses, and keep its number. It needs approval again.

```flow
You | Make a revision | Enter the changes.
Approver | Approve or reject | Check the changes.
You | Add expenses | Use the same number after final approval.
Note: Rejected? Edit and resubmit, or cancel.
```

## Make a revision

1. Open the purchase order.
2. Select **Revise**.
3. Change the total, currency, description, division, job or category.
4. Select the approvers that Turbo asks for.
5. Select **Save**.

Only the purchase order's owner can revise it. If Turbo doesn't ask for an approver, you are the approver. Approve it yourself.

## When you can't revise

You can't revise a purchase order if:

- It has expenses.
- It has child purchase orders.
- It shows **Manually created**.
- It came from a PO group.
- It has a closed job and you keep that job.

With expenses, you can't change its details. Otherwise, ask payables staff to cancel it. Then create a new one.

## If it has expenses

To spend more, check **Type**:

- **Cumulative**: Add the expense. If it goes over, create a [child PO](/docs/child-purchase-orders).
- **One-Time**: Ask payables staff to [convert it to Cumulative](/docs/purchase-orders#convert-a-one-time-po). Then follow **Cumulative**.
- **Recurring**: Ask payables staff to close it after its last commit. Then create another.

## Change the job

Select the new job, then its division and category. For a project, confirm the budget again. You can't add a job to a purchase order that has none, remove its job, or change the job of a child purchase order.

## Change the division or category

The division must have hours allocated to the job. The category must belong to the job.

A child purchase order keeps its parent's description and category.

## While it is under revision

- It shows **Under revision** to everyone.
- Select **Under revision** to see who must act next.
- Nobody can add expenses or print it.

## Change a revision before approval

1. Select **Edit Purchase Order**.
2. Change the details or replace the attachment.
3. Select **Save**.

You can replace the attachment only here. A change clears any first approval. To go back to the old total, enter it here. It still needs approval.

You can't cancel it now. Ask the approver to reject it. Then select **Cancel Purchase Order**.

## Approve a revision

1. Under **Purchase Orders**, open **Pending My Approval**.
2. Open the purchase order and check the changes.
3. Select **Approve** or **Reject**.

The email shows the previous values. A new total can require a second approval. After final approval, the purchase order is Active again. If the vendor has the old total, print a new copy.

## If a revision is rejected

Under **Purchase Orders**, open **Rejected (Mine)**. Open the purchase order and read the reason. Then:

- Select **Edit and Resubmit** to change it and send it again.
- Select **Cancel Purchase Order** to end it. You can't undo this.

You can't delete it.
