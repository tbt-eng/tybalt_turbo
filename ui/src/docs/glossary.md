# Glossary

Use these words in Turbo and in these docs. One word has one meaning.

| Word                  | Meaning                                                                                                                                             | Don't use       |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| Absorb                | Combine duplicate records into one record, the target. Turbo moves what used them to the target, then deletes them.                                 | merge           |
| Admin profile         | The record that Turbo keeps about you: your payroll ID, default branch and time off balances. See [Admin profiles](/docs/admin-profiles).           | staff profile   |
| Approver              | A person who approves a time sheet, expense, purchase order, project request or project value change. See **Who approves what**.                    |                 |
| Authorization PDF     | The signed PDF that authorizes a project or a value change. It must legally bind the company and the client.                                        |                 |
| Category              | A group of costs within one job.                                                                                                                    |                 |
| Child PO              | A One-Time PO for an expense that takes a Cumulative PO over its total. Its number is its parent's number plus a suffix, such as 2611-0379-01.      |                 |
| Committer             | Staff who commit time sheets and expenses. Commit is the final step.                                                                                |                 |
| Contract              | The agreement between the company and a client on the work and its rates. A project's rate sheet must match it.                                     |                 |
| Default branch        | The branch that an admin sets on your admin profile. Turbo uses it when a record has no job or PO.                                                  |                 |
| Expense               | The accounting record of a purchase. It records the actual amount spent.                                                                            |                 |
| Invoicing profile     | Who receives a client's invoices, and any invoice instructions. Many jobs can use one profile.                                                      | billing profile |
| Job                   | A project or a proposal. A proposal's number starts with P.                                                                                         |                 |
| Payables staff        | Staff who can cancel, close or convert a purchase order, and settle foreign-currency expenses.                                                      |                 |
| Procurement           | Staff who create and edit jobs, maintain clients, contacts and invoicing profiles, and propose project value changes. Purchase orders are separate. |                 |
| Profile               | Your own details: your name, your manager and your default division. See [Profile](/docs/profile).                                                  |                 |
| Project Authorization | The agreement that Turbo prints for signatures. You can use it as the authorization PDF, or use another document.                                   | PA              |
| Project Request       | A request to create a project or a subjob. Proposals don't need one.                                                                                | job request     |
| Proposal              | An offer of work to a client. Only procurement can create one.                                                                                      |                 |
| Proposed rates        | New rates in a project request. See [Proposed rates](/docs/proposed-rates).                                                                         |                 |
| Purchase order (PO)   | Permission for the company to make a purchase. A PO isn't an accounting record.                                                                     |                 |
| Rate sheet            | The hourly rates that a client pays for each role. A project uses one at a time. See [Rate sheets](/docs/rate-sheets).                              |                 |
| Rate sheet period     | The dates in which a project uses one rate sheet. See [Rate sheet periods](/docs/rate-sheet-periods).                                               |                 |
| Recall                | Take back a time sheet, expense or project request that you submitted, so you can edit it. Purchase orders have no recall.                          | unbundle        |
| Revise                | Make a new version of a record: new details for an approved purchase order, or new rates for a rate sheet.                                          |                 |
| Role                  | The work that you do on a job. The rate sheet sets the rate of each role.                                                                           |                 |
| Staff                 | People who work for the company.                                                                                                                    | employee        |
| Subjob                | A project under a parent project. It has its parent's client. Its number is its parent's number plus a suffix, such as 25-123-01.                   |                 |
| Time sheet            | Your time entries for one week.                                                                                                                     | timesheet       |
| Under revision        | A revised purchase order that waits for approval.                                                                                                   |                 |

## Who approves what

| Record                          | Approver                                                                                                                                             |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Time sheet or expense           | The manager you chose on your profile.                                                                                                               |
| Purchase order                  | A person you select. Their purchase order limit must cover the amount, unless the PO needs a second approver. Then the second approver's limit must. |
| Project request or value change | A person you select, whose project limit covers the full value.                                                                                      |
