import test from "node:test";
import assert from "node:assert/strict";
import MiniSearch from "minisearch";
import {
  rateSheetCodeRevision,
  rateSheetLabel,
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
    job_numbers: "24-321,P25-007",
    client_names: "Mobilia SA,BN Rail Corporation",
  },
  {
    id: "rev1",
    code: "2026-11-00004",
    name: "Client Rates",
    revision: 1,
    effective_date: "2026-12-01",
    job_numbers: "24-321,P25-007",
    client_names: "Mobilia SA,BN Rail Corporation",
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
