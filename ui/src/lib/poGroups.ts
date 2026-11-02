import { pb } from "$lib/pocketbase";
import type { PoGroupUploadField } from "$lib/poGroupCsv";
import type { PoApproversResponse } from "$lib/pocketbase-types";

// Purchase order groups (issue #168): one template purchase order copied to
// many staff members from a CSV upload. See app/routes/purchase_order_groups.go.

// super_po holders create groups. admin holders can look at them too.
export function canViewPOGroups(claims: string[]): boolean {
  return claims.includes("super_po") || claims.includes("admin");
}

export type PoGroupTemplate = {
  type: "One-Time" | "Cumulative";
  kind: string;
  payment_type: string;
  date: string;
  division: string;
  branch: string;
  description: string;
  total: number;
  currency: string;
  vendor: string;
  approver: string;
  priority_second_approver: string;
};

export type PoGroupRequest = {
  group_id?: string;
  template: PoGroupTemplate;
  upload_fields: PoGroupUploadField[];
  rows: Record<string, string>[];
};

export type PoGroupPreviewRow = {
  index: number;
  email: string;
  uid: string;
  name: string;
  total: number;
  description: string;
  division_code: string;
  branch_name: string;
  errors?: Record<string, string>;
};

export type PoGroupStaff = { uid: string; email: string; name: string };

export type PoGroupPreview = {
  rows: PoGroupPreviewRow[];
  count: number;
  error_count: number;
  total: number;
  currency_code: string;
  missing_active_staff: PoGroupStaff[];
  // The people who could approve every row. A second approver is needed only
  // when some row is over the second approval threshold.
  approvers: PoApproversResponse[];
  second_approvers: PoApproversResponse[];
  second_approval_required: boolean;
};

export type PoGroupSummary = {
  id: string;
  created: string;
  creator: string;
  creator_name: string;
  // The template's description, blank when each purchase order's came from
  // the upload.
  description: string;
  vendor_name: string;
  count: number;
  unapproved_count: number;
  total: number;
  currency_code: string;
};

// A group is named by its template description, or by its vendor when each
// purchase order has its own description.
export function poGroupTitle(group: PoGroupSummary): string {
  if (group.description) return group.description;
  return group.vendor_name ? `${group.vendor_name} purchase orders` : "PO Group";
}

// The parts of a PocketBase ClientResponseError these routes fill in.
export type PoGroupApiError = {
  status?: number;
  isAbort?: boolean;
  message?: string;
  response?: {
    code?: string;
    message?: string;
    data?: Record<string, { message: string }>;
    preview?: PoGroupPreview;
  };
};

export async function previewPoGroup(request: PoGroupRequest): Promise<PoGroupPreview> {
  return (await pb.send("/api/purchase_orders/groups/preview", {
    method: "POST",
    body: request,
    requestKey: "po-group-preview",
  })) as PoGroupPreview;
}

export async function createPoGroup(
  request: PoGroupRequest,
): Promise<{ group_id: string; count: number }> {
  return (await pb.send("/api/purchase_orders/groups", {
    method: "POST",
    body: request,
    requestKey: null,
  })) as { group_id: string; count: number };
}

export async function fetchPoGroups(): Promise<PoGroupSummary[]> {
  return (await pb.send("/api/purchase_orders/groups", { method: "GET" })) as PoGroupSummary[];
}

export async function fetchPoGroup(id: string): Promise<PoGroupSummary> {
  return (await pb.send(`/api/purchase_orders/groups/${id}`, {
    method: "GET",
  })) as PoGroupSummary;
}
