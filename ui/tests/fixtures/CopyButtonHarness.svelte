<script lang="ts">
  import DSCopyButton from "$lib/components/DSCopyButton.svelte";
  import DSRecordLabel from "$lib/components/DSRecordLabel.svelte";
  let rowClicks = $state(0);
</script>

<!-- A row that reacts to a click, as some lists do. -->
<!-- svelte-ignore a11y_click_events_have_key_events, a11y_no_static_element_interactions -->
<div data-row onclick={() => rowClicks++}>
  <div data-plain>2026-02-00001 (rev. 0) Client Rates</div>
  <div data-label>
    <DSRecordLabel label="2026-02-00001 (rev. 0)" description="Client Rates" href="#sheet">
      {#snippet trailing()}<DSCopyButton
          text="2026-02-00001"
          label="Copy code 2026-02-00001"
        />{/snippet}
    </DSRecordLabel>
  </div>
</div>
<output data-row-clicks>{rowClicks}</output>
