<script lang="ts">
  import AbsorbEditor from "$lib/components/AbsorbEditor.svelte";
  import type { ClientsResponse } from "$lib/pocketbase-types";
  import { clients } from "$lib/stores/clients";
  import { clientLabel } from "$lib/clientWorkspace";
  import { page } from "$app/stores";
  import { globalStore } from "$lib/stores/global";
  import { canAbsorbDirectly } from "$lib/absorbProposals";
</script>

<AbsorbEditor
  collectionName="clients"
  targetRecordId={$page.params.cid ?? ""}
  availableRecords={$clients.items}
  autoCompleteIndex={$clients.index as unknown as any}
  propose={!canAbsorbDirectly($globalStore.claims)}
>
  {#snippet recordSnippet(item: ClientsResponse)}
    {clientLabel(item)}
  {/snippet}
</AbsorbEditor>
