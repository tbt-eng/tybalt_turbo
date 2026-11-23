import { error } from "@sveltejs/kit";
import { globalStore } from "$lib/stores/global";
import { get } from "svelte/store";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { WorkflowRecord } from "$lib/clientEditors";
export const load: PageLoad = async ({ url }) => {
  const parentId = url.searchParams.get("parent") || "";
  const proposalId = url.searchParams.get("proposal") || "";
  const sourceId = parentId || proposalId;
  // One round trip: nothing here waits for another call.
  const [options, source] = await Promise.all([
    pb.send("/api/job_requests/options", { method: "GET" }),
    sourceId
      ? pb.send<WorkflowRecord>(`/api/jobs/${sourceId}/details`, { method: "GET" })
      : ({} as WorkflowRecord),
    globalStore.ensureClaims(),
  ]);
  if (!get(globalStore).claims.includes("time"))
    error(403, "The time permission is required to create a Project Request.");
  if (!options.enabled) error(403, "Request processing is disabled.");
  const request = {
    description: parentId ? "" : source.description || "",
    client: source.client?.id || "",
    contact: source.contact?.id || "",
    manager: source.manager?.id || "",
    alternate_manager: source.alternate_manager?.id || "",
    branch: source.branch_id || "",
    location: source.location || "",
    job_owner: source.job_owner?.id || "",
    parent: parentId,
    rate_sheet: source.rate_sheet?.id || "",
    invoicing_information: source.invoicing_information || "",
    time_and_materials: source.time_and_materials,
    proposal: proposalId || source.proposal_id || "",
    project_status: "Active",
    project_value: proposalId ? source.proposal_value : undefined,
    fn_agreement_answer: "",
    time_and_materials_set: typeof source.time_and_materials === "boolean",
  };
  return { options, request };
};
