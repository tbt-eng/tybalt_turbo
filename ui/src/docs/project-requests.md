# Project requests

Every new project and subjob starts as a request. An approver checks it, then procurement creates the job.

```flow
You | Create the request | Enter the details, division hours and approver.
You | Upload and submit | Upload one signed PDF that authorizes the project.
Approver | Approve or reject | Check for duplicates and confirm the checklist.
Procurement | Create the job | Check the request again, then create the job.
Note: Rejected? Edit the request and submit it again.
```

## Create a request

1. Select the plus button beside **Project Requests**. For a subjob, open the parent project and select the **Create Subjob Request** button.
2. Enter the details and division hours.
3. Select an invoicing profile and an approver. The approver's limit must cover the full project value.
4. Select **Save request**.

## Upload and submit

1. Get a signed document that legally binds the company and the client.
2. Upload it as one PDF.
3. Select **Submit**.

You can use **Print Project Authorization** or another document. The approver and procurement must accept it.

Include the rate appendix for time and materials work. Save your edits before you print or upload.

## Approve a request

1. Open **Request Approvals**. Open the request.
2. Review the possible duplicate jobs.
3. Confirm all seven checklist statements.
4. Select **Approve request**.

To decline, enter a rejection reason and select **Reject**.

## Create the job

1. Open **Job Create Queue**. Open the request.
2. Check the request, the PDF, existing jobs and other pending requests.
3. Confirm that the job doesn't duplicate existing work.
4. Select **Create Job**. Turbo assigns the job number.

## Change a request

- **You:** select **Recall for editing**. This works until procurement creates the job.
- **Approver:** enter a reason and select **Reject**.

Both actions clear the approval, so the request needs a new approval. You must edit a rejected request before you submit it again.

## Request statuses

| Status            | What happens next                      |
| ----------------- | -------------------------------------- |
| Draft             | You edit and submit it.                |
| Awaiting approval | The approver approves or rejects.      |
| Awaiting creation | Procurement creates the job.           |
| Rejected          | You edit it and submit it again.       |
| Recalled          | You can edit it, then submit it again. |

## Change a project's value

1. Open the job and select **Propose new value**. Only procurement can do this.
2. Enter the new value and select an approver.
3. Upload the new PDF.

The approver then opens **Value Approvals** and approves or rejects. The project value stays the same until the approver approves the change.
