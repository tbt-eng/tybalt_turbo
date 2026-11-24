import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { get } from "svelte/store";
import { globalStore } from "$lib/stores/global";
export const load: PageLoad = async () => {
  await globalStore.ensureClaims();
  if (!get(globalStore).claims.includes("payables_admin"))
    error(403, "Only payables staff can list expiring purchase orders.");
};
