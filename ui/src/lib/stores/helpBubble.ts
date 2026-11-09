/**
 * Help bubble store
 *
 * True when What's new has a week that this browser hasn't seen. Opening What's
 * new clears it. The newest commit is stored in the browser, so the database
 * keeps no per-user record, and the bubble shows once more on a new browser.
 */

import { writable } from "svelte/store";
import { latestWhatsNewCommit } from "$lib/whatsNewCommit";

const storageKey = "help_seen_whats_new";
const { subscribe, set } = writable(false);
// During a visit to What's new, the commit this browser saw before the visit, or
// null if it had seen none. Undefined when the visit has no new weeks to mark.
const seenBefore = writable<string | null | undefined>(undefined);
let latestCommit: Promise<string | null> | null = null;

// A dynamic import keeps the page's text out of the main bundle.
function loadLatestCommit() {
  latestCommit ??= import("../../docs/whats-new.md?raw")
    .then(({ default: source }) => latestWhatsNewCommit(source))
    .catch(() => {
      latestCommit = null;
      return null;
    });
  return latestCommit;
}

async function sync(inWhatsNew: boolean) {
  const commit = await loadLatestCommit();
  if (!commit) return;
  try {
    const stored = localStorage.getItem(storageKey);
    if (!inWhatsNew) seenBefore.set(undefined);
    if (inWhatsNew && stored !== commit) {
      localStorage.setItem(storageKey, commit);
      seenBefore.set(stored);
    }
    set(!inWhatsNew && stored !== commit);
  } catch {
    // Without storage the bubble could never clear, so don't show it.
    seenBefore.set(undefined);
    set(false);
  }
}

export const helpBubble = { subscribe, sync };
export const whatsNewSeenBefore = { subscribe: seenBefore.subscribe };
