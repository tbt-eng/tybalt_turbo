import test from "node:test";
import assert from "node:assert/strict";
import {
  canActivateRateSheets,
  canCorrectRateSheets,
  canCreateRateSheets,
  canReviseRateSheets,
} from "../src/lib/rateSheetPermissions.ts";

test("a user with only the time claim can change no rate sheet", () => {
  const claims = ["time"];
  assert.equal(canCreateRateSheets(claims), false);
  assert.equal(canActivateRateSheets(claims), false);
  assert.equal(canReviseRateSheets(claims), false);
  assert.equal(canCorrectRateSheets(claims), false);
});

test("procurement creates, activates and deactivates sheets, but does not revise", () => {
  const claims = ["time", "job"];
  assert.equal(canCreateRateSheets(claims), true);
  assert.equal(canActivateRateSheets(claims), true);
  assert.equal(canReviseRateSheets(claims), false);
  assert.equal(canCorrectRateSheets(claims), false);
});

test("rate_sheet_revise revises and renames only", () => {
  const claims = ["rate_sheet_revise"];
  assert.equal(canReviseRateSheets(claims), true);
  assert.equal(canCreateRateSheets(claims), false);
  assert.equal(canActivateRateSheets(claims), false);
  assert.equal(canCorrectRateSheets(claims), false);
});

test("admin revises and corrects rates, but activates only with job", () => {
  const claims = ["admin"];
  assert.equal(canReviseRateSheets(claims), true);
  assert.equal(canCorrectRateSheets(claims), true);
  assert.equal(canCreateRateSheets(claims), false);
  assert.equal(canActivateRateSheets(claims), false);
});
