import test from "node:test";
import assert from "node:assert/strict";
import {
  rateSheetCodeRevision,
  rateSheetLabel,
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
