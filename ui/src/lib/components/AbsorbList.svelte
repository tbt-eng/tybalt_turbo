<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import DsActionButton from "./DSActionButton.svelte";
  import type { AbsorbActionsResponse } from "$lib/pocketbase-types";
  import { goto } from "$app/navigation";
  import { getAbsorbRedirectUrl } from "$lib/utilities";
  import { jobsEditingEnabled, expensesEditingEnabled } from "$lib/stores/appConfig";
  import RateSheetLabel from "./rateSheets/RateSheetLabel.svelte";
  import { absorbedRateSheets } from "$lib/rateSheetAbsorb";
  import { rateSheets } from "$lib/stores/rateSheets";
  import DSCallout from "./DSCallout.svelte";
  import DSFieldList from "./DSFieldList.svelte";
  import DSField from "./DSField.svelte";

  // Helper: check if actions are disabled for a given absorb item
  function isAbsorbActionDisabled(item: AbsorbActionsResponse): boolean {
    const cn = item.collection_name;
    if (
      (cn === "clients" || cn === "client_contacts" || cn === "rate_sheets") &&
      !$jobsEditingEnabled
    )
      return true;
    if (cn === "vendors" && !$expensesEditingEnabled) return true;
    return false;
  }

  let {
    collectionName = undefined,
    returnTo,
    embedded = false,
    onresolved,
    targetLabel,
  }: {
    collectionName?: string;
    // The caller validates this optional client-workspace return route.
    returnTo?: string;
    // Inside another view, such as a modal: no heading, a short summary of
    // each action (record kept, records absorbed, references moved), and
    // Commit and Undo without a second confirmation, because the view that
    // holds the list is the confirmation.
    embedded?: boolean;
    // Called after a commit or an undo instead of leaving the page.
    onresolved?: (outcome: "committed" | "undone", item: AbsorbActionsResponse) => void;
    // The name of an action's kept record, when the caller knows it.
    targetLabel?: (item: AbsorbActionsResponse) => string | undefined;
  } = $props();

  // In embedded mode, Commit takes focus once the action loads.
  let commitBox = $state<HTMLElement | undefined>();
  $effect(() => {
    if (embedded && !loading && commitBox) commitBox.querySelector("button")?.focus();
  });

  // The names of the absorbed records, from their saved copies.
  function absorbedNames(records: unknown): string[] {
    if (!Array.isArray(records)) return [];
    return records.map((record) => {
      const { id, name, alias } = (record ?? {}) as { id?: string; name?: string; alias?: string };
      const label = name || id || "";
      return alias?.trim() ? `${label} (${alias.trim()})` : label;
    });
  }

  let errors = $state<Record<string, { message: string }>>({});
  let items = $state<AbsorbActionsResponse[]>([]);
  let loading = $state<boolean>(true);

  // Track per-item confirmation states
  type ConfirmState = {
    undo: boolean;
    commit: boolean;
  };
  let confirms = $state<Record<string, ConfirmState>>({});

  function setConfirm(id: string, key: keyof ConfirmState, value: boolean) {
    const current = confirms[id] ?? { undo: false, commit: false };
    confirms = { ...confirms, [id]: { ...current, [key]: value } };
  }

  async function loadItems() {
    loading = true;
    try {
      const list = await pb.collection("absorb_actions").getFullList({
        sort: "-created",
        ...(collectionName ? { filter: `collection_name='${collectionName}'` } : {}),
      });
      items = list;
      // A rate sheet absorb names its sheets by code, so load the sheets.
      if (list.some((item) => item.collection_name === "rate_sheets")) rateSheets.init();
    } catch (error: unknown) {
      if (error instanceof Error) {
        errors = { global: { message: error.message } };
      } else {
        errors = { global: { message: "Failed to load absorb actions" } };
      }
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadItems();
  });

  async function undoAbsorb(item: AbsorbActionsResponse) {
    try {
      await pb.send(`/api/${item.collection_name}/undo_absorb`, { method: "POST" });
      if (onresolved) return onresolved("undone", item);
      // After successful undo, redirect to appropriate page
      const url = returnTo || (await getAbsorbRedirectUrl(item.collection_name, item.target_id));
      goto(url);
    } catch (error: unknown) {
      if (error instanceof Error) {
        errors = { global: { message: error.message } };
      } else {
        errors = { global: { message: "Failed to undo absorb action" } };
      }
    } finally {
      setConfirm(item.id, "undo", false);
    }
  }

  async function commitAbsorb(item: AbsorbActionsResponse) {
    try {
      await pb.collection("absorb_actions").delete(item.id);
      if (onresolved) return onresolved("committed", item);
      // After successful commit, redirect to appropriate page
      const url = returnTo || (await getAbsorbRedirectUrl(item.collection_name, item.target_id));
      goto(url);
    } catch (error: unknown) {
      if (error instanceof Error) {
        errors = { global: { message: error.message } };
      } else {
        errors = { global: { message: "Failed to commit absorb action" } };
      }
    } finally {
      setConfirm(item.id, "commit", false);
    }
  }

  function formatCollectionLabel(name: string): string {
    if (name === "client_contacts") return "Client Contacts";
    if (name === "clients") return "Clients";
    if (name === "vendors") return "Vendors";
    if (name === "rate_sheets") return "Rate Sheets";
    return name;
  }

  // The target of a rate sheet absorb, from the rate sheets store. It is
  // undefined until the store loads, and for any other collection.
  function rateSheetTarget(item: AbsorbActionsResponse) {
    if (item.collection_name !== "rate_sheets") return undefined;
    return $rateSheets.items.find((sheet) => sheet.id === item.target_id);
  }

  function countAbsorbed(records: unknown): number | null {
    if (Array.isArray(records)) return records.length;
    return null;
  }

  function countUpdatedRefs(updated: unknown): number | null {
    try {
      if (updated && typeof updated === "object") {
        // updated[table][column][recordId] = oldValue
        const level1 = Object.values(updated as Record<string, any>);
        let count = 0;
        for (const table of level1) {
          for (const column of Object.values(table as Record<string, any>)) {
            count += Object.keys(column as Record<string, any>).length;
          }
        }
        return count;
      }
    } catch (_) {
      // ignore
    }
    return null;
  }

  function entriesOf(obj: unknown): [string, unknown][] {
    if (obj && typeof obj === "object") {
      return Object.entries(obj as Record<string, unknown>);
    }
    return [];
  }

  function formatPrevValue(value: unknown): string {
    if (
      value === null ||
      typeof value === "string" ||
      typeof value === "number" ||
      typeof value === "boolean"
    ) {
      return String(value);
    }
    try {
      return JSON.stringify(value);
    } catch (_) {
      return String(value);
    }
  }

  // Track expanded/collapsed state per item per table (collapsed by default)
  let expandedTables = $state<Record<string, Record<string, boolean>>>({});
  function isTableExpanded(itemId: string, tableName: string): boolean {
    return !!expandedTables[itemId]?.[tableName];
  }
  function setTableExpanded(itemId: string, tableName: string, value: boolean) {
    const current = expandedTables[itemId] ?? {};
    expandedTables = { ...expandedTables, [itemId]: { ...current, [tableName]: value } };
  }
  function toggleTable(itemId: string, tableName: string) {
    setTableExpanded(itemId, tableName, !isTableExpanded(itemId, tableName));
  }
</script>

{#if embedded}
  <div class="flex flex-col gap-3">
    {#if loading}
      <p>Loading...</p>
    {:else if items.length === 0}
      <DSCallout tone="info">No pending absorb action.</DSCallout>
    {/if}
    {#each loading ? [] : items as item (item.id)}
      {@const names = absorbedNames(item.absorbed_records)}
      <DSFieldList>
        <DSField label="Kept" value={targetLabel?.(item) || item.target_id} />
        <DSField label={`Absorbed (${names.length})`}>
          <ul>
            {#each names as name, index (index)}<li>{name}</li>{/each}
          </ul>
        </DSField>
        <DSField
          label="References moved"
          value={String(countUpdatedRefs(item.updated_references) ?? 0)}
        />
      </DSFieldList>
      {#if isAbsorbActionDisabled(item)}
        <DSCallout tone="warning"
          >Undo and commit are disabled during a system transition.</DSCallout
        >
      {:else}
        <div class="flex flex-wrap justify-end gap-2">
          <DsActionButton action={() => undoAbsorb(item)} variant="secondary">Undo</DsActionButton>
          <span bind:this={commitBox}>
            <DsActionButton action={() => commitAbsorb(item)} variant="primary"
              >Commit</DsActionButton
            >
          </span>
        </div>
      {/if}
    {/each}
    {#if errors.global}
      <DSCallout tone="danger">{errors.global.message}</DSCallout>
    {/if}
  </div>
{:else}
  <div class="flex w-full flex-col gap-4 p-4">
    <h1 class="text-xl font-bold">Absorb Actions</h1>

    {#if loading}
      <div class="text-neutral-500">Loading...</div>
    {:else if items.length === 0}
      <div class="rounded-sm border border-neutral-200 p-4 text-neutral-600">
        No pending absorb actions.
      </div>
    {:else}
      <ul class="flex flex-col gap-4">
        {#each items as item}
          {@const target = rateSheetTarget(item)}
          <li class="rounded-lg border-2 border-yellow-500 bg-yellow-50 p-4">
            <div class="mb-2 flex items-center justify-between">
              <div>
                <h2 class="text-lg font-bold text-yellow-800">
                  Pending Absorb — {formatCollectionLabel(item.collection_name)}
                </h2>
                {#if target}
                  <p class="text-sm text-yellow-800/80">
                    Target: <RateSheetLabel sheet={target} />
                  </p>
                {:else}
                  <p class="text-sm text-yellow-800/80">Target ID: {item.target_id}</p>
                {/if}
              </div>
              <div class="text-xs text-yellow-800/70">
                Created: {new Date(item.created).toLocaleString()}
              </div>
            </div>

            <div class="mb-3 grid grid-cols-1 gap-2 md:grid-cols-3">
              <div class="rounded-sm bg-yellow-100 p-2 text-sm">
                <div class="font-semibold">Summary</div>
                <div>Absorbed records: {countAbsorbed(item.absorbed_records) ?? "n/a"}</div>
                <div>Updated references: {countUpdatedRefs(item.updated_references) ?? "n/a"}</div>
              </div>
              <div class="rounded-sm bg-yellow-100 p-2 text-sm md:col-span-2">
                <div class="font-semibold">Details</div>
                <div class="mt-1 grid grid-cols-1 gap-2 md:grid-cols-2">
                  <div>
                    <div class="text-xs font-medium">Absorbed Records</div>
                    {#if item.collection_name === "rate_sheets"}
                      <!-- The absorbed sheets are deleted, so they have no page to link to. -->
                      <ul class="mt-1 flex flex-col gap-1">
                        {#each absorbedRateSheets(item.absorbed_records) as sheet (sheet.id)}
                          <li><RateSheetLabel {sheet} link={false} /></li>
                        {/each}
                      </ul>
                    {:else}
                      <pre
                        class="mt-1 max-h-48 overflow-auto rounded-sm bg-yellow-200/50 p-2 text-xs">
{JSON.stringify(item.absorbed_records, null, 2)}
                  </pre>
                    {/if}
                  </div>
                  <div>
                    <div class="text-xs font-medium">Updated References</div>
                    <div class="mt-1 flex flex-col gap-3">
                      {#each entriesOf(item.updated_references) as [tableName, columns]}
                        <div class="rounded-sm border border-yellow-200">
                          <button
                            class="flex w-full items-center justify-between bg-yellow-200/60 px-2 py-1 text-left text-xs font-semibold"
                            onclick={() => toggleTable(item.id, tableName as string)}
                            aria-expanded={isTableExpanded(item.id, tableName as string)}
                          >
                            <span>{tableName}</span>
                            <span>{isTableExpanded(item.id, tableName as string) ? "▾" : "▸"}</span>
                          </button>
                          {#if isTableExpanded(item.id, tableName as string)}
                            <div class="p-2">
                              {#each entriesOf(columns) as [columnName, rows]}
                                <div class="mb-3">
                                  <div class="text-xs font-semibold">{columnName}</div>
                                  <table class="mt-1 w-full table-auto border-collapse text-xs">
                                    <thead>
                                      <tr>
                                        <th class="border-b border-yellow-300 px-2 py-1 text-left"
                                          >row id</th
                                        >
                                        <th class="border-b border-yellow-300 px-2 py-1 text-left"
                                          >previous value</th
                                        >
                                      </tr>
                                    </thead>
                                    <tbody>
                                      {#each entriesOf(rows) as [rowId, prevValue]}
                                        <tr>
                                          <td class="border-b border-yellow-100 px-2 py-1 font-mono"
                                            >{rowId}</td
                                          >
                                          <td class="border-b border-yellow-100 px-2 py-1 break-all"
                                            >{formatPrevValue(prevValue)}</td
                                          >
                                        </tr>
                                      {/each}
                                    </tbody>
                                  </table>
                                </div>
                              {/each}
                            </div>
                          {/if}
                        </div>
                      {/each}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex flex-wrap gap-2">
              {#if isAbsorbActionDisabled(item)}
                <p class="text-sm text-yellow-800">
                  Undo and commit are disabled during a system transition.
                </p>
              {:else if confirms[item.id]?.undo}
                <div class="flex flex-col gap-2">
                  <p class="font-bold text-red-600">Confirm undo? This cannot be reversed.</p>
                  <div class="flex gap-2">
                    <DsActionButton action={() => undoAbsorb(item)} color="red"
                      >Confirm Undo</DsActionButton
                    >
                    <DsActionButton action={() => setConfirm(item.id, "undo", false)}
                      >Cancel</DsActionButton
                    >
                  </div>
                </div>
              {:else if confirms[item.id]?.commit}
                <div class="flex flex-col gap-2">
                  <p class="font-bold text-red-600">Confirm commit? This cannot be reversed.</p>
                  <div class="flex gap-2">
                    <DsActionButton action={() => commitAbsorb(item)} color="red"
                      >Confirm Commit</DsActionButton
                    >
                    <DsActionButton action={() => setConfirm(item.id, "commit", false)}
                      >Cancel</DsActionButton
                    >
                  </div>
                </div>
              {:else}
                <DsActionButton action={() => setConfirm(item.id, "undo", true)} color="yellow"
                  >Undo</DsActionButton
                >
                <DsActionButton action={() => setConfirm(item.id, "commit", true)} color="green"
                  >Commit</DsActionButton
                >
              {/if}
            </div>
          </li>
        {/each}
      </ul>
    {/if}

    {#if errors.global}
      <div class="text-red-600">{errors.global.message}</div>
    {/if}
  </div>
{/if}
