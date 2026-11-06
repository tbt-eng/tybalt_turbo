<script lang="ts">
  import type { PageData } from "./$types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import PurchaseOrdersEditor from "$lib/components/PurchaseOrdersEditor.svelte";
  let { data }: { data: PageData } = $props();
</script>

<div class="flex justify-end px-2 pt-2">
  <DSHelp doc="purchase-orders#change-a-po" label="Help" />
</div>
<PurchaseOrdersEditor {data} />
