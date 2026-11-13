# Proposed rates

If no active rate sheet has the rates in the client's contract, propose new rates in the project request. The approver approves them with the request.

```flow
You | Propose the rates | Start from an active sheet and change the rates.
Approver | Approve or reject | Compare the rates with the sheet that you started from.
Procurement | Create the job | Turbo activates the new rate sheet for the job.
Note: Rejected? Change the rates and submit the request again.
```

## Propose new rates

1. In the project request, select **Propose new rates**.
2. Select the sheet in **Start from**.
3. Enter the **Name** of the new sheet and the **Reason**.
4. Change the rates. Turbo marks each changed role.
5. Select **Save request**.

The reason tells why the client needs these rates, such as a contract schedule. You must enter it before you submit.

## Rates that match a sheet

Turbo won't save rates that match an active sheet. It names that sheet. Select **Use** and the sheet name to choose that sheet instead.

## Approve proposed rates

The approver sees each rate beside the rate of the sheet that the requester started from. Changed roles show **Changed**.

In the checklist, the approver confirms that the client accepts the proposed rates.

## After approval

Procurement creates the job and can change the name of the new sheet first. Turbo then activates the sheet, and the job uses it. Other projects can then use it too.

Until then, no **Rate Sheet** list shows the sheet. A rejected request keeps its proposed rates. If you delete the request, Turbo deletes them.
