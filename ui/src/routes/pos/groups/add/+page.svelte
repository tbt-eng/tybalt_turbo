<script lang="ts">
  // A super_po holder fills out one purchase order as a template, chooses which
  // fields come from a CSV upload, uploads one row per staff member, reviews
  // what will be created, and creates one purchase order per row (issue #168).
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSAutoComplete from "$lib/components/DSAutoComplete.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCollapsible from "$lib/components/DSCollapsible.svelte";
  import DSCurrencyInput from "$lib/components/DSCurrencyInput.svelte";
  import DSDateInput from "$lib/components/DSDateInput.svelte";
  import DSList from "$lib/components/DSList.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DSSelector from "$lib/components/DSSelector.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import DSToggle from "$lib/components/DSToggle.svelte";
  import DsCheck from "$lib/components/DsCheck.svelte";
  import DsEditingDisabledBanner from "$lib/components/DsEditingDisabledBanner.svelte";
  import DsFileSelect from "$lib/components/DsFileSelect.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";
  import VendorSelector from "$lib/components/VendorSelector.svelte";
  import { localIsoDate } from "$lib/localDate";
  import {
    newPoGroupId,
    poGroupColumns,
    poGroupUploadFieldOptions,
    readPoGroupCsv,
    type PoGroupUploadField,
  } from "$lib/poGroupCsv";
  import {
    createPoGroup,
    previewPoGroup,
    type PoGroupApiError,
    type PoGroupPreview,
    type PoGroupPreviewRow,
    type PoGroupRequest,
    type PoGroupTemplate,
  } from "$lib/poGroups";
  import type { BranchesResponse, PoApproversResponse } from "$lib/pocketbase-types";
  import { expensesEditingEnabled } from "$lib/stores/appConfig";
  import { branches, branchesAllowedFor } from "$lib/stores/branches";
  import { currencies } from "$lib/stores/currencies";
  import { divisions } from "$lib/stores/divisions";
  import { expenditureKinds } from "$lib/stores/expenditureKinds";
  import { globalStore } from "$lib/stores/global";
  import { DATE_INPUT_MIN, downloadCsvHeaders, formatCurrencyAmount } from "$lib/utilities";

  divisions.init();
  branches.init();
  currencies.init();
  expenditureKinds.init();

  const typeOptions = [
    {
      id: "Cumulative",
      label: "Cumulative",
      description: "allows several expenses up to the total",
    },
    { id: "One-Time", label: "One-Time", description: "for a single expense, closes after use" },
  ];
  const paymentTypeOptions = [
    { id: "OnAccount", label: "On Account", description: "charged directly to a vendor account" },
    { id: "Expense", label: "Expense", description: "paid out of pocket and reimbursed" },
    {
      id: "CorporateCreditCard",
      label: "Corporate Credit Card",
      description: "paid with a corporate card",
    },
  ];
  // Group purchase orders have no job, so the project kind, which needs one,
  // is left out.
  const kindOptions = $derived(
    $expenditureKinds.items
      .filter((kind) => kind.name !== "project")
      .map((kind) => ({ id: kind.id, label: kind.en_ui_label, description: kind.description })),
  );

  const canCreateGroups = $derived($globalStore.claims.includes("super_po"));

  let template = $state<PoGroupTemplate>({
    type: "Cumulative",
    kind: "",
    payment_type: "OnAccount",
    date: localIsoDate(),
    division: "",
    branch: "",
    description: "",
    total: 0,
    currency: "",
    vendor: "",
    approver: "",
    priority_second_approver: "",
  });
  let fromUpload = $state<Record<PoGroupUploadField, boolean>>({
    total: true,
    description: false,
    division: false,
    branch: false,
  });
  let upload = $state({ file: "" as File | string });
  let fileText = $state("");
  let errors = $state<Record<string, { message: string }>>({});
  let preview = $state<PoGroupPreview | null>(null);
  let previewedKey = $state("");
  let reviewing = $state(false);
  let creating = $state(false);
  // Made once per page so that a repeated create finds the group it made.
  const groupId = newPoGroupId();

  // Fill each default once, when its store loads, so that clearing a field
  // afterwards sticks.
  const defaulted = { kind: false, currency: false, division: false };
  $effect(() => {
    const kind = $expenditureKinds.items.find((row) => row.name === "capital")?.id;
    if (kind && !defaulted.kind) {
      defaulted.kind = true;
      template.kind = kind;
    }
    const currency = $currencies.items.find((row) => row.code === "CAD")?.id;
    if (currency && !defaulted.currency) {
      defaulted.currency = true;
      template.currency = currency;
    }
    const division = $globalStore.profile.default_division;
    if (division && !defaulted.division) {
      defaulted.division = true;
      template.division = division;
    }
  });

  $effect(() => {
    const file = upload.file;
    if (file instanceof File) {
      file.text().then((text) => (fileText = text));
    } else {
      fileText = "";
    }
  });

  const uploadFields = $derived(
    poGroupUploadFieldOptions.map((option) => option.id).filter((id) => fromUpload[id]),
  );
  const columns = $derived(poGroupColumns(uploadFields));
  const csv = $derived(fileText === "" ? null : readPoGroupCsv(fileText, uploadFields));
  const request = $derived<PoGroupRequest>({
    template: { ...template, total: Number(template.total) || 0 },
    upload_fields: uploadFields,
    rows: csv?.rows ?? [],
  });
  const requestKey = $derived(JSON.stringify(request));
  const readyToReview = $derived(
    csv !== null && csv.rows.length > 0 && csv.missingColumns.length === 0,
  );
  const previewIsCurrent = $derived(preview !== null && previewedKey === requestKey);
  const canCreate = $derived(
    canCreateGroups &&
      $expensesEditingEnabled &&
      previewIsCurrent &&
      preview !== null &&
      preview.count > 0 &&
      preview.error_count === 0 &&
      template.approver !== "" &&
      !creating,
  );
  const approverName = $derived.by(() => {
    const approver = preview?.approvers.find((row) => row.id === template.approver);
    return approver ? `${approver.given_name} ${approver.surname}` : "";
  });
  const readyMessage = $derived(
    preview === null
      ? ""
      : [
          `${preview.count} purchase orders totalling ${formatCurrencyAmount(preview.total, preview.currency_code)}.`,
          approverName ? `${approverName} will get one email asking them to approve them.` : "",
        ]
          .filter((part) => part)
          .join(" "),
  );
  // The creator sets the branch, so offer only the branches their claims allow.
  const availableBranches = $derived(branchesAllowedFor($branches.items, $globalStore.claimIds));
  const selectedCurrencyCode = $derived(
    $currencies.items.find((row) => row.id === template.currency)?.code ?? "CAD",
  );
  // Template errors the form has no input to show beside.
  const otherErrors = $derived(
    Object.entries(errors).filter(
      ([field]) =>
        ![
          "description",
          "total",
          "currency",
          "vendor",
          "division",
          "branch",
          "approver",
          "priority_second_approver",
          "global",
        ].includes(field),
    ),
  );

  // The review lists who can approve every row. Drop a choice that the latest
  // review no longer offers, as the PO form does when the amount changes.
  $effect(() => {
    if (!previewIsCurrent || preview === null) return;
    const offers = (list: PoApproversResponse[], id: string) => list.some((row) => row.id === id);
    if (template.approver !== "" && !offers(preview.approvers, template.approver)) {
      template.approver = "";
    }
    if (
      template.priority_second_approver !== "" &&
      !(
        preview.second_approval_required &&
        offers(preview.second_approvers, template.priority_second_approver)
      )
    ) {
      template.priority_second_approver = "";
    }
  });

  // Review again once the template or the upload stops changing. Every row is
  // checked on the server, so wait for a pause in typing. A newer review
  // cancels the one in flight. Without a usable upload there is nothing to
  // review, and errors from an earlier review no longer apply.
  $effect(() => {
    const key = requestKey;
    if (!readyToReview || !canCreateGroups) {
      preview = null;
      errors = {};
      return;
    }
    const timer = setTimeout(() => review(key), 750);
    return () => clearTimeout(timer);
  });

  function responseErrors(error: PoGroupApiError): Record<string, { message: string }> {
    const body = error.response ?? {};
    if (body.data && Object.keys(body.data).length > 0) {
      return body.data;
    }
    return { global: { message: body.message ?? error.message ?? "Something went wrong." } };
  }

  async function review(key: string) {
    reviewing = true;
    try {
      const result = await previewPoGroup(JSON.parse(key));
      if (key !== requestKey) return;
      preview = result;
      previewedKey = key;
      errors = {};
    } catch (caught) {
      const error = caught as PoGroupApiError;
      if (error.isAbort) return;
      preview = null;
      errors = responseErrors(error);
    } finally {
      reviewing = false;
    }
  }

  async function create() {
    if (!canCreate) return;
    creating = true;
    try {
      await createPoGroup({ ...request, group_id: groupId });
      await goto(resolve(`/pos/groups/${groupId}`));
    } catch (caught) {
      const error = caught as PoGroupApiError;
      const body = error.response ?? {};
      if (body.code === "group_exists") {
        await goto(resolve(`/pos/groups/${groupId}`));
      } else if (body.code === "rows_invalid" && body.preview) {
        preview = body.preview;
        previewedKey = requestKey;
      } else {
        errors = responseErrors(error);
      }
    } finally {
      creating = false;
    }
  }

  function rowLine(row: PoGroupPreviewRow): number {
    return csv?.lines[row.index] ?? row.index + 2;
  }

  function errorLabel(field: string): string {
    if (field === "uid" || field === "email") return "Staff";
    if (field === "priority_second_approver") return "Priority second approver";
    if (field === "global") return "Row";
    const name = field.replaceAll("_", " ");
    return name.charAt(0).toUpperCase() + name.slice(1);
  }
</script>

{#if !$expensesEditingEnabled}
  <DsEditingDisabledBanner
    message="Purchase order editing is currently disabled during a system transition."
  />
{/if}

<div class="flex w-full flex-col gap-4 p-2">
  <DSPageHeader title="Create PO Group" />

  {#if !canCreateGroups}
    <DSCallout tone="danger" title="super_po claim needed">
      Only super_po holders can create PO groups. You can look at existing groups on the PO Groups
      page.
    </DSCallout>
  {:else}
    <DSCallout tone="info" title="How a PO group works">
      <ol class="list-decimal space-y-1 pl-5">
        <li>
          Fill out one purchase order as a template. Every person in the upload gets their own copy:
          they own it, spend against it and see it in My Purchase Orders, just as if they had
          created it.
        </li>
        <li>
          Choose the fields that differ from person to person, usually the total, and upload a CSV
          with one row per person. Download the template CSV to get the right columns.
        </li>
        <li>
          Choose the approver. Only people who can approve every row's amount are offered, and a
          second approver is asked for when some amounts need one.
        </li>
        <li>
          Review the result. Every row is checked the way a normal purchase order is: the person is
          active and the values are valid. Nothing is created until every row passes, and then all
          of them are created at once.
        </li>
        <li>
          The approver gets one email for the whole group and approves each purchase order on its
          own from Pending My Approval. Each one gets its PO number when it's approved. A second
          approver, when amounts need one, also gets just one email for the group.
        </li>
      </ol>
      <p class="mt-2">
        A group can't be undone in one step, so check the review carefully. After it's created, an
        unapproved purchase order can still be edited or deleted by its owner.
      </p>
    </DSCallout>

    <DSSection title="1. Template">
      <div class="flex w-full flex-col gap-2">
        <p>
          Everything here is copied to every purchase order, except the fields you take from the
          upload in step 2.
        </p>
        <DSToggle
          bind:value={template.type}
          label="Type"
          options={typeOptions}
          showOptionDescriptions={true}
          fullWidth={true}
        />
        <DSToggle
          bind:value={template.kind}
          label="Kind"
          options={kindOptions}
          showOptionDescriptions={true}
          fullWidth={true}
        />
        <DSToggle
          bind:value={template.payment_type}
          label="Payment"
          options={paymentTypeOptions}
          showOptionDescriptions={true}
          fullWidth={true}
        />
        <span class="flex w-full items-center gap-2">
          <label for="date">Date</label>
          <DSDateInput
            class="flex-1 md:flex-none"
            name="date"
            min={DATE_INPUT_MIN}
            bind:value={template.date}
          />
        </span>
        {#if !fromUpload.division && $divisions.index !== null}
          <DSAutoComplete
            bind:value={template.division}
            index={$divisions.index}
            {errors}
            fieldName="division"
            uiName="Division"
          >
            {#snippet resultTemplate(item)}{item.code} - {item.name}{/snippet}
          </DSAutoComplete>
        {/if}
        {#if !fromUpload.branch}
          <DSSelector
            bind:value={template.branch}
            items={availableBranches}
            {errors}
            fieldName="branch"
            uiName="Branch"
            helperText="Every purchase order uses this branch, not each person's default branch."
          >
            {#snippet optionTemplate(item: BranchesResponse)}{item.name}{/snippet}
          </DSSelector>
        {/if}
        {#if !fromUpload.description}
          <DSTextInput
            bind:value={template.description}
            {errors}
            fieldName="description"
            uiName="Description"
          />
        {/if}
        <DSCurrencyInput
          bind:amount={template.total}
          bind:currency={template.currency}
          items={$currencies.items}
          {errors}
          amountFieldName="total"
          currencyFieldName="currency"
          uiName="Total"
          disabledAmount={fromUpload.total}
          helperText={fromUpload.total
            ? `Each person's total comes from the upload, in ${selectedCurrencyCode}.`
            : `Max in ${selectedCurrencyCode} for each person, including all taxes and shipping.`}
        />
        <VendorSelector bind:value={template.vendor} {errors} />
      </div>
    </DSSection>

    <DSSection title="2. Upload">
      <div class="flex w-full flex-col gap-2">
        <p>
          The upload is a CSV file with one row per person. Its <strong>email</strong> column must match
          each person's Turbo email; that person owns the purchase order. Tick the fields that come from
          the upload instead of the template. Other columns, such as names, are ignored, so you can keep
          them for your own reference.
        </p>
        {#each poGroupUploadFieldOptions as option (option.id)}
          <DsCheck
            bind:value={fromUpload[option.id]}
            errors={{}}
            fieldName={`from_upload_${option.id}`}
            uiName={`${option.label}: ${option.hint}`}
          />
        {/each}
        <div>
          <DSActionButton
            icon="mdi:download"
            action={() => downloadCsvHeaders("po-group-template.csv", columns)}
          >
            Download template CSV ({columns.join(", ")})
          </DSActionButton>
        </div>
        <DsFileSelect bind:record={upload} errors={{}} fieldName="file" uiName="CSV file" />
        {#if csv !== null && csv.missingColumns.length > 0}
          <DSCallout tone="danger" title="Missing columns">
            The CSV needs these columns: {csv.missingColumns.join(", ")}. Download the template to
            see the expected layout.
          </DSCallout>
        {:else if csv !== null && csv.rows.length === 0}
          <DSCallout tone="warning">The CSV has no rows below the header.</DSCallout>
        {/if}
      </div>
    </DSSection>

    {#if otherErrors.length > 0}
      <DSCallout tone="danger" title="Fix the template">
        <ul class="list-disc pl-5">
          {#each otherErrors as [field, error] (field)}
            <li>{errorLabel(field)}: {error.message}</li>
          {/each}
        </ul>
      </DSCallout>
    {/if}

    {#if preview !== null}
      <DSSection title="3. Approvers">
        <div class="flex w-full flex-col gap-2">
          {#if preview.approvers.length === 0}
            <DSCallout tone="warning" title="No one can approve every row">
              Fix any problems in the review below. If the rows are right, no approver's limits
              cover every amount and division, so split the upload into groups that one person can
              approve.
            </DSCallout>
          {:else}
            <DSSelector
              bind:value={template.approver}
              items={preview.approvers}
              {errors}
              fieldName="approver"
              uiName="Approver"
              helperText="Only people who can approve every row are listed. They get one email for the whole group."
            >
              {#snippet optionTemplate(item: PoApproversResponse)}
                {item.given_name} {item.surname}
              {/snippet}
            </DSSelector>
          {/if}
          {#if preview.second_approval_required}
            {#if preview.second_approvers.length === 0}
              <DSCallout tone="warning" title="No one can give second approval">
                Some amounts need a second approver, but no one can give final approval for all of
                them. Split the upload so each group's amounts share a second approver.
              </DSCallout>
            {:else}
              <DSSelector
                bind:value={template.priority_second_approver}
                items={preview.second_approvers}
                {errors}
                fieldName="priority_second_approver"
                uiName="Priority Second Approver"
                clear={true}
                helperText="Some amounts are over the second approval threshold. Only people who can give final approval for all of them are listed, and they also get one email for the whole group."
              >
                {#snippet optionTemplate(item: PoApproversResponse)}
                  {item.given_name} {item.surname}
                {/snippet}
              </DSSelector>
            {/if}
          {/if}
        </div>
      </DSSection>

      <DSSection title="4. Review">
        <div class="flex w-full flex-col gap-2">
          {#if !previewIsCurrent || reviewing}
            <DSCallout tone="info">Checking the latest changes…</DSCallout>
          {:else if preview.error_count > 0}
            <DSCallout
              tone="danger"
              title={`${preview.error_count} of ${preview.count} rows have problems`}
            >
              Fix them in the template or the CSV and upload it again. Nothing is created until
              every row is ready.
            </DSCallout>
          {:else if template.approver === ""}
            <DSCallout tone="info" title="Choose an approver">
              Every row is ready. Choose the approver above to create the purchase orders.
            </DSCallout>
          {:else}
            <DSCallout tone="success" title="Ready to create">
              {readyMessage}
            </DSCallout>
          {/if}

          {#if preview.missing_active_staff.length > 0}
            <DSCollapsible
              title={`${preview.missing_active_staff.length} active staff aren't in the upload`}
              collapsed={true}
            >
              <ul class="columns-1 sm:columns-2">
                {#each preview.missing_active_staff as staff (staff.uid)}
                  <li>{[staff.name, staff.email].filter((part) => part).join(" · ")}</li>
                {/each}
              </ul>
            </DSCollapsible>
          {/if}

          <DSList
            items={preview.rows.map((row) => ({ ...row, id: String(row.index) }))}
            inListHeader={`${preview.count} rows`}
          >
            {#snippet anchor(row)}Line {rowLine(row)}{/snippet}
            {#snippet headline(row)}{row.name || row.email}{/snippet}
            {#snippet byline(row)}
              {formatCurrencyAmount(row.total, preview?.currency_code ?? "CAD")}
            {/snippet}
            {#snippet line1(row)}
              {[row.name ? row.email : "", row.division_code, row.branch_name, row.description]
                .filter((part) => part)
                .join(" · ")}
            {/snippet}
            {#snippet line2(row)}
              {#if row.errors}
                <span class="flex flex-wrap gap-1">
                  {#each Object.entries(row.errors) as [field, message] (field)}
                    <DsLabel tone="danger" size="sm">{errorLabel(field)}: {message}</DsLabel>
                  {/each}
                </span>
              {/if}
            {/snippet}
          </DSList>
        </div>
      </DSSection>
    {/if}

    <div class="flex w-full flex-col gap-2">
      <span class="flex w-full gap-2">
        <DSActionButton action={create} disabled={!canCreate} loading={creating}>
          {preview !== null && preview.count > 0
            ? `Create ${preview.count} purchase orders`
            : "Create purchase orders"}
        </DSActionButton>
        <DSActionButton action="/pos/groups">Cancel</DSActionButton>
      </span>
      {#if errors.global}
        <DSCallout tone="danger">{errors.global.message}</DSCallout>
      {/if}
    </div>
  {/if}
</div>
