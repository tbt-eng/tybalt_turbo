<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_roles.md","sha256":"eb9ab5eef5d2579220bfeb1cefaceca7f15ae344fa337315ad804dfeb2681c8f","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_roles.md

## Signed-in staff can list rate roles

Lines 7-11

Asked as: list rate roles; see all rate roles; who can view rate roles; rate roles not showing

Any signed-in person can list the records of the rate roles collection (rate_roles). A request without a signed-in user is refused.

## Signed-in staff can view one rate role

Lines 13-17

Asked as: view one rate role; open a rate role; show rate role details; single rate role record

Any signed-in person can view a single rate role record (rate_roles). A request without a signed-in user is refused.

## Nobody can create rate roles through the API

Lines 19-21

Asked as: add a new rate role; create rate role via API; cannot add rate role; insert rate role record

No rule allows creating a rate role record (rate_roles) through the API. Only a PocketBase superuser can do it. Turbo's own server code can still create records.

## Nobody can update rate roles through the API

Lines 23-25

Asked as: edit a rate role; change rate role rate; update rate role via API; rate role cannot be modified

No rule allows updating a rate role record (rate_roles) through the API. Only a PocketBase superuser can do it. Turbo's own server code can still update records.

## Nobody can delete rate roles through the API

Lines 27-29

Asked as: remove a rate role; delete rate role; get rid of rate role; cannot delete rate role

No rule allows deleting a rate role record (rate_roles) through the API. Only a PocketBase superuser can do it. Turbo's own server code can still delete records.

## Rate role names must be unique and are required

Lines 33-44

Asked as: duplicate rate role name; rate role name taken; rate role needs a name; same name as another rate role

Every rate role (rate_roles) needs a name (name). Two rate roles cannot have the same name, because a unique index covers that field.
