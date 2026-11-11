import assert from "node:assert/strict";
import test from "node:test";
import { canShareTimesheet } from "../src/lib/timesheets/actions.ts";

const approver = "approver1";

test("the approver can share a submitted time sheet in any state", () => {
  assert.equal(canShareTimesheet({ approver, submitted: true }, approver), true);
  assert.equal(
    canShareTimesheet({ approver, submitted: true, approved: "2026-11-01" }, approver),
    true,
  );
  assert.equal(
    canShareTimesheet(
      { approver, submitted: true, approved: "2026-11-01", committed: "2026-11-02" },
      approver,
    ),
    true,
  );
  assert.equal(
    canShareTimesheet({ approver, submitted: true, rejected: "2026-11-01" }, approver),
    true,
  );
});

test("a time sheet without a submitted field counts as submitted", () => {
  assert.equal(canShareTimesheet({ approver }, approver), true);
});

test("nobody but the approver can share", () => {
  assert.equal(canShareTimesheet({ uid: "owner1", approver, submitted: true }, "owner1"), false);
  assert.equal(canShareTimesheet({ approver, submitted: true }, "someone"), false);
});

test("a person who approves their own time sheet can share it", () => {
  assert.equal(canShareTimesheet({ uid: approver, approver, submitted: true }, approver), true);
});

test("an unsubmitted time sheet can't be shared", () => {
  assert.equal(canShareTimesheet({ approver, submitted: false }, approver), false);
});

test("a signed-out viewer can't share a time sheet that has no approver", () => {
  assert.equal(canShareTimesheet({ approver: "", submitted: true }, ""), false);
  assert.equal(canShareTimesheet({ submitted: true }, ""), false);
});
