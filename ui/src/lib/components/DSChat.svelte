<script module lang="ts">
  export interface ChatSourceLink {
    number: number;
    // What the caller gets when the user opens the source.
    key: string;
    label: string;
    href: string;
    external: boolean;
  }

  export interface ChatTurnView {
    key: number;
    question: string;
    // The answer as HTML that the caller has made safe.
    html: string;
    running: boolean;
    // What is happening, while the turn runs and has no text yet.
    activity?: string;
    // Why there is no answer, or why it stopped.
    message?: string;
    tone?: "info" | "warning" | "danger";
    sources: ChatSourceLink[];
    rating: "" | "up" | "down";
    noteSent: boolean;
    // False hides the feedback row, for a turn that has no answer to rate.
    canRate: boolean;
  }
</script>

<script lang="ts">
  import { onMount, tick } from "svelte";
  import DSActionButton from "./DSActionButton.svelte";
  import DSCallout from "./DSCallout.svelte";
  import DSDoc from "./DSDoc.svelte";
  import DSLargeInput from "./DSLargeInput.svelte";
  import DSTextInput from "./DSTextInput.svelte";
  import LoadingAnimation from "./LoadingAnimation.svelte";

  // A conversation with an assistant, as a column that fills its space: each
  // question with its answer as it arrives, the sources of the answer and a way
  // to rate it, and the box for the next question at the bottom. The caller
  // owns the data and the calls to the server.
  let {
    turns,
    notice,
    askLabel = "Ask a follow-up",
    closeLabel = "Close",
    onask,
    onrate,
    onnote,
    onopen,
    onnavigate,
    oncopy,
    onstop,
    onclose,
  }: {
    turns: ChatTurnView[];
    // A line above the conversation, such as what the app records.
    notice?: string;
    // The placeholder of the box at the bottom.
    askLabel?: string;
    // The label of the button that ends the conversation.
    closeLabel?: string;
    // The user asked the next question. Without this, there is no box.
    onask?: (question: string) => void;
    onrate?: (key: number, rating: "up" | "down") => void;
    onnote?: (key: number, note: string) => void;
    // The user opened a source, from the list or from a link in the answer.
    onopen?: (key: number, source: string) => void;
    // Goes to a page of the app that a source or the answer links to. With
    // this set, the link does not navigate by itself.
    onnavigate?: (href: string) => void;
    // The user asked to copy the answer. The caller copies it.
    oncopy?: (key: number) => void;
    onstop?: () => void;
    onclose?: () => void;
  } = $props();

  const running = $derived(turns.some((turn) => turn.running));
  let notes = $state<Record<number, string>>({});
  let draft = $state("");
  let input: HTMLInputElement | undefined = $state();
  let scroller: HTMLDivElement | undefined = $state();
  let viewHeight = $state(0);

  // The user opened a source. A link to a page of the app goes through the
  // caller, which shows that page in place of the conversation. A click with a
  // key held, which opens a new tab, is left to the browser.
  function follow(event: MouseEvent, key: number, source: string, href: string, external: boolean) {
    onopen?.(key, source);
    const plainClick =
      event.button === 0 && !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey;
    if (external || !onnavigate || !plainClick) return;
    event.preventDefault();
    onnavigate(href);
  }

  // A link inside the answer goes to a page of the app.
  function answerClick(event: MouseEvent, key: number) {
    const href = (event.target as HTMLElement).closest("a")?.getAttribute("href");
    if (href) follow(event, key, href, href, false);
  }

  function submit(event: SubmitEvent) {
    event.preventDefault();
    const question = draft.trim();
    if (!question || running) return;
    draft = "";
    onask?.(question);
  }

  // The person can type the next question at once. A phone would open its
  // keyboard over the answer, so only a device with a mouse takes the focus.
  onMount(() => {
    if (window.matchMedia("(pointer: fine)").matches) input?.focus();
  });

  // A new question moves to the top of the view, and its answer arrives under
  // it. The last turn is as tall as the view, so there is room to scroll to it.
  let shownKey = -1;
  $effect(() => {
    const key = turns.at(-1)?.key;
    if (key === undefined || key === shownKey) return;
    shownKey = key;
    tick().then(() =>
      scroller?.querySelector(`[data-turn="${key}"]`)?.scrollIntoView({ block: "start" }),
    );
  });
</script>

<div class="flex min-h-0 min-w-0 flex-1 flex-col">
  <div
    class="flex flex-wrap items-start justify-between gap-2 border-b border-neutral-200 bg-white p-3"
    data-print-hide
  >
    <p class="min-w-0 flex-1 text-sm text-neutral-600">{notice}</p>
    {#if onclose}
      <DSActionButton variant="plain" icon="mdi:close" action={onclose}>{closeLabel}</DSActionButton
      >
    {/if}
  </div>

  <div bind:this={scroller} bind:clientHeight={viewHeight} class="min-h-0 flex-1 overflow-y-auto">
    <div class="space-y-6 p-4">
      {#each turns as turn, index (turn.key)}
        <section
          class="scroll-mt-4 space-y-2"
          data-turn={turn.key}
          aria-busy={turn.running}
          style:min-height={index === turns.length - 1
            ? `${Math.max(viewHeight - 32, 0)}px`
            : undefined}
        >
          <h2 class="text-lg font-semibold break-words">{turn.question}</h2>

          {#if turn.html}
            <!-- The links inside the answer take the keys. This only sees their clicks. -->
            <!-- svelte-ignore a11y_click_events_have_key_events, a11y_no_static_element_interactions -->
            <div onclick={(event) => answerClick(event, turn.key)}>
              <DSDoc blocks={[{ kind: "html", html: turn.html }]} />
            </div>
          {/if}

          {#if turn.running && !turn.html}
            <p class="flex items-center gap-2 text-sm text-neutral-600" role="status">
              <span class="flex h-5 w-5 items-center"><LoadingAnimation /></span>
              {turn.activity || "Thinking…"}
            </p>
          {/if}

          {#if turn.message}
            <DSCallout tone={turn.tone ?? "warning"}>{turn.message}</DSCallout>
          {/if}

          {#if turn.sources.length > 0}
            <div class="space-y-1 text-sm">
              <h3 class="font-semibold">Sources</h3>
              <ol class="space-y-1">
                {#each turn.sources as source (source.key)}
                  <li class="flex gap-2">
                    <span class="text-neutral-600 tabular-nums">{source.number}.</span>
                    <!-- eslint-disable svelte/no-navigation-without-resolve -- callers pass any href -->
                    <a
                      class="link min-w-0 break-words"
                      href={source.href}
                      target={source.external ? "_blank" : undefined}
                      rel={source.external ? "noopener" : undefined}
                      onclick={(event) =>
                        follow(event, turn.key, source.key, source.href, source.external)}
                      >{source.label}</a
                    >
                    <!-- eslint-enable svelte/no-navigation-without-resolve -->
                  </li>
                {/each}
              </ol>
            </div>
          {/if}

          {#if turn.canRate && !turn.running}
            <div class="flex flex-wrap items-center gap-2" data-print-hide>
              <DSActionButton
                icon={turn.rating === "up" ? "mdi:thumb-up" : "mdi:thumb-up-outline"}
                title={turn.rating === "up"
                  ? "You marked this answer as helpful"
                  : "This answer helped"}
                color={turn.rating === "up" ? "green" : "neutral"}
                transparentBackground={turn.rating !== "up"}
                action={() => onrate?.(turn.key, "up")}
              />
              <DSActionButton
                icon={turn.rating === "down" ? "mdi:thumb-down" : "mdi:thumb-down-outline"}
                title={turn.rating === "down"
                  ? "You marked this answer as not helpful"
                  : "This answer did not help"}
                color={turn.rating === "down" ? "red" : "neutral"}
                transparentBackground={turn.rating !== "down"}
                action={() => onrate?.(turn.key, "down")}
              />
              {#if oncopy}
                <DSActionButton
                  icon="mdi:content-copy"
                  title="Copy this answer"
                  color="neutral"
                  transparentBackground
                  action={() => oncopy(turn.key)}
                />
              {/if}
            </div>
            {#if turn.rating === "down"}
              {#if turn.noteSent}
                <p class="text-sm text-neutral-600" role="status">
                  Thank you. Turbo saved your note.
                </p>
              {:else}
                <form
                  class="flex flex-wrap items-end gap-2"
                  onsubmit={(event) => {
                    event.preventDefault();
                    onnote?.(turn.key, notes[turn.key] ?? "");
                  }}
                >
                  <div class="min-w-0 flex-1">
                    <DSTextInput
                      bind:value={notes[turn.key]}
                      errors={{}}
                      fieldName={`chat-note-${turn.key}`}
                      uiName="What was wrong?"
                      placeholder="Optional"
                    />
                  </div>
                  <DSActionButton variant="secondary" type="submit">Send</DSActionButton>
                </form>
              {/if}
            {/if}
          {/if}
        </section>
      {/each}
    </div>
  </div>

  {#if onask}
    <form class="border-t border-neutral-200 bg-white p-3" data-print-hide onsubmit={submit}>
      <DSLargeInput
        bind:value={draft}
        bind:element={input}
        label={askLabel}
        icon="mdi:message-question-outline"
        enterkeyhint="send"
        spellcheck
      >
        {#snippet actions()}
          {#if running && onstop}
            <DSActionButton variant="plain" icon="mdi:stop-circle-outline" action={onstop}
              >Stop</DSActionButton
            >
          {:else}
            <DSActionButton variant="plain" type="submit" disabled={!draft.trim() || running}
              >Ask</DSActionButton
            >
          {/if}
        {/snippet}
      </DSLargeInput>
    </form>
  {/if}
</div>
