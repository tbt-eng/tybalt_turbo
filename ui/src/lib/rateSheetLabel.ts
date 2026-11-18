// How Turbo names a rate sheet. The code is the label that people read, print
// and search for, and every revision of a sheet shares it; the name is an
// optional description. See the Rate sheet codes section of
// descriptions/job_rate_sheet_history.md.

export interface RateSheetLabelFields {
  code?: string | null;
  name?: string | null;
  revision?: number | null;
}

// The code and revision, such as "2026-11-00004 (rev. 2)". A sheet that has no
// code yet, such as one that is not saved, shows its name in its place.
export function rateSheetCodeRevision(sheet: RateSheetLabelFields): string {
  const code = (sheet.code ?? "").trim() || (sheet.name ?? "").trim();
  return `${code} (rev. ${sheet.revision ?? 0})`;
}

// The sheet as plain text, the same as the server's label
// (hooks.RateSheetLabel): the code and revision, then the name in quotes when
// there is one, such as 2026-11-00004 (rev. 2) "Client Rates". Use it where
// only text fits, such as a dropdown option or a button.
export function rateSheetLabel(sheet: RateSheetLabelFields): string {
  const name = (sheet.name ?? "").trim();
  const label = rateSheetCodeRevision(sheet);
  if (!name || !(sheet.code ?? "").trim()) return label;
  return `${label} "${name}"`;
}

// The fields of rate_sheets_augmented that a search finds a sheet by. A job
// number or a client name finds every revision of the code that the job used,
// now or in its rate sheet history. The rateSheets store indexes these fields
// for the pickers.
export const rateSheetSearchFields = [
  "code",
  "name",
  "effective_date",
  "job_numbers",
  "client_names",
];

// The hint in a rate sheet search box: what the search finds a sheet by.
export const rateSheetSearchPlaceholder = "Code, name, date, job or client";

// The text that a list search matches: the code, the name, the revision, the
// effective date, the job numbers and the client names.
export function rateSheetSearchText(
  sheet: RateSheetLabelFields & {
    effective_date?: string | null;
    job_numbers?: string | null;
    client_names?: string | null;
  },
): string {
  return [
    rateSheetLabel(sheet),
    sheet.effective_date ?? "",
    sheet.job_numbers ?? "",
    sheet.client_names ?? "",
  ].join(" ");
}
