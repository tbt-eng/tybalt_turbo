<script lang="ts">
  import type { ExpensesPageData } from "$lib/svelte-types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import ExpensesEditor from "$lib/components/ExpensesEditor.svelte";
  let { data }: { data: ExpensesPageData } = $props();
</script>

<ExpensesEditor {data}>
  {#snippet help()}
    <DSHelp doc="expenses#enter-an-expense-without-a-po" label="Help" />
  {/snippet}
</ExpensesEditor>
