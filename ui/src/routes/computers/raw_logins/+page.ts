import type { ComputerRawLoginRow } from "$lib/computers";
import { loadComputerPage } from "$lib/computersLoad";
import type { PageLoad } from "./$types";

// The rejected reports, a page at a time, without the reports themselves.
export const load: PageLoad = async ({ url }) => {
  return loadComputerPage<ComputerRawLoginRow>("/api/computers/raw_logins", url);
};
