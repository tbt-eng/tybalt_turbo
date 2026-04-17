<script lang="ts">
  import { untrack } from "svelte";
  import {
    DATE_INPUT_MIN,
    applyDefaultDivisionOnce,
    applyDefaultRoleOnce,
    createJobCategoriesSync,
    dateInputMaxMonthsAhead,
    formatJobLabel,
  } from "$lib/utilities";
  import { divisions } from "$lib/stores/divisions";
  import { jobs } from "$lib/stores/jobs";
  import { rateRoles } from "$lib/stores/rateRoles";
  import { timeTypes } from "$lib/stores/time_types";
  import { branchesAllowedFor, branches as branchesStore } from "$lib/stores/branches";
  import { timeEditingDisabledMessage, timeEditingEnabled } from "$lib/stores/appConfig";
  import { pb } from "$lib/pocketbase";
  import DsTextInput from "$lib/components/DSTextInput.svelte";
  import DsSelector from "$lib/components/DSSelector.svelte";
  import DsDateInput from "$lib/components/DSDateInput.svelte";
  import DsAutoComplete from "./DSAutoComplete.svelte";
  import DsActionButton from "./DSActionButton.svelte";
  import DsEditingDisabledBanner from "./DsEditingDisabledBanner.svelte";
  import Icon from "@iconify/svelte";
  import { globalStore } from "$lib/stores/global";
  import { authStore } from "$lib/stores/auth";
  import { refreshVisibleWorkRecords } from "$lib/stores/workRecords";
  import { goto } from "$app/navigation";
  import type { TimeEntriesPageData, WorkRecordOption } from "$lib/svelte-types";
  import type {
    TimeTypesRecord,
    DivisionsRecord,
    CategoriesResponse,
    BranchesResponse,
  } from "$lib/pocketbase-types";

  // -----------------------------------------------------------------------
  // FEATURE FLAG: explicit branch picker on time entries
  // -----------------------------------------------------------------------
  // The backend (app/hooks/time_entries.go cleanTimeEntry) now resolves
  // time_entries.branch using purchase-order-style precedence:
  //
  //   1. If `job` is set, branch is forced to the job's branch.
  //   2. If `job` is blank and `branch` is blank, branch defaults from the
  //      caller's admin_profiles.default_branch.
  //   3. If `job` is blank and `branch` is already set, that explicit branch
  //      is preserved.
  //
  // Rule (3) means a user *can* — via the API — create a jobless time entry
  // against a branch that is not their default branch (subject to corporate
  // branch claim gating). The UI does not yet expose this. Whether we want to
  // expose it is still under discussion with stakeholders, so the picker is
  // shipped behind this in-code feature flag.
  //
  // To enable the picker for local testing, flip this constant to `true` and
  // restart the dev server. There is intentionally NO env-var / appConfig /
  // user-claim path to enable this — it must be a code change reviewed via PR
  // so that the discussion happens before any user sees the new field.
  //
  // When the flag is `false`:
  //   - the <DsSelector> for branch is not rendered
  //   - the auto-resolve / pin-tracking effects bail out immediately
  //   - the editor sends no `branch` field, so the backend falls through to
  //     rule (2) and the user's default branch is used (current behavior)
  //
  // Known issue:
  //   - On edit, the loaded record still includes its stored `branch`, and
  //     save() currently sends the whole record back to PocketBase.
  //   - That means a jobless entry created elsewhere (API / copy_to_tomorrow)
  //     with an explicit non-default branch can be silently preserved even
  //     while the picker is hidden.
  //   - If we keep that behavior long-term, the picker should become visible
  //     whenever such a record exists so the user can see and change the
  //     branch that is being resubmitted.
  //
  // When the flag is `true`:
  //   - the picker mirrors the PO editor: visible only when no job is set,
  //     populated with branches the caller is allowed to see (via
  //     allowed_claims), defaulting to the caller's default branch but
  //     "pinned" once the user picks something else in this editor session.
  // -----------------------------------------------------------------------
  const EXPLICIT_BRANCH_PICKER_ENABLED = false;

  // initialize the stores, noop if already initialized
  jobs.init();
  divisions.init();
  timeTypes.init();
  rateRoles.init();
  // The branches store is only consumed when the explicit branch picker is
  // enabled, but `init()` is idempotent and cheap so we always call it to
  // keep the wiring simple. If the flag is off the store data is unused.
  if (EXPLICIT_BRANCH_PICKER_ENABLED) {
    branchesStore.init();
  }

  let { data }: { data: TimeEntriesPageData } = $props();

  const trainingTokensInDescriptionWhileRegularHours = $derived.by(() => {
    if (item.time_type !== undefined && item.description !== undefined) {
      const lowercase = item.description.toLowerCase().trim();
      const lowercaseTokens = lowercase.split(/\s+/);
      return (
        hasTimeType(["R"]) &&
        (["training", "train", "orientation", "course", "whmis", "learning"].some((token) =>
          lowercaseTokens.includes(token),
        ) ||
          ["working at heights", "first aid"].some((token) => lowercase.includes(token)))
      );
    }
  });

  const jobNumbersInDescription = $derived.by(() => {
    if (item.description !== undefined) {
      const lowercase = item.description.toLowerCase().trim();
      // look for any instances of XX-YYY where XX is a number between 15 and
      // 40 and YYY is a zero-padded number between 1 and 999 then return true
      // if any are found
      return /(1[5-9]|2[0-9]|3[0-9]|40)-(\d{3})/.test(lowercase);
    }
  });

  const isWorkTime = $derived(hasTimeType(["R", "RT"]));
  const linkedWorkRecord = $derived.by(() => data.linked_work_record);
  const linkedWorkRecordLocked = $derived.by(() => (item.work_record_id ?? "") !== "");
  const isLegacyWorkRecordUpgrade = $derived.by(
    () => data.editing && (item.work_record ?? "") !== "" && !linkedWorkRecord,
  );

  let errors = $state({} as any);
  // Use untrack to capture initial value for form state (component recreates on navigation)
  let item = $state(untrack(() => data.item));
  const dateInputMax = dateInputMaxMonthsAhead(15);
  let compatibleWorkRecords = $state([] as WorkRecordOption[]);
  let loadingCompatibleWorkRecords = $state(false);
  let compatibleWorkRecordsMessage = $state("");
  let compatibleWorkRecordsRequestId = 0;
  const compatibleWorkRecordSelectOptions = $derived.by(() => [
    { id: "", work_record: "Select a work record" },
    ...compatibleWorkRecords,
  ]);

  // Track allowed division IDs for the selected job (null = show all divisions, empty Set = job has no allocations)
  let jobDivisionIds = $state<Set<string> | null>(null);
  let jobHasNoAllocations = $derived(jobDivisionIds !== null && jobDivisionIds.size === 0);
  const associatedJobDivisionCodes = $derived.by(() => {
    const selectedJobDivisionIds = jobDivisionIds;
    if (selectedJobDivisionIds === null || selectedJobDivisionIds.size === 0) {
      return [];
    }

    return $divisions.items
      .filter((division) => selectedJobDivisionIds.has(division.id))
      .map((division) => division.code)
      .filter((code): code is string => !!code)
      .sort((a, b) => a.localeCompare(b));
  });
  let jobAllocationsRequestId = 0;

  // Fetch job divisions when job changes
  $effect(() => {
    const jobId = item.job;
    if (!jobId || jobId === "") {
      // Invalidate any in-flight allocation request so stale responses are ignored.
      jobAllocationsRequestId++;
      jobDivisionIds = null; // No job = show all divisions
      return;
    }

    const requestId = ++jobAllocationsRequestId;
    // Clear the previous job's allocation state until this job's allocations load.
    jobDivisionIds = null;
    // Fetch job details to get allocations.
    // /api/jobs/:id does not include allocations, so we must use /details.
    pb.send(`/api/jobs/${jobId}/details`, { method: "GET" })
      .then((job: { allocations?: Array<{ division?: { id?: string } | string }> }) => {
        if (requestId !== jobAllocationsRequestId) return;

        const allocationDivisionIds =
          job.allocations
            ?.map((allocation) =>
              typeof allocation.division === "string"
                ? allocation.division
                : allocation.division?.id,
            )
            .filter((id): id is string => !!id) ?? [];

        if (allocationDivisionIds.length > 0) {
          jobDivisionIds = new Set(allocationDivisionIds);
        } else {
          jobDivisionIds = new Set(); // Job has no allocations = empty set (will show no divisions)
        }
      })
      .catch(() => {
        if (requestId !== jobAllocationsRequestId) return;
        jobDivisionIds = null;
      });
  });

  // ---------------------------------------------------------------------
  // Branch picker state (gated by EXPLICIT_BRANCH_PICKER_ENABLED)
  // ---------------------------------------------------------------------
  // Mirrors the pin-tracking + auto-resolve logic in PurchaseOrdersEditor.
  // We declare the state unconditionally so the component shape is stable
  // regardless of the flag, but every effect bails out as a no-op when the
  // flag is off.
  //
  //   branchPinnedInSession         — true once the user manually changes
  //                                   the branch away from the auto-derived
  //                                   value; while pinned, job changes will
  //                                   not overwrite it.
  //   lastAutoBranch                — the last value we auto-assigned, so we
  //                                   can distinguish "auto" from "manual".
  //   branchChangeWatchInitialized  — guard so the very first observation of
  //                                   item.branch doesn't count as a manual
  //                                   change.
  //   lastObservedJobForAuto        — tracks job transitions so we can
  //                                   re-derive only when job actually
  //                                   changes (not on every render).
  //   branchLookupRequestId         — monotonically-increasing token used to
  //                                   discard stale async job lookups.
  let branchPinnedInSession = $state(false);
  let lastAutoBranch = $state("");
  let branchChangeWatchInitialized = $state(false);
  let lastObservedJobForAuto = $state<string | null>(null);
  let branchLookupRequestId = 0;

  // Caller's default branch from their admin profile. Used as the fallback
  // when no job is set and the user has not yet picked a branch.
  const creatorDefaultBranch = $derived.by(() => $globalStore.profile.default_branch ?? "");

  // Users only see branches they can actually post against (e.g. the
  // "corporate" branch is hidden unless the caller holds the corporate_branch
  // claim), plus the branch already on the record.
  const availableBranches = $derived.by(() =>
    EXPLICIT_BRANCH_PICKER_ENABLED
      ? branchesAllowedFor($branchesStore.items, $globalStore.claimIds, item.branch)
      : [],
  );

  // Async helper: resolve the branch implied by a job id, falling back to
  // the supplied default if the job has no branch or the lookup fails.
  // Mirrors PurchaseOrdersEditor.resolveDerivedBranch.
  async function resolveDerivedBranchForTimeEntry(
    jobId: string,
    fallbackDefaultBranch: string,
  ): Promise<string> {
    if (jobId !== "") {
      try {
        const job = await pb.collection("jobs").getOne(jobId, { requestKey: null });
        return job.branch ?? "";
      } catch (error) {
        console.error("Error loading job branch for time entry:", error);
      }
    }
    return fallbackDefaultBranch;
  }

  // Effect 1 — detect manual branch changes and "pin" them.
  // Only runs while no job is set (when a job is set the branch is forced
  // to the job's branch and pinning is meaningless). The first observation
  // of a non-empty branch is treated as the auto-assigned value, not a
  // manual change.
  $effect(() => {
    if (!EXPLICIT_BRANCH_PICKER_ENABLED) return;
    const branch = item.branch ?? "";
    if ((item.job ?? "") !== "") {
      return;
    }
    if (!branchChangeWatchInitialized) {
      branchChangeWatchInitialized = true;
      return;
    }
    if (branch === "" || branch === lastAutoBranch) {
      return;
    }
    branchPinnedInSession = true;
  });

  // Effect 2 — auto-resolve branch from job (or from default branch when
  // no job is set and the branch hasn't been pinned). Uses a request-id
  // token so a slow job lookup can't clobber a newer selection.
  $effect(() => {
    if (!EXPLICIT_BRANCH_PICKER_ENABLED) return;
    const jobId = item.job ?? "";
    const fallbackDefaultBranch = creatorDefaultBranch;
    const branch = item.branch ?? "";
    const pinned = branchPinnedInSession;

    if (jobId !== "") {
      // A job is selected — branch must follow the job. Clear any pin so
      // that subsequently clearing the job won't preserve a stale value.
      branchPinnedInSession = false;
      lastObservedJobForAuto = jobId;
      const requestId = ++branchLookupRequestId;
      resolveDerivedBranchForTimeEntry(jobId, fallbackDefaultBranch).then((derivedBranch) => {
        if (requestId !== branchLookupRequestId || derivedBranch === "") {
          return;
        }
        if (item.branch === derivedBranch) {
          return;
        }
        lastAutoBranch = derivedBranch;
        item.branch = derivedBranch;
      });
      return;
    }

    if (pinned) {
      lastObservedJobForAuto = jobId;
      return;
    }

    // No job, not pinned. Auto-populate when branch starts empty, and
    // auto-switch when the job has just been *cleared* (so a user who had
    // been working under a job's branch is bounced back to their default
    // rather than keeping the now-orphaned job branch).
    const firstObservation = lastObservedJobForAuto === null;
    const jobChanged = !firstObservation && lastObservedJobForAuto !== jobId;
    lastObservedJobForAuto = jobId;

    if (!jobChanged && branch !== "") {
      return;
    }

    const requestId = ++branchLookupRequestId;
    resolveDerivedBranchForTimeEntry(jobId, fallbackDefaultBranch).then((derivedBranch) => {
      if (requestId !== branchLookupRequestId || derivedBranch === "") {
        return;
      }
      if (item.branch === derivedBranch) {
        return;
      }
      lastAutoBranch = derivedBranch;
      item.branch = derivedBranch;
    });
  });

  // Default division from caller's profile if creating and empty
  $effect(() => applyDefaultDivisionOnce(item, data.editing));

  // Default role from caller's profile if creating, has a job, and role is empty
  $effect(() => {
    if (item.job && item.job !== "") {
      applyDefaultRoleOnce(item, data.editing);
    }
  });

  // given a list of time type codes, return true if the item's time type is in
  // the list
  function hasTimeType(typelist: string[]) {
    if ($timeTypes.items.length === 0) {
      return false;
    }
    return typelist
      .map((c) => {
        const foundTimeType = $timeTypes.items.find((t) => t.code === c);
        return foundTimeType ? foundTimeType.id : null;
      })
      .includes(item.time_type);
  }

  const hasNoDefaultRole = $derived(!$globalStore.profile.default_role);
  let showRoleInfo = $state(false);

  let categories = $state([] as CategoriesResponse[]);
  const syncCategoriesForJob = createJobCategoriesSync((rows) => {
    categories = rows;
  });

  // Watch for changes to the job and fetch categories accordingly
  $effect(() => {
    syncCategoriesForJob(item.job);
  });

  $effect(() => {
    const legacyUpgrade = isLegacyWorkRecordUpgrade;
    const jobId = item.job ?? "";
    const date = item.date ?? "";
    const entryId = data.id ?? "";

    if (!legacyUpgrade || jobId === "" || date === "") {
      compatibleWorkRecordsRequestId++;
      compatibleWorkRecords = [];
      compatibleWorkRecordsMessage = "";
      loadingCompatibleWorkRecords = false;
      return;
    }

    const requestId = ++compatibleWorkRecordsRequestId;
    loadingCompatibleWorkRecords = true;
    compatibleWorkRecordsMessage = "";
    const params = new URLSearchParams({ job: jobId, date });
    if (entryId !== "") {
      params.set("exclude_entry", entryId);
    }

    pb.send(`/api/work_records/compatible_for_time_entry?${params.toString()}`, {
      method: "GET",
    })
      .then((rows: WorkRecordOption[]) => {
        if (requestId !== compatibleWorkRecordsRequestId) return;
        compatibleWorkRecords = rows;
        if (rows.length === 0) {
          compatibleWorkRecordsMessage =
            "No compatible unlinked work records were found for this job and date.";
        } else if ((item.work_record_id ?? "") === "" && rows.length === 1) {
          item.work_record_id = rows[0].id;
        }
      })
      .catch((error) => {
        if (requestId !== compatibleWorkRecordsRequestId) return;
        console.error("Error loading compatible work records:", error);
        compatibleWorkRecords = [];
        compatibleWorkRecordsMessage = "Failed to load compatible work records.";
      })
      .finally(() => {
        if (requestId !== compatibleWorkRecordsRequestId) return;
        loadingCompatibleWorkRecords = false;
      });
  });

  async function save() {
    // set the uid from the authStore We do it here rather than the
    // backend because the
    // OnRecordBeforeCreateRequest("time_entries").Add() hook runs *AFTER*
    // schema validation (!)
    // https://github.com/pocketbase/pocketbase/discussions/2881 so we
    // need a correct value in the payload just to make it to the hook
    const userId = $authStore?.model?.id;
    if (!userId) {
      errors = { global: { message: "You must be logged in to save time entries" } };
      return;
    }
    item.uid = userId;

    // set a dummy value for week_ending to satisfy the schema non-empty
    // requirement. This will be changed in the backend to the correct
    // value every time a record is saved
    item.week_ending = "2006-01-02";

    // if the job is empty, clear job-related fields
    if (item.job === "") {
      item.category = "";
      item.work_record = "";
      item.work_record_id = "";
      item.role = "";
    }

    const payload: Record<string, unknown> = { ...item };
    delete payload.expand;

    if ((item.work_record_id ?? "") !== "") {
      payload.work_record = "";
    } else if (!isLegacyWorkRecordUpgrade) {
      payload.work_record = "";
    }

    try {
      if (data.editing && data.id !== null) {
        // update the item
        await pb.collection("time_entries").update(data.id, payload);
      } else {
        // create a new item
        await pb.collection("time_entries").create(payload);
      }

      // submission was successful, clear the errors
      errors = {};

      await refreshVisibleWorkRecords($globalStore.user_claims_summary.claims.includes("report"));

      // redirect to the list page
      goto("/time/entries/list");
    } catch (error: any) {
      errors = error?.data?.data ?? { global: { message: "Failed to save time entry" } };
    }
  }
</script>

{#if !$timeEditingEnabled}
  <DsEditingDisabledBanner message={timeEditingDisabledMessage} />
{/if}
<form
  class="flex w-full flex-col items-center gap-2 p-2 max-lg:[&_button]:text-base max-lg:[&_input]:text-base max-lg:[&_label]:text-base max-lg:[&_select]:text-base max-lg:[&_textarea]:text-base"
>
  <div class="flex w-full flex-col gap-2">
    <div class="flex w-full flex-wrap items-center justify-between gap-2">
      <h1 class="text-xl font-bold text-neutral-800">
        {#if data.editing}
          Editing Time Entry
        {:else}
          Create Time Entry
        {/if}
      </h1>
      {#if !data.editing}
        <DsActionButton action="/time/work-records" color="green">
          Need a linked entry? Start from Work Records
        </DsActionButton>
      {/if}
    </div>

    {#if !data.editing}
      <div class="w-full rounded-sm border border-green-300 bg-green-50 p-3 text-sm text-green-900">
        When a time entry needs a work record, start from the work record list or detail page so
        the link is created the right way.
      </div>
    {/if}

    {#if linkedWorkRecord && linkedWorkRecordLocked}
      <div class="w-full rounded-sm border border-cyan-300 bg-cyan-50 p-3 text-sm text-cyan-950">
        <div class="flex flex-wrap items-center gap-2">
          <DsLabel color="cyan">Linked Work Record {linkedWorkRecord.work_record}</DsLabel>
          {#if linkedWorkRecord.job_number}
            <span>{linkedWorkRecord.job_number}</span>
          {/if}
          <span>{linkedWorkRecord.hours_on_site + linkedWorkRecord.hours_travel_time} hours</span>
        </div>
        <div class="mt-1 opacity-80">{linkedWorkRecord.work_description}</div>
      </div>
    {/if}

    {#if isLegacyWorkRecordUpgrade}
      <div class="w-full rounded-sm border border-amber-300 bg-amber-50 p-3 text-sm text-amber-900">
        This legacy time entry still uses a text work record. Choose an existing first-class work
        record before saving any changes.
      </div>
    {/if}
  </div>

  <span class="flex w-full gap-2">
    <label for="date">Date</label>
    <DsDateInput
      class="flex-1 md:flex-none"
      name="date"
      min={DATE_INPUT_MIN}
      max={dateInputMax}
      bind:value={item.date}
      disabled={linkedWorkRecordLocked}
    />
  </span>
  {#snippet optionTemplate(item: TimeTypesRecord | DivisionsRecord)}
    {item.code} - {item.name}
  {/snippet}
  <DsSelector
    bind:value={item.time_type as string}
    items={$timeTypes.items}
    {errors}
    {optionTemplate}
    fieldName="time_type"
    uiName="Time Type"
    disabled={linkedWorkRecordLocked}
  />
  {#if trainingTokensInDescriptionWhileRegularHours}
    <span class="flex w-full gap-2 bg-red-200 text-red-600">
      ^Should you choose training instead?
    </span>
  {/if}

  <!-- ----------------------------------------------------------------- -->
  <!-- FEATURE-FLAGGED branch picker (EXPLICIT_BRANCH_PICKER_ENABLED).   -->
  <!-- Shown only when no job is selected. Once a job is picked the      -->
  <!-- backend forces the branch to the job's branch and the picker is  -->
  <!-- meaningless, so we hide it. When the flag is off this entire     -->
  <!-- block is compiled away by the conditional and the editor sends  -->
  <!-- no branch field, falling back to the user's default branch on   -->
  <!-- the server.                                                     -->
  <!-- ----------------------------------------------------------------- -->
  {#if EXPLICIT_BRANCH_PICKER_ENABLED && (item.job ?? "") === ""}
    <DsSelector
      bind:value={item.branch as string}
      items={availableBranches}
      {errors}
      fieldName="branch"
      uiName="Branch"
    >
      {#snippet optionTemplate(item: BranchesResponse)}
        {item.name}
      {/snippet}
    </DsSelector>
  {/if}

  <!----------------------------------------------->
  <!-- FIELDS VISIBLE ONLY FOR R or RT TimeTypes -->
  <!----------------------------------------------->
  {#if isWorkTime}
    {#if $divisions.index !== null}
      <DsAutoComplete
        bind:value={item.division as string}
        index={$divisions.index}
        {errors}
        fieldName="division"
        uiName="Division"
        filter={jobDivisionIds ? (div) => jobDivisionIds!.has(div.id) : undefined}
      >
        {#snippet resultTemplate(item)}{item.code} - {item.name}{/snippet}
      </DsAutoComplete>
    {/if}
    {#if $jobs.index !== null}
      <DsAutoComplete
        bind:value={item.job as string}
        index={$jobs.index}
        {errors}
        fieldName="job"
        uiName="Job"
        disabled={linkedWorkRecordLocked}
      >
        {#snippet resultTemplate(item)}{formatJobLabel(item)}{/snippet}
      </DsAutoComplete>
    {/if}
    {#if item.job && item.job !== "" && associatedJobDivisionCodes.length > 0}
      <span class="w-full text-sm text-neutral-600">
        Associated divisions: {associatedJobDivisionCodes.join(", ")}
      </span>
    {/if}
    {#if jobHasNoAllocations}
      <span class="flex w-full gap-2 bg-red-200 text-red-600">
        This job has no division allocations configured. Contact an administrator.
      </span>
    {/if}
    {#if item.job && item.job !== "" && $rateRoles.index !== null}
      <div class="relative flex w-full items-center gap-1">
        {#if hasNoDefaultRole}
          <button
            type="button"
            class="shrink-0 text-blue-500 hover:text-blue-700"
            aria-label="Role info"
            onclick={() => (showRoleInfo = !showRoleInfo)}
          >
            <Icon icon="mdi:information-outline" width="18px" />
          </button>
        {/if}
        <div class="flex-1">
          <DsAutoComplete
            bind:value={item.role as string}
            index={$rateRoles.index}
            {errors}
            fieldName="role"
            uiName="Role"
          >
            {#snippet resultTemplate(item)}{item.name}{/snippet}
          </DsAutoComplete>
        </div>
      </div>
      {#if showRoleInfo}
        <div
          class="fixed inset-0 z-50"
          role="button"
          tabindex="-1"
          onclick={() => (showRoleInfo = false)}
          onkeydown={(e) => e.key === "Escape" && (showRoleInfo = false)}
        >
          <div
            class="mx-auto mt-32 max-w-sm rounded-md border border-blue-200 bg-blue-50 p-4 text-sm text-blue-800 shadow-lg"
            role="presentation"
            onclick={(e) => e.stopPropagation()}
          >
            <p class="font-semibold">Rate Sheet Role</p>
            <p class="mt-1">
              This is the rate sheet role you are using to charge against the job for this time
              entry. If you do not know your role, speak with your manager or the project manager
              for this job. You may also set a default role in your profile so this field is
              automatically populated. Set the default role by clicking on your email address at the
              bottom of the navigation menu to access your profile settings.
            </p>
            <button
              type="button"
              class="mt-2 text-xs font-medium text-blue-600 hover:text-blue-800"
              onclick={() => (showRoleInfo = false)}
            >
              Got it
            </button>
          </div>
        </div>
      {/if}
    {/if}
  {/if}

  {#if item.job !== "" && categories.length > 0}
    <DsSelector
      bind:value={item.category as string}
      items={categories}
      {errors}
      fieldName="category"
      uiName="Category"
      clear={true}
    >
      {#snippet optionTemplate(item: CategoriesResponse)}
        {item.name}
      {/snippet}
    </DsSelector>
  {/if}

  {#if linkedWorkRecord && linkedWorkRecordLocked}
    <div class="flex w-full flex-col gap-2">
      <span class="flex w-full gap-2">
        <span class="font-medium">Work Record</span>
        <span>{linkedWorkRecord.work_record}</span>
      </span>
      <div class="text-sm text-neutral-600">
        Job, date, hours, and link ownership come from the work record. Edit the description here
        if needed.
      </div>
    </div>
  {:else if isLegacyWorkRecordUpgrade}
    <DsTextInput
      bind:value={item.work_record as string}
      {errors}
      fieldName="work_record"
      uiName="Legacy Work Record"
      disabled={true}
    />
    <DsSelector
      bind:value={item.work_record_id as string}
      items={compatibleWorkRecordSelectOptions}
      {errors}
      fieldName="work_record_id"
      uiName="Existing Work Record"
      disabled={loadingCompatibleWorkRecords}
    >
      {#snippet optionTemplate(item: WorkRecordOption)}
        {item.work_record}
      {/snippet}
    </DsSelector>
    {#if loadingCompatibleWorkRecords}
      <span class="w-full text-sm text-neutral-600">Loading compatible work records...</span>
    {:else if compatibleWorkRecordsMessage}
      <div class="w-full rounded-sm border border-amber-300 bg-amber-50 p-2 text-sm text-amber-900">
        {compatibleWorkRecordsMessage}
        <div class="mt-2">
          <DsActionButton action="/time/work-records" color="green">
            Go to Work Records
          </DsActionButton>
        </div>
      </div>
    {/if}
  {/if}

  <!--------------------------------------------------->
  <!-- END FIELDS VISIBLE ONLY FOR R or RT TimeTypes -->
  <!--------------------------------------------------->

  <!-- TODO: The item.job === undefined below is predecated on the
  autocomplete clearing the property. Right now we are using a text field
  so it will never show up after being set once -->
  {#if !hasTimeType(["OR", "OW", "OTO"])}
    <DsTextInput
      bind:value={item.hours as number}
      {errors}
      fieldName="hours"
      uiName="Hours"
      type="number"
      step={0.5}
      min={0}
      max={18}
      disabled={linkedWorkRecordLocked}
    />
  {/if}

  {#if item.division && isWorkTime}
    <DsTextInput
      bind:value={item.meals_hours as number}
      {errors}
      fieldName="meals_hours"
      uiName="Meals Hours"
      type="number"
      step={0.5}
      min={0}
      max={3}
    />
  {/if}

  {#if !hasTimeType(["OR", "OW", "OTO", "RB"])}
    <DsTextInput
      bind:value={item.description as string}
      {errors}
      fieldName="description"
      uiName="Description"
    />
  {/if}
  {#if jobNumbersInDescription}
    <span class="flex w-full gap-2 bg-red-200 text-red-600">
      Job numbers are not allowed in the description. Enter jobs numbers in the appropriate field
      and create one time entry per job.
    </span>
  {/if}

  {#if hasTimeType(["OTO"])}
    <DsTextInput
      bind:value={item.payout_request_amount as number}
      {errors}
      fieldName="payout_request_amount"
      uiName="$"
      type="number"
    />
  {/if}

  <div class="flex w-full flex-col gap-2 {errors.global !== undefined ? 'bg-red-200' : ''}">
    <span class="flex w-full gap-2">
      {#if !jobNumbersInDescription}
        <DsActionButton action={save}>Save</DsActionButton>
      {/if}
      <DsActionButton action="/time/entries/list">Cancel</DsActionButton>
    </span>
    {#if errors.global !== undefined}
      <span class="text-red-600">{errors.global.message}</span>
    {/if}
  </div>
</form>
