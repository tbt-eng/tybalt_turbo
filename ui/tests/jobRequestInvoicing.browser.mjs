// Procurement sets the invoicing profile of a request in the Job Create Queue
// from the request page. Holders of job see that a profile is missing.
import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

const profiles = [
  { id: "prof1", client: "client1", contact: "contact1", name: "Head office", created: "" },
  { id: "prof2", client: "client1", contact: "contact1", name: "Field office", created: "" },
];

await runBrowserHarness(
  async (page, origin) => {
    const calls = [];
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    await page.route("**/api/**", (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const json = (body, status = 200) =>
        route.fulfill({ status, contentType: "application/json", body: JSON.stringify(body) });
      const list = (items) =>
        json({ page: 1, perPage: 500, totalItems: items.length, totalPages: 1, items });
      if (path.includes("/collections/app_config/records"))
        return list([
          {
            id: "jobs",
            key: "jobs",
            value: { require_job_request_approval: true, create_edit_absorb: true },
          },
        ]);
      if (path === "/api/collections/client_invoicing_information/records") return list(profiles);
      if (path === "/api/collections/client_contacts/records")
        return list([{ id: "contact1", client: "client1", given_name: "Pat", surname: "Billing" }]);
      if (path === "/api/job_requests/request1/invoicing_profile") {
        calls.push({ method: request.method(), body: request.postDataJSON() });
        return json({});
      }
      return json({ items: [] });
    });
    const show = (permissions, related) =>
      page.evaluate(
        ([permissions, related]) => {
          const data = window.requestData("Awaiting creation", permissions, related);
          data.detail.request.client = "client1";
          window.setRequestData(data);
        },
        [permissions, related],
      );
    await page.goto(`${origin}/tests/fixtures/jobRequestReview.html`);
    await page.waitForFunction(() => !!window.setRequestData);

    const missing = page.getByText("No invoicing profile.", { exact: true });
    const set = page.getByRole("button", { name: "Set profile", exact: true });
    // A viewer without the job claim sees neither the warning nor the button.
    await show({ can_create: false, can_set_invoicing: false }, {});
    await page.getByText("Invoicing", { exact: true }).waitFor();
    assert.equal(await missing.count(), 0);
    assert.equal(await set.count(), 0);

    // Procurement sees both.
    await page.evaluate(() => window.setClaims(["job"]));
    await show({ can_create: true, can_set_invoicing: true }, {});
    await missing.waitFor();
    await set.click();
    const dialog = page.getByRole("dialog", { name: "Invoicing profile" });
    await dialog.waitFor();
    assert.ok(
      await dialog.getByText("The approval stays valid. The job takes this profile.").isVisible(),
    );
    const create = dialog.getByRole("link", { name: "Create an invoicing profile", exact: true });
    assert.equal(
      await create.getAttribute("href"),
      "/clients/client1/invoicing/add?return_to=%2Fclients%2Fclient1%2Fdetails",
    );
    assert.equal(await create.getAttribute("target"), "_blank");
    const save = dialog.getByRole("button", { name: "Save profile", exact: true });
    // Nothing to save until the choice differs from the request's profile.
    assert.equal(await save.isDisabled(), true);
    await dialog.getByLabel("Invoicing profile").selectOption("prof2");
    await save.click();
    await dialog.waitFor({ state: "detached" });
    assert.deepEqual(calls, [
      { method: "POST", body: { version: 3, invoicing_information: "prof2" } },
    ]);

    // With a profile, the button changes it, and an empty choice clears it.
    await page.evaluate(() => {
      const data = window.requestData(
        "Awaiting creation",
        { can_create: true, can_set_invoicing: true },
        {
          invoicing_information: { id: "prof2", client: "client1", name: "Field office" },
        },
      );
      data.detail.request.client = "client1";
      data.detail.request.invoicing_information = "prof2";
      window.setRequestData(data);
    });
    assert.equal(await missing.count(), 0);
    await page.getByRole("button", { name: "Change profile", exact: true }).click();
    await dialog.waitFor();
    await dialog.getByLabel("Invoicing profile").selectOption("");
    await dialog.getByRole("button", { name: "Save profile", exact: true }).click();
    await dialog.waitFor({ state: "detached" });
    assert.deepEqual(calls[1], { method: "POST", body: { version: 3, invoicing_information: "" } });
    assert.deepEqual(errors, []);
  },
  {
    services: {
      "$app/navigation":
        "export const goto = (url) => { window.lastNavigation = url; }; export const beforeNavigate = () => {}; export const invalidateAll = async () => {}; export const replaceState = (url, state) => history.replaceState(state, '', url);",
    },
  },
);
console.log("Job request invoicing browser checks passed.");
