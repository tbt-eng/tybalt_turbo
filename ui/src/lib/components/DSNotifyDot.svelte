<script lang="ts">
  // A small red dot beside a link, for something new that the person hasn't opened.
  let { label }: { label: string } = $props();
</script>

<span
  role="img"
  aria-label={label}
  title={label}
  class="inline-block h-2 w-2 shrink-0 rounded-full bg-red-500"
></span>
