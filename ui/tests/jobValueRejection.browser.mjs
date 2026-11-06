import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

await runBrowserHarness(async (page, origin) => {
  let items = ["A", "B"].map((id) => ({
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
  let failRejection = true;
  page.on("pageerror", (error) => errors.push(error.message));
  await page.route("**/api/**", (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;
    if (path.endsWith("/reject")) {
      actions.push({ path, body: request.postDataJSON() });
      if (failRejection) {
        failRejection = false;
        return route.fulfill({
          status: 400,
          contentType: "application/json",
          body: JSON.stringify({ message: "Rejection failed. Try again." }),
        });
      }
      items = items.filter((item) => path !== `/api/job_attachments/${item.id}/reject`);
      return route.fulfill({ contentType: "application/json", body: "{}" });
    }
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ items: path === "/api/job_attachments" ? items : [] }),
    });
  });
  await page.goto(`${origin}/tests/fixtures/projectRequests.html?values&queue`);
  const row = (id) => page.getByRole("listitem").filter({ hasText: `Project ${id}` });
  const dialog = page.getByRole("dialog", { name: "Reject value change" });
  const reason = dialog.getByLabel("Rejection reason", { exact: true });
  const submit = dialog.getByRole("button", { name: "Reject", exact: true });
  const open = (id) => row(id).getByRole("button", { name: "Reject", exact: true }).click();

  // Each opening starts with an empty reason, and rejection needs one.
  await open("B");
  assert.equal(await submit.isDisabled(), true);
  await reason.fill("Not this one yet.");
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();
  await open("B");
  assert.equal(await reason.inputValue(), "");
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();

  // A failed rejection keeps the reason, so the approver can try again.
  await open("A");
  assert.equal(await reason.inputValue(), "");
  await reason.fill("A needs the signed appendix.");
  await submit.click();
  await dialog.getByText("Rejection failed. Try again.").waitFor();
  assert.equal(await reason.inputValue(), "A needs the signed appendix.");
  await submit.click();
  await row("A").waitFor({ state: "hidden" });

  await open("B");
  assert.equal(await reason.inputValue(), "");
  await reason.fill("B has the wrong value.");
  await submit.click();
  await row("B").waitFor({ state: "hidden" });
  assert.deepEqual(
    actions,
    ["A", "A", "B"].map((id) => ({
      path: `/api/job_attachments/${id}/reject`,
      body: {
        version: 1,
        document_hash: `hash${id}`,
        review_token: `token${id}`,
        rejection_reason: id === "A" ? "A needs the signed appendix." : "B has the wrong value.",
      },
    })),
  );
  assert.deepEqual(errors, []);
});
console.log(
  "Value rejection checks passed: a reason per opening, failed retry, and queue refresh.",
);
