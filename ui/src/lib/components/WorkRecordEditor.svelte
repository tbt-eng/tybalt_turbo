<script lang="ts">
  import { untrack } from "svelte";
  import { goto } from "$app/navigation";
  import { authStore } from "$lib/stores/auth";
  import { globalStore } from "$lib/stores/global";
  import { jobs } from "$lib/stores/jobs";
  import { profiles } from "$lib/stores/profiles";
  import { workRecordTypes } from "$lib/stores/workRecordTypes";
  import { pb } from "$lib/pocketbase";
  import DsActionButton from "$lib/components/DSActionButton.svelte";
  import DsAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import type {
    WorkRecordCopyDraftResponse,
    WorkRecordConsumableResponse,
    WorkRecordCreateRequest,
    WorkRecordDetailsResponse,
    WorkRecordSaveConsumableRequest,
    WorkRecordSaveRequest,
    WorkRecordSaveResponse,
  } from "$lib/svelte-types";
  import { refreshVisibleWorkRecords } from "$lib/stores/workRecords";
  import { DATE_INPUT_MIN, dateInputMaxMonthsAhead } from "$lib/utilities";

  jobs.init();
  profiles.init();
  workRecordTypes.init();

  type WorkRecordForm = {
    id?: string;
    uid: string;
    creator: string;
    date: string;
    job: string;
    type: string;
    location: string;
    work_description: string;
    report_to: string;
    hours_on_site: string | number;
    hours_travel_time: string | number;
    distance_travelled_km: string | number;
    field_book_number: string;
    field_book_page_number: string | number;
    sub_contractor: string;
    equipment: string;
    supplies: string;
    vehicle_type: string;
    is_passenger: boolean;
    company_vehicle_unit_number: string;
  };

  type EditorConsumable = WorkRecordConsumableResponse & {
    blocked: boolean;
  };

  const vehicleTypeOptions = [
    { id: "company", label: "Company" },
    { id: "personal", label: "Personal" },
  ];

  const assigneeModeOptions = [
    { id: "you", name: "For you" },
    { id: "someone_else", name: "For someone else" },
  ];
  const workRecordTypeOptions = $derived([
    { id: "", name: "Select a work record type" },
    ...$workRecordTypes.items.map((item) => ({ id: item.id, name: item.name })),
  ]);

  let {
    item,
    editing = false,
    id = null,
    detail = null,
    copyContext = null,
    cancelAction = "/time/work-records",
  }: {
    item: WorkRecordForm;
    editing?: boolean;
    id?: string | null;
    detail?: WorkRecordDetailsResponse | null;
    copyContext?: WorkRecordCopyDraftResponse | null;
    cancelAction?: string;
  } = $props();

  let form = $state(untrack(() => item));
  let errors = $state({} as Record<string, { message: string }>);
  let saving = $state(false);
  let loadingConsumables = $state(false);
  let consumables = $state([] as EditorConsumable[]);
  let assigneeMode = $state<"you" | "someone_else">("you");
  let consumableWarning = $state("");
  let consumablesRequestId = 0;

  const dateInputMax = dateInputMaxMonthsAhead(15);
  const currentUserId = $derived($authStore?.model?.id ?? "");
  const canAssignOthers = $derived($globalStore.user_claims_summary.claims.includes("work_record"));
  const canViewAll = $derived($globalStore.user_claims_summary.claims.includes("report"));
  const canEdit = $derived(detail?.can_edit ?? true);
  const isCopyMode = $derived(copyContext !== null);
  const showSaveButton = $derived(!editing || canEdit);
  const hasDefaultType = $derived(($globalStore.profile.default_work_record_type ?? "") !== "");
  const visibleConsumables = $derived(consumables.filter((item) => !item.blocked));
  const blockedConsumables = $derived(consumables.filter((item) => item.blocked));

  function buildEmptyConsumable(row: any): EditorConsumable {
    const expandedConsumable = Array.isArray(row.expand?.consumable)
      ? row.expand.consumable[0]
      : row.expand?.consumable;

    return {
      entry_id: "",
      consumable_id: row.consumable,
      consumable_name: expandedConsumable?.name ?? "",
      input_kind: expandedConsumable?.input_kind ?? "number",
      unit_label: row.unit_label ?? "",
      sort_order: row.sort_order ?? 0,
      quantity_number: 0,
      boolean_selected: false,
      blocked: false,
    };
  }

  function hasConsumableValue(item: EditorConsumable) {
    return item.input_kind === "boolean" ? item.boolean_selected : Number(item.quantity_number) > 0;
  }

  async function loadTypeConsumables(typeId: string) {
    if (typeId === "") {
      consumables = [];
      consumableWarning = "";
      return;
    }

    const requestId = ++consumablesRequestId;
    loadingConsumables = true;
    try {
      const rows = await pb.collection("work_records_type_consumables").getFullList<any>({
        filter: pb.filter("type = {:type} && active = true", { type: typeId }),
        expand: "consumable",
        sort: "sort_order,consumable",
        requestKey: null,
      });
      if (requestId !== consumablesRequestId) return;

      const nextVisible = rows.map(buildEmptyConsumable);
      const existingByConsumable = new Map(consumables.map((row) => [row.consumable_id, row]));
      const mergedVisible = nextVisible.map((row) => {
        const existing = existingByConsumable.get(row.consumable_id);
        if (!existing) return row;
        return {
          ...row,
          entry_id: existing.entry_id,
          quantity_number: existing.quantity_number,
          boolean_selected: existing.boolean_selected,
        };
      });

      const blocked = consumables
        .filter(
          (row) => !nextVisible.some((visible) => visible.consumable_id === row.consumable_id),
        )
        .map((row) => ({ ...row, blocked: true }));

      consumables = [...mergedVisible, ...blocked];
      consumableWarning =
        blocked.length > 0
          ? "Some existing consumables are not allowed for the selected type. Clear them before saving."
          : "";
    } catch (error) {
      if (requestId !== consumablesRequestId) return;
      console.error("loading work record consumables", error);
      consumableWarning = "Failed to load consumables for the selected type.";
    } finally {
      if (requestId !== consumablesRequestId) return;
      loadingConsumables = false;
    }
  }

  function buildPayload() {
    return {
      uid: resolvedUID(),
      creator: resolvedCreator(),
      date: form.date,
      job: form.job,
      type: form.type,
      location: form.location,
      work_description: form.work_description,
      report_to: form.report_to,
      hours_on_site: Number(form.hours_on_site || 0),
      hours_travel_time: Number(form.hours_travel_time || 0),
      distance_travelled_km: Number(form.distance_travelled_km || 0),
      field_book_number: form.field_book_number || "",
      field_book_page_number:
        `${form.field_book_page_number ?? ""}`.trim() === ""
          ? ""
          : Number(form.field_book_page_number),
      sub_contractor: form.sub_contractor || "",
      equipment: form.equipment || "",
      supplies: form.supplies || "",
      vehicle_type: form.vehicle_type,
      is_passenger: Boolean(form.is_passenger),
      company_vehicle_unit_number: form.company_vehicle_unit_number || "",
    };
  }

  function resolvedCreator() {
    return currentUserId;
  }

  function resolvedUID() {
    if (!canAssignOthers || assigneeMode === "you" || form.uid === "") {
      return currentUserId;
    }
    return form.uid;
  }

  function buildSaveConsumables(): WorkRecordSaveConsumableRequest[] {
    return consumables.map((consumable) => ({
      entry_id: consumable.entry_id,
      consumable_id: consumable.consumable_id,
      input_kind: consumable.input_kind,
      quantity_number: Number(consumable.quantity_number || 0),
      boolean_selected: Boolean(consumable.boolean_selected),
    }));
  }

  function buildSaveRequest(): WorkRecordSaveRequest {
    const payload = buildPayload();
    return {
      date: payload.date,
      job: payload.job,
      type: payload.type,
      location: payload.location,
      work_description: payload.work_description,
      report_to: payload.report_to,
      hours_on_site: payload.hours_on_site,
      hours_travel_time: payload.hours_travel_time,
      distance_travelled_km: payload.distance_travelled_km,
      field_book_number: payload.field_book_number,
      field_book_page_number: payload.field_book_page_number as number | "",
      sub_contractor: payload.sub_contractor,
      equipment: payload.equipment,
      supplies: payload.supplies,
      vehicle_type: payload.vehicle_type,
      is_passenger: payload.is_passenger,
      company_vehicle_unit_number: payload.company_vehicle_unit_number,
      consumables: buildSaveConsumables(),
    };
  }

  function buildCreateRequest(): WorkRecordCreateRequest {
    const payload = buildPayload();
    return {
      uid: payload.uid,
      date: payload.date,
      job: payload.job,
      type: payload.type,
      location: payload.location,
      work_description: payload.work_description,
      report_to: payload.report_to,
      hours_on_site: payload.hours_on_site,
      hours_travel_time: payload.hours_travel_time,
      distance_travelled_km: payload.distance_travelled_km,
      field_book_number: payload.field_book_number,
      field_book_page_number: payload.field_book_page_number as number | "",
      sub_contractor: payload.sub_contractor,
      equipment: payload.equipment,
      supplies: payload.supplies,
      vehicle_type: payload.vehicle_type,
      is_passenger: payload.is_passenger,
      company_vehicle_unit_number: payload.company_vehicle_unit_number,
      consumables: buildSaveConsumables().map((consumable) => ({
        ...consumable,
        entry_id: "",
      })),
      copy_from_work_record_id: copyContext?.copy_from_work_record_id,
    };
  }

  function applySaveErrors(error: any) {
    const nextErrors = error?.data?.data ?? {};
    if (Object.keys(nextErrors).length === 0) {
      errors = {
        global: {
          message: error?.data?.message ?? error?.message ?? "Failed to save work record",
        },
      };
      return;
    }

    if (nextErrors.type?.code === "incompatible_consumable_entries" && !nextErrors.global) {
      nextErrors.global = { message: nextErrors.type.message };
    }

    errors = nextErrors;
  }

  async function save() {
    if (saving) return;
    saving = true;
    errors = {};

    try {
      if (!currentUserId) {
        errors = { global: { message: "You must be logged in to save work records" } };
        return;
      }

      if (editing && id) {
        const saved = (await pb.send(`/api/work_records/${encodeURIComponent(id)}/save`, {
          method: "POST",
          body: buildSaveRequest(),
        })) as WorkRecordSaveResponse;
        await refreshVisibleWorkRecords(canViewAll);
        await goto(`/time/work-records/${saved.id}`);
        return;
      }

      const saved = (await pb.send("/api/work_records/save", {
        method: "POST",
        body: buildCreateRequest(),
      })) as WorkRecordSaveResponse;
      await refreshVisibleWorkRecords(canViewAll);
      await goto(`/time/work-records/${saved.id}`);
    } catch (error: any) {
      applySaveErrors(error);
    } finally {
      saving = false;
    }
  }

  $effect(() => {
    form = item;
    assigneeMode =
      canAssignOthers && form.uid !== "" && currentUserId !== "" && form.uid !== currentUserId
        ? "someone_else"
        : "you";

    if (detail?.consumables?.length) {
      consumables = detail.consumables.map((row) => ({ ...row, blocked: false }));
      consumableWarning = "";
    } else if (copyContext?.consumables?.length) {
      consumables = copyContext.consumables.map((row) => ({ ...row, blocked: false }));
      consumableWarning = "";
    } else {
      consumables = [];
      consumableWarning = "";
    }
  });

  $effect(() => {
    if (!editing && hasDefaultType && form.type === "") {
      form.type = $globalStore.profile.default_work_record_type;
    }
  });

  $effect(() => {
    void loadTypeConsumables(form.type ?? "");
  });
</script>

<form class="flex w-full flex-col gap-3 p-4">
  <div class="flex items-center justify-between gap-2">
    <h1 class="text-2xl font-bold">{editing ? "Edit Work Record" : "Create Work Record"}</h1>
    <div class="flex gap-2">
      <DsActionButton action={cancelAction}>Cancel</DsActionButton>
      {#if showSaveButton}
        <DsActionButton action={save} color="green" loading={saving}>Save</DsActionButton>
      {/if}
    </div>
  </div>

  {#if isCopyMode && copyContext}
    <div class="rounded-sm border border-blue-300 bg-blue-50 p-3 text-sm text-blue-900">
      Copying from {copyContext.source_work_record}. Worker and job stay locked to the source
      record.
    </div>
  {/if}

  {#if detail}
    <div class="rounded-sm bg-neutral-100 p-3 text-sm">
      <div class="flex flex-wrap items-center gap-2">
        <span class="font-semibold">Number:</span>
        <span>{detail.work_record}</span>
        {#if detail.approved}
          <DsLabel color="green">Approved</DsLabel>
        {/if}
      </div>
    </div>
  {:else}
    <div class="rounded-sm bg-neutral-100 p-3 text-sm">
      Number will be assigned after the first successful save.
    </div>
  {/if}

  {#if errors.global}
    <div class="rounded-sm border border-red-300 bg-red-50 p-3 text-sm text-red-700">
      {errors.global.message}
    </div>
  {/if}

  {#if canAssignOthers}
    <DsSelector
      bind:value={assigneeMode}
      items={assigneeModeOptions}
      {errors}
      fieldName="assignee_mode"
      uiName="Who Is This For?"
      disabled={isCopyMode}
    >
      {#snippet optionTemplate(item)}{item.name}{/snippet}
    </DsSelector>
  {/if}

  {#if canAssignOthers && assigneeMode === "someone_else" && $profiles.index !== null}
    <DsAutoComplete
      bind:value={form.uid}
      index={$profiles.index}
      {errors}
      fieldName="uid"
      uiName="Assignee"
      disabled={isCopyMode}
    >
      {#snippet resultTemplate(item)}{item.surname}, {item.given_name}{/snippet}
    </DsAutoComplete>
  {/if}

  <DsDateInput bind:value={form.date} min={DATE_INPUT_MIN} max={dateInputMax} />

  {#if $jobs.index !== null}
    <DsAutoComplete
      bind:value={form.job}
      index={$jobs.index}
      {errors}
      fieldName="job"
      uiName="Job"
      disabled={isCopyMode}
    >
      {#snippet resultTemplate(item)}{item.number} - {item.description}{/snippet}
    </DsAutoComplete>
  {/if}

  <DsTextInput bind:value={form.location} {errors} fieldName="location" uiName="Location" />
  <DsTextInput
    bind:value={form.work_description}
    {errors}
    fieldName="work_description"
    uiName="Work Description"
  />
  <DsTextInput bind:value={form.report_to} {errors} fieldName="report_to" uiName="Report To" />
  <DsTextInput
    bind:value={form.hours_on_site}
    {errors}
    fieldName="hours_on_site"
    uiName="Hours On Site"
    type="number"
    min={0}
    step={0.5}
  />
  <DsTextInput
    bind:value={form.hours_travel_time}
    {errors}
    fieldName="hours_travel_time"
    uiName="Hours Travel Time"
    type="number"
    min={0}
    step={0.5}
  />
  <DsTextInput
    bind:value={form.distance_travelled_km}
    {errors}
    fieldName="distance_travelled_km"
    uiName="Distance Travelled (km)"
    type="number"
    min={0}
    step={0.1}
  />
  <DsTextInput
    bind:value={form.field_book_number}
    {errors}
    fieldName="field_book_number"
    uiName="Field Book Number"
  />
  <DsTextInput
    bind:value={form.field_book_page_number}
    {errors}
    fieldName="field_book_page_number"
    uiName="Field Book Page Number"
    type="number"
    min={1}
    step={1}
  />
  <DsTextInput
    bind:value={form.sub_contractor}
    {errors}
    fieldName="sub_contractor"
    uiName="Sub Contractor"
  />
  <DsTextInput bind:value={form.equipment} {errors} fieldName="equipment" uiName="Equipment" />
  <DsTextInput bind:value={form.supplies} {errors} fieldName="supplies" uiName="Supplies" />

  <div class="flex flex-col gap-1">
    <div class="flex flex-wrap items-center gap-4 max-[640px]:flex-col max-[640px]:items-start">
      <div class="flex items-center gap-2">
        <span>Transport to Site</span>
        <DSToggle
          bind:value={form.vehicle_type}
          options={vehicleTypeOptions}
          ariaLabel="Transport to Site"
        />
      </div>
      <div class="flex items-center gap-2">
        <input
          id="is_passenger"
          type="checkbox"
          name="is_passenger"
          bind:checked={form.is_passenger}
          class="h-4 w-4"
        />
        <label for="is_passenger">I was a passenger</label>
      </div>
    </div>
    {#if errors.vehicle_type}
      <span class="text-red-600">{errors.vehicle_type.message}</span>
    {/if}
    {#if errors.is_passenger}
      <span class="text-red-600">{errors.is_passenger.message}</span>
    {/if}
  </div>

  <DsTextInput
    bind:value={form.company_vehicle_unit_number}
    {errors}
    fieldName="company_vehicle_unit_number"
    uiName="Company Vehicle Unit Number"
  />

  <div class="rounded-sm bg-neutral-50 p-4">
    <div class="flex items-center justify-between gap-2">
      <h2 class="text-lg font-semibold">Consumables</h2>
      {#if loadingConsumables}
        <span class="text-sm text-neutral-600">Loading…</span>
      {/if}
    </div>

    <p class="mt-1 text-sm text-neutral-600">
      If a needed consumable is not listed, record it in the work record notes.
    </p>

    <div class="mt-3">
      <DsSelector
        bind:value={form.type}
        items={workRecordTypeOptions}
        {errors}
        fieldName="type"
        uiName="Work Record Type"
      >
        {#snippet optionTemplate(item)}{item.name}{/snippet}
      </DsSelector>
    </div>

    {#if consumableWarning !== ""}
      <div class="mt-3 rounded-sm border border-amber-300 bg-amber-50 p-3 text-sm text-amber-900">
        {consumableWarning}
      </div>
    {/if}

    {#if visibleConsumables.length === 0 && blockedConsumables.length === 0 && !loadingConsumables}
      <div class="mt-3 text-sm text-neutral-600">No consumables for this type.</div>
    {/if}

    {#if visibleConsumables.length > 0}
      <div class="mt-3 grid gap-3">
        {#each visibleConsumables as consumable (consumable.consumable_id)}
          <div class="rounded-sm border border-neutral-200 p-3">
            <div class="font-medium">
              {consumable.consumable_name}
              {#if consumable.unit_label}
                <span class="text-sm text-neutral-500">({consumable.unit_label})</span>
              {/if}
            </div>

            {#if consumable.input_kind === "boolean"}
              <DsCheck
                bind:value={consumable.boolean_selected}
                errors={{}}
                fieldName={`consumable_${consumable.consumable_id}`}
                uiName="Selected"
              />
            {:else}
              <DsTextInput
                bind:value={consumable.quantity_number}
                errors={{}}
                fieldName={`consumable_${consumable.consumable_id}`}
                uiName="Quantity"
                type="number"
                min={0}
                step={0.1}
              />
            {/if}
          </div>
        {/each}
      </div>
    {/if}

    {#if blockedConsumables.length > 0}
      <div class="mt-4 space-y-3">
        <h3 class="font-semibold text-amber-900">Incompatible Existing Consumables</h3>
        {#each blockedConsumables as consumable (consumable.consumable_id)}
          <div class="rounded-sm border border-amber-300 bg-amber-50 p-3">
            <div class="font-medium text-amber-900">{consumable.consumable_name}</div>
            <div class="mt-1 text-sm text-amber-800">
              Clear this value before saving the selected work record type.
            </div>

            {#if consumable.input_kind === "boolean"}
              <DsCheck
                bind:value={consumable.boolean_selected}
                errors={{}}
                fieldName={`blocked_${consumable.consumable_id}`}
                uiName="Selected"
              />
            {:else}
              <DsTextInput
                bind:value={consumable.quantity_number}
                errors={{}}
                fieldName={`blocked_${consumable.consumable_id}`}
                uiName="Quantity"
                type="number"
                min={0}
                step={0.1}
              />
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</form>
