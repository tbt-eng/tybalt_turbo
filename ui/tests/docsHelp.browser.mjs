import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

const failure = "Turbo can't load this help.";

await runBrowserHarness(async (page, origin) => {
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  const open = async (doc) => {
    await page.goto(`${origin}/tests/fixtures/docsHelp.html?doc=${encodeURIComponent(doc)}`);
    await page.getByRole("button", { name: /^Help/ }).click();
    return page.getByRole("dialog");
  };
  const fullPage = (panel) => panel.getByRole("link", { name: "Open the full page", exact: true });

  // A valid reference shows its section and links to the full page.
  let panel = await open("project-requests#approve-a-request");
  await panel.getByText("Approve request").waitFor();
  assert.equal(await panel.getByRole("heading").count(), 0);
  assert.equal(await panel.getByText(failure).count(), 0);
  assert.equal(
    await fullPage(panel).getAttribute("href"),
    "/docs/project-requests#approve-a-request",
  );

  // A reference with no matching section says so instead of showing nothing.
  panel = await open("project-requests#no-such-section");
  await panel.getByText(failure).waitFor();
  assert.equal(await fullPage(panel).count(), 1);

  // So does a docs chunk that a new deployment removed.
  await page.route("**/src/lib/docs.ts*", (route) => route.abort());
  panel = await open("project-requests#approve-a-request");
  await panel.getByText(failure).waitFor();
  assert.equal(await fullPage(panel).count(), 1);
  assert.deepEqual(errors, []);
});
console.log("Docs help checks passed: a section, a bad reference, and a failed load.");
