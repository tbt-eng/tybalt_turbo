import { redirect } from "@sveltejs/kit";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { RequestDetail } from "$lib/jobRequests";
// The form of a request that its creator can still change. Anyone else goes
// to the request's page, as does the creator of a submitted request. can_edit
// is false while request processing is off.
export const load: PageLoad = async ({ params }) => {
  const detail = await pb.send<RequestDetail>(`/api/job_requests/${params.id}`, { method: "GET" });
  if (!detail.permissions.can_edit) redirect(307, `/job-requests/${params.id}`);
  const options = await pb.send("/api/job_requests/options", { method: "GET" });
  return { detail, options };
};
