import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

await runBrowserHarness(
  async (page, origin) => {
    const actions = [];
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    await page.route("**/api/**", (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      if (path.includes("/collections/app_config/records"))
        return route.fulfill({
          contentType: "application/json",
          body: JSON.stringify({
            page: 1,
            perPage: 500,
            totalItems: 1,
            totalPages: 1,
            items: [{ id: "jobs", key: "jobs", value: { require_job_request_approval: true } }],
          }),
        });
      if (path.startsWith("/api/job_requests/request1/"))
        actions.push({ path, body: request.postDataJSON() });
      return route.fulfill({
        contentType: "application/json",
        body: JSON.stringify({ items: [], job_id: "job1" }),
      });
    });
    await page.goto(`${origin}/tests/fixtures/jobRequestReview.html`);

    const boxes = page.getByRole("checkbox");
    const approve = page.getByRole("button", { name: "Approve request", exact: true });
    const create = page.getByRole("button", { name: "Create Job", exact: true });
    await approve.waitFor();

    // The approver confirms the seven statements.
    assert.equal(await boxes.count(), 7);
    assert.equal(await approve.isDisabled(), true);
    for (const box of await boxes.all()) await box.check();
    assert.equal(await approve.isDisabled(), false);
    await approve.click();
    await page.waitForFunction(() => !document.querySelector("button[disabled]"));
    assert.equal(actions[0].path, "/api/job_requests/request1/approve");

    // Approval keeps the version and review token. A job holder who also
    // approved the request still starts from empty boxes.
    await page.evaluate(() =>
      window.setRequestData(window.requestData("Awaiting creation", { can_create: true })),
    );
    await create.waitFor();
    assert.equal(await boxes.count(), 7);
    for (const box of await boxes.all()) assert.equal(await box.isChecked(), false);
    assert.equal(await create.isDisabled(), true);

    // Each of the seven statements holds Create Job back.
    for (let skipped = 0; skipped < 7; skipped++) {
      for (let index = 0; index < 7; index++) await boxes.nth(index).setChecked(index !== skipped);
      assert.equal(await create.isDisabled(), true, `box ${skipped} was not required`);
    }
    await boxes.last().check();
    assert.equal(await create.isDisabled(), false);

    await create.click();
    await page.waitForFunction(() => !!window.lastNavigation);
    assert.equal(actions[1].path, "/api/job_requests/request1/create");
    assert.equal(actions[1].body.checklist_version, 2);
    assert.equal(Object.keys(actions[1].body.checklist_answers).length, 7);
    assert.ok(Object.values(actions[1].body.checklist_answers).every((answer) => answer === true));
    assert.equal(await page.evaluate(() => window.lastNavigation), "/jobs/job1/details");
    assert.deepEqual(errors, []);
  },
  {
    services: {
      "$app/navigation":
        "export const goto = (url) => { window.lastNavigation = url; }; export const beforeNavigate = () => {}; export const invalidateAll = async () => {}; export const replaceState = (url, state) => history.replaceState(state, '', url);",
    },
  },
);
console.log("Job request review browser checks passed.");
