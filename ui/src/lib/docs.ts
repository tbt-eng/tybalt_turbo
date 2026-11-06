// The user docs, loaded from src/docs. Import this module only from the docs
// routes, or dynamically, so that the docs and the Markdown parser stay out of
// the main bundle.
import { findSection, orderDocPages, parseDoc, type DocPage } from "$lib/docsFormat";

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

export function findDoc(slug: string): DocPage | undefined {
  return docPages.find((page) => page.slug === slug);
}

export function findDocSection(ref: string) {
  return findSection(docPages, ref);
}
