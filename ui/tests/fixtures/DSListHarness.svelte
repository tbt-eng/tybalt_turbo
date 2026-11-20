<script lang="ts">
  import DSList from "$lib/components/DSList.svelte";
  import rows from "./dsList.json";
  let listState = $state({ searchTerm: "", page: 1, perPage: 10 });
  // A stand-in for a server that pages: it matches names only, so a match on
  // the email would show that DSList filtered the page itself.
  let serverState = $state({ searchTerm: "", page: 1, perPage: 10 });
  const serverMatches = $derived(
    rows.filter((row) => row.name.toLowerCase().includes(serverState.searchTerm.toLowerCase())),
  );
  const serverPage = $derived(
    serverMatches.slice(
      (serverState.page - 1) * serverState.perPage,
      serverState.page * serverState.perPage,
    ),
  );
</script>

<section aria-label="Paged list">
  <DSList
    items={rows}
    search
    pagination
    searchTerm={listState.searchTerm}
    page={listState.page}
    perPage={listState.perPage}
    searchPlaceholder="Filter contacts"
    searchText={(item) => `${item.name} ${item.email}`}
    onStateChange={(next) => (listState = next)}
  >
    {#snippet headline(item)}<span data-row>{item.name}</span>{/snippet}
  </DSList>
</section>
<section aria-label="Default list">
  <DSList items={rows} search>
    {#snippet headline(item)}<span data-row>{item.name}</span>{/snippet}
  </DSList>
</section>
<section aria-label="Grouped list">
  <DSList items={rows} search groupField="group" groupSort="ASC">
    {#snippet headline(item)}<span data-row>{item.name}</span>{/snippet}
  </DSList>
</section>
<section aria-label="Server list">
  <DSList
    items={serverPage}
    search
    pagination
    totalItems={serverMatches.length}
    page={serverState.page}
    perPage={serverState.perPage}
    onStateChange={(next) => (serverState = next)}
  >
    {#snippet headline(item)}<span data-row>{item.name}</span>{/snippet}
  </DSList>
</section>
