<script lang="ts">
  import DSHelp from "$lib/components/DSHelp.svelte";
  import type { ExpensesListData } from "$lib/svelte-types";
  import ExpensesList from "$lib/components/ExpensesList.svelte";
  let { data }: { data: ExpensesListData } = $props();
</script>

<ExpensesList inListHeader="Expenses Approved by Me" {data} endpoint="/api/expenses/approved">
  {#snippet searchBarExtra()}
    <DSHelp doc="commit-expenses#see-the-expenses-you-approved" label="Help" />
  {/snippet}
</ExpensesList>
