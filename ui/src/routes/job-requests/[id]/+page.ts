import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { RequestDetail } from "$lib/jobRequests";
// The request's page shows it read only. The form is on its edit page.
export const load: PageLoad = async ({ params }) => {
  const detail = await pb.send<RequestDetail>(`/api/job_requests/${params.id}`, { method: "GET" });
  const duplicates = await pb.send(`/api/job_requests/${params.id}/duplicates`, { method: "GET" });
  return { detail, duplicates };
};
