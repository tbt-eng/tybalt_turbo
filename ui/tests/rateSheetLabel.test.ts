import test from "node:test";
import assert from "node:assert/strict";
import MiniSearch from "minisearch";
import {
  rateSheetCodeRevision,
  rateSheetLabel,
  rateSheetListMatch,
  rateSheetMatchLine,
  rateSheetPickerMatch,
  rateSheetSearchFields,
  rateSheetSearchText,
} from "../src/lib/rateSheetLabel.ts";

test("a sheet shows its code and revision", () => {
  assert.equal(
    rateSheetCodeRevision({ code: "2026-11-00004", name: "Client Rates", revision: 2 }),
    "2026-11-00004 (rev. 2)",
  );
});

test("the text label matches the server's label", () => {
  // hooks.RateSheetLabel: code, revision, then the name in quotes.
  assert.equal(
    rateSheetLabel({ code: "2026-11-00004", name: " Client Rates ", revision: 2 }),
    '2026-11-00004 (rev. 2) "Client Rates"',
  );
  assert.equal(
    rateSheetLabel({ code: "2026-11-00004", name: "", revision: 0 }),
    "2026-11-00004 (rev. 0)",
  );
  assert.equal(rateSheetLabel({ code: "2026-11-00004", revision: 0 }), "2026-11-00004 (rev. 0)");
});

test("a sheet with no code shows its name in the code's place", () => {
  assert.equal(rateSheetLabel({ code: "", name: "Draft", revision: 0 }), "Draft (rev. 0)");
});

test("the search text holds the code, the name and the effective date", () => {
  const text = rateSheetSearchText({
    code: "2026-11-00004",
    name: "Client Rates",
    revision: 1,
    effective_date: "2026-12-01",
  });
  for (const part of ["2026-11-00004", "Client Rates", "rev. 1", "2026-12-01"]) {
    assert.ok(text.includes(part), part);
  }
});

// Two revisions of one code carry the same jobs and clients, from
// rate_sheets_augmented, so a job on either finds both.
const sheets = [
  {
    id: "rev0",
    code: "2026-11-00004",
    name: "Client Rates",
    revision: 0,
    effective_date: "2026-01-01",
    job_numbers: "24-321\nP25-007",
    client_names: "Mobilia SA\nBN Rail Corporation",
  },
  {
    id: "rev1",
    code: "2026-11-00004",
    name: "Client Rates",
    revision: 1,
    effective_date: "2026-12-01",
    job_numbers: "24-321\nP25-007",
    client_names: "Mobilia SA\nBN Rail Corporation",
  },
  {
    id: "other",
    code: "2026-11-00005",
    name: "",
    revision: 0,
    effective_date: "2026-12-01",
    job_numbers: "",
    client_names: "",
  },
];

test("the list search finds a sheet by a job number or a client name", () => {
  for (const term of ["24-321", "4-32", "p25", "mobilia", "rail corp"]) {
    const found = sheets.filter((s) => rateSheetSearchText(s).toLowerCase().includes(term));
    assert.deepEqual(
      found.map((s) => s.id),
      ["rev0", "rev1"],
      term,
    );
  }
});

test("a picker's search finds a sheet by a job number or a client name", () => {
  // The rateSheets store's index, with collectionStore's AND and
  // DSAutoComplete's prefix search.
  const index = new MiniSearch({
    fields: rateSheetSearchFields,
    searchOptions: { combineWith: "AND" },
  });
  index.addAll(sheets);
  for (const term of [
    "24-321",
    "24-32",
    "321",
    "P25-007",
    "Mobilia",
    "rail",
    "BN Rail",
    "2026-11-00004",
  ]) {
    const found = index.search(term, { prefix: true }).map((r) => r.id);
    assert.deepEqual(found.sort(), ["rev0", "rev1"], term);
  }
  assert.deepEqual(index.search("Sepulchi", { prefix: true }), []);
});

// The match line, with each search's own rule. The list keeps a sheet when its
// search text holds the term; a picker keeps the results of the store's index.
type Sheet = (typeof sheets)[number];

const matchSheets = [
  {
    id: "jobs",
    code: "2026-11-00006",
    name: "Many Jobs",
    revision: 0,
    effective_date: "2026-03-01",
    job_numbers: "24-100\n25-321\n91-8930\n91-8931\n91-8932\n91-8933\n91-8934",
    client_names: "Lakehead Surveys, Inc.\nMobilia SA\nInc Partners",
  },
  ...sheets,
];

function listLines(term: string): Record<string, string> {
  const lines: Record<string, string> = {};
  for (const sheet of matchSheets) {
    if (!rateSheetSearchText(sheet).toLowerCase().includes(term.toLowerCase())) continue;
    lines[sheet.id] = rateSheetMatchLine(sheet, rateSheetListMatch(term));
  }
  return lines;
}

function pickerLines(term: string): Record<string, string> {
  const index = new MiniSearch<Sheet>({
    fields: rateSheetSearchFields,
    searchOptions: { combineWith: "AND" },
  });
  index.addAll(matchSheets);
  const byId = new Map(matchSheets.map((s) => [s.id, s]));
  const lines: Record<string, string> = {};
  for (const result of index.search(term, { prefix: true })) {
    lines[result.id] = rateSheetMatchLine(byId.get(result.id)!, rateSheetPickerMatch(result.match));
  }
  return lines;
}

test("a job match names the job", () => {
  assert.deepEqual(listLines("24-321"), { rev0: "Job 24-321", rev1: "Job 24-321" });
  assert.deepEqual(pickerLines("P25-0"), { rev0: "Job P25-007", rev1: "Job P25-007" });
});

test("a client match names the client", () => {
  assert.deepEqual(listLines("mobilia"), {
    jobs: "Client Mobilia SA",
    rev0: "Client Mobilia SA",
    rev1: "Client Mobilia SA",
  });
  assert.deepEqual(pickerLines("BN Rail"), {
    rev0: "Client BN Rail Corporation",
    rev1: "Client BN Rail Corporation",
  });
});

test("a client name with a comma is one entry", () => {
  assert.deepEqual(listLines("surveys, inc"), { jobs: "Client Lakehead Surveys, Inc." });
  // The picker matches by word, and "inc" is a word of two clients.
  assert.deepEqual(pickerLines("Surveys"), { jobs: "Client Lakehead Surveys, Inc." });
  assert.deepEqual(pickerLines("Surveys, Inc"), {
    jobs: "Client Lakehead Surveys, Inc. • Client Inc Partners",
  });
});

test("a search can match a job and a client at once", () => {
  assert.deepEqual(pickerLines("24-321 mobilia"), {
    jobs: "Job 24-100 • Job 25-321 • Client Mobilia SA",
    rev0: "Job 24-321 • Client Mobilia SA",
    rev1: "Job 24-321 • Client Mobilia SA",
  });
  const line = rateSheetMatchLine(
    { job_numbers: "24-321", client_names: "Client 24 Ltd" },
    rateSheetListMatch("24"),
  );
  assert.equal(line, "Job 24-321 • Client Client 24 Ltd");
});

test("more than three matches show three, then +N more", () => {
  assert.deepEqual(listLines("91-893"), {
    jobs: "Job 91-8930 • Job 91-8931 • Job 91-8932 +2 more",
  });
  assert.deepEqual(pickerLines("91"), { jobs: "Job 91-8930 • Job 91-8931 • Job 91-8932 +2 more" });
});

test("a match on the code, the name or the date shows no line", () => {
  for (const term of ["2026-11-00004", "client rates", "2026-12-01"]) {
    for (const lines of [listLines(term), pickerLines(term)]) {
      assert.ok(Object.keys(lines).length > 0, term);
      for (const line of Object.values(lines)) assert.equal(line, "", term);
    }
  }
});

test("an empty search shows no line", () => {
  for (const sheet of matchSheets) {
    assert.equal(rateSheetMatchLine(sheet, rateSheetListMatch("")), "");
    assert.equal(rateSheetMatchLine(sheet, rateSheetPickerMatch(undefined)), "");
  }
});

test("a picker's split term names each job that it matched", () => {
  // MiniSearch finds "24-321" by the words 24 and 321, here of 24-100 and
  // 25-321, so the line names both jobs. The list needs the whole text, so it
  // does not find the sheet.
  assert.deepEqual(pickerLines("24-321").jobs, "Job 24-100 • Job 25-321");
  assert.equal(listLines("24-321").jobs, undefined);
});

test("a list search does not span two fields", () => {
  // The date and the first job number are on lines of their own, and a search
  // box holds no line break.
  assert.deepEqual(listLines("2026-03-01 24"), {});
  assert.deepEqual(listLines("2026-03-01"), { jobs: "" });
});

test("the line agrees with why each search found a sheet", () => {
  const terms = ["24", "24-3", "321", "91-89", "inc", "mobilia sa", "rail", "2026", "rates", "p25"];
  for (const term of terms) {
    // A picker result has a line when the search matched a word of a job
    // number or a client name.
    const index = new MiniSearch<Sheet>({
      fields: rateSheetSearchFields,
      searchOptions: { combineWith: "AND" },
    });
    index.addAll(matchSheets);
    const lines = pickerLines(term);
    for (const result of index.search(term, { prefix: true })) {
      const byJobOrClient = Object.values(result.match)
        .flat()
        .some((field) => field === "job_numbers" || field === "client_names");
      assert.equal(lines[result.id] !== "", byJobOrClient, `${term} ${result.id}`);
    }
    // A list row that only a job or a client finds has a line.
    for (const [id, line] of Object.entries(listLines(term))) {
      const sheet = matchSheets.find((s) => s.id === id)!;
      const withoutJobs = { ...sheet, job_numbers: "", client_names: "" };
      if (!rateSheetSearchText(withoutJobs).toLowerCase().includes(term)) {
        assert.notEqual(line, "", `${term} ${id}`);
      }
    }
  }
});
