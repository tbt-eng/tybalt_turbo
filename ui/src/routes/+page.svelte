<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import DSLargeInput from "$lib/components/DSLargeInput.svelte";
  import { randomHelpQuestion } from "$lib/helpQuestions";
  import { helpChatEnabled } from "$lib/stores/appConfig";
  import { helpChat } from "$lib/stores/helpChat";

  const suggestion = randomHelpQuestion();
  let question = $state("");

  // The model is ready sooner when it has seen the shared prompt. Most visits
  // to this page ask nothing, so only a person who goes to the box warms it.
  let warmed = false;
  function warm() {
    if (warmed) return;
    warmed = true;
    helpChat.warm();
  }

  // Help shows the conversation. The question stays out of the URL.
  function ask(event: KeyboardEvent) {
    const text = question.trim();
    if (event.key !== "Enter" || event.isComposing || !text) return;
    helpChat.ask(text, "/");
    goto(resolve("/docs#ask"));
  }
</script>

<!-- Full-screen flex container to center content -->
<div class="mt-16 flex w-full flex-col items-center justify-center gap-6 bg-white px-4 md:px-0">
  <!-- Logo -->
  <img src="/logo.svg" alt="Tybalt 𝕋𝕌ℝ𝔹𝕆 Logo" class="h-40 w-40 select-none" />

  <!-- Welcome message card -->
  <div class="w-full max-w-xl rounded-lg bg-neutral-100 p-4 text-center shadow-sm">
    <h1 class="mb-2 text-2xl font-semibold text-neutral-800">Welcome to 𝕋𝕌ℝ𝔹𝕆</h1>
    {#if $helpChatEnabled}
      <DSLargeInput
        bind:value={question}
        label="Ask Help a question"
        {suggestion}
        enterkeyhint="send"
        spellcheck
        onkeydown={ask}
        onfocus={warm}
      />
    {:else}
      <p class="text-neutral-600">
        Use the navigation menu to get started managing time sheets, purchase orders, expenses and
        more.
      </p>
    {/if}
    <p class="mt-3"><a class="link" href={resolve("/docs/whats-new")}>What's New</a></p>
  </div>
</div>
