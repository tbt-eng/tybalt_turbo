import { goto } from "$app/navigation";
import { resolve } from "$app/paths";
import { pb } from "$lib/pocketbase";
import { expenseTitle } from "$lib/expenseTitle";
import type { Flash } from "$lib/flash";

type QueueExpense = Parameters<typeof expenseTitle>[0] & { id: string };

// IDs of the expenses the viewer can still approve, newest first. The pending
// list keeps rejected expenses until their owner recalls them, so skip those;
// they can't be approved. 200 is the endpoint's page cap.
export async function expensesLeftToReview(): Promise<string[]> {
  const pending: { data: { id: string; rejected: string }[] } = await pb.send(
    "/api/expenses/pending?limit=200",
    { method: "GET" },
  );
  return pending.data.filter((e) => e.rejected === "").map((e) => e.id);
}

// Opens the next expense pending the viewer's approval, or the Approved By Me
// list when none are left, with a flash naming the one just decided. Pushes a
// history entry so Back walks through the expenses just decided.
export async function advanceExpenseApprovalQueue(
  expense: QueueExpense,
  verb: "approved" | "rejected",
) {
  const nextId = (await expensesLeftToReview()).find((id) => id !== expense.id);

  const done = `${verb === "approved" ? "Approved" : "Rejected"} ${expenseTitle(expense)}`;
  const flash: Flash = {
    tone: verb === "approved" ? "success" : "info",
    message: nextId ? done : `${done} · nothing left to review`,
    link: { href: resolve(`/expenses/${expense.id}/details`), label: "View" },
  };
  const url = nextId ? resolve(`/expenses/${nextId}/details`) : resolve("/expenses/approved");
  await goto(url, { state: { flash } });
}
