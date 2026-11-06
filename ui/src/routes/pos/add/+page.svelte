<script lang="ts">
  import type { PageData } from "./$types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import PurchaseOrdersEditor from "$lib/components/PurchaseOrdersEditor.svelte";
  let { data }: { data: PageData } = $props();
</script>

<PurchaseOrdersEditor {data}>
  {#snippet help()}
    <DSHelp doc="purchase-orders#create-a-po" label="Help" />
  {/snippet}
</PurchaseOrdersEditor>
