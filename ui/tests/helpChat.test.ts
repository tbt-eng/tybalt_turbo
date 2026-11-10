import test from "node:test";
import assert from "node:assert/strict";
import {
  answerText,
  renderAnswer,
  sourceLinks,
  type ChatSource,
} from "../src/lib/helpChatFormat.ts";
import {
  createEventParser,
  HelpRequestError,
  helpUnavailable,
  streamHelp,
  type HelpEvent,
} from "../src/lib/helpChatStream.ts";

const sources: ChatSource[] = [
  {
    kind: "chunk",
    ref: "S2",
    id: "app/hooks/time_entries.go:24",
    path: "app/hooks/time_entries.go",
    start: 24,
    end: 53,
    type: "code",
  },
  { kind: "doc", url: "/docs/time#submit-a-week" },
  {
    kind: "chunk",
    ref: "S5",
    id: "descriptions/time_entries.md:1",
    path: "descriptions/time_entries.md",
    start: 1,
    end: 9,
    heading: "Time entries > Rules",
    type: "notes",
  },
];

test("an answer renders Markdown, and links only to the user docs", () => {
  const html = renderAnswer(
    "Select **Submit**. See [Time](/docs/time#submit-a-week).\n\n1. Open **Entries**.\n2. Select `Save`.\n",
    null,
  );
  assert.match(html, /<strong>Submit<\/strong>/);
  assert.match(html, /<a class="link" href="\/docs\/time#submit-a-week">Time<\/a>/);
  assert.match(html, /<ol>\s*<li>Open <strong>Entries<\/strong>\.<\/li>/);
  assert.match(html, /<code>Save<\/code>/);
});

test("text from the model can't become markup or a link out of Turbo", () => {
  const attacks = [
    `<script>alert(1)</script>`,
    `<img src=x onerror="alert(1)">`,
    `Hello <b onmouseover="alert(1)">there</b>`,
    `[Click](javascript:alert(1))`,
    `[Click](https://evil.example/docs/time)`,
    `[Click](//evil.example)`,
    `[Click](/docs/time" onclick="alert(1))`,
    `[Click](/expenses/list)`,
    `![alt text](https://evil.example/pixel.png)`,
    `<https://evil.example>`,
    `https://evil.example/auto`,
    `<a href="/docs/time">raw</a>`,
    `<iframe src="https://evil.example"></iframe>`,
    "```html\n<script>alert(1)</script>\n```",
    `[x]: javascript:alert(1)\n\n[x]`,
  ];
  const allowedTags = new Set([
    "p",
    "strong",
    "em",
    "ol",
    "ul",
    "li",
    "code",
    "pre",
    "a",
    "sup",
    "br",
  ]);
  for (const attack of attacks) {
    const html = renderAnswer(attack, sources);
    // Every real tag in the output is one that Markdown makes, and the only
    // attributes are a docs link's and a code block's language.
    for (const [tag, name, attributes] of html.matchAll(/<\/?([a-zA-Z][\w-]*)([^>]*)>/g)) {
      assert.ok(allowedTags.has(name.toLowerCase()), `${attack} gave ${tag}`);
      assert.match(
        attributes,
        /^(| class="link" href="\/docs\/[a-z0-9-]+(#[a-z0-9-]+)?"| class="language-[a-z]+")$/,
        `${attack} gave ${tag}`,
      );
    }
  }
  // The text is still there to read.
  assert.match(
    renderAnswer(`<script>alert(1)</script>`, null),
    /&lt;script&gt;alert\(1\)&lt;\/script&gt;/,
  );
  assert.match(renderAnswer(`[Click](https://evil.example)`, null), />Click</);
});

test("chunk marks become source numbers when the server has sent the sources", () => {
  const text = "It clears the field [S2]. The rule is here [S5, S2][S9]. No source [S1].";
  // While the answer streams, the marks are left out.
  assert.equal(
    renderAnswer(text, null),
    "<p>It clears the field . The rule is here . No source .</p>\n",
  );
  assert.equal(
    renderAnswer(text, sources),
    "<p>It clears the field <sup>1</sup>. The rule is here <sup>3</sup><sup>1</sup>. No source .</p>\n",
  );
  assert.equal(answerText(text), "It clears the field. The rule is here. No source.");
});

test("sources become links: docs pages in Turbo, and chunks on GitHub", () => {
  const links = sourceLinks(sources, { time: "Time" });
  assert.deepEqual(links, [
    {
      number: 1,
      key: "S2",
      label: "Code: app/hooks/time_entries.go",
      href: "https://github.com/tbt-eng/tybalt_turbo/blob/main/app/hooks/time_entries.go#L24-L53",
      external: true,
    },
    {
      number: 2,
      key: "/docs/time#submit-a-week",
      label: "Time",
      href: "/docs/time#submit-a-week",
      external: false,
    },
    {
      number: 3,
      key: "S5",
      label: "Design note: Time entries > Rules",
      href: "https://github.com/tbt-eng/tybalt_turbo/blob/main/descriptions/time_entries.md#L1-L9",
      external: true,
    },
  ]);
  // A route folder has brackets in its name.
  const [route] = sourceLinks(
    [
      {
        kind: "chunk",
        ref: "S1",
        path: "ui/src/routes/docs/[slug]/+page.svelte",
        start: 1,
        end: 6,
        type: "code",
      },
    ],
    {},
  );
  assert.equal(
    route.href,
    "https://github.com/tbt-eng/tybalt_turbo/blob/main/ui/src/routes/docs/%5Bslug%5D/%2Bpage.svelte#L1-L6",
  );
  // A source that is not what the server sends is left out.
  assert.deepEqual(
    sourceLinks(
      [
        { kind: "doc", url: "https://evil.example" },
        { kind: "doc", url: "/expenses/list" },
        { kind: "chunk", ref: "S1", path: "../../etc/passwd?x=<script>" },
        { kind: "chunk", path: "app/main.go" },
        { kind: "doc" },
      ],
      {},
    ),
    [],
  );
  assert.equal(sourceLinks([{ kind: "doc", url: "/docs/glossary" }], {})[0].label, "glossary");
});

test("the event parser joins events that arrive in pieces", () => {
  const events: HelpEvent[] = [];
  const parse = createEventParser((event) => events.push(event));
  parse('event: meta\ndata: {"chat":"abc","turn":"t1"}\n\nevent: te');
  assert.equal(events.length, 1);
  parse('xt\ndata: {"delta":"Sel');
  parse('ect "}\n\n');
  parse(
    'event: text\r\ndata: {"delta":"Submit"}\r\n\r\n: a comment\n\nevent: bad\ndata: not json\n\n',
  );
  parse('event: done\ndata: {"status":"answered","message":""}\n\n');
  assert.deepEqual(events, [
    { event: "meta", data: { chat: "abc", turn: "t1" } },
    { event: "text", data: { delta: "Select " } },
    { event: "text", data: { delta: "Submit" } },
    { event: "done", data: { status: "answered", message: "" } },
  ]);
});

function streamResponse(pieces: string[]): Response {
  const encoder = new TextEncoder();
  const body = new ReadableStream<Uint8Array>({
    start(controller) {
      for (const piece of pieces) controller.enqueue(encoder.encode(piece));
      controller.close();
    },
  });
  return new Response(body, { status: 200, headers: { "Content-Type": "text/event-stream" } });
}

test("streamHelp posts the question with the token and reads every event", async () => {
  const calls: { url: string; init: RequestInit }[] = [];
  const request = (async (url: string, init: RequestInit) => {
    calls.push({ url, init });
    // The last event has no blank line after it, and a character is split.
    return streamResponse([
      'event: text\ndata: {"delta":"caf',
      'é"}\n\nevent: done\ndata: {"status":"answered"}',
    ]);
  }) as unknown as typeof fetch;
  const events: HelpEvent[] = [];
  await streamHelp(
    request,
    "http://x/api/help/ask",
    "the-token",
    { question: "Q?" },
    new AbortController().signal,
    (event) => events.push(event),
  );
  assert.deepEqual(events, [
    { event: "text", data: { delta: "café" } },
    { event: "done", data: { status: "answered" } },
  ]);
  assert.equal(calls[0].url, "http://x/api/help/ask");
  assert.equal(calls[0].init.method, "POST");
  assert.deepEqual(calls[0].init.headers, {
    Authorization: "the-token",
    "Content-Type": "application/json",
  });
  assert.equal(calls[0].init.body, '{"question":"Q?"}');
});

test("a refusal before the answer is an error with the server's message", async () => {
  const refuse = (status: number, body: string) =>
    (async () => new Response(body, { status })) as unknown as typeof fetch;
  const ask = (request: typeof fetch) =>
    streamHelp(request, "http://x", "t", {}, new AbortController().signal, () => {});
  await assert.rejects(
    ask(refuse(429, '{"message":"Help is still answering your last question."}')),
    (error) => {
      assert.ok(error instanceof HelpRequestError);
      assert.equal(error.status, 429);
      assert.equal(error.message, "Help is still answering your last question.");
      return true;
    },
  );
  await assert.rejects(ask(refuse(502, "<html>Bad Gateway</html>")), (error) => {
    assert.ok(error instanceof HelpRequestError);
    assert.equal(error.message, helpUnavailable);
    return true;
  });
});
