<script lang="ts">
  import { onMount } from "svelte";
  import type { PageData } from "./$types";
  import { shortDate, trimmedOrEmpty } from "$lib/utilities";

  let { data }: { data: PageData } = $props();

  const item = $derived(data.item);
  const statusText = $derived(item.approved ? "Approved" : "Open");
  const totalHours = $derived(item.hours_on_site + item.hours_travel_time);
  const fieldBookNumber = $derived(trimmedOrEmpty(item.field_book_number));
  const activeConsumables = $derived(
    item.consumables.filter((c) =>
      c.input_kind === "boolean" ? c.boolean_selected : c.quantity_number > 0,
    ),
  );
  const printableNotes = $derived(item.notes.filter((note) => trimmedOrEmpty(note.note) !== ""));

  type FieldRow = {
    label: string;
    value: string;
  };

  function optionalRow(label: string, value: string | null | undefined): FieldRow | null {
    const normalized = trimmedOrEmpty(value);
    if (normalized === "") return null;
    return { label, value: normalized };
  }

  function numericRow(label: string, value: number, suffix = ""): FieldRow | null {
    if (!Number.isFinite(value)) return null;
    return { label, value: `${value}${suffix}` };
  }

  const summaryRows = $derived(
    [
      optionalRow("Worker", item.owner_name),
      optionalRow("Type", item.type_name),
      optionalRow("Status", statusText),
      optionalRow("Job Number", item.job_number),
      optionalRow("Client", item.client_name),
      optionalRow("Project", item.job_description),
    ].filter((row): row is FieldRow => row !== null),
  );

  const clientContextRows = $derived(
    [
      optionalRow("Client Reference", item.client_reference_number),
      optionalRow("Client PO", item.client_po),
    ].filter((row): row is FieldRow => row !== null),
  );

  const fieldRows = $derived(
    [
      optionalRow("Location", item.location),
      optionalRow("Report To", item.report_to),
      optionalRow("Sub Contractor", item.sub_contractor),
      optionalRow("Field Book", fieldBookNumber),
      optionalRow(
        "Field Book Page",
        fieldBookNumber !== "" && item.field_book_page_number > 0
          ? String(item.field_book_page_number)
          : "",
      ),
      optionalRow("Vehicle", item.vehicle_type),
      optionalRow("Passenger", item.is_passenger ? "Yes" : ""),
      optionalRow("Company Unit", item.company_vehicle_unit_number),
      numericRow("Mileage", item.distance_travelled_km, " km"),
      numericRow("On-Site Time", item.hours_on_site, " hours"),
      numericRow("Travel Time", item.hours_travel_time, " hours"),
      numericRow("Total Time", totalHours, " hours"),
    ].filter((row): row is FieldRow => row !== null),
  );

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
  <header class="border-b-2 border-black pb-4">
    <div class="header-top">
      <div class="brand-block">
        <img src="/logo.svg" alt="TBT Engineering Limited" class="logo" />
        <div>
          <div class="company-name">TBT Engineering Limited</div>
          <div class="document-type">Work Record</div>
        </div>
      </div>
      <div class="header-meta">
        <div class="meta-label">Date</div>
        <div>{shortDate(item.date, true)}</div>
      </div>
    </div>
    <div class="record-number">{item.work_record}</div>
  </header>

  <section class="section-grid">
    <section class="print-card">
      <div class="section-label">Summary</div>
      <dl class="row-list">
        {#each summaryRows as row}
          <div class="row-item">
            <dt>{row.label}</dt>
            <dd>{row.value}</dd>
          </div>
        {/each}
      </dl>
    </section>

    {#if clientContextRows.length > 0}
      <section class="print-card">
        <div class="section-label">Job / Client Context</div>
        <dl class="row-list">
          {#each clientContextRows as row}
            <div class="row-item">
              <dt>{row.label}</dt>
              <dd>{row.value}</dd>
            </div>
          {/each}
        </dl>
      </section>
    {/if}
  </section>

  <section class="print-card">
    <div class="section-label">Field / Work</div>
    <dl class="detail-grid">
      {#each fieldRows as row}
        <div class="row-item">
          <dt>{row.label}</dt>
          <dd>{row.value}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <section class="print-card">
    <div class="section-label">Work Description</div>
    <div class="body-copy">{trimmedOrEmpty(item.work_description) || "No work description provided."}</div>
  </section>

  {#if activeConsumables.length > 0}
    <section class="print-card">
      <div class="section-label">Consumables</div>
      <div class="chips">
        {#each activeConsumables as consumable (consumable.consumable_id)}
          <div class="chip">
            <span class="chip-label">{consumable.consumable_name}</span>
            <span class="chip-value">
              {#if consumable.input_kind === "boolean"}
                Recorded
              {:else}
                {consumable.quantity_number}{#if trimmedOrEmpty(consumable.unit_label)}
                  &nbsp;{consumable.unit_label}
                {/if}
              {/if}
            </span>
          </div>
        {/each}
      </div>
    </section>
  {/if}

  {#if printableNotes.length > 0}
    <section class="print-card">
      <div class="section-label">Notes</div>
      <div class="notes-list">
        {#each printableNotes as note (note.id)}
          <div class="note-item">
            <div class="note-meta">
              <span>{trimmedOrEmpty(note.author_name) || note.uid}</span>
              <span>{note.created}</span>
            </div>
            <div class="body-copy">{note.note}</div>
          </div>
        {/each}
      </div>
    </section>
  {/if}
</article>

<style>
  @page {
    margin: 1.25cm;
    size: auto;
  }

  .print-page {
    box-shadow: none;
  }

  .header-top {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: flex-start;
  }

  .brand-block {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .logo {
    width: 3rem;
    height: 3rem;
  }

  .company-name {
    font-size: 1.25rem;
    font-weight: 700;
    letter-spacing: 0.02em;
  }

  .document-type,
  .section-label,
  .meta-label {
    text-transform: uppercase;
    letter-spacing: 0.18em;
    font-size: 0.72rem;
    font-weight: 700;
    color: #525252;
  }

  .header-meta {
    min-width: 8rem;
    text-align: right;
    font-size: 0.95rem;
  }

  .record-number {
    margin-top: 0.75rem;
    font-size: 2rem;
    font-weight: 700;
  }

  .section-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }

  .print-card {
    border: 1px solid black;
    padding: 1rem;
    margin-top: 1rem;
    break-inside: avoid;
  }

  .row-list,
  .detail-grid {
    margin-top: 0.75rem;
    display: grid;
    gap: 0.75rem;
  }

  .detail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .row-item {
    display: grid;
    gap: 0.2rem;
  }

  dt {
    font-size: 0.8rem;
    font-weight: 700;
    color: #404040;
  }

  dd {
    margin: 0;
    font-size: 0.95rem;
    line-height: 1.4;
  }

  .body-copy {
    margin-top: 0.75rem;
    white-space: pre-wrap;
    line-height: 1.6;
    font-size: 0.95rem;
  }

  .chips {
    margin-top: 0.75rem;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 0.75rem;
  }

  .chip {
    border: 1px solid #d4d4d4;
    padding: 0.75rem;
    display: flex;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .chip-label {
    font-weight: 600;
  }

  .chip-value {
    text-align: right;
  }

  .notes-list {
    margin-top: 0.75rem;
    display: grid;
    gap: 0.75rem;
  }

  .note-item {
    border-top: 1px solid #d4d4d4;
    padding-top: 0.75rem;
  }

  .note-item:first-child {
    border-top: none;
    padding-top: 0;
  }

  .note-meta {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    font-size: 0.8rem;
    font-weight: 600;
    color: #525252;
  }

  @media (max-width: 768px) {
    .section-grid,
    .detail-grid {
      grid-template-columns: 1fr;
    }

    .header-top,
    .note-meta,
    .chip {
      display: block;
    }

    .header-meta,
    .chip-value {
      margin-top: 0.5rem;
      text-align: left;
    }
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
