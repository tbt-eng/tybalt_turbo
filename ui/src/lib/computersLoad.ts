import { error } from "@sveltejs/kit";
import { get } from "svelte/store";
import { goto } from "$app/navigation";
import { resolve } from "$app/paths";
import { canViewComputers, computerListSearch } from "$lib/computers";
import { pb } from "$lib/pocketbase";
import { globalStore } from "$lib/stores/global";

// Stops a Computers page for a user without the computers or admin claim. The
// server refuses the data too.
export async function requireComputersAccess() {
  await globalStore.ensureClaims();
  if (!canViewComputers(get(globalStore).claims)) {
    error(403, "The computers or admin claim is required to see computers.");
  }
}

// Loads every row of a server-searched list (the sign-ins or the rejected
// reports) that match the search in the URL. Sign-ins are deleted after 30
// days and Turbo keeps one rejected report per computer name, so the server
// sends them all.
export async function loadComputerList<T>(path: string, url: URL) {
  await requireComputersAccess();
  const q = url.searchParams.get("q")?.trim() ?? "";
  const items = (await pb.send(path, { method: "GET", query: q ? { q } : {} })) as T[];
  return { items, q };
}

// Moves a server-searched list to a new search by changing the URL, so the
// page's load runs again once typing pauses. A change of spaces alone loads
// nothing.
export function computerListNavigator(
  route: "/computers/logins" | "/computers/raw_logins",
  initialQ: string,
) {
  let lastQ = initialQ;
  let timer: ReturnType<typeof setTimeout> | undefined;
  return (q: string) => {
    if (q.trim() === lastQ) return;
    lastQ = q.trim();
    clearTimeout(timer);
    timer = setTimeout(
      () =>
        // eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolved; the rule allows no search after it
        goto(resolve(route) + computerListSearch(q), {
          replaceState: true,
          keepFocus: true,
          noScroll: true,
        }),
      300,
    );
  };
}
