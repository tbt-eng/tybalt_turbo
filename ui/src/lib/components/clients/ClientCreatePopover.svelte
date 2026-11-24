<script lang="ts">
  import DSRequiredMark from "$lib/components/DSRequiredMark.svelte";
  // A small form that creates a client from its name alone, for a form that
  // picks a client (as the client or as the job owner). Creating a client
  // needs the job claim, so the caller shows its button to job holders only.
  import { pb } from "$lib/pocketbase";
  import DSPopover from "$lib/components/DSPopover.svelte";

  let {
    show = $bindable(),
    onCreated,
  }: {
    show: boolean;
    onCreated?: (client: { id: string; name: string }) => void;
  } = $props();

  let name = $state("");
  let submitting = $state(false);
  let error = $state<string | null>(null);
  let wasOpen = false;

  $effect(() => {
    if (show && !wasOpen) {
      name = "";
      error = null;
    }
    wasOpen = show;
  });

  async function createClient() {
    if (!name.trim()) {
      error = "Name is required";
      return;
    }
    submitting = true;
    error = null;
    try {
      // The clients store takes the new client from its realtime subscription.
      const created = await pb.collection("clients").create({ name: name.trim() });
      onCreated?.({ id: created.id, name: created.name });
      show = false;
    } catch (e) {
      const failure = e as { message?: string; data?: { data?: { name?: { message?: string } } } };
      error = failure?.data?.data?.name?.message ?? failure?.message ?? "Failed to create client";
    } finally {
      submitting = false;
    }
  }
</script>

<DSPopover
  bind:show
  title="Add New Client"
  {error}
  {submitting}
  submitLabel="Create Client"
  onSubmit={createClient}
  onCancel={() => (show = false)}
>
  <div class="flex flex-col gap-1">
    <label class="text-sm font-semibold" for="new_client_name">Name<DSRequiredMark /></label>
    <input
      id="new_client_name"
      aria-required="true"
      type="text"
      class="rounded-sm border border-neutral-300 px-2 py-1"
      placeholder="Client name"
      bind:value={name}
      disabled={submitting}
    />
  </div>
</DSPopover>
