import type { ComputerLoginRow } from "$lib/computers";
import { loadComputerPage } from "$lib/computersLoad";
import type { PageLoad } from "./$types";

// Every sign-in, a page at a time: a sign-in is stored at each sign-in, so
// the list is too long to load at once.
export const load: PageLoad = async ({ url }) => {
  return loadComputerPage<ComputerLoginRow>("/api/computers/logins", url);
};
