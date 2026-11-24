import assert from "node:assert/strict";
import { mkdirSync } from "node:fs";
import { runBrowserHarness } from "./browserHarness.mjs";

// A PO list labels an Active purchase order that is stale, or a Recurring one
// near or past its end date. The payables lists of exactly those purchase
// orders turn the labels off.
await runBrowserHarness(async (page, origin) => {
  const errors = [];
  // The list subscribes to realtime events, and PocketBase has no server to
  // reach here: it reports "Something went wrong." for that and nothing else.
  page.on("pageerror", (error) => {
    if (error.message !== "Something went wrong.") errors.push(error.message);
  });
  await page.route("**/api/**", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: "{}" }),
  );
  const labels = async (id) => {
    const row = page.locator("li", { hasText: `Vendor ${id}` });
    await row.waitFor();
    const found = [];
    for (const name of ["Stale", "Expiring", "Expired"])
      if (await row.getByText(name, { exact: true }).count()) found.push(name);
    return found;
  };

  await page.goto(`${origin}/tests/fixtures/poAgeLabels.html`);
  assert.deepEqual(await labels("0001"), []);
  assert.deepEqual(await labels("0002"), ["Stale"]);
  assert.deepEqual(await labels("0003"), ["Expiring"]);
  assert.deepEqual(await labels("0004"), ["Stale", "Expired"]);
  assert.deepEqual(await labels("0005"), []);
  assert.deepEqual(await labels("0006"), []);
  if (process.env.SCREENSHOT_DIR) {
    mkdirSync(process.env.SCREENSHOT_DIR, { recursive: true });
    await page.screenshot({
      path: `${process.env.SCREENSHOT_DIR}/po-age-labels.png`,
      fullPage: true,
    });
  }

  await page.goto(`${origin}/tests/fixtures/poAgeLabels.html?no-labels`);
  assert.deepEqual(await labels("0002"), []);
  assert.deepEqual(await labels("0004"), []);

  assert.deepEqual(errors, []);
});
console.log("PO age labels: ok");
