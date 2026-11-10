<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_uncommit_reasons.md","sha256":"7811221ded73bf76f8762208c3f6f1e040fafa96505baaea3dcc1a1f59b72fb0","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_uncommit_reasons.md

## Nobody can list expense uncommit reasons through the API

Lines 7-9

Asked as: list uncommit reasons; see all uncommit reasons; uncommit reasons api list; why can't I see reasons

The expense_uncommit_reasons collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one expense uncommit reason through the API

Lines 11-13

Asked as: view one uncommit reason; open uncommit reason record; uncommit reason details; single reason not showing

The expense_uncommit_reasons collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser.

## Nobody can create, update or delete expense uncommit reasons through the API

Lines 15-25

Asked as: edit uncommit reason; cannot change uncommit reason; add new uncommit reason; delete uncommit reason api

The expense_uncommit_reasons collection has no create, update or delete rule. Nobody can create, change or delete a record through the API, except a PocketBase superuser. Turbo's own server code is not limited by these rules.

## Default outcome of an uncommit reason has only two allowed values

Lines 31-33

Asked as: uncommit reason default outcome; default outcome options; uncommit_only or uncommit_and_reject; what is default outcome; uncommit reason name required

Each expense uncommit reason must have a default outcome (default_outcome). The value is either uncommit_and_reject or uncommit_only, and only one can be chosen. The name is also required.
