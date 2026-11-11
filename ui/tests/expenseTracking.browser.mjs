import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

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
const row = (id, fields) => ({
  id,
  uid: "u1",
  given_name: "Alex",
  surname: "Adams",
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
  total: 0,
  settled_total: 0,
  currency_code: "CAD",
  ...fields,
});
// The server already turns a missing currency into CAD.
const rows = [
  row("cadrow", { description: "Hotel in Ottawa", total: 210.5 }),
  row("usdsettled", {
    description: "Conference in Boston",
    total: 400,
    currency_code: "USD",
    settled_total: 548.12,
  }),
  row("eurunsettled", { description: "Rail pass", total: 90, currency_code: "EUR" }),
  row("mileagerow", {
    description: "Site visit",
    payment_type: "Mileage",
    distance: 120,
    total: 87.6,
  }),
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
        path === `/api/expenses/tracking/${weekEnding}`
          ? rows
          : path.endsWith("/auth-refresh")
            ? { token, record }
            : path.endsWith("/records")
              ? { page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] }
              : {};
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        headers: { "access-control-allow-origin": "*" },
        body: JSON.stringify(body),
      });
    });

    await page.goto(`${origin}/expenses/tracking/${weekEnding}`);
    // The description and the amounts share the row's headline line.
    const rowText = async (description) => {
      const item = page.locator(".headline_wrapper").filter({ hasText: description });
      await item.waitFor();
      return (await item.innerText()).replace(/\s+/g, " ");
    };

    const cad = await rowText("Hotel in Ottawa");
    assert.match(cad, /\$210\.50/);
    assert.doesNotMatch(cad, /Settled CAD/);

    const usd = await rowText("Conference in Boston");
    assert.match(usd, /US\$400\.00/);
    assert.match(usd, /Settled CAD: \$548\.12/);

    // A foreign expense with no settled total shows its own currency only.
    const eur = await rowText("Rail pass");
    assert.match(eur, /€90\.00/);
    assert.doesNotMatch(eur, /Settled CAD/);

    const mileage = await rowText("Site visit");
    assert.match(mileage, /120 km \/ \$87\.60/);
    assert.doesNotMatch(mileage, /Settled CAD/);

    if (process.env.SCREENSHOT) await page.screenshot({ path: process.env.SCREENSHOT });
    assert.deepEqual(errors, []);
    console.log(
      "Expense tracking week: CAD, foreign settled, foreign unsettled and mileage totals",
    );
  },
  { app: true },
);
