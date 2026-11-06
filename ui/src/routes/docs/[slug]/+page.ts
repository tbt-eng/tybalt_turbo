import { error } from "@sveltejs/kit";
import { findDoc } from "$lib/docs";

export const load = ({ params }) => {
  const page = findDoc(params.slug);
  if (!page) error(404, "Page not found");
  return { page };
};
