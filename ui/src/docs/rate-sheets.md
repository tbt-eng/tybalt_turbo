# Rate sheets

A rate sheet lists the hourly rate and the overtime rate that a client pays for each role. Every project uses one rate sheet. Proposals don't use one.

## How Turbo uses a rate sheet

Each time entry on a job has a role. Turbo multiplies the hours by the rate for that role on the job's rate sheet. WIP and the job summaries show the result.

If the sheet has no rate for the role, Turbo uses the person's default charge out rate.

## Choose a rate sheet for a project

1. Open the project request or the job.
2. Select a sheet in **Rate Sheet**. Turbo lists only active sheets.
3. Save your changes.

A project must have a rate sheet. Choose with care. A change to a project's rate sheet reprices all past time on the job. Most staff can't make this change.

## Find a rate sheet

1. Select **Rate Sheets** under **Business**.
2. Select the name of a sheet.

The list shows the revision, the effective date, the number of jobs and **Active** or **Inactive**. To share the rates, select **Print** or **Download CSV** on the sheet.

## Create a rate sheet

Procurement creates rate sheets from an existing sheet.

1. Find a sheet that is close to the new one.
2. Select **Use as Template**.
3. Enter the **Name** and the **Effective Date**.
4. Change the **Rate** and the **Overtime Rate** for each role.
5. Select **Save**.

A rate is a whole number. Turbo saves the new sheet as Inactive.

## Set the overtime rates

1. Enter a multiplier above **Overtime Rate**.
2. Select **×**.

Turbo sets each overtime rate to the rate times the multiplier. Check the overtime rates before you save.

## Activate a rate sheet

1. Open the sheet.
2. Select **Active**.

A sheet must have a rate for every role. If roles are missing, Turbo shows **Add Entry**. Select the **Role**, enter the rates and select **Add Entry** for each one. Only procurement can activate a sheet.

## Deactivate a rate sheet

Select **Inactive** on the sheet. Only procurement can do this.

New projects can't select an inactive sheet. Jobs that use the sheet keep it, and Turbo still prices their time with it.

## Revise a rate sheet

A revision is a new version of an active sheet. It has the same name and the next revision number. Only admins and staff with revision access can revise.

1. Open the sheet and select **Revise**.
2. Enter the **Effective Date**. It can't be before the last revision's date.
3. Change the rates and select **Save**.
4. Activate the new revision.

## After a revision

Activation of a revision makes the older revisions inactive. You can't activate an older revision again.

Jobs keep the revision that they use. A job doesn't move to the new revision.

## Correct a rate

This is an exceptional fix, such as for a typing error. Don't use it to change rates. [Revise the sheet](/docs/rate-sheets#revise-a-rate-sheet) instead. Only admins can do this.

1. Open the sheet.
2. Change the **Rate** or the **Overtime Rate** of a role.
3. Select **Save** on that row.

The change reprices all past time on every job that uses the sheet.

## Your role on a time entry

Select your **Role** on each time entry for a job. Ask your manager or the project manager if you don't know it. To fill it in for you, set your **Default Role** on your [profile](/docs/profile).
