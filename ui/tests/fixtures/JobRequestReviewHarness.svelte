<script lang="ts">
  import Page from "../../src/routes/job-requests/[id]/+page.svelte";
  type RequestData = import("svelte").ComponentProps<typeof Page>["data"];
  const initial = (window as unknown as { requestData: () => RequestData }).requestData();
  let data = $state(initial);
  // The page keeps its state across a reload, as invalidateAll does in the app.
  (window as unknown as { setRequestData: (next: RequestData) => void }).setRequestData = (next) =>
    (data = next);
</script>

<Page {data} />
