# Time Sheets UI Actions (Current)

This note documents where review actions currently appear in the UI.

## Surfaces

- Main list: `/time/sheets/list` (`ui/src/routes/time/sheets/list/+page.svelte`)
- Pending queue: `/time/sheets/pending` (`ui/src/routes/time/sheets/pending/+page.svelte`)
- Approved By Me: `/time/sheets/approved` (`ui/src/routes/time/sheets/approved/+page.svelte`)
- Shared with Me: `/time/sheets/shared` (`ui/src/routes/time/sheets/shared/+page.svelte`)
- Details page: `/time/sheets/{id}/details` (`ui/src/routes/time/sheets/[id]/details/+page.svelte`)

## Action Placement

- Main list (`/time/sheets/list`) currently includes inline actions:
  - `Recall` (unbundle)
  - `Approve` (when not already approved)
  - `Reject`
  - `Share` (only on a row where the signed-in person is the approver, i.e. a person who approves their own time sheet)
- Pending queue (`/time/sheets/pending`) is detail-first:
  - list rows expose `Share` and `Details`
  - approve/reject is performed from the details page
- Approved By Me (`/time/sheets/approved`): list rows expose `Share`
- Details page (`/time/sheets/{id}/details`) includes review actions based on state:
  - pending: `Recall`, `Approve`, `Reject`
  - approved/not committed: `Commit` (requires `commit` claim or `showAllUi`), `Reject`
  - rejected: `Recall`
  - any state: `Share` for the approver

## Important Clarification

Unlike the current purchase-order and expense list policy, timesheets are not fully details-first yet: `/time/sheets/list` still exposes inline `Approve`/`Reject` controls.

## Sharing

Work on a person's time sheet may have been done under another manager. The person's own manager approves the time sheet and shares it so that the other manager can see the hours worked for them. The other manager finds it under Shared with Me and can read it, but can't approve or reject it.

Only the approver of the time sheet can share it. The create and delete rules of `time_sheet_reviewers` allow only the approver to add or remove a viewer, and they have no state condition, so a committed or rejected time sheet can still be shared. The UI shows the share control when `canShareTimesheet` (`ui/src/lib/timesheets/actions.ts`) is true: the signed-in person is the approver and the time sheet is submitted, because a viewer can open only a submitted time sheet. Adding a viewer queues a `timesheet_shared` notification to that viewer (`app/hooks/hooks.go`). `ShareModal.svelte` adds and removes viewers.

## Rejected time sheets in lists

Pending My Approval, Approved By Me and Shared with Me show a rejected time sheet with `TimesheetRejection.svelte`: a red Rejected label, the rejector's name (`rejector_name` from `app/routes/tallies.sql`), the date and the reason. A time sheet rejected after approval keeps its `approved` value, so it stays on Approved By Me, with the rejection in place of the green Approved label, until its owner recalls it. Recall deletes the time sheet (`app/routes/unbundle_timesheet.go`).
