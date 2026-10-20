import { pb } from "$lib/pocketbase";

// The vendor invoice behind an expense and every expense that shares it, from
// GET /api/expenses/invoice/{id}. See the vendor invoice section of
// descriptions/expense_documents.md.

export type VendorInvoicePart = {
  id: string;
  po_number: string;
  date: string;
  total: number;
  for_name: string;
  entered_by_name: string;
  submitted: boolean;
  approved: string;
  committed: string;
  rejected: string;
  // False when the viewer may not open this expense; it is listed without a link.
  can_open: boolean;
};

export type VendorInvoiceView = {
  document: string;
  vendor_invoice_number: string;
  vendor_invoice_date: string;
  vendor_invoice_total: number;
  currency_code: string;
  recorded_by_name: string;
  vendor_invoice_recorded: string;
  allocated: number;
  parts: VendorInvoicePart[];
};

export function loadVendorInvoice(expenseId: string) {
  return pb.send<VendorInvoiceView>(`/api/expenses/invoice/${expenseId}`, {
    method: "GET",
    requestKey: null,
  });
}

// Expenses have no status field; this follows the same order as the expense
// pages: rejected, committed, approved, submitted, then draft.
export function partStatus(part: VendorInvoicePart) {
  if (part.rejected) return "Rejected";
  if (part.committed) return "Committed";
  if (part.approved) return "Approved";
  if (part.submitted) return "Submitted";
  return "Draft";
}

// A part waiting for or past approval locks the invoice details.
export function partLocksDetails(part: VendorInvoicePart) {
  return part.submitted && !part.rejected;
}

export function invoiceIsBlank(view: VendorInvoiceView) {
  return (
    view.vendor_invoice_number === "" &&
    view.vendor_invoice_date === "" &&
    Math.round(view.vendor_invoice_total * 100) === 0
  );
}
