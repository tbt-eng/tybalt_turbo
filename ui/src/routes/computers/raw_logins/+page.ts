import type { ComputerRawLoginRow } from "$lib/computers";
import { loadComputerList } from "$lib/computersLoad";
import type { PageLoad } from "./$types";

// Every rejected report that matches the search in the URL, without the
// reports themselves.
export const load: PageLoad = async ({ url }) => {
  return loadComputerList<ComputerRawLoginRow>("/api/computers/raw_logins", url);
};
