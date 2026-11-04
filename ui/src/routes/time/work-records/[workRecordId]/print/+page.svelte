<script lang="ts">
  import { onMount } from "svelte";
  import type { PageData } from "./$types";
  import { formatDateTime, shortDate, trimmedOrEmpty } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  const item = $derived(data.item);
  const totalWorkerHours = $derived(
    item.subjects.reduce(
      (total, worker) => total + worker.hours_on_site + worker.hours_travel_time,
      0,
    ),
  );
  const activeConsumables = $derived(
    item.consumables.filter((c) =>
      c.input_kind === "boolean" ? c.boolean_selected : c.quantity_number > 0,
    ),
  );
  const consumablePairs = $derived.by(() => {
    const pairs: Array<
      [(typeof activeConsumables)[number] | null, (typeof activeConsumables)[number] | null]
    > = [];
    for (let i = 0; i < activeConsumables.length; i += 2) {
      pairs.push([activeConsumables[i], activeConsumables[i + 1] ?? null]);
    }
    return pairs;
  });
  const printableNotes = $derived(item.notes.filter((note) => trimmedOrEmpty(note.note) !== ""));

  function displayValue(value: string | number | null | undefined, fallback = ""): string {
    if (value === null || value === undefined) return fallback;
    if (typeof value === "number") return value === 0 ? fallback : String(value);
    return trimmedOrEmpty(value) || fallback;
  }

  function consumableValue(c: (typeof activeConsumables)[number]): string {
    if (c.input_kind === "boolean") return "Yes";
    const unit = trimmedOrEmpty(c.unit_label);
    return unit ? `${c.quantity_number} ${unit}` : String(c.quantity_number);
  }

  function workerVehicle(worker: (typeof item.subjects)[number]): string {
    const bits = [
      displayValue(worker.vehicle_type),
      worker.is_passenger ? "passenger" : "",
      displayValue(worker.company_vehicle_unit_number)
        ? `unit ${worker.company_vehicle_unit_number}`
        : "",
      worker.distance_travelled_km ? `${worker.distance_travelled_km} km` : "",
    ].filter(Boolean);
    return bits.join(" / ");
  }

  onMount(() => {
    document.title = item.work_record ? `Work Record ${item.work_record}` : "Work Record";

    const timer = setTimeout(() => window.print(), 300);
    return () => clearTimeout(timer);
  });
</script>

<svelte:head>
  <meta name="robots" content="noindex" />
</svelte:head>

<div class="screen-only mx-auto max-w-5xl px-4 pt-4 pb-0">
  <div
    class="rounded-sm border border-neutral-300 bg-neutral-100 px-4 py-2 text-sm text-neutral-700"
  >
    The print dialog should open automatically. If it does not, use your browser's print command to
    print this work record or save it as a PDF.
  </div>
</div>

<article class="print-page mx-auto my-4 max-w-5xl bg-white p-8 text-black">
  <div style="overflow: hidden; border-bottom: 2px solid black; padding-bottom: 0.75rem;">
    <table style="border-collapse: collapse;">
      <tbody>
        <tr>
          <td style="vertical-align: middle; padding-right: 0.75rem;">
            <img src="/logo.svg" alt="TBT Engineering Limited" style="width: 3rem; height: 3rem;" />
          </td>
          <td style="vertical-align: middle;">
            <div style="font-size: 1.1rem; font-weight: 700; letter-spacing: 0.03em;">
              TBT ENGINEERING
            </div>
            <div
              style="font-size: 0.7rem; font-weight: 700; letter-spacing: 0.15em; text-transform: uppercase; color: #525252;"
            >
              Consulting Group
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div style="overflow: hidden; margin-top: 0.5rem;">
      <div style="float: left; font-size: 1.6rem; font-weight: 700;">
        {displayValue(item.work_record, "-")}
      </div>
      <div
        style="float: right; font-size: 1.1rem; font-weight: 700; letter-spacing: 0.05em; padding-top: 0.4rem;"
      >
        Work Record
      </div>
    </div>
  </div>

  <table style="width: 100%; border-collapse: collapse; margin-top: 0.75rem;">
    <colgroup>
      <col style="width: 14%;" />
      <col style="width: 36%;" />
      <col style="width: 14%;" />
      <col style="width: 36%;" />
    </colgroup>
    <tbody>
      <tr>
        <td class="label-cell">Job Number</td>
        <td class="value-cell">{displayValue(item.job_number)}</td>
        <td class="label-cell">Date</td>
        <td class="value-cell">{shortDate(item.date, true)}</td>
      </tr>
      <tr>
        <td class="label-cell">Client</td>
        <td class="value-cell">{displayValue(item.client_name)}</td>
        <td class="label-cell">Workers</td>
        <td class="value-cell"
          >{item.subjects.map((worker) => worker.worker_name || worker.uid).join(", ")}</td
        >
      </tr>
      <tr>
        <td class="label-cell">Client Project No</td>
        <td class="value-cell">{displayValue(item.client_reference_number)}</td>
        <td class="label-cell">Work Performed</td>
        <td class="value-cell">{displayValue(item.type_name)}</td>
      </tr>
      <tr>
        <td class="label-cell">Project Desc</td>
        <td colspan="3" class="value-cell">{displayValue(item.job_description)}</td>
      </tr>
      <tr>
        <td class="label-cell">Location</td>
        <td colspan="3" class="value-cell">{displayValue(item.location)}</td>
      </tr>
      <tr>
        <td class="label-cell">Report To</td>
        <td class="value-cell">{displayValue(item.report_to)}</td>
        <td class="label-cell">Total Worker Hours</td>
        <td class="value-cell">{totalWorkerHours ? `${totalWorkerHours} hrs` : ""}</td>
      </tr>
      <tr>
        <td class="label-cell">Sub Contractor</td>
        <td class="value-cell">{displayValue(item.sub_contractor)}</td>
        <td class="label-cell">Lane km</td>
        <td class="value-cell">{item.lane_kms ? `${item.lane_kms} km` : ""}</td>
      </tr>
      <tr>
        <td class="label-cell">Client PO</td>
        <td class="value-cell">{displayValue(item.client_po)}</td>
        <td class="label-cell">Field Book</td>
        <td class="value-cell">
          {displayValue(item.field_book_number)}
          {#if item.field_book_page_number}
            / page {item.field_book_page_number}
          {/if}
        </td>
      </tr>
    </tbody>
  </table>

  <div class="print-section">
    <div class="section-label">Workers</div>
    <table style="width: 100%; border-collapse: collapse; margin-top: 0.4rem;">
      <thead>
        <tr>
          <th class="worker-heading">Worker</th>
          <th class="worker-heading numeric">On Site</th>
          <th class="worker-heading numeric">Travel</th>
          <th class="worker-heading numeric">Total</th>
          <th class="worker-heading">Vehicle Usage</th>
        </tr>
      </thead>
      <tbody>
        {#each item.subjects as worker (worker.id)}
          <tr>
            <td class="worker-cell">{displayValue(worker.worker_name || worker.uid)}</td>
            <td class="worker-cell numeric">{worker.hours_on_site || ""}</td>
            <td class="worker-cell numeric">{worker.hours_travel_time || ""}</td>
            <td class="worker-cell numeric">
              {worker.hours_on_site + worker.hours_travel_time || ""}
            </td>
            <td class="worker-cell">{workerVehicle(worker)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <div class="print-section">
    <div class="section-label">Details</div>
    <div class="body-text" style="min-height: 2rem;">
      {displayValue(item.work_description, "No work description provided.")}
    </div>
  </div>

  <div class="print-section">
    <div class="section-label">Equipment</div>
    <div class="body-text" style="min-height: 1rem;">{displayValue(item.equipment)}</div>
  </div>

  <div class="print-section">
    <div class="section-label">Supplies</div>
    <div class="body-text" style="min-height: 1rem;">{displayValue(item.supplies)}</div>
  </div>

  {#if activeConsumables.length > 0}
    <div class="print-section">
      <div class="section-label">Consumables</div>
      <table style="width: 100%; border-collapse: collapse; margin-top: 0.4rem;">
        <colgroup>
          <col style="width: 30%;" />
          <col style="width: 20%;" />
          <col style="width: 30%;" />
          <col style="width: 20%;" />
        </colgroup>
        <tbody>
          {#each consumablePairs as [left, right]}
            <tr>
              <td class="consumable-name">{left?.consumable_name ?? ""}</td>
              <td class="consumable-value">{left ? consumableValue(left) : ""}</td>
              <td class="consumable-name">{right?.consumable_name ?? ""}</td>
              <td class="consumable-value">{right ? consumableValue(right) : ""}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  {#if printableNotes.length > 0}
    <div class="print-section">
      <div class="section-label">Notes</div>
      {#each printableNotes as note, i (note.id)}
        <div class="note-item" style={i > 0 ? "border-top: 1px solid #e5e5e5;" : ""}>
          <div class="note-meta">
            <span>{trimmedOrEmpty(note.author_name) || note.uid}</span>
            <span style="float: right;">{formatDateTime(note.created)}</span>
          </div>
          <div
            style="margin-top: 0.2rem; white-space: pre-wrap; font-size: 0.8rem; line-height: 1.5;"
          >
            {note.note}
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <div
    style="margin-top: 1.5rem; padding-top: 0.75rem; border-top: 2px solid black; overflow: hidden;"
  >
    <div style="float: right; font-size: 0.85rem;">
      <span style="font-weight: 700;">Work Record Number:</span>
      {displayValue(item.work_record, "-")}
    </div>
    <div style="font-size: 0.7rem; color: #525252;">
      TBT Engineering Limited &nbsp; 1918 Yonge Street, Thunder Bay, ON &nbsp; P7E 6T9 &nbsp; (807)
      624-5100
    </div>
  </div>
</article>

<style>
  @page {
    margin: 1.25cm;
    size: auto;
  }

  .print-page {
    box-shadow: none;
  }

  .section-label {
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.15em;
    color: #525252;
  }

  .label-cell {
    padding: 0.35rem 0.5rem;
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #525252;
    vertical-align: bottom;
  }

  .value-cell {
    padding: 0.35rem 0.5rem;
    font-size: 0.85rem;
    vertical-align: bottom;
    border-bottom: 1px solid black;
  }

  .print-section {
    border: 1px solid black;
    margin-top: 1rem;
    padding: 0.75rem 1rem;
    break-inside: avoid;
  }

  .body-text {
    margin-top: 0.4rem;
    white-space: pre-wrap;
    line-height: 1.6;
    font-size: 0.85rem;
  }

  .worker-heading,
  .worker-cell,
  .consumable-name,
  .consumable-value {
    padding: 0.3rem 0.5rem;
    font-size: 0.78rem;
    border-bottom: 1px solid #e5e5e5;
  }

  .worker-heading {
    text-align: left;
    font-weight: 700;
    color: #525252;
  }

  .numeric,
  .consumable-value {
    text-align: right;
  }

  .consumable-name {
    font-weight: 600;
  }

  .note-item {
    margin-top: 0.5rem;
    padding-top: 0.5rem;
    overflow: hidden;
  }

  .note-meta {
    font-size: 0.7rem;
    font-weight: 600;
    color: #525252;
    overflow: hidden;
  }

  @media print {
    .screen-only {
      display: none !important;
    }

    .print-page {
      margin: 0 !important;
      max-width: 100% !important;
      padding: 0.5cm !important;
      background: white !important;
    }
  }
</style>
