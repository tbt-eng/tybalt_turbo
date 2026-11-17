<script module>
  let idCounter = $state(0);
</script>

<script lang="ts" generics="T">
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
    maxlength = undefined,
    helperText = "", // a line under the field explaining it
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
    maxlength?: number;
    helperText?: string;
  } = $props();
</script>

<div class="flex w-full flex-col gap-2 {errors[fieldName] !== undefined ? 'bg-red-200' : ''}">
  <span class="flex w-full gap-2 {rows ? 'flex-col' : ''}">
    <label for={`text-input-${thisId}`}>{uiName}</label>
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
        {disabled}></textarea>
    {:else}
      <input
        class="flex-1 rounded border border-neutral-300 px-1 {disabled
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
