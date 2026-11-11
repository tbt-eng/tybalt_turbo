<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/helpChatStatus.ts","sha256":"28aade9e8e91f7fb51f079095bf77a616866c9febd06b4c9c5f6da685564c18a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/helpChatStatus.ts

## Success notice says who has Help chat

Lines 72-146

Asked as: help chat is on; success notice help chat; who has help chat; help chat enabled for everyone; test the key

When Help chat is enabled and has no problems, the page shows a success notice "Help chat is on". With a deploy group it names the group and its member count ("1 member" or "N members"). With no deploy group it says "every user". The notice suggests Test the key to check that the provider accepts the key.

## Invalid help_chat setting turns Help chat off

Lines 79-87

Asked as: help_chat setting not valid; invalid help chat setting; help chat setting error; help chat off for everyone; app config danger notice

When the help_chat setting is not valid (settings_error), the App Config page shows a danger notice titled "The help_chat setting is not valid". The notice says "Help chat is off for everyone" and gives the error. Turbo shows only this notice and skips the other checks.

## Help chat is off without a help_chat row

Lines 88-96

Asked as: no help_chat row; help chat missing from app config; help chat off no row; app_config row missing; help chat not in config

When app_config has no help_chat row (has_row false), the App Config page shows an info notice "Help chat is off" with the text "app_config has no help_chat row." Turbo shows no other notice.

## Help chat is off when enabled is false

Lines 97-105

Asked as: help chat enabled false; help chat turned off; help chat disabled setting; enabled set to false; help chat switched off

When the help_chat row has enabled set to false, the App Config page shows an info notice "Help chat is off" with the text "The help_chat row has enabled set to false." Turbo shows no other notice.

## Several problems show as separate notices

Lines 107-137

Asked as: several help chat warnings; multiple help chat notices; help chat problems listed; why no success notice; help chat several errors

When Help chat is enabled, the page lists one notice for each problem found: missing key, deploy group missing or empty, and failed search index. The success notice does not appear while any problem exists.

## Missing server key for the provider in use

Lines 108-115

Asked as: no key for the provider; help chat not set up; missing api key help chat; server key secret missing; help chat answers not set up

When Help chat is enabled and the key in use has no secret set on the server, the page shows a danger notice "The server has no key for the provider". The notice says to set the named secret and restart the app. Until then, a question gets "Help chat is not set up."

## Deploy group that does not exist

Lines 116-122

Asked as: deploy group does not exist; deploy group not found; help chat wrong deploy group; deploy group id missing; help chat group not found

When a deploy group is configured but no deploy group has that id, the page shows a warning "The deploy group does not exist". Help chat is off for everyone.

## Deploy group with no members

Lines 123-129

Asked as: deploy group has no members; empty deploy group; help chat group no users; deploy group zero members; help chat nobody can use

When the deploy group exists but has 0 members, the page shows a warning "The deploy group has no members". Help chat is off for everyone until the group has a member.

## Search index failed to load

Lines 130-136

Asked as: search index did not load; search index failed; help chat index error; corpus failed to load; search index broken

When the corpus state is "failed", the page shows a danger notice "The search index did not load" with the corpus error as its text.

## Cache read share of input tokens

Lines 151-155

Asked as: cache read percent; cache hit share; cache read share tokens; why no cache percent; cache read rate zero

The cache read share is cache-read tokens divided by the sum of input, cache-read and cache-write tokens, rounded to a whole percent. When that sum is 0 or less, the share is null and no percent shows.

## How app_config values display as text

Lines 158-162

Asked as: empty string shows quotes; app config value display; undefined value empty text; app config value as json; how config value shows

A string value shows as is, but an empty string shows as "" (two quote marks). An undefined value shows as empty text. Any other value shows as JSON.
