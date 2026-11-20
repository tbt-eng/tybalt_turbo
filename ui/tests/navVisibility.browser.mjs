import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// The sidebar's links and buttons for a user, from the real nav config.
await runBrowserHarness(async (page, origin) => {
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.goto(`${origin}/tests/fixtures/navVisibility.html`);
  await page.waitForFunction(() => window.ready);
  const labels = (ctx, showAllUi = false) =>
    page.evaluate(([c, s]) => window.visibleLabels(c, s), [ctx, showAllUi]);
  const has = (list, ...wanted) =>
    wanted.forEach((label) => assert.ok(list.includes(label), label));
  const lacks = (list, ...unwanted) =>
    unwanted.forEach((label) => assert.ok(!list.includes(label), label));

  // Show All UI shows every link, child and button whatever the user's claims,
  // including features that are off.
  const everything = await page.evaluate(() => window.configLabels);
  for (const ctx of [
    {},
    { claims: ["admin"] },
    { claims: ["super_po", "kpi"], isBranchManager: true, jobRequestsEnabled: true },
  ]) {
    assert.deepEqual(await labels(ctx, true), everything);
  }
  has(
    everything,
    "Purchase Orders/PO Groups",
    "Purchase Orders/PO Groups+New PO Group",
    "Business/WIP/Branch",
    "Business/Job Create Queue",
    "Settings/Currencies",
  );

  // No claims: only the links everyone may use.
  const none = await labels({});
  has(
    none,
    "Time Management/Entries",
    "Purchase Orders/My Purchase Orders",
    "Business/Jobs",
    "Business/Jobs+Jobs Map",
    "Business/WIP/My WIP",
    "Business/Claims",
  );
  lacks(
    none,
    "Purchase Orders/PO Groups",
    "Business/Jobs+New Project or Proposal",
    "Business/WIP/Branch",
    "Business/Abandoned Requests",
    "Business/PO Approvers",
    "Expenses/Uncommit Requests",
  );
  assert.ok(!none.some((label) => label.startsWith("Reports/")), "Reports is hidden");
  assert.ok(!none.some((label) => label.startsWith("Settings/")), "Settings is hidden");

  // PO Groups: super_po holders create them, admin holders only look.
  has(await labels({ claims: ["super_po"] }), "Purchase Orders/PO Groups+New PO Group");
  const admin = await labels({ claims: ["admin"] });
  has(admin, "Purchase Orders/PO Groups", "Settings/Currencies", "Business/Abandoned Requests");
  lacks(admin, "Purchase Orders/PO Groups+New PO Group");

  // A branch manager sees the branch reports and nothing else in Reports.
  const manager = await labels({ isBranchManager: true });
  has(manager, "Business/WIP/Branch", "Business/WIP/Division");
  assert.deepEqual(
    manager.filter((label) => label.startsWith("Reports/")),
    ["Reports/Branch Hours"],
  );
  has(await labels({ claims: ["kpi"] }), "Reports/KPI", "Reports/Branch Hours");

  // Project request links need the feature on as well as the claim.
  lacks(await labels({ claims: ["time", "job"] }), "Business/Project Requests");
  const requests = await labels({ claims: ["time"], jobRequestsEnabled: true });
  has(requests, "Business/Project Requests");
  lacks(requests, "Business/Request Approvals", "Business/Job Create Queue");
  lacks(await labels({ claims: ["job"], jobRequestsEnabled: true }), "Business/Job Create Queue");

  // Rate Sheets: procurement, admins and revisers always see it. Requesters
  // and approvers see it while project requests are on.
  const rateSheets = "Business/Rate Sheets";
  lacks(none, rateSheets);
  has(admin, rateSheets);
  has(await labels({ claims: ["job"] }), rateSheets);
  has(await labels({ claims: ["rate_sheet_revise"] }), rateSheets);
  has(requests, rateSheets);
  has(await labels({ claims: ["project_authorization"], jobRequestsEnabled: true }), rateSheets);
  lacks(await labels({ claims: ["time", "project_authorization"] }), rateSheets);

  // One plus starts a project: beside Jobs, never beside Project Requests.
  // Procurement always sees it; staff with time see it while requests are on.
  const newProject = "Business/Jobs+New Project or Proposal";
  has(requests, newProject);
  lacks(await labels({ claims: ["time"] }), newProject);
  has(await labels({ claims: ["job"] }), newProject);
  has(await labels({ claims: ["job"], jobRequestsEnabled: true }), newProject);
  assert.ok(
    !everything.some((label) => label.startsWith("Business/Project Requests+")),
    "Project Requests has no plus",
  );
  has(
    await labels({
      claims: ["job"],
      jobRequestsEnabled: true,
      attentionCounts: { "/job-requests?queue=creation": 0 },
    }),
    "Business/Job Create Queue",
  );

  assert.deepEqual(errors, []);
});
console.log("nav visibility browser tests passed");
