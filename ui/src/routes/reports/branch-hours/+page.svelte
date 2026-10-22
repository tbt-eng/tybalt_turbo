<script lang="ts">
  import { untrack } from "svelte";
  import { afterNavigate, replaceState } from "$app/navigation";
  import { page } from "$app/state";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSDateInput from "$lib/components/DSDateInput.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSSelector from "$lib/components/DSSelector.svelte";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import BranchHoursTable from "$lib/components/reports/BranchHoursTable.svelte";
  import { pb } from "$lib/pocketbase";
  import {
    branchHoursDirections,
    branchHoursProjects,
    branchHoursTotals,
    formatHours,
    yearToDate,
    type BranchHoursDirection,
    type BranchHoursReport,
  } from "$lib/reports/branchHours";
  import { canViewWIPReports } from "$lib/reports/wipReports";
  import { authStore } from "$lib/stores/auth";
  import { branches } from "$lib/stores/branches";
  import { globalStore } from "$lib/stores/global";
  import { downloadCSV, shortDate } from "$lib/utilities";

  const userId = $derived($authStore?.model?.id ?? "");
  const managedBranch = $derived($branches.items.find((item) => item.manager === userId));
  const allowed = $derived(canViewWIPReports($globalStore.claims, !!managedBranch));

  // The selection lives in the URL hash so a report can be shared as a link.
  const range = yearToDate();
  let branch = $state("");
  let startDate = $state(range.start);
  let endDate = $state(range.end);
  let direction = $state<BranchHoursDirection>("received");
  let selectionReady = $state(false);
  // SvelteKit cannot read or replace the URL until its router has started.
  let routerReady = $state(false);
  afterNavigate(() => (routerReady = true));
  let search = $state("");
  let report = $state<BranchHoursReport | null>(null);
  let error = $state("");
  let retry = $state(0);
  const requestId = $props.id();

  const dateError = $derived(
    !startDate || !endDate
      ? "Select both dates."
      : startDate > endDate
        ? "The start date must be on or before the end date."
        : "",
  );
  const totals = $derived(branchHoursTotals(report?.rows ?? []));
  const projects = $derived(branchHoursProjects(report?.rows ?? [], direction, search));
  const directionCount = $derived(branchHoursProjects(report?.rows ?? [], direction).length);

  function hashFor(changes: Partial<Record<"branch" | "from" | "to" | "view", string>>) {
    return `#${new URLSearchParams({
      branch,
      from: startDate,
      to: endDate,
      view: direction,
      ...changes,
    }).toString()}`;
  }

  const tabs = $derived(
    branchHoursDirections.map((item) => ({
      label: report ? `${item.label} · ${formatHours(totals[item.key])} h` : item.label,
      href: hashFor({ view: item.key }),
      active: item.key === direction,
    })),
  );

  function readHash() {
    const params = new URLSearchParams(page.url.hash.slice(1));
    const view = params.get("view");
    if (branchHoursDirections.some((item) => item.key === view)) {
      direction = view as BranchHoursDirection;
    }
    startDate = params.get("from") ?? startDate;
    endDate = params.get("to") ?? endDate;
    branch = params.get("branch") ?? branch;
  }

  function writeHash() {
    // eslint-disable-next-line svelte/no-navigation-without-resolve -- same page, only the hash changes
    replaceState(hashFor({}), page.state);
  }

  // Tab links and shared links change the hash after the page has opened.
  $effect(() => {
    void page.url.hash;
    if (untrack(() => selectionReady)) untrack(readHash);
  });

  // Open on the branch the user manages, or else their own branch.
  $effect(() => {
    if (selectionReady || !routerReady || !$branches.initialized || !userId) return;
    let current = true;
    const finish = (fallback: string) =>
      untrack(() => {
        if (!current) return;
        branch = fallback;
        readHash();
        selectionReady = true;
        writeHash();
      });
    if (managedBranch) {
      finish(managedBranch.id);
      return;
    }
    pb.send<{ default_branch: string }>("/api/users/defaults", {
      method: "GET",
      requestKey: `${requestId}-defaults`,
    })
      .then((result) => finish(result.default_branch))
      .catch(() => finish(""));
    return () => {
      current = false;
    };
  });

  function reportQuery(format: "json" | "csv") {
    return new URLSearchParams({ branch, start_date: startDate, end_date: endDate, format });
  }

  $effect(() => {
    const permitted = allowed;
    void branch;
    void startDate;
    void endDate;
    void retry;
    report = null;
    error = "";
    if (!selectionReady || !permitted || !branch || dateError) return;
    untrack(writeHash);
    let current = true;
    const controller = new AbortController();
    pb.send<BranchHoursReport>(`/api/branch_hours?${reportQuery("json")}`, {
      method: "GET",
      requestKey: requestId,
      signal: controller.signal,
    })
      .then((result) => {
        if (current) report = result;
      })
      .catch((reason) => {
        if (current)
          error =
            reason?.status === 403
              ? "You do not have permission to view this report."
              : "The report could not be loaded.";
      });
    return () => {
      current = false;
      controller.abort();
    };
  });

  async function downloadReport() {
    if (!report) return;
    const name = report.branch.replace(/[^A-Za-z0-9]+/g, "_").replace(/^_|_$/g, "");
    await downloadCSV(
      `${pb.baseUrl}/api/branch_hours?${reportQuery("csv")}`,
      `branch_hours_${name}_${startDate}_${endDate}.csv`,
    );
  }
</script>

<svelte:head>
  <title>Branch Hours</title>
</svelte:head>

<main class="space-y-4 p-4">
  <DSPageHeader title="Branch Hours">
    {#snippet context()}
      Regular and training hours by project and person, between a branch and the other branches.
    {/snippet}
    {#snippet actions()}
      {#if report && report.rows.length > 0}
        <DSActionButton
          action={downloadReport}
          icon="mdi:file-download-outline"
          title="Download CSV"
          color="green"
        />
      {/if}
    {/snippet}
  </DSPageHeader>

  {#if !allowed}
    <DSCallout tone="warning">
      You need to be a branch manager or hold the kpi or admin claim to view this report.
    </DSCallout>
  {:else}
    <div class="flex flex-wrap items-center gap-4">
      <div class="w-72">
        <DSSelector
          bind:value={() => branch, (value) => (branch = String(value))}
          items={$branches.items}
          errors={{}}
          fieldName="branch"
          uiName="Branch"
        >
          {#snippet optionTemplate(item)}{item.name}{/snippet}
        </DSSelector>
      </div>
      <label class="flex items-center gap-2">
        From
        <DSDateInput bind:value={startDate} max={endDate || undefined} />
      </label>
      <label class="flex items-center gap-2">
        To
        <DSDateInput bind:value={endDate} min={startDate || undefined} />
      </label>
    </div>

    {#if dateError}
      <DSCallout tone="warning">{dateError}</DSCallout>
    {:else if selectionReady && !branch}
      <DSCallout tone="info">Select a branch to view its hours.</DSCallout>
    {:else if error}
      <DSCallout tone="danger">
        {error}
        {#snippet actions()}
          <DSActionButton action={() => retry++} title="Try again">Try again</DSActionButton>
        {/snippet}
      </DSCallout>
    {:else if !report}
      <p role="status">Loading branch hours…</p>
    {:else}
      <p>
        {report.branch}, {shortDate(startDate, true)} to {shortDate(endDate, true)}.
        {#if direction === "received"}
          Hours that other branches' staff worked for {report.branch}.
        {:else if direction === "supplied"}
          Hours that {report.branch} staff worked for other branches.
        {:else}
          Hours that {report.branch} staff worked for {report.branch}, including time with no
          project.
        {/if}
      </p>
      <div>
        <DSTabBar {tabs} />
        <div class="space-y-2 py-2">
          <DSTextInput
            bind:value={search}
            errors={{}}
            fieldName="search"
            uiName="Search"
            placeholder="Project, client, person or branch"
          />
          {#if directionCount === 0}
            <p role="status">No hours in this period.</p>
          {:else if projects.length === 0}
            <p role="status">No projects match your search.</p>
          {:else}
            {#key direction}
              <BranchHoursTable
                {projects}
                {direction}
                title={`${report.branch} ${direction} hours`}
              />
            {/key}
          {/if}
        </div>
      </div>
    {/if}
  {/if}
</main>
