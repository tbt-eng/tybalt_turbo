// The request page shows an unawarded proposal to everyone, and lets a job
// holder award it: from the notice, or when Create Job is refused for it.
import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

const proposal = (status) => ({
  proposal: { id: "prop1", number: "P25-001", description: "Bridge study", status },
});

await runBrowserHarness(
  async (page, origin) => {
    const calls = [];
    const errors = [];
    // What the next award and the next creation answer.
    let awardFails = false;
    let createRefused = false;
    page.on("pageerror", (error) => errors.push(error.message));
    await page.route("**/api/**", (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const json = (body, status = 200) =>
        route.fulfill({ status, contentType: "application/json", body: JSON.stringify(body) });
      if (path.includes("/collections/app_config/records"))
        return json({
          page: 1,
          perPage: 500,
          totalItems: 1,
          totalPages: 1,
          items: [{ id: "jobs", key: "jobs", value: { require_job_request_approval: true } }],
        });
      if (path === "/api/collections/jobs/records/prop1") {
        calls.push({ path, method: request.method(), body: request.postDataJSON() });
        if (awardFails)
          return json(
            {
              status: 400,
              message: "proposal value or time and materials required",
              data: {
                proposal_value: { code: "value_required_for_status", message: "needs a value" },
              },
            },
            400,
          );
        return json({ id: "prop1", status: "Awarded" });
      }
      if (path === "/api/job_requests/request1/create") {
        calls.push({ path, method: request.method(), body: request.postDataJSON() });
        if (createRefused) {
          createRefused = false;
          const message = "Proposal P25-001 is Submitted. Award it before you create the job.";
          return json(
            {
              status: 400,
              message,
              data: {
                proposal: { code: "proposal_not_awarded", message, data: { proposal_id: "prop1" } },
              },
            },
            400,
          );
        }
        return json({ id: "job1", job_id: "job1" });
      }
      return json({ items: [] });
    });
    const show = (status) =>
      page.evaluate(
        (related) =>
          window.setRequestData(
            window.requestData("Awaiting creation", { can_create: true }, related),
          ),
        proposal(status),
      );
    await page.goto(`${origin}/tests/fixtures/jobRequestReview.html`);
    await page.waitForFunction(() => !!window.setRequestData);
    await show("Submitted");

    // Everyone sees the status. Only a job holder gets the button.
    const notice = page.getByText("Proposal P25-001 is Submitted", { exact: true });
    const award = page.getByRole("button", { name: "Award proposal", exact: true });
    await notice.waitFor();
    assert.equal(await award.count(), 0);
    await page.evaluate(() => window.setClaims(["job"]));
    await award.waitFor();

    // A refused award shows the reason and a link to the proposal.
    awardFails = true;
    await award.click();
    const failed = page.getByText("Cannot award the proposal", { exact: true });
    await failed.waitFor();
    assert.ok(await page.getByText(/needs a value/).isVisible());
    assert.equal(
      await page.getByRole("link", { name: "Edit the proposal" }).getAttribute("href"),
      "/jobs/prop1/edit",
    );
    assert.deepEqual(calls[0], {
      path: "/api/collections/jobs/records/prop1",
      method: "PATCH",
      body: { status: "Awarded" },
    });

    // A successful award clears the failure. The reload shows it as Awarded.
    awardFails = false;
    await award.click();
    await failed.waitFor({ state: "detached" });
    assert.equal(calls.length, 2);
    await show("Awarded");
    await notice.waitFor({ state: "detached" });
    assert.equal(await award.count(), 0);

    // Create Job refused for the proposal offers to award it and continue.
    await show("Submitted");
    for (const box of await page.getByRole("checkbox").all()) await box.check();
    createRefused = true;
    await page.getByRole("button", { name: "Create Job", exact: true }).click();
    const confirm = page.getByRole("button", { name: "Award and create", exact: true });
    await confirm.waitFor();
    assert.ok(await page.getByText(/Set it to Awarded and continue\?/).isVisible());
    await confirm.click();
    await page.waitForFunction(() => !!window.lastNavigation);
    assert.deepEqual(
      calls.slice(2).map((call) => `${call.method} ${call.path}`),
      [
        "POST /api/job_requests/request1/create",
        "PATCH /api/collections/jobs/records/prop1",
        "POST /api/job_requests/request1/create",
      ],
    );
    // The retry sends the same digest and answers as the refused creation.
    assert.deepEqual(calls[4].body, calls[2].body);
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
console.log("Job request award browser checks passed.");
