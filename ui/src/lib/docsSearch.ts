// Searches the user docs in the browser, with no request to the server. This file
// has no Vite or Svelte imports, so tests can run it directly.
import MiniSearch from "minisearch";
import type { DocBlock, DocPage } from "$lib/docsFormat";

// One searchable piece of the docs: a page's intro or one of its sections.
export interface DocSearchEntry {
  id: string;
  slug: string;
  // Empty for a page's intro.
  sectionId: string;
  page: string;
  // Empty for a page's intro.
  section: string;
  text: string;
}

export interface ExcerptPart {
  text: string;
  match: boolean;
}

export interface DocSearchHit extends Omit<DocSearchEntry, "text"> {
  excerpt: ExcerptPart[];
}

const maxHits = 8;
const excerptLength = 160;
const excerptLead = 50;

// Words that a question holds but that don't say what it is about.
const fillerWords = new Set(
  "a an and are can do does how i in is it me my of on or the to what when where which who why you your".split(
    " ",
  ),
);

function htmlText(html: string): string {
  return html
    .replace(/<\/?(?:strong|em|code|del|a)\b[^>]*>/g, "")
    .replace(/<[^>]+>/g, " ")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&amp;/g, "&")
    .replace(/\s+/g, " ")
    .trim();
}

function blocksText(blocks: DocBlock[]): string {
  return blocks
    .map((block) =>
      block.kind === "html"
        ? htmlText(block.html)
        : [
            ...block.steps.map((step) => `${step.actor}: ${step.step}. ${htmlText(step.detail)}`),
            ...block.notes.map(htmlText),
          ].join(" "),
    )
    .join(" ");
}

export function docSearchEntries(pages: DocPage[]): DocSearchEntry[] {
  return pages.flatMap((page) => [
    {
      id: page.slug,
      slug: page.slug,
      sectionId: "",
      page: page.title,
      section: "",
      text: blocksText(page.intro),
    },
    ...page.sections.map((section) => ({
      id: `${page.slug}#${section.id}`,
      slug: page.slug,
      sectionId: section.id,
      page: page.title,
      section: section.title,
      text: blocksText(section.blocks),
    })),
  ]);
}

// The glossary lists words that the docs don't use, such as "timesheet", beside
// the word that they use instead. A search for the first finds the second.
export function glossarySynonyms(pages: DocPage[]): [RegExp, string][] {
  const glossary = pages.find((page) => page.slug === "glossary");
  const html = [glossary?.intro ?? [], ...(glossary?.sections ?? []).map((s) => s.blocks)]
    .flat()
    .map((block) => (block.kind === "html" ? block.html : ""))
    .join("");
  const rows = html.matchAll(
    /<tr>\s*<td[^>]*>(.*?)<\/td>\s*<td[^>]*>.*?<\/td>\s*<td[^>]*>(.*?)<\/td>\s*<\/tr>/gs,
  );
  return [...rows].flatMap(([, word, replaced]) => {
    const use = htmlText(word)
      .replace(/\s*\(.*\)$/, "")
      .toLowerCase();
    return htmlText(replaced)
      .split(",")
      .map((value) => value.trim().toLowerCase())
      .filter(Boolean)
      .map((value): [RegExp, string] => [
        new RegExp(`\\b${value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}s?\\b`, "g"),
        use,
      ]);
  });
}

// "orders" and "order" are one word to the search.
function singular(term: string): string {
  return term.length > 3 && term.endsWith("s") && !term.endsWith("ss") ? term.slice(0, -1) : term;
}

// The query in lower case, with the glossary's words in place of the words that
// the docs don't use.
function normalize(query: string, synonyms: [RegExp, string][]): string {
  return synonyms
    .reduce((value, [replaced, use]) => value.replace(replaced, use), query.toLowerCase())
    .replace(/\s+/g, " ")
    .trim();
}

export function queryTerms(query: string, synonyms: [RegExp, string][] = []): string[] {
  const terms = normalize(query, synonyms).match(/[a-z0-9]+/g) ?? [];
  const useful = terms.filter((term) => !fillerWords.has(term));
  return (useful.length ? useful : terms).map(singular);
}

// The part of the text around the first match, with each match marked.
export function excerpt(text: string, terms: string[]): ExcerptPart[] {
  const pattern = terms.length ? new RegExp(`\\b(?:${terms.join("|")})[a-z0-9]*`, "gi") : /(?!)/g;
  const first = text.search(pattern);
  let start = Math.max(0, (first < 0 ? 0 : first) - excerptLead);
  if (start > 0) start = text.indexOf(" ", start) + 1;
  let end = Math.min(text.length, start + excerptLength);
  if (end < text.length) end = Math.max(text.lastIndexOf(" ", end), start);
  const shown = `${start > 0 ? "… " : ""}${text.slice(start, end)}${end < text.length ? " …" : ""}`;

  const parts: ExcerptPart[] = [];
  let position = 0;
  for (const match of shown.matchAll(pattern)) {
    if (match.index > position)
      parts.push({ text: shown.slice(position, match.index), match: false });
    parts.push({ text: match[0], match: true });
    position = match.index + match[0].length;
  }
  if (position < shown.length) parts.push({ text: shown.slice(position), match: false });
  return parts;
}

// Returns the search function for a set of pages. A hit must hold every word of
// the query. When no entry does, a hit can hold any of them.
export function createDocSearch(pages: DocPage[]): (query: string) => DocSearchHit[] {
  const synonyms = glossarySynonyms(pages);
  // A page's intro has no section title, so its page title stands in. A search
  // for a page's title then puts that page first.
  const index = new MiniSearch<DocSearchEntry & { title: string }>({
    fields: ["title", "page", "text"],
    storeFields: ["slug", "sectionId", "page", "section", "text"],
    processTerm: (term) => singular(term.toLowerCase()),
  });
  index.addAll(
    docSearchEntries(pages).map((entry) => ({ ...entry, title: entry.section || entry.page })),
  );

  return (query) => {
    const terms = queryTerms(query, synonyms);
    if (!terms.length) return [];
    const options = {
      boost: { title: 3, page: 1.5 },
      prefix: true,
      fuzzy: (term: string) => (term.length > 4 ? 0.2 : false),
    };
    const text = terms.join(" ");
    let results = index.search(text, { ...options, combineWith: "AND" });
    if (!results.length) results = index.search(text, options);
    // A title that holds the whole query goes first: "who approves" finds "Who
    // approves what" before the sections that only say "approve".
    const phrase = normalize(query, synonyms);
    const holdsPhrase = (result: (typeof results)[number]) =>
      ((result.section || result.page) as string).toLowerCase().includes(phrase);
    return [...results.filter(holdsPhrase), ...results.filter((result) => !holdsPhrase(result))]
      .slice(0, maxHits)
      .map((result) => ({
        id: result.id as string,
        slug: result.slug as string,
        sectionId: result.sectionId as string,
        page: result.page as string,
        section: result.section as string,
        excerpt: excerpt(result.text as string, terms),
      }));
  };
}
