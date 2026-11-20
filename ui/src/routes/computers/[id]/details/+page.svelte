<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSMenuItem from "$lib/components/DSMenuItem.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import {
    bootDriveLabel,
    computerAssignment,
    formatBytes,
    lastUserLabel,
    localDay,
    localTime,
    modelLabel,
    networkAdapters,
    systemTypeLabel,
  } from "$lib/computers";
  import { getApiErrorMessage } from "$lib/errors";
  import type { ComputerLoginsResponse } from "$lib/pocketbase-types";
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import { formatDateTime, shortDate } from "$lib/utilities";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  const computer = $derived(data.computer);
  const assignment = $derived(computerAssignment(computer));
  const adapters = $derived(networkAdapters(computer.network_config));

  let busy = $state(false);
  let showRetire = $state(false);
  let retireError = $state<string | null>(null);

  function day(datetime: string): string {
    const local = localDay(datetime);
    return local ? shortDate(local, true) : "";
  }

  async function act(path: "assign" | "retire" | "unretire") {
    busy = true;
    try {
      await pb.send(`/api/computers/${computer.id}/${path}`, {
        method: "POST",
        // Assign refuses if another person signed in since the page loaded.
        body: path === "assign" ? { user: computer.last_user } : undefined,
      });
      showRetire = false;
      await invalidateAll();
    } catch (error) {
      const message = getApiErrorMessage(error, "The change failed");
      if (path === "retire") retireError = message;
      else globalStore.addError(message);
    } finally {
      busy = false;
    }
  }

  function openRetire() {
    retireError = null;
    showRetire = true;
  }
</script>

<svelte:head>
  <title>{computer.computer_name || computer.serial}</title>
</svelte:head>

<div class="mx-auto space-y-6 p-4">
  <DSPageHeader title={computer.computer_name || computer.serial} subtitle={modelLabel(computer)}>
    {#snippet badges()}
      {#if computer.retired}
        <DsLabel tone="neutral">Retired {day(computer.retired)}</DsLabel>
      {/if}
    {/snippet}
    {#snippet context()}
      First seen {day(computer.created)}
    {/snippet}
    {#snippet actions()}
      {#if computer.retired}
        <DSActionButton action={() => act("unretire")} loading={busy} color="blue">
          Un-retire
        </DSActionButton>
      {:else if assignment.canAssign}
        <DSActionButton action={() => act("assign")} loading={busy} color="green">
          Assign to {lastUserLabel(computer)}
        </DSActionButton>
      {/if}
      <DSHelp doc="computers#computer-details" label="Help" />
    {/snippet}
    {#snippet menu()}
      {#if !computer.retired}
        <DSMenuItem label="Retire" icon="mdi:archive-outline" action={openRetire} tone="danger" />
      {/if}
    {/snippet}
  </DSPageHeader>

  <DSCardGrid>
    <DSCard title="Assignment">
      <DSFieldList>
        {#if assignment.state === "unassigned"}
          <DSField label="Assigned to">
            <DsLabel tone="warning" size="sm">Not assigned</DsLabel>
          </DSField>
        {:else}
          <DSField label="Assigned to" value={computer.assigned_to_name} />
          <DSField label="Assigned" value={day(computer.assigned_time)} />
          {#if assignment.state === "other"}
            <DSField label="Last user">
              <DsLabel tone="warning" size="sm">Not the assignee</DsLabel>
            </DSField>
          {/if}
        {/if}
      </DSFieldList>
    </DSCard>

    <DSCard title="Last sign-in">
      <DSFieldList>
        <DSField label="When" value={formatDateTime(computer.last_login.replace(" ", "T"))} />
        <DSField
          label="User"
          value={lastUserLabel(computer)}
          hint={computer.last_user ? "" : "no Turbo user matched"}
        />
        <DSField label="Email" value={computer.email} />
        <DSField label="UPN" value={computer.upn} />
      </DSFieldList>
    </DSCard>

    <DSCard title="Hardware">
      <DSFieldList>
        <DSField label="Serial" value={computer.serial} />
        <DSField label="Manufacturer" value={computer.mfg} />
        <DSField label="Model" value={computer.model} />
        <DSField label="Type" value={systemTypeLabel(computer.system_type)} />
        <DSField label="RAM" value={formatBytes(computer.ram)} />
        <DSField label="Boot drive" value={bootDriveLabel(computer)} />
      </DSFieldList>
    </DSCard>

    <DSCard title="Windows">
      <DSFieldList>
        <DSField label="Version" value={computer.os_version} />
        <DSField label="Architecture" value={computer.os_arch} />
        <DSField label="SKU" value={computer.os_sku || null} />
        <DSField label="Radiator version" value={computer.radiator_version || null} />
      </DSFieldList>
    </DSCard>

    {#each adapters as adapter (adapter.mac)}
      <DSCard title={`Network ${adapter.mac}`}>
        <DSFieldList>
          {#each adapter.fields as field (field.label)}
            <DSField label={field.label} value={field.value} />
          {/each}
        </DSFieldList>
      </DSCard>
    {/each}
  </DSCardGrid>

  <DSSection title="Sign-ins">
    {#snippet titleNote()}
      the newest 20
    {/snippet}
    <DSList items={data.logins} emptyMessage="No sign-ins.">
      {#snippet anchor(login: ComputerLoginsResponse)}
        {day(login.created)}
      {/snippet}
      {#snippet headline(login: ComputerLoginsResponse)}
        {`${login.given_name} ${login.surname}`.trim() || "Unknown"}
      {/snippet}
      {#snippet byline(login: ComputerLoginsResponse)}
        {localTime(login.created)}
      {/snippet}
      {#snippet line1(login: ComputerLoginsResponse)}
        {#if !login.user}
          <DsLabel tone="neutral" size="sm">No Turbo user matched</DsLabel>
        {/if}
      {/snippet}
    </DSList>
  </DSSection>
</div>

<DSPopover
  bind:show={showRetire}
  title="Retire Computer"
  subtitle={computer.computer_name || computer.serial}
  error={retireError}
  submitting={busy}
  submitLabel="Retire"
  danger={true}
  onSubmit={() => act("retire")}
>
  <p>The computer moves to the Retired list. You can un-retire it later.</p>
</DSPopover>
