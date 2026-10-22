<script lang="ts">
  import { page } from "$app/state";
  import { replaceState } from "$app/navigation";
  import { globalStore } from "$lib/stores/global";
  import DSMessageBar from "./DSMessageBar.svelte";

  // The app's message area, rendered once by the root layout above every page.
  // It shows errors from the global store, then the flash message the previous
  // page left in page state (see $lib/flash), each with a way to dismiss it.
  //
  // This decides what to show; DSMessageBar decides how it looks. The two are
  // separate because DS components own all colour and styling and hold no app
  // state, while this one reads the store and page state and has no styling
  // of its own; see ui/CLAUDE.md.
  const flash = $derived(page.state.flash);

  function dismissFlash() {
    replaceState("", { ...page.state, flash: undefined });
  }
</script>

{#each $globalStore.errorMessages as error (error.id)}
  <DSMessageBar tone="danger" onDismiss={() => globalStore.dismissError(error.id)}>
    {error.message}
  </DSMessageBar>
{/each}

{#if flash}
  <DSMessageBar tone={flash.tone} onDismiss={dismissFlash}>
    {flash.message}
    {#if flash.link}
      <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- the sender resolves it -->
      · <a class="link" href={flash.link.href}>{flash.link.label}</a>
    {/if}
  </DSMessageBar>
{/if}
