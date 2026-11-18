import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";
await runBrowserHarness(async (page, origin) => {
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.context().grantPermissions(["clipboard-read", "clipboard-write"], { origin });
  await page.goto(`${origin}/tests/fixtures/copyButton.html`);
  const button = page.getByRole("button", { name: "Copy code 2026-02-00001", exact: true });
  await button.waitFor();
  // The button keeps the line's height.
  const lineHeight = async (selector) =>
    page.locator(selector).evaluate((element) => element.getBoundingClientRect().height);
  assert.equal(await lineHeight("[data-label]"), await lineHeight("[data-plain]"));

  await button.click();
  assert.equal(await button.getAttribute("data-copy-state"), "copied");
  assert.equal(await page.evaluate(() => navigator.clipboard.readText()), "2026-02-00001");
  assert.match(await page.locator("[aria-live]").innerText(), /Copied/);
  // The click neither follows the link nor reaches the row.
  assert.equal(new URL(page.url()).hash, "");
  assert.equal(await page.locator("[data-row-clicks]").innerText(), "0");
  // The mark goes back after a moment.
  await page.waitForFunction(
    () => document.querySelector("[data-copy-state]")?.getAttribute("data-copy-state") === "idle",
  );

  // Enter and Space copy too.
  for (const key of ["Enter", "Space"]) {
    await page.evaluate(() => navigator.clipboard.writeText(""));
    await button.focus();
    await page.keyboard.press(key);
    await page.waitForFunction(
      () =>
        document.querySelector("[data-copy-state]")?.getAttribute("data-copy-state") === "copied",
    );
    assert.equal(await page.evaluate(() => navigator.clipboard.readText()), "2026-02-00001", key);
  }

  // A browser that refuses the copy shows a failure, and nothing throws.
  await page.evaluate(() => {
    Object.defineProperty(navigator, "clipboard", {
      configurable: true,
      value: { writeText: () => Promise.reject(new Error("denied")) },
    });
  });
  await page.waitForFunction(
    () => document.querySelector("[data-copy-state]")?.getAttribute("data-copy-state") === "idle",
  );
  await button.click();
  assert.equal(await button.getAttribute("data-copy-state"), "failed");
  assert.match(await page.locator("[aria-live]").innerText(), /Copy failed/);

  // A page without a clipboard fails the same way.
  await page.evaluate(() => {
    Object.defineProperty(navigator, "clipboard", { value: undefined });
  });
  await page.waitForFunction(
    () => document.querySelector("[data-copy-state]")?.getAttribute("data-copy-state") === "idle",
  );
  await button.click();
  assert.equal(await button.getAttribute("data-copy-state"), "failed");

  // It doesn't print.
  await page.emulateMedia({ media: "print" });
  assert.equal(await button.isVisible(), false);
  assert.deepEqual(errors, []);
  console.log("DSCopyButton copy, keyboard, failure, propagation and print checks passed");
});
