import { pb } from "$lib/pocketbase";
import { fetchApprovedByMe, type ApprovedByMe } from "$lib/jobValue.svelte";
import type { PageLoad } from "./$types";
// Load every page of the queue so the list can search and page all of it.
async function fetchQueue(queue: string) {
  const fetchPage = (page: number) =>
    pb.send(`/api/job_requests?queue=${encodeURIComponent(queue)}&page=${page}&perPage=100`, {
      method: "GET",
    });
  const first = await fetchPage(1);
  const rest = await Promise.all(
    Array.from({ length: Math.max(0, first.totalPages - 1) }, (_, i) => fetchPage(i + 2)),
  );
  return [first, ...rest].flatMap((result) => result.items ?? []);
}
export const load: PageLoad = async ({ url }) => {
  const queue = url.searchParams.get("queue") || "mine";
  // The Approved tab lists this approver's own approvals on jobs, which come
  // from the job attachments, not from the requests.
  if (queue === "history") {
    return { items: [], approvedItems: [], historyItems: await fetchApprovedByMe(), queue };
  }
  // The Approvals tab also lists the requests this approver approved that
  // still wait for a job, as the PO approvals page does for second approval.
  const [items, approvedItems] = await Promise.all([
    fetchQueue(queue),
    queue === "approval" ? fetchQueue("approved") : Promise.resolve([]),
  ]);
  return { items, approvedItems, historyItems: [] as ApprovedByMe[], queue };
};
