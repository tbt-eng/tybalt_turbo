import { error } from "@sveltejs/kit";
import { get } from "svelte/store";
import { goto } from "$app/navigation";
import { resolve } from "$app/paths";
import {
  canViewComputers,
  computerListSearch,
  computerListState,
  type ComputerListState,
  type ComputerPage,
} from "$lib/computers";
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

// Loads the page of a server-paged list (the sign-ins or the rejected
// reports) that the URL names. The server moves a page past the end back to
// the last page.
export async function loadComputerPage<T>(path: string, url: URL) {
  await requireComputersAccess();
  const state = computerListState(url.searchParams);
  const result = (await pb.send(path, {
    method: "GET",
    query: { q: state.q, page: state.page, perPage: state.perPage },
  })) as ComputerPage<T>;
  return { result, state };
}

// Moves a server-paged list to a new search or page by changing the URL, so
// the page's load runs again. A change of search waits until typing pauses.
export function computerListNavigator(
  route: "/computers/logins" | "/computers/raw_logins",
  initial: ComputerListState,
) {
  let lastQ = initial.q;
  let timer: ReturnType<typeof setTimeout> | undefined;
  return (state: ComputerListState) => {
    clearTimeout(timer);
    const go = () =>
      // eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolved; the rule allows no search after it
      goto(resolve(route) + computerListSearch(state), {
        replaceState: true,
        keepFocus: true,
        noScroll: true,
      });
    const searchChanged = state.q.trim() !== lastQ;
    lastQ = state.q.trim();
    if (searchChanged) timer = setTimeout(go, 300);
    else go();
  };
}
