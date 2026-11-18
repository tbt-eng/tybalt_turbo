<script lang="ts">
  import type { Snippet } from "svelte";

  // A record named by its identifier, such as a code, with an optional
  // description in muted text after it: "2026-11-00004 (rev. 2) Client Rates".
  // The identifier is a link when href is set. It is inline, so it fits in a
  // sentence, a list headline or a DSField.
  let {
    label,
    description = "",
    href,
    trailing,
  }: {
    label: string;
    description?: string | null;
    href?: string;
    // A small control right after the identifier, such as a DSCopyButton. It
    // stays on the identifier's line.
    trailing?: Snippet;
  } = $props();

  const muted = $derived((description ?? "").trim());
</script>

{#snippet identifier()}
  {#if href}
    <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- callers pass any href -->
    <a class="link" {href}>{label}</a>
  {:else}
    {label}
  {/if}
{/snippet}

<span class="break-words">
  {#if trailing}
    <span class="whitespace-nowrap"
      >{@render identifier()}<span class="ml-1">{@render trailing()}</span></span
    >
  {:else}
    {@render identifier()}
  {/if}
  {#if muted}
    <span class="text-neutral-600">{muted}</span>
  {/if}
</span>
