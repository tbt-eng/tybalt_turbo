<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/payroll_report_week_endings.md","sha256":"f84f8e266ebe1a0d5684a1141a2f6a2c6d2020e5d3e7b78bb4e6e2f0bddddbb7","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/payroll_report_week_endings.md

## Payroll report week endings cannot be created, changed or deleted

Lines 3-5

Asked as: edit payroll report; cannot change payroll week endings; read-only payroll view; delete payroll week ending

The payroll_report_week_endings collection is a read-only view. The API has no create, update or delete rule for it, so nobody can change it through the API.

## Who can list payroll report week endings

Lines 7-11

Asked as: who can see payroll report; list payroll periods access; report claim needed; payroll report permission

Only a person with the report claim (the 'report' claim) can list payroll_report_week_endings records. Everyone else is refused.

## Who can view one payroll report week ending

Lines 13-17

Asked as: view single payroll period; open one payroll week; payroll period access denied; see one week ending

Only a person with the report claim (the 'report' claim) can view a single payroll_report_week_endings record. Everyone else is refused.

## Payroll periods are two-week periods counted from 2025-03-01

Lines 22-39

Asked as: which dates are payroll periods; fortnight payroll dates; why week not listed; two week pay cycle start

A week ending counts as a payroll period ending only when it falls a whole multiple of 14 days after 2025-03-01 (JULIANDAY difference % 14 = 0). Turbo lists a period only if a committed time sheet has that week ending, or a committed expense has that pay period ending (pay_period_ending, not empty). Other dates are not listed as periods themselves.

## Period record id is 'p' plus the date digits

Lines 40-45

Asked as: payroll period id format; period id p2025; sort order payroll weeks; newest payroll period first

Each period's id is 'p' followed by the week ending with the dashes removed. For example, 2025-03-01 gives p20250301. The list is sorted by week ending, newest first.

## Committed time sheet count includes the earlier week

Lines 46-57

Asked as: timesheet count wrong week; committed timesheets missing from payroll; timesheet counted wrong period; late week timesheets payroll

A committed time sheet counts toward the period when its week ending is a period ending. A committed time sheet whose week ending is not a period ending counts toward the period ending 7 days later. Uncommitted time sheets are not counted. A period with none shows 0.

## Committed expense count uses pay period ending

Lines 58-67

Asked as: expense count payroll period; expenses missing in payroll; pay period ending expenses total; why expenses not counted

Committed expenses are counted by their pay period ending (pay_period_ending). Only expenses with a non-empty pay period ending that is a period ending (multiple of 14 days from 2025-03-01) are counted. A period with none shows 0.

## Placeholder payroll ID means nine digits starting with 9

Lines 68-74

Asked as: placeholder payroll number; nine digit payroll id; temporary payroll id meaning; payroll id starts with 9

Staff count as having a placeholder payroll ID (admin profile payroll_id) when it is exactly 9 characters: a 9 followed by 8 digits. A person with no admin profile or an empty payroll ID is not counted as a placeholder.

## Placeholder time counts include committed time amendments

Lines 68-85

Asked as: amendments in placeholder time; amended time payroll placeholder; committed amendment placeholder hours; placeholder hours missing amendment

Placeholder time rows come from committed time sheets, and from committed time amendments that have a committed week ending (committed_week_ending), where the person has a placeholder payroll ID. The amendment's committed week ending is used as its week.

## Placeholder time split into week 1 and week 2 counts

Lines 86-109

Asked as: week 1 vs week 2 placeholder; placeholder time split weeks; which week placeholder hours; previous week placeholder count

Placeholder time rows are assigned to a payroll period. A row whose week ending is the period ending counts as week 2 (placeholder_payroll_id_week2_time_count). A row from the week before, 7 days earlier, counts as week 1 (placeholder_payroll_id_week1_time_count). A period with none shows 0.

## Placeholder expense count

Lines 110-121

Asked as: placeholder expense total; expenses with temporary payroll id; placeholder id expense count; placeholder claims payroll count

The placeholder expense count (placeholder_payroll_id_expense_count) counts committed expenses on a period ending, with a non-empty pay period ending, whose owner has a placeholder payroll ID. A period with none shows 0.

## Zero fills missing counts in the report

Lines 122-139

Asked as: blank payroll counts; payroll row shows zero; missing payroll numbers; empty counts in payroll report

Each period row shows the committed time sheet count, committed expense count and the three placeholder counts. Any count that has no matching records shows 0.
