// The proposed rates of a project request, as the request editor and the
// rates table hold them. See the Proposed rates section of
// descriptions/job_rate_sheet_history.md.

// One role's rates on a proposed rate sheet, with the rates of the sheet that
// the proposal started from when they are known. A rate is null until the
// requester enters it, for a role that neither the source sheet nor the saved
// proposal has a rate for.
export interface ProposedRate {
  role: string;
  roleName: string;
  rate: number | null;
  overtime_rate: number | null;
  source_rate?: number;
  source_overtime_rate?: number;
}

export interface ProposedRates {
  name: string;
  reason: string;
  based_on: string;
  rows: ProposedRate[];
}

export function emptyProposedRates(): ProposedRates {
  return { name: "", reason: "", based_on: "", rows: [] };
}

// An empty number input binds null.
const empty = (value: number | string | null | undefined) =>
  value === null || value === undefined || value === "";

// A role has no rate while its rate or its overtime rate is empty.
export function rateMissing(row: ProposedRate): boolean {
  return empty(row.rate) || empty(row.overtime_rate);
}

// The rows with each overtime rate set to the rate times the multiplier,
// rounded to the cent. A row with no rate keeps its overtime rate, so an empty
// row stays empty. A multiplier that is not a positive number changes nothing.
// The proposed rates editor and the rate sheet copy page both use it.
export function withOvertimeMultiplier<
  T extends { rate: number | null; overtime_rate: number | null },
>(rows: T[], multiplier: number | null): T[] {
  const factor = Number(multiplier);
  if (empty(multiplier) || !Number.isFinite(factor) || factor <= 0) return rows;
  return rows.map((row) => {
    const rate = Number(row.rate);
    if (empty(row.rate) || !Number.isFinite(rate)) return row;
    return { ...row, overtime_rate: Math.round(rate * factor * 100) / 100 };
  });
}

// A rate is changed when it differs from the source sheet's, or when the
// source sheet has no rate for the role.
export function rateChanged(row: ProposedRate): boolean {
  return (
    row.source_rate === undefined ||
    Number(row.rate) !== row.source_rate ||
    Number(row.overtime_rate) !== row.source_overtime_rate
  );
}

// The proposed_rates object of a request save. A role with no rate is left
// out, so an incomplete proposal can still be saved; submission names the
// roles that have none. The request editor refuses a save while a role has
// one of its two rates (incompleteRateErrors), so a typed rate is never left
// out.
export function proposedRatesPayload(proposal: ProposedRates) {
  return {
    name: proposal.name,
    reason: proposal.reason,
    based_on: proposal.based_on,
    entries: proposal.rows
      .filter((row) => !rateMissing(row))
      .map((row) => ({
        role: row.role,
        rate: Number(row.rate),
        overtime_rate: Number(row.overtime_rate),
      })),
  };
}

// A rate of a role in the rates table: its rate, its overtime rate, or the
// role itself.
export type RateField = "rate" | "overtime_rate" | "role";

// The key of a rate's error among the rates editor's errors. The rates table
// shows it on the role's row.
export function rateErrorKey(role: string, field: RateField): string {
  return `proposed_rates.rates.${role}.${field}`;
}

// An error for each role that has one of its two rates and not the other. A
// save leaves out a role with no rate, which would lose the rate typed, so the
// request editor refuses the save instead. A role with both rates empty has
// no error.
export function incompleteRateErrors(rows: ProposedRate[]): Record<string, { message: string }> {
  const errors: Record<string, { message: string }> = {};
  for (const row of rows) {
    const noRate = empty(row.rate);
    const noOvertimeRate = empty(row.overtime_rate);
    if (noRate && !noOvertimeRate) {
      errors[rateErrorKey(row.role, "rate")] = { message: "Enter the rate." };
    } else if (!noRate && noOvertimeRate) {
      errors[rateErrorKey(row.role, "overtime_rate")] = { message: "Enter the overtime rate." };
    }
  }
  return errors;
}

// The rows with an empty row added for each rate role that has none, in role
// order, or null when every role has a row. A proposal needs a rate for every
// role, and the source sheet, or a proposal saved before a role was added, can
// lack one.
export function withEveryRole(
  rows: ProposedRate[],
  roles: { id: string; name: string }[],
): ProposedRate[] | null {
  const known = new Set(rows.map((row) => row.role));
  const missing = roles.filter((role) => !known.has(role.id));
  if (missing.length === 0) return null;
  return [
    ...rows,
    ...missing.map((role) => ({
      role: role.id,
      roleName: role.name,
      rate: null,
      overtime_rate: null,
    })),
  ].sort((a, b) => a.roleName.localeCompare(b.roleName));
}

// The fields of the rates editor that show their own error after a refused
// save, under the server's keys.
const proposedRatesFields = ["name", "reason", "based_on"].map(
  (field) => `proposed_rates.${field}`,
);

type FieldErrors = Record<string, { message?: string } | undefined>;

// The field errors of a refused request save that belong to the rates
// editor's fields.
export function proposedRatesFieldErrors(error: unknown): Record<string, { message: string }> {
  const failed = error as { data?: { data?: FieldErrors }; response?: { data?: FieldErrors } };
  const fields: FieldErrors = (failed?.data ?? failed?.response)?.data ?? {};
  return Object.fromEntries(
    proposedRatesFields
      .filter((key) => fields[key]?.message)
      .map((key) => [key, { message: fields[key]?.message ?? "" }]),
  );
}
