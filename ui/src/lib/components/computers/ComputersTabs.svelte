<script lang="ts">
  import { resolve } from "$app/paths";
  import DSTabBar from "$lib/components/DSTabBar.svelte";

  // The tabs of the Computers pages.
  let { active }: { active: "list" | "retired" | "logins" | "raw_logins" } = $props();
</script>

<DSTabBar
  tabs={[
    { label: "Computers", href: resolve("/computers/list"), active: active === "list" },
    { label: "Retired", href: resolve("/computers/retired"), active: active === "retired" },
    { label: "Logins", href: resolve("/computers/logins"), active: active === "logins" },
    {
      label: "Raw Logins",
      href: resolve("/computers/raw_logins"),
      active: active === "raw_logins",
    },
  ]}
/>
