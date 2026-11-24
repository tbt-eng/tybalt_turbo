<script module>
  let idCounter = $state(0);
</script>

<script lang="ts" generics="T">
  import DSRequiredMark from "./DSRequiredMark.svelte";

  // get an id for this instance from the counter in the module context then
  // increment it so the next instance gets a different id
  const thisId = idCounter;
  idCounter += 1;

  let {
    value = $bindable(),
    type = "text",
    step = undefined,
    min = undefined,
    max = undefined,
    errors,
    fieldName,
    uiName,
    placeholder,
    disabled = false,
    rows = undefined, // a text box of this many lines, with the label above it
    stacked = false, // a bold label above a full-width box
    maxlength = undefined,
    helperText = "", // a line under the field explaining it
    required = false, // a red asterisk after the label; the server still validates
  }: {
    value: string | number;
    type?: "text" | "number" | "password";
    step?: number;
    min?: number;
    max?: number;
    errors: Record<string, { message: string }>;
    fieldName: string;
    uiName: string;
    placeholder?: string;
    disabled?: boolean;
    rows?: number;
    stacked?: boolean;
    maxlength?: number;
    helperText?: string;
    required?: boolean;
  } = $props();
</script>

<div class="flex w-full flex-col gap-2 {errors[fieldName] !== undefined ? 'bg-red-200' : ''}">
  <span class="flex w-full {rows ? 'flex-col gap-2' : stacked ? 'flex-col gap-1' : 'gap-2'}">
    <label class={stacked ? "text-sm font-semibold" : ""} for={`text-input-${thisId}`}
      >{uiName}{#if required}<DSRequiredMark />{/if}</label
    >
    {#if rows}
      <textarea
        class="w-full rounded-sm border border-neutral-300 p-2 {disabled
          ? 'cursor-not-allowed opacity-50'
          : ''}"
        {rows}
        maxlength={maxlength ?? null}
        id={`text-input-${thisId}`}
        name={fieldName}
        placeholder={placeholder || uiName}
        bind:value
        aria-required={required || undefined}
        {disabled}></textarea>
    {:else}
      <input
        class="flex-1 rounded border border-neutral-300 {stacked ? 'px-2 py-1' : 'px-1'} {disabled
          ? 'opacity-50'
          : ''} {disabled ? 'cursor-not-allowed' : ''}"
        {type}
        step={step || null}
        min={min || null}
        max={max || null}
        maxlength={maxlength ?? null}
        id={`text-input-${thisId}`}
        name={fieldName}
        placeholder={placeholder || uiName}
        bind:value
        aria-required={required || undefined}
        {disabled}
      />
    {/if}
  </span>
  {#if errors[fieldName] !== undefined}
    <span class="text-red-600">{errors[fieldName].message}</span>
  {/if}
  {#if helperText}
    <span class="text-sm text-neutral-600">{helperText}</span>
  {/if}
</div>
