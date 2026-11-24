<script lang="ts">
  import ClientWorkspaceShell from "$lib/components/ClientWorkspaceShell.svelte";
  import type { PageData } from "./$types";
  import AbsorbEditor from "$lib/components/AbsorbEditor.svelte";
  import type { ClientContactsResponse } from "$lib/pocketbase-types";
  import MiniSearch from "minisearch";
  import { page } from "$app/stores";
  import { globalStore } from "$lib/stores/global";
  import { canAbsorbDirectly } from "$lib/absorbProposals";
  import { pickerDropdownLimit } from "$lib/components/DSPicker.svelte";
  export let data: PageData;

  const availableRecords: ClientContactsResponse[] = data.contacts;
  let autoCompleteIndex: MiniSearch<ClientContactsResponse> | null = null;

  // A dropdown for a short list, a search for a longer one, as in DSPicker.
  if (availableRecords.length > pickerDropdownLimit) {
    autoCompleteIndex = new MiniSearch<ClientContactsResponse>({
      // include "id" to allow searching by id when a result is selected
      fields: ["id", "given_name", "surname", "email"],
      storeFields: ["id", "given_name", "surname", "email"],
      searchOptions: {
        combineWith: "AND",
      },
    });
    autoCompleteIndex.addAll(availableRecords);
  }
</script>

<ClientWorkspaceShell clientId={data.client.id} clientName={data.client.name}>
  <AbsorbEditor
    collectionName="client_contacts"
    returnTo={data.returnTo}
    targetRecordId={$page.params.kid!}
    {availableRecords}
    autoCompleteIndex={autoCompleteIndex as unknown as any}
    propose={!canAbsorbDirectly($globalStore.claims)}
    proposal={data.proposal}
  >
    {#snippet recordSnippet(item: ClientContactsResponse)}
      {item.given_name} {item.surname} — {item.email}
    {/snippet}
  </AbsorbEditor>
</ClientWorkspaceShell>
