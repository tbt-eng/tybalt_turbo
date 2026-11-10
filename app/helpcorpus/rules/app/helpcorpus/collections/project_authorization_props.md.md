<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/project_authorization_props.md","sha256":"441cbac474a8fa2f76fbbb50bf60ef3364196f9ce6b034329af0bd7117683811","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/project_authorization_props.md

## Who can list or view project authorization props

Lines 7-17

Asked as: who can view project authorization props; list authorization props access; non-admin cannot see auth props; project auth props permissions

Only an active admin (an admin profile with active = true and the admin claim) can list project_authorization_props records. The same people can view one record. Other signed-in people and anonymous requests cannot.

## Who can create or update project authorization props

Lines 19-29

Asked as: who can edit authorization props; create project auth props; update auth props refused; admin can change auth props

Only an active admin (an admin profile with active = true and the admin claim) can create a project_authorization_props record. Only the same people can update one. Everyone else is refused.

## Nobody can delete project authorization props through the API

Lines 31-33

Asked as: delete project authorization props; remove auth props record; cannot delete auth props; delete props via API

No delete rule exists for project_authorization_props, so nobody can delete a record through the API, not even an admin. Only a PocketBase superuser can, and Turbo's own server code is not limited by this rule. Deleting the linked user claim still deletes the record (cascadeDelete on user_claim).

## Each user claim has only one project authorization prop

Lines 42-49

Asked as: duplicate authorization props error; second auth props for user; one props per user claim; user claim already has props

A unique index (idx_pa_props_claim) allows only one project_authorization_props record for each user claim (user_claim). Each record must link to exactly one user claim, and it is required. A second record for the same claim is refused.

## Maximum value must be a whole number of at least 1

Lines 43-43

Asked as: max value rejected; max_value zero error; max value decimal refused; max value must be whole number

The max_value field is required. It must be a whole number from 1 up to 9007199254740991. A value of 0 or a decimal is refused.
