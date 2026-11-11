<script lang="ts">
  import { pb } from "$lib/pocketbase";
  import { globalStore } from "$lib/stores/global";
  import { appConfig } from "$lib/stores/appConfig";
  import { formatDateTime } from "$lib/utilities";
  import {
    cacheReadShare,
    configValueText,
    helpChatNotices,
    type HelpChatKeyTest,
    type HelpChatStatus,
  } from "$lib/helpChatStatus";
  import DSActionButton from "$lib/components/DSActionButton.svelte";
  import DSCallout from "$lib/components/DSCallout.svelte";
  import DSCard from "$lib/components/DSCard.svelte";
  import DSCardGrid from "$lib/components/DSCardGrid.svelte";
  import DSField from "$lib/components/DSField.svelte";
  import DSFieldList from "$lib/components/DSFieldList.svelte";
  import DSPageHeader from "$lib/components/DSPageHeader.svelte";
  import DSSection from "$lib/components/DSSection.svelte";
  import DsLabel from "$lib/components/DsLabel.svelte";

  // This page shows settings and changes none. A PocketBase superuser changes
  // an app_config row in the PocketBase admin. The status of Help chat comes
  // from /api/help/status, which only an admin can call.
  const isAdmin = $derived($globalStore.claims.includes("admin"));

  let status = $state<HelpChatStatus | null>(null);
  let statusError = $state("");
  let loading = $state(false);
  let keyTest = $state<HelpChatKeyTest | null>(null);
  let testing = $state(false);

  function errorText(error: unknown, fallback: string): string {
    const message = (error as { response?: { message?: unknown } })?.response?.message;
    return typeof message === "string" && message ? message : fallback;
  }

  async function loadStatus() {
    loading = true;
    statusError = "";
    // A key test result belongs to the status that it was made with.
    keyTest = null;
    try {
      status = (await pb.send("/api/help/status", {
        method: "GET",
        requestKey: "help_status",
      })) as HelpChatStatus;
    } catch (error) {
      statusError = errorText(error, "Turbo can't read the status of Help chat.");
    } finally {
      loading = false;
    }
  }

  // One request to the provider. It uses no tokens. It runs only on a press.
  async function testTheKey() {
    testing = true;
    keyTest = null;
    try {
      keyTest = (await pb.send("/api/help/status/test_key", {
        method: "POST",
        requestKey: "help_status_test_key",
      })) as HelpChatKeyTest;
    } catch (error) {
      keyTest = {
        ok: false,
        provider: "",
        model: "",
        message: errorText(error, "Turbo can't test the key."),
      };
    } finally {
      testing = false;
    }
  }

  let asked = false;
  $effect(() => {
    if (isAdmin && !asked) {
      asked = true;
      loadStatus();
    }
  });

  const notices = $derived(status ? helpChatNotices(status) : []);
  const cacheShare = $derived(status ? cacheReadShare(status.last_day) : null);
  const rows = $derived([...$appConfig.items].sort((a, b) => a.key.localeCompare(b.key)));

  function properties(value: unknown): [string, unknown][] {
    if (value === null || typeof value !== "object" || Array.isArray(value)) return [];
    return Object.entries(value as Record<string, unknown>);
  }

  const corpusText = { not_loaded: "Not loaded yet", loaded: "Loaded", failed: "Failed" };
</script>

<div class="mx-auto space-y-8 p-4">
  <DSPageHeader title="App Config" subtitle="Turbo's settings and the status of Help chat">
    {#snippet context()}
      This page changes nothing. A PocketBase superuser changes a setting in the PocketBase admin.
    {/snippet}
  </DSPageHeader>

  {#if !isAdmin}
    <DSCallout tone="danger" title="Admins only">
      This page is visible to people with the admin claim.
    </DSCallout>
  {:else}
    <DSSection title="Help chat">
      {#snippet actions()}
        <DSActionButton variant="secondary" action={loadStatus} {loading}>Refresh</DSActionButton>
        <DSActionButton
          variant="secondary"
          action={testTheKey}
          loading={testing}
          disabled={!status}
        >
          Test the key
        </DSActionButton>
      {/snippet}

      {#if statusError}
        <DSCallout tone="danger" title="No status">{statusError}</DSCallout>
      {/if}

      {#if status}
        <div class="space-y-4">
          {#each notices as notice (notice.title)}
            <DSCallout tone={notice.tone} title={notice.title}>{notice.text}</DSCallout>
          {/each}

          {#if keyTest}
            <DSCallout
              tone={keyTest.ok ? "success" : "danger"}
              title={keyTest.ok ? "The key test passed" : "The key test failed"}
            >
              {keyTest.message}
            </DSCallout>
          {/if}

          <DSCardGrid>
            <DSCard title="Setting in use">
              <DSFieldList>
                <DSField label="Switch">
                  <DsLabel tone={status.settings.enabled ? "success" : "neutral"}>
                    {status.settings.enabled ? "On" : "Off"}
                  </DsLabel>
                </DSField>
                <DSField label="Provider" value={status.settings.provider} />
                <DSField label="Model" value={status.settings.model} />
                <DSField label="Effort" value={status.settings.effort} showEmpty />
                <DSField
                  label="Questions per hour, each user"
                  value={status.settings.questions_per_hour}
                />
              </DSFieldList>
            </DSCard>

            <DSCard title="Provider keys">
              <DSFieldList>
                {#each status.keys as key (key.provider)}
                  <DSField label={key.secret} hint={key.in_use ? "in use" : ""}>
                    <DsLabel tone={key.set ? "success" : key.in_use ? "danger" : "neutral"}>
                      {key.set ? "Set" : "Not set"}
                    </DsLabel>
                  </DSField>
                {/each}
              </DSFieldList>
            </DSCard>

            <DSCard title="Who gets Help chat">
              <DSFieldList>
                {#if status.deploy_group}
                  <DSField
                    label="Deploy group"
                    value={status.deploy_group.found ? status.deploy_group.name : "No such group"}
                  />
                  <DSField label="Members" value={status.deploy_group.members} />
                  <DSField label="Group id" value={status.deploy_group.id} />
                {:else}
                  <DSField label="Deploy group" value="None: every user" />
                {/if}
              </DSFieldList>
            </DSCard>

            <DSCard title="Search index">
              <DSFieldList>
                <DSField label="State" value={corpusText[status.corpus.state]} />
                {#if status.corpus.state === "not_loaded"}
                  <DSField label="Note" value="The first question after a start loads it." />
                {/if}
                {#if status.corpus.state === "loaded"}
                  <DSField label="Chunks" value={status.corpus.chunks.toLocaleString()} />
                  <DSField label="Corpus" value={status.corpus.hash.slice(0, 12)} />
                {/if}
                <DSField label="Error" value={status.corpus.error} />
              </DSFieldList>
            </DSCard>

            <DSCard title="Last 24 hours">
              <DSFieldList>
                <DSField label="Questions" value={status.last_day.questions.toLocaleString()} />
                <DSField label="Users who asked" value={status.last_day.users.toLocaleString()} />
                <DSField label="Turns that ended in an error">
                  <DsLabel tone={status.last_day.errors > 0 ? "danger" : "neutral"}>
                    {status.last_day.errors.toLocaleString()}
                  </DsLabel>
                </DSField>
                <DSField
                  label="Input tokens read from the cache"
                  value={cacheShare === null ? null : `${cacheShare}%`}
                  hint={`${status.last_day.tokens_cache_read.toLocaleString()} read, ${status.last_day.tokens_cache_write.toLocaleString()} written, ${status.last_day.tokens_in.toLocaleString()} not cached`}
                  showEmpty
                />
                <DSField
                  label="Output tokens"
                  value={status.last_day.tokens_out.toLocaleString()}
                />
              </DSFieldList>
            </DSCard>

            <DSCard title="Last results">
              <DSFieldList>
                <DSField
                  label="The model last answered"
                  value={formatDateTime(status.last_answered)}
                  showEmpty
                />
                {#if status.last_error}
                  <DSField label="Last error" value={formatDateTime(status.last_error.at)} />
                  <DSField label="Its message" value={status.last_error.error} />
                {:else}
                  <DSField label="Last error" value="None" />
                {/if}
              </DSFieldList>
            </DSCard>

            {#if status.deploy_groups.length > 0}
              <DSCard title="All deploy groups" span="full">
                <DSFieldList columns={2}>
                  {#each status.deploy_groups as group (group.id)}
                    <DSField
                      label={group.name}
                      value={group.id}
                      hint={group.members === 1 ? "1 member" : `${group.members} members`}
                    />
                  {/each}
                </DSFieldList>
              </DSCard>
            {/if}
          </DSCardGrid>
        </div>
      {/if}
    </DSSection>

    <DSSection title="All settings">
      {#snippet titleNote()}
        the rows of app_config
      {/snippet}
      <!-- The store also sets an error when only its live updates fail. The rows
           are then on the page, so the error shows only with no rows. -->
      {#if $appConfig.error && rows.length === 0}
        <DSCallout tone="danger" title="No settings">{$appConfig.error}</DSCallout>
      {/if}
      <DSCardGrid>
        {#each rows as row (row.id)}
          <DSCard title={row.key}>
            <DSFieldList>
              <DSField label="About" value={row.description} />
              {#each properties(row.value) as [name, value] (name)}
                <DSField label={name} value={configValueText(value)} />
              {:else}
                <DSField label="Value" value={configValueText(row.value)} />
              {/each}
              <DSField label="Changed" value={formatDateTime(row.updated)} />
            </DSFieldList>
          </DSCard>
        {/each}
      </DSCardGrid>
    </DSSection>
  {/if}
</div>
