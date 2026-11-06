# Client invoices

A client invoice records an invoice that the company sends to a client. Invoicing staff create it and link the billed-back costs to it. One invoice can cover several jobs.

```flow
Invoicing staff | Create a draft | Take the number from the Teams group as usual.
Invoicing staff | Issue the invoice | Do this when you send it to the client.
Invoicing staff | Add billed-back costs | Select each cost and enter the amount billed.
```

## Create an invoice

1. Open **Invoices** under **Business**.
2. Select **New invoice**.
3. Select the **Client**.
4. Enter the **Number**, **Date** and **Total (CAD)**.
5. Select **Save draft**.

A number has the form YYMMNNN, such as 2609001. Turbo refuses a number that another invoice uses. A draft can have no number.

## Change or delete a draft

Open the draft. Select **Edit** to change it. Select **Delete draft** to delete it, then select **Delete** to confirm. Its number becomes free for another invoice.

## Issue an invoice

Do this when you send the invoice to the client. The invoice needs a number, date and total.

1. Open the draft.
2. Select **Issue**.
3. Select **Issue** again to confirm.

You can't edit or delete an issued invoice.

## Bill a cost from the invoice

1. Open the issued invoice.
2. Select **Add billed-back cost**.
3. Select the **Job** and the **Cost**.
4. Enter the **Amount billed** in Canadian dollars, before tax.
5. Select **Add**.

Turbo lists unbilled committed expenses and Active or Closed POs on the client's jobs.

## Bill a cost from an expense or PO

1. Open the committed expense, or the Active or Closed PO.
2. Select **Bill on invoice** under **Billed to client**.
3. Select the **Invoice**. Turbo lists only issued invoices.
4. Enter the **Amount billed** in Canadian dollars, before tax.
5. Select **Bill**.

## Remove a billed-back cost

1. Open the invoice.
2. Select **Remove** on the cost.
3. Select **Remove** again to confirm.

You can then bill the cost again.

## See the invoice for a cost

Open the expense or PO. **Billed to client** shows the invoice, its date and the amount billed.
