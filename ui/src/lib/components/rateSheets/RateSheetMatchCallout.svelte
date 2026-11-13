<script lang="ts">
  import { resolve } from "$app/paths";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import type { Snippet } from "svelte";
  import type { RateSheetMatch } from "$lib/rateSheetMatches";

  // A notice that lists active rate sheets matched against a sheet: exact
  // matches (same rates) or near matches ("differs from X in 2 roles").
  // actions holds buttons, such as one that chooses the matching sheet.
  let {
    tone,
    title,
    intro,
    matches,
    actions,
  }: {
    tone: "info" | "warning" | "danger";
    title?: string;
    intro?: string;
    matches: RateSheetMatch[];
    actions?: Snippet;
  } = $props();

  // The text after a near match's link. It is built here because Svelte trims
  // a space at the start of a block.
  function differs(match: RateSheetMatch): string {
    if (match.differing_roles <= 0) return "";
    return ` in ${match.differing_roles} role${match.differing_roles === 1 ? "" : "s"}`;
  }
</script>

<DSCallout {tone} {title} {actions}>
  {#if intro}<p>{intro}</p>{/if}
  <ul class="list-disc pl-5">
    {#each matches as match (match.id)}
      <li>
        {#if match.differing_roles > 0}Differs from{/if}
        <a class="link" href={resolve("/rate-sheets/[id]/details", { id: match.id })}
          >{match.name} (rev. {match.revision})</a
        >{differs(match)}
      </li>
    {/each}
  </ul>
</DSCallout>
