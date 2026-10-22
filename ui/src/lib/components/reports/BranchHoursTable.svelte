<script lang="ts">
  import { resolve } from "$app/paths";
  import { SvelteSet } from "svelte/reactivity";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import {
    formatHours,
    type BranchHoursDirection,
    type BranchHoursProject,
  } from "$lib/reports/branchHours";
  let {
    projects,
    direction,
    title,
  }: {
    projects: BranchHoursProject[];
    direction: BranchHoursDirection;
    title: string;
  } = $props();
  // Supplied time goes to other branches; received time comes from them.
  const branchHeading = $derived(
    direction === "supplied" ? "Worked for" : direction === "received" ? "Staff from" : "",
  );
  const total = $derived(projects.reduce((sum, project) => sum + project.hours, 0));
  const expanded = new SvelteSet<string>();
  function toggle(key: string) {
    if (expanded.has(key)) expanded.delete(key);
    else expanded.add(key);
  }
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex (Keyboard users need to scroll the table.) -->
<div class="overflow-x-auto" role="region" aria-label={title} tabindex="0">
  <table class="w-full text-sm tabular-nums">
    <caption class="sr-only">{title}</caption>
    <thead class="bg-neutral-100">
      <tr class="border-b border-neutral-300 whitespace-nowrap">
        <th scope="col" class="w-8 px-2 py-2"><span class="sr-only">Show people</span></th>
        <th scope="col" class="px-2 py-2 text-left">Project</th>
        <th scope="col" class="px-2 py-2 text-left">Client</th>
        {#if branchHeading}<th scope="col" class="px-2 py-2 text-left">{branchHeading}</th>{/if}
        <th scope="col" class="px-2 py-2 text-right">People</th>
        <th scope="col" class="px-2 py-2 text-right">Hours</th>
      </tr>
    </thead>
    <tbody>
      {#each projects as project (project.key)}
        {@const open = expanded.has(project.key)}
        <tr
          class="h-9 border-b border-neutral-200 hover:bg-neutral-50"
          data-testid={`branch-hours-${project.key}`}
        >
          <td class="px-2 py-0">
            <DSActionButton
              variant="plain"
              icon={open ? "mdi:chevron-down" : "mdi:chevron-right"}
              title={`${open ? "Hide" : "Show"} people on ${project.jobNumber || "no-project time"}`}
              action={() => toggle(project.key)}
            />
          </td>
          <th scope="row" class="px-2 py-0 text-left font-normal">
            {#if project.jobId}
              <a
                class="link font-semibold"
                href={resolve("/jobs/[jid]/details", { jid: project.jobId })}>{project.jobNumber}</a
              >
              <span class="text-neutral-600">{project.jobDescription}</span>
            {:else}
              <span class="text-neutral-600">No project</span>
            {/if}
          </th>
          <td class="px-2 py-0"
            ><span class="block max-w-64 truncate" title={project.client}>{project.client}</span
            ></td
          >
          {#if direction === "supplied"}
            <td class="px-2 py-0 whitespace-nowrap">{project.workBranch}</td>
          {:else if direction === "received"}
            <td class="px-2 py-0">{project.staffBranches.join(", ")}</td>
          {/if}
          <td class="px-2 py-0 text-right">{project.people.length}</td>
          <td class="px-2 py-0 text-right font-semibold whitespace-nowrap"
            >{formatHours(project.hours)}</td
          >
        </tr>
        {#if open}
          {#each project.people as person (person.uid)}
            <tr class="h-8 border-b border-neutral-100 bg-neutral-50 text-neutral-700">
              <td></td>
              <td class="py-0 pr-2 pl-6" colspan="2">{person.givenName} {person.surname}</td>
              {#if direction === "received"}
                <td class="px-2 py-0">{person.staffBranch}</td>
              {:else if direction === "supplied"}
                <td></td>
              {/if}
              <td></td>
              <td class="px-2 py-0 text-right whitespace-nowrap">{formatHours(person.hours)}</td>
            </tr>
          {/each}
        {/if}
      {/each}
    </tbody>
    <tfoot>
      <tr class="h-9 border-t border-neutral-300 font-semibold">
        <td></td>
        <th scope="row" class="px-2 py-0 text-left" colspan={branchHeading ? 4 : 3}>Total</th>
        <td class="px-2 py-0 text-right whitespace-nowrap">{formatHours(total)}</td>
      </tr>
    </tfoot>
  </table>
</div>
