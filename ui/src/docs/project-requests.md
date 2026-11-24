# Project requests

Every new project and subjob starts as a request. An approver checks it, then procurement creates the job. A request can [find a sheet by rates](/docs/proposed-rates#find-a-sheet-by-rates) or [propose new rates](/docs/proposed-rates).

```flow
You | Create the request | Start from a proposal if one exists. Enter the details, division hours and approver.
You | Upload and submit | Upload one signed PDF that authorizes the project.
Approver | Approve or reject | Confirm the checklist.
Procurement | Create or reject | Check for duplicates and confirm the checklist.
Note: Rejected? Edit the request and submit it again.
```

## Start from a proposal

This is the preferred start.

1. [Find the proposal](/docs/proposals#find-a-proposal).
2. Select the yellow briefcase button on its row.
3. Check the copied details and the project value. Add what is missing.
4. Enter the division hours and select an approver.
5. Select **Save request**.

Only an In Progress, Submitted or Awarded proposal has the button. Procurement must award it before you submit.

## Create a request

Do this only when no proposal exists.

1. Select the plus button beside **Jobs**. For a subjob, select **Create Subjob Request** on the parent project.
2. Select **Request a project without a proposal**.
3. Enter the details and division hours.
4. Select an invoicing profile and an approver. The list shows approvers whose limit covers the project value.
5. Select **Save request**.

## Upload and submit

1. Get a signed document that legally binds the company and the client.
2. Upload it as one PDF.
3. Select **Submit**.

You can use **Print Project Authorization** or another document. The approver and procurement must accept it.

Include the rate appendix for time and materials work. Save your edits before you print or upload.

## Approve a request

1. Open **Request Approvals**. Open the request.
2. Confirm all checklist statements.
3. Select **Approve request**.

To decline, enter a rejection reason and select **Reject**.

Requests that you approved stay below your list, under **Approved By Me, Awaiting Job Creation**, until Procurement creates the job.

## Create the job

1. Open **Job Create Queue**. Open the request.
2. Check the request, PDF and **Duplicate-work review**.
3. Confirm the checklist.
4. Select **Create Job**.

The review lists exact matches and the client's open projects with an award or time in the last year.

For [proposed rates](proposed-rates#after-approval), you can change **Rate sheet name** first.

To reject a duplicate, enter a reason. Select **Reject**.

## Change a request

- **You:** select **Recall for editing**. This works until procurement creates the job.
- **Approver:** enter a reason and select **Reject**.
- **Procurement:** after approval, enter a reason and select **Reject**.

These actions clear the approval, so the request needs a new approval. You must edit a rejected request before you submit it again.

## Delete a request

You can delete your own request if its status is Draft, Recalled or Rejected. Recall a submitted request first.

1. Open the request.
2. Select **Delete request**.
3. Select **Delete**.

Turbo removes the request and its PDF. You can't undo this.

## Request statuses

| Status            | What happens next                                   |
| ----------------- | --------------------------------------------------- |
| Draft             | You edit and submit it, or delete it.               |
| Awaiting approval | The approver approves or rejects.                   |
| Awaiting creation | Procurement creates the job or rejects the request. |
| Rejected          | You edit it and submit it again, or delete it.      |
| Recalled          | You edit it and submit it again, or delete it.      |

## Change a project's value

1. Open the job and select **Propose new value**. Only procurement can do this.
2. Enter the new value and select an approver.
3. Upload the new PDF.
4. Select **Save proposed value**.

The approver then opens **Value Approvals** and approves or rejects. The project value stays the same until the approver approves the change.

## Your approvals

**Approved** lists jobs and value changes you approved in the last 30 days.
