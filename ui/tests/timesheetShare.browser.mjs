import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// The signed-in person approves ts1 and ts3. ts2 and ts4 have another approver.
const record = {
  id: "approver1",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "approver@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");

const tally = (id, uid, approver, extra = {}) => ({
  id,
  uid,
  approver,
  week_ending: "2026-11-07",
  submitted: true,
  approved: "",
  rejected: "",
  rejection_reason: "",
  committed: "",
  committer: "",
  committer_name: "",
  approver_name: "",
  given_name: "Ann",
  surname: "Smith",
  shared_reviewer_count: 0,
  work_total_hours: 40,
  non_work_total_hours: 0,
  rb_hours: 0,
  payout_request_amount: 0,
  job_numbers: [],
  divisions: [],
  division_names: [],
  time_types: [],
  time_type_names: [],
  off_rotation_dates: [],
  off_week_dates: [],
  payout_request_dates: [],
  bank_entry_dates: [],
  ...extra,
});

const managers = [
  { id: "approver1", surname: "Approver", given_name: "Pat" },
  { id: "mgra", surname: "Able", given_name: "Alex" },
  { id: "mgrb", surname: "Baker", given_name: "Blair" },
];

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

    // The server, as the browser sees it.
    let reviewers = [];
    let nextReviewer = 1;
    const created = [];
    const deleted = [];
    await page.route("**/api/**", async (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const method = request.method();
      const json = (body, status = 200) =>
        route.fulfill({
          status,
          contentType: "application/json",
          headers: { "access-control-allow-origin": "*" },
          body: JSON.stringify(body),
        });
      if (path.endsWith("/auth-refresh")) return json({ token, record });
      if (path === "/api/time_sheets/tallies/approved") {
        return json([
          tally("ts1", "ann", "approver1", { approved: "2026-11-08" }),
          tally("ts2", "ann", "someoneelse", { approved: "2026-11-08" }),
        ]);
      }
      if (path === "/api/time_sheets/tallies/pending") {
        return json([tally("ts1", "ann", "approver1", { shared_reviewer_count: 2 })]);
      }
      if (path === "/api/time_sheets/tallies") {
        return json([tally("ts3", "approver1", "approver1"), tally("ts4", "approver1", "boss")]);
      }
      if (path === "/api/time_sheets/ts1/reviewers") {
        return json({ owner_name: "Ann Smith", week_ending: "2026-11-07", reviewers });
      }
      if (path === "/api/time_sheets/ts1/details") {
        return json({
          timeSheet: {
            id: "ts1",
            uid: "ann",
            approver: "approver1",
            week_ending: "2026-11-07",
            submitted: true,
            approved: "",
            rejected: "",
            rejection_reason: "",
            committed: "",
          },
          items: [],
          approverInfo: {},
          sharedReviewerCount: reviewers.length,
        });
      }
      if (path === "/api/collections/managers/records") {
        return json({
          page: 1,
          perPage: 500,
          totalItems: managers.length,
          totalPages: 1,
          items: managers,
        });
      }
      if (path === "/api/collections/time_sheet_reviewers/records" && method === "POST") {
        const body = request.postDataJSON();
        created.push(body);
        const manager = managers.find((m) => m.id === body.reviewer);
        const row = {
          id: `rev${nextReviewer++}`,
          time_sheet: body.time_sheet,
          reviewer: body.reviewer,
          reviewed: "",
          surname: manager.surname,
          given_name: manager.given_name,
        };
        reviewers = [...reviewers, row];
        return json(row);
      }
      if (
        path.startsWith("/api/collections/time_sheet_reviewers/records/") &&
        method === "DELETE"
      ) {
        const id = path.split("/").pop();
        deleted.push(id);
        reviewers = reviewers.filter((r) => r.id !== id);
        return route.fulfill({ status: 204, headers: { "access-control-allow-origin": "*" } });
      }
      if (path.endsWith("/records")) {
        return json({ page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] });
      }
      return json({});
    });

    const share = page.getByRole("button", { name: "Share with another manager", exact: true });
    const dialog = page.getByRole("dialog");

    // Approved By Me: only the row that the signed-in person approves can be shared.
    await page.goto(`${origin}/time/sheets/approved`);
    await page.getByText("Approved Time Sheets").waitFor();
    await share.first().waitFor();
    assert.equal(await share.count(), 1);
    await share.click();
    await dialog.getByRole("heading", { name: "Share Ann Smith's time sheet" }).waitFor();
    assert.equal(await dialog.getByRole("heading", { name: "Viewers" }).count(), 0);

    // The approver is not offered as a viewer.
    const manager = dialog.getByRole("combobox", { name: "Manager" });
    await manager.locator("option").nth(1).waitFor({ state: "attached" });
    assert.deepEqual(await manager.locator("option").allInnerTexts(), [
      "Able, Alex",
      "Baker, Blair",
    ]);
    // Nothing is picked until the person picks a manager.
    assert.equal(await manager.evaluate((select) => select.selectedIndex), -1);
    await manager.selectOption("mgrb");
    await dialog.getByRole("button", { name: "Add Viewer", exact: true }).click();
    await dialog.locator("li", { hasText: "Baker, Blair" }).waitFor();
    assert.deepEqual(created, [{ time_sheet: "ts1", reviewer: "mgrb" }]);
    // A viewer is not offered again.
    assert.deepEqual(await manager.locator("option").allInnerTexts(), ["Able, Alex"]);
    await page.getByText("Shared (1)").waitFor();

    await dialog.getByRole("button", { name: "Remove Viewer", exact: true }).click();
    await dialog.locator("li", { hasText: "Baker, Blair" }).waitFor({ state: "detached" });
    assert.deepEqual(deleted, ["rev1"]);
    await page.getByText("Shared (1)").waitFor({ state: "detached" });
    await dialog.getByRole("button", { name: "Close", exact: true }).click();
    await dialog.waitFor({ state: "hidden" });

    // Pending My Approval: Share sits beside Details, and the badge shows the count.
    await page.goto(`${origin}/time/sheets/pending`);
    await share.first().waitFor();
    assert.equal(await share.count(), 1);
    await page.getByRole("button", { name: "Details", exact: true }).waitFor();
    await page.getByText("Shared (2)").waitFor();

    // Sheets: your own time sheets can be shared only where you are the approver.
    await page.goto(`${origin}/time/sheets/list`);
    await share.first().waitFor();
    assert.equal(await share.count(), 1);

    // Details: the approver shares, and the badge follows the viewer count.
    await page.goto(`${origin}/time/sheets/ts1/details`);
    await share.waitFor();
    await share.click();
    await dialog.getByRole("combobox", { name: "Manager" }).selectOption("mgra");
    await dialog.getByRole("button", { name: "Add Viewer", exact: true }).click();
    await dialog.locator("li", { hasText: "Able, Alex" }).waitFor();
    assert.deepEqual(created.at(-1), { time_sheet: "ts1", reviewer: "mgra" });
    await dialog.getByRole("button", { name: "Close", exact: true }).click();
    await page.getByText("Shared (1)").waitFor();

    assert.deepEqual(errors, []);
    console.log(
      "Time sheet share checks passed: Approved By Me, Pending My Approval, Sheets and details show Share only to the approver; add and remove viewers; badge counts.",
    );
  },
  { app: true },
);
