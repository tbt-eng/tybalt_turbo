import { docPages } from "$lib/docs";

export const load = () => ({
  pages: docPages.map(({ slug, title }) => ({ slug, title })),
});
