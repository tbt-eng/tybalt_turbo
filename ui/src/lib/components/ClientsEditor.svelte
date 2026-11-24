<script lang="ts">
  import DSRequiredLegend from "$lib/components/DSRequiredLegend.svelte";
  import { pb } from "$lib/pocketbase";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import { beforeNavigate, goto } from "$app/navigation";
  import type { ClientsPageData } from "$lib/svelte-types";
  import DsActionButton from "./DSActionButton.svelte";
  import { globalStore } from "$lib/stores/global";
  import PostalFields from "$lib/components/clients/PostalFields.svelte";
  import { editableFields, clientFields, type FieldErrors } from "$lib/clientEditors";
  import { workflowError } from "$lib/clientEditors";
  import { untrack } from "svelte";

  let {
    data,
    onSaved,
    onCancel,
  }: {
    data: ClientsPageData;
    onSaved?: (id: string) => void | Promise<void>;
    onCancel?: () => void;
  } = $props();
  let errors = $state<FieldErrors>({});
  let item = $state(untrack(() => editableFields(data.item, clientFields)));
  let saved = $state(untrack(() => JSON.stringify(item)));
  let busy = $state(false);
  let savedForNavigation = false;
  let error = $state("");
  const dirty = $derived(JSON.stringify(item) !== saved);
  beforeNavigate(({ cancel }) => {
    if (busy && !savedForNavigation) {
      cancel();
      return;
    }
    if (dirty && !window.confirm("Discard unsaved client changes?")) cancel();
  });
  function cancelEdit() {
    if (dirty && !window.confirm("Discard unsaved client changes?")) return;
    saved = JSON.stringify(item);
    if (onCancel) onCancel();
    else goto(data.id ? `/clients/${data.id}/details` : "/clients/list");
  }
  async function save(event: Event) {
    event.preventDefault();
    if (busy || !$globalStore.claims.includes("job")) return;
    errors = {};
    if (!item.name.trim()) errors.name = { message: "This field is required." };
    if (Object.keys(errors).length) return;
    busy = true;
    savedForNavigation = false;
    error = "";
    try {
      const record =
        data.editing && data.id
          ? await pb.collection("clients").update(data.id, item)
          : await pb.collection("clients").create(item);
      saved = JSON.stringify(item);
      savedForNavigation = true;
      if (onSaved) await onSaved(record.id);
      else await goto(`/clients/${record.id}/details`);
    } catch (e: any) {
      errors = e?.data?.data || {};
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
</script>

<form
  class="flex w-full min-w-0 flex-col gap-4 p-2 [&_input]:min-w-0 [&_select]:min-w-0"
  onsubmit={save}
>
  <h2 class="text-xl font-bold">{data.editing ? "Edit client information" : "Create client"}</h2>
  <DSRequiredLegend />
  {#if error}<p role="alert" class="text-red-700">{error}</p>{/if}
  <fieldset
    disabled={busy || !$globalStore.claims.includes("job")}
    class="flex min-w-0 flex-col gap-3"
  >
    <DsTextInput
      bind:value={item.name}
      {errors}
      fieldName="name"
      uiName="Name"
      required
      placeholder="Official or numbered company name"
    />
    <DsTextInput
      bind:value={item.alias}
      {errors}
      fieldName="alias"
      uiName="Alias (optional)"
      placeholder="Name the client is commonly known by"
    />
    <DsTextInput bind:value={item.phone} {errors} fieldName="phone" uiName="Phone (optional)" />
    <PostalFields bind:item {errors} />
  </fieldset>
  <div class="flex gap-2">
    <DsActionButton type="submit" loading={busy} disabled={!$globalStore.claims.includes("job")}
      >Save client</DsActionButton
    >
    <DsActionButton action={cancelEdit} disabled={busy}>Cancel</DsActionButton>
  </div>
</form>
