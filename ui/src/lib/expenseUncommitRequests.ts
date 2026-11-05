import { shortDate } from "$lib/utilities";
import type {
  ExpenseUncommitBlock,
  ExpenseUncommitOutcome,
  ExpenseUncommitRequestRow,
} from "$lib/svelte-types";

export const expenseUncommitOutcomeOptions: { id: ExpenseUncommitOutcome; name: string }[] = [
  { id: "uncommit_and_reject", name: "Uncommit and reject" },
  { id: "uncommit_only", name: "Uncommit only" },
];

// Keep in step with IsPayrollExpensePaymentType in the backend.
export function isPayrollPaymentType(paymentType: string): boolean {
  return ["Expense", "PersonalReimbursement", "Allowance", "Mileage", "Meals"].includes(
    paymentType,
  );
}

export function expenseUncommitOutcomeLabel(outcome: ExpenseUncommitOutcome): string {
  return outcome === "uncommit_and_reject" ? "Uncommit and reject" : "Uncommit only";
}

export function expenseUncommitReasonText(row: ExpenseUncommitRequestRow): string {
  const reason = row.reason_name || "Direct uncommit";
  return row.note ? `${reason}: ${row.note}` : reason;
}

export function expenseUncommitBlockLabel(block: ExpenseUncommitBlock): string {
  return block.code === "payroll_closed" ? "Payroll closed" : "Too old";
}

export function expenseUncommitBlockText(block: ExpenseUncommitBlock): string {
  if (block.code === "too_old") {
    return `This expense was committed more than ${block.limit_days} days ago.`;
  }
  if (block.pay_period_ending && block.pay_date) {
    return `Payroll for the period ending ${shortDate(block.pay_period_ending, true)} may already have run (pay day ${shortDate(block.pay_date, true)}).`;
  }
  return "This expense is paid on payroll and has no pay period, so it may already have been paid.";
}

// PocketBase datetimes have a space where Date.parse wants a T.
export function daysSince(datetime: string, now = Date.now()): number {
  return Math.floor((now - Date.parse(datetime.replace(" ", "T"))) / 86_400_000);
}
