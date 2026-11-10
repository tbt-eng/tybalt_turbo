<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_types.md","sha256":"c9438e11498dad9bdf467b87396b2ae765533a0fb03017d2079864a11128a7a0","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_types.md

## Who can list or view time types

Lines 7-17

Asked as: who can see time types; list time types access; view time type record; not signed in time types; anyone can read time types

Any signed-in person can list time types and view one time type record. A request with no signed-in person is refused.

## Only holders of the tt claim can create time types

Lines 19-24

Asked as: create time type permission; tt claim needed; add new time type; who can make time types; other staff cannot create time type

A person who is signed in and holds the tt claim (the claim named 'tt') can create a time type record. Other signed-in people cannot create one.

## Nobody can update or delete time types through the API

Lines 26-32

Asked as: cannot edit time type; delete time type api; update time type forbidden; change time type through api; superuser only time type changes

No one can update or delete a time type through the API, except a PocketBase superuser. Turbo's own server code can still do it.

## Time type code must be unique

Lines 39-51

Asked as: duplicate time type code; code already taken; unique time type code; same code twice; time type code clash

Every time type needs a code (code) of at least 1 character. Two time types cannot have the same code, because of a unique index.

## Time type field requirements

Lines 40-43

Asked as: time type required fields; name too short time type; allowed_fields json size; time type description optional; required_fields limit

A time type needs a name (name) of at least 2 characters. It also needs allowed_fields, a JSON value up to 2,000,000 bytes. The description (description) and required_fields (JSON up to 2,000,000 bytes) are optional.
