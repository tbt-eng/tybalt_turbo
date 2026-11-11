// Parses the user docs in src/docs. This file has no Vite or Svelte imports, so
// tests can run it directly. See src/docs/README.md for the format.
// Outside the docs routes, import only types from here, because it loads the
// Markdown parser. Load anything else through a dynamic import of $lib/docs.
import { Marked, type Token, type TokensList } from "marked";

export interface FlowStep {
  actor: string;
  step: string;
  detail: string;
}

export type DocBlock =
  { kind: "html"; html: string } | { kind: "flow"; steps: FlowStep[]; notes: string[] };

export interface DocSection {
  id: string;
  title: string;
  // The last commit that a What's new week covers.
  commit?: string;
  blocks: DocBlock[];
  words: number;
}

export interface DocPage {
  slug: string;
  title: string;
  intro: DocBlock[];
  sections: DocSection[];
  words: number;
}

const markdown = new Marked({
  renderer: {
    link({ href, tokens }) {
      return `<a class="link" href="${href}">${this.parser.parseInline(tokens)}</a>`;
    },
  },
});

function inline(text: string): string {
  return markdown.parseInline(text) as string;
}

const commitsUrl = "https://github.com/tbt-eng/tybalt_turbo/commit";

export function commitUrl(commit: string): string {
  return `${commitsUrl}/${commit}`;
}

export function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/['’]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

export function countWords(text: string): number {
  return text.match(/[A-Za-z][\w'’-]*/g)?.length ?? 0;
}

// Each line is "Actor | Step | Detail" or "Note: text".
export function parseFlow(text: string): { steps: FlowStep[]; notes: string[] } {
  const steps: FlowStep[] = [];
  const notes: string[] = [];
  for (const line of text.split("\n").map((value) => value.trim())) {
    if (!line) continue;
    if (line.startsWith("Note:")) {
      notes.push(inline(line.slice("Note:".length).trim()));
      continue;
    }
    const parts = line.split("|").map((value) => value.trim());
    if (parts.length !== 3 || parts.some((value) => !value)) {
      throw new Error(`A flow line needs "Actor | Step | Detail": ${line}`);
    }
    steps.push({ actor: parts[0], step: parts[1], detail: inline(parts[2]) });
  }
  return { steps, notes };
}

function toBlocks(tokens: Token[]): DocBlock[] {
  const blocks: DocBlock[] = [];
  let pending: Token[] = [];
  const flush = () => {
    if (!pending.length) return;
    blocks.push({ kind: "html", html: markdown.parser(pending as TokensList) });
    pending = [];
  };
  for (const token of tokens) {
    if (token.type === "code" && token.lang === "flow") {
      flush();
      blocks.push({ kind: "flow", ...parseFlow(token.text) });
    } else if (token.type !== "space") {
      pending.push(token);
    }
  }
  flush();
  return blocks;
}

export function parseDoc(slug: string, source: string): DocPage {
  let title = "";
  const introTokens: Token[] = [];
  const sectionTokens: { title: string; commit?: string; tokens: Token[] }[] = [];
  for (const token of markdown.lexer(source)) {
    if (token.type === "heading" && token.depth === 1 && !title) {
      title = token.text;
    } else if (token.type === "heading" && token.depth === 2) {
      // A What's new heading ends with its last commit: "... 2026 (to af436ec9)".
      const match = token.text.match(/^(.*) \(to ([0-9a-f]{8,40})\)$/);
      sectionTokens.push({
        title: match?.[1] ?? token.text,
        commit: match?.[2],
        tokens: [],
      });
    } else {
      (sectionTokens.at(-1)?.tokens ?? introTokens).push(token);
    }
  }
  if (!title) throw new Error(`Doc ${slug} needs a "# Title" line`);
  return {
    slug,
    title,
    intro: toBlocks(introTokens),
    sections: sectionTokens.map((section) => ({
      id: slugify(section.title),
      title: section.title,
      commit: section.commit,
      blocks: toBlocks(section.tokens),
      words: countWords(section.tokens.map((token) => token.raw).join(" ")),
    })),
    words: countWords(source),
  };
}

// The order of the pages in Help. src/docs/README.md has the rule:
//   1. Start here.
//   2. The pages of tasks, in the order of Turbo's main menu: Time, Purchase
//      Orders, Expenses, Business. In Business, a proposal comes before the
//      project request that it can lead to.
//   3. The pages about you: your profile.
//   4. Reference: the glossary, then What's new.
// A test fails when a page is in neither this list nor childDocPages.
export const docPageOrder = [
  "start-here",
  "time",
  "purchase-orders",
  "expenses",
  "proposals",
  "project-requests",
  "rate-sheets",
  "profile",
  "glossary",
  "whats-new",
];

// Pages that Help lists under another page, with a shorter label. The children
// of one page show in the order of this list.
const childDocPages: Record<string, { parent: string; label: string }> = {
  "review-time": { parent: "time", label: "Review and track" },
  "revise-a-purchase-order": { parent: "purchase-orders", label: "Revisions" },
  "find-a-purchase-order": { parent: "purchase-orders", label: "Find and print" },
  "child-purchase-orders": { parent: "purchase-orders", label: "Child POs" },
  "commit-expenses": { parent: "expenses", label: "Review and commit" },
  "admin-profiles": { parent: "profile", label: "Admin profiles" },
  "rate-sheet-periods": { parent: "rate-sheets", label: "Periods" },
  "rate-sheet-absorb": { parent: "rate-sheets", label: "Duplicates" },
};

// Help lists pages in the order of docPageOrder. A child page follows its
// parent. A page that the list lacks goes before What's new, by title.
export function orderDocPages(pages: DocPage[]): DocPage[] {
  const unlisted = docPageOrder.indexOf("whats-new") - 0.5;
  const place = ({ slug }: DocPage) => {
    const at = docPageOrder.indexOf(slug);
    return at < 0 ? unlisted : at;
  };
  const children = Object.keys(childDocPages);
  return pages
    .filter(({ slug }) => !childDocPages[slug])
    .sort((a, b) => place(a) - place(b) || a.title.localeCompare(b.title))
    .flatMap((page) => [
      page,
      ...pages
        .filter(({ slug }) => childDocPages[slug]?.parent === page.slug)
        .sort((a, b) => children.indexOf(a.slug) - children.indexOf(b.slug)),
    ]);
}

// What Help shows for a page in its side navigation.
export function docNavEntry({ slug, title }: Pick<DocPage, "slug" | "title">) {
  const child = childDocPages[slug];
  return { slug, label: child?.label ?? title, indent: Boolean(child) };
}

// A help popover sits over a form, so its links must not leave unsaved work.
export function withNewTabLinks(blocks: DocBlock[]): DocBlock[] {
  const newTab = (html: string) =>
    html.replaceAll('<a class="link" ', '<a class="link" target="_blank" rel="noopener" ');
  return blocks.map((block) =>
    block.kind === "html"
      ? { ...block, html: newTab(block.html) }
      : {
          ...block,
          steps: block.steps.map((step) => ({ ...step, detail: newTab(step.detail) })),
          notes: block.notes.map(newTab),
        },
  );
}

// Reads a help reference, such as "project-requests#approve-a-request".
export function findSection(
  pages: DocPage[],
  ref: string,
): { page: DocPage; section: DocSection } | null {
  const [slug, id] = ref.split("#");
  const page = pages.find((candidate) => candidate.slug === slug);
  const section = page?.sections.find((candidate) => candidate.id === id);
  return page && section ? { page, section } : null;
}
