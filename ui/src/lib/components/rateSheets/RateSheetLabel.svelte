<script lang="ts">
  import { resolve } from "$app/paths";
  import DSCopyButton from "$lib/components/DSCopyButton.svelte";
  import DSRecordLabel from "$lib/components/DSRecordLabel.svelte";
  import {
    rateSheetCodeRevision,
    rateSheetCopyText,
    type RateSheetLabelFields,
  } from "$lib/rateSheetLabel";

  // A rate sheet as Turbo shows it: its code and revision, linked to the
  // sheet's page, then its name in muted text when it has one. Use
  // rateSheetLabel() from $lib/rateSheetLabel where only text fits.
  let {
    sheet,
    link = true,
    copy = false,
  }: {
    sheet: RateSheetLabelFields & { id?: string };
    // Link the code to the sheet's page. A request's proposed rates have no
    // page until the job is created.
    link?: boolean;
    // Show a button after the code and revision that copies the code alone,
    // to paste in a rate sheet search.
    copy?: boolean;
  } = $props();

  const code = $derived(rateSheetCopyText(sheet));
</script>

{#snippet copyButton()}
  <DSCopyButton text={code} label={`Copy code ${code}`} />
{/snippet}

<DSRecordLabel
  label={rateSheetCodeRevision(sheet)}
  description={sheet.code ? sheet.name : ""}
  href={link && sheet.id ? resolve("/rate-sheets/[id]/details", { id: sheet.id }) : undefined}
  trailing={copy && code ? copyButton : undefined}
/>
