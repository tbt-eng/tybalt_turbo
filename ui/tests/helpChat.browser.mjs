import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

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

// The answer holds Markdown, a docs link, a chunk mark and an attack.
const answer = [
  "Open **Entries** and select **Submit** on the week [S1]. ",
  'See [Time](/docs/time#submit-a-week). <img src="x" onerror="window.__attack = 1">',
];
const sources = [
  {
    kind: "chunk",
    ref: "S1",
    id: "app/hooks/time_entries.go:24",
    path: "app/hooks/time_entries.go",
    start: 24,
    end: 53,
    type: "code",
  },
  { kind: "doc", url: "/docs/time#submit-a-week" },
];
const answered = (turn) =>
  sse([
    ["meta", { chat: "chat00000000001", turn }],
    ["status", { state: "searching", detail: "submit week" }],
    ["text", { delta: answer[0] }],
    ["text", { delta: answer[1] }],
    ["sources", sources],
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

    // The server, as the browser sees it. Each test step sets what ask returns.
    let chatEnabled = true;
    let askReply = () => ({
      status: 200,
      contentType: "text/event-stream",
      body: answered("turn0000000001"),
    });
    const asks = [];
    const feedback = [];
    let warmed = 0;
    await page.route("**/api/**", async (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const json = (body, status = 200) =>
        route.fulfill({
          status,
          contentType: "application/json",
          headers: { "access-control-allow-origin": "*" },
          body: JSON.stringify(body),
        });
      if (path === "/api/help/ask") {
        asks.push({ body: request.postDataJSON(), authorization: request.headers().authorization });
        const reply = await askReply();
        return route
          .fulfill({ headers: { "access-control-allow-origin": "*" }, ...reply })
          .catch(() => {});
      }
      if (path.startsWith("/api/help/turns/")) {
        feedback.push({ path, body: request.postDataJSON() });
        return route.fulfill({ status: 204, headers: { "access-control-allow-origin": "*" } });
      }
      if (path === "/api/help/warm") {
        warmed += 1;
        return route.fulfill({ status: 204, headers: { "access-control-allow-origin": "*" } });
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

    const search = page.getByRole("textbox", { name: /Search Help/ });
    const composer = page.getByRole("textbox", { name: "Ask a follow-up" });
    const docsTab = page.locator('a[href="#docs"]');
    const askTab = page.locator('a[href="#ask"]');
    const actionRow = page.locator("li[data-action]");
    // The hits are links in the list under the search bar, not in the side menu.
    const hits = page
      .getByRole("search")
      .locator("xpath=following-sibling::div")
      .locator("ul > li > a");
    const turns = page.locator("section[data-turn]");
    const pageTitle = (name) => page.getByRole("heading", { level: 1, name, exact: true });
    const focused = () => page.evaluate(() => document.activeElement?.getAttribute("aria-label"));
    const waitFor = async (check, what) => {
      const deadline = Date.now() + 5000;
      while (Date.now() < deadline) {
        if (await check()) return;
        await page.waitForTimeout(50);
      }
      assert.fail(`timed out: ${what}`);
    };
    // How far a turn's top is below the top of the conversation's view.
    const turnOffset = (index) =>
      turns.nth(index).evaluate((section) => {
        const scroller = section.closest(".overflow-y-auto");
        return section.getBoundingClientRect().top - scroller.getBoundingClientRect().top;
      });

    await page.goto(`${origin}/docs`);
    await pageTitle("Start here").waitFor();
    await page.getByRole("textbox", { name: "Search Help, or ask a question" }).waitFor();
    await waitFor(() => warmed > 0, "the page asks the server to get the model ready");
    // Before the first question, Help is the docs with a search bar, and no tabs.
    assert.equal(await askTab.count(), 0);

    // Typing shows the Ask row first, with the hits under it. Nothing is sent.
    // The pointer rests where the hits appear: a row that appears under a
    // pointer at rest does not become the chosen row.
    await page.mouse.move(600, 260);
    await search.fill("How do I submit my week?");
    await actionRow.waitFor();
    await page.waitForTimeout(300);
    assert.match(await actionRow.textContent(), /Ask Help: How do I submit my week\?/);
    assert.equal(await actionRow.getAttribute("data-active"), "");
    assert.ok((await hits.count()) > 0, "the docs hits still show");
    assert.equal(asks.length, 0);

    // The arrow keys and a pointer that moves reach a hit, and come back.
    await search.press("ArrowDown");
    assert.equal(await actionRow.getAttribute("data-active"), null);
    assert.equal(await hits.first().getAttribute("data-active"), "");
    await search.press("ArrowUp");
    assert.equal(await actionRow.getAttribute("data-active"), "");
    await hits.first().hover();
    assert.equal(await hits.first().getAttribute("data-active"), "");
    await actionRow.hover();
    assert.equal(await actionRow.getAttribute("data-active"), "");

    // Enter asks. The conversation shows in the Ask tab, with the box for the
    // next question at the bottom, ready to type in.
    await search.press("Enter");
    await turns.first().getByRole("heading", { name: "How do I submit my week?" }).waitFor();
    await turns.first().getByText("Sources").waitFor();
    assert.equal(await pageTitle("Start here").count(), 0);
    assert.equal(await search.count(), 0);
    assert.equal(await askTab.getAttribute("aria-current"), "page");
    assert.equal(await docsTab.getAttribute("aria-current"), null);
    assert.equal(await focused(), "Ask a follow-up");
    assert.deepEqual(asks[0].body, {
      question: "How do I submit my week?",
      chat: "",
      route: "/docs",
      surface: "help",
    });
    assert.equal(asks[0].authorization, token);
    await page.getByText("Help writes these answers with AI").waitFor();
    // The box is under the conversation, at the bottom of the screen.
    const composerBox = await composer.boundingBox();
    const turnBox = await turns.first().boundingBox();
    assert.ok(composerBox.y > turnBox.y && composerBox.y > 600, "the box is at the bottom");

    // The answer is Markdown. Text from the model does not become markup.
    const first = turns.first();
    assert.deepEqual(await first.locator("strong").allTextContents(), ["Entries", "Submit"]);
    assert.equal(await first.locator("img, script, iframe").count(), 0);
    assert.equal(await page.evaluate(() => window.__attack), undefined);
    assert.match(await first.textContent(), /<img src="x" onerror=/);
    assert.equal(await first.locator("sup").textContent(), "1");
    assert.doesNotMatch(await first.textContent(), /\[S1\]/);

    // The sources: the code on GitHub in a new tab, and the docs page in Turbo.
    const sourceLinks = first.locator("ol a");
    assert.deepEqual(await sourceLinks.allTextContents(), [
      "Code: app/hooks/time_entries.go",
      "Time",
    ]);
    assert.equal(
      await sourceLinks.first().getAttribute("href"),
      "https://github.com/tbt-eng/tybalt_turbo/blob/main/app/hooks/time_entries.go#L24-L53",
    );
    assert.equal(await sourceLinks.first().getAttribute("target"), "_blank");
    assert.equal(await sourceLinks.nth(1).getAttribute("href"), "/docs/time#submit-a-week");
    assert.equal(await sourceLinks.nth(1).getAttribute("target"), null);

    // Feedback: a rating, a note for a thumbs down, a change of mind, a copy.
    const turnPath = "/api/help/turns/turn0000000001/feedback";
    await first.getByRole("button", { name: "This answer did not help" }).click();
    await waitFor(() => feedback.length === 1, "the rating is sent");
    assert.deepEqual(feedback[0], { path: turnPath, body: { rating: "down" } });
    await first.getByRole("button", { name: "You marked this answer as not helpful" }).waitFor();
    await first
      .getByRole("textbox", { name: "What was wrong?" })
      .fill("It named the wrong button.");
    await first.getByRole("button", { name: "Send" }).click();
    await first.getByText("Thank you. Turbo saved your note.").waitFor();
    assert.deepEqual(feedback[1], { path: turnPath, body: { note: "It named the wrong button." } });
    await first.getByRole("button", { name: "This answer helped" }).click();
    await waitFor(() => feedback.length === 3, "the new rating is sent");
    assert.deepEqual(feedback[2].body, { rating: "up" });
    assert.equal(await first.getByRole("textbox", { name: "What was wrong?" }).count(), 0);
    await first.getByRole("button", { name: "Copy this answer" }).click();
    await waitFor(() => feedback.length === 4, "the copy is sent");
    assert.deepEqual(feedback[3].body, { copied: true });

    // The next question goes in the box at the bottom and continues the chat.
    // The new question moves to the top of the view.
    askReply = () => ({
      status: 200,
      contentType: "text/event-stream",
      body: answered("turn0000000002"),
    });
    assert.equal(await page.getByRole("button", { name: "Ask", exact: true }).isDisabled(), true);
    await composer.fill("And who approves it?");
    await composer.press("Enter");
    await turns.nth(1).getByText("Sources").waitFor();
    assert.equal(asks[1].body.chat, "chat00000000001");
    assert.equal(await turns.count(), 2);
    assert.equal(await composer.inputValue(), "");
    await waitFor(
      async () => Math.abs((await turnOffset(1)) - 16) < 4,
      "the new turn is at the top",
    );

    // A refusal from the server shows its message, with nothing to rate.
    askReply = () => ({
      status: 429,
      contentType: "application/json",
      body: JSON.stringify({
        message: "You have asked a lot of questions in the last hour. Try again later.",
      }),
    });
    await composer.fill("One more?");
    await page.getByRole("button", { name: "Ask", exact: true }).click();
    await turns.nth(2).getByText("You have asked a lot of questions in the last hour.").waitFor();
    assert.equal(await turns.nth(2).getByRole("button").count(), 0);

    // A turn that fails on the server shows the message that it sends.
    askReply = () => ({
      status: 200,
      contentType: "text/event-stream",
      body: sse([
        ["meta", { chat: "chat00000000001", turn: "turn0000000003" }],
        ["done", { status: "error", message: "Help can't answer now. Try again in a minute." }],
      ]),
    });
    await composer.fill("Does it fail?");
    await composer.press("Enter");
    await turns.nth(3).getByText("Help can't answer now. Try again in a minute.").waitFor();

    // Stop ends an answer that is still on its way. No second question can
    // start while one runs.
    let release;
    askReply = () =>
      new Promise((resolve) => {
        release = () =>
          resolve({
            status: 200,
            contentType: "text/event-stream",
            body: answered("turn0000000004"),
          });
      });
    await composer.fill("A slow question?");
    await composer.press("Enter");
    await turns.nth(4).getByText("Thinking…").waitFor();
    const sent = asks.length;
    assert.equal(await page.getByRole("button", { name: "Ask", exact: true }).count(), 0);
    await composer.fill("Another one?");
    await composer.press("Enter");
    await page.waitForTimeout(100);
    assert.equal(asks.length, sent);
    assert.equal(await turns.count(), 5);
    await page.getByRole("button", { name: "Stop" }).click();
    await turns.nth(4).getByText("You stopped this answer.").waitFor();
    assert.equal(await page.getByRole("button", { name: "Stop" }).count(), 0);
    release();
    await composer.fill("");

    // The Docs tab shows the docs and the search bar. The conversation stays.
    await docsTab.click();
    await pageTitle("Start here").waitFor();
    assert.equal(await docsTab.getAttribute("aria-current"), "page");
    assert.equal(await turns.count(), 0);
    assert.equal(await composer.count(), 0);
    await search.waitFor();
    assert.equal(new URL(page.url()).pathname, "/docs");
    // A page from the side menu shows in the Docs tab.
    await page
      .locator('nav[aria-label="Help"]')
      .getByRole("link", { name: "Time", exact: true })
      .click();
    await pageTitle("Time").waitFor();
    assert.equal(await docsTab.getAttribute("aria-current"), "page");
    // The Ask tab brings the conversation back, with its box ready.
    await askTab.click();
    await turns.first().waitFor();
    assert.equal(await turns.count(), 5);
    assert.equal(await pageTitle("Time").count(), 0);
    assert.equal(await focused(), "Ask a follow-up");
    assert.equal(
      await first.getByRole("button", { name: "You marked this answer as helpful" }).count(),
      1,
    );

    // A phone keeps the conversation in the screen's width.
    await page.setViewportSize({ width: 375, height: 800 });
    assert.equal(
      await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth),
      true,
    );
    await page.setViewportSize({ width: 1100, height: 720 });

    // A docs link in the answer opens the page at its section, in the Docs tab.
    const feedbackBefore = feedback.length;
    await first.locator('a[href="/docs/time#submit-a-week"]').first().click();
    await page.waitForURL("**/docs/time#submit-a-week");
    await pageTitle("Time").waitFor();
    assert.equal(await turns.count(), 0);
    assert.equal(await docsTab.getAttribute("aria-current"), "page");
    await page.waitForFunction(() => {
      const section = document.querySelector("#submit-a-week");
      const scroller = section?.closest(".overflow-y-auto");
      if (!section || !scroller) return false;
      const offset = section.getBoundingClientRect().top - scroller.getBoundingClientRect().top;
      return offset >= 0 && offset < scroller.clientHeight;
    });
    await waitFor(() => feedback.length > feedbackBefore, "the opened source is sent");
    assert.deepEqual(feedback.at(-1), {
      path: turnPath,
      body: { opened: "/docs/time#submit-a-week" },
    });

    // A question from the search bar continues the same chat, in the Ask tab.
    askReply = () => ({
      status: 200,
      contentType: "text/event-stream",
      body: answered("turn0000000005"),
    });
    await search.fill("And from the search bar?");
    await actionRow.waitFor();
    await search.press("Enter");
    await turns.nth(5).getByText("Sources").waitFor();
    assert.deepEqual(asks.at(-1).body, {
      question: "And from the search bar?",
      chat: "chat00000000001",
      route: "/docs/time",
      surface: "help",
    });
    assert.equal(await askTab.getAttribute("aria-current"), "page");

    // End chat removes the conversation and the tabs. The docs show again, and
    // the next question starts a new chat.
    await page.getByRole("button", { name: "End chat" }).click();
    await pageTitle("Time").waitFor();
    assert.equal(await askTab.count(), 0);
    assert.equal(await docsTab.count(), 0);
    await search.fill("A new chat?");
    await actionRow.waitFor();
    await search.press("Enter");
    await turns.first().getByText("Sources").waitFor();
    assert.equal(await turns.count(), 1);
    assert.equal(asks.at(-1).body.chat, "");
    await page.getByRole("button", { name: "End chat" }).click();
    await pageTitle("Time").waitFor();

    // With the chat off, Help has its search and no Ask row, and asks nothing.
    chatEnabled = false;
    const asked = asks.length;
    await page.goto(`${origin}/docs`);
    await pageTitle("Start here").waitFor();
    await page.getByRole("textbox", { name: "Search Help", exact: true }).waitFor();
    await search.fill("How do I submit my week?");
    await hits.first().waitFor();
    assert.equal(await actionRow.count(), 0);
    await search.press("Enter");
    await page.waitForURL("**/docs/**");
    assert.equal(asks.length, asked);

    assert.deepEqual(errors, []);
    console.log(
      "Help chat checks passed: Ask row, bottom box, streamed answer, safe markup, sources, feedback, follow-up, errors, stop, tabs, links, phone and the off switch.",
    );
  },
  { app: true },
);
