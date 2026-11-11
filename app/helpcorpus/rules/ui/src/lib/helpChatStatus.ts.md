<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/helpChatStatus.ts","sha256":"2983c72c16f95c95a69b74fb1f0f318df787da76ba0634b1495b813d002fdf34","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/helpChatStatus.ts

## Member count wording on App Config: 1 member vs N members

Lines 72-74

Asked as: member count wording; 1 member or members; deploy group member count; members says 0 members

When the App Config page names the members of the deploy group, it writes "1 member" for exactly one member. For any other count, including 0, it writes "N members" (for example "0 members" or "5 members").

## Help chat is off when the help_chat setting is not valid

Lines 79-87

Asked as: help_chat setting not valid; help chat off for everyone error; help_chat settings error; danger notice help chat setting

On the App Config page, if the help_chat setting has an error (settings_error), Turbo shows a danger notice titled "The help_chat setting is not valid". The notice says "Help chat is off for everyone." and gives the error. Turbo shows no other notice in this case.

## Help chat is off when app_config has no help_chat row

Lines 88-96

Asked as: no help_chat row; help chat is off info; app_config missing help_chat row; help chat off missing row

If app_config has no help_chat row (has_row is false), the App Config page shows an info notice "Help chat is off" with the text "app_config has no help_chat row." It shows no other notice.

## Help chat is off when enabled is false

Lines 97-105

Asked as: help chat enabled false; help chat switched off; help chat disabled row; help_chat enabled set to false

If the help_chat row has enabled set to false, the App Config page shows an info notice "Help chat is off" with the text "The help_chat row has enabled set to false." Turbo does not check the key, deploy group or search index in this case.

## Notice when the server has no key for the provider

Lines 108-117

Asked as: no key for provider; help chat not set up; fly secrets set key; api key missing help chat

When Help chat is enabled and the key for the provider in use is not set, the page shows a danger notice "The server has no key for the provider". It tells the admin to set the secret with "fly secrets set", which restarts the app, and not with the Fly dashboard, because a secret set there did not reach the app. Until then, a question gets "Help chat is not set up."

## Notice when the deploy group does not exist

Lines 118-124

Asked as: deploy group does not exist; deploy group id not found; help chat deploy group missing; wrong deploy group id warning

If a deploy group is set but no deploy group has that id (found is false), the page shows a warning "The deploy group does not exist". It says Help chat is off for everyone.

## Notice when the deploy group has no members

Lines 125-131

Asked as: deploy group has no members; empty deploy group help chat; help chat group zero members; who gets help chat nobody

If the deploy group exists but has 0 members, the page shows a warning "The deploy group has no members". It says Help chat is off for everyone until the group has a member.

## Notice when the search index failed to load

Lines 132-138

Asked as: search index did not load; help chat corpus failed; search index failed to load; help chat search index error

If the corpus state is "failed", the page shows a danger notice "The search index did not load" with the corpus error as its text. The state "not_loaded" gives no notice.

## Problem notices all show together; success only when none

Lines 139-139

Asked as: help chat problem notices; help chat warnings together; success notice only when no problems; help chat all problems shown

When Help chat is enabled, the page shows every problem notice together (key, deploy group, search index). It shows the success notice only when there are no problems.

## Success notice says who Help chat is on for

Lines 141-148

Asked as: help chat is on success; help chat on for everyone; test the key button help chat; help chat enabled for whom

With no problems, the page shows a success notice "Help chat is on". If a deploy group is set, the notice names the group and its member count ("1 member" or "N members"). With no deploy group it says "every user". The notice also suggests Test the key to check that the provider accepts the key.

## Cache read share of input tokens

Lines 153-157

Asked as: cache read percent; cache hit share of input tokens; cache read share blank; token cache percentage null

Turbo computes the cache read share as cache-read tokens divided by the sum of input, cache-read and cache-write tokens, as a whole percent (rounded). It shows nothing (null) when that sum is 0 or less.

## How App Config shows a property value as text

Lines 160-164

Asked as: app config property value display; empty string shows as quotes; app config blank property text; how app config shows values

A string property shows as is, and an empty string shows as "" (two quote marks). An undefined value shows as empty text. Any other value shows as JSON.
