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
      if (path === "/api/job_requests/request1/submit")
        return route.fulfill({
          status: 400,
          contentType: "application/json",
          body: JSON.stringify({ message: "Select an approver before you submit." }),
        });
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

    // The approver confirms the approval checklist and has no duplicate review.
    const duplicateReview = page.getByRole("heading", { name: "Duplicate-work review" });
    assert.equal(await boxes.count(), 6);
    assert.equal(await duplicateReview.count(), 0);
    assert.equal(await approve.isDisabled(), true);
    for (const box of await boxes.all()) await box.check();
    assert.equal(await approve.isDisabled(), false);
    await approve.click();
    await page.waitForFunction(() => !document.querySelector("button[disabled]"));
    assert.equal(actions[0].path, "/api/job_requests/request1/approve");
    assert.equal(actions[0].body.checklist_version, 1);
    assert.equal(Object.keys(actions[0].body.checklist_answers).length, 6);

    // Approval keeps the version and review token. A job holder answers the
    // creation checklist, from empty boxes even if they also approved the request.
    await page.evaluate(() =>
      window.setRequestData(window.requestData("Awaiting creation", { can_create: true })),
    );
    await create.waitFor();
    assert.equal(await boxes.count(), 7);
    assert.equal(await duplicateReview.count(), 1);
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

    // The creator's draft shows read only, as for everyone else. The form is
    // on the edit page, so every action here works on saved data.
    await page.evaluate(() => {
      window.lastNavigation = undefined;
      const data = window.requestData("Draft", {
        can_edit: true,
        can_submit: true,
        can_delete: true,
      });
      data.detail.request.description = "Inspect the bridge";
      data.detail.related = {
        manager: { name: "Morgan Manager" },
        alternate_manager: { name: "Avery Alternate" },
        job_owner: { name: "Owner Corp" },
      };
      window.setRequestData(data);
    });
    const editRequest = page.getByRole("button", { name: "Edit request", exact: true });
    await editRequest.waitFor();
    assert.equal(await page.locator("form").count(), 0, "no form on the request's page");
    assert.equal(await page.getByRole("button", { name: "Save request" }).count(), 0);
    assert.equal(await page.getByText("Inspect the bridge", { exact: true }).count(), 1);
    const cards = await page.locator("#app").innerText();
    assert.match(cards, /Alternate manager\s+Avery Alternate/);
    assert.match(cards, /Job owner\s+Owner Corp/);
    assert.doesNotMatch(cards, /Invoicing profile/, "an empty optional field is hidden");
    const submit = page.getByRole("button", { name: "Submit", exact: true });
    for (const button of [
      submit,
      page.getByRole("button", { name: "Replace PDF", exact: true }),
      page.getByRole("button", { name: "Remove PDF", exact: true }),
      page.getByRole("button", { name: "Print Project Authorization", exact: true }),
    ])
      assert.equal(await button.isDisabled(), false);
    assert.equal(await page.getByText(/unsaved changes/i).count(), 0);
    await editRequest.click();
    assert.equal(await page.evaluate(() => window.lastNavigation), "/job-requests/request1/edit");

    // A refused submit says why beside the workflow buttons, not at the top.
    await submit.click();
    const workflow = page
      .locator("section")
      .filter({ has: page.getByRole("heading", { name: "Workflow actions", exact: true }) });
    await workflow
      .getByRole("alert")
      .filter({ hasText: "Select an approver before you submit." })
      .waitFor();
    assert.equal(
      await page.getByRole("alert").count(),
      1,
      "the refusal shows only in Workflow actions",
    );
    assert.equal(actions.at(-1).path, "/api/job_requests/request1/submit");

    // A rejected request has no Submit until it changes. The header offers Edit.
    await page.evaluate(() => {
      window.lastNavigation = undefined;
      const data = window.requestData("Rejected", { can_edit: true, can_delete: true });
      data.detail.request.rejection_reason = "Wrong client";
      window.setRequestData(data);
    });
    await page.getByText("Rejection reason: Wrong client", { exact: false }).waitFor();
    assert.equal(await submit.count(), 0);
    assert.equal(await editRequest.count(), 1, "only the header offers Edit");
    await editRequest.click();
    assert.equal(await page.evaluate(() => window.lastNavigation), "/job-requests/request1/edit");

    // A viewer who can't edit gets no Edit.
    await page.evaluate(() =>
      window.setRequestData(window.requestData("Draft", { can_delete: false })),
    );
    await page.getByText("Untitled request", { exact: true }).waitFor();
    assert.equal(await editRequest.count(), 0);
    assert.equal(await page.getByRole("button", { name: /Upload PDF|Replace PDF/ }).count(), 0);

    // Every viewer reads who acts next, by the person's name. The fixture
    // signs in as user1, so these people are someone else until a case says so.
    const people = {
      creator: { id: "creator1", given_name: "Casey", surname: "Creator" },
      approver: { id: "approver1", given_name: "Avery", surname: "Approver" },
    };
    const showRequest = (status, permissions, change = {}) =>
      page.evaluate(
        ({ status, permissions, people, change }) => {
          const data = window.requestData(status, permissions);
          data.detail.request.creator = change.creator ?? people.creator.id;
          data.detail.request.approver = change.approver ?? people.approver.id;
          data.detail.related = people;
          if (change.reason) data.detail.approval_blocked_reason = change.reason;
          window.setRequestData(data);
        },
        { status, permissions, people, change },
      );
    const workflowButtons = workflow.getByRole("button", {
      name: /Approve request|Create Job|Submit|Recall for editing|Reject/,
    });
    const helpTarget = async () => {
      await workflow.getByRole("button", { name: /^Help/ }).click();
      const href = await page
        .getByRole("dialog")
        .getByRole("link", { name: "Open the full page", exact: true })
        .getAttribute("href");
      await page.keyboard.press("Escape");
      return href;
    };

    // Someone who opens another person's request from the All tab has no
    // action. The section still says who acts next and helps with the statuses.
    await showRequest("Awaiting approval", {});
    await workflow.getByText("Waiting for approval by Avery Approver.", { exact: true }).waitFor();
    assert.equal(await workflowButtons.count(), 0, "a viewer with no permission has no buttons");
    assert.equal(await helpTarget(), "/docs/project-requests#request-statuses");

    // The selected approver reads that it waits for them.
    await showRequest("Awaiting approval", { can_approve: true }, { approver: "user1" });
    await workflow.getByText("Waiting for your approval.", { exact: true }).waitFor();
    assert.equal(await approve.count(), 1);
    assert.equal(await helpTarget(), "/docs/project-requests#approve-a-request");

    // An approver whose limit is too low gets the server's reason, not a button.
    const reason = "Your approval limit does not cover this project value.";
    await showRequest("Awaiting approval", {}, { approver: "user1", reason });
    await workflow.getByText(reason, { exact: true }).waitFor();
    assert.equal(await approve.count(), 0);
    assert.equal(await workflow.getByText("Waiting for your approval.").count(), 1);

    // Other viewers never get a reason, and the page shows none without one.
    await showRequest("Awaiting approval", {});
    await workflow.getByText("Waiting for approval by Avery Approver.", { exact: true }).waitFor();
    assert.equal(await workflow.getByText(reason).count(), 0);

    // A creator who can only recall still gets help to change the request.
    await showRequest("Awaiting creation", { can_recall: true }, { creator: "user1" });
    await workflow
      .getByText("Approved. Waiting for procurement to create the job.", { exact: true })
      .waitFor();
    assert.equal(await helpTarget(), "/docs/project-requests#change-a-request");

    // Someone else's draft and rejected request name the creator.
    await showRequest("Draft", {});
    await workflow
      .getByText("Casey Creator has not submitted this request.", { exact: true })
      .waitFor();
    await showRequest("Rejected", { can_edit: true, can_delete: true }, { creator: "user1" });
    await workflow
      .getByText("Rejected. Correct the request, then submit it again.", { exact: true })
      .waitFor();
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
