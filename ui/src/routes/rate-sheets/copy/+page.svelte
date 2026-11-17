<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import { rateSheetCodeRevision, rateSheetLabel } from "$lib/rateSheetLabel";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import RateSheetMatchCallout from "$lib/components/rateSheets/RateSheetMatchCallout.svelte";
  import OvertimeMultiplier from "$lib/components/rateSheets/OvertimeMultiplier.svelte";
  import { globalStore } from "$lib/stores/global";
  import {
    duplicateRateSheetMatches,
    type RateSheetComparison,
    type RateSheetMatch,
  } from "$lib/rateSheetMatches";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";
  import { untrack } from "svelte";

  let { data }: { data: PageData } = $props();

  const sourceSheet = untrack(() => data.sourceSheet);
  const sourceEntries = untrack(() => data.sourceEntries);
  const isRevising = untrack(() => data.isRevising);
  const nextRevision = untrack(() => data.nextRevision);
  const source = {
    code: sourceSheet.code as string,
    name: sourceSheet.name as string,
    revision: sourceSheet.revision as number,
  };

  let errors = $state({} as Record<string, { message: string }>);
  let saveError = $state("");
  let saving = $state(false);
  const entryError = (index: number, field: string) =>
    errors[`entries.${index}.${field}`] as { message: string } | undefined;
  const entriesError = () => errors.entries as { message: string } | undefined;
  // Errors of fields that the form has no input for: a revision's name, which
  // it keeps from its sheet, the revision and the code. The server sends them
  // under these keys.
  const sheetErrors = $derived(
    (isRevising ? ["name", "revision", "code"] : ["revision", "code"])
      .map((key) => errors[key]?.message)
      .filter((message): message is string => !!message),
  );
  // Whether a field or callout on this page shows the server's error for this
  // key. A key that none shows, such as previous_effective_date, needs the
  // message beside Save.
  const shownOnPage = (key: string) =>
    ["name", "effective_date", "revision", "code", "entries"].includes(key) ||
    /^entries\.\d+\.(rate|overtime_rate|role)$/.test(key);

  // The name is an optional description. A revision keeps its sheet's name
  // and code, and a new sheet gets its code when it is saved.
  let name = $state(isRevising ? sourceSheet.name : "");
  let effective_date = $state("");

  // Staged entries - editable copies of source entries, sorted by role name
  let stagedEntries = $state(
    sourceEntries
      .map((e) => ({
        role: e.role,
        roleName: e.expand?.role?.name ?? "Unknown",
        rate: e.rate,
        overtime_rate: e.overtime_rate,
      }))
      .sort((a, b) => a.roleName.localeCompare(b.roleName)),
  );

  // Compare the staged entries with the active sheets as the user edits them.
  // Exact matches would block the save; near matches are information. The
  // compare endpoint needs the job, rate_sheet_revise or admin claim, so other
  // users see only the result of the save.
  const canCompare = $derived(
    ["job", "rate_sheet_revise", "admin"].some((claim) => $globalStore.claims.includes(claim)),
  );
  const noComparison: RateSheetComparison = { exact: [], near: [] };
  let comparison = $state<RateSheetComparison>(noComparison);
  let compareSequence = 0;
  // Sheets that blocked the last save because they have the same rates.
  let blockingMatches = $state<RateSheetMatch[] | null>(null);

  $effect(() => {
    // A revision is compared without the active revision of its own sheet,
    // which it replaces.
    const body = JSON.stringify({
      ...(isRevising ? { code: sourceSheet.code } : {}),
      entries: stagedEntries.map((entry) => ({
        role: entry.role,
        rate: Number(entry.rate),
        overtime_rate: Number(entry.overtime_rate),
      })),
    });
    const sequence = ++compareSequence;
    // An edit replaces the result of the last save.
    blockingMatches = null;
    if (!canCompare || stagedEntries.length === 0) {
      comparison = noComparison;
      return;
    }
    const timer = setTimeout(async () => {
      try {
        const response = await pb.send("/api/rate_sheets/compare", {
          method: "POST",
          body,
          requestKey: null,
        });
        if (sequence === compareSequence) {
          comparison = { exact: response.exact ?? [], near: response.near ?? [] };
        }
      } catch (error) {
        console.error("Failed to compare rate sheets:", error);
        if (sequence === compareSequence) comparison = noComparison;
      }
    }, 400);
    return () => clearTimeout(timer);
  });

  async function save(event: Event) {
    event.preventDefault();
    saveError = "";
    errors = {};
    blockingMatches = null;

    // Validate
    if (!effective_date) {
      errors = { effective_date: { message: "Effective date is required" } };
      return;
    }

    saving = true;

    try {
      // Create rate sheet and all entries atomically via custom endpoint
      const response = await pb.send("/api/rate_sheets", {
        method: "POST",
        // A revision names the sheet it revises by its code and keeps the
        // sheet's name.
        body: JSON.stringify({
          ...(isRevising ? { code: sourceSheet.code } : { name: name.trim() }),
          effective_date,
          revision: isRevising ? nextRevision : 0,
          entries: stagedEntries.map((entry) => ({
            role: entry.role,
            rate: entry.rate,
            overtime_rate: entry.overtime_rate,
          })),
        }),
      });

      // Success - redirect to details page
      saving = false;
      goto(`/rate-sheets/${response.id}/details`);
    } catch (error: any) {
      console.error("Failed to save:", error);
      blockingMatches = duplicateRateSheetMatches(error, "entries");
      if (blockingMatches?.length) {
        // The callout below the entries names the matching sheets.
        saving = false;
        return;
      }
      if (error.data?.data) {
        errors = error.data.data;
      }
      // Field errors show at their fields; the message adds nothing to them.
      // A refusal that no field shows keeps the message, so a refused save
      // always shows something.
      const keys = Object.keys(errors);
      saveError =
        keys.length > 0 && keys.every(shownOnPage)
          ? ""
          : (error.data?.message ?? error.message ?? "Failed to create rate sheet");
      saving = false;
    }
  }
</script>

<div class="p-4">
  <div class="mb-4 flex items-center gap-2">
    <h1 class="text-xl font-bold">
      {isRevising ? "Revise Rate Sheet" : "New Rate Sheet"}
    </h1>
    <div class="ml-auto">
      {#if isRevising}
        <DSHelp doc="rate-sheets#revise-a-rate-sheet" label="Help" />
      {:else}
        <DSHelp doc="rate-sheets#create-a-rate-sheet" label="Help" />
      {/if}
    </div>
  </div>
  <p class="mb-6 text-neutral-600">
    {isRevising
      ? `Revising ${rateSheetCodeRevision(source)}`
      : `Creating a new rate sheet based on ${rateSheetLabel(source)}`}
  </p>

  <form onsubmit={save} class="flex flex-col gap-4">
    {#if isRevising}
      <!-- A revision keeps its sheet's code and name. -->
      <DSFieldList>
        <DSField
          label="New revision"
          value={rateSheetCodeRevision({ code: source.code, revision: nextRevision })}
        />
        <DSField label="Name" value={name} />
      </DSFieldList>
    {:else}
      <DsTextInput
        bind:value={name}
        {errors}
        fieldName="name"
        uiName="Name"
        placeholder="Optional description, such as the client or contract"
        helperText="Optional. Turbo gives the new sheet a code when you save it."
      />
    {/if}

    <!-- Effective Date field -->
    <div
      class="flex w-full flex-col gap-2 {errors.effective_date !== undefined ? 'bg-red-200' : ''}"
    >
      <span class="flex w-full gap-2">
        <label for="effective_date">Effective Date</label>
        <DsDateInput
          class="flex-1 rounded-sm border border-neutral-300 px-1 md:flex-none"
          id="effective_date"
          name="effective_date"
          bind:value={effective_date}
        />
      </span>
      {#if errors.effective_date !== undefined}
        <span class="text-red-600">{errors.effective_date.message}</span>
      {/if}
    </div>

    <!-- Entries Table -->
    <div class="mt-4">
      <h2 class="mb-2 text-lg font-semibold">Rate Entries</h2>
      <p class="mb-3 text-sm text-neutral-600">
        Edit the rates below. All entries will be created when you save.
      </p>

      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <!-- Multiplier row above headers -->
            <tr>
              <th></th>
              <th></th>
              <th class="px-3 py-2 text-left">
                <OvertimeMultiplier bind:rows={stagedEntries} />
              </th>
            </tr>
            <tr class="border-b border-neutral-300">
              <th class="px-3 py-2 text-left">Role</th>
              <th class="px-3 py-2 text-left">Rate</th>
              <th class="px-3 py-2 text-left">Overtime Rate</th>
            </tr>
          </thead>
          <tbody>
            {#each stagedEntries as entry, i}
              <tr class="border-b border-neutral-200">
                <td class="px-3 py-2 text-neutral-700">{entry.roleName}</td>
                <td class="px-3 py-2">
                  <input
                    type="number"
                    min="1"
                    step="1"
                    bind:value={stagedEntries[i].rate}
                    class={`w-28 rounded-sm border px-2 py-1 ${
                      entryError(i, "rate") ? "border-red-500 bg-red-50" : "border-neutral-300"
                    }`}
                  />
                  {#if entryError(i, "rate")}
                    <div class="mt-1 text-xs text-red-600">{entryError(i, "rate")?.message}</div>
                  {/if}
                </td>
                <td class="px-3 py-2">
                  <input
                    type="number"
                    min="1"
                    step="0.01"
                    bind:value={stagedEntries[i].overtime_rate}
                    class={`w-28 rounded-sm border px-2 py-1 ${
                      entryError(i, "overtime_rate")
                        ? "border-red-500 bg-red-50"
                        : "border-neutral-300"
                    }`}
                  />
                  {#if entryError(i, "overtime_rate")}
                    <div class="mt-1 text-xs text-red-600">
                      {entryError(i, "overtime_rate")?.message}
                    </div>
                  {/if}
                </td>
              </tr>
              {#if entryError(i, "role")}
                <tr class="border-b border-neutral-200">
                  <td class="px-3 pb-2 text-xs text-red-600" colspan="3">
                    {entryError(i, "role")?.message}
                  </td>
                </tr>
              {/if}
            {/each}
          </tbody>
        </table>
      </div>

      {#if stagedEntries.length === 0}
        <p class="mt-2 text-neutral-500">
          The source rate sheet has no entries. The new sheet will be created empty.
        </p>
      {/if}
      {#if entriesError()}
        <p class="mt-2 text-red-600">{entriesError()?.message}</p>
      {/if}
    </div>

    <!-- Comparison with the active sheets -->
    {#if blockingMatches && blockingMatches.length > 0}
      <RateSheetMatchCallout
        tone="danger"
        title="Not saved"
        intro="An active sheet has the same rates. Use that sheet, or change a rate."
        matches={blockingMatches}
      />
    {:else if comparison.exact.length > 0}
      <RateSheetMatchCallout
        tone="warning"
        title="Same rates as an active sheet"
        intro="Turbo won't save a sheet with the same rates. Use that sheet, or change a rate."
        matches={comparison.exact}
      />
    {/if}
    {#if comparison.near.length > 0}
      <RateSheetMatchCallout
        tone="info"
        intro="Close to these active sheets:"
        matches={comparison.near}
      />
    {/if}

    {#if sheetErrors.length > 0}
      <DSCallout tone="danger" title="Not saved">
        {#each sheetErrors as message (message)}<p>{message}</p>{/each}
      </DSCallout>
    {/if}

    <!-- Save/Cancel buttons -->
    <div class="mt-4 flex items-center gap-2">
      <DsActionButton type="submit" disabled={saving}>
        {saving ? "Saving..." : "Save"}
      </DsActionButton>
      <DsActionButton action="/rate-sheets/list">Cancel</DsActionButton>

      {#if saveError}
        <span class="text-red-600">{saveError}</span>
      {/if}
    </div>
  </form>

  <!-- Back link -->
  <div class="mt-6">
    <a href="/rate-sheets/list" class="text-blue-600 hover:underline">
      &larr; Back to Rate Sheets
    </a>
  </div>
</div>
