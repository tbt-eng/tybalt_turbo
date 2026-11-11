/**
 * App configuration store
 *
 * Loads configuration values from the app_config collection and provides
 * derived stores for specific feature flags like jobsEditingEnabled.
 */

import { writable, derived, get } from "svelte/store";
import { pb } from "$lib/pocketbase";
import type { UnsubscribeFunc } from "pocketbase";
import type { AppConfigResponse } from "$lib/pocketbase-types";
import { deployGroupAllows } from "$lib/deployGroups";

interface AppConfigState {
  items: AppConfigResponse[];
  // The ids of the deploy groups that the signed-in user is in. See
  // $lib/deployGroups.
  deployGroups: string[];
  loading: boolean;
  initialized: boolean;
  error: string | null;
}

const store = writable<AppConfigState>({
  items: [],
  deployGroups: [],
  loading: false,
  initialized: false,
  error: null,
});

let unsubscribeFunc: UnsubscribeFunc | null = null;
let unsubscribeGroupsFunc: UnsubscribeFunc | null = null;
let loadPromise: Promise<void> | null = null;
// The user that the deploy groups were loaded for. Sign-out doesn't reset this
// store, so a second user in the same tab must not keep the first user's groups.
let loadedFor = "";

/**
 * Initialize the store by loading all config records and setting up
 * realtime subscription for updates.
 */
async function init() {
  // Check authentication before proceeding
  if (!pb.authStore.token || !pb.authStore.record) return;

  // If already initialized for this user, return (idempotent)
  if (get(store).initialized && loadedFor === pb.authStore.record.id) return;

  // A caller that arrives during a load waits for that load
  loadPromise ??= load().finally(() => {
    loadPromise = null;
  });
  return loadPromise;
}

// Reads the deploy groups of the signed-in user. The collection's rule returns
// only the user's own rows. On an error the user is in no group, so a feature
// that names a group stays hidden.
async function loadDeployGroups(): Promise<string[]> {
  try {
    const rows = await pb
      .collection("deploy_group_memberships")
      .getFullList<{ deploy_group: string }>({ requestKey: "deploy_groups_init" });
    return rows.map((row) => row.deploy_group);
  } catch {
    return [];
  }
}

async function load() {
  store.update((state) => ({ ...state, loading: true }));
  const uid = pb.authStore.record?.id ?? "";

  try {
    // Load all config records
    const [items, deployGroups] = await Promise.all([
      pb.collection("app_config").getFullList<AppConfigResponse>({
        requestKey: "app_config_init",
      }),
      loadDeployGroups(),
    ]);
    loadedFor = uid;

    store.update((state) => ({
      ...state,
      items,
      deployGroups,
      loading: false,
      initialized: true,
      error: null,
    }));

    // Set up realtime subscription
    if (unsubscribeFunc) {
      unsubscribeFunc();
    }

    unsubscribeFunc = await pb.collection("app_config").subscribe("*", async (e) => {
      if (e.action === "create") {
        store.update((state) => ({
          ...state,
          items: [...state.items, e.record as unknown as AppConfigResponse],
        }));
      } else if (e.action === "update") {
        store.update((state) => ({
          ...state,
          items: state.items.map((item) =>
            item.id === e.record.id ? (e.record as unknown as AppConfigResponse) : item,
          ),
        }));
      } else if (e.action === "delete") {
        store.update((state) => ({
          ...state,
          items: state.items.filter((item) => item.id !== e.record.id),
        }));
      }
    });

    // A change to the user's memberships shows without a reload. The server
    // sends only the user's own rows. A failure here must not undo the config
    // load above: the groups then change on the next reload.
    if (unsubscribeGroupsFunc) {
      unsubscribeGroupsFunc();
      unsubscribeGroupsFunc = null;
    }
    try {
      unsubscribeGroupsFunc = await pb
        .collection("deploy_group_memberships")
        .subscribe("*", async () => {
          const deployGroups = await loadDeployGroups();
          store.update((state) => ({ ...state, deployGroups }));
        });
    } catch {
      // Keep the groups that were loaded.
    }
  } catch (error) {
    store.update((state) => ({
      ...state,
      loading: false,
      initialized: false,
      error: error instanceof Error ? error.message : "Failed to load config",
    }));
  }
}

/**
 * Clean up subscription when store is no longer needed
 */
function unsubscribe() {
  if (unsubscribeFunc) {
    unsubscribeFunc();
    unsubscribeFunc = null;
  }
  if (unsubscribeGroupsFunc) {
    unsubscribeGroupsFunc();
    unsubscribeGroupsFunc = null;
  }
  loadedFor = "";
  store.update((state) => ({
    ...state,
    initialized: false,
    items: [],
    deployGroups: [],
    loading: false,
    error: null,
  }));
}

/**
 * Helper to get a config value by domain key
 */
function getConfigValue(items: AppConfigResponse[], domainKey: string): unknown | null {
  const config = items.find((item) => item.key === domainKey);
  if (!config) return null;
  return config.value;
}

/**
 * Helper to get a boolean property from a domain config
 * Returns defaultValue if domain or property doesn't exist
 */
function getConfigBool(
  items: AppConfigResponse[],
  domainKey: string,
  property: string,
  defaultValue: boolean,
): boolean {
  const value = getConfigValue(items, domainKey);
  if (value === null || typeof value !== "object") return defaultValue;

  const obj = value as Record<string, unknown>;
  if (!(property in obj)) return defaultValue;

  const propValue = obj[property];
  if (typeof propValue !== "boolean") return defaultValue;

  return propValue;
}

// Derived store for jobs editing enabled
// Reads from app_config where key="jobs", checks value.create_edit_absorb
// Defaults to true (fail-open) if config is missing
export const jobsEditingEnabled = derived(store, ($store) => {
  return getConfigBool($store.items, "jobs", "create_edit_absorb", true);
});

// Derived store for expenses/PO/vendor editing enabled
// Reads from app_config where key="expenses", checks value.create_edit_absorb
// Defaults to true (fail-open) if config is missing
export const expensesEditingEnabled = derived(store, ($store) => {
  return getConfigBool($store.items, "expenses", "create_edit_absorb", true);
});

// Derived store for time entry editing and selected timesheet mutations.
// Reads from app_config where key="time", checks value.create_edit.
// Defaults to true (fail-open) if config is missing.
export const timeEditingEnabled = derived(store, ($store) => {
  return getConfigBool($store.items, "time", "create_edit", true);
});

// Derived store for hidden legacy PO create/update flow.
// Reads from app_config where key="purchase_orders", checks
// value.enable_legacy_po_create_update. Defaults to false (fail-closed).
export const legacyPOCreateUpdateEnabled = derived(store, ($store) => {
  return getConfigBool($store.items, "purchase_orders", "enable_legacy_po_create_update", false);
});

export const timeEditingDisabledMessage =
  "Time entry and time amendment creation/editing, plus timesheet submission/approval, are temporarily disabled during a system transition.";

// Export the store with init and unsubscribe methods
export const appConfig = {
  subscribe: store.subscribe,
  init,
  unsubscribe,
};

export const jobRequestsEnabled = derived(store, ($store) =>
  getConfigBool($store.items, "jobs", "require_job_request_approval", false),
);

// Help chat answers questions in Help. Reads from app_config where
// key="help_chat", checks value.enabled. Defaults to false (fail-closed): with
// no row, Help has its search and no Ask. When the row names a deploy group,
// only the users in that group get Ask.
export const helpChatEnabled = derived(
  store,
  ($store) =>
    getConfigBool($store.items, "help_chat", "enabled", false) &&
    deployGroupAllows(getConfigValue($store.items, "help_chat"), $store.deployGroups),
);
