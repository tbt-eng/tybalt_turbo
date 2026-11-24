import assert from "node:assert/strict";
import test from "node:test";
import { expiringCutoff, poEndState, poIsStale, staleCutoff } from "../src/lib/poAge.ts";

const now = new Date(2026, 10, 23, 12, 0, 0); // 23 November 2026, local time
const po = {
  status: "Active",
  type: "Cumulative",
  approved: "2026-11-01 10:00:00.000Z",
  second_approval: "",
  end_date: "",
};

test("an Active PO is stale 30 days after its approval", () => {
  assert.equal(poIsStale(po, now), false);
  assert.equal(poIsStale({ ...po, approved: "2026-10-01 10:00:00.000Z" }, now), true);
  assert.equal(
    poIsStale({ ...po, approved: "2026-10-01 10:00:00.000Z", status: "Closed" }, now),
    false,
  );
});

test("a late second approval does not keep a PO from being stale", () => {
  const twice = {
    ...po,
    approved: "2026-10-01 10:00:00.000Z",
    second_approval: "2026-11-20 10:00:00.000Z",
  };
  assert.equal(poIsStale(twice, now), true);
});

test("the stale cutoff is a PocketBase datetime", () => {
  assert.match(staleCutoff(now), /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}Z$/);
});

test("only an Active Recurring PO with an end date can expire", () => {
  assert.equal(poEndState({ ...po, end_date: "2026-11-30" }, now), null);
  const recurring = { ...po, type: "Recurring" };
  assert.equal(poEndState(recurring, now), null);
  assert.equal(poEndState({ ...recurring, end_date: "2026-11-30", status: "Closed" }, now), null);
});

test("a Recurring PO is expiring within 30 days and expired after its end date", () => {
  const recurring = { ...po, type: "Recurring" };
  assert.equal(expiringCutoff(30, now), "2026-12-23");
  assert.equal(poEndState({ ...recurring, end_date: "2026-12-24" }, now), null);
  assert.equal(poEndState({ ...recurring, end_date: "2026-12-23" }, now), "expiring");
  assert.equal(poEndState({ ...recurring, end_date: "2026-11-23" }, now), "expiring");
  assert.equal(poEndState({ ...recurring, end_date: "2026-11-22" }, now), "expired");
});
