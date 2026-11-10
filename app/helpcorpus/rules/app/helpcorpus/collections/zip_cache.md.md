<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/zip_cache.md","sha256":"80c6230c8bd200570c6d6cc290cddbc0d2d3bf018ff7ef05c22ce358ab067b5c","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/zip_cache.md

## Who can list or view zip cache records

Lines 7-17

Asked as: who can list zip cache; zip cache access denied; can't view zip cache records; report claim needed for zip

Only a person who holds the report claim can list zip cache records (zip_cache) through the API. The same person can view one record. Everyone else is refused.

## Nobody can create, update or delete zip cache records through the API

Lines 19-29

Asked as: can't edit zip cache via api; delete zip cache record; zip cache read only; create zip cache api forbidden

No API rule lets anyone create, update or delete a zip cache record (zip_cache), except a PocketBase superuser. Turbo's own server code can still do these actions.

## Zip cache field limits

Lines 31-41

Asked as: zip file too large; zip upload size limit; max size zip cache; required fields zip cache

Each zip cache record requires a key, a class and one zip file. The zip file can be at most 536870912 bytes (512 MB), and only one file can be attached. The manifest (a JSON field) is optional.
