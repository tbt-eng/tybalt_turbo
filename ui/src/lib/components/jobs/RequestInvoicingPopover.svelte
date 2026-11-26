<script lang="ts">
  import MiniSearch from "minisearch";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSPicker from "$lib/components/DSPicker.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import { pb } from "$lib/pocketbase";
  import { resolve } from "$app/paths";
  import { untrack } from "svelte";
  import { workflowError } from "$lib/clientEditors";
  import { invoicingProfileLabel } from "$lib/utilities";
  import type {
    ClientContactsResponse,
    ClientInvoicingInformationResponse,
  } from "$lib/pocketbase-types";

  // Procurement sets, changes or clears the invoicing profile of a request that
  // waits for job creation. Only the profile changes: the request keeps its
  // status and its approval. A new profile is made in the client's own profile
  // editor, in a new tab, so this page stays open; Check again reloads the list.
  let {
    show = $bindable(),
    request,
    onchanged,
  }: {
    show: boolean;
    request: { id: string; version: number; client: string; invoicing_information?: string };
    onchanged: () => void | Promise<void>;
  } = $props();

  type Profile = ClientInvoicingInformationResponse & { label: string };

  // The parent mounts this dialog only while it is open, so each open starts
  // from the request's profile and loads the client's profiles again.
  let profile = $state(untrack(() => request.invoicing_information ?? ""));
  let profiles = $state<Profile[]>([]);
  let loading = $state(true);
  let loadError = $state(false);
  let submitting = $state(false);
  let error = $state<string | null>(null);
  let loadId = 0;

  const index = $derived.by(() => {
    const search = new MiniSearch<Profile>({
      fields: ["label"],
      storeFields: ["label"],
      searchOptions: { prefix: true, combineWith: "AND" },
    });
    search.addAll(profiles);
    return search;
  });
  const unchanged = $derived(profile === (request.invoicing_information ?? ""));

  async function load() {
    const id = ++loadId;
    loading = true;
    loadError = false;
    try {
      const filter = pb.filter("client={:client}", { client: request.client });
      const [contacts, rows] = await Promise.all([
        pb
          .collection("client_contacts")
          .getFullList<ClientContactsResponse>({ filter, requestKey: null }),
        pb
          .collection("client_invoicing_information")
          .getFullList<ClientInvoicingInformationResponse>({ filter, requestKey: null }),
      ]);
      if (id !== loadId) return;
      profiles = rows.map((row) => ({ ...row, label: invoicingProfileLabel(row, rows, contacts) }));
    } catch {
      if (id === loadId) loadError = true;
    } finally {
      if (id === loadId) loading = false;
    }
  }
  load();

  async function save() {
    submitting = true;
    error = null;
    try {
      await pb.send(`/api/job_requests/${request.id}/invoicing_profile`, {
        method: "POST",
        body: { version: request.version, invoicing_information: profile },
      });
      show = false;
      await onchanged();
    } catch (e) {
      error = workflowError(e);
    } finally {
      submitting = false;
    }
  }

  // The profile editor of this client. It returns to the client's page.
  const createHref = $derived(
    `${resolve("/clients/[cid]/invoicing/add", { cid: request.client })}?return_to=${encodeURIComponent(
      resolve("/clients/[cid]/details", { cid: request.client }),
    )}`,
  );
</script>

<DSPopover
  bind:show
  title="Invoicing profile"
  submitLabel="Save profile"
  {submitting}
  {error}
  submitDisabled={loading || loadError || unchanged}
  onSubmit={save}
>
  {#if loading}
    <p role="status">Loading invoicing profiles…</p>
  {:else if loadError}
    <p role="alert">Could not load this client's invoicing profiles.</p>
  {:else}
    <DSPicker
      bind:value={profile}
      items={profiles}
      {index}
      empty
      errors={{}}
      fieldName="invoicing_information"
      uiName="Invoicing profile"
    >
      {#snippet option(row)}{row.label}{/snippet}
    </DSPicker>
  {/if}
  <p>The approval stays valid. The job takes this profile.</p>
  <div class="flex flex-wrap items-center gap-2">
    <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- query on a resolved path -->
    <a class="link" href={createHref} target="_blank" rel="noopener">Create an invoicing profile</a>
    <DSActionButton variant="plain" disabled={loading} action={load}>Check again</DSActionButton>
  </div>
</DSPopover>
