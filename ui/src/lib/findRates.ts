// Find by rates in a project request: the requester enters the rates they know
// and Turbo lists the active sheets that have every one of them. It is a tool,
// not a state of the request: it ends with Use on a sheet, which chooses it,
// or with the rates carried into proposed rates. See the Find a sheet by rates
// section of descriptions/job_rate_sheet_history.md.

import { type ProposedRate, rateErrorKey } from "./proposedRates.ts";

// An empty number input binds null.
const empty = (value: number | string | null | undefined) =>
  value === null || value === undefined || value === "";

// A rate as the server compares it: in whole cents, with a small tolerance,
// as hooks.HasAtMostTwoDecimals does (1.15 * 100 is 114.99999999999999).
const hasAtMostTwoDecimals = (value: number) => {
  const cents = value * 100;
  return Math.abs(cents - Math.round(cents)) < 1e-6;
};

// One role's rates that the requester entered. overtime_rate is null when
// they entered none: a sheet then matches on the rate alone.
export interface FoundRate {
  role: string;
  roleName: string;
  rate: number;
  overtime_rate: number | null;
}

// An active sheet that has every rate entered, as POST
// /api/job_requests/find_rates returns it, newest code first.
export interface FoundRateSheet {
  id: string;
  code: string;
  name: string;
  revision: number;
  effective_date: string;
  job_count: number;
}

// Whether the requester entered anything in a row.
export function rateEntered(row: ProposedRate): boolean {
  return !empty(row.rate) || !empty(row.overtime_rate);
}

// An error for each rate entered that a search can't use, under rateErrorKey,
// as the rates table shows them: an overtime rate without the role's rate, a
// rate that is not a whole number of at least 1, or an overtime rate below 1
// or with more than two decimal places. These are the rules of a sheet's
// rates, which the server checks too.
export function findRatesErrors(rows: ProposedRate[]): Record<string, { message: string }> {
  const errors: Record<string, { message: string }> = {};
  for (const row of rows) {
    if (empty(row.rate)) {
      if (!empty(row.overtime_rate)) {
        errors[rateErrorKey(row.role, "rate")] = {
          message: "Enter the rate, or clear the overtime rate.",
        };
      }
      continue;
    }
    const rate = Number(row.rate);
    if (!(rate >= 1)) {
      errors[rateErrorKey(row.role, "rate")] = { message: "Rate must be at least 1." };
    } else if (!Number.isInteger(rate)) {
      errors[rateErrorKey(row.role, "rate")] = { message: "Rate must be a whole number." };
    }
    if (empty(row.overtime_rate)) continue;
    const overtime = Number(row.overtime_rate);
    if (!(overtime >= 1)) {
      errors[rateErrorKey(row.role, "overtime_rate")] = {
        message: "Overtime rate must be at least 1.",
      };
    } else if (!hasAtMostTwoDecimals(overtime)) {
      errors[rateErrorKey(row.role, "overtime_rate")] = {
        message: "Overtime rate must have at most two decimal places.",
      };
    }
  }
  return errors;
}

// The rates entered, in row order. A row with no rate is left out, so call it
// only when findRatesErrors finds none.
export function foundRates(rows: ProposedRate[]): FoundRate[] {
  return rows
    .filter((row) => !empty(row.rate))
    .map((row) => ({
      role: row.role,
      roleName: row.roleName,
      rate: Number(row.rate),
      overtime_rate: empty(row.overtime_rate) ? null : Number(row.overtime_rate),
    }));
}

// The body of the search. An overtime rate is sent only when entered.
export function findRatesBody(rows: ProposedRate[]) {
  return {
    entries: foundRates(rows).map((rate) =>
      rate.overtime_rate === null
        ? { role: rate.role, rate: rate.rate }
        : { role: rate.role, rate: rate.rate, overtime_rate: rate.overtime_rate },
    ),
  };
}

// The key of a search, and its body: "" when there is nothing to search, with
// no rate entered or with an error to fix first. A result is shown only while
// the rates still have the key it was made for.
export function findRatesKey(rows: ProposedRate[]): string {
  if (Object.keys(findRatesErrors(rows)).length > 0) return "";
  const body = findRatesBody(rows);
  return body.entries.length === 0 ? "" : JSON.stringify(body);
}

// The line above the matches.
export function findResultMessage(count: number): string {
  if (count === 0) return "No active sheet has these rates.";
  return `${count} ${count === 1 ? "sheet matches" : "sheets match"}`;
}

// The rows of proposed rates with the carried rates in place: each carried
// rate replaces its role's rate, and its overtime rate when one was entered.
// Otherwise the row keeps its overtime rate, the starting sheet's for a new
// source, which the requester can change or set with the multiplier. A role
// with no row gets one. The other rows are unchanged. The rows stay in role
// order.
export function withCarriedRates(rows: ProposedRate[], carried: FoundRate[]): ProposedRate[] {
  if (carried.length === 0) return rows;
  const byRole = new Map(carried.map((rate) => [rate.role, rate]));
  const result = rows.map((row) => {
    const rate = byRole.get(row.role);
    if (!rate) return row;
    byRole.delete(row.role);
    return {
      ...row,
      rate: rate.rate,
      overtime_rate: rate.overtime_rate ?? row.overtime_rate,
    };
  });
  for (const rate of byRole.values()) {
    result.push({
      role: rate.role,
      roleName: rate.roleName,
      rate: rate.rate,
      overtime_rate: rate.overtime_rate,
    });
  }
  return result.sort((a, b) => a.roleName.localeCompare(b.roleName));
}

// Where a search stands, for the status line under the choice of rates.
export interface FindStatus {
  // The rates have no rate entered.
  empty: boolean;
  // A rate entered has an error to correct first.
  invalid: boolean;
  // The last search for these rates failed.
  failed: boolean;
  // The number of matches for these rates, or null until the search answers.
  count: number | null;
}

// The status line's text. It is always one short line, so the line keeps
// its height as the search goes on.
export function findStatusText(status: FindStatus): string {
  if (status.empty) return "Enter the rates you know.";
  if (status.invalid) return "Correct the rates marked below.";
  if (status.failed) return "Turbo could not search the rate sheets.";
  if (status.count === null) return "Searching…";
  return findResultMessage(status.count);
}
