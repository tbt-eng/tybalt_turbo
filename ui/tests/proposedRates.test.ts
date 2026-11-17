import test from "node:test";
import assert from "node:assert/strict";
import {
  type ProposedRate,
  type ProposedRates,
  incompleteRateErrors,
  proposedRatesFieldErrors,
  proposedRatesPayload,
  rateChanged,
  rateErrorKey,
  rateMissing,
  withEveryRole,
  withOvertimeMultiplier,
  withSourceRates,
  withoutRowErrors,
} from "../src/lib/proposedRates.ts";

const roles = [
  { id: "principal", name: "Principal" },
  { id: "drone", name: "Survey Drone Pilot" },
  { id: "biologist", name: "Biologist" },
];

const sourceRow = (role: string, roleName: string, rate: number): ProposedRate => ({
  role,
  roleName,
  rate,
  overtime_rate: rate * 1.3,
  source_rate: rate,
  source_overtime_rate: rate * 1.3,
});

const proposal = (rows: ProposedRate[]): ProposedRates => ({
  name: "Client Rates",
  reason: "Contract schedule B",
  based_on: "standard",
  rows,
});

test("a role with no row gets an empty row, in role order", () => {
  const rows = [sourceRow("principal", "Principal", 250), sourceRow("biologist", "Biologist", 120)];
  const merged = withEveryRole(rows, roles);
  assert.deepEqual(
    merged?.map((row) => row.role),
    ["biologist", "principal", "drone"],
  );
  const drone = merged?.find((row) => row.role === "drone");
  assert.deepEqual(drone, {
    role: "drone",
    roleName: "Survey Drone Pilot",
    rate: null,
    overtime_rate: null,
  });
});

test("rows that cover every role are left as they are", () => {
  const rows = roles.map((role) => sourceRow(role.id, role.name, 100));
  assert.equal(withEveryRole(rows, roles), null);
  // The roles are not loaded yet.
  assert.equal(withEveryRole([], []), null);
});

test("a role with an empty rate or overtime rate has no rate", () => {
  const row = sourceRow("principal", "Principal", 250);
  assert.equal(rateMissing(row), false);
  assert.equal(rateMissing({ ...row, rate: null }), true);
  assert.equal(rateMissing({ ...row, overtime_rate: null }), true);
  // Zero is a value; the server refuses it.
  assert.equal(rateMissing({ ...row, rate: 0 }), false);
});

test("an empty row counts as changed and does not throw", () => {
  assert.equal(
    rateChanged({ role: "drone", roleName: "Drone", rate: null, overtime_rate: null }),
    true,
  );
  assert.equal(rateChanged(sourceRow("principal", "Principal", 250)), false);
});

test("the save leaves out roles with no rate and sends numbers only", () => {
  const rows = withEveryRole([sourceRow("principal", "Principal", 250)], roles) ?? [];
  rows[0] = { ...rows[0], rate: 140 }; // Biologist: rate only
  const payload = proposedRatesPayload(proposal(rows));
  assert.deepEqual(payload.entries, [{ role: "principal", rate: 250, overtime_rate: 325 }]);
  for (const entry of payload.entries) {
    assert.ok(Number.isFinite(entry.rate) && Number.isFinite(entry.overtime_rate));
  }
  assert.equal(payload.reason, "Contract schedule B");
});

test("a role with only one of its rates gets an error on the empty rate", () => {
  const rows = withEveryRole([sourceRow("principal", "Principal", 250)], roles) ?? [];
  const biologist = rows.findIndex((row) => row.role === "biologist");
  const drone = rows.findIndex((row) => row.role === "drone");
  rows[biologist] = { ...rows[biologist], rate: 140 };
  rows[drone] = { ...rows[drone], overtime_rate: 130 };
  assert.deepEqual(incompleteRateErrors(rows), {
    [rateErrorKey("biologist", "overtime_rate")]: { message: "Enter the overtime rate." },
    [rateErrorKey("drone", "rate")]: { message: "Enter the rate." },
  });
  // A role with both rates, or with neither, has no error.
  const complete = withEveryRole([sourceRow("principal", "Principal", 250)], roles) ?? [];
  assert.deepEqual(incompleteRateErrors(complete), {});
  // A rate the user cleared is empty too.
  const cleared = {
    ...sourceRow("principal", "Principal", 250),
    overtime_rate: "" as unknown as null,
  };
  assert.deepEqual(Object.keys(incompleteRateErrors([cleared])), [
    rateErrorKey("principal", "overtime_rate"),
  ]);
});

test("the field errors of a refused save are those of the editor's fields", () => {
  const data = {
    "proposed_rates.reason": { code: "value_required", message: "Reason is required." },
    "proposed_rates.entries.0.rate": { code: "x", message: "Entry rate must be at least 1." },
  };
  const expected = { "proposed_rates.reason": { message: "Reason is required." } };
  assert.deepEqual(proposedRatesFieldErrors({ data: { message: "Check", data } }), expected);
  assert.deepEqual(proposedRatesFieldErrors({ response: { data } }), expected);
  assert.deepEqual(proposedRatesFieldErrors(new Error("offline")), {});
  assert.deepEqual(proposedRatesFieldErrors(undefined), {});
});

test("an error on one rate of a refused save goes on its role's row", () => {
  // The save sent Biologist and then Principal; the keys count in that order.
  const entries = [{ role: "biologist" }, { role: "principal" }];
  const data = {
    "proposed_rates.reason": { code: "value_required", message: "reason is required" },
    "proposed_rates.entries.0.rate": {
      code: "validation_not_whole_number",
      message: "entry rate must be a whole number",
    },
    "proposed_rates.entries.1.overtime_rate": {
      code: "validation_negative_number",
      message: "entry overtime_rate must be at least 1",
    },
    "proposed_rates.entries.1.role": { code: "x", message: "entry role must be unique" },
    // No entry was sent in this place, so no row can show it.
    "proposed_rates.entries.7.rate": { code: "x", message: "entry rate must be at least 1" },
  };
  const error = { message: "Check the request fields.", data: { message: "Check", data } };
  assert.deepEqual(proposedRatesFieldErrors(error, entries), {
    "proposed_rates.reason": { message: "reason is required" },
    [rateErrorKey("biologist", "rate")]: { message: "Rate must be a whole number." },
    [rateErrorKey("principal", "overtime_rate")]: { message: "Overtime rate must be at least 1." },
    [rateErrorKey("principal", "role")]: { message: "Role must be unique." },
  });

  // The general message keeps the errors that no row shows.
  const general = withoutRowErrors(error, entries) as {
    message: string;
    data: { message: string; data: object };
  };
  assert.equal(general.message, "Check the request fields.");
  assert.equal(general.data.message, "Check");
  assert.deepEqual(Object.keys(general.data.data), [
    "proposed_rates.reason",
    "proposed_rates.entries.7.rate",
  ]);
  // An error with no row errors is left as it is.
  const other = new Error("offline");
  assert.equal(withoutRowErrors(other, entries), other);
  assert.equal(withoutRowErrors(error, []), error);
});

test("the multiplier sets each overtime rate to the cent and leaves empty rows empty", () => {
  const rows = withEveryRole([sourceRow("principal", "Principal", 251)], roles) ?? [];
  const applied = withOvertimeMultiplier(rows, 1.3);
  assert.equal(applied.find((row) => row.role === "principal")?.overtime_rate, 326.3);
  for (const row of applied.filter((row) => row.role !== "principal")) {
    assert.equal(row.rate, null);
    assert.equal(row.overtime_rate, null);
  }
  // A row with no rate keeps the overtime rate it has.
  const partial = { ...sourceRow("drone", "Drone", 100), rate: null };
  assert.equal(withOvertimeMultiplier([partial], 1.5)[0].overtime_rate, 130);
  // Rounding is to the cent.
  assert.equal(withOvertimeMultiplier([sourceRow("x", "X", 117)], 1.333)[0].overtime_rate, 155.96);
});

test("a multiplier that is not a positive number changes nothing", () => {
  const rows = [sourceRow("principal", "Principal", 250)];
  for (const multiplier of [null, 0, -1, Number.NaN]) {
    assert.equal(withOvertimeMultiplier(rows, multiplier), rows);
  }
});

test("an empty row of a saved proposal gets the source rates, not the rates", () => {
  // A saved proposal with a rate for Principal only; the source has a rate
  // for Principal and Biologist, and none for the drone pilot.
  const saved = { ...sourceRow("principal", "Principal", 250), rate: 260 };
  const rows = withEveryRole([saved], roles) ?? [];
  const source = [
    { role: "principal", rate: 250, overtime_rate: 325 },
    { role: "biologist", rate: 120, overtime_rate: 156 },
  ];
  const filled = withSourceRates(rows, source);
  assert.deepEqual(
    filled?.find((row) => row.role === "biologist"),
    {
      role: "biologist",
      roleName: "Biologist",
      rate: null,
      overtime_rate: null,
      source_rate: 120,
      source_overtime_rate: 156,
    },
  );
  // A row that has source rates keeps them, and its rates.
  assert.equal(
    filled?.find((row) => row.role === "principal"),
    saved,
  );
  // A role that the source lacks keeps no source rates.
  assert.equal(filled?.find((row) => row.role === "drone")?.source_rate, undefined);
  // Once filled, nothing changes.
  assert.equal(withSourceRates(filled ?? [], source), null);
});
