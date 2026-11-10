<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/expense_tracking.go","sha256":"9a66cedfaece15a246f920d40b973bf5ce331111ad4075b48ee86ee647911a68","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/expense_tracking.go

## Who can view expense tracking

Lines 19-54

Asked as: who can see expense tracking; no access to expense tracking; expense tracking permission denied; cannot open expense tracking page

Staff with the report claim, the commit claim (committers) or the admin claim can view expense tracking. Everyone else gets the message "you are not authorized to view expense tracking".

## Only committers can view the expense commit queue

Lines 56-75

Asked as: expense commit queue not allowed; cannot open commit queue; commit queue access denied; who can commit expenses

Only a committer (the commit claim) can view the expense commit queue. Report holders and admins without the commit claim get the message "you are not authorized to view the expense commit queue".

## Committed expense counts by committed week ending

Lines 84-94

Asked as: committed expenses by week; expense totals per week; weekly expense counts; week ending expense summary

Expense tracking counts committed expenses for each committed week ending (committed_week_ending), newest week first. It skips expenses that have no committed value or no committed week ending.

## Count of expenses with placeholder payroll ID

Lines 88-90

Asked as: placeholder payroll ID; fake payroll number; generated payroll id count; staff missing real payroll number

Each week's count also gives how many of its expenses belong to staff whose payroll ID in the admin profile is a generated placeholder (placeholder_payroll_id_expense_count). The test is GeneratedPayrollPlaceholderSQLCondition, defined in another file.

## Expense tracking counts fail with a query error

Lines 96-101

Asked as: expense tracking error; expense counts not loading; failed to execute query expenses; expense report query failed

When the query for committed expense counts fails, Turbo returns the error "failed to execute query". Otherwise it returns one row per committed week ending, each with the committed count and the placeholder payroll ID count.

## Weekly committed expense list is for report, commit or admin holders

Lines 152-155

Asked as: weekly committed list access; who can view weekly expenses; weekly expense list not authorized; committed list permission

The weekly list of committed expenses runs the same viewer check as expense tracking. A person without the report, commit or admin claim gets "you are not authorized to view expense tracking" and no list.

## Weekly expense list needs a valid date

Lines 157-163

Asked as: week ending date missing; wrong date format expenses; committedWeekEnding required; expense week date error

The list of committed expenses for a week requires a committed week ending in the form YYYY-MM-DD. If it is missing, Turbo says "committedWeekEnding is required". If the format is wrong, Turbo says "invalid date format (YYYY-MM-DD)".

## Weekly expense list contents and order

Lines 165-229

Asked as: expenses for the whole week; weekly list everyone's expenses; expense list sort order; all staff expenses for week

The weekly list shows every committed expense in the organization with that committed week ending, not only the viewer's own. It sorts by staff surname, given name, then expense date.

## Allowance types shown as text

Lines 187-192

Asked as: allowance text order; meals listed on expense; breakfast lunch dinner text; which allowances show

Turbo builds the allowance text (allowance_str) from the expense's allowance types. It lists Breakfast, Lunch, Dinner and Lodging in that fixed order, separated by spaces.

## Currency defaults to CAD when missing

Lines 205-211

Asked as: currency missing on expense; expense shows CAD by default; zero settled total; no exchange rate expense

When an expense has no matching currency, Turbo reports currency code CAD, symbol CAD and rate 1. The settled total shows 0 when it is empty.

## Weekly committed expense list fails with a query error

Lines 231-238

Asked as: weekly expense list error; committed expenses won't load; week expenses query failed; weekly list not loading

When the query for the weekly list of committed expenses fails, Turbo returns the error "failed to execute query". Otherwise it returns the matching committed expenses for the requested week ending.

## Commit queue needs the commit claim, not report or admin

Lines 245-248

Asked as: report user cannot commit; admin cannot open commit queue; need commit permission; commit queue not authorized

Before it builds the expense commit queue, Turbo checks for the commit claim (committer). A report or admin holder without the commit claim gets "you are not authorized to view the expense commit queue".

## Commit queue shows blank names when a person is missing

Lines 254-265

Asked as: blank name in commit queue; missing approver name; no name shown on expense; empty committer field

In the expense commit queue, a missing submitter name (given_name, surname) shows as blank. A missing approver, committer, rejector or settler name also shows as blank. Missing approver and committer IDs show as an empty string.

## Commit queue shows who settled the expense

Lines 265-265

Asked as: who settled expense; settler name missing; settled by which person; who paid the expense queue

Each commit queue row names the person who settled the expense (settler_name). The weekly committed list does not include that name.

## Commit queue rows carry a phase label

Lines 266-271

Asked as: expense status label; approved or submitted status; what phase is expense; expense stage in queue

Each commit queue row gets a phase (phase). It is 'Committed' when the expense has a committed value. Otherwise it is 'Approved' when approved is set, 'Submitted' when it is submitted but not approved, and 'Unsubmitted' in every other case. The queue only selects approved, uncommitted expenses, so its rows show 'Approved'.

## Commit queue shows allowances as text

Lines 274-280

Asked as: allowances in commit queue; meal allowance text queue; attachment blank in queue; allowance column empty

Each commit queue row lists the expense's allowance types as one text (allowance_str). It joins Breakfast, Lunch, Dinner and Lodging in that order, separated by spaces. A type the expense does not have is left out. The attachment shows blank when the expense has no attachment document.

## Commit queue defaults for job, client and purchase order

Lines 281-283

Asked as: blank job number queue; no client shown on expense; missing PO number; job description empty

In the commit queue, the job number, job description, client name and PO number show as blank when the expense has none. The client comes from the expense's job.

## Commit queue currency and settled total defaults

Lines 284-292

Asked as: no currency symbol queue; rate shows 1; settled total zero queue; PO blank on expense

In the commit queue, a missing currency shows as code 'CAD' and symbol 'CAD'. Its icon and rate date are blank, and its rate is 1. A missing settled total (settled_total) shows as 0. The PO number shows blank when the expense has no purchase order.

## Commit queue lists approved, uncommitted expenses

Lines 304-307

Asked as: why expense not in queue; expenses missing from commit queue; which expenses appear to commit; approved expenses list

The expense commit queue lists only expenses that are submitted and approved, and not committed or rejected.

## Foreign-currency card and on-account expenses wait for settlement

Lines 308-312

Asked as: foreign currency expense missing; on account expense not showing; credit card expense hidden; waiting for settlement expense

The commit queue leaves out expenses in a currency other than CAD, with payment type OnAccount or CorporateCreditCard, that have no settlement yet. They appear after they are settled. An expense with no currency counts as CAD, so it is not held back.

## Commit queue sort order

Lines 313-320

Asked as: commit queue order; queue sorted by name; expense list sorting; who shows first in queue

The commit queue sorts by staff surname, then given name. The sort also has a phase ordering, but the filter leaves only approved expenses.
