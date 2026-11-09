// The user docs, loaded from src/docs. Import this module only from the docs
// routes, or dynamically, so that the docs and the Markdown parser stay out of
// the main bundle.
import {
  findSection,
  orderDocPages,
  parseDoc,
  withNewTabLinks,
  type DocPage,
  type DocSection,
} from "$lib/docsFormat";
import { createDocSearch } from "$lib/docsSearch";

const files = import.meta.glob(["/src/docs/*.md", "!/src/docs/README.md"], {
  query: "?raw",
  import: "default",
  eager: true,
}) as Record<string, string>;

export const docPages: DocPage[] = orderDocPages(
  Object.entries(files)
    .map(([path, source]) => ({
      slug: path.replace("/src/docs/", "").replace(".md", ""),
      source,
    }))
    .map(({ slug, source }) => parseDoc(slug, source)),
);

// The search that the Help search bar runs on each keystroke.
export const searchDocs = createDocSearch(docPages);

export function findDoc(slug: string): DocPage | undefined {
  return docPages.find((page) => page.slug === slug);
}

// The section that a help popover shows.
export function findHelpSection(ref: string): DocSection | null {
  const section = findSection(docPages, ref)?.section;
  return section ? { ...section, blocks: withNewTabLinks(section.blocks) } : null;
}
