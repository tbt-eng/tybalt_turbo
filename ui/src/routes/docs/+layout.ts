import { docPages } from "$lib/docs";
import { docNavEntry } from "$lib/docsFormat";

export const load = () => ({
  pages: docPages.map(docNavEntry),
});
