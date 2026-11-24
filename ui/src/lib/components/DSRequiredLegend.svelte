<script lang="ts">
  import DSRequiredMark from "./DSRequiredMark.svelte";

  // The words after the asterisk. A form that saves an incomplete draft says
  // what the fields are required for, such as "Required to submit".
  let { label = "Required" }: { label?: string } = $props();
</script>

<!-- The one-line key to DSRequiredMark, near the top of a form that marks its
required fields. -->
<p class="text-sm text-neutral-600"><DSRequiredMark /> {label}</p>
