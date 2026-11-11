import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// A time sheet rejected after its approval stays on Approved By Me. Its row
// shows the rejection, not the approval. Pending My Approval and Shared with Me
// show a rejection the same way.
// Set SCREENSHOT to a .png path to save the Approved By Me list.
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

const tally = (id, given_name, surname, extra = {}) => ({
  id,
  uid: `owner-${id}`,
  approver: "approver1",
  week_ending: "2026-11-07",
  submitted: true,
  approved: "",
  rejected: "",
  rejection_reason: "",
  rejector_name: "",
  committed: "",
  committer: "",
  committer_name: "",
  approver_name: "Pat Approver",
  given_name,
  surname,
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

const rejection = {
  rejected: "2026-11-09 15:30:00.000Z",
  rejection_reason: "Wrong job on Tuesday",
  rejector_name: "Chris Committer",
};

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

    await page.route("**/api/**", async (route) => {
      const path = new URL(route.request().url()).pathname;
      const json = (body) =>
        route.fulfill({
          status: 200,
          contentType: "application/json",
          headers: { "access-control-allow-origin": "*" },
          body: JSON.stringify(body),
        });
      if (path.endsWith("/auth-refresh")) return json({ token, record });
      if (path === "/api/time_sheets/tallies/approved") {
        return json([
          tally("ts1", "Ann", "Smith", { approved: "2026-11-08 10:00:00.000Z" }),
          tally("ts2", "Ben", "Jones", { approved: "2026-11-08 10:00:00.000Z", ...rejection }),
          tally("ts3", "Cara", "Lee", {
            approved: "2026-11-08 10:00:00.000Z",
            committed: "2026-11-10 08:00:00.000Z",
            committer: "committer1",
            committer_name: "Dana Payroll",
          }),
        ]);
      }
      if (path === "/api/time_sheets/tallies/pending") {
        return json([
          tally("ts4", "Eve", "Park"),
          tally("ts5", "Finn", "Ray", { ...rejection, rejector_name: "Pat Approver" }),
        ]);
      }
      if (path === "/api/time_sheets/tallies/shared") {
        return json([
          tally("ts6", "Gail", "Wu", { approver: "boss", approved: "2026-11-08 10:00:00.000Z" }),
          tally("ts7", "Hal", "Ng", { approver: "boss", ...rejection }),
        ]);
      }
      if (path.endsWith("/records")) {
        return json({ page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] });
      }
      return json({});
    });

    const row = (name) => page.locator("li.contents", { hasText: name });

    // Approved By Me
    await page.goto(`${origin}/time/sheets/approved`);
    await page.getByText("Approved Time Sheets").waitFor();
    await row("Ben Jones").waitFor();

    // Approved: green Approved with its date, no rejection.
    const approved = await row("Ann Smith").innerText();
    assert.match(approved, /Approved\s+on Nov 8, 2026/);
    assert.doesNotMatch(approved, /Rejected/);

    // Rejected after approval: red Rejected with who, when and why, and no Approved.
    const rejected = row("Ben Jones");
    const rejectedText = await rejected.innerText();
    assert.match(rejectedText, /Rejected\s+by Chris Committer\s+on Nov 9, 2026/);
    assert.match(rejectedText, /Wrong job on Tuesday/);
    assert.doesNotMatch(rejectedText, /Approved/);
    const label = rejected.getByText("Rejected", { exact: true });
    assert.match(await label.getAttribute("class"), /\bbg-red-100\b/);

    // Committed: Approved and Committed, no rejection.
    const committed = await row("Cara Lee").innerText();
    assert.match(committed, /Approved\s+on Nov 8, 2026/);
    assert.match(committed, /Committed\s+on Nov 10, 2026 by Dana Payroll/);
    assert.doesNotMatch(committed, /Rejected/);

    if (process.env.SCREENSHOT) {
      await page.locator("ul").first().screenshot({ path: process.env.SCREENSHOT });
    }

    // Pending My Approval: a time sheet rejected before approval says so.
    await page.goto(`${origin}/time/sheets/pending`);
    await row("Finn Ray").waitFor();
    const pendingRejected = await row("Finn Ray").innerText();
    assert.match(pendingRejected, /Rejected\s+by Pat Approver\s+on Nov 9, 2026/);
    assert.match(pendingRejected, /Wrong job on Tuesday/);
    assert.doesNotMatch(await row("Eve Park").innerText(), /Rejected/);

    // Shared with Me: a rejected time sheet says so.
    await page.goto(`${origin}/time/sheets/shared`);
    await row("Hal Ng").waitFor();
    const sharedRejected = await row("Hal Ng").innerText();
    assert.match(sharedRejected, /Rejected\s+by Chris Committer\s+on Nov 9, 2026/);
    assert.match(sharedRejected, /Wrong job on Tuesday/);
    assert.doesNotMatch(await row("Gail Wu").innerText(), /Rejected/);

    assert.deepEqual(errors, []);
    console.log(
      "Time sheet rejection checks passed: Approved By Me shows Rejected (by, on, reason) instead of Approved; Pending My Approval and Shared with Me show Rejected.",
    );
  },
  { app: true },
);
