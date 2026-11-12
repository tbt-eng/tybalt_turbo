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
    const apiCalls = [];
    await page.route("**/api/**", (route) => {
      const path = new URL(route.request().url()).pathname;
      apiCalls.push(path);
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

    const input = page.getByRole("textbox", { name: "Search Help" });
    const results = page
      .getByRole("search")
      .locator("xpath=following-sibling::div")
      .locator("ul a");
    const status = page.getByRole("status");
    const pageTitle = (name) => page.getByRole("heading", { level: 1, name, exact: true });

    await page.goto(`${origin}/docs`);
    await pageTitle("Start here").waitFor();
    await input.waitFor();

    // Typing replaces the page with hits, and sends nothing to the server.
    const callsBefore = apiCalls.length;
    await input.fill("timesheet");
    await results.first().waitFor();
    assert.equal(await pageTitle("Start here").count(), 0);
    assert.match((await status.textContent()).trim(), /^\d+ results?$/);
    assert.ok((await results.count()) <= 8);
    assert.match(await results.first().getAttribute("href"), /^\/docs\/time(#|$)/);
    assert.ok((await results.first().locator("mark").count()) > 0);
    assert.equal(apiCalls.length, callsBefore);

    // Clear and Escape bring the page back.
    await page.getByRole("button", { name: "Clear" }).click();
    await pageTitle("Start here").waitFor();
    assert.equal(await input.inputValue(), "");
    await input.fill("recall");
    await results.first().waitFor();
    await input.press("Escape");
    await pageTitle("Start here").waitFor();

    // One letter is not a search, and a search with no hit says so.
    await input.fill("t");
    assert.equal(await pageTitle("Start here").count(), 1);
    await input.fill("zzzzqq");
    await status.getByText("Nothing matches. Try other words.").waitFor();
    assert.equal(await results.count(), 0);

    // The arrow keys move the highlight, and Enter opens that hit at its section.
    await input.fill("make a revision");
    await results.first().waitFor();
    assert.equal(await results.first().getAttribute("data-active"), "");
    await input.press("ArrowDown");
    assert.equal(await results.nth(1).getAttribute("data-active"), "");
    await input.press("ArrowUp");
    const href = await results.first().getAttribute("href");
    assert.equal(href, "/docs/revise-a-purchase-order#make-a-revision");
    await input.press("Enter");
    await page.waitForURL(`**${href}`);
    await pageTitle("Revise a purchase order").waitFor();
    assert.equal(await input.inputValue(), "");
    const section = page.locator("#make-a-revision");
    await section.waitFor();
    const top = async () =>
      section.evaluate((element) => {
        const scroller = element.closest(".overflow-y-auto");
        return element.getBoundingClientRect().top - scroller.getBoundingClientRect().top;
      });
    await page.waitForFunction(
      () => {
        const element = document.querySelector("#make-a-revision");
        const scroller = element?.closest(".overflow-y-auto");
        if (!element || !scroller) return false;
        const offset = element.getBoundingClientRect().top - scroller.getBoundingClientRect().top;
        return offset >= 0 && offset < scroller.clientHeight;
      },
      undefined,
      { timeout: 5000 },
    );
    assert.ok((await top()) >= 0);

    // The hits start at the top, and Escape returns the page to where it was.
    const scroller = page
      .locator("#make-a-revision")
      .locator("xpath=ancestor::div[contains(@class,'overflow-y-auto')][1]");
    await scroller.evaluate((element) => (element.scrollTop = 120));
    await input.fill("approve");
    await results.first().waitFor();
    assert.equal(
      await page
        .locator(".overflow-y-auto")
        .last()
        .evaluate((e) => e.scrollTop),
      0,
    );
    await input.press("Escape");
    await pageTitle("Revise a purchase order").waitFor();
    assert.equal(await scroller.evaluate((element) => element.scrollTop), 120);

    // A hit on the page that is open also returns to that page, at the section.
    await input.fill("can't revise");
    await results.first().waitFor();
    assert.equal(
      await results.first().getAttribute("href"),
      "/docs/revise-a-purchase-order#when-you-cant-revise",
    );
    await results.first().click();
    await page.waitForURL("**/docs/revise-a-purchase-order#when-you-cant-revise");
    await page.locator("#when-you-cant-revise").waitFor();
    assert.equal(await input.inputValue(), "");

    // A phone keeps the bar and the hits in the screen's width.
    await page.setViewportSize({ width: 375, height: 800 });
    await input.fill("purchase order");
    await results.first().waitFor();
    assert.equal(
      await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth),
      true,
    );
    assert.deepEqual(errors, []);
    console.log(
      "Docs search checks passed: instant hits, no server call, clear, keys, section links and phone layout.",
    );
  },
  { app: true },
);
