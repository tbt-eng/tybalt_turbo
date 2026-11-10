<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expenditure_kinds.md","sha256":"4230ced71896edb07c43db96a7b6f31ba230042eafcf45caf9b254c07f15ce15","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expenditure_kinds.md

## Who can list or view expenditure kinds

Lines 7-17

Asked as: who can see expenditure kinds; list expenditure kinds api; view expenditure kind without signing in; anonymous access to expenditure kinds

Any signed-in person can list expenditure kinds and view a single one through the API (list rule and view rule: `@request.auth.id != ""`). A request without a sign-in is refused.

## Nobody can create, update or delete expenditure kinds through the API

Lines 19-29

Asked as: cannot edit expenditure kind; add new expenditure kind api; delete expenditure kind via api; why can't I change expenditure kinds

Creating, updating and deleting expenditure kinds have no rule, so nobody can do them through the API except a PocketBase superuser. Turbo's own server code can still change these records.

## Expenditure kind field limits

Lines 35-39

Asked as: expenditure kind label too short; required ui label error; ui order must be number; expenditure kind id length

The UI label (`en_ui_label`) is required and must have at least 3 characters. The display order (`ui_order`) is required and must be a whole number. The id must be exactly 15 characters, lowercase letters and digits only.

## Expenditure kind name must be unique

Lines 36-49

Asked as: duplicate expenditure kind name; name already taken expenditure kind; same expenditure kind twice; expenditure kind name conflict

Two expenditure kinds cannot have the same name (`name`), because a unique index covers that field. The name is also required.

## Optional job and second approval threshold fields

Lines 40-43

Asked as: allow job flag expenditure kind; second approval limit optional; expenditure kind approval threshold; is job allowed on expenditure kind

An expenditure kind has an optional flag `allow_job` and an optional number `second_approval_threshold`. Neither is required.
