import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { runBrowserHarness } from "./browserHarness.mjs";

const record = {
  id: "docsuser",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "docs@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");
const whatsNew = readFileSync(new URL("../src/docs/whats-new.md", import.meta.url), "utf8");
const commits = [...whatsNew.matchAll(/^## .* \(to ([0-9a-f]{8,40})\)$/gm)].map(
  (match) => match[1],
);
const latest = commits[0];
const storageKey = "help_seen_whats_new";

await runBrowserHarness(
  async (page, origin) => {
    page.setDefaultTimeout(15000);
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    await page.addInitScript(
      ({ token, record }) =>
        localStorage.setItem("pocketbase_auth", JSON.stringify({ token, record })),
      { token, record },
    );
    await page.route("**/api/**", (route) => {
      const path = new URL(route.request().url()).pathname;
      const body = path.endsWith("/auth-refresh")
        ? { token, record }
        : path.endsWith("/records")
          ? { page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] }
          : {};
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        headers: { "access-control-allow-origin": "*" },
        body: JSON.stringify(body),
      });
    });

    const dot = page.locator('[role="img"][aria-label="New in Help"]:visible');
    const helpNav = page.locator('nav[aria-label="Help"]');
    const whatsNewLink = helpNav.getByRole("link", { name: /What's new/ });
    const newLabels = page.locator("section").getByText("New", { exact: true });
    const seen = () => page.evaluate((key) => localStorage.getItem(key), storageKey);
    const setSeen = (value) =>
      page.evaluate(([key, value]) => localStorage.setItem(key, value), [storageKey, value]);

    // A browser that has never opened What's new shows the bubble.
    await page.goto(`${origin}/`);
    await dot.first().waitFor();
    assert.equal(await seen(), null);

    // Opening Help keeps it, and points to What's new in the side column.
    await page.locator('a[href="/docs"]:visible').first().click();
    await page.waitForURL("**/docs");
    await whatsNewLink.locator('[role="img"]').waitFor();
    assert.equal(await seen(), null);

    // Opening What's new clears both dots, saves the newest week and marks only it.
    await whatsNewLink.click();
    await page.waitForURL("**/docs/whats-new");
    await newLabels.first().waitFor();
    assert.equal(await dot.count(), 0);
    assert.equal(await seen(), latest);
    assert.equal(await newLabels.count(), 1);

    // Leaving and coming back in the same session marks nothing.
    await helpNav.getByRole("link", { name: "Start here" }).click();
    await page.waitForURL("**/docs/start-here");
    await whatsNewLink.click();
    await page.waitForURL("**/docs/whats-new");
    await page.waitForTimeout(500);
    assert.equal(await newLabels.count(), 0);

    // It stays cleared after a reload.
    await page.goto(`${origin}/`);
    await page.locator('a[href="/docs"]:visible').first().waitFor();
    await page.waitForTimeout(1000);
    assert.equal(await dot.count(), 0);

    // A newer week brings it back, and a commit no longer on the page marks every week.
    await setSeen("00000000");
    await page.reload();
    await dot.first().waitFor();
    await page.goto(`${origin}/docs/whats-new`);
    await newLabels.first().waitFor();
    assert.equal(await newLabels.count(), commits.length);
    assert.equal(await dot.count(), 0);

    // A commit still on the page marks only the weeks after it.
    if (commits.length > 1) {
      await setSeen(commits[1]);
      await page.reload();
      await newLabels.first().waitFor();
      assert.equal(await newLabels.count(), 1);
    }

    // On a phone the top bar shows the bubble.
    await page.setViewportSize({ width: 375, height: 800 });
    await page.evaluate((key) => localStorage.removeItem(key), storageKey);
    await page.goto(`${origin}/`);
    await dot.first().waitFor();
    assert.deepEqual(errors, []);
    console.log(
      "Help bubble checks passed: first visit, What's new, labels, reload, new week and phone.",
    );
  },
  { app: true },
);
