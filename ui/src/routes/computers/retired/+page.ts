import type { ComputerRow } from "$lib/computers";
import { requireComputersAccess } from "$lib/computersLoad";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
  await requireComputersAccess();
  const items = (await pb.send("/api/computers", {
    method: "GET",
    query: { retired: 1 },
  })) as ComputerRow[];
  return { items };
};
