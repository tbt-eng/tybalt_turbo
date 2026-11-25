<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { tick, untrack } from "svelte";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSModal from "$lib/components/DSModal.svelte";
  import DSMergeGroup, { type MergeRow } from "$lib/components/DSMergeGroup.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsEditingDisabledBanner from "$lib/components/DsEditingDisabledBanner.svelte";
  import AbsorbList from "$lib/components/AbsorbList.svelte";
  import { getApiErrorMessage } from "$lib/errors";
  import { shortDate } from "$lib/utilities";
  import { jobsEditingEnabled, expensesEditingEnabled } from "$lib/stores/appConfig";
  import {
    absorbIntoTarget,
    absorbSuggestionCollections,
    absorbSuggestionLevels,
    absorbSuggestionRecordHref,
    dismissAbsorbSuggestion,
    fetchAbsorbSuggestions,
    type AbsorbSuggestionBucket,
    type AbsorbSuggestionCollection,
    type AbsorbSuggestionLevel,
    type AbsorbSuggestionMember,
  } from "$lib/absorbSuggestions";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  // One bucket as the person works on it. The list keeps its order for the
  // session: a finished bucket leaves its place and nothing re-sorts.
  type BucketState = {
    key: string;
    confidence: number;
    suggestedId: string;
    members: AbsorbSuggestionMember[];
    keepId: string;
    picked: string[];
    busy: boolean;
    error: string;
  };

  // The absorb that waits for a commit or an undo. bucketKey is set when it
  // was made on this page, so the bucket can be updated after the commit.
  type Pending = {
    bucketKey: string | null;
    targetName: string;
    absorbed: string[];
    created: string;
    byYou: boolean;
  };

  type TabState = {
    loading: boolean;
    error: string;
    buckets: BucketState[];
    dismissed: number;
    pending: Pending | null;
  };

  const emptyTab = (): TabState => ({
    loading: true,
    error: "",
    buckets: [],
    dismissed: 0,
    pending: null,
  });

  let tabs = $state<Record<AbsorbSuggestionCollection, TabState>>({
    clients: emptyTab(),
    vendors: emptyTab(),
  });
  let levelChoice = $state<AbsorbSuggestionLevel>(untrack(() => data.level));
  let modalOpen = $state(false);
  // Remounts the pending action in the modal, so it loads again.
  let modalRun = $state(0);

  // Derived, so a new data object with the same values changes nothing: a
  // change of tab must not reload the lists.
  const tab = $derived(data.tab);
  const urlLevel = $derived(data.level);
  const current = $derived(tabs[tab]);
  const editingEnabled = $derived(
    tab === "clients" ? $jobsEditingEnabled : $expensesEditingEnabled,
  );
  const countLabel = $derived(tab === "clients" ? "Jobs" : "POs and expenses");
  const movedLabel = $derived(tab === "clients" ? "jobs" : "POs and expenses");

  function pageHref(collection: AbsorbSuggestionCollection, level: AbsorbSuggestionLevel) {
    return `${resolve("/absorb/suggestions")}?tab=${collection}&level=${level}`;
  }

  function navigate(collection: AbsorbSuggestionCollection, level: AbsorbSuggestionLevel) {
    // eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolved in pageHref
    goto(pageHref(collection, level), { noScroll: true, keepFocus: true });
  }

  // The level control changes the URL; the URL changes the list.
  $effect(() => {
    const level = levelChoice;
    if (level !== untrack(() => urlLevel))
      navigate(
        untrack(() => tab),
        level,
      );
  });

  $effect(() => {
    const level = urlLevel;
    untrack(() => {
      levelChoice = level;
      for (const collection of absorbSuggestionCollections) load(collection, level);
    });
  });

  async function load(collection: AbsorbSuggestionCollection, level: AbsorbSuggestionLevel) {
    tabs[collection] = emptyTab();
    try {
      const response = await fetchAbsorbSuggestions(collection, level);
      if (level !== data.level) return; // a later level is loading
      const action = response.pending_action;
      tabs[collection] = {
        loading: false,
        error: "",
        buckets: response.buckets.map(toState),
        dismissed: response.dismissed,
        pending: action
          ? {
              bucketKey: null,
              targetName: action.target_name || action.target_id,
              absorbed: action.absorbed,
              created: action.created,
              byYou: false,
            }
          : null,
      };
    } catch (error) {
      tabs[collection] = {
        ...emptyTab(),
        loading: false,
        error: getApiErrorMessage(error, "Suggestions could not be loaded"),
      };
    }
  }

  function toState(bucket: AbsorbSuggestionBucket): BucketState {
    return {
      key: bucket.key,
      confidence: bucket.confidence,
      suggestedId: bucket.target_id,
      members: bucket.members,
      keepId: bucket.target_id,
      // Nothing is ticked at first: the person ticks each duplicate.
      picked: [],
      busy: false,
      error: "",
    };
  }

  function groupId(collection: AbsorbSuggestionCollection, key: string) {
    return `bucket-${collection}-${key}`;
  }

  function rows(collection: AbsorbSuggestionCollection, bucket: BucketState): MergeRow[] {
    return bucket.members.map((member) => ({
      id: member.id,
      label: member.name,
      detail: member.alias ? `(${member.alias})` : undefined,
      count: member.references,
      href: absorbSuggestionRecordHref(collection, member.id),
      locked: member.proposal ? "In a proposal" : undefined,
    }));
  }

  function subtitle(bucket: BucketState) {
    const keep = bucket.members.find((member) => member.id === bucket.keepId);
    const moved = bucket.picked.reduce(
      (sum, id) => sum + (bucket.members.find((member) => member.id === id)?.references ?? 0),
      0,
    );
    const match = `${Math.round(bucket.confidence * 100)}% match`;
    if (bucket.picked.length === 0 || !keep) return match;
    return `${match} · ${bucket.picked.length} to absorb into ${keep.name}, ${moved} ${movedLabel} move`;
  }

  // Focus goes to the bucket after index, or the one before when it was last.
  async function focusAfter(collection: AbsorbSuggestionCollection, nextKey: string | undefined) {
    await tick();
    if (nextKey) document.getElementById(groupId(collection, nextKey))?.focus();
  }

  function neighbourKey(collection: AbsorbSuggestionCollection, key: string) {
    const list = tabs[collection].buckets;
    const index = list.findIndex((bucket) => bucket.key === key);
    return list[index + 1]?.key ?? list[index - 1]?.key;
  }

  async function absorb(collection: AbsorbSuggestionCollection, bucket: BucketState) {
    bucket.busy = true;
    bucket.error = "";
    const picked = [...bucket.picked];
    try {
      await absorbIntoTarget(collection, bucket.keepId, picked);
      const keep = bucket.members.find((member) => member.id === bucket.keepId);
      tabs[collection].pending = {
        bucketKey: bucket.key,
        targetName: keep?.name ?? bucket.keepId,
        absorbed: bucket.members.filter((m) => picked.includes(m.id)).map((m) => m.name),
        created: "",
        byYou: true,
      };
      openModal();
    } catch (error) {
      bucket.error = getApiErrorMessage(error, "The records could not be absorbed");
    } finally {
      bucket.busy = false;
    }
  }

  async function dismiss(collection: AbsorbSuggestionCollection, bucket: BucketState) {
    bucket.busy = true;
    bucket.error = "";
    try {
      await dismissAbsorbSuggestion(
        collection,
        bucket.members.map((member) => member.id),
      );
      const next = neighbourKey(collection, bucket.key);
      tabs[collection].buckets = tabs[collection].buckets.filter((b) => b.key !== bucket.key);
      tabs[collection].dismissed += 1;
      focusAfter(collection, next);
    } catch (error) {
      bucket.error = getApiErrorMessage(error, "The bucket could not be dismissed");
      bucket.busy = false;
    }
  }

  function openModal() {
    modalRun += 1;
    modalOpen = true;
  }

  function resolved(outcome: "committed" | "undone") {
    const collection = tab;
    const pending = tabs[collection].pending;
    modalOpen = false;
    tabs[collection].pending = null;
    const bucket = tabs[collection].buckets.find((b) => b.key === pending?.bucketKey);
    if (!bucket) return;
    if (outcome === "undone") {
      // The bucket is as it was, with its ticks.
      focusAfter(collection, bucket.key);
      return;
    }
    const next = neighbourKey(collection, bucket.key);
    const absorbed = new Set(bucket.picked);
    const moved = bucket.members
      .filter((member) => absorbed.has(member.id))
      .reduce((sum, member) => sum + member.references, 0);
    bucket.members = bucket.members
      .filter((member) => !absorbed.has(member.id))
      .map((member) =>
        member.id === bucket.keepId ? { ...member, references: member.references + moved } : member,
      );
    bucket.picked = [];
    // A bucket keeps the records left unticked; one record alone is done.
    if (bucket.members.length < 2) {
      tabs[collection].buckets = tabs[collection].buckets.filter((b) => b.key !== bucket.key);
    }
    focusAfter(collection, next);
  }

  function selectTab(collection: AbsorbSuggestionCollection) {
    modalOpen = false;
    navigate(collection, data.level);
  }

  function tabLabel(collection: AbsorbSuggestionCollection) {
    const name = collection === "clients" ? "Clients" : "Vendors";
    const state = tabs[collection];
    return state.loading || state.error ? name : `${name} (${state.buckets.length})`;
  }
</script>

<svelte:head>
  <title>Absorb Suggestions</title>
</svelte:head>

<div class="flex flex-col">
  <DSTabBar
    tabs={absorbSuggestionCollections.map((collection) => ({
      label: tabLabel(collection),
      href: pageHref(collection, data.level),
      active: collection === tab,
      onselect: () => selectTab(collection),
    }))}
  >
    {#snippet rightContent()}
      <div class="flex items-center gap-2">
        <DSToggle
          bind:value={levelChoice}
          options={absorbSuggestionLevels}
          ariaLabel="How close the names must be"
        />
        <DSHelp doc="absorb-suggestions#find-duplicates" label="Help" />
      </div>
    {/snippet}
  </DSTabBar>

  <div class="flex flex-col gap-3 p-2">
    {#if !editingEnabled}
      <DsEditingDisabledBanner
        message="Absorb is turned off during a system transition. You can still mark buckets Not duplicates."
      />
    {/if}

    <!-- While the modal is open it shows the action; the notice is for later. -->
    {#if current.pending && !modalOpen}
      <DSCallout
        tone="warning"
        title={current.pending.byYou ? "Pending absorb by you" : "Pending absorb"}
      >
        {current.pending.targetName} kept, {current.pending.absorbed.join(", ")} absorbed{current
          .pending.created
          ? ` on ${shortDate(current.pending.created, true)}`
          : ""}. Commit or undo it before the next absorb.
        {#snippet actions()}
          <DsActionButton action={openModal} variant="primary">Review</DsActionButton>
        {/snippet}
      </DSCallout>
    {/if}

    {#if current.error}
      <DSCallout tone="danger">{current.error}</DSCallout>
    {:else if current.loading}
      <p>Finding duplicates...</p>
    {:else if current.buckets.length === 0}
      <DSCallout tone="info">
        No suggestions left at this level.{current.dismissed > 0
          ? ` ${current.dismissed} marked Not duplicates.`
          : ""}
      </DSCallout>
    {/if}

    {#if !current.loading}
      {#each current.buckets as bucket (bucket.key)}
        <DSMergeGroup
          id={groupId(tab, bucket.key)}
          title={`${bucket.members.length} ${tab}`}
          subtitle={subtitle(bucket)}
          rows={rows(tab, bucket)}
          bind:keepId={bucket.keepId}
          bind:picked={bucket.picked}
          suggestedId={bucket.suggestedId}
          {countLabel}
          disabled={bucket.busy}
        >
          {#snippet actions()}
            <DsActionButton
              action={() => dismiss(tab, bucket)}
              variant="secondary"
              disabled={bucket.busy || current.pending?.bucketKey === bucket.key}
              >Not duplicates</DsActionButton
            >
            <DsActionButton
              action={() => absorb(tab, bucket)}
              variant="primary"
              disabled={bucket.picked.length === 0 ||
                current.pending !== null ||
                !editingEnabled ||
                bucket.busy}>Absorb</DsActionButton
            >
          {/snippet}
          {#snippet notice()}
            {#if bucket.error}
              <DSCallout tone="danger">{bucket.error}</DSCallout>
            {/if}
          {/snippet}
        </DSMergeGroup>
      {/each}
    {/if}
  </div>
</div>

<DSModal bind:show={modalOpen} title="Pending absorb" size="md">
  {#snippet help()}
    <DSHelp doc="absorb-suggestions#commit-or-undo" label="Help" />
  {/snippet}
  {#key modalRun}
    <AbsorbList
      collectionName={tab}
      embedded
      onresolved={resolved}
      targetLabel={() => current.pending?.targetName}
    />
  {/key}
</DSModal>
