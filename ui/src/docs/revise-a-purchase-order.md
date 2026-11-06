# Revise a purchase order

Change the total or currency of an approved purchase order and keep its number. It needs approval again.

```flow
You | Revise the total | Enter the new total or currency.
Approver | Approve or reject | Check the new total and currency.
You | Add expenses | Use the same number after final approval.
Note: Rejected? Edit and resubmit, or cancel.
```

## Revise the total

1. Open the purchase order.
2. Select **Revise Total**.
3. Enter the new total or select a new currency.
4. Select the approvers that Turbo asks for.
5. Select **Save**.

Only the person who created the purchase order can revise it. If Turbo doesn't ask for an approver, you are the approver. Approve it yourself.

## When you can't revise

You can't revise a purchase order if:

- It has expenses.
- It has child purchase orders.
- It shows **Manually created**.
- It has a closed job. **Save** shows an error.

To change other details on a purchase order with no expenses, ask payables staff to cancel it. Then create a new one.

## If it has expenses

Check its **Type**:

- **Cumulative**: Add the expense. If it goes over, select **Create Child PO**. Set its total to the full expense. One expense can't use two purchase orders.
- **One-Time**: Ask payables staff to [convert it to Cumulative](/docs/purchase-orders#convert-a-one-time-po). Then add the expense.
- **Recurring**: Ask payables staff to close it after its last commit. Then create another.

## While it is under revision

- It shows **Under revision** to everyone.
- Select **Under revision** to see who must act next.
- Nobody can add expenses or print it.

## Change a revision before approval

1. Select **Edit Purchase Order**.
2. Change the total, currency or approvers.
3. Select **Save**.

A change clears any first approval. To go back to the old total, enter it here. It still needs approval.

You can't cancel it now. Ask the approver to reject it. Then select **Cancel Purchase Order**.

## Approve a revision

1. Under **Purchase Orders**, open **Pending My Approval**.
2. Open the purchase order and check the new total and currency.
3. Select **Approve** or **Reject**.

A new total can require a second approval. After final approval, the purchase order is Active again. If the vendor has the old total, print a new copy.

## If a revision is rejected

Under **Purchase Orders**, open **Rejected (Mine)**. Open the purchase order and read the reason. Then:

- Select **Edit and Resubmit** to change it and send it again.
- Select **Cancel Purchase Order** to end it. You can't undo this.

You can't delete it.
