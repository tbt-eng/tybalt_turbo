import type { PageLoad } from "./$types";
import { pb } from "$lib/pocketbase";
import { error } from "@sveltejs/kit";
import { expensesLeftToReview } from "$lib/expenseApprovalQueue";

export const load: PageLoad = async ({ params }) => {
  let expense;
  try {
    expense = await pb.send(`/api/expenses/details/${params.eid}`, {
      method: "GET",
    });
  } catch (err) {
    console.error(`loading expense details: ${err}`);
    throw error(404, `Failed to load expense details`);
  }

  // How many expenses the viewer has left to review, including this one, when
  // this one is waiting on them. The header shows it; it's not worth failing
  // the page over.
  let leftToReview = 0;
  const pendingMyApproval =
    expense.approver === pb.authStore.record?.id &&
    expense.submitted &&
    expense.approved === "" &&
    expense.rejected === "" &&
    expense.committed === "";
  if (pendingMyApproval) {
    try {
      leftToReview = (await expensesLeftToReview()).length;
    } catch (err) {
      console.error(`counting expenses left to review: ${err}`);
    }
  }

  return { expense, leftToReview };
};
