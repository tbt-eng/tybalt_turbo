// The App Config page in a real browser, with the server stubbed. Set
// SCREENSHOT_DIR to keep a picture of each state.
import assert from "node:assert/strict";
import path from "node:path";
import { runBrowserHarness } from "./browserHarness.mjs";

const record = {
  id: "adminuser",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "admin@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");

const configRows = [
  {
    id: "helpchatconfig1",
    key: "help_chat",
    value: { enabled: true, deploy_group: "group0000000001" },
    description: "Help chat",
    updated: "2026-11-11 14:00:00.000Z",
  },
  {
    id: "expensesconfig1",
    key: "expenses",
    value: {
      create_edit_absorb: true,
      no_po_expense_limit: 100,
      po_expense_allowed_excess: { mode: "lesser_of", percent: 5, value: 100 },
    },
    description: "",
    updated: "2026-10-01 09:30:00.000Z",
  },
];

const goodStatus = {
  settings_error: "",
  has_row: true,
  settings: {
    enabled: true,
    deploy_group: "group0000000001",
    provider: "anthropic",
    model: "claude-haiku-5-5",
    effort: "low",
    max_output_tokens: 8000,
    max_question_chars: 2000,
    questions_per_hour: 30,
    max_tool_calls: 6,
    history_turns: 10,
    turn_timeout_seconds: 60,
  },
  keys: [
    { provider: "anthropic", secret: "ANTHROPIC_API_KEY", set: true, in_use: true },
    { provider: "openai", secret: "OPENAI_API_KEY", set: false, in_use: false },
    { provider: "openrouter", secret: "OPENROUTER_API_KEY", set: false, in_use: false },
  ],
  deploy_group: { id: "group0000000001", name: "help chat testers", found: true, members: 4 },
  deploy_groups: [
    { id: "group0000000001", name: "help chat testers", found: true, members: 4 },
    { id: "group0000000002", name: "invoice pilot", found: true, members: 1 },
  ],
  corpus: { state: "loaded", chunks: 6880, hash: "0123456789abcdef0123456789abcdef", error: "" },
  last_day: {
    questions: 12,
    users: 3,
    errors: 1,
    tokens_in: 4000,
    tokens_cache_read: 90000,
    tokens_cache_write: 6000,
    tokens_out: 5200,
  },
  last_answered: "2026-11-11 15:20:00.000Z",
  last_error: { at: "2026-11-11 13:05:00.000Z", error: "the answer took too long" },
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

    let claims = ["admin"];
    let status = goodStatus;
    let keyTest = {
      ok: true,
      provider: "anthropic",
      model: "claude-haiku-5-5",
      message: "The key works and the provider knows the model.",
    };
    const calls = { status: 0, test: 0 };
    await page.route("**/api/**", async (route) => {
      const request = route.request();
      const url = new URL(request.url());
      const json = (body, code = 200) =>
        route.fulfill({
          status: code,
          contentType: "application/json",
          headers: { "access-control-allow-origin": "*" },
          body: JSON.stringify(body),
        });
      const list = (items) =>
        json({ page: 1, perPage: 500, totalItems: items.length, totalPages: 1, items });
      if (url.pathname === "/api/help/status") {
        calls.status += 1;
        assert.equal(request.method(), "GET");
        return json(status);
      }
      if (url.pathname === "/api/help/status/test_key") {
        calls.test += 1;
        assert.equal(request.method(), "POST");
        return json(keyTest);
      }
      if (url.pathname.endsWith("/auth-refresh")) return json({ token, record });
      if (url.pathname === `/api/collections/user_claims_summary/records/${record.id}`) {
        return json({ id: record.id, claims });
      }
      if (url.pathname === "/api/collections/app_config/records") return list(configRows);
      if (url.pathname.endsWith("/records")) return list([]);
      return json({});
    });

    const shot = async (name) => {
      if (!process.env.SCREENSHOT_DIR) return;
      await page.screenshot({
        path: path.join(process.env.SCREENSHOT_DIR, `${name}.png`),
        fullPage: true,
      });
    };
    const text = () => page.locator("main, body").first().innerText();

    // An admin sees the status and every setting.
    await page.goto(`${origin}/app_config`);
    await page.getByRole("heading", { level: 1, name: "App Config" }).waitFor();
    await page.getByText("Help chat is on", { exact: true }).waitFor();
    let body = await text();
    assert.match(body, /the 4 members of help chat testers/);
    assert.match(body, /ANTHROPIC_API_KEY/);
    assert.match(body, /claude-haiku-5-5/);
    assert.match(body, /90%/);
    assert.match(body, /the answer took too long/);
    assert.match(body, /invoice pilot/);
    // Every row of app_config, with each property.
    assert.match(body, /no_po_expense_limit/);
    assert.match(body, /"mode":"lesser_of"/);
    assert.match(body, /deploy_group/);
    assert.equal(await page.getByText("No settings", { exact: true }).count(), 0);
    assert.equal(calls.status, 1);
    // Opening the page does not test the key.
    assert.equal(calls.test, 0);
    await shot("on");
    await page.getByRole("heading", { level: 2, name: "All settings" }).scrollIntoViewIfNeeded();
    await shot("all-settings");

    // The key test runs on a press, and shows the provider's answer.
    await page.getByRole("button", { name: "Test the key" }).click();
    await page.getByText("The key test passed", { exact: true }).waitFor();
    assert.equal(calls.test, 1);
    keyTest = {
      ok: false,
      provider: "anthropic",
      model: "claude-haiku-5-5",
      message: "The provider refused the key in ANTHROPIC_API_KEY (HTTP 401). invalid x-api-key",
    };
    await page.getByRole("button", { name: "Test the key" }).click();
    await page.getByText("The key test failed", { exact: true }).waitFor();
    assert.match(await text(), /refused the key/);
    await shot("key-test-failed");

    // A bad setting and a missing key show after a refresh.
    status = {
      ...goodStatus,
      settings_error: 'app_config help_chat: json: unknown field "deploy_grp"',
      deploy_group: null,
    };
    await page.getByRole("button", { name: "Refresh" }).click();
    await page.getByText("The help_chat setting is not valid", { exact: true }).waitFor();
    assert.match(await text(), /deploy_grp/);
    status = {
      ...goodStatus,
      keys: goodStatus.keys.map((key) => ({ ...key, set: false })),
    };
    await page.getByRole("button", { name: "Refresh" }).click();
    await page.getByText("The server has no key for the provider", { exact: true }).waitFor();
    assert.equal(await page.getByText("Help chat is on", { exact: true }).count(), 0);
    await shot("no-key");

    // At phone width nothing is wider than the screen.
    await page.setViewportSize({ width: 390, height: 844 });
    await page.waitForTimeout(100);
    const overflow = await page.evaluate(
      () => document.documentElement.scrollWidth - document.documentElement.clientWidth,
    );
    assert.ok(overflow <= 0, `the page is ${overflow}px wider than the phone`);
    await shot("phone");
    await page.setViewportSize({ width: 1000, height: 900 });

    // A user with no admin claim sees no status, and the page asks for none.
    claims = [];
    const before = calls.status;
    await page.goto(`${origin}/app_config`);
    await page.getByText("Admins only", { exact: true }).waitFor();
    await page.waitForTimeout(300);
    assert.equal(calls.status, before);
    assert.equal(await page.getByText("ANTHROPIC_API_KEY").count(), 0);
    assert.equal(await page.getByText("no_po_expense_limit").count(), 0);

    assert.deepEqual(errors, []);
    console.log(
      "App Config checks passed: status, settings, key test, bad setting, missing key, phone width and admins only.",
    );
  },
  { app: true },
);
