<script lang="ts">
  import type { ExpensesListData } from "$lib/svelte-types";
  import DSHelp from "$lib/components/DSHelp.svelte";
  import ExpensesList from "$lib/components/ExpensesList.svelte";
  let { data }: { data: ExpensesListData } = $props();
</script>

<div class="flex justify-end px-2 pt-2">
  <DSHelp doc="expenses#approve-an-expense" label="Help" />
</div>
<ExpensesList
  inListHeader="Expenses Pending My Approval"
  {data}
  endpoint="/api/expenses/pending"
  showInlineAllowanceApprove
/>
