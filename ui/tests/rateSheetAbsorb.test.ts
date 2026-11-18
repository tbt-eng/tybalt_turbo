import test from "node:test";
import assert from "node:assert/strict";
import { absorbedRateSheets } from "../src/lib/rateSheetAbsorb.ts";

test("absorbed sheets are named by code, revision and name", () => {
  const stored = [
    {
      id: "sheet2",
      code: "2026-11-00007",
      name: "",
      revision: 1,
      active: true,
      _owned_rows: { rate_sheet_entries: [{ id: "entry1", rate: 100 }] },
    },
    { id: "sheet1", code: "2026-11-00007", name: "", revision: 0, active: false },
    { id: "sheet3", code: "2026-11-00004", name: "Client Rates", revision: 0 },
  ];
  assert.deepEqual(absorbedRateSheets(stored), [
    { id: "sheet3", code: "2026-11-00004", name: "Client Rates", revision: 0 },
    { id: "sheet1", code: "2026-11-00007", name: "", revision: 0 },
    { id: "sheet2", code: "2026-11-00007", name: "", revision: 1 },
  ]);
});

test("records that are not stored sheets are left out", () => {
  assert.deepEqual(absorbedRateSheets(null), []);
  assert.deepEqual(absorbedRateSheets({ id: "sheet1" }), []);
  assert.deepEqual(absorbedRateSheets([null, "sheet1", { code: "2026-11-00004" }]), []);
});
