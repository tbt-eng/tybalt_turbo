<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/deployGroups.ts","sha256":"414c3a7aaea7df4cc3ae2978e74888097c98217fe33a71292c8456b03e7138fe","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/deployGroups.ts

## Feature with no deploy group is shown to everyone

Lines 16-18

Asked as: feature visible to everyone; no deploy group set; feature showing for all users; deploy group empty so shown

The UI treats a feature as allowed by its deploy group when the feature's app_config value is missing or is not an object. It does the same when deploy_group is absent, null or an empty string. In these cases only the feature's own switch decides (deployGroupAllows).

## Feature with a deploy group is shown only to group members

Lines 17-21

Asked as: feature hidden from my team; only group members see feature; user not in deploy group; feature missing for some users

When a feature's app_config row names a deploy_group, the UI shows the feature only to a user whose deploy group ids include that group id. A user outside the group does not see it, even if the feature's switch is on. The server does the same check and is the real gate (utilities.InDeployGroup).

## Deploy group that is not a text id lets nobody in

Lines 19-20

Asked as: feature hidden from everyone; deploy group number not working; feature never shows up; wrong deploy group id

If deploy_group is set but is not a string, such as a number, the UI shows the feature to nobody. An id that matches no group also lets nobody in, because it has no members.
