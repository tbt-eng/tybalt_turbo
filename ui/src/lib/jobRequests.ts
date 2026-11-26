// Shared request rules. The server checks these rules again on each action.
import type { WorkflowRecord } from "./clientEditors";

export interface RequestChecklistQuestion {
  key: string;
  text: string;
}
export interface RequestDetail {
  request: WorkflowRecord;
  allocations: WorkflowRecord[];
  attachment?: WorkflowRecord;
  related: Record<string, any>;
  review_token: string;
  // The current questions for the request's stage: approval, or job creation
  // once approved. Each must be answered true to continue.
  checklist: { version: number; questions: RequestChecklistQuestion[] };
  permissions: Record<string, boolean>;
  enabled?: boolean;
  // Why the selected approver has no Approve button, such as a limit below the
  // project value. Only that approver gets it; it is empty for everyone else.
  approval_blocked_reason?: string;
}
// Who acts next on a request, as one sentence for the request page. The
// viewer's role comes from their id. Only procurement's turn depends on a
// permission, because procurement is a claim, not a person on the request.
// A missing name leaves the name out, never a role in its place.
export function requestNextStep(next: {
  status: string;
  creatorName: string;
  approverName: string;
  viewerIsCreator: boolean;
  viewerIsApprover: boolean;
  canCreate: boolean;
  // The number of a linked proposal that is not Awarded yet. The job cannot
  // be created until procurement awards it.
  unawardedProposal?: string;
}): string {
  const creator = next.creatorName.trim(),
    approver = next.approverName.trim();
  switch (next.status) {
    case "Draft":
      if (next.viewerIsCreator) return "You have not submitted this request.";
      return creator ? `${creator} has not submitted this request.` : "Not submitted yet.";
    case "Recalled":
      if (next.viewerIsCreator) return "Recalled. You have not submitted it again.";
      return creator
        ? `Recalled. ${creator} has not submitted it again.`
        : "Recalled. Not submitted again yet.";
    case "Awaiting approval":
      if (next.viewerIsApprover) return "Waiting for your approval.";
      return approver ? `Waiting for approval by ${approver}.` : "Waiting for approval.";
    case "Awaiting creation":
      if (next.unawardedProposal)
        return next.canCreate
          ? `Approved. Award proposal ${next.unawardedProposal}, then create the job.`
          : `Approved. Waiting for procurement to award proposal ${next.unawardedProposal} and create the job.`;
      return next.canCreate
        ? "Approved. Waiting for you or another procurement user to create the job."
        : "Approved. Waiting for procurement to create the job.";
    case "Rejected":
      if (next.viewerIsCreator) return "Rejected. Correct the request, then submit it again.";
      return creator
        ? `Rejected. Waiting for ${creator} to correct it.`
        : "Rejected. Waiting for a correction.";
  }
  return "";
}
// A request can link a proposal only in these statuses; the server refuses
// the others at save.
export function proposalAcceptsRequest(status?: string): boolean {
  return status === "In Progress" || status === "Submitted" || status === "Awarded";
}
export interface RequestAllocation {
  division: string;
  hours: number;
}
// A request started from a proposal or a parent job copies that job's
// allocations (from /api/jobs/{id}/details). A save refuses an inactive
// division, so those rows are left out. `divisions` is the request options'
// list, the same one the division picker offers.
export function copiedAllocations(
  source: WorkflowRecord[] | undefined,
  divisions: WorkflowRecord[] | undefined,
): RequestAllocation[] {
  const active = new Set(
    (divisions || []).filter((record) => record.active).map((record) => record.id),
  );
  return (source || [])
    .map((row) => ({
      division: row.division?.id || "",
      hours: Number(row.hours) || 0,
    }))
    .filter((row) => active.has(row.division));
}
// Rows with no division chosen yet don't count, so a proposal chosen later
// may still fill the list.
export function hasAllocations(rows: WorkflowRecord[]): boolean {
  return rows.some((row) => !!row.division);
}
export function emptyChecklist(questions: RequestChecklistQuestion[]): Record<string, boolean> {
  return Object.fromEntries(questions.map((question) => [question.key, false]));
}
// A contact's email address, or a note that the client has none.
export function contactEmail(contact?: WorkflowRecord): string {
  return contact?.email || (contact?.no_email ? "No email address" : "");
}
export function fullName(record?: WorkflowRecord): string {
  return record?.name || [record?.given_name, record?.surname].filter(Boolean).join(" ");
}
export function address(record?: WorkflowRecord): string {
  return [
    record?.address,
    record?.city,
    record?.province_state,
    record?.postal_code,
    record?.country,
  ]
    .filter(Boolean)
    .join(", ");
}
export function billingAddress(billing?: WorkflowRecord, client?: WorkflowRecord): string {
  return address(billing) || address(client);
}
export function wholeCAD(value: unknown): string {
  return typeof value === "number" && Number.isInteger(value) && value > 0
    ? new Intl.NumberFormat("en-CA", {
        style: "currency",
        currency: "CAD",
        maximumFractionDigits: 0,
      }).format(value)
    : "Not provided";
}
export function printMissing(detail: RequestDetail): string[] {
  const r = detail.request,
    d = detail.related;
  const missing: string[] = [];
  if (!d.client?.name?.trim()) missing.push("Client name");
  if (!r.description?.trim()) missing.push("Description");
  if (
    !r.location?.trim() ||
    !/^[23456789CFGHJMPQRVWX]{8}\+[23456789CFGHJMPQRVWX]{2,3}$/i.test(r.location.trim())
  )
    missing.push("Valid Plus Code location");
  for (const [label, contact] of [
    ["Project contact", d.contact],
    ["Billing contact", d.billing_contact],
  ] as const) {
    if (!fullName(contact).trim()) missing.push(`${label} name`);
    if (!contact?.no_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(contact?.email || ""))
      missing.push(`${label} email`);
  }
  if (!Number.isInteger(r.project_value) || r.project_value <= 0)
    missing.push("Positive whole-dollar project value");
  if (r.time_and_materials_set !== true || typeof r.time_and_materials !== "boolean")
    missing.push("Fee basis");
  // Proposed rates stay inactive until the job is created. They are ready
  // when the server finds a rate for every role, as at submission. A chosen
  // sheet must be active.
  const sheetReady =
    d.rate_sheet_proposed === true
      ? Array.isArray(d.rate_sheet_missing_roles) && d.rate_sheet_missing_roles.length === 0
      : d.rate_sheet?.active === true;
  if (!d.rate_sheet?.id || !sheetReady) missing.push("Rate sheet");
  if (
    r.time_and_materials === true &&
    (!Array.isArray(d.rate_sheet_entries) ||
      !d.rate_sheet_entries.length ||
      d.rate_sheet_entries.some(
        (entry: WorkflowRecord) =>
          !Number.isFinite(entry.rate) || !Number.isFinite(entry.overtime_rate),
      ))
  )
    missing.push("Client billing rates");
  return missing;
}
// The server adds this count only for an active job holder while requests are enabled.
export function canViewJobCreateQueue(
  claims: string[],
  enabled: boolean,
  counts: Record<string, number>,
): boolean {
  return enabled && claims.includes("job") && Object.hasOwn(counts, "/job-requests?queue=creation");
}
