<script lang="ts">
  import { resolve } from "$app/paths";
  import AbsorbList from "$lib/components/AbsorbList.svelte";
  import DSTabBar from "$lib/components/DSTabBar.svelte";
</script>

<DSTabBar
  tabs={[
    { label: "Pending", href: resolve("/absorb/actions"), active: true },
    { label: "History", href: resolve("/absorb/history"), active: false },
  ]}
/>
<AbsorbList />
