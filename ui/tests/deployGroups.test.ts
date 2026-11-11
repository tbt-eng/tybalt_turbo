import test from "node:test";
import assert from "node:assert/strict";
import { deployGroupAllows } from "../src/lib/deployGroups.ts";

test("a feature with no deploy group is for every user", () => {
  assert.equal(deployGroupAllows(null, []), true);
  assert.equal(deployGroupAllows({ enabled: true }, []), true);
  assert.equal(deployGroupAllows({ enabled: true, deploy_group: "" }, []), true);
  assert.equal(deployGroupAllows({ enabled: true, deploy_group: null }, ["group1"]), true);
});

test("a feature with a deploy group is for the users in it", () => {
  const value = { enabled: true, deploy_group: "group1" };
  assert.equal(deployGroupAllows(value, ["group1"]), true);
  assert.equal(deployGroupAllows(value, ["group2", "group1"]), true);
  assert.equal(deployGroupAllows(value, ["group2"]), false);
  assert.equal(deployGroupAllows(value, []), false);
});

test("a deploy group that is not an id lets nobody in", () => {
  assert.equal(deployGroupAllows({ deploy_group: 7 }, ["7"]), false);
  assert.equal(deployGroupAllows({ deploy_group: ["group1"] }, ["group1"]), false);
});
