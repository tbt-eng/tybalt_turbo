// Turns a Help chat answer into safe HTML, and its sources into links. This
// file has no Vite or Svelte imports, so tests can run it directly.
//
// A language model writes the answer, so the text is not trusted. It can hold
// anything that a user or a document got the model to repeat. Nothing in it
// becomes markup except the few Markdown forms below.
import { Marked } from "marked";

// What the server says an answer cites. See Source in app/helpchat/turn.go.
export interface ChatSource {
  kind: "doc" | "chunk";
  // A page of the user docs, such as /docs/time#submit-a-week.
  url?: string;
  // The mark in the answer for a chunk, such as S3.
  ref?: string;
  id?: string;
  path?: string;
  start?: number;
  end?: number;
  heading?: string;
  // notes or code
  type?: string;
}

export interface SourceLink {
  // The source's number in the list under the answer, from 1.
  number: number;
  // What to send as feedback when the user opens it.
  key: string;
  label: string;
  href: string;
  // An external link opens in a new tab.
  external: boolean;
}

const repoUrl = "https://github.com/tbt-eng/tybalt_turbo/blob/main";

// A link in an answer may go only to a page of the user docs.
const docsLink = /^\/docs\/[a-z0-9-]+(#[a-z0-9-]+)?$/;

// [S3] or [S3, S4]
const chunkMarks = /\[(S\d+(?:\s*,\s*S\d+)*)\]/g;

export function escapeHtml(text: string): string {
  return text
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

const markdown = new Marked({
  renderer: {
    // Raw HTML in the answer shows as text.
    html({ text }) {
      return escapeHtml(text);
    },
    link({ href, tokens }) {
      const text = this.parser.parseInline(tokens);
      return docsLink.test(href) ? `<a class="link" href="${href}">${text}</a>` : text;
    },
    image({ text }) {
      return escapeHtml(text);
    },
  },
});

// Renders an answer. While the answer streams, sources is null and the chunk
// marks are left out. When the server has sent the sources, each mark that it
// kept becomes the number of that source in the list.
export function renderAnswer(text: string, sources: ChatSource[] | null): string {
  const numbers = new Map<string, number>();
  (sources ?? []).forEach((source, index) => {
    if (source.kind === "chunk" && source.ref) numbers.set(source.ref, index + 1);
  });
  const html = markdown.parse(text, { async: false, gfm: true, breaks: false });
  return html.replace(chunkMarks, (_match, refs: string) =>
    refs
      .split(",")
      .map((ref) => numbers.get(ref.trim()))
      .filter((number) => number !== undefined)
      .map((number) => `<sup>${number}</sup>`)
      .join(""),
  );
}

// The plain text of an answer, for Copy: no chunk marks.
export function answerText(text: string): string {
  return text
    .replace(chunkMarks, "")
    .replace(/[ \t]+([.,;:])/g, "$1")
    .trim();
}

// Makes the list of links under an answer. docTitles maps a docs page's slug to
// its title.
export function sourceLinks(
  sources: ChatSource[],
  docTitles: Record<string, string>,
): SourceLink[] {
  return sources.flatMap((source, index): SourceLink[] => {
    const number = index + 1;
    if (source.kind === "doc" && source.url && docsLink.test(source.url)) {
      const slug = source.url.slice("/docs/".length).split("#")[0];
      return [
        {
          number,
          key: source.url,
          label: docTitles[slug] ?? slug,
          href: source.url,
          external: false,
        },
      ];
    }
    if (
      source.kind === "chunk" &&
      source.ref &&
      source.path &&
      /^[\w./()[\]+@-]+$/.test(source.path)
    ) {
      const lines = source.start && source.end ? `#L${source.start}-L${source.end}` : "";
      const path = source.path.split("/").map(encodeURIComponent).join("/");
      const kind = source.type === "notes" ? "Design note" : "Code";
      const where = source.heading && source.type === "notes" ? source.heading : source.path;
      return [
        {
          number,
          key: source.ref,
          label: `${kind}: ${where}`,
          href: `${repoUrl}/${path}${lines}`,
          external: true,
        },
      ];
    }
    return [];
  });
}
