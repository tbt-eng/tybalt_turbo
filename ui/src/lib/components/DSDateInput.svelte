<script lang="ts">
  // A date input without a label; the caller writes the <label>. With
  // `required`, the input sets aria-required and the caller puts DSRequiredMark
  // in its label. It never sets the native required attribute, which would
  // block the save with a browser bubble; the server validates.
  let {
    value = $bindable(""),
    id,
    name,
    class: className = "",
    min = undefined,
    max = undefined,
    disabled = false,
    required = false,
  }: {
    value: string;
    id?: string;
    name?: string;
    class?: string;
    min?: string;
    max?: string;
    disabled?: boolean;
    required?: boolean;
  } = $props();
</script>

<input
  type="date"
  class={className}
  id={id ?? name}
  {name}
  min={min ?? undefined}
  max={max ?? undefined}
  bind:value
  {disabled}
  aria-required={required || undefined}
/>
