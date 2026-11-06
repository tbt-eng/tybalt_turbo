# Proposals

A proposal is an offer of work to a client. Only procurement can create one. When the client awards the work, the proposal can lead to a project.

```flow
Procurement | Create the proposal | Enter the details. Turbo assigns a number that starts with P.
Procurement | Set it to Submitted | Do this after the company submits the proposal to the client.
Procurement | Award it | Do this after the client awards the work.
You | Create a project request | Start the request from the proposal.
Note: Client said no? Procurement sets the status to **Not Awarded**.
```

## Create a proposal

Only procurement can do this. If you need a proposal, ask procurement.

1. Select the plus button beside **Jobs**.
2. Select **Create a proposal**.
3. Enter the details, both proposal dates and at least one division.
4. Select **Save**. Turbo assigns the number.

## Find a proposal

1. Open **Jobs**.
2. Select **Proposals**.
3. Search for the number, client or description.

## Change a proposal

Procurement opens the proposal and selects **Edit Job**.

The proposal's manager or alternate manager can select **Edit allocations** to change the division hours.

Nobody can change a **Cancelled** proposal.

## Set the status

Procurement does this.

- **Award:** on a **Submitted** proposal, select **Award proposal**.
- **Other statuses:** select **Edit Job**, then change **Status** and save.
- **No Bid or Cancelled:** Turbo asks for a comment and a confirmation. Turbo then saves the status.

Submitted, Awarded and Not Awarded need a **Proposal Value** or **Time and Materials**.

## Charge time to a proposal

You can charge time to an **In Progress** or **Submitted** proposal.

After the award, you need a project that links to the proposal. You can then charge time to the proposal until 30 days after the project's award date. After that, charge time to the project.

## Start a project

Anyone who enters time can [create a project request from a proposal](/docs/project-requests#start-from-a-proposal). One proposal can lead to more than one project.

## Proposal statuses

| Status      | What it means                                 |
| ----------- | --------------------------------------------- |
| In Progress | Procurement prepares it.                      |
| Submitted   | The company submitted it to the client.       |
| Awarded     | The client awarded the work.                  |
| Not Awarded | The client didn't award the work.             |
| No Bid      | The company decided not to submit a proposal. |
| Cancelled   | It stopped. This status is final.             |
