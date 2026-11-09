import test from "node:test";
import assert from "node:assert/strict";
import { readdirSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { parseDoc, type DocPage } from "../src/lib/docsFormat.ts";
import {
  createDocSearch,
  docSearchEntries,
  excerpt,
  glossarySynonyms,
  queryTerms,
} from "../src/lib/docsSearch.ts";

const docsDir = fileURLToPath(new URL("../src/docs/", import.meta.url));
const realPages: DocPage[] = readdirSync(docsDir)
  .filter((file) => file.endsWith(".md") && file !== "README.md")
  .map((file) => parseDoc(file.replace(".md", ""), readFileSync(`${docsDir}${file}`, "utf8")));

const pages = [
  parseDoc(
    "time",
    "# Time\n\nA time sheet holds your entries.\n\n```flow\nYou | Submit the week | Turbo creates a time sheet.\nNote: Rejected? Select **Recall**.\n```\n\n## Submit a week\n\n1. Open **Entries**.\n2. Select **Submit** on the week.\n\n## See your time off\n\n**Time Off** shows your available time off.\n",
  ),
  parseDoc(
    "expenses",
    "# Expenses\n\nAn expense records what you spent.\n\n## Submit an expense\n\nSelect **Submit**. Tom & Ann's manager approves it.\n",
  ),
  parseDoc(
    "glossary",
    "# Glossary\n\n| Word | Meaning | Don't use |\n| --- | --- | --- |\n| Time sheet | Your time entries for one week. | timesheet |\n| Project Request (PR) | A request. | job request, intake |\n| Staff | People. | |\n",
  ),
];
const search = createDocSearch(pages);

test("each intro and each section is one entry, as plain text", () => {
  const entries = docSearchEntries(pages);
  assert.deepEqual(
    entries.map((entry) => entry.id),
    [
      "time",
      "time#submit-a-week",
      "time#see-your-time-off",
      "expenses",
      "expenses#submit-an-expense",
      "glossary",
    ],
  );
  const intro = entries[0];
  assert.equal(intro.section, "");
  assert.equal(
    intro.text,
    "A time sheet holds your entries. You: Submit the week. Turbo creates a time sheet. Rejected? Select Recall.",
  );
  assert.match(entries[4].text, /Tom & Ann's manager/);
});

test("a search finds a section and marks the match", () => {
  const [hit] = search("time off");
  assert.equal(hit.id, "time#see-your-time-off");
  assert.equal(hit.page, "Time");
  assert.equal(hit.section, "See your time off");
  assert.deepEqual(
    hit.excerpt.filter((part) => part.match).map((part) => part.text),
    ["Time", "Off", "time", "off"],
  );
  assert.equal(
    hit.excerpt.map((part) => part.text).join(""),
    "Time Off shows your available time off.",
  );
});

test("a search matches the start of a word, and forgives a wrong letter", () => {
  assert.equal(search("expen")[0].slug, "expenses");
  assert.equal(search("expanse")[0].slug, "expenses");
});

test("a question finds the same hits as its key words", () => {
  const ids = (query: string) => search(query).map((hit) => hit.id);
  assert.deepEqual(ids("how do I submit my expense"), ids("submit expense"));
  assert.equal(ids("submit expense")[0], "expenses#submit-an-expense");
  assert.deepEqual(queryTerms("How do I…"), ["how", "do", "i"]);
  assert.deepEqual(queryTerms("the Orders and statuses of a business"), [
    "order",
    "statuse",
    "business",
  ]);
  assert.equal(search("time sheets entries")[0].id, "time");
});

test("a word that the glossary replaces finds the glossary word", () => {
  const synonyms = glossarySynonyms(pages);
  assert.deepEqual(queryTerms("my Timesheets", synonyms), ["time", "sheet"]);
  assert.deepEqual(queryTerms("job request", synonyms), ["project", "request"]);
  assert.deepEqual(queryTerms("intake", synonyms), ["project", "request"]);
  assert.ok(search("timesheet").some((hit) => hit.id === "time"));
});

test("a search with no match in every word falls back to any word", () => {
  assert.equal(search("expense zebra")[0].slug, "expenses");
  assert.deepEqual(search("zebra"), []);
  assert.deepEqual(search("  "), []);
});

test("an excerpt shows the text around the first match", () => {
  const text = `${"word ".repeat(40)}the needle is here ${"word ".repeat(40)}`.trim();
  const parts = excerpt(text, ["needle"]);
  const shown = parts.map((part) => part.text).join("");
  assert.ok(shown.startsWith("… word"));
  assert.ok(shown.endsWith("word …"));
  assert.ok(shown.length <= 164);
  assert.deepEqual(
    parts.filter((part) => part.match).map((part) => part.text),
    ["needle"],
  );
  assert.deepEqual(excerpt("Short text.", ["absent"]), [{ text: "Short text.", match: false }]);
});

test("every hit in the real docs points to a real page and section", () => {
  const realSearch = createDocSearch(realPages);
  for (const query of ["timesheet", "purchase order", "recall", "approve", "how do I submit"]) {
    const hits = realSearch(query);
    assert.ok(hits.length > 0 && hits.length <= 8, query);
    for (const hit of hits) {
      const page = realPages.find((candidate) => candidate.slug === hit.slug);
      assert.ok(page, hit.id);
      if (hit.sectionId)
        assert.ok(
          page.sections.some((s) => s.id === hit.sectionId),
          hit.id,
        );
    }
  }
  assert.equal(realSearch("timesheet")[0].slug, "time");
  assert.equal(realSearch("purchase orders")[0].id, "purchase-orders");
  assert.equal(realSearch("expense")[0].id, "expenses");
  assert.equal(realSearch("who approves")[0].id, "glossary#who-approves-what");
});
