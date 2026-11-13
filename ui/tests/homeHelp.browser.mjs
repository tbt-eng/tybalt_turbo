import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// The home page's Help box: the example question that types itself, Tab to use
// it, and Enter to ask it in Help. Set SCREENSHOT_DIR to keep pictures.
const record = {
  id: "docsuser",
  collectionName: "users",
  collectionId: "_pb_users_auth_",
  email: "docs@example.test",
  verified: true,
};
const token = [
  "eyJhbGciOiJIUzI1NiJ9",
  Buffer.from(JSON.stringify({ id: record.id, exp: 4102444800 })).toString("base64url"),
  "test-signature",
].join(".");

const sse = (events) =>
  events.map(([event, data]) => `event: ${event}\ndata: ${JSON.stringify(data)}\n\n`).join("");
const answered = sse([
  ["meta", { chat: "chat00000000001", turn: "turn0000000001" }],
  ["text", { delta: "Open **Entries** and select **Submit** on the week." }],
  ["sources", []],
  ["done", { status: "answered", message: "" }],
]);

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

    let chatEnabled = true;
    const asks = [];
    let warmed = 0;
    await page.route("**/api/**", async (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const cors = { "access-control-allow-origin": "*" };
      const json = (body) =>
        route.fulfill({
          contentType: "application/json",
          headers: cors,
          body: JSON.stringify(body),
        });
      if (path === "/api/help/ask") {
        asks.push(request.postDataJSON());
        return route.fulfill({ contentType: "text/event-stream", headers: cors, body: answered });
      }
      if (path === "/api/help/warm") {
        warmed += 1;
        return route.fulfill({ status: 204, headers: cors });
      }
      if (path.endsWith("/auth-refresh")) return json({ token, record });
      if (path === "/api/collections/app_config/records") {
        const items = chatEnabled
          ? [{ id: "helpchatconfig1", key: "help_chat", value: { enabled: true }, description: "" }]
          : [];
        return json({ page: 1, perPage: 500, totalItems: items.length, totalPages: 1, items });
      }
      if (path.endsWith("/records")) {
        return json({ page: 1, perPage: 500, totalItems: 0, totalPages: 1, items: [] });
      }
      return json({});
    });

    const box = page.getByRole("textbox", { name: "Ask Help a question" });
    const hint = page.locator("[aria-hidden='true'] > span.truncate");
    const useButton = page.getByRole("button", { name: "Tab" });
    const whatsNew = page.getByRole("link", { name: "What's New" });
    const welcome = page.getByText("Use the navigation menu to get started");
    const shot = async (name) => {
      if (process.env.SCREENSHOT_DIR)
        await page.screenshot({ path: `${process.env.SCREENSHOT_DIR}/${name}.png` });
    };

    await page.goto(origin);
    await box.waitFor();
    // The box replaces the line about the menu. It has no placeholder: the
    // example types itself after a moment.
    assert.equal(await welcome.count(), 0);
    assert.equal(await box.getAttribute("placeholder"), "");
    assert.equal(await hint.count(), 0);
    assert.equal(warmed, 0, "a visit that asks nothing does not warm the model");

    // The model gets ready when the person goes to the box.
    await box.focus();
    for (let i = 0; i < 40 && warmed === 0; i++) await page.waitForTimeout(50);
    assert.equal(warmed, 1, "the first focus warms the model");
    await box.press("Enter");
    assert.equal(asks.length, 0, "Enter in the empty box asks nothing");

    // The example types itself, one character after another, to a question.
    await hint.waitFor();
    const early = await hint.textContent();
    await useButton.waitFor();
    const question = await hint.textContent();
    assert.ok(early.length < question.length && question.startsWith(early), "it types");
    assert.match(question, /^(How|When|What) .+\?$/);
    await shot("home-hint");

    // Text from the person hides the example. An empty box shows it again.
    await box.fill("x");
    assert.equal(await hint.count(), 0);
    assert.equal(await useButton.count(), 0);
    await box.fill("");
    assert.equal(await hint.textContent(), question);

    // Enter still asks nothing: the example is not in the box.
    await box.press("Enter");
    assert.equal(asks.length, 0);

    // Tab puts the example in the box and the focus stays. The next Tab moves on.
    await box.press("Tab");
    assert.equal(await box.inputValue(), question);
    assert.equal(await box.evaluate((input) => input === document.activeElement), true);
    await box.press("Tab");
    assert.equal(await whatsNew.evaluate((link) => link === document.activeElement), true);
    assert.equal(warmed, 1);

    // The button does the same, for a phone.
    await box.fill("");
    await useButton.click();
    assert.equal(await box.inputValue(), question);

    // Enter asks in Help, on the Ask tab. The question is not in the URL.
    await box.press("Enter");
    await page.waitForURL("**/docs#ask");
    await page.locator("section[data-turn]").getByRole("heading", { name: question }).waitFor();
    await page.getByText("select Submit on the week").waitFor();
    assert.equal(await page.locator('a[href="#ask"]').getAttribute("aria-current"), "page");
    assert.deepEqual(asks, [{ question, chat: "", route: "/", surface: "help" }]);
    await shot("home-asked");

    // A person who asked for less motion gets the example with no typing.
    await page.emulateMedia({ reducedMotion: "reduce" });
    await page.goto(origin);
    await box.waitFor();
    await useButton.waitFor({ timeout: 1000 });
    await page.emulateMedia({ reducedMotion: "no-preference" });

    // With the chat off, the home page has its line about the menu and no box.
    chatEnabled = false;
    await page.goto(origin);
    await welcome.waitFor();
    assert.equal(await box.count(), 0);
    assert.equal(await whatsNew.getAttribute("href"), "/docs/whats-new");
    await shot("home-off");

    assert.deepEqual(errors, []);
    console.log(
      "Home Help checks passed: typed example, Tab and button to use it, Enter to ask in Help, warm on focus, less motion and the off switch.",
    );
  },
  { app: true },
);
