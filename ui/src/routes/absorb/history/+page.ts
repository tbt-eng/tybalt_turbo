import type { PageLoad } from "./$types";
import { fetchAbsorbHistory } from "$lib/absorbHistory";

export const load: PageLoad = async () => {
  return { items: await fetchAbsorbHistory() } as const;
};
