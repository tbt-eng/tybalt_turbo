<script lang="ts">
  import DSRequiredLegend from "$lib/components/DSRequiredLegend.svelte";
  import { beforeNavigate } from "$app/navigation";
  import { untrack } from "svelte";
  import MiniSearch from "minisearch";
  import { pb } from "$lib/pocketbase";
  import { ClientResponseError } from "pocketbase";
  import { globalStore } from "$lib/stores/global";
  import DsList from "$lib/components/DSList.svelte";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSModal from "$lib/components/DSModal.svelte";
  import DSCollapsible from "$lib/components/DSCollapsible.svelte";
  import DSPicker from "$lib/components/DSPicker.svelte";
  import DSPopover from "$lib/components/DSPopover.svelte";
  import { invoicingProfileLabel } from "$lib/utilities";
  import type {
    ClientContactsResponse,
    ClientInvoicingInformationResponse,
  } from "$lib/pocketbase-types";
  import ContactEditor from "./ContactEditor.svelte";
  import InvoicingInstructions from "./InvoicingInstructions.svelte";
  import {
    addressText,
    billingName,
    canCreateInvoicingProfile,
    contactName,
    countLabel,
  } from "$lib/clientWorkspace";
  import {
    editableFields,
    invoicingFields,
    workflowError,
    type FieldErrors,
    type WorkflowRecord,
  } from "$lib/clientEditors";

  let {
    client,
    contacts: initialContacts,
    record = {},
    projectContact = "",
    onSaved,
    onCancel,
    onDeleted,
    onReplaced,
  }: {
    client: WorkflowRecord;
    contacts: WorkflowRecord[];
    record?: WorkflowRecord;
    projectContact?: string;
    onSaved: (
      record: WorkflowRecord,
      mode?: "create" | "update" | "copy" | "reuse",
    ) => void | Promise<void>;
    onCancel: () => void;
    // Given, the editor offers to delete an unused profile and to replace it.
    onDeleted?: () => void | Promise<void>;
    // Receives the profile that took over this profile's jobs and requests.
    onReplaced?: (replacement: WorkflowRecord) => void | Promise<void>;
  } = $props();
  // A blank billing_name bills the official name, as does "name".
  let item = $state(
    untrack(() => {
      const fields = editableFields(record, invoicingFields);
      if (fields.billing_name !== "alias") fields.billing_name = "";
      return fields;
    }),
  );
  let saved = $state(untrack(() => JSON.stringify(item)));
  let contacts = $state(untrack(() => [...initialContacts]));
  let contactDraft = $state<WorkflowRecord | null>(null);
  let contactDirty = $state(false);
  let contactBusy = $state(false);
  let savedForNavigation = false;
  let errors = $state<FieldErrors>({});
  let error = $state("");
  let busy = $state(false);
  let completed = $state<{
    record: WorkflowRecord;
    mode: "create" | "update" | "copy" | "reuse";
    unchanged: boolean;
  } | null>(null);
  let decision = $state<{ jobs: number; requests: number } | null>(null);
  const helpId = $props.id();
  const explanations = {
    choose: "Choose how to save your changes.",
    create: "Create a profile with this contact and these instructions.",
    update: "Apply these changes everywhere this profile is used.",
    copy: "Keep the original profile unchanged and create a new profile from these fields.",
    cancel: "Return to editing without saving. Your changes stay in the form.",
    leave: "Leave this editor without saving these changes.",
  };
  // A new profile has one way to save; an existing one offers a choice.
  const defaultHelp = untrack(() => (record.id ? explanations.choose : explanations.create));
  let actionHelp = $state(defaultHelp);
  // Show the explanation for the action under the pointer or keyboard focus.
  const explain = (text: string) => () => (actionHelp = text);
  function cancelDecision() {
    if (busy) return;
    decision = null;
    actionHelp = defaultHelp;
  }
  // Holders of time can create profiles; only holders of job change shared ones.
  const canManage = $derived($globalStore.claims.includes("job"));
  const canCreate = $derived(canCreateInvoicingProfile($globalStore.claims));
  const canSave = $derived(record.id ? canManage : canCreate);
  const matchingProfiles = $derived(
    ((client.invoicing_profiles || []) as (WorkflowRecord & { id: string })[]).filter(
      (profile) =>
        profile.client === client.id &&
        profile.contact === item.contact &&
        profile.id !== record.id,
    ),
  );
  const dirty = $derived(JSON.stringify(item) !== saved || contactDirty);
  const recipient = $derived(contacts.find((contact) => contact.id === item.contact));
  const clientAlias = $derived(String(client.alias || "").trim());
  const contactIndex = $derived.by(() => {
    const index = new MiniSearch<WorkflowRecord>({
      fields: ["given_name", "surname", "email", "phone"],
      storeFields: ["given_name", "surname", "email"],
      searchOptions: { combineWith: "AND" },
    });
    index.addAll(contacts);
    return index;
  });
  beforeNavigate(({ cancel }) => {
    if (contactBusy || decision || (busy && !savedForNavigation)) {
      cancel();
      return;
    }
    if (dirty && !window.confirm("Discard unsaved invoicing changes?")) cancel();
  });
  // Contacts need not have an email address.
  function contactLabel(contact: WorkflowRecord) {
    return [contactName(contact), contact.email].filter(Boolean).join(" — ");
  }
  function cancelEdit() {
    if (dirty && !window.confirm("Discard unsaved invoicing changes?")) return;
    saved = JSON.stringify(item);
    contactDirty = false;
    onCancel();
  }
  function contactSaved(contact: WorkflowRecord) {
    contacts = [...contacts.filter((row) => row.id !== contact.id), contact];
    item.contact = contact.id;
    contactDirty = false;
    contactDraft = null;
  }
  function validate() {
    errors = {};
    if (!contacts.some((contact) => contact.id === item.contact))
      errors.contact = { message: "Select a contact for this client." };
    if (item.name.length > 120) errors.name = { message: "Use 120 characters or fewer." };
    return Object.keys(errors).length === 0;
  }
  async function returnToWork() {
    if (!completed) return;
    error = "";
    try {
      await onSaved(completed.record, completed.mode);
    } catch {
      error =
        "Your profile is ready, but the return page could not open. Retry the return; this will not save again.";
    }
  }
  async function retryReturn() {
    if (busy) return;
    busy = true;
    try {
      await returnToWork();
    } finally {
      busy = false;
    }
  }
  async function finish(
    result: WorkflowRecord,
    mode: "create" | "update" | "copy" | "reuse",
    unchanged = false,
  ) {
    saved = JSON.stringify(item);
    decision = null;
    completed = { record: result, mode, unchanged };
    savedForNavigation = true;
    await returnToWork();
  }
  async function reuseProfile(profile: WorkflowRecord) {
    if (
      busy ||
      completed ||
      record.id ||
      !canCreate ||
      !matchingProfiles.some((row) => row.id === profile.id)
    )
      return;
    if (
      [item.name, item.billing_name, item.fax, item.invoicing_instructions].some((value) =>
        String(value).trim(),
      ) &&
      !window.confirm("Use this existing profile and discard the unsaved profile fields?")
    )
      return;
    busy = true;
    error = "";
    try {
      await finish(profile, "reuse");
    } finally {
      busy = false;
    }
  }
  async function persist(mode: "create" | "update" | "copy") {
    const body = { ...editableFields(item, invoicingFields), client: client.id };
    const result =
      mode === "update"
        ? await pb.collection("client_invoicing_information").update(record.id, body)
        : await pb.collection("client_invoicing_information").create(body);
    await finish(result, mode);
  }
  async function save(event?: Event, copy = false) {
    event?.preventDefault();
    if (busy || completed || decision || !canSave || (copy && !canManage) || !validate()) return;
    busy = true;
    savedForNavigation = false;
    error = "";
    try {
      if (record.id && !copy && JSON.stringify(item) === saved) {
        await finish(record, "update", true);
        return;
      }
      if (record.id && !copy) {
        // Refresh usage before confirmation; jobs may have linked this profile since it opened.
        const current = await pb.send<WorkflowRecord>(`/api/clients/${client.id}`, {
          method: "GET",
          requestKey: null,
        });
        const profile = current.invoicing_profiles?.find(
          (profile: WorkflowRecord) => profile.id === record.id,
        );
        if (!profile)
          throw new Error("This invoicing profile is no longer available. Reload the page.");
        const jobs = Number(profile.job_count || 0);
        const requests = Number(profile.request_count || 0);
        if (jobs > 1) {
          decision = { jobs, requests };
          actionHelp = explanations.choose;
          return;
        }
      }
      await persist(record.id ? (copy ? "copy" : "update") : "create");
    } catch (e: any) {
      errors = e?.data?.data || {};
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
  // A profile in use can't be deleted; replacing it moves its jobs and
  // requests to another profile of the client, then deletes it.
  const inUse = $derived(Number(record.job_count || 0) + Number(record.request_count || 0) > 0);
  const otherProfiles = $derived(
    ((client.invoicing_profiles || []) as ClientInvoicingInformationResponse[]).filter(
      (profile) => profile.client === client.id && profile.id !== record.id,
    ),
  );
  const profileIndex = $derived.by(() => {
    const index = new MiniSearch<ClientInvoicingInformationResponse>({
      fields: ["name", "invoicing_instructions", "fax"],
      storeFields: ["name", "contact", "invoicing_instructions", "created"],
      searchOptions: { combineWith: "AND", prefix: true },
    });
    index.addAll(otherProfiles);
    return index;
  });
  function profileLabel(profile: ClientInvoicingInformationResponse) {
    return invoicingProfileLabel(
      profile,
      otherProfiles,
      contacts as unknown as ClientContactsResponse[],
    );
  }
  let replacing = $state(false);
  let replacement = $state("");
  let replaceErrors = $state<FieldErrors>({});
  let replaceError = $state<string | null>(null);
  // Leave the editor without the unsaved-changes prompt: the profile is gone.
  async function leaveAfter(action: () => void | Promise<void>) {
    saved = JSON.stringify(item);
    contactDirty = false;
    savedForNavigation = true;
    await action();
  }
  async function deleteProfile() {
    if (busy || inUse || !canManage || !record.id) return;
    if (!window.confirm("Delete this invoicing profile?")) return;
    busy = true;
    savedForNavigation = false;
    error = "";
    try {
      await pb.collection("client_invoicing_information").delete(record.id);
      await leaveAfter(() => onDeleted?.());
    } catch (e) {
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
  function openReplace() {
    if (busy || !canManage) return;
    replacement = "";
    replaceErrors = {};
    replaceError = null;
    replacing = true;
  }
  async function replaceProfile() {
    if (busy || !canManage || !record.id) return;
    if (!replacement) {
      replaceErrors = { replacement: { message: "Select the invoicing profile to use instead." } };
      return;
    }
    busy = true;
    savedForNavigation = false;
    replaceErrors = {};
    replaceError = null;
    try {
      await pb.send(`/api/clients/${client.id}/invoicing_profiles/${record.id}/replace`, {
        method: "POST",
        body: { replacement },
      });
      const chosen = otherProfiles.find((profile) => profile.id === replacement) ?? {
        id: replacement,
      };
      await leaveAfter(() => onReplaced?.(chosen));
    } catch (e: unknown) {
      const response = e instanceof ClientResponseError ? e.data : {};
      replaceErrors = (response?.data as FieldErrors | undefined) ?? {};
      replaceError = response?.message ?? "The profile could not be replaced. Try again.";
    } finally {
      busy = false;
    }
  }
  async function saveDecision(mode: "update" | "copy") {
    if (busy || !decision || !canManage) return;
    busy = true;
    error = "";
    try {
      await persist(mode);
    } catch (e: any) {
      errors = e?.data?.data || {};
      error = workflowError(e);
    } finally {
      busy = false;
    }
  }
</script>

{#if completed}
  <div class="space-y-3">
    <p role="status">
      {completed.mode === "reuse"
        ? "Using the existing invoicing profile."
        : completed.unchanged
          ? "No changes to save."
          : "Invoicing profile saved."}
    </p>
    {#if error}<p role="alert" class="text-red-700">{error}</p>{/if}
    <DsActionButton action={retryReturn} disabled={busy}>Return to your work</DsActionButton>
  </div>
{:else if contactDraft}
  <p class="mb-4 text-sm text-neutral-600">Save the contact to return to your invoicing profile.</p>
  {#key contactDraft.id || "new"}
    <ContactEditor
      clientId={client.id}
      record={contactDraft}
      guardNavigation={false}
      bind:dirty={contactDirty}
      bind:busy={contactBusy}
      onSaved={contactSaved}
      onCancel={() => {
        contactDraft = null;
        contactDirty = false;
      }}
    />
  {/key}
{:else}
  <form class="flex min-w-0 flex-col gap-4 [&_input]:min-w-0 [&_select]:min-w-0" onsubmit={save}>
    <h2 class="text-xl font-bold">
      {record.id ? "Edit invoicing profile" : "Add invoicing profile"}
    </h2>
    {#if record.id}<p class="text-sm text-neutral-600">
        Used by {countLabel(record.job_count || 0, "job")}{record.request_count > 0
          ? ` and ${countLabel(record.request_count, "request")}`
          : ""}. Changes apply wherever this profile is used.
      </p>{/if}
    <DSRequiredLegend />
    {#if error}<p role="alert" class="text-red-700">{error}</p>{/if}
    <fieldset disabled={busy || !!decision || !canSave} class="flex min-w-0 flex-col gap-3">
      <DsTextInput
        bind:value={item.name}
        {errors}
        fieldName="name"
        uiName="Profile name"
        placeholder="For example, Capital projects"
      />
      {#if contacts.length > 10}
        <DsAutoComplete
          bind:value={item.contact}
          index={contactIndex}
          {errors}
          fieldName="contact"
          uiName="Send invoices to"
          required
        >
          {#snippet resultTemplate(row)}{contactLabel(row as WorkflowRecord)}{/snippet}
        </DsAutoComplete>
      {:else}
        <DsSelector
          bind:value={item.contact}
          items={[
            { id: "", label: "Select a contact" },
            ...contacts.map((contact) => ({
              id: String(contact.id),
              label: contactLabel(contact),
            })),
          ]}
          {errors}
          fieldName="contact"
          uiName="Send invoices to"
          required
        >
          {#snippet optionTemplate(row)}{row.label}{/snippet}
        </DsSelector>
      {/if}
      <div class="flex flex-wrap gap-2">
        {#if canManage}
          <DsActionButton
            action={() => {
              contactDraft = {};
            }}>Add contact</DsActionButton
          >
          {#if recipient}<DsActionButton
              action={() => {
                contactDraft = recipient;
              }}>Edit contact</DsActionButton
            >{/if}
        {/if}
        {#if projectContact && contacts.some((contact) => contact.id === projectContact)}
          <DsActionButton
            action={() => {
              item.contact = projectContact;
            }}>Same as project contact</DsActionButton
          >
        {/if}
      </div>
      {#if !record.id && recipient && matchingProfiles.length}
        <!-- svelte-ignore a11y_no_noninteractive_element_interactions (Enter in profile search must not save a new profile.) -->
        <section
          aria-label="Existing profiles for this contact"
          onkeydown={(event) => {
            if (event.key === "Enter" && event.target instanceof HTMLInputElement)
              event.preventDefault();
          }}
          class="space-y-2 rounded-sm border border-neutral-300 p-3"
        >
          <h3 class="font-semibold">Existing profiles for this contact</h3>
          <p class="text-sm text-neutral-600">
            Use an existing profile, or enter different instructions below to create another.
          </p>
          {#key item.contact}
            <DsList
              items={matchingProfiles}
              search={matchingProfiles.length > 10}
              pagination={matchingProfiles.length > 10}
              perPage={10}
              searchPlaceholder="Find an existing profile"
              searchText={(profile) =>
                [profile.name, profile.invoicing_instructions, profile.fax]
                  .filter(Boolean)
                  .join(" ")}
            >
              {#snippet headline(profile)}{profile.name || contactName(recipient)}{/snippet}
              {#snippet line1(profile)}
                <InvoicingInstructions text={profile.invoicing_instructions || ""} />
                {#if profile.fax}<p class="text-sm">Fax: {profile.fax}</p>{/if}
              {/snippet}
              {#snippet actions(profile)}<DsActionButton
                  action={() => reuseProfile(profile)}
                  disabled={busy || !canCreate}>Use this profile</DsActionButton
                >{/snippet}
            </DsList>
          {/key}
        </section>
      {/if}
      {#if clientAlias || item.billing_name === "alias"}
        <fieldset class="flex flex-col gap-1">
          <legend class="mb-1">Name on Project Authorizations and invoices</legend>
          <label class="flex items-center gap-2">
            <input type="radio" name="billing_name" value="" bind:group={item.billing_name} />
            {client.name}
            <span class="text-sm text-neutral-600">(name)</span>
          </label>
          <label class="flex items-center gap-2">
            <input
              type="radio"
              name="billing_name"
              value="alias"
              bind:group={item.billing_name}
              disabled={!clientAlias}
            />
            {clientAlias || "No alias"}
            <span class="text-sm text-neutral-600">(alias)</span>
          </label>
          {#if errors.billing_name}<p class="text-red-600">{errors.billing_name.message}</p>{/if}
        </fieldset>
      {:else}
        <p class="text-sm text-neutral-600">
          Project Authorizations and invoices use the client name. To use a different name, add an
          alias to the client.
        </p>
      {/if}
      {#if recipient}
        <DSCard title="Invoice address">
          <div class="text-sm">
            <p>{billingName(client as { name: string; alias?: string }, item)}</p>
            <p class="whitespace-pre-line">
              {addressText(recipient) || addressText(client) || "No postal address entered"}
            </p>
            {#if !addressText(recipient)}<p class="text-neutral-600">
                Uses the client address.
              </p>{/if}
          </div>
        </DSCard>
      {/if}
      <label class="flex flex-col gap-2" for="invoicing-instructions"
        >Invoice instructions
        <textarea
          id="invoicing-instructions"
          name="invoicing_instructions"
          class="rounded-sm border border-neutral-300 p-2"
          rows="4"
          bind:value={item.invoicing_instructions}></textarea>
      </label>
      {#if errors.invoicing_instructions}<p class="text-red-600">
          {errors.invoicing_instructions.message}
        </p>{/if}
      <DsTextInput bind:value={item.fax} {errors} fieldName="fax" uiName="Fax" />
    </fieldset>
    <div class="flex flex-wrap gap-2">
      <DsActionButton
        type="submit"
        disabled={busy || !!decision || !canSave}
        describedBy={helpId + "-form"}
        onexplain={explain(record.id ? explanations.update : explanations.create)}
        >Save invoicing profile</DsActionButton
      >
      {#if record.id && canManage}
        <DsActionButton
          action={() => save(undefined, true)}
          disabled={busy || !!decision || !canManage}
          color="neutral"
          describedBy={helpId + "-form"}
          onexplain={explain(explanations.copy)}>Save as new profile</DsActionButton
        >
      {/if}
      <DsActionButton
        action={cancelEdit}
        disabled={busy || !!decision}
        color="neutral"
        describedBy={helpId + "-form"}
        onexplain={explain(explanations.leave)}>Cancel</DsActionButton
      >
    </div>
    <p id={helpId + "-form"} class="min-h-10 text-sm text-neutral-600" aria-live="polite">
      {actionHelp}
    </p>
    {#if record.id && onDeleted && canManage}
      <DSCollapsible title="More actions">
        <div class="flex flex-wrap gap-2">
          {#if onReplaced && otherProfiles.length}
            <DsActionButton action={openReplace} disabled={busy || !!decision}
              >Replace with…</DsActionButton
            >
          {/if}
          <DsActionButton action={deleteProfile} disabled={busy || !!decision || inUse} color="red"
            >Delete invoicing profile</DsActionButton
          >
        </div>
        {#if inUse}<p class="text-sm text-neutral-600">
            This profile is in use and cannot be deleted.
            {#if onReplaced && otherProfiles.length}Use <strong>Replace with…</strong> to move its jobs
              and requests to another profile.{/if}
          </p>{/if}
      </DSCollapsible>
    {/if}
  </form>
{/if}

{#if replacing}
  <DSPopover
    bind:show={replacing}
    title="Replace invoicing profile"
    subtitle="Move the jobs and requests that use this profile to another profile of this client, then delete this profile."
    error={replaceError}
    submitting={busy}
    submitLabel="Replace profile"
    submitDisabled={!replacement}
    danger
    onSubmit={replaceProfile}
    onCancel={() => {
      if (!busy) replacing = false;
    }}
  >
    <DSPicker
      bind:value={replacement}
      items={otherProfiles}
      index={profileIndex}
      errors={replaceErrors}
      fieldName="replacement"
      uiName="Replace with"
      placeholder="Find an invoicing profile"
      disabled={busy}
      required
      empty
    >
      {#snippet option(profile)}{profileLabel(profile)}{/snippet}
    </DSPicker>
    <DSCallout tone="warning">
      <p>
        {countLabel(Number(record.job_count || 0), "job")} and {countLabel(
          Number(record.request_count || 0),
          "request",
        )} will move to the chosen profile. This profile is then deleted. This cannot be undone.
      </p>
    </DSCallout>
  </DSPopover>
{/if}

<DSModal
  show={!!decision}
  title="Save invoicing profile"
  dismissible={!busy}
  onclose={cancelDecision}
>
  {#if decision}
    <div class="space-y-2">
      <p>
        This profile is used by {countLabel(decision.jobs, "job")}{decision.requests > 0
          ? ` and ${countLabel(decision.requests, "request")}`
          : ""}.
      </p>
      <p>Update the shared profile, or save your changes as a new profile?</p>
      {#if error}<p role="alert" class="text-red-700">{error}</p>{/if}
      <p id={helpId + "-choice"} class="min-h-10 text-sm text-neutral-600" aria-live="polite">
        {actionHelp}
      </p>
    </div>
  {/if}
  {#snippet footer()}
    <DsActionButton
      variant="secondary"
      action={cancelDecision}
      disabled={busy}
      describedBy={helpId + "-choice"}
      onexplain={explain(explanations.cancel)}>Cancel</DsActionButton
    >
    <DsActionButton
      variant="secondary"
      action={() => saveDecision("copy")}
      disabled={busy}
      describedBy={helpId + "-choice"}
      onexplain={explain(explanations.copy)}>Save as new profile</DsActionButton
    >
    <DsActionButton
      variant="primary"
      action={() => saveDecision("update")}
      disabled={busy}
      describedBy={helpId + "-choice"}
      onexplain={explain(explanations.update)}>Update shared profile</DsActionButton
    >
  {/snippet}
</DSModal>
