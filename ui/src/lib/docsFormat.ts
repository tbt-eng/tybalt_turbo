// Parses the user docs in src/docs. This file has no Vite or Svelte imports, so
// tests can run it directly. See src/docs/README.md for the format.
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
  const sectionTokens: { title: string; tokens: Token[] }[] = [];
  for (const token of markdown.lexer(source)) {
    if (token.type === "heading" && token.depth === 1 && !title) {
      title = token.text;
    } else if (token.type === "heading" && token.depth === 2) {
      sectionTokens.push({ title: token.text, tokens: [] });
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
      blocks: toBlocks(section.tokens),
      words: countWords(section.tokens.map((token) => token.raw).join(" ")),
    })),
    words: countWords(source),
  };
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
