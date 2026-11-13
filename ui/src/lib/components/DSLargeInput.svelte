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
    onfocus,
    actions,
    suggestion,
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
    onfocus?: () => void;
    // Buttons at the right end, inside the box.
    actions?: Snippet;
    // An example of what to type. It types itself into the empty box, in place
    // of the placeholder, a moment after the box shows. Tab puts it in the box.
    suggestion?: string;
  } = $props();

  const typeAfter = 2000;
  const typeEvery = 45;

  // The number of the suggestion's characters that show.
  let typed = $state(0);
  let focused = $state(false);
  const hint = $derived(suggestion && value === "" ? suggestion.slice(0, typed) : "");

  $effect(() => {
    const text = suggestion;
    typed = 0;
    if (!text) return;
    // A person who asked for less motion gets the text with no typing.
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      typed = text.length;
      return;
    }
    let count = 0;
    let timer = setTimeout(function step() {
      typed = ++count;
      if (count < text.length) timer = setTimeout(step, typeEvery);
    }, typeAfter);
    return () => clearTimeout(timer);
  });

  function fill() {
    if (!suggestion) return;
    value = suggestion;
    element?.focus();
  }

  function keydown(event: KeyboardEvent) {
    // Tab fills the box only while the suggestion shows. The next Tab moves on.
    const plainTab =
      event.key === "Tab" && !event.shiftKey && !event.altKey && !event.ctrlKey && !event.metaKey;
    if (hint && plainTab) {
      event.preventDefault();
      fill();
      return;
    }
    onkeydown?.(event);
  }
</script>

<div
  class="flex items-center gap-2 rounded-sm border border-neutral-400 bg-white px-3 py-2 focus-within:border-blue-700 focus-within:ring-1 focus-within:ring-blue-700"
>
  <Icon {icon} width="24px" class="shrink-0 text-neutral-500" />
  <div class="relative flex min-w-0 flex-1">
    <input
      bind:this={element}
      type="text"
      {enterkeyhint}
      autocomplete="off"
      {spellcheck}
      aria-label={label}
      placeholder={suggestion ? "" : label}
      class="min-w-0 flex-1 bg-transparent text-lg outline-none placeholder:text-neutral-500"
      {value}
      oninput={(event) => {
        value = event.currentTarget.value;
        oninput?.(event);
      }}
      onkeydown={keydown}
      onfocus={() => {
        focused = true;
        onfocus?.();
      }}
      onblur={() => (focused = false)}
    />
    {#if hint}
      <!-- The box's own caret shows while it has the focus. -->
      <div
        aria-hidden="true"
        class="pointer-events-none absolute inset-0 flex items-center text-lg text-neutral-500"
      >
        <span class="truncate">{hint}</span>
        {#if !focused}
          <span class="ml-px h-5 w-px shrink-0 animate-pulse bg-neutral-500"></span>
        {/if}
      </div>
    {/if}
  </div>
  {@render actions?.()}
</div>
