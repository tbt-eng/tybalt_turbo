// When an Active purchase order is stale or near its end. Payables staff get
// company-wide lists of these (the stale and expiring scopes of
// /api/purchase_orders/visible); every other list shows a label on the
// purchase order. The rules here match those scopes.
// This file imports nothing so that the unit tests can load it.

export const STALE_AFTER_DAYS = 30;
export const EXPIRING_WITHIN_DAYS = 30;

const DAY_MS = 24 * 60 * 60 * 1000;

type AgedPO = {
  status: string;
  type: string;
  approved: string;
  second_approval: string;
  end_date: string;
};

// An approval before this moment is stale. PocketBase stores a datetime as
// "YYYY-MM-DD HH:MM:SS.sssZ" in UTC.
export function staleCutoff(now = new Date()): string {
  return new Date(now.getTime() - STALE_AFTER_DAYS * DAY_MS).toISOString().replace("T", " ");
}

// A calendar date on the user's own device, as YYYY-MM-DD.
function localDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

// A Recurring purchase order that ends on or before this date is expiring.
export function expiringCutoff(days = EXPIRING_WITHIN_DAYS, now = new Date()): string {
  const target = new Date(now);
  target.setHours(0, 0, 0, 0);
  target.setDate(target.getDate() + days);
  return localDate(target);
}

export function poIsStale(po: AgedPO, now = new Date()): boolean {
  const cutoff = staleCutoff(now);
  return (
    po.status === "Active" &&
    ((po.second_approval !== "" && po.second_approval < cutoff) || po.approved < cutoff)
  );
}

// "expired" when the end date is in the past, "expiring" when it is today or
// within the window, and null for every other purchase order.
export function poEndState(po: AgedPO, now = new Date()): "expiring" | "expired" | null {
  if (po.status !== "Active" || po.type !== "Recurring" || po.end_date === "") return null;
  if (po.end_date < localDate(now)) return "expired";
  return po.end_date <= expiringCutoff(EXPIRING_WITHIN_DAYS, now) ? "expiring" : null;
}
