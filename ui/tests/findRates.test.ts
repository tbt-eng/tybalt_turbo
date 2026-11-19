import test from "node:test";
import assert from "node:assert/strict";
import {
  findRatesBody,
  findRatesErrors,
  findRatesKey,
  findResultMessage,
  findStatusText,
  foundRates,
  rateEntered,
  withCarriedRates,
} from "../src/lib/findRates.ts";
import { type ProposedRate, rateChanged, rateErrorKey } from "../src/lib/proposedRates.ts";

const row = (
  role: string,
  roleName: string,
  rate: number | null,
  overtime_rate: number | null = null,
): ProposedRate => ({ role, roleName, rate, overtime_rate });

const sourceRow = (role: string, roleName: string, rate: number, overtime: number) => ({
  role,
  roleName,
  rate,
  overtime_rate: overtime,
  source_rate: rate,
  source_overtime_rate: overtime,
});

test("only the rates entered are searched, with an overtime rate only when entered", () => {
  const rows = [
    row("admin", "Administration Staff", 111),
    row("drone", "Survey Drone Pilot", null),
    row("principal", "Principal", 4321, 5617.3),
  ];
  assert.deepEqual(findRatesBody(rows), {
    entries: [
      { role: "admin", rate: 111 },
      { role: "principal", rate: 4321, overtime_rate: 5617.3 },
    ],
  });
  assert.equal(findRatesKey(rows), JSON.stringify(findRatesBody(rows)));
  assert.deepEqual(
    rows.map((r) => rateEntered(r)),
    [true, false, true],
  );
});

test("no rate entered gives no search", () => {
  assert.equal(findRatesKey([row("admin", "Administration Staff", null)]), "");
  assert.equal(findRatesKey([]), "");
});

test("an input error stops the search and shows on its rate", () => {
  const cases: [ProposedRate, string, string][] = [
    [row("a", "A", null, 150), "rate", "Enter the rate, or clear the overtime rate."],
    [row("a", "A", -5), "rate", "Rate must be at least 1."],
    [row("a", "A", 0), "rate", "Rate must be at least 1."],
    [row("a", "A", 110.5), "rate", "Rate must be a whole number."],
    [row("a", "A", 110, 0.5), "overtime_rate", "Overtime rate must be at least 1."],
    [
      row("a", "A", 110, 143.001),
      "overtime_rate",
      "Overtime rate must have at most two decimal places.",
    ],
  ];
  for (const [input, field, message] of cases) {
    const errors = findRatesErrors([input]);
    assert.deepEqual(errors, { [rateErrorKey("a", field as "rate")]: { message } });
    assert.equal(findRatesKey([input]), "");
  }
  // Two decimals are fine, even where a float is not exact.
  assert.deepEqual(findRatesErrors([row("a", "A", 110, 1.15)]), {});
});

test("the line above the matches", () => {
  assert.equal(findResultMessage(0), "No active sheet has these rates.");
  assert.equal(findResultMessage(1), "1 sheet matches");
  assert.equal(findResultMessage(3), "3 sheets match");
});

test("carried rates replace the starting sheet's rates and show as changed", () => {
  const source = [
    sourceRow("admin", "Administration Staff", 90, 117),
    sourceRow("principal", "Principal", 250, 325),
    sourceRow("designer", "Senior Designer", 160, 208),
  ];
  const carried = foundRates([
    row("admin", "Administration Staff", 111),
    row("principal", "Principal", 4321, 5617.3),
    row("designer", "Senior Designer", null),
  ]);
  const rows = withCarriedRates(source, carried);
  // A rate with no overtime rate entered keeps the sheet's overtime rate.
  assert.deepEqual(rows[0], { ...source[0], rate: 111, overtime_rate: 117 });
  assert.deepEqual(rows[1], { ...source[1], rate: 4321, overtime_rate: 5617.3 });
  // A role not entered is unchanged.
  assert.equal(rows[2], source[2]);
  assert.deepEqual(
    rows.map((r) => rateChanged(r)),
    [true, true, false],
  );
});

test("a carried role that the starting sheet lacks gets a row, in role order", () => {
  const rows = withCarriedRates(
    [sourceRow("principal", "Principal", 250, 325)],
    [{ role: "drone", roleName: "Survey Drone Pilot", rate: 140, overtime_rate: null }],
  );
  assert.deepEqual(
    rows.map((r) => [r.role, r.rate, r.overtime_rate]),
    [
      ["principal", 250, 325],
      ["drone", 140, null],
    ],
  );
  assert.equal(withCarriedRates(rows, []), rows);
});

test("the status line says where the search stands", () => {
  const none = { empty: false, invalid: false, failed: false, count: null };
  assert.equal(findStatusText({ ...none, empty: true }), "Enter the rates you know.");
  assert.equal(findStatusText({ ...none, invalid: true }), "Correct the rates marked below.");
  assert.equal(
    findStatusText({ ...none, failed: true }),
    "Turbo could not search the rate sheets.",
  );
  assert.equal(findStatusText(none), "Searching…");
  assert.equal(findStatusText({ ...none, count: 0 }), "No active sheet has these rates.");
  assert.equal(findStatusText({ ...none, count: 3 }), "3 sheets match");
});
