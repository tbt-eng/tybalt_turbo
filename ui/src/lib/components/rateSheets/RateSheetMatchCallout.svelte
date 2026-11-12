<script lang="ts">
  import { resolve } from "$app/paths";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import type { RateSheetMatch } from "$lib/rateSheetMatches";

  // A notice that lists active rate sheets matched against a sheet: exact
  // matches (same rates) or near matches ("differs from X in 2 roles").
  let {
    tone,
    title,
    intro,
    matches,
  }: {
    tone: "info" | "warning" | "danger";
    title?: string;
    intro?: string;
    matches: RateSheetMatch[];
  } = $props();
</script>

<DSCallout {tone} {title}>
  {#if intro}<p>{intro}</p>{/if}
  <ul class="list-disc pl-5">
    {#each matches as match (match.id)}
      <li>
        {#if match.differing_roles > 0}Differs from{/if}
        <a class="link" href={resolve("/rate-sheets/[id]/details", { id: match.id })}
          >{match.name} (rev. {match.revision})</a
        >{#if match.differing_roles > 0}
          in {match.differing_roles} role{match.differing_roles === 1 ? "" : "s"}{/if}
      </li>
    {/each}
  </ul>
</DSCallout>
