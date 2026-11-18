<script lang="ts">
  import Icon from "@iconify/svelte";
  import { onDestroy } from "svelte";
  import { copyText } from "$lib/clipboard";

  // A small icon button that copies `text` to the clipboard. It sits inline
  // beside the value it copies, such as a code, in a sentence, a list row or a
  // DSField, and keeps the line's height. A check mark shows for a moment
  // after a copy, and a screen reader hears "Copied". A browser with no
  // clipboard shows a failure mark instead. It never prints.
  let {
    text,
    label,
  }: {
    text: string;
    // The accessible name and the tooltip, such as "Copy code 2026-02-00001".
    label: string;
  } = $props();

  const resetMs = 1500;
  let status = $state<"idle" | "copied" | "failed">("idle");
  let timer: ReturnType<typeof setTimeout> | undefined;

  const icons = {
    idle: "mdi:clipboard-text-outline",
    copied: "mdi:check",
    failed: "mdi:alert-circle-outline",
  } as const;
  const tones = {
    idle: "text-neutral-500 hover:text-neutral-800",
    copied: "text-green-700",
    failed: "text-red-700",
  } as const;
  const messages = { idle: "", copied: "Copied", failed: "Copy failed" } as const;

  // The button may sit in a row or a card that reacts to a click, so the
  // click stops here.
  async function copy(event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    const copied = await copyText(text);
    status = copied ? "copied" : "failed";
    clearTimeout(timer);
    timer = setTimeout(() => (status = "idle"), resetMs);
  }

  onDestroy(() => clearTimeout(timer));
</script>

<span data-print-hide>
  <!-- The ::before box widens the touch target without changing the line. -->
  <button
    type="button"
    aria-label={label}
    title={status === "idle" ? label : messages[status]}
    data-copy-state={status}
    onclick={copy}
    class="relative inline-flex size-4 cursor-pointer rounded-sm align-[-0.15em] font-normal before:absolute before:-inset-x-1 before:-inset-y-2 before:content-[''] focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-blue-600 {tones[
      status
    ]}"
  >
    <Icon icon={icons[status]} width="16" height="16" aria-hidden="true" />
  </button>
  <span class="sr-only" aria-live="polite">{messages[status]}</span>
</span>
