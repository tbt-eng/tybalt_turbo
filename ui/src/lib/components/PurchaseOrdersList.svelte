<script lang="ts">
  import Icon from "@iconify/svelte";
  import DsFileLink from "$lib/components/DsFileLink.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import { pb } from "$lib/pocketbase";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { type UnsubscribeFunc } from "pocketbase";
  import { proxySubscriptionWithLoader } from "$lib/utilities";
  import DsList from "$lib/components/DSList.svelte";
  import type { PurchaseOrdersResponse } from "$lib/pocketbase-types";
  import type { PurchaseOrdersListData } from "$lib/svelte-types";
  import type { VisiblePurchaseOrderResponse } from "$lib/poVisibility";
  import { globalStore } from "$lib/stores/global";
  import { authStore } from "$lib/stores/auth";
  import POEndPopover from "$lib/components/POEndPopover.svelte";
  import { poMayBeCancelledBy, poMayBeClosedBy } from "$lib/poVisibility";
  import { poEndState, poIsStale, STALE_AFTER_DAYS } from "$lib/poAge";
  import {
    formatCurrencyAmount,
    formatJobLabel,
    pocketBaseFileHref,
    shortDate,
    trimmedOrEmpty,
  } from "$lib/utilities";
  import { onMount, onDestroy, untrack, type Snippet } from "svelte";
  import { expensesEditingEnabled } from "$lib/stores/appConfig";
  import { fetchPendingPO, fetchVisiblePO, isPORevision } from "$lib/poVisibility";
  import PORevisionStatus from "$lib/components/PORevisionStatus.svelte";
  // import { toastStore, type ToastSettings } from "@skeletonlabs/skeleton";

  const collectionId = "purchase_orders";

  // Load the initial data
  let {
    inListHeader,
    data,
    showOwner = false,
    showRemaining = false,
    refreshItems,
    shouldRefreshOnEvent,
    hideWhenEmpty = false,
    filter,
    showTemplateAction = false,
    ageLabels = true,
    searchBarExtra,
  }: {
    inListHeader?: string;
    data: PurchaseOrdersListData;
    showOwner?: boolean;
    showRemaining?: boolean;
    refreshItems?: (() => Promise<VisiblePurchaseOrderResponse[]>) | undefined;
    shouldRefreshOnEvent?:
      | ((record: PurchaseOrdersResponse, currentItems: VisiblePurchaseOrderResponse[]) => boolean)
      | undefined;
    hideWhenEmpty?: boolean;
    filter?: ((item: VisiblePurchaseOrderResponse) => boolean) | undefined;
    showTemplateAction?: boolean;
    // Stale, Expiring and Expired labels on Active purchase orders. The
    // payables lists of exactly those purchase orders turn them off.
    ageLabels?: boolean;
    searchBarExtra?: Snippet;
  } = $props();
  let items = $state(untrack(() => data.items));
  let createdItemIsVisible = $state(untrack(() => data.createdItemIsVisible));
  let pendingApprovalIds = $state(new Set<string>());
  let remainingInfoOpenId = $state<string | null>(null);
  const visibleItems = $derived.by(() => {
    const currentItems = items ?? [];
    return filter ? currentItems.filter(filter) : currentItems;
  });

  // Subscribe to the base collection but update the items from the augmented
  // view
  let unsubscribeFunc: UnsubscribeFunc;
  async function refreshPendingApprovalIds(): Promise<void> {
    try {
      const pending = (await pb.send("/api/purchase_orders/pending", {
        method: "GET",
      })) as VisiblePurchaseOrderResponse[];
      pendingApprovalIds = new Set((pending ?? []).map((po) => po.id));
    } catch (error) {
      // Fail closed for action buttons when pending-state sync fails.
      pendingApprovalIds = new Set<string>();
      console.error("Error loading pending purchase orders:", error);
    }
  }

  onMount(async () => {
    await refreshPendingApprovalIds();
    if (items === undefined) {
      return;
    }
    if (data.realtime_source === "pending") {
      unsubscribeFunc = await proxySubscriptionWithLoader<
        PurchaseOrdersResponse,
        VisiblePurchaseOrderResponse
      >(
        items,
        "purchase_orders",
        (id: string) => fetchPendingPO(id),
        (newItems) => {
          items = newItems;
          void refreshPendingApprovalIds();
        },
        createdItemIsVisible,
      );
    } else if (data.realtime_source === "visible") {
      unsubscribeFunc = await proxySubscriptionWithLoader<
        PurchaseOrdersResponse,
        VisiblePurchaseOrderResponse
      >(
        items,
        "purchase_orders",
        (id: string) => fetchVisiblePO(id),
        (newItems) => {
          items = newItems;
          void refreshPendingApprovalIds();
        },
        createdItemIsVisible,
      );
    } else if (data.realtime_source === "none" && refreshItems) {
      unsubscribeFunc = await pb
        .collection("purchase_orders")
        .subscribe<PurchaseOrdersResponse>("*", async (e) => {
          const currentItems = items ?? [];
          if (shouldRefreshOnEvent && !shouldRefreshOnEvent(e.record, currentItems)) {
            return;
          }
          try {
            items = await refreshItems();
            void refreshPendingApprovalIds();
          } catch (error) {
            console.error("Error refreshing purchase orders:", error);
          }
        });
    }
  });
  onDestroy(async () => {
    if (unsubscribeFunc) {
      unsubscribeFunc();
    }
  });

  // Set to true to show all buttons regardless of user permissions or PO
  // status. This is used for testing purposes.
  const deactivateButtonHiding = $state(false);

  function poIsRejected(po: VisiblePurchaseOrderResponse): boolean {
    return po.status === "Unapproved" && po.rejected !== "";
  }

  function poDisplayStatus(po: VisiblePurchaseOrderResponse): string {
    return poIsRejected(po) ? "Rejected" : po.status;
  }

  function poStatusColor(po: VisiblePurchaseOrderResponse): "green" | "gray" | "yellow" | "red" {
    if (po.status === "Active") return "green";
    if (po.status === "Closed" || po.status === "Cancelled") return "gray";
    if (poIsRejected(po)) return "red";
    return "yellow";
  }

  function rejectionReasonPreview(reason: string, maxLen = 60): string {
    const trimmed = reason.trim();
    if (trimmed.length <= maxLen) return trimmed;
    return `${trimmed.slice(0, maxLen - 1)}…`;
  }

  function poMayBeApprovedOrRejectedByUser(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    if (po.status !== "Unapproved" || po.rejected !== "") {
      return false;
    }
    return pendingApprovalIds.has(po.id);
  }

  function poMayBeCancelledByUser(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    return poMayBeCancelledBy(po, $globalStore.claims, $authStore?.model?.id ?? "");
  }

  function poMayBeClosedByUser(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    return poMayBeClosedBy(po, $globalStore.claims, $authStore?.model?.id ?? "");
  }

  // Closing and cancelling are final, so each opens a confirmation first.
  let endAction = $state<"close" | "cancel">("cancel");
  let endPO = $state<{ id: string; po_number: string }>({ id: "", po_number: "" });
  let showEndPopover = $state(false);
  function confirmEnd(action: "close" | "cancel", po: VisiblePurchaseOrderResponse) {
    endAction = action;
    endPO = { id: po.id, po_number: po.po_number };
    showEndPopover = true;
  }

  function poMayBeDeletedByUser(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    const currentUserId = $authStore?.model?.id ?? "";
    // A revised purchase order keeps its po_number and cannot be deleted.
    return po.status === "Unapproved" && po.po_number === "" && po.uid === currentUserId;
  }

  function poIsOwnedByCurrentUser(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    const currentUserId = $authStore?.model?.id ?? "";
    return po.uid === currentUserId;
  }

  function poMayBeApprovedFromListByOwner(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    return poIsOwnedByCurrentUser(po) && poMayBeApprovedOrRejectedByUser(po);
  }

  function poMayBeEditedByOwner(po: VisiblePurchaseOrderResponse): boolean {
    if (deactivateButtonHiding) return true;
    return poIsOwnedByCurrentUser(po) && po.status === "Unapproved" && po.rejected === "";
  }

  function toggleRemainingInfo(id: string): void {
    remainingInfoOpenId = remainingInfoOpenId === id ? null : id;
  }

  async function del(id: string): Promise<void> {
    // return immediately if items is not an array
    if (!Array.isArray(items)) return;

    try {
      await pb.collection("purchase_orders").delete(id);
      await globalStore.refreshAttentionCounts();
    } catch (error: any) {
      globalStore.addError(error?.response?.message);
    }
  }

  async function approve(id: string): Promise<void> {
    try {
      await pb.send(`/api/purchase_orders/${id}/approve`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      await globalStore.refreshAttentionCounts();
      await goto(resolve("/pos/pending"));
    } catch (error: any) {
      globalStore.addError(error?.response?.message);
    }
  }
</script>

{#if !hideWhenEmpty || visibleItems.length > 0}
  <DsList items={visibleItems} search={true} {inListHeader} {searchBarExtra}>
    {#snippet anchor(item: VisiblePurchaseOrderResponse)}
      <span class="flex flex-col items-center gap-2">
        {#if item.status === "Active"}
          <DsLabel style="inverted" color="green">
            <a href={`/pos/${item.id}/details`} class="hover:underline">
              {item.po_number}
            </a>
          </DsLabel>
        {:else if item.status === "Unapproved"}
          <a href={`/pos/${item.id}/details`} class="text-blue-600 hover:underline">
            {item.po_number || shortDate(item.date, true)}
          </a>
        {:else}
          <a href={`/pos/${item.id}/details`} class="text-blue-600 hover:underline">
            {item.po_number}
          </a>
        {/if}
      </span>
    {/snippet}

    {#snippet headline({ vendor_name, vendor_alias }: VisiblePurchaseOrderResponse)}
      <span class="flex items-center gap-2">
        {vendor_name}
        {#if trimmedOrEmpty(vendor_alias)}
          ({trimmedOrEmpty(vendor_alias)})
        {/if}
      </span>
    {/snippet}

    {#snippet byline(item: VisiblePurchaseOrderResponse)}
      <span class="flex items-center gap-2">
        {formatCurrencyAmount(item.total, item.currency_code)}
        {#if showRemaining && item.status === "Active"}
          <span class="flex items-center gap-1">
            <span>
              Provisional Remaining: {formatCurrencyAmount(
                item.remaining_amount,
                item.currency_code,
              )}
            </span>
            <button
              type="button"
              class="inline-flex items-center text-slate-500 hover:text-slate-700"
              aria-label="Explain provisional remaining"
              aria-expanded={remainingInfoOpenId === item.id}
              onclick={() => toggleRemainingInfo(item.id)}
            >
              <Icon icon="mdi:information-outline" width="16px" />
            </button>
            {#if remainingInfoOpenId === item.id}
              <span class="text-sm text-slate-500">includes uncommitted expenses</span>
            {/if}
          </span>
        {/if}
        {#if item.legacy_manual_entry}
          <DsLabel color="cyan">Manually created</DsLabel>
        {/if}
        {#if isPORevision(item)}
          <PORevisionStatus poId={item.id} />
        {:else if item.status !== "Active"}
          <DsLabel color={poStatusColor(item)}>
            {poDisplayStatus(item)}
          </DsLabel>
        {/if}
        {#if ageLabels}
          {@const endState = poEndState(item)}
          {#if endState === "expired"}
            <DsLabel tone="warning" title={`Ended on ${shortDate(item.end_date)}`}>Expired</DsLabel>
          {:else if endState === "expiring"}
            <DsLabel tone="warning" title={`Ends on ${shortDate(item.end_date)}`}>Expiring</DsLabel>
          {/if}
          {#if poIsStale(item)}
            <DsLabel tone="warning" title={`Approved more than ${STALE_AFTER_DAYS} days ago`}>
              Stale
            </DsLabel>
          {/if}
        {/if}
      </span>
    {/snippet}

    {#snippet line1({ description }: VisiblePurchaseOrderResponse)}
      {description}
    {/snippet}

    {#snippet line2({
      job_number,
      client_name,
      job_description,
      category_name,
    }: VisiblePurchaseOrderResponse)}
      {#if job_number !== ""}
        <span class="flex items-center gap-1">
          {formatJobLabel({ job_number, client_name, job_description })}
          {#if category_name !== ""}
            <DsLabel color="teal">{category_name}</DsLabel>
          {/if}
        </span>
      {/if}
    {/snippet}
    {#snippet line3(item: VisiblePurchaseOrderResponse)}
      <span class="flex items-center gap-1">
        {#if showOwner && item.uid_name}
          <span class="text-sm text-slate-600">Owner: {item.uid_name}</span>
        {/if}
        {#if poIsRejected(item)}
          <DsLabel color="red" title={`${shortDate(item.rejected)}: ${item.rejection_reason}`}>
            <Icon icon="mdi:cancel" width="20px" class="inline-block" />
            {item.rejector_name} ({shortDate(item.rejected)}){#if item.rejection_reason}
              - {rejectionReasonPreview(item.rejection_reason)}{/if}
          </DsLabel>
        {/if}
        <!-- if the item is cumulative, show the sigma icon -->
        {#if item.type === "Cumulative"}
          <DsLabel color="cyan">
            <!-- <Icon icon="mdi:chart-bell-curve-cumulative" width="24px" class="inline-block" /> -->
            <Icon icon="mdi:sigma" width="24px" class="inline-block" />
          </DsLabel>
        {/if}
        <!-- if the item is recurring, show the frequency -->
        {#if item.type === "Recurring"}
          <DsLabel color="cyan">
            <!-- Perhaps use Japanese character for monthly payment-->
            <Icon icon="mdi:recurring-payment" width="24px" class="inline-block" />
            {item.frequency} until {shortDate(item.end_date, true)}
          </DsLabel>
        {/if}
        {#if item.status === "Unapproved" && item.second_approval === "" && item.approved !== ""}
          <!-- Approved, but second approval is required -->
          <span class="flex items-center gap-1">
            /
            <Icon icon="mdi:timer-sand" width="24px" class="inline-block" />
            {item.priority_second_approver_name}
          </span>
        {/if}
        {#if item.attachment}
          <a href={pocketBaseFileHref(collectionId, item.id, item.attachment)} target="_blank">
            <DsFileLink filename={item.attachment as string} />
          </a>
        {/if}
      </span>
    {/snippet}

    {#snippet actions(item: VisiblePurchaseOrderResponse)}
      {#if $expensesEditingEnabled}
        {#if showTemplateAction}
          <DsActionButton
            action={`/pos/add?template=${item.id}`}
            icon="mdi:content-copy"
            title="Copy to New Purchase Order"
            color="blue"
          />
        {/if}
        {#if item.status === "Active"}
          <DsActionButton
            action={`/expenses/add/${item.id}`}
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
        {#if item.status === "Unapproved"}
          {#if poIsRejected(item)}
            {#if poIsOwnedByCurrentUser(item)}
              <DsActionButton
                action={`/pos/${item.id}/edit`}
                icon="mdi:edit-outline"
                title="Edit and Resubmit"
                color="blue"
              />
            {/if}
          {:else if poMayBeEditedByOwner(item)}
            <DsActionButton
              action={`/pos/${item.id}/edit`}
              icon="mdi:edit-outline"
              title="Edit"
              color="blue"
            />
          {/if}
          {#if poMayBeApprovedFromListByOwner(item)}
            <DsActionButton
              action={() => approve(item.id)}
              icon="mdi:approve"
              title="Approve"
              color="green"
            />
          {/if}
          {#if poMayBeDeletedByUser(item)}
            <DsActionButton
              action={() => del(item.id)}
              icon="mdi:delete"
              title="Delete"
              color="red"
            />
          {/if}
        {/if}
      {/if}
    {/snippet}
  </DsList>
{/if}

<POEndPopover bind:show={showEndPopover} action={endAction} po={endPO} />
