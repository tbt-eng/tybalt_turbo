import { pb } from "$lib/pocketbase";
import type { InvoicesResponse } from "$lib/pocketbase-types";

// Client invoices and the costs billed back on them. Drafts and lines are
// written through the invoices and invoice_lines collections; the routes here
// issue invoices and return joined rows. See descriptions/invoices.md.

export type InvoiceWithClient = InvoicesResponse<{ client?: { id: string; name: string } }>;

// An expense or PO billed, or billable, on an invoice.
export type InvoiceCost = {
  kind: "expense" | "po";
  cost_id: string;
  job_id: string;
  job_number: string;
  po_number: string;
  description: string;
  vendor_name: string;
  date: string;
  cost_total: number;
  currency_code: string;
};

export type InvoiceLine = InvoiceCost & {
  id: string;
  amount: number;
  created: string;
  creator_name: string;
  // False when the viewer may not open the expense or PO; it is listed without a link.
  can_open: boolean;
};

export type BillableJob = {
  id: string;
  number: string;
  description: string;
  billable_count: number;
};

// The invoice an expense or PO was billed on, for its details page.
export type CostInvoice = {
  line: {
    id: string;
    amount: number;
    invoice_id: string;
    invoice_number: string;
    invoice_date: string;
  } | null;
  can_bill: boolean;
  can_open_invoice: boolean;
  client_id: string;
  cost_total: number;
  currency_code: string;
};

export function canViewInvoices(claims: string[]) {
  return claims.includes("invoicing") || claims.includes("report");
}

export function issueInvoice(id: string) {
  return pb.send<InvoicesResponse>(`/api/invoices/${id}/issue`, { method: "POST" });
}

export function loadInvoiceLines(id: string) {
  return pb.send<InvoiceLine[]>(`/api/invoices/${id}/lines`, { method: "GET", requestKey: null });
}

export function loadBillableJobs(invoiceId: string) {
  return pb.send<BillableJob[]>(`/api/invoices/${invoiceId}/billable`, {
    method: "GET",
    requestKey: null,
  });
}

export function loadBillableCosts(invoiceId: string, jobId: string) {
  return pb.send<InvoiceCost[]>(`/api/invoices/${invoiceId}/billable`, {
    method: "GET",
    query: { job: jobId },
    requestKey: null,
  });
}

export function loadCostInvoice(kind: "expense" | "po", id: string) {
  const path = kind === "expense" ? "expenses" : "purchase_orders";
  return pb.send<CostInvoice>(`/api/${path}/client_invoice/${id}`, {
    method: "GET",
    requestKey: null,
  });
}

// Adds a line billing one expense or PO on an issued invoice.
export function billCost(
  invoiceId: string,
  kind: "expense" | "po",
  costId: string,
  amount: number,
) {
  return pb.collection("invoice_lines").create({
    invoice: invoiceId,
    kind,
    [kind === "expense" ? "expense" : "purchase_order"]: costId,
    amount,
  });
}

// How a cost is named in a list: its PO number, then what it was for.
export function costLabel(cost: InvoiceCost) {
  const name = cost.kind === "po" ? `PO ${cost.po_number}` : cost.description || "Expense";
  if (cost.kind === "expense" && cost.po_number) return `${name} (PO ${cost.po_number})`;
  return name;
}
