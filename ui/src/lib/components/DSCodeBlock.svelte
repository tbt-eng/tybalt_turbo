<script lang="ts">
  import DSCopyButton from "./DSCopyButton.svelte";

  // Preformatted text, such as JSON, in a box that scrolls when the text is
  // long. The text is shown as text, never as HTML, so it is safe for data
  // that came from outside Turbo. A copy button sits at the top right.
  let {
    text,
    label = "Copy",
  }: {
    text: string;
    // The accessible name of the copy button.
    label?: string;
  } = $props();
</script>

<div class="relative">
  <pre
    class="max-h-96 overflow-auto rounded-sm border border-neutral-300 bg-neutral-50 p-2 pr-8 text-xs break-words whitespace-pre-wrap">{text}</pre>
  <div class="absolute top-1 right-1">
    <DSCopyButton {text} {label} />
  </div>
</div>
