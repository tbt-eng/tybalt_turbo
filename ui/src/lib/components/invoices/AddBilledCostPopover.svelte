<script lang="ts">
  import DSPopover from "$lib/components/DSPopover.svelte";
  import DSSelector from "$lib/components/DSSelector.svelte";
  import DSTextInput from "$lib/components/DSTextInput.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import { formatCurrencyAmount, saveErrorMessage, shortDate } from "$lib/utilities";
  import {
    billCost,
    costLabel,
    loadBillableCosts,
    loadBillableJobs,
    type BillableJob,
    type InvoiceCost,
  } from "$lib/invoices";

  // Bills one of the client's unbilled costs on an issued invoice: pick a job,
  // then one of its committed expenses or Active or Closed POs, and enter the
  // amount billed, net of tax.
  let {
    show = $bindable(),
    invoiceId,
    onAdded,
  }: {
    show: boolean;
    invoiceId: string;
    onAdded: () => void;
  } = $props();

  let jobs = $state([] as BillableJob[]);
  let costs = $state([] as (InvoiceCost & { id: string })[]);
  let jobId = $state("");
  let costId = $state("");
  let amount = $state("" as string | number);
  let error = $state<string | null>(null);
  let submitting = $state(false);

  const cost = $derived(costs.find((c) => c.id === costId));

  $effect(() => {
    if (!show) return;
    jobId = "";
    costId = "";
    amount = "";
    error = null;
    costs = [];
    loadBillableJobs(invoiceId)
      .then((rows) => (jobs = rows))
      .catch((err) => (error = saveErrorMessage(err, "Could not load the client's jobs.")));
  });

  async function chooseJob(id: string) {
    costId = "";
    costs = [];
    if (!id) return;
    try {
      // DSSelector needs an id; a cost is unique by kind and id.
      costs = (await loadBillableCosts(invoiceId, id)).map((c) => ({
        ...c,
        id: `${c.kind}:${c.cost_id}`,
      }));
    } catch (err) {
      error = saveErrorMessage(err, "Could not load the job's costs.");
    }
  }

  $effect(() => {
    chooseJob(jobId);
  });

  async function submit() {
    if (!cost) return;
    submitting = true;
    error = null;
    try {
      await billCost(invoiceId, cost.kind, cost.cost_id, Number(amount));
      show = false;
      onAdded();
    } catch (err) {
      error = saveErrorMessage(err, "The cost could not be billed.");
    } finally {
      submitting = false;
    }
  }
</script>

<DSPopover
  bind:show
  title="Add a billed-back cost"
  subtitle="Pick an unbilled cost on one of this client's jobs, and enter what was billed for it in Canadian dollars, before tax."
  submitLabel="Add"
  {error}
  {submitting}
  submitDisabled={!cost || !(Number(amount) > 0)}
  onSubmit={submit}
>
  {#if jobs.length === 0 && !error}
    <p>No job of this client has an unbilled committed expense or Active or Closed PO.</p>
  {:else}
    <DSSelector bind:value={jobId} items={jobs} errors={{}} fieldName="job" uiName="Job">
      {#snippet optionTemplate(job)}
        {job.number} · {job.description} ({job.billable_count})
      {/snippet}
    </DSSelector>
  {/if}

  {#if jobId && costs.length > 0}
    <DSSelector bind:value={costId} items={costs} errors={{}} fieldName="cost" uiName="Cost">
      {#snippet optionTemplate(c)}
        {costLabel(c)} · {formatCurrencyAmount(c.cost_total, c.currency_code)}
      {/snippet}
    </DSSelector>
  {/if}

  {#if cost}
    <DSFieldList>
      <DSField label="Vendor" value={cost.vendor_name} />
      <DSField label="Date" value={cost.date ? shortDate(cost.date, true) : ""} />
      <DSField label="Cost" value={formatCurrencyAmount(cost.cost_total, cost.currency_code)} />
    </DSFieldList>
    <DSTextInput
      bind:value={amount}
      type="number"
      step={0.01}
      min={0}
      errors={{}}
      fieldName="amount"
      uiName="Amount billed"
    />
  {/if}
</DSPopover>
