<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import Icon from "@iconify/svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsFileLink from "$lib/components/DsFileLink.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DSSearchList from "$lib/components/DSSearchList.svelte";
  import PORevisionStatus from "$lib/components/PORevisionStatus.svelte";
  import { isPORevision, poMayBeCancelledBy, poMayBeClosedBy } from "$lib/poVisibility";
  import POEndPopover from "$lib/components/POEndPopover.svelte";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import { pb } from "$lib/pocketbase";
  import type { POSearchApiResponse } from "$lib/stores/poSearch";
  import { poSearch } from "$lib/stores/poSearch";
  import { expensesEditingEnabled } from "$lib/stores/appConfig";
  import { globalStore } from "$lib/stores/global";
  import {
    formatCurrencyAmount,
    formatJobLabel,
    pocketBaseFileHref,
    shortDate,
    trimmedOrEmpty,
  } from "$lib/utilities";
  import { page } from "$app/stores";

  const collectionId = "purchase_orders";
  type SearchStatus = POSearchApiResponse["status"];

  const allStatusOptions: { id: SearchStatus; label: string }[] = [
    { id: "Active", label: "Active" },
    { id: "Closed", label: "Closed" },
    { id: "Cancelled", label: "Cancelled" },
  ];

  // Only show toggle options for statuses the user actually has visibility into.
  // The backend already filters by permission, so we derive from loaded data.
  const statusOptions = $derived(
    allStatusOptions.filter((opt) =>
      $poSearch.items.some((po: POSearchApiResponse) => po.status === opt.id),
    ),
  );

  let selectedStatus = $state<SearchStatus>("Active");
  const sourceExpense = $derived($page.url.searchParams.get("source_expense") ?? "");

  poSearch.init();

  function poStatusColor(
    status: POSearchApiResponse["status"],
  ): "green" | "gray" | "yellow" | "red" {
    if (status === "Active") return "green";
    if (status === "Closed" || status === "Cancelled") return "gray";
    return "yellow";
  }

  const viewerId = pb.authStore.record?.id ?? "";
  const poMayBeCancelledByUser = (po: POSearchApiResponse) =>
    poMayBeCancelledBy(po, $globalStore.claims, viewerId);
  const poMayBeClosedByUser = (po: POSearchApiResponse) =>
    poMayBeClosedBy(po, $globalStore.claims, viewerId);

  // Closing and cancelling are final, so each opens a confirmation first.
  let endAction = $state<"close" | "cancel">("cancel");
  let endPO = $state<{ id: string; po_number: string }>({ id: "", po_number: "" });
  let showEndPopover = $state(false);
  function confirmEnd(action: "close" | "cancel", po: POSearchApiResponse) {
    endAction = action;
    endPO = { id: po.id, po_number: po.po_number };
    showEndPopover = true;
  }

  // A purchase order under revision was Active, so it is found under Active.
  const statusFilter = $derived(
    (po: POSearchApiResponse) =>
      po.status === selectedStatus || (selectedStatus === "Active" && isPORevision(po)),
  );
</script>

{#if $poSearch.index !== null}
  <DSSearchList
    index={$poSearch.index}
    filter={statusFilter}
    inListHeader={selectedStatus}
    fieldName="purchase_order"
    uiName="search purchase orders..."
    collectionName="purchase_orders_search"
  >
    {#snippet searchBarExtra()}
      <div class="flex items-center gap-2 max-[639px]:w-full max-[639px]:flex-wrap">
        <DSToggle
          bind:value={selectedStatus}
          options={statusOptions}
          ariaLabel="PO status filter"
        />
        <DSHelp doc="find-a-purchase-order#search-for-a-po" label="Help" />
      </div>
    {/snippet}

    {#snippet anchor(item: POSearchApiResponse)}
      {#if item.status === "Active"}
        <DsLabel style="inverted" color="green">
          <a href={`/pos/${item.id}/details`} class="hover:underline">
            {item.po_number}
          </a>
        </DsLabel>
      {:else}
        <a href={`/pos/${item.id}/details`} class="text-blue-600 hover:underline">
          {item.po_number}
        </a>
      {/if}
    {/snippet}

    {#snippet headline({ vendor_name, vendor_alias }: POSearchApiResponse)}
      <span class="flex items-center gap-2">
        {vendor_name}
        {#if trimmedOrEmpty(vendor_alias)}
          ({trimmedOrEmpty(vendor_alias)})
        {/if}
      </span>
    {/snippet}

    {#snippet byline(item: POSearchApiResponse)}
      <span class="flex items-center gap-2">
        {formatCurrencyAmount(item.total, item.currency_code)}
        {#if item.legacy_manual_entry}
          <DsLabel color="cyan">Manually created</DsLabel>
        {/if}
        {#if isPORevision(item)}
          <PORevisionStatus poId={item.id} />
        {:else if item.status !== "Active"}
          <DsLabel color={poStatusColor(item.status)}>
            {item.status}
          </DsLabel>
        {/if}
      </span>
    {/snippet}

    {#snippet line1({ description }: POSearchApiResponse)}
      {description}
    {/snippet}

    {#snippet line2({
      job_number,
      client_name,
      job_description,
      category_name,
    }: POSearchApiResponse)}
      {#if job_number !== ""}
        <span class="flex items-center gap-1">
          {formatJobLabel({ job_number, client_name, job_description })}
          {#if category_name !== ""}
            <DsLabel color="teal">{category_name}</DsLabel>
          {/if}
        </span>
      {/if}
    {/snippet}

    {#snippet line3(item: POSearchApiResponse)}
      <span class="flex items-center gap-1">
        {#if item.uid_name}
          <span class="text-sm text-slate-600">Owner: {item.uid_name}</span>
        {/if}
        {#if item.type === "Cumulative"}
          <DsLabel color="cyan">
            <Icon icon="mdi:sigma" width="24px" class="inline-block" />
          </DsLabel>
        {/if}
        {#if item.type === "Recurring"}
          <DsLabel color="cyan">
            <Icon icon="mdi:recurring-payment" width="24px" class="inline-block" />
            {item.frequency} until {shortDate(item.end_date, true)}
          </DsLabel>
        {/if}
        {#if item.attachment}
          <a href={pocketBaseFileHref(collectionId, item.id, item.attachment)} target="_blank">
            <DsFileLink filename={item.attachment} />
          </a>
        {/if}
      </span>
    {/snippet}

    {#snippet actions(item: POSearchApiResponse)}
      {#if $expensesEditingEnabled && item.status === "Active"}
        <DsActionButton
          action={`/expenses/add/${item.id}${sourceExpense ? `?source_expense=${encodeURIComponent(sourceExpense)}` : ""}`}
          icon="mdi:add-bold"
          title="Create Expense"
          color="green"
        />
        {#if poMayBeCancelledByUser(item)}
          <DsActionButton
            action={() => confirmEnd("cancel", item)}
            icon="mdi:cancel"
            title="Cancel"
            color="orange"
          />
        {/if}
        {#if poMayBeClosedByUser(item)}
          <DsActionButton
            action={() => confirmEnd("close", item)}
            icon="mdi:curtains-closed"
            title="Close Purchase Order"
            color="orange"
          />
        {/if}
      {/if}
    {/snippet}
  </DSSearchList>
{/if}

<POEndPopover bind:show={showEndPopover} action={endAction} po={endPO} />
