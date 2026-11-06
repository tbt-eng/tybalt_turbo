import { error } from "@sveltejs/kit";
import { findDoc } from "$lib/docs";

// Help opens on Start here.
export const load = () => {
  const page = findDoc("start-here");
  if (!page) error(404, "Page not found");
  return { page };
};
