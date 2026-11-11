import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// A Cumulative PO's remaining amount and an overflow are in the PO's currency.
const record = {
  id: "poexpuser",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "poexp@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");

const currencies = [
  { id: "curcad", code: "CAD", symbol: "$", icon: "", rate: 1, rate_date: "", ui_sort: 0 },
  {
    id: "curusd",
    code: "USD",
    symbol: "US$",
    icon: "",
    rate: 1.37,
    rate_date: "2026-11-09",
    ui_sort: 1,
  },
];
const po = (id, fields) => ({
  id,
  uid: record.id,
  po_number: `2611-${id}`,
  status: "Active",
  legacy_manual_entry: false,
  type: "Cumulative",
  division: "",
  description: "Site supplies",
  payment_type: "OnAccount",
  vendor: "",
  job: "",
  category: "",
  kind: "",
  recurring_expected_occurrences: 0,
  recurring_remaining_occurrences: 0,
  ...fields,
});
const pos = {
  pocad: po("pocad", { currency: "curcad", currency_code: "CAD", remaining_amount: 1234.5 }),
  pousd: po("pousd", { currency: "curusd", currency_code: "USD", remaining_amount: 400 }),
};
// An expense being edited carries its PO only as an expanded record, whose
// currency is an id, so the editor looks the code up.
const editedExpense = {
  id: "expusd",
  collectionName: "expenses",
  collectionId: "expenses",
  uid: record.id,
  creator: record.id,
  date: "2026-11-09",
  division: "",
  description: "Site supplies",
  payment_type: "OnAccount",
  purchase_order: "pousd",
  currency: "curusd",
  vendor: "",
  job: "",
  category: "",
  kind: "",
  total: 100,
  settled_total: 0,
  expand: { purchase_order: { ...pos.pousd, collectionName: "purchase_orders" } },
};
const overflows = {
  pocad: { po_total: 5000, overflow_amount: 1234.5 },
  pousd: { po_total: 2000, overflow_amount: 25.75 },
};
const list = (items) => ({ page: 1, perPage: 500, totalItems: items.length, totalPages: 1, items });
const overflowError = (poId) => ({
  status: 400,
  message: "cumulative expenses exceed purchase order total",
  data: {
    total: {
      code: "cumulative_po_overflow",
      message: "cumulative expenses exceed purchase order total",
      data: { purchase_order: poId, po_number: pos[poId].po_number, ...overflows[poId] },
    },
  },
});

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
    let overflowPO = "";
    await page.route("**/api/**", (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const reply = (body, status = 200) =>
        route.fulfill({
          status,
          contentType: "application/json",
          headers: { "access-control-allow-origin": "*" },
          body: JSON.stringify(body),
        });
      const visible = path.match(/^\/api\/purchase_orders\/visible\/(\w+)$/);
      if (visible) return reply(pos[visible[1]]);
      if (request.method() === "POST" && path === "/api/expenses") {
        return reply(overflowError(overflowPO), 400);
      }
      if (request.method() === "PATCH" && path === `/api/expenses/${editedExpense.id}`) {
        return reply(overflowError("pousd"), 400);
      }
      if (path === `/api/collections/expenses/records/${editedExpense.id}`) {
        return reply(editedExpense);
      }
      if (path === "/api/collections/currencies/records") return reply(list(currencies));
      if (path.endsWith("/auth-refresh")) return reply({ token, record });
      if (path.endsWith("/records")) return reply(list([]));
      return reply({});
    });

    const text = async (locator) => {
      await locator.waitFor();
      return (await locator.innerText()).replace(/\s+/g, " ");
    };
    const remaining = () => text(page.getByText("Provisional Remaining:"));
    const overflowMessage = async () => {
      await page.getByRole("button", { name: "Save" }).click();
      const modalText = page.getByText("exceeds total of purchase order");
      const message = await text(modalText);
      // The modal comes before the form, whose own Cancel button is second.
      await page.getByRole("button", { name: "Cancel" }).first().click();
      await modalText.waitFor({ state: "hidden" });
      return message;
    };

    overflowPO = "pocad";
    await page.goto(`${origin}/expenses/add/pocad`);
    // The first page waits for Vite to prepare the app.
    await page.getByText("Provisional Remaining:").waitFor({ timeout: 60000 });
    assert.equal(
      await remaining(),
      "Provisional Remaining: $1,234.50 (includes uncommitted expenses)",
    );
    assert.match(await overflowMessage(), /purchase order 2611-pocad by \$1,234\.50\./);

    overflowPO = "pousd";
    await page.goto(`${origin}/expenses/add/pousd`);
    assert.equal(
      await remaining(),
      "Provisional Remaining: US$400.00 (includes uncommitted expenses)",
    );
    assert.match(await overflowMessage(), /purchase order 2611-pousd by US\$25\.75\./);

    await page.goto(`${origin}/expenses/${editedExpense.id}/edit`);
    await page.getByText("PO 2611-pousd").waitFor();
    assert.match(await overflowMessage(), /purchase order 2611-pousd by US\$25\.75\./);

    if (process.env.SCREENSHOT) await page.screenshot({ path: process.env.SCREENSHOT });
    assert.deepEqual(errors, []);
    console.log(
      "Cumulative PO: remaining and overflow amounts in the PO currency (CAD, USD, edit)",
    );
  },
  { app: true },
);
