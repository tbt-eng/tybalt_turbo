import { get, writable } from "svelte/store";
import { pb } from "$lib/pocketbase";
import type { WorkRecordListRow, WorkRecordSearchResponse } from "$lib/svelte-types";
import { authStore } from "./auth";

type WorkRecordStoreState = {
  items: WorkRecordListRow[];
  loading: boolean;
  error: string | null;
  initialized: boolean;
  lastRefresh: Date;
  maxAge: number;
};

const WORK_RECORDS_MAX_AGE = 20 * 60 * 1000;

function createWorkRecordStore(scope: "mine" | "all") {
  let refreshPromise: Promise<void> | null = null;

  const store = writable<WorkRecordStoreState>({
    items: [],
    loading: false,
    error: null,
    initialized: false,
    lastRefresh: new Date(0),
    maxAge: WORK_RECORDS_MAX_AGE,
  });

  const load = async () => {
    const response = (await pb.send(`/api/work_records?scope=${scope}`, {
      method: "GET",
      requestKey: null,
    })) as WorkRecordSearchResponse;

    store.update((state) => ({
      ...state,
      items: response.items,
      error: null,
      initialized: true,
      lastRefresh: new Date(),
    }));
  };

  const reset = () => {
    store.set({
      items: [],
      loading: false,
      error: null,
      initialized: false,
      lastRefresh: new Date(0),
      maxAge: WORK_RECORDS_MAX_AGE,
    });
  };

  const refresh = async () => {
    if (!get(authStore)?.isValid) {
      reset();
      return;
    }

    if (refreshPromise) return refreshPromise;

    store.update((state) => ({ ...state, loading: true }));
    refreshPromise = load()
      .catch((error: unknown) => {
        store.update((state) => ({
          ...state,
          error: error instanceof Error ? error.message : "Failed to load work records",
        }));
      })
      .finally(() => {
        store.update((state) => ({ ...state, loading: false }));
        refreshPromise = null;
      });

    return refreshPromise;
  };

  const refreshIfStale = async () => {
    if (!get(authStore)?.isValid) {
      reset();
      return;
    }

    const state = get(store);
    if (!state.initialized) {
      await refresh();
      return;
    }

    if (Date.now() - state.lastRefresh.getTime() >= state.maxAge) {
      await refresh();
    }
  };

  const init = async () => {
    await refreshIfStale();
  };

  return {
    subscribe: store.subscribe,
    init,
    refresh,
    refreshIfStale,
    reset,
    unsubscribe: reset,
  };
}

export const myWorkRecords = createWorkRecordStore("mine");
export const allWorkRecords = createWorkRecordStore("all");

export async function refreshVisibleWorkRecords(canViewAll: boolean) {
  const refreshes = [myWorkRecords.refresh()];
  if (canViewAll) {
    refreshes.push(allWorkRecords.refresh());
  }
  await Promise.all(refreshes);
}
