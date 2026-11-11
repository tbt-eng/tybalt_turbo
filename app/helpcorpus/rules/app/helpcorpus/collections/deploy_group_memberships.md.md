<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/deploy_group_memberships.md","sha256":"319450b7f9e5edbb3a522129a5b22e93eb2eca9ee61ab252b90bb03c2a311459","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/deploy_group_memberships.md

## Who can list deploy group memberships

Lines 7-11

Asked as: see other people's memberships; list my deploy groups; which groups am I in; who else is in a group; my group memberships list

A signed-in person can list only the deploy group membership records whose user (`user`) is that person. Nobody sees another person's memberships through the API.

## Who can view one deploy group membership

Lines 13-17

Asked as: open someone else's membership; view a colleague's group membership; can't see a membership record; membership record not visible; look up another user's group

A signed-in person can view a deploy group membership record only when its user (`user`) is that person. Records of other people are not visible.

## Nobody can create deploy group memberships through the API

Lines 19-21

Asked as: add someone to a deploy group; API add user to group; create membership via API; can't add member to group; assign user to deploy group

Nobody can create a deploy group membership through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update deploy group memberships through the API

Lines 23-25

Asked as: change someone's group membership; edit membership through API; can't update group membership; move user to another group; modify membership record

Nobody can update a deploy group membership through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete deploy group memberships through the API

Lines 27-29

Asked as: remove someone from a deploy group; delete membership via API; take user out of group; can't delete group membership; kick member out of group

Nobody can delete a deploy group membership through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Membership is removed when its group or user is deleted

Lines 36-37

Asked as: membership vanished after group deleted; group deleted user removed; why did membership disappear; deleting user removes group membership; orphaned membership after deletion

Each membership needs one deploy group (`deploy_group`) and one user (`user`). Deleting the deploy group or the user also deletes the membership (cascadeDelete).

## A user can join a deploy group only once

Lines 43-45

Asked as: user added to group twice; duplicate group membership; same person listed twice in group; can't join group again; already a member error

A unique index on user and deploy group (`idx_deploy_group_memberships_user_group`) means the same user cannot have two membership records in the same deploy group.
