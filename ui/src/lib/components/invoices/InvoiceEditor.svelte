<script lang="ts">
  import { resolve } from "$app/paths";
  import { goto } from "$app/navigation";
  import { untrack } from "svelte";
  import { pb } from "$lib/pocketbase";
  import { clients } from "$lib/stores/clients";
  import { clientLabel } from "$lib/clientWorkspace";
  import { saveErrorMessage } from "$lib/utilities";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";

  // Creates or edits a draft client invoice. The admin team chooses the
  // number as today; the server checks its YYMMNNN form and that no other
  // invoice uses it. See descriptions/invoices.md.

  type InvoiceDraft = { client: string; number: string; date: string; total: number | string };

  let { item: initial, id = null }: { item: InvoiceDraft; id?: string | null } = $props();

  clients.init();

  let item = $state(untrack(() => ({ ...initial })));
  let errors = $state({} as Record<string, { message: string }>);
  let saveError = $state("");
  let saving = $state(false);

  // Fields without their own error line: the date input has none.
  const otherErrors = $derived(
    Object.entries(errors)
      .filter(([field]) => !["client", "number", "total"].includes(field))
      .map(([, error]) => error.message),
  );

  async function save(event: Event) {
    event.preventDefault();
    saving = true;
    saveError = "";
    const fields = { ...item, total: item.total === "" ? 0 : Number(item.total) };
    try {
      const saved = id
        ? await pb.collection("invoices").update(id, fields)
        : await pb.collection("invoices").create(fields);
      errors = {};
      goto(resolve(`/invoices/${saved.id}/details`));
    } catch (error: unknown) {
      const data = (error as { response?: { data?: Record<string, { message: string }> } })
        ?.response?.data;
      errors = data ?? {};
      if (!data || Object.keys(data).length === 0) {
        saveError = saveErrorMessage(error, "The invoice could not be saved.");
      }
    } finally {
      saving = false;
    }
  }
</script>

<div class="flex w-full flex-col gap-4 p-2">
  <DSPageHeader
    title={id ? `Edit invoice ${initial.number || "draft"}` : "New invoice"}
    subtitle="Take the number from the Teams group as usual. Turbo checks that no other invoice uses it."
  />

  <form class="flex w-full flex-col gap-2" onsubmit={save}>
    {#if $clients.index !== null}
      <DsAutoComplete
        bind:value={item.client}
        index={$clients.index}
        {errors}
        fieldName="client"
        uiName="Client"
      >
        {#snippet resultTemplate(client)}{clientLabel({
            name: client.name,
            alias: client.alias,
          })}{/snippet}
      </DsAutoComplete>
    {/if}

    <DSTextInput
      bind:value={item.number}
      {errors}
      fieldName="number"
      uiName="Number"
      placeholder="YYMMNNN, e.g. 2609001"
    />

    <span class="flex w-full items-center gap-2">
      <label for="invoice-date">Date</label>
      <DsDateInput
        id="invoice-date"
        name="date"
        class="flex-1 md:flex-none"
        bind:value={item.date}
      />
    </span>

    <DSTextInput
      bind:value={item.total}
      type="number"
      step={0.01}
      min={0}
      {errors}
      fieldName="total"
      uiName="Total (CAD)"
    />

    {#if otherErrors.length > 0 || saveError}
      <DSCallout tone="danger">
        {[...otherErrors, saveError].filter(Boolean).join(" ")}
      </DSCallout>
    {/if}

    <span class="flex w-full gap-2">
      <DSActionButton type="submit" variant="primary" loading={saving}>Save draft</DSActionButton>
      <DSActionButton
        variant="secondary"
        action={id ? `/invoices/${id}/details` : "/invoices/list"}
      >
        Cancel
      </DSActionButton>
    </span>
  </form>
</div>
