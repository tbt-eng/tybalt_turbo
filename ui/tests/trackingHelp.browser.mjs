import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// The Tracking lists have no search bar or tab bar, so Help sits at the far
// right of the blue header of the list.
const record = {
  id: "trackuser",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "track@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");

const weekEnding = "2026-11-07";
const expense = {
  id: "exp1",
  uid: "u1",
  given_name: "Alex",
  surname: "Adams",
  description: "Hotel in Ottawa",
  approved: "2026-11-03 12:00:00.000Z",
  rejected: "",
  committed: "2026-11-05 12:00:00.000Z",
  approver_name: "Pat Lee",
  rejector_name: "",
  rejection_reason: "",
  date: "2026-11-02",
  payment_type: "OnAccount",
  distance: 0,
  cc_last_4_digits: "",
  attachment: "",
  job_number: "",
  division_code: "BM",
  division_name: "Building",
  category: "",
  category_name: "",
  vendor: "",
  vendor_name: "",
  vendor_alias: "",
  total: 210.5,
  settled_total: 0,
  currency_code: "CAD",
};
const bodies = {
  "/api/time_sheets/tracking_counts": [
    {
      week_ending: weekEnding,
      submitted_count: 2,
      approved_count: 5,
      committed_count: 12,
      rejected_count: 1,
    },
  ],
  "/api/expenses/tracking_counts": [
    {
      committed_week_ending: weekEnding,
      committed_count: 8,
      placeholder_payroll_id_expense_count: 1,
    },
  ],
  [`/api/expenses/tracking/${weekEnding}`]: [expense],
};

const screens = [
  { url: "/time/tracking", title: "Download the time of a week", shot: "time-tracking-weeks" },
  {
    url: "/expenses/tracking",
    title: "Download the expenses of a week",
    shot: "expense-tracking-weeks",
  },
  {
    url: `/expenses/tracking/${weekEnding}`,
    title: "Track committed expenses",
    shot: "expense-tracking-week",
  },
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
    await page.route("**/api/**", (route) => {
      const path = new URL(route.request().url()).pathname;
      const body =
        bodies[path] ??
        (path.endsWith("/auth-refresh")
          ? { token, record }
          : path.endsWith("/records")
            ? { page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] }
            : {});
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        headers: { "access-control-allow-origin": "*" },
        body: JSON.stringify(body),
      });
    });

    const shots = process.env.SCREENSHOT_DIR;
    for (const width of [1000, 375]) {
      await page.setViewportSize({ width, height: width === 375 ? 812 : 900 });
      for (const { url, title, shot } of screens) {
        await page.goto(`${origin}${url}`);
        const header = page.locator("li.inlistheader").first();
        // Wait for a row, so the list has its final layout.
        await page.locator("li.contents").first().waitFor();
        const help = header.getByRole("button", { name: /^Help/ });
        const [headerBox, helpBox] = [await header.boundingBox(), await help.boundingBox()];
        assert.ok(
          headerBox.x + headerBox.width - (helpBox.x + helpBox.width) < 20,
          `${url} at ${width}px: Help is at the right of the header`,
        );
        assert.ok(helpBox.y >= headerBox.y && helpBox.y < headerBox.y + headerBox.height);
        // The page does not scroll sideways on a phone.
        assert.ok(
          await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth),
          `${url} at ${width}px: no sideways scroll`,
        );
        if (shots) {
          await page.screenshot({
            path: `${shots}/${shot}-${width === 375 ? "375" : "desktop"}.png`,
            fullPage: true,
          });
        }
        await help.click();
        const panel = page.getByRole("dialog", { name: title, exact: true });
        await panel.getByRole("link", { name: "Open the full page", exact: true }).waitFor();
        if (shots) {
          await page.screenshot({
            path: `${shots}/${shot}-${width === 375 ? "375" : "desktop"}-open.png`,
          });
        }
        await page.keyboard.press("Escape");
        await panel.waitFor({ state: "hidden" });
      }
    }

    assert.deepEqual(errors, []);
    console.log(
      "Tracking help checks passed: the time and expense week lists and the expense week page show Help at the right of the list header, at desktop and phone widths, and open the right section.",
    );
  },
  { app: true },
);
