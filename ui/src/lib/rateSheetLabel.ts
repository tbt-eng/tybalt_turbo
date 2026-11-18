// How Turbo names a rate sheet. The code is the label that people read, print
// and search for, and every revision of a sheet shares it; the name is an
// optional description. See the Rate sheet codes section of
// descriptions/job_rate_sheet_history.md.

import MiniSearch, { type MatchInfo } from "minisearch";

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
// effective date, the job numbers and the client names. A line break
// separates them, as it separates the entries of job_numbers and
// client_names, so a search matches inside one of them and never spans two.
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
  ].join("\n");
}

type RateSheetMatchFields = { job_numbers?: string | null; client_names?: string | null };

// Whether a search matched an entry of job_numbers or client_names.
export type RateSheetEntryMatch = (field: "job_numbers" | "client_names", entry: string) => boolean;

// The Rate Sheets list's rule. The list keeps a sheet when its
// rateSheetSearchText, lowercase, holds the search text, and that text can
// only fall inside one entry, so an entry matches when it holds the search
// text.
export function rateSheetListMatch(searchTerm: string): RateSheetEntryMatch {
  const term = searchTerm.toLowerCase();
  return (_field, entry) => term !== "" && entry.toLowerCase().includes(term);
}

// A picker's rule. MiniSearch splits the search and each field into words and
// keeps a sheet when every search word starts a word of the sheet. Each
// result's match lists the sheet's words that it matched, by field. An entry
// matches when one of its words is in that list for its field. So "24-321"
// finds a sheet by the words 24 and 321, and on a sheet with 24-100 and
// 25-321 it matches both.
export function rateSheetPickerMatch(match: MatchInfo | undefined): RateSheetEntryMatch {
  const tokenize = MiniSearch.getDefault("tokenize");
  const processTerm = MiniSearch.getDefault("processTerm");
  return (field, entry) =>
    match !== undefined &&
    tokenize(entry, field).some((word: string) => match[processTerm(word, field)]?.includes(field));
}

// Why a search found a sheet: the job numbers and client names that it
// matched, such as "Job 24-321 • Client Mobilia SA". It shows three, then
// "+N more". It is empty when none matched, so for an empty search or one
// that matched only the code, the name or the date.
export function rateSheetMatchLine(
  sheet: RateSheetMatchFields,
  matches: RateSheetEntryMatch,
): string {
  const entries = (text?: string | null) => (text ?? "").split("\n").filter((e) => e !== "");
  const found = [
    ...entries(sheet.job_numbers)
      .filter((number) => matches("job_numbers", number))
      .map((number) => `Job ${number}`),
    ...entries(sheet.client_names)
      .filter((name) => matches("client_names", name))
      .map((name) => `Client ${name}`),
  ];
  const shown = found.slice(0, 3).join(" • ");
  return found.length > 3 ? `${shown} +${found.length - 3} more` : shown;
}
