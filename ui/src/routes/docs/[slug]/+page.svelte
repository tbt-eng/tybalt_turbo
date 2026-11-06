<script lang="ts">
  import DSDocPage from "$lib/components/DSDocPage.svelte";
  let { data } = $props();
</script>

<DSDocPage page={data.page} />
