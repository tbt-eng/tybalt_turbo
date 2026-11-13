# Duplicate rate sheets

Two rate sheets are duplicates when they have the same roles, with the same rate and overtime rate for each role. Absorb them into one sheet.

## Find duplicate rate sheets

1. Select **Rate Sheets** under **Business**.
2. Select **Duplicates**.

Each group shows sheets with identical rates. **Target** marks the sheet with the earliest effective date. Turbo tells you why it can't absorb a sheet.

## Absorb a duplicate rate sheet

Only admins, and staff with both absorb access and revision access, can do this.

1. Select **Absorb identical sheets** beside the target.
2. In **Select Record**, select each sheet to absorb.
3. Select **Absorb**.
4. On the **Absorb Actions** page, select **Commit**.
5. Select **Confirm Commit**.

To keep the sheets, select **Undo** instead of **Commit**.

## What the absorb changes

Jobs, their [rate sheet periods](/docs/rate-sheet-periods) and project requests move to the target. Their rates don't change. Until you commit or undo the absorb, the periods of those jobs can't change.

Turbo deletes the absorbed sheets. The target shows their names after **Formerly**.

## Sheets that Turbo can't absorb

- A sheet with an effective date before the target's date.
- A sheet with revisions, if a job, a rate sheet period or a project request uses one of them.
- An active sheet, if the target is inactive. Activate the target first.
