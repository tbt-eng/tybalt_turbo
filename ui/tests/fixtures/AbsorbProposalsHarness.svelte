<script lang="ts">
  import AbsorbEditor from "$lib/components/AbsorbEditor.svelte";
  import AbsorbProposalQueue from "$lib/components/AbsorbProposalQueue.svelte";
  import type { AbsorbProposalRow } from "$lib/absorbProposals";
  import { globalStore } from "$lib/stores/global";

  // ?view=propose, review, queue (with &claims=absorb for a holder).
  const params = new URLSearchParams(location.search);
  const view = params.get("view") ?? "propose";
  // The test stubs the global store with a writable one.
  (globalStore as unknown as { set: (state: { claims: string[] }) => void }).set({
    claims: params.get("claims")?.split(",") ?? [],
  });

  const clients = [
    { id: "clienta00000001", name: "Acme Ltd" },
    { id: "clientb00000001", name: "Acme Limited" },
    { id: "clientc00000001", name: "Other Co" },
  ];
  const proposal = (over: Partial<AbsorbProposalRow>): AbsorbProposalRow => ({
    id: "proposal0000001",
    collection_name: "clients",
    target_id: "clienta00000001",
    target_name: "Acme Ltd",
    client_id: "",
    client_name: "",
    requester: "requester000001",
    requester_name: "Robin Requester",
    reason: "Entered twice",
    status: "pending",
    decided_at: "",
    decided_by_name: "",
    decision_comment: "",
    created: "2026-11-20 15:00:00.000Z",
    in_progress: false,
    records: [{ record_id: "clientb00000001", record_name: "Acme Limited" }],
    ...over,
  });
  const queue = [
    proposal({}),
    proposal({
      id: "proposal0000002",
      requester: "someoneelse0001",
      requester_name: "Sam Other",
      target_name: "Other Co",
      reason: "Same vendor",
    }),
  ];
  (window as unknown as { reloads: number }).reloads = 0;
</script>

{#if view === "queue"}
  <AbsorbProposalQueue
    status="pending"
    items={queue}
    viewerId="requester000001"
    onchange={async () => {
      (window as unknown as { reloads: number }).reloads += 1;
    }}
  />
{:else}
  <AbsorbEditor
    collectionName="clients"
    targetRecordId="clienta00000001"
    availableRecords={clients as never[]}
    propose={view === "propose"}
    proposal={view === "review" ? proposal({}) : null}
  >
    {#snippet recordSnippet(item: { name: string })}
      {item.name}
    {/snippet}
  </AbsorbEditor>
{/if}
