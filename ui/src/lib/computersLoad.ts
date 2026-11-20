import { error } from "@sveltejs/kit";
import { get } from "svelte/store";
import { canViewComputers } from "$lib/computers";
import { globalStore } from "$lib/stores/global";

// Stops a Computers page for a user without the computers or admin claim. The
// server refuses the data too.
export async function requireComputersAccess() {
  await globalStore.ensureClaims();
  if (!canViewComputers(get(globalStore).claims)) {
    error(403, "The computers or admin claim is required to see computers.");
  }
}
