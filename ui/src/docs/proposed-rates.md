# Proposed rates

If no active rate sheet has the rates in the client's contract, propose new rates in the project request. The approver approves them with the request.

```flow
You | Propose the rates | Start from an active sheet and change the rates.
Approver | Approve or reject | Compare the rates with the sheet that you started from.
Procurement | Create the job | Turbo activates the new rate sheet for the job.
Note: Rejected? Change the rates and submit the request again.
```

## Find a sheet by rates

1. In the project request, select **Find by rates**.
2. Enter the rates you know. The overtime rate is optional.
3. Select **Use** beside a sheet.
4. Select **Save request**.

Turbo lists the active sheets that have all your rates. If none fits, select **Propose new rates with these rates**. Then select the sheet in **Start from**. Your rates replace its rates.

## Propose new rates

1. In the project request, select **Propose new rates**.
2. Select the sheet in **Start from**.
3. Enter the **Reason**. The **Name** of the new sheet is optional.
4. Change the rates. Turbo marks each changed role.
5. Select **Save request**.

The reason tells why the client needs these rates, such as a contract schedule. Turbo won't save the rates without it.

## Set the overtime rates

1. Enter a multiplier above **Overtime Rate**.
2. Select **×**.

Turbo sets each overtime rate to the rate times the multiplier. A role with no rate stays empty.

## Roles with no rate

Turbo shows every role. If the starting sheet has no rate for a role, it shows **No rate**. Enter its rates before you submit. Turbo saves the request without them. Turbo won't save the request while a role has only one rate. That role shows **No rate** or **No overtime rate**. Enter both rates, or leave both empty.

## Rates that match a sheet

Turbo won't save rates that match an active sheet. It shows the code of that sheet. Select **Use** and the code to choose that sheet instead.

## Rates close to a sheet

As you change rates, Turbo lists up to three active sheets that differ in five roles or fewer. A role with no rate doesn't count. Select a sheet's code to see its rates. If one fits the client's contract, select **Use** and its code. Turbo then saves the request with that sheet, without the proposed rates.

## Approve proposed rates

The approver sees each rate beside the starting sheet's rate. Changed roles show **Changed**. Turbo won't approve the request while a role shows **No rate**.

Below the rates, the approver selects **I approve the creation of this new rate sheet and its rates**. Approval needs it.

After a recall or a rejection, the approver must approve the rates again.

## After approval

Procurement sees who approved the rates. They create the job and can rename the new sheet first. Turbo then activates the sheet, and the job uses it. Other projects can then use it too.

Until then, no **Rate Sheet** list shows the sheet. A rejected request keeps its proposed rates. If you delete the request, Turbo deletes them.
