<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { Snippet } from "svelte";
  import DSActionButton from "./DSActionButton.svelte";

  // Draws one app-level message: a full-width strip above the page, so it
  // reads as the app speaking rather than as part of the record below. For
  // notices that belong to the page itself, use DSCallout.
  //
  // This is the look only. MessageBar, in the layout, decides which messages
  // to show (errors from the global store, the flash in page state) and
  // renders one of these for each. The two are separate because DS components
  // own all colour and styling and hold no app state; see ui/CLAUDE.md.
  let {
    tone,
    onDismiss,
    children,
  }: {
    tone: "info" | "success" | "danger";
    onDismiss?: () => void;
    children: Snippet;
  } = $props();

  const tones = {
    info: { bar: "border-blue-200 bg-blue-50 text-blue-900", icon: "mdi:information-outline" },
    success: {
      bar: "border-green-200 bg-green-50 text-green-900",
      icon: "mdi:check-circle-outline",
    },
    danger: { bar: "border-red-200 bg-red-50 text-red-900", icon: "mdi:alert-circle-outline" },
  } as const;
</script>

<div class="flex w-full items-center gap-3 border-b px-4 py-2 {tones[tone].bar}" role="status">
  <Icon icon={tones[tone].icon} width="20px" class="shrink-0" />
  <div class="min-w-0 flex-1">{@render children()}</div>
  {#if onDismiss}
    <DSActionButton
      action={onDismiss}
      icon="mdi:close"
      title="Dismiss"
      color="neutral"
      transparentBackground
    />
  {/if}
</div>
