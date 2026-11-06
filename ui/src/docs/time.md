# Time

A time entry records your hours. A time sheet holds your entries for one week. Your manager approves the time sheet.

```flow
You | Enter your time | Add an entry for your hours.
You | Submit the week | Turbo creates a time sheet.
Your manager | Approve or reject | Check the time sheet.
Committer | Commit it | Turbo locks the time sheet.
Note: Rejected? Select **Recall**, edit the entries and submit the week again.
```

## Enter time

1. Select the plus button beside **Entries**.
2. Select the **Time Type**. The other fields change with the type.
3. Enter the details and hours.
4. Select **Save**.

## Submit a week

1. Open **Entries**.
2. Select **Submit** on the week.

Turbo moves the entries to a time sheet and sends it to your manager. You can find it under **Sheets**.

## Approve a time sheet

1. Under **Time Management**, open **Pending My Approval**.
2. Select **Details** on the time sheet.
3. Select **Approve**, or select **Reject** and enter a reason.

## Change a time sheet

Select **Recall** on a pending or rejected time sheet. The entries return to **Entries**. Edit them and submit the week again.

You can't recall an approved time sheet. Your manager or a committer can reject it. You can't change a committed time sheet. Staff with access can correct it under **Amendments**.

## Time sheet statuses

| Status    | What happens next                                 |
| --------- | ------------------------------------------------- |
| Pending   | Your manager approves or rejects it.              |
| Approved  | A committer commits or rejects it.                |
| Rejected  | You recall it, edit the entries and submit again. |
| Committed | It is final.                                      |

## See your time off

**Time Off** shows your available time off.
