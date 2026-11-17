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
  // The code of the saved proposed rates. Turbo gives it on the first save,
  // and a save never sends it.
  code?: string;
  // An optional description.
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

// The label of a role that has no rate, which says what is missing: "No
// overtime rate" when only the overtime rate is empty, otherwise "No rate".
// It is "" for a role that has both rates.
export function missingRateLabel(row: ProposedRate): string {
  if (!rateMissing(row)) return "";
  return empty(row.overtime_rate) && !empty(row.rate) ? "No overtime rate" : "No rate";
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

// The rows with the source sheet's rates added to each row that has none, or
// null when no row changes. Only the source rates change, never the rates. A
// saved proposal has source rates only for the roles it has a rate for, so
// the empty row of another role shows the source sheet's rates this way.
export function withSourceRates(
  rows: ProposedRate[],
  source: { role: string; rate: number; overtime_rate: number }[],
): ProposedRate[] | null {
  const byRole = new Map(source.map((entry) => [entry.role, entry]));
  let changed = false;
  const result = rows.map((row) => {
    const entry = byRole.get(row.role);
    if (row.source_rate !== undefined || !entry) return row;
    changed = true;
    return { ...row, source_rate: entry.rate, source_overtime_rate: entry.overtime_rate };
  });
  return changed ? result : null;
}

// The fields of the rates editor that show their own error after a refused
// save, under the server's keys.
const proposedRatesFields = ["name", "reason", "based_on"].map(
  (field) => `proposed_rates.${field}`,
);

type FieldErrors = Record<string, { message?: string } | undefined>;

type FailedSave = {
  message?: string;
  data?: { message?: string; data?: FieldErrors };
  response?: { message?: string; data?: FieldErrors };
};

// The server keys an error on one rate of the entries that a save sent, such
// as proposed_rates.entries.2.rate, where 2 is the place of the entry in the
// save, not in the table.
const entryErrorKey = /^proposed_rates\.entries\.(\d+)\.(rate|overtime_rate|role)$/;

// The server's message on one rate, in the words of the table, such as
// "Overtime rate must be at least 1." for "entry overtime_rate must be at
// least 1".
function rateMessage(message: string): string {
  const text = message
    .replace(/^entry /i, "")
    .replaceAll("overtime_rate", "overtime rate")
    .trim();
  const sentence = text.charAt(0).toUpperCase() + text.slice(1);
  return /[.!?]$/.test(sentence) ? sentence : `${sentence}.`;
}

// Each error of a refused save on one rate of a sent entry, with its server
// key and its key in the rates table.
function rowErrors(fields: FieldErrors, entries: { role: string }[]) {
  return Object.entries(fields).flatMap(([key, value]) => {
    const match = entryErrorKey.exec(key);
    const entry = match ? entries[Number(match[1])] : undefined;
    if (!match || !entry || !value?.message) return [];
    const field = match[2] as RateField;
    return [{ key, rowKey: rateErrorKey(entry.role, field), message: rateMessage(value.message) }];
  });
}

// The field errors of a refused request save that belong to the rates
// editor's fields. entries are the entries that the save sent, in the order
// sent: an error on one of their rates goes on its role's row, under
// rateErrorKey.
export function proposedRatesFieldErrors(
  error: unknown,
  entries: { role: string }[] = [],
): Record<string, { message: string }> {
  const failed = error as FailedSave | undefined;
  const fields: FieldErrors = (failed?.data ?? failed?.response)?.data ?? {};
  return Object.fromEntries([
    ...proposedRatesFields
      .filter((key) => fields[key]?.message)
      .map((key) => [key, { message: fields[key]?.message ?? "" }]),
    ...rowErrors(fields, entries).map((row) => [row.rowKey, { message: row.message }]),
  ]);
}

// The refused save without the errors that the rates editor shows, for the
// general error message: those of its fields (name, reason, based_on) and
// those that the rates table shows on a row. The editor shows each where it
// belongs, and the keys mean nothing to the user. The other errors stay.
export function withoutShownErrors(error: unknown, entries: { role: string }[] = []): unknown {
  const failed = error as FailedSave | undefined;
  const body = failed?.data ?? failed?.response;
  if (!body?.data) return error;
  const placed = new Set([
    ...proposedRatesFields.filter((key) => body.data?.[key]?.message),
    ...rowErrors(body.data, entries).map((row) => row.key),
  ]);
  if (placed.size === 0) return error;
  const data = Object.fromEntries(Object.entries(body.data).filter(([key]) => !placed.has(key)));
  return { message: failed?.message, data: { ...body, data } };
}
