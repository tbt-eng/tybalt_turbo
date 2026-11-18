// How Turbo names a rate sheet. The code is the label that people read, print
// and search for, and every revision of a sheet shares it; the name is an
// optional description. See the Rate sheet codes section of
// descriptions/job_rate_sheet_history.md.

import MiniSearch, { type SearchResult } from "minisearch";

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

// The text that a copy of the sheet's code takes, such as "2026-11-00004":
// the code alone, without the revision or the name, as a rate sheet search
// takes it. Empty when the sheet has no code yet.
export function rateSheetCopyText(sheet: RateSheetLabelFields): string {
  return (sheet.code ?? "").trim();
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

// The entries of job_numbers or client_names, which a line break separates.
function entries(text?: string | null): string[] {
  return (text ?? "").split("\n").filter((e) => e !== "");
}

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
// keeps a sheet when every search word starts a word of the sheet. A result's
// queryTerms are the search words, and its match lists the sheet's words that
// they matched, by field. When an entry has every search word, as the start
// of one of its words, only such entries match: "91-89" names 91-8901, not
// 91-8801. When none has them all, an entry matches when one of its words is
// in match for its field: "24-321" on a sheet with 24-100 and 25-321 names
// both, since together they are why the search found it.
export function rateSheetPickerMatch(
  sheet: RateSheetMatchFields,
  result: SearchResult | undefined,
): RateSheetEntryMatch {
  if (result === undefined) return () => false;
  const tokenize = MiniSearch.getDefault("tokenize");
  const processTerm = MiniSearch.getDefault("processTerm");
  const words = (field: string, entry: string): string[] =>
    tokenize(entry, field).map((word: string) => processTerm(word, field));
  const hasEvery: RateSheetEntryMatch = (field, entry) => {
    const own = words(field, entry);
    return result.queryTerms.every((term) => own.some((word) => word.startsWith(term)));
  };
  const hasOne: RateSheetEntryMatch = (field, entry) =>
    words(field, entry).some((word) => result.match[word]?.includes(field));
  const anyHasEvery =
    entries(sheet.job_numbers).some((n) => hasEvery("job_numbers", n)) ||
    entries(sheet.client_names).some((c) => hasEvery("client_names", c));
  return anyHasEvery ? hasEvery : hasOne;
}

// One kind of match in why a search found a sheet: the label, said once and
// plural when more than one entry of the kind matched, such as "Jobs"; the
// first three entries; and how many more matched.
export interface RateSheetMatchGroup {
  label: string;
  shown: string[];
  more: number;
}

const matchKinds = [
  { field: "job_numbers", one: "Job", many: "Jobs" },
  { field: "client_names", one: "Client", many: "Clients" },
] as const;

// Why a search found a sheet: the job numbers that it matched, then the client
// names, as one group for each kind that matched, such as Jobs 91-8901 ·
// 91-8902 · 91-8903 +5 more, then Client Lakehead Surveys, Inc. It is empty
// when none matched, so for an empty search or one that matched only the
// code, the name or the date. RateSheetMatchLine.svelte shows it.
export function rateSheetMatchGroups(
  sheet: RateSheetMatchFields,
  matches: RateSheetEntryMatch,
): RateSheetMatchGroup[] {
  return matchKinds.flatMap(({ field, one, many }) => {
    const found = entries(sheet[field]).filter((entry) => matches(field, entry));
    if (found.length === 0) return [];
    return [
      {
        label: found.length === 1 ? one : many,
        shown: found.slice(0, 3),
        more: Math.max(0, found.length - 3),
      },
    ];
  });
}
