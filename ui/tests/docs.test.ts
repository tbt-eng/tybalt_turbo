import test from "node:test";
import assert from "node:assert/strict";
import { readdirSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { checkDocStyle, replacedWords } from "../scripts/check-docs-style.mjs";
import { needsDocsWarning } from "../scripts/check-docs-commit.mjs";
import {
  commitUrl,
  countWords,
  docNavEntry,
  findSection,
  orderDocPages,
  parseDoc,
  parseFlow,
  slugify,
  withNewTabLinks,
  type DocPage,
} from "../src/lib/docsFormat.ts";

const glossary = readFileSync(new URL("../src/docs/glossary.md", import.meta.url), "utf8");

const maxSectionWords = 60;
const maxPageWords = 600;
const maxFlowSteps = 5;

const uiRoot = fileURLToPath(new URL("../", import.meta.url));
const docsDir = `${uiRoot}src/docs/`;
const sources = readdirSync(docsDir)
  .filter((file) => file.endsWith(".md") && file !== "README.md")
  .map((file) => ({
    slug: file.replace(".md", ""),
    source: readFileSync(`${docsDir}${file}`, "utf8"),
  }));
const pages: DocPage[] = sources.map(({ slug, source }) => parseDoc(slug, source));

test("slugify makes the reference that DSHelp uses", () => {
  assert.equal(slugify("Change a project's value"), "change-a-projects-value");
});

test("countWords ignores Markdown punctuation", () => {
  assert.equal(countWords("1. Select **Approve request**. | --- |"), 3);
});

test("parseFlow reads steps and notes, and refuses a bad line", () => {
  const flow = parseFlow("You | Submit | Upload **one** PDF.\nNote: Rejected? Edit it.");
  assert.deepEqual(flow.steps, [
    { actor: "You", step: "Submit", detail: "Upload <strong>one</strong> PDF." },
  ]);
  assert.deepEqual(flow.notes, ["Rejected? Edit it."]);
  assert.throws(() => parseFlow("You | Submit"), /Actor \| Step \| Detail/);
});

test("parseDoc splits a page into intro and sections", () => {
  const page = parseDoc("x", "# Title\n\nIntro.\n\n## First step\n\nDo it.\n");
  assert.equal(page.title, "Title");
  assert.equal(page.intro.length, 1);
  assert.equal(page.sections[0].id, "first-step");
  assert.equal(page.sections[0].words, 2);
});

test("a What's new heading keeps its commit out of the title and reference", () => {
  const [section] = parseDoc(
    "x",
    "# T\n\n## October 25 to 31, 2026 (to 5c35b0a0)\n\n- Item.\n",
  ).sections;
  assert.equal(section.title, "October 25 to 31, 2026");
  assert.equal(section.id, "october-25-to-31-2026");
  assert.equal(section.commit, "5c35b0a0");
  assert.equal(commitUrl("5c35b0a0"), "https://github.com/tbt-eng/tybalt_turbo/commit/5c35b0a0");
  const [plain] = parseDoc("x", "# T\n\n## Create a request\n\nDo it.\n").sections;
  assert.equal(plain.title, "Create a request");
  assert.equal(plain.commit, undefined);
});

test("docs links render with the shared link style", () => {
  const page = parseDoc("x", "# T\n\n[Glossary](/docs/glossary)\n");
  const [block] = page.intro;
  assert.ok(block.kind === "html" && block.html.includes('<a class="link" href="/docs/glossary">'));
  const [popoverBlock] = withNewTabLinks(page.intro);
  assert.ok(
    popoverBlock.kind === "html" &&
      popoverBlock.html.includes(
        '<a class="link" target="_blank" rel="noopener" href="/docs/glossary">',
      ),
  );
});

for (const page of pages) {
  test(`${page.slug} stays brief`, () => {
    // What's new is limited by its 12 weeks, and popovers don't show its weeks.
    if (page.slug === "whats-new") return;
    assert.ok(page.words <= maxPageWords, `${page.words} words; the limit is ${maxPageWords}`);
    for (const section of page.sections) {
      assert.ok(
        section.words <= maxSectionWords,
        `"${section.title}" has ${section.words} words; the limit is ${maxSectionWords}`,
      );
    }
    for (const block of [page.intro, ...page.sections.map((s) => s.blocks)].flat()) {
      if (block.kind === "flow") assert.ok(block.steps.length <= maxFlowSteps);
    }
  });
}

test("section headings on a page are unique", () => {
  for (const page of pages) {
    const ids = page.sections.map((section) => section.id);
    assert.equal(new Set(ids).size, ids.length, page.slug);
  }
});

test("links between docs resolve", () => {
  for (const { slug, source } of sources) {
    for (const [, target] of source.matchAll(/\]\(\/docs\/([^)]+)\)/g)) {
      const ok = target.includes("#")
        ? findSection(pages, target) !== null
        : pages.some((page) => page.slug === target);
      assert.ok(ok, `${slug} links to /docs/${target}`);
    }
  }
});

test("DSHelp references resolve", () => {
  const files = readdirSync(`${uiRoot}src`, { recursive: true, encoding: "utf8" })
    .filter((file) => file.endsWith(".svelte"))
    .map((file) => `src/${file}`);
  let references = 0;
  for (const file of files) {
    const source = readFileSync(`${uiRoot}${file}`, "utf8");
    assert.doesNotMatch(source, /<DSHelp\b[^>]*?\bdoc=\{/, `${file} needs a plain text reference`);
    for (const [, ref] of source.matchAll(/<DSHelp\b[^>]*?\bdoc="([^"]+)"/g)) {
      references += 1;
      assert.ok(findSection(pages, ref), `${file} refers to ${ref}`);
    }
  }
  assert.ok(references > 0, "expected at least one DSHelp doc reference");
});

test("the style check finds what it should", () => {
  const messages = (source: string, replaced: string[] = []) =>
    checkDocStyle(source, replaced).map((problem) => problem.message);
  assert.deepEqual(messages("# Title\n\n1. Select **Save**.\n\nThis is fine.\n"), []);
  assert.match(messages("1. The approver selects it.\n")[0], /start this step with a verb/);
  assert.match(messages("1. Open it, then close it.\n")[0], /one instruction/);
  assert.match(messages("The PO is approved by a person.\n")[0], /active voice/);
  assert.match(messages(`1. Select ${"word ".repeat(21).trim()}.\n`)[0], /limit is 20/);
  assert.match(messages(`${"Word ".repeat(26).trim()}.\n`)[0], /limit is 25/);
  assert.match(messages("Check your timesheet.\n", ["timesheet"])[0], /glossary word/);
  assert.deepEqual(messages("```bash\nThe command is approved by nobody at all.\n```\n"), []);
});

test("the glossary lists its replaced words", () => {
  const replaced = replacedWords(glossary);
  for (const word of [
    "job request",
    "billing profile",
    "PA",
    "unbundle",
    "timesheet",
    "employee",
  ]) {
    assert.ok(replaced.includes(word), word);
  }
});

test("the commit check warns only about user-facing changes without docs", () => {
  const route = ["ui/src/routes/pos/add/+page.svelte"];
  assert.equal(needsDocsWarning(route, "Change a PO screen"), true);
  assert.equal(needsDocsWarning(["app/hooks/expenses.go"], "Change a rule"), true);
  assert.equal(needsDocsWarning(route, "Change a screen\n\nDocs: none\n"), false);
  assert.equal(needsDocsWarning([...route, "ui/src/docs/expenses.md"], "Change"), false);
  assert.equal(needsDocsWarning(["app/hooks/expenses_test.go"], "Add a test"), false);
  assert.equal(needsDocsWarning(["ui/src/routes/dev/components/+page.svelte"], "Sample"), false);
  assert.equal(needsDocsWarning(["README.md"], "Words"), false);
});

test("Help lists Start here first and What's new last", () => {
  const titles = orderDocPages(pages).map((page) => page.slug);
  assert.equal(titles[0], "start-here");
  assert.equal(titles.at(-1), "whats-new");
  const middle = orderDocPages(pages)
    .slice(1, -1)
    .filter((page) => !docNavEntry(page).indent);
  assert.deepEqual(
    middle.map((page) => page.title),
    middle.map((page) => page.title).sort((a, b) => a.localeCompare(b)),
  );
});

test("Help lists every page, with a child page after its parent", () => {
  const ordered = orderDocPages(pages);
  assert.equal(ordered.length, pages.length);
  const slugs = ordered.map((page) => page.slug);
  assert.equal(slugs.indexOf("revise-a-purchase-order"), slugs.indexOf("purchase-orders") + 1);
  assert.deepEqual(docNavEntry({ slug: "revise-a-purchase-order", title: "Revise" }), {
    slug: "revise-a-purchase-order",
    label: "Revisions",
    indent: true,
  });
  assert.deepEqual(docNavEntry({ slug: "time", title: "Time" }), {
    slug: "time",
    label: "Time",
    indent: false,
  });
});

test("What's new has complete weeks, newest first, within 12 weeks", () => {
  const source = sources.find(({ slug }) => slug === "whats-new")?.source ?? "";
  const dayMs = 24 * 60 * 60 * 1000;
  const months =
    "January February March April May June July August September October November December".split(
      " ",
    );
  const weeks = [...source.matchAll(/^## (.+)$/gm)].map(([, heading]) => heading);
  assert.ok(weeks.length > 0);
  const starts = weeks.map((heading) => {
    const match = heading.match(/^(\w+) (\d+) to (?:(\w+) )?(\d+), (\d{4}) \(to [0-9a-f]{8}\)$/);
    assert.ok(match, `"${heading}" needs a date range`);
    const [, month, day, endMonth, endDay, year] = match;
    const start = new Date(Date.UTC(Number(year), months.indexOf(month), Number(day)));
    const end = new Date(Date.UTC(Number(year), months.indexOf(endMonth ?? month), Number(endDay)));
    assert.equal(start.getUTCDay(), 0, `"${heading}" must start on a Sunday`);
    assert.equal(end.getUTCDay(), 6, `"${heading}" must end on a Saturday`);
    // The next day has started in UTC only after the Saturday ends.
    assert.ok(end.getTime() + dayMs <= Date.now(), `"${heading}" hasn't ended yet`);
    return start.getTime();
  });
  const week = 7 * dayMs;
  starts.forEach((start, index) => {
    if (index > 0) assert.ok(start < starts[index - 1], "newest week first");
  });
  assert.ok(starts[0] - starts.at(-1)! < 12 * week, "only the last 12 weeks");

  for (const block of source.split(/^## /m).slice(1)) {
    const items = block.split("\n").filter((line) => line.startsWith("- "));
    assert.ok(items.length > 0, `${block.split("\n")[0]} has no items`);
    for (const item of items) assert.match(item, /\(v\d+\.\d+\.\d+(?:, v\d+\.\d+\.\d+)*\)$/, item);
  }
});
