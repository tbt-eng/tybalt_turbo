import assert from "node:assert/strict";
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

    const nav = page.locator('nav[aria-label="Help"]');
    const revisions = nav.getByRole("link", { name: "Revisions", exact: true });
    const parent = nav.getByRole("link", { name: "Purchase orders", exact: true });
    const margin = (link) =>
      link.evaluate((element) => getComputedStyle(element.closest("li")).marginLeft);

    // The column lists the page under its parent, with the short label.
    await page.goto(`${origin}/docs/purchase-orders`);
    await revisions.waitFor();
    const labels = await nav.getByRole("link").allTextContents();
    const trimmed = labels.map((label) => label.trim());
    assert.equal(trimmed.indexOf("Revisions"), trimmed.indexOf("Purchase orders") + 1);
    assert.equal(trimmed.includes("Revise a purchase order"), false);
    assert.equal(await margin(revisions), "12px");
    assert.equal(await margin(parent), "0px");

    // The link keeps the page's URL and title, and marks the current page.
    await revisions.click();
    await page.waitForURL("**/docs/revise-a-purchase-order");
    await page.getByRole("heading", { name: "Revise a purchase order" }).first().waitFor();
    assert.equal(await revisions.getAttribute("aria-current"), "page");
    assert.equal(await parent.getAttribute("aria-current"), null);

    // A phone wraps the links in rows, so nothing indents.
    await page.setViewportSize({ width: 375, height: 800 });
    assert.equal(await margin(revisions), "0px");
    assert.deepEqual(errors, []);
    console.log("Docs nav checks passed: child placement, label, indent, URL and phone layout.");
  },
  { app: true },
);
