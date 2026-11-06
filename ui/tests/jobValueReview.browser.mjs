import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

await runBrowserHarness(async (page, origin) => {
  let items = ["A", "B", "C"].map((id) => ({
    id,
    job: `job${id}`,
    job_number: `Project ${id}`,
    status: "Pending",
    project_value: 10000,
    assigned_approver: "reviewer",
    version: 1,
    document_hash: `hash${id}`,
    current_review_token: `token${id}`,
    uploaded: "2026-10-14 10:00:00.000Z",
  }));
  const actions = [];
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.route("**/api/**", (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;
    if (path.endsWith("/approve")) {
      actions.push({ path, body: request.postDataJSON() });
      items = items.filter((item) => path !== `/api/job_attachments/${item.id}/approve`);
      // While A was approved, B was corrected with a replacement PDF. C is
      // unchanged.
      if (path === "/api/job_attachments/A/approve") {
        items = items.map((item) =>
          item.id === "B" ? { ...item, version: 2, document_hash: "hashB2" } : item,
        );
      }
      return route.fulfill({ contentType: "application/json", body: "{}" });
    }
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ items: path === "/api/job_attachments" ? items : [] }),
    });
  });
  await page.goto(`${origin}/tests/fixtures/projectRequests.html?values&queue`);
  const row = (id) => page.getByRole("listitem").filter({ hasText: `Project ${id}` });
  const dialog = page.getByRole("dialog", { name: "Approve value change" });
  const reviewed = dialog.getByRole("checkbox");
  const submit = dialog.getByRole("button", { name: "Approve value", exact: true });
  const open = (id) => row(id).getByRole("button", { name: "Approve", exact: true }).click();

  // Approval needs the confirmation, and each opening asks for it again.
  await open("C");
  assert.equal(await submit.isDisabled(), true);
  await reviewed.check();
  assert.equal(await submit.isDisabled(), false);
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();
  await open("C");
  assert.equal(await reviewed.isChecked(), false);
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();

  await open("A");
  await reviewed.check();
  await submit.click();
  await row("A").waitFor({ state: "hidden" });

  // B's approval must send its new revision, not the one the list first showed.
  await open("B");
  assert.equal(await reviewed.isChecked(), false);
  await reviewed.check();
  await submit.click();
  await row("B").waitFor({ state: "hidden" });
  assert.deepEqual(
    actions.map((action) => [action.path, action.body.version, action.body.document_hash]),
    [
      ["/api/job_attachments/A/approve", 1, "hashA"],
      ["/api/job_attachments/B/approve", 2, "hashB2"],
    ],
  );
  assert.deepEqual(errors, []);
});
console.log(
  "Value review checks passed: each opening asks again, and approval sends the new revision.",
);
