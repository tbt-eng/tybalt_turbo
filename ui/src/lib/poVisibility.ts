import { pb } from "$lib/pocketbase";
import type { PurchaseOrdersAugmentedResponse } from "$lib/pocketbase-types";

export type VisiblePOScope =
  | "all"
  | "mine"
  | "active"
  | "rejected"
  | "stale"
  | "expiring"
  | "approved_by_me_awaiting_second"
  | "group";

// A revision is an Active purchase order sent back for approval. It keeps its
// po_number, which an Unapproved purchase order otherwise doesn't have yet.
export function isPORevision(po: {
  status: string;
  po_number: string;
  legacy_manual_entry?: boolean;
}): boolean {
  return po.status === "Unapproved" && po.po_number !== "" && !po.legacy_manual_entry;
}

// The same test as a collection API filter.
export const poRevisionFilter =
  "(status='Unapproved' && po_number!='' && legacy_manual_entry=false)";

export function poRevisionWaitingOn(
  po: Pick<
    VisiblePurchaseOrderResponse,
    "rejected" | "approved" | "uid_name" | "approver_name" | "priority_second_approver_name"
  >,
): { name: string; reason: string } {
  if (po.rejected !== "") {
    return { name: po.uid_name || "The creator", reason: "To resubmit or cancel it" };
  }
  if (po.approved !== "") {
    return {
      name: po.priority_second_approver_name || "A second approver",
      reason: "Second approval",
    };
  }
  return { name: po.approver_name || "The approver", reason: "First approval" };
}

export type VisiblePurchaseOrderResponse = PurchaseOrdersAugmentedResponse & {
  branch: string;
  branch_code: string;
  branch_name: string;
  canceller_name: string;
  closer_name: string;
  covered_within_project_budget: boolean;
  expenses_total: number;
  kind_label: string;
  kind_name: string;
  po_group: string;
  recurring_expected_occurrences: number;
  recurring_remaining_occurrences: number;
  remaining_amount: number;
  print_max_amount: number;
  // Set only when one purchase order is fetched by id.
  revisable: boolean;
};

export async function fetchVisiblePOs(
  scope: VisiblePOScope,
  beforeDate?: string,
  limit?: number,
): Promise<VisiblePurchaseOrderResponse[]> {
  const params = new URLSearchParams({ scope });
  if (scope === "stale" && beforeDate) {
    params.set("stale_before", beforeDate);
  }
  if (scope === "expiring" && beforeDate) {
    params.set("expiring_before", beforeDate);
  }
  if (limit !== undefined) {
    params.set("limit", String(limit));
  }

  return (await pb.send(`/api/purchase_orders/visible?${params.toString()}`, {
    method: "GET",
  })) as VisiblePurchaseOrderResponse[];
}

// The visible purchase orders in a purchase order group.
export async function fetchGroupPOs(groupId: string): Promise<VisiblePurchaseOrderResponse[]> {
  const params = new URLSearchParams({ scope: "group", group: groupId });
  return (await pb.send(`/api/purchase_orders/visible?${params.toString()}`, {
    method: "GET",
  })) as VisiblePurchaseOrderResponse[];
}

export async function fetchVisiblePO(id: string): Promise<VisiblePurchaseOrderResponse> {
  return (await pb.send(`/api/purchase_orders/visible/${id}`, {
    method: "GET",
  })) as VisiblePurchaseOrderResponse;
}

export async function fetchPendingPO(id: string): Promise<VisiblePurchaseOrderResponse> {
  return (await pb.send(`/api/purchase_orders/pending/${id}`, {
    method: "GET",
  })) as VisiblePurchaseOrderResponse;
}

// Closing and cancelling are final. The owner of a purchase order and payables
// staff (payables_admin) may do either; the server makes the same checks.
type EndablePO = { status: string; type: string; uid: string; committed_expenses_count: number };
function mayEndPO(po: EndablePO, claims: string[], userId: string): boolean {
  return po.status === "Active" && (claims.includes("payables_admin") || po.uid === userId);
}
// A purchase order with expenses can't be cancelled. The lists know only the
// committed count; the server refuses one with an expense that isn't committed.
export function poMayBeCancelledBy(po: EndablePO, claims: string[], userId: string): boolean {
  return mayEndPO(po, claims, userId) && po.committed_expenses_count === 0;
}
// Only a Recurring or Cumulative purchase order with a committed expense can
// be closed by hand. The server checks pending expenses and says what remains.
export function poMayBeClosedBy(po: EndablePO, claims: string[], userId: string): boolean {
  return mayEndPO(po, claims, userId) && po.type !== "One-Time" && po.committed_expenses_count > 0;
}
