import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
// Load every page of the queue so the list can search and page all of it.
export const load: PageLoad = async ({ url }) => {
  const queue = url.searchParams.get("queue") || "mine";
  const fetchPage = (page: number) =>
    pb.send(`/api/job_requests?queue=${encodeURIComponent(queue)}&page=${page}&perPage=100`, {
      method: "GET",
    });
  const first = await fetchPage(1);
  const rest = await Promise.all(
    Array.from({ length: Math.max(0, first.totalPages - 1) }, (_, i) => fetchPage(i + 2)),
  );
  return { items: [first, ...rest].flatMap((result) => result.items ?? []), queue };
};
