<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { Snippet } from "svelte";

  // A large one-line input with an icon, for the main thing to type on a page,
  // such as a search or a question.
  let {
    value = $bindable(""),
    label,
    icon = "mdi:magnify",
    enterkeyhint = "search",
    spellcheck = false,
    element = $bindable(),
    oninput,
    onkeydown,
    actions,
  }: {
    value?: string;
    // The placeholder, and the name that a screen reader says.
    label: string;
    icon?: string;
    // The label of the Enter key on a phone's keyboard.
    enterkeyhint?: "search" | "send" | "go" | "done";
    spellcheck?: boolean;
    // The input element, for a caller that moves the focus to it.
    element?: HTMLInputElement;
    // Runs after value holds the new text.
    oninput?: (event: Event & { currentTarget: HTMLInputElement }) => void;
    onkeydown?: (event: KeyboardEvent) => void;
    // Buttons at the right end, inside the box.
    actions?: Snippet;
  } = $props();
</script>

<div
  class="flex items-center gap-2 rounded-sm border border-neutral-400 bg-white px-3 py-2 focus-within:border-blue-700 focus-within:ring-1 focus-within:ring-blue-700"
>
  <Icon {icon} width="24px" class="shrink-0 text-neutral-500" />
  <input
    bind:this={element}
    type="text"
    {enterkeyhint}
    autocomplete="off"
    {spellcheck}
    aria-label={label}
    placeholder={label}
    class="min-w-0 flex-1 bg-transparent text-lg outline-none placeholder:text-neutral-500"
    {value}
    oninput={(event) => {
      value = event.currentTarget.value;
      oninput?.(event);
    }}
    {onkeydown}
  />
  {@render actions?.()}
</div>
