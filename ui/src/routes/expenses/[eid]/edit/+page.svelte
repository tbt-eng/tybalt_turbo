<script lang="ts">
  import type { ExpensesPageData } from "$lib/svelte-types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import ExpensesEditor from "$lib/components/ExpensesEditor.svelte";
  let { data }: { data: ExpensesPageData } = $props();
</script>

<div class="flex justify-end px-2 pt-2">
  <DSHelp doc="expenses#change-an-expense" label="Help" />
</div>
<ExpensesEditor {data} />
