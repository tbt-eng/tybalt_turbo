import type { ComputerLoginRow } from "$lib/computers";
import { loadComputerList } from "$lib/computersLoad";
import type { PageLoad } from "./$types";

// Every sign-in that matches the search in the URL.
export const load: PageLoad = async ({ url }) => {
  return loadComputerList<ComputerLoginRow>("/api/computers/logins", url);
};
