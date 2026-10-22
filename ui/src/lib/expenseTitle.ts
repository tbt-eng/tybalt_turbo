import { allExpensePaymentTypeOptions } from "$lib/expensePaymentTypes";
import { formatCurrencyAmount } from "$lib/utilities";

// A short name for an expense, which has no number and whose description can
// run to paragraphs: who was paid (or what kind of claim it is) and how much.
export function expenseTitle(expense: {
  vendor_name?: string;
  payment_type: string;
  distance?: number;
  total: number;
  currency_code: string;
}): string {
  if (expense.payment_type === "Mileage") {
    return `Mileage · ${expense.distance} km`;
  }
  const payee =
    expense.vendor_name?.trim() ||
    allExpensePaymentTypeOptions.find((o) => o.id === expense.payment_type)?.name ||
    expense.payment_type;
  return `${payee} · ${formatCurrencyAmount(expense.total, expense.currency_code)}`;
}
